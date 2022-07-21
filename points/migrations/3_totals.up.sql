CREATE EXTENSION hstore;
CREATE TABLE customer_totals (customer hstore);
CREATE INDEX customer_totals_idx ON customer_totals USING GIN (customer);
