package services

import (
	"sort"
	"strings"

	"github.com/c-bata/go-prompt"
)

func services(current []string) []string {

	s, e := dc.Services()
	if e != nil {
		newError(e)
		return nil
	}

	var keys []string
	for k := range s {
		exist := false
		for _, c := range current {
			if c == k {
				exist = true
			}
		}
		if !exist {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	return keys
}

func suggest(current []string) []prompt.Suggest {

	var sug []prompt.Suggest
	for _, k := range services(current) {
		sug = append(sug, prompt.Suggest{
			Text:        k,
			Description: strings.Title(strings.ReplaceAll(k, "-", " ")),
		})
	}

	return sug
}

func suggestStartWith(current []string, prefix string) []prompt.Suggest {

	var sug []prompt.Suggest
	for _, k := range services(current) {
		if strings.HasPrefix(k, prefix) {
			sug = append(sug, prompt.Suggest{
				Text:        k,
				Description: strings.Title(strings.ReplaceAll(k, "-", " ")),
			})
		}
	}

	return sug
}
