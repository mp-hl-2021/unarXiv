CREATE TABLE IF NOT EXISTS Accounts (
    Id serial PRIMARY KEY,
    Login text not null,
    Password text not null
);

CREATE TABLE IF NOT EXISTS Articles (
    Id text PRIMARY KEY,
    Title text,
    Abstract text,
    LastUpdateTimestamp bigint,
    FullDocumentURL text
);

CREATE TABLE IF NOT EXISTS ArticlesFTS (
    Id text PRIMARY KEY references Articles(Id),
    TextData tsvector
);
CREATE INDEX IF NOT EXISTS idx_articles_fts_gin ON ArticlesFTS USING gin (TextData);

CREATE TABLE IF NOT EXISTS AuthorsOfArticles (
    ArticleId text REFERENCES Articles (Id),
    AuthorName text
);


CREATE TABLE IF NOT EXISTS AccountArticleRelations (
    UserId integer REFERENCES Accounts (Id),
    ArticleId text REFERENCES Articles (Id),
    IsSubscribed boolean,
    LastAccess bigint
);

CREATE TABLE IF NOT EXISTS AccountSearchRelations (
    UserId integer REFERENCES Accounts (Id),
    Search text,
    IsSubscribed boolean,
    LastAccess bigint
);

CREATE TABLE IF NOT EXISTS CrawlerConfig (
    RootURL text primary key,
    DesiredArticleCount integer
);
INSERT INTO CrawlerConfig VALUES ('http://arxiv.org/', 1000) ON CONFLICT DO NOTHING;

CREATE TABLE IF NOT EXISTS CrawlStatus (
    URL text not null primary key,
    Visited boolean not null,
    LastAccess bigint,
    LastHTTPStatus integer
);

INSERT INTO CrawlStatus (URL, Visited) VALUES ('http://arxiv.org/', false) ON CONFLICT DO NOTHING;

