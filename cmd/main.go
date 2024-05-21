package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

const (
	localhost    = "127.0.0.1"
	port         = "8333"
	waitInterval = 10 * time.Second
)

const (
	msgHello         = "1. Hello! Here are my encryption methods.\n"
	msgPublicKey     = "2. Hello! Here is my public key\n"
	msgSecretKey     = "3. Here is encrypted secret-key.\n"
	msgSampleMessage = "5. This is encrypted sample message.\n"
	msgVerified      = "6. Verified sample msg. All OK.\n"
)

func main() {
	listener, remotePort := bindPort()
	if listener == nil {
		return
	}

	done := make(chan struct{})
	go func() {
		if err := startServer(listener); err != nil {
			fmt.Println(err)
		}
		close(done)
	}()

	go func() {
		if err := startClient(remotePort); err != nil {
			fmt.Println(err)
		}
		close(done)
	}()

	<-done
}

func startServer(listener net.Listener) error {
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		if err := processIncoming(conn); err != nil {
			fmt.Println("handshake failed, try again.", err)
		} else {
			fmt.Println("handshake success.")
			return nil
		}
	}
}

func startClient(remotePort uint16) error {
	for {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", localhost, remotePort))
		if err == nil {
			if err := doHandshake(conn); err != nil {
				fmt.Println(err)
			} else {
				return nil
			}
		}

		fmt.Println("Waiting for remote")
		time.Sleep(waitInterval)
	}
}

func bindPort() (net.Listener, uint16) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", localhost, port))
	if err != nil {
		fmt.Println("Error binding port", err)
		return nil, 0
	}

	return listener, 8333
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

func doHandshake(conn net.Conn) error {
	messages := []string{msgHello, msgSecretKey, msgSampleMessage}
	for _, msg := range messages {
		if _, err := conn.Write([]byte(msg)); err != nil {
			return err
		}

		if err := readLine(conn); err != nil {
			return err
		}
	}

	return nil
}

func processIncoming(conn net.Conn) error {
	messages := []string{msgPublicKey, msgVerified}
	for _, msg := range messages {
		if _, err := conn.Write([]byte(msg)); err != nil {
			return err
		}

		if err := readLine(conn); err != nil {
			return err
		}
	}

	return nil
}
