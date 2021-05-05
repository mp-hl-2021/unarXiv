package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/mp-hl-2021/unarXiv/internal/domain/model"
	"github.com/mp-hl-2021/unarXiv/internal/usecases"
	"log"
	"net/http"
	"strconv"
	"strings"
	//"github.com/dgrijalva/jwt-go"
)

type HttpApi struct {
	usecases usecases.Interface
}

func New(usecases usecases.Interface) *HttpApi {
	return &HttpApi{
		usecases: usecases,
	}
}

type contextKey string

const (
	contextKeyUserId = contextKey("userId")
	bearer           = "Bearer"
)

func (a *HttpApi) Router() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/register", a.postRegister).Methods(http.MethodPost)
	router.HandleFunc("/login", a.postLogin).Methods(http.MethodPost)

	// offset is optional, should be passed as "?offset=smth"
	router.Path("/search/{query}").HandlerFunc(a.extractAuth(a.getSearch)).Methods(http.MethodGet)

	router.HandleFunc("/articles/{articleId}", a.extractAuth(a.getArticle)).Methods(http.MethodGet)

	router.HandleFunc("/history/searches", a.extractAuth(a.getSearchHistory)).Methods(http.MethodGet)
	router.HandleFunc("/history/articles", a.extractAuth(a.getArticlesHistory)).Methods(http.MethodGet)

	router.HandleFunc("/updates/searches", a.extractAuth(a.getSearchQueriesUpdates)).Methods(http.MethodGet)
	router.HandleFunc("/updates/articles", a.extractAuth(a.getArticlesUpdates)).Methods(http.MethodGet)

	router.Path("/subscriptions/articles/{articleId}").
		HandlerFunc(a.extractAuth(a.getArticleSubscriptionStatus)).Methods(http.MethodGet)
	router.Path("/subscriptions/articles/{articleId}").
		HandlerFunc(a.extractAuth(a.postArticleSubscriptionStatus)).Methods(http.MethodPost)
	router.Path("/subscriptions/articles/{articleId}").
		HandlerFunc(a.extractAuth(a.deleteArticleSubscriptionStatus)).Methods(http.MethodDelete)

	router.Path("/subscriptions/searches/{query}").
		HandlerFunc(a.extractAuth(a.getSearchQuerySubscriptionStatus)).Methods(http.MethodGet)
	router.Path("/subscriptions/searches/{query}").
		HandlerFunc(a.extractAuth(a.postSearchQuerySubscriptionStatus)).Methods(http.MethodPost)
	router.Path("/subscriptions/searches/{query}").
		HandlerFunc(a.extractAuth(a.deleteSearchQuerySubscriptionStatus)).Methods(http.MethodDelete)

	return router
}

func respondWithJSON(w http.ResponseWriter, object interface{}, status int) error {
	if status != http.StatusOK {
		w.WriteHeader(status)
	}

	w.Header().Set("Content-Type", "application/json")

	if object != nil { // non-empty body
		if err := json.NewEncoder(w).Encode(object); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}
	}
	return nil
}

func userIdFromRequest(r *http.Request) (model.UserId, bool) {
	userId, ok := r.Context().Value(contextKeyUserId).(model.UserId)
	if userId == "" {
		ok = false
	}
	return userId, ok
}

func userIdPtrFromRequest(r *http.Request) *model.UserId {
	userId, ok := userIdFromRequest(r)
	if !ok {
		return nil
	}
	return &userId
}

// postRegister handles request for a new account creation.
func (a *HttpApi) postRegister(w http.ResponseWriter, r *http.Request) {
	var registerRequest AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&registerRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	authToken, err := a.usecases.Register(usecases.AuthRequest{
		Login:    registerRequest.Login,
		Password: registerRequest.Password,
	})
	if err != nil { // todo: map domain errors to http error codes
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error happened in usecases.Register: %v", err)
		return
	}

	if err := respondWithJSON(w, renderAuthToken(authToken), http.StatusCreated); err != nil {
		log.Printf("Error happened while responding to PostRegister: %v", err)
	}
}

