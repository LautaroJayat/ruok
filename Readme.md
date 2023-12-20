# RUOK Scheduler

![Testing](https://github.com/back-end-labs/ruok/actions/workflows/test.yaml/badge.svg?event=push&branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/back-end-labs/ruok)](https://goreportcard.com/report/github.com/helm/helm)
![Version](https://img.shields.io/badge/version-unstable-blue)
[![License](https://img.shields.io/github/license/back-end-labs/ruok)](/LICENSE)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fback-end-labs%2Fruok.svg?type=shield&issueType=license)](https://app.fossa.com/projects/git%2Bgithub.com%2Fback-end-labs%2Fruok?ref=badge_shield&issueType=license)

<div align="center">
    <img  src="./assets/horizontal_logo.png" alt="RUOK Logo" />
</div>

RUOK Scheduler is an open-source tool designed for hassle-free service monitoring. Keep a close eye on your infrastructure effortlessly with our intuitive scheduler.

- [RUOK Scheduler](#ruok-scheduler)
  - [1. Introduction](#1-introduction)
    - [1.1 Purpose](#11-purpose)
    - [1.2 Why RUOK Scheduler?](#12-why-ruok-scheduler)
  - [2. Getting Started](#2-getting-started)
    - [2.1 Building from Source](#21-building-from-source)
    - [2.2 Preparing the Database](#21-preparing-the-database)
    - [2.3 Starting RUOK Scheduler](#23-starting-ruok-scheduler)
  - [3. Configurations](#3-configurations)
    - [3.1 DB user](#31-db-user)
    - [3.2 DB Password](#32-db-password)
    - [3.3 DB Host](#33-db-host)
    - [3.4 DB Port](#34-db-port)
    - [3.5 DB Name](#35-db-name)
    - [3.6 Application Name](#36-application-name)
    - [3.7 SSL Mode](#37-ssl-mode)
    - [3.8 Client Cert Password](#38-client-cert-password)
    - [3.9 Polling Interval](#38-client-cert-password)
    - [3.10 Max Number of Jobs](#310-max-number-of-jobs)
  - [4. Job Configuration](#4-job-configuration)
  - [5. HTTP API](#5-http-api)
    - [5.1 Create Jobs](#51-create-jobs)
    - [5.2 Update Jobs](#52-update-jobs)
    - [5.3 List Jobs](#53-list-jobs)
    - [5.4 List Job Executions](#54-list-job-executions)
    - [5.5 Get Instance Info](#55-get-instance-info)
  - [6. Cron Specification](#6-cron-specification)
  - [7. License](#7-license)

## 1. Introduction

### 1.1 Purpose

RUOK Scheduler serves the purpose of transforming a PostgreSQL database into a reliable and efficient broker for a backend service monitoring system. It simplifies the process of scheduling and monitoring services, offering a straightforward solution for scenarios where complex deployments are unnecessary.

### 1.2 Why RUOK Scheduler?

In many cases, deploying and managing monitoring systems can be complex and resource-intensive. RUOK Scheduler aims to address this challenge by providing a simple, yet effective, solution for users who prioritize ease of use and minimal configuration overhead.

## 2. Getting Started

### 2.1 Building from Source

To build RUOK Scheduler from source, use the following command:

```bash
make build
```

This will compile and output a single binary named `ruok` in the root directory of this project.

As the cli is implemented using [cobra](https://github.com/spf13/cobra) you can explore the available commands just by running `./ruok help`, which will output something like the following message:

```sh
./ruok help

# Turn your postgres database into a backend service monitor.
# Receive notifications via http, slack, sqs/sns, and much more!
#
# Usage:
#   ruok [flags]
#   ruok [command]
#
# Available Commands:
#   completion  Generate the autocompletion script for the specified shell
#   help        Help about any command
#   setupdb     Runs all migrations needed to setup postgres to work with ruok
#   start       Starts the scheduler main process
#   version     Print the version of ruok
#
# Flags:
#   -h, --help   help for ruok
#
# Use "ruok [command] --help" for more information about a command.
```

### 2.1 Preparing the Database

Before running the scheduler you will need to setup a database. This database will store information about the jobs (which endpoints to monitor, what to do in case of failure...) and also a registry of the execution results.

To migrate the database you will need to run:

```bash
./ruok setupdb
```

This command will do the following:

1. Create a new schema named `ruok`.
2. Create a function to get a timestamp in microseconds.
3. Create a function to check TLS status and version.
4. Create jobs and job_results tables.
5. Set [Row Security Policies](https://www.postgresql.org/docs/current/ddl-rowsecurity.html) for tables mentioned above.
6. Create two roles one for the scheduler and one to handle crud operations when using multiple instances of the scheduler.

To see what resources will be created, please check sql files for the [migrations command](./cmd/migrate/migrations/)

As this command needs to create and manage several resources, the postgres user/role provided to the cli must the correct permissions.

You can set those by exporting/using the following envs:

```bash
export DB_PASS=correct_pass # Database password (default: password)
export DB_USER=correct_user # Database user (default: user)
export DB_HOST=your_host    # Database host (default: localhost)
export DB_PORT=your_port    # Database port (default: 5432)
export DB_NAME=your_db      # Database name (default: db1)

./ruok setupdb
```

### 2.3 Starting RUOK Scheduler

Again, make sure you are using a user with correct privileges for this db client.

Our recommendation is to use ruok provided roles. To do it, you must execute the following SQL commands:

```sql
-- create a new user for ruok
CREATE ROLE new_ruok_user WITH LOGIN PASSWORD 'ruok_user_pass';

--- set the role created with './ruok migrate'
GRANT RUOK_SCHEDULER_ROLE TO new_ruok_user;
```

After setting the user, you can run the following commands to start the `ruok`:

```bash
# use the user created above
export DB_USER=ruok_user_pass

# use the pass for the user created above
export DB_PASS=new_ruok_user

# use a name to identify this instance
export APP_NAME=some_name
export DB_HOST=your_host
export DB_PORT=your_port
export DB_NAME=your_db

./ruok start
```

## 3 Configurations

All configurations for RUOK Scheduler are expected as environment variables. Below are the configurations along with their respective environment variables:

### 3.1 DB user

Use this environment to provide a user for the ruok client. We recommend a user with the provided `RUOK_SCHEDULER_ROLE` set, but you can craft your own.

```bash
DB_PASS                 # Database password (default: password)
```

### 3.2 DB Password

Use this environment variable to provide the password for the user above.

```bash
DB_USER                 # Database user (default: user)
```

### 3.3 DB Host

Use this environment to provide the hostname of your postgres deployment. It will be used by the `ruok` process to interact with it.

```bash
DB_HOST                 # Database host (default: localhost)
```

### 3.4 DB Port

Use this environment to provide the port where `ruok` can reach your db.

```bash
DB_PORT                 # Database port (default: 5432)
```

### 3.5 DB Name

Use this environment to indicate which db the process needs to interact with.
Make sure `./ruok setupdb` and `./ruok start` are pointing to the same db.

```bash
DB_NAME                 # Database name (default: db1)
```

### 3.6 Application Name

Is recommended to use different application names for different `ruok` processes interacting with the same db. This ensures each process will only mess with their own jobs.

Only alphanumeric and low dashes are allowed.

```bash
APP_NAME                # Application name (default: application1)
```

### 3.7 SSL Mode

At the moment we only support `require` and `disable`.

```bash
DB_SSLMode              # SSL mode for the database (default: disable)
```

If `require` is provided, the application expects SSL certificates in `/app` directory with the following names.

```bash
ls /app
# ca-cert.pem
# client-cert.pem
# client-key.pem
```

### 3.8 Client Cert Password

Use this environment to specify the password for the encryped client key while using SSL for DB connections.

```bash
DB_SSL_PASS             # SSL password for the database (default: clientpass)
```

### 3.9 Polling Interval

Use this environment to set how frequent you want the instance to check if there are jobs pending to be claimed.

```bash
POLL_INTERVAL_SECONDS   # Polling interval in seconds (default: 60)
```

### 3.10 Max Number of Jobs

Use this environment to set the maximum amount of jobs the instance can manage.
After setting this, it doesn't matter how many jobs are pending to be claimed, the instance will never exceed this limit.

```bash
MAX_JOBS                # Maximum number of jobs (default: 10000)
```

## 4. Job Configuration

If you are setting jobs for `ruok`, those need specific configurations.
Jobs are stored in `ruok.jobs` table and follows this structure.

```bash
# A human friendly name to identify the job
job_name

# The service endpoint to monitor
endpoint

# The HTTP method to use for monitoring (e.g., GET, POST)
httpmethod

# A JSON string containing headers for the HTTP request
headers

# An array of HTTP status codes indicating a successful response
successStatuses

# A cron expression specifying the job schedule
cronexp

# the channel used to alert in case the service is down ("http" at the moment)
alert_strategy

# the endpoint where alerts must be sent
alert_endpoint

# http method used to interact with the endpoint above
alert_method

# a JSON string containing the headers needed
alert_headers_string

# a string representing the payload to send
alert_payload
```

## 5. HTTP API

Each instance of ruok implements an http api to perform common operations.

The server listens on `http://localhost:8080` by default.

### 5.1 Create Jobs

```bash
# endpoint
POST /v1/jobs

# Example body
{
    "name:" "Service 1",
    "cronexp": "*/5 * * * * *",
    "maxRetries": 1,
    "endpoint": "http://example.com",
    "httpmethod": "GET",
    "headers": "",
    "successStatuses": [
        200,
        201
    ],
    "alertStrategy": "http",
    "alertMethod": "POST",
    "alertEndpoint": "http://alert.me/now",
    "alertPayload": "An error occurred with your endpoint",
    "alertHeaders": ""
}
```

### 5.2 Update Jobs

```bash
# endpoint
PUT /v1/jobs/:id

# path params
id --> the id of the job to update

# Example body
{
    "name:" "Service 1 - beta",
    "cronexp": "*/5 * * * * *",
    "maxRetries": 1,
    "endpoint": "http://example.com",
    "httpmethod": "GET",
    "headers": "",
    "successStatuses": [
        200,
        201
    ],
    "alertStrategy": "http",
    "alertMethod": "POST",
    "alertEndpoint": "http://alert.me/now",
    "alertPayload": "An error occurred with your endpoint",
    "alertHeaders": ""
}
```

### 5.3 List Jobs

```bash
# endpoint
GET /v1/jobs?limit=int&offset=int

# query params
limit  --> how many jobs should appear in the result
offset --> how many should skip
```

### 5.4 List Job Executions

```bash
# endpoint
GET /v1/jobs/:id?limit=int&offset=int

# path params
id --> the id of the related job

# query params
limit  --> how many jobs should appear in the result
offset --> how many should skip
```

### 5.5 Get Instance Info

```bash
# endpoint
GET /v1/instance
```

## 6. Cron Specification

RUOK Scheduler uses the [cron expression specification outlined in Wikipedia's CRON expression](https://en.wikipedia.org/wiki/Cron#CRON_expression). Behind the scenes, it leverages the [gorhill/cronexpr package](https://github.com/gorhill/cronexpr) for cron expression handling.

## 7. License

This software is released under [Apache License 2.0](./LICENSE.md). For more details, refer to the source code and documentation.

Feel free to explore the source code and adapt RUOK Scheduler to meet your specific monitoring needs. If you encounter any issues or have suggestions for improvement, please contribute to the project. Happy monitoring!

_RUOK Scheduler: Simple, Open, Reliable Service Monitoring._
