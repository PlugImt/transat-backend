-- +goose Up
CREATE TABLE planning_sport (
    id SERIAL PRIMARY KEY,
    day_of_week VARCHAR(10) NOT NULL,
    activity VARCHAR(50) NOT NULL,  
    place VARCHAR(50),               
    start_time TIME NOT NULL,       
    end_time TIME NOT NULL
);

INSERT INTO planning_sport (day_of_week, activity, place, start_time, end_time) VALUES
--Monday
('Monday', 'Pilates', 'J059', '18:00:00', '19:00:00'),
('Monday', 'Basket H', 'Gymnase', '19:00:00', '21:00:00'),
('Monday', 'POMP''IMT', 'Tatami', '19:30:00', '21:00:00'),

--Tuesday
('Tuesday', 'Rugby F', NULL, '18:00:00', '20:30:00'),
('Tuesday', 'Handball (F&H)', 'Gymnase', '18:30:00', '21:30:00'),
('Tuesday', 'Foot M', 'Stade', '20:00:00', '22:30:00'),

--Wednesday
('Wednesday', 'Volley', 'Gymnase', '18:00:00', '21:00:00'),

--Thursday
('Thursday', 'Tennis', 'Terrain IMT', '13:00:00', '14:30:00'),
('Thursday', 'Athlétisme', 'Stade Moulin Boisseau', '14:30:00', '16:00:00'),
('Thursday', 'Escalade', 'EL CAP', '14:30:00', '16:30:00'),
('Thursday', 'Judo & JJB', 'Tatami', '15:30:00', '17:00:00'),
('Thursday', 'Ultimate', 'Gymnase', '17:00:00', '19:00:00'),
('Thursday', 'Natation', 'Piscine Carquefou', '17:30:00', '18:30:00'),
('Thursday', 'Badminton', 'Gymnase', '19:00:00', '21:00:00'),
('Thursday', 'POMP''IMT', 'Tatami', '20:00:00', '21:30:00');
-- +goose Down
DROP TABLE IF EXISTS planning_sport;