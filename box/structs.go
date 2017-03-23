package box

import (
	"encoding/json"
	"time"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Login string `json:"login"`
}

type File struct {
	ID                string     `json:"id"`                  // Box’s unique string identifying this file.
	SequenceID        string     `json:"sequence_id"`         // A unique ID for use with the /events endpoint.
	ETag              string     `json:"etag"`                // A unique string identifying the version of this file.
	SHA1              string     `json:"sha1"`                // The sha1 hash of this file.
	Name              string     `json:"name"`                // The name of this file.
	Description       string     `json:"description"`         // The description of this file.
	Size              int        `json:"size"`                // Size of this file in bytes.
	PathCollection    Collection `json:"path_collection"`     // The path of folders to this item, starting at the root.
	CreatedAt         time.Time  `json:"created_at"`          // When this file was created on Box’s servers.
	ModifiedAt        time.Time  `json:"modified_at"`         // When this file was last updated on the Box servers.
	ContentCreatedAt  time.Time  `json:"content_created_at"`  // When the content of this file was created.
	ContentModifiedAt time.Time  `json:"content_modified_at"` // When the content of this file was last modified.
	CreatedBy         User       `json:"created_by"`          // The user who first created file.
	ModifiedBy        User       `json:"modified_by"`         // The user who last updated this file.
	OwnedBy           User       `json:"owned_by"`            // The user who owns this file.
	Parent            *Folder    `json:"parent"`              // The folder containing this file.
	ItemStatus        string     `json:"item_status"`         // Whether this item is deleted or not.
	VersionNumber     string     `json:"version_number"`      // The version of the file.
	CommentCount      int        `json:"comment_count"`       // The number of comments on a file.
	Tags              []string   `json:"tags"`                // All tags applied to this file.
	Extension         string     `json:"extension"`           // Indicates the suffix, when available, on the file.
}

type Folder struct {
	ID                string     `json:"id"`                  // The folder’s ID.
	SequenceID        string     `json:"sequence_id"`         // A unique ID for use with the /events endpoint.
	ETag              string     `json:"etag"`                // A unique string identifying the version of this folder.
	Name              string     `json:"name"`                // The name of this folder.
	Description       string     `json:"description"`         // The description of this folder.
	Size              int        `json:"size"`                // Size of this file in bytes.
	PathCollection    Collection `json:"path_collection"`     // The path of folders to this item, starting at the root.
	CreatedAt         time.Time  `json:"created_at"`          // The time the folder was created.
	ModifiedAt        time.Time  `json:"modified_at"`         // The time the folder or its contents were last modified.
	ContentCreatedAt  time.Time  `json:"content_created_at"`  // The time the folder or its contents were originally created (according to the uploader).
	ContentModifiedAt time.Time  `json:"content_modified_at"` // The time the folder or its contents were last modified (according to the uploader).
	CreatedBy         User       `json:"created_by"`          // The user who created this folder.
	ModifiedBy        User       `json:"modified_by"`         // The user who last modified this folder.
	OwnedBy           User       `json:"owned_by"`            // The user who owns this file.
	Parent            *Folder    `json:"parent"`              // The folder that contains this one.
	ItemStatus        string     `json:"item_status"`         // Whether this item is deleted or not.
	Tags              []string   `json:"tags"`                // All tags applied to this file.
	HasCollaborations bool       `json:"has_collaborations"`  // Whether this folder has any collaborators.
	SyncStatus        string     `json:"sync_status"`         // Whether this folder will be synced by the Box sync clients or not. Can be
}

type Collection struct {
	Count   int               `json:"total_count"`
	Entries []json.RawMessage `json:"entries"`
	Limit   int               `json:"limit"`
	Offset  int               `json:"offset"`
}

type FolderContents struct {
	ID      string
	Files   []File
	Folders []Folder
}

type Event struct {
	EventID   string          `json:"event_id"`
	CreatedBy User            `json:"created_by"`
	EventType string          `json:"event_type"`
	SessionID string          `json:"session_id"`
	Source    json.RawMessage `json:"source"`
}

type EventCollection struct {
	ChunkSize          int     `json:"chunk_size"`
	NextStreamPosition int     `json:"next_stream_position"`
	Entries            []Event `json:"entries"`
}

type UploadAttributes struct {
	Name   string       `json:"name"`
	Parent UploadParent `json:"parent"`
}

type UploadParent struct {
	ID string `json:"id"`
}
