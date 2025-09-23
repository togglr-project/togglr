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
(gen_random_uuid(), 'User Tags', 'user-tags', '#10B981', 'User-defined tags', 'nocopy'),
(gen_random_uuid(), 'Experiment', 'experiment', '#3B82F6', 'Feature participates in experiment (A/B, bandit, etc.)', 'system'),
-- (gen_random_uuid(), 'Bandit', 'bandit', '#8B5CF6', 'Feature controlled by multi-armed bandit algorithm', 'system'),
-- (gen_random_uuid(), 'Contextual Bandit', 'contextual-bandit', '#6366F1', 'Feature controlled by contextual bandit', 'system'),
-- (gen_random_uuid(), 'ML-Driven', 'ml-driven', '#10B981', 'Feature rollout managed by ML model', 'system'),

(gen_random_uuid(), 'Critical', 'critical', '#DC2626', 'Critical feature, excluded from algorithms', 'system'),
(gen_random_uuid(), 'Auto-Disable', 'auto-disable', '#F97316', 'Feature automatically disabled on high error rate', 'system'),
(gen_random_uuid(), 'Guarded', 'guarded', '#F59E0B', 'Feature requires manual approval for changes', 'system'),

(gen_random_uuid(), 'UI/UX', 'ui-ux', '#06B6D4', 'UI or UX related feature', 'system'),
(gen_random_uuid(), 'Backend', 'backend', '#4B5563', 'Backend logic feature', 'system'),
(gen_random_uuid(), 'Infra', 'infra', '#9CA3AF', 'Infrastructure feature', 'system')
-- (gen_random_uuid(), 'Ads Campaign', 'ads-campaign', '#EC4899', 'Advertising campaign feature', 'system'),
-- (gen_random_uuid(), 'Pricing', 'pricing', '#84CC16', 'Pricing or discount related feature', 'system')
ON CONFLICT (slug) DO NOTHING;
