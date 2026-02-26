package communication

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"gitee.com/lishimeng/event-bus/cmd/server/proc"
	"gitee.com/lishimeng/event-bus/internal/message"
	"gitee.com/lishimeng/event-bus/sdk"
	"github.com/lishimeng/app-starter/server"
	"github.com/lishimeng/go-log"
)

func apiPublish(ctx server.Context) {
	var err error
	var req sdk.Request
	var resp sdk.Resp
	log.Info("api publish")

	err = ctx.C.ReadJSON(&req)
	if err != nil {
		log.Info("read json fail")
		log.Info(err)
		resp.Code = http.StatusBadRequest
		resp.Msg = "read json fail"
		ctx.Json(resp)
		return
	}

	var biz message.BizMessage
	if len(req.Payload) > 0 { // 覆盖biz字段
		var bs []byte
		bs, err = base64.StdEncoding.DecodeString(req.Payload)
		if err != nil {
			log.Info("base64 decode fail")
			log.Info(err)
			resp.Code = http.StatusBadRequest
			resp.Msg = "base64 decode fail"
			ctx.Json(resp)
			return
		}
		err = json.Unmarshal(bs, &biz)
		if err != nil {
			log.Info("json unmarshal fail")
			log.Info(err)
			resp.Code = http.StatusBadRequest
			resp.Msg = "json unmarshal fail"
			ctx.Json(resp)
			return
		}
	}
	var opts []proc.MessageCreateFunc
	if len(req.ReferId) > 0 {
		opts = append(opts, proc.WithParentId(req.ReferId))
	}
	msg, err := proc.Create(req.Route, req.Biz, opts...)
	if err != nil {
		log.Info("create fail")
		log.Info(err)
		resp.Code = http.StatusInternalServerError
		resp.Msg = "create fail"
		ctx.Json(resp)
		return
	}
	proc.Publish(msg)
	resp.Code = http.StatusOK
	resp.Msg = msg.RequestId
	ctx.Json(resp)
}
