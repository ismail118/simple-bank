CREATE TABLE "users" (
     "username" varchar PRIMARY KEY,
     "hashed_password" varchar NOT NULL,
     "full_name" varchar NOT NULL,
     "email" varchar UNIQUE NOT NULL,
     "created_at" timestamptz NOT NULL DEFAULT (now()),
     "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "accounts" ("owner", "currency");

ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username") ON DELETE CASCADE;

-- CREATE UNIQUE INDEX ON "accounts" ("owner", "currency");
ALTER TABLE "accounts" ADD CONSTRAINT "owner_currency_key" UNIQUE ("owner", "currency");