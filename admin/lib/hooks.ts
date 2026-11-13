import { keepPreviousData, useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import type {
  CreateCategoryRequest,
  CreateItemRequest,
  ReservationItem,
  ReservationTreeItem,
  UpdateCategoryRequest,
  UpdateItemRequest,
  UpdateReservationItemMessagesRequest,
} from "./api";
import {
  bassineApi,
  clubsApi,
  eventsApi,
  menuApi,
  notificationsApi,
  photoAlbumsApi,
  photoboothEventsApi,
  photoboothReviewApi,
  photoReportsApi,
  reservationApi,
  reviewsApi,
  rolesApi,
  statsApi,
  traqApi,
  usersApi,
} from "./api";
import type {
  ActiveUsersPeriod,
  ActivityHeatmapRange,
  Club,
  Event,
  NotificationAudience,
  NotificationCategory,
  PhotoAlbumsFilter,
  PhotoboothEventSettings,
  PhotoReportsFilter,
  SendNotificationRequest,
  TraqArticle,
  User,
} from "./types";

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
    refetchInterval: 60_000,
  });
};

export const useAppUsageStats = () => {
  return useQuery({
    queryKey: ["app-usage-stats"],
    queryFn: statsApi.getAppUsage,
    refetchInterval: 60_000,
  });
};

export const useGlobalApiStats = () => {
  return useQuery({
    queryKey: ["api-stats-global"],
    queryFn: statsApi.getGlobal,
    refetchInterval: 60_000,
  });
};

export const useEndpointApiStats = () => {
  return useQuery({
    queryKey: ["api-stats-endpoints"],
    queryFn: statsApi.getEndpoints,
    refetchInterval: 60_000,
  });
};

export const useActiveUsersOverTime = (period: ActiveUsersPeriod) => {
  return useQuery({
    queryKey: ["active-users-over-time", period],
    queryFn: () => statsApi.getActiveUsers(period),
    refetchInterval: 60_000,
    placeholderData: keepPreviousData,
  });
};

export const useActivityHeatmap = (range: ActivityHeatmapRange) => {
  return useQuery({
    queryKey: ["activity-heatmap", range],
    queryFn: () => statsApi.getActivityHeatmap(range),
    refetchInterval: 60_000,
    placeholderData: keepPreviousData,
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

// Traq hooks
export const useTraqArticles = () => {
  return useQuery({
    queryKey: ["traq-articles"],
    queryFn: traqApi.getAllArticles,
  });
};

export const useCreateTraqArticle = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: traqApi.createArticle,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["traq-articles"] });
    },
  });
};

export const useUpdateTraqArticle = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: Partial<TraqArticle> }) =>
      traqApi.updateArticle(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["traq-articles"] });
    },
  });
};

export const useDeleteTraqArticle = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: traqApi.deleteArticle,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["traq-articles"] });
    },
  });
};

export const useTraqTypes = () => {
  return useQuery({
    queryKey: ["traq-types"],
    queryFn: traqApi.getAllTypes,
  });
};

export const useCreateTraqType = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: traqApi.createType,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["traq-types"] });
    },
  });
};

export const useUpdateTraqType = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: { name: string } }) =>
      traqApi.updateType(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["traq-types"] });
      queryClient.invalidateQueries({ queryKey: ["traq-articles"] });
    },
  });
};

export const useDeleteTraqType = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: traqApi.deleteType,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["traq-types"] });
    },
  });
};

// Notifications hooks
export const useNotificationRecipients = (
  audience: NotificationAudience | null,
  service?: NotificationCategory,
) => {
  return useQuery({
    queryKey: ["notification-recipients", audience, service],
    queryFn: () =>
      notificationsApi.send({
        title: "preview",
        message: "",
        audience: audience as NotificationAudience,
        service,
        dryRun: true,
      }),
    enabled: audience !== null,
  });
};

export const useSendNotification = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (request: SendNotificationRequest) => notificationsApi.send(request),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["push-messages"] }),
  });
};

/** Versions of the app in use, newest first. */
export const useAppVersions = (enabled: boolean) =>
  useQuery({
    queryKey: ["app-versions"],
    queryFn: notificationsApi.appVersions,
    enabled,
    staleTime: 60_000,
  });

