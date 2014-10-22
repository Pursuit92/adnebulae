package main

import (
	"fmt"
	"github.com/Pursuit92/adnebulae"
	"os"
)

var conn *adnebulae.AdNebulae

func init() {
	var err error
	conn, err = adnebulae.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
