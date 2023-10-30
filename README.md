# Bento 🍱

A simple backend to manage the sweet `.env`. It is still under development but once complete it should:

- Separate environment variables per project (Create a new bento)
- Separate different scopes for each project (development, producition, etc...)
- Edit specific environment variables based on project and scope
- Delete specific environment variables based on project and scope
- Download the sweet `.env` file.

## Pre-requisites

1. Install libsql-migrate

`go install github.com/juancwu/libsql-migrate@latest`

## Current Plan

Read the current plan for this project in [here](/docs/PLANNING.md)

> This backend should be treated as a service which any client that is able to make HTTP request can tap into and use.

The stack being used in this project is:

- Golang
- PlanetScale
- Chi
