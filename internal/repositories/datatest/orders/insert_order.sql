INSERT INTO users(id, name, email, password, role, status)
VALUES (1000, 'Quang Thuy', 'qthuy1000@gmail.com', 'password', 'users', 'activated');

INSERT INTO products(id, name, description, price, quantity, author_id, category_id)
VALUES (1000, 'iPhone 14 1000', 'abc', 1200, 20, 1000, 1),
(1001, 'iPhone 14 1001', 'abc', 1200, 20, 1000, 1);

INSERT INTO orders(id, user_id, status, total_price)
VALUES (1000, 1000, 'Created', 1200),
(1001, 1000, 'Created', 1200);

INSERT INTO order_items(id, order_id, product_id, price, quantity)
VALUES (1000, 1000, 1000, 1200, 1),
(1001, 1000, 1001, 1200, 1),
(1002, 1001, 1000, 1200, 1),
(1003, 1001, 1001, 1200, 1);
