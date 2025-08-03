
'use client'; // This component is interactive
import { useState, useEffect } from 'react';
import CustomMenu from '@/components/CustomMenu'; // Assuming you have this

interface Props {
  friends: Friend[];
  setActive: (id: string) => void;
}

export default function FriendsList({ friends, setActive }: Props) {
  const [menu, setMenu] = useState<{
    visible: boolean;
    x: number;
    y: number;
    selectedId: string | null;
  }>({
    visible: false,
    x: 0,
    y: 0,
    selectedId: null,
  });

  useEffect(() => {
    const handleClickAway = () => setMenu({ ...menu, visible: false });
    window.addEventListener('click', handleClickAway);
    return () => window.removeEventListener('click', handleClickAway);
  });

  const handleContextMenu = (event: React.MouseEvent, friendId: string) => {
    event.preventDefault();
    setMenu({
      visible: true,
      x: event.clientX,
      y: event.clientY,
      selectedId: friendId,
    });
  };

  return (
    <>
      <ul>
        {friends.map((f) => (
          <li
            key={f.friend}
            onContextMenu={(e) => handleContextMenu(e, f.friend)}
            className="flex w-full cursor-pointer justify-between text-xl text-gray-300"
          >
            <button
              type="button"
              onClick={() => setActive(f.friend)}
              className="flex w-full items-center gap-2"
            >
              <img
                className="h-7 w-7 rounded-2xl"
                src={`${process.env.NEXT_PUBLIC_API}/media/avatars/${f.friend}.jpg`}
                alt={`${f.friend}'s avatar`}
              />
              {f.friend}
            </button>
          </li>
        ))}
      </ul>

      {/* 5. Render the menu */}
      {menu.visible && menu.selectedId && (
        <CustomMenu x={menu.x} y={menu.y} selectedId={menu.selectedId} />
      )}
    </>
  );
}
