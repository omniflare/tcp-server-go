import { useState, useEffect, useRef } from 'react';

export default function ChatApp() {
    const [messages, setMessages] = useState([]);
    const [inputValue, setInputValue] = useState('');
    const [connected, setConnected] = useState(false);
    const [currentRoom, setCurrentRoom] = useState('');
    const [username, setUsername] = useState('anonymous');
    const wsRef = useRef(null);
    const messagesEndRef = useRef(null);

    useEffect(() => {
        connectWebSocket();
        return () => {
            if (wsRef.current) {
                wsRef.current.close();
            }
        };
    }, []);

    const connectWebSocket = () => {
        const ws = new WebSocket('ws://localhost:8081/ws');

        ws.onopen = () => {
            setConnected(true);
            addMessage('System', 'Connected to chat server');
        };

        ws.onclose = () => {
            setConnected(false);
            addMessage('System', 'Disconnected from chat server');
            // Attempt to reconnect after 2 seconds
            setTimeout(connectWebSocket, 2000);
        };

        ws.onmessage = (event) => {
            addMessage('Server', event.data);
            // Handle server responses
            if (event.data.includes('welcome to')) {
                const roomName = event.data.split('welcome to ')[1];
                setCurrentRoom(roomName.trim());
            }
        };

        wsRef.current = ws;
    };

    const addMessage = (sender, text) => {
        setMessages(prev => [...prev, { sender, text, timestamp: new Date() }]);
    };

    const scrollToBottom = () => {
        messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
    };

    useEffect(scrollToBottom, [messages]);

    const handleSubmit = (e) => {
        e.preventDefault();
        if (!inputValue.trim() || !connected) return;

        const command = inputValue.trim();

        // Handle commands
        if (command.startsWith('/name')) {
            setUsername(command.split(' ')[1] || 'anonymous');
        } else if (command.startsWith('/join')) {
            const roomName = command.split(' ')[1];
            if (roomName) setCurrentRoom(roomName);
        }

        wsRef.current.send(command);
        setInputValue('');
    };

    return (
        <div className="max-w-2xl mx-auto p-4">
            <div className="mb-4">
                <div className={`p-2 rounded ${connected ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}`}>
                    Status: {connected ? 'Connected' : 'Disconnected'}
                    {currentRoom && ` - Room: ${currentRoom}`}
                </div>
            </div>

            <div className="bg-gray-100 p-4 rounded mb-4">
                <h3 className="font-bold mb-2">Commands:</h3>
                <ul className="space-y-1 text-sm">
                    <li>/name &lt;username&gt; - Set your username</li>
                    <li>/join &lt;room&gt; - Join a chat room</li>
                    <li>/rooms - List available rooms</li>
                    <li>/msg &lt;message&gt; - Send a message</li>
                    <li>/quit - Disconnect from chat</li>
                </ul>
            </div>

            <div className="h-96 border rounded p-4 mb-4 overflow-y-auto bg-white">
                {messages.map((msg, idx) => (
                    <div key={idx} className="mb-2">
                        <span className="text-gray-500 text-xs">
                            {msg.timestamp.toLocaleTimeString()}
                        </span>
                        <span className="ml-2 font-semibold">{msg.sender}:</span>
                        <span className="ml-2">{msg.text}</span>
                    </div>
                ))}
                <div ref={messagesEndRef} />
            </div>

            <form onSubmit={handleSubmit} className="flex gap-2">
                <input
                    type="text"
                    value={inputValue}
                    onChange={(e) => setInputValue(e.target.value)}
                    placeholder="Type command or message..."
                    className="flex-1 p-2 border rounded"
                    disabled={!connected}
                />
                <button
                    type="submit"
                    disabled={!connected}
                    className="px-4 py-2 bg-blue-500 text-white rounded disabled:bg-gray-300"
                >
                    Send
                </button>
            </form>
        </div>
    );
}