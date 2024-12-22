package main

import (
    "bufio"
    "log"
    "net"
    "net/http"
    "github.com/gorilla/websocket"
)

type WStoTCPBridge struct {
    tcpAddr  string
    upgrader websocket.Upgrader
}

func NewBridge(tcpAddr string) *WStoTCPBridge {
    return &WStoTCPBridge{
        tcpAddr: tcpAddr,
        upgrader: websocket.Upgrader{
            ReadBufferSize:  1024,
            WriteBufferSize: 1024,
            CheckOrigin: func(r *http.Request) bool {
                return true
            },
        },
    }
}

func (b *WStoTCPBridge) setupHTTP() {
    // Handle WebSocket connections
    http.HandleFunc("/ws", b.handleWebSocket)
    
    // Start HTTP server
    go func() {
        log.Printf("Starting WebSocket server on :8081")
        if err := http.ListenAndServe(":8081", nil); err != nil {
            log.Fatal(err)
        }
    }()
}

func (b *WStoTCPBridge) handleWebSocket(w http.ResponseWriter, r *http.Request) {
    ws, err := b.upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("WebSocket upgrade error: %v", err)
        return
    }
    defer ws.Close()

    // Connect to TCP server
    tcpConn, err := net.Dial("tcp", b.tcpAddr)
    if err != nil {
        log.Printf("TCP connection error: %v", err)
        return
    }
    defer tcpConn.Close()

    // Create error channel
    errChan := make(chan error, 2)

    // Forward WebSocket -> TCP
    go func() {
        for {
            _, message, err := ws.ReadMessage()
            if err != nil {
                errChan <- err
                return
            }
            
            // Add newline for TCP protocol
            message = append(message, '\n')
            
            if _, err := tcpConn.Write(message); err != nil {
                errChan <- err
                return
            }
        }
    }()

    // Forward TCP -> WebSocket
    go func() {
        reader := bufio.NewReader(tcpConn)
        for {
            message, err := reader.ReadString('\n')
            if err != nil {
                errChan <- err
                return
            }
            
            if err := ws.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
                errChan <- err
                return
            }
        }
    }()

    // Wait for any error
    <-errChan
}