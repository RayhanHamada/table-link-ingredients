/*
 Navicat Premium Data Transfer

 Target Server Type    : PostgreSQL
 Target Server Version : 160006 (160006)
 File Encoding         : 65001

 Date: 19/02/2025 16:23:53
*/


-- ----------------------------
-- Table structure for tm_ingredient
-- ----------------------------
DROP TABLE IF EXISTS "public"."tm_ingredient";
CREATE TABLE "public"."tm_ingredient" (
  "uuid" uuid NOT NULL,
  "name" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "cause_alergy" bool NOT NULL,
  "type" int4 NOT NULL,
  "status" int4 DEFAULT 0,
  "created_at" timestamp(6),
  "updated_at" timestamp(6),
  "deleted_at" timestamp(6)
)
;
ALTER TABLE "public"."tm_ingredient" OWNER TO "postgres";
COMMENT ON COLUMN "public"."tm_ingredient"."type" IS '0 (none), 1 (veggie), 2 (vegan)';
COMMENT ON COLUMN "public"."tm_ingredient"."status" IS '0 (inactive), 1 (active)';

-- ----------------------------
-- Records of tm_ingredient
-- ----------------------------
BEGIN;
INSERT INTO "public"."tm_ingredient" ("uuid", "name", "cause_alergy", "type", "status", "created_at", "updated_at", "deleted_at") VALUES ('8666d298-517b-45f3-8566-378cb5c8738c', 'Chicken', 'f', 0, 1, '2025-02-19 09:16:45.151435', NULL, NULL);
INSERT INTO "public"."tm_ingredient" ("uuid", "name", "cause_alergy", "type", "status", "created_at", "updated_at", "deleted_at") VALUES ('e97c6a5f-c541-4f3f-84d4-953c3eabe686', 'Pork', 'f', 0, 1, '2025-02-19 09:16:45.151435', NULL, NULL);
INSERT INTO "public"."tm_ingredient" ("uuid", "name", "cause_alergy", "type", "status", "created_at", "updated_at", "deleted_at") VALUES ('9a3a33cf-7144-4c5d-a0c6-fd8e894a0db5', 'Radish', 'f', 2, 1, '2025-02-19 09:16:45.151435', NULL, NULL);
INSERT INTO "public"."tm_ingredient" ("uuid", "name", "cause_alergy", "type", "status", "created_at", "updated_at", "deleted_at") VALUES ('b2752259-090e-4e6e-a9a1-2f47d538d833', 'Egg', 't', 1, 1, '2025-02-19 09:16:45.151435', NULL, NULL);
COMMIT;

-- ----------------------------
-- Table structure for tm_item
-- ----------------------------
DROP TABLE IF EXISTS "public"."tm_item";
CREATE TABLE "public"."tm_item" (
  "uuid" uuid NOT NULL,
  "name" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "price" numeric(10,2) NOT NULL,
  "status" int4 DEFAULT 0,
  "created_at" timestamp(6),
  "updated_at" timestamp(6),
  "deleted_at" timestamp(6)
)
;
ALTER TABLE "public"."tm_item" OWNER TO "postgres";

-- ----------------------------
-- Records of tm_item
-- ----------------------------
BEGIN;
INSERT INTO "public"."tm_item" ("uuid", "name", "price", "status", "created_at", "updated_at", "deleted_at") VALUES ('07419b87-4702-49f9-83aa-f9b489f64b14', 'Chicken Pork', 30000.00, 1, '2025-02-19 09:19:02.37464', NULL, NULL);
INSERT INTO "public"."tm_item" ("uuid", "name", "price", "status", "created_at", "updated_at", "deleted_at") VALUES ('d290fc2b-6a32-4bfe-98c7-25b9884c5245', 'Chicken Pork with Radish', 35000.00, 1, '2025-02-19 09:19:02.37464', NULL, NULL);
INSERT INTO "public"."tm_item" ("uuid", "name", "price", "status", "created_at", "updated_at", "deleted_at") VALUES ('7cc760a3-393b-493c-b780-3cfd7afd1cf9', 'Salad Egg', 20000.00, 1, '2025-02-19 09:19:02.37464', NULL, NULL);
COMMIT;

