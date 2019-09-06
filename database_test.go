package main

import (
	"flag"
	"os"
	"testing"
	"time"

	arango "github.com/arangodb/go-driver"
	arangohttp "github.com/arangodb/go-driver/http"
)

func TestIntegrationAll(t *testing.T) {
	testMap := map[string]func(t *testing.T){
		"testTimetableModelCreate":   testTimetableModelCreate,
		"testTimetableModelSave":     testTimetableModelSave,
		"testTimetableModelFetchAll": testTimetableModelFetchAll}
	flag.Parse()
	if !testing.Short() {
		for name, testFunc := range testMap {
			InitDatabase()
			t.Run(name, testFunc)
			tearDownDatabase()
		}

	} else {
		t.Skip("Skipping integration tests")
	}
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

func testTimetableModelCreate(t *testing.T) {
	model := new(TimetableModel)
	if err := model.Create(); err != nil {
		t.Fatal(err)
	}
}

func testTimetableModelSave(t *testing.T) {
	timetable := NewTimetable("this")
	timetable.Insert(&Task{Id: "123", RunAt: time.Now().String()})
	model := new(TimetableModel)
	if _, err := model.Save(timetable); err != nil {
		t.Fatal(err)
	}
}

func testTimetableModelFetchAll(t *testing.T) {
	model := new(TimetableModel)
	if _, err := model.Save(NewTimetable("key1")); err != nil {
		t.Fatal(err)
	}
	timetables, err := model.FetchAll()
	if err != nil {
		t.Fatal(err)
	}
	if timetables[0].(*Timetable).Key != "key1" {
		t.Fatal("expected timetable key to be 'key1' was ", timetables[0].(*Timetable).Key)
	}
}
