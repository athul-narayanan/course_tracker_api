CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    role VARCHAR(100) NOT NULL UNIQUE,
    priority INT NOT NULL
);

INSERT INTO roles (role, priority)
VALUES
    ('Admin',        1),
    ('User',         2);

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    firstname VARCHAR(100) NOT NULL,
    lastname VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    role VARCHAR(100) NOT NULL DEFAULT 'User'
        REFERENCES roles(role) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS universities (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    province VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS fields (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS specializations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    field_id INT NOT NULL REFERENCES fields(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(name, field_id)
);

INSERT INTO universities (name, province) VALUES
-- Ontario (ON)
('Algoma University', 'Ontario'),
('University of Toronto', 'Ontario'),
('Toronto Metropolitan University', 'Ontario'),
('York University', 'Ontario'),
('Queen''s University', 'Ontario'),
('Western University', 'Ontario'),
('University of Waterloo', 'Ontario'),
('University of Guelph', 'Ontario'),
('University of Ottawa', 'Ontario'),
('Carleton University', 'Ontario'),
('McMaster University', 'Ontario'),
('Lakehead University', 'Ontario'),
('Laurentian University', 'Ontario'),
('Trent University', 'Ontario'),
('Nipissing University', 'Ontario'),
('Ontario Tech University', 'Ontario'),
('Brock University', 'Ontario'),
('Wilfrid Laurier University', 'Ontario'),

-- Quebec (QC)
('McGill University', 'Quebec'),
('Concordia University', 'Quebec'),
('Université de Montréal', 'Quebec'),
('Université Laval', 'Quebec'),
('École Polytechnique de Montréal', 'Quebec'),
('Université du Québec à Montréal (UQAM)', 'Quebec'),
('HEC Montréal', 'Quebec'),
('Université de Sherbrooke', 'Quebec'),
('Université du Québec en Outaouais', 'Quebec'),
('Université du Québec à Trois-Rivières', 'Quebec'),
('Université du Québec à Chicoutimi', 'Quebec'),
('Institut National de la Recherche Scientifique (INRS)', 'Quebec'),
('École de technologie supérieure (ETS)', 'Quebec'),

-- British Columbia (BC)
('University of British Columbia', 'British Columbia'),
('Simon Fraser University', 'British Columbia'),
('University of Victoria', 'British Columbia'),
('University of Northern British Columbia', 'British Columbia'),
('Thompson Rivers University', 'British Columbia'),
('Royal Roads University', 'British Columbia'),
('Kwantlen Polytechnic University', 'British Columbia'),
('Capilano University', 'British Columbia'),
('Emily Carr University of Art and Design', 'British Columbia'),
('Vancouver Island University', 'British Columbia'),
('British Columbia Institute of Technology (BCIT)', 'British Columbia'),

-- Alberta (AB)
('University of Alberta', 'Alberta'),
('University of Calgary', 'Alberta'),
('University of Lethbridge', 'Alberta'),
('MacEwan University', 'Alberta'),
('Mount Royal University', 'Alberta'),
('Athabasca University', 'Alberta'),
('Northern Alberta Institute of Technology (NAIT)', 'Alberta'),

-- Manitoba (MB)
('University of Manitoba', 'Manitoba'),
('University of Winnipeg', 'Manitoba'),
('Brandon University', 'Manitoba'),
('Canadian Mennonite University', 'Manitoba'),
('Université de Saint-Boniface', 'Manitoba'),

-- Saskatchewan (SK)
('University of Saskatchewan', 'Saskatchewan'),
('University of Regina', 'Saskatchewan'),
('First Nations University of Canada', 'Saskatchewan'),

-- Nova Scotia (NS)
('Dalhousie University', 'Nova Scotia'),
('Saint Mary''s University', 'Nova Scotia'),
('Mount Saint Vincent University', 'Nova Scotia'),
('Acadia University', 'Nova Scotia'),
('Cape Breton University', 'Nova Scotia'),
('St. Francis Xavier University', 'Nova Scotia'),
('Atlantic School of Theology', 'Nova Scotia'),

-- New Brunswick (NB)
('University of New Brunswick', 'New Brunswick'),
('St. Thomas University', 'New Brunswick'),
('Mount Allison University', 'New Brunswick'),
('Université de Moncton', 'New Brunswick'),

-- Newfoundland & Labrador (NL)
('Memorial University of Newfoundland', 'Newfoundland and Labrador'),

-- Prince Edward Island (PEI)
('University of Prince Edward Island', 'Prince Edward Island'),

-- Territories (YT, NT, NU)
('Yukon University', 'Yukon'),
('Aurora College', 'Northwest Territories'),
('Nunavut Arctic College', 'Nunavut');

INSERT INTO fields (name) VALUES
('Computer Science & IT'),
('Business & Management'),
('Engineering'),
('Health & Medicine'),
('Life Sciences'),
('Environmental Science'),
('Psychology & Social Science'),
('Education & Teaching'),
('Law & Public Policy'),
('Arts, Design & Humanities'),
('Mathematics & Statistics'),
('Physical Sciences'),
('Media & Communications'),
('Agriculture & Food Science'),
('Architecture & Urban Planning'),
('Hospitality & Tourism'),
('Sports & Physical Education'),
('Veterinary & Animal Science'),
('Mining & Geological Studies'),
('Marine & Ocean Studies');


INSERT INTO specializations (name, field_id)
SELECT s.name, f.id FROM (VALUES
('Software Engineering'),
('Computer Science'),
('Data Science'),
('Cybersecurity'),
('Artificial Intelligence'),
('Cloud Computing'),
('Machine Learning'),
('Information Systems'),
('Computer Networks'),
('Human Computer Interaction')
) AS s(name)
JOIN fields f ON f.name = 'Computer Science & IT';

INSERT INTO specializations (name, field_id)
SELECT s.name, f.id FROM (VALUES
('MBA'),
('Finance'),
('Accounting'),
('Marketing'),
('International Business'),
('Business Analytics'),
('Human Resource Management'),
('Supply Chain Management')
) AS s(name)
JOIN fields f ON f.name = 'Business & Management';

INSERT INTO specializations (name, field_id)
SELECT s.name, f.id FROM (VALUES
('Civil Engineering'),
('Mechanical Engineering'),
('Electrical Engineering'),
('Chemical Engineering'),
('Industrial Engineering'),
('Biomedical Engineering'),
('Aerospace Engineering'),
('Materials Engineering'),
('Mechatronics Engineering'),
('Petroleum Engineering')
) AS s(name)
JOIN fields f ON f.name = 'Engineering';

INSERT INTO specializations (name, field_id)
SELECT s.name, f.id FROM (VALUES
('Nursing'),
('Public Health'),
('Pharmacy'),
('Physiotherapy'),
('Medical Laboratory Science'),
('Nutrition & Dietetics')
) AS s(name)
JOIN fields f ON f.name = 'Health & Medicine';

INSERT INTO specializations (name, field_id)
SELECT s.name, f.id FROM (VALUES
('Biology'),
('Biotechnology'),
('Microbiology'),
('Neuroscience'),
('Biochemistry'),
('Genetics')
) AS s(name)
JOIN fields f ON f.name = 'Life Sciences';

INSERT INTO specializations (name, field_id)
SELECT s.name, f.id FROM (VALUES
('Psychology'),
('Clinical Psychology'),
('Sociology'),
('Social Work')
) AS s(name)
JOIN fields f ON f.name = 'Psychology & Social Science';

INSERT INTO specializations (name, field_id)
SELECT s.name, f.id FROM (VALUES
('Architecture'),
('Urban Planning'),
('Landscape Architecture'),
('Interior Architecture')
) AS s(name)
JOIN fields f ON f.name = 'Architecture & Urban Planning';

INSERT INTO specializations (name, field_id)
SELECT s.name, f.id FROM (VALUES
('Environmental Science'),
('Climate Studies'),
('Marine Biology'),
('Fisheries & Aquaculture'),
('Agricultural Science'),
('Mining Engineering')
) AS s(name)
JOIN fields f ON f.name IN (
    'Environmental Science',
    'Marine & Ocean Studies',
    'Agriculture & Food Science',
    'Mining & Geological Studies'
);


