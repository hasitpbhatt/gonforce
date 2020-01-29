#!/bin/bash

run_test () {
  package=$1
  exit_code=$2

  cd $package
  $pwd/gonforce
  if [ $? -ne $exit_code ]
  then
    echo "$package test case failed"
    exit 1
  else
    echo "$package test case passed"
  fi
  cd ..
}

pwd=$(pwd)
echo $pwd
go build -o gonforce

cd testdata

run_test "no_gonforce_file" 1
run_test "invalid_gonforce_file" 1
run_test "valid" 0

cd ..

rm gonforce