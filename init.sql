-- Create the 'users' table
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(32) NOT NULL,
    password VARCHAR(32) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO users (username, password) VALUES
    ('user1', 'password1'),
    ('user2', 'password2'),
    ('user3', 'password3');
