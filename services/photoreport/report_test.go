package photoreport

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"testing"
	"unicode/utf8"

	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"html/template"

	"bytes"

	"github.com/plugimt/transat-backend/i18n"
	"github.com/plugimt/transat-backend/internal/database/dbtest"
	"github.com/plugimt/transat-backend/models"
)

type sentMail struct {
	to       string
	language string
	subject  string
	data     EmailData
}

// fakeMailer records every message and fails the recipients it is told to.
type fakeMailer struct {
	mu      sync.Mutex
	sent    []sentMail
	calls   map[string]int
	failFor map[string]bool
	panics  map[string]bool
}

func newMailer() *fakeMailer {
	return &fakeMailer{calls: map[string]int{}, failFor: map[string]bool{}, panics: map[string]bool{}}
}

func (m *fakeMailer) SendEmail(mail models.Email, data interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.calls[mail.Recipient]++
	if m.panics[mail.Recipient] {
		panic("boom")
	}
	if m.failFor[mail.Recipient] {
		return errors.New("provider unavailable")
	}
	m.sent = append(m.sent, sentMail{to: mail.Recipient, language: mail.Language, subject: mail.SubjectKey, data: data.(EmailData)})
	return nil
}

func (m *fakeMailer) recipients() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	var out []string
	for _, s := range m.sent {
		out = append(out, s.to)
	}
	sort.Strings(out)
	return out
}

func (m *fakeMailer) set(who string, fail bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.failFor[who] = fail
}

type fx struct {
	t   *testing.T
	db  *sql.DB
	svc *Service
	m   *fakeMailer
}

func newFx(t *testing.T) *fx {
	db := dbtest.Open(t)
	m := newMailer()
	return &fx{t: t, db: db, m: m, svc: NewService(db, m, "https://admin.example.test/")}
}

func (f *fx) exec(q string, args ...any) {
	f.t.Helper()
	if _, err := f.db.Exec(q, args...); err != nil {
		f.t.Fatalf("%s: %v", q, err)
	}
}

func (f *fx) user(email, first string, roles ...string) {
	f.t.Helper()
	f.exec(`INSERT INTO newf (email, password, first_name, last_name, graduation_year, formation_name)
		VALUES ($1, 'x', $2, 'Doe', 2027, 'FISE')`, email, first)
	for _, r := range roles {
		f.exec(`INSERT INTO roles (name) VALUES ($1) ON CONFLICT (name) DO NOTHING`, r)
		f.exec(`INSERT INTO newf_roles (email, id_roles) SELECT $1, id_roles FROM roles WHERE name = $2`, email, r)
	}
}

func (f *fx) club(name string) int64 {
	f.t.Helper()
	var id int64
	if err := f.db.QueryRow(`INSERT INTO clubs (name, picture) VALUES ($1, 'p') RETURNING id_clubs`, name).Scan(&id); err != nil {
		f.t.Fatal(err)
	}
	return id
}

func (f *fx) album(title string, club, parent *int64) int64 {
	f.t.Helper()
	var id int64
	err := f.db.QueryRow(`INSERT INTO photo_albums (root_folder_id, title, raw_folder_name, drive_url, visibility, confidence,
			rule, last_seen_run, club_id, parent_album_id)
		VALUES ($1, $2, $2, $3, 'published', 1, 'R', 1, $4, $5) RETURNING id_photo_albums`,
		"root-"+title, title, "https://drive.example/album-"+title, club, parent).Scan(&id)
	if err != nil {
		f.t.Fatal(err)
	}
	return id
}

func (f *fx) photo(id string, album *int64) {
	f.t.Helper()
	f.exec(`INSERT INTO photos (drive_file_id, album_id, name, mime_type, media_type, sort_at, sort_at_source, web_url, last_seen_run)
		VALUES ($1, $2, $3, 'image/jpeg', 'image', now(), 'x', $4, 1)`, id, album, "IMG_"+id+".jpg", "https://drive.example/file-"+id)
}

func (f *fx) submit(email, photo string, reason Reason, comment string) (Receipt, error) {
	return f.svc.Submit(context.Background(), Submission{PhotoID: photo, Reason: reason, Comment: comment, ReporterEmail: email})
}

func (f *fx) deliver() int {
	f.t.Helper()
	n, err := f.svc.DeliverDue(context.Background())
	if err != nil {
		f.t.Fatal(err)
	}
	return n
}

