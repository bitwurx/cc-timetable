package main

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestTimetableDelay(t *testing.T) {
	timetable := NewTimetable("test")
	now := time.Now().Add(time.Minute * 5)
	if _, err := timetable.Delay(); err == nil {
		t.Fatal("expected empty schedule error")
	}
	if err := timetable.Insert(&Task{RunAt: "test"}); err != nil {
		t.Fatal(err)
	}
	delay, err := timetable.Delay()
	if err == nil {
		t.Fatal("expected time parse error")
	}
	timetable.Remove("test")
	if err := timetable.Insert(&Task{RunAt: now.Format(time.RFC3339)}); err != nil {
		t.Fatal(err)
	}
	delay, err = timetable.Delay()
	if err != nil {
		t.Fatal(err)
	}
	if delay != 5 {
		t.Fatalf("expected delay to be 5 minutes, got %d", delay)
	}
}

func TestTimetableRemove(t *testing.T) {
	timetable := NewTimetable("test")
	if err := timetable.Remove("abc123"); err == nil {
		t.Fatal("expected not found error")
	}
	if err := timetable.Insert(&Task{Id: "abc123"}); err != nil {
		t.Fatal(err)
	}
	if err := timetable.Remove("abc123"); err != nil {
		t.Fatal(err)
	}
	if len(timetable.List()) != 0 {
		t.Fatal("expected there to be no tasks")
	}
}

func TestTimetableInsert(t *testing.T) {
	timetable := NewTimetable("test")
	runAt := time.Now().Format(time.RFC3339)
	if err := timetable.Insert(&Task{RunAt: runAt}); err != nil {
		t.Fatal(err)
	}
	if err := timetable.Insert(&Task{RunAt: runAt}); err == nil {
		t.Fatal("expected schedule conflict error")
	}
	if timetable.List()[0].RunAt != runAt {
		t.Fatal("unexpected task run at time")
	}
}

func TestTimetableList(t *testing.T) {
	timetable := NewTimetable("test")
	now := time.Now()
	if err := timetable.Insert(&Task{RunAt: now.Format(time.RFC3339)}); err != nil {
		t.Fatal(err)
	}
	now = now.Add(time.Minute * 5)
	if err := timetable.Insert(&Task{RunAt: now.Format(time.RFC3339)}); err != nil {
		t.Fatal(err)
	}
	tasks := timetable.List()
	if len(tasks) != 2 {
		t.Fatal("expected tasks count to be 2")
	}
}

func TestTimetableNext(t *testing.T) {
	now := time.Now()
	tasks := []*Task{
		{RunAt: now.Add(time.Minute * 3).Format(time.RFC3339)},
		{RunAt: now.Format(time.RFC3339)},
		{RunAt: now.Add(time.Minute * 5).Format(time.RFC3339)},
	}
	timetable := NewTimetable("test")
	if task := timetable.Next(); task != nil {
		t.Fatal("expected task to be nil")
	}
	for _, task := range tasks {
		if err := timetable.Insert(task); err != nil {
			t.Fatal(err)
		}
	}
	task := timetable.Next()
	if task.RunAt != tasks[1].RunAt {
		t.Fatal("got unexpected next task")
	}
}

func TestTimetableSave(t *testing.T) {
	var model Model
	if testing.Short() {
		model = MockModel{}
	} else {
		model = &TimetableModel{}
	}
	timetable := NewTimetable("test")
	timetable.Insert(&Task{Id: "xyz", RunAt: time.Now().Format(time.RFC3339)})
	if _, err := timetable.Save(model); err != nil {
		t.Fatal(err)
	}
}

func TestTimetableMarshalJSON(t *testing.T) {
	runAt1 := time.Now().Format(time.RFC3339)
	runAt2 := time.Now().Add(time.Second * 30).Format(time.RFC3339)
	timetable := NewTimetable("test")
	timetable.Insert(&Task{Id: "123", RunAt: runAt1})
	timetable.Insert(&Task{Id: "321", RunAt: runAt2})
	data, err := json.Marshal(timetable)
	if err != nil {
		t.Fatal(err)
	}
	s1 := `{"_key":"test","schedule":[{"_key":"123","runAt":"%s"},{"_key":"321","runAt":"%s"}]}`
	s2 := `{"_key":"test","schedule":[{"_key":"321","runAt":"%s"},{"_key":"123","runAt":"%s"}]}`
	if string(data) != fmt.Sprintf(s1, runAt1, runAt2) && string(data) != fmt.Sprintf(s2, runAt2, runAt1) {
		t.Fatal("got unexpected value from timetable json marshal")
	}
}

func TestTimetableUnmarshalJSON(t *testing.T) {
	timetable := new(Timetable)
	runAt := time.Now().Format(time.RFC3339)
	b := []byte(fmt.Sprintf(`{"_key":"test","schedule":[{"_key":"123","runAt":"%s"}]}`, runAt))
	if err := json.Unmarshal(b, timetable); err != nil {
		t.Fatal(err)
	}
	if timetable.Key != "test" {
		t.Fatal("expected timetable key to be 'test'")
	}
	if timetable.List()[0].RunAt != runAt {
		t.Fatalf("expected timetable run at to be %s", runAt)
	}
}
