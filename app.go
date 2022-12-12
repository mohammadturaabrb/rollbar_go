package main

import (
	// "fmt"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/rollbar/rollbar-go"
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

	rollbar.SetToken("db0eeaf619584284873f3ba1fbc012e1")
	rollbar.SetEnvironment("production")
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
