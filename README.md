# Concord Timetable

Concord Timetable maintains a schedule of tasks that require point in time execution.

### Usage
To build the docker image run:

`make build`

To run the full test suite run:

`make test`

To run the short (dependency free) test suite run:

`make test-short`

### JSON-RPC 2.0 HTTP API - Method Reference

This service uses the [JSON-RPC 2.0 Spec](http://www.jsonrpc.org/specification) over HTTP for its API.

---
#### delay(key) : get the time until next task execution
---

#### Returns:
(*Number*) the amount of minutes until the next task is scheduled

---
#### get(key) : get a timetable by key
---

#### Parameters:

key - (*String*) the time table key.

#### Returns:
(*Object*) the timetable with the associated key

---
#### getAll() : get all timetables
---

#### Returns:
(*Array*) the list of all existing timetables

---
#### insert(key, id, runAt) : adds a task to a timetable schedule
---

#### Parameters:

key - (*String*) the resource key for the task.

id - (*String*) the id of the task.

runAt - (*String*) the execution point in time of the task.

#### Returns:
(*Number*) 0 on success or -1 on failure 

---
#### remove(key, runAt) - remove a task from a timetable
---

#### Parameters:

key - (*String*) the timetable key.

runAt - (*String*) the execution point in time of the task.

#### Returns:
(*Number*) 0 on success or -1 on failure

