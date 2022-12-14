package main

import (
	// "fmt"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/rollbar/rollbar-go"
	"github.com/gorilla/mux"
)

// func homeHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "text/html; charset=utf-8")
// 	fmt.Fprint(w, "<h1> Welcome to my webiste </h1>")
// }

// func contactHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "text/html; charset=utf-8")
// 	fmt.Fprint(w, "<h1> Contact Page </h1> <p>To get in touch, email me at <a href=\"mailto:support@rollbar.com\">support@rollbar.com</a>.</p>")
// }

func main() {
	// http.HandleFunc("/", homeHandler)
	// http.HandleFunc("/contact", contactHandler)
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/feeling", feelingHandler)
	http.ListenAndServe(":3000", nil)

	router := mux.NewRouter()
	router := ("/signup", signup).Methods("POST")


	rollbar.SetToken("db0eeaf619584284873f3ba1fbc012e1")
	rollbar.SetEnvironment("production")
}

func signup(w http.ResponseWriter, r *http.Request) {
	// Get the signup data from the request body
	decoder := json.NewDecoder(r.Body)
	var userData UserData
	err := decoder.Decode(&userData)
	if err != nil {
	  // Return an error response if the signup data is invalid
	  return
	}
  
	// Validate the signup data
	if userData.Username == "" || userData.Password == "" {
	  // Return an error response if the signup data is invalid
	  return
	}
  
	// Hash the password for security
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if err != nil {
	  // Return an error response if there is a problem with the password
	  return
	}
  
	// Store the signup data in a database
	db, err := sql.Open("mysql", "user:password@/database")
	if err != nil {
	  // Return an error response if there is a problem with the database
	  return
	}
	defer db.Close()
  
	stmt, err := db.Prepare("INSERT INTO users(username, password) VALUES(?, ?)")
	if err != nil {
	  // Return an error response if there is a problem with the database
	  return
	}
	defer stmt.Close()
  
	_, err = stmt.Exec(userData.Username, hashedPassword)
	if err != nil {
	  // Return an error response if there is a problem with the database
  
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
