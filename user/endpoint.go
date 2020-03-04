package user

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Route ...
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes Route Slice
type Routes []Route

var userRoutes = Routes{
	Route{
		"GetUserList",
		"Get",
		"/",
		nil, //GetUsers,
	},
	Route{
		"CreateUser",
		"POST",
		"/",
		CreateUser,
	},
	Route{
		"GerUser",
		"GET",
		"/{uuid}",
		GetUser,
	},
	Route{
		"UpdateUser",
		"POST",
		"/{uuid}",
		nil, //UpdateUser,
	},
}

// NewRouter is mux router constructor
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range userRoutes {
		var handler http.Handler
		handler = http.Handler(route.HandlerFunc)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

// UsersEndpoint ...
func UsersEndpoint(w http.ResponseWriter, r *http.Request) {
	m := NewRouter()
	m.ServeHTTP(w, r)
}

// YourselfEndpoint ...
func YourselfEndpoint(w http.ResponseWriter, r *http.Request) {
	type Info struct {
		ContentType string `json:"contentType"`
		Method      string `json:"method"`
	}

	info := &Info{}

	switch r.Header.Get("Content-Type") {
	case "application/json":
		// only application/json
		log.Printf("ContentType: %s", r.Header.Get("Content-Type"))
		info.ContentType = r.Header.Get("Content-Type")
	default:
		// TODO: error handler
		log.Printf("ContentType: %s", r.Header.Get("Content-Type"))
		info.ContentType = r.Header.Get("Content-Type")
	}

	switch r.Method {
	case http.MethodGet:
		// GET
		log.Printf("ContentType: %s", r.Method)
		info.Method = r.Method
	case http.MethodPost:
		// POST
		log.Printf("ContentType: %s", r.Method)
		info.Method = r.Method
	case http.MethodPut:
		// PUT
		log.Printf("ContentType: %s", r.Method)
		info.Method = r.Method
	case http.MethodDelete:
		// DELETE
		log.Printf("ContentType: %s", r.Method)
		info.Method = r.Method
	default:
		http.Error(w, "405 - Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	res, err := json.Marshal(info)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "JSON Marshal Error.\n")
		log.Printf("JSON Marshal Error: %+v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)

	return
}

// GetUser ...
func GetUser(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	params := mux.Vars(r)
	uuid := params["uuid"]
	user := &User{
		UUID: uuid,
	}

	if err := user.Get(c); err != nil {
		log.Printf("err: %+v", err)
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, "Not Fountd Error.\n")
		log.Printf("Not Fountd Error: %s", uuid)
		return
	}

	res, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "JSON Marshal Error.\n")
		log.Printf("JSON Marshal Error: %+v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)

	return
}

// CreateUser is a function to create user information
func CreateUser(w http.ResponseWriter, r *http.Request) {
	user := &User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		// リクエストの json がでコードできなかった時
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "JSON Decord Error.\n")
		log.Printf("JSON Decord Error: %+v", err)
		return
	}
	if user.UUID == "" {
		// json 内に UUID プロパティがなかった時
		// UUID の生成
		u, err := uuid.NewRandom()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, "Generate UUID Error.\n")
			log.Printf("Generate UUID Error: %+v", err)
			return
		}
		user.UUID = u.String()
	}

	c := r.Context()

	// firestore に既に登録されているか
	if err := user.Get(c); err == nil {
		// エラーじゃないときは多分いる
		// TODO: code = NotFound desc で判定
		w.WriteHeader(http.StatusForbidden)
		io.WriteString(w, fmt.Sprintf("Exist User UUID: %s", user.UUID))
		log.Printf("Exist User UUID: %s", user.UUID)
		return
	}

	if err := user.Create(c); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "Create User Error.\n")
		log.Printf("Create User Error: %+v", err)
		return
	}

	log.Printf("Crented User UUID: %s", user.UUID)

	res, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "JSON Marshal Error.\n")
		log.Printf("JSON Marshal Error: %+v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)

	return
}
