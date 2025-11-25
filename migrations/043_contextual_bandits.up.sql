-- Add contextual bandit algorithms
INSERT INTO algorithms (slug, name, kind, description, default_settings) VALUES
    ('lin_ucb', 'Linear UCB', 'contextual_bandit', 
     'Linear Upper Confidence Bound algorithm. Uses linear regression with UCB exploration to learn optimal actions based on user context (country, device, etc.)',
     '{"alpha": 1.0, "feature_dim": 32}'::jsonb),
    
    ('contextual_thompson', 'Contextual Thompson Sampling', 'contextual_bandit',
     'Thompson Sampling with linear model. Uses Bayesian approach to balance exploration and exploitation based on user context.',
     '{"prior_variance": 1.0, "feature_dim": 32}'::jsonb),
    
    ('contextual_epsilon', 'Contextual Epsilon-Greedy', 'contextual_bandit',
     'Epsilon-Greedy with linear model. Simple exploration strategy that uses context features for predictions.',
     '{"epsilon": 0.1, "feature_dim": 32}'::jsonb)
ON CONFLICT (slug) DO UPDATE SET
    name = EXCLUDED.name,
    kind = EXCLUDED.kind,
    description = EXCLUDED.description,
    default_settings = EXCLUDED.default_settings,
    updated_at = NOW();