// PostLogin handles login request for existing user.
func (a *HttpApi) postLogin(w http.ResponseWriter, r *http.Request) {
	var authRequest AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&authRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	authToken, err := a.usecases.Login(usecases.AuthRequest{
		Login:    authRequest.Login,
		Password: authRequest.Password,
	})
	if err != nil { // todo: map domain errors to http error codes
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error happened in usecases.Login: %v", err)
		return
	}

	if err := respondWithJSON(w, renderAuthToken(authToken), http.StatusOK); err != nil {
		log.Printf("Error happened while responding to PostLogin: %v", err)
	}
}

func (a *HttpApi) getSearch(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error happened while parsing form params: %v", err)
		return
	}
	var searchQueryRequest = model.SearchQuery{}
	searchQueryRequest.Query = mux.Vars(r)["query"]
	if strOffset := r.Form.Get("offset"); len(strOffset) != 0 {
		offset, err := strconv.Atoi(strOffset)
		if err != nil || offset < 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		searchQueryRequest.Offset = uint32(offset)
	}

	result, err := a.usecases.Search(searchQueryRequest, userIdPtrFromRequest(r))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error happened in usecases.Search: %v", err)
		return
	}

	if err := respondWithJSON(w, renderSearchResults(result), http.StatusOK); err != nil {
		log.Printf("Error happened while responding to GetSearch: %v", err)
	}
}

func (a *HttpApi) getArticle(w http.ResponseWriter, r *http.Request) {
	var articleId model.ArticleId
	articleId = model.ArticleId(mux.Vars(r)["articleId"])

	result, err := a.usecases.AccessArticle(articleId, userIdPtrFromRequest(r))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error happened in usecases.AccessArticle: %v", err)
		return
	}

	if err := respondWithJSON(w, renderArticle(result), http.StatusOK); err != nil {
		log.Printf("Error happened while responding to GetArticle: %v", err)
	}
}

func (a *HttpApi) getArticlesHistory(w http.ResponseWriter, r *http.Request) {
	result, err := a.usecases.GetArticleHistory("0") // TODO extract auth from headers
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error happened in usecases.GetArticlesHistory: %v", err)
		return
	}

	if err := respondWithJSON(w, renderUserArticleHistory(result), http.StatusOK); err != nil {
		log.Printf("Error happened while responding to GetArticlesHistory: %v", err)
	}
}

