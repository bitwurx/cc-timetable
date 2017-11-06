package main

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/bitwurx/jrpc2"
)

const (
	TimetableNotFoundCode jrpc2.ErrorCode = -32002 // timetable not found json rpc 2.0 error code.
)

const (
	TimetableNotFoundMsg jrpc2.ErrorMsg = "Timetable not found" // timetable not found json rpc 2.0 error message.
)

// ApiV1 is the version 1 implementation of the rpc methods.
type ApiV1 struct {
	// model the priority timetable database model.
	// timetables is a represetation of timetables by key.
	model      Model
	timetables map[string]*Timetable
}

// DelayParams contains the rpc parameters for the Delay method.
type DelayParams struct {
	// Key is the timetable key.
	Key *string `json:"key"`
}

// FromPositional parses the key from the positional parameters.
func (params *DelayParams) FromPositional(args []interface{}) error {
	if len(args) != 1 {
		return errors.New("key parameter is required")
	}
	key := args[0].(string)
	params.Key = &key

	return nil
}

// Delay returns the time until the next scheduled point in time execution.
func (api *ApiV1) Delay(params json.RawMessage) (interface{}, *jrpc2.ErrorObject) {
	p := new(DelayParams)
	if err := jrpc2.ParseParams(params, p); err != nil {
		return nil, err
	}
	if p.Key == nil {
		return nil, &jrpc2.ErrorObject{
			Code:    jrpc2.InvalidParamsCode,
			Message: jrpc2.InvalidParamsMsg,
			Data:    "timetable key is required",
		}
	}
	timetable, ok := api.timetables[*p.Key]
	if !ok {
		return nil, &jrpc2.ErrorObject{
			Code:    TimetableNotFoundCode,
			Message: TimetableNotFoundMsg,
		}
	}
	delay, err := timetable.Delay()
	if err != nil {
		return nil, &jrpc2.ErrorObject{
			Code:    -32099,
			Message: jrpc2.ServerErrorMsg,
			Data:    err.Error(),
		}
	}
	return delay, nil
}

// GetParams contains the rpc parameters for the Get method.
type GetParams struct {
	// Key is the timetable key.
	Key *string `json:"key"`
}

// FromPositional parses the key from the positional parameters.
func (params *GetParams) FromPositional(args []interface{}) error {
	if len(args) != 1 {
		return errors.New("key parameter is required")
	}
	key := args[0].(string)
	params.Key = &key

	return nil
}

// Get returns a timetable by key.  An error is returned if the timetable
//  does not exist.
func (api *ApiV1) Get(params json.RawMessage) (interface{}, *jrpc2.ErrorObject) {
	p := new(GetParams)
	if err := jrpc2.ParseParams(params, p); err != nil {
		return nil, err
	}
	if p.Key == nil {
		return nil, &jrpc2.ErrorObject{
			Code:    jrpc2.InvalidParamsCode,
			Message: jrpc2.InvalidParamsMsg,
			Data:    "timetable key is required",
		}
	}
	timetable, ok := api.timetables[*p.Key]
	if !ok {
		return nil, &jrpc2.ErrorObject{
			Code:    TimetableNotFoundCode,
			Message: TimetableNotFoundMsg,
		}
	}
	return timetable, nil
}

// GetAll returns all existing timetables.
func (api *ApiV1) GetAll(params json.RawMessage) (interface{}, *jrpc2.ErrorObject) {
	timetables := make([]*Timetable, 0)
	for _, timetable := range api.timetables {
		timetables = append(timetables, timetable)
	}
	return timetables, nil
}

// InsertParams contains the rpc parameters for Insert method.
type InsertParams struct {
	// Key is the timetable key.
	// Id is the id of the task.
	// RunAt in the execution point of time of the task
	Key   *string `json:"key"`
	Id    *string `json:"id"`
	RunAt *string `json:"runAt"`
}

// FromPositional parse the key, id, and runAt from the positional parameters.
func (params *InsertParams) FromPositional(args []interface{}) error {
	if len(args) != 3 {
		return errors.New("key, id, and runAt parameters are required")
	}
	key := args[0].(string)
	id := args[1].(string)
	runAt := args[2].(string)
	params.Key = &key
	params.Id = &id
	params.RunAt = &runAt

	return nil
}

