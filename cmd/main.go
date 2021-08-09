package main

import (
	"fmt"
	"os"
)

var (
)



func main() {
	if err := App.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}