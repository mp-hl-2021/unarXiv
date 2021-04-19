package memory

import (
    "github.com/mp-hl-2021/unarXiv/internal/domain"
    "github.com/mp-hl-2021/unarXiv/internal/domain/model"
    "github.com/mp-hl-2021/unarXiv/internal/interface/utils"
    "sync"
    "time"
)

type HistoryRepo struct {
    articleAccess map[model.UserId]map[model.ArticleId]uint64
    searchAccess  map[model.UserId]map[string]uint64
    mutex         *sync.Mutex
}

func NewHistoryRepo() *HistoryRepo {
    return &HistoryRepo{
        articleAccess: make(map[model.UserId]map[model.ArticleId]uint64),
        searchAccess:  make(map[model.UserId]map[string]uint64),
        mutex:         &sync.Mutex{},
    }
}

func (h *HistoryRepo) ArticleAccessOccurred(userId model.UserId, articleId model.ArticleId) error {
    h.mutex.Lock()
    defer h.mutex.Unlock()
    if _, ok := h.articleAccess[userId]; !ok {
        h.articleAccess[userId] = make(map[model.ArticleId]uint64)
    }
    h.articleAccess[userId][articleId] = utils.Uint64Time(time.Now())
    return nil
}

func (h *HistoryRepo) GetArticleLastAccessTimestamp(userId model.UserId, articleId model.ArticleId) (uint64, error) {
    h.mutex.Lock()
    defer h.mutex.Unlock()
    if _, ok := h.articleAccess[userId]; !ok {
        return 0, domain.NeverAccessed
    }
    if t, ok := h.articleAccess[userId][articleId]; !ok {
        return 0, domain.NeverAccessed
    } else {
        return t, nil
    }
}

func (h *HistoryRepo) SearchAccessOccurred(userId model.UserId, query string) error {
    h.mutex.Lock()
    defer h.mutex.Unlock()
    if _, ok := h.searchAccess[userId]; !ok {
        h.searchAccess[userId] = make(map[string]uint64)
    }
    h.searchAccess[userId][query] = utils.Uint64Time(time.Now())
    return nil
}

func (h *HistoryRepo) GetSearchLastAccessTimestamp(userId model.UserId, query string) (uint64, error) {
    h.mutex.Lock()
    defer h.mutex.Unlock()
    if _, ok := h.searchAccess[userId]; !ok {
        return 0, domain.NeverAccessed
    }
    if t, ok := h.searchAccess[userId][query]; !ok {
        return 0, domain.NeverAccessed
    } else {
        return t, nil
    }
}

func (h *HistoryRepo) GetSearchHistory(userId model.UserId) ([]string, error) {
    h.mutex.Lock()
    defer h.mutex.Unlock()
    if history, ok := h.searchAccess[userId]; !ok {
        return []string{}, nil
    } else {
        result := make([]string, 0, len(history))
        for k := range history {
            result = append(result, k)
        }
        return result, nil
    }
}

func (h *HistoryRepo) ClearSearchHistory(userId model.UserId) error {
    delete(h.searchAccess, userId)
    return nil
}

func (h *HistoryRepo) GetArticleHistory(userId model.UserId) ([]model.ArticleId, error) {
    h.mutex.Lock()
    defer h.mutex.Unlock()
    if history, ok := h.articleAccess[userId]; !ok {
        return []model.ArticleId{}, nil
    } else {
        result := make([]model.ArticleId, 0, len(history))
        for k := range history {
            result = append(result, k)
        }
        return result, nil
    }
}

func (h *HistoryRepo) ClearArticleHistory(userId model.UserId) error {
    delete(h.articleAccess, userId)
    return nil
}
