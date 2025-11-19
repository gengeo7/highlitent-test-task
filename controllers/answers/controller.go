package answers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gengeo7/highlitent/apierror"
	"github.com/gengeo7/highlitent/middleware"
	answersService "github.com/gengeo7/highlitent/services/answers"
	answersStorage "github.com/gengeo7/highlitent/storage/answers"
	"github.com/gengeo7/highlitent/types/answers"
	"github.com/gengeo7/highlitent/utils"
)

const BaseRoute string = "/answers"

type AnswersController struct {
	Storage answersStorage.Storage
}

func NewAnswersController(storage answersStorage.Storage) *AnswersController {
	return &AnswersController{Storage: storage}
}

func (ac *AnswersController) RegisterController(mux *http.ServeMux) {
	mux.Handle(
		"POST /questions/{id}/answers",
		middleware.Chain(
			http.HandlerFunc(ac.postAnswer),
			middleware.Timeout(5*time.Second),
			middleware.ValidateJson[answers.AnswerDto](),
		),
	)

	mux.Handle(
		fmt.Sprintf("GET %s/{id}", BaseRoute),
		middleware.Chain(
			http.HandlerFunc(ac.getAnswer),
			middleware.Timeout(5*time.Second),
		),
	)

	mux.Handle(
		fmt.Sprintf("DELETE %s/{id}", BaseRoute),
		middleware.Chain(
			http.HandlerFunc(ac.deleteAnswer),
			middleware.Timeout(5*time.Second),
		),
	)
}

func (ac *AnswersController) getAnswer(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendResponse(nil, apierror.NewApiError(http.StatusBadRequest, "некоректный id", nil), w, r)
		return
	}
	answer, err := answersService.GetAnswer(r.Context(), ac.Storage, id)
	utils.SendResponse(&utils.Response{Data: answer, Status: http.StatusOK}, err, w, r)
}

func (ac *AnswersController) deleteAnswer(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendResponse(nil, apierror.NewApiError(http.StatusBadRequest, "некоректный id", nil), w, r)
		return
	}
	err = answersService.DeleteAnswer(r.Context(), ac.Storage, id)
	utils.SendResponse(nil, err, w, r)
}

func (ac *AnswersController) postAnswer(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendResponse(nil, apierror.NewApiError(http.StatusBadRequest, "некоректный id", nil), w, r)
		return
	}
	dto := middleware.DtoFromContext[answers.AnswerDto](r.Context())
	answer, err := answersService.CreateAnswer(r.Context(), ac.Storage, dto, id)
	utils.SendResponse(&utils.Response{Data: answer, Status: http.StatusCreated}, err, w, r)
}
