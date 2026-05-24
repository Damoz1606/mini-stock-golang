-- Set comment to schema: "public"
COMMENT ON SCHEMA "public" IS 'MiniStock public schema';
-- Create "categories" table
CREATE TABLE "public"."categories" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" text NOT NULL,
  "createdAt" bigint NOT NULL,
  "updatedAt" bigint NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uq_categories_name" UNIQUE ("name")
);
-- Create "attributes" table
CREATE TABLE "public"."attributes" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "categoryId" uuid NOT NULL,
  "name" text NOT NULL,
  "description" text NULL,
  "createdAt" bigint NOT NULL,
  "updatedAt" bigint NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uq_attributes_category_name" UNIQUE ("categoryId", "name"),
  CONSTRAINT "fk_attributes_category" FOREIGN KEY ("categoryId") REFERENCES "public"."categories" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create "items" table
CREATE TABLE "public"."items" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "categoryId" uuid NOT NULL,
  "name" text NOT NULL,
  "description" text NULL,
  "createdAt" bigint NOT NULL,
  "updatedAt" bigint NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uq_items_category_name" UNIQUE ("categoryId", "name"),
  CONSTRAINT "fk_items_category" FOREIGN KEY ("categoryId") REFERENCES "public"."categories" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create "attributes_items" table
CREATE TABLE "public"."attributes_items" (
  "categoryId" uuid NOT NULL,
  "attributeId" uuid NOT NULL,
  "itemId" uuid NOT NULL,
  PRIMARY KEY ("categoryId", "attributeId", "itemId"),
  CONSTRAINT "uq_attributes_items_attribute_item" UNIQUE ("attributeId", "itemId"),
  CONSTRAINT "fk_attributes_items_attribute" FOREIGN KEY ("attributeId") REFERENCES "public"."attributes" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_attributes_items_category" FOREIGN KEY ("categoryId") REFERENCES "public"."categories" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_attributes_items_item" FOREIGN KEY ("itemId") REFERENCES "public"."items" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
