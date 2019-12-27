package UserAPI


import (
	"net/http"
	"io/ioutil"
	"encoding/json"

	"github.com/buyubaya/go-blog/interfaces"
	"github.com/buyubaya/go-blog/helpers"
)


type GetUsersRequest struct {
	UserUIDs []string `json:"userUIDs"`
}


func GetUsers(a interfaces.IApp) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			helpers.ERROR(w, http.StatusInternalServerError, err)
			return
		}


		getUsersRequest := &GetUsersRequest{}
		err = json.Unmarshal(body, getUsersRequest)
		if err != nil {
			helpers.ERROR(w, http.StatusInternalServerError, err)
			return
		}


		// GET USERS
		userInfos, err := (*a.GetFirebaseApp()).GetUserByUIDs(getUsersRequest.UserUIDs)
		if err != nil {
			helpers.ERROR(w, http.StatusInternalServerError, err)
			return
		}


		helpers.JSON(w, http.StatusOK, userInfos)
	}
}