import axios from "axios";
import type {
  BassineScore,
  BassineScoreHistory,
  Club,
  DashboardStats,
  Event,
  MenuItem,
  MenuItemReview,
  UpdateBassineScoreRequest,
  User,
} from "./types";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:3000";

const api = axios.create({
  baseURL: API_BASE_URL,
});

// Add auth token to requests
api.interceptors.request.use((config) => {
  if (typeof window !== "undefined") {
    const token = localStorage.getItem("adminToken");
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
  }
  return config;
});

export const authApi = {
  login: async (email: string, password: string) => {
    const response = await api.post("/auth/login", { email, password });
    return response.data;
  },
  verify: async (token: string) => {
    const response = await api.get("/newf/me", {
      headers: { Authorization: `Bearer ${token}` },
    });
    return response.data;
  },
};

export const usersApi = {
  getAll: async (): Promise<User[]> => {
    const response = await api.get("/admin/users");
    return response.data;
  },
  create: async (user: Partial<User>) => {
    const filteredUser = Object.fromEntries(
      Object.entries(user).filter(
        ([, value]) => value !== "" && value !== null && value !== undefined,
      ),
    );
    const response = await api.post("/admin/users", filteredUser);
    return response.data;
  },
  update: async (email: string, user: Partial<User>) => {
    try {
      const response = await api.patch(`/admin/users/${encodeURIComponent(email)}`, user);
      return response.data;
    } catch (error: unknown) {
      console.error("Error updating user:", error);
      throw error;
    }
  },
  updateRoles: async (email: string, roles: string[]) => {
    const response = await api.patch(`/admin/users/${encodeURIComponent(email)}/roles`, { roles });
    return response.data;
  },
  deleteUser: async (email: string) => {
    const response = await api.delete(`/admin/users/${encodeURIComponent(email)}`);
    return response.data;
  },
  validateUser: async (email: string) => {
    const response = await api.post(`/admin/users/${encodeURIComponent(email)}/validate`);
    return response.data;
  },
};

export const eventsApi = {
  getAll: async (): Promise<Event[]> => {
    const response = await api.get("/admin/events");
    return response.data.sort(
      (a: Event, b: Event) => new Date(a.start_date).getTime() - new Date(b.start_date).getTime(),
    );
  },
  create: async (event: Partial<Event>) => {
    const filteredEvent = Object.fromEntries(
      Object.entries(event).filter(
        ([, value]) => value !== "" && value !== null && value !== undefined,
      ),
    );
    const response = await api.post("/admin/events", filteredEvent);
    return response.data;
  },
  update: async (id: number, event: Partial<Event>) => {
    const response = await api.patch(`/admin/events/${id}`, event);
    return response.data;
  },
  delete: async (id: number) => {
    const response = await api.delete(`/admin/events/${id}`);
    return response.data;
  },
};

export const clubsApi = {
  getAll: async (): Promise<Club[]> => {
    const response = await api.get("/admin/clubs");
    return response.data;
  },
  create: async (club: Partial<Club>) => {
    const filteredClub = Object.fromEntries(
      Object.entries(club).filter(
        ([, value]) => value !== "" && value !== null && value !== undefined,
      ),
    );
    const response = await api.post("/admin/clubs", filteredClub);
    return response.data;
  },
  update: async (id: number, club: Partial<Club>) => {
    const response = await api.patch(`/admin/clubs/${id}`, club);
    return response.data;
  },
  delete: async (id: number) => {
    const response = await api.delete(`/admin/clubs/${id}`);
    return response.data;
  },
};

export const statsApi = {
  getDashboard: async (): Promise<DashboardStats> => {
    const response = await api.get("/statistics/dashboard");
    return response.data;
  },
};

export const rolesApi = {
  getAll: async (): Promise<{ id_roles: number; name: string }[]> => {
    const response = await api.get("/admin/roles");
    return response.data;
  },
};

export const menuApi = {
  getAll: async (): Promise<MenuItem[]> => {
    const response = await api.get("/admin/menu");
    return response.data;
  },
  delete: async (id: number) => {
    const response = await api.delete(`/admin/menu/${id}`);
    return response.data;
  },
  getReviews: async (id: number): Promise<MenuItemReview[]> => {
    const response = await api.get(`/admin/menu/${id}/reviews`);
    return response.data;
  },
  deleteReview: async (id: number, email: string) => {
    const response = await api.delete(`/admin/menu/${id}/reviews/${encodeURIComponent(email)}`);
    return response.data;
  },
};

export const bassineApi = {
  getScores: async (): Promise<BassineScore[]> => {
    const response = await api.get("/admin/bassine/scores");
    return response.data;
  },
  updateScore: async (request: UpdateBassineScoreRequest) => {
    const response = await api.post("/admin/bassine/update-score", request);
    return response.data;
  },
  getHistory: async (email: string): Promise<BassineScoreHistory[]> => {
    const response = await api.get(`/admin/bassine/history/${encodeURIComponent(email)}`);
    return response.data;
  },
};
