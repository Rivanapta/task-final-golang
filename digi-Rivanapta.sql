
-- DDL
CREATE TABLE public.accounts (
	account_id int8 GENERATED ALWAYS AS IDENTITY( INCREMENT BY 1 MINVALUE 1 MAXVALUE 9223372036854775807 START 1 CACHE 1 NO CYCLE) NOT NULL,
	"name" varchar NOT NULL,
	balance int8 NOT NULL,
	referral_account_id int8 NULL,
	CONSTRAINT account_id PRIMARY KEY (account_id),
	CONSTRAINT fk_referral_account FOREIGN KEY (referral_account_id) REFERENCES public.accounts(account_id)
);

CREATE TABLE public.auths (
	auth_id int8 GENERATED ALWAYS AS IDENTITY NOT NULL,
	account_id int8 NOT NULL,
	username varchar NOT NULL,
	"password" varchar NOT NULL,
	CONSTRAINT auth_id PRIMARY KEY (auth_id),
	CONSTRAINT auth_username UNIQUE (username),
	CONSTRAINT auths_unique UNIQUE (account_id)
);


CREATE TABLE public.transaction_categories (
	transaction_category_id int8 GENERATED ALWAYS AS IDENTITY NOT NULL,
	"name" varchar NOT NULL,
	CONSTRAINT transaction_categories_pk PRIMARY KEY (transaction_category_id)
);


CREATE TABLE public."transaction" (
	transaction_id int8 GENERATED ALWAYS AS IDENTITY NOT NULL,
	transaction_category_id int8 NULL,
	account_id int8 NULL,
	from_account_id int8 NULL,
	to_account_id int8 NULL,
	amount int8 NULL,
	transaction_date timestamp NULL,
	CONSTRAINT transaction_pk PRIMARY KEY (transaction_id),
	CONSTRAINT transaction_category_id FOREIGN KEY (transaction_category_id) REFERENCES public."transaction"(transaction_id)
);