-- ----------------------------
-- Table structure for tm_item_ingredient
-- ----------------------------
DROP TABLE IF EXISTS "public"."tm_item_ingredient";
CREATE TABLE "public"."tm_item_ingredient" (
  "uuid_item" uuid NOT NULL,
  "uuid_ingredient" uuid NOT NULL
)
;
ALTER TABLE "public"."tm_item_ingredient" OWNER TO "postgres";

-- ----------------------------
-- Records of tm_item_ingredient
-- ----------------------------
BEGIN;
INSERT INTO "public"."tm_item_ingredient" ("uuid_item", "uuid_ingredient") VALUES ('07419b87-4702-49f9-83aa-f9b489f64b14', '8666d298-517b-45f3-8566-378cb5c8738c');
INSERT INTO "public"."tm_item_ingredient" ("uuid_item", "uuid_ingredient") VALUES ('07419b87-4702-49f9-83aa-f9b489f64b14', 'e97c6a5f-c541-4f3f-84d4-953c3eabe686');
INSERT INTO "public"."tm_item_ingredient" ("uuid_item", "uuid_ingredient") VALUES ('d290fc2b-6a32-4bfe-98c7-25b9884c5245', '8666d298-517b-45f3-8566-378cb5c8738c');
INSERT INTO "public"."tm_item_ingredient" ("uuid_item", "uuid_ingredient") VALUES ('d290fc2b-6a32-4bfe-98c7-25b9884c5245', 'e97c6a5f-c541-4f3f-84d4-953c3eabe686');
INSERT INTO "public"."tm_item_ingredient" ("uuid_item", "uuid_ingredient") VALUES ('d290fc2b-6a32-4bfe-98c7-25b9884c5245', '9a3a33cf-7144-4c5d-a0c6-fd8e894a0db5');
INSERT INTO "public"."tm_item_ingredient" ("uuid_item", "uuid_ingredient") VALUES ('7cc760a3-393b-493c-b780-3cfd7afd1cf9', '9a3a33cf-7144-4c5d-a0c6-fd8e894a0db5');
INSERT INTO "public"."tm_item_ingredient" ("uuid_item", "uuid_ingredient") VALUES ('7cc760a3-393b-493c-b780-3cfd7afd1cf9', 'b2752259-090e-4e6e-a9a1-2f47d538d833');
COMMIT;

-- ----------------------------
-- Function structure for uuid_generate_v1
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."uuid_generate_v1"();
CREATE OR REPLACE FUNCTION "public"."uuid_generate_v1"()
  RETURNS "pg_catalog"."uuid" AS '$libdir/uuid-ossp', 'uuid_generate_v1'
  LANGUAGE c VOLATILE STRICT
  COST 1;
ALTER FUNCTION "public"."uuid_generate_v1"() OWNER TO "postgres";

-- ----------------------------
-- Function structure for uuid_generate_v1mc
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."uuid_generate_v1mc"();
CREATE OR REPLACE FUNCTION "public"."uuid_generate_v1mc"()
  RETURNS "pg_catalog"."uuid" AS '$libdir/uuid-ossp', 'uuid_generate_v1mc'
  LANGUAGE c VOLATILE STRICT
  COST 1;
ALTER FUNCTION "public"."uuid_generate_v1mc"() OWNER TO "postgres";

-- ----------------------------
-- Function structure for uuid_generate_v3
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."uuid_generate_v3"("namespace" uuid, "name" text);
CREATE OR REPLACE FUNCTION "public"."uuid_generate_v3"("namespace" uuid, "name" text)
  RETURNS "pg_catalog"."uuid" AS '$libdir/uuid-ossp', 'uuid_generate_v3'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."uuid_generate_v3"("namespace" uuid, "name" text) OWNER TO "postgres";

