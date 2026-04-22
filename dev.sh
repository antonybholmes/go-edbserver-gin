python scripts/update_version.py
#export $(cat consts.env | xargs)
APP_ENV=development go run --tags "fts5" .
