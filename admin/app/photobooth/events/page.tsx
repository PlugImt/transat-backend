import { redirect } from "next/navigation";

// The events now live in the Photobooth page.
export default function PhotoboothEventsRedirect() {
  redirect("/photobooth");
}
