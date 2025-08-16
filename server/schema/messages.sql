-- wip

CREATE TABLE IF NOT EXISTS messages (
    id BIGINT NOT NULL,
    conversation_id BIGINT NOT NULL,
    author_id BIGINT NOT NULL,
    content TEXT,
    created_at BIGINT NOT NULL,
    edited_at BIGINT,
    deleted BOOLEAN DEFAULT false NOT NULL,
    PRIMARY KEY (id, conversation_id)
    -- todo: add message attachments
);

-- Distribute over Citus shards
SELECT create_distributed_table('messages', 'conversation_id');

-- For searching messages by created_at (pagination)
CREATE INDEX idx_messages_conv_created ON messages(conversation_id, created_at DESC);
-- For searching messages by user ID
CREATE INDEX idx_messages_conv_author ON messages(conversation_id, author_id);