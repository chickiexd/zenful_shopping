-- Drop the many-to-many relationship join table first
DROP TABLE IF EXISTS IngredientMeasurementUnits;

-- Remove the foreign key column in Ingredients that links to FoodGroups
ALTER TABLE Ingredients DROP CONSTRAINT IF EXISTS FK_FoodGroup;
ALTER TABLE Ingredients DROP COLUMN IF EXISTS FoodGroupID;

-- Drop main tables in reverse order of their dependencies to ensure integrity
DROP TABLE IF EXISTS Ingredients;
DROP TABLE IF EXISTS MeasurementUnits;
DROP TABLE IF EXISTS FoodGroups;

