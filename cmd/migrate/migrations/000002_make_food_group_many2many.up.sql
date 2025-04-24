-- Add a new table for the join between ingredients and food_groups
CREATE TABLE ingredient_food_groups (
    ingredient_id INT NOT NULL,
    food_group_id INT NOT NULL,
    PRIMARY KEY (ingredient_id, food_group_id),
    FOREIGN KEY (ingredient_id) REFERENCES ingredients (ingredient_id),
    FOREIGN KEY (food_group_id) REFERENCES food_groups (food_group_id)
);

-- Remove the foreign key constraint from ingredients table
ALTER TABLE ingredients DROP CONSTRAINT ingredients_food_group_id_fkey;

-- Drop the food_group_id column from ingredients table
ALTER TABLE ingredients DROP COLUMN food_group_id;

