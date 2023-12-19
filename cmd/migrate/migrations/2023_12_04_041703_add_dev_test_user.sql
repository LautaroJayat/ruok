-- A role that allows to drop when testing
DO
$do$
BEGIN
   IF EXISTS (
      SELECT FROM pg_catalog.pg_roles
      WHERE rolname = 'RUOK_SEED_AND_DROP') THEN
      RAISE NOTICE 'Role "RUOK_SEED_AND_DROP" already exists. Skipping.';
   ELSE
      BEGIN   -- nested block
         CREATE ROLE RUOK_SEED_AND_DROP WITH NOLOGIN;
      EXCEPTION
         WHEN duplicate_object THEN
            RAISE NOTICE 'Role "RUOK_SEED_AND_DROP" was just created by a concurrent transaction. Skipping.';
      END;
   END IF;
END
$do$;

GRANT INSERT,DELETE ON ruok.jobs to RUOK_SEED_AND_DROP;
GRANT INSERT,DELETE ON ruok.job_RESULTS to RUOK_SEED_AND_DROP;


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
-- A default role for dropping, seeding, and so on
GRANT USAGE on SCHEMA ruok to RUOK_SEED_AND_DROP;
GRANT RUOK_SEED_AND_DROP TO testing_user;

DROP POLICY IF EXISTS testing_user_delete_job_results ON ruok.job_results;
CREATE POLICY testing_user_delete_job_results ON ruok.job_results FOR DELETE TO RUOK_SEED_AND_DROP USING (true);
DROP POLICY IF EXISTS testing_user_insert_job_results ON ruok.job_results;
CREATE POLICY testing_user_insert_job_results ON ruok.job_results FOR INSERT TO RUOK_SEED_AND_DROP WITH CHECK (true);

DROP POLICY IF EXISTS testing_user_delete_jobs ON ruok.jobs;
CREATE POLICY testing_user_delete_jobs ON ruok.jobs FOR DELETE TO RUOK_SEED_AND_DROP USING (true);
DROP POLICY IF EXISTS testing_user_insert_jobs ON ruok.jobs;
CREATE POLICY testing_user_insert_jobs ON ruok.jobs FOR INSERT TO RUOK_SEED_AND_DROP WITH CHECK (true);


-- A role to login as application1
DO
$do$
BEGIN
   IF EXISTS (
      SELECT FROM pg_catalog.pg_roles
      WHERE rolname = 'application1') THEN
      RAISE NOTICE 'Role "application1" already exists. Skipping.';
   ELSE
      BEGIN   -- nested block
        CREATE ROLE application1 WITH LOGIN PASSWORD 'password';
      EXCEPTION
         WHEN duplicate_object THEN
            RAISE NOTICE 'Role "application1" was just created by a concurrent transaction. Skipping.';
      END;
   END IF;
END
$do$;
GRANT RUOK_SCHEDULER_ROLE TO application1;


-- A role to login as job_manager
DO
$do$
BEGIN
   IF EXISTS (
      SELECT FROM pg_catalog.pg_roles
      WHERE rolname = 'job_manager_1') THEN
      RAISE NOTICE 'Role "application1" already exists. Skipping.';
   ELSE
      BEGIN   -- nested block
        CREATE ROLE job_manager_1 WITH LOGIN PASSWORD 'password';
      EXCEPTION
         WHEN duplicate_object THEN
            RAISE NOTICE 'Role "job_manager_1" was just created by a concurrent transaction. Skipping.';
      END;
   END IF;
END
$do$;
GRANT RUOK_SCHEDULER_ROLE TO job_manager_1;
