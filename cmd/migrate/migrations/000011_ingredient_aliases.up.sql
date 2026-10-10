-- alternative names (stored normalized) that resolve to an existing ingredient during recipe import
CREATE TABLE ingredient_aliases (
    alias TEXT PRIMARY KEY,
    ingredient_id INTEGER NOT NULL REFERENCES ingredients(ingredient_id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW()
);
