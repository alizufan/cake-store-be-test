package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rotisserie/eris"
	"github.com/zufzuf/cake-store/libs/util"
	"github.com/zufzuf/cake-store/repository"
	"github.com/zufzuf/cake-store/schema"
	"github.com/zufzuf/cake-store/service"
)

type CakeService interface {
	Find(ctx context.Context, id int) (*schema.Cake, error)
	FindAll(ctx context.Context, fil *service.FindAllRequest) ([]schema.Cake, error)
	Insert(ctx context.Context, req *service.CakeRequest) error
	Update(ctx context.Context, req *service.CakeRequest) error
	Delete(ctx context.Context, id int) error
}

type Cake struct {
	Service CakeService
}

func (h *Cake) FindCake(rw http.ResponseWriter, r *http.Request) {
	var (
		ctx   = r.Context()
		id, _ = strconv.Atoi(chi.URLParam(r, "id"))
	)

	res, err := h.Service.Find(ctx, id)
	if err != nil {
		if !eris.Is(err, repository.ErrRecordNotFound) {
			util.ErrHTTPResponse(ctx, rw, err)
			return
		}
	}

	msg := "not found"
	if res != nil {
		msg = "found"
	}

	util.HTTPResponse(rw, http.StatusOK, "search cake "+msg, res)
}

func (h *Cake) FindAllCake(rw http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
		q   = r.URL.Query()
		fil = service.FindAllRequest{
			Title:       q.Get("title"),
			Description: q.Get("description"),
		}
	)

	res, err := h.Service.FindAll(ctx, &fil)
	if err != nil {
		if !eris.Is(err, repository.ErrRecordNotFound) {
			util.ErrHTTPResponse(ctx, rw, err)
			return
		}
	}

	msg := "not found"
	if len(res) > 0 {
		msg = "found"
	}

	util.HTTPResponse(rw, http.StatusOK, "search cakes "+msg, res)
}

func (h *Cake) AddCake(rw http.ResponseWriter, r *http.Request) {
	var (
		ctx  = r.Context()
		body = service.CakeRequest{}
	)

	if ok := JSONDecodeValidation(rw, r.Body, &body); !ok {
		return
	}

	if err := h.Service.Insert(ctx, &body); err != nil {
		util.ErrHTTPResponse(ctx, rw, err)
		return
	}

	util.HTTPResponse(rw, http.StatusOK, "adding new cake", map[string]int{
		"id": body.ID,
	})
}

func (h *Cake) UpdateCake(rw http.ResponseWriter, r *http.Request) {
	var (
		ctx   = r.Context()
		id, _ = strconv.Atoi(chi.URLParam(r, "id"))
		body  = service.CakeRequest{}
	)

	if ok := JSONDecodeValidation(rw, r.Body, &body); !ok {
		return
	}

	body.ID = id
	if err := h.Service.Update(ctx, &body); err != nil {
		if !eris.Is(err, repository.ErrRecordNotFound) {
			util.ErrHTTPResponse(ctx, rw, err)
			return
		}
	}

	util.HTTPResponse(rw, http.StatusOK, "updating a cake", map[string]int{
		"id": id,
	})
}

func (h *Cake) DeleteCake(rw http.ResponseWriter, r *http.Request) {
	var (
		ctx   = r.Context()
		id, _ = strconv.Atoi(chi.URLParam(r, "id"))
	)

	if err := h.Service.Delete(ctx, id); err != nil {
		if !eris.Is(err, repository.ErrRecordNotFound) {
			util.ErrHTTPResponse(ctx, rw, err)
			return
		}
	}

	util.HTTPResponse(rw, http.StatusOK, "deleting a cake", map[string]int{
		"id": id,
	})
}

func JSONDecodeValidation(rw http.ResponseWriter, reqBody io.ReadCloser, data any) bool {
	if err := json.NewDecoder(reqBody).Decode(data); err != nil {
		var (
			code = http.StatusBadRequest
			msg  = "parse request body, an error occured"
		)
		util.ErrorHTTPResponse(rw, code, msg, nil)
		return false
	}

	if errs := util.Validation(data); len(errs) > 0 {
		var (
			code = http.StatusUnprocessableEntity
			msg  = "unprocessable request body, an error occured"
		)
		util.ErrorHTTPResponse(rw, code, msg, map[string]any{
			"validation": errs,
		})
		return false
	}

	return true
}
