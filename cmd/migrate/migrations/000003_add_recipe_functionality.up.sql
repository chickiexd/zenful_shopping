-- Create meal_types table
CREATE TABLE meal_types (
    meal_type_id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

-- Some default meal types
INSERT INTO meal_types (name) VALUES
  ('breakfast'),
  ('main'),
  ('appetizer'),
  ('snack'),
  ('dessert');

-- Create recipes table
CREATE TABLE recipes (
    recipe_id     SERIAL PRIMARY KEY,
    title         TEXT NOT NULL,
    description   TEXT,
    public        BOOLEAN DEFAULT FALSE,
    cook_time     INTEGER,
    servings      INTEGER,
    image_path    TEXT,
    meal_type_id  INTEGER REFERENCES meal_types(meal_type_id),
    created_at    TIMESTAMP DEFAULT now(),
    updated_at    TIMESTAMP DEFAULT now()
);

-- Create recipe_has_ingredients table (many-to-many + quantity)
CREATE TABLE recipe_ingredients (
    recipe_id      INTEGER NOT NULL REFERENCES recipes(recipe_id),
    ingredient_id  INTEGER NOT NULL REFERENCES ingredients(ingredient_id),
    measurement_unit_id  INTEGER NOT NULL REFERENCES measurement_units(measurement_unit_id),
    quantity       DOUBLE PRECISION NOT NULL DEFAULT 1.0,
    PRIMARY KEY (recipe_id, ingredient_id)
);

-- Create instructions table
CREATE TABLE instructions (
    instruction_id SERIAL PRIMARY KEY,
    recipe_id      INTEGER NOT NULL REFERENCES recipes(recipe_id),
    step_number    INTEGER NOT NULL,
    description    TEXT NOT NULL
);

