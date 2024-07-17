CREATE TABLE public.user_token (
	username varchar(256) NOT NULL,
	"token" varchar(256) NULL,
	CONSTRAINT user_token_pkey PRIMARY KEY (username)
);
CREATE INDEX user_token_token_idx ON public.user_token USING btree (token);