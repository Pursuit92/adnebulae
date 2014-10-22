package main

import "os"

func main() {
	cli := NewCLI()
	cli.Subs(AnCliCommands...)
	cli.Run(os.Args)
}
