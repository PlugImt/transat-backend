"use client";

import { format, formatDistanceToNow } from "date-fns";
import { fr } from "date-fns/locale";
import {
  AlertTriangle,
  ArrowLeft,
  Camera,
  CheckCircle2,
  ExternalLink,
  ImageOff,
  Images,
  Loader2,
  Trash2,
  X,
} from "lucide-react";
import { useEffect, useRef, useState } from "react";
import toast from "react-hot-toast";
import Modal from "@/components/Modal";
import { fetchAuthorizedImage, photoboothReviewApi } from "@/lib/api";
import {
  useApprovePhotobooth,
  useDeletePhotoboothPicture,
  usePhotoboothPending,
  usePhotoboothPictures,
} from "@/lib/hooks";
import type { PhotoboothPendingFolder, PhotoboothPendingPicture } from "@/lib/types";

/** Pictures uploaded more recently than this suggest the evening is still going on. */
const RECENT_MS = 60 * 60 * 1000;

const eventDay = (date: string) =>
  format(new Date(`${date}T12:00:00`), "EEEE d MMMM yyyy", { locale: fr });

const ago = (value: string) =>
  formatDistanceToNow(new Date(value), { addSuffix: true, locale: fr });

const stillGoing = (f: PhotoboothPendingFolder) =>
  f.open_shoots > 0 ||
  (f.last_upload_at !== null && Date.now() - new Date(f.last_upload_at).getTime() < RECENT_MS);

export default function PhotoboothPage() {
  const { data: folders = [], isLoading, isError } = usePhotoboothPending();
  const [openId, setOpenId] = useState<string | null>(null);
  const open = folders.find((f) => f.folder_id === openId) ?? null;

  // The folder left the list (approved elsewhere, renamed on Drive): back to the list.
  useEffect(() => {
    if (openId && !isLoading && !open) setOpenId(null);
  }, [openId, open, isLoading]);

  if (open) return <FolderReview folder={open} onBack={() => setOpenId(null)} />;

  return (
    <div className="mx-auto max-w-7xl p-6 lg:p-8">
      <header className="mb-6 animate-rise">
        <h1 className="text-2xl font-bold tracking-tight text-gray-900">Photobooth</h1>
        <p className="mt-1 max-w-3xl text-sm text-gray-500">
          Les photos du photobooth arrivent dans le dossier Drive « photobooth_temp » de l'événement
          et restent invisibles dans l'application. En fin de soirée, vérifiez-les, supprimez celles
          qui ne doivent pas être publiées, puis approuvez : le dossier devient « Photobooth » et
          les photos apparaissent dans l'application.
        </p>
      </header>

      {isError && (
        <p className="rounded-2xl border border-red-200 bg-red-50 p-6 text-center text-sm text-red-700">
          Impossible de charger les photos en attente.
        </p>
      )}

      {isLoading && (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {Array.from({ length: 3 }, (_, i) => (
            // biome-ignore lint/suspicious/noArrayIndexKey: static placeholders
            <div key={i} className="h-40 animate-pulse rounded-2xl bg-gray-100" />
          ))}
        </div>
      )}

      {!isLoading && !isError && folders.length === 0 && (
        <div className="rounded-2xl border border-dashed border-gray-300 py-16 text-center animate-rise">
          <Camera className="mx-auto h-10 w-10 text-gray-400" />
          <h3 className="mt-3 text-sm font-semibold text-gray-900">Rien à valider</h3>
          <p className="mt-1 text-sm text-gray-500">
            Les photos d'une soirée avec le photobooth apparaîtront ici.
          </p>
        </div>
      )}

      <ul className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {folders.map((f, i) => (
          <li
            key={f.folder_id}
            className="animate-rise"
            style={{ animationDelay: `${Math.min(i, 12) * 35}ms` }}
          >
            <button
              type="button"
              onClick={() => setOpenId(f.folder_id)}
              className="w-full rounded-2xl border border-gray-200 bg-white p-5 text-left shadow-sm transition-all duration-200 hover:-translate-y-0.5 hover:shadow-md"
            >
              <div className="flex items-start justify-between gap-3">
                <div className="min-w-0">
                  <h2 className="truncate text-base font-semibold text-gray-900">{f.event_name}</h2>
                  <p className="mt-0.5 text-sm capitalize text-gray-500">
                    {eventDay(f.event_date)}
                  </p>
                </div>
                <span className="flex shrink-0 items-center gap-1.5 rounded-full bg-blue-50 px-2.5 py-1 text-xs font-medium text-blue-700">
                  <Images className="h-3.5 w-3.5" />
                  {f.pictures}
                </span>
              </div>
              <p className="mt-4 text-xs text-gray-500">
                {f.shoots} passage{f.shoots > 1 ? "s" : ""}
                {f.last_upload_at && <> · dernière photo {ago(f.last_upload_at)}</>}
              </p>
              {stillGoing(f) && (
                <p className="mt-3 flex items-center gap-1.5 text-xs font-medium text-amber-700">
                  <AlertTriangle className="h-3.5 w-3.5" />
                  La soirée semble encore en cours
                </p>
              )}
            </button>
          </li>
        ))}
      </ul>
    </div>
  );
}

