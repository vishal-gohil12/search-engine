# Go Search Engine Research Prototype

This repository is a research-oriented prototype of core search engine components implemented in Go.
It is designed to show practical understanding of information retrieval basics and web crawling fundamentals, not to be a production search platform.

## Executive Summary

This project explores the first stages of a search engine pipeline:

- collecting pages from the web (crawler),
- processing text (normalization + stemming),
- building an inverted index,
- and scoring terms with TF-IDF.

For HR reviewers, this demonstrates problem solving in backend systems, algorithms, and data processing.
For engineering reviewers, the code shows working implementations of URL normalization, robots.txt checks, tokenization, stemming, postings-list indexing, and TF-IDF math.

## Research Goal

The main research question is:
"How can a minimal Go codebase implement the essential building blocks of search, from crawl to relevance scoring?"

The current implementation focuses on correctness of core ideas before optimization.

## What Is Implemented

The project currently has two active tracks:

- `crawler.go`: depth-limited same-domain crawler with robots.txt compliance checks.
- `main.go`: text processing, inverted index creation, term lookup, and TF-IDF scoring utilities.

These are implemented in one repository as a learning/research foundation; they are not yet fully integrated into one end-to-end pipeline.

## Architecture At A Glance

```text
Web Page URL
   -> crawler (robots-aware, depth-limited, same-host filter)
   -> extracted links

Text Documents (currently independent input path)
   -> tokenize + normalize + stem
   -> inverted index (token -> []docID)
   -> search(term) and TF-IDF(term, docID)
```

## Technical Details

### 1. Crawling (`crawler.go`)

- Starts from `https://go.dev`.
- Fetches and parses `robots.txt` using `github.com/temoto/robotstxt`.
- Uses recursive depth-limited crawl (`depth` currently set to 2).
- Extracts links from `<a href="...">` tags via `golang.org/x/net/html` tokenizer.
- Normalizes links with `net/url` (`ResolveReference`, fragment removed).
- Restricts traversal to same host.
- Adds a polite delay (`500ms`) between page visits.

### 2. Text Processing (`main.go`)

`tokenize(content string)` performs:

- lowercase conversion,
- punctuation cleanup with regex `[^\w\s]`,
- whitespace tokenization,
- English stemming via `github.com/kljensen/snowball`.

### 3. Inverted Index (`main.go`)

`Add(documents []Document)` builds:

- `map[string][]int` where key is token and value is posting list of document IDs.
- Repeated terms are stored as repeated doc IDs, enabling term-frequency counting later.

### 4. Search (`main.go`)

`Search(item string) []int`:

- lowercases the query term,
- returns matching posting list,
- returns empty slice if token is not present.

### 5. TF-IDF (`main.go`)

`TF_IDF(token string, docID int) float64` computes:

- `TF = termCountInDoc / totalTermsInDoc`
- `IDF = ln(totalDocs / docsContainingTerm)`
- `score = TF * IDF`

Guard clauses return `0` when data is missing or invalid.

## How To Run

From the `search-engine` folder:

```bash
go run .
```

Current executable entrypoint is in `crawler.go` (`main()`).
Typical current output begins with:

```text
Crawling (Depth 2): https://go.dev
```

## Code Map

- `crawler.go`: crawler entrypoint + crawling helpers.
- `main.go`: document model + indexing/retrieval/scoring logic.
- `go.mod`: module/dependency definitions.



