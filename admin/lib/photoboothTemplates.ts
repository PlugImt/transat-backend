import type { PhotoboothAspect, PhotoboothLayout, PhotoboothZone } from "./types";

/*
 * Templates of the photobooth's layouts, drawn from the geometry the booths publish. Same drawing as the frame
 * upload page of transat.dev (/photobooth/frames/:token), so organisers get the same files from both.
 */

/** 2 mm at 300 dpi: what the printer may trim off each edge. */
const BLEED_PX = 24;
const MAX_BYTES = 14 * 1024 * 1024;
const MAX_PX = 6000;
const MAX_ASPECT_DRIFT = 0.02;

export const ASPECT_LABELS: Record<PhotoboothAspect, string> = {
  RATIO_4_3: "4:3",
  RATIO_3_2: "3:2",
  RATIO_16_9: "16:9",
};

/** The vector templates of every layout, hosted with the frame upload page. */
export const templatesPdfUrl = (aspect: PhotoboothAspect) =>
  `https://transat.dev/photobooth-templates/gabarits-photos-${ASPECT_LABELS[aspect].replace(":", "-")}.pdf`;

export const isStrips = (l: PhotoboothLayout) => l.type.startsWith("STRIP_");

export const shotCount = (l: PhotoboothLayout) => new Set(l.photos.map((p) => p.shot ?? 1)).size;

const fileStem = (event: string, l: PhotoboothLayout) =>
  `${event} - ${l.title}`
    .normalize("NFD")
    .replace(/[̀-ͯ]/g, "")
    .replace(/[^A-Za-z0-9 _-]+/g, "")
    .trim()
    .replace(/\s+/g, "_");

/**
 * The template of one layout as an SVG at the canvas' pixel size: the photo zones (printed over, so nothing drawn
 * there shows), the zone kept free for the title, the trim margin and, for strips, the cut line.
 */
export function templateSvg(l: PhotoboothLayout, labels = true): string {
  const font = Math.round(Math.min(l.width, l.height) / 26);
  const zone = (z: PhotoboothZone) => {
    const label = labels
      ? `<text x="${z.x + z.w / 2}" y="${z.y + z.h / 2}" font-size="${font}" font-family="Poppins, Arial, sans-serif" font-weight="700" fill="#be185d" text-anchor="middle" dominant-baseline="middle">Photo ${z.shot ?? 1}</text>
         <text x="${z.x + z.w / 2}" y="${z.y + z.h / 2 + font * 1.2}" font-size="${font * 0.55}" font-family="Poppins, Arial, sans-serif" fill="#be185d" text-anchor="middle" dominant-baseline="middle">${z.w} × ${z.h} px · x ${z.x}, y ${z.y}</text>`
      : "";
    return `<rect x="${z.x}" y="${z.y}" width="${z.w}" height="${z.h}" fill="#fbcfe8" stroke="#db2777" stroke-width="4"/>${label}`;
  };
  const t = l.title_zone;
  const title = t
    ? `<rect x="${t.x + 8}" y="${t.y + 8}" width="${t.w - 16}" height="${t.h - 16}" fill="none" stroke="#2563eb" stroke-width="4" stroke-dasharray="18 12"/>
       ${labels ? `<text x="${t.x + t.w / 2}" y="${t.y + t.h / 2}" font-size="${font * 0.8}" font-family="Poppins, Arial, sans-serif" font-weight="600" fill="#2563eb" text-anchor="middle" dominant-baseline="middle">Zone titre / logo</text>` : ""}`
    : "";
  const cut = isStrips(l)
    ? `<line x1="${l.width / 2}" y1="0" x2="${l.width / 2}" y2="${l.height}" stroke="#111" stroke-width="3" stroke-dasharray="14 10"/>`
    : "";
  const bleed = `<rect x="${BLEED_PX}" y="${BLEED_PX}" width="${l.width - 2 * BLEED_PX}" height="${l.height - 2 * BLEED_PX}" fill="none" stroke="#dc2626" stroke-width="3" stroke-dasharray="10 8"/>`;
  return `<svg xmlns="http://www.w3.org/2000/svg" width="${l.width}" height="${l.height}" viewBox="0 0 ${l.width} ${l.height}"><rect width="${l.width}" height="${l.height}" fill="#ffffff"/>${title}${l.photos.map(zone).join("")}${cut}${bleed}</svg>`;
}

const toBlob = (canvas: HTMLCanvasElement) =>
  new Promise<Blob>((resolve, reject) =>
    canvas.toBlob((b) => (b ? resolve(b) : reject(new Error("png"))), "image/png"),
  );

async function templatePng(l: PhotoboothLayout): Promise<Blob> {
  const url = URL.createObjectURL(new Blob([templateSvg(l)], { type: "image/svg+xml" }));
  try {
    const img = new Image();
    img.src = url;
    await img.decode();
    const canvas = document.createElement("canvas");
    canvas.width = l.width;
    canvas.height = l.height;
    canvas.getContext("2d")?.drawImage(img, 0, 0, l.width, l.height);
    return await toBlob(canvas);
  } finally {
    URL.revokeObjectURL(url);
  }
}

/** A white layer with transparent holes where the photos go: the starting point for a design in any editor. */
async function maskPng(l: PhotoboothLayout): Promise<Blob> {
  const canvas = document.createElement("canvas");
  canvas.width = l.width;
  canvas.height = l.height;
  const ctx = canvas.getContext("2d");
  if (ctx) {
    ctx.fillStyle = "#ffffff";
    ctx.fillRect(0, 0, l.width, l.height);
    for (const z of l.photos) ctx.clearRect(z.x, z.y, z.w, z.h);
  }
  return toBlob(canvas);
}

function saveBlob(blob: Blob, name: string) {
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = name;
  document.body.appendChild(a);
  a.click();
  a.remove();
  setTimeout(() => URL.revokeObjectURL(url), 10_000);
}

export const downloadTemplate = async (event: string, l: PhotoboothLayout) =>
  saveBlob(await templatePng(l), `${fileStem(event, l)}_gabarit.png`);

export const downloadMask = async (event: string, l: PhotoboothLayout) =>
  saveBlob(await maskPng(l), `${fileStem(event, l)}_masque.png`);

/** The server's checks, done before sending so the message comes at once. Null when the file is fine. */
export async function checkFrame(file: File, l: PhotoboothLayout): Promise<string | null> {
  if (file.type !== "image/png" && !file.name.toLowerCase().endsWith(".png")) {
    return "Le fond doit être un PNG.";
  }
  if (file.size > MAX_BYTES)
    return "L'image dépasse 14 Mo : exportez-la avec une compression PNG plus forte.";
  let w: number;
  let h: number;
  try {
    const bmp = await createImageBitmap(file);
    w = bmp.width;
    h = bmp.height;
    bmp.close();
  } catch {
    return "Ce fichier n'est pas une image PNG lisible.";
  }
  const want = l.width / l.height;
  if (Math.abs(w / h - want) / want > MAX_ASPECT_DRIFT) {
    return `Ce collage attend une image de ${l.width} × ${l.height} px (ou les mêmes proportions), celle-ci fait ${w} × ${h} px.`;
  }
  if (w < l.width / 2)
    return `L'image est trop petite (${w} × ${h} px) : exportez-la en ${l.width} × ${l.height} px.`;
  if (w > MAX_PX || h > MAX_PX)
    return `L'image est trop grande : exportez-la en ${l.width} × ${l.height} px.`;
  return null;
}
