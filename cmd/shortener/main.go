package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
)

func mainPage(rw http.ResponseWriter, r *http.Request) {
	//body := fmt.Sprintf("Method: %s\r\n", r.Method)
	//body += "Header ===============\r\n"
	//for k, v := range r.Header {
	//	body += fmt.Sprintf("%s: %v\r\n", k, v)
	//}
	//body += "Query parameters ===============\r\n"
	//if err := r.ParseForm(); err != nil {
	//	rw.Write([]byte(err.Error()))
	//	return
	//}
	//for k, v := range r.Form {
	//	body += fmt.Sprintf("%s: %v\r\n", k, v)
	//}
	rw.WriteHeader(http.StatusCreated)
	//rw.Write([]byte(body))
}

func apiPage(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("Это страница /api."))
}

// Abs возвращает абсолютное значение.
// Например: 3.1 => 3.1, -3.14 => 3.14, -0 => 0.
func Abs(value float64) float64 {
	return math.Abs(value)
}

// Client — пользователь в системе.
type Client struct {
	FirstName string
	LastName  string
}

// FullName возвращает имя и фамилию пользователя.
func (u Client) FullName() string {
	return u.FirstName + " " + u.LastName
}

// Relationship определяет положение в семье.
type Relationship string

// Person описывает конкретного человека в семье.
type Person struct {
	FirstName string
	LastName  string
	Age       int
}

// Возможные роли в семье.
const (
	Father      = Relationship("father")
	Mother      = Relationship("mother")
	Child       = Relationship("child")
	GrandMother = Relationship("grandMother")
	GrandFather = Relationship("grandFather")
)

// Family описывает семью.
type Family struct {
	Members map[Relationship]Person
}

var (
	// ErrRelationshipAlreadyExists возвращает ошибку, если роль уже занята.
	ErrRelationshipAlreadyExists = errors.New("relationship already exists")
)

// AddNew добавляет нового члена семьи.
// Если в семье ещё нет людей, создаётся пустая мапа.
// Если роль уже занята, метод выдаёт ошибку.
func (f *Family) AddNew(r Relationship, p Person) error {
	if f.Members == nil {
		f.Members = map[Relationship]Person{}
	}
	if _, ok := f.Members[r]; ok {
		return ErrRelationshipAlreadyExists
	}
	f.Members[r] = p
	return nil
}

// User — основной объект для теста.
type User struct {
	ID        string
	FirstName string
	LastName  string
}

func StatusHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	// намеренно добавлена ошибка в JSON
	rw.Write([]byte(`{"status":"ok}`))
}

// UserViewHandler — хендлер, который нужно протестировать.
func UserViewHandler(users map[string]User) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			//http.Error(rw, "user_id is empty", http.StatusBadRequest)
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(`{"message": "user_id is empty"}`))
			return
		}

		user, ok := users[userID]
		if !ok {
			//http.Error(rw, "user not found", http.StatusNotFound)
			rw.WriteHeader(http.StatusNotFound)
			rw.Write([]byte(`{"message": "user not found"}`))
			return
		}

		jsonUser, err := json.Marshal(user)
		if err != nil {
			//http.Error(rw, "can't provide a json. internal error",
			//	http.StatusInternalServerError)
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(`{"message": "can't provide a json. internal error"}`))
			return
		}

		rw.WriteHeader(http.StatusOK)
		rw.Write(jsonUser)
	}
}

func main() {
	// 1 task
	// 2 fragment
	v := Abs(3)
	fmt.Println(v)
	// 2 fragment
	u := Client{
		FirstName: "Misha",
		LastName:  "Popov",
	}
	fmt.Println(u.FullName())
	// 3 fragment
	f := Family{}
	err := f.AddNew(Father, Person{
		FirstName: "Misha",
		LastName:  "Popov",
		Age:       56,
	})
	fmt.Println(f, err)
	err = f.AddNew(Father, Person{
		FirstName: "Drug",
		LastName:  "Mishi",
		Age:       57,
	})
	fmt.Println(f, err)

	// 4 task
	users := make(map[string]User)
	u1 := User{
		ID:        "u1",
		FirstName: "Misha",
		LastName:  "Popov",
	}
	u2 := User{
		ID:        "u2",
		FirstName: "Sasha",
		LastName:  "Popov",
	}
	users["u1"] = u1
	users["u2"] = u2

	mux := http.NewServeMux()
	mux.HandleFunc(`/api/`, apiPage)
	mux.HandleFunc(`/status`, StatusHandler)
	mux.HandleFunc(`/users`, UserViewHandler(users))
	mux.HandleFunc(`/`, mainPage)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
