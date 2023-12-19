CREATE ROLE admin;  -- Administrator

-- Create table for Jobs
CREATE TABLE IF NOT EXISTS ruok.jobs (
	id bigserial PRIMARY KEY,
	job_name text NOT NULL,
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

ALTER TABLE ruok.jobs ENABLE ROW LEVEL SECURITY;
CREATE POLICY admin_all_jobs ON ruok.jobs TO admin USING (true) WITH CHECK (true);


-- Create table for Job Executions
CREATE TABLE IF NOT EXISTS ruok.job_results (
	id uuid PRIMARY KEY NOT NULL,
	job_id uuid NOT NULL,
	job_name text NOT NULL,
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

ALTER TABLE ruok.job_results ENABLE ROW LEVEL SECURITY;
CREATE POLICY admin_all_job_results ON ruok.jobs TO admin USING (true) WITH CHECK (true);
