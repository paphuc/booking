package middleware

import (
	"net/http"

	"booking/configs"
	"booking/internal/pkg/auth"
	"booking/internal/pkg/glog"
	"booking/internal/pkg/respond"
)

func Auth(h http.HandlerFunc, em *configs.ErrorMessage) http.HandlerFunc {
	logger := glog.New().WithField("package", "middleware")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenpath := auth.ExtractToken(r)
		if tokenpath == "" {
			logger.Infof("The request does not contain token")
			respond.JSON(w, http.StatusUnauthorized, &em.InvalidValue.FailedAuthentication)
			return
		}
		_, err := auth.IsAuthorized(tokenpath)

		if err != nil {
			logger.Errorf("Not authorized, error: ", err)
			respond.JSON(w, http.StatusUnauthorized, &em.InvalidValue.FailedAuthentication)
			return
		}

		h.ServeHTTP(w, r)
	})
}
