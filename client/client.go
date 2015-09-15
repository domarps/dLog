package client

import (
	"bufio"
	"fmt"
	"github.com/johnmcconnell/proto"
	con "gitlab-beta.engr.illinois.edu/mcconne7/dlog/config"
	"io"
	"log"
	"net"
	"strings"
)

const (
	// BufferSize ...
	BufferSize = 512
	// NetworkProtocol ...
	NetworkProtocol = "tcp"
)

// Client ...
type Client struct {
	Writer   io.Writer
	Reader   io.Reader
	Logger   *log.Logger
	Protocol proto.Protocol
}

// New ...
func New(w io.Writer, r io.Reader, p proto.Protocol, l *log.Logger) *Client {
	c := Client{
		Writer:   w,
		Reader:   r,
		Protocol: p,
		Logger:   l,
	}

	return &c
}

// Run ...
func (c *Client) Run(url string) error {
	c.Logger.Printf(
		"Connection to: [%v]\n",
		url,
	)

	R := bufio.NewReader(c.Reader)

	for {
		line, err := R.ReadString('\n')

		if err == io.EOF {
			return c.RunLine(url, line)
		}

		if err != nil {
			return err
		}

		err = c.RunLine(url, line)

		if err != nil {
			return err
		}
	}

	return nil
}

// RunLine ...
func (c *Client) RunLine(url string, line string) error {
	tokens := strings.Fields(line)

	if len(tokens) == 0 {
		return fmt.Errorf(
			"Command has zero tokens [%v]",
			tokens,
		)
	}

	c.Logger.Println(
		"Running:",
		url,
		tokens,
	)

	switch tokens[0] {
	case "status":
		return c.Status(url, tokens)
	case "echo":
		return c.Echo(url, tokens)
	case "sys":
		return c.Sys(url, tokens)
	case "cast":
		return c.Forward(url, tokens)
	default:
		return fmt.Errorf(
			"Command not found [%v]",
			tokens[0],
		)
	}

	return nil
}

// Status ...
func (c *Client) Status(url string, str []string) error {
	conn, err := net.Dial(NetworkProtocol, url)

	if err != nil {
		return err
	}

	c.Logger.Printf(
		"Connected to [%v]\n",
		url,
	)

	Decoder := c.Protocol.NewReader(conn)
	Encoder := c.Protocol.NewWriter(conn)

	defer conn.Close()

	code := []byte{
		con.HelloCode,
	}

	c.Logger.Printf(
		"Writing: [%v]",
		code,
	)

	_, err = Encoder.Write(code)

	if err != nil {
		return err
	}

	_, err = Encoder.Write(nil)

	if err != nil {
		return err
	}

	b := make([]byte, BufferSize)

	for {
		n, err := Decoder.Read(b)

		if err != nil && err != io.EOF && err != proto.ErrEOM {
			c.Logger.Fatalln(
				err.Error(),
			)
		}

		c.Writer.Write(b[:n])

		if err == io.EOF || err == proto.ErrEOM {
			c.Writer.Write(
				[]byte("\n"),
			)

			break
		}

		if err != nil {
			return err
		}
	}

	return nil
}

// Echo ...
func (c *Client) Echo(url string, str []string) error {
	conn, err := net.Dial(NetworkProtocol, url)

	if err != nil {
		return err
	}

	c.Logger.Printf(
		"Connected to [%v]\n",
		url,
	)

	Decoder := c.Protocol.NewReader(conn)
	Encoder := c.Protocol.NewWriter(conn)

	defer conn.Close()

	code := []byte{
		con.EchoCode,
	}

	c.Logger.Printf(
		"Writing: [%v]",
		code,
	)

	_, err = Encoder.Write(code)

	if err != nil {
		return err
	}

	if len(str) < 2 {
		return fmt.Errorf(
			"echo needs at least one argument",
		)
	}

	_, err = Encoder.Write(
		[]byte(strings.Join(str[1:], " ")),
	)

	if err != nil {
		return err
	}

	_, err = Encoder.Write(nil)

	if err != nil {
		return err
	}

	b := make([]byte, BufferSize)

	for {
		n, err := Decoder.Read(b)

		if err != nil && err != io.EOF && err != proto.ErrEOM {
			c.Logger.Fatalln(
				err.Error(),
			)
		}

		c.Writer.Write(b[:n])

		if err == io.EOF || err == proto.ErrEOM {
			c.Writer.Write(
				[]byte("\n"),
			)

			break
		}

		if err != nil {
			return err
		}
	}

	return nil
}

// Sys ...
func (c *Client) Sys(url string, str []string) error {
	conn, err := net.Dial(NetworkProtocol, url)

	if err != nil {
		return err
	}

	c.Logger.Printf(
		"Connected to [%v]\n",
		url,
	)

	Decoder := c.Protocol.NewReader(conn)
	Encoder := c.Protocol.NewWriter(conn)

	defer conn.Close()

	code := []byte{
		con.SystemCode,
	}

	c.Logger.Printf(
		"Writing: [%v]",
		code,
	)

	_, err = Encoder.Write(code)

	if err != nil {
		return err
	}

	_, err = Encoder.Write(
		[]byte(strings.Join(str[1:], " ")),
	)

	if err != nil {
		return err
	}

	_, err = Encoder.Write(nil)

	if err != nil {
		return err
	}

	b := make([]byte, BufferSize)

	for {
		n, err := Decoder.Read(b)

		if err != nil && err != io.EOF && err != proto.ErrEOM {
			c.Logger.Fatalln(
				err.Error(),
			)
		}

		c.Writer.Write(b[:n])

		if err == io.EOF || err == proto.ErrEOM {
			c.Writer.Write(
				[]byte("\n"),
			)

			break
		}

		if err != nil {
			return err
		}
	}

	return nil
}

// Forward ...
func (c *Client) Forward(url string, str []string) error {
	conn, err := net.Dial(NetworkProtocol, url)

	if err != nil {
		return err
	}

	c.Logger.Printf(
		"Connected to [%v]\n",
		url,
	)

	Decoder := c.Protocol.NewReader(conn)
	Encoder := c.Protocol.NewWriter(conn)

	defer conn.Close()

	code := []byte{
		con.CastCode,
	}

	c.Logger.Printf(
		"Writing Cast: [%v]",
		code,
	)

	_, err = Encoder.Write(code)

	if err != nil {
		return err
	}

	switch str[1] {
	case "status":
		code[0] = con.HelloCode
	case "echo":
		code[0] = con.EchoCode
	case "sys":
		code[0] = con.SystemCode
	default:
		return fmt.Errorf(
			"Cannot cast the command:[%v]",
			str[1],
		)
	}

	c.Logger.Printf(
		"Writing Code: [%v]",
		code,
	)

	_, err = Encoder.Write(code)

	if err != nil {
		return err
	}

	if len(str) > 2 {
		_, err = Encoder.Write(
			[]byte(strings.Join(str[2:], " ")),
		)
	}

	if err != nil {
		return err
	}

	_, err = Encoder.Write(nil)

	if err != nil {
		return err
	}

	b := make([]byte, BufferSize)

	for {
		n, err := Decoder.Read(b)

		if err != nil && err != io.EOF && err != proto.ErrEOM {
			c.Logger.Fatalln(
				err.Error(),
			)
		}

		c.Writer.Write(b[:n])

		if err == io.EOF || err == proto.ErrEOM {
			c.Writer.Write(
				[]byte("\n"),
			)

			break
		}

		if err != nil {
			return err
		}
	}

	return nil
}
