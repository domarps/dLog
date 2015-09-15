package echo

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
// echo code
func (s *Handler) Handle(in io.Reader, out io.Writer) {
	s.Logger.Println("Handling echo")

	for {
		n, err := in.Read(s.B)

		if err == proto.ErrEOM || err == io.EOF {
			break
		}

		if err != nil {
			s.Logger.Fatalln(err.Error())
			return
		}

		_, err = out.Write(s.B[:n])

		if err != nil {
			s.Logger.Fatalln(err.Error())
			return
		}
	}

	// Finish message
	_, err := out.Write(nil)

	if err != nil {
		s.Logger.Fatalln(err.Error())
		return
	}

	s.Logger.Printf(
		"Finished\n",
	)
}
