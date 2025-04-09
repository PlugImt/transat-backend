-- Add friendship table to manage connections between users
CREATE TABLE realcampus_friendships
(
    id_friendship   SERIAL,
    user_id         VARCHAR(100) NOT NULL,                   -- Reference to user email in newf table
    friend_id       VARCHAR(100) NOT NULL,                   -- Reference to friend's email in newf table
    status          VARCHAR(20)  NOT NULL DEFAULT 'PENDING', -- PENDING, ACCEPTED, REJECTED, BLOCKED
    request_date    TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    acceptance_date TIMESTAMP,
    PRIMARY KEY (id_friendship),
    UNIQUE (user_id, friend_id),
    FOREIGN KEY (user_id) REFERENCES newf (email) ON DELETE CASCADE,
    FOREIGN KEY (friend_id) REFERENCES newf (email) ON DELETE CASCADE,
    CHECK (status IN ('PENDING', 'ACCEPTED', 'REJECTED')),
    CHECK (user_id <> friend_id)                             -- Prevent self-friendship
);

-- Create an index to speed up friendship lookups
CREATE INDEX idx_realcampus_friendships_users ON realcampus_friendships (user_id, friend_id);

-- First, create a privacy level enum
CREATE TYPE realcampus_privacy_level AS ENUM ('PUBLIC', 'PRIVATE');

-- Create a posts table that extends the files functionality
CREATE TABLE realcampus_posts
(
    id_post       SERIAL,
    id_file_1     INTEGER                  NOT NULL,
    id_file_2     INTEGER                  NOT NULL,
    author_email  VARCHAR(100)             NOT NULL,
    location      VARCHAR(200),
    privacy       realcampus_privacy_level NOT NULL DEFAULT 'PRIVATE',
    creation_date TIMESTAMP                NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id_post),
    FOREIGN KEY (id_file_1) REFERENCES files (id_files) ON DELETE CASCADE,
    FOREIGN KEY (id_file_2) REFERENCES files (id_files) ON DELETE CASCADE,
    FOREIGN KEY (author_email) REFERENCES newf (email) ON DELETE CASCADE
);

-- Add index for post lookups
CREATE INDEX idx_realcampus_posts_files_1 ON realcampus_posts (id_file_1);
CREATE INDEX idx_realcampus_posts_files_2 ON realcampus_posts (id_file_2);

-- Create table for post reactions (using images)
CREATE TABLE realcampus_post_reactions
(
    id_reaction   SERIAL,
    id_post       INTEGER      NOT NULL,
    id_file       INTEGER      NOT NULL, -- The reaction image
    author_email  VARCHAR(100) NOT NULL, -- Who created this reaction
    creation_date TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id_reaction),
    UNIQUE (id_post, author_email),      -- One reaction per user per post
    FOREIGN KEY (id_post) REFERENCES realcampus_posts (id_post) ON DELETE CASCADE,
    FOREIGN KEY (id_file) REFERENCES files (id_files) ON DELETE CASCADE,
    FOREIGN KEY (author_email) REFERENCES newf (email) ON DELETE CASCADE
);

-- Add index for reaction lookups
CREATE INDEX idx_realcampus_post_reactions_post ON realcampus_post_reactions (id_post);