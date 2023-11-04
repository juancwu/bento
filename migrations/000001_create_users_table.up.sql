CREATE TABLE IF NOT EXISTS users (
    id INTEGER NOT NULL AUTO_INCREMENT,
    email VARCHAR(255) NOT NULL,
    gh_id INTEGER NOT NULL,
    object_id VARCHAR(12) NOT NULL,
    created_at DATETIME(6) NOT NULL,
    PRIMARY KEY (id),
    UNIQUE KEY idx_gh_id (gh_id),
    UNIQUE KEY idx_object_id (object_id)
);
