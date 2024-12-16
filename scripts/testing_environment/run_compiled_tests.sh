#! /bin/bash
# this script assumes that we are running in a test environment
echo "Running setting up database for tests"

# this file creates the testing container
su postgres -c "initdb"

su postgres -c "postgres > $LOGFILE 2>&1 &"
sleep 3
su postgres -c "psql -U postgres -f $BACKEND_DIR/infra/1_test.sql  && echo "postgres databases and users created" "
$BACKEND_DIR/build/db_tools -server &> $DB_TOOLS_LOGFILE &
sleep 2
# loops through all test files and runs them, to run a singel test
# you just have to call the file like this test/blah.test
echo "Running unit tests"
for file in $BACKEND_DIR/build/test/unit/*
do
echo "__________________ $file ____________________________"
$file
done
echo "Running db_state tests"
for file in $BACKEND_DIR/build/test/db_state/*
do
echo "__________________ $file ____________________________"
$file
done
