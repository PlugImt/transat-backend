"use client";

import { format } from "date-fns";
import { fr } from "date-fns/locale";
import { CheckCircle2, ExternalLink, Loader2 } from "lucide-react";
import Modal from "@/components/Modal";
import { AuthImage, FolderReview } from "@/components/photobooth/Review";
import { photoboothReviewApi } from "@/lib/api";
import { usePhotoboothAlbumPictures, usePhotoboothPending } from "@/lib/hooks";
import type { PhotoboothAlbumFolder } from "@/lib/types";

const eventDay = (date: string) =>
  format(new Date(`${date}T12:00:00`), "EEEE d MMMM yyyy", { locale: fr });

/**
 * The photobooth's pictures of a gallery album: those awaiting approval (reviewed and approved right here) and
 * those already published in the album.
 */
export function AlbumPhotoboothModal({
  albumId,
  title,
  onClose,
}: {
  albumId: number | null;
  title: string;
  onClose: () => void;
}) {
  const { data: folders = [], isLoading, isError } = usePhotoboothAlbumPictures(albumId);
  const { data: pending = [] } = usePhotoboothPending();

  return (
    <Modal
      isOpen={albumId !== null}
      onClose={onClose}
      title={`Photobooth · ${title}`}
      maxWidth="4xl"
    >
      {isLoading && (
        <div className="flex justify-center py-12 text-gray-400">
          <Loader2 className="h-6 w-6 animate-spin" />
        </div>
      )}
      {isError && (
        <p className="rounded-2xl border border-red-200 bg-red-50 p-6 text-center text-sm text-red-700">
          Impossible de charger les photos du photobooth.
        </p>
      )}
      {!isLoading && !isError && folders.length === 0 && (
        <p className="py-10 text-center text-sm text-gray-500">
          Aucune photo du photobooth dans cet album.
        </p>
      )}
      <div className="space-y-6">
        {folders.map((f) => {
          const review = f.pending ? pending.find((p) => p.folder_id === f.folder_id) : undefined;
          return review ? (
            <div
              key={f.folder_id}
              className="rounded-2xl border border-amber-200 bg-amber-50/40 p-4 sm:p-5"
            >
              <FolderReview folder={review} embedded />
            </div>
          ) : (
            <ApprovedFolder key={f.folder_id} folder={f} />
          );
        })}
      </div>
    </Modal>
  );
}

function ApprovedFolder({ folder: f }: { folder: PhotoboothAlbumFolder }) {
  return (
    <section>
      <h3 className="flex flex-wrap items-center gap-2 text-sm font-semibold text-gray-900">
        <CheckCircle2 className="h-4 w-4 text-emerald-600" />
        {f.pending ? "En attente d'approbation" : "Publiées dans l'album"} · {f.pictures.length}{" "}
        photo
        {f.pictures.length > 1 ? "s" : ""}
        <span className="font-normal capitalize text-gray-500">· {eventDay(f.event_date)}</span>
        <a
          href={f.drive_url}
          target="_blank"
          rel="noreferrer"
          className="inline-flex items-center gap-1 font-normal text-blue-600 hover:underline"
        >
          Drive <ExternalLink className="h-3.5 w-3.5" />
        </a>
      </h3>
      <ul className="mt-3 grid grid-cols-3 gap-2 sm:grid-cols-5">
        {f.pictures.map((p) => (
          <li key={p.id}>
            <AuthImage
              path={photoboothReviewApi.picturePath(p.id, "thumb")}
              alt={p.name}
              className="aspect-square w-full rounded-xl"
              fit="cover"
            />
          </li>
        ))}
      </ul>
    </section>
  );
}
