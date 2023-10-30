#!/bin/bash

URL=$(cat .env | grep BENTO_DB_CONN | sed 's/^[^=]*=//')
libsql-migrate up --url $URL --path ./migrations
