# TCP Chat Application

A barebones TCP-based chat application with a React frontend. This project demonstrates the fundamentals of TCP networking in Go while providing a modern web interface for user interaction.

## Overview

This project consists of two main components:

1. **TCP Server** (Go)
   - Handles raw TCP connections
   - Manages chat rooms and clients
   - Implements basic chat commands
   - Provides WebSocket bridge for web clients

2. **Web Client** (React)
   - Modern user interface for the chat
   - Real-time message updates
   - Room management
   - Connection status monitoring

## Getting Started

### Prerequisites

- Go 1.19 or higher
- Node.js 16 or higher
- Make

### Installation

```bash
# Clone the repository
git clone <repository-url>

# Start both server and client
make run
```

### Available Chat Commands

- `/name <username>` - Set your display name
- `/join <room>` - Join a chat room
- `/rooms` - List all available rooms
- `/msg <message>` - Send a message to your current room
- `/quit` - Disconnect from the chat

## Project Structure

```
.
├── server/         # Go TCP server
│   ├── main.go
│   ├── client.go
│   ├── server.go
│   ├── room.go
│   └── model.go
├── client/         # React frontend
│   ├── src/
│   ├── public/
│   └── package.json
├── Makefile
└── README.md
```

## How It Works

The application uses raw TCP connections for the core chat functionality. When a client connects (either through telnet or the web interface), they can interact with the chat server using text-based commands.

For web users, a WebSocket bridge translates between WebSocket connections and TCP connections, allowing seamless integration with the TCP server.

## Development Status

⚠️ **Note**: This is a barebones implementation focused on demonstrating TCP networking concepts. For a more scalable and feature-rich chat application, check out my upcoming project built with NestJS and GraphQL (coming soon).

## Future Development

I am currently working on a more robust and scalable version of this chat application using:
- NestJS for backend services
- GraphQL for efficient data querying
- Real-time subscriptions
- Proper user authentication
- Message persistence
- And much more!

Stay tuned for the improved version!

## Running with Telnet

You can also connect to the chat server directly using telnet:

```bash
telnet localhost 8080
```

This allows you to interact with the server using raw TCP connections.

## License

MIT License - feel free to use this code for learning and development!