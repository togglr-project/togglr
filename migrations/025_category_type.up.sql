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