/** Latest pushes with their delivery counts (no text: the server keeps none of it in the summaries). */
export const useRecentPushMessages = () =>
  useQuery({
    queryKey: ["push-messages"],
    queryFn: () => notificationsApi.recent(8),
    refetchInterval: 15_000,
    retry: false,
  });

// Photo reports hooks
export const usePhotoReports = (filter: PhotoReportsFilter) => {
  return useQuery({
    queryKey: ["photo-reports", filter],
    queryFn: () => photoReportsApi.list(filter),
    placeholderData: keepPreviousData,
  });
};

export const usePhotoReport = (id: number | null) => {
  return useQuery({
    queryKey: ["photo-report", id],
    queryFn: () => photoReportsApi.get(id as number),
    enabled: id !== null,
  });
};

export const useRetryPhotoReport = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: photoReportsApi.retry,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["photo-reports"] });
      queryClient.invalidateQueries({ queryKey: ["photo-report"] });
    },
  });
};

export const useDecidePhotoReport = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      id,
      decision,
      note,
    }: {
      id: number;
      decision: "approve" | "reject";
      note: string;
    }) => photoReportsApi.decide(id, decision, note),
    onSuccess: (report) => {
      queryClient.setQueryData(["photo-report", report.id], report);
      queryClient.invalidateQueries({ queryKey: ["photo-reports"] });
      queryClient.invalidateQueries({ queryKey: ["photo-report"] });
    },
  });
};

export const useSetAlbumClub = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ albumId, clubId }: { albumId: number; clubId: number | null }) =>
      photoReportsApi.setAlbumClub(albumId, clubId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["photo-reports"] });
      queryClient.invalidateQueries({ queryKey: ["photo-report"] });
    },
  });
};

// Photobooth review hooks
export const usePhotoboothPending = () =>
  useQuery({
    queryKey: ["photobooth-pending"],
    queryFn: photoboothReviewApi.pending,
    refetchInterval: 60_000,
  });

export const usePhotoboothPictures = (folderId: string | null) =>
  useQuery({
    queryKey: ["photobooth-pictures", folderId],
    queryFn: () => photoboothReviewApi.pictures(folderId as string),
    enabled: folderId !== null,
    refetchInterval: 60_000,
  });

export const useApprovePhotobooth = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: photoboothReviewApi.approve,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["photobooth-pending"] });
      queryClient.invalidateQueries({ queryKey: ["photobooth-pictures"] });
    },
  });
};

export const useDeletePhotoboothPicture = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: photoboothReviewApi.deletePicture,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["photobooth-pending"] });
      queryClient.invalidateQueries({ queryKey: ["photobooth-pictures"] });
    },
  });
};

// Photobooth events hooks
export const usePhotoboothCatalog = () =>
  useQuery({ queryKey: ["photobooth-catalog"], queryFn: photoboothEventsApi.catalog });

export const usePhotoboothEvents = () =>
  useQuery({ queryKey: ["photobooth-events"], queryFn: photoboothEventsApi.list });

export const useSavePhotoboothEvent = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, settings }: { id: number | null; settings: PhotoboothEventSettings }) =>
      id === null ? photoboothEventsApi.create(settings) : photoboothEventsApi.update(id, settings),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["photobooth-events"] }),
  });
};

export const useDeletePhotoboothEvent = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: photoboothEventsApi.remove,
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["photobooth-events"] }),
  });
};

export const usePhotoboothFrame = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, type, file }: { id: number; type: string; file: File | null }) =>
      file
        ? photoboothEventsApi.uploadFrame(id, type, file).then(() => undefined)
        : photoboothEventsApi.removeFrame(id, type),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["photobooth-events"] }),
  });
};

// Photo albums hooks
export const usePhotoAlbums = (filter: PhotoAlbumsFilter) =>
  useQuery({
    queryKey: ["photo-albums", filter],
    queryFn: () => photoAlbumsApi.list(filter),
    placeholderData: keepPreviousData,
  });

export const useSetAlbumHidden = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, hidden }: { id: number; hidden: boolean }) =>
      photoAlbumsApi.setHidden(id, hidden),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["photo-albums"] });
    },
  });
};
