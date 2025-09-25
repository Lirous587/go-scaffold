-- 安装扩展
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- 用户表
CREATE TABLE users
(
    "id"            BIGSERIAL PRIMARY KEY,
    "email"         VARCHAR(255) UNIQUE NOT NULL,
    "password_hash" VARCHAR(255),
    "name"          VARCHAR(255)        NOT NULL,
    "github_id"     VARCHAR(100) UNIQUE,
    "last_login_at" TIMESTAMP WITH TIME ZONE
);

-- 用户表索引
CREATE INDEX idx_users_email ON users (email);
CREATE INDEX idx_users_github_id ON users (github_id) WHERE github_id IS NOT NULL;

-- labels 表
CREATE TABLE "labels"
(
    "id"          BIGSERIAL PRIMARY KEY,
    "name"        VARCHAR(30) NOT NULL,
    "description" VARCHAR(60)
);

-- articles 表
CREATE TABLE "articles"
(
    "id"            BIGSERIAL PRIMARY KEY,
    "created_at"    TIMESTAMP   NOT NULL DEFAULT now(),
    "updated_at"    TIMESTAMP   NOT NULL DEFAULT now(),
    "deleted_at"    TIMESTAMP,
    "title"         VARCHAR(30) NOT NULL,
    "description"   VARCHAR(160),
    "content"       TEXT        NOT NULL,
    "preview_theme" VARCHAR(20) NOT NULL DEFAULT 'cyanosis',
    "code_theme"    VARCHAR(20) NOT NULL DEFAULT 'github',
    "img_url"       VARCHAR(100),
    "visited_times" BIGINT      NOT NULL DEFAULT 1,
    "priority"      INT2        NOT NULL DEFAULT 50
);
CREATE INDEX "idx_articles_priority" ON "articles" ("priority");
CREATE INDEX "idx_articles_deleted_at" ON "articles" ("deleted_at");
CREATE INDEX "idx_articles_title" ON "articles" ("title");

-- articles_labels 表
CREATE TABLE "articles_labels"
(
    "article_id" BIGINT NOT NULL,
    "label_id"   BIGINT NOT NULL,
    PRIMARY KEY ("article_id", "label_id"),
    CONSTRAINT "fk_articles_labels_article" FOREIGN KEY ("article_id") REFERENCES "articles" ("id"),
    CONSTRAINT "fk_articles_labels_label" FOREIGN KEY ("label_id") REFERENCES "labels" ("id")
);
CREATE INDEX "fk_articles_labels_label" ON "articles_labels" ("label_id");

-- maxims 表
CREATE TABLE "maxims"
(
    "id"         BIGSERIAL PRIMARY KEY,
    "created_at" TIMESTAMP    NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMP    NOT NULL DEFAULT now(),
    "deleted_at" TIMESTAMP,
    "content"    VARCHAR(255) NOT NULL,
    "author"     VARCHAR(30),
    "note"       VARCHAR(255),
    "avatar_url" VARCHAR(255) NOT NULL
);
CREATE INDEX "idx_maxims_deleted_at" ON "maxims" ("deleted_at");

-- moments 表
CREATE TABLE "moments"
(
    "id"          BIGSERIAL PRIMARY KEY,
    "created_at"  TIMESTAMP   NOT NULL DEFAULT now(),
    "updated_at"  TIMESTAMP   NOT NULL DEFAULT now(),
    "deleted_at"  TIMESTAMP,
    "title"       VARCHAR(30) NOT NULL,
    "content"     TEXT        NOT NULL,
    "location"    VARCHAR(30),
    "coordinates" POINT,
    "cover_url"   VARCHAR(255)
);
CREATE INDEX " idx_moments_deleted_at" ON "moments" ("deleted_at");

-- friendlinks 表
CREATE TYPE friendlink_status AS ENUM (
    'pending', -- 待审核
    'approved' -- 已通过
    );
CREATE TABLE "friendlinks"
(
    "id"          BIGSERIAL PRIMARY KEY,
    "created_at"  TIMESTAMP         NOT NULL DEFAULT now(),
    "updated_at"  TIMESTAMP         NOT NULL DEFAULT now(),
    "description" VARCHAR(80)       NOT NULL,
    "site_name"   VARCHAR(80)       NOT NULL,
    "url"         VARCHAR(120)      NOT NULL,
    "logo"        VARCHAR(120)      NOT NULL,
    "status"      friendlink_status NOT NULL DEFAULT 'pending',
    "email"       VARCHAR(80)       NOT NULL,
    -- 联合唯一约束：同一URL在同一状态下只能有一条记录
    CONSTRAINT unique_url_status UNIQUE (url, status)
);
CREATE INDEX "idx_friendlinks_url" ON "friendlinks" ("url");
CREATE INDEX "idx_friendlinks_status" ON "friendlinks" ("status");


-- img 表
CREATE TABLE img
(
    "id"          BIGSERIAL PRIMARY KEY,
    "path"        VARCHAR(120)                   NOT NULL,
    "description" VARCHAR(60),
    "created_at" TIMESTAMP                      NOT NULL DEFAULT now(),
    "updated_at"  TIMESTAMP  NOT NULL,
    "deleted_at"  TIMESTAMP,
    "category_id" integer
);

CREATE UNIQUE INDEX idx_img_path ON img (path);
CREATE INDEX idx_img_deleted_at ON img (deleted_at);
CREATE INDEX idx_img_description_trgm ON img USING gin (description gin_trgm_ops);


CREATE TABLE
    img_category
(
    "id"         BIGSERIAL PRIMARY KEY,
    "title"      VARCHAR(10) NOT NULL UNIQUE,
    "prefix"     VARCHAR(20) NOT NULL,
    "created_at" TIMESTAMP   NOT NULL DEFAULT now()
);


ALTER TABLE img
    ADD CONSTRAINT fk_img_category
        FOREIGN KEY (category_id) REFERENCES img_category (id);