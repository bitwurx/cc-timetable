package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	arango "github.com/arangodb/go-driver"
	arangohttp "github.com/arangodb/go-driver/http"
)

const (
	CollectionTimetables = "timetables" // the name of the timetables database collection.
)

var db arango.Database // package local arango database instance.

// DocumentMeta contains meta data for an arango document
type DocumentMeta struct {
	Id arango.DocumentID
}

// Model contains methods for interacting with database collections.
type Model interface {
	Create() error
	FetchAll() ([]interface{}, error)
	Save(interface{}) (DocumentMeta, error)
}

// TimetableModel represents a priority queue collection model.
type TimetableModel struct{}

// Create creates the timetables collection in the arangodb database.
func (model *TimetableModel) Create() error {
	_, err := db.CreateCollection(nil, CollectionTimetables, nil)
	if err != nil && arango.IsConflict(err) {
		return nil
	}
	return err
}

// FetchAll gets all documents from the timetables collection.
func (model *TimetableModel) FetchAll() ([]interface{}, error) {
	timetables := make([]interface{}, 0)
	query := fmt.Sprintf("FOR q IN %s RETURN q", CollectionTimetables)
	cursor, err := db.Query(nil, query, nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()
	for {
		t := new(Timetable)
		_, err := cursor.ReadDocument(nil, t)
		if arango.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, err
		}
		timetables = append(timetables, t)
	}
	return timetables, nil
}

// Query runs the AQL query against the timetables model collection.
func (model *TimetableModel) Save(table interface{}) (DocumentMeta, error) {
	var meta arango.DocumentMeta
	var doc struct {
		Key      string  `json:"_key"`
		Schedule []*Task `json:"schedule"`
	}
	col, err := db.Collection(nil, CollectionTimetables)
	if err != nil {
		return DocumentMeta{}, err
	}
	data, err := json.Marshal(table.(*Timetable))
	if err != nil {
		return DocumentMeta{}, err
	}
	if err := json.Unmarshal(data, &doc); err != nil {
		return DocumentMeta{}, err
	}
	meta, err = col.CreateDocument(nil, doc)
	if arango.IsConflict(err) {
		patch := map[string]interface{}{
			"schedule": doc.Schedule,
		}
		meta, err = col.UpdateDocument(nil, doc.Key, patch)
		if err != nil {
			return DocumentMeta{}, err
		}
	} else if err != nil {
		return DocumentMeta{}, err
	}
	return DocumentMeta{Id: meta.ID}, nil
}

// InitDatabase connects to the arangodb and creates the collections from the
// provided models.
func InitDatabase() {
	host := os.Getenv("ARANGODB_HOST")
	name := os.Getenv("ARANGODB_NAME")
	user := os.Getenv("ARANGODB_USER")
	pass := os.Getenv("ARANGODB_PASS")

	conn, err := arangohttp.NewConnection(
		arangohttp.ConnectionConfig{Endpoints: []string{host}},
	)
	if err != nil {
		panic(err)
	}
	client, err := arango.NewClient(arango.ClientConfig{
		Connection:     conn,
		Authentication: arango.BasicAuthentication(user, pass),
	})
	if err != nil {
		panic(err)
	}

	for {
		if exists, err := client.DatabaseExists(nil, name); err == nil {
			if !exists {
				db, err = client.CreateDatabase(nil, name, nil)
			} else {
				db, err = client.Database(nil, name)
			}
			if err == nil {
				break
			}
		}
		time.Sleep(time.Second * 1)
	}

	models := []Model{
		&TimetableModel{},
	}
	for _, model := range models {
		if err := model.Create(); err != nil {
			panic(err)
		}
	}
}
