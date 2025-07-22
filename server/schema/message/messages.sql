-- wip

CREATE TABLE messages (
    id BIGINT PRIMARY KEY,
    conversation_id BIGINT NOT NULL,
    sender_id BIGINT NOT NULL,
    content TEXT,
    created_at BIGINT NOT NULL,
    edited_at BIGINT,
    deleted BOOLEAN DEFAULT false NOT NULL
    -- todo: add message attachments
);