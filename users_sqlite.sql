PRAGMA journal_mode = WAL;
PRAGMA foreign_keys = ON;

DROP TABLE IF EXISTS user_permissions;
CREATE TABLE user_permissions (
    id INTEGER PRIMARY KEY ASC, 
    public_id TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT "",
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL);
CREATE INDEX roles_name_idx ON user_permissions (name);

INSERT INTO user_permissions (public_id, name, description) VALUES('uwkrk2ljj387', 'Super', 'Superuser');
INSERT INTO user_permissions (public_id, name, description) VALUES('iz4kbfy3z0a3', 'Admin', 'Administrator');
INSERT INTO user_permissions (public_id, name, description) VALUES('loq75e7zqcbl', 'User', 'User');
INSERT INTO user_permissions (public_id, name, description) VALUES('kflynb03pxbj', 'Login', 'Can login');
INSERT INTO user_permissions (public_id, name, description) VALUES('og1o5d0p0mjy', 'RDF', 'Can view RDF lab data');

DROP TABLE IF EXISTS user_roles;
CREATE TABLE user_roles (
    id INTEGER PRIMARY KEY ASC, 
    public_id TEXT NOT NULL UNIQUE, 
    name TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT "",
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL);
CREATE INDEX user_roles_name_idx ON user_roles (name);

INSERT INTO user_roles (public_id, name) VALUES('p1gbjods0h90', 'Super');
INSERT INTO user_roles (public_id, name) VALUES('mk4bgg4w43fp', 'Admin');
INSERT INTO user_roles (public_id, name) VALUES('3xvte0ik4aq4', 'User');
-- INSERT INTO user_roles (public_id, name) VALUES('UZuAVHDGToa4F786IPTijA==', 'GetDNA');
INSERT INTO user_roles (public_id, name) VALUES('x4ewk9papip2', 'Login');
INSERT INTO user_roles (public_id, name) VALUES('kh2yynyheqhv', 'RDF');

DROP TABLE IF EXISTS user_roles_permissions;
CREATE TABLE user_roles_permissions (
    id INTEGER PRIMARY KEY ASC, 
    role_id INTEGER NOT NULL,
    permission_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(role_id, permission_id),
    FOREIGN KEY(role_id) REFERENCES user_roles(id) ON DELETE CASCADE,
    FOREIGN KEY(permission_id) REFERENCES user_permissions(id) ON DELETE CASCADE);
CREATE INDEX roles_permissions_role_id_idx ON user_roles_permissions (role_id, permission_id);

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

DROP TABLE IF EXISTS users;
CREATE TABLE users (
    id INTEGER PRIMARY KEY ASC, 
    public_id TEXT NOT NULL UNIQUE,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL DEFAULT '',
    first_name TEXT NOT NULL DEFAULT '',
    last_name TEXT NOT NULL DEFAULT '',
    -- use the epoch as a default no --
    email_verified_at TIMESTAMP DEFAULT '1970-01-01' NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL);
CREATE INDEX users_public_id_idx ON users (public_id);
-- CREATE INDEX name ON users (first_name, last_name);
CREATE INDEX users_username_idx ON users (username);
CREATE INDEX users_email_idx ON users (email);

CREATE TRIGGER users_updated_trigger AFTER UPDATE ON users
BEGIN
      update users SET updated_at = CURRENT_TIMESTAMP WHERE id=NEW.id;
END;

DROP TABLE IF EXISTS users_roles;
CREATE TABLE users_roles (
    id INTEGER PRIMARY KEY ASC, 
    user_id INTEGER NOT NULL,
    role_id INTEGER NOT NULL, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(user_id, role_id),
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(role_id) REFERENCES user_roles(id) ON DELETE CASCADE);
CREATE INDEX users_roles_user_id_idx ON users_roles (user_id, role_id);


CREATE TABLE users_sessions(
  id INTEGER PRIMARY KEY ASC,
  public_id TEXT NOT NULL,
  session_id INTEGER NOT NULL UNIQUE,
  FOREIGN KEY(public_id) REFERENCES users(public_id)
);
CREATE INDEX users_sessions_public_id_idx ON users_sessions (public_id);
CREATE INDEX users_sessions_session_id_idx ON users_sessions (session_id);


DROP TABLE IF EXISTS api_keys;
CREATE TABLE api_keys (
    id INTEGER PRIMARY KEY ASC, 
    user_id INTEGER NOT NULL,
    key TEXT NOT NULL, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(user_id, key),
    FOREIGN KEY(user_id) REFERENCES users(id));
CREATE INDEX api_keys_key_idx ON api_keys (key);