package forward

import (
	"github.com/johnmcconnell/proto"
	"gitlab-beta.engr.illinois.edu/mcconne7/dlog/server/neighbors"
	"io"
	"io/ioutil"
	"log"
	"net"
)

const (
	// NetworkProtocol ...
	NetworkProtocol = "tcp"
)

// Handler handles system requests
type Handler struct {
	B         []byte
	Protocol  proto.Protocol
	Neighbors neighbors.Neighbors
	Logger    *log.Logger
}

// Handle this command is in charge of running
// forwarding calls
func (s *Handler) Handle(in io.Reader, out io.Writer) {
	s.Logger.Println("Handling forward")

	bs, err := ioutil.ReadAll(in)

	if err != nil && err != proto.ErrEOM {
		s.Logger.Println(err.Error())
		return
	}

	s.Logger.Printf(
		"Casting:\n[%v]",
		bs,
	)

	conn, err := net.Dial(NetworkProtocol, "127.0.0.1:3000")

	if err != nil {
		s.Logger.Println(err.Error())
		return
	}

	Encoder := s.Protocol.NewWriter(conn)
	_, err = Encoder.Write(bs)
	Encoder.Write(nil)

	if err != nil {
		s.Logger.Println(err.Error())
		return
	}

	Decoder := s.Protocol.NewReader(conn)

	for {
		n, err := Decoder.Read(s.B)

		if err == io.EOF {
			break
		}

		if err == proto.ErrEOM {
			out.Write(s.B[:n])
			break
		}

		if err != nil {
			s.Logger.Println(err.Error())
			return
		}

		s.Logger.Printf(
			"Reading:\n%v\n",
			string(s.B[:n]),
		)

		out.Write(s.B[:n])
	}

	out.Write(nil)

	s.Logger.Println(
		"Finished Status",
	)
}

// DialAll this command is in charge of
// running system calls
// returns the multi writer of all the connections
func (s *Handler) DialAll() ([]net.Conn, io.Reader, io.Writer) {
	L := len(s.Neighbors)

	conns := make([]net.Conn, L)
	ws := make([]io.Writer, L)
	rs := make([]io.Reader, L)

	for i, n := range s.Neighbors {
		conn, err := net.Dial(NetworkProtocol, n.URL)

		if err != nil {
			s.Logger.Println(err.Error())
			continue
		}

		conns[i] = conn
		ws[i] = conn
		rs[i] = conn
	}

	return conns, io.MultiReader(rs...), io.MultiWriter(ws...)
}
