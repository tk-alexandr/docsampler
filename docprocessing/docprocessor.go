package docprocessing

import (
	//"os"
	"log"
	"github.com/unidoc/unioffice/document"
	//"strings"
	//"bytes"
)
//GetDocumentVariables receives path to a saved file and looking for variables in it
func GetDocumentVariables (path string) ([]string, error) {
	doc, err := document.Open(path)
	if err != nil {
		log.Printf("error opening document: %s", err)
		return nil, err
	}
	
	variables := make([]string, 0)

	varFlag := false
	for _, paragraph := range doc.Paragraphs() {
		for _, r := range paragraph.Runs(){
			text := r.Text()
						
			if text == "}}" {
				varFlag = false
			}

			if varFlag {
				variables = append(variables, text)
			}

			if text == "{{" {
				varFlag = true
			}
		}
	}

	return variables, nil
}

//SetDocumentVariables receives path and values to set in file and returns link on processed document
func SetDocumentVariables(filePath string, values map[string]string) (string, error) {
	doc, err := document.Open(filePath)
	if err != nil {
		log.Printf("error opening document: %s", err)
		return "", err
	}

	varFlag := false
	for _, paragraph := range doc.Paragraphs() {
		for _, r := range paragraph.Runs(){
			text := r.Text()
						
			if text == "}}" {
				varFlag = false
				r.ClearContent()
			}

			if varFlag {
				r.ClearContent()
				r.AddText(values[text])
			}

			if text == "{{" {
				varFlag = true
				r.ClearContent()
			}
		}
	}

	doc.SaveToFile("edit-document.docx")

	return "edit-document.docx", nil
}