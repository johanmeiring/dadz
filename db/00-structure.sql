create table users
(
    id      serial
        constraint users_pk
            primary key,
    name    text not null,
    api_key text not null
);

create index users_api_key_index
    on users (api_key);

create table jokes
(
    id        serial
        constraint jokes_pk
            primary key,
    intro     text not null,
    punchline text not null,
    user_id   integer
        constraint jokes_users_id_fk
            references users
            on update cascade on delete set null
);
