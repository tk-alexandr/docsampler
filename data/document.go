package data

import (
	"time"
	"log"
)

//Document database interface
type Document struct {
	ID int					`json:"id"`
	Name string				`json:"name"`
	Path string				`json:"path"`
	Variables []string		`json:"variables"`
	Time time.Time			`json:"time"`
}

//Save document to database
func (doc *Document) Save () {
	open()
	defer close()

	sqlStatement := `
	INSERT INTO documents (name, path, time)
	VALUES ($1, $2, $3)
	RETURNING id`

	id := 0
	err := Db.QueryRow(sqlStatement, doc.Name, doc.Path, doc.Time).Scan(&id)

	if err != nil {
		log.Fatal(err)
	}

	saveFieldStatement := `
	INSERT INTO fields (document_id, text)
	VALUES ($1, $2)
	`

	for _, val := range doc.Variables {
		_, err = Db.Exec(saveFieldStatement, id, val)
		if err != nil {
			log.Fatal(err)
		}
	}
}