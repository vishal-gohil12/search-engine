package main

import (
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/kljensen/snowball"
)

type Document struct {
	ID      int
	Content string
}

type InvertedIndex struct {
	Index map[string][]int
}

func tokenize(content string) []string {
	text := strings.ToLower(content)
	regx := regexp.MustCompile(`[^\w\s]`)
	text = regx.ReplaceAllString(text, " ")

	steammedTokens := []string{}
	for _, token := range strings.Fields(text) {
		stemmed, err := snowball.Stem(token, "english", true)
		if err == nil {
			steammedTokens = append(steammedTokens, stemmed)
		} else {
			steammedTokens = append(steammedTokens, token)
		}
	}

	return steammedTokens
}

func (ii *InvertedIndex) Add(documents []Document) map[string][]int {
	tokens := make(map[string][]int)
	for _, doc := range documents {
		tokensInDoc := tokenize(doc.Content)
		for _, token := range tokensInDoc {
			tokens[token] = append(tokens[token], doc.ID)
		}
	}
	ii.Index = tokens
	return tokens
}

func (ii *InvertedIndex) Search(item string) []int {
	token := strings.ToLower(item)
	if docIDs, found := ii.Index[token]; found {
		return docIDs
	}
	return []int{}
}

func (ii *InvertedIndex) TF_IDF(token string, docID int) float64 {
	if ii == nil || ii.Index == nil {
		return 0
	}

	token = strings.ToLower(token)

	postings, found := ii.Index[token]
	if !found {
		return 0
	}

	termCountInDoc := 0
	for _, id := range postings {
		if id == docID {
			termCountInDoc++
		}
	}
	if termCountInDoc == 0 {
		return 0
	}

	totalTermsInDoc := 0
	docsSet := make(map[int]struct{})
	for _, ids := range ii.Index {
		for _, id := range ids {
			docsSet[id] = struct{}{}
			if id == docID {
				totalTermsInDoc++
			}
		}
	}
	if totalTermsInDoc == 0 {
		return 0
	}

	docFreqSet := make(map[int]struct{})
	for _, id := range postings {
		docFreqSet[id] = struct{}{}
	}

	tf := float64(termCountInDoc) / float64(totalTermsInDoc)
	idf := math.Log(float64(len(docsSet)) / float64(len(docFreqSet)))

	return tf * idf
}

func main() {
	documents := []Document{
		{ID: 1, Content: "The quick brown fox jumps over the lazy dog."},
		{ID: 2, Content: "The lazy dog is sleeping."},
		{ID: 3, Content: "The fox is quick and clever."},
		{ID: 4, Content: "Dogs are loyal animals."},
		{ID: 5, Content: "Foxes are wild animals."},
	}
	ii := &InvertedIndex{}
	index := ii.Add(documents)
	fmt.Println("Inverted Index:")
	for token, docIDs := range index {
		fmt.Printf("%s: %v\n", token, docIDs)
	}

	ii.Search("the")
	fmt.Println("\nSearch results for 'the':", ii.Search("the"))
	ii.Search("fox")
	fmt.Println("Search results for 'fox':", ii.Search("fox"))
	ii.Search("dog")
	fmt.Println("Search results for 'dog':", ii.Search("dog"))

	fmt.Printf("\nTF-IDF for 'fox' in Document 1: %.4f\n", ii.TF_IDF("fox", 1))
	fmt.Printf("TF-IDF for 'lazy' in Document 2: %.4f\n", ii.TF_IDF("lazy", 2))
	fmt.Printf("TF-IDF for 'quick' in Document 3: %.4f\n", ii.TF_IDF("quick", 3))
	fmt.Printf("TF-IDF for 'wild' in Document 5: %.4f\n", ii.TF_IDF("wild", 5))

}
