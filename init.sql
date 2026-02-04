CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Разрешить подключения со всех адресов
ALTER SYSTEM SET listen_addresses = '*';

-- Перезагрузить конфигурацию
SELECT pg_reload_conf();

-- Ждем немного для применения изменений
SELECT pg_sleep(2);

-- Создаем базу данных, если она не создалась через переменные окружения
-- (Это резервный вариант)
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_database WHERE datname = current_database()) THEN
        PERFORM dblink_exec('dbname=postgres', 'CREATE DATABASE "25-12courses"');
    END IF;
END
$$;

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

create table if not exists "session" (
	id_session uuid primary key,
	user_name varchar(50) not null unique,
	refresh_token varchar(500) not null unique,
	is_revoked boolean not null,
	created_at TIMESTAMP DEFAULT (now()) not null,
	expires_at TIMESTAMP  not null
);

create table if not exists "role" (
	id_role UUID primary key, 
	role varchar(50) not null unique
);

CREATE TABLE IF NOT EXISTS "user" (
    id_user serial PRIMARY KEY,
    user_name varchar(100) not null unique,
    password varchar(100) not null unique,
	role UUID references role(id_role)
);

CREATE TABLE IF NOT EXISTS passport (
    id_passport UUID PRIMARY KEY,
    place_birth varchar(100) not null,
    citizenship varchar(100) not null,
    gender varchar(10) not null,
    seria INTEGER not null,
    number INTEGER not null,
    passport_given varchar(50) not null,
    date_given date not null,
    code varchar(7) not null,
	UNIQUE(seria, number)
);

CREATE TABLE IF NOT EXISTS registrationaddress (
    id_regaddress UUID PRIMARY KEY,
    mail_index INTEGER not null,
    region varchar(100) not null,
    city varchar(100) not null,
    street varchar(50) not null,
    house varchar(50) not null,
    building varchar(50) not null,
    apartment varchar(50) not null
);

CREATE TABLE IF NOT EXISTS leveleducation (
    id_leveleducation UUID PRIMARY KEY,
    education varchar(100) not null unique
);

CREATE TABLE IF NOT EXISTS educationlistener (
    id_educationlistener UUID PRIMARY KEY,
    diplom_seria INTEGER not null,
    diplom_number INTEGER not null,
    date_given DATE not null,
    city varchar(255) not null,
    region varchar(255) not null,
    educational_institution varchar(255) not null,
    speciality varchar(255) not null,
    level_education uuid not null,
	UNIQUE(diplom_seria, diplom_number)
);

CREATE TABLE IF NOT EXISTS placework (
    id_placework UUID PRIMARY KEY,
    name_company varchar(255) not null,
    job_title varchar(255) not null,
    all_experience INTEGER not null,
    job_title_experience INTEGER not null
);

CREATE TABLE IF NOT EXISTS divisionseducation (
    id_divisionseducation UUID PRIMARY KEY,
    divisions varchar(100) UNIQUE not null
);

CREATE TABLE IF NOT EXISTS educationtypes (
    id_educationtype UUID PRIMARY KEY,
    type_name varchar(100) UNIQUE not null
);

CREATE TABLE IF NOT EXISTS programeducation (
    id_programeducation UUID PRIMARY KEY,
    name_prof_education varchar(100) not null unique,
    time_education INTEGER not null,
    individual_price REAL not null,
    group_price REAL not null,
    campus_price REAL not null,
    id_educationtype UUID REFERENCES educationtypes(id_educationtype) not null,
    id_divisionseducation UUID REFERENCES divisionseducation(id_divisionseducation) not null
);


create table if not exists accurateprogram (
	id_listener uuid not null,
	name_prof_education varchar(100) not null,
    time_education INTEGER not null,
    individual_price REAL not null,
    group_price REAL not null,
    campus_price REAL not null,
    educationtype varchar(255) not null,
    divisionseducation varchar(255) not null
);

CREATE TABLE IF NOT EXISTS legal_entity (
	id_legalentity uuid primary key,
	name_company varchar(255) not null, 
	inn varchar(12) not null,
	kpp varchar(12) not null, 
	ogrn varchar(15) not null,
	phone varchar(20) not null, 
	email varchar(255) not null,
	first_name varchar(100) not null,
    second_name varchar(100) not null,
    middle_name varchar(100),
	id_regaddress UUID REFERENCES registrationaddress(id_regaddress) not null
);

