CREATE TABLE user_total_debit (
  username varchar(256) primary key,
  total_debit bigint
);
CREATE INDEX user_total_debit_idx ON public.user_total_debit USING BTREE (total_debit);
