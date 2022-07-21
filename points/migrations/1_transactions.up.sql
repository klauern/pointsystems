
CREATE TABLE transactions (
    id BIGSERIAL PRIMARY KEY,
    customer integer NOT NULL,
    payer TEXT NOT NULL,
    amount integer NOT NULL,
    entered_time timestamp with time zone NOT NULL
);

CREATE TABLE customers (
    id BIGSERIAL PRIMARY KEY,
    name text NOT NULL
);

CREATE index customer_id_idx on transactions (customer);

CREATE VIEW customer_transactions AS
    SELECT customers.id as customer_id, name, customer, payer, amount, entered_time
    FROM customers, transactions
    WHERE customers.id = transactions.customer;