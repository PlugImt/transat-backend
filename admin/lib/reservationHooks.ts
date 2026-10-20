import { keepPreviousData, useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { reservationApi } from "./api";
import type { ReservationFilters } from "./reservationTypes";

// Prefix shared with the category tree: the existing mutations that invalidate it refresh these views too.
const KEY = "reservation-tree";

/** Everything the reservations screens show is derived from these: one refresh after any action. */
const useRefreshReservations = () => {
  const queryClient = useQueryClient();
  return () => queryClient.invalidateQueries({ queryKey: [KEY] });
};

export const useReservationCatalogue = () =>
  useQuery({
    queryKey: [KEY, "catalogue"],
    queryFn: reservationApi.getCatalogue,
    // The statuses are live: a loan can start while the panel is open.
    refetchInterval: 60_000,
  });

export const useReservationStats = () =>
  useQuery({
    queryKey: [KEY, "stats"],
    queryFn: reservationApi.getStats,
    refetchInterval: 60_000,
  });

export const useAdminReservations = (
  filters: ReservationFilters,
  options?: { enabled?: boolean },
) =>
  useQuery({
    queryKey: [KEY, "list", filters],
    queryFn: () => reservationApi.list(filters),
    placeholderData: keepPreviousData,
    enabled: options?.enabled ?? true,
  });

export const useAdminWaitlist = (
  itemId: number | undefined,
  startDate?: string,
  options?: { enabled?: boolean },
) =>
  useQuery({
    queryKey: [KEY, "waitlist", itemId, startDate ?? null],
    queryFn: () => reservationApi.getWaitlist(itemId as number, startDate),
    enabled: itemId !== undefined && (options?.enabled ?? true),
  });

export const useAdminCancelReservation = () => {
  const refresh = useRefreshReservations();
  return useMutation({ mutationFn: reservationApi.cancel, onSettled: refresh });
};

export const useMarkItemReturned = () => {
  const refresh = useRefreshReservations();
  return useMutation({ mutationFn: reservationApi.markReturned, onSettled: refresh });
};

export const useRemoveWaitlistEntry = () => {
  const refresh = useRefreshReservations();
  return useMutation({ mutationFn: reservationApi.removeWaitlistEntry, onSettled: refresh });
};

export const useAdvanceWaitlist = () => {
  const refresh = useRefreshReservations();
  return useMutation({
    mutationFn: ({ itemId, startDate }: { itemId: number; startDate?: string }) =>
      reservationApi.advanceWaitlist(itemId, startDate),
    onSettled: refresh,
  });
};
