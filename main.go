package main

import (
	"fmt"
	"ksetup/pkg/config"
)

func main() {
	conf, err := config.Load("./example/config.yaml")
	if err != nil {
		panic(err)
	}

	fmt.Println(conf)
}
