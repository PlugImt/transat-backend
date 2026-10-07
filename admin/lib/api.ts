import type {
  ActiveUsersPeriod,
  ActiveUsersPoint,
  ActivityHeatmapRange,
  ActivityHourPoint,
  BassineScore,
  BassineScoreHistory,
  Club,
  DashboardStats,
  EndpointApiStat,
  Event,
  GlobalApiStats,
  MenuItem,
  MenuItemReview,
  NotificationRecipients,
  SendNotificationRequest,
  TraqArticle,
  TraqType,
  UpdateBassineScoreRequest,
  User,
} from "./types";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:3001";

// API response shapes are validated by their existing consumers, not this transport layer.
// biome-ignore lint/suspicious/noExplicitAny: Preserve the current unvalidated JSON boundary.
type UncheckedJson = any;
type ApiResponse<T = UncheckedJson> = { data: T };
type RequestConfig = { params?: Record<string, string | number | undefined> };

const request = async <T = UncheckedJson>(
  method: string,
  path: string,
  body?: unknown,
  config: RequestConfig = {},
  headers: Record<string, string> = {},
): Promise<ApiResponse<T>> => {
  const url = new URL(path.replace(/^\/+/, ""), `${API_BASE_URL.replace(/\/+$/, "")}/`);
  for (const [key, value] of Object.entries(config.params ?? {})) {
    if (value !== undefined) url.searchParams.set(key, String(value));
  }

  const token = typeof window !== "undefined" ? localStorage.getItem("adminToken") : null;
  const requestHeaders: Record<string, string> = {
    Accept: "application/json",
    ...(body !== undefined ? { "Content-Type": "application/json" } : {}),
    ...(token ? { Authorization: `Bearer ${token}` } : {}),
    ...headers,
  };

  let response: Response;
  try {
    response = await fetch(url, {
      method,
      headers: requestHeaders,
      body: body === undefined ? undefined : JSON.stringify(body),
    });
  } catch (cause) {
    throw new Error("Network request failed", { cause });
  }

  const text = await response.text();
  let data: unknown;
  try {
    data = text ? JSON.parse(text) : undefined;
  } catch {
    data = text;
  }

  if (!response.ok) {
    const serverMessage =
      data && typeof data === "object" && "error" in data
        ? String(data.error)
        : `Request failed with status ${response.status}`;
    const error = new Error(serverMessage) as Error & {
      response: { status: number; data: unknown };
    };
    error.response = { status: response.status, data };
    throw error;
  }

  return { data: data as T };
};

const api = {
  get: <T = UncheckedJson>(path: string, config?: RequestConfig) =>
    request<T>("GET", path, undefined, config),
  post: <T = UncheckedJson>(path: string, body?: unknown) => request<T>("POST", path, body),
  patch: <T = UncheckedJson>(path: string, body?: unknown) => request<T>("PATCH", path, body),
  delete: <T = UncheckedJson>(path: string, config?: { data?: unknown }) =>
    request<T>("DELETE", path, config?.data),
};

function asArray<T>(data: T[] | null | undefined): T[] {
  return data ?? [];
}

export const authApi = {
  login: async (email: string, password: string) => {
    const response = await api.post("/auth/login", { email, password });
    return response.data;
  },
  verify: async (token: string) => {
    const response = await request("GET", "/newf/me", undefined, undefined, {
      Authorization: `Bearer ${token}`,
    });
    return response.data;
  },
};

export const usersApi = {
  getAll: async (): Promise<User[]> => {
    const response = await api.get("/admin/users");
    return asArray<User>(response.data);
  },
  create: async (user: Partial<User>) => {
    const filteredUser = Object.fromEntries(
      Object.entries(user).filter(
        ([, value]) => value !== "" && value !== null && value !== undefined,
      ),
    );
    const response = await api.post("/admin/users", filteredUser);
    return response.data;
  },
  update: async (email: string, user: Partial<User>) => {
    try {
      const response = await api.patch(`/admin/users/${encodeURIComponent(email)}`, user);
      return response.data;
    } catch (error: unknown) {
      console.error("Error updating user:", error);
      throw error;
    }
  },
  updateRoles: async (email: string, roles: string[]) => {
    const response = await api.patch(`/admin/users/${encodeURIComponent(email)}/roles`, { roles });
    return response.data;
  },
  deleteUser: async (email: string) => {
    const response = await api.delete(`/admin/users/${encodeURIComponent(email)}`);
    return response.data;
  },
  validateUser: async (email: string) => {
    const response = await api.post(`/admin/users/${encodeURIComponent(email)}/validate`);
    return response.data;
  },
};

