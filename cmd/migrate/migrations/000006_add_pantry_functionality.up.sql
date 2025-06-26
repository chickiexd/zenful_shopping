CREATE TABLE pantry_ingredients (
    pantry_ingredient_id SERIAL PRIMARY KEY,
    ingredient_id INTEGER NOT NULL,
    FOREIGN KEY (ingredient_id) REFERENCES ingredients(ingredient_id)
);
