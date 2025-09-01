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