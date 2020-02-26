package todo

import (
	"errors"
	"github.com/rs/xid"
	"sync"
)

var (
	list []Todo
	mtx  sync.RWMutex
	once sync.Once
)

func init() {
	once.Do(initialiseList)
}

func initialiseList() {
	list = []Todo{}
}

type Todo struct {
	ID       string `json:"id"`
	Message  string `json:"message"`
	Complete bool   `json:"complete"`

}

func Get() []Todo {
	return list
}
func Add(message string) string {
	t := newTodo(message)
	mtx.Lock()
	list = append(list, t)
	mtx.Unlock()
	return t.ID
}
func Delete(id string) error {
	location, err := findTodoLocation(id)
	if err != nil {
		return err
	}
	removeElementByLocation(location)
	return nil
}

func removeElementByLocation(i int) {
	mtx.Lock()
	list = append(list[:i], list[i+1:]...)
	mtx.Unlock()
}
func complete(id string) error {
	location, err := findTodoLocation(id)
	if err != nil {
		return err
	}
	setTodoCompleteByLocation(location)
	return nil
}

func setTodoCompleteByLocation(location int) {
	mtx.Lock()
	list[location].Complete = true
	mtx.Unlock()
}

func findTodoLocation(id string) (int, error) {
	mtx.RLock()
	defer mtx.RUnlock()
	for i, t := range list {
		if isMatchingID(t.ID, id) {
			return i, nil
		}
	}
	return 0, errors.New("Could not Find todo based on id")
}

func isMatchingID(a string, b string) bool {
	return a == b

}

func newTodo(msg string) Todo {
	return Todo{
		ID:       xid.New().String(),
		Message:  msg,
		Complete: false,
	}
}
