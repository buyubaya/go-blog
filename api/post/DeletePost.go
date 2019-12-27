package PostAPI


import (
	"net/http"
	"strconv"
	"time"
	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/gorilla/context"
	"firebase.google.com/go/auth"

	"github.com/buyubaya/go-blog/models"
	"github.com/buyubaya/go-blog/interfaces"
	"github.com/buyubaya/go-blog/helpers"
)


func DeletePost(a interfaces.IApp) (http.HandlerFunc) {
	return func (w http.ResponseWriter, r *http.Request) {
		token := context.Get(r, "Token")
		uid := token.(*auth.Token).UID

		vars := mux.Vars(r)
		pid, err := strconv.ParseUint(vars["postID"], 10, 64)
		if err != nil {
			helpers.ERROR(w, http.StatusInternalServerError, err)
			return
		}


		// GET OLD POST
		oldPost := &models.Post{}
		err = a.GetDB().Debug().Model(&models.Post{}).Where("id = ?", pid).Find(&oldPost).Error
		if err != nil {
			helpers.ERROR(w, http.StatusInternalServerError, err)
			return
		}



		// DELETE POST
		err = a.GetDB().Debug().Model(&models.Post{}).Where("id = ?", pid).Delete(&models.Post{}).Error
		if err != nil {
			helpers.ERROR(w, http.StatusInternalServerError, err)
			return
		}


		// DELETE CATEGORY
		err = a.GetDB().Debug().Model(&models.PostFieldCategory{}).
			Where("post_id = ?", pid).
			Delete(&models.PostFieldCategory{}).Error
		if err != nil {
			helpers.ERROR(w, http.StatusInternalServerError, err)
			return
		}


		// CHANGELOG
		changelog, err := json.Marshal(oldPost)
		
		tx := a.GetDB().Begin()
		// ACTIVITY
		activity := &models.Activity{
			Action: "DELETE_POST",
			Changelog: string(changelog),
			CreatedAt: time.Now(),
			CreatedBy: uid,
		}
		
		err = tx.Model(&models.Activity{}).
			Create(activity).Error
		if err != nil {
			tx.Rollback()
			helpers.ERROR(w, http.StatusInternalServerError, err)
			return
		}

		// OBJECT_ACTIVITY
		objectActivity := &models.ObjectFieldActivity{
			ObjectID: "post::" + strconv.FormatUint(pid, 10),
			ActivityID: activity.ID,
		}

		err = tx.Model(&models.ObjectFieldActivity{}).
			Create(objectActivity).Error
		if err != nil {
			tx.Rollback()
			helpers.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		tx.Commit()


		helpers.JSON(w, http.StatusOK, pid)
	}
}
