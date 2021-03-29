package httpapi

import (
    "encoding/json"
    "github.com/gorilla/mux"
	"github.com/mp-hl-2021/unarXiv/internal/usecases"
    "log"
    "net/http"
    "strconv"

    //"github.com/dgrijalva/jwt-go"
)


type UnarXivApi struct {
    Core usecases.Interface
}

func NewApi(useCases usecases.Interface) *UnarXivApi {
    return &UnarXivApi{
        Core: useCases,
    }
}

func (a *UnarXivApi) Router() http.Handler {
    router := mux.NewRouter()

    router.HandleFunc("/register", a.postRegister).Methods(http.MethodPost)
    router.HandleFunc("/login", a.postLogin).Methods(http.MethodPost)

    // offset is optional, should be passed as "?offset=smth"
    router.Path("/search/{query}").HandlerFunc(a.getSearch).Methods(http.MethodGet)

    router.HandleFunc("/articles/{articleId}", a.getArticle).Methods(http.MethodGet)

    router.HandleFunc("/history/searches", a.getSearchHistory).Methods(http.MethodGet)
    router.HandleFunc("/history/articles", a.getArticlesHistory).Methods(http.MethodGet)

    router.HandleFunc("/updates/searches", a.getSearchQueriesUpdates).Methods(http.MethodGet)
    router.HandleFunc("/updates/articles", a.getArticlesUpdates).Methods(http.MethodGet)

    router.Path("/subscriptions/articles/{articleId}").
        HandlerFunc(a.getArticleSubscriptionStatus).Methods(http.MethodGet)
    router.Path("/subscriptions/articles/{articleId}").
        HandlerFunc(a.setArticleSubscriptionStatusMaker(true)).Methods(http.MethodPost)
    router.Path("/subscriptions/articles/{articleId}").
        HandlerFunc(a.setArticleSubscriptionStatusMaker(false)).Methods(http.MethodDelete)

    router.Path("/subscriptions/searches/{query}").
        HandlerFunc(a.getSearchQuerySubscriptionStatus).Methods(http.MethodGet)
    router.Path("/subscriptions/searches/{query}").
        HandlerFunc(a.setSearchQuerySubscriptionStatusMaker(true)).Methods(http.MethodPost)
    router.Path("/subscriptions/searches/{query}").
        HandlerFunc(a.setSearchQuerySubscriptionStatusMaker(false)).Methods(http.MethodDelete)

    return router
}

func respondWithJSON(w http.ResponseWriter, object interface{}, status int) error {
    w.Header().Set("Content-Type", "application/json")

    if object != nil { // non-empty body
        if err := json.NewEncoder(w).Encode(object); err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            return err
        }
    }

    if status != http.StatusOK {
        w.WriteHeader(status)
    }
    return nil
}

// postRegister handles request for a new account creation.
func (a *UnarXivApi) postRegister(w http.ResponseWriter, r *http.Request) {
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

    if err := respondWithJSON(w, &authData, http.StatusCreated); err != nil {
        log.Printf("Error happened while responding to PostRegister: %v", err)
    }
}

// PostLogin handles login request for existing user.
func (a *UnarXivApi) postLogin(w http.ResponseWriter, r *http.Request) {
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

    if err := respondWithJSON(w, &authData, http.StatusOK); err != nil {
        log.Printf("Error happened while responding to PostLogin: %v", err)
    }
}

func (a *UnarXivApi) getSearch(w http.ResponseWriter, r *http.Request) {
    var searchQueryRequest = usecases.SearchQueryRequest{}
    if err := r.ParseForm(); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        log.Printf("Error happened while parsing form params: %v", err)
        return
    }
    searchQueryRequest.Query = r.Form.Get("query")
    if strOffset := r.Form.Get("offset"); len(strOffset) != 0 {
        offset, err := strconv.Atoi(strOffset)
        if err != nil || offset < 0 {
            w.WriteHeader(http.StatusBadRequest)
            return
        }
        searchQueryRequest.Offset = uint32(offset)
    }
    searchQueryRequest.AuthData = &usecases.DummyAuthenticationData // TODO extract auth from headers

    response, err := a.Core.Search(&searchQueryRequest)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        log.Printf("Error happened in Core.Search: %v", err)
        return
    }

    if err := respondWithJSON(w, response, http.StatusOK); err != nil {
        log.Printf("Error happened while responding to GetSearch: %v", err)
    }
}

func (a *UnarXivApi) getArticle(w http.ResponseWriter, r *http.Request) {
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

    if err := respondWithJSON(w, &response, http.StatusOK); err != nil {
        log.Printf("Error happened while responding to GetArticle: %v", err)
    }
}

