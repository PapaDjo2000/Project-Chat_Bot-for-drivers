CREATE TABLE IF NOT EXISTS pr.reports (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES pr.users (chat_id) ON DELETE CASCADE,
    date TIMESTAMP NOT NULL DEFAULT NOW(),
    request JSONB,
    response JSONB
);