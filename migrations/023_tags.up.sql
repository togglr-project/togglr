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
        CHECK (kind IN ('system', 'user', 'domain')),
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

-- All project tags with category info
create or replace view v_project_tags as
select
    t.id as tag_id,
    t.project_id,
    t.name as tag_name,
    t.slug as tag_slug,
    t.color,
    t.description,
    c.name as category_name,
    c.slug as category_slug,
    c.kind as category_kind
from tags t
join categories c on c.id = t.category_id;

-- ======================================
-- Preset Categories
-- ======================================
insert into categories (id, name, slug, color, description, kind)
values
    -- Safety (only one system category)
    (gen_random_uuid(), 'Safety', 'safety', '#F59E0B', 'Safety and governance related tags', 'system'),

    -- Domain categories
    (gen_random_uuid(), 'UI/UX', 'ui-ux', '#06B6D4', 'UI or UX related features', 'domain'),
    (gen_random_uuid(), 'Backend', 'backend', '#4B5563', 'Backend logic features', 'domain'),
    (gen_random_uuid(), 'Infra', 'infra', '#9CA3AF', 'Infrastructure features', 'domain')
--     (gen_random_uuid(), 'Experiment', 'experiment', '#3B82F6', 'Features used in experiments (A/B, bandit, etc.)', 'system')
on conflict (slug) do nothing;

-- ======================================
-- Functions & Triggers
-- ======================================

-- Insert default safety tags for each project
create or replace function init_project_safety_tags()
    returns trigger as $$
declare
    safety_id uuid;
begin
    select id into safety_id from categories where slug = 'safety';

    insert into tags (project_id, category_id, name, slug, color, description)
    values
        (new.id, safety_id, 'Critical', 'critical', '#DC2626', 'Critical feature, excluded from algorithms'),
        (new.id, safety_id, 'Auto-Disable', 'auto-disable', '#F97316', 'Feature auto-disabled on high error rate'),
        (new.id, safety_id, 'Guarded', 'guarded', '#F59E0B', 'Feature requires manual approval for changes');
    return new;
end;
$$ language plpgsql;

-- Trigger: after project creation
create trigger trg_init_project_safety_tags
    after insert on projects
    for each row
execute function init_project_safety_tags();
