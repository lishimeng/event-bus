package ddd

import (
	"github.com/lishimeng/app-starter/server"
	"github.com/lishimeng/event-bus/cmd/server/ddd/admin"
	"github.com/lishimeng/event-bus/cmd/server/ddd/communication"
)

func Router(root server.Router) {
	api := root.Path("/api").Path("/v1")
	admin.Router(api.Path("/admin"))
	communication.Router(api.Path("/communication"))
}
