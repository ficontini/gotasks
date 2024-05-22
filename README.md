# Gotasks API

Gotasks is a RESTful API built with Go using the Fiber framework. It provides endpoints for managing users, tasks, projects, and authentication.

## Table of Contents
- [Installation] (#installation)
- [Usage] (#usage)

## Installation
1. Clone the repository
```bash
git clone https://github.com/ficontini/gotasks.git
```
2. Rename .env.example to .env and fill with your environment variables
3. Setup database
3.1. DynamoDB 
```bash
make deploy
```
4. Seeding the database
```bash
make seed
```
## Usage
```bash
make run 
```
