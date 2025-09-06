import { useEffect } from 'react';

export interface KeyboardShortcut {
  key: string;
  ctrlKey?: boolean;
  metaKey?: boolean;
  shiftKey?: boolean;
  altKey?: boolean;
}

/**
 * Custom hook for handling keyboard shortcuts
 * @param shortcut - The keyboard shortcut configuration
 * @param callback - The function to call when the shortcut is pressed
 * @param dependencies - Optional dependencies array for the callback
 */
export function useKeyboardShortcut(
  shortcut: KeyboardShortcut,
  callback: () => void,
  dependencies: React.DependencyList = []
) {
  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      const {
        key,
        ctrlKey = false,
        metaKey = false,
        shiftKey = false,
        altKey = false,
      } = shortcut;

      // Check if all modifiers match
      const modifiersMatch =
        event.ctrlKey === ctrlKey &&
        event.metaKey === metaKey &&
        event.shiftKey === shiftKey &&
        event.altKey === altKey;

      // Check if key matches (case insensitive)
      const keyMatches = event.key.toLowerCase() === key.toLowerCase();

      if (modifiersMatch && keyMatches) {
        event.preventDefault();
        callback();
      }
    };

    document.addEventListener('keydown', handleKeyDown);
    return () => document.removeEventListener('keydown', handleKeyDown);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [shortcut.key, shortcut.ctrlKey, shortcut.metaKey, shortcut.shiftKey, shortcut.altKey, ...dependencies]);
}

/**
 * Multiple keyboard shortcuts hook
 * @param shortcuts - Array of shortcut configurations with callbacks
 */
export function useKeyboardShortcuts(
  shortcuts: Array<KeyboardShortcut & { callback: () => void }>,
  dependencies: React.DependencyList = []
) {
  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      for (const shortcut of shortcuts) {
        const {
          key,
          ctrlKey = false,
          metaKey = false,
          shiftKey = false,
          altKey = false,
          callback,
        } = shortcut;

        // Check if all modifiers match
        const modifiersMatch =
          event.ctrlKey === ctrlKey &&
          event.metaKey === metaKey &&
          event.shiftKey === shiftKey &&
          event.altKey === altKey;

        // Check if key matches (case insensitive)
        const keyMatches = event.key.toLowerCase() === key.toLowerCase();

        if (modifiersMatch && keyMatches) {
          event.preventDefault();
          callback();
          break; // Stop after first match
        }
      }
    };

    document.addEventListener('keydown', handleKeyDown);
    return () => document.removeEventListener('keydown', handleKeyDown);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [...shortcuts.map(s => [s.key, s.ctrlKey, s.metaKey, s.shiftKey, s.altKey]).flat(), ...dependencies]);
}