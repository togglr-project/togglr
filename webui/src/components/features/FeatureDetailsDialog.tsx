import React from 'react';
import { Box, Chip, CircularProgress, Dialog, DialogActions, DialogContent, DialogTitle, Button, Typography } from '@mui/material';
import { useQuery } from '@tanstack/react-query';
import apiClient from '../../api/apiClient';
import type { Feature, FeatureDetailsResponse } from '../../generated/api/client';

export interface FeatureDetailsDialogProps {
  open: boolean;
  onClose: () => void;
  feature: Feature | null;
}

const FeatureDetailsDialog: React.FC<FeatureDetailsDialogProps> = ({ open, onClose, feature }) => {
  const { data: featureDetails, isLoading, error } = useQuery<FeatureDetailsResponse>({
    queryKey: ['feature-details', feature?.id],
    queryFn: async () => {
      const res = await apiClient.getFeature(feature!.id);
      return res.data as FeatureDetailsResponse;
    },
    enabled: open && !!feature?.id,
  });

  const getVariantName = (id: string) => {
    const arr = featureDetails?.variants || [];
    const found = arr.find(v => v.id === id);
    return found ? (found.name || found.id) : id;
    };

  return (
    <Dialog open={open} onClose={onClose} fullWidth maxWidth="md">
      <DialogTitle className="gradient-text-purple">Feature Details</DialogTitle>
      <DialogContent>
        {!feature || isLoading ? (
          <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
            <CircularProgress />
          </Box>
        ) : error ? (
          <Typography color="error">Failed to load feature details.</Typography>
        ) : featureDetails ? (
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, mt: 1 }}>
            <Box>
              <Typography variant="h6">{featureDetails.feature.name}</Typography>
              <Typography variant="body2" color="text.secondary">Key: {featureDetails.feature.key}</Typography>
              {featureDetails.feature.description && (
                <Typography variant="body2" sx={{ mt: 1 }}>{featureDetails.feature.description}</Typography>
              )}
              <Box sx={{ mt: 1, display: 'flex', gap: 1, flexWrap: 'wrap' }}>
                <Chip size="small" label={`id: ${featureDetails.feature.id}`} variant="outlined" />
                <Chip size="small" label={`kind: ${featureDetails.feature.kind}`} />
                <Chip size="small" label={`default: ${featureDetails.feature.default_variant}`} />
                <Chip size="small" label={featureDetails.feature.enabled ? 'enabled' : 'disabled'} color={featureDetails.feature.enabled ? 'success' : 'default'} />
              </Box>
            </Box>

            <Box>
              <Typography variant="subtitle1" sx={{ mb: 1 }}>Variants</Typography>
              {featureDetails.variants && featureDetails.variants.length > 0 ? (
                <Box>
                  {featureDetails.variants.map((v) => (
                    <Box key={v.id} sx={{ display: 'grid', gridTemplateColumns: { xs: '2fr 1fr' }, gap: 1, mb: 0.5 }}>
                      <Typography variant="body2">{v.name}</Typography>
                      <Typography variant="body2" color="text.secondary">{v.rollout_percent}%</Typography>
                    </Box>
                  ))}
                </Box>
              ) : (
                <Typography variant="body2" color="text.secondary">No variants</Typography>
              )}
            </Box>

            <Box>
              <Typography variant="subtitle1" sx={{ mb: 1 }}>Rules</Typography>
              {featureDetails.rules && featureDetails.rules.length > 0 ? (
                <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1 }}>
                  {featureDetails.rules.sort((a, b) => a.priority - b.priority).map((r) => (
                    <Box key={r.id} sx={{ border: '1px solid', borderColor: 'divider', borderRadius: 1, p: 1.5 }}>
                      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 1 }}>
                        <Chip size="small" label={`priority: ${r.priority}`} />
                        <Chip size="small" label={`target: ${getVariantName(r.flag_variant_id)}`} />
                        <Chip size="small" label={`id: ${r.id}`} variant="outlined" />
                      </Box>
                      <Typography variant="body2" color="text.secondary" sx={{ mb: 0.5 }}>Conditions:</Typography>
                      {r.conditions.map((c, idx) => (
                        <Box key={idx} sx={{ display: 'grid', gridTemplateColumns: { xs: '1fr 1fr' , md: '1.2fr 0.8fr 1.5fr' }, gap: 1, mb: 0.5, alignItems: 'center' }}>
                          <Typography variant="body2">{c.attribute}</Typography>
                          <Typography variant="body2" color="text.secondary">{c.operator}</Typography>
                          <Typography variant="body2" sx={{ fontFamily: 'monospace' }}>
                            {typeof c.value === 'string' ? c.value : JSON.stringify(c.value)}
                          </Typography>
                        </Box>
                      ))}
                    </Box>
                  ))}
                </Box>
              ) : (
                <Typography variant="body2" color="text.secondary">No rules</Typography>
              )}
            </Box>
          </Box>
        ) : null}
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>Close</Button>
      </DialogActions>
    </Dialog>
  );
};

export default FeatureDetailsDialog;
