package admin

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/lishimeng/app-starter"
	"github.com/lishimeng/app-starter/server"
	"github.com/lishimeng/event-bus/internal/db"
	"github.com/lishimeng/go-log"
)

type ChannelReq struct {
	Code       string           `json:"code,omitempty"`
	Name       string           `json:"name,omitempty"`
	Route      string           `json:"route,omitempty"`
	Category   db.RouteCategory `json:"category,omitempty"`
	PrivateKey string           `json:"privateKey,omitempty"` // base64 pem
	PublicKey  string           `json:"publicKey,omitempty"`  // base64 pem
	Callback   string           `json:"callback,omitempty"`
}

func apiChannelConfig(ctx server.Context) {

	var err error
	var req ChannelReq
	var resp app.ResponseWrapper

	err = ctx.C.ReadJSON(&req)
	if err != nil {
		log.Info("req must be json")
		log.Info(err)
		resp.Code = http.StatusBadRequest
		ctx.Json(resp)
		return
	}

	if len(req.Name) == 0 {
		resp.Code = http.StatusBadRequest
		resp.Message = "name is required"
		ctx.Json(resp)
		return
	}
	if len(req.Route) == 0 {
		resp.Code = http.StatusBadRequest
		resp.Message = "route is required"
		ctx.Json(resp)
		return
	}
	if req.Category != db.PublishTo && req.Category != db.Subscribe {
		resp.Code = http.StatusBadRequest
		resp.Message = "category not support"
		ctx.Json(resp)
		return
	}

	var ch db.ChannelConfig
	isNew := req.Code == ""
	if isNew {
		var list []db.ChannelConfig
		_, err = app.GetOrm().Context.QueryTable(new(db.ChannelConfig)).
			Filter("router", req.Route).
			Filter("category", req.Category).
			All(&list)
		if err != nil {
			resp.Code = http.StatusInternalServerError
			resp.Message = err.Error()
			ctx.Json(resp)
			return
		}
		if len(list) > 0 {
			resp.Code = http.StatusBadRequest
			resp.Message = "channel already exists"
			ctx.Json(resp)
			return
		}
		ch.Code = genId()
		ch.Router = req.Route
	} else {
		var list []db.ChannelConfig
		_, err = app.GetOrm().Context.QueryTable(new(db.ChannelConfig)).
			Filter("code", req.Code).
			Filter("router", req.Route).
			All(&list)
		if err != nil {
			resp.Code = http.StatusInternalServerError
			resp.Message = err.Error()
			ctx.Json(resp)
			return
		}
		if len(list) == 0 {
			resp.Code = http.StatusBadRequest
			resp.Message = "channel not found"
			ctx.Json(resp)
			return
		}
		ch = list[0]
	}
	ch.Name = req.Name
	ch.Callback = req.Callback
	ch.Category = req.Category

	touchSecurity := len(req.PrivateKey) > 0 || len(req.PublicKey) > 0
	if touchSecurity {
		ch.UseSecurity = 1
		var secret db.ChannelSecurity
		var raw []byte
		if len(req.PrivateKey) > 0 {
			raw, err = base64.StdEncoding.DecodeString(stripBase64Noise(req.PrivateKey))
			if err != nil {
				resp.Code = http.StatusBadRequest
				resp.Message = err.Error()
				ctx.Json(resp)
				return
			}
			var pemBytes []byte
			pemBytes, err = normalizePrivateKeyToPEM(raw)
			if err != nil {
				resp.Code = http.StatusBadRequest
				resp.Message = err.Error()
				ctx.Json(resp)
				return
			}
			secret.RsaKey = string(pemBytes)
		}
		if len(req.PublicKey) > 0 {
			raw, err = base64.StdEncoding.DecodeString(stripBase64Noise(req.PublicKey))
			if err != nil {
				resp.Code = http.StatusBadRequest
				resp.Message = err.Error()
				ctx.Json(resp)
				return
			}
			var pemBytes []byte
			pemBytes, err = normalizePublicKeyToPEM(raw)
			if err != nil {
				resp.Code = http.StatusBadRequest
				resp.Message = err.Error()
				ctx.Json(resp)
				return
			}
			secret.RsaPem = string(pemBytes)
		}
		var secJSON []byte
		secJSON, err = json.Marshal(secret)
		if err != nil {
			resp.Code = http.StatusInternalServerError
			resp.Message = err.Error()
			ctx.Json(resp)
			return
		}
		ch.Security = base64.StdEncoding.EncodeToString(secJSON)
	} else if isNew {
		var secJSON []byte
		secJSON, err = json.Marshal(struct{}{})
		if err != nil {
			resp.Code = http.StatusInternalServerError
			resp.Message = err.Error()
			ctx.Json(resp)
			return
		}
		ch.Security = base64.StdEncoding.EncodeToString(secJSON)
		ch.UseSecurity = 0
	}

	if isNew {
		_, err = app.GetOrm().Context.Insert(&ch)
	} else {
		if touchSecurity {
			_, err = app.GetOrm().Context.Update(&ch,
				"name", "callback", "category", "use_security", "security")
		} else {
			_, err = app.GetOrm().Context.Update(&ch, "name", "callback", "category")
		}
	}
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		ctx.Json(resp)
		return
	}
	resp.Code = http.StatusOK
	resp.Data = ch
	ctx.Json(resp)
}

func genId() string {
	md5Ctx := md5.New()
	content := fmt.Sprintf("%d%d", time.Now().UnixMicro(), rand.Int63())
	md5Ctx.Write([]byte(content))
	bs := md5Ctx.Sum(nil)
	return fmt.Sprintf("%02x", bs)
}
