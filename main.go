package main

import (
	"Documents/Project/Ebookie/config"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pirmd/epub"
	"github.com/pkg/errors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Sets the style
type Styles struct {
	cursor      lipgloss.Style
	choices     lipgloss.Style
	highlighted lipgloss.Style
}

func initialModel() model {
	cfg, err := config.ParseConfig()
	if err != nil {
		log.Fatal(err)
	}

	ebookDir := cfg.Settings.EbookDir

	return model{
		title:       find(ebookDir, ".epub"),
		choices:     listEpubs(ebookDir),
		cursor:      ">",
		min:         0,
		max:         0,
		Height:      0,
		highlighted: 0,
		Styles:      DefaultStyles(),
	}
}

// DefaultStyles defines the default styling for the file picker.
func DefaultStyles() Styles {
	return DefaultStylesWithRenderer(lipgloss.DefaultRenderer())
}

func DefaultStylesWithRenderer(r *lipgloss.Renderer) Styles {
	return Styles{
		cursor:      r.NewStyle().Foreground(lipgloss.Color("212")),
		choices:     r.NewStyle(),
		highlighted: r.NewStyle().Foreground(lipgloss.Color("500")).Bold(true),
	}
}

func find(root, ext string) []string {
	var filename []string
	filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ext {
			filename = append(filename, s)
		}
		return nil
	})
	return filename
}

func listEpubs(directory string) []string {
	var sr []string
	for _, sr2 := range find(directory, ".epub") {
		sr2, err := epub.GetMetadataFromFile(sr2)
		if err != nil {
			errors.Cause(err)
		}
		sr = append(sr, sr2.Title...)
	}
	return sr
}

// MAIN MODEL
type model struct {
	// epub title to be displayed
	title []string

	// directory of the file
	choices []string

	cursor string // Which item is pointed out

	highlighted int

	min    int
	max    int
	Height int

	Styles Styles
}

// Runs on start up
func (m model) Init() tea.Cmd {
	return nil
}

// UPDATE=handle incoming events
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Height = 20
		m.max = m.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			m.highlighted--
			if m.highlighted < 0 {
				m.highlighted = 0
			}
			if m.highlighted < m.min {
				m.min--
				m.max--
			}

		case "down", "j":
			m.highlighted++
			if m.highlighted >= len(m.choices) {
				m.highlighted = len(m.choices) - 1
			}
			if m.highlighted > m.max {
				m.min++
				m.max++
			}

			// Open selected file using foliate
		case "right", "l":
			if m.cursor != "" {
				go func() {
					err := exec.Command("foliate", m.title[m.highlighted]).Run()
					if err != nil {
						//         // Handle errors
						fmt.Println("Error opening Foliate:", err)
					}
				}()
			}

		}

	}
	return m, nil
}

// view
func (m model) View() string {
	var s strings.Builder

	for i, items := range m.choices {
		if i < m.min {
			continue
		}
		if i > m.max {
			break
		}

		if m.highlighted == i {
			highlighted := fmt.Sprint(m.Styles.highlighted.Render(items))
			s.WriteString(m.Styles.cursor.Render(m.cursor) + m.Styles.highlighted.Render(highlighted))
			s.WriteRune('\n')
			continue
		}
		s.WriteString(fmt.Sprintf("%s\n", m.Styles.choices.Render(items)))
		s.WriteRune('\n')

	}
	return s.String()
}

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatalf("err: %w", err)
	}
	defer f.Close()

	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		fmt.Print("Alas, there has been an error", err)
		os.Exit(1)
	}
}
