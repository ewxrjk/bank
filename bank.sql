CREATE TABLE users (
  user TEXT PRIMARY KEY,
  password TEXT
);

CREATE TABLE sessions (
  user TEXT,
  id TEXT PRIMARY KEY,
  expires INTEGER
);

CREATE TABLE transactions (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  time INTEGER,
  user TEXT,
  description TEXT,
  origin TEXT,
  destination TEXT,
  amount INTEGER,
  origin_balance_after INTEGER,
  destination_balance_after INTEGER
);

CREATE TABLE accounts (
  user TEXT PRIMARY KEY,
  balance INTEGER
);

INSERT INTO users (user,password) VALUES("richard","");
INSERT INTO accounts (user,balance) VALUES("richard",0);