create table if not exists contractor (
	id_contractor uuid primary key,
	first_name varchar(100) not null,
    second_name varchar(100) not null,
    middle_name varchar(100),
	contact_phone varchar(20) UNIQUE not null,
    email varchar(50) UNIQUE not null,
	id_passport UUID REFERENCES passport(id_passport),
    id_regaddress UUID REFERENCES registrationaddress(id_regaddress) not null
);

CREATE TABLE IF NOT EXISTS listener (
    id_listener UUID PRIMARY KEY,
    first_name varchar(100) not null,
    second_name varchar(100) not null,
    middle_name varchar(100),
    date_of_birth DATE not null,
    snils varchar(14) UNIQUE not null,
    contact_phone varchar(20) UNIQUE not null,
    email varchar(50) UNIQUE not null,
    id_passport UUID REFERENCES passport(id_passport),
    id_regaddress UUID REFERENCES registrationaddress(id_regaddress) not null,
    id_educationlistener UUID REFERENCES educationlistener(id_educationlistener),
    id_placework UUID REFERENCES placework(id_placework),
    id_legalentity uuid references legal_entity(id_legalentity),
    id_contractor uuid references contractor(id_contractor) ON DELETE SET NULL,
    looting_education boolean DEFAULT false
);

CREATE TABLE IF NOT EXISTS enrollmentlistener (
    id_listener UUID REFERENCES listener(id_listener) not null,
    id_programeducation UUID REFERENCES programeducation(id_programeducation) not null,
    start_date DATE not null,
    end_date DATE not null,
    current_price decimal(10,2) not null,
    is_active  boolean not null,
	group_number varchar(50) not null,
	type_of_retraining varchar(50) not null,
    PRIMARY KEY (id_listener, id_programeducation)
);

create table if not exists executor (
	id_executor uuid primary key,
	status varchar(255) not null,
	first_name varchar(100) not null,
    second_name varchar(100) not null,
    middle_name varchar(100)
);



--дашборд

CREATE OR REPLACE VIEW v_total_listeners AS
SELECT COUNT(*) AS total_listeners FROM listener;


CREATE OR REPLACE VIEW v_total_programs AS
SELECT COUNT(*) AS total_programs FROM programeducation;


CREATE OR REPLACE VIEW v_active_enrollments AS
SELECT COUNT(*) AS active_enrollments
FROM enrollmentlistener
WHERE start_date <= now() AND end_date >= now();


CREATE OR REPLACE VIEW v_programs_ending_soon AS
SELECT 
    p.name_prof_education,
    e.end_date,
    COUNT(e.id_listener) AS total_listeners
FROM enrollmentlistener e
JOIN programeducation p ON e.id_programeducation = p.id_programeducation
WHERE e.end_date BETWEEN now() AND now() + interval '30 days'
GROUP BY p.name_prof_education, e.end_date
ORDER BY e.end_date ASC;


SELECT * FROM v_total_listeners;
SELECT * FROM v_total_programs;
SELECT * FROM v_active_enrollments;

--аудит

CREATE TABLE IF NOT EXISTS audit_log (
    id_audit serial PRIMARY KEY,
    table_name varchar(100) NOT NULL,
    operation varchar(10) NOT NULL,
    user_name varchar(50), 
    changed_at TIMESTAMP DEFAULT now(),
    old_data jsonb,
    new_data jsonb
);

CREATE OR REPLACE FUNCTION audit_all_tables()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        INSERT INTO audit_log(table_name, operation, user_name, new_data)
        VALUES(TG_TABLE_NAME, 'INSERT', current_user, row_to_json(NEW));
        RETURN NEW;
    ELSIF TG_OP = 'UPDATE' THEN
        INSERT INTO audit_log(table_name, operation, user_name, old_data, new_data)
        VALUES(TG_TABLE_NAME, 'UPDATE', current_user, row_to_json(OLD), row_to_json(NEW));
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        INSERT INTO audit_log(table_name, operation, user_name, old_data)
        VALUES(TG_TABLE_NAME, 'DELETE', current_user, row_to_json(OLD));
        RETURN OLD;
    END IF;
END;
$$ LANGUAGE plpgsql;


DO $$
DECLARE
    r RECORD;
