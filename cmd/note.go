package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"time"
)

func note(noteName string) error {
	notesDir := getNotesDir()
	if err := ensureNotesDir(notesDir); err != nil {
		return fmt.Errorf("failed to get notes directory: %v", err)
	}
	editor := getEditor()
	err := checkFzfInstalled()
	if err != nil {
		return err
	}

	notePath, err := getNotePath(noteName, notesDir)
	if err != nil {
		// If error message is "not selection cancelled", we exit gracefully
		if err.Error() == "note selection cancelled" {
			os.Exit(0)
		}
		return err
	}

	cmd := exec.Command(editor, notePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to open editor: %v", err)
	}

	return nil
}

func getMatches(noteName string) ([]string, error) {
	files, err := filepath.Glob(getNotesDir() + string(os.PathSeparator) + "*" + noteName + "*.md")
	if err != nil {
		return nil, fmt.Errorf("failed to read notes directory: %v", err)
	}

	return files, nil
}

type ByDate []string

func (a ByDate) Len() int      { return len(a) }
func (a ByDate) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool {
	dateI := a[i][:11]
	dateJ := a[j][:11]
	layout := "02-Jan-2006"
	timeI, errI := time.Parse(layout, dateI)
	timeJ, errJ := time.Parse(layout, dateJ)
	if errI != nil || errJ != nil {
		return a[i] < a[j] // Fallback to string comparison if parsing fails
	}
	return timeJ.Before(timeI)
}

// List all files in notes directory matching *.md
func getAllNotePaths(directory string) ([]string, error) {
	files, err := filepath.Glob(directory + string(os.PathSeparator) + "*.md")
	if err != nil {
		return nil, err
	}

	filepaths := make([]string, 0, len(files))
	for _, file := range files {
		filepaths = append(filepaths, filepath.Base(file))
	}
	sortedFiles := ByDate(filepaths)
	sort.Sort(sortedFiles)

	return sortedFiles, nil
}

func getNotePath(noteName string, notesDir string) (string, error) {
	var matches []string
	var err error
	if noteName == "" {
		matches, err = getAllNotePaths(notesDir)
	} else {
		matches, err = getMatches(noteName)
	}
	if err != nil {
		return "", err
	}

	if len(matches) == 1 {
		return matches[0], nil
	}

	// If multiple matches, the user selects one using fzf
	if len(matches) > 1 {
		fzfCmd := exec.Command("fzf", "--no-sort", "--prompt", "Select a note: ")

		// Create a pipe to send matches to fzf
		stdin, err := fzfCmd.StdinPipe()
		if err != nil {
			return "", fmt.Errorf("failed to create stdin pipe: %v", err)
		}

		// Send matches to fzf in a goroutine
		go func() {
			defer stdin.Close()
			for _, match := range matches {
				fmt.Fprintln(stdin, filepath.Base(match))
			}
		}()

		// Get output from fzf
		output, err := fzfCmd.Output()
		if err != nil {
			// If error is exit status 130, it means the user cancelled the selection
			// We handle this case gracefully
			if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() == 130 {
				return "", fmt.Errorf("note selection cancelled")
			}
			return "", fmt.Errorf("failed to select note: %v", err)
		}

		selectedNote := string(output)
		selectedNote = selectedNote[:len(selectedNote)-1] // Remove newline character

		// Find the full path for the selected note
		for _, match := range matches {
			if filepath.Base(match) == selectedNote {
				return match, nil
			}
		}
		return "", fmt.Errorf("selected note not found")
	}

	currentDate := time.Now().Format("02-Jan-2006")
	return getNotesDir() + string(os.PathSeparator) + currentDate + " " + noteName + ".md", nil
}

func ensureNotesDir(notesDir string) error {
	if _, err := os.Stat(notesDir); os.IsNotExist(err) {
		return os.MkdirAll(notesDir, 0755)
	}
	return nil
}

func getNotesDir() string {
	notesPath := os.Getenv("NOTES_DIRECTORY")
	if notesPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			panic("Could not determine home directory")
		}
		notesPath = homeDir + string(os.PathSeparator) + ".notes"
	}
	return notesPath
}

func getEditor() string {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		fmt.Println("WARNING: $EDITOR not set, defaulting to 'vi'")
		editor = "vi"
	}
	return editor
}

func checkFzfInstalled() error {
	_, err := exec.LookPath("fzf")
	if err != nil {
		return fmt.Errorf("fzf not found in PATH. Please install fzf to enable note searching")
	}
	return nil
}
