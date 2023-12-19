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

GRANT SELECT,INSERT,UPDATE ON ruok.jobs to RUOK_SCHEDULER_ROLE;
GRANT SELECT,INSERT ON ruok.job_results to RUOK_SCHEDULER_ROLE;

GRANT EXECUTE ON FUNCTION ruok.get_ssl_conn_version(text) to RUOK_SCHEDULER_ROLE;
GRANT EXECUTE ON FUNCTION ruok.micro_unix_now() to RUOK_SCHEDULER_ROLE;

DROP POLICY IF EXISTS scheduler_select_jobs ON ruok.jobs; 
CREATE POLICY scheduler_select_jobs ON ruok.jobs FOR SELECT TO RUOK_SCHEDULER_ROLE USING (
	(
      claimed_by IS null 
      AND 
      jobs.status = 'pending to be claimed'
   )
   OR	claimed_by = current_setting('application_name')
);

DROP POLICY IF EXISTS scheduler_update_jobs ON ruok.jobs; 
CREATE POLICY scheduler_update_jobs ON ruok.jobs FOR UPDATE TO RUOK_SCHEDULER_ROLE USING (
	(
      claimed_by IS null 
      AND 
      jobs.status = 'pending to be claimed'
   )
   OR	claimed_by = current_setting('application_name')
);

DROP POLICY IF EXISTS scheduler_insert_jobs ON ruok.jobs; 
CREATE POLICY scheduler_insert_jobs ON ruok.jobs FOR INSERT TO RUOK_SCHEDULER_ROLE WITH CHECK (true);

DROP POLICY IF EXISTS scheduler_insert_job_results ON ruok.job_results; 
CREATE POLICY scheduler_insert_job_results ON ruok.job_results FOR INSERT TO RUOK_SCHEDULER_ROLE WITH CHECK (true);

DROP POLICY IF EXISTS scheduler_select_job_results ON ruok.job_results; 
CREATE POLICY scheduler_select_job_results ON ruok.job_results FOR SELECT TO RUOK_SCHEDULER_ROLE USING (
	claimed_by = current_setting('application_name')
);


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

GRANT EXECUTE ON FUNCTION ruok.get_ssl_conn_version(text) to RUOK_JOBS_MANAGER;
GRANT EXECUTE ON FUNCTION ruok.micro_unix_now() to RUOK_JOBS_MANAGER;
GRANT EXECUTE ON FUNCTION pg_notify(text, text) to RUOK_JOBS_MANAGER;

DROP POLICY IF EXISTS jobs_manager_insert_jobs ON ruok.jobs;
CREATE POLICY jobs_manager_insert_jobs ON ruok.jobs FOR INSERT TO RUOK_JOBS_MANAGER WITH CHECK (true);

DROP POLICY IF EXISTS jobs_manager_select_jobs ON ruok.jobs;
CREATE POLICY jobs_manager_select_jobs ON ruok.jobs FOR SELECT TO RUOK_JOBS_MANAGER USING (true);

DROP POLICY IF EXISTS jobs_manager_update_jobs ON ruok.jobs;
CREATE POLICY jobs_manager_update_jobs ON ruok.jobs FOR UPDATE TO RUOK_JOBS_MANAGER USING (true);

DROP POLICY IF EXISTS jobs_manager_select_job_results ON ruok.job_results;
CREATE POLICY jobs_manager_select_job_results ON ruok.job_results FOR SELECT TO RUOK_JOBS_MANAGER USING (true);






