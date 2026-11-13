export interface User {
  id_newf: number;
  email: string;
  first_name: string;
  last_name: string;
  phone_number?: string;
  profile_picture?: string;
  graduation_year?: number;
  formation_name?: string;
  campus?: string;
  language: string;
  password_updated_date?: string;
  creation_date?: string;
  verification_code?: string;
  verification_code_expiration?: string;
  roles?: string[];
}

export interface Event {
  id_events: number;
  name: string;
  description: string;
  link: string;
  start_date: string;
  end_date: string;
  location: string;
  creation_date: string;
  picture: string;
  creator: string;
  id_club: number;
  attendee_count?: number;
}

export interface Club {
  id_clubs: number;
  name: string;
  picture: string;
  description: string;
  location: string;
  link: string;
  member_count?: number;
}

export interface ClubWithResponsible extends Club {
  responsible?: {
    first_name: string;
    last_name: string;
  };
}

export interface DashboardStats {
  totalUsers: number;
  unverifiedUsers: number;
  /** Accounts no push can reach: notifications refused, or the app never opened on a phone. */
  usersWithoutNotifications?: number;
  totalEvents: number;
  totalClubs: number;
  userGrowth: { date: string; count: number; cumulativeCount: number }[];
  activeUsers?: { dau: number; wau: number; mau: number; yau: number };
  dailyActiveUsers?: { date: string; count: number }[];
}

export type ActiveUsersPeriod = "day" | "week" | "month" | "quarter" | "year";

export interface ActiveUsersPoint {
  date: string;
  count: number;
}

export type ActivityHeatmapRange = "week" | "month" | "year" | "all";

export interface ActivityHourPoint {
  dayOfWeek: number; // ISO day of week: 1 = lundi ... 7 = dimanche
  hour: number; // 0-23
  count: number;
}

export interface GlobalApiStats {
  total_request_count: number;
  global_avg_duration_ms: number;
  global_min_duration_ms: number;
  global_max_duration_ms: number;
  global_success_rate_percent: number;
  first_request: string;
  last_request: string;
  success_count: number;
  error_count: number;
}

export interface EndpointApiStat {
  endpoint: string;
  method: string;
  request_count: number;
  avg_duration_ms: number;
  min_duration_ms: number;
  max_duration_ms: number;
  success_rate_percent: number;
  success_count: number;
  error_count: number;
  first_request: string;
  last_request: string;
}

// Types pour la gestion des erreurs
export interface ApiError {
  message?: string;
  response?: {
    data?: {
      error?: string;
    };
  };
}

export interface EventWithClubName extends Event {
  club_name?: string;
}

// Menu types
export interface MenuItem {
  id_restaurant_articles: number;
  name: string;
  first_time_served: string;
  last_time_served?: string;
  last_served?: string;
  average_rating: number;
  total_ratings: number;
  times_served: number;
}

export interface MenuItemReview {
  email: string;
  note: number;
  comment: string;
  date: string;
  first_name: string;
  last_name: string;
  profile_picture: string;
  dish_id?: number;
  dish_name?: string;
}

// Bassine/Games types
export interface BassineScore {
  id: number;
  user_email: string;
  user_first_name: string;
  user_last_name: string;
  current_score: number;
  total_games_played: number;
  creation_date: string;
  last_updated: string;
}

export interface BassineScoreHistory {
  id: number;
  user_email: string;
  score_change: number;
  new_total: number;
  game_date: string;
  notes?: string;
  admin_email?: string;
}

export interface UpdateBassineScoreRequest {
  userEmail: string;
  scoreChange: number;
  notes?: string;
}

// Traq types
export interface TraqArticle {
  id_traq: number;
  name: string;
  disabled: boolean;
  limited: boolean;
  alcohol: number;
  out_of_stock: boolean;
  creation_date: string;
  picture: string;
  description: string;
  price: number;
  price_half: number;
  traq_type: string;
}

export interface TraqType {
  id_traq_types: number;
  name: string;
}

// Notifications
export type NotificationAudience =
  | { type: "all" }
  | { type: "club"; clubId: number }
  | { type: "campus"; campus: string }
  | { type: "users"; emails: string[] }
  | { type: "cohort"; formation?: string; graduationYear?: number }
  /** Devices running an app older than belowVersion (only those devices are notified). */
  | { type: "app_version"; belowVersion: string; platform?: "ios" | "android" };

export type NotificationCategory =
  | "RESTAURANT"
  | "TRAQ"
  | "EVENTS"
  | "EVENT_REMINDERS"
  | "RESERVATIONS";

// Screen opened when the notification is tapped; types and ids match the backend NavigationTarget.
export interface NotificationNavigation {
  type: "event" | "club" | "restaurant" | "service" | "url";
  id?: string;
}

export interface SendNotificationRequest {
  title: string;
  message: string;
  audience: NotificationAudience;
  service?: NotificationCategory;
  navigation?: NotificationNavigation;
  /** Public https image shown in the notification (Android large picture, iOS attachment). */
  imageUrl?: string;
  dryRun?: boolean;
}

export interface NotificationRecipients {
  users: number;
  devices: number;
  withoutDevice: string[];
  /** The same people as withoutDevice, by name. */
  withoutDeviceNames?: string[];
  sent: boolean;
  // Set when the message was queued for the delivery worker; `id` follows it.
  queued?: boolean;
  id?: number;
}

