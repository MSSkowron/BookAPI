create table users
(
    id bigint primary key generated always as identity,
    email    varchar(50) NOT NULL,
	password varchar(256) NOT NULL,
    first_name varchar(50) NOT NULL,
	last_name varchar(50) NOT NULL,
	age smallint NOT NULL
);

create table books (
    id bigint primary key generated always as identity,
    author varchar(100) NOT NULL, 
    title varchar(100) NOT NULL
);