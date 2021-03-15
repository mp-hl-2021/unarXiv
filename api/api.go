package api

import (
    "encoding/json"
    "github.com/gorilla/mux"
    "github.com/mp-hl-2021/unarXiv/usecases"
    "log"
    "net/http"
    "strconv"

    //"github.com/dgrijalva/jwt-go"
)


type UnarXivApi struct {
    Core usecases.UseCasesInterface
}

func NewApi(useCases usecases.UseCasesInterface) *UnarXivApi {
    return &UnarXivApi{
        Core: useCases,
    }
}

func (a *UnarXivApi) Router() http.Handler {
    router := mux.NewRouter()

    router.HandleFunc("/register", a.PostRegister).Methods(http.MethodPost)
    router.HandleFunc("/login", a.PostLogin).Methods(http.MethodPost)

    // offset is optional, should be passed as "?offset=smth"
    router.Path("/search").Queries("query", "{query}").
        HandlerFunc(a.GetSearch).Methods(http.MethodGet)

    router.HandleFunc("/articles/{articleId}", a.GetArticle).Methods(http.MethodGet)

    router.HandleFunc("/history/searches", a.GetSearchHistory).Methods(http.MethodGet)
    router.HandleFunc("/history/articles", a.GetArticlesHistory).Methods(http.MethodGet)

    router.HandleFunc("/updates/searches", a.GetSearchQueriesUpdates).Methods(http.MethodGet)
    router.HandleFunc("/updates/articles", a.GetArticlesUpdates).Methods(http.MethodGet)

    router.Path("/subscriptions/articles/{articleId}").
        HandlerFunc(a.GetArticleSubscriptionStatus).Methods(http.MethodGet)
    router.Path("/subscriptions/articles/{articleId}").Queries("subscribe", "{subscribe:(?:true|false)}").
        HandlerFunc(a.PostArticleSubscriptionStatus).Methods(http.MethodPost)

    router.Path("/subscriptions/searches").Queries("query", "{query}").
        HandlerFunc(a.GetSearchQuerySubscriptionStatus).Methods(http.MethodGet)
    router.Path("/subscriptions/searches").Queries("query", "{query}", "subscribe", "{subscribe:(?:true|false)}").
        HandlerFunc(a.PostSearchQuerySubscriptionStatus).Methods(http.MethodPost)

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
func (a *UnarXivApi) PostRegister(w http.ResponseWriter, r *http.Request) {
    var registerRequest usecases.AuthenticationRequest
    if err := json.NewDecoder(r.Body).Decode(&registerRequest); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    authData, err := a.Core.Register(&registerRequest)
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
func (a *UnarXivApi) PostLogin(w http.ResponseWriter, r *http.Request) {
    var authRequest usecases.AuthenticationRequest
    if err := json.NewDecoder(r.Body).Decode(&authRequest); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    authData, err := a.Core.Login(&authRequest)
    if err != nil { // todo: map domain errors to http error codes
        w.WriteHeader(http.StatusInternalServerError)
        log.Printf("Error happened in Core.Login: %v", err)
        return
    }

    if err := RespondWithJSON(w, &authData, http.StatusOK); err != nil {
        log.Printf("Error happened while responding to PostLogin: %v", err)
    }
}

func (a *UnarXivApi) GetSearch(w http.ResponseWriter, r *http.Request) {
    var searchQueryRequest = usecases.SearchQueryRequest{}
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
    searchQueryRequest.AuthData = &usecases.DummyAuthenticationData // TODO extract auth from headers

    response, err := a.Core.Search(&searchQueryRequest)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        log.Printf("Error happened in Core.Search: %v", err)
        return
    }

    if err := RespondWithJSON(w, response, http.StatusOK); err != nil {
        log.Printf("Error happened while responding to GetSearch: %v", err)
    }
}

func (a *UnarXivApi) GetArticle(w http.ResponseWriter, r *http.Request) {
    var articleRequest usecases.AccessArticleRequest
    if err := r.ParseForm(); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        log.Printf("Error happened while parsing form params: %v", err)
        return
    }
    articleRequest.ArticleId = r.Form.Get("articleId")
    articleRequest.AuthData = &usecases.DummyAuthenticationData // TODO extract auth from headers

    response, err := a.Core.AccessArticle(&articleRequest)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        log.Printf("Error happened in Core.AccessArticle: %v", err)
        return
    }

    if err := RespondWithJSON(w, &response, http.StatusOK); err != nil {
        log.Printf("Error happened while responding to GetArticle: %v", err)
    }
}

func (a *UnarXivApi) GetArticlesHistory(w http.ResponseWriter, r *http.Request) {
    authData := usecases.DummyAuthenticationData // TODO extract auth from headers
    response, err := a.Core.GetArticlesHistory(&authData)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        log.Printf("Error happened in Core.GetArticlesHistory: %v", err)
        return
    }

    if err := RespondWithJSON(w, &response, http.StatusOK); err != nil {
        log.Printf("Error happened while responding to GetArticlesHistory: %v", err)
    }
}

