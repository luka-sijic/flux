"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import axios from "axios";

export default function Login() {
  const router = useRouter();
  //const [username, setUsername] = useState<string>("");
  const [username, setUsername] = useState<string>("");
  const [password, setPassword] = useState<string>("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const res = await axios.post("http://127.0.0.1:8081/login", {
      username,
      password,
    });
    if (res.status === 200) {
      router.push("/ws");
    } else {
      console.log("Error: failed to login");
    }
  };

  return (
    <div className="flex items-center justify-center min-h-screen">
      <form
        onSubmit={handleSubmit}
        className="flex flex-col max-w-md space-y-4"
      >
        <input
          type="text"
          value={username}
          placeholder="Username"
          onChange={(e) => setUsername(e.target.value)}
          //placeholder={username ? "Send a message…" : "Set your username…"}
          minLength={2}
          maxLength={25}
          required
          className="border rounded px-3 py-2 focus:outline-none focus:ring"
        />
        <input
          type="text"
          value={password}
          placeholder="Password"
          onChange={(e) => setPassword(e.target.value)}
          //placeholder={username ? "Send a message…" : "Set your username…"}
          minLength={2}
          maxLength={25}
          required
          className="border rounded px-3 py-2 focus:outline-none focus:ring"
        />
        <button
          type="submit"
          className="self-start px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition"
        >
          Submit
        </button>
      </form>
    </div>
  );
}
