-- Drop tables in the correct order to avoid foreign key constraint issues
DROP TABLE IF EXISTS enrollments;
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS users;

-- Drop extensions 
-- DROP EXTENSION IF EXISTS postgis;
-- DROP EXTENSION IF EXISTS "uuid-ossp";
