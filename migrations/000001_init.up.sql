create table users (
  id serial not null unique,
  username varchar(255) not null,
  password_hash varchar(255) not null
);