package main

import (
	"log"
	"net"
)

const PORT = ":8888"

func main() {
	s := newServer()
	go s.run()

	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("unable to start server: %s", err.Error())
	}

	defer listener.Close()
	log.Printf("started server on " + PORT)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("unable to accept connection: %s", err.Error())
			continue
		}

		go s.newClient(conn)
	}
}

