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
    Id integer primary key,
    DesiredArticleCount integer,
    RootURL text
);
INSERT INTO CrawlerConfig VALUES (0, 50, 'https://arxiv.org/') ON CONFLICT DO NOTHING;

CREATE TABLE IF NOT EXISTS CrawlStatus (
    URL text not null primary key,
    LastAccess bigint not null,
    LastHTTPStatus integer not null
);
