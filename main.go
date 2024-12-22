package main

import (
	"log"
	"net"
)

func main() {
	s := newServer()

	go s.run()
	listner, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("unable to start servers %s ", err.Error())
	}

	defer listner.Close()
	log.Print("started server on port :8080")

	for {
		conn, err := listner.Accept()
		if err != nil {
			log.Println("unable to accept connection ", err.Error())
			continue
		}

		go s.newClient(conn)
	}
}
