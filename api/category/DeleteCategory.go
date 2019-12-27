package CategoryAPI


import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/buyubaya/go-blog/models"
	"github.com/buyubaya/go-blog/interfaces"
	"github.com/buyubaya/go-blog/helpers"
)


func DeleteCategory(a interfaces.IApp) (http.HandlerFunc) {
	return func (w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		cid, err := strconv.ParseUint(vars["categoryID"], 10, 64)
		if err != nil {
			helpers.ERROR(w, http.StatusInternalServerError, err)
			return
		}


		tx := a.GetDB().Begin()

		// DELETE POST FIELD CATEGORY
		err = a.GetDB().Debug().Model(&models.PostFieldCategory{}).
			Where("category_id = ?", cid).
			Delete(&models.PostFieldCategory{}).Error
		if err != nil {
			helpers.ERROR(w, http.StatusInternalServerError, err)
			tx.Rollback()
			return
		}

		// DELETE CATEGORY
		err = a.GetDB().Debug().Model(&models.Category{}).Where("id = ?", cid).Delete(&models.Category{}).Error
		if err != nil {
			helpers.ERROR(w, http.StatusInternalServerError, err)
			tx.Rollback()
			return
		}

		tx.Commit()


		helpers.JSON(w, http.StatusOK, cid)
	}
}
