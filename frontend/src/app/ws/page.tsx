"use client";
import axios from "axios";
import { useState, useEffect, useRef } from "react";
import { Heart, Plus, Handshake } from "lucide-react";
import { useRouter } from "next/navigation";
import toast from "react-hot-toast";
import refreshToken from "@/lib/auth";
axios.defaults.withCredentials = true;

type WSMessage = {
  type: "join" | "chat" | "users" | "log" | "ping";
  user?: string;
  user1?: string;
  content?: string;
  users?: Record<string, string>;
  log?: string[];
};

type Friend = {
  friend: string;
};

type Message = {
  username: string;
  message: string;
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
  const [log, setLog] = useState<Message[]>([]);
  const [active, setActive] = useState<string>("");
  const [liked, setLiked] = useState<Set<number>>(new Set());
  const [users, setUsers] = useState<Record<string, string>>({});

  useEffect(() => {
    const fetchUsername = async () => {
      try {
        const res = await axios.get("http://127.0.0.1:8081/me");
        console.log(res);
        if (res.status === 200) {
          setUsername(res.data);
        } else {
          console.log("can't find username");
        }
      } catch (err: any) {
        if (err.response?.status === 500) {
          const ok = refreshToken();
          if (ok) {
            router.replace("/ws");
          }
        }
        console.log(err);
      }
    };
    fetchUsername();
  }, []);

  useEffect(() => {
    if (!username) return;
    const fetchFriends = async () => {
      try {
        const response = await axios.get(
          `http://127.0.0.1:8081/friend/${username}`
        );
        if (response.status === 200) {
          console.log("RESPONSE: ", response.data);
          setFriends(response.data);
          console.log(response.data.length);
          if (response.data.length > 0) {
            setActive(response.data[0].friend);
          }
        }
      } catch (err) {
        console.log("NO FRIENDS");
      }
    };
    fetchFriends();
  }, [username]);

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
        //setLog(data.log ?? []);
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
      JSON.stringify({ type: "chat", user: username, user2: active, content })
    );
    setInputValue("");
  };

  const friendRequests = async () => {
    void router.push("/requests");
  };

  /*const getMessages = async (friend: string) => {
    try {
      console.log("Active username: ", friend);
      setActive(friend);
      const response = await axios.get(
        `http://127.0.0.1:8081/friend/${username}/${friend}`
      );
      console.log("MESSAGE GET RESPONSE: ", response);
      if (response.status === 200) {
        setMessages([]);
        setLog(response.data);
      }
    } catch (err) {
      console.log("ERROR");
    }
  };*/

  useEffect(() => {
    if (!active) return;
    const fetchMessages = async () => {
      try {
        console.log("Active username: ", active);
        const response = await axios.get(
          `http://127.0.0.1:8081/friend/${username}/${active}`
        );
        console.log("MESSAGE GET RESPONSE: ", response);
        if (response.status === 200) {
          setMessages([]);
          setLog(response.data);
        }
      } catch (err) {
        console.log("ERROR");
      }
    };
    fetchMessages();
  }, [active]);

  const likeMessage = async (id: number) => {
    setLiked((prev) => {
      const next = new Set(prev);
      if (next.has(id)) {
        next.delete(id);
      } else {
        next.add(id);
      }
      return next;
    });
  };

  return (
    <div className="flex h-screen overflow-x-hidden">
      <div className="w-64 border border-white rounded p-4 bg-black text-white">
        <div className="flex justify-between">
          <h2 className="text-xl font-semibold mb-4">Friends</h2>
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
            friends.map((f, i) => (
              <li
                key={i}
                className="flex justify-between w-22 text-gray-300 text-xl cursor-pointer"
              >
                <button
                  type="button"
                  onClick={() => setActive(f.friend)}
                  className="flex items-center gap-2 cursor-pointer"
                >
                  <img
                    className="w-7 h-7 rounded-2xl"
                    src="https://avatars.akamai.steamstatic.com/7d88fb593b5030a1d1d2cfb8b05d282bc07fc389_full.jpg"
                  />
                  {f.friend}
                </button>
              </li>
            ))
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
          <p>Load chat for: {active}</p>
          <ul>
            {log.map((m, i) => (
              <li key={i} className="flex gap-1">
                <Heart
                  onClick={() => likeMessage(i)}
                  className="w-3 h3"
                  fill={liked.has(i) ? "pink" : "none"}
                />
                {m.username}: {m.message}
              </li>
            ))}
            {messages.map((m, i) => (
              <li key={i}>
                <strong>{m.user}:</strong> {m.content}
              </li>
            ))}
          </ul>
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
