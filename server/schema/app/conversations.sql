CREATE TABLE IF NOT EXISTS conversations (
    id BIGINT PRIMARY KEY,
    type VARCHAR(10) CHECK (type IN ('DM', 'GROUP')),
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,   -- update on new message
    group_name TEXT,              -- for groups only
);

CREATE TABLE IF NOT EXISTS conversation_participants (
    conversation_id BIGINT NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL,
    joined_at BIGINT NOT NULL,
    role VARCHAR(10) DEFAULT 'MEMBER' CHECK (role IN ('MEMBER', 'ADMIN')) NOT NULL, -- group roles
    PRIMARY KEY (conversation_id, user_id)
);

CREATE TABLE IF NOT EXISTS direct_conversations (
    user1 BIGINT NOT NULL,
    user2 BIGINT NOT NULL,
    conversation_id BIGINT NOT NULL UNIQUE REFERENCES conversations(id) ON DELETE CASCADE,
    created_at BIGINT NOT NULL,
    PRIMARY KEY (user1, user2),
    CHECK (user1 < user2)
);