package main

import (
	"flag"
	"log"
)

func main() {
	confFile := flag.String("config", "/etc/adnebulae/config.toml", "Config file location")
	dev := flag.Bool("devmode", false, "Reload templates on change")
	flag.Parse()

	srv, err := NewServer(*confFile)
	if err != nil {
		log.Fatal(err)
	}
	if *dev {
		srv.WatchChanges()
	}

	log.Fatal(srv.Run())
}