func (a *UnarXivApi) getArticlesHistory(w http.ResponseWriter, r *http.Request) {
    authData := usecases.DummyAuthenticationData // TODO extract auth from headers
    response, err := a.Core.GetArticlesHistory(&authData)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        log.Printf("Error happened in Core.GetArticlesHistory: %v", err)
        return
    }

    if err := respondWithJSON(w, &response, http.StatusOK); err != nil {
        log.Printf("Error happened while responding to GetArticlesHistory: %v", err)
    }
}

func (a *UnarXivApi) getSearchHistory(w http.ResponseWriter, r *http.Request) {
    authData := usecases.DummyAuthenticationData // TODO extract auth from headers
    response, err := a.Core.GetSearchHistory(&authData)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        log.Printf("Error happened in Core.GetSearchHistory: %v", err)
        return
    }

    if err := respondWithJSON(w, &response, http.StatusOK); err != nil {
        log.Printf("Error happened while responding to GetSearchHistory: %v", err)
    }
}

func (a *UnarXivApi) getSearchQueriesUpdates(w http.ResponseWriter, r *http.Request) {
    authData := usecases.DummyAuthenticationData // TODO extract auth from headers
    response, err := a.Core.GetSearchQueriesUpdates(&authData)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        log.Printf("Error happened in Core.GetSearchQueriesUpdates: %v", err)
        return
    }

    if err := respondWithJSON(w, &response, http.StatusOK); err != nil {
        log.Printf("Error happened while responding to GetSearchQueriesUpdates: %v", err)
    }
}

func (a *UnarXivApi) getArticlesUpdates(w http.ResponseWriter, r *http.Request) {
    authData := usecases.DummyAuthenticationData // TODO extract auth from headers
    response, err := a.Core.GetArticlesUpdates(&authData)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        log.Printf("Error happened in Core.GetArticlesUpdates: %v", err)
        return
    }

    if err := respondWithJSON(w, &response, http.StatusOK); err != nil {
        log.Printf("Error happened while responding to GetArticlesUpdates: %v", err)
    }
}

func (a *UnarXivApi) getArticleSubscriptionStatus(w http.ResponseWriter, r *http.Request) {
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

    if err := respondWithJSON(w, &response, http.StatusOK); err != nil {
        log.Printf("Error happened while responding to GetArticlesSubscriptionStatus: %v", err)
    }
}

func (a *UnarXivApi) setArticleSubscriptionStatusMaker(status bool) func (w http.ResponseWriter, r *http.Request) {
    return func (w http.ResponseWriter, r *http.Request) {
        var setArticleSubscriptionStatusRequest usecases.SetArticleSubscriptionStatusRequest

        if err := r.ParseForm(); err != nil {
            w.WriteHeader(http.StatusBadRequest)
            log.Printf("Error happened while parsing form params: %v", err)
            return
        }
        setArticleSubscriptionStatusRequest.ArticleId = r.Form.Get("articleId")
        setArticleSubscriptionStatusRequest.Subscribe = status
        setArticleSubscriptionStatusRequest.AuthData = usecases.DummyAuthenticationData // TODO extract auth from headers

        response, err := a.Core.SetArticleSubscriptionStatus(&setArticleSubscriptionStatusRequest)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            log.Printf("Error happened in Core.SetArticleSubscriptionStatus(status=%v): %v", status, err)
            return
        }

        if err := respondWithJSON(w, &response, http.StatusAccepted); err != nil {
            log.Printf("Error happened while responding to SetArticleSubscriptionStatus(status=%v): %v", status, err)
        }
    }
}

func (a *UnarXivApi) getSearchQuerySubscriptionStatus(w http.ResponseWriter, r *http.Request) {
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

    if err := respondWithJSON(w, &response, http.StatusOK); err != nil {
        log.Printf("Error happened while responding to GetSearchQuerySubscriptionStatus: %v", err)
    }
}

func (a *UnarXivApi) setSearchQuerySubscriptionStatusMaker(status bool) func (w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        var setSearchQuerySubscriptionStatusRequest usecases.SetSearchQuerySubscriptionStatusRequest
        if err := r.ParseForm(); err != nil {
            w.WriteHeader(http.StatusBadRequest)
            log.Printf("Error happened while parsing form params: %v", err)
            return
        }
        setSearchQuerySubscriptionStatusRequest.Query = r.Form.Get("query")
        setSearchQuerySubscriptionStatusRequest.Subscribe = status
        setSearchQuerySubscriptionStatusRequest.AuthData = usecases.DummyAuthenticationData // TODO extract auth from headers

        response, err := a.Core.SetSearchQuerySubscriptionStatus(&setSearchQuerySubscriptionStatusRequest)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            log.Printf("Error happened in Core.SetSearchQuerySubscriptionStatus(status=%v): %v", status, err)
            return
        }

        if err := respondWithJSON(w, &response, http.StatusAccepted); err != nil {
            log.Printf("Error happened while responding to SetSearchQuerySubscriptionStatus(status=%v) %v", status, err)
        }
    }
}
