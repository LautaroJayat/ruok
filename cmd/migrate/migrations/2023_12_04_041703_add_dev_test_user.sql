-- A role that allows to drop when testing
DO
$do$
BEGIN
   IF EXISTS (
      SELECT FROM pg_catalog.pg_roles
      WHERE rolname = 'RUOK_DELETE_TESTING_DATA') THEN
      RAISE NOTICE 'Role "RUOK_DELETE_TESTING_DATA" already exists. Skipping.';
   ELSE
      BEGIN   -- nested block
         CREATE ROLE RUOK_DELETE_TESTING_DATA WITH NOLOGIN;
      EXCEPTION
         WHEN duplicate_object THEN
            RAISE NOTICE 'Role "RUOK_DELETE_TESTING_DATA" was just created by a concurrent transaction. Skipping.';
      END;
   END IF;
END
$do$;

GRANT DELETE ON ruok.jobs to RUOK_DELETE_TESTING_DATA;
GRANT DELETE ON ruok.job_RESULTS to RUOK_DELETE_TESTING_DATA;


-- A role that allows to drop when testing
DO
$do$
BEGIN
   IF EXISTS (
      SELECT FROM pg_catalog.pg_roles
      WHERE rolname = 'testing_user') THEN
      RAISE NOTICE 'Role "TESTING_USER" already exists. Skipping.';
   ELSE
      BEGIN   -- nested block
        CREATE ROLE testing_user WITH LOGIN PASSWORD 'password';
      EXCEPTION
         WHEN duplicate_object THEN
            RAISE NOTICE 'Role "TESTING_USER" was just created by a concurrent transaction. Skipping.';
      END;
   END IF;
END
$do$;
-- A default role for testing
GRANT RUOK_DELETE_TESTING_DATA, RUOK_JOBS_MANAGER, RUOK_SCHEDULER_ROLE TO testing_user;