package main

import (
	"flag"
	"os"
	"testing"

	arango "github.com/arangodb/go-driver"
	arangohttp "github.com/arangodb/go-driver/http"
)

func TestMain(m *testing.M) {
	flag.Parse()
	if !testing.Short() {
		InitDatabase()
	}
	result := m.Run()
	if !testing.Short() {
		tearDownDatabase()
	}
	os.Exit(result)
}

func tearDownDatabase() {
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
	if db, err := client.Database(nil, name); err != nil {
		panic(err)
	} else {
		if err = db.Remove(nil); err != nil {
			panic(err)
		}
	}
}

type MockModel struct{}

func (m MockModel) Create() error {
	return nil
}

func (m MockModel) FetchAll() ([]interface{}, error) {
	return make([]interface{}, 0), nil
}

func (m MockModel) Save(interface{}) (DocumentMeta, error) {
	return DocumentMeta{}, nil
}
