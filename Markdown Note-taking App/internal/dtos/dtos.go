package dtos

import "github.com/freitasmatheusrn/markdown-note-taking-app/internal/entity"

type SaveNoteInputDTO struct {
	Name    string
	Content []byte
}

type NoteOutputDTO struct {
	ID        int
	Name      string
	Content   []byte
	HasErrors bool
	Matches   []entity.Match
}
