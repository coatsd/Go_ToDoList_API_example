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

var toDoList struct {
	sync.Mutex
	items []toDo
}
// returns the index value, a pointer to the toDo item found, and an error
func findItem(id int) (int, *toDo, error) {
	var toDoPt *toDo
	index := -1
	for i, td := range toDoList.items {
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

func ToDoHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		items, getErr := getToDo(r)
		if getErr != nil {
			handleError(w, getErr)
		} else {
			w.WriteHeader(http.StatusOK)
			toDoJson, _ := json.Marshal(*items)
			w.Write([]byte(toDoJson))
		}
	case "POST":
		item, postErr := postToDo(r)
		if postErr != nil {
			handleError(w, postErr)
		} else {
			w.WriteHeader(http.StatusOK)
			toDoJson, _ := json.Marshal(*item)
			w.Write([]byte(toDoJson))
		}
	case "PUT":
		item, putErr := putToDo(r)
		if putErr != nil {
			handleError(w, putErr)
		} else {
			w.WriteHeader(http.StatusOK)
			toDoJson, _ := json.Marshal(*item)
			w.Write([]byte(toDoJson))
		}
	case "PATCH":
		item, patchErr := patchToDo(r)
		if patchErr != nil {
			handleError(w, patchErr)
		} else {
			w.WriteHeader(http.StatusOK)
			toDoJson, _ := json.Marshal(*item)
			w.Write([]byte(toDoJson))
		}
	case "DELETE":
		item, delErr := deleteToDo(r)
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

func getToDo(r *http.Request) (*[]toDo, error) {
	toDoList.Lock()
	defer toDoList.Unlock()

	return &toDoList.items, nil
}

func postToDo(r *http.Request) (*toDo, error) {
	toDoPt, err := extractToDoBody(r)
	if err != nil {
		return toDoPt, err
	}
	toDoList.Lock()
	defer toDoList.Unlock()
	toDoList.items = append(toDoList.items, *toDoPt)
	return &toDoList.items[len(toDoList.items) - 1], nil
}

func putToDo(r *http.Request) (*toDo, error) {
	toDoPt, err := extractToDoBody(r)
	if err != nil {
		return toDoPt, err
	}
	toDoList.Lock()
	defer toDoList.Unlock()
	index, toDoItem, findErr := findItem(toDoPt.Id)
	if findErr != nil {
		return toDoItem, findErr
	}
	toDoList.items[index] = *toDoPt
	return &toDoList.items[index], nil
}

func patchToDo(r *http.Request) (*toDo, error) {
	toDoPt, err := extractToDoBody(r)
	if err != nil {
		return toDoPt, err
	}
	toDoList.Lock()
	defer toDoList.Unlock()
	index, toDoItem, findErr := findItem(toDoPt.Id)
	if findErr != nil {
		return toDoItem, findErr
	}
	toDoList.items[index] = *toDoPt
	return &toDoList.items[index], nil
}

func deleteToDo(r *http.Request) (*toDo, error) {
	toDoPt, err := extractToDoBody(r)
	if err != nil {
		return toDoPt, err
	}
	toDoList.Lock()
	defer toDoList.Unlock()
	index, toDoItem, findErr := findItem(toDoPt.Id)
	if findErr != nil {
		return toDoItem, findErr
	}
	toDoList.items = append(toDoList.items[:index], toDoList.items[index+1:]...)
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

func handleRequests() {
	fmt.Println("Request handler initialized. ctrl+c to exit")
	http.HandleFunc("/api/todo/", ToDoHandler)
	err := http.ListenAndServe(":10000", nil)
	if err != nil {
		panic(err)
	}
}

func main() {
	toDoList.items = []toDo{
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

	handleRequests()
}