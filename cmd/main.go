package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

const (
	LOCALHOST    = "127.0.0.1"
	WAITINTERVAL = 10 * time.Second
)

func main() {
	// Bind port for listening
	listener, remotePort := bindPort()
	if listener == nil {
		return
	}

	// Listen and connect to remote in parallel (Exit on Ctrl+C)
	done := make(chan struct{})
	go func() {
		if err := listen(listener); err != nil {
			fmt.Println(err)
		}
		close(done)
	}()

	go func() {
		if err := connect(remotePort); err != nil {
			fmt.Println(err)
		}
		close(done)
	}()

	<-done
}

// Makes a connection to remote/peer. Tries until connection succeeds.
// Calls handshake on a succeeded remote and breaks loop.
func connect(remotePort uint16) error {
	for {
		// Make a connection
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", LOCALHOST, remotePort))
		if err == nil {
			// Start communication
			if err := doHandshake(conn); err != nil {
				fmt.Println(err)
			}
			return nil
		}

		fmt.Println("Waiting for remote")
		// Wait some time and try again
		time.Sleep(WAITINTERVAL)
	}
}

// Listens for connections from another peer. Processes incoming connection and breaks loop.
func listen(listener net.Listener) error {
	for {
		// addr is not used as we know it is localhost
		conn, _ := listener.Accept()
		if err := processIncoming(conn); err != nil {
			fmt.Println("handshake failed, try again.", err)
		} else {
			fmt.Println("handshake success.")
			return nil
		}
	}
}

// Returns listener and remote port
func bindPort() (net.Listener, uint16) {
	// Bind to port 4000
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:4000", LOCALHOST))
	if err == nil {
		return listener, 4001 // If succeeds then remote is port 4001
	}

	// If fails, then 4000 is already running another instance so use 4001
	listener, err = net.Listen("tcp", fmt.Sprintf("%s:4001", LOCALHOST))
	if err != nil {
		return nil, 0
	}
	return listener, 4000 // In this case remote port is 4000
}

func readLine(conn io.Reader) error {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return err
	}
	fmt.Print(string(buf[:n]))
	return nil
}

// Initiate handshake
func doHandshake(conn net.Conn) error {
	_, err := conn.Write([]byte("1. Hello! Here are my encryption methods.\n"))
	if err != nil {
		return err
	}
	if err := readLine(conn); err != nil {
		return err
	}
	_, err = conn.Write([]byte("3. Here is encrypted secret-key.\n"))
	if err != nil {
		return err
	}
	if err := readLine(conn); err != nil {
		return err
	}
	_, err = conn.Write([]byte("5. This is encrypted sample message.\n"))
	if err != nil {
		return err
	}
	return readLine(conn)
}

// Process incoming handshake request
func processIncoming(conn net.Conn) error {
	if err := readLine(conn); err != nil {
		return err
	}
	_, err := conn.Write([]byte("2. Hello! Here is my public key\n"))
	if err != nil {
		return err
	}
	if err := readLine(conn); err != nil {
		return err
	}
	_, err = conn.Write([]byte("4. Got secret-key.\n"))
	if err != nil {
		return err
	}
	if err := readLine(conn); err != nil {
		return err
	}
	_, err = conn.Write([]byte("6. Verified sample msg. All OK.\n"))
	return err
}
