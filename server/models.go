package main

import "net"

const (
	CMD_NICK commandId = iota
	CMD_JOIN
	CMD_ROOMS
	CMD_MSG
	CMD_QUIT
)

type commandId int

type command struct {
	id     commandId
	client *client
	args   []string
}

type client struct {
	conn     net.Conn
	name     string
	room     *room
	commands chan<- command
}

type room struct {
	name    string
	members map[net.Addr]*client
}
type server struct {
	rooms    map[string]*room
	commands chan command
}
