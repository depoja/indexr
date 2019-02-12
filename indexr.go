package indexr

// TODO: Update map usages to slices

// IndexEntry includes the resulting documents and the match occurrences
type IndexEntry struct {
	docs        Docs        // the docs this n-gram entry occurs
	occurrences map[int]int // the occurrences for this n-gram (deleted when 0)
}

// Index contains the definition of the n-gram index
type Index struct {
	n       int                  // n-gram type (2, 3, 4 etc...)
	entries map[Ngram]IndexEntry // the index entries
	docs    map[int]bool         // all the document ids
}

// New creates a new Index
func New(n int) *Index {
	return &Index{
		n:       n,
		docs:    map[int]bool{},
		entries: map[Ngram]IndexEntry{},
	}
}

// Add adds a new document to the Index
func (i *Index) Add(id int, str string) int {
	ngrams := getNgrams(str, i.n)

	for _, ngram := range ngrams {
		entry, exists := i.entries[ngram]

		// Index entry does not exists for this n-gram
		if !exists {
			entry = IndexEntry{}
			entry.docs = append(entry.docs, id)
			entry.occurrences = map[int]int{}
			entry.occurrences[id] = 1
		}

		// Index entry already exists for this n-gram
		if exists {
			if _, occurred := entry.occurrences[id]; occurred {
				entry.occurrences[id] = entry.occurrences[id] + 1
			} else {
				entry.docs = append(entry.docs, id)
				entry.occurrences[id] = 1
			}
		}
		// Save the index entry
		i.entries[ngram] = entry
	}

	i.docs[id] = true
	return id
}

// Delete removes a document from the Index
func (i *Index) Delete(id int, str string) {
	ngrams := getNgrams(str, i.n)

	for _, ngram := range ngrams {
		entry, exists := i.entries[ngram]

		// Index entry does not exists for this n-gram, nothing to delete
		if !exists {
			return
		}

		// Index entry already exists for this n-gram
		if exists {
			occurrences, occurred := entry.occurrences[id]

			// If multiple occurrences decrement occurrence, otherwise remove
			if occurred && occurrences > 1 {
				entry.occurrences[id] = entry.occurrences[id] - 1
			} else {
				delete(entry.occurrences, id)
				for i, v := range entry.docs {
					if v == id {
						entry.docs = append(entry.docs[:i], entry.docs[i+1:]...)
					}
				}
			}

			// If no documents present for the index entry remove it, otherwise update it
			if len(entry.docs) == 0 {
				delete(i.entries, ngram) // TODO: Also remove docs with no ngrams
			} else {
				i.entries[ngram] = entry
			}
		}
	}
}

// Query searches for a given string and returns the matching document ids
func (i *Index) Query(str string) Docs {
	// If no n-grams for the string, return all docs
	ngrams := getNgrams(str, i.n)
	if len(ngrams) == 0 {
		return i.getDocs()
	}

	// Get the first n-gram as an intersection base (get the documents matching at least the first n-gram)
	entry, exists := i.entries[ngrams[0]]
	if !exists {
		return nil
	}
	result := entry.docs

	// For the rest of the n-grams intersect their results, returning the entries that contain them
	ngrams = ngrams[1:]
	for _, ngram := range ngrams {
		entry, exists := i.entries[ngram]
		if !exists {
			return nil
		}
		result = intersect(result, entry.docs)
	}

	// Sort the entries
	return sortDocs(result)
}

func (i *Index) getDocs() Docs {
	return mapToDocs(i.docs)
}
