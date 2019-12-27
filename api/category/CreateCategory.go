package CategoryAPI


import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"errors"

	"github.com/buyubaya/go-blog/interfaces"
	"github.com/buyubaya/go-blog/models"
	"github.com/buyubaya/go-blog/helpers"
)


type CreateCategoryRequest struct {
	Name string `json:"name"`
	Image *string `json:"image"`
	BackgroundColor *string `json:"backgroundColor"`
}


func CreateCategory(a interfaces.IApp) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			helpers.ERROR(w, http.StatusInternalServerError, err)
			return
		}


		createCategoryRequest := &CreateCategoryRequest{}
		err = json.Unmarshal(body, createCategoryRequest)
		if err != nil {
			helpers.ERROR(w, http.StatusInternalServerError, err)
			return
		}


		// VALIDATE
		if (createCategoryRequest.Name == "") {
			if err != nil {
				helpers.ERROR(w, http.StatusInternalServerError, errors.New("Category Name is required"))
				return
			}
		}


		// ADD
		newCategory := &models.Category{
			Name: createCategoryRequest.Name,
		}
		if (createCategoryRequest.BackgroundColor != nil) {
			newCategory.BackgroundColor = createCategoryRequest.BackgroundColor
		}

		err = a.GetDB().Model(&models.Category{}).Create(newCategory).Error
		if err != nil {
			helpers.ERROR(w, http.StatusInternalServerError, err)
			return
		}


		helpers.JSON(w, http.StatusOK, newCategory)
	}
}