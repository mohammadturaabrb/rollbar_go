package main

import (
    "os"

    middleware "rollbar_go/middleware"
    routes "rollbar_go/routes"

    "github.com/gin-gonic/gin"
    _ "github.com/heroku/x/hmetrics/onload"
	"github.com/rollbar/rollbar-go"
)

func main() {

	rollbar.SetToken("db0eeaf619584284873f3ba1fbc012e1")
	rollbar.SetEnvironment("production")
	rollbar.SetServerRoot("https://github.com/mohammadturaabrb/rollbar_go")
	rollbar.SetCodeVersion("1.0.0")

    port := os.Getenv("PORT")

    if port == "" {
        port = "8000"
    }

    router := gin.New()
    router.Use(gin.Logger())
    routes.UserRoutes(router)

    router.Use(middleware.Authentication())

    // API-2
    router.GET("/api-1", func(c *gin.Context) {

        c.JSON(200, gin.H{"success": "Access granted for api-1"})

    })

    // API-1
    router.GET("/api-2", func(c *gin.Context) {
        c.JSON(200, gin.H{"success": "Access granted for api-2"})
    })

    router.Run(":" + port)
}




























// func homePage(res http.ResponseWriter, req *http.Request) {
// 	http.ServeFile(res, req, "index.html")
// }

// // var feelings []Feeling = []Feeling{
// // 	Feeling{
// // 		Name:        "happy",
// // 		Description: "a feeling of joy and pleasure",
// // 		Filename:    "how-to-be-happy.jpeg",
// // 	},
// // 	Feeling{
// // 		Name:        "sad",
// // 		Description: "a feeling of unhappiness and melancholy",
// // 		Filename:    "sad.jpg",
// // 	},
// // 	Feeling{
// // 		Name:        "angry",
// // 		Description: "a feeling of extreme annoyance or frustration",
// // 		Filename:    "angry.jpg",
// // 	},
// // 	Feeling{
// // 		Name:        "confused",
// // 		Description: "a feeling of being perplexed or uncertain",
// // 		Filename:    "confused.jpg",
// // 	},
// // 	Feeling{
// // 		Name:        "surprised",
// // 		Description: "a feeling of being astonished or shocked",
// // 		Filename:    "surprised.jpg",
// // 	},
// // }

// // func indexHandler(w http.ResponseWriter, r *http.Request) {
// // 	t, err := template.ParseFiles("input.html")
// // 	if err != nil {
// // 		rollbar.Error(rollbar.ERR, err)
// // 		http.Error(w, err.Error(), http.StatusInternalServerError)
// // 		return
// // 	}

// // 	err = t.Execute(w, nil)
// // 	if err != nil {
// // 		rollbar.Error(rollbar.ERR, err)

// // 		http.Error(w, err.Error(), http.StatusInternalServerError)
// // 	}
// // }
// // func feelingHandler(w http.ResponseWriter, r *http.Request) {
// // 	err := r.ParseForm()
// // 	if err != nil {
// // 		rollbar.Error(rollbar.ERR, err)

// // 		http.Error(w, err.Error(), http.StatusInternalServerError)
// // 		return
// // 	}

// // 	feeling := r.FormValue("feeling")

// // 	var f Feeling
// // 	for _, v := range feelings {
// // 		if v.Name == feeling {
// // 			f = v
// // 			break
// // 		}
// // 	}

// // 	img, err := ioutil.ReadFile(f.Filename)
// // 	if err != nil {
// // 		rollbar.Error(rollbar.ERR, err)

// // 		http.Error(w, err.Error(), http.StatusInternalServerError)
// // 		return
// // 	}

// // 	w.Header().Set("Content-Type", "image/jpeg")

// // 	_, err = w.Write(img)
// // 	if err != nil {
// // 		rollbar.Error(rollbar.ERR, err)

// // 		http.Error(w, err.Error(), http.StatusInternalServerError)
// // 	}
// // }

