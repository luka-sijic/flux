"use client";

import { useState, useEffect, useRef } from "react";
import { Circle } from "lucide-react";

type WSMessage = {
  type: "join" | "chat" | "users" | "log" | "ping";
  user?: string;
  content?: string;
  users?: Record<string, string>;
  log?: string[];
};

export default function Home() {
  const connection = useRef<WebSocket | null>(null);
  const [username, setUsername] = useState<string>("");
  const [inputValue, setInputValue] = useState<string>("");
  const [messages, setMessages] = useState<WSMessage[]>([]);
  const [log, setLog] = useState<string[]>([]);
  const [users, setUsers] = useState<Record<string, string>>({});

  useEffect(() => {
    const socket = new WebSocket("ws://127.0.0.1:8080/ws");

    socket.onopen = () => {
      console.log("Connection opened");
    };

    socket.onmessage = (event) => {
      const data = JSON.parse(event.data) as WSMessage;

      switch (data.type) {
        case "log":
          console.log(data);
          setLog(data.log ?? []);
        case "users":
          // data.users is string[] or undefined → default to []
          setUsers(data.users ?? {});
          break;

        case "chat":
          // append the new WSMessage to our array
          setMessages((prev) => [...prev, data]);
          break;
        case "ping":
          socket.send(JSON.stringify({ type: "pong" }));
          break;
      }
    };

    connection.current = socket;
    return () => {
      socket.close();
    };
  }, []);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (
      !connection.current ||
      connection.current.readyState !== WebSocket.OPEN
    ) {
      console.error("WebSocket not open");
      return;
    }

    // First submit uses inputValue as the username
    if (!username) {
      const name = inputValue.trim();
      if (!name) return;
      connection.current.send(JSON.stringify({ type: "join", user: name }));
      setUsername(name);
      setInputValue("");
      return;
    }

    // Afterwards, send chat messages
    const content = inputValue.trim();
    if (!content) return;
    connection.current.send(
      JSON.stringify({ type: "chat", user: username, content })
    );
    setInputValue("");
  };

  return (
    <div className="flex h-screen">
      {/* Chat Column */}
      <div className="flex flex-col flex-1 p-4">
        <div className="flex-1 overflow-y-auto border rounded p-3 whitespace-pre-wrap bg-black text-white">
          {log.map((m, i) => (
            <p key={i}>{m}</p>
          ))}
          {messages.map((m, i) => (
            <p key={i}>
              <strong>{m.user}:</strong> {m.content}
            </p>
          ))}
        </div>

        <form onSubmit={handleSubmit} className="flex mt-4 space-x-2">
          <input
            type="text"
            value={inputValue}
            onChange={(e) => setInputValue(e.target.value)}
            placeholder={username ? "Send a message…" : "Set your username…"}
            className="flex-1 border rounded px-3 py-2 focus:outline-none focus:ring"
          />
          <button
            type="submit"
            className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition"
          >
            {username ? "Send" : "Join"}
          </button>
        </form>
      </div>

      {/* Users Column */}
      <div className="w-64 border-l p-4 bg-black">
        <h2 className="text-xl font-semibold mb-2">Users</h2>
        <ul>
          {Object.entries(users).map(([key, value]) => (
            <li key={key} className="flex items-center gap-1">
              <Circle
                size={12}
                stroke="none"
                fill="currentColor"
                className={
                  value === "active" ? "text-green-500" : "text-red-500"
                }
              />
              <span>{key}</span>
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
}
