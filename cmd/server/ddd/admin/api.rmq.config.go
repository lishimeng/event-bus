package admin

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/lishimeng/app-starter"
	"github.com/lishimeng/app-starter/persistence"
	"github.com/lishimeng/app-starter/server"
	"github.com/lishimeng/event-bus/internal/db"
	"github.com/lishimeng/event-bus/providers/RocketMqProvider"
	"github.com/lishimeng/go-log"
)

type rmqConfigGetResp struct {
	Configured bool                        `json:"configured"`
	Config     *RocketMqProvider.RmqConfig `json:"config,omitempty"`
}

func apiGetRmqConfig(ctx server.Context) {
	var resp app.ResponseWrapper
	var list []db.SysConfig
	_, err := app.GetOrm().Context.QueryTable(new(db.SysConfig)).
		Filter("name", db.SysRmqConfig).All(&list)
	if err != nil {
		log.Info(err)
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		ctx.Json(resp)
		return
	}
	if len(list) == 0 || strings.TrimSpace(list[0].Config) == "" {
		resp.Code = http.StatusOK
		resp.Data = rmqConfigGetResp{Configured: false}
		ctx.Json(resp)
		return
	}
	raw, err := base64.StdEncoding.DecodeString(list[0].Config)
	if err != nil {
		resp.Code = http.StatusOK
		resp.Data = rmqConfigGetResp{Configured: true}
		resp.Message = "config 字段 Base64 解码失败: " + err.Error()
		ctx.Json(resp)
		return
	}
	var cfg RocketMqProvider.RmqConfig
	err = json.Unmarshal(raw, &cfg)
	if err != nil {
		resp.Code = http.StatusOK
		resp.Data = rmqConfigGetResp{Configured: true}
		resp.Message = "JSON 解析失败: " + err.Error()
		ctx.Json(resp)
		return
	}
	resp.Code = http.StatusOK
	resp.Data = rmqConfigGetResp{Configured: true, Config: &cfg}
	ctx.Json(resp)
}

func apiSaveRmqConfig(ctx server.Context) {
	var req RocketMqProvider.RmqConfig
	var resp app.ResponseWrapper

	err := ctx.C.ReadJSON(&req)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = err.Error()
		ctx.Json(resp)
		return
	}
	if strings.TrimSpace(req.Endpoint) == "" {
		resp.Code = http.StatusBadRequest
		resp.Message = "endpoint 必填"
		ctx.Json(resp)
		return
	}

	raw, err := json.Marshal(&req)
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		ctx.Json(resp)
		return
	}
	encoded := base64.StdEncoding.EncodeToString(raw)
	err = saveSysConfigNamed(db.SysRmqConfig, encoded)
	if err != nil {
		log.Info("save rmq_config failed: %v", err)
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		ctx.Json(resp)
		return
	}
	//_ = sysCfg.ReloadFromDB()
	resp.Code = http.StatusOK
	resp.Message = "已写入 sys_config；RocketMQ 连接需重启服务后生效"
	ctx.Json(resp)
}

func saveSysConfigNamed(name string, configB64 string) (err error) {
	err = app.GetOrm().Transaction(func(tx persistence.TxContext) (e error) {
		var tmpList []db.SysConfig
		_, e = tx.Context.QueryTable(new(db.SysConfig)).
			Filter("name", name).All(&tmpList)
		if e != nil {
			return
		}
		if len(tmpList) == 0 {
			var item db.SysConfig
			item.Name = name
			item.Config = configB64
			item.Status = 1
			_, e = tx.Context.Insert(&item)
		} else {
			item := tmpList[0]
			item.Config = configB64
			_, e = tx.Context.Update(&item, "config")
		}
		return
	})
	return
}
