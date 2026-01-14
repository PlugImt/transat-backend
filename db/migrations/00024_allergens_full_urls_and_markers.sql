-- +goose Up
-- +goose StatementBegin

-- Ensure allergens picture_url uses full absolute URLs
UPDATE allergens
SET picture_url = CONCAT('https://toast-js.ew.r.appspot.com/', picture_url)
WHERE picture_url IS NOT NULL
  AND picture_url <> ''
  AND picture_url NOT LIKE 'http%';

-- Add a flag to distinguish allergen entries from label/marker entries
ALTER TABLE allergens
    ADD COLUMN IF NOT EXISTS is_marker BOOLEAN NOT NULL DEFAULT FALSE;

-- Insert label/marker pictograms (pictosLabel) into the allergens table as markers
INSERT INTO allergens (name, description, description_en, picture_url, is_marker)
VALUES
    ('ardoise',
     'Suggestion du chef',
     'Chef''s suggestion',
     'https://toast-js.ew.r.appspot.com/views/assets/img/pictosLabel/ardoise.svg',
     TRUE),
    ('formule',
     'Plat inclus dans une formule',
     'Dish included in a menu deal',
     'https://toast-js.ew.r.appspot.com/views/assets/img/pictosLabel/formule.svg',
     TRUE),
    ('vitalite',
     'Option vitalité',
     'Vitality option',
     'https://toast-js.ew.r.appspot.com/views/assets/img/pictosLabel/vitalite.svg',
     TRUE),
    ('vegetarien',
     'Plat végétarien',
     'Vegetarian dish',
     'https://toast-js.ew.r.appspot.com/views/assets/img/pictosLabel/vegetarien.svg',
     TRUE),
    ('bio',
     'Produit issu de l''agriculture biologique',
     'Organic product',
     'https://toast-js.ew.r.appspot.com/views/assets/img/pictosLabel/bio.svg',
     TRUE),
    ('local',
     'Produit local',
     'Local product',
     'https://toast-js.ew.r.appspot.com/views/assets/img/pictosLabel/local.svg',
     TRUE),
    ('saison',
     'Produit de saison',
     'Seasonal product',
     'https://toast-js.ew.r.appspot.com/views/assets/img/pictosLabel/saison.svg',
     TRUE),
    ('equitable',
     'Produit issu du commerce équitable',
     'Fair trade product',
     'https://toast-js.ew.r.appspot.com/views/assets/img/pictosLabel/equitable.svg',
     TRUE),
    ('weightWatcher',
     'O Api Déj',
     'O Api Déj',
     'https://toast-js.ew.r.appspot.com/views/assets/img/pictosLabel/weightWatcher.svg',
     TRUE),
    ('peche',
     'Pêche durable',
     'Sustainable fishing',
     'https://toast-js.ew.r.appspot.com/views/assets/img/pictosLabel/peche.svg',
     TRUE),
    ('france',
     'Produit certifié Origine France',
     'Certified French origin product',
     'https://toast-js.ew.r.appspot.com/views/assets/img/pictosLabel/france.svg',
     TRUE)
ON CONFLICT (name) DO UPDATE
    SET description    = EXCLUDED.description,
        description_en = EXCLUDED.description_en,
        picture_url    = EXCLUDED.picture_url,
        is_marker      = EXCLUDED.is_marker;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE allergens
    DROP COLUMN IF EXISTS is_marker;
-- Note: We intentionally do NOT delete the inserted markers to avoid data loss.
-- +goose StatementEnd

