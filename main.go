package main

import (
	"strconv"
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
	staticFilesHandler := http.FileServer(http.Dir("./example"))
	http.Handle("/example/", http.StripPrefix("/example/", staticFilesHandler))

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
	//http.Redirect(writer, request, "/doclist", http.StatusSeeOther)
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
	data := make(map[string]string) 

	id, err := strconv.Atoi(request.FormValue("id"))
	if err != nil {
		log.Printf("Cannot convert id to int: %s", err)
		return
	}
	
	doc, err := documentManager.Load(id)
	if err != nil {
		return
	}

	for _, val := range doc.Variables {
		data[val] = request.FormValue(val)
	}

	resFilePath, err := docprocessing.SetDocumentVariables(doc.Path, data)
	if err != nil {
		return
	}

	Openfile, err := os.Open(resFilePath)
	
	if err != nil {
		http.Error(writer, "File not found.", 404)
		return
	}
	defer Openfile.Close()

	
	FileHeader := make([]byte, 512)
	
	Openfile.Read(FileHeader)
	
	FileContentType := http.DetectContentType(FileHeader)

	FileStat, _ := Openfile.Stat()                     
	FileSize := strconv.FormatInt(FileStat.Size(), 10) 

	writer.Header().Set("Content-Disposition", "attachment; filename="+resFilePath)
	writer.Header().Set("Content-Type", FileContentType)
	writer.Header().Set("Content-Length", FileSize)

	Openfile.Seek(0, 0)
	
	io.Copy(writer, Openfile)
}