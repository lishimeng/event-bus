package communication

import (
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
	var msg message.Message
	proc.Publish()
}
