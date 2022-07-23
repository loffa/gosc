package main

import (
	"github.com/loffa/gosc"
	"log"
)

func main() {
	cli, err := gosc.NewClient("127.0.0.1:1234")
	if err != nil {
		log.Fatalln(err)
	}
	res, err := cli.CallMessage("/hello")
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Server responded %v\n", res)
}
