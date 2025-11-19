-- +goose Up
create table answers (
    id bigserial primary key,
    question_id bigserial not null references questions(id) on delete cascade,
    user_id uuid not null,
    text text not null,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp
);
create index idx_answers_question_id on answers(question_id);

-- +goose Down
drop index if exists idx_answers_question_id;
drop table if exists answers;
