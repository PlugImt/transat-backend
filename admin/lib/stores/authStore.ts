import { create } from "zustand";
import { persist } from "zustand/middleware";

export interface AuthUser {
  email: string;
  roles: string[];
}

interface AuthState {
  user: AuthUser | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;

  // Actions
  setAuth: (user: AuthUser, token: string) => void;
  logout: () => void;
  setLoading: (loading: boolean) => void;
  updateUser: (updates: Partial<AuthUser>) => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      token: null,
      isAuthenticated: false,
      isLoading: true,

      setAuth: (user, token) =>
        set({
          user,
          token,
          isAuthenticated: true,
          isLoading: false,
        }),

      logout: () => {
        set({
          user: null,
          token: null,
          isAuthenticated: false,
          isLoading: false,
        });
        // Clear localStorage
        localStorage.removeItem("adminToken");
      },

      setLoading: (loading) => set({ isLoading: loading }),

      updateUser: (updates) => {
        const currentUser = get().user;
        if (currentUser) {
          set({ user: { ...currentUser, ...updates } });
        }
      },
    }),
    {
      name: "admin-auth-storage",
      partialize: (state) => ({
        user: state.user,
        token: state.token,
        isAuthenticated: state.isAuthenticated,
      }),
    },
  ),
);
