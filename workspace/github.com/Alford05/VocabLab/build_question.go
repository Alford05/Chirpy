package main

import (
	"math/rand/v2"
)

func buildQuestions(cards []Card) []Question {
	rand.Shuffle(len(cards), func(i, j int) { cards[i], cards[j] = cards[j], cards[i] })
	terms := make([]string, len(cards))
	for i, c := range cards {
		terms[i] = c.Term
	}

	qs := make([]Question, 0, len(cards))
	for _, c := range cards {
		pool := make([]string, 0, len(terms)-1)
		for _, t := range terms {
			if t != c.Term {
				pool = append(pool, t)
			}
		}

		rand.Shuffle(len(pool), func(i, j int) { pool[i], pool[j] = pool[j], pool[i] })
		n := 3
		if len(pool) < 3 {
			n = len(pool)
		}

		choices := append([]string{c.Term}, pool[:n]...)
		rand.Shuffle(len(choices), func(i, j int) { choices[i], choices[j] = choices[j], choices[i] })

		ans := 0
		for idx, t := range choices {
			if t == c.Term {
				ans = idx
				break
			}
		}

		qs = append(qs, Question{
			Definition:  c.Definition,
			Choices:     choices,
			AnswerIdx:   ans,
			CorrectTerm: c.Term,
		})
	}
	return qs
}
