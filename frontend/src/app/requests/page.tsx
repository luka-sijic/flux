"use client";

import { useEffect, useState, FormEvent } from "react";
import axios from "axios";
import { useRouter } from "next/navigation";
import { Check, X, Home } from "lucide-react";
import useUser, { getFriendRequests } from "@/hooks/useHook";

export default function Requests() {
  const router = useRouter();
  const { username } = useUser();
  const { friends, setFriends } = getFriendRequests();

  if (!username) return <p>User not logged in</p>;
  //if (!friends) return <p>No friend getFriendRequests</p>

  const friendResponse = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const action = (e.nativeEvent as SubmitEvent)
      .submitter as HTMLButtonElement;

    // collect the form fields
    const form = e.currentTarget;
    const data = new FormData(form);
    data.set("action", action.value);

    const response = await axios.post(
      process.env.NEXT_PUBLIC_API + "friend/respond",
      data
    );
    if (response.status === 200) {
      setFriends((prev) =>
        prev.filter((f) => f.friend !== data.get("friendId"))
      );
    }
  };

  const handleBack = () => {
    void router.push("/ws");
  };

  return (
    <div>
      <div className="flex w-60 m-4 gap-2">
        <Home className="cursor-pointer" onClick={handleBack} />
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
                <form onSubmit={friendResponse} className="ml-auto flex gap-3">
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
