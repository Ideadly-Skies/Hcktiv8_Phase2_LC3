-- Drop tables if they already exist
DROP TABLE IF EXISTS OrderItems CASCADE;
DROP TABLE IF EXISTS Orders CASCADE;
DROP TABLE IF EXISTS Carts CASCADE;
DROP TABLE IF EXISTS Products CASCADE;
DROP TABLE IF EXISTS Users CASCADE;

-- Create Users table
CREATE TABLE Users (
    user_id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    jwt_token VARCHAR(255)
);

-- Create Products table
CREATE TABLE Products (
    product_id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    description TEXT,
    price DECIMAL(10, 2)
);

-- Create Carts table, which contains user_id and product_id as foreign keys
CREATE TABLE Carts (
    cart_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES Users(user_id),
    product_id INTEGER REFERENCES Products(product_id),
    quantity INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create Orders table, with a foreign key reference to Users
CREATE TABLE Orders (
    order_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES Users(user_id),
    total_price DECIMAL(10,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create OrderItems table to record individual items in each order, with foreign key references to Orders and Products
CREATE TABLE OrderItems (
    order_item_id SERIAL PRIMARY KEY,
    order_id INTEGER REFERENCES Orders(order_id),
    product_id INTEGER REFERENCES Products(product_id),
    quantity INTEGER,
    price DECIMAL(10,2)
);

-- Insert sample data into Users table
INSERT INTO Users (name, email, password, jwt_token) 
VALUES 
('Alice Johnson', 'alice.johnson@example.com', 'hashed_password1', 'jwt_token1'),
('Bob Smith', 'bob.smith@example.com', 'hashed_password2', 'jwt_token2');

-- Insert sample data into Products table
INSERT INTO Products (name, description, price) 
VALUES 
('Product1', 'Product1 Description', 100.00),
('Product2', 'Product2 Description', 200.00),
('Product3', 'Product3 Description', 300.00);

-- Insert sample data into Carts table
INSERT INTO Carts (user_id, product_id, quantity, created_at) 
VALUES 
(1, 1, 2, '2023-09-09 10:00:00'),
(1, 2, 1, '2023-09-09 10:05:00'),
(2, 3, 3, '2023-09-09 10:10:00');

-- Insert sample data into Orders table
INSERT INTO Orders (user_id, total_price, created_at) 
VALUES 
(1, 300.00, '2023-09-10 11:00:00'),
(2, 900.00, '2023-09-10 11:05:00');

-- Insert sample data into OrderItems table
INSERT INTO OrderItems (order_id, product_id, quantity, price) 
VALUES 
(1, 1, 2, 100.00),
(1, 2, 1, 200.00),
(2, 3, 3, 300.00);