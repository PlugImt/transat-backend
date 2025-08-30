import axios from 'axios';
import { User, Event, Club, DashboardStats } from './types';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:3000';

const api = axios.create({
  baseURL: API_BASE_URL,
});

// Add auth token to requests
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('adminToken');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export const authApi = {
  login: async (email: string, password: string) => {
    const response = await api.post('/auth/login', { email, password });
    return response.data;
  },
  verify: async (token: string) => {
    const response = await api.get('/newf/me', {
      headers: { Authorization: `Bearer ${token}` }
    });
    return response.data;
  },
};

export const usersApi = {
  getAll: async (): Promise<User[]> => {
    const response = await api.get('/admin/users');
    return response.data;
  },
  create: async (user: Partial<User>) => {
    const filteredUser = Object.fromEntries(
      Object.entries(user).filter(([, value]) => 
        value !== '' && value !== null && value !== undefined
      )
    );
    const response = await api.post('/admin/users', filteredUser);
    return response.data;
  },
  update: async (email: string, user: Partial<User>) => {
    try {
      const response = await api.patch(`/admin/users/${email}`, user);
      return response.data;
    } catch (error: unknown) {
      console.error('Error updating user:', error);
      throw error;
    }
  },
  updateRoles: async (email: string, roles: string[]) => {
    const response = await api.patch(`/admin/users/${email}/roles`, { roles });
    return response.data;
  },
  deleteUser: async (email: string) => {
    const response = await api.delete(`/admin/users/${email}`);
    return response.data;
  },
  validateUser: async (email: string) => {
    const response = await api.post(`/admin/users/${email}/validate`);
    return response.data;
  },
};

export const eventsApi = {
  getAll: async (): Promise<Event[]> => {
    const response = await api.get('/admin/events');
    return response.data.sort((a: Event, b: Event) => new Date(a.start_date).getTime() - new Date(b.start_date).getTime());
  },
  create: async (event: Partial<Event>) => {
    const filteredEvent = Object.fromEntries(
      Object.entries(event).filter(([, value]) => 
        value !== '' && value !== null && value !== undefined
      )
    );
    const response = await api.post('/admin/events', filteredEvent);
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
    const response = await api.get('/admin/clubs');
    return response.data;
  },
  create: async (club: Partial<Club>) => {
    const filteredClub = Object.fromEntries(
      Object.entries(club).filter(([, value]) => 
        value !== '' && value !== null && value !== undefined
      )
    );
    const response = await api.post('/admin/clubs', filteredClub);
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
    const response = await api.get('/statistics/dashboard');
    return response.data;
  },
};

export const rolesApi = {
  getAll: async (): Promise<{id_roles: number, name: string}[]> => {
    const response = await api.get('/admin/roles');
    return response.data;
  },
};