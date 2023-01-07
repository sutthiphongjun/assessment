#!/bin/bash

docker stop expensetracking
docker rm expensetracking

docker build  -t expensetracking:v1 .

docker run -d -p 2565:2565 --name expensetracking -e PORT=:2565 -e DATABASE_URL=postgres://pupffhjj:cTLk0BZ4OkVPGze0vhiED7wOZjO5ZMyN@tiny.db.elephantsql.com/pupffhjj expensetracking:v1

docker ps