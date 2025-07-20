interface Props {
  friends: Friend[];
  setActive: (id: string) => void;
}

export default function FriendsList({ friends, setActive }: Props) {
  return (
    <ul>
      {friends.length > 0 ? (
        friends.map((f, i) => (
          <li
            key={i}
            className="flex justify-between w-22 text-gray-300 text-xl cursor-pointer"
          >
            <button
              type="button"
              onClick={() => setActive(f.friend)}
              className="flex items-center gap-2 cursor-pointer"
            >
              <img
                className="w-7 h-7 rounded-2xl"
                src="https://avatars.akamai.steamstatic.com/7d88fb593b5030a1d1d2cfb8b05d282bc07fc389_full.jpg"
              />
              {f.friend}
            </button>
          </li>
        ))
      ) : (
        <li>No friends found</li>
      )}
    </ul>
  );
}
