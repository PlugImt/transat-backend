// Command photobooth-devserver runs the photobooth API on a local database with an in-memory Google Drive, so a
// booth can be tested end to end without Drive credentials:
//
//	docker run -d --rm --name pb-pg -e POSTGRES_PASSWORD=pg -e POSTGRES_DB=transat_dev -p 54329:5432 postgres:18-alpine
//	go run ./scripts/photobooth-devserver
//	adb reverse tcp:8099 tcp:8099   # the booth then uses http://localhost:8099 (debug builds only)
//
// GET /fake-drive/tree shows the folders and files the booth produced.
package main

import (
	"context"
	"crypto/md5" // #nosec G501 -- Google Drive identifies files by their MD5 checksum
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	photoboothHandler "github.com/plugimt/transat-backend/handlers/photobooth"
	"github.com/plugimt/transat-backend/internal/database"
	"github.com/plugimt/transat-backend/routes"
	pb "github.com/plugimt/transat-backend/services/photobooth"
)

type file struct {
	pb.FileInfo
	data []byte
}

type memDrive struct {
	public   string
	mu       sync.Mutex
	seq      int
	folders  map[string]pb.Folder
	parent   map[string]string
	files    map[string]*file
	sessions map[string]pb.FileInfo
}

func (d *memDrive) id(p string) string { d.seq++; return fmt.Sprintf("%s%012d", p, d.seq) }

func (d *memDrive) ListChildFolders(_ context.Context, parent string) ([]pb.Folder, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	var out []pb.Folder
	for id, p := range d.parent {
		if p == parent {
			out = append(out, d.folders[id])
		}
	}
	return out, nil
}

func (d *memDrive) ListChildFilesNamed(_ context.Context, parent, name string) ([]pb.FileInfo, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	var out []pb.FileInfo
	for _, f := range d.files {
		if f.Parents[0] == parent && f.Name == name && !f.Trashed {
			out = append(out, f.FileInfo)
		}
	}
	return out, nil
}

func (d *memDrive) CreateFolder(_ context.Context, parent, name string) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	id := d.id("folder")
	d.folders[id] = pb.Folder{ID: id, Name: name}
	d.parent[id] = parent
	return id, nil
}

func (d *memDrive) StartUpload(_ context.Context, parent, name, mime string, size int64) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	sid := d.id("session")
	d.sessions[sid] = pb.FileInfo{Name: name, Parents: []string{parent}, Size: size, MimeType: mime}
	return d.public + "/fake-drive/upload/" + sid, nil
}

func (d *memDrive) receive(sid string, data []byte) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	s, ok := d.sessions[sid]
	if !ok {
		return "", fmt.Errorf("unknown session")
	}
	if int64(len(data)) != s.Size {
		return "", fmt.Errorf("size %d, announced %d", len(data), s.Size)
	}
	delete(d.sessions, sid)
	sum := md5.Sum(data) // #nosec G401 -- same checksum as Google Drive, not a security use
	f := &file{FileInfo: s, data: data}
	f.ID, f.MD5 = d.id("file"), hex.EncodeToString(sum[:])
	d.files[f.ID] = f
	return f.ID, nil
}

func (d *memDrive) GetFile(_ context.Context, id string) (pb.FileInfo, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if f, ok := d.files[id]; ok {
		return f.FileInfo, nil
	}
	if f, ok := d.folders[id]; ok {
		return pb.FileInfo{ID: f.ID, Name: f.Name}, nil
	}
	return pb.FileInfo{}, pb.ErrNotFound
}

func (d *memDrive) Trash(_ context.Context, id string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if f, ok := d.files[id]; ok {
		f.Trashed = true
		log.Printf("fake-drive: trashed %s", f.Name)
	}
	return nil
}

func (d *memDrive) Download(_ context.Context, id string) (io.ReadCloser, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	f, ok := d.files[id]
	if !ok || f.Trashed {
		return nil, pb.ErrNotFound
	}
	return io.NopCloser(strings.NewReader(string(f.data))), nil
}

func (d *memDrive) Thumbnail(context.Context, pb.FileInfo, int) (io.ReadCloser, string, error) {
	return nil, "", pb.ErrNotFound
}

func (d *memDrive) tree() string {
	d.mu.Lock()
	defer d.mu.Unlock()
	var lines []string
	var walk func(id, indent string)
	walk = func(id, indent string) {
		var kids []string
		for k, p := range d.parent {
			if p == id {
				kids = append(kids, k)
			}
		}
		sort.Slice(kids, func(a, b int) bool { return d.folders[kids[a]].Name < d.folders[kids[b]].Name })
		for _, k := range kids {
			lines = append(lines, indent+d.folders[k].Name+"/")
			walk(k, indent+"  ")
		}
		var names []string
		for _, f := range d.files {
			if f.Parents[0] == id {
				state := ""
				if f.Trashed {
					state = "  [trashed]"
				}
				names = append(names, fmt.Sprintf("%s%s  (%d kB)%s", indent, f.Name, len(f.data)/1024, state))
			}
		}
		sort.Strings(names)
		lines = append(lines, names...)
	}
	walk("root", "")
	return strings.Join(lines, "\n") + "\n"
}

func env(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func main() {
	port := env("PORT", "8099")
	dsn := env("DSN", "host=localhost port=54329 user=postgres password=pg dbname=transat_dev sslmode=disable")
	db, err := database.Open(dsn, "db/migrations")
	if err != nil {
		log.Fatal(err)
	}
	secret := os.Getenv("PHOTOBOOTH_DEV_SECRET")
	if secret == "" {
		b := make([]byte, 32)
		_, _ = rand.Read(b)
		secret = base64.RawURLEncoding.EncodeToString(b)
	}
	keys, err := pb.ParseDeviceKeys("booth-dev:" + secret)
	if err != nil {
		log.Fatal(err)
	}
	drive := &memDrive{
		public: env("PUBLIC_URL", "http://localhost:"+port), folders: map[string]pb.Folder{"root": {ID: "root", Name: "Photos"}},
		parent: map[string]string{}, files: map[string]*file{}, sessions: map[string]pb.FileInfo{},
	}
	svc := pb.NewService(db, drive, "root", env("SHARE_BASE_URL", "http://localhost:5173/photobooth"))

	app := fiber.New(fiber.Config{StrictRouting: true, BodyLimit: 100 << 20})
	app.Use(logger.New())
	app.Put("/fake-drive/upload/:sid", func(c *fiber.Ctx) error {
		id, err := drive.receive(c.Params("sid"), append([]byte(nil), c.Body()...))
		if err != nil {
			return c.Status(400).SendString(err.Error())
		}
		return c.JSON(fiber.Map{"id": id})
	})
	app.Get("/fake-drive/tree", func(c *fiber.Ctx) error { return c.SendString(drive.tree()) })
	app.Post("/fake-drive/trash-job", func(c *fiber.Ctx) error {
		n, err := svc.RunTrash(c.UserContext())
		if err != nil {
			return err
		}
		return c.JSON(fiber.Map{"done": n})
	})
	routes.SetupPhotoboothRoutes(app, photoboothHandler.NewHandler(svc, pb.NewVerifier(keys)))

	// #nosec G706 -- local development tool: printing the generated key is its purpose
	log.Printf("Booth key for the admin panel: booth-dev.%s", secret)
	log.Fatal(app.Listen(":" + port))
}
