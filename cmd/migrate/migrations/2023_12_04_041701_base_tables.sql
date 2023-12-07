
-- Create table for Jobs
CREATE TABLE IF NOT EXISTS ruok.jobs (
	id bigserial PRIMARY KEY,
	cron_exp_string text NOT NULL,
	endpoint text NOT NULL,
	httpmethod text NOT NULL,
	max_retries smallint DEFAULT 1,
	last_execution bigint,
	should_execute_at bigint,
	last_response_at bigint,
	last_message text,
	last_status_code int,
	headers_string text,
	success_statuses int[] NOT NULL,
	succeeded text,
	tls_client_cert text,
	alert_strategy text,
	alert_endpoint text,
	alert_method text,
	alert_headers_string text,
	alert_payload text,
	status text,
	claimed_by text,
	created_at bigint DEFAULT ruok.micro_unix_now() NOT NULL,
	updated_at bigint,
	deleted_at bigint
);

-- Create table for Job Executions
CREATE TABLE IF NOT EXISTS ruok.job_results (
	id bigserial PRIMARY KEY,
	job_id bigint,
	cron_exp_string varchar NOT NULL,
	endpoint varchar NOT NULL,
	httpmethod varchar NOT NULL,
	max_retries smallint DEFAULT 1,
	execution_time bigint NOT NULL,
	should_execute_at bigint NOT NULL,
	last_response_at bigint NOT NULL,
	last_message varchar,
	last_status_code int NOT NULL,
	success_statuses int[] NOT NULL,
	succeeded text NOT NULL,
	tls_client_cert varchar,
	status varchar NOT NULL,
	claimed_by varchar NOT NULL,
	created_at bigint DEFAULT ruok.micro_unix_now() NOT NULL,
	deleted_at bigint
);



-- Create role for scheduler
CREATE ROLE RUOK_SCHEDULER_ROLE WITH NOLOGIN;
GRANT USAGE ON SCHEMA ruok to RUOK_SCHEDULER_ROLE;
GRANT SELECT,UPDATE ON ruok.jobs to RUOK_SCHEDULER_ROLE;
GRANT SELECT,INSERT ON ruok.job_results to RUOK_SCHEDULER_ROLE;
GRANT EXECUTE ON FUNCTION ruok.get_ssl_conn_version(text) to RUOK_SCHEDULER_ROLE;
GRANT EXECUTE ON FUNCTION ruok.micro_unix_now() to RUOK_SCHEDULER_ROLE;

CREATE ROLE RUOK_JOBS_MANAGER WITH NOLOGIN;
GRANT SELECT,UPDATE,INSERT,DELETE ON ruok.jobs to RUOK_JOBS_MANAGER;
GRANT SELECT,INSERT ON ruok.job_results to RUOK_JOBS_MANAGER;
GRANT EXECUTE ON FUNCTION ruok.get_ssl_conn_version(text) to RUOK_JOBS_MANAGER;
GRANT EXECUTE ON FUNCTION ruok.micro_unix_now() to RUOK_JOBS_MANAGER;
GRANT EXECUTE ON FUNCTION pg_notify(text, text) to RUOK_JOBS_MANAGER;



