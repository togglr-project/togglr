-- Hill Climbing Optimizer
INSERT INTO algorithms (slug, name, description, kind, default_settings)
VALUES (
    'hill_climb',
    'Hill Climbing',
    'Gradually adjusts a numeric feature value up or down to improve reward.',
    'optimizer',
    '{
         "step": 0.05,
         "direction": 1
    }'::jsonb
);

-- Simulated Annealing
INSERT INTO algorithms (slug, name, description, kind, default_settings)
VALUES (
    'simulated_annealing',
    'Simulated Annealing',
    'Explores parameter space with random jumps that decrease over time to avoid local minima.',
    'optimizer',
    '{
        "temp": 1.0,
        "cooling": 0.95,
        "step_scale": 0.1
    }'::jsonb
);

-- Cross-Entropy Method (CEM)
INSERT INTO algorithms (slug, name, description, kind, default_settings)
VALUES (
    'cem',
    'Cross-Entropy Method',
    'Samples candidate values, keeps top-performing percentiles, and refines the search distribution.',
    'optimizer',
    '{
        "population_size": 20,
        "elite_fraction": 0.2
    }'::jsonb
);

-- Bayesian Optimization (Gaussian Process)
INSERT INTO algorithms (slug, name, description, kind, default_settings)
VALUES (
    'bayes_opt',
    'Bayesian Optimization',
    'Models reward as a Gaussian Process and selects values maximizing expected improvement.',
    'optimizer',
    '{
        "noise": 0.01
    }'::jsonb
);

-- PID Controller
INSERT INTO algorithms (slug, name, description, kind, default_settings)
VALUES (
    'pid_controller',
    'PID Controller',
    'Regulates a numeric feature to maintain a target metric using proportional-integral-derivative feedback.',
    'optimizer',
    '{
        "kp": 0.2,
        "ki": 0.05,
        "kd": 0.01
    }'::jsonb
);
