package main

import (
	"os"
	"time"

	arango "github.com/arangodb/go-driver"
	arangohttp "github.com/arangodb/go-driver/http"
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

	models := []Model{}
	for _, model := range models {
		if err := model.Create(); err != nil {
			panic(err)
		}
	}
}
