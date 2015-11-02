
-- NEED TO CREATE A DIFFERENT POSTGRES USER FRIST


CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;

CREATE TABLE menuitems (
  menuitemid uuid default uuid_generate_v4() NOT NULL,
  restaurantid uuid NOT NULL,
  name text NOT NULL,
  description text,
  price integer,
  imageurl text,
  created timestamp default transaction_timestamp() NOT NULL
);