func (f *fx) due() {
	f.exec(`UPDATE photo_report_deliveries SET next_attempt_at = now() - interval '1 second'`)
}

const photoID = "PHOTO-0000000001"

// standard: a club with one owner, two admins, a reporter and one photo of the club's album.
func standard(t *testing.T) (*fx, int64) {
	f := newFx(t)
	f.user("owner@x.test", "Olive", "clubphoto_respo", "NEWF")
	f.user("admin1@x.test", "Ada", "ADMIN")
	f.user("admin2@x.test", "Ben", "ADMIN")
	f.user("reporter@x.test", "Rita", "NEWF")
	club := f.club("Club Photo")
	album := f.album("Gala", &club, nil)
	f.photo(photoID, &album)
	return f, club
}

func TestEveryReasonReachesTheClubOwnerAndAllAdmins(t *testing.T) {
	for _, reason := range []Reason{ReasonAppearsInPhoto, ReasonInappropriate, ReasonPoorQuality, ReasonOther} {
		t.Run(string(reason), func(t *testing.T) {
			f, club := standard(t)
			rec, err := f.submit("reporter@x.test", photoID, reason, "Merci de vérifier")
			if err != nil || rec.Duplicate {
				t.Fatalf("submit: %v %+v", err, rec)
			}
			if n := f.deliver(); n != 3 {
				t.Fatalf("delivered %d, want 3", n)
			}
			got := fmt.Sprint(f.m.recipients())
			if want := "[admin1@x.test admin2@x.test owner@x.test]"; got != want {
				t.Fatalf("recipients %s, want %s", got, want)
			}
			for _, s := range f.m.sent {
				d := s.data
				if d.ReasonCode != string(reason) || d.Comment != "Merci de vérifier" || d.PhotoID != photoID ||
					d.AlbumTitle != "Gala" || d.ClubName != "Club Photo" || d.ClubID != club || !d.HasClub {
					t.Fatalf("unexpected payload %+v", d)
				}
				if d.ReporterMail != "reporter@x.test" || d.ReporterName != "Rita Doe" || d.ReporterYear != 2027 || d.Formation != "FISE" {
					t.Fatalf("reporter data %+v", d)
				}
				if d.SubjectTitle != "Gala" || d.SubmittedAt == "" || d.AdminURL != fmt.Sprintf("https://admin.example.test/photo-reports?report=%d", rec.ID) {
					t.Fatalf("references %+v", d)
				}
				if d.IsClubOwner != (s.to == "owner@x.test") {
					t.Fatalf("%s got the wrong role: %+v", s.to, d)
				}
			}
			r, err := f.svc.Get(context.Background(), rec.ID)
			if err != nil || r.DeliveryStatus != "sent" || r.Recipients != 3 || r.RecipientsSent != 3 {
				t.Fatalf("status %+v %v", r, err)
			}
		})
	}
}

func TestPhotosWithoutAClubBelongToTheDefaultClub(t *testing.T) {
	f, defaultClub := standard(t)
	album := f.album("Sans club", nil, nil)
	f.photo("PHOTO-0000000002", &album)
	f.photo("PHOTO-0000000003", nil) // not even in an album

	for _, id := range []string{"PHOTO-0000000002", "PHOTO-0000000003"} {
		rec, err := f.submit("reporter@x.test", id, ReasonOther, "details")
		if err != nil {
			t.Fatal(err)
		}
		if r, _ := f.svc.Get(context.Background(), rec.ID); r.ClubName != "Club Photo" || r.ClubID == nil || *r.ClubID != defaultClub {
			t.Fatalf("%s: club %+v, want the default club", id, r)
		}
	}
	f.deliver()
	if got := fmt.Sprint(f.m.recipients()); got != "[admin1@x.test admin1@x.test admin2@x.test admin2@x.test owner@x.test owner@x.test]" {
		t.Fatalf("recipients %s, want the photo club's owner and the administrators", got)
	}
}

func TestAnAlbumsOwnClubBeatsTheDefault(t *testing.T) {
	f, _ := standard(t)
	f.user("sport@x.test", "Sam", "sport_respo")
	sport := f.club("Sport")
	album := f.album("Tournoi", &sport, nil)
	f.photo("PHOTO-0000000002", &album)

	if _, err := f.submit("reporter@x.test", "PHOTO-0000000002", ReasonInappropriate, ""); err != nil {
		t.Fatal(err)
	}
	f.deliver()
	if got := fmt.Sprint(f.m.recipients()); got != "[admin1@x.test admin2@x.test sport@x.test]" {
		t.Fatalf("recipients %s: the default club's owner must not be told", got)
	}
}

