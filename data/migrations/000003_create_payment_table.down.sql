DROP TABLE IF EXISTS "payments";
DROP TABLE IF EXISTS "payment_details";

ALTER TABLE "orders"
DROP COLUMN payment_detail_id,
DROP COLUMN price;
