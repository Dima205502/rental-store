CREATE TABLE Users (
    id SERIAL PRIMARY KEY,
    nickname TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL, 
    verify BOOL DEFAULT FALSE NOT NULL
);

CREATE TABLE Email_token(
    user_id INT REFERENCES Users(id),
    token TEXT NOT NULL
);

CREATE TABLE Sessions(
    user_id INT REFERENCES Users(id),
    token TEXT NOT NULL
);