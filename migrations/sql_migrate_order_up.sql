CREATE TABLE IF NOT EXISTS orders (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id int(255) NOT NULL,
    product_id INT NOT NULL,
    ordered_count INT NOT NULL,
    price FLOAT NOT NULL
    );

INSERT INTO orders (user_id, product_id, ordered_count, price)
VALUES
(5, 3, 1, 10),
(6, 5, 1, 20),
(9, 2, 1, 75);


