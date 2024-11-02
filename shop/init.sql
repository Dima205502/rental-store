CREATE TABLE things (
    id SERIAL PRIMARY KEY, 
    owner TEXT NOT NULL,
    type TEXT NOT NULL,
    description TEXT,
    price INT NOT NULL CHECK(price >= 0),
    available BOOLEAN DEFAULT TRUE
);

CREATE TABLE taken_things(
    thing_id INT REFERENCES things(id),
    buyer TEXT NOT NULL,
    finish_time TIMESTAMP NOT NULL
);