
DROP TABLE IF EXISTS users_roles; 
DROP TABLE IF EXISTS roles_permissions;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS api_keys;
DROP TABLE IF EXISTS users;

CREATE TABLE permissions (
    id INTEGER PRIMARY KEY NOT NULL AUTO_INCREMENT, 
    uuid VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL UNIQUE,
    description VARCHAR(255) NOT NULL DEFAULT "",
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL);
CREATE INDEX roles_name_idx ON permissions (name);

INSERT INTO permissions (uuid, name, description) VALUES('uwkrk2ljj387', 'Super', 'Superuser');
INSERT INTO permissions (uuid, name, description) VALUES('iz4kbfy3z0a3', 'Admin', 'Administrator');
INSERT INTO permissions (uuid, name, description) VALUES('loq75e7zqcbl', 'User', 'User');
INSERT INTO permissions (uuid, name, description) VALUES('kflynb03pxbj', 'Login', 'Can login');
INSERT INTO permissions (uuid, name, description) VALUES('og1o5d0p0mjy', 'RDF', 'Can view RDF lab data');



CREATE TABLE roles (
    id INTEGER PRIMARY KEY NOT NULL AUTO_INCREMENT, 
    uuid VARCHAR(255) NOT NULL UNIQUE, 
    name VARCHAR(255) NOT NULL UNIQUE,
    description VARCHAR(255) NOT NULL DEFAULT "",
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL);
CREATE INDEX roles_name_idx ON roles (name);

INSERT INTO roles (uuid, name) VALUES('p1gbjods0h90', 'Super');
INSERT INTO roles (uuid, name) VALUES('mk4bgg4w43fp', 'Admin');
INSERT INTO roles (uuid, name) VALUES('3xvte0ik4aq4', 'User');
-- INSERT INTO roles (uuid, name) VALUES('UZuAVHDGToa4F786IPTijA==', 'GetDNA');
INSERT INTO roles (uuid, name) VALUES('x4ewk9papip2', 'Signin');
INSERT INTO roles (uuid, name) VALUES('kh2yynyheqhv', 'RDF');


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
    uuid VARCHAR(255) NOT NULL UNIQUE,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL DEFAULT '',
    first_name VARCHAR(255) NOT NULL DEFAULT '',
    last_name VARCHAR(255) NOT NULL DEFAULT '',
    is_locked BOOLEAN NOT NULL DEFAULT false,
    email_verified_at DATETIME DEFAULT '1000-01-01' NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL);
CREATE INDEX users_uuid_idx ON users (uuid);
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
INSERT INTO users (uuid, username, email, is_locked, email_verified_at) VALUES (
    '25bhmb459eg7',
    'root',
    'edb-root@antonyholmes.dev',
    true,
    now()
);

INSERT INTO users (uuid, username, email, password, is_locked, email_verified_at) VALUES (
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
INSERT INTO api_keys (user_id, api_key) VALUES (1, '4715057f-0b11-49d0-8a7b-296a2248046d');

-- RDF api key --
INSERT INTO api_keys (user_id, api_key) VALUES (2, 'f80e8d48-112b-4760-8efa-9754d3469f6b');
INSERT INTO api_keys (user_id, api_key) VALUES (2, '887af980-995b-46c3-80d7-6223491e398f'); 
 