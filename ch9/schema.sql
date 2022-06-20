CREATE TABLE IF NOT EXISTS users(
  user_id varchar(32) NOT NULL,
  user_name varchar(100) NOT NULL,
  created_at timestamp with time zone,
  CONSTRAINT pk_users PRIMARY KEY (user_id));
