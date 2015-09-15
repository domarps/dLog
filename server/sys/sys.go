package sys

import (
	"github.com/johnmcconnell/proto"
	"io"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
)

const (
	// BufferSize is the size of the buffer
	BufferSize = 512
)

// Handler handles system requests
type Handler struct {
	Logger *log.Logger
}

// Handle this command is in charge of running
// system calls
func (s *Handler) Handle(in io.Reader, out io.Writer) {
	s.Logger.Println("Handling system")

	bs, err := ioutil.ReadAll(in)

	if err == proto.ErrEOM {
		command := string(bs)
		err = s.Exec(command, in, out)
	}

	if err != nil {
		s.Logger.Fatalln(err.Error())
		return
	}

	out.Write(nil)

	s.Logger.Println(
		"Finished Status",
	)
}

// Exec this command is in charge of
// running system calls
func (s *Handler) Exec(command string, in io.Reader, out io.Writer) error {
	tokens := strings.Fields(command)

	s.Logger.Printf(
		"Command [%v]\n",
		tokens,
	)

	c := exec.Command(tokens[0], tokens[1:]...)

	c.Stdout = out

	return c.Run()
}
