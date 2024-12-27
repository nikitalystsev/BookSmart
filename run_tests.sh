#!/bin/sh

. ./.env

export UNIT_TESTS_IS_SUCCESS=0
if go test -v -shuffle on ./internal/tests_for_testing/unitTests/
then
  export UNIT_TESTS_IS_SUCCESS=1
fi

export INTEGRATION_TESTS_IS_SUCCESS=0
if go test -v -shuffle on ./internal/tests_for_testing/integrationTests/ --parallel 2
then
  export INTEGRATION_TESTS_IS_SUCCESS=1
fi


go test -v -shuffle on ./internal/tests_for_testing/e2eTest/

unset UNIT_TESTS_IS_SUCCESS
unset INTEGRATION_TESTS_IS_SUCCESS

cp -R ./internal/tests_for_testing/allure-report/history ./internal/tests_for_testing/allure-results/
rm -rf ./internal/tests_for_testing/allure-report/
cp ./internal/tests_for_testing/environment.properties ./internal/tests_for_testing/allure-results
cd ./internal/tests_for_testing/ && allure generate && allure open
