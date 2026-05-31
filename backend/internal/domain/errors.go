package domain

import "errors"

var (
	ErrInvalidID = errors.New("ID cannot be empty.")
 
	ErrNoteNotFound       = errors.New("Note not found.")
	ErrOriginNoteNotFound = errors.New("Origin note not found.")
	ErrInvalidNoteID      = errors.New("Invalid note ID.")
	ErrSameIds            = errors.New("Origin and destination note cannot be the same.")
	ErrNoteAlreadyExists  = errors.New("Note with the same ID already exists.")
	ErrEmptyTitle         = errors.New("The title cannot be empty.")
	ErrEmptyTag           = errors.New("Tag value cannot be empty.")
	ErrTagNotFound        = errors.New("Tag not found in the note.")
	ErrNilNote            = errors.New("Note cannot be nil.")
	ErrInvalidStatus      = errors.New("invalid status")
	ErrInvalidPriority    = errors.New("invalid priority")
	ErrDestNoteNotFound   = errors.New("Destination note not found.")

	ErrLinkNotFound      = errors.New("Link not found.")
	ErrInvalidLinkID     = errors.New("Invalid link ID.")
	ErrLinkAlreadyExists = errors.New("Link with the same ID already exists.")
	ErrNilLink           = errors.New("Link cannot be nil.")

	ErrTagAlreadyExists = errors.New("Tag already exists in the note.")
)
