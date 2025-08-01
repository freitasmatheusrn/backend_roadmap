package usecase

import (
	"github.com/freitasmatheusrn/markdown-note-taking-app/internal/dtos"
	"github.com/freitasmatheusrn/markdown-note-taking-app/internal/repositories"
)

type GetNoteUseCase struct {
	NoteRepository repositories.NoteRepositoryInterface
}

func NewGetNoteUseCase(noteRepository repositories.NoteRepositoryInterface) *GetNoteUseCase {
	return &GetNoteUseCase{
		NoteRepository: noteRepository,
	}
}

func (uc *GetNoteUseCase) Execute(id int) (*dtos.NoteOutputDTO, error) {
	note, err := uc.NoteRepository.GetNote(id)
	if err != nil {
		return nil, err
	}
	dto := &dtos.NoteOutputDTO{
		ID:        note.ID,
		Name:      note.Name,
		Content:   note.Content,
		HasErrors: note.HasErrors,
		Matches:   note.Matches,
	}
	return dto, nil
}
