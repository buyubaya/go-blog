package ActivityAPI


import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/buyubaya/go-blog/models"
	"github.com/buyubaya/go-blog/interfaces"
	"github.com/buyubaya/go-blog/helpers"
)


func DeleteActivity(a interfaces.IApp) (http.HandlerFunc) {
	return func (w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		aid, err := strconv.ParseUint(vars["activityID"], 10, 64)
		if err != nil {
			helpers.ERROR(w, http.StatusInternalServerError, err)
			return
		}


		// DELETE ACTIVITY
		tx := a.GetDB().Begin()

		err = tx.Debug().Model(&models.Activity{}).
			Where("id = ?", aid).
			Delete(&models.Activity{}).Error
		if err != nil {
			tx.Rollback()
			helpers.ERROR(w, http.StatusInternalServerError, err)
			return
		}

		err = tx.Debug().Model(&models.ObjectFieldActivity{}).
			Where("activity_id = ?", aid).
			Delete(&models.ObjectFieldActivity{}).Error
		if err != nil {
			tx.Rollback()
			helpers.ERROR(w, http.StatusInternalServerError, err)
			return
		}

		tx.Commit()


		helpers.JSON(w, http.StatusOK, aid)
	}
}
