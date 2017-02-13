package box

import (
	"time"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Login string `json:"login"`
}

type File struct {
	fileId            string      `json:"id,omitempty"`                  // Box’s unique string identifying this file.
	SequenceId        string      `json:"sequence_id,omitempty"`         // A unique ID for use with the /events endpoint.
	ETag              string      `json:"etag,omitempty"`                // A unique string identifying the version of this file.
	Sha1              string      `json:"sha1,omitempty"`                // The sha1 hash of this file.
	Name              string      `json:"name,omitempty"`                // The name of this file.
	Description       string      `json:"description,omitempty"`         // The description of this file.
	Size              int         `json:"size,omitempty"`                // Size of this file in bytes.
	PathCollection    *Collection `json:"path_collection,omitempty"`     // The path of folders to this item, starting at the root.
	CreatedAt         *string     `json:"created_at,omitempty"`          // When this file was created on Box’s servers.
	ModifiedAt        *string     `json:"modified_at,omitempty"`         // When this file was last updated on the Box servers.
	ThrashedAt        *string     `json:"thrashed_at,omitempty"`         // When this file was last moved to the trash.
	PurgedAt          *string     `json:"purged_at,omitempty"`           // When this file will be permanently deleted.
	ContentCreatedAt  *string     `json:"content_created_at,omitempty"`  // When the content of this file was created.
	ContentModifiedAt *string     `json:"content_modified_at,omitempty"` // When the content of this file was last modified.
	CreatedBy         *Entity     `json:"created_by,omitempty"`          // The user who first created file.
	ModifiedBy        *Entity     `json:"modified_by,omitempty"`         // The user who last updated this file.
	OwnedBy           *Entity     `json:"owned_by,omitempty"`            // The user who owns this file.
	Parent            *Entity     `json:"parent,omitempty"`              // The folder containing this file.
	ItemStatus        string      `json:"item_status,omitempty"`         // Whether this item is deleted or not.
	VersionNumber     string      `json:"version_number,omitempty"`      // The version of the file.
	CommentCount      int         `json:"comment_count,omitempty"`       // The number of comments on a file.
	Tags              []string    `json:"tags,omitempty"`                // All tags applied to this file.
	Extension         string      `json:"extension,omitempty"`           // Indicates the suffix, when available, on the file.
	//SharedLink        *SharedObject `json:"shared_link,omitempty"`         // The shared link object for this file.
	//Lock              *BoxLock      `json:"lock,omitempty"`                // The lock held on the file.
	//Permissions       *Permission   `json:"permissions,omitempty"`         // The permissions that the current user has on this file.
}

type Folder struct {
	folderId          string      `json:"id,omitempty"`                  // The folder’s ID.
	SequenceId        string      `json:"sequence_id,omitempty"`         // A unique ID for use with the /events endpoint.
	ETag              string      `json:"etag,omitempty"`                // A unique string identifying the version of this folder.
	Name              string      `json:"name,omitempty"`                // The name of this folder.
	Description       string      `json:"description,omitempty"`         // The description of this folder.
	Size              int         `json:"size,omitempty"`                // Size of this file in bytes.
	PathCollection    *Collection `json:"path_collection,omitempty"`     // The path of folders to this item, starting at the root.
	CreatedAt         *string     `json:"created_at,omitempty"`          // The time the folder was created.
	ModifiedAt        *string     `json:"modified_at,omitempty"`         // The time the folder or its contents were last modified.
	ThrashedAt        *string     `json:"thrashed_at,omitempty"`         // The time the folder or its contents were put in the trash.
	PurgedAt          *string     `json:"purged_at,omitempty"`           // The time the folder or its contents were purged from the trash.
	ContentCreatedAt  *string     `json:"content_created_at,omitempty"`  // The time the folder or its contents were originally created (according to the uploader).
	ContentModifiedAt *string     `json:"content_modified_at,omitempty"` // The time the folder or its contents were last modified (according to the uploader).
	CreatedBy         *Entity     `json:"created_by,omitempty"`          // The user who created this folder.
	ModifiedBy        *Entity     `json:"modified_by,omitempty"`         // The user who last modified this folder.
	OwnedBy           *Entity     `json:"owned_by,omitempty"`            // The user who owns this file.
	Parent            *Entity     `json:"parent,omitempty"`              // The folder that contains this one.
	ItemStatus        string      `json:"item_status,omitempty"`         // Whether this item is deleted or not.
	Tags              []string    `json:"tags,omitempty"`                // All tags applied to this file.
	HasCollaborations bool        `json:"has_collaborations,omitempty"`  // Whether this folder has any collaborators.
	SyncStatus        string      `json:"sync_status,omitempty"`         // Whether this folder will be synced by the Box sync clients or not. Can be
	ItemCollection    *Collection `json:"item_collection,omitempty"`     // A collection of mini file and folder objects contained in this folder.
	//FolderUploadEmail *UploadEmail  `json:"folder_upload_email,omitempty"` // The upload email address for this folder. Null if not set.
	//SharedLink        *SharedObject `json:"shared_link,omitempty"`         // The shared link for this folder. Null if not set..
	//Permissions       *Permission   `json:"permissions,omitempty"`         // The permissions that the current user has on this file.
}

//  Represents both mini folder and mini file.
type Entity struct {
	SequenceId string `json:"sequence_id,omitempty"` // A unique ID for use with the /events endpoint.
	Name       string `json:"name,omitempty"`        // The name of the entity.
	Id         string `json:"id,omitempty"`          // The id of the entity.
	ETag       string `json:"etag,omitempty"`        // A unique string identifying the version of this entity.
	Type       string `json:"type,omitempty"`        // Type of entity
}

type Collection struct {
	Count  int      `json:"total_count,omitempty"`
	Entry  []Entity `json:"entries,omitempty"`
	Limit  int      `json:"limit,omitempty"`
	Offset int      `json:"offset,omitempty"`
}

type BoxTime time.Time

type FolderEntity struct {
	FolderId string `json:"id,omitempty"`
	Files    []File
	Folders  []Folder
}
