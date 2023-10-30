#!/bin/bash

if [ $# -gt 0 ]; then
    libsql-migrate gen "$1" --path ./migrations
else
    echo "No migration file name provided."
    exit 1
fi
