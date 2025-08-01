package handlers


type Handlers struct {
	Note *NoteHandler
}

func NewHandlers(note *NoteHandler) *Handlers {
	return &Handlers{
		Note: note,
	}
}
