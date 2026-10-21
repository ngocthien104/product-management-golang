CREATE DATABASE IF NOT EXISTS mini_product_management
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci;

USE mini_product_management;

CREATE TABLE IF NOT EXISTS products (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  slug VARCHAR(120) NOT NULL UNIQUE,
  description TEXT,
  price DECIMAL(10, 2) NOT NULL,
  category VARCHAR(30) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_products_category (category),
  INDEX idx_products_name (name)
);

INSERT INTO products (name, slug, description, price, category)
VALUES
('Wireless Mouse', 'wireless-mouse', 'A simple wireless mouse', 19.99, 'electronics'),
('Go Programming Book', 'go-programming-book', 'Beginner-friendly Go book', 29.50, 'books'),
('Cotton T-Shirt', 'cotton-t-shirt', 'Everyday cotton t-shirt', 12.00, 'clothing')
ON DUPLICATE KEY UPDATE name = VALUES(name);
