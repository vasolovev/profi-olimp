# Educational Institution API

A RESTful API service for managing an educational institution's students and academic groups, built using clean architecture principles.

## Table of Contents

- [Features](#features)
- [Architecture](#architecture)
- [API Documentation](#api-documentation)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Configuration](#configuration)
  - [Running the Application](#running-the-application)
- [Database Schema](#database-schema)
- [API Testing](#api-testing)

## Features

- Create, read, update, and delete students
- Create, read, update, and delete academic groups
- Search students by name or group name
- Search groups by name
- Hierarchical group structure (groups can have subgroups)
- Protection against deleting groups with subgroups

## Architecture

This project follows the clean architecture pattern with the following layers:

1. **Entities** - Core business objects
   - Student
   - Group

2. **Use Cases** - Application business rules
   - StudentUseCase
   - GroupUseCase

3. **Controllers/Adapters** - Interface adapters
   - HTTP REST API controllers
   - PostgreSQL repositories

4. **Frameworks & Drivers** - External frameworks and tools
   - Fiber web framework
   - PostgreSQL database

## API Documentation

The API is documented using Swagger. When the application is running, you can access the Swagger UI at:

```
http://127.0.0.1:8080/swagger/index.html
```

## Database Schema

### Groups Table

```sql
CREATE TABLE groups (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    parent_id INTEGER NULL REFERENCES groups(id)
);
```

### Students Table

```sql
CREATE TABLE students (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    group_id INTEGER NOT NULL REFERENCES groups(id)
);
```

## API Testing

You can test the API using curl or any API testing tool like Postman. Here are some example requests:

### Create a Group

```bash
curl -X POST http://localhost:8080/groups \
  -H 'Content-Type: application/json' \
  -d '{"name": "Computer Science"}'
```

### Create a Student

```bash
curl -X POST http://localhost:8080/students \
  -H 'Content-Type: application/json' \
  -d '{"name": "John Doe", "email": "john@example.com", "group_id": 1}'
```

### Get All Students

```bash
curl -X GET http://localhost:8080/students
```

### Search Students by Name or Group Name

```bash
curl -X GET 'http://localhost:8080/students?query=john'
```

### Update a Student

```bash
curl -X PUT http://localhost:8080/students/1 \
  -H 'Content-Type: application/json' \
  -d '{"name": "John Smith", "group_id": 2}'
```

### Delete a Student

```bash
curl -X DELETE http://localhost:8080/students/1
```