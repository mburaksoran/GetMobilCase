CREATE TABLE IF NOT EXISTS products (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    stock_count INT NOT NULL,
    price FLOAT NOT NULL
    );

INSERT INTO products (name, stock_count, price)
VALUES
('product-name-1', 3, 10),
('product-name-2', 5, 20),
('product-name-3', 2, 5);
