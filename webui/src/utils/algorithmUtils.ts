import type { FeatureAlgorithm } from '../generated/api/client';

export const getFirstEnabledAlgorithm = (algorithms?: FeatureAlgorithm[]): FeatureAlgorithm | null => {
  if (!algorithms || !Array.isArray(algorithms)) {
    return null;
  }
  
  return algorithms.find(algorithm => algorithm.enabled === true) || null;
};

export const getFirstEnabledAlgorithmSlug = (algorithms?: FeatureAlgorithm[]): string | null => {
  const algorithm = getFirstEnabledAlgorithm(algorithms);
  return algorithm?.algorithm_slug || null;
};
