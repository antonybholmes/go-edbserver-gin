
DROP TABLE IF EXISTS users_roles; 
DROP TABLE IF EXISTS roles_permissions;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS api_keys;
DROP TABLE IF EXISTS users;

CREATE TABLE permissions (
    id INTEGER PRIMARY KEY NOT NULL AUTO_INCREMENT, 
    public_id VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL UNIQUE,
    description VARCHAR(255) NOT NULL DEFAULT "",
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL);
CREATE INDEX roles_name_idx ON permissions (name);

INSERT INTO permissions (public_id, name, description) VALUES('01997350-f1db-734a-aabc-b738772a9d0c', 'Super', 'Superuser');
INSERT INTO permissions (public_id, name, description) VALUES('01997351-06c7-7f0d-b026-c51376a044ee', 'Admin', 'Administrator');
INSERT INTO permissions (public_id, name, description) VALUES('01997351-16e5-70f6-b869-ba08cdac4c85', 'User', 'User');
INSERT INTO permissions (public_id, name, description) VALUES('01997351-2586-7e76-8a34-db50b222d47a', 'Signin', 'User can sign in');
INSERT INTO permissions (public_id, name, description) VALUES('01997351-36fe-7e77-b06b-8222ab057601', 'RDF', 'User can view RDF lab data');



CREATE TABLE roles (
    id INTEGER PRIMARY KEY NOT NULL AUTO_INCREMENT, 
    public_id VARCHAR(255) NOT NULL UNIQUE, 
    name VARCHAR(255) NOT NULL UNIQUE,
    description VARCHAR(255) NOT NULL DEFAULT "",
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL);
CREATE INDEX roles_name_idx ON roles (name);

INSERT INTO roles (public_id, name) VALUES('01997351-4c4b-7900-adb9-eeb64f772ed7', 'Super');
INSERT INTO roles (public_id, name) VALUES('01997351-7d4f-72bc-aba9-1fbd2d5d41a2', 'Admin');
INSERT INTO roles (public_id, name) VALUES('01997351-8b67-758e-b798-f200f70c653b', 'User');
-- INSERT INTO roles (public_id, name) VALUES('UZuAVHDGToa4F786IPTijA==', 'GetDNA');
INSERT INTO roles (public_id, name) VALUES('01997351-9ba1-7a6e-9560-fc03b3098665', 'Login');
INSERT INTO roles (public_id, name) VALUES('01997351-aac6-7fcc-b886-b3f0585b90d8', 'RDF');


CREATE TABLE roles_permissions (
    id INTEGER PRIMARY KEY NOT NULL AUTO_INCREMENT, 
    role_id INTEGER NOT NULL,
    permission_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(role_id, permission_id),
    FOREIGN KEY(role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY(permission_id) REFERENCES permissions(id) ON DELETE CASCADE);
CREATE INDEX roles_permissions_role_id_idx ON roles_permissions (role_id, permission_id);

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
    id INTEGER PRIMARY KEY NOT NULL AUTO_INCREMENT, 
    public_id VARCHAR(255) NOT NULL UNIQUE,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL DEFAULT '',
    first_name VARCHAR(255) NOT NULL DEFAULT '',
    last_name VARCHAR(255) NOT NULL DEFAULT '',
    is_locked BOOLEAN NOT NULL DEFAULT false,
    email_verified_at DATETIME DEFAULT '1000-01-01' NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL);
CREATE INDEX users_public_id_idx ON users (public_id);
-- CREATE INDEX name ON users (first_name, last_name);
CREATE INDEX users_username_idx ON users (username);
CREATE INDEX users_email_idx ON users (email);

 
CREATE TABLE users_roles (
    id INTEGER PRIMARY KEY NOT NULL AUTO_INCREMENT, 
    user_id INTEGER NOT NULL,
    role_id INTEGER NOT NULL, 
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(user_id, role_id),
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(role_id) REFERENCES roles(id) ON DELETE CASCADE);
CREATE INDEX users_roles_user_id_idx ON users_roles (user_id, role_id);


DROP TABLE IF EXISTS api_keys; 
CREATE TABLE api_keys (
    id INTEGER PRIMARY KEY NOT NULL AUTO_INCREMENT, 
    user_id INTEGER NOT NULL,
    api_key VARCHAR(255) NOT NULL, 
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(user_id, api_key),
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE);
CREATE INDEX api_keys_api_key_idx ON api_keys (api_key);
 

-- the superuser me --
INSERT INTO users (public_id, username, email, is_locked, email_verified_at) VALUES (
    '25bhmb459eg7',
    'root',
    'edb-root@antonyholmes.dev',
    true,
    now()
);

INSERT INTO users (public_id, username, email, password, is_locked, email_verified_at) VALUES (
    'fr87kybn5q14',
    'rdf',
    'rdf@antonyholmes.dev',
    '$2a$10$su3OksRXYrpx6JYoYyN0heK8UnOXjCDorYvqlYAZ5Kov8y7L5Ze4O',
    true,
    now()
);

-- su group --
INSERT INTO users_roles (user_id, role_id) VALUES (1, 1);

-- RDF member of RDF role -
INSERT INTO users_roles (user_id, role_id) VALUES (2, 3);
INSERT INTO users_roles (user_id, role_id) VALUES (2, 4);
INSERT INTO users_roles (user_id, role_id) VALUES (2, 5);
-- default key --
INSERT INTO api_keys (user_id, api_key) VALUES (1, '01997352-5d48-7f9d-b244-5cee8a0239dd');

-- RDF api key --
INSERT INTO api_keys (user_id, api_key) VALUES (2, '01997352-4858-7909-a3a2-5c0a0329f083');
INSERT INTO api_keys (user_id, api_key) VALUES (2, '01997352-71f9-7bd1-b288-a885227ce00a'); 
 