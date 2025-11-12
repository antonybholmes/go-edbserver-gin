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


CREATE TABLE groups (
    id SERIAL PRIMARY KEY, 
    public_id TEXT NOT NULL UNIQUE, 
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

INSERT INTO groups (public_id, name, description) VALUES('019a74b7-6459-70f0-9762-daa347d07f50', 'SuperUsers', 'Superusers');
INSERT INTO groups (public_id, name, description) VALUES('019a74b7-85c8-7330-9402-995f07a24fee', 'Admins', 'Administrators');
INSERT INTO groups (public_id, name, description) VALUES('019a74b7-ed63-7022-9489-9a6f64ac7f21', 'Users', 'Standard users');
INSERT INTO groups (public_id, name, description) VALUES('019a74e1-2669-7046-a1ed-ed28c1a1419f', 'LoginUsers', 'Users who can login');
INSERT INTO groups (public_id, name, description) VALUES('019a750c-c751-72b2-af19-e05fdb5ade15', 'RDFLabMembers', 'For viewers of RDF data');


CREATE TABLE permissions (
    id SERIAL PRIMARY KEY, 
    public_id TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL);
CREATE INDEX permissions_name_idx ON permissions (name);
CREATE TRIGGER permissions_updated_trigger
    BEFORE UPDATE
    ON
        permissions
    FOR EACH ROW
EXECUTE PROCEDURE update_at_updated();

INSERT INTO permissions (public_id, name, description) VALUES('01997350-f1db-734a-aabc-b738772a9d0c', '*', 'All permissions');
INSERT INTO permissions (public_id, name, description) VALUES('01997351-06c7-7f0d-b026-c51376a044ee', 'read:*', 'User has read access');
INSERT INTO permissions (public_id, name, description) VALUES('01997351-16e5-70f6-b869-ba08cdac4c85', 'write:*', 'User has write access');
INSERT INTO permissions (public_id, name, description) VALUES('01997351-2586-7e76-8a34-db50b222d47a', 'web:login', 'User can sign in');
INSERT INTO permissions (public_id, name, description) VALUES('019a7893-12e2-7d3a-ab13-89cc3cc43336', 'rdf:read:*', 'For viewers of RDF data');
 



CREATE TABLE roles (
    id SERIAL PRIMARY KEY, 
    public_id TEXT NOT NULL UNIQUE, 
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


INSERT INTO roles (public_id, name) VALUES('01997351-4c4b-7900-adb9-eeb64f772ed7', 'SuperAccess');
INSERT INTO roles (public_id, name) VALUES('01997351-7d4f-72bc-aba9-1fbd2d5d41a2', 'AdminAccess');
INSERT INTO roles (public_id, name) VALUES('01997351-8b67-758e-b798-f200f70c653b', 'ReadOnlyUser');
-- INSERT INTO roles (public_id, name) VALUES('UZuAVHDGToa4F786IPTijA==', 'GetDNA');
INSERT INTO roles (public_id, name) VALUES('01997351-9ba1-7a6e-9560-fc03b3098665', 'Login');
INSERT INTO roles (public_id, name) VALUES('01997351-aac6-7fcc-b886-b3f0585b90d8', 'RDFLabReadOnly');


CREATE TABLE group_roles (
    id SERIAL PRIMARY KEY, 
    group_id INTEGER NOT NULL,
    role_id INTEGER NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(group_id, role_id),
    FOREIGN KEY(group_id) REFERENCES groups(id) ON DELETE CASCADE,
    FOREIGN KEY(role_id) REFERENCES roles(id) ON DELETE CASCADE);
CREATE INDEX group_roles_group_id_idx ON group_roles (group_id, role_id);
CREATE TRIGGER group_roles_updated_trigger
    BEFORE UPDATE
    ON
        group_roles
    FOR EACH ROW
EXECUTE PROCEDURE update_at_updated();

-- super/user are both part of the admin group
INSERT INTO group_roles (group_id, role_id, description) VALUES(1, 1, 'Superuser all access');

-- INSERT INTO role_permissions (role_id, permission_id) VALUES(1, 2);
INSERT INTO group_roles (group_id, role_id, description) VALUES(2, 2, 'Admin all access');

-- standard
INSERT INTO group_roles (group_id, role_id, description) VALUES(3, 3, 'Standard user role');

-- login
INSERT INTO group_roles (group_id, role_id, description) VALUES(4, 4, 'Login access');

-- rdf
INSERT INTO group_roles (group_id, role_id, description) VALUES(5, 5, 'RDF access');



CREATE TABLE role_permissions (
    id SERIAL PRIMARY KEY, 
    role_id INTEGER NOT NULL,
    permission_id INTEGER NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(role_id, permission_id),
    FOREIGN KEY(role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY(permission_id) REFERENCES permissions(id) ON DELETE CASCADE);
CREATE INDEX role_permissions_role_id_idx ON role_permissions (role_id, permission_id);
CREATE TRIGGER role_permissions_updated_trigger
    BEFORE UPDATE
    ON
        role_permissions
    FOR EACH ROW
EXECUTE PROCEDURE update_at_updated();

-- super/user admin
INSERT INTO role_permissions (role_id, permission_id, description) VALUES(1, 1, 'Superuser all access');

-- INSERT INTO role_permissions (role_id, permission_id) VALUES(1, 2);
INSERT INTO role_permissions (role_id, permission_id, description) VALUES(2, 1, 'Admin all access');

-- standard
INSERT INTO role_permissions (role_id, permission_id, description) VALUES(3, 2, 'Standard user read access');

-- users can login
INSERT INTO role_permissions (role_id, permission_id, description) VALUES(4, 4, 'Login access');
INSERT INTO role_permissions (role_id, permission_id, description) VALUES(4, 2, 'Login user read access');


-- rdf
INSERT INTO role_permissions (role_id, permission_id, description) VALUES(5, 5, 'RDF read access');


DROP TABLE IF EXISTS users;
CREATE TABLE users (
    id SERIAL PRIMARY KEY, 
    public_id TEXT NOT NULL UNIQUE,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL DEFAULT '',
    first_name TEXT NOT NULL DEFAULT '',
    last_name TEXT NOT NULL DEFAULT '',
    is_locked BOOLEAN NOT NULL DEFAULT false,
    email_verified_at TIMESTAMP DEFAULT 'epoch' NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL);
CREATE INDEX users_public_id_idx ON users (public_id);
-- CREATE INDEX name ON users (first_name, last_name);
CREATE INDEX users_username_idx ON users (username);
CREATE INDEX users_email_idx ON users (email);
CREATE TRIGGER users_updated_trigger
    BEFORE UPDATE
    ON
        users
    FOR EACH ROW
EXECUTE PROCEDURE update_at_updated();

CREATE TABLE user_groups (
    id SERIAL PRIMARY KEY, 
    user_id INTEGER NOT NULL,
    group_id INTEGER NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(user_id, group_id),
    FOREIGN KEY(user_id) REFERENCES users(id),
    FOREIGN KEY(group_id) REFERENCES groups(id));
CREATE INDEX user_groups_user_id_idx ON user_groups (user_id, group_id);
CREATE TRIGGER user_groups_updated_trigger
    BEFORE UPDATE
    ON
        user_groups
    FOR EACH ROW
EXECUTE PROCEDURE update_at_updated();

INSERT INTO user_groups (user_id, group_id, description) VALUES(1, 1, 'Superuser all access');
INSERT INTO user_groups (user_id, group_id, description) VALUES(2, 2, 'Admin all access');



CREATE TABLE user_roles (
    id SERIAL PRIMARY KEY, 
    user_id INTEGER NOT NULL,
    role_id INTEGER NOT NULL, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(user_id, role_id),
    FOREIGN KEY(user_id) REFERENCES users(id),
    FOREIGN KEY(role_id) REFERENCES roles(id));
CREATE INDEX users_roles_user_id_idx ON users_roles (user_id, role_id);
CREATE TRIGGER users_roles_updated_trigger
    BEFORE UPDATE
    ON
        users_roles
    FOR EACH ROW
EXECUTE PROCEDURE update_at_updated();

 

DROP TABLE IF EXISTS api_keys; 
CREATE TABLE api_keys (
    id SERIAL PRIMARY KEY, 
    user_id INTEGER NOT NULL,
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


-- the superuser me --
INSERT INTO users (public_id, username, email, email_verified_at) VALUES (
    '01997349-0995-733e-8d23-eb14136f0486',
    'root',
    'edb-root@antonyholmes.dev',
    now()
);

-- su group --
INSERT INTO users_roles (user_id, role_id) VALUES (1, 1);

