INSERT INTO food_groups (food_group_id, name, description)
VALUES
  (1, 'Vegetables', 'Fresh vegetables like carrots, lettuce, and peppers'),
  (2, 'Fruits', 'Fresh fruits like apples, bananas, and berries'),
  (3, 'Herbs', 'Fresh herbs like basil, parsley, and cilantro'),
  (4, 'Dairy', 'Milk, yogurt, butter, and other dairy products'),
  (5, 'Meat', 'Fresh meat such as beef, pork, and chicken'),
  (6, 'Fish', 'Fresh or packaged seafood like salmon and tuna'),
  (7, 'Ice Cream', 'Frozen desserts'),
  (8, 'Frozen Vegetables', 'Pre-packaged frozen vegetables'),
  (9, 'Condiments', 'Items like ketchup, mustard, mayonnaise'),
  (10, 'Spices', 'Dry spices such as pepper, cumin, and paprika'),
  (11, 'Baking', 'Flour, sugar, baking powder, and similar items'),
  (12, 'Canned Good', 'Canned beans, tomatoes, soups, etc.'),
  (13, 'Oils', 'Cooking oils like olive oil, vegetable oil');

INSERT INTO shopping_lists (name)
VALUES
  ('Produce'),
  ('Chilled'),
  ('Frozen'),
  ('Other');

INSERT INTO shopping_list_food_groups (shopping_list_id, food_group_id)
SELECT sl.shopping_list_id, fg.food_group_id
FROM shopping_lists sl, food_groups fg
WHERE
  (sl.name = 'Produce' AND fg.name IN ('Vegetables', 'Fruits', 'Herbs')) OR
  (sl.name = 'Chilled' AND fg.name IN ('Dairy', 'Meat', 'Fish')) OR
  (sl.name = 'Frozen' AND fg.name IN ('Ice Cream', 'Frozen Vegetables')) OR
  (sl.name = 'Other' AND fg.name IN ('Condiments', 'Spices', 'Baking', 'Canned Goods', 'Oils'));

INSERT INTO measurement_units (measurement_unit_id, name) VALUES
(1, 'gram'),
(2, 'kilogram'),
(3, 'milliliter'),
(4, 'liter'),
(5, 'teaspoon'),
(6, 'tablespoon'),
(7, 'cup'),
(8, 'piece');


INSERT INTO ingredients (ingredient_id, name, description) VALUES
(1, 'Salt', 'Granulated seasoning used in almost every recipe'),
(2, 'Olive Oil', 'Versatile cooking oil made from olives'),
(3, 'Chicken Breast', 'Lean cut of chicken meat'),
(4, 'White Rice', 'Staple grain used in many cuisines'),
(5, 'Whole Milk', 'Full-fat milk from cows'),
(6, 'Egg', 'Common protein-rich ingredient'),
(7, 'Tomato', 'Red fruit often used as a vegetable in cooking'),
(8, 'Garlic', 'Aromatic bulb used in savory dishes'),
(9, 'Onion', 'Essential base ingredient for many meals'),
(10, 'Butter', 'Solid fat made by churning milk or cream'),
(11, 'Sugar', 'Sweetener used in desserts and baking'),
(12, 'Black Pepper', 'Common spice used for heat and flavor'),
(13, 'Ketchup', 'Sweet and tangy tomato-based condiment'),
(14, 'Mayonnaise', 'Creamy condiment made from eggs and oil'),
(15, 'Flour', 'Powdered grain used in baking'),
(16, 'Baking Powder', 'Leavening agent used in baking'),
(17, 'Vegetable Oil', 'Neutral cooking oil from plants'),
(18, 'Canned Tomatoes', 'Preserved tomatoes used in sauces'),
(19, 'Cheddar Cheese', 'Sharp-tasting yellow cheese'),
(20, 'Tuna (Canned)', 'Canned fish, often used in salads or sandwiches'),
(21, 'Frozen Peas', 'Pre-packaged peas, stored frozen'),
(22, 'Vanilla Extract', 'Concentrated vanilla flavoring used in desserts'),
(23, 'Carrot', 'Crunchy root vegetable, often used raw or cooked'),
(24, 'Broccoli', 'Green vegetable with edible florets'),
(25, 'Bell Pepper', 'Sweet and colorful pepper, not spicy'),
(26, 'Spinach', 'Leafy green, rich in iron and nutrients'),
(27, 'Zucchini', 'Green summer squash, soft texture when cooked'),
(28, 'Cucumber', 'Cool, crisp vegetable often eaten raw'),
(29, 'Celery', 'Crisp stalk used in soups, salads, or raw'),
(30, 'Potato', 'Starchy root vegetable, very versatile'),
(31, 'Sweet Potato', 'Sweet, orange-fleshed root vegetable'),
(32, 'Green Beans', 'Slim green pods, eaten whole');

INSERT INTO ingredient_measurement_units (ingredient_id, measurement_unit_id) VALUES
(1, 1), (1, 5), (1, 6),
(2, 3), (2, 6), (2, 7),
(3, 1), (3, 2), (3, 8),
(4, 1), (4, 7),
(5, 3), (5, 4), (5, 7),
(6, 8),
(7, 1), (7, 8),
(8, 1), (8, 5), (8, 6),
(9, 1), (9, 8),
(10, 1), (10, 6), (10, 7),
(11, 1), (11, 6), (11, 7),
(12, 1), (12, 5), (12, 6),
(13, 3), (13, 5), (13, 6),
(14, 3), (14, 5), (14, 6),
(15, 1), (15, 7),
(16, 1), (16, 5),
(17, 3), (17, 6), (17, 7),
(18, 8),
(19, 1), (19, 8),
(20, 8),
(21, 1), (21, 7),
(22, 3), (22, 5),
(23, 1), (23, 8),
(24, 1), (24, 8),
(25, 1), (25, 8),
(26, 1), (26, 7),
(27, 1), (27, 8),
(28, 1), (28, 8),
(29, 1), (29, 8),
(30, 1), (30, 8),
(31, 1), (31, 8),
(32, 1), (32, 7);


INSERT INTO ingredient_food_groups (ingredient_id, food_group_id) VALUES
(1, 10),
(2, 13),
(3, 5),
(4, 11),
(5, 4),
(6, 5),
(7, 2),
(8, 3),
(9, 1),
(10, 4),
(11, 11),
(12, 10),
(13, 9),
(14, 9),
(15, 11),
(16, 11),
(17, 13),
(18, 12),
(19, 4),
(20, 6),
(21, 8),
(22, 11),
(23, 1),
(24, 1),
(25, 1),
(26, 1),
(27, 1),
(28, 1),
(29, 1),
(30, 1),
(31, 1),
(32, 1);

