package main

import (
	"code.google.com/p/go.crypto/bcrypt"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	session    *mgo.Session
	collection *mgo.Collection
)

var cachedTemplates = map[string]*template.Template{}
var cachedMutex sync.Mutex

var funcs = template.FuncMap{
	"reverse": reverse,
}

func T(name string) *template.Template {
	cachedMutex.Lock()
	defer cachedMutex.unlock()

	//this is a two value assignment test for existence of a key
	//this returns a template and true or false, then, if true return t
	//see http://blog.golang.org/go-maps-in-action  - also has a good note on the use of "_"
	if t, ok := cachedTemplates[name]; ok {
		return t
	}

	t := template.New("_base.html").Funcs(funcs)

	t = template.Must(t.ParseFiles(
		"templates/_base.html",
		filepath.join("templates", name),
	))

}

var login = parseTemplate(
	"templates/_base.html",
	"templates/login.html",
)

type user struct {
	ID       bson.objectID `bson:"_id, omitempty"`
	Username string
	Password []byte
}

type Kitten struct {
	Id      bson.ObjectId `bson:"_id" json:"id"`
	Name    string        `bson:"Name" json:"name"`
	Picture string        `bson:"Picture" json:"picture"`
}

type KittenJSON struct {
	Kitten Kitten `json:"kitten"`
}

type KittensJSON struct {
	Kittens []Kitten `json:"kittens"`
}

func (u *user) SetPassword(password string) {
	hpass, err := bcrypt.GenerateFromPassword([]byet(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err) //this is a panic because bcrypt errors on invalid cost
	}
	u.Password = hpass
}

//Login validates and returns a user object if they exist in the database.
func Login(ctx *Context, username, password string) (u *user, err error) {
	err = ctx.C("users").Find(bson.M{"username": username}).One(&u)
	if err != nil {
		return
	}

	err = bcrypt.compareHashAndPassword(u.Password, []byte(password))
	if err != nil {
		u = nil
	}
	return
}

func login(w http.Responsewriter, req *httpRequest, ctx *context) (err error) {
	username, password := req.FormValue("username"), req.FormValue("password")

	//log in the user
	user, err := Login(ctx, username, password)

	_ = user //"_" believe is just a blank identifier
	return
}

func loginForm(w http.ResponseWriter, req *http.Request, ctx *Context) (err error) {
	err = login.Execute(w, nil)
}

func CreateKittenHandler(w http.ResponseWriter, r *http.Request) {

	var kittenJSON KittenJSON

	//decode income kitten from json
	err := json.NewDecoder(r.Body).Decode(&kittenJSON)
	if err != nil {
		panic(err)
	}
	kitten := kittenJSON.Kitten

	// Generate a random dimension for the kitten
	width := rand.Int() % 400
	height := rand.Int() % 400
	if width < 100 {
		width += 100
	}
	if height < 100 {
		height += 100
	}
	kitten.Picture = fmt.Sprintf("http://placekitten.com/%d/%d", width, height)

	// Store the new kitten in the database
	// First, let's get a new id
	obj_id := bson.NewObjectId()
	kitten.Id = obj_id

	err = collection.Insert(&kitten)
	if err != nil {
		panic(err)
	} else {
		log.Printf("Inserted new kitten %s with name %s", kitten.Id, kitten.Name)
	}

	j, err := json.Marshal(KittenJSON{Kitten: kitten})
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func KittensHandler(w http.ResponseWriter, r *http.Request) {
	// lets build up the kittens slice
	var mykittens []Kitten

	iter := collection.Find(nil).Iter()
	result := Kitten{}
	for iter.Next(&result) {
		mykittens = append(mykittens, result)
	}

	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(KittensJSON{Kittens: mykittens})
	if err != nil {
		panic(err)
	}

	w.Write(j)
	log.Println("Provided json")
}

func DeleteKittenHandler(w http.ResponseWriter, r *http.Request) {
	//grap the kitten's id from the incoming url
	var err error
	vars := mux.Vars(r)
	//id := vars["id"]

	//remove it from database
	id := bson.ObjectIdHex(vars["id"])
	err = collection.Remove(bson.M{"_id": id})
	if err != nil {
		log.Printf("Could not find kitten %s to delete", id)
	}
	w.WriteHeader(http.StatusNoContent)
}

func UpdateKittenHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	//grap the kitten's id from the incoming url
	vars := mux.Vars(r)
	id := bson.ObjectIdHex(vars["id"])

	//decode the incoming kitten json
	var kittenJSON KittenJSON
	err = json.NewDecoder(r.Body).Decode(&kittenJSON)
	if err != nil {
		panic(err)
	}

	// update the database
	err = collection.Update(bson.M{"_id": id},
		bson.M{
			"name": kittenJSON.Kitten.Name,
			"_id":  id, "picture": kittenJSON.Kitten.Picture,
		})
	if err == nil {
		log.Printf("Updated kitten %s name to %s", id, kittenJSON.Kitten.Name)
	} else {
		panic(err)
	}
	w.WriteHeader(http.StatusNoContent)
}

func KittensHandler2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"kittens": [
    {"id": 1, "name": "Bobby", "picture": "http://placekitten.com/200/200"},
    {"id": 2, "name": "Wally", "picture": "http://placekitten.com/200/200"}
  ]}`))
}

func main() {
	log.Print("starting server")
	ip := os.Getenv("IP")
	address := ip + GetPort()
	uri := "mongodb://johng:admin@ds055709.mongolab.com:55709/db_gem_01"
	log.Print("The uri is: " + uri)
	var err error

	r := mux.NewRouter()
	r.HandleFunc("/api/kittens", KittensHandler).Methods("GET")
	r.HandleFunc("/api/kittens", CreateKittenHandler).Methods("POST")
	http.Handle("/api/", r)

	//fs := http.FileServer(http.Dir("public"))

	http.Handle("/", http.FileServer(http.Dir("public")))

	log.Println("Starting mongo db session")

	maxWait := time.Duration(5 * time.Second)
	session, err = mgo.DialWithTimeout(uri, maxWait)
	log.Print(session)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	collection = session.DB("db_gem_01").C("kitten")

	log.Print("Listening on 8080")
	listen_err := http.ListenAndServe(address, nil)
	if listen_err != nil {
		log.Fatal("ListenAndServe Error: ", listen_err)
		return
	}

}

func GetPort() string {
	var port = os.Getenv("PORT")
	//set a default port if there is nothing in the environment
	if port == "" {
		port = "4747"
		fmt.Println("INFO: No PORT environemnt variable detected, default to " + port)
	}
	return ":" + port
}
