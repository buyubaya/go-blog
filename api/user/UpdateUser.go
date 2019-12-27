package UserAPI


import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"time"
	"errors"

	"github.com/gorilla/mux"

	"github.com/buyubaya/go-blog/models"
	"github.com/buyubaya/go-blog/interfaces"
	"github.com/buyubaya/go-blog/helpers"
	"firebase.google.com/go/auth"
)


type UpdateUserRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
	DisplayName string `json:"displayName"`
	PhotoURL string `json:"photoURL"`
}


func UpdateUser(a interfaces.IApp) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		uid := vars["UID"]

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			helpers.ERROR(w, http.StatusInternalServerError, err)
			return
		}


		updateUserRequest := &UpdateUserRequest{}
		err = json.Unmarshal(body, updateUserRequest)
		if err != nil {
			helpers.ERROR(w, http.StatusInternalServerError, err)
			return
		}


		// UPDATE FIREBASE USER
		newFirebaseUser := (&auth.UserToUpdate{})

		if (updateUserRequest.Email != "") {
			newFirebaseUser = newFirebaseUser.Email(updateUserRequest.Email)
		}
		if (updateUserRequest.Password != "" && updateUserRequest.ConfirmPassword != "") {
			if (updateUserRequest.Password == updateUserRequest.ConfirmPassword) {
				newFirebaseUser = newFirebaseUser.Password(updateUserRequest.Password)
			} else {
				helpers.ERROR(w, http.StatusInternalServerError, errors.New("Password and Confirm Pasword does not match"))
				return
			}
		}
		if (updateUserRequest.DisplayName != "") {
			newFirebaseUser = newFirebaseUser.DisplayName(updateUserRequest.DisplayName)
		}
		if (updateUserRequest.PhotoURL != "") {
			newFirebaseUser = newFirebaseUser.PhotoURL(updateUserRequest.PhotoURL)
		}


		firebaseUser, err := (*a.GetFirebaseApp()).UpdateUser(uid, newFirebaseUser)
		if err != nil {
			helpers.ERROR(w, http.StatusInternalServerError, err)
			return
		}


		// UPDATE APP USER
		newAppUser := &models.User{
			UpdatedAt: time.Now(),
		}

		err = a.GetDB().Debug().Model(&models.User{}).Where("firebase_uid = ?", firebaseUser.UID).Updates(newAppUser).Error
		if err != nil {
			helpers.ERROR(w, http.StatusInternalServerError, err)
			return
		}


		helpers.JSON(w, http.StatusOK, firebaseUser)
	}
}