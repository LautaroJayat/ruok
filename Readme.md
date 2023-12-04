# RUOK Scheduler

<div align="center">
    <img width="400" src="./assets/big_ruok_logo.svg" alt="RUOK Logo" />
</div>

RUOK Scheduler is an open-source tool designed for hassle-free service monitoring. Keep a close eye on your infrastructure effortlessly with our intuitive scheduler.

## Introduction

### Purpose

RUOK Scheduler serves the purpose of transforming a PostgreSQL database into a reliable and efficient broker for a backend service monitoring system. It simplifies the process of scheduling and monitoring services, offering a straightforward solution for scenarios where complex deployments are unnecessary.

### Why RUOK Scheduler?

In many cases, deploying and managing monitoring systems can be complex and resource-intensive. RUOK Scheduler aims to address this challenge by providing a simple, yet effective, solution for users who prioritize ease of use and minimal configuration overhead.

## Getting Started

### Building from Source

To build RUOK Scheduler from source, use the following command:

```bash
make build
```

### Starting RUOK Scheduler

After building, run the scheduler using:

```bash
./ruok
```

### Preparing the Database

Before running the scheduler, migrate the database with:

```bash
./ruok -migrate
```

This command sets up the necessary schema for tables and functions in the PostgreSQL database.

### Configurations

All configurations for RUOK Scheduler are expected as environment variables. Below are the configurations along with their respective environment variables:

```bash
STORAGE_KIND: Storage type (default: postgres)
DB_PROTOCOL: Database protocol (default: postgresql)
DB_PASS: Database password (default: password)
DB_USER: Database user (default: user)
DB_HOST: Database host (default: localhost)
DB_PORT: Database port (default: 5432)
DB_NAME: Database name (default: db1)
APP_NAME: Application name (default: application1)
DB_SSLMode: SSL mode for the database (default: disable)
DB_SSL_PASS: SSL password for the database (default: clientpass)
POLL_INTERVAL_SECONDS: Polling interval in seconds (default: 60)
MAX_JOBS: Maximum number of jobs (default: 10000)
```

### Job Configuration

Jobs in RUOK Scheduler require specific configurations. Each job should have the following parameters:

```bash
endpoint: The service endpoint to monitor
httpmethod: The HTTP method to use for monitoring (e.g., GET, POST)
headers: A JSON string containing headers for the HTTP request
successStatuses: An array of HTTP status codes indicating a successful response
cronexp: A cron expression specifying the job schedule
```

### Cron Specification

RUOK Scheduler uses the [cron expression specification outlined in Wikipedia's CRON expression](https://en.wikipedia.org/wiki/Cron#CRON_expression). Behind the scenes, it leverages the [gorhill/cronexpr package](https://github.com/gorhill/cronexpr) for cron expression handling.

## License

This software is released under the MIT License. For more details, refer to the source code and documentation.

Feel free to explore the source code and adapt RUOK Scheduler to meet your specific monitoring needs. If you encounter any issues or have suggestions for improvement, please contribute to the project. Happy monitoring!

_RUOK Scheduler: Simple, Open, Reliable Service Monitoring._
