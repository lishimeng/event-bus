package proc

import (
	"encoding/json"

	"github.com/lishimeng/go-log"
)

func Publish(m Message) {
	// TODO
	log.Info("publish message")
	bs, _ := json.Marshal(m)
	log.Info(string(bs))
}