func (a *HttpApi) getSearchHistory(w http.ResponseWriter, r *http.Request) {
	userId, ok := userIdFromRequest(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	result, err := a.usecases.GetSearchHistory(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error happened in usecases.GetSearchHistory: %v", err)
		return
	}

	if err := respondWithJSON(w, renderUserSearchHistory(result), http.StatusOK); err != nil {
		log.Printf("Error happened while responding to GetSearchHistory: %v", err)
	}
}

func (a *HttpApi) getSearchQueriesUpdates(w http.ResponseWriter, r *http.Request) {
	userId, ok := userIdFromRequest(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	response, err := a.usecases.GetSearchUpdates(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error happened in usecases.GetSearchQueriesUpdates: %v", err)
		return
	}

	if err := respondWithJSON(w, response, http.StatusOK); err != nil {
		log.Printf("Error happened while responding to GetSearchQueriesUpdates: %v", err)
	}
}

func (a *HttpApi) getArticlesUpdates(w http.ResponseWriter, r *http.Request) {
	userId, ok := userIdFromRequest(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	result, err := a.usecases.GetArticleUpdates(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error happened in usecases.GetArticlesUpdates: %v", err)
		return
	}

	response := make([]ArticleMetaResponse, len(result))
	for i := range result {
		response[i] = renderArticleMeta(result[i])
	}

	if err := respondWithJSON(w, response, http.StatusOK); err != nil {
		log.Printf("Error happened while responding to GetArticlesUpdates: %v", err)
	}
}

func (a *HttpApi) getArticleSubscriptionStatus(w http.ResponseWriter, r *http.Request) {
	articleId := model.ArticleId(mux.Vars(r)["articleId"])
	userId, ok := userIdFromRequest(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	result, err := a.usecases.CheckArticleSubscription(userId, articleId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error happened in usecases.GetArticleSubscriptionStatus: %v", err)
		return
	}

	if result != nil {
		err = respondWithJSON(w, renderUserArticleSubscription(*result), http.StatusAccepted)
	} else {
		err = respondWithJSON(w, struct{}{}, http.StatusAccepted)
	}

	if err != nil {
		log.Printf("Error happened while responding to GetArticlesSubscriptionStatus: %v", err)
	}
}

func (a *HttpApi) postArticleSubscriptionStatus(w http.ResponseWriter, r *http.Request) {
	articleId := model.ArticleId(mux.Vars(r)["articleId"])
	userId, ok := userIdFromRequest(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	result, err := a.usecases.SubscribeForArticle(userId, articleId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error happened in usecases.PostArticleSubscriptionStatus: %v", err)
		return
	}

	if err := respondWithJSON(w, renderUserArticleSubscription(result), http.StatusAccepted); err != nil {
		log.Printf("Error happened while responding to PostArticleSubscriptionStatus: %v", err)
	}
}

func (a *HttpApi) deleteArticleSubscriptionStatus(w http.ResponseWriter, r *http.Request) {
	articleId := model.ArticleId(mux.Vars(r)["articleId"])

	userId, ok := userIdFromRequest(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err := a.usecases.UnsubscribeFromArticle(userId, articleId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error happened in usecases.PostArticleSubscriptionStatus: %v", err)
		return
	}

	if err := respondWithJSON(w, struct{}{}, http.StatusAccepted); err != nil {
		log.Printf("Error happened while responding to PostArticleSubscriptionStatus: %v", err)
	}
}

func (a *HttpApi) getSearchQuerySubscriptionStatus(w http.ResponseWriter, r *http.Request) {
	query := mux.Vars(r)["query"]
	userId, ok := userIdFromRequest(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	result, err := a.usecases.CheckSearchSubscription(userId, query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error happened in usecases.GetSearchQuerySubscriptionStatus: %v", err)
		return
	}

	if result != nil {
		err = respondWithJSON(w, renderUserSearchSubscription(*result), http.StatusOK)
	} else {
		err = respondWithJSON(w, struct{}{}, http.StatusOK)
	}

	if err != nil {
		log.Printf("Error happened while responding to GetSearchQuerySubscriptionStatus: %v", err)
	}
}

func (a *HttpApi) postSearchQuerySubscriptionStatus(w http.ResponseWriter, r *http.Request) {
	query := mux.Vars(r)["query"]

	userId, ok := userIdFromRequest(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	result, err := a.usecases.SubscribeForSearch(userId, query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error happened in usecases.PostSearchQuerySubscriptionStatus: %v", err)
		return
	}

	if err := respondWithJSON(w, renderUserSearchSubscription(result), http.StatusAccepted); err != nil {
		log.Printf("Error happened while responding to PostSearchQuerySubscriptionStatus %v", err)
	}
}

func (a *HttpApi) deleteSearchQuerySubscriptionStatus(w http.ResponseWriter, r *http.Request) {
	query := mux.Vars(r)["query"]
	userId, ok := userIdFromRequest(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	err := a.usecases.UnsubscribeFromSearch(userId, query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error happened in usecases.PostSearchQuerySubscriptionStatus: %v", err)
		return
	}

	if err := respondWithJSON(w, struct{}{}, http.StatusAccepted); err != nil {
		log.Printf("Error happened while responding to PostSearchQuerySubscriptionStatus %v", err)
	}
}

func extractTokenFromAuthHeader(r *http.Request) (usecases.AuthToken, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", nil
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != bearer {
		return "", errors.New("incorrect authorization header format")
	}
	return usecases.AuthToken(parts[1]), nil
}

func (a *HttpApi) extractIdFromHeader(r *http.Request) (model.UserId, error) {
	token, err := extractTokenFromAuthHeader(r)
	if err != nil {
		return "", err
	}
	if token == "" {
		return "", nil
	}
	userId, err := a.usecases.Decode(token)
	if err != nil {
		return "", err
	}
	return userId, nil
}

func (a *HttpApi) extractAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := a.extractIdFromHeader(r)
		if err != nil || userId == "" {
			handler(w, r)
			return
			//TODO handle errors?
		}
		handler(w, r.WithContext(context.WithValue(r.Context(), contextKeyUserId, userId)))
	}
}
