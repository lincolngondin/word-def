package main

import (
	"encoding/xml"
	"errors"
	"io"
	"os"
)

type LexicalResource struct {
	Lexicons []*Lexicon
}

func newLexicalResource() *LexicalResource {
	return &LexicalResource{
		Lexicons: make([]*Lexicon, 0),
	}
}

type Lexicon struct {
	Id                  string
	Label               string
	Language            string
	Email               string
	License             string
	Version             string
	Requires            []Requires
	LexicalEntrys       []*LexicalEntry
	Synsets             []*Synset
	SyntacticBehaviours []*SyntacticBehaviour
}

func newLexicon() *Lexicon {
	return &Lexicon{
		Id:                  "",
		Label:               "",
		Language:            "",
		Email:               "",
		License:             "",
		Version:             "",
		Requires:            make([]Requires, 0),
		LexicalEntrys:       make([]*LexicalEntry, 0),
		Synsets:             make([]*Synset, 0),
		SyntacticBehaviours: make([]*SyntacticBehaviour, 0),
	}
}

type Requires struct {
	Id      string
	Version string
}

type LexicalEntry struct {
	Id                 string
	Lemma              *Lemma
	Forms              []Form
	Senses             []*Sense
	SyntaticBehaviours []SyntacticBehaviour
}

func NewLexicalEntry() *LexicalEntry {
	return &LexicalEntry{
		Id:                 "",
		Lemma:              nil,
		Forms:              make([]Form, 0),
		Senses:             make([]*Sense, 0),
		SyntaticBehaviours: make([]SyntacticBehaviour, 0),
	}
}

type Synset struct {
	Id              string
	ILI             string
	Definitions     []Definition
	ILIDefinitions  *ILIDefinition
	SynsetRelations []*SynsetRelation
	Examples        []Example
}

func NewSynset() *Synset {
	return &Synset{
		Id:              "",
		ILI:             "",
		Definitions:     make([]Definition, 0),
		ILIDefinitions:  nil,
		SynsetRelations: make([]*SynsetRelation, 0),
		Examples:        make([]Example, 0),
	}
}

type Definition string

type ILIDefinition string

type Lemma struct {
	WrittenForm    string
	PartOfSpeech   rune
	Pronunciations []Pronunciation
	Tags           []Tag
}

func NewLemma() *Lemma {
	return &Lemma{
		WrittenForm:    "",
		PartOfSpeech:   0,
		Pronunciations: make([]Pronunciation, 0),
		Tags:           make([]Tag, 0),
	}
}

type Form struct {
	WrittenForm    string
	Pronunciations []Pronunciation
	Tags           []Tag
}

func NewForm() *Form {
	return &Form{
		WrittenForm:    "",
		Pronunciations: make([]Pronunciation, 0),
		Tags:           make([]Tag, 0),
	}
}

type Sense struct {
	Id string
	// reference to an Synset
	Synset         *Synset
	SenseRelations []*SenseRelation
	Examples       []Example
	Counts         []Count
}

func NewSense() *Sense {
    return &Sense {
        Id: "",
        Synset: nil,
        SenseRelations: make([]*SenseRelation, 0),
        Examples: make([]Example, 0),
        Counts: make([]Count, 0),
    }
}

type RelationType int8

const ()

func GetPartOfSpeech(pos rune) string {
	switch pos {
	case 'n':
		return "Noun"
    case 'v':
        return "Verb"
    case 'a':
        return "Adjective"
    case 'r':
        return "Adverb"
    case 's':
        return "Adjective Satellite"
    case 'c':
        return "Conjuction"
    case 'p':
        return "Adposition"
    case 'x':
        return "Other"
	default:
		return "Unknown"
	}
}

const (
	RelationTypeAntonym = RelationType(iota)
	RelationTypeAlso
	RelationTypeParticiple
	RelationTypePertainym
	RelationTypeDerivation
	RelationTypeDomainTopic
	RelationTypeHasDomainTopic
	RelationTypeDomainRegion
	RelationTypeHasDomainRegion
	RelationTypeExemplifies
	RelationTypeIsExemplifiedBy
	RelationTypeSimilar
	RelationTypeOther
	RelationTypeSimpleAspectIp
)

