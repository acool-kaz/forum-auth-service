CREATE TABLE IF NOT EXISTS users (
    id serial not null unique, 
    first_name varchar(255) not null, 
    last_name varchar(255) not null, 
    email varchar(255) not null unique, 
    username varchar(255) not null unique, 
    password varchar(255) not null
);