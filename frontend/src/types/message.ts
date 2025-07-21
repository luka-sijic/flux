type WSMessage = {
  type: "join" | "chat" | "users" | "log" | "ping";
  user?: string;
  user1?: string;
  content?: string;
  users?: Record<string, string>;
  log?: string[];
};

type Message = {
  username: string;
  message: string;
};

type User = {
  id: string;
  username: string;
}

type Friend = {
  friend: string;
};