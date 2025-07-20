import axios from "axios";

const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API,
  withCredentials: true,
});

export const fetchUsername = async (): Promise<string | null> => {
  try {
    const res = await api.get('/me');
    if (res.status === 200) {
      return res.data;
    } else {
      return null;
    }
  } catch (err: any) {
    if (err.response?.status === 500) {
      const ok = await refreshToken();
      if (ok) {
        return "true";
      }
    }
    return null;
  }
};

export const addFriendReq = async (friend: string): Promise<Number> => {
  try {
    const res = await api.post(`/friend`, {friend});
    console.log(res.status);
    return res.status;
  } catch (err: any) {
    if (axios.isAxiosError(err) && err.response) {
      return 500;
    }
    return 500;
  }
};

export const fetchMessages = async (user1: string, user2: string): Promise<Message[] | null> => {
  try {
    const res = await api.get(`/friend/${user1}/${user2}`);
    if (res.status === 200) {
      return res.data;
    } else {
      return null;
    } 
  } catch (err: any) {
    console.log("Issue fetching messages");
    return null;
  }
};

export const refreshToken = async (): Promise<boolean> => {
  try {
    const res = await api.get('/refresh');
    if (res.status === 200) {
      console.log("Token refreshed");
      return true;
    } else {
      console.log("token failed to refresh");
      return false;
    }
  } catch (err: any) {
    console.log("refresh token failed to refresh");
    return false;
  }
};

export const fetchFriends = async (username: string): Promise<Friend[] | null> => {
  try {
    const res = await api.get(`/friend/${username}`);
    if (res.status === 200) {
      return res.data;
    } else {
      return null;
    }
  } catch (err: any) {
    return null;
  }
};
