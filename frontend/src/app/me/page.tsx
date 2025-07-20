"use client";

import { useState, useEffect } from "react";
import axios from "axios";

//axios.defaults.withCredentials = true

export default function Me() {
    const [me, setMe] = useState<string>('');

    useEffect(() => {
        const getMe = async () => {
            const respon2 = await axios.get("http://127.0.0.1:8081/cookie")
            console.log(respon2)
            const response = await axios.get("http://127.0.0.1:8081/me", {withCredentials: true})
            console.log(response);
            if (response.status === 200) {
                setMe(response.data);
            } else {
                console.log("Cookie problem :D");
            }
        }
        getMe();
    }, []);
    return (
        <div>{me}</div>
    )
}