export const eventsApi = {
  getAll: async (): Promise<Event[]> => {
    const response = await api.get("/admin/events");
    return asArray<Event>(response.data).sort(
      (a, b) => new Date(a.start_date).getTime() - new Date(b.start_date).getTime(),
    );
  },
  create: async (event: Partial<Event>) => {
    const filteredEvent = Object.fromEntries(
      Object.entries(event).filter(
        ([, value]) => value !== "" && value !== null && value !== undefined,
      ),
    );
    const response = await api.post("/admin/events", filteredEvent);
    return response.data;
  },
  update: async (id: number, event: Partial<Event>) => {
    const response = await api.patch(`/admin/events/${id}`, event);
    return response.data;
  },
  delete: async (id: number) => {
    const response = await api.delete(`/admin/events/${id}`);
    return response.data;
  },
};

export const clubsApi = {
  getAll: async (): Promise<Club[]> => {
    const response = await api.get("/admin/clubs");
    return asArray<Club>(response.data);
  },
  getOwners: async (id: number): Promise<User[]> => {
    const response = await api.get(`/club/${id}`);
    const data = response.data;
    if (Array.isArray(data.responsibles)) {
      return data.responsibles;
    }
    if (data.responsible) {
      return [data.responsible];
    }
    return [];
  },
  create: async (club: Partial<Club>) => {
    const filteredClub = Object.fromEntries(
      Object.entries(club).filter(
        ([, value]) => value !== "" && value !== null && value !== undefined,
      ),
    );
    const response = await api.post("/admin/clubs", filteredClub);
    return response.data;
  },
  update: async (id: number, club: Partial<Club>) => {
    const response = await api.patch(`/admin/clubs/${id}`, club);
    return response.data;
  },
  delete: async (id: number) => {
    const response = await api.delete(`/admin/clubs/${id}`);
    return response.data;
  },
  addOwner: async (clubId: number, email: string) => {
    const response = await api.post(`/club/${clubId}/respo`, { email });
    return response.data;
  },
  removeOwner: async (clubId: number, email: string) => {
    const response = await api.delete(`/club/${clubId}/respo`, {
      data: { email },
    });
    return response.data;
  },
};

export const statsApi = {
  getDashboard: async (): Promise<DashboardStats> => {
    const response = await api.get("/statistics/dashboard");
    return response.data;
  },
  getGlobal: async (): Promise<GlobalApiStats | null> => {
    const response = await api.get("/statistics/global");
    return response.data?.statistics ?? null;
  },
  getEndpoints: async (): Promise<EndpointApiStat[]> => {
    const response = await api.get("/statistics/endpoints");
    return asArray<EndpointApiStat>(response.data?.statistics);
  },
  getActiveUsers: async (period: ActiveUsersPeriod): Promise<ActiveUsersPoint[]> => {
    const response = await api.get("/statistics/active-users", { params: { period } });
    return asArray<ActiveUsersPoint>(response.data?.data);
  },
  getActivityHeatmap: async (range: ActivityHeatmapRange): Promise<ActivityHourPoint[]> => {
    const response = await api.get("/statistics/activity-heatmap", { params: { range } });
    return asArray<ActivityHourPoint>(response.data?.data);
  },
};

export const rolesApi = {
  getAll: async (): Promise<{ id_roles: number; name: string }[]> => {
    const response = await api.get("/admin/roles");
    return asArray<{ id_roles: number; name: string }>(response.data);
  },
};

export const notificationsApi = {
  send: async (request: SendNotificationRequest): Promise<NotificationRecipients> => {
    const response = await api.post("/admin/notifications/send", request);
    return response.data;
  },
};

export const menuApi = {
  getAll: async (): Promise<MenuItem[]> => {
    const response = await api.get("/admin/menu");
    return asArray<MenuItem>(response.data);
  },
  delete: async (id: number) => {
    const response = await api.delete(`/admin/menu/${id}`);
    return response.data;
  },
  getReviews: async (id: number): Promise<MenuItemReview[]> => {
    const response = await api.get(`/admin/menu/${id}/reviews`);
    return asArray<MenuItemReview>(response.data);
  },
  deleteReview: async (id: number, email: string) => {
    const response = await api.delete(`/admin/menu/${id}/reviews/${encodeURIComponent(email)}`);
    return response.data;
  },
};

export const reviewsApi = {
  getAll: async (userEmail?: string): Promise<MenuItemReview[]> => {
    const params = userEmail ? `?user_email=${encodeURIComponent(userEmail)}` : "";
    const response = await api.get(`/admin/reviews${params}`);
    return asArray<MenuItemReview>(response.data);
  },
};

