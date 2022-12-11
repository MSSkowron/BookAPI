CREATE TABLE IF NOT EXISTS users (
		id INT GENERATED ALWAYS AS IDENTITY,
		email varchar(50) NOT NULL,
		password varchar(256) NOT NULL,
		first_name varchar(50) NOT NULL,
		last_name varchar(50) NOT NULL,
		age smallint NOT NULL,
		PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS authors (
		id INT GENERATED ALWAYS AS IDENTITY,
		first_name varchar(50) NOT NULL,
		last_name varchar(50) NOT NULL,
		PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS books (
		id INT GENERATED ALWAYS AS IDENTITY,
		author_id  INT,
		isbn varchar(100) NOT NULL,
		title varchar(100) NOT NULL,
		CONSTRAINT fk_author 
			FOREIGN KEY(author_id)
				REFERENCES authors(id)
);

INSERT INTO authors (first_name, last_name) values ('mateusz', 'skowron');
INSERT INTO books (author_id, isbn, title) values (1, '123123123123123', 'Franklin pisze w Go');