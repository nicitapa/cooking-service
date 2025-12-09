CREATE TABLE categories (
                            id SERIAL PRIMARY KEY,
                            name TEXT NOT NULL
);

CREATE TABLE recipes (
                         id SERIAL PRIMARY KEY,
                         title TEXT NOT NULL,
                         description TEXT,
                         instructions TEXT,
                         image_url TEXT,
                         category_id INTEGER REFERENCES categories(id),
                         created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE ingredients (
                             id SERIAL PRIMARY KEY,
                             name TEXT NOT NULL,
                             unit TEXT
);

CREATE TABLE recipe_ingredients (
                                    recipe_id INTEGER REFERENCES recipes(id) ON DELETE CASCADE,
                                    ingredient_id INTEGER REFERENCES ingredients(id),
                                    quantity NUMERIC,
                                    PRIMARY KEY (recipe_id, ingredient_id)
);

CREATE TABLE tags (
                      id SERIAL PRIMARY KEY,
                      name TEXT NOT NULL
);

CREATE TABLE recipe_tags (
                             recipe_id INTEGER REFERENCES recipes(id) ON DELETE CASCADE,
                             tag_id INTEGER REFERENCES tags(id),
                             PRIMARY KEY (recipe_id, tag_id)
);
