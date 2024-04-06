package main

import (
	"errors"
	"github.com/godra-y/go-project/pkg/api/model"
	"github.com/godra-y/go-project/pkg/api/validator"
	"log"
	"net/http"
)

func (app *application) createCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string `json:"name"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		log.Println(err)
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}

	category := &model.Category{
		Name: input.Name,
	}

	err = app.models.Category.Insert(category)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"category": category}, nil)
}

func (app *application) getCategoriesList(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string
		model.Filters
	}
	v := validator.New()
	qs := r.URL.Query()

	input.Name = app.readStrings(qs, "name", "")

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readStrings(qs, "sort", "id")

	input.Filters.SortSafeList = []string{
		"id", "name",
		"-id", "-name",
	}

	if model.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	categories, metadata, err := app.models.Category.GetAll(input.Name, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"categories": categories, "metadata": metadata}, nil)
}

func (app *application) getCategoryHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	category, err := app.models.Category.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"category": category}, nil)
}

func (app *application) updateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	category, err := app.models.Category.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):

			app.notFoundResponse(w, r)
		default:

			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Name *string `json:"name"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Name != nil {
		category.Name = *input.Name
	}

	v := validator.New()

	if model.ValidateCategory(v, category); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Category.Update(category)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"category": category}, nil)
}

func (app *application) deleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Category.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "success"}, nil)
}
