CREATE EXTENSION IF NOT EXISTS pg_cron;

CREATE OR REPLACE FUNCTION update_at_updated()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

DROP TABLE IF EXISTS users_roles;
DROP TABLE IF EXISTS user_groups;
DROP TABLE IF EXISTS roles_permissions;
DROP TABLE IF EXISTS group_roles;
DROP TABLE IF EXISTS groups;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS api_keys;
DROP TABLE IF EXISTS users;


DROP TABLE IF EXISTS group_roles;
DROP TABLE IF EXISTS user_groups;
DROP TABLE IF EXISTS groups;
CREATE TABLE groups (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    name TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL);
CREATE INDEX groups_name_idx ON groups (name);
CREATE TRIGGER groups_updated_trigger
    BEFORE UPDATE
    ON
        groups
    FOR EACH ROW
EXECUTE PROCEDURE update_at_updated();

INSERT INTO groups (name, description) VALUES('superusers', 'Superusers');
INSERT INTO groups (name, description) VALUES('admins', 'Administrators');
-- INSERT INTO groups (name, description) VALUES('users', 'Standard users');
INSERT INTO groups (name, description) VALUES('login', 'Users who can login');
INSERT INTO groups (name, description) VALUES('rdf', 'For viewers of RDF data');

DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS resources;
CREATE TABLE resources (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    name TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL);
CREATE INDEX resources_name_idx ON resources (name);
CREATE TRIGGER resources_updated_trigger
    BEFORE UPDATE
    ON
        resources
    FOR EACH ROW
EXECUTE PROCEDURE update_at_updated();

INSERT INTO resources (name, description) VALUES('*', 'All resources');
INSERT INTO resources (name, description) VALUES('web', 'Web access');
INSERT INTO resources (name, description) VALUES('rdf', 'RDF access');
 

DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS actions;
CREATE TABLE actions (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    name TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL);
CREATE INDEX actions_name_idx ON actions (name);
CREATE TRIGGER actions_updated_trigger
    BEFORE UPDATE
    ON
        actions
    FOR EACH ROW
EXECUTE PROCEDURE update_at_updated();

INSERT INTO actions (name, description) VALUES('*', 'All actions');
INSERT INTO actions (name, description) VALUES('read', 'Read access');
INSERT INTO actions (name, description) VALUES('write', 'Write access');
INSERT INTO actions (name, description) VALUES('delete', 'Delete access');
INSERT INTO actions (name, description) VALUES('login', 'Login access');
INSERT INTO actions (name, description) VALUES('view', 'View resources');




DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS permissions;
CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    name TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    resource_id UUID NOT NULL,
    action_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(resource_id, action_id),
    FOREIGN KEY(resource_id) REFERENCES resources(id) ON DELETE CASCADE,
    FOREIGN KEY(action_id) REFERENCES actions(id) ON DELETE CASCADE);
CREATE INDEX permissions_resource_action_idx ON permissions (resource_id, action_id);
CREATE TRIGGER permissions_updated_trigger
    BEFORE UPDATE
    ON
        permissions
    FOR EACH ROW
EXECUTE PROCEDURE update_at_updated();

INSERT INTO permissions (resource_id, action_id, name) 
SELECT r.id, a.id, '*'
FROM resources r, actions a
WHERE r.name = '*' AND a.name = '*';

INSERT INTO permissions (resource_id, action_id, name) 
SELECT r.id, a.id, 'login'
FROM resources r, actions a
WHERE r.name = 'web' AND a.name = 'login';

INSERT INTO permissions (resource_id, action_id, name) 
SELECT r.id, a.id, 'rdf-viewer'
FROM resources r, actions a
WHERE r.name = 'rdf' AND a.name = 'view';


DROP TABLE IF EXISTS group_roles;
DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS roles;
CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    name TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL);
CREATE INDEX roles_name_idx ON roles (name);

CREATE TRIGGER roles_updated_trigger
    BEFORE UPDATE
    ON
        roles
    FOR EACH ROW
EXECUTE PROCEDURE update_at_updated();


INSERT INTO roles (name) VALUES('root');
INSERT INTO roles (name) VALUES('admin');
INSERT INTO roles (name) VALUES('login');
INSERT INTO roles (name) VALUES('rdf-viewer');


