CREATE TABLE IF NOT EXISTS users(
    id serial  PRIMARY KEY,
    username  varchar(20) unique NOT NULL,
    email  varchar(254) unique NOT NULL,
    password varchar(60) NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL
);