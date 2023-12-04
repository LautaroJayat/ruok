-- Create schema for RUOK
CREATE SCHEMA IF NOT EXISTS ruok;

-- Create function to get time in unix microseconds
CREATE OR REPLACE FUNCTION ruok.micro_unix_now() RETURNS BIGINT AS 
$$
    BEGIN 
	    RETURN (SELECT (EXTRACT(epoch FROM now()) * 1000)::bigint);
    END;
$$
LANGUAGE plpgsql;

-- Create function to get SSL status and version of application
CREATE OR REPLACE FUNCTION ruok.get_ssl_conn_version(app_name text)
RETURNS TABLE(ssl_active boolean, ssl_version text)	AS
$$
    SELECT ssl,version FROM pg_stat_ssl
    JOIN pg_stat_activity
	ON pg_stat_ssl.pid = pg_stat_activity.pid
	    WHERE application_name = $1
	LIMIT 1
$$
LANGUAGE SQL;
