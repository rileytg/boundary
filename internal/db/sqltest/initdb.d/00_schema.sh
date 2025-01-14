#!/usr/bin/env bash
#
# This script is designed to run as an init script in a postgres docker container.
# The docker-entrypoint.sh script does not traverse a nested directory
# structue to load sql files at container start. However, it does allow for
# running .sh files. This script should run first to execute the boundary sql
# migration files in the correct order.
#
# See Initialization Scripts https://hub.docker.com/_/postgres
set -e
shopt -s globstar

## Taken from postgres docker image's docker-entrypoint.sh
## See: https://github.com/docker-library/postgres/blob/517c64f87e6661366b415df3f2273c76cea428b0/docker-entrypoint.sh#L175-L187
# Execute sql script, passed via stdin (or -f flag of pqsl)
# usage: docker_process_sql [psql-cli-args]
#    ie: docker_process_sql --dbname=mydb <<<'INSERT ...'
#    ie: docker_process_sql -f my-file.sql
#    ie: docker_process_sql <my-file.sql
docker_process_sql() {
	local query_runner=( psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --no-password )
	if [ -n "$POSTGRES_DB" ]; then
		query_runner+=( --dbname "$POSTGRES_DB" )
	fi

	PGHOST= PGHOSTADDR= "${query_runner[@]}" "$@"
}

# Run migrations in order.
for file in $(ls -v /migrations/**/*.sql); do
  echo running "$file"
  docker_process_sql -f "$file"
done
