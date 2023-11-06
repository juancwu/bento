CREATE TABLE gh_access_tokens (
    id INTEGER NOT NULL AUTO_INCREMENT,
    token VARCHAR(255) NOT NULL,
    user_id INTEGER NOT NULL,

    PRIMARY KEY(id)
);