func TestWithoutADefaultClubPhotosGoToAdministratorsOnly(t *testing.T) {
	f, _ := standard(t)
	f.exec(`UPDATE clubs SET name = 'Autre club'`) // the photo club does not exist
	f.photo("PHOTO-0000000003", nil)
	if _, err := f.submit("reporter@x.test", "PHOTO-0000000003", ReasonOther, "details"); err != nil {
		t.Fatal(err)
	}
	f.deliver()
	if got := fmt.Sprint(f.m.recipients()); got != "[admin1@x.test admin2@x.test]" {
		t.Fatalf("recipients %s, want administrators only", got)
	}
	for _, s := range f.m.sent {
		if s.data.HasClub || s.data.IsClubOwner {
			t.Fatalf("no club expected: %+v", s.data)
		}
	}

	f.exec(`UPDATE clubs SET name = 'club PHOTO'`) // matched without regard to case
	f.photo("PHOTO-0000000004", nil)
	if _, err := f.submit("reporter@x.test", "PHOTO-0000000004", ReasonOther, "details"); err != nil {
		t.Fatal(err)
	}
	f.deliver()
	if f.m.calls["owner@x.test"] != 1 {
		t.Fatalf("owner mailed %d times, want 1", f.m.calls["owner@x.test"])
	}
}

func TestClubWithoutOwnerAndSubAlbumInheritance(t *testing.T) {
	f := newFx(t)
	f.user("admin@x.test", "Ada", "ADMIN")
	f.user("reporter@x.test", "Rita", "NEWF")
	f.user("other@x.test", "Otto", "otherclub_respo") // owner of another club: must never be told
	club := f.club("Club Photo")
	parent := f.album("Gala", &club, nil)
	child := f.album("Gala - scène", nil, &parent)
	f.photo(photoID, &child)

	rec, err := f.submit("reporter@x.test", photoID, ReasonInappropriate, "")
	if err != nil {
		t.Fatal(err)
	}
	f.deliver()
	if got := fmt.Sprint(f.m.recipients()); got != "[admin@x.test]" {
		t.Fatalf("recipients %s: a club without owner falls back to administrators", got)
	}
	if r, _ := f.svc.Get(context.Background(), rec.ID); r.ClubName != "Club Photo" || r.Recipients != 1 {
		t.Fatalf("the sub-album inherits its parent's club: %+v", r)
	}
}

func TestOwnerWhoIsAlsoAdminReceivesOneEmail(t *testing.T) {
	f, _ := standard(t)
	f.user("both@x.test", "Bo", "clubphoto_respo", "ADMIN")
	if _, err := f.submit("reporter@x.test", photoID, ReasonOther, "details"); err != nil {
		t.Fatal(err)
	}
	f.deliver()
	if f.m.calls["both@x.test"] != 1 {
		t.Fatalf("both@ mailed %d times", f.m.calls["both@x.test"])
	}
}

func TestRecipientLanguageIsUsed(t *testing.T) {
	f, _ := standard(t)
	f.exec(`UPDATE newf SET language = (SELECT id_languages FROM languages WHERE code = 'en') WHERE email = 'admin1@x.test'`)
	if _, err := f.submit("reporter@x.test", photoID, ReasonOther, "details"); err != nil {
		t.Fatal(err)
	}
	f.deliver()
	for _, s := range f.m.sent {
		want := "fr"
		if s.to == "admin1@x.test" {
			want = "en"
		}
		if s.language != want || s.subject != "photo_report_email.subject" {
			t.Fatalf("%s: language %q subject %q", s.to, s.language, s.subject)
		}
	}
}

