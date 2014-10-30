package main

import (
	"fmt"
	"os"

	"github.com/Pursuit92/adnebulae"
)

var conn *adnebulae.AdNebulae

func start() {
	var err error
	conn, err = adnebulae.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
