CREATE TABLE users (
	id INTEGER PRIMARY KEY, 
	username VARCHAR(50), 
	password VARCHAR(50)
);

CREATE TABLE items (
	id INTEGER PRIMARY KEY, 
	title VARCHAR(50), 
	url VARCHAR(80),
	uid INTEGER,
	FOREIGN KEY(uid) REFERENCES users(id)
);

CREATE TABLE sessions (
	id INTEGER PRIMARY KEY,
	sid VARCHAR(50),
	valid INTEGER
);