func TestValidation(t *testing.T) {
	f, _ := standard(t)
	f.exec(`INSERT INTO photos (drive_file_id, name, mime_type, media_type, sort_at, sort_at_source, web_url, last_seen_run, deleted_at)
		VALUES ('DELETED-000001', 'd.jpg', 'image/jpeg', 'image', now(), 'x', 'u', 1, now())`)
	f.exec(`INSERT INTO photos (drive_file_id, name, mime_type, media_type, sort_at, sort_at_source, web_url, last_seen_run)
		VALUES ('DOCUMENT-00001', 'd.pdf', 'application/pdf', 'other', now(), 'x', 'u', 1)`)

	cases := []struct {
		name    string
		photo   string
		reason  Reason
		comment string
		email   string
		want    error
	}{
		{"unknown reason", photoID, "spam", "", "reporter@x.test", ErrInvalidReason},
		{"empty reason", photoID, "", "", "reporter@x.test", ErrInvalidReason},
		{"short photo id", "abc", ReasonOther, "details", "reporter@x.test", ErrInvalidPhotoID},
		{"empty photo id", "", ReasonOther, "details", "reporter@x.test", ErrInvalidPhotoID},
		{"unknown photo", "NOPE-0000000000", ReasonOther, "details", "reporter@x.test", ErrPhotoNotFound},
		{"deleted photo", "DELETED-000001", ReasonOther, "details", "reporter@x.test", ErrPhotoNotFound},
		{"not an image or video", "DOCUMENT-00001", ReasonOther, "details", "reporter@x.test", ErrPhotoNotFound},
		{"comment too long", photoID, ReasonOther, strings.Repeat("é", MaxCommentRunes+1), "reporter@x.test", ErrCommentTooLong},
		{"huge comment", photoID, ReasonOther, strings.Repeat("a", 100000), "reporter@x.test", ErrCommentTooLong},
		{"unknown reporter", photoID, ReasonOther, "details", "ghost@x.test", ErrUnknownReporter},
		{"other without comment", photoID, ReasonOther, "", "reporter@x.test", ErrCommentRequired},
		{"other with a blank comment", photoID, ReasonOther, " \x00\r\n\t ", "reporter@x.test", ErrCommentRequired},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if _, err := f.submit(c.email, c.photo, c.reason, c.comment); !errors.Is(err, c.want) {
				t.Fatalf("err = %v, want %v", err, c.want)
			}
		})
	}
	var n int
	_ = f.db.QueryRow(`SELECT count(*) FROM photo_reports`).Scan(&n)
	if n != 0 {
		t.Fatalf("%d reports stored from invalid submissions", n)
	}

	// Only "other" needs a comment.
	for _, reason := range []Reason{ReasonAppearsInPhoto, ReasonInappropriate, ReasonPoorQuality} {
		if _, err := f.submit("reporter@x.test", photoID, reason, ""); err != nil {
			t.Fatalf("%s without comment: %v", reason, err)
		}
	}

	// The limit itself is accepted; control characters are dropped, line breaks kept.
	if _, err := f.submit("reporter@x.test", photoID, ReasonOther, strings.Repeat("é", MaxCommentRunes)); err != nil {
		t.Fatalf("a comment of exactly %d characters: %v", MaxCommentRunes, err)
	}
	if _, err := f.submit("reporter@x.test", photoID, ReasonPoorQuality, "  a\x00b\r\nc\x07  "); err != nil {
		t.Fatal(err)
	}
	var stored string
	_ = f.db.QueryRow(`SELECT comment FROM photo_reports WHERE reason = 'poor_quality' AND comment <> ''`).Scan(&stored)
	if stored != "ab\nc" {
		t.Fatalf("stored %q", stored)
	}
}

func TestDuplicateSubmissionsAreAbsorbedButDistinctReportsAreKept(t *testing.T) {
	f, _ := standard(t)
	first, err := f.submit("reporter@x.test", photoID, ReasonInappropriate, "Texte  du\nsignalement")
	if err != nil {
		t.Fatal(err)
	}
	f.deliver()
	mails := len(f.m.sent)

	again, err := f.submit("reporter@x.test", photoID, ReasonInappropriate, " texte du signalement ")
	if err != nil || !again.Duplicate || again.ID != first.ID {
		t.Fatalf("double tap: %+v %v, want the first report back", again, err)
	}
	f.deliver()
	if len(f.m.sent) != mails {
		t.Fatal("a duplicate must not send new emails")
	}

	other, _ := f.submit("reporter@x.test", photoID, ReasonPoorQuality, "Texte du signalement")
	other2, _ := f.submit("reporter@x.test", photoID, ReasonInappropriate, "Un autre commentaire")
	f.photo("PHOTO-0000000002", nil)
	other3, _ := f.submit("reporter@x.test", "PHOTO-0000000002", ReasonInappropriate, "Texte du signalement")
	for _, r := range []Receipt{other, other2, other3} {
		if r.Duplicate || r.ID == first.ID {
			t.Fatalf("a different reason, comment or photo is a different report: %+v", r)
		}
	}

	f.user("friend@x.test", "Fay", "NEWF")
	if r, _ := f.submit("friend@x.test", photoID, ReasonInappropriate, "Texte du signalement"); r.Duplicate {
		t.Fatal("another user's report is never a duplicate")
	}

	f.exec(`UPDATE photo_reports SET created_at = now() - interval '25 hours'`)
	if r, _ := f.submit("reporter@x.test", photoID, ReasonInappropriate, "Texte du signalement"); r.Duplicate {
		t.Fatal("after a day the same report is allowed again")
	}
}

