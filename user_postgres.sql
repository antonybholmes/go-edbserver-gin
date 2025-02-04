CREATE OR REPLACE FUNCTION updated_at_updated()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

DROP TABLE IF EXISTS users_roles;
DROP TABLE IF EXISTS user_roles_permissions;
DROP TABLE IF EXISTS user_permissions;
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS users;

CREATE TABLE user_permissions (
    id SERIAL PRIMARY KEY, 
    public_id VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL UNIQUE,
    description VARCHAR(255) NOT NULL DEFAULT '',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL);
CREATE INDEX user_permissions_name_idx ON user_permissions (name);

INSERT INTO user_permissions (public_id, name, description) VALUES('uwkrk2ljj387', 'Super', 'Superuser');
INSERT INTO user_permissions (public_id, name, description) VALUES('iz4kbfy3z0a3', 'Admin', 'Administrator');
INSERT INTO user_permissions (public_id, name, description) VALUES('loq75e7zqcbl', 'User', 'User');
INSERT INTO user_permissions (public_id, name, description) VALUES('kflynb03pxbj', 'Login', 'Can login');
INSERT INTO user_permissions (public_id, name, description) VALUES('og1o5d0p0mjy', 'RDF', 'Can view RDF lab data');



CREATE TABLE user_roles (
    id SERIAL PRIMARY KEY, 
    public_id VARCHAR(255) NOT NULL UNIQUE, 
    name VARCHAR(255) NOT NULL UNIQUE,
    description VARCHAR(255) NOT NULL DEFAULT '',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL);
CREATE INDEX user_roles_name_idx ON user_roles (name);
CREATE TRIGGER user_roles_updated_trigger
    BEFORE UPDATE
    ON
        user_roles
    FOR EACH ROW
EXECUTE PROCEDURE updated_at_updated();


INSERT INTO user_roles (public_id, name) VALUES('p1gbjods0h90', 'Super');
INSERT INTO user_roles (public_id, name) VALUES('mk4bgg4w43fp', 'Admin');
INSERT INTO user_roles (public_id, name) VALUES('3xvte0ik4aq4', 'User');
-- INSERT INTO user_roles (public_id, name) VALUES('UZuAVHDGToa4F786IPTijA==', 'GetDNA');
INSERT INTO user_roles (public_id, name) VALUES('x4ewk9papip2', 'Login');
INSERT INTO user_roles (public_id, name) VALUES('kh2yynyheqhv', 'RDF');


CREATE TABLE user_roles_permissions (
    id SERIAL PRIMARY KEY, 
    role_id INTEGER NOT NULL,
    permission_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(role_id, permission_id),
    FOREIGN KEY(role_id) REFERENCES user_roles(id),
    FOREIGN KEY(permission_id) REFERENCES user_permissions(id));
CREATE INDEX roles_permissions_role_id_idx ON user_roles_permissions (role_id, permission_id);
CREATE TRIGGER user_roles_permissions_updated_trigger
    BEFORE UPDATE
    ON
        user_roles_permissions
    FOR EACH ROW
EXECUTE PROCEDURE updated_at_updated();

-- super/user admin
INSERT INTO user_roles_permissions (role_id, permission_id) VALUES(1, 1);
INSERT INTO user_roles_permissions (role_id, permission_id) VALUES(1, 2);
INSERT INTO user_roles_permissions (role_id, permission_id) VALUES(2, 2);

--
-- standard
INSERT INTO user_roles_permissions (role_id, permission_id) VALUES(3, 3);

-- users can login
INSERT INTO user_roles_permissions (role_id, permission_id) VALUES(4, 4);

-- rdf
INSERT INTO user_roles_permissions (role_id, permission_id) VALUES(5, 5);

 
CREATE TABLE users (
    id SERIAL PRIMARY KEY, 
    public_id VARCHAR(255) NOT NULL UNIQUE,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL DEFAULT '',
    first_name VARCHAR(255) NOT NULL DEFAULT '',
    last_name VARCHAR(255) NOT NULL DEFAULT '',
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
EXECUTE PROCEDURE updated_at_updated();
 

CREATE TABLE users_roles (
    id SERIAL PRIMARY KEY, 
    user_id INTEGER NOT NULL,
    role_id INTEGER NOT NULL, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(user_id, role_id),
    FOREIGN KEY(user_id) REFERENCES users(id),
    FOREIGN KEY(role_id) REFERENCES user_roles(id));
CREATE INDEX users_roles_user_id_idx ON users_roles (user_id, role_id);
CREATE TRIGGER users_roles_updated_trigger
    BEFORE UPDATE
    ON
        users_roles
    FOR EACH ROW
EXECUTE PROCEDURE updated_at_updated();


-- the superuser me --
INSERT INTO users (public_id, username, email, email_verified_at) VALUES (
    '25bhmb459eg7',
    'root',
    'antony@antonyholmes.dev',
    now()
);

-- su group --
INSERT INTO users_roles (user_id, role_id) VALUES (1, 1);