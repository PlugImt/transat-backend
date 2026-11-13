"use client";

import { format } from "date-fns";
import { fr } from "date-fns/locale";
import {
  ChevronLeft,
  ChevronRight,
  Copy,
  ExternalLink,
  Eye,
  EyeOff,
  FolderOpen,
  ImageOff,
  Images,
  Loader2,
  Search,
  ShieldQuestion,
} from "lucide-react";
import { useEffect, useState } from "react";
import toast from "react-hot-toast";
import { fetchAuthorizedImage } from "@/lib/api";
import { usePhotoAlbums, useSetAlbumHidden } from "@/lib/hooks";
import type { PhotoAlbum, PhotoAlbumStatus } from "@/lib/types";

const TABS: { id: PhotoAlbumStatus | ""; label: string; icon: typeof Eye }[] = [
  { id: "", label: "Tous", icon: Images },
  { id: "shown", label: "Visibles", icon: Eye },
  { id: "hidden", label: "Masqués", icon: EyeOff },
  { id: "review", label: "À vérifier", icon: ShieldQuestion },
];

const COUNT_KEY = { "": "all", shown: "shown", hidden: "hidden", review: "review" } as const;

const dateLabel = (album: PhotoAlbum) => {
  if (album.event_date) {
    return format(new Date(`${album.event_date}T12:00:00`), "d MMMM yyyy", { locale: fr });
  }
  if (album.first_taken_at) {
    return `Photos du ${format(new Date(album.first_taken_at), "d MMMM yyyy", { locale: fr })}`;
  }
  return "Date inconnue";
};

export default function PhotoAlbumsPage() {
  const [status, setStatus] = useState<PhotoAlbumStatus | "">("");
  const [input, setInput] = useState("");
  const [q, setQ] = useState("");
  const [page, setPage] = useState(1);

  // Search as you type, without one request per keystroke.
  useEffect(() => {
    const t = setTimeout(() => {
      setQ(input.trim());
      setPage(1);
    }, 250);
    return () => clearTimeout(t);
  }, [input]);

  const { data, isLoading, isError, isFetching } = usePhotoAlbums({ q, status, page });
  const albums = data?.albums ?? [];
  const total = data?.total ?? 0;
  const pageSize = data?.page_size ?? 50;
  const first = (page - 1) * pageSize;
  const last = Math.min(first + pageSize, total);

  return (
    <div className="mx-auto max-w-7xl p-6 lg:p-8">
      <header className="mb-6 flex flex-wrap items-end justify-between gap-4 animate-rise">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-gray-900">Albums photo</h1>
          <p className="mt-1 max-w-3xl text-sm text-gray-500">
            Les albums sont créés automatiquement depuis le Drive Photos. Un album masqué disparaît,
            avec ses photos, de l'application et de sa page transat.dev ; la synchronisation du
            Drive ne le réaffiche jamais.
          </p>
        </div>
        <label className="relative w-full sm:w-72">
          <Search className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
          <input
            type="search"
            value={input}
            onChange={(e) => setInput(e.target.value)}
            placeholder="Rechercher un album…"
            aria-label="Rechercher un album"
            className="w-full rounded-xl border border-gray-300 bg-white py-2 pl-9 pr-3 text-sm text-gray-900"
          />
        </label>
      </header>

      <div
        role="tablist"
        aria-label="Statut"
        className="mb-6 inline-flex max-w-full gap-1 overflow-x-auto rounded-2xl border border-gray-200 bg-white p-1"
      >
        {TABS.map(({ id, label, icon: Icon }) => {
          const active = status === id;
          return (
            <button
              key={id || "all"}
              type="button"
              role="tab"
              aria-selected={active}
              onClick={() => {
                setStatus(id);
                setPage(1);
              }}
              className={`flex shrink-0 items-center gap-2 rounded-xl px-4 py-2 text-sm font-medium transition-all duration-200 ${
                active ? "bg-blue-600 text-white shadow" : "text-gray-600 hover:bg-gray-100"
              }`}
            >
              <Icon className="h-4 w-4" />
              {label}
              {data && (
                <span
                  className={`rounded-full px-2 py-0.5 text-xs ${active ? "bg-white/20" : "bg-gray-100 text-gray-600"}`}
                >
                  {data.counts[COUNT_KEY[id]]}
                </span>
              )}
            </button>
          );
        })}
      </div>

      {isError && (
        <p className="rounded-2xl border border-red-200 bg-red-50 p-6 text-center text-sm text-red-700">
          Impossible de charger les albums.
        </p>
      )}

      {isLoading && (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
          {Array.from({ length: 8 }, (_, i) => (
            // biome-ignore lint/suspicious/noArrayIndexKey: static placeholders
            <div key={i} className="h-72 animate-pulse rounded-2xl bg-gray-100" />
          ))}
        </div>
      )}

      {!isLoading && !isError && albums.length === 0 && (
        <div className="rounded-2xl border border-dashed border-gray-300 py-16 text-center animate-rise">
          <FolderOpen className="mx-auto h-10 w-10 text-gray-400" />
          <h3 className="mt-3 text-sm font-semibold text-gray-900">Aucun album</h3>
          <p className="mt-1 text-sm text-gray-500">
            {q ? "Aucun album ne correspond à cette recherche." : "Rien dans cette catégorie."}
          </p>
        </div>
      )}

      <ul
        className={`grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 transition-opacity ${isFetching ? "opacity-70" : ""}`}
      >
        {albums.map((album, i) => (
          <li
            key={album.id}
            className="animate-rise"
            style={{ animationDelay: `${Math.min(i, 12) * 35}ms` }}
          >
            <AlbumCard album={album} />
          </li>
        ))}
      </ul>

      {total > pageSize && (
        <div className="mt-6 flex items-center justify-between text-sm text-gray-600">
          <span>
            {first + 1}–{last} sur {total}
          </span>
          <div className="flex gap-2">
            <button
              type="button"
              disabled={page === 1}
              onClick={() => setPage(page - 1)}
              className="rounded-xl border border-gray-300 bg-white p-2 transition-colors hover:bg-gray-50 disabled:opacity-40"
              aria-label="Page précédente"
            >
              <ChevronLeft className="h-4 w-4" />
            </button>
            <button
              type="button"
              disabled={last >= total}
              onClick={() => setPage(page + 1)}
              className="rounded-xl border border-gray-300 bg-white p-2 transition-colors hover:bg-gray-50 disabled:opacity-40"
              aria-label="Page suivante"
            >
              <ChevronRight className="h-4 w-4" />
            </button>
          </div>
        </div>
      )}
    </div>
  );
}

