-- Step 1: Add the food_group_id column back to the ingredients table
ALTER TABLE ingredients ADD COLUMN food_group_id INT;

-- Step 2: Add the foreign key constraint back to the ingredients table
ALTER TABLE ingredients 
ADD CONSTRAINT ingredients_food_group_id_fkey FOREIGN KEY (food_group_id) 
REFERENCES food_groups (food_group_id) ON DELETE SET NULL;

-- Step 3: Drop the ingredient_food_groups join table
DROP TABLE IF EXISTS ingredient_food_groups;

