-- ========================================
-- Categories (global, cross-project)
-- ========================================
CREATE TABLE categories (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    name varchar(100) NOT NULL UNIQUE,
    slug varchar(100) NOT NULL UNIQUE,
    description varchar(300),
    color varchar(7), -- hex like #AABBCC
    kind varchar(20) DEFAULT 'user' NOT NULL
        CHECK (kind IN ('system', 'user', 'nocopy')),
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz DEFAULT now() NOT NULL
);

-- ========================================
-- Tags (project-local)
-- ========================================
CREATE TABLE tags (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    project_id uuid NOT NULL REFERENCES projects ON DELETE CASCADE,
    category_id uuid REFERENCES categories ON DELETE SET NULL,
    name varchar(100) NOT NULL,
    slug varchar(100) NOT NULL,
    color varchar(7), -- hex like #AABBCC
    description varchar(300),
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz DEFAULT now() NOT NULL,
    UNIQUE (project_id, slug)
);

CREATE INDEX idx_tags_project_name ON tags (project_id, name);
CREATE INDEX idx_tags_category ON tags (category_id);

-- ========================================
-- Feature <-> Tag (many-to-many)
-- ========================================
CREATE TABLE feature_tags (
    feature_id uuid NOT NULL REFERENCES features ON DELETE CASCADE,
    tag_id uuid NOT NULL REFERENCES tags ON DELETE CASCADE,
    created_at timestamptz DEFAULT now() NOT NULL,
    PRIMARY KEY (feature_id, tag_id)
);

CREATE INDEX idx_feature_tags_tag_id ON feature_tags (tag_id);
CREATE INDEX idx_feature_tags_feature_id ON feature_tags (feature_id);

-- ========================================
-- View: features with categories (via tags)
-- ========================================
CREATE VIEW v_feature_categories AS
SELECT
    f.id AS feature_id,
    f.key AS feature_key,
    f.project_id AS project_id,
    t.id AS tag_id,
    t.name AS tag_name,
    t.slug AS tag_slug,
    c.id AS category_id,
    c.name AS category_name,
    c.slug AS category_slug
FROM features f
JOIN feature_tags ft ON f.id = ft.feature_id
JOIN tags t ON ft.tag_id = t.id
JOIN categories c ON t.category_id = c.id;


INSERT INTO categories (id, name, slug, color, description, kind) VALUES
-- special/meta
-- (gen_random_uuid(), 'User Tags', 'user-tags', '#10B981', 'User-defined tags', 'nocopy'),

-- safety / governance
(gen_random_uuid(), 'Critical', 'critical', '#DC2626', 'Critical feature, excluded from algorithms', 'system'),
(gen_random_uuid(), 'Auto-Disable', 'auto-disable', '#F97316', 'Feature automatically disabled on high error rate', 'system'),
(gen_random_uuid(), 'Guarded', 'guarded', '#F59E0B', 'Feature requires manual approval for changes', 'system')
-- (gen_random_uuid(), 'Security', 'security', '#7C3AED', 'Security or access control feature', 'system'),

-- domains
-- (gen_random_uuid(), 'UI/UX', 'ui-ux', '#06B6D4', 'UI or UX related feature', 'system'),
-- (gen_random_uuid(), 'Backend', 'backend', '#4B5563', 'Backend logic feature', 'system'),
-- (gen_random_uuid(), 'Infra', 'infra', '#9CA3AF', 'Infrastructure or DevOps feature', 'system'),
-- (gen_random_uuid(), 'Ads Campaign', 'ads-campaign', '#EC4899', 'Advertising campaign feature', 'system'),
-- (gen_random_uuid(), 'Pricing', 'pricing', '#84CC16', 'Pricing or discount related feature', 'system'),
-- (gen_random_uuid(), 'Compliance', 'compliance', '#0EA5E9', 'Regulatory or compliance-related feature', 'system')
ON CONFLICT (slug) DO NOTHING;

CREATE OR REPLACE FUNCTION init_project_tags(p_project_id uuid)
    RETURNS void AS $$
BEGIN
    INSERT INTO tags (project_id, category_id, name, slug, color, description)
    SELECT
        p_project_id,
        c.id,
        c.name,
        c.slug,
        c.color,
        c.description
    FROM categories c
    WHERE c.kind <> 'nocopy'
    ON CONFLICT (project_id, slug) DO NOTHING;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION trigger_init_project_tags()
    RETURNS trigger AS $$
BEGIN
    PERFORM init_project_tags(NEW.id);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_init_project_tags
    AFTER INSERT ON projects
    FOR EACH ROW
EXECUTE FUNCTION trigger_init_project_tags();

CREATE OR REPLACE FUNCTION init_category_tags(p_category_id uuid)
    RETURNS void AS $$
BEGIN
    INSERT INTO tags (project_id, category_id, name, slug, color, description)
    SELECT
        p.id,
        c.id,
        c.name,
        c.slug,
        c.color,
        c.description
    FROM projects p
             CROSS JOIN categories c
    WHERE c.id = p_category_id
      AND c.kind <> 'nocopy'
    ON CONFLICT (project_id, slug) DO NOTHING;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION trigger_init_category_tags()
    RETURNS trigger AS $$
BEGIN
    PERFORM init_category_tags(NEW.id);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_init_category_tags
    AFTER INSERT ON categories
    FOR EACH ROW
EXECUTE FUNCTION trigger_init_category_tags();
