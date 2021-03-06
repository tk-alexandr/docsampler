package data

import (
	"encoding/json"
	"log"
)

//DocumentManager mager needed to manipulate documents
type DocumentManager struct {}

//GetDocumentListJSON returns documents list from database in Json format
func (manager *DocumentManager) GetDocumentListJSON () (string, error) {
	//TODO Create new type for Json encoding

	docs, err := manager.GetDocumentList()
	if err != nil {
		return "", err
	}

	jsonData, err := json.Marshal(docs)
	if err != nil {
		panic(err)
	}
	return string(jsonData), nil
}

//GetDocumentList returns list of documents from database
func (manager *DocumentManager) GetDocumentList () ([]Document, error) {
	open()
	defer close()

	stmtDocs := `
	SELECT id, name, time 
	FROM documents
	ORDER BY time
	`

	stmtFields := `
	SELECT text 
	FROM fields 
	WHERE document_id = $1
	`	

	rowsDocs, err := Db.Query(stmtDocs)
	if err != nil {
		log.Printf("Error in GetDocumentList, getting docs: %s", err)
		return nil, err
	}

	var docs []Document

	for rowsDocs.Next() {
		var doc Document

		rowsDocs.Scan(&doc.ID, &doc.Name, &doc.Time)

		rowsFields, err := Db.Query(stmtFields, doc.ID)
		if err != nil {
			log.Printf("Error in GetDocumentList, getting fields: %s", err)
		}

		for rowsFields.Next() {
			var field string
			rowsFields.Scan(&field)

			doc.Variables = append(doc.Variables, field)
		}
		docs = append(docs, doc)
	}
	
	return docs, nil
}

//Load gets document from database with id
func (manager *DocumentManager) Load(id int) (Document, error) {
	open()
	defer close()
	
	stmtDoc := `
	SELECT id, name, path, time 
	FROM documents
	WHERE id = $1
	`
 
	stmtFields := `
	SELECT text 
	FROM fields 
	WHERE document_id = $1
	`

	var doc Document

	err := Db.QueryRow(stmtDoc, id).Scan(&doc.ID, &doc.Name, &doc.Path, &doc.Time)
	if err != nil {
		log.Printf("Error in loading document: %s", err)
		return doc, err
	}

	fieldsRows, err := Db.Query(stmtFields, id)
	if err != nil {
		log.Printf("Error in loading document: %s", err)
		return doc, err
	}

	for fieldsRows.Next() {
		var field string
		fieldsRows.Scan(&field)

		doc.Variables = append(doc.Variables, field)
	}
	

	return doc, nil
}

