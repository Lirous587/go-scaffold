-- 安装扩展
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- 用户表
CREATE TABLE public.users
(
    id            bigserial   NOT NULL PRIMARY KEY,
    nickname      varchar(20) NOT NULL,
    email         varchar(80) NOT NULL UNIQUE,
    github_id     varchar(60) NULL UNIQUE,
    --     google_id     varchar(60) NULL UNIQUE,
    password_hash text        NULL,
    created_at    timestamp   NOT NULL DEFAULT now(),
    updated_at    timestamp   NOT NULL DEFAULT now(),
    last_login_at TIMESTAMP WITH TIME ZONE
);
CREATE INDEX IF NOT EXISTS idx_users_created_at ON public.users (created_at);
CREATE INDEX IF NOT EXISTS idx_users_nickname ON public.users (nickname);
CREATE INDEX idx_users_github_id ON public.users (github_id) WHERE github_id IS NOT NULL;

-- img_categories
CREATE TABLE public.img_categories
(
    id         bigserial PRIMARY KEY,
    title      varchar(10) NOT NULL UNIQUE,
    prefix     varchar(20) NOT NULL,
    created_at timestamp   NOT NULL DEFAULT now()
);

-- img 表
CREATE TABLE public.imgs
(
    id          bigserial PRIMARY KEY,
    path        varchar(120) NOT NULL,
    description varchar(60),
    created_at  timestamp    NOT NULL DEFAULT now(),
    updated_at  timestamp    NOT NULL,
    deleted_at  timestamp,
    category_id bigint REFERENCES img_categories (id)
);
CREATE UNIQUE INDEX idx_img_path ON public.imgs (path);
CREATE INDEX idx_img_deleted_at ON public.imgs (deleted_at);
CREATE INDEX idx_img_description_trgm ON public.imgs USING gin (description gin_trgm_ops);