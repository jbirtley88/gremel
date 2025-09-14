package main

import (
	"log"

	"github.com/jbirtley88/gremel/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
