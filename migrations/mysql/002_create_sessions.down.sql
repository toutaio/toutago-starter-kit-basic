-- Drop sessions table
DROP INDEX idx_sessions_expires_at ON sessions;
DROP INDEX idx_sessions_user_id ON sessions;
DROP TABLE IF EXISTS sessions;
