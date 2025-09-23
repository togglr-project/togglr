ALTER TABLE categories
    ADD COLUMN IF NOT EXISTS kind varchar(20) DEFAULT 'user' NOT NULL
        CHECK (kind IN ('system', 'user'));

INSERT INTO categories (id, name, slug, color, description, kind) VALUES
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
