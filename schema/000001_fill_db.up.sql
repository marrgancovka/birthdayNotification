CREATE TABLE IF NOT EXISTS employees (
    id SERIAL PRIMARY KEY,
    name TEXT CONSTRAINT name_length CHECK (char_length(name) <= 50) NOT NULL,
    surname TEXT CONSTRAINT surname_length CHECK(char_length(surname) <= 50) NOT NULL,
    birth_date DATE NOT NULL
);

CREATE TABLE IF NOT EXISTS "users" (
    id SERIAL PRIMARY KEY,
    email TEXT CONSTRAINT email_length CHECK ( char_length(email) <= 50 ) NOT NULL UNIQUE,
    password TEXT CONSTRAINT password_length CHECK (char_length(password) <= 64) NOT NULL,
    id_employee INT NOT NULL,
    FOREIGN KEY (id_employee) REFERENCES employees(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS subscriptions (
    id_from INT NOT NULL,
    id_to INT NOT NULL,
    CONSTRAINT id_subscription PRIMARY KEY (id_from, id_to),
    FOREIGN KEY (id_from) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (id_to) REFERENCES users(id) ON DELETE CASCADE
)
