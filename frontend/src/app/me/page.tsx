"use client";

import { useWebSocket } from "@/hooks/useWebSocket";
import useUser from "@/hooks/useHook";

export default function Me() {
  const { user, loading } = useUser();
  const { isOpen, messages, send } = useWebSocket(
    process.env.NEXT_PUBLIC_WS + "/ws"
  );

  if (loading) return <p>Loading</p>;
  if (!isOpen) return <p>Websocket not open</p>;

  return (
    <div>
      Welcome {user?.username} {user?.id}
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
          send({ type: "chat" });
        }}
        disabled={!isOpen}
      >
        Ping
      </button>
    </div>
  );
}
