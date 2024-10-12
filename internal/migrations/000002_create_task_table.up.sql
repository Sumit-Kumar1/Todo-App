CREATE TABLE IF NOT EXISTS users(
    user_id TEXT NOT NULL PRIMARY KEY, 
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE CHECK (email LIKE '%'), 
    password TEXT NOT NULL);