package admin

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"gitee.com/lishimeng/event-bus/internal/db"
	"github.com/lishimeng/app-starter"
	"github.com/lishimeng/app-starter/server"
	"github.com/lishimeng/go-log"
)

type ChannelReq struct {
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
	if req.Category != db.Publish && req.Category != db.Subscriber {
		resp.Code = http.StatusBadRequest
		resp.Message = "category not support"
		ctx.Json(resp)
		return
	}

	var ch db.ChannelConfig

	ch, err = svsGetChannel(req.Route)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = err.Error()
		ctx.Json(resp)
		return
	}
	ch.Code = genId()
	ch.Name = req.Name
	ch.Callback = req.Callback
	ch.Category = req.Category
	ch.Router = req.Route

	var bs []byte
	var secret db.ChannelSecurity
	if len(req.PrivateKey) > 0 {
		bs, err = base64.StdEncoding.DecodeString(req.PrivateKey)
		if err != nil {
			resp.Code = http.StatusBadRequest
			resp.Message = err.Error()
			ctx.Json(resp)
			return
		}
		secret.RsaKey = string(bs)
	}
	if len(req.PublicKey) > 0 {
		bs, err = base64.StdEncoding.DecodeString(req.PublicKey)
		if err != nil {
			resp.Code = http.StatusBadRequest
			resp.Message = err.Error()
			ctx.Json(resp)
			return
		}
		secret.RsaPem = string(bs)
	}
	bs, err = json.Marshal(secret)
	ch.Security = base64.StdEncoding.EncodeToString(bs)
	_, err = app.GetOrm().Context.Insert(&ch)
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

func svsGetChannel(router string) (ch db.ChannelConfig, err error) {
	var list []db.ChannelConfig
	_, err = app.GetOrm().Context.
		QueryTable(new(db.ChannelConfig)).
		Filter("router", router).All(&list)
	if err != nil {
		return
	}
	if len(list) == 0 {
		err = errors.New("no channel config")
		return
	}
	ch = list[0]
	return
}

func genId() string {
	md5Ctx := md5.New()
	content := fmt.Sprintf("%d%d", time.Now().UnixMicro(), rand.Int63())
	md5Ctx.Write([]byte(content))
	bs := md5Ctx.Sum(nil)
	return fmt.Sprintf("%02x", bs)
}
