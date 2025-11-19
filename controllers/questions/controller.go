package questions

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gengeo7/highlitent/apierror"
	"github.com/gengeo7/highlitent/middleware"
	questionsService "github.com/gengeo7/highlitent/services/questions"
	questionsStorage "github.com/gengeo7/highlitent/storage/questions"
	"github.com/gengeo7/highlitent/types/questions"
	"github.com/gengeo7/highlitent/utils"
)

const BaseRoute string = "/questions"

type QuestionsController struct {
	Storage questionsStorage.Storage
}

func NewQuestionsController(storage questionsStorage.Storage) *QuestionsController {
	return &QuestionsController{Storage: storage}
}

func (qc *QuestionsController) RegisterController(mux *http.ServeMux) {
	mux.Handle(
		fmt.Sprintf("GET %s", BaseRoute),
		middleware.Chain(
			http.HandlerFunc(qc.getAllQuestions),
			middleware.Timeout(5*time.Second),
		),
	)

	mux.Handle(
		fmt.Sprintf("POST %s", BaseRoute),
		middleware.Chain(
			http.HandlerFunc(qc.newQuestion),
			middleware.Timeout(5*time.Second),
			middleware.ValidateJson[questions.QuestionDto](),
		),
	)

	mux.Handle(
		fmt.Sprintf("GET %s/{id}", BaseRoute),
		middleware.Chain(
			http.HandlerFunc(qc.getQuestionWithAnswers),
			middleware.Timeout(5*time.Second),
		),
	)

	mux.Handle(
		fmt.Sprintf("DELETE %s/{id}", BaseRoute),
		middleware.Chain(
			http.HandlerFunc(qc.deleteQuestion),
			middleware.Timeout(5*time.Second),
		),
	)
}

func (qc *QuestionsController) getAllQuestions(w http.ResponseWriter, r *http.Request) {
	questions, err := questionsService.GetAllQuestions(r.Context(), qc.Storage)
	utils.SendResponse(&utils.Response{Data: questions, Status: http.StatusOK}, err, w, r)
}

func (qc *QuestionsController) newQuestion(w http.ResponseWriter, r *http.Request) {
	dto := middleware.DtoFromContext[questions.QuestionDto](r.Context())
	question, err := questionsService.CreateQuestion(r.Context(), qc.Storage, dto)
	utils.SendResponse(&utils.Response{Data: question, Status: http.StatusCreated}, err, w, r)
}

func (qc *QuestionsController) getQuestionWithAnswers(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendResponse(nil, apierror.NewApiError(http.StatusBadRequest, "некоректный id", nil), w, r)
		return
	}
	questionWithAnswers, err := questionsService.GetQuestionWithAnswers(r.Context(), qc.Storage, id)
	utils.SendResponse(&utils.Response{Data: questionWithAnswers, Status: http.StatusOK}, err, w, r)
}

func (qc *QuestionsController) deleteQuestion(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendResponse(nil, apierror.NewApiError(http.StatusBadRequest, "некоректный id", nil), w, r)
		return
	}
	err = questionsService.DeleteQuestion(r.Context(), qc.Storage, id)
	utils.SendResponse(nil, err, w, r)
}
