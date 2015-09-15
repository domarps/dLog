package hello

import (
	"github.com/johnmcconnell/proto"
	"io"
	"log"
)

// Handler handles system requests
type Handler struct {
	B      []byte
	Logger *log.Logger
}

// Handle this command is in charge of running
// the status call
func (s *Handler) Handle(in io.Reader, out io.Writer) {
	s.Logger.Println("Handling status")

	status := `{ "status": "running", "running": true }`

	for {
		_, err := in.Read(s.B)
		if err == proto.ErrEOM {
			break
		}

		if err != nil {
			s.Logger.Fatalln(err.Error())
			return
		}
	}

	s.Logger.Printf(
		"Writing [%v]\n",
		status,
	)

	_, err := out.Write(
		[]byte(status),
	)

	if err != nil {
		s.Logger.Fatalln(err.Error())
		return
	}

	_, err = out.Write(
		nil,
	)

	if err != nil {
		s.Logger.Fatalln(err.Error())
		return
	}

	s.Logger.Println(
		"Finished Status",
	)
}
