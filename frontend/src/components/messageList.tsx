import { useState } from "react";
import { Heart } from "lucide-react";

interface Props {
  log: Message[];
}

export default function MessageList({ log }: Props) {
  const [liked, setLiked] = useState<Set<number>>(new Set());

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
      {/*{messages.map((m, i) => (
              <li key={i}>
                <strong>{m.user}:</strong> {m.content}
              </li>
      ))*/}
    </ul>
  );
}
