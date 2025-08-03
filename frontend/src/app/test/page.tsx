'use client'; // This must be a Client Component

import { useState, useEffect } from 'react';
import CustomMenu from '@/components/CustomMenu';

export default function HomePage() {
  const [menu, setMenu] = useState<{
    visible: boolean;
    x: number;
    y: number;
  }>({
    visible: false,
    x: 0,
    y: 0,
  });

  // This effect closes the menu when you click away
  useEffect(() => {
    const handleClickAway = () => setMenu({ ...menu, visible: false });
    window.addEventListener('click', handleClickAway);
    return () => {
      window.removeEventListener('click', handleClickAway);
    };
  });

  const handleContextMenu = (event: React.MouseEvent) => {
    // 1. Prevent the default browser context menu
    event.preventDefault();

    // 2. Set the menu's position and make it visible
    setMenu({
      visible: true,
      x: event.clientX, // Use clientX/Y for viewport-relative coordinates
      y: event.clientY,
    });
  };

  return (
    <div
      onContextMenu={handleContextMenu}
      className="flex h-screen w-full items-center justify-center bg-gray-100"
    >
      <p className="select-none text-gray-500">Right-click anywhere on this page.</p>
      {menu.visible && <CustomMenu x={menu.x} y={menu.y} />}
    </div>
  );
}