DROP TABLE IF EXISTS role_permissions;
CREATE TABLE role_permissions (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    role_id UUID NOT NULL,
    permission_id UUID NOT NULL,
    name TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(role_id, permission_id),
    FOREIGN KEY(role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY(permission_id) REFERENCES permissions(id) ON DELETE CASCADE);
CREATE INDEX role_permissions_role_permission_idx ON role_permissions (role_id, permission_id);
CREATE TRIGGER role_permissions_updated_trigger
    BEFORE UPDATE
    ON
        role_permissions
    FOR EACH ROW
EXECUTE PROCEDURE update_at_updated();

-- super/user admin
INSERT INTO role_permissions (role_id, permission_id, name) 
SELECT r.id, p.id, 'Superuser all access'
FROM roles r, permissions p
WHERE r.name = 'root' AND p.name = '*';

-- super/user can login
INSERT INTO role_permissions (role_id, permission_id, name) 
SELECT r.id, p.id, 'Admin login access'
FROM roles r, permissions p
WHERE r.name = 'admin' AND p.name = '*';

-- super/user can login
INSERT INTO role_permissions (role_id, permission_id, name) 
SELECT r.id, p.id, 'User can login'
FROM roles r, permissions p
WHERE r.name = 'login' AND p.name = 'login';


INSERT INTO role_permissions (role_id, permission_id, name) 
SELECT r.id, p.id, 'rdf-viewer'
FROM roles r, permissions p
WHERE r.name = 'rdf-viewer' AND p.name = 'rdf-viewer';

 

DROP TABLE IF EXISTS group_roles;
CREATE TABLE group_roles (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    group_id UUID NOT NULL,
    role_id UUID NOT NULL,
    name TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(group_id, role_id),
    FOREIGN KEY(group_id) REFERENCES groups(id) ON DELETE CASCADE,
    FOREIGN KEY(role_id) REFERENCES roles(id) ON DELETE CASCADE);
CREATE INDEX group_roles_group_role_idx ON group_roles (group_id, role_id);
CREATE TRIGGER group_roles_updated_trigger
    BEFORE UPDATE
    ON
        group_roles
    FOR EACH ROW
EXECUTE PROCEDURE update_at_updated();

INSERT INTO group_roles (group_id, role_id, name)
SELECT g.id, r.id, 'superusers'
FROM groups g, roles r
WHERE g.name = 'superusers' AND r.name = 'root' ON CONFLICT DO NOTHING;

INSERT INTO group_roles (group_id, role_id, name)
SELECT g.id, r.id, 'admins'
FROM groups g, roles r
WHERE g.name = 'admins' AND r.name = 'admin' ON CONFLICT DO NOTHING;

INSERT INTO group_roles (group_id, role_id, name)
SELECT g.id, r.id, 'login'
FROM groups g, roles r
WHERE g.name = 'login' AND r.name = 'login' ON CONFLICT DO NOTHING;

INSERT INTO group_roles (group_id, role_id, name)
SELECT g.id, r.id, 'rdf-viewer'
FROM groups g, roles r
WHERE g.name = 'rdf' AND r.name = 'rdf-viewer' ON CONFLICT DO NOTHING;


 

-- fix original users table
DROP TABLE IF EXISTS api_keys;
ALTER TABLE IF EXISTS users DROP column id;
ALTER TABLE IF EXISTS users ADD COLUMN id UUID PRIMARY KEY DEFAULT uuidv7();
ALTER TABLE IF EXISTS users DROP column public_id;



DROP TABLE IF EXISTS users;
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL DEFAULT '',
    name TEXT NOT NULL DEFAULT '',
    first_name TEXT NOT NULL DEFAULT '',
    last_name TEXT NOT NULL DEFAULT '',
    is_locked BOOLEAN NOT NULL DEFAULT false,
    email_verified_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL);
-- CREATE INDEX name ON users (first_name, last_name);
CREATE INDEX users_username_idx ON users (username);
CREATE INDEX users_email_idx ON users (email);
CREATE TRIGGER users_updated_trigger
    BEFORE UPDATE
    ON
        users
    FOR EACH ROW
EXECUTE PROCEDURE update_at_updated();

-- the superuser me --
INSERT INTO users (username, email, email_verified_at) VALUES (
    'root',
    'edb-root@antonyholmes.dev',
    now()
);

DROP TABLE IF EXISTS auth_providers;
CREATE TABLE auth_providers (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(name));
-- CREATE INDEX name ON users (first_name, last_name);
CREATE INDEX auth_providers_name_idx ON auth_providers (name);
CREATE TRIGGER auth_providers_updated_trigger
    BEFORE UPDATE
    ON
        auth_providers
    FOR EACH ROW
EXECUTE PROCEDURE update_at_updated();

INSERT INTO auth_providers (name) VALUES ('edb');
INSERT INTO auth_providers (name) VALUES ('auth0');
INSERT INTO auth_providers (name) VALUES ('google');
INSERT INTO auth_providers (name) VALUES ('github');

