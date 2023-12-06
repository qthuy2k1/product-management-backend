INSERT INTO product_categories (id, name, description)
VALUES (2001, 'Smartphone', 'A mobile device'),
    (2002, 'Laptop', 'A portable computer');

INSERT INTO users
VALUES (1001, 'John Doe', 'doe@example.com', '123123', 'user', 'activated'),
(1002, 'Steve Job', 'job@example.com', '212212', 'user', 'activated');
  
INSERT INTO products (id, name, description, price, quantity, category_id, author_id)
VALUES (1001, 'Macbook Air M1 16GB', 'A macbook air with Apple M1 chip, 16GB of RAM, 512GB SSD', 1200, 10, 2002, 1001),
(1002, 'iPhone 13', 'An iPhone with Apple A15 chipset, 4GB RAM, 128GB storage', 800, 10, 2001, 1002);
