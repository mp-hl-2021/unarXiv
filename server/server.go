package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/mp-hl-2021/unarXiv/api"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

type UnarXivServer struct {
	Core api.UnarXiv
}

func NewServer(unarXiv api.UnarXiv) *UnarXivServer {
	return &UnarXivServer{
		Core: unarXiv,
	}
}

func (server *UnarXivServer) Router() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/register", server.PostRegister).Methods(http.MethodPost)
	router.HandleFunc("/login", server.PostLogin).Methods(http.MethodPost)

	// offset is optional, should be passed as "?offset=smth"
	router.Path("/search").Queries("query", "{query}").
		HandlerFunc(server.GetSearch).Methods(http.MethodGet)

	router.HandleFunc("/articles/{articleId}", server.GetArticle).Methods(http.MethodGet)

	router.HandleFunc("/history/searches", server.GetSearchHistory).Methods(http.MethodGet)
	router.HandleFunc("/history/articles", server.GetArticlesHistory).Methods(http.MethodGet)

	router.HandleFunc("/updates/searches", server.GetSearchQueriesUpdates).Methods(http.MethodGet)
	router.HandleFunc("/updates/articles", server.GetArticlesUpdates).Methods(http.MethodGet)

	router.Path("/subscriptions/articles/{articleId}").
		HandlerFunc(server.GetArticleSubscriptionStatus).Methods(http.MethodGet)
	router.Path("/subscriptions/articles/{articleId}").Queries("subscribe", "{subscribe:(true|false)}").
		HandlerFunc(server.PostArticleSubscriptionStatus).Methods(http.MethodGet)

	router.Path("/subscriptions/searches").Queries("query", "{query}").
		HandlerFunc(server.GetSearchQuerySubscriptionStatus).Methods(http.MethodGet)
	router.Path("/subscriptions/searches").Queries("query", "{query}", "subscribe", "{subscribe:(true|false)}").
		HandlerFunc(server.PostSearchQuerySubscriptionStatus).Methods(http.MethodPost)

	return router
}

func RespondWithJSON(w http.ResponseWriter, object interface{}, status int) error {
	w.Header().Set("Content-Type", "application/json")

	if object != nil { // non-empty body
		if err := json.NewEncoder(w).Encode(object); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}
	}

	w.WriteHeader(status)
	return nil
}

// postRegister handles request for a new account creation.
func (server *UnarXivServer) PostRegister(w http.ResponseWriter, r *http.Request) {
	var registerRequest api.AuthenticationRequest
	if err := json.NewDecoder(r.Body).Decode(&registerRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	authData, err := server.Core.Register(&registerRequest)
	if err != nil { // todo: map domain errors to http error codes
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error happened in Core.Register: %v", err)
		return
	}

	if err := RespondWithJSON(w, authData, http.StatusCreated); err != nil {
		log.Printf("Error happened while responding to PostRegister: %v", err)
	}
}

// PostLogin handles login request for existing user.
func (server *UnarXivServer) PostLogin(w http.ResponseWriter, r *http.Request) {
	var authRequest api.AuthenticationRequest
	if err := json.NewDecoder(r.Body).Decode(&authRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	authData, err := server.Core.Login(&authRequest)
	if err != nil { // todo: map domain errors to http error codes
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error happened in Core.Login: %v", err)
		return
	}

	if err := RespondWithJSON(w, authData, http.StatusOK); err != nil {
		log.Printf("Error happened while responding to PostLogin: %v", err)
	}
}

func (server *UnarXivServer) GetSearch(w http.ResponseWriter, r *http.Request) {
	var searchQueryRequest = api.SearchQueryRequest{}
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error happened while parsing form params: %v", err)
		return
	}
	searchQueryRequest.Query = r.Form.Get("query")
	// todo: extract auth ?
	if err := json.NewDecoder(r.Body).Decode(&searchQueryRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := server.Core.Search(&searchQueryRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error happened in Core.Search: %v", err)
		return
	}

	if err := RespondWithJSON(w, response, http.StatusCreated); err != nil {
		log.Printf("Error happened while responding to PostLogin: %v", err)
	}
}

func (server *UnarXivServer) GetArticle(w http.ResponseWriter, r *http.Request) {
	var articleRequest api.AccessArticleRequest
	if err := json.NewDecoder(r.Body).Decode(&articleRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := server.Core.AccessArtice(articleRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	WriteJson(&w, []byte(response))
}

func (server *UnarXivServer) GetArticlesHistory(w http.ResponseWriter, r *http.Request) {
	authData := api.AuthenticationData("jwt") // TODO extract data from headers
	response, err := server.Core.GetArticlesHistory(authData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	WriteJson(&w, []byte(response))
}

func (server *UnarXivServer) GetSearchHistory(w http.ResponseWriter, r *http.Request) {
	authData := api.AuthenticationData("jwt") // TODO extract data from headers
	response, err := server.Core.GetSearchHistory(authData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	WriteJson(&w, []byte(response))
}

func (server *UnarXivServer) GetSearchQueriesUpdates(w http.ResponseWriter, r *http.Request) {
	authData := api.AuthenticationData("jwt") // TODO extract data from headers
	response, err := server.Core.GetSearchQueriesUpdates(authData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	WriteJson(&w, []byte(response))
}

func (server *UnarXivServer) GetArticlesUpdates(w http.ResponseWriter, r *http.Request) {
	authData = api.AuthenticationData("jwt") // TODO extract data from headers
	response, err := server.Core.GetArticlesUpdates(authData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	WriteJson(&w, []byte(response))
}

func (server *UnarXivServer) GetArticleSubscriptionStatus(w http.ResponseWriter, r *http.Request) {
	var artRequest api.SetArticleSubscriptionStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&artRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := server.Core.GetArticleSubscriptionStatus(artRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	WriteJson(&w, []byte(response))
}

func (server *UnarXivServer) PostArticleSubscriptionStatus(w http.ResponseWriter, r *http.Request) {
	var artRequest api.SetArticleSubscriptionStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&artRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := server.Core.AccessArtice(artRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (server *UnarXivServer) GetSearchQuerySubscriptionStatus(w http.ResponseWriter, r *http.Request) {
	var s api.SetSearchQuerySubscriptionStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := server.Core.AccessArtice(s)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	WriteJson(&w, []byte(response))
}

func (server *UnarXivServer) PostSearchQuerySubscriptionStatus(w http.ResponseWriter, r *http.Request) {
	var s api.AccessArticleRequest
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := server.Core.SetSearchQuerySubscriptionStatus(S)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
