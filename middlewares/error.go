package middlewares

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
	models "github.com/toolboxmodels"
)

// RespondWithError encapsulates a generic error response
// TODO: create a version which uses the toolbox/logger. Start using the new one. Stop using this.
func RespondWithError(w http.ResponseWriter, code int, message string, messageCode string) {
	data, err := json.Marshal(models.ErrorResponse{
		Error: models.ErrorResponseFormat{
			Code:        code,
			Message:     message,
			MessageCode: messageCode,
		},
	})
	if err != nil {
		logrus.Error("Error while Marshaling error response using json.Marshal(...)")
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte("Internal Server Error")); err != nil {
			logrus.WithError(err).Error("Failed to write HTTP response")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if _, err := w.Write(data); err != nil {
		logrus.WithError(err).Error("Failed to write HTTP response")
	}
}
