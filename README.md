# Overview
Backend for my submission to CVWO's assignment.

# Features
- Account System powered by JWT
- Projects with tasks
- Multiple users per project
- Tagging for tasks

# TODO
- Track projects (with CRUD) *Done*
- Track tasks (with CRUD) *Done*
- Tagging for tasks and projects
- Add users to projects *Done*

# Setup
There are three environment variables required:
A random key for JWT generation: "jwtkey"
The port the server will run on: "PORT" 
The database url string: "DATABASE_URL"