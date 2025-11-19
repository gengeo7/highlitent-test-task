-- +goose Up
create table questions (
    id bigserial primary key,
    text text not null,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp
);
create index idx_questions_id on questions(id);

-- +goose Down
drop index if exists idx_questions_id;
drop table if exists questions;

