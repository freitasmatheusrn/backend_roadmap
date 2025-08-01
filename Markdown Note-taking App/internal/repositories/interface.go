package repositories

import "github.com/freitasmatheusrn/markdown-note-taking-app/internal/entity"


type NoteRepositoryInterface interface {
	Save(*entity.Note) (*entity.Note, error)
	GetNote(id int) (*entity.Note, error)
	ListAll() ([]*entity.Note, error) 
}
