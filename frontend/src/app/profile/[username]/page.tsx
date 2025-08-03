import { Button } from "@/components/ui/button";

import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Home, UserPlus, Calendar } from "lucide-react";
import axios from 'axios';
import { notFound } from 'next/navigation';
import Link from "next/link";


type Props = {
  params: {
    username: string;
  };
};

async function getUserProfile(username: string) : Promise<User | null> {
  const response = await axios.get(process.env.NEXT_PUBLIC_API + '/profile/' + username);
  if (Number(response.data) === 1) {
    const user: User = {
      id: "1",
      username: username,
    };
    return user;
  } 
  return null;   
};


export default async function Profile({ params }: Props) {
  const user = await getUserProfile(params.username);
  console.log(user);
  if (!user) {
    notFound()
  } 

  const friends = [
    { id: 1, friend: "talent" },
    { id: 2, friend: "WOW" },
  ]
 
  return (
    <div className="min-h-screen bg-black text-white p-4">
      <div className="flex w-60 m-4 gap-2">
        <Link href="/"><Home className="cursor-pointer"/></Link>
      </div>
      <div className="max-w-md mx-auto space-y-6">
        {/* Profile Header */}
        <div className="flex flex-col items-center space-y-4 pt-8">
          <Avatar className="w-32 h-32 border-4 border-gray-700">
            <AvatarImage src={`${process.env.NEXT_PUBLIC_API}/media/avatars/${user.username}.jpg`} alt="Profile picture" />
          </Avatar>

          <div className="text-center">
            <h1 className="text-2xl font-bold">{user.username}</h1>
            <p className="text-gray-400">@{user.username}</p>
          </div>

          <Button className="bg-blue-600 hover:bg-blue-700 text-white px-8 py-2 rounded-lg flex items-center gap-2">
            <UserPlus className="w-4 h-4" />
            Add Friend
          </Button>
        </div>

        {/* Member Since Box */}
        <Card className="bg-zinc-950 border-zinc-950">
          <CardHeader className="pb-3">
            <div className="flex items-center gap-2 text-gray-300">
              <Calendar className="w-4 h-4" />
              <span className="text-sm font-medium">Member Since</span>
            </div>
          </CardHeader>
          <CardContent className="pt-0">
            <p className="text-white font-semibold">March 15, 2020</p>
          </CardContent>
        </Card>

        {/* Friends List Box */}
        <Card className="bg-zinc-950 border-zinc-950">
          <CardHeader className="pb-3">
            <h3 className="text-lg font-semibold text-white">Friends ({friends?.length})</h3>
          </CardHeader>
          <CardContent className="pt-0">
            <div className="space-y-3 max-h-64 overflow-y-auto">
              {friends?.map((friend) => (
                <div
                  key={friend.id}
                  className="flex items-center gap-3 p-2 rounded-lg hover:bg-gray-800 transition-colors"
                >
                  <Avatar className="w-10 h-10">
                    <AvatarImage src={`${process.env.NEXT_PUBLIC_API}/media/avatars/${friend.friend}.jpg`} alt="Profile picture" />
                      <AvatarFallback className="bg-gray-700 text-white text-sm">
                      {friend.friend
                        .split(" ")
                        .map((n) => n[0])
                        .join("")}
                    </AvatarFallback>
                  </Avatar>
                  <div className="flex-1 min-w-0">
                    <p className="text-white font-medium truncate">{friend.friend}</p>
                    <p className="text-gray-400 text-sm truncate">@{friend.friend}</p>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )}
  
