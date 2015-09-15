package main

import (
	"fmt"
	"github.com/johnmcconnell/proto/qik"
	"gitlab-beta.engr.illinois.edu/mcconne7/dlog/client"
	"log"
	"os"
)

var (
	// Logger ...
	Logger = log.New(
		os.Stderr,
		"Client: ",
		log.LstdFlags|log.Lshortfile,
	)

	// Protocol ...
	Protocol = qik.NewProtocol()
)

func usage() string {
	return `dgrep_client url command`
}

func main() {
	c := client.New(
		os.Stdout,
		os.Stdin,
		Protocol,
		Logger,
	)

	if len(os.Args) < 1 {
		fmt.Println(usage())
		os.Exit(-1)
	}

	url := os.Args[1]

	err := c.Run(url)

	if err != nil {
		c.Logger.Fatalln(
			err.Error(),
		)
	}
}