function AlbumCard({ album }: { album: PhotoAlbum }) {
  const setHidden = useSetAlbumHidden();
  const pending = setHidden.isPending;

  const toggle = () => {
    const hide = !album.hidden;
    if (
      hide &&
      !window.confirm(
        `Masquer « ${album.title} » ? L'album et ses ${album.photo_count} photos disparaîtront de l'application et de transat.dev.`,
      )
    ) {
      return;
    }
    setHidden.mutate(
      { id: album.id, hidden: hide },
      {
        onSuccess: () => toast.success(hide ? "Album masqué" : "Album de nouveau visible"),
        onError: () => toast.error("Impossible de modifier l'album"),
      },
    );
  };

  const copy = async () => {
    try {
      await navigator.clipboard.writeText(album.share_url);
      toast.success("Lien copié");
    } catch {
      toast.error("Impossible de copier le lien");
    }
  };

  const badge = album.hidden
    ? { label: "Masqué", className: "bg-gray-900/80 text-white" }
    : album.visibility === "published"
      ? { label: "Visible", className: "bg-emerald-100 text-emerald-800" }
      : { label: "À vérifier", className: "bg-amber-100 text-amber-800" };

  return (
    <article className="flex h-full flex-col overflow-hidden rounded-2xl border border-gray-200 bg-white shadow-sm transition-shadow duration-200 hover:shadow-md">
      <div className="relative">
        <Cover
          path={album.cover_path}
          alt={album.title}
          dimmed={!album.shown}
          className="aspect-[4/3] w-full"
        />
        <span
          className={`absolute left-3 top-3 rounded-full px-2.5 py-1 text-xs font-medium shadow-sm ${badge.className}`}
        >
          {badge.label}
        </span>
      </div>
      <div className="flex flex-1 flex-col gap-1 p-4">
        <h2 className="line-clamp-2 text-sm font-semibold text-gray-900" title={album.title}>
          {album.title}
        </h2>
        <p className="text-xs text-gray-500">
          {dateLabel(album)} · {album.photo_count} photo{album.photo_count > 1 ? "s" : ""}
        </p>
        {album.category && <p className="truncate text-xs text-gray-400">{album.category}</p>}
        {album.hidden && album.hidden_by && (
          <p className="truncate text-xs text-gray-500" title={album.hidden_by}>
            Masqué par {album.hidden_by}
            {album.hidden_at &&
              ` le ${format(new Date(album.hidden_at), "d MMM yyyy", { locale: fr })}`}
          </p>
        )}
        {!album.hidden && album.visibility !== "published" && (
          <p className="text-xs text-amber-700">
            Le classement automatique n'a pas publié cet album : il n'est pas visible.
          </p>
        )}

        <div className="mt-auto flex items-center gap-2 pt-3">
          <button
            type="button"
            onClick={toggle}
            disabled={pending}
            className={`flex flex-1 items-center justify-center gap-1.5 rounded-xl px-3 py-2 text-xs font-medium transition-colors disabled:opacity-60 ${
              album.hidden
                ? "bg-blue-600 text-white hover:bg-blue-700"
                : "border border-gray-300 text-gray-700 hover:bg-gray-50"
            }`}
          >
            {pending ? (
              <Loader2 className="h-3.5 w-3.5 animate-spin" />
            ) : album.hidden ? (
              <Eye className="h-3.5 w-3.5" />
            ) : (
              <EyeOff className="h-3.5 w-3.5" />
            )}
            {album.hidden ? "Afficher" : "Masquer"}
          </button>
          <button
            type="button"
            onClick={copy}
            className="rounded-xl border border-gray-300 p-2 text-gray-600 transition-colors hover:bg-gray-50"
            aria-label="Copier le lien transat.dev"
            title={album.share_url}
          >
            <Copy className="h-4 w-4" />
          </button>
          {album.shown ? (
            <a
              href={album.share_url}
              target="_blank"
              rel="noopener noreferrer"
              className="rounded-xl border border-gray-300 p-2 text-gray-600 transition-colors hover:bg-gray-50"
              aria-label="Ouvrir la page transat.dev"
              title="Ouvrir la page transat.dev"
            >
              <ExternalLink className="h-4 w-4" />
            </a>
          ) : (
            <span
              className="cursor-not-allowed rounded-xl border border-gray-200 p-2 text-gray-300"
              title="La page n'est accessible que pour un album visible"
            >
              <ExternalLink className="h-4 w-4" />
            </span>
          )}
        </div>
      </div>
    </article>
  );
}

