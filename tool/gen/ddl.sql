CREATE TABLE public.$Domain
(
    id          bigserial primary key ,
    title       varchar   NOT NULL,
    description varchar   NULL,
    created_at  timestamptz(6)    NOT NULL,
    updated_at  timestamptz(6)    NOT NULL,
    deleted_at  timestamptz(6)    NULL
);