package main

import (
	"context"
	"fmt"
	"os"

	"time"

	"github.com/lishimeng/app-starter"
	"github.com/lishimeng/app-starter/persistence"
	"github.com/lishimeng/event-bus/cmd/server/ddd"
	"github.com/lishimeng/event-bus/internal/etc"
	"github.com/lishimeng/go-log"
)

import _ "github.com/lib/pq"

func main() {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	err := _main()
	if err != nil {
		fmt.Println(err)
	}
	time.Sleep(time.Millisecond * 200)
}

func _main() (err error) {
	application := app.New()

	err = application.Start(func(ctx context.Context, builder *app.ApplicationBuilder) error {

		var err error
		err = builder.LoadConfig(&etc.Config, nil)
		if err != nil {
			return err
		}

		dbConfig := persistence.PostgresConfig{
			UserName:  etc.Config.Db.User,
			Password:  etc.Config.Db.Password,
			Host:      etc.Config.Db.Host,
			Port:      etc.Config.Db.Port,
			DbName:    etc.Config.Db.Database,
			InitDb:    true,
			AliasName: "default",
			SSL:       etc.Config.Db.Ssl,
		}
		log.Info("database config: %s[%d]-->%s", dbConfig.Host, dbConfig.Port, dbConfig.DbName)

		var dbLogEnabled = os.Getenv("DB_LOG_ENABLED")
		var webLogEnabled = os.Getenv("WEB_LOG_ENABLED")

		builder.
			EnableDatabase(dbConfig.Build(), ddd.Tables()...).
			ComponentBefore(ddd.BeforeWeb).
			EnableWeb(etc.Config.Web.Listen, ddd.Router).
			PrintVersion()
		if dbLogEnabled == "1" {
			builder.EnableDatabaseLog()
		}
		if webLogEnabled == "1" {
			builder.SetWebLogLevel("DEBUG")
		}
		return err
	}, func(s string) {
		log.Info(s)
	})

	return
}
