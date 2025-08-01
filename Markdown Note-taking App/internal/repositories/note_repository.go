package repositories

import (
	"database/sql"

	"github.com/freitasmatheusrn/markdown-note-taking-app/internal/entity"
)

type NoteRepository struct {
	DB *sql.DB
}

func NewNoteRepository(db *sql.DB) *NoteRepository {
	return &NoteRepository{DB: db}
}

func (r *NoteRepository) Save(note *entity.Note) (*entity.Note, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	err = tx.QueryRow(`
		INSERT INTO notes (name, content, has_errors)
		VALUES ($1, $2, $3)
		RETURNING id
	`, note.Name, note.Content, note.HasErrors).Scan(&note.ID)
	if err != nil {
		return nil, err
	}
	for _, match := range note.Matches {
		_, err = tx.Exec(`
			INSERT INTO matches (note_id, message, sentence)
			VALUES ($1, $2, $3)
		`, note.ID, match.Message, match.Sentence)
		if err != nil {
			return nil, err
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return note, nil
}

func (r *NoteRepository) GetNote(id int) (*entity.Note, error) {
	query := `SELECT id, name, content, has_errors FROM notes WHERE id = $1`

	var n entity.Note
	err := r.DB.QueryRow(query, id).Scan(&n.ID, &n.Name, &n.Content, &n.HasErrors)
	if err != nil {
		return nil, err
	}
	matchesQuery := `SELECT message, sentence, note_id FROM matches WHERE note_id = $1`
	rows, err := r.DB.Query(matchesQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []entity.Match
	for rows.Next() {
		var match entity.Match
		err := rows.Scan(&match.Message, &match.Sentence, &match.NoteID)
		if err != nil {
			return nil, err
		}
		matches = append(matches, match)
	}

	n.Matches = matches
	return &n, nil
}

func (r *NoteRepository) ListAll() ([]*entity.Note, error) {
	rows, err := r.DB.Query(`
		SELECT
			n.id,
			n.name,
			n.content,
			n.has_errors,
			m.message,
			m.sentence
		FROM notes n
		LEFT JOIN matches m ON m.note_id = n.id
		ORDER BY n.id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []*entity.Note
	var currentNote *entity.Note
	var lastNoteID int

	for rows.Next() {
		var (
			noteID   int
			name     string
			content  []byte
			hasErr   bool
			message  sql.NullString
			sentence sql.NullString
		)

		if err := rows.Scan(&noteID, &name, &content, &hasErr, &message, &sentence); err != nil {
			return nil, err
		}

		if currentNote == nil || noteID != lastNoteID {
			currentNote = &entity.Note{
				ID:        noteID,
				Name:      name,
				Content:   content,
				HasErrors: hasErr,
				Matches:   []entity.Match{},
			}
			notes = append(notes, currentNote)
			lastNoteID = noteID
		}

		if message.Valid && sentence.Valid {
			currentNote.Matches = append(currentNote.Matches, entity.Match{
				Message:  message.String,
				Sentence: sentence.String,
			})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notes, nil
}
