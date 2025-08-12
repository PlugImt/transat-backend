-- +goose Up
-- +goose StatementBegin
DROP TABLE articles CASCADE;
DROP TABLE articles_types CASCADE;
DROP TABLE bassine CASCADE;
DROP TABLE caps CASCADE;
DROP TABLE courses CASCADE;
DROP TABLE games CASCADE;
DROP TABLE players_bassine CASCADE;
DROP TABLE players_caps CASCADE;
DROP TABLE realcampus_posts CASCADE;
DROP TABLE realcampus_friendships CASCADE;
DROP TABLE realcampus_post_reactions CASCADE;
DROP TABLE rooms CASCADE;
DROP TABLE support CASCADE;
DROP TABLE support_media CASCADE;
DROP TABLE washing_machines CASCADE;
DROP TABLE washing_machines_reports CASCADE;






-- +goose StatementEnd
