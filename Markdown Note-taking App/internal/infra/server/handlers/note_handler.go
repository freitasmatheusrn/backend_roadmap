package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/freitasmatheusrn/markdown-note-taking-app/internal/dtos"
	"github.com/freitasmatheusrn/markdown-note-taking-app/internal/infra/web/views"
	"github.com/freitasmatheusrn/markdown-note-taking-app/internal/repositories"
	"github.com/freitasmatheusrn/markdown-note-taking-app/internal/usecase"
	"github.com/freitasmatheusrn/markdown-note-taking-app/pkg"
	"github.com/go-chi/chi/v5"
)

type NoteHandler struct {
	NoteRepository repositories.NoteRepositoryInterface
}

func NewNoteHandler(NoteRepository repositories.NoteRepositoryInterface) *NoteHandler {
	return &NoteHandler{
		NoteRepository: NoteRepository,
	}
}

func (h *NoteHandler) Home(w http.ResponseWriter, r *http.Request){
	uc := usecase.NewGetListNotesUseCase(h.NoteRepository)
	output, err := uc.Execute()
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	ctx := context.Background()
	err = views.Home(output).Render(ctx, w)
	if err != nil {
		http.Error(w, "Could not render page", http.StatusInternalServerError)
	}

}

func (h *NoteHandler) UploadNote(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Could not parse multipart form", http.StatusBadRequest)
		return
	}
	file, _, err := r.FormFile("content")
	if err != nil {
		http.Error(w, "Could not get uploaded file", http.StatusBadRequest)
		return
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Could not read file content", http.StatusInternalServerError)
		return
	}
	content = pkg.MdToHTML(content)
	input := dtos.SaveNoteInputDTO{
		Name:    r.FormValue("name"),
		Content: content,
	}
	uc := usecase.NewSaveNoteUseCase(h.NoteRepository)
	output, err := uc.Execute(input)
	if err != nil {
		http.Error(w, "Could not save file", http.StatusInternalServerError)
	}
	http.Redirect(w, r, fmt.Sprintf("/notes/%d", output.ID), http.StatusSeeOther)

}
func (h *NoteHandler) ShowNote(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	fmt.Println(idStr)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	uc := usecase.NewGetNoteUseCase(h.NoteRepository)
	output, err := uc.Execute(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := context.Background()
	err = views.ShowNote(output).Render(ctx, w)
	if err != nil {
		http.Error(w, "Could not render page", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}
