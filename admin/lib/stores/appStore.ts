import { create } from 'zustand';
import { User, Club, Event } from '../types';

interface AppState {
  // UI State
  sidebarOpen: boolean;
  commandPaletteOpen: boolean;
  currentPage: string;
  
  // Modals
  userModalOpen: boolean;
  editingUser: User | null;
  clubModalOpen: boolean;
  editingClub: Club | null;
  eventModalOpen: boolean;
  editingEvent: Event | null;
  
  // Search & Filters
  globalSearch: string;
  activeFilters: Record<string, unknown>;
  
  // Actions
  setSidebarOpen: (open: boolean) => void;
  toggleSidebar: () => void;
  setCommandPaletteOpen: (open: boolean) => void;
  toggleCommandPalette: () => void;
  setCurrentPage: (page: string) => void;
  setGlobalSearch: (search: string) => void;
  setFilter: (key: string, value: unknown) => void;
  clearFilters: () => void;
  
  // Modal actions
  openUserModal: (user?: User) => void;
  closeUserModal: () => void;
  openClubModal: (club?: Club) => void;
  closeClubModal: () => void;
  openEventModal: (event?: Event) => void;
  closeEventModal: () => void;
}

export const useAppStore = create<AppState>((set) => ({
  // Initial state
  sidebarOpen: false,
  commandPaletteOpen: false,
  currentPage: '',
  
  // Modals
  userModalOpen: false,
  editingUser: null,
  clubModalOpen: false,
  editingClub: null,
  eventModalOpen: false,
  editingEvent: null,
  
  globalSearch: '',
  activeFilters: {},
  
  // Actions
  setSidebarOpen: (open) => set({ sidebarOpen: open }),
  
  toggleSidebar: () => set((state) => ({ 
    sidebarOpen: !state.sidebarOpen 
  })),
  
  setCommandPaletteOpen: (open) => set({ 
    commandPaletteOpen: open 
  }),
  
  toggleCommandPalette: () => set((state) => ({ 
    commandPaletteOpen: !state.commandPaletteOpen 
  })),
  
  setCurrentPage: (page) => set({ currentPage: page }),
  
  setGlobalSearch: (search) => set({ globalSearch: search }),
  
  setFilter: (key, value) => set((state) => ({
    activeFilters: { ...state.activeFilters, [key]: value }
  })),
  
  clearFilters: () => set({ activeFilters: {} }),
  
  // Modal actions
  openUserModal: (user) => set({ 
    userModalOpen: true, 
    editingUser: user || null 
  }),
  
  closeUserModal: () => {
    set({ userModalOpen: false });
    setTimeout(() => {
      set({ editingUser: null });
    }, 150);
  },

  // Club modal actions
  openClubModal: (club) => set({ 
    clubModalOpen: true, 
    editingClub: club || null 
  }),

  closeClubModal: () => {
    set({ clubModalOpen: false });
    setTimeout(() => {
      set({ editingClub: null });
    }, 150);
  },

  // Event modal actions
  openEventModal: (event) => set({ 
    eventModalOpen: true, 
    editingEvent: event || null 
  }),

  closeEventModal: () => {
    set({ eventModalOpen: false });
    setTimeout(() => {
      set({ editingEvent: null });
    }, 150);
  },
}));