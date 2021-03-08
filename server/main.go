package main

import (
    "context"
    "net"

    "google.golang.org/grpc"

    "github.com/mp-hl-2021/unarXiv/api"
)

func main() {
    l, err := net.Listen("tcp", "localhost:8081")
    if err != nil {
        panic(err)
    }

    grpcServer := grpc.NewServer()
    api.RegisterUnarXivServer(grpcServer, &UnarXivServer{})

    if err := grpcServer.Serve(l); err != nil {
        panic(err)
    }
}

type UnarXivServer struct {
    api.UnimplementedUnarXivServer
}

func (e *UnarXivServer) Register(_ context.Context, r *api.AuthenticationRequest) (*api.AuthenticationData, error) {
    return &api.AuthenticationData{Jwt: "placeholder"}, nil
}

func (e *UnarXivServer) Login(_ context.Context, r *api.AuthenticationRequest) (*api.AuthenticationData, error) {
    return &api.AuthenticationData{Jwt: "placeholder"}, nil
}

func (e *UnarXivServer) Search(_ context.Context, r *api.SearchQueryRequest) (*api.SearchQueryResponse, error) {
    return &api.SearchQueryResponse{TotalMatchesCount: 0}, nil
}

func (e *UnarXivServer) SetArticleSubscriptionStatus(_ context.Context, r *api.SetArticleSubscriptionStatusRequest) (*api.ArticleSubscriptionStatus, error) {
    return &api.ArticleSubscriptionStatus{AbsId: "0", IsSubscribedNow: false}, nil
}

func (e *UnarXivServer) GetArticleSubscriptionStatus(_ context.Context, r *api.GetArticleSubscriptionStatusRequest) (*api.ArticleSubscriptionStatus, error) {
    return &api.ArticleSubscriptionStatus{AbsId: "0", IsSubscribedNow: false}, nil
}

func (e *UnarXivServer) SetSearchQuerySubscriptionStatus(_ context.Context, r *api.SetSearchQuerySubscriptionStatusRequest) (*api.SearchQuerySubscriptionStatus, error) {
    return &api.SearchQuerySubscriptionStatus{Query: "placeholder", IsSubscribedNow: false}, nil
}

func (e *UnarXivServer) GetSearchQuerySubscriptionStatus(_ context.Context, r *api.GetSearchQuerySubscriptionStatusRequest) (*api.SearchQuerySubscriptionStatus, error) {
    return &api.SearchQuerySubscriptionStatus{Query: "placeholder", IsSubscribedNow: false}, nil
}

func (e *UnarXivServer) GetArticlesUpdates(_ context.Context, r *api.AuthenticationData) (*api.GetArticlesUpdatesResponse, error) {
    return &api.GetArticlesUpdatesResponse{}, nil
}

func (e *UnarXivServer) AccessArticle(_ context.Context, r *api.AccessArticleRequest) (*api.AccessArticleResponse, error) {
    return &api.AccessArticleResponse{Article: &api.ArticleMetaInfo{Title: "placeholder", Authors: "placeholder", AbsId: "0", Abstract: "placeholder", LastUpdateTimestamp: 0}}, nil
}