BEGIN
    FOR r IN 
        SELECT table_name 
        FROM information_schema.tables 
        WHERE table_schema = 'public' 
          AND table_type = 'BASE TABLE'
          AND table_name <> 'audit_log'
    LOOP
        EXECUTE format('
            CREATE TRIGGER %I
            AFTER INSERT OR UPDATE OR DELETE ON %I
            FOR EACH ROW EXECUTE FUNCTION audit_all_tables();
        ', 'trg_audit_' || r.table_name, r.table_name);
    END LOOP;
END $$;


--дата на записях
CREATE OR REPLACE FUNCTION check_enrollment_dates()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.start_date < CURRENT_DATE THEN
        RAISE EXCEPTION 'Start date cannot be in the past';
    END IF;
    IF NEW.end_date <= NEW.start_date THEN
        RAISE EXCEPTION 'End date must be after start date';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_check_enrollment_dates
BEFORE INSERT OR UPDATE ON enrollmentlistener
FOR EACH ROW EXECUTE FUNCTION check_enrollment_dates();



--деактиваия записи
CREATE OR REPLACE PROCEDURE deactivate_finished_enrollments()
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE enrollmentlistener
    SET is_active = FALSE
    WHERE end_date < CURRENT_DATE
      AND is_active = TRUE;
END;
$$;


CREATE OR REPLACE PROCEDURE shuffle_program_order()
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE programeducation
    SET name_prof_education = name_prof_education
    WHERE time_education > 0;
END;
$$;

CREATE OR REPLACE PROCEDURE shuffle_education_type()
LANGUAGE plpgsql
AS $$
BEGIN
    WITH tmp AS (
        SELECT id_educationtype
        FROM educationtypes
        ORDER BY random()
    )
    UPDATE educationtypes l
    SET type_name = l.type_name
    FROM tmp
    WHERE l.id_educationtype = tmp.id_educationtype;
END;
$$;


CREATE OR REPLACE VIEW admin_table_sizes AS
SELECT
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname || '.' || tablename)) AS total_size,
    pg_size_pretty(pg_relation_size(schemaname || '.' || tablename)) AS table_size,
    pg_size_pretty(pg_total_relation_size(schemaname || '.' || tablename) 
                   - pg_relation_size(schemaname || '.' || tablename)) AS indexes_size
FROM pg_tables
WHERE schemaname NOT IN ('pg_catalog', 'information_schema')
ORDER BY pg_total_relation_size(schemaname || '.' || tablename) DESC;


--дашборд админа
CREATE VIEW admin_active_sessions AS
SELECT
    pid,
    application_name,
    client_addr,
    client_port,
    state,
    query_start
FROM pg_stat_activity
WHERE state = 'active'
ORDER BY query_start DESC;



CREATE VIEW admin_db_stats AS
SELECT
    d.datname AS database_name,
    pg_size_pretty(pg_database_size(d.datname)) AS total_size,
    s.numbackends AS active_connections,
    s.xact_commit AS committed_transactions,
    s.xact_rollback AS rolledback_transactions,
    s.blks_read AS disk_blocks_read,
    s.blks_hit AS buffer_hits
FROM pg_database d
LEFT JOIN pg_stat_database s ON d.datname = s.datname
WHERE d.datname = '25-12courses'
ORDER BY pg_database_size(d.datname) DESC;






INSERT INTO role (id_role, role)
VALUES 
('c7253f36-b89b-47df-b9e2-174eb4fbb227', 'admin'),
('236e78b7-0988-471b-b366-11638036803e', 'worker'),
('46ccb5d9-8f41-4677-a500-945154be25ef', 'accountant')
ON CONFLICT (id_role) DO NOTHING;

INSERT INTO "user" (user_name, password, role)
VALUES (
    'admin',
    '$2a$10$4n4jQDtxjCwpfDbJFAeqnOqfQfchn6jME7bIeOU8hDacRL7eP0zG6',
    'c7253f36-b89b-47df-b9e2-174eb4fbb227'
);




