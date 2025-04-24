CREATE TABLE shopping_lists (
    shopping_list_id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE shopping_list_items (
    shopping_list_item_id SERIAL PRIMARY KEY,
    shopping_list_id INTEGER NOT NULL REFERENCES shopping_lists(shopping_list_id),
    ingredient_id INTEGER NOT NULL REFERENCES ingredients(ingredient_id),
    measurement_unit_id INTEGER NOT NULL REFERENCES measurement_units(measurement_unit_id),
    quantity DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE shopping_list_food_groups (
    shopping_list_id INTEGER NOT NULL REFERENCES shopping_lists(shopping_list_id),
    food_group_id INTEGER NOT NULL REFERENCES food_groups(food_group_id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (shopping_list_id, food_group_id)
);

-- INSERT INTO shopping_lists (name)
-- VALUES
--   ('Produce'),
--   ('Chilled'),
--   ('Frozen'),
--   ('Other');
--
-- INSERT INTO shopping_list_food_groups (shopping_list_id, food_group_id)
-- SELECT sl.shopping_list_id, fg.food_group_id
-- FROM shopping_lists sl, food_groups fg
-- WHERE
--   (sl.name = 'Produce' AND fg.name IN ('Vegetables', 'Fruits', 'Herbs')) OR
--   (sl.name = 'Chilled' AND fg.name IN ('Dairy', 'Meat', 'Fish')) OR
--   (sl.name = 'Frozen' AND fg.name IN ('Ice Cream', 'Frozen Vegetables')) OR
--   (sl.name = 'Other' AND fg.name IN ('Condiments', 'Spices', 'Baking', 'Canned Goods', 'Oils'));
