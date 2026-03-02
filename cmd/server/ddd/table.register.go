package ddd

import (
	"github.com/lishimeng/event-bus/internal/db"
)

func Tables() (t []any) {
	t = append(t)
	t = append(t, db.CommonModels...)
	return
}