-- Passport
INSERT INTO passport (id_passport, place_birth, citizenship, gender, seria, number, passport_given, date_given, code)
VALUES
(gen_random_uuid(), 'Москва', 'Россия', 'Мужской', 1234, 111111, 'ОВД Москвы', '2010-01-01', '770-101'),
(gen_random_uuid(), 'Санкт-Петербург', 'Россия', 'Женский', 2345, 222222, 'ОВД СПБ', '2012-02-02', '780-002'),
(gen_random_uuid(), 'Казань', 'Россия', 'Мужской', 3456, 333333, 'ОВД Казань', '2015-03-03', '160003'),
(gen_random_uuid(), 'Новосибирск', 'Россия', 'Женский', 4567, 444444, 'ОВД Новосибирск', '2018-04-04', '540-004'),
(gen_random_uuid(), 'Сочи', 'Россия', 'Мужской', 5678, 555555, 'ОВД Сочи', '2020-05-05', '230-005');

-- Registration Address
INSERT INTO registrationaddress (id_regaddress, mail_index, region, city, street, house, building, apartment)
VALUES
(gen_random_uuid(), 101000, 'Москва', 'Москва', 'Ленина', '1', 'А', '10'),
(gen_random_uuid(), 102000, 'Санкт-Петербург', 'Санкт-Петербург', 'Невский', '2', 'Б', '20'),
(gen_random_uuid(), 603000, 'Казань', 'Казань', 'Баумана', '3', 'В', '30'),
(gen_random_uuid(), 630000, 'Новосибирск', 'Новосибирск', 'Ленина', '4', 'Г', '40'),
(gen_random_uuid(), 354000, 'Сочи', 'Сочи', 'Пушкина', '5', 'Д', '50');

-- Level Education
INSERT INTO leveleducation (id_leveleducation, education)
VALUES
(gen_random_uuid(), 'Бакалавр'),
(gen_random_uuid(), 'Специалитет'),
(gen_random_uuid(), 'Магистр'),
(gen_random_uuid(), 'Кандидат наук'),
(gen_random_uuid(), 'Среднее специальное'),
(gen_random_uuid(), 'Среднее образование');

-- Education Listener
INSERT INTO educationlistener (id_educationlistener, diplom_seria, diplom_number, date_given, city, region, educational_institution, speciality, level_education)
VALUES
(gen_random_uuid(), 101120, 1111321, '2010-06-01', 'Москва', 'Москва', 'МГУ', 'Информатика', (SELECT id_leveleducation FROM leveleducation LIMIT 1 OFFSET 0)),
(gen_random_uuid(), 202320, 2223212, '2011-06-01', 'Санкт-Петербург', 'Санкт-Петербург', 'СПбГУ', 'Математика', (SELECT id_leveleducation FROM leveleducation LIMIT 1 OFFSET 1)),
(gen_random_uuid(), 303120, 7434521, '2012-06-01', 'Казань', 'Казань', 'КФУ', 'Физика', (SELECT id_leveleducation FROM leveleducation LIMIT 1 OFFSET 2)),
(gen_random_uuid(), 404320, 4442344, '2013-06-01', 'Новосибирск', 'Новосибирск', 'НГУ', 'Химия', (SELECT id_leveleducation FROM leveleducation LIMIT 1 OFFSET 3)),
(gen_random_uuid(), 505310, 5554325, '2014-06-01', 'Сочи', 'Сочи', 'Сочинский университет', 'Биология', (SELECT id_leveleducation FROM leveleducation LIMIT 1 OFFSET 4));

-- Place Work
INSERT INTO placework (id_placework, name_company, job_title, all_experience, job_title_experience)
VALUES
(gen_random_uuid(), 'Компания1', 'Инженер', 10, 5),
(gen_random_uuid(), 'Компания2', 'Менеджер', 12, 6),
(gen_random_uuid(), 'Компания3', 'Разработчик', 8, 3),
(gen_random_uuid(), 'Компания4', 'Аналитик', 15, 7),
(gen_random_uuid(), 'Компания5', 'Дизайнер', 9, 4);

-- Divisions Education
INSERT INTO divisionseducation (id_divisionseducation, divisions)
VALUES
(gen_random_uuid(), 'Лингвистический центр'),
(gen_random_uuid(), 'Центр прикладных технологий');

-- Education Types
INSERT INTO educationtypes (id_educationtype, type_name)
VALUES
(gen_random_uuid(), 'Очная'),
(gen_random_uuid(), 'Заочная'),
(gen_random_uuid(), 'Дистанционная'),
(gen_random_uuid(), 'Вечерняя'),
(gen_random_uuid(), 'Смешанная');

