CREATE TABLE IF NOT EXISTS "users" (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_digest VARCHAR(255) NOT NULL
);

-- testpassword
INSERT INTO "users" (email, password_digest) VALUES ('admin@example.com', '$2a$14$e85TNgTUn8BWZIj.NdJxXeOv/AOm4lZ3pj27BcgO.Qt/orWDILw52');