/** The album's cover, read with the panel's token so hidden albums keep theirs. */
function Cover({
  path,
  alt,
  dimmed,
  className = "",
}: {
  path: string | null;
  alt: string;
  dimmed: boolean;
  className?: string;
}) {
  const [src, setSrc] = useState<string | null>(null);
  const [failed, setFailed] = useState(false);

  useEffect(() => {
    if (!path) return;
    const controller = new AbortController();
    let url: string | null = null;
    setFailed(false);
    fetchAuthorizedImage(path, controller.signal)
      .then((u) => {
        url = u;
        setSrc(u);
      })
      .catch(() => {
        if (!controller.signal.aborted) setFailed(true);
      });
    return () => {
      controller.abort();
      if (url) URL.revokeObjectURL(url);
      setSrc(null);
    };
  }, [path]);

  if (!path || failed) {
    return (
      <div
        className={`flex flex-col items-center justify-center gap-2 bg-gray-100 text-gray-400 ${className}`}
      >
        <ImageOff className="h-7 w-7" />
        <span className="text-xs">Pas de couverture</span>
      </div>
    );
  }
  return (
    <div className={`relative overflow-hidden bg-gray-100 ${className}`}>
      {!src && <div className="absolute inset-0 animate-pulse bg-gray-200" />}
      {src && (
        // biome-ignore lint/performance/noImgElement: an object URL from an authorised fetch
        <img
          src={src}
          alt={alt}
          className={`h-full w-full object-cover transition-[opacity,filter] duration-500 animate-fade-in ${dimmed ? "opacity-60 grayscale" : ""}`}
        />
      )}
    </div>
  );
}
