-- Create post_versions table for version history
CREATE TABLE post_versions (
    id SERIAL PRIMARY KEY,
    post_id INTEGER NOT NULL,
    version INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    excerpt TEXT,
    author_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_post_versions_post FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    CONSTRAINT fk_post_versions_author FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE (post_id, version)
);

-- Create indexes
CREATE INDEX idx_post_versions_post_id ON post_versions(post_id);
CREATE INDEX idx_post_versions_version ON post_versions(post_id, version DESC);
