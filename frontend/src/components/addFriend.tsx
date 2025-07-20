"use client";

import { useState } from "react";
import toast from "react-hot-toast";
import { addFriendReq } from "@/services/service";
import { Plus, Handshake } from "lucide-react";
import { useRouter } from "next/navigation";

export default function AddFriendForm() {
  const router = useRouter();
  const [friend, setFriend] = useState<string>("");
  const [addFriend, setAddFriend] = useState<Boolean>(false);

  const addNewFriend = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const res = await addFriendReq(friend);
      if (res === 200) {
        toast.success("Friend request sent");
        setFriend("");
        console.log("friend addded");
      } else {
        toast.error("User not found 🤔");
      }
    } catch (err) {
      setFriend("");
      console.log("No user exists");
    }
  };

  return (
    <div>
      <div className="flex justify-between">
        <h2 className="text-xl font-semibold mb-4">Friends</h2>
        <form
          onSubmit={(e) => {
            e.preventDefault();
            setAddFriend(!addFriend);
          }}
        >
          <button type="submit">
            <Plus className="cursor-pointer" />
          </button>
          <button
            onClick={() => {
              void router.push("/requests");
            }}
          >
            <Handshake size={20} className="cursor-pointer" />
          </button>
        </form>
      </div>
      {addFriend ? (
        <form onSubmit={addNewFriend}>
          <input
            value={friend}
            onChange={(e) => setFriend(e.target.value)}
            placeholder="add friend..."
            className="bg-gray-800 rounded p-1"
          />
        </form>
      ) : (
        ""
      )}
    </div>
  );
}
