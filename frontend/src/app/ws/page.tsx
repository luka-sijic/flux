"use client";
import axios from "axios";
import { useState, useEffect, useRef } from "react";
import { Circle, Plus, Handshake } from "lucide-react";
import { useRouter } from "next/navigation";
import toast from "react-hot-toast";
import { jwtDecode } from "jwt-decode";

type WSMessage = {
  type: "join" | "chat" | "users" | "log" | "ping";
  user?: string;
  content?: string;
  users?: Record<string, string>;
  log?: string[];
};

type Friend = {
  friend: string;
};

export default function Home() {
  const router = useRouter();
  const connection = useRef<WebSocket | null>(null);
  const [username, setUsername] = useState<string>("");
  const [friends, setFriends] = useState<Friend[]>([]);
  const [inputValue, setInputValue] = useState<string>("");
  const [friend, setFriendValue] = useState<string>("");
  const [messages, setMessages] = useState<WSMessage[]>([]);
  const [addFriend, setAddFriend] = useState<Boolean>(false);
  const [log, setLog] = useState<string[]>([]);
  const [users, setUsers] = useState<Record<string, string>>({});

  useEffect(() => {
    const fetchFriends = async () => {
      try {
        const response = await axios.get(
          `http://127.0.0.1:8081/friend/${username}`
        );
        if (response.status === 200) {
          setFriends(response.data);
          console.log(response.data);
        }
      } catch (err) {
        console.log("NO FRIENDS");
      }
    };
    fetchFriends();
  }, []);

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
          const byteLen = new TextEncoder().encode(JSON.stringify(data)).length;
          console.log("Length: " + byteLen);
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

  const addNewFriend = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const res = await axios.post("http://127.0.0.1:8081/friend", { friend });
      toast.success("Friend request sent");
      setFriendValue("");
      if (res.status === 200) {
        console.log("friend addded");
      }
    } catch (err) {
      if (axios.isAxiosError(err) && err.response) {
        if (err.response.status === 404) {
          toast.error("User not found 🤔", {
            className: "bg-red-400",
          });
          setFriendValue("");
          console.log("No user exists");
        }
      }
    }
  };

  const handleFriend = async (e: React.FormEvent) => {
    e.preventDefault();
    setAddFriend(!addFriend);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (
      !connection.current ||
      connection.current.readyState !== WebSocket.OPEN
    ) {
      console.error("WebSocket not open");
      return;
    }

    // Afterwards, send chat messages
    const content = inputValue.trim();
    if (content.length) if (!content) return;
    connection.current.send(
      JSON.stringify({ type: "chat", user: username, content })
    );
    setInputValue("");
  };

  const friendRequests = async () => {
    void router.push("/requests");
  };

  return (
    <div className="flex h-screen overflow-x-hidden">
      <div className="w-64 border border-white rounded p-4 bg-black text-white">
        <div className="flex justify-between">
          <h2 className="text-xl font-semibold mb-2">Friends</h2>
          <form onSubmit={handleFriend}>
            <button type="submit">
              <Plus className="cursor-pointer" />
            </button>
            <button onClick={friendRequests}>
              <Handshake size={20} className="cursor-pointer" />
            </button>
          </form>
        </div>
        {addFriend ? (
          <form onSubmit={addNewFriend}>
            <input
              value={friend}
              onChange={(e) => setFriendValue(e.target.value)}
              placeholder="add friend..."
              className="bg-gray-800 rounded p-1"
            />
          </form>
        ) : (
          ""
        )}
        {/* 1) Users column */}
        {/*<ul>
          {Object.keys(users).length === 0 ? (
            <p className="text-gray-500 py-2">No users connected</p>
          ) : (
            Object.entries(users).map(([key, value]) => (
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
            ))
          )}
        </ul>*/}
        <ul>
          {friends.length > 0 ? (
            friends.map((f, i) => <li key={i}>{f.friend}</li>)
          ) : (
            <li>No friends found</li>
          )}
        </ul>
      </div>
      {/* Chat Column */}
      <div className="flex flex-col flex-1 p-4 min-h-0">
        {/* Messages: now fills full height and only scrolls vertically */}
        <div
          className="
        flex-1                /* take up all remaining vertical space */
        overflow-y-auto       /* vertical scrolling as needed */
        overflow-x-hidden     /* hide any horizontal overflow */
        border rounded p-3
        whitespace-pre-wrap
        break-words           /* wrap long words */
        break-all             /* break really long tokens if needed */
        bg-black text-white
      "
        >
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
            placeholder="Send message"
            minLength={2}
            maxLength={200}
            required
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
    </div>
  );
}
