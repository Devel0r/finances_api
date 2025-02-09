-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    balance DECIMAL(15,2) DEFAULT 0
);
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    operation_type VARCHAR(20) NOT NULL,
    description TEXT, 
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    related_user_id INT REFERENCES users(id)
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE transactions;
DROP TABLE users; 
-- +goose StatementEnd