func (a *UnarXivApi) GetSearchHistory(w http.ResponseWriter, r *http.Request) {
    authData := usecases.DummyAuthenticationData // TODO extract auth from headers
    response, err := a.Core.GetSearchHistory(&authData)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        log.Printf("Error happened in Core.GetSearchHistory: %v", err)
        return
    }

    if err := RespondWithJSON(w, &response, http.StatusOK); err != nil {
        log.Printf("Error happened while responding to GetSearchHistory: %v", err)
    }
}

func (a *UnarXivApi) GetSearchQueriesUpdates(w http.ResponseWriter, r *http.Request) {
    authData := usecases.DummyAuthenticationData // TODO extract auth from headers
    response, err := a.Core.GetSearchQueriesUpdates(&authData)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        log.Printf("Error happened in Core.GetSearchQueriesUpdates: %v", err)
        return
    }

    if err := RespondWithJSON(w, &response, http.StatusOK); err != nil {
        log.Printf("Error happened while responding to GetSearchQueriesUpdates: %v", err)
    }
}

func (a *UnarXivApi) GetArticlesUpdates(w http.ResponseWriter, r *http.Request) {
    authData := usecases.DummyAuthenticationData // TODO extract auth from headers
    response, err := a.Core.GetArticlesUpdates(&authData)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        log.Printf("Error happened in Core.GetArticlesUpdates: %v", err)
        return
    }

    if err := RespondWithJSON(w, &response, http.StatusOK); err != nil {
        log.Printf("Error happened while responding to GetArticlesUpdates: %v", err)
    }
}

func (a *UnarXivApi) GetArticleSubscriptionStatus(w http.ResponseWriter, r *http.Request) {
    var getArticleSubscriptionStatusRequest usecases.GetArticleSubscriptionStatusRequest
    if err := r.ParseForm(); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        log.Printf("Error happened while parsing form params: %v", err)
        return
    }
    getArticleSubscriptionStatusRequest.ArticleId = r.Form.Get("articleId")
    getArticleSubscriptionStatusRequest.AuthData = usecases.DummyAuthenticationData // TODO extract auth from headers

    response, err := a.Core.GetArticleSubscriptionStatus(&getArticleSubscriptionStatusRequest)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        log.Printf("Error happened in Core.GetArticleSubscriptionStatus: %v", err)
        return
    }

    if err := RespondWithJSON(w, &response, http.StatusOK); err != nil {
        log.Printf("Error happened while responding to GetArticlesSubscriptionStatus: %v", err)
    }
}

func (a *UnarXivApi) PostArticleSubscriptionStatus(w http.ResponseWriter, r *http.Request) {
    var setArticleSubscriptionStatusRequest usecases.SetArticleSubscriptionStatusRequest

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
    setArticleSubscriptionStatusRequest.AuthData = usecases.DummyAuthenticationData // TODO extract auth from headers

    response, err := a.Core.SetArticleSubscriptionStatus(&setArticleSubscriptionStatusRequest)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        log.Printf("Error happened in Core.SetArticleSubscriptionStatus: %v", err)
        return
    }

    if err := RespondWithJSON(w, &response, http.StatusAccepted); err != nil {
        log.Printf("Error happened while responding to PostArticleSubscriptionStatus: %v", err)
    }
}

func (a *UnarXivApi) GetSearchQuerySubscriptionStatus(w http.ResponseWriter, r *http.Request) {
    var getSearchQuerySubscriptionStatusRequest usecases.GetSearchQuerySubscriptionStatusRequest
    if err := r.ParseForm(); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        log.Printf("Error happened while parsing form params: %v", err)
        return
    }
    getSearchQuerySubscriptionStatusRequest.Query = r.Form.Get("query")
    getSearchQuerySubscriptionStatusRequest.AuthData = usecases.DummyAuthenticationData // TODO extract auth from headers

    response, err := a.Core.GetSearchQuerySubscriptionStatus(&getSearchQuerySubscriptionStatusRequest)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        log.Printf("Error happened in Core.GetSearchQuerySubscriptionStatus: %v", err)
        return
    }

    if err := RespondWithJSON(w, &response, http.StatusOK); err != nil {
        log.Printf("Error happened while responding to GetSearchQuerySubscriptionStatus: %v", err)
    }
}

func (a *UnarXivApi) PostSearchQuerySubscriptionStatus(w http.ResponseWriter, r *http.Request) {
    var setSearchQuerySubscriptionStatusRequest usecases.SetSearchQuerySubscriptionStatusRequest
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
    setSearchQuerySubscriptionStatusRequest.AuthData = usecases.DummyAuthenticationData // TODO extract auth from headers

    response, err := a.Core.SetSearchQuerySubscriptionStatus(&setSearchQuerySubscriptionStatusRequest)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        log.Printf("Error happened in Core.SetSearchQuerySubscriptionStatus: %v", err)
        return
    }

    if err := RespondWithJSON(w, &response, http.StatusAccepted); err != nil {
        log.Printf("Error happened while responding to PostSearchQuerySubscriptionStatus %v", err)
    }
}
