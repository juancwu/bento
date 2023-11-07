CREATE TABLE oauth_states (
    id INTEGER NOT NULL AUTO_INCREMENT,
    state_id VARCHAR(12) NOT NULL,
    flow VARCHAR(3) NOT NULL,
    redirect TEXT NOT NULL,
    port SMALLINT NOT NULL,
    expires_at DATETIME NOT NULL,

    PRIMARY KEY(id),
    UNIQUE KEY idx_state_id (state_id)
);
