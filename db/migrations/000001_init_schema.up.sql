CREATE TABLE "accounts"
(
    "id" BIGSERIAL PRIMARY KEY,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT(now()),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT(now()),
    "owner" VARCHAR(255) NOT NULL,
    "balance" BIGINT NOT NULL,
    "currency" VARCHAR(255) NOT NULL
);

CREATE TABLE "entries"
(
    "id" BIGSERIAL PRIMARY KEY,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT(now()),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT(now()),
    "amount" BIGINT NOT NULL,
    "account_id" BIGINT NOT NULL
);

CREATE TABLE "payments"
(
    "id" BIGSERIAL PRIMARY KEY,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT(now()),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT(now()),
    "amount" BIGINT NOT NULL,
    "from_account_id" BIGINT NOT NULL,
    "to_account_id" BIGINT NOT NULL
);

ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "accounts"("id");
ALTER TABLE "payments" ADD FOREIGN KEY ("from_account_id") REFERENCES "accounts"("id");
ALTER TABLE "payments" ADD FOREIGN KEY ("to_account_id") REFERENCES "accounts"("id");

CREATE INDEX "accounts_index_0" ON "accounts" ("owner");
CREATE INDEX "entries_index_1" ON "entries" ("account_id");
CREATE INDEX "payments_index_2" ON "payments" ("from_account_id");
CREATE INDEX "payments_index_3" ON "payments" ("to_account_id");
CREATE INDEX "payments_index_4" ON "payments" ("from_account_id", "to_account_id");