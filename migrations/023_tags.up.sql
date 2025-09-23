-- ========================================
-- Categories (global, cross-project)
-- ========================================
CREATE TABLE categories (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    name varchar(100) NOT NULL UNIQUE,
    slug varchar(100) NOT NULL UNIQUE,
    description varchar(300),
    color varchar(7), -- hex like #AABBCC
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz DEFAULT now() NOT NULL
);

-- ========================================
-- Tags (project-local)
-- ========================================
CREATE TABLE tags (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    project_id uuid NOT NULL REFERENCES projects ON DELETE CASCADE,
    name varchar(100) NOT NULL,
    slug varchar(100) NOT NULL,
    color varchar(7), -- hex like #AABBCC
    description varchar(300),
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz DEFAULT now() NOT NULL,
    UNIQUE (project_id, slug)
);

-- индекс для быстрого поиска по проекту и имени
CREATE INDEX idx_tags_project_name ON tags (project_id, name);

-- ========================================
-- Tag <-> Category (many-to-many)
-- ========================================
CREATE TABLE tag_categories (
    tag_id uuid NOT NULL REFERENCES tags ON DELETE CASCADE,
    category_id uuid NOT NULL REFERENCES categories ON DELETE CASCADE,
    PRIMARY KEY (tag_id, category_id)
);

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
    t.id AS tag_id,
    t.name AS tag_name,
    c.id AS category_id,
    c.name AS category_name
FROM features f
JOIN feature_tags ft ON f.id = ft.feature_id
JOIN tags t ON ft.tag_id = t.id
JOIN tag_categories tc ON tc.tag_id = t.id
JOIN categories c ON tc.category_id = c.id;
