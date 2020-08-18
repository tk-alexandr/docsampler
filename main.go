package main

import (
	"docsampler/docprocessing"
	"path/filepath"
	"fmt"
	"io"
	"net/http"
	"os"
	"log"
	"time"	

	"docsampler/data"

	_ "github.com/lib/pq"
)

var documentManager data.DocumentManager

func main() {	
	http.HandleFunc("/create_sample", createSampleHandler)
	http.HandleFunc("/doclist", docListHandler)
	http.HandleFunc("/generate", generateHandler)

	//TODO delete if examples don't needed
	http.Handle("/example/", http.StripPrefix("/example/",http.FileServer(http.Dir("./example"))))

	http.ListenAndServe(":2222", nil)
}

func createSampleHandler(writer http.ResponseWriter, request *http.Request) {
	
	request.ParseMultipartForm(32 << 20)
	uploaded, handler, err := request.FormFile("file")
	if err != nil {
		log.Print(err)
		return
	}
	defer uploaded.Close()
		
	if filepath.Ext(handler.Filename) != ".docx" {
		log.Print("File extension needs to be .docx")
		return
	}

	path := "./app_data/" + handler.Filename
	
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Print(err)
		return
	}
	
	io.Copy(file, uploaded)
	file.Close()

	name := request.FormValue("name")
	
	docVars, err := docprocessing.GetDocumentVariables(path)
	if err != nil {
		log.Print(err)
		//TODO delete file
		return
	}
	
	doc := data.Document { 
		Name: name,
		Path: path,
		Time: time.Now(),
		Variables: docVars,
	}
	doc.Save()
	//return message
	http.Redirect(writer, request, "/doclist", http.StatusSeeOther)
}

func docListHandler(writer http.ResponseWriter, request *http.Request) {
	docsJSON, err := documentManager.GetDocumentListJSON()
	if err != nil {
		log.Print(err)
		return
	}

	writer.Header().Set("Access-Control-Allow-Origin", "*")

	fmt.Fprint(writer, docsJSON)
}

func generateHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprint(writer, "Hello generate handler")
}