package UserAPI


import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"errors"

	"github.com/buyubaya/go-blog/models"
	"github.com/buyubaya/go-blog/interfaces"
	"github.com/buyubaya/go-blog/helpers"
	"firebase.google.com/go/auth"
)


type CreateUserRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}


func CreateUser(a interfaces.IApp) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			helpers.ERROR(w, http.StatusInternalServerError, err)
			return
		}


		createUserRequest := &CreateUserRequest{}
		err = json.Unmarshal(body, createUserRequest)
		if err != nil {
			helpers.ERROR(w, http.StatusInternalServerError, err)
			return
		}

		if (createUserRequest.Email == "") {
			helpers.ERROR(w, http.StatusInternalServerError, errors.New("Email is required"))
			return
		}

		if (createUserRequest.Password == "") {
			helpers.ERROR(w, http.StatusInternalServerError, errors.New("Password is required"))
			return
		}

		if (createUserRequest.ConfirmPassword == "") {
			helpers.ERROR(w, http.StatusInternalServerError, errors.New("Confirm Password is required"))
			return
		}

		if (createUserRequest.Password != createUserRequest.ConfirmPassword) {
			helpers.ERROR(w, http.StatusInternalServerError, errors.New("Password and Confirm Password does not match"))
			return
		}


		// CREATE FIREBASE USER
		newFirebaseUser := (&auth.UserToCreate{}).
			Email(createUserRequest.Email).
			Password(createUserRequest.Password)

		firebaseUser, err := (*a.GetFirebaseApp()).CreateUser(newFirebaseUser)
		if err != nil {
			helpers.ERROR(w, http.StatusInternalServerError, err)
			return
		}


		// CREATE APP USER
		newAppUser := &models.User{
			FirebaseUID: firebaseUser.UID,
		}

		err = a.GetDB().Debug().Model(&models.User{}).Create(newAppUser).Error
		if err != nil {
			helpers.ERROR(w, http.StatusInternalServerError, err)
			return
		}


		helpers.JSON(w, http.StatusOK, newAppUser)
	}
}