package ddd

import (
	"gitee.com/lishimeng/event-bus/cmd/server/ddd/admin"
	"github.com/lishimeng/app-starter/server"
)

func Router(root server.Router) {
	admin.Router(root.Path("/api").Path("/v1"))
}
