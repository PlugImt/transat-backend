/**
 * Performance utilities for React optimization
 */

import { useCallback, useRef } from "react";

// Debounce function for performance optimization
export function debounce<T extends (...args: unknown[]) => void>(
  func: T,
  wait: number,
): (...args: Parameters<T>) => void {
  let timeout: NodeJS.Timeout | null = null;

  return (...args: Parameters<T>) => {
    if (timeout) {
      clearTimeout(timeout);
    }

    timeout = setTimeout(() => {
      func(...args);
    }, wait);
  };
}

// Throttle function for performance optimization
export function throttle<T extends (...args: unknown[]) => void>(
  func: T,
  wait: number,
): (...args: Parameters<T>) => void {
  let timeout: NodeJS.Timeout | null = null;
  let previous = 0;

  return (...args: Parameters<T>) => {
    const now = Date.now();
    const remaining = wait - (now - previous);

    if (remaining <= 0 || remaining > wait) {
      if (timeout) {
        clearTimeout(timeout);
        timeout = null;
      }
      previous = now;
      func(...args);
    } else if (!timeout) {
      timeout = setTimeout(() => {
        previous = Date.now();
        timeout = null;
        func(...args);
      }, remaining);
    }
  };
}

// Hook for stable callback references
export function useStableCallback<T extends (...args: unknown[]) => unknown>(
  callback: T,
  dependencies: unknown[],
): T {
  const callbackRef = useRef<T>(callback);
  const depsRef = useRef(dependencies);

  // Update callback if dependencies changed
  if (dependencies.some((dep, i) => dep !== depsRef.current[i])) {
    callbackRef.current = callback;
    depsRef.current = dependencies;
  }

  return useCallback((...args: Parameters<T>) => {
    return callbackRef.current(...args);
  }, []) as T;
}

// Intersection Observer hook for lazy loading
export function useIntersectionObserver(
  callback: (entries: IntersectionObserverEntry[]) => void,
  options?: IntersectionObserverInit,
) {
  const observer = useRef<IntersectionObserver | null>(null);

  const observe = useCallback(
    (element: Element) => {
      if (observer.current) {
        observer.current.disconnect();
      }

      observer.current = new IntersectionObserver(callback, options);
      observer.current.observe(element);
    },
    [callback, options],
  );

  const disconnect = useCallback(() => {
    if (observer.current) {
      observer.current.disconnect();
    }
  }, []);

  return { observe, disconnect };
}

// Memory usage monitoring (development only)
export function logMemoryUsage(label: string) {
  if (process.env.NODE_ENV === "development" && "memory" in performance) {
    const memory = (
      performance as {
        memory?: {
          usedJSHeapSize: number;
          totalJSHeapSize: number;
          jsHeapSizeLimit: number;
        };
      }
    ).memory;
    if (memory) {
      console.log(`[${label}] Memory usage:`, {
        used: `${Math.round(memory.usedJSHeapSize / 1024 / 1024)} MB`,
        total: `${Math.round(memory.totalJSHeapSize / 1024 / 1024)} MB`,
        limit: `${Math.round(memory.jsHeapSizeLimit / 1024 / 1024)} MB`,
      });
    }
  }
}

// Performance timing helper
export function measurePerformance(name: string, fn: () => void) {
  if (process.env.NODE_ENV === "development") {
    const start = performance.now();
    fn();
    const end = performance.now();
    console.log(`[Performance] ${name}: ${(end - start).toFixed(2)}ms`);
  } else {
    fn();
  }
}

// Component render tracking (development only)
export function useRenderTracker(componentName: string) {
  const renderCount = useRef(0);

  if (process.env.NODE_ENV === "development") {
    renderCount.current++;
    console.log(`[Render] ${componentName} rendered ${renderCount.current} times`);
  }
}

// Image preloader for better UX
export function preloadImages(urls: string[]): Promise<undefined[]> {
  return Promise.all(
    urls.map(
      (url) =>
        new Promise<void>((resolve, reject) => {
          const img = new Image();
          img.onload = () => resolve();
          img.onerror = reject;
          img.src = url;
        }),
    ),
  );
}
