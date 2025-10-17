-- 安装扩展
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- 用户表
CREATE TABLE public.users
(
    id            bigserial      NOT NULL PRIMARY KEY,
    nickname      varchar(20)    NOT NULL,
    email         varchar(80)    NOT NULL UNIQUE,
    github_id     varchar(60)    NULL UNIQUE,
    --     google_id     varchar(60) NULL UNIQUE,
    password_hash text           NULL,
    created_at    timestamptz(6) NOT NULL DEFAULT now(),
    updated_at    timestamptz(6) NOT NULL DEFAULT now(),
    last_login_at timestamptz(6) NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_users_created_at ON public.users (created_at);
CREATE INDEX IF NOT EXISTS idx_users_nickname ON public.users (nickname);
CREATE INDEX idx_users_github_id ON public.users (github_id) WHERE github_id IS NOT NULL;
