package main

import (
	"log"
	"net"
)

func main() {
    // Start your existing TCP server
    s := newServer()
    go s.run()

    // Create and start the bridge
    bridge := NewBridge(":8080")
    bridge.setupHTTP()

    // Continue with your TCP server setup
    listener, err := net.Listen("tcp", ":8080")
    if err != nil {
        log.Fatalf("unable to start server: %s", err.Error())
    }
    defer listener.Close()
    log.Print("started server on port :8080")

    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Println("unable to accept connection", err.Error())
            continue
        }
        go s.newClient(conn)
    }
}
