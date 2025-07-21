import {useState, useEffect} from "react";
import {fetchFriends, fetchMessages, fetchUser, getFriendReq} from "../services/service";

export default function useUser() {
    const [user, setUser] = useState<User>();
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        fetchUser()
            .then((u) => setUser(u))
            .finally(() => setLoading(false))
    }, []);

    return {user, loading};
}

export function getFriends(id: string | null) {
    const [friends, setFriends] = useState<Friend[] | null>([]);
    const [active, setActive] = useState<string | null>(null);

    useEffect(() => {
        if (!id) return;
        fetchFriends(id)
            .then((f) => {
                setFriends(f)
                if (f && f.length > 0) setActive(f[0].friend);
            })
            .catch(() => setFriends([]))
    }, [id]);

    return {friends, active, setActive};
}

export function getMessages(username: string | null, active: string | null) {
    const [log, setLog] = useState<Message[] | null>([]);
    const [messages, setMessages] = useState<WSMessage[] | null>([]);

    useEffect(() => {
        if (!active || !username) return;
        fetchMessages(username, active)
            .then((l) => {
                setLog(l);
            })
            .catch(() => setLog([]));
    }, [active]);

    return {log, setLog};
}

export function getFriendRequests() {
    const [friends, setFriends] = useState<Friend[]>([]);

    useEffect(() => {
        getFriendReq()
            .then((f) => setFriends(f))
            .catch(() => setFriends([]));
    }, []);

    return {friends, setFriends};
}