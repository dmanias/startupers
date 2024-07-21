-- Creating the users table
CREATE TABLE IF NOT EXISTS users
(
    id            UUID PRIMARY KEY,
    name          VARCHAR(255) NOT NULL,
    email         TEXT         NULL UNIQUE, -- Assuming email will be stored as a simple string
    roles         TEXT[],                   -- Postgres supports array types, which can be used for storing roles
    password_hash BYTEA        NOT NULL,    -- BYTEA for binary data
    enabled       BOOLEAN      NOT NULL,
    date_created  TIMESTAMPTZ  NOT NULL,
    date_updated  TIMESTAMPTZ  NOT NULL
);

-- Creating the ideas queries table
CREATE TABLE IF NOT EXISTS ideas
(
    id            UUID PRIMARY KEY,
    user_id       UUID         NOT NULL,
    title         VARCHAR(255) NOT NULL,
    description   TEXT,
    category      VARCHAR(255),
    tags          TEXT[],
    privacy       VARCHAR(20)  NOT NULL    DEFAULT 'private',
    collaborators UUID[],
    avatar_url    TEXT,
    stage         VARCHAR(50),
    inspiration   TEXT,
    date_created  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    date_updated  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users (id)
);

-- Creating the moderators table
CREATE TABLE IF NOT EXISTS moderators
(
    id           UUID PRIMARY KEY,
    name         VARCHAR(255) NOT NULL,
    instruction  TEXT         NOT NULL,
    date_created TIMESTAMPTZ  NOT NULL,
    date_updated TIMESTAMPTZ  NOT NULL
);
--Ideas' posts
CREATE TABLE IF NOT EXISTS posts
(
    id           UUID PRIMARY KEY,
    idea_id      UUID NOT NULL,
    author_id    UUID NOT NULL,
    content      TEXT NOT NULL,
    owner_type   VARCHAR(20) NOT NULL DEFAULT 'idea',
    date_created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    date_updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY (idea_id) REFERENCES ideas (id),
    FOREIGN KEY (author_id) REFERENCES users (id)
);

---- Challenges for ideas
CREATE TABLE IF NOT EXISTS challenges
(
    id           UUID PRIMARY KEY,
    idea_id      UUID         NOT NULL,
    moderator_id UUID         NOT NULL,
    answer       TEXT,
    photo_url    TEXT,
    date_created TIMESTAMPTZ  NOT NULL,
    date_updated TIMESTAMPTZ  NOT NULL,
    FOREIGN KEY (idea_id) REFERENCES ideas (id),
    FOREIGN KEY (moderator_id) REFERENCES moderators (id)
);

-- Creating the Ai queries table
CREATE TABLE IF NOT EXISTS ais
(
    id           UUID PRIMARY KEY,           -- Using SERIAL for auto-incrementing integer ID
    name         VARCHAR(255) NOT NULL,
    query        TEXT         NOT NULL,      -- Storing the actual AI query
--    userid UUID NOT NULL REFERENCES users(id), -- Link to users table
    userid       UUID REFERENCES users (id), -- Link to users table
    date_created TIMESTAMPTZ  NOT NULL,
    date_updated TIMESTAMPTZ  NOT NULL
);

-- Creating indexes (assuming you still want them)
CREATE INDEX IF NOT EXISTS idx_users_id ON users (id);
CREATE INDEX IF NOT EXISTS idx_moderators_id ON moderators (id);
CREATE INDEX IF NOT EXISTS idx_ais_userid ON ais (userid); -- Index on userid for quick lookup

-- Creating a view for the users table
-- Note: Adjusting the view to include all relevant fields you might need.
CREATE OR REPLACE VIEW users_view AS
SELECT id, name, email, enabled, date_created, date_updated
FROM users;
