package main

import (
	"net/http"
	"zenful_shopping_backend/utils"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"env":     app.config.env,
	}
	if err := utils.WriteJSON(w, http.StatusOK, data); err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
	}
}
