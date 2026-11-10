export interface PanelAccess {
  email: string;
  roles: string[];
  admin: boolean;
  /** Top-level paths the account may open ("dashboard", "photo-reports"…); ["*"] for administrators. */
  sections: string[];
}

/** First segment of a panel path: "/photo-reports?report=3" -> "photo-reports". */
export const sectionOf = (path: string) =>
  path.replace(/^\/+/, "").split(/[/?#]/)[0] || "dashboard";

export const canOpen = (access: PanelAccess | null, path: string) =>
  !!access && (access.sections.includes("*") || access.sections.includes(sectionOf(path)));

/** Where to land when the requested page is not allowed. */
export const homeOf = (access: PanelAccess | null) =>
  !access || access.sections.includes("*") || access.sections.includes("dashboard")
    ? "/dashboard"
    : `/${access.sections[0] ?? "dashboard"}`;

/** Keeps a deep link (such as the report link of an email) through the login page. Only panel paths are kept. */
export const safeNext = (value: string | null) =>
  value?.startsWith("/") && !value.startsWith("//") && !value.startsWith("/login") ? value : null;
