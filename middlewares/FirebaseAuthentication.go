package middlewares


import (
	"net/http"

	"github.com/gorilla/context"

	"github.com/buyubaya/go-blog/interfaces"
	"github.com/buyubaya/go-blog/helpers"
)


func FirebaseAuthentication(a interfaces.IApp) (func(http.Handler) http.Handler) {
	return func (next http.Handler) http.Handler {
	    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			idToken := r.Header.Get("Authorization")
			token, err := (*a.GetFirebaseApp()).VerifyIDToken(idToken)
			
			if err != nil {
				helpers.ERROR(w, http.StatusUnauthorized, err)
				return
			}

			context.Set(r, "Token", token)
			
	        next.ServeHTTP(w, r)
	    })
	}
}