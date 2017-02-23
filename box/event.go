package box

import (
	"encoding/json"
	"net/url"
)

const (
	EventTypeItemCreate               = "ITEM_CREATE"
	EventTypeItemUpload               = "ITEM_UPLOAD"
	EventTypeCommentCreate            = "COMMENT_CREATE"
	EventTypeCommentDelete            = "COMMENT_DELETE"
	EventTypeItemDownload             = "ITEM_DOWNLOAD"
	EventTypeItemPreview              = "ITEM_PREVIEW"
	EventTypeContentAccess            = "CONTENT_ACCESS"
	EventTypeItemMove                 = "ITEM_MOVE"
	EventTypeItemCopy                 = "ITEM_COPY"
	EventTypeTaskAssignmentCreate     = "TASK_ASSIGNMENT_CREATE"
	EventTypeTaskCreate               = "TASK_CREATE"
	EventTypeLockCreate               = "LOCK_CREATE"
	EventTypeLockDestroy              = "LOCK_DESTROY"
	EventTypeItemTrash                = "ITEM_TRASH"
	EventTypeItemUndeleteViaTrash     = "ITEM_UNDELETE_VIA_TRASH"
	EventTypeCollabAddCollaborator    = "COLLAB_ADD_COLLABORATOR"
	EventTypeCollabRoleChange         = "COLLAB_ROLE_CHANGE"
	EventTypeCollabInviteCollaborator = "COLLAB_INVITE_COLLABORATOR"
	EventTypeCollabRemoveCollaborator = "COLLAB_REMOVE_COLLABORATOR"
	EventTypeItemSync                 = "ITEM_SYNC"
	EventTypeItemUnsync               = "ITEM_UNSYNC"
	EventTypeItemRename               = "ITEM_RENAME"
	EventTypeItemSharedCreate         = "ITEM_SHARED_CREATE"
	EventTypeItemSharedUnshare        = "ITEM_SHARED_UNSHARE"
	EventTypeItemShared               = "ITEM_SHARED"
	EventTypeItemMakeCurrentVersion   = "ITEM_MAKE_CURRENT_VERSION"
	EventTypeTagItemCreate            = "TAG_ITEM_CREATE"
	EventTypeEnableTwoFactorAuth      = "ENABLE_TWO_FACTOR_AUTH"
	EventTypeMasterInviteAccept       = "MASTER_INVITE_ACCEPT"
	EventTypeMasterInviteReject       = "MASTER_INVITE_REJECT"
	EventTypeAccessGranted            = "ACCESS_GRANTED"
	EventTypeAccessRevoked            = "ACCESS_REVOKED"
	EventTypeGroupAddUser             = "GROUP_ADD_USER"
	EventTypeGroupRemoveUser          = "GROUP_REMOVE_USER"

	StreamPositionNow = "now"
)

func (c *client) GetEvents(streamPosition string) (*EventCollection, error) {
	body, err := c.Get("/events?stream_type=all&stream_position=" +
		url.QueryEscape(streamPosition))
	if err != nil {
		return nil, err
	}
	var events EventCollection
	err = json.Unmarshal(body, &events)
	if err != nil {
		return nil, err
	}
	return &events, nil
}
