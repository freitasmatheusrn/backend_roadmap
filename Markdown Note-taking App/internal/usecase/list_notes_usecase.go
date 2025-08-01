package usecase

import (
	"github.com/freitasmatheusrn/markdown-note-taking-app/internal/dtos"
	"github.com/freitasmatheusrn/markdown-note-taking-app/internal/repositories"
)

type GetListNotesUseCase struct {
	NoteRepository repositories.NoteRepositoryInterface
}

func NewGetListNotesUseCase(noteRepository repositories.NoteRepositoryInterface) *GetListNotesUseCase {
	return &GetListNotesUseCase{
		NoteRepository: noteRepository,
	}
}

func (uc *GetListNotesUseCase) Execute() ([]dtos.NoteOutputDTO, error) {
	notes, err := uc.NoteRepository.ListAll()
	if err != nil {
		return nil, err
	}
	dtosSlice := make([]dtos.NoteOutputDTO, 0)
	for _, note := range notes{
		dto := dtos.NoteOutputDTO{
			ID:        note.ID,
			Name:      note.Name,
			Content:   note.Content,
			HasErrors: note.HasErrors,
			Matches:   note.Matches,
		}
		dtosSlice = append(dtosSlice, dto)

	}
	return dtosSlice, nil
}
