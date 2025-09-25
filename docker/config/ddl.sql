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

-- img 表
CREATE TABLE img
(
    "id"          BIGSERIAL PRIMARY KEY,
    "path"        VARCHAR(120) NOT NULL,
    "description" VARCHAR(60),
    "created_at"  TIMESTAMP    NOT NULL DEFAULT now(),
    "updated_at"  TIMESTAMP    NOT NULL,
    "deleted_at"  TIMESTAMP,
    "category_id" BIGINT
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