func TestConcurrentIdenticalSubmissionsStoreOneReport(t *testing.T) {
	f, _ := standard(t)
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = f.submit("reporter@x.test", photoID, ReasonOther, "same")
		}()
	}
	wg.Wait()
	var n int
	_ = f.db.QueryRow(`SELECT count(*) FROM photo_reports`).Scan(&n)
	if n != 1 {
		t.Fatalf("%d reports stored, want 1", n)
	}
}

func TestRateLimit(t *testing.T) {
	f, _ := standard(t)
	for i := 0; i < burstLimit; i++ {
		if _, err := f.submit("reporter@x.test", photoID, ReasonOther, fmt.Sprint("c", i)); err != nil {
			t.Fatalf("report %d: %v", i, err)
		}
	}
	_, err := f.submit("reporter@x.test", photoID, ReasonOther, "one too many")
	var limited *RateLimitedError
	if !errors.As(err, &limited) || limited.RetryAfter <= 0 || limited.RetryAfter > burstWindow {
		t.Fatalf("err = %v, want a rate limit with a wait", err)
	}
	// An exact repeat is still absorbed, not counted as spam.
	if r, err := f.submit("reporter@x.test", photoID, ReasonOther, "c0"); err != nil || !r.Duplicate {
		t.Fatalf("repeat while limited: %+v %v", r, err)
	}
	f.user("friend@x.test", "Fay", "NEWF")
	if _, err := f.submit("friend@x.test", photoID, ReasonOther, "mine"); err != nil {
		t.Fatalf("the limit is per user: %v", err)
	}

	f.exec(`UPDATE photo_reports SET created_at = now() - interval '11 minutes'`)
	if _, err := f.submit("reporter@x.test", photoID, ReasonOther, "later"); err != nil {
		t.Fatalf("after the burst window: %v", err)
	}
}

func TestDailyLimit(t *testing.T) {
	f, _ := standard(t)
	for i := 0; i < dayLimit; i++ {
		if _, err := f.submit("reporter@x.test", photoID, ReasonOther, fmt.Sprint("c", i)); err != nil {
			f.exec(`UPDATE photo_reports SET created_at = now() - interval '1 hour'`)
			if _, err = f.submit("reporter@x.test", photoID, ReasonOther, fmt.Sprint("c", i)); err != nil {
				t.Fatalf("report %d: %v", i, err)
			}
		}
	}
	f.exec(`UPDATE photo_reports SET created_at = now() - interval '1 hour'`)
	var limited *RateLimitedError
	if _, err := f.submit("reporter@x.test", photoID, ReasonOther, "over the day"); !errors.As(err, &limited) {
		t.Fatalf("err = %v, want the daily limit", err)
	}
}

