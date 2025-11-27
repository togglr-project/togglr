-- Remove constraint
ALTER TABLE public.feature_algorithms DROP CONSTRAINT IF EXISTS chk_algorithm_type;

-- Remove custom_algorithm_id column
DROP INDEX IF EXISTS idx_feature_algorithms_custom;
ALTER TABLE public.feature_algorithms DROP COLUMN IF EXISTS custom_algorithm_id;

-- Drop custom algorithm stats table
DROP TABLE IF EXISTS monitoring.custom_algorithm_stats;

-- Drop custom algorithms table
DROP TABLE IF EXISTS public.custom_algorithms;

