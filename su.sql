-- the superuser me --
INSERT INTO users (public_id, username, email, email_verified_at) VALUES (
    '25bhmb459eg7',
    'root',
    'antony@antonyholmes.dev',
    now()
);

-- su group --
INSERT INTO users_roles (user_id, role_id) VALUES (1, 1);