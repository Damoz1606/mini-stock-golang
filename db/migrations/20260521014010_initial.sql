-- Create "types" table
CREATE TABLE "public"."types" (
  "id" uuid NOT NULL,
  "name" character varying(255) NOT NULL,
  PRIMARY KEY ("id")
);
-- Create "features" table
CREATE TABLE "public"."features" (
  "id" uuid NOT NULL,
  "type_id" uuid NOT NULL,
  "feature" character varying(255) NOT NULL,
  "description" text NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "features_type_id_fkey" FOREIGN KEY ("type_id") REFERENCES "public"."types" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_features_type_id" to table: "features"
CREATE INDEX "idx_features_type_id" ON "public"."features" ("type_id");
-- Create "items" table
CREATE TABLE "public"."items" (
  "id" uuid NOT NULL,
  "type_id" uuid NOT NULL,
  "name" character varying(255) NOT NULL,
  "description" text NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "items_type_id_fkey" FOREIGN KEY ("type_id") REFERENCES "public"."types" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_items_type_id" to table: "items"
CREATE INDEX "idx_items_type_id" ON "public"."items" ("type_id");
-- Create "features_items" table
CREATE TABLE "public"."features_items" (
  "feature_id" uuid NOT NULL,
  "item_id" uuid NOT NULL,
  CONSTRAINT "features_items_feature_id_item_id_key" UNIQUE ("feature_id", "item_id"),
  CONSTRAINT "features_items_feature_id_fkey" FOREIGN KEY ("feature_id") REFERENCES "public"."features" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "features_items_item_id_fkey" FOREIGN KEY ("item_id") REFERENCES "public"."items" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_features_items_feature_id" to table: "features_items"
CREATE INDEX "idx_features_items_feature_id" ON "public"."features_items" ("feature_id");
-- Create index "idx_features_items_item_id" to table: "features_items"
CREATE INDEX "idx_features_items_item_id" ON "public"."features_items" ("item_id");
