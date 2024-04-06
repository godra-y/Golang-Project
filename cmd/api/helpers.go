package main

import (
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (app *application) readIDParam(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	param := vars["id"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}
