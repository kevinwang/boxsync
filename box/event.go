package box

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
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

func (c *client) GetLongPollURL() (string, error) {
	body, err := c.Options("/events")
	if err != nil {
		return "", err
	}
	var resp LongPollURLResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return "", err
	}
	if resp.ChunkSize != 1 {
		return "", errors.New("Long poll chunk size is not 1")
	}
	return resp.Entries[0].URL, nil
}

func (c *client) GetEventStream(longPollURL, streamPosition string, quit <-chan struct{}) (<-chan Event, <-chan error, error) {
	eventStream := make(chan Event)
	errorStream := make(chan error)

	if streamPosition == "now" || streamPosition == "" {
		collection, err := c.GetEvents("now")
		if err != nil {
			return nil, nil, err
		}
		streamPosition = strconv.Itoa(collection.NextStreamPosition)
	}

	go func() {
		for {
			select {
			case <-quit:
				break
			default:
			}

			body, err := c.GetByURL(longPollURL + "&streamPosition=" + streamPosition)
			if err != nil {
				errorStream <- err
				break
			}
			var resp LongPollResponse
			err = json.Unmarshal(body, &resp)
			if err != nil {
				errorStream <- err
				break
			}
			if resp.Message != "new_change" {
				errorStream <- err
				break
			}
			collection, err := c.GetEvents(streamPosition)
			if err != nil {
				errorStream <- err
				break
			}
			for _, event := range collection.Entries {
				eventStream <- event
			}
			streamPosition = strconv.Itoa(collection.NextStreamPosition)
		}
	}()

	return eventStream, errorStream, nil
}
