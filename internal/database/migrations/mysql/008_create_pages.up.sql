-- Create pages table
CREATE TABLE pages (
    id VARCHAR(36) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    content TEXT NOT NULL,
    author_id VARCHAR(36) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    `order` INT NOT NULL DEFAULT 0,
    deleted_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_pages_author FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT chk_pages_status CHECK (status IN ('draft', 'published', 'archived'))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create indexes for common queries
CREATE INDEX idx_pages_slug ON pages(slug);
CREATE INDEX idx_pages_author_id ON pages(author_id);
CREATE INDEX idx_pages_status ON pages(status);
CREATE INDEX idx_pages_order ON pages(`order`);
CREATE INDEX idx_pages_deleted_at ON pages(deleted_at);
