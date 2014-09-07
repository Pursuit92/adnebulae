package main

import (
	"flag"
	an "github.com/Pursuit92/adnebulae"
	"log"
)

func main() {
	confFile := flag.String("config", "/etc/adnebulae/config.toml", "Config file location")
	dev := flag.Bool("devmode", false, "Reload templates on change")
	flag.Parse()

	srv, err := an.NewServer(*confFile)
	if err != nil {
		log.Fatal(err)
	}
	if *dev {
		srv.WatchChanges()
	}

	log.Fatal(srv.Run())
}

