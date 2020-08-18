package docprocessing

import (
	//"os"
	"log"
	"github.com/unidoc/unioffice/document"
	//"strings"
	//"bytes"
)

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