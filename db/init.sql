create table if not exists urls
(
    id serial
    constraint urls_pk
    primary key,
    original_url varchar(255) not null,
    short_url varchar(50),
    user_id varchar(50),
    correlation_id varchar(100)
    );

create unique index if not exists urls_short_url_uindex
    on urls (short_url);

create unique index if not exists urls_correlation_id_uindex
    on urls (correlation_id);

create unique index if not exists urls_original_url_uindex
    on urls (original_url);