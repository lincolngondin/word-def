package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

type wordDefinition struct {
	definition  string
	useExamples []string
}

func newWordDefinition() wordDefinition {
	return wordDefinition{
		definition:  "",
		useExamples: make([]string, 0),
	}
}

type ISearch interface {
    Search(query string) (*Word, error)
}

type Word struct {
    writtenForm string 
    partOfSpeech string
    definitions []wordDefinition
}

type Dictionary struct {
    synsetIdToDefinition map[string]wordDefinition
    dictionary map[string]dictionaryEntry
}

func NewDictionary() *Dictionary {
    return &Dictionary{
        synsetIdToDefinition: make(map[string]wordDefinition, 100000),
        dictionary: make(map[string]dictionaryEntry, 100000),
    }
}

func (dicti *Dictionary) Search(query string) (*Word, error) {
    dictEntry, ok := dicti.dictionary[query]
    if !ok {
        return nil, errors.New("Word not found!")
    }

    word := &Word{
        writtenForm: dictEntry.writtenForm,
        partOfSpeech: dictEntry.partOfSpeech,
        definitions: make([]wordDefinition, 0, len(dictEntry.synsetsId)),
    }

    for _, synsetId := range dictEntry.synsetsId {
        wordDef, ok := dicti.synsetIdToDefinition[synsetId]
        if !ok {
            continue
        }
        newDef := wordDefinition{
            definition: wordDef.definition,
            useExamples: make([]string, len(wordDef.useExamples)),
        }
        copy(newDef.useExamples, wordDef.useExamples)
        word.definitions = append(word.definitions, newDef)
    }
    return word, nil
}

type dictionaryEntry struct {
	writtenForm  string
	partOfSpeech string
	synsetsId    []string
}

func newDictionaryEntry() dictionaryEntry {
	return dictionaryEntry{
		writtenForm:  "",
		partOfSpeech: "",
		synsetsId:    make([]string, 0),
	}
}


/*
func showWordDefinition(word string) {
	dictEntry, ok := dictionary[word]
    if !ok {
        fmt.Printf("Word not found!\n")
        return
    }
    
	fmt.Printf("%s (%s): \n", dictEntry.writtenForm, dictEntry.partOfSpeech)
	for i, synsetId := range dictEntry.synsetsId {
		def, ok := synsetIdToDefinition[synsetId]
        if !ok {
            continue
        }

		fmt.Printf("df %d: %s\n\n", i+1, def.definition)
		if len(def.useExamples) > 0 {
			fmt.Printf("Examples: \n")
		}
		for _, example := range def.useExamples {
			fmt.Printf(" - %s\n\n", example)
		}
	}
}
*/

func main() {
	file, err := os.Open("wn.xml")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	xmlDecoder := xml.NewDecoder(file)

    WordDict := NewDictionary()

	var nextDictionaryEntry dictionaryEntry = newDictionaryEntry()
	var nextLexicalEntryId string
	fmt.Println(nextLexicalEntryId)
	var nextSynsetId string
	var nextExamples []string = make([]string, 0)
	var nextDefinition string
	var insideExample bool = false
	var insideDefinition bool = false

loop:
	for {
		tok, err := xmlDecoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		switch v := tok.(type) {
		case xml.StartElement:
			elementName := v.Name.Local
			if elementName == "LexicalEntry" {
				nextDictionaryEntry = newDictionaryEntry()
				for _, attr := range v.Attr {
					if attr.Name.Local == "id" {
						nextLexicalEntryId = attr.Value
					}
				}
			} else if elementName == "Lemma" {
				for _, attr := range v.Attr {
					if attr.Name.Local == "writtenForm" {
						nextDictionaryEntry.writtenForm = attr.Value
					} else if attr.Name.Local == "partOfSpeech" {
						nextDictionaryEntry.partOfSpeech = attr.Value
					}
				}
			} else if elementName == "Sense" {
				for _, attr := range v.Attr {
					if attr.Name.Local == "synset" {
						nextDictionaryEntry.synsetsId = append(nextDictionaryEntry.synsetsId, attr.Value)
					}
				}
			} else if elementName == "Synset" {
				nextExamples = make([]string, 0)
				for _, attr := range v.Attr {
					if attr.Name.Local == "id" {
						nextSynsetId = attr.Value
					}
				}
			} else if elementName == "Example" {
				insideExample = true
			} else if elementName == "Definition" {
				insideDefinition = true
			}

		case xml.EndElement:
			elementName := v.Name.Local
			if elementName == "LexicalEntry" {
				WordDict.dictionary[nextDictionaryEntry.writtenForm] = nextDictionaryEntry
			} else if elementName == "Lemma" {
			} else if elementName == "Example" {
				insideExample = false
			} else if elementName == "Definition" {
				insideDefinition = false
			} else if elementName == "Synset" {
				WordDict.synsetIdToDefinition[nextSynsetId] = wordDefinition{
					definition:  nextDefinition,
					useExamples: nextExamples,
				}
			}
		case xml.CharData:
			if insideExample {
				nextExamples = append(nextExamples, string(v))
			}
			if insideDefinition {
				nextDefinition = string(v)
			}
		case xml.Comment:
		case xml.ProcInst:
		case xml.Directive:
		default:
			fmt.Println("Invalid token!")
			break loop
		}
	}

    initApplication(WordDict)
}