type SenseRelation struct {
	// reference to an Sense
	Target  *Sense
	RelType RelationType
}

func NewSenseRelation(target *Sense, reltype string) *SenseRelation{
    return &SenseRelation{
        Target: target,
        RelType: RelationTypeAlso,
    }
}

type Example string

type Count string

type Pronunciation string

type Tag struct {
	Category string
	Value    string
}

type SyntacticBehaviour struct {
	SubCategorizationFrame string
}

type SynsetRelation struct {
	// reference to an Synset
	Target  *Synset
	RelType RelationType
}

func NewSynsetRelation(target *Synset, reltype string) *SynsetRelation{
    return &SynsetRelation{
        Target: target,
        RelType: RelationTypeAlso,
    }

}

func ParseLexicalXML(filename string) (*LexicalResource, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	xmlDecoder := xml.NewDecoder(file)
	var lexicalResource *LexicalResource = newLexicalResource()

	var insideLexicon bool = false
	var insideLexicalEntry bool = false
	var insideLemma bool = false
	var insideForm bool = false
	var insidePronunciation = false
	var insideTag bool = false
	var insideDefinition bool = false
	var insideILIDefinition bool = false
    var insideSense bool = false
    var insideSynset bool = false
    var insideExample bool = false
    var insideCount bool = false

	var nextLexicon *Lexicon
	var nextRequires Requires
	var nextLexicalEntry *LexicalEntry
	var nextLemma *Lemma
	var nextSyntacticBehaviour *SyntacticBehaviour
	var nextPronunciation Pronunciation
	var nextTag *Tag
	var nextForm *Form
	var nextSynset *Synset
	var nextDefinition Definition
	var nextILIDefinition ILIDefinition
    var nextSense *Sense
    var nextExample Example
    var nextCount Count

    var tempSynsetIdToSynset map[string]*Synset = make(map[string]*Synset, 100000)
    var tempSenseIdToSynsetId map[string]string = make(map[string]string, 100000)
    var tempSenseIdToLinkedsSenseRelation map[string][]*SenseRelation = make(map[string][]*SenseRelation, 100000)
    var tempSynsetIdToLinkedsSynsetRelation map[string][]*SynsetRelation = make(map[string][]*SynsetRelation, 100000)
    var tempSenseIDToSense map[string]*Sense = make(map[string]*Sense, 10000)

	for {
		nextToken, decodeErr := xmlDecoder.Token()
		if decodeErr == io.EOF {
			break
		}
		if decodeErr != nil {
			return nil, errors.New("Invalid xml file!")
		}
		switch v := nextToken.(type) {
		case xml.StartElement:
			elementName := v.Name.Local
			if elementName == "Lexicon" {
				insideLexicon = true
				nextLexicon = newLexicon()
				for _, attr := range v.Attr {
					if attr.Name.Local == "id" {
						nextLexicon.Id = attr.Value
					} else if attr.Name.Local == "label" {
						nextLexicon.Label = attr.Value
					} else if attr.Name.Local == "language" {
						nextLexicon.Language = attr.Value
					} else if attr.Name.Local == "email" {
						nextLexicon.Email = attr.Value
					} else if attr.Name.Local == "license" {
						nextLexicon.License = attr.Value
					} else if attr.Name.Local == "version" {
						nextLexicon.Version = attr.Value
					}
				}
			} else if elementName == "Requires" {
				nextRequires = Requires{}
				for _, attr := range v.Attr {
					if attr.Name.Local == "id" {
						nextRequires.Id = attr.Value
					} else if attr.Name.Local == "version" {
						nextRequires.Version = attr.Value
					}
				}
			} else if elementName == "LexicalEntry" {
				insideLexicalEntry = true
				nextLexicalEntry = NewLexicalEntry()
				for _, attr := range v.Attr {
					if attr.Name.Local == "id" {
						nextLexicalEntry.Id = attr.Value
					}
				}
			} else if elementName == "Lemma" {
				insideLemma = true
				nextLemma = NewLemma()
				for _, attr := range v.Attr {
					if attr.Name.Local == "writtenForm" {
						nextLemma.WrittenForm = attr.Value
					} else if attr.Name.Local == "partOfSpeech" {
						nextLemma.PartOfSpeech = []rune(attr.Value)[0]
					}
				}
			} else if elementName == "Form" {
				insideForm = true
				nextForm = NewForm()
				for _, attr := range v.Attr {
					if attr.Name.Local == "writtenForm" {
						nextForm.WrittenForm = attr.Value
					}
				}

			} else if elementName == "Tag" {
				insideTag = true
				nextTag = &Tag{}
				for _, attr := range v.Attr {
					if attr.Name.Local == "category" {
						nextTag.Category = attr.Value
					}
				}
			} else if elementName == "Pronunciation" {
				insidePronunciation = true

			} else if elementName == "SyntacticBehaviour" {
				nextSyntacticBehaviour = &SyntacticBehaviour{}
				for _, attr := range v.Attr {
					if attr.Name.Local == "subcategorizationFrame" {
						nextSyntacticBehaviour.SubCategorizationFrame = attr.Value
					}
				}
				if insideLexicalEntry {
					nextLexicalEntry.SyntaticBehaviours = append(nextLexicalEntry.SyntaticBehaviours, *nextSyntacticBehaviour)
				} else if insideLexicon {
					nextLexicon.SyntacticBehaviours = append(nextLexicon.SyntacticBehaviours, nextSyntacticBehaviour)
				}
			} else if elementName == "Synset" {
                insideSynset = true
				nextSynset = NewSynset()
				for _, attr := range v.Attr {
					if attr.Name.Local == "id" {
						nextSynset.Id = attr.Value
					} else if attr.Name.Local == "ili" {
						nextSynset.ILI = attr.Value
					}

				}
			} else if elementName == "Definition" {
				insideDefinition = true
			} else if elementName == "ILIDefinition" {
				insideILIDefinition = true
			} else if elementName == "SynsetRelation" {
                var relType string
                var target string
				for _, attr := range v.Attr {
					if attr.Name.Local == "target" {
                        target = attr.Value
					} else if attr.Name.Local == "relType" {
                        relType = attr.Value
					}
				}
                var newSynsetRelation = NewSynsetRelation(nil, relType)
                nextSynset.SynsetRelations = append(nextSynset.SynsetRelations, newSynsetRelation)
                _, ok := tempSynsetIdToLinkedsSynsetRelation[target]
                if !ok {
                    tempSynsetIdToLinkedsSynsetRelation[target] = make([]*SynsetRelation, 0)
                }
                tempSynsetIdToLinkedsSynsetRelation[target] = append(tempSynsetIdToLinkedsSynsetRelation[target], newSynsetRelation)
			} else if elementName == "Example" {
                insideExample = true

			} else if elementName == "Sense" {
                insideSense = true
                nextSense = NewSense()
                var synsetId string
				for _, attr := range v.Attr {
					if attr.Name.Local == "id" {
                        nextSense.Id = attr.Value
                        tempSenseIDToSense[attr.Value] = nextSense
					} else if attr.Name.Local == "synset" {
                        nextSense.Synset = nil
                        synsetId = attr.Value
					}
				}
                tempSenseIdToSynsetId[nextSense.Id] = synsetId

            } else if elementName == "SenseRelation" {
                var relType string
                var target string
				for _, attr := range v.Attr {
					if attr.Name.Local == "relType" {
                        relType = attr.Value
					}else if attr.Name.Local == "target" {
                        target = attr.Value
					}
				}
                var newSenseRelation = NewSenseRelation(nil, relType)
                nextSense.SenseRelations = append(nextSense.SenseRelations, newSenseRelation)
                _, ok := tempSenseIdToLinkedsSenseRelation[target]
                if !ok {
                    tempSenseIdToLinkedsSenseRelation[target] = make([]*SenseRelation, 0)
                }
                tempSenseIdToLinkedsSenseRelation[target] = append(tempSenseIdToLinkedsSenseRelation[target], newSenseRelation)

            } else if elementName == "Count" {
                insideCount = true
            }

		case xml.EndElement:
			elementName := v.Name.Local
			if elementName == "Lexicon" {
				insideLexicon = false
				lexicalResource.Lexicons = append(lexicalResource.Lexicons, nextLexicon)
			} else if elementName == "Requires" {
				nextLexicon.Requires = append(nextLexicon.Requires, nextRequires)
			} else if elementName == "LexicalEntry" {
				insideLexicalEntry = false
				nextLexicon.LexicalEntrys = append(nextLexicon.LexicalEntrys, nextLexicalEntry)
			} else if elementName == "Lemma" {
				insideLemma = false
				nextLexicalEntry.Lemma = nextLemma
			} else if elementName == "Form" {
				insideForm = false
				nextLexicalEntry.Forms = append(nextLexicalEntry.Forms, *nextForm)
			} else if elementName == "Tag" {
				insideTag = false
				if insideLemma {
					nextLemma.Tags = append(nextLemma.Tags, *nextTag)
				} else if insideForm {
					nextForm.Tags = append(nextForm.Tags, *nextTag)
				}
			} else if elementName == "Pronunciation" {
				insidePronunciation = false
				if insideLemma {
					nextLemma.Pronunciations = append(nextLemma.Pronunciations, nextPronunciation)
				}
				if insideForm {
					nextForm.Pronunciations = append(nextForm.Pronunciations, nextPronunciation)
				}
			} else if elementName == "Synset" {
                insideSynset = false
				nextLexicon.Synsets = append(nextLexicon.Synsets, nextSynset)
                tempSynsetIdToSynset[nextSynset.Id] = nextSynset
			} else if elementName == "Definition" {
				insideDefinition = false
				nextSynset.Definitions = append(nextSynset.Definitions, nextDefinition)
			} else if elementName == "ILIDefinition" {
				insideILIDefinition = false
				nextSynset.ILIDefinitions = new(ILIDefinition)
				*nextSynset.ILIDefinitions = nextILIDefinition
			} else if elementName == "SynsetRelation" {

			} else if elementName == "Example" {
                insideExample = false
                if insideSense {
                    nextSense.Examples = append(nextSense.Examples, nextExample)
                } else if insideSynset {
                    nextSynset.Examples = append(nextSynset.Examples, nextExample)
                }
			} else if elementName == "Sense" {
                insideSense = false
                nextLexicalEntry.Senses = append(nextLexicalEntry.Senses, nextSense)
            } else if elementName == "Count" {
                insideCount = false
                nextSense.Counts = append(nextSense.Counts, nextCount)
            }

		case xml.CharData:
			if insidePronunciation {
				nextPronunciation = Pronunciation(v)
			}
			if insideTag {
				nextTag.Value = string(v)
			}
			if insideDefinition {
				nextDefinition = Definition(v)
			}
			if insideILIDefinition {
				nextILIDefinition = ILIDefinition(v)
			}
            if insideExample {
                nextExample = Example(v)
            }
            if insideCount {
                nextCount = Count(v)
            }

		case xml.Comment:
		case xml.ProcInst:
		case xml.Directive:
		default:
			return nil, errors.New("Invalid XML token!")
		}
	}
    
    // fill the field Target in SenseRelation
    for senseID, senseRelations := range tempSenseIdToLinkedsSenseRelation {
        for _, senseRL := range senseRelations {
            senseRL.Target = tempSenseIDToSense[senseID]
        }
    }

    // fill the field Target in SynsetRelation
    for synsetID, synsetsRelation := range tempSynsetIdToLinkedsSynsetRelation {
        for _, synsetRL := range synsetsRelation {
            synsetRL.Target = tempSynsetIdToSynset[synsetID]
        }
    }

    // fill the field Synset in Sense
    for senseID, synsetID := range tempSenseIdToSynsetId {
        sense := tempSenseIDToSense[senseID]
        sense.Synset = tempSynsetIdToSynset[synsetID]
    }

	return lexicalResource, nil
}
