-- Create page_versions table for version history
CREATE TABLE page_versions (
    id SERIAL PRIMARY KEY,
    page_id INTEGER NOT NULL,
    version INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    author_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_page_versions_page FOREIGN KEY (page_id) REFERENCES pages(id) ON DELETE CASCADE,
    CONSTRAINT fk_page_versions_author FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE (page_id, version)
);

-- Create indexes
CREATE INDEX idx_page_versions_page_id ON page_versions(page_id);
CREATE INDEX idx_page_versions_version ON page_versions(page_id, version DESC);
