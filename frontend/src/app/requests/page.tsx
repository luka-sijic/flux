"use client";

import { useEffect, useState, FormEvent } from "react";
import axios from "axios";
import { useRouter } from "next/navigation";
import { Check, X, Home } from "lucide-react";
import toast from 'react-hot-toast';

type Friends = {
  friend: string;
};

axios.defaults.withCredentials = true

export default function Requests() {
  const router = useRouter();
  const [username, setUsername] = useState<string>('');
  const [friends, setFriends] = useState<Friends[]>([]);
  const [token, setToken] = useState<string>("");

  useEffect(() => {
    const fetchUsername = async () => {
      try {
        const res = await axios.get('http://127.0.0.1:8081/me', {withCredentials: true});
        console.log(res);
        if (res.status === 200) {
          setUsername(res.data);
        } else {
          console.log("can't find username");
          router.replace('/auth');
        }
      } catch (err) {
        console.log(err);
      }
    }
    fetchUsername();
  }, []);

  useEffect(() => {
    const fetchFriendReq = async () => {
      const response = await axios.get("http://127.0.0.1:8081/friend");
      console.log(response);
      if (response.status === 200) {
        setFriends(response.data);
      }
    };
    fetchFriendReq();
  }, []);

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

  const handleBack = () => {
    void router.push('/ws');
  }

  return (
    <div>
      <div className="flex w-60 m-4 gap-2">
        <Home 
          className="cursor-pointer"
          onClick={handleBack}
        />
        <h2>Friend Requests</h2>
      </div>
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
              <li key={i} className="flex items-center m-2">
                <span className="truncate">{f.friend}</span>
                  <form 
                    onSubmit={friendResponse} 
                    className="ml-auto flex gap-3"
                  >
                    <input type="hidden" name="friendId" value={f.friend} />
                    <button value="accepted">
                      <Check className="bg-green-400 cursor-pointer" />
                    </button>
                    <button value="rejected">
                      <X className="bg-red-400 cursor-pointer" />
                    </button>
                  </form>
              </li>
            ))
          )}
        </ul>
      </div>
    </div>
  );
}
