create table if not exists urls
(
    id serial
    constraint urls_pk
    primary key,
    original_url varchar(255) not null,
    short_url varchar(50),
    user_id varchar(50)
    );

create unique index if not exists urls_short_url_uindex
    on urls (short_url);