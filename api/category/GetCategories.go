package CategoryAPI


import (
	"net/http"

	"github.com/buyubaya/go-blog/interfaces"
	"github.com/buyubaya/go-blog/models"
	"github.com/buyubaya/go-blog/helpers"
)


func GetCategories(a interfaces.IApp) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {

		categories := []models.Category{}


		err := a.GetDB().Model(&models.Category{}).Find(&categories).Error
		if err != nil {
			helpers.ERROR(w, http.StatusInternalServerError, err)
			return
		}


		helpers.JSON(w, http.StatusOK, categories)
	}
}