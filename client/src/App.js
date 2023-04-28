import { useEffect, useState, useCallback } from 'react';
import { useWebSocket } from 'react-use-websocket/dist/lib/use-websocket';
import './App.css';

function App() {

    const [messages, setMessages] = useState([])
    const [socketUrl] = useState("ws://localhost:9000/ws/chat/1");
    const { sendMessage, lastMessage } = useWebSocket(socketUrl, {share: true})

    useEffect(() => {
            if (lastMessage !== null) {
                setMessages(prev => prev.concat(lastMessage))
            }
    }, [lastMessage, setMessages])

    const renderMessages = () => {
        return messages.map((message, i) =>
            <div key={i}>{message.data}</div>
        )
    }

    const handleSendMessage = useCallback((msg) => sendMessage(msg), [sendMessage]);

    const handleSubmit = event => {
        event.preventDefault()

        handleSendMessage(event.target[0].value)
        event.target[0].value = ""
    }

    return (
        <div className="App">
            <header className="App-header">
                {renderMessages()}
                <form onSubmit={handleSubmit}>
                    <input type="text" name="send-message"></input>
                </form>
            </header>
        </div>
    );
}

export default App;
