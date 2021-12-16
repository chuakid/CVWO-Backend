# Overview
Backend for my submission to CVWO's assignment.
Currently deployed to https://cvwo.herokuapp.com/

# Features
- Account System powered by JWT
- Projects with tasks
- Multiple users per project
- Tagging for tasks

# Development
## Requirements
- postgresql (version >=12)
- Go (version >= 1.16)

## Environment variables
- jwtkey: key for JWT generation, can be any random string
- PORT: The port the server will run on (defaults to 8000)
- DATABASE_URL: database connection string for postgresql
  - postgresql://[user[:password]@][netloc][:port][/dbname][?param1=value1&...]
  - Example: postgresql://user@localhost
## Setup
1. `git clone` the repository
2. `go get -u ./...` to get all dependencies
3. `go *.go` to spin up a development server

`Gow` is recommended to watch the folder and rebuild on save 

# Deployment 
`go build` generates a binary based on the current platform and architecture. It is then executable using `./executablename`.


# TODO
- Track projects (with CRUD) *Done*
- Track tasks (with CRUD) *Done*
- Tagging for tasks and projects
- Add users to projects *Done*