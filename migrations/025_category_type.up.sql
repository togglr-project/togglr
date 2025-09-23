ALTER TABLE categories
    ADD COLUMN category_type varchar(30) NOT NULL DEFAULT 'user'
        CHECK (category_type IN ('domain', 'safety', 'user'));

UPDATE categories
SET category_type = 'safety'
WHERE slug IN ('critical', 'auto-disable', 'guarded', 'security');

UPDATE categories
SET category_type = 'domain'
WHERE slug IN ('ui-ux', 'backend', 'infra', 'pricing', 'ads-campaign', 'compliance');

UPDATE categories
SET category_type = 'user'
WHERE slug = 'user-tags';

CREATE OR REPLACE FUNCTION enforce_safety_category_tags()
    RETURNS trigger LANGUAGE plpgsql AS $$
DECLARE
    cat_type text;
    existing_count int;
BEGIN
    -- Узнаём тип категории
    SELECT category_type INTO cat_type
    FROM categories
    WHERE id = NEW.category_id;

    IF cat_type = 'safety' THEN
        -- Проверяем, сколько тегов уже есть в этой категории
        SELECT COUNT(*) INTO existing_count
        FROM tags
        WHERE category_id = NEW.category_id;

        IF existing_count >= 1 THEN
            RAISE EXCEPTION 'Safety category % can only have one tag', NEW.category_id;
        END IF;
    END IF;

    RETURN NEW;
END;
$$;

CREATE TRIGGER trg_enforce_safety_category_tags
    BEFORE INSERT ON tags
    FOR EACH ROW
EXECUTE FUNCTION enforce_safety_category_tags();

CREATE OR REPLACE FUNCTION prevent_tag_reassign_to_safety()
    RETURNS trigger LANGUAGE plpgsql AS $$
DECLARE
    cat_type text;
BEGIN
    SELECT category_type INTO cat_type
    FROM categories
    WHERE id = NEW.category_id;

    IF cat_type = 'safety' AND NEW.category_id <> OLD.category_id THEN
        RAISE EXCEPTION 'Cannot move a tag into a safety category';
    END IF;

    RETURN NEW;
END;
$$;

CREATE TRIGGER trg_prevent_tag_reassign_to_safety
    BEFORE UPDATE ON tags
    FOR EACH ROW
    WHEN (NEW.category_id IS DISTINCT FROM OLD.category_id)
EXECUTE FUNCTION prevent_tag_reassign_to_safety();
