-- Jeu de données de démonstration des réservations.
--
--   psql "$DATABASE_URL" -f db/seeds/reservations_demo.sql
--
-- Il crée des éléments « [Démo] ... » (salles à créneaux, objets prêtés) et au moins 6 réservations par utilisateur
-- existant : passées, en cours, à venir, prêts rendus, prêts en cours et prêts en retard (1, 5 et 7 jours et plus).
-- Les rappels déjà échus sont marqués comme envoyés pour ne rien notifier au démarrage du serveur.
--
-- Rejouable : le jeu précédent est supprimé d'abord. Rien d'autre n'est touché : pour tout retirer,
--   DELETE FROM reservation_element WHERE name LIKE '[Démo] %'; DELETE FROM reservation_category WHERE name LIKE '[Démo] %';
-- (les réservations, rappels, files d'attente et annulations de ces éléments partent en cascade).
-- Prérequis : migrations appliquées jusqu'à la 00044 (lancer le serveur une fois).

\set ON_ERROR_STOP on

BEGIN;

DELETE FROM reservation_element WHERE name LIKE '[Démo] %';
DELETE FROM reservation_category WHERE name LIKE '[Démo] %';

INSERT INTO clubs (name, picture)
SELECT 'Club Démo', 'https://placehold.co/200x200'
WHERE NOT EXISTS (SELECT 1 FROM clubs);

CREATE TEMP TABLE _users ON COMMIT DROP AS
SELECT email, row_number() OVER (ORDER BY email) AS rn FROM newf;

-- s: nombre d'éléments à créneaux ; l: nombre d'objets prêtés (un prêt en cours pour un utilisateur sur trois).
CREATE TEMP TABLE _p ON COMMIT DROP AS
SELECT date_trunc('second', now() AT TIME ZONE 'UTC') AS now_utc,
       date_trunc('day', now() AT TIME ZONE 'UTC')    AS today,
       (SELECT min(id_clubs) FROM clubs)              AS club,
       6                                              AS s,
       GREATEST(8, CEIL((SELECT count(*) FROM _users) / 3.0))::int AS l;

INSERT INTO reservation_category (name, id_parent_category, id_clubs)
SELECT n, NULL, club FROM _p, unnest(ARRAY['[Démo] Salles et équipements', '[Démo] Matériel à emprunter']) AS n;

INSERT INTO reservation_element (name, description, location, slot, id_clubs, id_reservation_category)
SELECT '[Démo] ' || t.name, t.description, t.location, TRUE, p.club,
       (SELECT id_reservation_category FROM reservation_category WHERE name = '[Démo] Salles et équipements')
FROM _p p, unnest(
        ARRAY['Salle Ampère', 'Salle Turing', 'Studio son', 'Terrain de padel', 'Salle de sport', 'Cabine d''enregistrement'],
        ARRAY['Salle de réunion 12 places', 'Salle informatique', 'Studio d''enregistrement', 'Terrain extérieur', 'Salle de musculation', 'Cabine insonorisée'],
        ARRAY['Bâtiment A', 'Bâtiment B', 'Maison des élèves', 'Campus', 'Gymnase', 'Maison des élèves']
     ) AS t(name, description, location);

INSERT INTO reservation_element (name, description, location, slot, id_clubs, id_reservation_category)
SELECT '[Démo] ' || (ARRAY['Câble HDMI', 'Rallonge électrique', 'Enceinte Bluetooth', 'Appareil photo', 'Perceuse', 'Jeu de société',
                           'Vidéoprojecteur', 'Casque audio'])[(i - 1) % 8 + 1]
           || CASE WHEN i > 8 THEN ' n°' || ((i - 1) / 8 + 1) ELSE '' END,
       'Objet à emprunter', 'Local associatif', FALSE, p.club,
       (SELECT id_reservation_category FROM reservation_category WHERE name = '[Démo] Matériel à emprunter')
FROM _p p, generate_series(1, (SELECT l FROM _p)) AS i;

CREATE TEMP TABLE _slot_items ON COMMIT DROP AS
SELECT id_reservation_element AS id, (row_number() OVER (ORDER BY id_reservation_element) - 1) AS idx
FROM reservation_element WHERE name LIKE '[Démo] %' AND slot;

CREATE TEMP TABLE _loan_items ON COMMIT DROP AS
SELECT id_reservation_element AS id, (row_number() OVER (ORDER BY id_reservation_element) - 1) AS idx
FROM reservation_element WHERE name LIKE '[Démo] %' AND NOT slot;

-- Créneaux passés (3) et à venir (2) pour chacun. Chaque utilisateur occupe une case distincte (élément, heure, jour)
-- de la grille : jamais deux réservations sur le même créneau d'un élément.
INSERT INTO reservation (email, id_reservation_element, start_date, end_date, created_at)
SELECT u.email, si.id, g.start_date, g.start_date + INTERVAL '1 hour', (g.start_date - INTERVAL '3 days') AT TIME ZONE 'UTC'
FROM _users u
CROSS JOIN _p p
CROSS JOIN (VALUES (0, -1, 2), (1, -1, 4), (2, -1, 6), (3, 1, 1), (4, 1, 4)) AS k(k, sign, base)
JOIN _slot_items si ON si.idx = (u.rn - 1 + k.k) % p.s
CROSS JOIN LATERAL (
    SELECT p.today + make_interval(days => k.sign * (k.base + 10 * ((u.rn - 1) / (p.s * 12))::int),
                                   hours => 7 + (((u.rn - 1) / p.s) % 12)::int) AS start_date
) g;

-- Créneaux du jour (avant et après l'heure actuelle), pour que le planning d'aujourd'hui ne soit pas vide.
INSERT INTO reservation (email, id_reservation_element, start_date, end_date, created_at)
SELECT u.email, si.id, g.start_date, g.start_date + INTERVAL '1 hour', (g.start_date - INTERVAL '3 days') AT TIME ZONE 'UTC'
FROM _users u
CROSS JOIN _p p
JOIN _slot_items si ON si.idx = (u.rn - 1 + 5) % p.s
CROSS JOIN LATERAL (SELECT p.today + make_interval(hours => 7 + (((u.rn - 1) / p.s) % 12)::int) AS start_date) g
WHERE u.rn - 1 < p.s * 12 AND g.start_date <> date_trunc('hour', p.now_utc);

-- Créneaux en cours : un par élément, pour les premiers utilisateurs.
INSERT INTO reservation (email, id_reservation_element, start_date, end_date, created_at)
SELECT u.email, si.id, date_trunc('hour', p.now_utc), date_trunc('hour', p.now_utc) + INTERVAL '1 hour',
       (date_trunc('hour', p.now_utc) - INTERVAL '2 days') AT TIME ZONE 'UTC'
FROM _users u CROSS JOIN _p p JOIN _slot_items si ON si.idx = u.rn - 1
WHERE u.rn <= p.s;

-- Un prêt rendu pour chacun, sur des fenêtres de deux jours qui ne se chevauchent jamais pour un même objet.
INSERT INTO reservation (email, id_reservation_element, start_date, end_date, created_at)
SELECT u.email, li.id, g.start_date, g.start_date + INTERVAL '2 days', g.start_date AT TIME ZONE 'UTC'
FROM _users u
CROSS JOIN _p p
JOIN _loan_items li ON li.idx = (u.rn - 1) % p.l
CROSS JOIN LATERAL (
    SELECT p.today + INTERVAL '9 hours' - make_interval(days => 12 + 3 * ((u.rn - 1) / p.l)::int) AS start_date
) g;

-- Prêts en cours, jamais deux sur un même objet : récent, puis en retard de 1, 5 et 7 jours et plus.
INSERT INTO reservation (email, id_reservation_element, start_date, end_date, created_at)
SELECT u.email, li.id, g.start_date, NULL, g.start_date AT TIME ZONE 'UTC'
FROM _users u
CROSS JOIN _p p
JOIN _loan_items li ON li.idx = u.rn - 1
CROSS JOIN LATERAL (
    SELECT p.now_utc - CASE u.rn % 4 WHEN 0 THEN INTERVAL '3 hours'
                                     WHEN 1 THEN INTERVAL '1 day 6 hours'
                                     WHEN 2 THEN INTERVAL '5 days 4 hours'
                                     ELSE INTERVAL '8 days 2 hours' END AS start_date
) g
WHERE u.rn <= p.l;

-- Les rappels déjà échus (prêts, créneaux dans moins d'une heure) sont considérés comme envoyés.
INSERT INTO reservation_reminders (email, id_reservation_element, start_date, milestone)
SELECT r.email, r.id_reservation_element, r.start_date, 'slot_1h'
FROM reservation r
JOIN reservation_element re ON re.id_reservation_element = r.id_reservation_element
CROSS JOIN _p p
WHERE re.name LIKE '[Démo] %' AND re.slot AND r.start_date > p.now_utc AND r.start_date - INTERVAL '1 hour' <= p.now_utc
ON CONFLICT DO NOTHING;
INSERT INTO reservation_reminders (email, id_reservation_element, start_date, milestone)
SELECT r.email, r.id_reservation_element, r.start_date, m.name
FROM reservation r
JOIN reservation_element re ON re.id_reservation_element = r.id_reservation_element
CROSS JOIN _p p
CROSS JOIN (VALUES ('loan_d1', INTERVAL '1 day'), ('loan_d5', INTERVAL '5 days'), ('loan_d7', INTERVAL '7 days')) AS m(name, after)
WHERE re.name LIKE '[Démo] %' AND NOT re.slot AND r.end_date IS NULL AND r.start_date + m.after <= p.now_utc
ON CONFLICT DO NOTHING;

COMMIT;

-- Récapitulatif.
SELECT CASE
           WHEN re.slot AND r.end_date <= p.now_utc THEN '1. créneau passé'
           WHEN re.slot AND r.start_date <= p.now_utc THEN '2. créneau en cours'
           WHEN re.slot THEN '3. créneau à venir'
           WHEN r.end_date IS NOT NULL THEN '4. prêt rendu'
           WHEN r.start_date > p.now_utc - INTERVAL '1 day' THEN '5. prêt en cours'
           WHEN r.start_date > p.now_utc - INTERVAL '5 days' THEN '6. prêt en retard (1 j et plus)'
           WHEN r.start_date > p.now_utc - INTERVAL '7 days' THEN '7. prêt en retard (5 j et plus)'
           ELSE '8. prêt en retard (7 j et plus)'
       END AS etat, count(*) AS reservations
FROM reservation r
JOIN reservation_element re ON re.id_reservation_element = r.id_reservation_element
CROSS JOIN (SELECT date_trunc('second', now() AT TIME ZONE 'UTC') AS now_utc) p
WHERE re.name LIKE '[Démo] %'
GROUP BY 1 ORDER BY 1;

SELECT count(*) AS utilisateurs, min(c) AS minimum_par_utilisateur, max(c) AS maximum_par_utilisateur
FROM (SELECT n.email, count(r.email) AS c
      FROM newf n
      LEFT JOIN reservation r ON r.email = n.email
           AND r.id_reservation_element IN (SELECT id_reservation_element FROM reservation_element WHERE name LIKE '[Démo] %')
      GROUP BY n.email) per_user;
