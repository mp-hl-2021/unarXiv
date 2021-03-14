package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/mp-hl-2021/unarXiv/api"
	"net"
	"net/http"
)

type Api struct {
	 UseCases api.UseCasesInterface
}

func NewApi(a api.UseCasesInterface) *Api {
	return &Api{
		UseCases: a,
	}
}


func (a *Api) Router() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/register", a.PostRegister).Methods(http.MethodPost)
	router.HandleFunc("/login", a.PostLogin).Methods(http.MethodPost)

	router.HandleFunc("/search", a.GetSearch).Methods(http.MethodGet)

    router.HandleFunc("/articles/{articleId}", a.GetArticle).Methods(http.MethodGet)

	router.HandleFunc("/history/search", a.GetSeachHistory).Methods(http.MethodGet)
	router.HandleFunc("/history/articles", a.GetArticlesHistory).Methods(http.MethodGet)

    router.HandleFunc("/updates/search", a.GetSearchQueriesUpdates).Methods(http.MethodGet)
    router.HandleFunc("/updates/articles", a.GetArticlesUpdates).Methods(http.MethodGet)

	router.HandleFunc("/subscriptions/articles/{articleId}",
		a.GetArticleSubscriptionStatus).Methods(http.MethodGet)
	router.HandleFunc("/subscriptions/articles/{articleId}",
		a.PostArticleSubscriptionStatus).Methods(http.MethodPost)

	router.HandleFunc("/subscriptions/queries",
		a.GetSearchQuerySubscriptionStatus).Methods(http.MethodGet)
	router.HandleFunc("/subscriptions/queries",
		a.PostSearchQuerySubscriptionStatus).Methods(http.MethodPost)

	return router
}

func WriteJson(w *http.ResponseWriter, content []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(content)
	w.WriteHeader(http.StatusOK)
}

// postRegister handles request for a new account creation.
func (a *Api) PostRegister(w http.ResponseWriter, r *http.Request) {
    var authRequest api.AuthenticationRequest
    if err := json.NewDecoder(r.Body).Decode(&authRequest); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    token, err := a.Register(authRequest)
    if err != nil { // todo: map domain errors to http error codes
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
}

// PostLogin handles login request for existing user.
func (a *Api) PostLogin(w http.ResponseWriter, r *http.Request) {
    var m api.AuthenticationRequest
    if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    token, err := a.UseCases.Login(authRequest)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

	w.Header().Set("Content-Type", "application/jwt")
	w.Write([]byte(token))
	w.WriteHeader(http.StatusOK)
}


func (a *Api) GetSearch(w http.ResponseWriter, r *http.Request) {
	var s api.SearchQueryResponse
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := a.UseCases.Search(s)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	WriteJson(&w, []byte(response))
}

func (a *Api) GetArticle(w http.ResponseWriter, r *http.Request) {
	var artRequest api.AccessArticleRequest
	if err := json.NewDecoder(r.Body).Decode(&artRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := a.UseCases.AccessArtice(artRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	WriteJson(&w, []byte(response))
}

func (a *Api) GetArticlesHistory(w http.ResponseWriter, r *http.Request) {
	authData := api.AuthenticationData("jwt") // TODO extract data from headers
	response, err := a.UseCases.GetArticlesHistory(authData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	WriteJson(&w, []byte(response))
}

func (a *Api) GetSearchHistory(w http.ResponseWriter, r *http.Request) {
	authData := api.AuthenticationData("jwt") // TODO extract data from headers
	response, err := a.UseCases.GetSearchHistory(authData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	WriteJson(&w, []byte(response))
}


func (a *Api) GetSearchQueriesUpdates(w http.ResponseWriter, r *http.Request) {
	authData := api.AuthenticationData("jwt") // TODO extract data from headers
	response, err := a.UseCases.GetSearchQueriesUpdates(authData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	WriteJson(&w, []byte(response))
}

func (a *Api) GetArticlesUpdates(w http.ResponseWriter, r *http.Request) {
	authData = api.AuthenticationData("jwt") // TODO extract data from headers
	response, err := a.UseCases.GetArticlesUpdates(authData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	WriteJson(&w, []byte(response))
}

func (a *Api) GetArticleSubscriptionStatus(w http.ResponseWriter, r *http.Request) {
	var artRequest api.SetArticleSubscriptionStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&artRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := a.UseCases.GetArticleSubscriptionStatus(artRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	WriteJson(&w, []byte(response))
}

func (a *Api) PostArticleSubscriptionStatus(w http.ResponseWriter, r *http.Request) {
	var artRequest api.SetArticleSubscriptionStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&artRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := a.UseCases.AccessArtice(artRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *Api) GetSearchQuerySubscriptionStatus(w http.ResponseWriter, r *http.Request) {
	var s api.SetSearchQuerySubscriptionStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := a.UseCases.AccessArtice(s)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	WriteJson(&w, []byte(response))
}

func (a *Api) PostSearchQuerySubscriptionStatus(w http.ResponseWriter, r *http.Request) {
	var s api.AccessArticleRequest
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := a.UseCases.SetSearchQuerySubscriptionStatus(S)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}