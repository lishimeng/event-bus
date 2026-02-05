package admin

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"gitee.com/lishimeng/event-bus/internal/db"
	"gitee.com/lishimeng/event-bus/internal/tls/cypher"
	"github.com/lishimeng/app-starter"
	"github.com/lishimeng/app-starter/persistence"
	"github.com/lishimeng/app-starter/server"
	"github.com/lishimeng/go-log"
)

type Req struct {
	Key string `json:"key"` // pem再转一次base64, 变成一行
}

func apiConfigLocalSecret(ctx server.Context) {

	var req Req
	var resp app.Response
	var err error

	err = ctx.C.ReadJSON(&req)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = err.Error()
		return
	}
	if len(req.Key) == 0 {
		resp.Code = http.StatusBadRequest
		resp.Message = "key is empty"
		return
	}
	bs, err := base64.StdEncoding.DecodeString(req.Key) // 转回pem原文
	if err != nil {
		log.Info("base64 decode key")
		log.Info(err)
		resp.Code = http.StatusBadRequest
		resp.Message = err.Error()
		return
	}
	var localSecret db.LocalSecurity
	_, err = cypher.LoadPrivateKey(bs)
	if err != nil {
		log.Info("parse private key failed")
		log.Info(err)
		resp.Code = http.StatusBadRequest
		resp.Message = err.Error()
		return
	}
	localSecret.RsaKey = req.Key
	bs, err = json.Marshal(localSecret)
	if err != nil {
		log.Info("local secret marshal failed")
		log.Info(err)
		resp.Code = http.StatusBadRequest
		resp.Message = err.Error()
		return
	}
	config := base64.StdEncoding.EncodeToString(bs) // 保存base64格式pem
	err = _saveLocalSecret(config)
	if err != nil {
		log.Info("save local secret failed")
		log.Info(err)
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		return
	}
	resp.Code = http.StatusOK
	ctx.Json(&resp)
}

func _saveLocalSecret(config string) (err error) {
	err = app.GetOrm().Transaction(func(tx persistence.TxContext) (e error) {
		// 检查数据库在是否存在
		var tmpList []db.SysConfig
		_, e = tx.Context.QueryTable(new(db.SysConfig)).
			Filter("name", db.SysLocalSecret).All(tmpList)
		if e != nil {
			return
		}
		if len(tmpList) == 0 { // insert
			var item db.SysConfig
			item.Name = db.SysLocalSecret
			item.Config = config
			item.Status = 1
			_, e = tx.Context.Insert(&item)
		} else { // update
			var item = tmpList[0]
			item.Config = config
			_, e = tx.Context.Update(&item, "config")
		}
		return
	})
	return
}
