CREATE TABLE users (
  user_id varchar(32) NOT NULL,
  user_name varchar(100) NOT NULL,
  created_at timestamp with time zone,
  CONSTRAINT pk_users PRIMARY KEY (user_id)
);

-- name: GetUser :one
SELECT * FROM users
WHERE user_id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY user_name;

-- name: CreateUser :one
INSERT INTO users (
  user_id, user_name, created_at
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE user_id = $1;

-- name: UpdateUser :exec
UPDATE users
set user_name = $2
WHERE user_id = $1;


CREATE TABLE products (
  product_no integer NOT NULL,
  name varchar(100) NOT NULL,
  price integer NOT NULL,
  CONSTRAINT pk_products PRIMARY KEY (product_no)
);

-- name: GetProduct :one
SELECT * FROM products
WHERE product_no = $1 LIMIT 1;

-- name: ListProducts :many
SELECT * FROM products
ORDER BY name;

-- name: CreateProduct :one
INSERT INTO products (
  product_no, name, price
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE product_no = $1;

-- name: UpdateProduct :exec
UPDATE products
set name = $2,
    price = $3
WHERE product_no = $1;
