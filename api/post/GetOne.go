package PostAPI


import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/buyubaya/go-blog/models"
	"github.com/buyubaya/go-blog/interfaces"
	"github.com/buyubaya/go-blog/helpers"
)


// GET ONE POST
func GetOne(a interfaces.IApp) (http.HandlerFunc) {
	return func (w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		pid, err := strconv.ParseUint(vars["postID"], 10, 64)
		if err != nil {
			helpers.ERROR(w, http.StatusInternalServerError, err)
			return
		}

		model := &models.Post{}

		post, err := model.FindPostByID(a.GetDB(), pid)
		if err != nil {
			helpers.ERROR(w, http.StatusInternalServerError, err)
			return
		}

		helpers.JSON(w, http.StatusOK, post)
	}
}