func TestEmailFailureKeepsTheReportAndRetriesOnlyWhatFailed(t *testing.T) {
	f, _ := standard(t)
	f.m.set("owner@x.test", true)
	rec, err := f.submit("reporter@x.test", photoID, ReasonAppearsInPhoto, "Merci")
	if err != nil {
		t.Fatal(err)
	}
	if n := f.deliver(); n != 2 {
		t.Fatalf("delivered %d, want the 2 administrators", n)
	}
	r, _ := f.svc.Get(context.Background(), rec.ID)
	if r.DeliveryStatus != "pending" || r.RecipientsSent != 2 {
		t.Fatalf("a failed email must not count as sent: %+v", r)
	}
	var failed *Delivery
	for i := range r.Deliveries {
		if r.Deliveries[i].Email == "owner@x.test" {
			failed = &r.Deliveries[i]
		}
	}
	if failed == nil || failed.Status != "failed" || failed.Attempts != 1 || failed.SentAt != nil || failed.LastError == nil {
		t.Fatalf("owner delivery %+v", failed)
	}

	// Not due yet: backoff.
	if f.deliver() != 0 || f.m.calls["owner@x.test"] != 1 {
		t.Fatal("a failed email waits for its retry delay")
	}

	f.m.set("owner@x.test", false)
	f.due()
	if n := f.deliver(); n != 1 {
		t.Fatalf("retry delivered %d, want only the owner", n)
	}
	if f.m.calls["admin1@x.test"] != 1 || f.m.calls["admin2@x.test"] != 1 || f.m.calls["owner@x.test"] != 2 {
		t.Fatalf("calls %v: nobody may be mailed twice", f.m.calls)
	}
	if r, _ = f.svc.Get(context.Background(), rec.ID); r.DeliveryStatus != "sent" {
		t.Fatalf("status %s", r.DeliveryStatus)
	}
}

func TestExhaustedAttemptsAreFailedUntilAnAdministratorRetries(t *testing.T) {
	f, _ := standard(t)
	f.m.set("owner@x.test", true)
	rec, _ := f.submit("reporter@x.test", photoID, ReasonOther, "details")
	for i := 0; i < MaxAttempts+2; i++ {
		f.due()
		f.deliver()
	}
	if f.m.calls["owner@x.test"] != MaxAttempts {
		t.Fatalf("owner tried %d times, want %d", f.m.calls["owner@x.test"], MaxAttempts)
	}
	r, _ := f.svc.Get(context.Background(), rec.ID)
	if r.DeliveryStatus != "failed" || r.RecipientsFailed != 1 {
		t.Fatalf("%+v", r)
	}

	f.m.set("owner@x.test", false)
	n, err := f.svc.RetryFailed(context.Background(), rec.ID)
	if err != nil || n != 1 {
		t.Fatalf("retry: %d %v", n, err)
	}
	if r, _ = f.svc.Get(context.Background(), rec.ID); r.DeliveryStatus != "sent" {
		t.Fatalf("status after retry: %s", r.DeliveryStatus)
	}
	if f.m.calls["admin1@x.test"] != 1 {
		t.Fatal("retrying must not resend what was accepted")
	}
}

func TestPanickingSenderIsAFailureNotACrash(t *testing.T) {
	f, _ := standard(t)
	f.m.panics["owner@x.test"] = true
	rec, _ := f.submit("reporter@x.test", photoID, ReasonOther, "details")
	if n := f.deliver(); n != 2 {
		t.Fatalf("delivered %d", n)
	}
	if r, _ := f.svc.Get(context.Background(), rec.ID); r.DeliveryStatus != "pending" {
		t.Fatalf("%+v", r)
	}
}

func TestWithoutMailerReportsStayPending(t *testing.T) {
	f, _ := standard(t)
	f.svc.mailer = nil
	rec, err := f.submit("reporter@x.test", photoID, ReasonOther, "details")
	if err != nil {
		t.Fatal(err)
	}
	f.deliver()
	r, _ := f.svc.Get(context.Background(), rec.ID)
	if r.DeliveryStatus != "pending" || r.Deliveries[0].Attempts != 0 {
		t.Fatalf("%+v", r)
	}
	f.svc.mailer = f.m
	if f.deliver() != 3 {
		t.Fatal("pending reports are sent once a mailer is available")
	}
}

func TestInterruptedSenderIsReclaimed(t *testing.T) {
	f, _ := standard(t)
	rec, _ := f.submit("reporter@x.test", photoID, ReasonOther, "details")
	f.exec(`UPDATE photo_report_deliveries SET status = 'sending', attempts = 1, lease_until = now() + interval '5 minutes'`)
	if f.deliver() != 0 {
		t.Fatal("a leased delivery belongs to its sender")
	}
	f.exec(`UPDATE photo_report_deliveries SET lease_until = now() - interval '1 minute'`)
	if f.deliver() != 3 {
		t.Fatal("an expired lease is taken over")
	}
	f.exec(`UPDATE photo_report_deliveries SET status = 'sending', attempts = $1, lease_until = now() - interval '1 minute' WHERE recipient_email = 'owner@x.test'`, MaxAttempts)
	f.deliver()
	if r, _ := f.svc.Get(context.Background(), rec.ID); r.DeliveryStatus != "failed" {
		t.Fatalf("a sender that died on its last attempt leaves the delivery failed: %+v", r)
	}
}

