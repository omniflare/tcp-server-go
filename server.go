package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
)

func newServer() *server {
	return &server{
		rooms:    make(map[string]*room),
		commands: make(chan command),
	}
}

func (s *server) run() {
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_NICK:
			s.name(cmd.client, cmd.args)
		case CMD_JOIN:
			s.join(cmd.client, cmd.args)
		case CMD_MSG:
			s.msg(cmd.client, cmd.args)
		case CMD_QUIT:
			s.quit(cmd.client)
		case CMD_ROOMS:
			s.listrooms(cmd.client)
		}
	}
}

func (s *server) newClient(conn net.Conn) {
	log.Printf("new client has connected now %s", conn.RemoteAddr().String())
	c := &client{
		conn:     conn,
		name:     "anonymus",
		commands: s.commands,
	}
	c.readInput()
}

func (s *server) name(c *client, args []string) {
	c.name = args[1]
	c.msg(fmt.Sprintf("so you're %s", c.name))
}

func (s *server) join(c *client, args []string) {
	if len(args) < 2 {
		c.msg("room name is required in usage :   /join ROOM_NAME")
		return
	}
	roomName := args[1]
	r, ok := s.rooms[roomName]
	if !ok {
		r = &room{
			name:    roomName,
			members: make(map[net.Addr]*client),
		}
		s.rooms[roomName] = r
	}
	r.members[c.conn.RemoteAddr()] = c

	s.quitCurrentRoom(c)
	c.room = r

	r.broadcast(c, fmt.Sprintf("%s joined the room", c.name))

	c.msg(fmt.Sprintf("welcome to %s", roomName))
}

func (s *server) quitCurrentRoom(c *client) {

	if c.room != nil {
		delete(c.room.members, c.conn.RemoteAddr())
		c.room.broadcast(c, fmt.Sprintf("%s left", c.name))
	}
}

func (s *server) msg(c *client, args []string) {
	if c.room == nil {
		c.err(errors.New("you must join a room first to send message "))
		return
	}
	c.room.broadcast(c, c.name+" : "+strings.Join(args[1:], " "))
}
func (s *server) quit(c *client) {

	log.Printf("quitting %s ", c.conn.RemoteAddr().String())
	s.quitCurrentRoom(c)

	c.msg("bye :)")
	c.conn.Close()
}
func (s *server) listrooms(c *client) {

	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}
	c.msg(fmt.Sprintf("available rooms are :  %s ", strings.Join(rooms, ", ")))
}
