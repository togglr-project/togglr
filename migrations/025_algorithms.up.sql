CREATE TABLE algorithms (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name varchar(100) NOT NULL UNIQUE,
    slug varchar(100) NOT NULL UNIQUE,
    description varchar(300),
    default_settings jsonb NOT NULL DEFAULT '{}'::jsonb,
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz DEFAULT now() NOT NULL
);

CREATE TABLE feature_algorithms (
    feature_id uuid NOT NULL REFERENCES features ON DELETE CASCADE,
    algorithm_id uuid NOT NULL REFERENCES algorithms ON DELETE CASCADE,
    settings jsonb,
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz DEFAULT now() NOT NULL,
    PRIMARY KEY (feature_id, algorithm_id)
);

INSERT INTO algorithms (id, name, slug, description, default_settings) VALUES
(gen_random_uuid(), 'Epsilon-Greedy', 'epsilon-greedy',
'Chooses mostly best variant with some random exploration',
'{"epsilon": 0.1}'),

(gen_random_uuid(), 'Thompson Sampling', 'thompson-sampling',
'Bayesian multi-armed bandit using beta distribution',
'{"prior_alpha": 1, "prior_beta": 1}'),

(gen_random_uuid(), 'Upper Confidence Bound', 'ucb',
'Selects variant with highest upper confidence bound',
'{"confidence": 2.0}');

CREATE OR REPLACE FUNCTION apply_default_algorithm_settings()
    RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    -- INSERT case: if settings are not passed → copy default_settings
    IF TG_OP = 'INSERT' THEN
        IF NEW.settings IS NULL THEN
            SELECT a.default_settings
            INTO NEW.settings
            FROM algorithms a
            WHERE a.id = NEW.algorithm_id;

            IF NEW.settings IS NULL THEN
                NEW.settings := '{}'::jsonb;
            END IF;
        END IF;
    END IF;

    -- UPDATE case: if algorithm_id is changed → reset to default_settings
    IF TG_OP = 'UPDATE' AND NEW.algorithm_id <> OLD.algorithm_id THEN
        SELECT a.default_settings
        INTO NEW.settings
        FROM algorithms a
        WHERE a.id = NEW.algorithm_id;

        IF NEW.settings IS NULL THEN
            NEW.settings := '{}'::jsonb;
        END IF;
    END IF;

    RETURN NEW;
END;
$$;

-- INSERT: apply default settings, if they are not passed
CREATE TRIGGER trg_apply_default_algorithm_settings_insert
    BEFORE INSERT ON feature_algorithms
    FOR EACH ROW
EXECUTE FUNCTION apply_default_algorithm_settings();

-- UPDATE: reset settings when algorithm_id is changed
CREATE TRIGGER trg_apply_default_algorithm_settings_update
    BEFORE UPDATE OF algorithm_id ON feature_algorithms
    FOR EACH ROW
EXECUTE FUNCTION apply_default_algorithm_settings();
