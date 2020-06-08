#!/bin/bash

execute () {
  $1
  if [ $? -eq 1 ]
  then
    exit 1
  fi
}

run_test () {
  package=$1
  exit_code=$2

  (
    cd "testdata/$package" || exit 1
    "$pwd/gonforce"
    if [ $? -ne "$exit_code" ]
    then
      echo "$package test case failed"
      exit 1
    else
      echo "$package test case passed"
    fi
    cd ../..
  )
}

execute "go test ./..."
execute "golint ./..."
execute "goimports -w -d $(find . -type f -name '*.go' -not -path './testdata/*')"
execute "gofmt -w -d $(find . -type f -name '*.go' -not -path './testdata/*')"
execute "go run main.go"

pwd=$(pwd)
echo "$pwd"
go build -o gonforce

run_test "no_gonforce_file" 1
run_test "invalid_gonforce_file" 1
run_test "using_blacklisted_import" 1
run_test "using_not_whitelisted_import" 1
run_test "using_whitelist_exception" 1
run_test "failed_whitelist_exception" 1
run_test "invalid_rule_structure" 1
run_test "invalid_import_in_subdir" 1
run_test "subfolder_using_invalid_import" 1
run_test "valid" 0

rm gonforce
echo "All tests passed..."
