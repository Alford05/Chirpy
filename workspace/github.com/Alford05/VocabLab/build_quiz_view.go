package main

import (
	"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func makeQuizView(win fyne.Window, quiz *Quiz) fyne.CanvasObject {
	defLabel := widget.NewRichText() // shows definition
	defLabel.Wrapping = fyne.TextWrapWord
	// defLabel.Alignment = fyne.TextAlignLeading

	defScroll := container.NewVScroll(defLabel)
	defScroll.SetMinSize(fyne.NewSize(500, 150))
	defInner := container.NewPadded(defScroll)
	defCard := widget.NewCard("Definition", "", defInner)

	title := widget.NewLabel("Vocab Trainer")
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	total := strconv.Itoa(len(quiz.Questions))
	scoreLabel := widget.NewLabel("Score: 0")             // shows score
	qIndexLabel := widget.NewLabel("Question 1/" + total) // shows progress

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
			for _, b := range btns {
				b.Disable()
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
				dialog.ShowInformation("Quiz Complete", result+"\nFinal score: "+strconv.Itoa(quiz.Score)+"/"+total, win)
				end := container.NewVBox(
					widget.NewLabel("Quiz Complete!"),
					widget.NewLabel("Score: "+strconv.Itoa(quiz.Score)+"/"+total),
					widget.NewButton("New Quiz", func() {
						// rebuild the start screen
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
							questions := buildQuestions(cards)
							newQ := Quiz{Questions: questions}
							win.SetContent(makeQuizView(win, &newQ))
						})
						win.SetContent(container.NewVBox(
							widget.NewLabel("Enter Vocabulary Words:"),
							wordEntry,
							startButton,
						))
					}),
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

	root := container.NewPadded(container.NewVBox(
		title,
		widget.NewSeparator(),
		container.NewGridWithColumns(2, qIndexLabel, scoreLabel),
		widget.NewSeparator(),
		defCard,
		widget.NewSeparator(),
		widget.NewLabel("Choose the correct term:"),
		choicesGrid,
	))
	// soft sky blue background:
	bg := canvas.NewRectangle(color.NRGBA{R: 224, G: 242, B: 254, A: 255})
	stack := container.NewStack(bg, root)
	stack.Resize(fyne.NewSize(600, 900))
	bg.Resize(stack.Size())

	// initial render
	if len(quiz.Questions) == 0 {
		return container.NewVBox(widget.NewLabel("No questions left to answer"))
	}
	renderQuestion(quiz, defLabel, qIndexLabel, btns)
	return stack
}

func renderQuestion(quiz *Quiz, defLabel *widget.RichText, qIndexLabel *widget.Label, btns []*widget.Button) {
	q := &quiz.Questions[quiz.Current]
	defLabel.Segments = []widget.RichTextSegment{
		&widget.TextSegment{
			Text: q.Definition,
			Style: widget.RichTextStyle{
				ColorName: theme.ColorRed,
			},
		},
	}
	defLabel.Refresh()
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
