"use client";

import { useCallback, useEffect, useState } from "react";

export type ThemeChoice = "light" | "dark" | "system";

const KEY = "transat-admin-theme";

/** Runs before the first paint (inlined in app/layout.tsx) so the page never flashes the wrong theme. */
export const themeBootScript = `(function(){try{var c=localStorage.getItem("${KEY}")||"system";var d=c==="dark"||(c==="system"&&window.matchMedia("(prefers-color-scheme: dark)").matches);document.documentElement.classList.toggle("dark",d);}catch(e){}})();`;

const readChoice = (): ThemeChoice => {
  try {
    const v = localStorage.getItem(KEY);
    return v === "light" || v === "dark" ? v : "system";
  } catch {
    return "system";
  }
};

const apply = (choice: ThemeChoice) => {
  const dark =
    choice === "dark" ||
    (choice === "system" && window.matchMedia("(prefers-color-scheme: dark)").matches);
  document.documentElement.classList.toggle("dark", dark);
};

/** Light, dark or the system's setting; remembered on this browser. */
export function useTheme() {
  const [choice, setChoice] = useState<ThemeChoice>("system");

  useEffect(() => {
    setChoice(readChoice());
  }, []);

  useEffect(() => {
    apply(choice);
    if (choice !== "system") return;
    const media = window.matchMedia("(prefers-color-scheme: dark)");
    const follow = () => apply("system");
    media.addEventListener("change", follow);
    return () => media.removeEventListener("change", follow);
  }, [choice]);

  const set = useCallback((next: ThemeChoice) => {
    setChoice(next);
    try {
      localStorage.setItem(KEY, next);
    } catch {
      // Private browsing: the choice lasts for this page only.
    }
  }, []);

  return { choice, setChoice: set };
}
