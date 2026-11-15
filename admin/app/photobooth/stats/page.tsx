import { redirect } from "next/navigation";

// The statistics now live in the Photobooth page.
export default function PhotoboothStatsRedirect() {
  redirect("/photobooth?tab=stats");
}
