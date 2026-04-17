package admin

import (
	"encoding/base64"
	"net/http"

	"github.com/lishimeng/app-starter"
	"github.com/lishimeng/app-starter/server"
	"github.com/lishimeng/event-bus/internal/db"
	"github.com/lishimeng/go-log"
)

// ChannelSummary 管理端列表；privateKey/publicKey 为 PEM 整段的 StdEncoding Base64（与 POST channel 一致），便于载入表单。
type ChannelSummary struct {
	Id            int    `json:"id"`
	Code          string `json:"code"`
	Name          string `json:"name"`
	Category      int    `json:"category"`
	CategoryLabel string `json:"categoryLabel"`
	Router        string `json:"router"`
	UseSecurity   int    `json:"useSecurity"`
	HasPrivateKey bool   `json:"hasPrivateKey"`
	HasPublicKey  bool   `json:"hasPublicKey"`
	Callback      string `json:"callback"`
	PrivateKey    string `json:"privateKey,omitempty"`
	PublicKey     string `json:"publicKey,omitempty"`
}

func apiListChannels(ctx server.Context) {
	var resp app.ResponseWrapper
	var list []db.ChannelConfig
	_, err := app.GetOrm().Context.QueryTable(new(db.ChannelConfig)).OrderBy("id").All(&list)
	if err != nil {
		log.Info(err)
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		ctx.Json(resp)
		return
	}
	out := make([]ChannelSummary, 0, len(list))
	for _, c := range list {
		var sec db.ChannelSecurity
		_ = sec.Unmarshal(c.Security)
		sum := ChannelSummary{
			Id:            c.Id,
			Code:          c.Code,
			Name:          c.Name,
			Category:      int(c.Category),
			CategoryLabel: c.Category.String(),
			Router:        c.Router,
			UseSecurity:   c.UseSecurity,
			HasPrivateKey: len(sec.RsaKey) > 0,
			HasPublicKey:  len(sec.RsaPem) > 0,
			Callback:      c.Callback,
		}
		if len(sec.RsaKey) > 0 {
			sum.PrivateKey = base64.StdEncoding.EncodeToString([]byte(sec.RsaKey))
		}
		if len(sec.RsaPem) > 0 {
			sum.PublicKey = base64.StdEncoding.EncodeToString([]byte(sec.RsaPem))
		}
		out = append(out, sum)
	}
	resp.Code = http.StatusOK
	resp.Data = out
	ctx.Json(resp)
}
