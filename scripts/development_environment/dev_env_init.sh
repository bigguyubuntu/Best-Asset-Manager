#! /bin/bash
echo "Dev env startup script"
# run podman ps -a if the name pgsql doesnt exist then run the infro command then the 
# create contaienr command. then run the  podman start pgsql command.
# If it exists then run podman start pgsql command anyways.

PASSWD=123456
DB_NAME=bam_db
TEST_DB_NAME=bam_test_db
CONTAINER_NAME=bam_db
DB_USERNAME=bam_backend
DB_TOOL_USERNAME=bam_admin
# where thye db migrations are
MIGRATIONS="$PWD/db/migration"
DB_PORT=5535

# these are run inside the machine the go app will run in, not inside the db container.
export DB_HOST=localhost
export BACKEND_HOST=http://127.0.0.1
export DB_PORT=$DB_PORT
export DB_USERNAME=$DB_USERNAME
export DB_TOOL_USERNAME=$DB_TOOL_USERNAME
export DB_PASSWORD=$PASSWD
export DB_NAME=$DB_NAME
export TEST_DB_NAME=$TEST_DB_NAME
export FRONTEND=http://localhost:3000
export ENV_TYPE=dev
export MIGRATIONS=$MIGRATIONS
export BACKEND_PORT=9000
export DB_TOOLS_PORT=4444

alias db_tools_test='ENV_TYPE=test go run $(pwd)/tools/db_tools/cmd/main.go'
alias db_tools='go run $(pwd)/tools/db_tools/cmd/main.go'
alias restart_bam='podman machine stop && podman machine start && podman start pgsql && db_tools f'

# run the test tool in the background if we decided to run integration tests locally
var=dbTools db_tools_test -server &> $(pwd)/my_db_tools_log.txt &

if ! podman images | grep -q $CONTAINER_NAME;then
echo "database container image doesnt exist, will create it"
echo "preparing database up migrations"
podman build -f infra/dev/db.dockerfile -t $CONTAINER_NAME
echo "container image created"
else
echo "database container image already exist"
fi


if ! podman ps -a | grep -q pgsql
then
echo "database container not found, will create one"
podman run -p $DB_PORT:5432 -e POSTGRES_PASSWORD=$PASSWD --name pgsql -d $CONTAINER_NAME
echo "database container created"
else
echo "database container already exist"
fi

podman start pgsql
 
echo "started database container."
echo "run 'podman ps' to confirm, if it is not there run 'podman start pgsql' to start it"


echo "Finished setting up dev ennvironment"
