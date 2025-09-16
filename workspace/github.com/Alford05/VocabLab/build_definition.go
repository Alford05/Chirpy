package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type dictResp struct {
	term string
	def  string
	err  error
}

func fetchDefinitions(terms []string) ([]Card, error) {
	ch := make(chan dictResp)
	sem := make(chan struct{}, 10) // limit concurrency

	for _, t := range terms {
		t := t
		go func() {
			sem <- struct{}{}
			def, err := fetchOne(t)
			<-sem
			ch <- dictResp{term: t, def: def, err: err}
		}()
	}

	var cards []Card
	var errs []error
	for range terms {
		r := <-ch
		if r.err == nil && r.def != "" {
			cards = append(cards, Card{Term: r.term, Definition: r.def})
		} else {
			if r.err != nil {
				errs = append(errs, r.err)
			}
		}
	}

	if len(cards) < 10 {
		return nil, fmt.Errorf("need at least 10 definitions, got %d", len(cards))
	}
	return cards, nil
}

func fetchOne(term string) (string, error) {
	u := "http://api.dictionaryapi.dev/api/v2/entries/en/" + url.QueryEscape(term)
	resp, err := http.Get(u)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("no entry for %s", term, resp.StatusCode)
	}

	var data []struct {
		Meanings []struct {
			Definitions []struct {
				Definition string `json:"definition"`
			} `json:"definitions"`
		} `json:"meanings"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	for _, entry := range data {
		for _, m := range entry.Meanings {
			for _, d := range m.Definitions {
				if d.Definition != "" {
					return d.Definition, nil
				}
			}
		}
	}
	return "", fmt.Errorf("no definition for %s", term)
}
