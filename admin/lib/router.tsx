"use client";

import { lazy, Suspense } from "react";
import ErrorBoundary from "@/components/ErrorBoundary";
import { FullPageLoading } from "@/components/LoadingSpinner";

// Route-based code splitting - lazy load all pages
const DashboardPage = lazy(() => import("@/app/dashboard/page"));
const UsersPage = lazy(() => import("@/app/users/page"));
const EventsPage = lazy(() => import("@/app/events/page"));
const ClubsPage = lazy(() => import("@/app/clubs/page"));
const LoginPage = lazy(() => import("@/app/login/page"));

interface RouteConfig {
  path: string;
  component: React.LazyExoticComponent<React.ComponentType<Record<string, unknown>>>;
  loadingMessage?: string;
  requiresAuth?: boolean;
}

// Define all routes with their lazy-loaded components
export const routes: RouteConfig[] = [
  {
    path: "/dashboard",
    component: DashboardPage,
    loadingMessage: "Chargement du tableau de bord...",
    requiresAuth: true,
  },
  {
    path: "/users",
    component: UsersPage,
    loadingMessage: "Chargement de la gestion des utilisateurs...",
    requiresAuth: true,
  },
  {
    path: "/events",
    component: EventsPage,
    loadingMessage: "Chargement de la gestion des événements...",
    requiresAuth: true,
  },
  {
    path: "/clubs",
    component: ClubsPage,
    loadingMessage: "Chargement de la gestion des clubs...",
    requiresAuth: true,
  },
  {
    path: "/login",
    component: LoginPage,
    loadingMessage: "Chargement de la page de connexion...",
    requiresAuth: false,
  },
];

// Higher-order component for route-level lazy loading with error boundaries
export function withRouteLoading<T extends Record<string, unknown>>(
  Component: React.LazyExoticComponent<React.ComponentType<T>>,
  loadingMessage?: string,
) {
  return function LazyRoute(props: T) {
    return (
      <ErrorBoundary level="page">
        <Suspense
          fallback={<FullPageLoading text={loadingMessage || "Chargement de la page..."} />}
        >
          <Component {...props} />
        </Suspense>
      </ErrorBoundary>
    );
  };
}

// Pre-wrapped lazy components for each route
export const LazyDashboard = withRouteLoading(DashboardPage, "Chargement du tableau de bord...");
export const LazyUsers = withRouteLoading(
  UsersPage,
  "Chargement de la gestion des utilisateurs...",
);
export const LazyEvents = withRouteLoading(
  EventsPage,
  "Chargement de la gestion des événements...",
);
export const LazyClubs = withRouteLoading(ClubsPage, "Chargement de la gestion des clubs...");
export const LazyLogin = withRouteLoading(LoginPage, "Chargement de la page de connexion...");

// Route helper functions
export function getRouteConfig(pathname: string): RouteConfig | undefined {
  return routes.find((route) => route.path === pathname);
}

export function isAuthRequired(pathname: string): boolean {
  const route = getRouteConfig(pathname);
  return route?.requiresAuth ?? true;
}

export function getLoadingMessage(pathname: string): string {
  const route = getRouteConfig(pathname);
  return route?.loadingMessage || "Chargement...";
}
