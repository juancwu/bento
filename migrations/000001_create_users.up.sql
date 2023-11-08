CREATE TABLE users (
    id INT NOT NULL AUTO_INCREMENT,
    login TEXT NOT NULL,
    email TEXT NOT NULL,
    provider_id INT NOT NULL,
    provider_token VARCHAR(255),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,

    PRIMARY KEY (id),
    UNIQUE KEY idx_provider_id (provider_id)
);
