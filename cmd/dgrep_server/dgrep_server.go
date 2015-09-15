package main

import (
	mux "github.com/johnmcconnell/mux/qikmux"
	"github.com/johnmcconnell/proto/qik"
	"github.com/joho/godotenv"
	con "gitlab-beta.engr.illinois.edu/mcconne7/dlog/config"
	"gitlab-beta.engr.illinois.edu/mcconne7/dlog/server"
	"gitlab-beta.engr.illinois.edu/mcconne7/dlog/server/echo"
	"gitlab-beta.engr.illinois.edu/mcconne7/dlog/server/forward"
	"gitlab-beta.engr.illinois.edu/mcconne7/dlog/server/hello"
	"gitlab-beta.engr.illinois.edu/mcconne7/dlog/server/neighbors"
	"gitlab-beta.engr.illinois.edu/mcconne7/dlog/server/sys"
	"log"
	"os"
)

const (
	// Type is the connection type
	Type = "tcp"
	// Host is the host name
	Host = "127.0.0.1"
)

var (
	// Logger ...
	Logger = log.New(
		os.Stderr,
		"Server: ",
		log.LstdFlags|log.Lshortfile,
	)

	// Protocol ...
	Protocol = qik.NewProtocol()

	// Handlers ...
	Handlers = map[byte]mux.Handler{
		con.HelloCode: &hello.Handler{
			B:      make([]byte, 64),
			Logger: Logger,
		},
		con.EchoCode: &echo.Handler{
			B:      make([]byte, 64),
			Logger: Logger,
		},
		con.SystemCode: &sys.Handler{
			Logger: Logger,
		},
		con.CastCode: &forward.Handler{
			Neighbors: neighbors.New(
				"127.0.0.1:3000",
			),
			B:        make([]byte, 512),
			Logger:   Logger,
			Protocol: Protocol,
		},
	}
)

func usage() string {
	return `dgrep_server port`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file")
		os.Exit(-1)
	}

	if len(os.Args) < 2 {
		log.Fatalln(usage())
		os.Exit(-1)
	}

	os.Getenv("NEIGHBORS")

	m := mux.Mux{
		Handlers: Handlers,
	}

	Port := os.Args[1]

	s := server.New(
		Type,
		Host,
		Port,
		m,
		Protocol,
		Logger,
	)

	s.Serve()
}
