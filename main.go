package main

import (
	"fmt"
	"errors"
	"net/http"
	"encoding/json"
	"sync"
	"io/ioutil"
)

type toDo struct {
	Id int `json:"id"`
	Title string `json:"title"`
	Desc string `json:"description"`
	IsDone bool `json:"isDone"`
}
func makeToDo(id int, title, desc string, isDone bool) toDo {
	return toDo{Id: id, Title: title, Desc: desc, IsDone: isDone,}
}
func newToDo(id int, title, desc string, isDone bool) *toDo {
	return &toDo{Id: id, Title: title, Desc: desc, IsDone: isDone,}
}

type toDoList struct {
	sync.Mutex
	items []toDo
}
func makeToDoList() toDoList {
	return toDoList{}
}
func (l *toDoList) fetchData() {
	l.items = []toDo{
		toDo{
			Id: 1,
			Title: "Test 1",
			Desc: "Description 1",
			IsDone: false,
		},
		toDo{
			Id: 2,
			Title: "Test 2",
			Desc: "Description 2",
			IsDone: false,
		},
		toDo{
			Id: 3,
			Title: "Test 3",
			Desc: "Description 3",
			IsDone: false,
		},
	}
}
// returns the index value, a pointer to the toDo item found, and an error
func (l *toDoList) findItem(id int) (int, *toDo, error) {
	var toDoPt *toDo
	index := -1
	for i, td := range l.items {
		if td.Id == id {
			index = i
			toDoPt = &td
			break
		}
	}
	if toDoPt == nil {
		return index, toDoPt, errors.New("ToDo item not found")
	}
	return index, toDoPt, nil
}

func (l *toDoList) ToDoHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		items, getErr := l.getToDo(r)
		if getErr != nil {
			handleError(w, getErr)
		} else {
			w.WriteHeader(http.StatusOK)
			toDoJson, _ := json.Marshal(*items)
			w.Write([]byte(toDoJson))
		}
	case "POST":
		item, postErr := l.postToDo(r)
		if postErr != nil {
			handleError(w, postErr)
		} else {
			w.WriteHeader(http.StatusOK)
			toDoJson, _ := json.Marshal(*item)
			w.Write([]byte(toDoJson))
		}
	case "PUT":
		item, putErr := l.putToDo(r)
		if putErr != nil {
			handleError(w, putErr)
		} else {
			w.WriteHeader(http.StatusOK)
			toDoJson, _ := json.Marshal(*item)
			w.Write([]byte(toDoJson))
		}
	case "PATCH":
		item, patchErr := l.patchToDo(r)
		if patchErr != nil {
			handleError(w, patchErr)
		} else {
			w.WriteHeader(http.StatusOK)
			toDoJson, _ := json.Marshal(*item)
			w.Write([]byte(toDoJson))
		}
	case "DELETE":
		item, delErr := l.deleteToDo(r)
		if delErr != nil {
			handleError(w, delErr)
		} else {
			w.WriteHeader(http.StatusOK)
			toDoJson, _ := json.Marshal(*item)
			w.Write([]byte(toDoJson))
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("The method is not supported by this resource"))
	}
}

func (l *toDoList) getToDo(r *http.Request) (*[]toDo, error) {
	l.Lock()
	defer l.Unlock()

	return &l.items, nil
}

func (l *toDoList) postToDo(r *http.Request) (*toDo, error) {
	toDoPt, err := extractToDoBody(r)
	if err != nil {
		return toDoPt, err
	}
	l.Lock()
	defer l.Unlock()
	l.items = append(l.items, *toDoPt)
	return &l.items[len(l.items) - 1], nil
}

func (l *toDoList) putToDo(r *http.Request) (*toDo, error) {
	toDoPt, err := extractToDoBody(r)
	if err != nil {
		return toDoPt, err
	}
	l.Lock()
	defer l.Unlock()
	index, toDoItem, findErr := l.findItem(toDoPt.Id)
	if findErr != nil {
		return toDoItem, findErr
	}
	l.items[index] = *toDoPt
	return &l.items[index], nil
}

// Identical to putToDo, but may modify it later.
// Patch updates record columns; Put replaces records completely, in theory.
// Therefore, it might be good to modify this to do the same.
func (l *toDoList) patchToDo(r *http.Request) (*toDo, error) {
	toDoPt, err := extractToDoBody(r)
	if err != nil {
		return toDoPt, err
	}
	l.Lock()
	defer l.Unlock()
	index, toDoItem, findErr := l.findItem(toDoPt.Id)
	if findErr != nil {
		return toDoItem, findErr
	}
	l.items[index] = *toDoPt
	return &l.items[index], nil
}

func (l *toDoList) deleteToDo(r *http.Request) (*toDo, error) {
	toDoPt, err := extractToDoBody(r)
	if err != nil {
		return toDoPt, err
	}
	l.Lock()
	defer l.Unlock()
	index, toDoItem, findErr := l.findItem(toDoPt.Id)
	if findErr != nil {
		return toDoItem, findErr
	}
	l.items = append(l.items[:index], l.items[index+1:]...)
	return toDoItem, nil
}

func handleError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(err.Error()))
}

// extracts the body of a request (r) and parses the data as a toDo struct
func extractToDoBody(r *http.Request) (*toDo, error) {
	var toDoPt *toDo
	body, readErr := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if readErr != nil {
		return toDoPt, readErr
	}
	var newToDo toDo
	if umErr := json.Unmarshal(body, &newToDo); umErr != nil {
		return toDoPt, umErr
	}
	return &newToDo, nil
}

func handleRequests(l toDoList) {
	fmt.Println("Request handler initialized. ctrl+c to exit")
	http.HandleFunc("/api/todo/", l.ToDoHandler)
	err := http.ListenAndServe(":10000", nil)
	if err != nil {
		panic(err)
	}
}

func main() {
	tdl := makeToDoList()
	tdl.fetchData()

	handleRequests(tdl)
}