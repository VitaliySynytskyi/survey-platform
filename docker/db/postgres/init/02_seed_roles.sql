-- Seed initial roles
INSERT INTO roles (name, created_at, updated_at) VALUES
('admin', NOW(), NOW()),
('user', NOW(), NOW())
ON CONFLICT (name) DO NOTHING; 