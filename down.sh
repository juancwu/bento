#!/bin/bash

URL=$(cat .env | grep BENTO_DB_CONN | sed 's/^[^=]*=//')
libsql-migrate down --url $URL --path ./migrations
