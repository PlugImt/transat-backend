"use client";

import { ImageOff } from "lucide-react";
import { useState } from "react";
import { photoPreviewFallbackUrl, photoPreviewUrl } from "@/lib/photoReports";

interface PhotoPreviewProps {
  photoId: string;
  alt: string;
  width: number;
  className?: string;
  /** The photo was deleted: Drive no longer serves it. */
  removed?: boolean;
}

/** The reported photo as Drive serves it, with Drive's other thumbnail address as a fallback. */
export function PhotoPreview({ photoId, alt, width, className = "", removed }: PhotoPreviewProps) {
  const [attempt, setAttempt] = useState(0);
  const [loaded, setLoaded] = useState(false);

  if (removed || attempt > 1) {
    return (
      <div
        className={`flex flex-col items-center justify-center gap-2 bg-gray-100 text-gray-400 ${className}`}
      >
        <ImageOff className="h-7 w-7" />
        <span className="text-xs">{removed ? "Photo supprimée" : "Aperçu indisponible"}</span>
      </div>
    );
  }
  const src =
    attempt === 0 ? photoPreviewUrl(photoId, width) : photoPreviewFallbackUrl(photoId, width);
  return (
    <div className={`relative overflow-hidden bg-gray-100 ${className}`}>
      {!loaded && <div className="absolute inset-0 animate-pulse bg-gray-200" />}
      {/* biome-ignore lint/performance/noImgElement: remote Drive thumbnails, not optimisable by next/image */}
      <img
        key={src}
        src={src}
        alt={alt}
        loading="lazy"
        referrerPolicy="no-referrer"
        onLoad={() => setLoaded(true)}
        onError={() => {
          setLoaded(false);
          setAttempt((a) => a + 1);
        }}
        className={`h-full w-full object-cover transition-opacity duration-500 ${loaded ? "opacity-100" : "opacity-0"}`}
      />
    </div>
  );
}