function FolderReview({ folder, onBack }: { folder: PhotoboothPendingFolder; onBack: () => void }) {
  const { data: pictures = [], isLoading, isError } = usePhotoboothPictures(folder.folder_id);
  const approve = useApprovePhotobooth();
  const remove = useDeletePhotoboothPicture();
  const [zoomed, setZoomed] = useState<PhotoboothPendingPicture | null>(null);
  const [toDelete, setToDelete] = useState<PhotoboothPendingPicture | null>(null);
  const [confirmApprove, setConfirmApprove] = useState(false);
  const going = stillGoing(folder);

  const doApprove = () =>
    approve.mutate(folder.folder_id, {
      onSuccess: (res) => {
        setConfirmApprove(false);
        toast.success(
          res.merged
            ? "Photos approuvées et ajoutées au dossier Photobooth de l'événement"
            : "Photos approuvées : elles apparaîtront dans l'application d'ici quelques minutes",
        );
        onBack();
      },
      onError: (e) => toast.error(e.message || "L'approbation a échoué"),
    });

  const doDelete = (p: PhotoboothPendingPicture) =>
    remove.mutate(p.id, {
      onSuccess: () => {
        setToDelete(null);
        if (zoomed?.id === p.id) setZoomed(null);
        toast.success("Photo supprimée");
      },
      onError: (e) => toast.error(e.message || "La suppression a échoué"),
    });

  return (
    <div className="mx-auto max-w-7xl p-6 lg:p-8">
      <button
        type="button"
        onClick={onBack}
        className="mb-4 flex items-center gap-1.5 text-sm text-gray-500 transition-colors hover:text-gray-900"
      >
        <ArrowLeft className="h-4 w-4" />
        Toutes les soirées
      </button>

      <header className="mb-6 flex flex-wrap items-end justify-between gap-4 animate-rise">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-gray-900">{folder.event_name}</h1>
          <p className="mt-1 text-sm text-gray-500">
            <span className="capitalize">{eventDay(folder.event_date)}</span> · {folder.pictures}{" "}
            photo{folder.pictures > 1 ? "s" : ""} en attente ·{" "}
            <a
              href={folder.drive_url}
              target="_blank"
              rel="noreferrer"
              className="inline-flex items-center gap-1 text-blue-600 hover:underline"
            >
              Ouvrir sur Drive <ExternalLink className="h-3.5 w-3.5" />
            </a>
          </p>
        </div>
        <button
          type="button"
          onClick={() => setConfirmApprove(true)}
          disabled={pictures.length === 0}
          className="flex items-center gap-2 rounded-xl bg-blue-600 px-4 py-2.5 text-sm font-semibold text-white shadow transition-colors hover:bg-blue-700 disabled:opacity-50"
        >
          <CheckCircle2 className="h-4 w-4" />
          Approuver et publier
        </button>
      </header>

      {going && (
        <p className="mb-6 flex items-start gap-2 rounded-2xl border border-amber-200 bg-amber-50 p-4 text-sm text-amber-800">
          <AlertTriangle className="mt-0.5 h-4 w-4 shrink-0" />
          {folder.open_shoots > 0
            ? "Des invités sont peut-être encore devant le photobooth."
            : "La dernière photo date de moins d'une heure : la soirée semble encore en cours."}{" "}
          Les photos prises après l'approbation iront dans un nouveau dossier à valider.
        </p>
      )}

      {isError && (
        <p className="rounded-2xl border border-red-200 bg-red-50 p-6 text-center text-sm text-red-700">
          Impossible de charger les photos.
        </p>
      )}

      <ul className="grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5">
        {isLoading &&
          Array.from({ length: 10 }, (_, i) => (
            // biome-ignore lint/suspicious/noArrayIndexKey: static placeholders
            <li key={i} className="aspect-[2/3] animate-pulse rounded-xl bg-gray-100" />
          ))}
        {pictures.map((p) => (
          <li key={p.id} className="group relative">
            <button
              type="button"
              onClick={() => setZoomed(p)}
              className="block w-full overflow-hidden rounded-xl"
              aria-label={`Agrandir ${p.name}`}
            >
              <AuthImage
                path={photoboothReviewApi.picturePath(p.id, "thumb")}
                alt={p.name}
                className="aspect-[2/3] w-full"
                fit="cover"
              />
            </button>
            {p.kind === "collage" && (
              <span className="pointer-events-none absolute left-2 top-2 rounded-full bg-black/60 px-2 py-0.5 text-[11px] font-medium text-white">
                Montage
              </span>
            )}
            <button
              type="button"
              onClick={() => setToDelete(p)}
              className="absolute right-2 top-2 rounded-full bg-black/60 p-1.5 text-white opacity-0 transition-opacity hover:bg-red-600 focus:opacity-100 group-hover:opacity-100"
              aria-label={`Supprimer ${p.name}`}
            >
              <Trash2 className="h-4 w-4" />
            </button>
          </li>
        ))}
      </ul>

      {zoomed && (
        <div
          className="fixed inset-0 z-50 flex items-center justify-center bg-black/85 p-4"
          role="dialog"
          aria-modal="true"
          aria-label={zoomed.name}
          onClick={() => setZoomed(null)}
          onKeyDown={(e) => e.key === "Escape" && setZoomed(null)}
        >
          <AuthImage
            path={photoboothReviewApi.picturePath(zoomed.id, "preview")}
            alt={zoomed.name}
            className="h-full max-h-[88vh] w-full max-w-5xl"
            fit="contain"
          />
          <div className="absolute right-4 top-4 flex gap-2">
            <button
              type="button"
              onClick={(e) => {
                e.stopPropagation();
                setToDelete(zoomed);
              }}
              className="flex items-center gap-1.5 rounded-full bg-white/10 px-3 py-2 text-sm text-white hover:bg-red-600"
            >
              <Trash2 className="h-4 w-4" /> Supprimer
            </button>
            <button
              type="button"
              onClick={() => setZoomed(null)}
              className="rounded-full bg-white/10 p-2 text-white hover:bg-white/20"
              aria-label="Fermer"
            >
              <X className="h-5 w-5" />
            </button>
          </div>
        </div>
      )}

      <Modal
        isOpen={toDelete !== null}
        onClose={() => setToDelete(null)}
        title="Supprimer cette photo ?"
      >
        <p className="text-sm text-gray-600">
          Elle est retirée du Drive et de la page de partage des invités. Cette action est
          définitive.
        </p>
        <div className="mt-6 flex justify-end gap-2">
          <button
            type="button"
            onClick={() => setToDelete(null)}
            className="rounded-xl border border-gray-300 bg-white px-4 py-2 text-sm text-gray-700 hover:bg-gray-50"
          >
            Annuler
          </button>
          <button
            type="button"
            disabled={remove.isPending}
            onClick={() => toDelete && doDelete(toDelete)}
            className="flex items-center gap-2 rounded-xl bg-red-600 px-4 py-2 text-sm font-semibold text-white hover:bg-red-700 disabled:opacity-50"
          >
            {remove.isPending && <Loader2 className="h-4 w-4 animate-spin" />}
            Supprimer
          </button>
        </div>
      </Modal>

      <Modal
        isOpen={confirmApprove}
        onClose={() => setConfirmApprove(false)}
        title="Publier les photos dans l'application ?"
      >
        <p className="text-sm text-gray-600">
          Les {pictures.length} photo{pictures.length > 1 ? "s" : ""} de « {folder.event_name} »
          deviendront visibles par tous les utilisateurs de l'application.
        </p>
        {going && (
          <p className="mt-3 flex items-start gap-2 text-sm text-amber-700">
            <AlertTriangle className="mt-0.5 h-4 w-4 shrink-0" />
            La soirée semble encore en cours.
          </p>
        )}
        <div className="mt-6 flex justify-end gap-2">
          <button
            type="button"
            onClick={() => setConfirmApprove(false)}
            className="rounded-xl border border-gray-300 bg-white px-4 py-2 text-sm text-gray-700 hover:bg-gray-50"
          >
            Annuler
          </button>
          <button
            type="button"
            disabled={approve.isPending}
            onClick={doApprove}
            className="flex items-center gap-2 rounded-xl bg-blue-600 px-4 py-2 text-sm font-semibold text-white hover:bg-blue-700 disabled:opacity-50"
          >
            {approve.isPending && <Loader2 className="h-4 w-4 animate-spin" />}
            Approuver et publier
          </button>
        </div>
      </Modal>
    </div>
  );
}

