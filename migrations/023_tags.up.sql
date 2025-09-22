-- tags table
CREATE TABLE feature_tags (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    project_id uuid NOT NULL REFERENCES projects ON DELETE CASCADE,
    name varchar(100) NOT NULL,
    slug varchar(100) NOT NULL,
    color varchar(7), -- hex like #AABBCC
    description varchar(300),
    parent_id uuid REFERENCES feature_tags ON DELETE SET NULL,
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz DEFAULT now() NOT NULL,
    UNIQUE (project_id, slug)
);

-- index for fast lookup
CREATE INDEX idx_feature_tags_project_name ON feature_tags (project_id, name);
CREATE INDEX idx_feature_tags_project_slug ON feature_tags (project_id, slug);

-- join table feature <-> tag
CREATE TABLE features_feature_tags (
    feature_id uuid NOT NULL REFERENCES features ON DELETE CASCADE,
    tag_id uuid NOT NULL REFERENCES feature_tags ON DELETE CASCADE,
    created_at timestamptz DEFAULT now() NOT NULL,
    PRIMARY KEY (feature_id, tag_id)
);

CREATE INDEX idx_fft_tag_id ON features_feature_tags (tag_id);
CREATE INDEX idx_fft_feature_id ON features_feature_tags (feature_id);

-- optional: tags for segments
CREATE TABLE segments_tags (
    segment_id uuid NOT NULL REFERENCES segments ON DELETE CASCADE,
    tag_id uuid NOT NULL REFERENCES feature_tags ON DELETE CASCADE,
    PRIMARY KEY (segment_id, tag_id)
);
