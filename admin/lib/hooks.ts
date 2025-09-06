import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { usersApi, eventsApi, clubsApi, statsApi, rolesApi, menuApi, bassineApi } from "./api";
import { User, Event, Club } from "./types";

// Export utility hooks
export * from "./hooks/useDebounce";
export * from "./hooks/useLocalStorage";
export * from "./hooks/useKeyboardShortcut";
export * from "./hooks/useToggle";
export * from "./hooks/useClickOutside";

// Users hooks
export const useUsers = () => {
  return useQuery({
    queryKey: ["users"],
    queryFn: usersApi.getAll,
  });
};

export const useCreateUser = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: usersApi.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["users"] });
    },
  });
};

export const useUpdateUser = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: ({ email, data }: { email: string; data: Partial<User> }) => 
      usersApi.update(email, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["users"] });
    },
  });
};

export const useDeleteUser = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: usersApi.deleteUser,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["users"] });
    },
  });
};

export const useValidateUser = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: usersApi.validateUser,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["users"] });
    },
  });
};

// Events hooks
export const useEvents = () => {
  return useQuery({
    queryKey: ["events"],
    queryFn: eventsApi.getAll,
  });
};

export const useCreateEvent = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: eventsApi.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["events"] });
    },
  });
};

export const useUpdateEvent = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: Partial<Event> }) => 
      eventsApi.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["events"] });
    },
  });
};

export const useDeleteEvent = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: eventsApi.delete,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["events"] });
    },
  });
};

// Clubs hooks
export const useClubs = () => {
  return useQuery({
    queryKey: ["clubs"],
    queryFn: clubsApi.getAll,
  });
};

export const useCreateClub = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: clubsApi.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["clubs"] });
    },
  });
};

export const useUpdateClub = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: Partial<Club> }) => 
      clubsApi.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["clubs"] });
    },
  });
};

export const useDeleteClub = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: clubsApi.delete,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["clubs"] });
    },
  });
};

// Dashboard stats hook
export const useDashboardStats = () => {
  return useQuery({
    queryKey: ["dashboard-stats"],
    queryFn: statsApi.getDashboard,
  });
};

// Roles hook
export const useRoles = () => {
  return useQuery({
    queryKey: ["roles"],
    queryFn: rolesApi.getAll,
  });
};

// Menu hooks
export const useMenuItems = () => {
  return useQuery({
    queryKey: ["menu-items"],
    queryFn: menuApi.getAll,
  });
};

export const useDeleteMenuItem = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: menuApi.delete,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["menu-items"] });
    },
  });
};

export const useMenuItemReviews = (menuItemId: number) => {
  return useQuery({
    queryKey: ["menu-item-reviews", menuItemId],
    queryFn: () => menuApi.getReviews(menuItemId),
    enabled: !!menuItemId,
  });
};

export const useDeleteMenuItemReview = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: ({ menuItemId, email }: { menuItemId: number; email: string }) => 
      menuApi.deleteReview(menuItemId, email),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ["menu-item-reviews", variables.menuItemId] });
      queryClient.invalidateQueries({ queryKey: ["menu-items"] });
    },
  });
};

// Bassine/Games hooks
export const useBassineScores = () => {
  return useQuery({
    queryKey: ["bassine-scores"],
    queryFn: bassineApi.getScores,
  });
};

export const useBassineHistory = (userEmail: string) => {
  return useQuery({
    queryKey: ["bassine-history", userEmail],
    queryFn: () => bassineApi.getHistory(userEmail),
    enabled: !!userEmail,
  });
};

export const useUpdateBassineScore = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: bassineApi.updateScore,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["bassine-scores"] });
    },
  });
};