func TestConcurrentDeliveryNeverMailsTwice(t *testing.T) {
	f, _ := standard(t)
	for i := 0; i < 4; i++ {
		f.user(fmt.Sprintf("extra%d@x.test", i), "X", "ADMIN")
	}
	if _, err := f.submit("reporter@x.test", photoID, ReasonOther, "details"); err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); _, _ = f.svc.DeliverDue(context.Background()) }()
	}
	wg.Wait()
	for who, n := range f.m.calls {
		if n != 1 {
			t.Fatalf("%s mailed %d times", who, n)
		}
	}
	if len(f.m.calls) != 7 {
		t.Fatalf("%d recipients mailed, want 7", len(f.m.calls))
	}
}

func TestAdminListFiltersAndAlbumClub(t *testing.T) {
	f, club := standard(t)
	a, _ := f.submit("reporter@x.test", photoID, ReasonInappropriate, "")
	f.deliver()
	f.m.set("owner@x.test", true)
	f.user("friend@x.test", "Fay", "NEWF")
	b, _ := f.submit("friend@x.test", photoID, ReasonPoorQuality, "")
	f.deliver()

	ctx := context.Background()
	all, total, err := f.svc.List(ctx, ListFilter{})
	if err != nil || total != 2 || len(all) != 2 || all[0].ID != b.ID {
		t.Fatalf("list %v %d %v", all, total, err)
	}
	if got, total, _ := f.svc.List(ctx, ListFilter{Reason: "inappropriate"}); total != 1 || got[0].ID != a.ID {
		t.Fatalf("reason filter: %v", got)
	}
	if got, _, _ := f.svc.List(ctx, ListFilter{Delivery: "pending"}); len(got) != 1 || got[0].ID != b.ID {
		t.Fatalf("delivery filter: %v", got)
	}
	if got, _, _ := f.svc.List(ctx, ListFilter{Limit: 1, Offset: 1}); len(got) != 1 || got[0].ID != a.ID {
		t.Fatalf("paging: %v", got)
	}
	if all[1].Reporter.Email == nil || *all[1].Reporter.Email != "reporter@x.test" || !all[1].PhotoAvailable || all[1].PhotoURL == "" {
		t.Fatalf("admin details: %+v", all[1])
	}

	// The photo disappearing from Drive does not erase the report.
	f.exec(`UPDATE photos SET deleted_at = now()`)
	if r, _ := f.svc.Get(ctx, a.ID); r.PhotoAvailable || r.PhotoName == "" {
		t.Fatalf("snapshot: %+v", r)
	}
	// Account deletion drops the email but keeps the report.
	f.exec(`DELETE FROM newf WHERE email = 'reporter@x.test'`)
	if r, _ := f.svc.Get(ctx, a.ID); r.Reporter.Email != nil || r.Reporter.Name != "Rita Doe" {
		t.Fatalf("after account deletion: %+v", r.Reporter)
	}

	var albumID int64
	_ = f.db.QueryRow(`SELECT id_photo_albums FROM photo_albums LIMIT 1`).Scan(&albumID)
	if err := f.svc.SetAlbumClub(ctx, albumID, nil); err != nil {
		t.Fatal(err)
	}
	if err := f.svc.SetAlbumClub(ctx, albumID, &club); err != nil {
		t.Fatal(err)
	}
	bad := int64(99999)
	if err := f.svc.SetAlbumClub(ctx, albumID, &bad); !errors.Is(err, ErrClubNotFound) {
		t.Fatalf("unknown club: %v", err)
	}
	if err := f.svc.SetAlbumClub(ctx, 99999, &club); !errors.Is(err, ErrAlbumNotFound) {
		t.Fatalf("unknown album: %v", err)
	}
	if _, err := f.svc.Get(ctx, 99999); !errors.Is(err, ErrReportNotFound) {
		t.Fatalf("unknown report: %v", err)
	}
}

var i18nOnce sync.Once

func initI18n() (err error) {
	i18nOnce.Do(func() { err = i18n.Init() })
	return err
}