-- ----------------------------
-- Function structure for uuid_generate_v4
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."uuid_generate_v4"();
CREATE OR REPLACE FUNCTION "public"."uuid_generate_v4"()
  RETURNS "pg_catalog"."uuid" AS '$libdir/uuid-ossp', 'uuid_generate_v4'
  LANGUAGE c VOLATILE STRICT
  COST 1;
ALTER FUNCTION "public"."uuid_generate_v4"() OWNER TO "postgres";

-- ----------------------------
-- Function structure for uuid_generate_v5
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."uuid_generate_v5"("namespace" uuid, "name" text);
CREATE OR REPLACE FUNCTION "public"."uuid_generate_v5"("namespace" uuid, "name" text)
  RETURNS "pg_catalog"."uuid" AS '$libdir/uuid-ossp', 'uuid_generate_v5'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."uuid_generate_v5"("namespace" uuid, "name" text) OWNER TO "postgres";

-- ----------------------------
-- Function structure for uuid_nil
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."uuid_nil"();
CREATE OR REPLACE FUNCTION "public"."uuid_nil"()
  RETURNS "pg_catalog"."uuid" AS '$libdir/uuid-ossp', 'uuid_nil'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."uuid_nil"() OWNER TO "postgres";

-- ----------------------------
-- Function structure for uuid_ns_dns
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."uuid_ns_dns"();
CREATE OR REPLACE FUNCTION "public"."uuid_ns_dns"()
  RETURNS "pg_catalog"."uuid" AS '$libdir/uuid-ossp', 'uuid_ns_dns'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."uuid_ns_dns"() OWNER TO "postgres";

-- ----------------------------
-- Function structure for uuid_ns_oid
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."uuid_ns_oid"();
CREATE OR REPLACE FUNCTION "public"."uuid_ns_oid"()
  RETURNS "pg_catalog"."uuid" AS '$libdir/uuid-ossp', 'uuid_ns_oid'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."uuid_ns_oid"() OWNER TO "postgres";

-- ----------------------------
-- Function structure for uuid_ns_url
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."uuid_ns_url"();
CREATE OR REPLACE FUNCTION "public"."uuid_ns_url"()
  RETURNS "pg_catalog"."uuid" AS '$libdir/uuid-ossp', 'uuid_ns_url'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."uuid_ns_url"() OWNER TO "postgres";

-- ----------------------------
-- Function structure for uuid_ns_x500
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."uuid_ns_x500"();
CREATE OR REPLACE FUNCTION "public"."uuid_ns_x500"()
  RETURNS "pg_catalog"."uuid" AS '$libdir/uuid-ossp', 'uuid_ns_x500'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."uuid_ns_x500"() OWNER TO "postgres";

-- ----------------------------
-- Primary Key structure for table tm_ingredient
-- ----------------------------
ALTER TABLE "public"."tm_ingredient" ADD CONSTRAINT "tm_ingredient_pkey" PRIMARY KEY ("uuid");

-- ----------------------------
-- Primary Key structure for table tm_item
-- ----------------------------
ALTER TABLE "public"."tm_item" ADD CONSTRAINT "tm_item_pkey" PRIMARY KEY ("uuid");

-- ----------------------------
-- Primary Key structure for table tm_item_ingredient
-- ----------------------------
ALTER TABLE "public"."tm_item_ingredient" ADD CONSTRAINT "tm_item_ingredient_pkey" PRIMARY KEY ("uuid_item", "uuid_ingredient");

-- ----------------------------
-- Foreign Keys structure for table tm_item_ingredient
-- ----------------------------
ALTER TABLE "public"."tm_item_ingredient" ADD CONSTRAINT "tm_item_ingredient_uuid_ingredient_fkey" FOREIGN KEY ("uuid_ingredient") REFERENCES "public"."tm_ingredient" ("uuid") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."tm_item_ingredient" ADD CONSTRAINT "tm_item_ingredient_uuid_item_fkey" FOREIGN KEY ("uuid_item") REFERENCES "public"."tm_item" ("uuid") ON DELETE NO ACTION ON UPDATE NO ACTION;
