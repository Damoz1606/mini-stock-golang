schema "public" {}

table "types" {
  schema = schema.public
  column "id" {
    type = uuid
    null = false
  }
  column "name" {
    type = varchar(255)
    null = false
  }
  primary_key {
    columns = [column.id]
  }
}

table "features" {
  schema = schema.public
  column "id" {
    type = uuid
    null = false
  }
  column "type_id" {
    type = uuid
    null = false
  }
  column "feature" {
    type = varchar(255)
    null = false
  }
  column "description" {
    type = text
    null = true
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "features_type_id_fkey" {
    columns     = [column.type_id]
    ref_columns = [table.types.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  index "idx_features_type_id" {
    columns = [column.type_id]
  }
}

table "items" {
  schema = schema.public
  column "id" {
    type = uuid
    null = false
  }
  column "type_id" {
    type = uuid
    null = false
  }
  column "name" {
    type = varchar(255)
    null = false
  }
  column "description" {
    type = text
    null = true
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "items_type_id_fkey" {
    columns     = [column.type_id]
    ref_columns = [table.types.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  index "idx_items_type_id" {
    columns = [column.type_id]
  }
}

table "features_items" {
  schema = schema.public
  column "feature_id" {
    type = uuid
    null = false
  }
  column "item_id" {
    type = uuid
    null = false
  }
  foreign_key "features_items_feature_id_fkey" {
    columns     = [column.feature_id]
    ref_columns = [table.features.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  foreign_key "features_items_item_id_fkey" {
    columns     = [column.item_id]
    ref_columns = [table.items.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  index "idx_features_items_feature_id" {
    columns = [column.feature_id]
  }
  index "idx_features_items_item_id" {
    columns = [column.item_id]
  }
  unique "features_items_feature_id_item_id_key" {
    columns = [column.feature_id, column.item_id]
  }
}
