"use client";

import { useWebSocket } from "@/hooks/useWebSocket";
import useUser from "@/hooks/useHook";

export default function Me() {
  const { username, loading } = useUser();
  const { isOpen, messages, send } = useWebSocket(
    process.env.NEXT_PUBLIC_WS + "/ws"
  );

  if (loading) return <p>Loading</p>;
  if (!username) return <p>Not logged in</p>;
  if (!isOpen) return <p>Websocket not open</p>;

  return (
    <div>
      Welcome {username}
      <h1>OK</h1>
      {messages.map((m, i) => (
        <p key={i}>
          {m.username}
          {m.content}
        </p>
      ))}
      <button
        className="bg-gray-400"
        onClick={() => {
            console.log("sending");
          send({ type: "ping" });
        }}
        disabled={!isOpen}
      >
        Ping
      </button>
    </div>
  );
}
