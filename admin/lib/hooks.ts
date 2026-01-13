import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  bassineApi,
  clubsApi,
  eventsApi,
  menuApi,
  reservationApi,
  reviewsApi,
  rolesApi,
  statsApi,
  usersApi,
} from "./api";
import type {
  CreateCategoryRequest,
  CreateItemRequest,
  ReservationItem,
  ReservationTreeItem,
  UpdateCategoryRequest,
  UpdateItemRequest,
  UpdateReservationItemMessagesRequest,
} from "./api";
import type { Club, Event, User } from "./types";

export * from "./hooks/useClickOutside";
// Export utility hooks
export * from "./hooks/useDebounce";
export * from "./hooks/useKeyboardShortcut";
export * from "./hooks/useLocalStorage";
export * from "./hooks/useToggle";

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
    mutationFn: ({ id, data }: { id: number; data: Partial<Event> }) => eventsApi.update(id, data),
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
    mutationFn: ({ id, data }: { id: number; data: Partial<Club> }) => clubsApi.update(id, data),
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

export const useClubOwners = (clubId: number, options?: { enabled?: boolean }) => {
  return useQuery({
    queryKey: ["club-owners", clubId],
    queryFn: () => clubsApi.getOwners(clubId),
    enabled: options?.enabled ?? true,
  });
};

export const useAddClubOwner = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ clubId, email }: { clubId: number; email: string }) =>
      clubsApi.addOwner(clubId, email),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: ["club-owners", variables.clubId] });
      queryClient.invalidateQueries({ queryKey: ["clubs"] });
    },
  });
};

export const useRemoveClubOwner = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ clubId, email }: { clubId: number; email: string }) =>
      clubsApi.removeOwner(clubId, email),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: ["club-owners", variables.clubId] });
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
      queryClient.invalidateQueries({ queryKey: ["all-reviews"] });
    },
  });
};

// Reviews hooks
export const useAllReviews = (userEmail?: string) => {
  return useQuery({
    queryKey: ["all-reviews", userEmail],
    queryFn: () => reviewsApi.getAll(userEmail),
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

// Reservation items hooks
export const useReservationItemsForClub = (clubId: number, options?: { enabled?: boolean }) => {
  return useQuery({
    queryKey: ["reservation-items", clubId],
    queryFn: () => reservationApi.getItemsForClub(clubId),
    enabled: options?.enabled ?? true,
  });
};

export const useUpdateReservationItemMessages = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      itemId,
      messages,
    }: {
      itemId: number;
      messages: UpdateReservationItemMessagesRequest;
    }) => reservationApi.updateItemMessages(itemId, messages),
    onSuccess: (_data, variables) => {
      // Invalidate the reservation items query for the club
      // We need to find which club this item belongs to, but for simplicity, invalidate all
      queryClient.invalidateQueries({ queryKey: ["reservation-items"] });
    },
  });
};

// Reservation tree hooks
export const useReservationTree = () => {
  return useQuery({
    queryKey: ["reservation-tree"],
    queryFn: reservationApi.getTree,
  });
};

export const useCreateReservationCategory = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (category: CreateCategoryRequest) => reservationApi.createCategory(category),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["reservation-tree"] });
    },
  });
};

export const useUpdateReservationCategory = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: UpdateCategoryRequest }) =>
      reservationApi.updateCategory(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["reservation-tree"] });
    },
  });
};

export const useDeleteReservationCategory = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: reservationApi.deleteCategory,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["reservation-tree"] });
    },
  });
};

export const useCreateReservationItem = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (item: CreateItemRequest) => reservationApi.createItem(item),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["reservation-tree"] });
    },
  });
};

export const useUpdateReservationItem = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: UpdateItemRequest }) =>
      reservationApi.updateItem(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["reservation-tree"] });
    },
  });
};

export const useDeleteReservationItem = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: reservationApi.deleteItem,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["reservation-tree"] });
    },
  });
};
