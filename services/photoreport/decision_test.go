package photoreport

import (
	"context"
	"errors"
	"testing"
)

type fakeRemover struct {
	trashed []string
	fail    error
}

func (r *fakeRemover) Trash(_ context.Context, id string) error {
	if r.fail != nil {
		return r.fail
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
