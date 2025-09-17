package main

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Card struct {
	Term       string
	Definition string
}

type Question struct {
	Definition  string
	Choices     []string //4 terms, 1 correct answer
	AnswerIdx   int
	CorrectTerm string
}

type Quiz struct {
	Questions []Question
	Current   int
	Score     int
}

func main() {
	newApp := app.New()
	newWindow := newApp.NewWindow("Vocab Trainer")
	newWindow.Resize(fyne.NewSize(600, 900))

	wordEntry := widget.NewMultiLineEntry()
	wordEntry.SetPlaceHolder("Enter one word per line, minimum 10 words")

	startButton := widget.NewButton("Start Quiz", func() {
		words := getWordList(wordEntry.Text)
		if len(words) < 10 {
			wordEntry.SetText("Please enter at least 10 words.")
			return
		}

		cards, err := fetchDefinitions(words)
		if err != nil {
			wordEntry.SetText("Error: " + err.Error())
			return
		}

		// Build Quiz
		questions := buildQuestions(cards)
		q := Quiz{Questions: questions, Current: 0, Score: 0}
		newWindow.SetContent(makeQuizView(newWindow, &q))
	})

	content := container.NewVBox(
		widget.NewLabel("Enter Vocabulary Words:"),
		wordEntry,
		startButton,
	)
	newWindow.SetContent(content)
	newWindow.ShowAndRun()
}

func getWordList(input string) []string {
	lines := strings.Split(input, "\n")
	var words []string
	seen := make(map[string]bool)
	for _, line := range lines {
		w := strings.TrimSpace(line)
		if w != "" && !seen[w] {
			seen[w] = true
			words = append(words, w)
		}
	}
	return words
}
