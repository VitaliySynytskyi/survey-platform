-- Seed an admin user with username 'admin' and password 'admin'
-- The password hash is for the string 'admin' generated with bcrypt cost factor 10

-- Insert admin user if not exists
INSERT INTO users (username, email, password_hash, first_name, last_name, is_active, created_at, updated_at)
VALUES 
    ('admin', 'admin@survey-app.com', '$2a$10$2WE.jwUKvHsmarx.OWEK8.2ZTKz.uSTPL8YkzOcKw4tCWgW33N4f6', 'Admin', 'User', true, NOW(), NOW())
ON CONFLICT (username) DO NOTHING
RETURNING id;

-- Get admin user ID
DO $$
DECLARE
    admin_user_id INTEGER;
    admin_role_id INTEGER;
BEGIN
    -- Get admin user ID
    SELECT id INTO admin_user_id FROM users WHERE username = 'admin';
    
    -- Get admin role ID
    SELECT id INTO admin_role_id FROM roles WHERE name = 'admin';
    
    -- Remove any existing roles for admin user (to avoid having both 'user' and 'admin' roles)
    DELETE FROM user_roles WHERE user_id = admin_user_id;
    
    -- Assign admin role to admin user
    IF admin_user_id IS NOT NULL AND admin_role_id IS NOT NULL THEN
        INSERT INTO user_roles (user_id, role_id, created_at)
        VALUES (admin_user_id, admin_role_id, NOW())
        ON CONFLICT (user_id, role_id) DO NOTHING;
    END IF;
END $$; 