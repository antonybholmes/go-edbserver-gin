rm data/users.db
cat users.sql | sqlite3 data/users.db
cat su.sql | sqlite3 data/users.db