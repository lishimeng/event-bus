package main

import (
	"fmt"

	"github.com/lishimeng/app-starter/buildscript"
)

func main() {
	err := buildscript.Generate(
		buildscript.Project{
			Namespace: "github.com/event-bus",
		},
		buildscript.Application{
			Name:    "server",
			AppPath: "cmd/server",
			HasUI:   true,
		},
	)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("success")
	}
}
