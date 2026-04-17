package admin

import (
	"net/http"
	"strconv"
	"time"

	"github.com/lishimeng/app-starter"
	"github.com/lishimeng/app-starter/server"
	"github.com/lishimeng/event-bus/internal/db"
	"github.com/lishimeng/go-log"
)

// dataRecordItem 列表项（camelCase JSON，避免嵌入类型序列化不一致）。
type dataRecordItem struct {
	ID         int    `json:"id"`
	Code       string `json:"code"`
	ReferCode  string `json:"referCode"`
	Source     string `json:"source"`
	Route      int    `json:"route"`
	RouteLabel string `json:"routeLabel"`
	Payload    string `json:"payload"`
	BizPayload string `json:"bizPayload"`
	CreateTime string `json:"createTime,omitempty"`
	UpdateTime string `json:"updateTime,omitempty"`
}

type recordsListResp struct {
	Total int64            `json:"total"`
	Items []dataRecordItem `json:"items"`
}

func toDataRecordItem(r db.DataRecord) dataRecordItem {
	it := dataRecordItem{
		ID:         r.Id,
		Code:       r.Code,
		ReferCode:  r.ReferCode,
		Source:     r.Source,
		Route:      int(r.Route),
		RouteLabel: r.Route.String(),
		Payload:    r.Payload,
		BizPayload: r.BizPayload,
	}
	if !r.CreateTime.IsZero() {
		it.CreateTime = r.CreateTime.Format(time.RFC3339)
	}
	if !r.UpdateTime.IsZero() {
		it.UpdateTime = r.UpdateTime.Format(time.RFC3339)
	}
	return it
}

func apiListDataRecords(ctx server.Context) {
	var resp app.ResponseWrapper
	limit := int64(50)
	offset := int64(0)
	if v := ctx.C.URLParam("limit"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil && n > 0 && n <= 200 {
			limit = n
		}
	}
	if v := ctx.C.URLParam("offset"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil && n >= 0 {
			offset = n
		}
	}
	qs := app.GetOrm().Context.QueryTable(new(db.DataRecord))
	if src := ctx.C.URLParam("source"); src != "" {
		qs = qs.Filter("source", src)
	}
	total, err := qs.Count()
	if err != nil {
		log.Info(err)
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		ctx.Json(resp)
		return
	}
	var list []db.DataRecord
	_, err = qs.OrderBy("-id").Limit(limit).Offset(offset).All(&list)
	if err != nil {
		log.Info(err)
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		ctx.Json(resp)
		return
	}
	items := make([]dataRecordItem, 0, len(list))
	for i := range list {
		items = append(items, toDataRecordItem(list[i]))
	}
	resp.Code = http.StatusOK
	resp.Data = recordsListResp{Total: total, Items: items}
	ctx.Json(resp)
}
