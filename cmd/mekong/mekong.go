package main

import (
	"log"
	"os"
)

func main() {
	cmd, err := newRootCmd(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	//if err := cmd.Execute(); err != nil {
	//	log.Fatal(err)
	//}
	cmd.Execute()
}