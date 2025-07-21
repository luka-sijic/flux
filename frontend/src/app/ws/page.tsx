"use client";
import { useState, useRef } from "react";
import { useRouter } from "next/navigation";
import useUser, { getFriends, getMessages } from "@/hooks/useHook";
import AddFriendForm from "@/components/addFriend";
import FriendsList from "@/components/friendsList";
import MessageList from "@/components/messageList";

export default function Home() {
  const router = useRouter();
  // New
  const { user, loading } = useUser();
  const { friends, active, setActive } = getFriends(user ? user.id : null);
  const { log, setLog } = getMessages(user ? user.username : null, active);

  const connection = useRef<WebSocket | null>(null);
  const [inputValue, setInputValue] = useState<string>("");
  // Todo
  const [messages, setMessages] = useState<WSMessage[]>([]);
  //const [users, setUsers] = useState<Record<string, string>>({});

  if (loading) return <p>Loading</p>;
  if (!user) return <p>Not logged in</p>;
  if (!friends) return <p>No friends</p>;
  if (!log) return <p>No messages</p>;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (
      !connection.current ||
      connection.current.readyState !== WebSocket.OPEN
    ) {
      console.error("WebSocket not open");
      return;
    }

    const content = inputValue.trim();
    if (content.length) if (!content) return;
    connection.current.send(
      JSON.stringify({ type: "chat", user: user.username, user2: active, content })
    );
    setInputValue("");
  };

  return (
    <div className="flex h-screen overflow-x-hidden">
      <div className="w-64 border border-gray rounded p-4 bg-black text-white">
        <AddFriendForm />
        <FriendsList friends={friends} setActive={setActive} />
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
          <MessageList log={log} />
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
            {user.username ? "Send" : "Join"}
          </button>
        </form>
      </div>
    </div>
  );
}
