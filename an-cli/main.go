package main

import (
	"fmt"
	"os"
)

func main() {
	cli := NewCLI()
	cli.Subs(AnCliCommands...)
	err := cli.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}

}
