package entity

type Note struct {
	ID        int
	Name      string
	Content   []byte
	HasErrors bool
	Matches   []Match
}

func NewNote(name string, content []byte) *Note {
	return &Note{
		Name:    name,
		Content: content,
	}
}
