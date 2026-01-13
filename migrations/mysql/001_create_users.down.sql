-- Drop users table
DROP INDEX idx_users_role ON users;
DROP INDEX idx_users_username ON users;
DROP INDEX idx_users_email ON users;
DROP TABLE IF EXISTS users;