DROP TABLE IF EXISTS user_auth_providers;
CREATE TABLE user_auth_providers (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    user_id UUID NOT NULL,
    auth_provider_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(user_id, auth_provider_id),
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(auth_provider_id) REFERENCES auth_providers(id) ON DELETE CASCADE);
CREATE TRIGGER users_updated_trigger
    BEFORE UPDATE
    ON
        user_auth_providers
    FOR EACH ROW
EXECUTE PROCEDURE update_at_updated();

DROP TABLE IF EXISTS user_groups;
CREATE TABLE user_groups (
    id UUID PRIMARY KEY DEFAULT uuidv7(), 
    user_id UUID NOT NULL,
    group_id UUID NOT NULL,
    name TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(user_id, group_id),
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(group_id) REFERENCES groups(id) ON DELETE CASCADE);
CREATE INDEX user_groups_user_group_idx ON user_groups (user_id, group_id);
CREATE INDEX user_groups_group_user_idx ON user_groups (group_id, user_id);
CREATE TRIGGER user_groups_updated_trigger
    BEFORE UPDATE
    ON
        user_groups
    FOR EACH ROW
EXECUTE PROCEDURE update_at_updated();


INSERT INTO user_groups (user_id, group_id, name) 
SELECT u.id, g.id, 'Add to superuser group'
FROM users u, groups g
WHERE u.username = 'root' AND g.name = 'superusers' ON CONFLICT DO NOTHING;

INSERT INTO user_groups (user_id, group_id, name) 
SELECT u.id, g.id, 'Add to superuser group'
FROM users u, groups g
WHERE u.username LIKE '%antony%' AND g.name = 'superusers' ON CONFLICT DO NOTHING;


INSERT INTO user_groups (user_id, group_id, name)
SELECT u.id, g.id, 'Add to login group'
FROM users u, groups g
WHERE g.name = 'login' ON CONFLICT DO NOTHING;

INSERT INTO user_groups (user_id, group_id, name)
SELECT u.id, g.id, 'Add to rdf group'
FROM users u, groups g
WHERE u.email LIKE '%columbia%' AND g.name = 'rdf' ON CONFLICT DO NOTHING;

CREATE TABLE IF NOT EXISTS public_keys (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    user_id UUID NOT NULL,
    key TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE);
CREATE TRIGGER public_keys_updated_trigger
    BEFORE UPDATE
    ON
        public_keys
    FOR EACH ROW
EXECUTE PROCEDURE update_at_updated();

-- VALUES(1, 1, 'Superuser all access');
-- INSERT INTO user_groups (user_id, group_id, name) VALUES(1, 4, 'Superuser can login');

-- INSERT INTO user_groups (user_id, group_id, name) VALUES(2, 2, 'Admin all access');
-- INSERT INTO user_groups (user_id, group_id, name) VALUES(2, 4, 'Admin can login');
-- INSERT INTO user_groups (user_id, group_id, name) VALUES(3, 2, 'Admin all access');
-- INSERT INTO user_groups (user_id, group_id, name) VALUES(3, 4, 'Admin can login');

-- insert INTO user_groups (user_id, group_id, name) SELECT id, 4, 'Add to Login group' FROM users ON CONFLICT DO NOTHING;
-- insert INTO user_groups (user_id, group_id, name) SELECT id, 5, 'Add to RDF group' FROM users ON CONFLICT DO NOTHING;


-- CREATE TABLE user_roles (
--     id SERIAL PRIMARY KEY, 
--     user_id INTEGER NOT NULL,
--     role_id INTEGER NOT NULL, 
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
--     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
--     UNIQUE(user_id, role_id),
--     FOREIGN KEY(user_id) REFERENCES users(id),
--     FOREIGN KEY(role_id) REFERENCES roles(id));
-- CREATE INDEX users_roles_user_id_idx ON users_roles (user_id, role_id);
-- CREATE TRIGGER users_roles_updated_trigger
--     BEFORE UPDATE
--     ON
--         users_roles
--     FOR EACH ROW
-- EXECUTE PROCEDURE update_at_updated();

 

DROP TABLE IF EXISTS api_keys; 
CREATE TABLE api_keys (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    user_id UUID NOT NULL,
    api_key TEXT NOT NULL, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(user_id, api_key),
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE);
CREATE INDEX api_keys_api_key_idx ON api_keys (api_key);
CREATE TRIGGER api_keys_updated_trigger
    BEFORE UPDATE
    ON
        api_keys
    FOR EACH ROW
EXECUTE PROCEDURE update_at_updated();


