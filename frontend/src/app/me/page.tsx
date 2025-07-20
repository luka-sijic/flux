"use client";

import useUser  from '../../hooks/useHook';

export default function Me() {
    const {username, loading} = useUser();

    if (loading) return <p>Loading</p>
    if (!username) return <p>Not logged in</p>

    return (
        <div>
            Welcome {username}
        </div>
    )
}