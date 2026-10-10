-- merge duplicate shopping list items (same list, ingredient and unit) by summing their quantities
UPDATE shopping_list_items sli
SET quantity = d.total
FROM (
    SELECT MIN(shopping_list_item_id) AS keep_id, SUM(quantity) AS total
    FROM shopping_list_items
    GROUP BY shopping_list_id, ingredient_id, measurement_unit_id
    HAVING COUNT(*) > 1
) d
WHERE sli.shopping_list_item_id = d.keep_id;

DELETE FROM shopping_list_items sli
USING shopping_list_items keep
WHERE sli.shopping_list_id = keep.shopping_list_id
  AND sli.ingredient_id = keep.ingredient_id
  AND sli.measurement_unit_id = keep.measurement_unit_id
  AND sli.shopping_list_item_id > keep.shopping_list_item_id;

ALTER TABLE shopping_list_items
    ADD CONSTRAINT shopping_list_items_list_ingredient_unit_key
    UNIQUE (shopping_list_id, ingredient_id, measurement_unit_id);

DELETE FROM pantry_ingredients p
USING pantry_ingredients keep
WHERE p.ingredient_id = keep.ingredient_id
  AND p.pantry_ingredient_id > keep.pantry_ingredient_id;

ALTER TABLE pantry_ingredients
    ADD CONSTRAINT pantry_ingredients_ingredient_id_key UNIQUE (ingredient_id);
