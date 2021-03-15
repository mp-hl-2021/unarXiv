package server

import (
    "encoding/json"
    "github.com/gorilla/mux"
    "github.com/mp-hl-2021/unarXiv/api"
    "github.com/mp-hl-2021/unarXiv/core"
    "log"
    "net/http"
    "strconv"

    //"github.com/dgrijalva/jwt-go"
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
    router.Path("/subscriptions/articles/{articleId}").Queries("subscribe", "{subscribe:(?:true|false)}").
        HandlerFunc(server.PostArticleSubscriptionStatus).Methods(http.MethodPost)

    router.Path("/subscriptions/searches").Queries("query", "{query}").
        HandlerFunc(server.GetSearchQuerySubscriptionStatus).Methods(http.MethodGet)
    router.Path("/subscriptions/searches").Queries("query", "{query}", "subscribe", "{subscribe:(?:true|false)}").
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

    if err := RespondWithJSON(w, &authData, http.StatusCreated); err != nil {
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

    if err := RespondWithJSON(w, &authData, http.StatusOK); err != nil {
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
    offset, err := strconv.Atoi(r.Form.Get("offset"))
    if err != nil && offset > 0 {
        var u32offset uint32
        u32offset = uint32(offset)
        searchQueryRequest.Offset = &u32offset
    }
    searchQueryRequest.AuthData = &core.DummyAuthenticationData // TODO extract auth from headers

    response, err := server.Core.Search(&searchQueryRequest)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        log.Printf("Error happened in Core.Search: %v", err)
        return
    }

    if err := RespondWithJSON(w, response, http.StatusOK); err != nil {
        log.Printf("Error happened while responding to GetSearch: %v", err)
    }
}

func (server *UnarXivServer) GetArticle(w http.ResponseWriter, r *http.Request) {
    var articleRequest api.AccessArticleRequest
    if err := r.ParseForm(); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        log.Printf("Error happened while parsing form params: %v", err)
        return
    }
    articleRequest.ArticleId = r.Form.Get("articleId")
    articleRequest.AuthData = &core.DummyAuthenticationData // TODO extract auth from headers

    response, err := server.Core.AccessArticle(&articleRequest)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        log.Printf("Error happened in Core.AccessArticle: %v", err)
        return
    }

    if err := RespondWithJSON(w, &response, http.StatusOK); err != nil {
        log.Printf("Error happened while responding to GetArticle: %v", err)
    }
}

func (server *UnarXivServer) GetArticlesHistory(w http.ResponseWriter, r *http.Request) {
    authData := core.DummyAuthenticationData // TODO extract auth from headers
    response, err := server.Core.GetArticlesHistory(&authData)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        log.Printf("Error happened in Core.GetArticlesHistory: %v", err)
        return
    }

    if err := RespondWithJSON(w, &response, http.StatusOK); err != nil {
        log.Printf("Error happened while responding to GetArticlesHistory: %v", err)
    }
}

func (server *UnarXivServer) GetSearchHistory(w http.ResponseWriter, r *http.Request) {
    authData := core.DummyAuthenticationData // TODO extract auth from headers
    response, err := server.Core.GetSearchHistory(&authData)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        log.Printf("Error happened in Core.GetSearchHistory: %v", err)
        return
    }

    if err := RespondWithJSON(w, &response, http.StatusOK); err != nil {
        log.Printf("Error happened while responding to GetSearchHistory: %v", err)
    }
}

func (server *UnarXivServer) GetSearchQueriesUpdates(w http.ResponseWriter, r *http.Request) {
    authData := core.DummyAuthenticationData // TODO extract auth from headers
    response, err := server.Core.GetSearchQueriesUpdates(&authData)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        log.Printf("Error happened in Core.GetSearchQueriesUpdates: %v", err)
        return
    }

    if err := RespondWithJSON(w, &response, http.StatusOK); err != nil {
        log.Printf("Error happened while responding to GetSearchQueriesUpdates: %v", err)
    }
}

func (server *UnarXivServer) GetArticlesUpdates(w http.ResponseWriter, r *http.Request) {
    authData := core.DummyAuthenticationData // TODO extract auth from headers
    response, err := server.Core.GetArticlesUpdates(&authData)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        log.Printf("Error happened in Core.GetArticlesUpdates: %v", err)
        return
    }

    if err := RespondWithJSON(w, &response, http.StatusOK); err != nil {
        log.Printf("Error happened while responding to GetArticlesUpdates: %v", err)
    }
}

