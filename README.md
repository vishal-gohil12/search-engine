# Search Engine in Go (Inverted Index + TF-IDF)

This project is a simple learning implementation of a search engine core in Go.
It includes:

- tokenization and normalization
- inverted index creation
- single-term search
- TF-IDF scoring for a term in a specific document

## Project Structure

- `main.go`: all logic (document model, indexing, search, TF-IDF, demo in `main`)
- `go.mod`: Go module definition (`search-engine`)

## How It Works

### 1. Tokenization

`tokenize(content string) []string`:

- converts text to lowercase
- removes punctuation with regex: `[^\w\s]`
- splits text into tokens using whitespace

Example:
`"The quick brown fox."` -> `["the", "quick", "brown", "fox"]`

### 2. Inverted Index

`Add(documents []Document)` builds:

- `map[string][]int`
- key = token
- value = list of document IDs where the token appears

Important detail:
If a token appears multiple times in a document, that document ID appears multiple times in the postings list.  
This is used by `TF_IDF` to count term frequency.

### 3. Search

`Search(item string) []int`:

- lowercases query token
- returns postings list from index
- returns empty slice if not found

### 4. TF-IDF

`TF_IDF(token string, docID int) float64` computes:

- **TF (Term Frequency)** = occurrences of `token` in `docID` / total terms in `docID`
- **IDF (Inverse Document Frequency)** = `ln(totalDocs / docsContainingToken)`
- **TF-IDF** = `TF * IDF`

Implementation details:

- safely returns `0` for missing/invalid cases
- counts unique docs using a set
- counts docs containing token using a set built from token postings

## Run

From the `search-engine` folder:

```bash
go run .
```

## Current Demo Documents

Defined in `main.go`:

1. `The quick brown fox jumps over the lazy dog.`
2. `The lazy dog is sleeping.`
3. `The fox is quick and clever.`

## Example Output (Search)

The current program prints:

- full inverted index
- search results for:
  - `the`
  - `fox`
  - `dog`

## Notes and Limitations

- Search is currently single-token only.
- No ranking pipeline yet (TF-IDF is implemented, but not used to sort search results).
- No stop-word removal (terms like `the`, `is` remain indexed).
- No stemming/lemmatization in current code path.
- No unit tests yet.

## Suggested Next Improvements

1. Rank search results by TF-IDF score.
2. Add multi-term query handling.
3. Add stop-word filtering and optional stemming.
4. Add tests for tokenization, index correctness, and TF-IDF math.
5. Split logic into packages (`index`, `search`, `scoring`) for cleaner structure.
