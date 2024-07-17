CREATE TYPE transaction_type AS ENUM ('CREDIT', 'DEBIT');

CREATE TABLE transactions (
	id bigserial primary key,
  username varchar(256),
  second_party varchar(256),
  amount int,
  "type" transaction_type
);
CREATE INDEX transactions_username_amount_idx ON public.transactions USING BTREE (username, amount);