export const bassineApi = {
  getScores: async (): Promise<BassineScore[]> => {
    const response = await api.get("/admin/bassine/scores");
    return asArray<BassineScore>(response.data);
  },
  updateScore: async (request: UpdateBassineScoreRequest) => {
    const response = await api.post("/admin/bassine/update-score", request);
    return response.data;
  },
  getHistory: async (email: string): Promise<BassineScoreHistory[]> => {
    const response = await api.get(`/admin/bassine/history/${encodeURIComponent(email)}`);
    return asArray<BassineScoreHistory>(response.data);
  },
};

export interface ReservationItem {
  id: number;
  name: string;
  slot: boolean;
  description?: string;
  location?: string;
  warning_message?: string;
  confirmation_message?: string;
}

export interface UpdateReservationItemMessagesRequest {
  warning_message?: string | null;
  confirmation_message?: string | null;
}

export interface ReservationCategory {
  type: "category";
  id: number;
  name: string;
  club_id: number;
  club_name: string;
  parent_id?: number;
  children?: ReservationCategory[];
  items: ReservationItem[];
}

export interface ReservationTreeItem {
  type: "category" | "club_items";
  id?: number;
  name?: string;
  club_id?: number;
  club_name?: string;
  parent_id?: number;
  children?: ReservationTreeItem[];
  items?: ReservationItem[];
}

export interface CreateCategoryRequest {
  name: string;
  id_club_parent?: number;
  id_category_parent?: number;
}

export interface UpdateCategoryRequest {
  name: string;
}

export interface CreateItemRequest {
  name: string;
  slot: boolean;
  description?: string;
  location?: string;
  id_club_parent?: number;
  id_category_parent?: number;
}

export interface UpdateItemRequest {
  name?: string;
  slot?: boolean;
  description?: string | null;
  location?: string | null;
  warning_message?: string | null;
  confirmation_message?: string | null;
}

export const traqApi = {
  getAllArticles: async (): Promise<TraqArticle[]> => {
    const response = await api.get("/traq/");
    return asArray<TraqArticle>(response.data);
  },
  createArticle: async (article: Partial<TraqArticle>) => {
    const response = await api.post("/traq/", article);
    return response.data;
  },
  updateArticle: async (id: number, article: Partial<TraqArticle>) => {
    const response = await api.patch(`/traq/${id}`, article);
    return response.data;
  },
  deleteArticle: async (id: number) => {
    const response = await api.delete(`/traq/${id}`);
    return response.data;
  },
  getAllTypes: async (): Promise<TraqType[]> => {
    const response = await api.get("/traq/types/");
    return asArray<TraqType>(response.data);
  },
  createType: async (type: { name: string }) => {
    const response = await api.post("/traq/types/", type);
    return response.data;
  },
  updateType: async (id: number, type: { name: string }) => {
    const response = await api.patch(`/traq/types/${id}`, type);
    return response.data;
  },
  deleteType: async (id: number) => {
    const response = await api.delete(`/traq/types/${id}`);
    return response.data;
  },
};

export const reservationApi = {
  getItemsForClub: async (clubId: number): Promise<ReservationItem[]> => {
    const response = await api.get(`/admin/clubs/${clubId}/reservation-items`);
    return asArray<ReservationItem>(response.data);
  },
  updateItemMessages: async (itemId: number, messages: UpdateReservationItemMessagesRequest) => {
    const response = await api.patch(`/admin/reservation-items/${itemId}/messages`, messages);
    return response.data;
  },
  getTree: async (): Promise<ReservationTreeItem[]> => {
    const response = await api.get("/admin/reservations/tree");
    return asArray<ReservationTreeItem>(response.data);
  },
  createCategory: async (category: CreateCategoryRequest) => {
    const response = await api.post("/reservation/category", category);
    return response.data;
  },
  updateCategory: async (id: number, category: UpdateCategoryRequest) => {
    const response = await api.patch(`/admin/reservations/categories/${id}`, category);
    return response.data;
  },
  deleteCategory: async (id: number) => {
    const response = await api.delete(`/admin/reservations/categories/${id}`);
    return response.data;
  },
  createItem: async (item: CreateItemRequest) => {
    const response = await api.post("/reservation/item", item);
    return response.data;
  },
  updateItem: async (id: number, item: UpdateItemRequest) => {
    const response = await api.patch(`/admin/reservations/items/${id}`, item);
    return response.data;
  },
  deleteItem: async (id: number) => {
    const response = await api.delete(`/admin/reservations/items/${id}`);
    return response.data;
  },
};
