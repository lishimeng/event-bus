package communication

import "github.com/lishimeng/app-starter/server"

func Router(root server.Router) {
	root.Post("/publish", apiPublish)
}
