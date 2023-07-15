CREATE TABLE "verify_email" (
    "id" bigserial PRIMARY KEY,
    "username" varchar NOT NULL,
    "email" varchar NOT NULL,
    "secret_code" varchar NOT NULL,
    "is_used" boolean NOT NULL DEFAULT 'false',
    "created_at" timestamptz NOT NULL DEFAULT 'now()',
    "expired_at" timestamptz NOT NULL DEFAULT 'now()'
);

ALTER TABLE "verify_email" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "users" ADD COLUMN "is_email_verify" boolean NOT NULL DEFAULT 'false';