package elevator 
import (
    "fmt"
    "sync"
)

/*
This was made as dumb and simple as possible -- this is NOT a real elevator
*/
type Call struct {
    Floor     int // the floor number
    Direction int // the direction (1 for up, -1 for down)
}

type ElevatorCallList struct {
    mu sync.Mutex
    Calls []Call
}

func NewElevatorCallList() *ElevatorCallList {
	calls := make([]Call, 0)
	return &ElevatorCallList{
		Calls: calls,
	}
}

func (e *ElevatorCallList) Len() int {
	return len(e.Calls)
}

func (e *ElevatorCallList) Push(c Call) error {
	// no duplicates
    e.mu.Lock()
    defer e.mu.Unlock()
	for _, v := range e.Calls {
		if v.Floor == c.Floor && v.Direction == c.Direction {
			return fmt.Errorf("duplicate call on list: %v", c)
		}
	}
	e.Calls = append(e.Calls, c)
	return nil
}

// we use this for maintenance call and other elevator call overrides
func (e *ElevatorCallList) Prepend(c Call) error {
    e.mu.Lock()
    defer e.mu.Unlock()
	// needs to be a slice
	// so we make one of length 1
	cl := []Call{c}
	// then append the existing calls
	e.Calls = append(cl, e.Calls...)
	return nil
}

// we'll treat the call list like a queue
func (e *ElevatorCallList) PopLeft() *Call {
    e.mu.Lock()
    defer e.mu.Unlock()
    if len(e.Calls) == 0 {
        return nil
    }
	v, ll := (e.Calls)[0], (e.Calls)[1:]
	e.Calls = ll
	return &v
}