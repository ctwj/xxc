import { api } from "@/lib/api";

export const authService = {
  login: async (username: string, password: string) => {
    return api.login(username, password);
  },

  logout: async () => {
    return api.logout();
  },

  getCurrentUser: async () => {
    return api.getCurrentUser();
  },

  register: async (username: string, email: string, password: string) => {
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/auth/register`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username, email, password }),
    });
    return response.json();
  },
};