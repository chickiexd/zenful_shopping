SELECT setval('ingredients_ingredient_id_seq', (SELECT MAX(ingredient_id) FROM ingredients) + 1);

SELECT setval('measurement_units_measurement_unit_id_seq', (SELECT MAX(measurement_unit_id) FROM measurement_units) + 1);

SELECT setval('food_groups_food_group_id_seq', (SELECT MAX(food_group_id) FROM food_groups) + 1);

SELECT setval('meal_types_meal_type_id_seq', (SELECT MAX(meal_type_id) FROM meal_types) + 1);

