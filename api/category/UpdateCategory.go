package CategoryAPI


import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/buyubaya/go-blog/interfaces"
	"github.com/buyubaya/go-blog/models"
	"github.com/buyubaya/go-blog/helpers"
)


type UpdateCategoryRequest struct {
	Name *string `json:"name"`
	Image *string `json:"image"`
	BackgroundColor *string `json:"backgroundColor"`
}


func UpdateCategory(a interfaces.IApp) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		cid, err := strconv.ParseUint(vars["categoryID"], 10, 64)
		if err != nil {
			helpers.ERROR(w, http.StatusInternalServerError, err)
			return
		}


		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			helpers.ERROR(w, http.StatusInternalServerError, err)
			return
		}


		updateCategoryRequest := &UpdateCategoryRequest{}
		err = json.Unmarshal(body, updateCategoryRequest)
		if err != nil {
			helpers.ERROR(w, http.StatusInternalServerError, err)
			return
		}


		// ADD
		newCategory := &models.Category{}
		change := false
		if (updateCategoryRequest.Name != nil) {
			newCategory.Name = *updateCategoryRequest.Name
			change = true
		}
		if (updateCategoryRequest.BackgroundColor != nil) {
			newCategory.BackgroundColor = updateCategoryRequest.BackgroundColor
			change = true
		}


		// VALIDATE
		if (change) {
			if err != nil {
				helpers.ERROR(w, http.StatusInternalServerError, errors.New("Nothing change"))
				return
			}
		}


		err = a.GetDB().Model(&models.Category{}).Where("id = ?", cid).Updates(newCategory).Error
		if err != nil {
			helpers.ERROR(w, http.StatusInternalServerError, err)
			return
		}


		helpers.JSON(w, http.StatusOK, newCategory)
	}
}