package photobooth

import (
	"context"
	"testing"
)

func pendingFolderOf(t *testing.T, svc *Service) PendingFolder {
	t.Helper()
	list, err := svc.ListPending(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Fatalf("%d pending folders, want 1: %+v", len(list), list)
	}
	return list[0]
}

func TestApprovalPublishesThePendingFolder(t *testing.T) {
	svc, d := setup(t)
	ctx := context.Background()
	b := booth{t, svc, d, "session-0001"}
	b.put(KindPhoto, 1, "p1")
	b.put(KindPhoto, 2, "p2")

	p := pendingFolderOf(t, svc)
	if p.EventName != event.Name || p.EventDate != event.Date || p.Pictures != 2 || p.OpenShoots != 1 {
		t.Fatalf("pending folder %+v", p)
	}
	pics, err := svc.PendingPictures(ctx, p.FolderID)
	if err != nil || len(pics) != 2 {
		t.Fatalf("pictures %+v (%v)", pics, err)
	}

	// The reviewer removes one picture, then approves the rest.
	if err := svc.DeletePendingPicture(ctx, pics[0].ID, "admin@test"); err != nil {
		t.Fatal(err)
	}
	waitLive(t, d, pics[1].Name)
	res, err := svc.Approve(ctx, p.FolderID, "admin@test")
	if err != nil {
		t.Fatal(err)
	}
	if res.Merged || res.FolderID != p.FolderID {
		t.Fatalf("approve %+v", res)
	}
	if info, _ := d.GetFile(ctx, p.FolderID); info.Name != PhotoboothFolderName {
		t.Fatalf("approved folder named %q", info.Name)
	}
	if list, _ := svc.ListPending(ctx); len(list) != 0 {
		t.Fatalf("still pending after approval: %+v", list)
	}
	if _, err := svc.Approve(ctx, p.FolderID, "admin@test"); status(err) != 409 {
		t.Fatalf("second approval: %v", err)
	}
	if err := svc.DeletePendingPicture(ctx, pics[1].ID, "admin@test"); status(err) != 409 {
		t.Fatalf("deleting a published picture through the review: %v", err)
	}

	// The party goes on: new pictures wait in a new pending folder, then merge into the published one.
	later := booth{t, svc, d, "session-0002"}
	later.put(KindPhoto, 1, "late")
	p2 := pendingFolderOf(t, svc)
	if p2.FolderID == p.FolderID {
		t.Fatal("new pictures went into the published folder")
	}
	res, err = svc.Approve(ctx, p2.FolderID, "admin@test")
	if err != nil {
		t.Fatal(err)
	}
	if !res.Merged || res.FolderID != p.FolderID || res.Moved != 1 {
		t.Fatalf("merge %+v", res)
	}
	if info, _ := d.GetFile(ctx, p2.FolderID); !info.Trashed {
		t.Fatal("the emptied pending folder was kept")
	}
	kids, _ := d.ListChildren(ctx, p.FolderID)
	if len(kids) != 2 {
		t.Fatalf("published folder holds %d pictures, want 2", len(kids))
	}
}

func TestFolderRenamedByHandIsNoLongerPending(t *testing.T) {
	svc, d := setup(t)
	ctx := context.Background()
	b := booth{t, svc, d, "session-0001"}
	b.put(KindPhoto, 1, "p1")
	p := pendingFolderOf(t, svc)
	if err := d.Rename(ctx, p.FolderID, "Photobooth"); err != nil {
		t.Fatal(err)
	}
	if list, _ := svc.ListPending(ctx); len(list) != 0 {
		t.Fatalf("a folder renamed on Drive is still listed: %+v", list)
	}
	if _, err := svc.Approve(ctx, p.FolderID, "admin@test"); status(err) != 409 {
		t.Fatalf("approving a folder already renamed: %v", err)
	}
}
