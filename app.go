package main

import (
	// "fmt"
	"encoding/json"
	"errors"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rollbar/rollbar-go"
)

func main() {
	// http.HandleFunc("/", homeHandler)
	// http.HandleFunc("/contact", contactHandler)
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/feeling", feelingHandler)
	http.ListenAndServe(":3000", nil)

	r := mux.NewRouter()
	r.HandleFunc("/signup", SignupHandler).Methods("POST")
	r.HandleFunc("/login", LoginHandler).Methods("POST")

	rollbar.SetToken("db0eeaf619584284873f3ba1fbc012e1")
	rollbar.SetEnvironment("production")
	rollbar.SetServerRoot()
	rollbar.Error(rollbar.ERR)
}

// SignupHandler is a handler for handling signup requests
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get the signup data
	decoder := json.NewDecoder(r.Body)
	var data SignupData
	err := decoder.Decode(&data)
	if err != nil {
		rollbar.Error(rollbar.ERR, err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the signup data
	if data.Username == "" || data.Password == "" || data.Email == "" {
		rollbar.Error(rollbar.ERR, errors.New("Missing username, password, or email"))
		http.Error(w, "Missing username, password, or email", http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		rollbar.Error(rollbar.ERR, errors.New("Error hashing password"))
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	// Save the user to the database
	user := User{
		Username: data.Username,
		Password: string(hashedPassword),
		Email:    data.Email,
	}
	err = user.Save()
	if err != nil {
		rollbar.Error(rollbar.ERR, errors.New("Error saving user to database"))
		http.Error(w, "Error saving user to database", http.StatusInternalServerError)
		return
	}

	// Create a JSON response with the user data
	userData := struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
	jsonResponse, err := json.Marshal(userData)
	if err != nil {
		rollbar.Error(rollbar.ERR, err)
		http.Error(w, "Error marshalling user data", http.StatusInternalServerError)
		return
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

// LoginHandler is a handler for handling login requests
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get the login data
	decoder := json.NewDecoder(r.Body)
	var data LoginData
	err := decoder.Decode(&data)
	if err != nil {
		rollbar.Error(rollbar.ERR, err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the login data
	if data.Username == "" || data.Password == "" {
		rollbar.Error(rollbar.ERR, err)
		http.Error(w, "Missing username or password", http.StatusBadRequest)
		return
	}

	// Get the user from the database
	user, err := User.GetByUsername(data.Username)
	if err != nil {
		rollbar.Error(rollbar.ERR, err)
		http.Error(w, "Error getting user from database", http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Check if the password is correct
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	if err != nil {
		rollbar.Error(rollbar.ERR, err)
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Create a JSON response with the user data
	userData := struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
	jsonResponse, err := json.Marshal(userData)
	if err != nil {
		rollbar.Error(rollbar.ERR, err)
		http.Error(w, "Error marshalling user data", http.StatusInternalServerError)
		return
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

type Feeling struct {
	Name        string
	Description string
	Filename    string
}

var feelings []Feeling = []Feeling{
	Feeling{
		Name:        "happy",
		Description: "a feeling of joy and pleasure",
		Filename:    "how-to-be-happy.jpeg",
	},
	Feeling{
		Name:        "sad",
		Description: "a feeling of unhappiness and melancholy",
		Filename:    "sad.jpg",
	},
	Feeling{
		Name:        "angry",
		Description: "a feeling of extreme annoyance or frustration",
		Filename:    "angry.jpg",
	},
	Feeling{
		Name:        "confused",
		Description: "a feeling of being perplexed or uncertain",
		Filename:    "confused.jpg",
	},
	Feeling{
		Name:        "surprised",
		Description: "a feeling of being astonished or shocked",
		Filename:    "surprised.jpg",
	},
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// parse the HTML template
	t, err := template.ParseFiles("index.html")
	if err != nil {
		rollbar.Error(rollbar.ERR, err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// execute the template, writing the result to the response
	err = t.Execute(w, nil)
	if err != nil {
		rollbar.Error(rollbar.ERR, err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func feelingHandler(w http.ResponseWriter, r *http.Request) {
	// parse the request form
	err := r.ParseForm()
	if err != nil {
		rollbar.Error(rollbar.ERR, err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// get the feeling from the form
	feeling := r.FormValue("feeling")

	// find the corresponding Feeling struct
	var f Feeling
	for _, v := range feelings {
		if v.Name == feeling {
			f = v
			break
		}
	}

	// read the image file
	img, err := ioutil.ReadFile(f.Filename)
	if err != nil {
		rollbar.Error(rollbar.ERR, err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// set the content type of the response to the type of the image
	w.Header().Set("Content-Type", "image/jpeg")

	// write the image to the response
	_, err = w.Write(img)
	if err != nil {
		rollbar.Error(rollbar.ERR, err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
