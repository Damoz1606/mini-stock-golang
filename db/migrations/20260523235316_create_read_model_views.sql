CREATE VIEW v_categories AS 
SELECT id, name FROM categories;

CREATE VIEW v_attributes AS 
SELECT id, "categoryId", name, description FROM attributes;

CREATE VIEW v_items AS 
SELECT id, "categoryId", name, description FROM items;
