// Shapes of the reservation administration endpoints (backend: handlers/admin/reservations_overview.go).

export type ItemStatus = "active" | "maintenance" | "inactive";

export type ItemAvailability = "available" | "borrowed" | "occupied" | "maintenance" | "inactive";

export interface AdminHolder {
  email: string;
  first_name: string;
  last_name: string;
  start_date: string;
  end_date?: string;
  /** Only meaningful for the loan of an item without slot. */
  elapsed_seconds: number;
}

export interface AdminItem {
  id: number;
  name: string;
  slot: boolean;
  description?: string;
  location?: string;
  status: ItemStatus;
  status_note?: string;
  club_id: number;
  club_name: string;
  category_id?: number;
  /** "Parent / Child", empty at the root of a club. */
  category_path: string;
  warning_message?: string;
  confirmation_message?: string;
  availability: ItemAvailability;
  current?: AdminHolder;
  next?: AdminHolder;
  upcoming_count: number;
  waitlist_count: number;
}

export type ReservationStatus = "upcoming" | "ongoing" | "active" | "past" | "cancelled";

/** Filter of the list: "active" is a loan in progress or a slot being used, "overdue" a loan older than a day. */
export type ReservationStatusFilter = "" | "active" | "upcoming" | "past" | "cancelled" | "overdue";

export type ReminderMilestone = "slot_1h" | "loan_d1" | "loan_d5" | "loan_d7";

export type CancelKind = "cancelled" | "cancelled_by_admin" | "returned_by_admin";

export interface AdminReservation {
  key: string;
  item_id: number;
  item_name: string;
  slot: boolean;
  club_name: string;
  email: string;
  first_name: string;
  last_name: string;
  start_date: string;
  end_date?: string;
  status: ReservationStatus;
  elapsed_seconds?: number;
  /** 1, 5 or 7 for a loan in progress that is that many days old. */
  overdue_level: 0 | 1 | 5 | 7;
  reminders: ReminderMilestone[];
  cancelled_at?: string;
  cancelled_by?: string;
  cancel_kind?: CancelKind;
  cancellable: boolean;
}

export interface ReservationFilters {
  status?: ReservationStatusFilter;
  kind?: "slot" | "loan" | "";
  item_id?: number;
  q?: string;
  /** YYYY-MM-DD or an ISO instant. */
  from?: string;
  to?: string;
  overdue_days?: number;
  order?: "asc" | "desc";
  limit?: number;
  offset?: number;
}

export interface ReservationPage {
  items: AdminReservation[];
  total: number;
  limit: number;
  offset: number;
}

export interface ReservationStats {
  active_loans: number;
  ongoing_slots: number;
  upcoming_slots: number;
  overdue_1d: number;
  overdue_5d: number;
  overdue_7d: number;
  waitlist_entries: number;
  waitlist_awaiting_response: number;
  cancelled_30d: number;
  items: number;
  unavailable_items: number;
}

export interface AdminWaitlistEntry {
  id: number;
  position: number;
  email: string;
  first_name: string;
  last_name: string;
  slot_start?: string;
  joined_at: string;
  notified_at?: string;
  /** When the person who was told loses their turn; the scheduler then tells the next one. */
  turn_expires_at?: string;
}

export interface AdminWaitlist {
  entries: AdminWaitlistEntry[];
  free: boolean;
  awaiting_response: boolean;
}

export interface CancelReservationRequest {
  item_id: number;
  /** UTC wall time "YYYY-MM-DD HH:mm:ss" or an ISO instant; required for an item with slots. */
  start_date?: string;
}

export interface CancelReservationResult {
  message: string;
  owner: string;
  cancelled_by_admin: boolean;
  waitlist_notified: boolean;
}
