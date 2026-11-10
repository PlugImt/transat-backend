"use client";

import { usePathname, useRouter } from "next/navigation";
import { memo, useEffect } from "react";
import { canOpen, homeOf } from "@/lib/access";
import { authApi } from "@/lib/api";
import { useAppStore } from "@/lib/stores/appStore";
import { useAuthStore } from "@/lib/stores/authStore";
import ErrorBoundary from "./ErrorBoundary";
import { CommandPalette } from "./LazyComponents";
import { FullPageLoading } from "./LoadingSpinner";
import Sidebar from "./Sidebar";

interface LayoutProps {
  children: React.ReactNode;
}

function Layout({ children }: LayoutProps) {
  const router = useRouter();
  const pathname = usePathname();

  const { isLoading, setAuth, logout, setLoading, access, setAccess } = useAuthStore();

  const { setCurrentPage } = useAppStore();

  useEffect(() => {
    // Set current page for navigation state
    setCurrentPage(pathname);
  }, [pathname, setCurrentPage]);

  useEffect(() => {
    const checkAuth = async () => {
      if (pathname === "/login") {
        setLoading(false);
        return;
      }

      // The page asked for (with its query, such as an email's ?report=) survives the login.
      const here = `${pathname}${window.location.search}`;
      const toLogin = () => router.push(`/login?next=${encodeURIComponent(here)}`);

      const token = localStorage.getItem("adminToken");
      if (!token) {
        logout();
        toLogin();
        return;
      }

      try {
        const [userData, panelAccess] = await Promise.all([
          authApi.verify(token),
          authApi.access(token),
        ]);
        setAccess(panelAccess);
        // Assuming the API returns user data with email and roles
        setAuth(
          {
            email: userData.email,
            roles: userData.roles || [],
            firstName: userData.first_name,
            lastName: userData.last_name,
          },
          token,
        );
      } catch {
        logout();
        toLogin();
      }
    };

    checkAuth();
  }, [router, pathname, setAuth, logout, setLoading, setAccess]);

  // A page outside the account's sections sends it to its own home page.
  const allowed = pathname === "/login" || canOpen(access, pathname);
  useEffect(() => {
    if (!isLoading && access && !allowed && access.sections.length > 0)
      router.replace(homeOf(access));
  }, [isLoading, access, allowed, router]);

  if (isLoading) {
    return <FullPageLoading text="Chargement de l'administration..." />;
  }

  if (pathname === "/login") {
    return <>{children}</>;
  }

  if (access && access.sections.length === 0) {
    return (
      <NoAccess
        email={access.email}
        onLogout={() => {
          logout();
          router.push("/login");
        }}
      />
    );
  }

  if (!allowed) {
    return <FullPageLoading text="Redirection…" />;
  }

  return (
    <ErrorBoundary level="page">
      <div className="lg:flex h-screen bg-[var(--page)]">
        <Sidebar />
        <main className="flex-1 overflow-auto lg:ml-0">
          <ErrorBoundary level="section">{children}</ErrorBoundary>
        </main>
        {/* The palette searches admin-only data. */}
        {access?.admin && <CommandPalette />}
      </div>
    </ErrorBoundary>
  );
}

function NoAccess({ email, onLogout }: { email: string; onLogout: () => void }) {
  return (
    <div className="min-h-screen flex items-center justify-center bg-[var(--page)] p-6">
      <div className="max-w-md w-full rounded-2xl border border-gray-200 bg-white p-8 text-center shadow-sm animate-rise">
        <h1 className="text-xl font-semibold text-gray-900">Accès non autorisé</h1>
        <p className="mt-2 text-sm text-gray-600">
          Le compte {email} n'a pas accès au panel d'administration. Demandez à un administrateur de
          vous donner le rôle adéquat.
        </p>
        <button
          type="button"
          onClick={onLogout}
          className="mt-6 inline-flex rounded-lg bg-blue-600 px-4 py-2 text-sm font-semibold text-white hover:bg-blue-700"
        >
          Se déconnecter
        </button>
      </div>
    </div>
  );
}

export default memo(Layout);
