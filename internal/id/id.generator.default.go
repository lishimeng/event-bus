package id

import "fmt"

var defaultGenerator *Snowflake

func init() {

	var err error
	defaultGenerator, err = NewSnowflake(1)
	if err != nil {
		panic(err)
	}
}

func GenId() string {
	id := defaultGenerator.Generate()
	return fmt.Sprintf("%016x", id)
}
