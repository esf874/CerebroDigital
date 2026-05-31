package domain

import (
	"time"
	"strings"
	"reflect"
)

// NoteStatus representa el estado o progreso de una nota
type NoteStatus string

// NotePriority representa la prioridad de una nota
type NotePriority string

const (
	Pending  NoteStatus = "pending"
	Progress NoteStatus = "in_progress"
	Finished NoteStatus = "finished"
)

const (
	Low    NotePriority = "low"
	Medium NotePriority = "medium"
	High   NotePriority = "high"
)

type Note struct {
	ID      string `bson:"_id,omitempty"`
	Title   string `bson:"title"`
	Theme   string `bson:"theme,omitempty"`
	Content string `bson:"content"`

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`

	Status   NoteStatus   `bson:"status"`
	Priority NotePriority `bson:"priority,omitempty"`

	Tags []string `bson:"tags,omitempty"`
}

// NewNote crea una nueva nota inicializando sus valores por defecto y validando los campos obligatorios
func NewNote(id string, title, content string) (*Note, error) {
	if id == "" {
		return nil, ErrInvalidNoteID
	}

	if title == "" {
		return nil, ErrEmptyTitle
	}

	now := time.Now()
	return &Note{
		ID:        id,
		Title:     title,
		Content:   content,
		Status:    Pending, 
		CreatedAt: now,
		UpdatedAt: now,
		Tags:      []string{},
	}, nil
}


func (n *Note) AddTag(tag string) error {

	tag = strings.TrimSpace(tag)

	if tag == "" {
		return ErrEmptyTag
	}
	// Verificar si ya existe para evitar duplicados
	for _, t := range n.Tags {
		if t == tag {
			return ErrTagAlreadyExists
		}
	}
	n.Tags = append(n.Tags, tag)
	n.UpdatedAt = time.Now()
	return nil
}

func (n *Note) RemoveTag(value string) error {
	for i, tag := range n.Tags {
		if tag == value {
			n.Tags = append(n.Tags[:i], n.Tags[i+1:]...)
			n.UpdatedAt = time.Now()
			return nil
		}
	}
	return ErrTagNotFound
}

// NoteUpdate representa los campos que son actualizables, todos son opcionales para permitir actualizaciones parciales.
type NoteUpdate struct {
	Title    *string
	Content  *string
	Theme    *string
	Status   *NoteStatus
	Priority *NotePriority
	Tags     *[]string
}

func (n *Note) ApplyUpdate(update NoteUpdate) (bool, error) {
	cambiado := false

	if update.Title != nil {
		if *update.Title == "" {
			return false, ErrEmptyTitle
		}
		if n.Title != *update.Title {
			n.Title = *update.Title
			cambiado = true
		}
	}

	if update.Content != nil {
		if n.Content != *update.Content {
			n.Content = *update.Content
			cambiado = true
		}
	}

	if update.Theme != nil {
		if n.Theme != *update.Theme {
			n.Theme = *update.Theme
			cambiado = true
		}
	}

	if update.Status != nil {
		switch *update.Status {
		case Pending, Progress, Finished:
			if n.Status != *update.Status {
				n.Status = *update.Status
				cambiado = true
			}
		default:
			return false, ErrInvalidStatus
		}
	}

	if update.Priority != nil {
		switch *update.Priority {
		case Low, Medium, High:
			if n.Priority != *update.Priority {
				n.Priority = *update.Priority
				cambiado = true
			}
		default:
			return false, ErrInvalidPriority
		}
	}

	if update.Tags != nil {
		if !reflect.DeepEqual(n.Tags, *update.Tags) {
			// Reemplaza completamente los tags
			n.Tags = *update.Tags
			cambiado = true
		}
	}

	if cambiado {
		n.UpdatedAt = time.Now()
	}
	return cambiado, nil
}
