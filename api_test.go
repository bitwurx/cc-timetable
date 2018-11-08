package main

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/bitwurx/jrpc2"
)

func TestApiV1Delay(t *testing.T) {
	api := NewApiV1(&MockModel{}, jrpc2.NewServer("", ""))
	runAt := time.Now().Add(time.Minute * 5).Format(time.RFC3339)
	result, errObj := api.Insert([]byte(fmt.Sprintf(`{"key": "delay", "id": "abc123", "runAt": "%s"}`, runAt)))
	if errObj != nil {
		t.Fatal(errObj.Message)
	}
	if result != 0 {
		t.Fatal("expected result to be 0")
	}
	result, errObj = api.Delay([]byte(`{"key": "delay"}`))
	if result != 5 {
		t.Fatal("expected delay to be 5")
	}
}

func TestApiV1Get(t *testing.T) {
	api := NewApiV1(&MockModel{}, jrpc2.NewServer("", ""))
	runAt := time.Now().Format(time.RFC3339)
	result, errObj := api.Insert([]byte(fmt.Sprintf(`{"key": "get", "id": "abc123", "runAt": "%s"}`, runAt)))
	if errObj != nil {
		t.Fatal(errObj.Message)
	}
	if result != 0 {
		t.Fatal("expected result to be 0")
	}
	timetable, errObj := api.Get([]byte(`{"key": "get"}`))
	if errObj != nil {
		t.Fatal(errObj.Message)
	}
	m := make(map[string]interface{})
	data, err := json.Marshal(timetable)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}
	if m["_key"].(string) != "get" {
		t.Fatal("expected key to be 'get'")
	}
}

func TestApiV1GetAll(t *testing.T) {
	api := NewApiV1(&MockModel{}, jrpc2.NewServer("", ""))
	_, errObj := api.Insert([]byte(fmt.Sprintf(`{"key": "k1", "id": "abc123", "runAt": "%s"}`, time.Now().String())))
	if errObj != nil {
		t.Fatal(errObj.Message)
	}
	_, errObj = api.Insert([]byte(fmt.Sprintf(`{"key": "k2", "id": "abc123", "runAt": "%s"}`, time.Now().String())))
	if errObj != nil {
		t.Fatal(errObj.Message)
	}
	timetables, errObj := api.GetAll([]byte(`{"key": "getAll"}`))
	if errObj != nil {
		t.Fatal(errObj.Message)
	}
	var timtable []map[string]interface{}
	data, err := json.Marshal(timetables)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, &timtable); err != nil {
		t.Fatal(err)
	}
	keys := make([]string, 0)
	keys = append(keys, timtable[0]["_key"].(string))
	keys = append(keys, timtable[1]["_key"].(string))

	for _, k := range keys {
		if k != "k1" && k != "k2" {
			t.Fatal("got unexpected timetable key")
		}
	}
}

func TestApiV1Insert(t *testing.T) {
	api := NewApiV1(&MockModel{}, jrpc2.NewServer("", ""))
	runAt := time.Now().Format(time.RFC3339)
	result, err := api.Insert([]byte(fmt.Sprintf(`{"key": "get", "id": "abc123", "runAt": "%s"}`, runAt)))
	if err != nil {
		t.Fatal(err)
	}
	if result != 0 {
		t.Fatal("expected result to be 0")
	}
}

func TestApiV1Next(t *testing.T) {
	api := NewApiV1(&MockModel{}, jrpc2.NewServer("", ""))
	now := time.Now()
	_, errObj := api.Insert([]byte(fmt.Sprintf(`{"key": "k3", "id": "abc123", "runAt": "%s"}`, now.Format(time.RFC3339))))
	if errObj != nil {
		t.Fatal(errObj.Message)
	}
	_, errObj = api.Insert([]byte(fmt.Sprintf(`{"key": "k3", "id": "abc123", "runAt": "%s"}`, now.Add(time.Minute*5).Format(time.RFC3339))))
	if errObj != nil {
		t.Fatal(errObj.Message)
	}
	result, err := api.Next([]byte(fmt.Sprintf(`{"key": "k3"}`)))
	if err != nil {
		t.Fatal(err)
	}
	task := result.(*Task)
	if task.RunAt == now.String() {
		t.Fatal("expected run at time to be %s, got %s", task.RunAt, now.String())
	}
}

func TestApiV1Remove(t *testing.T) {
	api := NewApiV1(&MockModel{}, jrpc2.NewServer("", ""))
	runAt := time.Now().String()
	result, errObj := api.Insert([]byte(fmt.Sprintf(`{"key": "test1", "id": "abc321", "runAt": "%s"}`, runAt)))
	if errObj != nil {
		t.Fatal(errObj)
	}
	if result != 0 {
		t.Fatal("expected result to be 0")
	}
	result, errObj = api.Remove([]byte(`{"key": "test1", "id": "9g49g44"}`))
	if errObj != nil {
		t.Fatal(errObj.Message)
	}
	if result != -1 {
		t.Fatal("expected result to be -1")
	}
	result, errObj = api.Remove([]byte(fmt.Sprintf(`{"key": "test1", "id": "abc321"}`)))
	if errObj != nil {
		t.Fatal(errObj.Message)
	}
	if result != 0 {
		t.Fatal("expected result to be 0")
	}
	for _, task := range api.timetables["test1"].List() {
		if task.Id == "abc321" {
			t.Fatal("expected task with id 'abc321' to be removed")
		}
	}
}
