package main

import "fmt"

func main() {
	lr, err := ParseLexicalXML("wn.xml")
    if err != nil {
        fmt.Println("Error parsing the .xml dictionary file!")
    }
	dict := NewOpenEnglishDictionary(lr)

	initApplication(dict)
}
