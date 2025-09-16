package main

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func makeQuizView(win fyne.Window, quiz *Quiz) fyne.CanvasObject {
	defLabel := widget.NewLabel("")                 // shows definition
	scoreLabel := widget.NewLabel("Score: 0")       // shows score
	qIndexLabel := widget.NewLabel("Question 1/10") // shows progress

	// 4 choice buttons
	btns := []*widget.Button{
		widget.NewButton("", nil),
		widget.NewButton("", nil),
		widget.NewButton("", nil),
		widget.NewButton("", nil),
	}

	// click handler factory
	setHandler := func(i int) func() {
		return func() {
			if quiz.Current >= len(quiz.Questions) {
				return
			}
			q := &quiz.Questions[quiz.Current]
			correct := i == q.AnswerIdx
			if correct {
				quiz.Score++
				scoreLabel.SetText("Score: " + strconv.Itoa(quiz.Score))
			}
			// feedback
			result := "Incorrect! The correct answer is: " + q.CorrectTerm
			if correct {
				result = "Correct!"
			}

			// Advance the Quiz
			quiz.Current++
			if quiz.Current >= len(quiz.Questions) {
				// Quiz over
				dialog.ShowInformation("Quiz Complete", result+"\nFinal score: "+strconv.Itoa(quiz.Score)+"/"+strconv.Itoa(len(quiz.Questions)), win)
				end := container.NewVBox(
					widget.NewLabel("Quiz Complete!"),
					widget.NewLabel("Score: "+strconv.Itoa(quiz.Score)+"/"+strconv.Itoa(len(quiz.Questions))),
					widget.NewButton("Close", func() { win.Close() }),
				)
				win.SetContent(end)
				return
			}

			// Show feedback then next question
			dialog.ShowInformation("Result", result, win)
			renderQuestion(quiz, defLabel, qIndexLabel, btns)
		}
	}

	// attach handlers
	for i := range btns {
		idx := i
		btns[i].OnTapped = setHandler(idx)
	}
	// layout
	choicesGrid := container.NewGridWithRows(4, btns[0], btns[1], btns[2], btns[3])
	root := container.NewVBox(
		qIndexLabel,
		scoreLabel,
		widget.NewSeparator(),
		widget.NewLabel("Definition:"),
		defLabel,
		widget.NewSeparator(),
		widget.NewLabel("Choose the correct term:"),
		choicesGrid,
	)

	// initial render
	if len(quiz.Questions) == 0 {
		return container.NewVBox(widget.NewLabel("No questions left to answer"))
	}
	renderQuestion(quiz, defLabel, qIndexLabel, btns)
	return root
}

func renderQuestion(quiz *Quiz, defLabel, qIndexLabel *widget.Label, btns []*widget.Button) {
	q := &quiz.Questions[quiz.Current]
	defLabel.SetText(q.Definition)
	qIndexLabel.SetText("Question " + strconv.Itoa(quiz.Current+1) + "/" + strconv.Itoa(len(quiz.Questions)))
	for i := range btns {
		if i < len(q.Choices) {
			btns[i].SetText(q.Choices[i])
			btns[i].Enable()
			continue
		}
		btns[i].SetText("")
		btns[i].Disable()
	}
}
