package main

import "errors"

// represent an sense
type Def struct {
	Definitions []string
	UseExamples []string
}

type WordDefinition struct {
	WrittenForm  string
	PartOfSpeech string
	Definitions  []Def
}

type Word struct {
	WordDefinitions []WordDefinition
}

func NewWord() *Word {
	return &Word{
		WordDefinitions: make([]WordDefinition, 0),
	}
}

type Dictionary interface {
	Search(query string) (*Word, error)
}

type OpenEnglishDictionary struct {
	lx                 *LexicalResource
	wordToLexicalEntry map[string][]*LexicalEntry
	// link the alternatives names for a word inside the Form element and link to the name
	alternativeNames map[string]string
}

func NewOpenEnglishDictionary(lx *LexicalResource) *OpenEnglishDictionary {
	if lx == nil {
		return &OpenEnglishDictionary{lx: nil}
	}

	var wordToLexicalEntry map[string][]*LexicalEntry = make(map[string][]*LexicalEntry, 100000)
	var alternativeNames map[string]string = make(map[string]string, 10000)

	for _, lexicalEntry := range lx.Lexicons[0].LexicalEntrys {
		_, ok := wordToLexicalEntry[lexicalEntry.Lemma.WrittenForm]
		if !ok {
			wordToLexicalEntry[lexicalEntry.Lemma.WrittenForm] = make([]*LexicalEntry, 0)
		}
		// link the written form of the lemma to the lexicalEntry
		wordToLexicalEntry[lexicalEntry.Lemma.WrittenForm] = append(wordToLexicalEntry[lexicalEntry.Lemma.WrittenForm], lexicalEntry)

		for _, form := range lexicalEntry.Forms {
			alternativeNames[form.WrittenForm] = lexicalEntry.Lemma.WrittenForm
		}
	}
	return &OpenEnglishDictionary{
		lx:                 lx,
		wordToLexicalEntry: wordToLexicalEntry,
		alternativeNames:   alternativeNames,
	}
}

func (oe *OpenEnglishDictionary) Search(query string) (*Word, error) {
	finded, ok := oe.wordToLexicalEntry[query]
	if !ok {
		// search by the alternative names
		altName, okAlt := oe.alternativeNames[query]
		if !okAlt {
			return nil, errors.New("Word not found!")
		}
        finded, ok = oe.wordToLexicalEntry[altName]
        if !ok {
			return nil, errors.New("Word not found!")
        }
	}

	wordToReturn := NewWord()
	for _, v := range finded {
		var defs []Def = make([]Def, 0)
		for _, sense := range v.Senses {
			newDef := Def{
				Definitions: make([]string, len(sense.Synset.Definitions)),
				UseExamples: make([]string, len(sense.Synset.Examples)),
			}
			for i, definition := range sense.Synset.Definitions {
				newDef.Definitions[i] = string(definition)
			}
			for i, example := range sense.Synset.Examples {
				newDef.UseExamples[i] = string(example)
			}
			defs = append(defs, newDef)
		}
		wordToReturn.WordDefinitions = append(wordToReturn.WordDefinitions, WordDefinition{
			WrittenForm:  v.Lemma.WrittenForm,
			PartOfSpeech: GetPartOfSpeech(v.Lemma.PartOfSpeech),
			Definitions:  defs,
		})
	}

	return wordToReturn, nil
}
