# Gotasks API

Gotasks is a RESTful API built with Go using the Fiber framework. It provides endpoints for managing users, tasks, projects, and authentication.

# Table of Contents
- [Installation](#installation)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
  - [Authentication](#authentication)
  - [User Management](#user-management)
  - [Task Management](#task-management)
  - [Admin Operations](#admin-operations)
  - [Project Management](#project-management)

## Installation
1. Clone the repository
```
git clone https://github.com/ficontini/gotasks.git
```
2. Rename .env.example to .env and fill with your environment variables
3. Setup database
    1. DynamoDB 
    ```
    make deploy
    ```
4. Seeding the database
```
make seed
```
## Usage
```
make run 
```
## API Endpoints
### Authentication
* `POST /api/auth` : Authenticate a user
### User Management
* `POST /api/user` : Create a user
* `POST /api/v1/user/reset-password` : Reset the password of the authenticated user
* `GET /ap1/v1/user` : Get authenticated user
### Task Management
* `GET /api/v1/task`: Get all tasks associated with the authenticated user
* `GET /api/v1/task/all`: Get all tasks
* `GET /api/v1/task/:id`: Get a specific task
* `POST /api/v1/task/:id/assign`: Assign a task to the authenticated user
* `POST /api/v1/task/:id/complete`: Complete a task
* `POST /api/v1/task`: Create a task
### Admin Operations:
* `PUT /api/v1/admin/user/:id/enable`: Enable a user
* `PUT /api/v1/admin/user/:id/disable`: Disable a user
* `GET /api/v1/admin/user/:id`: Get a specific user
* `GET /api/v1/admin/user`: Get all users
* `POST /api/v1/admin/task`: Get all tasks 
* `DELETE /api/v1/admin/task/:id`: Delete a task 
* `POST /api/v1/admin/task/:id/assign`: Assign a task to a user 
### Project Management:
* `POST /project`: Create a project
* `POST /project/:id/task`: Assign an existing task to a project
* `GET /project/:id/task`: Get all tasks of a project