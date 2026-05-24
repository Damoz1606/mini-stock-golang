schema "public" {
  comment = "MiniStock public schema"
}

table "categories" {
  schema = schema.public
  column "id" {
    type    = uuid
    null    = false
    default = sql("gen_random_uuid()")
  }
  column "name" {
    type = text
    null = false
  }
  column "createdAt" {
    type = bigint
    null = false
  }
  column "updatedAt" {
    type = bigint
    null = false
  }
  primary_key {
    columns = [column.id]
  }
  unique "uq_categories_name" {
    columns = [column.name]
  }
}

table "attributes" {
  schema = schema.public
  column "id" {
    type    = uuid
    null    = false
    default = sql("gen_random_uuid()")
  }
  column "categoryId" {
    type = uuid
    null = false
  }
  column "name" {
    type = text
    null = false
  }
  column "description" {
    type = text
    null = true
  }
  column "createdAt" {
    type = bigint
    null = false
  }
  column "updatedAt" {
    type = bigint
    null = false
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_attributes_category" {
    columns     = [column.categoryId]
    ref_columns = [table.categories.column.id]
    on_delete   = CASCADE
  }
  unique "uq_attributes_category_name" {
    columns = [column.categoryId, column.name]
  }
}

table "items" {
  schema = schema.public
  column "id" {
    type    = uuid
    null    = false
    default = sql("gen_random_uuid()")
  }
  column "categoryId" {
    type = uuid
    null = false
  }
  column "name" {
    type = text
    null = false
  }
  column "description" {
    type = text
    null = true
  }
  column "createdAt" {
    type = bigint
    null = false
  }
  column "updatedAt" {
    type = bigint
    null = false
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_items_category" {
    columns     = [column.categoryId]
    ref_columns = [table.categories.column.id]
    on_delete   = CASCADE
  }
  unique "uq_items_category_name" {
    columns = [column.categoryId, column.name]
  }
}

table "attributes_items" {
  schema = schema.public
  column "categoryId" {
    type = uuid
    null = false
  }
  column "attributeId" {
    type = uuid
    null = false
  }
  column "itemId" {
    type = uuid
    null = false
  }
  primary_key {
    columns = [column.categoryId, column.attributeId, column.itemId]
  }
  foreign_key "fk_attributes_items_category" {
    columns     = [column.categoryId]
    ref_columns = [table.categories.column.id]
    on_delete   = CASCADE
  }
  foreign_key "fk_attributes_items_attribute" {
    columns     = [column.attributeId]
    ref_columns = [table.attributes.column.id]
    on_delete   = CASCADE
  }
  foreign_key "fk_attributes_items_item" {
    columns     = [column.itemId]
    ref_columns = [table.items.column.id]
    on_delete   = CASCADE
  }
  unique "uq_attributes_items_attribute_item" {
    columns = [column.attributeId, column.itemId]
  }
}
