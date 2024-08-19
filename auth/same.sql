create table users
(
    id            uuid                     default gen_random_uuid()         not null
        primary key,
    full_name     varchar(255)                                               not null,
    username      varchar(50)                                                not null
        unique,
    email         varchar(255)                                               not null
        unique,
    password_hash varchar(255)                                               not null,
    phone         varchar(20),
    image         varchar(255),
    role          varchar(50)              default 'user'::character varying not null,
    created_at    timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at    timestamp with time zone default CURRENT_TIMESTAMP,
    deleted_at    timestamp
);