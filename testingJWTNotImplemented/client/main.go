package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var mySigningKey = []byte("thisisalongphrasetomakethetokenmoresecure")

func homePage(w http.ResponseWriter, r *http.Request) {
	validToken, err := GenerateJWT()
	if err != nil {
		fmt.Fprintf(w, err.Error())
		panic(err)
	}

	client := &http.Client{}

	req, _ := http.NewRequest("GET", "http://localhost:3299/", nil)
	req.Header.Set("Token", validToken)
	res, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(w, "Error : %s", err.Error())
		panic(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Fprint(w, err.Error())
		panic(err)
	}

	fmt.Fprintf(w, string(body))

	// fmt.Fprintf(w, validToken)
}

func GenerateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorize"] = true
	claims["user"] = "Authorized Personnel"
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Errorf("Something went wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil

}

func handleRequests() {
	http.HandleFunc("/", homePage)

	log.Fatal(http.ListenAndServe(":3299", nil))
}

func main() {
	fmt.Println("This is the Client view")

	// no need here after initial test, its been done in homePage func
	// tokenString, err := GenerateJWT()

	// if err != nil {
	// 	fmt.Println("Error generating token string")
	// }

	// fmt.Printf("tokenString : %v\n\n", tokenString)

	handleRequests()

}
