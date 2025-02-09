-- +goose Up
-- +goose StatementBegin
INSERT INTO users (id, balance) VALUES
(1, 100.00),
(2, 250.50),
(3, 50.75),
(4, 500.00),
(5, 0.00),
(6, 750.25),
(7, 1200.00),
(8, 30.00),
(9, 999.99),
(10, 10.10);

INSERT INTO transactions (user_id, amount, operation_type, description, related_user_id) VALUES 
(1, 200.00, 'deposit', 'Пополнение счета', NULL),  
(2, -50.00, 'transfer', 'Перевод пользователю 3', 3),  
(3, 50.00, 'transfer', 'Получен перевод от пользователя 2', 2),  
(4, -300.00, 'transfer', 'Перевод пользователю 5', 5),  
(5, 300.00, 'transfer', 'Получен перевод от пользователя 4', 4);  

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd
