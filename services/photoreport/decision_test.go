package photoreport

import (
	"context"
	"errors"
	"testing"
	"time"
)

type fakeRemover struct {
	trashed []string
	fail    error
	refuse  map[string]bool // photos Drive refuses to trash
}

func (r *fakeRemover) Trash(_ context.Context, id string) error {
	if r.fail != nil {
		return r.fail
	}
	if r.refuse[id] {
		return errors.New("drive refused")
	}
	r.trashed = append(r.trashed, id)
	return nil
}

func (f *fx) deleted(photo string) bool {
	var gone bool
	if err := f.db.QueryRow(`SELECT deleted_at IS NOT NULL FROM photos WHERE drive_file_id = $1`, photo).Scan(&gone); err != nil {
		f.t.Fatal(err)
	}
	return gone
}

func TestApprovingDeletesThePhotoAndClosesItsOtherReports(t *testing.T) {
	f, _ := standard(t)
	ctx := context.Background()
	f.photo("PHOTO-0000000010", nil)
	first, err := f.submit("reporter@x.test", "PHOTO-0000000010", ReasonAppearsInPhoto, "")
	if err != nil {
		t.Fatal(err)
	}
	second, err := f.submit("reporter@x.test", "PHOTO-0000000010", ReasonInappropriate, "")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := f.svc.Decide(ctx, first.ID, DecisionApprove, "", "admin1@x.test"); !errors.Is(err, ErrRemovalUnavailable) {
		t.Fatalf("approve without a Drive writer: %v", err)
	}
	remover := &fakeRemover{fail: errors.New("drive down")}
	f.svc.SetRemover(remover)
	if _, err := f.svc.Decide(ctx, first.ID, DecisionApprove, "", "admin1@x.test"); err == nil {
		t.Fatal("a Drive failure must not record the decision")
	}
	if r, _ := f.svc.Get(ctx, first.ID); r.Status != "pending" || f.deleted("PHOTO-0000000010") {
		t.Fatalf("after a Drive failure: status %s", r.Status)
	}

	remover.fail = nil
	r, err := f.svc.Decide(ctx, first.ID, DecisionApprove, "  visible face  ", "owner@x.test")
	if err != nil {
		t.Fatal(err)
	}
	if r.Status != "approved" || r.DecidedBy == nil || *r.DecidedBy != "owner@x.test" || r.DecisionNote != "visible face" || r.DecidedAt == nil {
		t.Fatalf("decision not recorded: %+v", r)
	}
	if len(remover.trashed) != 1 || remover.trashed[0] != "PHOTO-0000000010" || !f.deleted("PHOTO-0000000010") {
		t.Fatalf("photo not removed: trashed %v", remover.trashed)
	}
	if other, _ := f.svc.Get(ctx, second.ID); other.Status != "approved" || other.PhotoAvailable {
		t.Fatalf("the other report of the photo: %+v", other)
	}
	if _, err := f.svc.Decide(ctx, first.ID, DecisionReject, "", "admin1@x.test"); !errors.Is(err, ErrAlreadyDecided) {
		t.Fatalf("second decision: %v", err)
	}
}

func TestRejectingKeepsThePhoto(t *testing.T) {
	f, _ := standard(t)
	ctx := context.Background()
	f.photo("PHOTO-0000000011", nil)
	rec, err := f.submit("reporter@x.test", "PHOTO-0000000011", ReasonPoorQuality, "")
	if err != nil {
		t.Fatal(err)
	}
	remover := &fakeRemover{}
	f.svc.SetRemover(remover)
	if _, err := f.svc.Decide(ctx, rec.ID, "delete", "", "admin1@x.test"); !errors.Is(err, ErrInvalidDecision) {
		t.Fatalf("unknown decision: %v", err)
	}
	r, err := f.svc.Decide(ctx, rec.ID, DecisionReject, "fine as is", "admin1@x.test")
	if err != nil {
		t.Fatal(err)
	}
	if r.Status != "rejected" || len(remover.trashed) != 0 || f.deleted("PHOTO-0000000011") {
		t.Fatalf("reject touched the photo: %+v trashed %v", r, remover.trashed)
	}
	counts, err := f.svc.Counts(ctx)
	if err != nil || counts["rejected"] != 1 || counts["pending"] != 0 {
		t.Fatalf("counts %v %v", counts, err)
	}
	list, total, err := f.svc.List(ctx, ListFilter{Status: "pending"})
	if err != nil || total != 0 || len(list) != 0 {
		t.Fatalf("pending list: %d %v", total, err)
	}
}

