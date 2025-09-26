CREATE TABLE public.$Domain
(
    id          bigserial primary key ,
    title       varchar   NOT NULL,
    description varchar   NULL,
    created_at  timetz    NOT NULL,
    updated_at  timetz    NOT NULL,
    deleted_at  timetz    NULL
);