// renderEmail does what EmailService.SendEmail does with the template and the translations, without the network.
func renderEmail(t *testing.T, lang string, data EmailData) (subject, body string) {
	t.Helper()
	if err := initI18n(); err != nil {
		t.Fatal(err)
	}
	loc := i18n.GetLocalizer(lang)
	raw, err := os.ReadFile("../../email_templates/email_template_photo_report.html")
	if err != nil {
		t.Fatal(err)
	}
	translate := func(key string, _ ...interface{}) template.HTML {
		msg, err := loc.Localize(&goi18n.LocalizeConfig{MessageID: key, TemplateData: data})
		if err != nil {
			return template.HTML("!!!" + key + "!!!")
		}
		return template.HTML(msg)
	}
	tmpl, err := template.New("report").Parse(string(raw))
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, struct {
		Data interface{}
		T    func(string, ...interface{}) template.HTML
	}{data, translate}); err != nil {
		t.Fatal(err)
	}
	subject, err = loc.Localize(&goi18n.LocalizeConfig{MessageID: "photo_report_email.subject", TemplateData: data})
	if err != nil {
		t.Fatal(err)
	}
	return subject, buf.String()
}

func TestEmailTemplateEscapesUserTextInEveryLanguage(t *testing.T) {
	hostile := EmailData{
		SubjectTitle: "Gala", ReportID: 7, ReasonCode: "other",
		Comment: `<script>alert(1)</script> <img src=x onerror=alert(2)>` + "\nligne 2",
		PhotoID: "ID<b>", PhotoName: `"><svg onload=alert(3)>.jpg`, PhotoURL: "javascript:alert(4)",
		AlbumID: 3, AlbumTitle: "<i>Gala</i>", AlbumURL: "https://drive.example/a?x=1&y=2",
		ClubName: "<u>Club</u>", HasClub: true, ReporterName: "<b>Eve</b>", ReporterMail: "eve@x.test",
		ReporterID: 9, Formation: "FISE", SubmittedAt: "2026-10-11 10:00 UTC",
	}
	for _, lang := range []string{"fr", "en", "de", "es", "pt"} {
		for _, reason := range []string{"appears_in_photo", "inappropriate", "poor_quality", "other"} {
			for _, owner := range []bool{true, false} {
				d := hostile
				d.ReasonCode, d.IsClubOwner = reason, owner
				subject, body := renderEmail(t, lang, d)
				for _, raw := range []string{"<script", "<img src=x", "<svg", "<b>Eve", "<i>Gala", "<u>Club", "javascript:"} {
					if strings.Contains(body, raw) {
						t.Fatalf("%s/%s: unescaped %q in the body", lang, reason, raw)
					}
				}
				if !strings.Contains(body, "&lt;script&gt;alert(1)&lt;/script&gt;") {
					t.Fatalf("%s: the comment should be shown, escaped", lang)
				}
				if strings.Contains(body, "!!!") || strings.Contains(subject, "!!!") {
					t.Fatalf("%s/%s: missing translation in %s", lang, reason, body)
				}
				if !strings.Contains(subject, "Gala") || !strings.HasPrefix(subject, "[Transat]") {
					t.Fatalf("%s: subject %q", lang, subject)
				}
			}
		}
	}
	// Without comment, event and club the email still reads well and links to the photo.
	_, body := renderEmail(t, "fr", EmailData{SubjectTitle: "x.jpg", ReasonCode: "other", PhotoName: "x.jpg", PhotoURL: "https://drive.example/f", SubmittedAt: "now"})
	for _, want := range []string{"Aucun commentaire", "Aucun club associ", "https://drive.example/f", "administrateurs"} {
		if !strings.Contains(body, want) {
			t.Fatalf("missing %q", want)
		}
	}
}

func TestSubjectCannotCarryLineBreaks(t *testing.T) {
	f, _ := standard(t)
	f.exec(`UPDATE photo_albums SET title = E'Gala\r\nBcc: evil@x.test ' || repeat('x', 200)`)
	if _, err := f.submit("reporter@x.test", photoID, ReasonOther, "details"); err != nil {
		t.Fatal(err)
	}
	f.deliver()
	title := f.m.sent[0].data.SubjectTitle
	if strings.ContainsAny(title, "\r\n") || utf8.RuneCountInString(title) > 80 {
		t.Fatalf("subject title %q", title)
	}
}
