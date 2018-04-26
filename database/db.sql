CREATE EXTENSION "pgcrypto";
CREATE EXTENSION "uuid-ossp";
SET timezone = 'UTC';

CREATE TABLE transaction_log(
	id uuid DEFAULT uuid_generate_v4()  NOT NULL PRIMARY KEY,
	source jsonb,
	txtype text,
	txid text,
	data jsonb,
	created_at timestamp with time zone NOT NULL default now(),
	unique(txtype,txid)
	);
CREATE INDEX index_transaction_log_id on transaction_log(id);
CREATE INDEX index_transaction_log_type on transaction_log(txtype);
CREATE INDEX index_transaction_log_txid on transaction_log(txid);