-- Program Education
INSERT INTO programeducation (id_programeducation, name_prof_education, time_education, individual_price, group_price, campus_price, id_educationtype, id_divisionseducation)
VALUES
(gen_random_uuid(), 'Разработка игровых продуктов на Unity', 256, 146000, 146000, 146000, (SELECT id_educationtype FROM educationtypes LIMIT 1 OFFSET 0), (SELECT id_divisionseducation FROM divisionseducation LIMIT 1 OFFSET 0)),
(gen_random_uuid(), 'Разработка кроссплатформенных мобильных приложений на Flutter', 256, 146000, 146000, 146000, (SELECT id_educationtype FROM educationtypes LIMIT 1 OFFSET 1), (SELECT id_divisionseducation FROM divisionseducation LIMIT 1 OFFSET 1)),
(gen_random_uuid(), 'Разработка корпоративных приложений на Java', 256, 146000, 146000, 146000, (SELECT id_educationtype FROM educationtypes LIMIT 1 OFFSET 2), (SELECT id_divisionseducation FROM divisionseducation LIMIT 1 OFFSET 1)),
(gen_random_uuid(), 'Python: первые шаги в программировании', 256, 146000, 146000, 146000, (SELECT id_educationtype FROM educationtypes LIMIT 1 OFFSET 3), (SELECT id_divisionseducation FROM divisionseducation LIMIT 1 OFFSET 1)),
(gen_random_uuid(), 'Математика для программистов Junior', 20, 20000, 20000, 20000, (SELECT id_educationtype FROM educationtypes LIMIT 1 OFFSET 4), (SELECT id_divisionseducation FROM divisionseducation LIMIT 1 OFFSET 1));

-- Listener (частично без работы или образования)
INSERT INTO listener (id_listener, first_name, second_name, middle_name, date_of_birth, snils, contact_phone, email, id_passport, id_regaddress, id_educationlistener, id_placework)
VALUES
(gen_random_uuid(), 'Иван', 'Иванов', 'Иванович', '1990-01-01', '123-456-789 00', '+79001234567', 'ivanov@mail.ru',
 (SELECT id_passport FROM passport LIMIT 1 OFFSET 0),
 (SELECT id_regaddress FROM registrationaddress LIMIT 1 OFFSET 0),
 (SELECT id_educationlistener FROM educationlistener LIMIT 1 OFFSET 0),
 (SELECT id_placework FROM placework LIMIT 1 OFFSET 0)),

(gen_random_uuid(), 'Петр', 'Петров', 'Петрович', '1991-02-02', '223-456-789 11', '+79011234567', 'petrov@mail.ru',
 (SELECT id_passport FROM passport LIMIT 1 OFFSET 1),
 (SELECT id_regaddress FROM registrationaddress LIMIT 1 OFFSET 1),
 NULL,
 (SELECT id_placework FROM placework LIMIT 1 OFFSET 1)),

(gen_random_uuid(), 'Сергей', 'Сергеев', 'Сергеевич', '1992-03-03', '323-456-789 22', '+79021234567', 'sergeev@mail.ru',
 (SELECT id_passport FROM passport LIMIT 1 OFFSET 2),
 (SELECT id_regaddress FROM registrationaddress LIMIT 1 OFFSET 2),
 (SELECT id_educationlistener FROM educationlistener LIMIT 1 OFFSET 2),
 NULL),

(gen_random_uuid(), 'Алексей', 'Алексеев', 'Алексеевич', '1993-04-04', '423-456-789 33', '+79031234567', 'alekseev@mail.ru',
 (SELECT id_passport FROM passport LIMIT 1 OFFSET 3),
 (SELECT id_regaddress FROM registrationaddress LIMIT 1 OFFSET 3),
 NULL,
 NULL),

(gen_random_uuid(), 'Дмитрий', 'Дмитриев', 'Дмитриевич', '1994-05-05', '523-456-789 44', '+79041234567', 'dmitriev@mail.ru',
 (SELECT id_passport FROM passport LIMIT 1 OFFSET 4),
 (SELECT id_regaddress FROM registrationaddress LIMIT 1 OFFSET 4),
 (SELECT id_educationlistener FROM educationlistener LIMIT 1 OFFSET 4),
 (SELECT id_placework FROM placework LIMIT 1 OFFSET 4));

 alter table legal_entity
 add status varchar(255);