// Insert adds the task to the timetable schedule.
func (api *ApiV1) Insert(params json.RawMessage) (interface{}, *jrpc2.ErrorObject) {
	p := new(InsertParams)
	if err := jrpc2.ParseParams(params, p); err != nil {
		return nil, err
	}
	if p.Key == nil {
		return nil, &jrpc2.ErrorObject{
			Code:    jrpc2.InvalidParamsCode,
			Message: jrpc2.InvalidParamsMsg,
			Data:    "task key is required",
		}
	}
	if p.Id == nil {
		return nil, &jrpc2.ErrorObject{
			Code:    jrpc2.InvalidParamsCode,
			Message: jrpc2.InvalidParamsMsg,
			Data:    "task id is required",
		}
	}
	if p.RunAt == nil {
		return nil, &jrpc2.ErrorObject{
			Code:    jrpc2.InvalidParamsCode,
			Message: jrpc2.InvalidParamsMsg,
			Data:    "task runAt is required",
		}
	}

	var timetable *Timetable
	var ok bool

	if timetable, ok = api.timetables[*p.Key]; !ok {
		timetable = NewTimetable(*p.Key)
		api.timetables[*p.Key] = timetable
	}
	if err := timetable.Insert(&Task{Id: *p.Id, RunAt: *p.RunAt}); err != nil {
		return nil, &jrpc2.ErrorObject{
			Code:    -32099,
			Message: jrpc2.ServerErrorMsg,
			Data:    err.Error(),
		}
	}
	if _, err := timetable.Save(api.model); err != nil {
		log.Println(err)
		return nil, &jrpc2.ErrorObject{
			Code:    -32099,
			Message: jrpc2.ServerErrorMsg,
			Data:    err.Error(),
		}
	}

	return 0, nil
}

// RemoveParams contains the rpc parameters for the Remove method.
type RemoveParams struct {
	// Key is queue id.
	// RunAt is the execution point in time of the task.
	Key   *string `json:"key"`
	RunAt *string `json:"runAt"`
}

// NextParams contains the rpc parameters for the Next method.
type NextParams struct {
	Key *string `json:"key"`
}

// FromPositional parse the key positional parameter.
func (params *NextParams) FromPositional(args []interface{}) error {
	if len(args) != 1 {
		return errors.New("key is required")
	}
	key := args[0].(string)
	params.Key = &key

	return nil
}

// Next returns the next scheduled task from the timetable.
func (api *ApiV1) Next(params json.RawMessage) (interface{}, *jrpc2.ErrorObject) {
	p := new(InsertParams)
	if err := jrpc2.ParseParams(params, p); err != nil {
		return nil, err
	}
	if p.Key == nil {
		return nil, &jrpc2.ErrorObject{
			Code:    jrpc2.InvalidParamsCode,
			Message: jrpc2.InvalidParamsMsg,
			Data:    "task key is required",
		}
	}
	timetable, ok := api.timetables[*p.Key]
	if !ok {
		return nil, &jrpc2.ErrorObject{
			Code:    TimetableNotFoundCode,
			Message: TimetableNotFoundMsg,
		}
	}
	return timetable.Next(), nil
}

// FromPositional parses the key and runAt from the positional
// parameters.
func (params *RemoveParams) FromPositional(args []interface{}) error {
	if len(args) != 2 {
		return errors.New("key, and runAt parameters are required")
	}
	key := args[0].(string)
	runAt := args[1].(string)
	params.Key = &key
	params.RunAt = &runAt

	return nil
}

// Remove removes the task from the timetable
func (api *ApiV1) Remove(params json.RawMessage) (interface{}, *jrpc2.ErrorObject) {
	p := new(RemoveParams)
	if err := jrpc2.ParseParams(params, p); err != nil {
		return nil, err
	}
	if p.Key == nil {
		return nil, &jrpc2.ErrorObject{
			Code:    jrpc2.InvalidParamsCode,
			Message: jrpc2.InvalidParamsMsg,
			Data:    "task key is required",
		}
	}
	if p.RunAt == nil {
		return nil, &jrpc2.ErrorObject{
			Code:    jrpc2.InvalidParamsCode,
			Message: jrpc2.InvalidParamsMsg,
			Data:    "task runAt is required",
		}
	}

	timetable, ok := api.timetables[*p.Key]
	if !ok {
		return nil, &jrpc2.ErrorObject{
			Code:    TimetableNotFoundCode,
			Message: TimetableNotFoundMsg,
		}
	}

	if err := timetable.Remove(*p.RunAt); err != nil {
		return -1, nil
	}
	if _, err := timetable.Save(api.model); err != nil {
		log.Println(err)
		return -1, &jrpc2.ErrorObject{
			Code:    -32099,
			Message: jrpc2.ServerErrorMsg,
			Data:    err.Error(),
		}
	}
	return 0, nil
}

// NewApiV1 returns a new api version 1 rpc api instance
func NewApiV1(model Model, s *jrpc2.Server) *ApiV1 {
	api := &ApiV1{model, make(map[string]*Timetable)}
	timetables, err := model.FetchAll()
	if err != nil {
		log.Fatal(err)
	}
	for _, timetable := range timetables {
		v, _ := timetable.(*Timetable)
		api.timetables[v.Key] = v
	}

	s.Register("delay", jrpc2.Method{Method: api.Delay})
	s.Register("get", jrpc2.Method{Method: api.Get})
	s.Register("getAll", jrpc2.Method{Method: api.GetAll})
	s.Register("insert", jrpc2.Method{Method: api.Insert})
	s.Register("next", jrpc2.Method{Method: api.Next})
	s.Register("remove", jrpc2.Method{Method: api.Remove})

	return api
}
