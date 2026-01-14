-- +goose Up
-- +goose StatementBegin

-- Update allergens to store only the filename (extract from full URL or path)
UPDATE allergens
SET picture_url = CASE
    -- Extract filename from full URL
    WHEN picture_url LIKE 'https://toast-js.ew.r.appspot.com/views/assets/img/pictosAllergenes/%' THEN
        SUBSTRING(picture_url FROM 'pictosAllergenes/(.+)$')
    WHEN picture_url LIKE 'https://toast-js.ew.r.appspot.com/views/assets/img/pictosLabel/%' THEN
        SUBSTRING(picture_url FROM 'pictosLabel/(.+)$')
    -- Extract filename from relative path
    WHEN picture_url LIKE 'views/assets/img/pictosAllergenes/%' THEN
        SUBSTRING(picture_url FROM 'pictosAllergenes/(.+)$')
    WHEN picture_url LIKE 'views/assets/img/pictosLabel/%' THEN
        SUBSTRING(picture_url FROM 'pictosLabel/(.+)$')
    -- Already just a filename, keep as is
    ELSE picture_url
END
WHERE picture_url IS NOT NULL AND picture_url <> '';

-- Update markers to store only the filename
UPDATE allergens
SET picture_url = CASE
    WHEN picture_url LIKE 'https://toast-js.ew.r.appspot.com/views/assets/img/pictosLabel/%' THEN
        SUBSTRING(picture_url FROM 'pictosLabel/(.+)$')
    WHEN picture_url LIKE 'views/assets/img/pictosLabel/%' THEN
        SUBSTRING(picture_url FROM 'pictosLabel/(.+)$')
    ELSE picture_url
END
WHERE is_marker = TRUE AND picture_url IS NOT NULL AND picture_url <> '';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Note: We don't restore full URLs in down migration as the local filenames are the new standard
-- +goose StatementEnd