/** An image served by the API behind the panel's token, loaded once it scrolls into view. */
function AuthImage({
  path,
  alt,
  className,
  fit,
}: {
  path: string;
  alt: string;
  className: string;
  fit: "cover" | "contain";
}) {
  const ref = useRef<HTMLDivElement>(null);
  const [visible, setVisible] = useState(false);
  const [src, setSrc] = useState<string | null>(null);
  const [failed, setFailed] = useState(false);

  useEffect(() => {
    const el = ref.current;
    if (!el) return;
    const io = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          setVisible(true);
          io.disconnect();
        }
      },
      { rootMargin: "300px" },
    );
    io.observe(el);
    return () => io.disconnect();
  }, []);

  useEffect(() => {
    if (!visible) return;
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
    };
  }, [path, visible]);

  return (
    <div
      ref={ref}
      className={`relative flex items-center justify-center overflow-hidden ${fit === "cover" ? "bg-gray-100" : ""} ${className}`}
    >
      {failed ? (
        <div className="flex flex-col items-center gap-2 text-gray-400">
          <ImageOff className="h-7 w-7" />
          <span className="text-xs">Aperçu indisponible</span>
        </div>
      ) : src ? (
        // biome-ignore lint/performance/noImgElement: blob URL from an authenticated request
        <img
          src={src}
          alt={alt}
          className={`h-full w-full ${fit === "cover" ? "object-cover" : "object-contain"}`}
        />
      ) : (
        <Loader2 className="h-6 w-6 animate-spin text-gray-400" />
      )}
    </div>
  );
}
