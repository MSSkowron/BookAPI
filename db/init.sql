CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email varchar(50),
		password varchar(256),
		first_name varchar(50),
		last_name varchar(50),
		age smallint
)