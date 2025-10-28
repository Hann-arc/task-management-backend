package errors

import "errors"

var (
	ErrUnauthorizedBoardAction = errors.New("unauthorized: only project owner can manage boards")
	ErrInvalidOrderIndex       = errors.New("invalid order_index: must be between 1 and max+1")
	ErrNoFieldsToUpdate        = errors.New("no fields to update")
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidUserData = errors.New("invalid user data")
	ErrEmailExists     = errors.New("email already exists")
)

var (
	ErrProjectNotFound       = errors.New("project not found")
	ErrUnauthorizedProject   = errors.New("unauthorized: only owner or member can access this project")
	ErrUnauthorizedOwnerOnly = errors.New("unauthorized: only project owner can perform this action")
)

var (
	ErrTaskNotFound     = errors.New("task not found")
	ErrUnauthorizedTask = errors.New("unauthorized: only project members can manage tasks")
	ErrInvalidTaskData  = errors.New("invalid task data")
	ErrBoardNotFound    = errors.New("board not found")
	ErrAssigneeNotFound = errors.New("assignee not found")
)

var (
	ErrProjectMemberNotFound = errors.New("project member not found")
	ErrAlreadyMember         = errors.New("user is already a member of this project")
	ErrCannotRemoveSelf      = errors.New("you cannot remove yourself from the project")
	ErrInviteeIsOwner        = errors.New("cannot invite project owner")
)

var (
	ErrInvitationNotFound = errors.New("invitation not found")
	ErrInvitationExpired  = errors.New("invitation has expired")
	ErrInvitationUsed     = errors.New("invitation already used")
	ErrCannotInviteSelf   = errors.New("cannot invite yourself")
)

var (
	ErrCommentNotFound    = errors.New("comment not found")
	ErrCannotReplyToReply = errors.New("cannot reply to a reply directly")
)
