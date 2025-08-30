'use client';

import { useEffect, memo } from 'react';
import { useRouter, usePathname } from 'next/navigation';
import Sidebar from './Sidebar';
import ErrorBoundary from './ErrorBoundary';
import { CommandPalette } from './LazyComponents';
import { FullPageLoading } from './LoadingSpinner';
import { authApi } from '@/lib/api';
import { useAuthStore } from '@/lib/stores/authStore';
import { useAppStore } from '@/lib/stores/appStore';

interface LayoutProps {
  children: React.ReactNode;
}

function Layout({ children }: LayoutProps) {
  const router = useRouter();
  const pathname = usePathname();
  
  const { 
    isLoading, 
    setAuth, 
    logout, 
    setLoading 
  } = useAuthStore();
  
  const { setCurrentPage } = useAppStore();

  useEffect(() => {
    // Set current page for navigation state
    setCurrentPage(pathname);
  }, [pathname, setCurrentPage]);

  useEffect(() => {
    const checkAuth = async () => {
      if (pathname === '/login') {
        setLoading(false);
        return;
      }

      const token = localStorage.getItem('adminToken');
      if (!token) {
        logout();
        router.push('/login');
        return;
      }

      try {
        const userData = await authApi.verify(token);
        // Assuming the API returns user data with email and roles
        setAuth({ 
          email: userData.email, 
          roles: userData.roles || [] 
        }, token);
      } catch {
        logout();
        router.push('/login');
      }
    };

    checkAuth();
  }, [router, pathname, setAuth, logout, setLoading]);

  if (isLoading) {
    return <FullPageLoading text="Chargement de l'administration..." />;
  }

  if (pathname === '/login') {
    return <>{children}</>;
  }

  return (
    <ErrorBoundary level="page">
      <div className="lg:flex h-screen bg-gray-50">
        <Sidebar />
        <main className="flex-1 overflow-auto lg:ml-0">
          <ErrorBoundary level="section">
            {children}
          </ErrorBoundary>
        </main>
        <CommandPalette />
      </div>
    </ErrorBoundary>
  );
}

export default memo(Layout);