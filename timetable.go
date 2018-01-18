package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"time"
)

// Task is a unit of work that is scheduled in the timetable.
type Task struct {
	// Id is the unique version 1 uuid assigned for task identification.
	// RunAt is the point in time to schedule for execution.
	Id    string `json:"_key"`
	RunAt string `json:"runAt"`
}

// Timetable keeps track of scheduled tasks for a given resource.
type Timetable struct {
	// Key is the task resource key.
	// schedule holds the tasks keyed on their run at times.
	Key      string
	schedule map[string]*Task
}

// Delay returns the time delay in minutes until the next scheduled task.
func (table *Timetable) Delay() (int, error) {
	if len(table.schedule) == 0 {
		return 0, errors.New("empty schedule")
	}

	tasks := make([]string, 0)
	for runAt, _ := range table.schedule {
		tasks = append(tasks, runAt)
	}
	sort.Strings(tasks)

	t, err := time.Parse(time.RFC3339, tasks[0])
	if err != nil {
		return 0, err
	}
	delay := int(t.Sub(time.Now()).Minutes())
	if delay > 0 {
		delay++
	}
	return delay, nil
}

// Insert adds the task to the schedule if the run at time is
// not already reserved.
func (table *Timetable) Insert(task *Task) error {
	if _, ok := table.schedule[task.RunAt]; ok {
		return errors.New("schedule conflict")
	}
	table.schedule[task.RunAt] = task
	return nil
}

// List returns all items in the schedule.
func (table *Timetable) List() []*Task {
	tasks := make([]*Task, 0)
	for _, task := range table.schedule {
		tasks = append(tasks, task)
	}
	return tasks
}

// Next returns the next task in the schedule
func (table *Timetable) Next() *Task {
	var next *time.Time
	for k := range table.schedule {
		t, _ := time.Parse(time.RFC3339, k)
		if next == nil {
			next = &t
			continue
		}
		if t.Before(*next) {
			next = &t
		}
	}
	return table.schedule[next.Format(time.RFC3339)]
}

// Remove deletes the task with the matching run at time from
// the timetable.
func (table *Timetable) Remove(runAt string) error {
	if _, ok := table.schedule[runAt]; !ok {
		return errors.New("not found")
	}
	delete(table.schedule, runAt)
	return nil
}

// Save writes the timetable to the database.
func (table *Timetable) Save(model Model) (DocumentMeta, error) {
	return model.Save(table)
}

// MarshalJSON serializes the timetable key and schedule.
func (table *Timetable) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString(
		fmt.Sprintf(`{"_key": "%s", "schedule": %s}`, table.Key, (func() string {
			tasks := bytes.NewBuffer([]byte("["))
			i := 0
			for _, task := range table.schedule {
				tasks.WriteString(
					fmt.Sprintf(`{"_key": "%s", "runAt": "%s"}`, task.Id, task.RunAt),
				)
				if i < (len(table.schedule) - 1) {
					tasks.WriteByte(',')
				}
				i++
			}
			tasks.WriteString("]")
			return tasks.String()
		})(),
		))
	return buf.Bytes(), nil
}

// UnmarshalJSON deserializes the stored timetable meta data into
// a timetable instance.
func (table *Timetable) UnmarshalJSON(b []byte) error {
	if table.schedule == nil {
		table.schedule = make(map[string]*Task)
	}
	data := make(map[string]interface{})
	json.Unmarshal(b, &data)
	table.Key = data["_key"].(string)
	for _, task := range data["schedule"].([]interface{}) {
		v, _ := task.(map[string]interface{})
		task := &Task{
			Id:    v["_key"].(string),
			RunAt: v["runAt"].(string),
		}
		table.schedule[task.RunAt] = task
	}
	return nil
}

// Newtimetable creates a new Timetable instance.
func NewTimetable(key string) *Timetable {
	return &Timetable{key, make(map[string]*Task)}
}
