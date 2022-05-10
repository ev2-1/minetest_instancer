package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"strings"
)

// A TelnetWriter can be used to print something at the other end
// of a telnet connection. It implements the io.Writer interface.
type TelnetWriter struct {
	conn net.Conn
}

// Write writes its parameter to the telnet connection.
// A trailing newline is always appended.
// It returns the number of bytes written and an error.
func (tw *TelnetWriter) Write(p []byte) (n int, err error) {
	return tw.conn.Write(append(p, '\n'))
}

var telnetCh = make(chan struct{})

type dest struct {
	conn     net.Conn
	destType string
}

var boradcastDests []dest

func telnetBroadcaster(port string) error {
	ln, err := net.Listen("tcp", "[::]:"+port)
	if err != nil {
		return err
	}
	defer ln.Close()

	log.Println("listen telnet", ln.Addr())

	for {
		select {
		case <-telnetCh:
			return nil
		default:
			conn, err := ln.Accept()
			if err != nil {
				log.Print(err)
				continue
			}

			s, err := bufio.NewReader(conn).ReadString('\n')
			i := int(math.Max(float64(len(s)-1), 1))
			if len(s) == 0 {
				return err
			}

			s = s[:i]

			boradcastDests = append(boradcastDests, dest{
				conn:     conn,
				destType: s,
			})
		}
	}
}

func telnetServer(port string) error {
	ln, err := net.Listen("tcp", "[::]:"+port)
	if err != nil {
		return err
	}
	defer ln.Close()

	log.Println("listen telnet", ln.Addr())

	for {
		select {
		case <-telnetCh:
			return nil
		default:
			conn, err := ln.Accept()
			if err != nil {
				log.Print(err)
				continue
			}

			go handleTelnet(conn)
		}
	}
}

func broadcast(str, dst string) {
	for _, dest := range boradcastDests {
		if dest.destType == dst {
			io.WriteString(dest.conn, str+"\r")
		}
	}
}

func handleTelnet(conn net.Conn) {
	prefix := fmt.Sprintf("[telnet %s] ", conn.RemoteAddr())
	fmt.Println(prefix, "connected")

	fmt.Println(prefix, "<->", "connect")

	defer fmt.Println(prefix, "<->", "disconnect")
	defer conn.Close()

	readString := func(delim byte) (string, error) {
		s, err := bufio.NewReader(conn).ReadString(delim)
		i := int(math.Max(float64(len(s)-1), 1))
		if len(s) == 0 {
			return s, err
		}

		s = s[:i]
		return s, err
	}

	writeString := func(s string) (n int, err error) {
		return io.WriteString(conn, s)
	}

	for {
		s, err := readString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}

			log.Print(err)
			continue
		}

		fmt.Println(prefix, "->", "command", s)

		if s == "\\quit" || s == "\\q" {
			return
		}

		result := handle(prefix, &TelnetWriter{conn: conn}, strings.Split(s, " "))
		if result != "\n" {
			fmt.Println(prefix, "<-", result)
			writeString(result + "\n")
		}
	}
}
