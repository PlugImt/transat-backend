"use client";

import { lazy, Suspense } from 'react';
import { FullPageLoading, ComponentLoading } from './LoadingSpinner';

// Lazy load page components for code splitting
export const LazyDashboardPage = lazy(() => import('@/app/dashboard/page'));
export const LazyUsersPage = lazy(() => import('@/app/users/page'));
export const LazyEventsPage = lazy(() => import('@/app/events/page'));
export const LazyClubsPage = lazy(() => import('@/app/clubs/page'));

// Lazy load large components
export const LazyDataTable = lazy(() => import('./DataTable'));
export const LazyUserModal = lazy(() => import('./UserModal'));
export const LazyCommandPalette = lazy(() => import('./CommandPalette'));

// Higher-order component for wrapping lazy components with Suspense
export function withLazyLoading<P extends object>(
  Component: React.LazyExoticComponent<React.ComponentType<P>>,
  fallback?: React.ReactNode
) {
  return function LazyWrapper(props: P) {
    return (
      <Suspense fallback={fallback || <ComponentLoading />}>
        <Component {...props} />
      </Suspense>
    );
  };
}

// Higher-order component for page-level lazy loading
export function withPageLazyLoading<P extends object>(
  Component: React.LazyExoticComponent<React.ComponentType<P>>,
  loadingText?: string
) {
  return function LazyPageWrapper(props: P) {
    return (
      <Suspense fallback={<FullPageLoading text={loadingText || "Chargement de la page..."} />}>
        <Component {...props} />
      </Suspense>
    );
  };
}

// Lazy wrapped components ready to use
export const DashboardPage = withPageLazyLoading(LazyDashboardPage, "Chargement du tableau de bord...");
export const UsersPage = withPageLazyLoading(LazyUsersPage, "Chargement de la gestion des utilisateurs...");
export const EventsPage = withPageLazyLoading(LazyEventsPage, "Chargement de la gestion des événements...");
export const ClubsPage = withPageLazyLoading(LazyClubsPage, "Chargement de la gestion des clubs...");

export const DataTable = withLazyLoading(LazyDataTable);
export const UserModal = withLazyLoading(LazyUserModal);
export const CommandPalette = withLazyLoading(LazyCommandPalette);