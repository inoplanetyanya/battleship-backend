create table users (
  id serial not null unique,
  name varchar(255) not null,
  password_hash varchar(255) not null
);