-- +goose Up
-- +goose StatementBegin

-- Core allergens catalog
CREATE TABLE IF NOT EXISTS allergens
(
    id_allergen    SERIAL PRIMARY KEY,
    name           VARCHAR(100) NOT NULL UNIQUE,
    description    TEXT,
    description_en TEXT,
    picture_url    TEXT
);

-- Link table between restaurant articles and allergens
CREATE TABLE IF NOT EXISTS restaurant_article_allergens
(
    id_restaurant_articles INTEGER NOT NULL,
    id_allergen            INTEGER NOT NULL,
    PRIMARY KEY (id_restaurant_articles, id_allergen),
    CONSTRAINT fk_article_allergens_article
        FOREIGN KEY (id_restaurant_articles)
            REFERENCES restaurant_articles (id_restaurant_articles)
            ON DELETE CASCADE,
    CONSTRAINT fk_article_allergens_allergen
        FOREIGN KEY (id_allergen)
            REFERENCES allergens (id_allergen)
            ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_article_allergens_article
    ON restaurant_article_allergens (id_restaurant_articles);

CREATE INDEX IF NOT EXISTS idx_article_allergens_allergen
    ON restaurant_article_allergens (id_allergen);

-- Seed allergens based on exemple_restaurant_page.html / loadingAllergenes
INSERT INTO allergens (name, description, description_en, picture_url)
VALUES
    ('Arachide',
     'Arachides et produits à base d''arachides',
     'Peanuts and peanut-based products',
     'views/assets/img/pictosAllergenes/Arachide.svg'),
    ('FruitaCoque',
     'Fruits à coques* et produits à base de ces fruits',
     'Tree nuts* and products based on these nuts',
     'views/assets/img/pictosAllergenes/FruitaCoque.svg'),
    ('Soja',
     'Soja et produits à base de soja',
     'Soy and soy-based products',
     'views/assets/img/pictosAllergenes/Soja.svg'),
    ('Crustace',
     'Crustacés et produits à base de crustacés',
     'Crustaceans and crustacean-based products',
     'views/assets/img/pictosAllergenes/Crustace.svg'),
    ('Celeri',
     'Céleri et produits à base de céleri',
     'Celery and celery-based products',
     'views/assets/img/pictosAllergenes/Celeri.svg'),
    ('Gluten',
     'Céréales contenant du gluten* et produits à base de ces céréales',
     'Cereals containing gluten* and products based on these cereals',
     'views/assets/img/pictosAllergenes/Gluten.svg'),
    ('Lupin',
     'Lupin et produits à base de lupin',
     'Lupin and lupin-based products',
     'views/assets/img/pictosAllergenes/Lupin.svg'),
    ('Mollusque',
     'Mollusques et produits à base de mollusques',
     'Mollusks and mollusk-based products',
     'views/assets/img/pictosAllergenes/Mollusque.svg'),
    ('Moutarde',
     'Moutarde et produits à base de moutarde',
     'Mustard and mustard-based products',
     'views/assets/img/pictosAllergenes/Moutarde.svg'),
    ('Oeuf',
     'Œufs et produits à base d''œufs',
     'Eggs and egg-based products',
     'views/assets/img/pictosAllergenes/Oeuf.svg'),
    ('Poisson',
     'Poissons et produits à base de poisson',
     'Fish and fish-based products',
     'views/assets/img/pictosAllergenes/Poisson.svg'),
    ('ProduitLaitier',
     'Lait et produits à base de lait (y compris lactose)',
     'Milk and milk-based products /including lactose',
     'views/assets/img/pictosAllergenes/ProduitLaitier.svg'),
    ('Sesame',
     'Graines de sésame et produits à base de graines de sésame',
     'Sesame seeds and sesame seed-based products',
     'views/assets/img/pictosAllergenes/Sesame.svg'),
    ('Sulfites',
     'Sulfites et anhydride sulfureux',
     'Sulphites and sulphur dioxide',
     'views/assets/img/pictosAllergenes/Sulfites.svg'),
    ('vide',
     'none',
     'none',
     'views/assets/img/pictosAllergenes/vide.svg'),
    ('undefined',
     'none',
     'none',
     'views/assets/img/pictosAllergenes/undefined.svg')
ON CONFLICT (name) DO NOTHING;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS restaurant_article_allergens;
DROP TABLE IF EXISTS allergens;
-- +goose StatementEnd

