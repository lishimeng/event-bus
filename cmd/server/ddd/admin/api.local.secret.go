package admin

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/lishimeng/app-starter"
	"github.com/lishimeng/app-starter/persistence"
	"github.com/lishimeng/app-starter/server"
	"github.com/lishimeng/event-bus/internal/db"
	"github.com/lishimeng/event-bus/internal/tls/cypher"
	"github.com/lishimeng/go-log"
)

type Req struct {
	Key string `json:"key"` // pem再转一次base64, 变成一行
}

// localSecretStatus 表示本地密钥是否已配置及其指纹。
type localSecretStatus struct {
	Configured     bool   `json:"configured"`
	KeyFingerprint string `json:"keyFingerprint,omitempty"` // 对库中整条 config 做 SHA256 前缀，便于核对是否变更（不含私钥明文）
}

// apiGetLocalSecret 返回本地密钥配置状态。
func apiGetLocalSecret(ctx server.Context) {
	var resp app.ResponseWrapper
	var list []db.SysConfig
	err := app.GetOrm().Model(new(db.SysConfig)).
		Equal("name", db.SysLocalSecret).Find(&list)
	if err != nil {
		log.Info(err)
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		ctx.Json(resp)
		return
	}
	if len(list) == 0 || len(list[0].Config) == 0 {
		resp.Code = http.StatusOK
		resp.Data = localSecretStatus{Configured: false}
		ctx.Json(resp)
		return
	}
	sum := sha256.Sum256([]byte(list[0].Config))
	resp.Code = http.StatusOK
	resp.Data = localSecretStatus{
		Configured:     true,
		KeyFingerprint: hex.EncodeToString(sum[:8]),
	}
	ctx.Json(resp)
}

// 更新本地密钥
func apiConfigLocalSecret(ctx server.Context) {

	var req Req
	var resp app.ResponseWrapper
	var err error

	err = ctx.C.ReadJSON(&req)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = err.Error()
		ctx.Json(resp)
		return
	}
	if len(req.Key) == 0 {
		resp.Code = http.StatusBadRequest
		resp.Message = "key is empty"
		ctx.Json(resp)
		return
	}
	bs, err := base64.StdEncoding.DecodeString(req.Key) // 转回pem原文
	if err != nil {
		log.Info("base64 decode key")
		log.Info(err)
		resp.Code = http.StatusBadRequest
		resp.Message = err.Error()
		ctx.Json(resp)
		return
	}
	var localSecret db.LocalSecurity
	_, err = cypher.LoadPrivateKey(bs)
	if err != nil {
		log.Info("parse private key failed")
		log.Info(err)
		resp.Code = http.StatusBadRequest
		resp.Message = err.Error()
		ctx.Json(resp)
		return
	}
	localSecret.RsaKey = req.Key
	bs, err = json.Marshal(localSecret)
	if err != nil {
		log.Info("local secret marshal failed")
		log.Info(err)
		resp.Code = http.StatusBadRequest
		resp.Message = err.Error()
		ctx.Json(resp)
		return
	}
	config := base64.StdEncoding.EncodeToString(bs) // 保存base64格式pem
	err = _saveLocalSecret(config)
	if err != nil {
		log.Info("save local secret failed")
		log.Info(err)
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		ctx.Json(resp)
		return
	}
	resp.Code = http.StatusOK
	resp.Message = "已保存"
	ctx.Json(resp)
}

func _saveLocalSecret(config string) (err error) {
	err = app.Transaction(func(tx persistence.TxContext) (e error) {
		var tmpList []db.SysConfig
		e = tx.Model(new(db.SysConfig)).
			Equal("name", db.SysLocalSecret).Find(&tmpList)
		if e != nil {
			return
		}
		if len(tmpList) == 0 {
			var item db.SysConfig
			item.Name = db.SysLocalSecret
			item.Config = config
			item.Status = 1
			e = tx.Create(&item)
		} else {
			item := tmpList[0]
			item.Config = config
			e = tx.Model(&item).Select("Config").Updates(&item)
		}
		return
	})
	return
}
