create table users
(
    id bigint primary key generated always as identity,
    created_at timestamptz default NOW() NOT NULL,
    email    varchar(50) unique NOT NULL,
	password varchar(256) NOT NULL,
    first_name varchar(50) NOT NULL,
	last_name varchar(50) NOT NULL,
	age smallint NOT NULL
);

create table books (
    id bigint primary key generated always as identity,
    created_at timestamptz default NOW() NOT NULL,
    author varchar(100) NOT NULL, 
    title varchar(100) NOT NULL
);