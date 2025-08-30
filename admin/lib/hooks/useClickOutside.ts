import { useEffect, useRef } from 'react';

/**
 * Custom hook for detecting clicks outside of a component
 * @param callback - Function to call when clicking outside
 * @returns Ref to attach to the element you want to detect outside clicks for
 */
export function useClickOutside<T extends HTMLElement = HTMLElement>(
  callback: () => void
) {
  const ref = useRef<T>(null);

  useEffect(() => {
    const handleClick = (event: MouseEvent) => {
      if (ref.current && !ref.current.contains(event.target as Node)) {
        callback();
      }
    };

    document.addEventListener('mousedown', handleClick);
    return () => document.removeEventListener('mousedown', handleClick);
  }, [callback]);

  return ref;
}