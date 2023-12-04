-- Create role for scheduler
DO
$do$
BEGIN
   IF EXISTS (
      SELECT FROM pg_catalog.pg_roles
      WHERE rolname = 'RUOK_SCHEDULER_ROLE') THEN
      RAISE NOTICE 'Role "RUOK_SCHEDULER_ROLE" already exists. Skipping.';
   ELSE
      BEGIN   -- nested block
         CREATE ROLE RUOK_SCHEDULER_ROLE WITH NOLOGIN;
      EXCEPTION
         WHEN duplicate_object THEN
            RAISE NOTICE 'Role "RUOK_SCHEDULER_ROLE" was just created by a concurrent transaction. Skipping.';
      END;
   END IF;
END
$do$;

GRANT USAGE ON SCHEMA ruok to RUOK_SCHEDULER_ROLE;
GRANT SELECT,UPDATE ON ruok.jobs to RUOK_SCHEDULER_ROLE;
GRANT SELECT,INSERT ON ruok.job_results to RUOK_SCHEDULER_ROLE;
GRANT USAGE ON SEQUENCE ruok.job_results_id_seq to RUOK_SCHEDULER_ROLE;
GRANT USAGE ON SEQUENCE ruok.jobs_id_seq to RUOK_SCHEDULER_ROLE;
GRANT EXECUTE ON FUNCTION ruok.get_ssl_conn_version(text) to RUOK_SCHEDULER_ROLE;
GRANT EXECUTE ON FUNCTION ruok.micro_unix_now() to RUOK_SCHEDULER_ROLE;


DO
$do$
BEGIN
   IF EXISTS (
      SELECT FROM pg_catalog.pg_roles
      WHERE rolname = 'RUOK_JOBS_MANAGER') THEN
      RAISE NOTICE 'Role "RUOK_JOBS_MANAGER" already exists. Skipping.';
   ELSE
      BEGIN   -- nested block
         CREATE ROLE RUOK_JOBS_MANAGER WITH NOLOGIN;
      EXCEPTION
         WHEN duplicate_object THEN
            RAISE NOTICE 'Role "RUOK_JOBS_MANAGER" was just created by a concurrent transaction. Skipping.';
      END;
   END IF;
END
$do$;

GRANT SELECT,UPDATE,INSERT,DELETE ON ruok.jobs to RUOK_JOBS_MANAGER;
GRANT SELECT,INSERT ON ruok.job_results to RUOK_JOBS_MANAGER;
GRANT USAGE ON SEQUENCE ruok.job_results_id_seq to RUOK_JOBS_MANAGER;
GRANT USAGE ON SEQUENCE ruok.jobs_id_seq to RUOK_JOBS_MANAGER;
GRANT EXECUTE ON FUNCTION ruok.get_ssl_conn_version(text) to RUOK_JOBS_MANAGER;
GRANT EXECUTE ON FUNCTION ruok.micro_unix_now() to RUOK_JOBS_MANAGER;
GRANT EXECUTE ON FUNCTION pg_notify(text, text) to RUOK_JOBS_MANAGER;



