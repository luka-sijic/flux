import axios from "axios";

axios.defaults.withCredentials = true

export async function refreshToken() {
    const response = await axios.get("http://127.0.0.1:8081/refresh");
    if (response.status === 200) {
        return true;
    }
    return false;
}

export default function GetUser() {
    const token = localStorage.getItem('token');
    if (!token) {
        return null;
    }
}