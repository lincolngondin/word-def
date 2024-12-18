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
	for _, wordDefinition := range word.WordDefinitions {
        builderString.WriteString(fmt.Sprintf("[blue::b]%s[-::-]([green]%s[-]):", wordDefinition.WrittenForm, wordDefinition.PartOfSpeech))
		if len(wordDefinition.Definitions) == 0 {
			builderString.WriteString("There's no definitions for this word!")
		}
		for i, def := range wordDefinition.Definitions {
			builderString.WriteString(fmt.Sprintf("\n[yellow]%d: %s[yellow]\n", i+1, def.Definitions[0]))
			if len(def.UseExamples) != 0 {
                builderString.WriteString("[red::u]Examples[-::-]: \n")
			}
			for _, example := range def.UseExamples {
				builderString.WriteString(fmt.Sprintf(" - [cyan]%s[-]\n", example))
			}
		}
        builderString.WriteString("\n")
	}
	return builderString.String()
}

func initApplication(dict Dictionary) {
	app := tview.NewApplication().EnableMouse(true)
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	textView := tview.NewTextView().SetDynamicColors(true)
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
