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
  | { type: "cohort"; formation?: string; graduationYear?: number };

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
  dryRun?: boolean;
}

export interface NotificationRecipients {
  users: number;
  devices: number;
  withoutDevice: string[];
  sent: boolean;
}
