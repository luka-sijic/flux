"use client";

import { useEffect, useState, FormEvent } from "react";
import axios from "axios";
import { useRouter } from "next/navigation";
import { Check, X } from "lucide-react";
import toast from 'react-hot-toast';

type Friends = {
  friend: string;
};

export default function Requests() {
  const router = useRouter();
  const [friends, setFriends] = useState<Friends[]>([]);
  const [token, setToken] = useState<string>("");

  useEffect(() => {
    const token = localStorage.getItem("token");
    if (!token) return router.replace("/auth");
    setToken(token);
  }, []);

  useEffect(() => {
    if (!token) return;
    const fetchFriendReq = async () => {
      const response = await axios.get("http://127.0.0.1:8081/friend", {
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
      });
      console.log(response);
      if (response.status === 200) {
        setFriends(response.data);
      }
    };
    fetchFriendReq();
  }, [token]);

  const friendResponse = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const action = (e.nativeEvent as SubmitEvent)
      .submitter as HTMLButtonElement;

    // collect the form fields
    const form = e.currentTarget;
    const data = new FormData(form);
    data.set("action", action.value);

    const response = await axios.post(
      "http://127.0.0.1:8081/friend/respond",
      data,
      {
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
      }
    );
    if (response.status === 200) {
      setFriends((prev) => prev.filter(f => f.friend !== data.get('friendId')));
    }
  };

  return (
    <div>
      <h2>Friend Requests</h2>
      <div className="bg-gray-700 p-2 m-2 w-64 rounded">
        {" "}
        {/* smaller box */}
        <ul className="space-y-1">
          {" "}
          {/* tiny gap between rows */}
          {friends.length === 0 ? (
            <p>No pending friend requests</p>
          ) : (
            friends.map((f, i) => (
              <li key={i} className="flex items-center">
                <span className="truncate">{f.friend}</span>{" "}
                {/* name on the left */}
                {/* icons on the far right, with space between them */}
                <div className="ml-auto flex gap-3">
                  <form onSubmit={friendResponse}>
                    <input type="hidden" name="friendId" value={f.friend} />
                    <button value="accepted">
                      <Check className="cursor-pointer" />
                    </button>
                    <button value="rejected">
                      <X className="cursor-pointer" />
                    </button>
                  </form>
                </div>
              </li>
            ))
          )}
        </ul>
      </div>
    </div>
  );
}
