package usecase

import (
	"fmt"

	"github.com/freitasmatheusrn/markdown-note-taking-app/internal/dtos"
	"github.com/freitasmatheusrn/markdown-note-taking-app/internal/entity"
	"github.com/freitasmatheusrn/markdown-note-taking-app/internal/repositories"
	"github.com/freitasmatheusrn/markdown-note-taking-app/pkg"
)

type SaveNoteUseCase struct {
	NoteRepository repositories.NoteRepositoryInterface
}

func NewSaveNoteUseCase(noteRepository repositories.NoteRepositoryInterface) *SaveNoteUseCase {
	return &SaveNoteUseCase{
		NoteRepository: noteRepository,
	}
}

func (uc *SaveNoteUseCase) Execute(input dtos.SaveNoteInputDTO) (*dtos.NoteOutputDTO, error) {
	note := entity.NewNote(input.Name, input.Content)
	pkg.CheckGrammar(note)
	note, err := uc.NoteRepository.Save(note)
	fmt.Println(note)
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
