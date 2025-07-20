import { useEffect, useRef, useState, useCallback } from "react";

export interface Message {
    type: string;
    username?: string;
    content?: string;
}

export function useWebSocket(url: string) {
    const socketRef = useRef<WebSocket>(null);
    const [isOpen, setIsOpen] = useState<Boolean>(false);
    const [messages, setMessages] = useState<Message[]>([]);

    useEffect(() => {
        const ws = new WebSocket(url);
        console.log(url);
        socketRef.current = ws;

        ws.onopen = () => {
            console.log("WS OPEN");
            setIsOpen(true);
        }
        ws.onclose = () => setIsOpen(false);
        ws.onmessage = (e) => {
            const msg: Message = JSON.parse(e.data);
            setMessages(prev => [...prev, msg]);
        }
        ws.onerror = (e) => {
            console.log("WS ERROR", e);
        }
        ws.onclose = (e) => {
            console.warn("WS CLOSED");
            setIsOpen(false);
        }

        return () => ws.close();
    }, [url]);

    const send = useCallback((msg: Message) => {
        if (socketRef.current?.readyState === WebSocket.OPEN) {
            socketRef.current.send(JSON.stringify(msg));
        }
    }, []);

    return { isOpen, messages, send};
}