// Delivery diagnostics of one queued message: counts and codes only.
export interface PushMessageSummary {
  id: number;
  kind?: string;
  category?: string;
  status: string;
  createdBy?: string;
  createdAt?: string;
  completedAt?: string;
  users?: number;
  eligibleDevices?: number;
  deliveries: Record<string, number>;
  errors?: Record<string, number>;
}

export type PhotoReportReason = "appears_in_photo" | "inappropriate" | "poor_quality" | "other";

export type PhotoReportDeliveryStatus = "sent" | "pending" | "failed" | "no_recipients";

export interface PhotoReportReporter {
  id: number;
  public_id: string | null;
  name: string;
  // Null once the reporter deleted their account.
  email: string | null;
  graduation_year: number | null;
  formation: string;
}

export interface PhotoReportDelivery {
  email: string;
  role: "club_owner" | "photo_club" | "admin";
  status: "pending" | "sending" | "sent" | "failed";
  attempts: number;
  sent_at: string | null;
  last_error: string | null;
}

export interface PhotoReport {
  id: number;
  created_at: string;
  reason: PhotoReportReason;
  comment: string;
  photo_id: string;
  photo_name: string;
  // Opens the file in Google Drive, where it can be deleted.
  photo_url: string;
  // False once the photo is gone from the Drive index.
  photo_available: boolean;
  album_id: number | null;
  album_title: string;
  album_url: string;
  club_id: number | null;
  club_name: string;
  reporter: PhotoReportReporter;
  delivery_status: PhotoReportDeliveryStatus;
  recipients: number;
  recipients_sent: number;
  recipients_failed: number;
  // Only on the detail response.
  deliveries?: PhotoReportDelivery[];
  status: PhotoReportStatus;
  decided_at: string | null;
  decided_by: string | null;
  decision_note: string;
}

export interface PhotoReportsPage {
  reports: PhotoReport[];
  total: number;
  counts: Record<PhotoReportStatus, number>;
  /** False when the server cannot delete photos from Drive: reports can then only be rejected. */
  can_remove: boolean;
}

export type PhotoReportStatus = "pending" | "approved" | "rejected";

export interface PhotoReportsFilter {
  status?: PhotoReportStatus | "";
  reason?: PhotoReportReason | "";
  delivery?: PhotoReportDeliveryStatus | "";
  limit: number;
  offset: number;
}

export interface AppVersionStat {
  version: string;
  platform: "ios" | "android" | string;
  devices: number;
  users: number;
  lastSeen: string | null;
}

export interface FaceRecognitionStats {
  available: boolean;
  users_enabled: number;
  users_with_selfie: number;
  users_with_photos: number;
  average_photos_by_user: number;
  rejected_matches: number;
  photos_total: number;
  photos_analysed: number;
  photos_failed: number;
  photos_pending: number;
  faces_found: number;
  last_scan_at: string | null;
  scoped: boolean;
}

export interface AppUsageStats {
  devices: {
    active_window_days: number;
    devices: number;
    users: number;
    platforms: { platform: "ios" | "android"; devices: number; users: number }[];
    versions: { platform: "ios" | "android"; version: string; devices: number; users: number }[];
  };
  faces: FaceRecognitionStats;
}

/** An event's photobooth pictures awaiting approval (Drive folder "photobooth_temp"). */
export interface PhotoboothPendingFolder {
  folder_id: string;
  event_name: string;
  event_date: string; // yyyy-mm-dd
  pictures: number;
  shoots: number;
  /** Shoots started within the last hour and not finished: guests may still be at the booth. */
  open_shoots: number;
  first_upload_at: string | null;
  last_upload_at: string | null;
  drive_url: string;
}

export interface PhotoboothPendingPicture {
  id: number;
  shoot_id: number;
  kind: "photo" | "collage";
  number: number;
  name: string;
  width: number | null;
  height: number | null;
  uploaded_at: string;
}

export interface PhotoboothApproveResult {
  folder_id: string;
  merged: boolean;
  moved: number;
}

// Photo albums (auto-created from the Photos Drive)
export type PhotoAlbumStatus = "shown" | "hidden" | "review";

export interface PhotoAlbum {
  id: number;
  parent_id: number | null;
  title: string;
  event_date: string | null;
  academic_year: string;
  category: string;
  photo_count: number;
  /** The classifier's decision: only "published" albums can be shown. */
  visibility: "published" | "needs_review" | "hidden";
  /** Hidden by an administrator or the Club Photo, whatever the classifier says. */
  hidden: boolean;
  hidden_at: string | null;
  hidden_by: string | null;
  first_taken_at: string | null;
  /** Visible in the app and on transat.dev. */
  shown: boolean;
  /** The album's transat.dev page (answers only while the album is shown). */
  share_url: string;
  /** Cover image for the panel, relative to the API (needs the panel's token). */
  cover_path: string | null;
  drive_url: string;
}

export interface PhotoAlbumsPage {
  albums: PhotoAlbum[];
  counts: { all: number; shown: number; hidden: number; review: number };
  page: number;
  page_size: number;
  total: number;
}

export interface PhotoAlbumsFilter {
  q: string;
  status: PhotoAlbumStatus | "";
  page: number;
}

/** One batch of approving every pending photo report. */
export interface ApproveAllResult {
  /** Photos removed from Drive and the app in this batch. */
  photos: number;
  /** Reports closed in this batch (a photo's other reports close with it). */
  approved: number;
  /** Reports Drive refused: still pending, to skip in the next batches. */
  failed: number[];
  /** Pending reports still to do; call again while above zero. */
  remaining: number;
}
