import { useCallback, useState } from "react";

/**
 * Custom hook for managing boolean toggle state
 * @param initialValue - Initial boolean value (default: false)
 * @returns [value, toggle, setTrue, setFalse] tuple
 */
export function useToggle(initialValue: boolean = false) {
  const [value, setValue] = useState<boolean>(initialValue);

  const toggle = useCallback(() => setValue((v) => !v), []);
  const setTrue = useCallback(() => setValue(true), []);
  const setFalse = useCallback(() => setValue(false), []);

  return [value, toggle, setTrue, setFalse] as const;
}
