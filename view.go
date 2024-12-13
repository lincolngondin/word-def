package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)


func generateTextToShow(word *Word) string {
    builderString := &strings.Builder{}
    builderString.WriteString(fmt.Sprintf("%s (%s):\n", word.writtenForm, word.partOfSpeech))
    if len(word.definitions) == 0 {
        builderString.WriteString("There's no definitions for this word!")
    }
    for i, def := range word.definitions {
        builderString.WriteString(fmt.Sprintf("%d: %s\n", i+1, def.definition))
        if len(def.useExamples) != 0 {
            builderString.WriteString("Examples: \n")
        }
        for _, example := range def.useExamples {
            builderString.WriteString(fmt.Sprintf(" - %s\n\n", example))
        }
    }
    return builderString.String()
}

func initApplication(dict *Dictionary) {
	app := tview.NewApplication().EnableMouse(true)
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	textView := tview.NewTextView()
	textView.SetBorder(true).SetTitle("Definition")


    textArea := tview.NewTextArea().SetLabel("Enter you search: ")
	textArea.SetBorder(true).SetBorderAttributes(tcell.AttrBold)
    textArea.SetChangedFunc(func() {
        textView.ScrollToBeginning()
        input := textArea.GetText()
        word, err := dict.Search(input)
        if err != nil {
            textView.SetText("Word not found!")
        } else {
            textView.SetText(generateTextToShow(word))
        }
    })

	flex.AddItem(textView, 0, 9, false)
	flex.AddItem(textArea, 0, 1, true)

	app.SetRoot(flex, true)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
