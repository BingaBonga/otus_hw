CREATE TABLE event (
    id UUID not null,
    title varchar(256) not null,
    start_date timestamp not null default now(),
    duration BIGINT not null,
    description varchar(4096),
    owner varchar(256) not null,
    remind_at BIGINT,
    is_send bool default false
);