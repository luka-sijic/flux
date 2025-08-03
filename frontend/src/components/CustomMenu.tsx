import Link from 'next/link';

type CustomMenuProps = {
  x: number;
  y: number;
  selectedId: string;
};

export default function CustomMenu({ x, y, selectedId }: CustomMenuProps) {
  const menuStyle = {
    top: y, // Y position from cursor
    left: x, // X position from cursor
  };

  return (
    <div
      className="absolute z-50 rounded-md bg-zinc-950 p-2 shadow-lg ring-1 ring-black ring-opacity-5"
      style={menuStyle}
    >
      <ul>
        <li>
          <Link
            href={`/profile/${selectedId}`}
            className="block w-full px-4 py-2 text-left text-sm text-white-700 hover:bg-gray-500"
          >
            Go to Profile
          </Link>
        </li>
        {/* Add other menu items here */}
      </ul>
    </div>
  );
}
