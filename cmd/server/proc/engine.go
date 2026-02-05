package proc

import "context"

func Start(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	default:
		one()
	}
}

func one() {
	//
	var m Message
	Callback(m)
}
