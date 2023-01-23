package main

import (
	"fmt"
	"log"
)

func main() {
	cfg := LoadConfig()
	fmt.Println(cfg)

	store, err := NewPostgresStore(cfg.Dbcfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	server := NewApiServer(cfg.Servercfg.Port, store)

	server.Run()

}
