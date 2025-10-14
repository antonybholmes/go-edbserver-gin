CREATE OR REPLACE FUNCTION update_at_updated()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

DROP TABLE IF EXISTS users_roles;
DROP TABLE IF EXISTS roles_permissions;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS api_keys;
DROP TABLE IF EXISTS users;

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

INSERT INTO permissions (public_id, name, description) VALUES('01997350-f1db-734a-aabc-b738772a9d0c', 'Super', 'Superuser');
INSERT INTO permissions (public_id, name, description) VALUES('01997351-06c7-7f0d-b026-c51376a044ee', 'Admin', 'Administrator');
INSERT INTO permissions (public_id, name, description) VALUES('01997351-16e5-70f6-b869-ba08cdac4c85', 'User', 'User');
INSERT INTO permissions (public_id, name, description) VALUES('01997351-2586-7e76-8a34-db50b222d47a', 'Signin', 'User can sign in');
INSERT INTO permissions (public_id, name, description) VALUES('01997351-36fe-7e77-b06b-8222ab057601', 'RDF', 'User can view RDF lab data');



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


INSERT INTO roles (public_id, name) VALUES('01997351-4c4b-7900-adb9-eeb64f772ed7', 'Super');
INSERT INTO roles (public_id, name) VALUES('01997351-7d4f-72bc-aba9-1fbd2d5d41a2', 'Admin');
INSERT INTO roles (public_id, name) VALUES('01997351-8b67-758e-b798-f200f70c653b', 'User');
-- INSERT INTO roles (public_id, name) VALUES('UZuAVHDGToa4F786IPTijA==', 'GetDNA');
INSERT INTO roles (public_id, name) VALUES('01997351-9ba1-7a6e-9560-fc03b3098665', 'Login');
INSERT INTO roles (public_id, name) VALUES('01997351-aac6-7fcc-b886-b3f0585b90d8', 'RDF');


CREATE TABLE roles_permissions (
    id SERIAL PRIMARY KEY, 
    role_id INTEGER NOT NULL,
    permission_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(role_id, permission_id),
    FOREIGN KEY(role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY(permission_id) REFERENCES permissions(id) ON DELETE CASCADE);
CREATE INDEX roles_permissions_role_id_idx ON roles_permissions (role_id, permission_id);
CREATE TRIGGER roles_permissions_updated_trigger
    BEFORE UPDATE
    ON
        roles_permissions
    FOR EACH ROW
EXECUTE PROCEDURE update_at_updated();

-- super/user admin
INSERT INTO roles_permissions (role_id, permission_id) VALUES(1, 1);
INSERT INTO roles_permissions (role_id, permission_id) VALUES(1, 2);
INSERT INTO roles_permissions (role_id, permission_id) VALUES(2, 2);

--
-- standard
INSERT INTO roles_permissions (role_id, permission_id) VALUES(3, 3);

-- users can login
INSERT INTO roles_permissions (role_id, permission_id) VALUES(4, 4);

-- rdf
INSERT INTO roles_permissions (role_id, permission_id) VALUES(5, 5);


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
 

CREATE TABLE users_roles (
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