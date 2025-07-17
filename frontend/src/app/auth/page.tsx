import { useState } from "react";
import axios from "axios";

export default function Login() {
  //const [username, setUsername] = useState<string>("");
  const [inputValue, setInputValue] = useState<string>("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const data = inputValue.trim();
    const res = await axios.post("http://127.0.0.1:8082/login", data);
    console.log(res);
    if (res.data.response == 200) {
        const {token} = res.data;
        localStorage.setItem('token', token);
    } else {
        console.log("Error: failed to login");
    }
  };

  return (
    <div>
      <form onSubmit={handleSubmit} className="flex mt-4 space-x-2">
        <input
          type="text"
          value={inputValue}
          onChange={(e) => setInputValue(e.target.value)}
          //placeholder={username ? "Send a message…" : "Set your username…"}
          minLength={2}
          maxLength={200}
          required
          className="flex-1 border rounded px-3 py-2 focus:outline-none focus:ring"
        />
        <button
          type="submit"
          className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition"
        >
        </button>
      </form>
    </div>
  );
}
