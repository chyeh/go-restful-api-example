SET NAMES 'UTF8';
CREATE TABLE IF NOT EXISTS recipe(
    r_id SERIAL PRIMARY KEY,
    r_name VARCHAR(512) NOT NULL,
    r_prep_time SMALLINT,
    r_difficulty SMALLINT,
    r_vegetarian BOOLEAN NOT NULL,
    r_rating REAL NOT NULL DEFAULT 0.0,
    r_rated_num INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS hellofresh_user(
    hu_id SERIAL PRIMARY KEY,
    hu_account VARCHAR(32) NOT NULL UNIQUE,
    hu_access_token VARCHAR(32) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS hellofresh_user_recipe(
    hur_hu_id INTEGER,
    hur_r_id INTEGER,
    CONSTRAINT pk_hellofresh_user_recipe PRIMARY KEY(hur_hu_id, hur_r_id),
    CONSTRAINT fk_hellofresh_user_recipe__hellofresh_user FOREIGN KEY
        (hur_hu_id) REFERENCES hellofresh_user(hu_id)
        ON DELETE CASCADE
        ON UPDATE RESTRICT,
    CONSTRAINT fk_hellofresh_user_recipe__recipe FOREIGN KEY
		(hur_r_id) REFERENCES recipe(r_id)
		ON DELETE CASCADE
		ON UPDATE RESTRICT
);
