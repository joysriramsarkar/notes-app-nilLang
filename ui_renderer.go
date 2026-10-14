package main

import (
	"fmt"
	"strings"
)

// Analytics metrics
type Analytics struct {
	TotalNotes     int `json:"totalNotes"`
	ActiveNotes    int `json:"activeNotes"`
	PinnedNotes    int `json:"pinnedNotes"`
	FavoriteNotes  int `json:"favoriteNotes"`
	ArchivedNotes  int `json:"archivedNotes"`
	TrashedNotes   int `json:"trashedNotes"`
	TotalWords     int `json:"totalWords"`
	TotalChars     int `json:"totalChars"`
	TotalTodos     int `json:"totalTodos"`
	CompletedTodos int `json:"completedTodos"`
	TodoProgress   int `json:"todoProgress"`
}

func CalculateAnalytics(notes []Note) Analytics {
	a := Analytics{
		TotalNotes: len(notes),
	}

	for _, n := range notes {
		if n.IsTrashed {
			a.TrashedNotes++
		} else if n.IsArchived {
			a.ArchivedNotes++
		} else {
			a.ActiveNotes++
			if n.IsPinned {
				a.PinnedNotes++
			}
			if n.IsFavorite {
				a.FavoriteNotes++
			}
			w := len(strings.Fields(n.Title)) + len(strings.Fields(n.Content))
			a.TotalWords += w
			a.TotalChars += len(n.Title) + len(n.Content)

			for _, c := range n.Checklist {
				a.TotalTodos++
				if c.Completed {
					a.CompletedTodos++
				}
			}
		}
	}

	if a.TotalTodos > 0 {
		a.TodoProgress = (a.CompletedTodos * 100) / a.TotalTodos
	}
	return a
}

func RenderASCIITree(notes []Note) string {
	var sb strings.Builder
	sb.WriteString("┌─────────────────────────────────────────────────────────┐\n")
	sb.WriteString("│       📝 NILANG NOTES APP — DESKTOP RUNNER v1.0         │\n")
	sb.WriteString("└─────────────────────────────────────────────────────────┘\n\n")

	for i, n := range notes {
		if n.IsTrashed || n.IsArchived {
			continue
		}
		pin := ""
		if n.IsPinned {
			pin = " 📌"
		}
		sb.WriteString(fmt.Sprintf(" [%d] %s%s\n", i+1, n.Title, pin))
		sb.WriteString(fmt.Sprintf("     📁 %s | 🎨 %s | 🏷️ %s\n", n.Category, n.Color, strings.Join(n.Tags, ", ")))
		if len(n.Checklist) > 0 {
			done := 0
			for _, c := range n.Checklist {
				if c.Completed {
					done++
				}
			}
			sb.WriteString(fmt.Sprintf("     ✅ Tasks: %d/%d completed\n", done, len(n.Checklist)))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}
