package photos

import (
	"context"
	"errors"
	"sort"
	"strings"
	"time"
)

var errNoExternal = errors.New("no reader for folders outside the shared drive")

// MediaKind is derived from the Drive mimeType: file extensions are unreliable (photobooth files have none).
type MediaKind string

const (
	MediaImage MediaKind = "image"
	MediaVideo MediaKind = "video"
	MediaOther MediaKind = "other" // RAW stored as octet-stream, sidecars…: indexed but never listed
)

func KindOfMime(mime string) MediaKind {
	switch {
	case strings.HasPrefix(mime, "image/"):
		return MediaImage
	case strings.HasPrefix(mime, "video/"):
		return MediaVideo
	default:
		return MediaOther
	}
}

// External reads folders that live outside the shared drive and are reached through a shortcut.
type External interface {
	GetItem(ctx context.Context, id string) (DriveItem, error)
	ListChildren(ctx context.Context, folderID string) ([]DriveItem, error)
}

type FolderInfo struct {
	ID            string
	ParentID      string // empty for a root or an external shortcut target
	Name          string // real Drive name
	EffectiveName string // name used for classification (the shortcut's label when reached through one)
	Path          string // logical path from the drive root, built from effective names
	ShortcutID    string
	WebURL        string
	ModifiedTime  time.Time
}

type MediaFile struct {
	DriveItem
	Kind     MediaKind
	FolderID string
}

type Snapshot struct {
	Folders []FolderInfo
	Media   []MediaFile // only files found while walking: external shortcut targets, or all files when items include them
	Results map[string]Result
}

// BuildSnapshot assembles the folder tree from a flat listing and classifies it.
// ext may be nil: shortcuts to folders outside the drive are then flagged unresolved.
// knownMedia gives the direct media count per folder when items holds no files (structure-only rebuilds).
func BuildSnapshot(ctx context.Context, items []DriveItem, driveID string, ext External) (*Snapshot, error) {
	inDrive := make(map[string]bool)
	for _, it := range items {
		if it.IsFolder() {
			inDrive[it.ID] = true
		}
	}

	childFolders := map[string][]DriveItem{}
	files := map[string][]DriveItem{}
	for _, it := range items {
		switch {
		case it.IsFolder(), it.IsShortcut() && it.ShortcutMime == MimeFolder:
			childFolders[it.ParentID] = append(childFolders[it.ParentID], it)
		case it.IsShortcut():
			// Shortcuts to files would duplicate a photo already indexed under its own ID.
		default:
			files[it.ParentID] = append(files[it.ParentID], it)
		}
	}

	snap := &Snapshot{}
	visited := map[string]bool{}

	var build func(it DriveItem, parentID, parentPath string) *Node
	build = func(it DriveItem, parentID, parentPath string) *Node {
		fi := FolderInfo{ID: it.ID, ParentID: parentID, Name: it.Name, EffectiveName: it.Name, WebURL: it.WebViewLink, ModifiedTime: it.ModifiedTime}
		node := &Node{ID: it.ID, Name: it.Name}
		path := joinPath(parentPath, it.Name)
		fi.Path = path

		childItems := childFolders[it.ID]
		fileItems := files[it.ID]

		if it.IsShortcut() {
			target := it.ShortcutTarget
			fi.ShortcutID, node.ShortcutID = it.ID, it.ID
			if inDrive[target] || visited[target] {
				node.AliasOfKnown = true
				snap.Folders = append(snap.Folders, fi)
				return node
			}
			var kids []DriveItem
			t, err := resolveExternal(ctx, ext, target)
			if err == nil {
				kids, err = ext.ListChildren(ctx, target)
			}
			if err != nil {
				node.TargetUnresolved = true
				snap.Folders = append(snap.Folders, fi)
				return node
			}
			visited[target] = true
			node.ID, fi.ID, fi.ParentID = target, target, ""
			fi.Name, fi.WebURL, fi.ModifiedTime = t.Name, t.WebViewLink, t.ModifiedTime
			for _, k := range kids {
				if k.IsFolder() {
					childItems = append(childItems, k)
				} else if !k.IsShortcut() {
					fileItems = append(fileItems, k)
				}
			}
		}
		visited[node.ID] = true

		for _, f := range fileItems {
			snap.Media = append(snap.Media, MediaFile{DriveItem: f, Kind: KindOfMime(f.MimeType), FolderID: node.ID})
		}
		snap.Folders = append(snap.Folders, fi)

		sort.Slice(childItems, func(a, b int) bool { return childItems[a].Name < childItems[b].Name })
		for _, c := range childItems {
			if c.IsFolder() && visited[c.ID] {
				continue
			}
			node.Children = append(node.Children, build(c, node.ID, path))
		}
		return node
	}

	var roots []*Node
	rootItems := childFolders[driveID]
	sort.Slice(rootItems, func(a, b int) bool { return rootItems[a].Name < rootItems[b].Name })
	for _, r := range rootItems {
		roots = append(roots, build(r, "", ""))
	}
	// Files directly at the drive root belong to no folder but stay in the gallery.
	for _, f := range files[driveID] {
		snap.Media = append(snap.Media, MediaFile{DriveItem: f, Kind: KindOfMime(f.MimeType)})
	}

	snap.Results = Classify(roots)
	return snap, nil
}

func joinPath(parent, name string) string {
	if parent == "" {
		return name
	}
	return parent + "/" + name
}

func resolveExternal(ctx context.Context, ext External, id string) (DriveItem, error) {
	if ext == nil {
		return DriveItem{}, errNoExternal
	}
	return ext.GetItem(ctx, id)
}
