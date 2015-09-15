package server

import (
	"bytes"
	"fmt"
	mux "github.com/johnmcconnell/mux/qikmux"
	"github.com/johnmcconnell/proto"
	"io"
	"log"
	"net"
)

const (
	// BufferSize the size of the buffer
	BufferSize = 512
)

// AcceptErrorFunc is a function that is called when
// the server errors when accepting a socket
// it is only called with non-nil errors
type AcceptErrorFunc func(error)

// Server serves requests
type Server struct {
	Type     string
	Host     string
	Port     string
	End      bool
	Mux      mux.Mux
	Protocol proto.Protocol
	Logger   *log.Logger
}

// New creates a new server
func New(t, h, p string, m mux.Mux, pr proto.Protocol, l *log.Logger) *Server {
	s := Server{
		Type:     t,
		Host:     h,
		Port:     p,
		End:      false,
		Mux:      m,
		Protocol: pr,
		Logger:   l,
	}

	return &s
}

// Serve starts the server to serve requests
func (s *Server) Serve() {
	// Listen for incoming connections.
	sock, err := net.Listen(s.Type, s.URL())

	if err != nil {
		s.Logger.Fatalln(
			"Could not serve:",
			err.Error(),
		)

		return
	}

	// Close the listener when the server closes.
	defer sock.Close()

	s.Logger.Printf(
		"%v begins on: %v\n",
		s,
		s.URL(),
	)

	for !s.End {
		// Listen for an incoming connection.
		conn, err := sock.Accept()

		if err != nil {
			s.Logger.Println(
				"Could not accept socket:",
				err.Error(),
			)

			return
		}

		// Handle connections in a new goroutine.
		go s.handleRequest(conn)
	}
}

// URL the url string used to open the connection
func (s *Server) URL() string {
	return s.Host + ":" + s.Port
}

// Halt gracefully halt the server
func (s *Server) Halt() {
	s.End = true
}

// String makes a printable string to see the server
func (s *Server) String() string {
	str := fmt.Sprintf(
		"Server [%v]: is ending? %v",
		s.URL(),
		s.End,
	)

	return str
}

// Handles incoming requests.
func (s *Server) handleRequest(conn net.Conn) {
	s.Logger.Println(
		"Begin Handling Request",
	)

	Decoder := s.Protocol.NewReader(conn)

	// Close the connection when you're done with it.
	defer func() {
		conn.Close()
		log.Println(
			"Finished Handling the Request",
		)
	}()

	// Make a buffer to hold incoming data.
	BS := make([]byte, BufferSize)

	s.Logger.Println("Reading first bytes")

	// Read the incoming connection into the buffer.
	N, err := Decoder.Read(BS)

	if err != nil && err != proto.ErrEOM && err != io.EOF {
		s.Logger.Println(
			"Count not read first bytes off connection",
			err.Error(),
		)

		return
	}

	if N == 0 {
		s.Logger.Println(
			"No bytes were read from connection",
			err.Error(),
		)
		return
	}

	// c is the command byte
	// r is the remaining bytes
	// in the buffer
	c, r := BS[0], BS[1:N]

	R := bytes.NewReader(r)

	// Combine into a multi-reader to re-read the previous bytes
	Decoder = io.MultiReader(R, Decoder)

	Encoder := s.Protocol.NewWriter(conn)

	s.Logger.Printf(
		"Selecting [%v]\n",
		c,
	)

	if !s.Mux.Select(c, Decoder, Encoder) {
		Encoder.Write(
			[]byte("Unable to process code"),
		)
	}
}
