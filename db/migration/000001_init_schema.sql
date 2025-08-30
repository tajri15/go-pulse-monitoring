CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password_hash" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "sites" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "url" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "health_checks" (
  "id" bigserial PRIMARY KEY,
  "site_id" bigint NOT NULL,
  "status_code" int,
  "response_time_ms" int,
  "is_up" boolean NOT NULL,
  "checked_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "sites" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "health_checks" ADD FOREIGN KEY ("site_id") REFERENCES "sites" ("id");

CREATE INDEX ON "sites" ("user_id");
CREATE INDEX ON "health_checks" ("site_id");