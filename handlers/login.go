package handlers


import (
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"time"
	"os"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/buyubaya/go-blog/interfaces"
	"github.com/buyubaya/go-blog/helpers"
)


type LoginRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}


func Login(a interfaces.IApp) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			helpers.ERROR(w, http.StatusInternalServerError, err)
			return
		}


		loginRequest := &LoginRequest{}
		err = json.Unmarshal(body, loginRequest)
		if err != nil {
			helpers.ERROR(w, http.StatusInternalServerError, err)
			return
		}


		// JWT
		// Create a new token object, specifying signing method and the claims
		// you would like it to contain.
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"email": loginRequest.Email,
			"password": loginRequest.Password,
			"exp": time.Now().Add(time.Hour * 1).Unix(),
		})

		// Sign and get the complete encoded token as a string using the secret
		tokenString, err := token.SignedString([]byte(os.Getenv("API_SECRET")))


		// PARSE TOKEN
		token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
		
			return []byte(os.Getenv("API_SECRET")), nil
		})
		
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			fmt.Println(claims["email"], claims["password"])
		} else {
			fmt.Println(err)
		}

		fmt.Println("LOGIN", loginRequest, tokenString)
	}
}