func (server *UnarXivServer) GetArticleSubscriptionStatus(w http.ResponseWriter, r *http.Request) {
    var getArticleSubscriptionStatusRequest api.GetArticleSubscriptionStatusRequest
    if err := r.ParseForm(); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        log.Printf("Error happened while parsing form params: %v", err)
        return
    }
    getArticleSubscriptionStatusRequest.ArticleId = r.Form.Get("articleId")
    getArticleSubscriptionStatusRequest.AuthData = core.DummyAuthenticationData // TODO extract auth from headers

    response, err := server.Core.GetArticleSubscriptionStatus(&getArticleSubscriptionStatusRequest)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        log.Printf("Error happened in Core.GetArticleSubscriptionStatus: %v", err)
        return
    }

    if err := RespondWithJSON(w, &response, http.StatusOK); err != nil {
        log.Printf("Error happened while responding to GetArticlesSubscriptionStatus: %v", err)
    }
}

func (server *UnarXivServer) PostArticleSubscriptionStatus(w http.ResponseWriter, r *http.Request) {
    var setArticleSubscriptionStatusRequest api.SetArticleSubscriptionStatusRequest

    if err := r.ParseForm(); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        log.Printf("Error happened while parsing form params: %v", err)
        return
    }
    setArticleSubscriptionStatusRequest.ArticleId = r.Form.Get("articleId")
    subscribe := r.Form.Get("subscribe")
    if subscribe == "true" {
        setArticleSubscriptionStatusRequest.Subscribe = true
    } else if subscribe == "false" {
        setArticleSubscriptionStatusRequest.Subscribe = false
    } else {
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    setArticleSubscriptionStatusRequest.AuthData = core.DummyAuthenticationData // TODO extract auth from headers

    response, err := server.Core.SetArticleSubscriptionStatus(&setArticleSubscriptionStatusRequest)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        log.Printf("Error happened in Core.SetArticleSubscriptionStatus: %v", err)
        return
    }

    if err := RespondWithJSON(w, &response, http.StatusAccepted); err != nil {
        log.Printf("Error happened while responding to PostArticleSubscriptionStatus: %v", err)
    }
}

func (server *UnarXivServer) GetSearchQuerySubscriptionStatus(w http.ResponseWriter, r *http.Request) {
    var getSearchQuerySubscriptionStatusRequest api.GetSearchQuerySubscriptionStatusRequest
    if err := r.ParseForm(); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        log.Printf("Error happened while parsing form params: %v", err)
        return
    }
    getSearchQuerySubscriptionStatusRequest.Query = r.Form.Get("query")
    getSearchQuerySubscriptionStatusRequest.AuthData = core.DummyAuthenticationData // TODO extract auth from headers

    response, err := server.Core.GetSearchQuerySubscriptionStatus(&getSearchQuerySubscriptionStatusRequest)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        log.Printf("Error happened in Core.GetSearchQuerySubscriptionStatus: %v", err)
        return
    }

    if err := RespondWithJSON(w, &response, http.StatusOK); err != nil {
        log.Printf("Error happened while responding to GetSearchQuerySubscriptionStatus: %v", err)
    }
}

func (server *UnarXivServer) PostSearchQuerySubscriptionStatus(w http.ResponseWriter, r *http.Request) {
    var setSearchQuerySubscriptionStatusRequest api.SetSearchQuerySubscriptionStatusRequest
    if err := r.ParseForm(); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        log.Printf("Error happened while parsing form params: %v", err)
        return
    }
    setSearchQuerySubscriptionStatusRequest.Query = r.Form.Get("query")
    subscribe := r.Form.Get("subscribe")
    if subscribe == "true" {
        setSearchQuerySubscriptionStatusRequest.Subscribe = true
    } else if subscribe == "false" {
        setSearchQuerySubscriptionStatusRequest.Subscribe = false
    } else {
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    setSearchQuerySubscriptionStatusRequest.AuthData = core.DummyAuthenticationData // TODO extract auth from headers

    response, err := server.Core.SetSearchQuerySubscriptionStatus(&setSearchQuerySubscriptionStatusRequest)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        log.Printf("Error happened in Core.SetSearchQuerySubscriptionStatus: %v", err)
        return
    }

    if err := RespondWithJSON(w, &response, http.StatusAccepted); err != nil {
        log.Printf("Error happened while responding to PostSearchQuerySubscriptionStatus %v", err)
    }
}
