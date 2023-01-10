package main

import (
	"context"
	"fmt"
	"encoding/json"
	"log"
	"net/http"
	"time"
	"github.com/gorilla/mux"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"

	"github.com/rollbar/rollbar-go"
)

var SECRET_KEY = []byte("gosecretkey")

type User struct {
	FirstName string `json:"firstname" bson:"firstname"`
	LastName  string `json:"lastname" bson:"lastname"`
	Email     string `json:"email" bson:"email"`
	Password  string `json:"password" bson:"password"`
}

var client *mongo.Client

const uri = "mongodb://admin:Rollbar123!@ac-fyf3jtj-shard-00-02.aftpwhu.mongodb.net:27017"
//const uri = "ac-fyf3jtj-shard-00-02.aftpwhu.mongodb.net:27017"

func getHash(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
		rollbar.Error(err)
	}
	return string(hash)
}

func GenerateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	tokenString, err := token.SignedString(SECRET_KEY)
	if err != nil {
		log.Println("Error in JWT token generation")
		rollbar.Error(err)
		return "", err
	}
	return tokenString, nil
}

func userSignup(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var user User
	json.NewDecoder(request.Body).Decode(&user)
	user.Password = getHash([]byte(user.Password))
	collection := client.Database("userData").Collection("user")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx, user)
	json.NewEncoder(response).Encode(result)
	fmt.Println(result)
}

func userLogin(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var user User
	var dbUser User
	json.NewDecoder(request.Body).Decode(&user)
	collection := client.Database("userData").Collection("user")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&dbUser)
	rollbar.Error(collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&dbUser))

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"` + err.Error() + `"}`))
		rollbar.Error(http.StatusInternalServerError, err)
		return
	}
	userPass := []byte(user.Password)
	dbPass := []byte(dbUser.Password)

	passErr := bcrypt.CompareHashAndPassword(dbPass, userPass)

	if passErr != nil {
		log.Println(passErr)
		response.Write([]byte(`{"response":"Wrong Password!"}`))
		return
	}
	jwtToken, err := GenerateJWT()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"` + err.Error() + `"}`))
		rollbar.Error(http.StatusInternalServerError, err)
		return
	}
	response.Write([]byte(`{"token":"` + jwtToken + `"}`))

}

func main() {

	rollbar.SetToken("db0eeaf619584284873f3ba1fbc012e1")
	rollbar.SetEnvironment("production")
	rollbar.SetServerRoot("https://github.com/mohammadturaabrb/rollbar_go")
	log.Println("Starting the application")

	router := mux.NewRouter()
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI(uri))

	router.HandleFunc("/user/login", userLogin).Methods("POST")
	router.HandleFunc("/user/signup", userSignup).Methods("POST")

	err := http.ListenAndServe(":3000", router)
	if err != nil {
			rollbar.Error(err)
			log.Fatalln("There's an error with the server", err)
	}
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