func TestApproveAllApprovesThePendingReportsInBatches(t *testing.T) {
	f, _ := standard(t)
	ctx := context.Background()
	f.user("reporter2@x.test", "Rob", "NEWF")
	for _, p := range []string{"PHOTO-0000000020", "PHOTO-0000000021", "PHOTO-0000000022", "PHOTO-0000000023", "PHOTO-0000000024"} {
		f.photo(p, nil)
	}
	submit := func(email, photo string, reason Reason) int64 {
		t.Helper()
		r, err := f.submit(email, photo, reason, "")
		if err != nil {
			t.Fatal(err)
		}
		return r.ID
	}
	a1 := submit("reporter@x.test", "PHOTO-0000000020", ReasonAppearsInPhoto)
	a2 := submit("reporter2@x.test", "PHOTO-0000000020", ReasonInappropriate) // same photo, another reason
	b := submit("reporter@x.test", "PHOTO-0000000021", ReasonAppearsInPhoto)
	refused := submit("reporter@x.test", "PHOTO-0000000022", ReasonAppearsInPhoto)
	other := submit("reporter2@x.test", "PHOTO-0000000023", ReasonPoorQuality)
	rejected := submit("reporter2@x.test", "PHOTO-0000000024", ReasonAppearsInPhoto)
	// Rejecting needs no Drive writer.
	if _, err := f.svc.Decide(ctx, rejected, DecisionReject, "", "admin1@x.test"); err != nil {
		t.Fatal(err)
	}

	if _, err := f.svc.ApproveAll(ctx, "", "admin1@x.test", nil, time.Minute); !errors.Is(err, ErrRemovalUnavailable) {
		t.Fatalf("without a Drive writer: %v", err)
	}
	if _, err := f.svc.ApproveAll(ctx, "nope", "admin1@x.test", nil, time.Minute); !errors.Is(err, ErrInvalidReason) {
		t.Fatalf("unknown reason: %v", err)
	}
	remover := &fakeRemover{refuse: map[string]bool{"PHOTO-0000000022": true}}
	f.svc.SetRemover(remover)
	// No time left: nothing is done, everything is still to do.
	res, err := f.svc.ApproveAll(ctx, "appears_in_photo", "admin1@x.test", nil, 0)
	if err != nil || res.Photos != 0 || res.Remaining != 3 {
		t.Fatalf("zero budget: %+v %v", res, err)
	}

	// One reason only: the poor-quality report stays; the other report of the same photo closes with it.
	res, err = f.svc.ApproveAll(ctx, "appears_in_photo", "owner@x.test", nil, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if res.Photos != 2 || res.Approved != 2 || len(res.Failed) != 1 || res.Failed[0] != refused || res.Remaining != 0 {
		t.Fatalf("approve all appears_in_photo: %+v", res)
	}
	for id, want := range map[int64]string{a1: "approved", a2: "approved", b: "approved", refused: "pending", other: "pending", rejected: "rejected"} {
		r, _ := f.svc.Get(ctx, id)
		if r.Status != want {
			t.Errorf("report %d: %s, want %s", id, r.Status, want)
		}
	}
	if r, _ := f.svc.Get(ctx, a1); r.DecidedBy == nil || *r.DecidedBy != "owner@x.test" {
		t.Fatalf("decided by: %+v", r.DecidedBy)
	}
	if !f.deleted("PHOTO-0000000020") || !f.deleted("PHOTO-0000000021") || f.deleted("PHOTO-0000000022") || f.deleted("PHOTO-0000000024") {
		t.Fatal("wrong photos deleted")
	}

	// Everything else, skipping the report Drive refused: only the poor-quality one is left to do.
	res, err = f.svc.ApproveAll(ctx, "", "admin1@x.test", []int64{refused}, time.Minute)
	if err != nil || res.Photos != 1 || res.Remaining != 0 || len(res.Failed) != 0 {
		t.Fatalf("second call: %+v %v", res, err)
	}
	if r, _ := f.svc.Get(ctx, refused); r.Status != "pending" {
		t.Fatal("a skipped report was decided")
	}
}
