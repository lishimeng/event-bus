package admin

import "github.com/lishimeng/app-starter/server"

func Router(root server.Router) {

	localSecret := root.Path("/local_secret")
	localSecret.Post("/", apiConfigLocalSecret)
}
