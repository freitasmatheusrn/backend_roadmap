package pkg

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/freitasmatheusrn/markdown-note-taking-app/internal/entity"
)

func CheckGrammar(note *entity.Note) error {
	text := string(note.Content)
	cleanText := CleanHTML(text)

	form := url.Values{}
	form.Set("text", cleanText)
	form.Set("language", "en-US")

	resp, err := http.PostForm("https://api.languagetool.org/v2/check", form)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return errors.New("erro na API LanguageTool: " + string(body))
	}

	var result struct {
		Matches []struct {
			Rule struct {
				IssueType string `json:"issueType"`
			} `json:"rule"`
			Context struct {
				Text string `json:"text"` 
			} `json:"context"`
			Sentence string `json:"sentence"` 
		} `json:"matches"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	var matches []entity.Match
	for _, m := range result.Matches {
		if m.Rule.IssueType == "misspelling" {
			matches = append(matches, entity.Match{
				Message:  "misspelling: " + m.Context.Text,
				Sentence: m.Sentence,
			})
		}
	}

	note.Matches = matches
	note.HasErrors = len(matches) > 0

	return nil
}
