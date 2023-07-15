drop table if exists verify_email CASCADE;

ALTER TABLE "users" DROP COLUMN "is_email_verify";