import React, { useState } from 'react';
import {
  Box,
  Paper,
  Typography,
  Button,
  CircularProgress,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  MenuItem,
  IconButton,
  Grid,
  Chip,
  Switch,
  FormControlLabel,
  RadioGroup,
  Radio,
} from '@mui/material';
import { Add as AddIcon, Delete as DeleteIcon, Flag as FlagIcon } from '@mui/icons-material';
import { useParams } from 'react-router-dom';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import AuthenticatedLayout from '../components/AuthenticatedLayout';
import PageHeader from '../components/PageHeader';
import apiClient from '../api/apiClient';
import type { Feature, FeatureKind, Project } from '../generated/api/client';

interface ProjectResponse { project: Project }

const kindOptions: FeatureKind[] = ['boolean', 'multivariant'];

interface VariantFormItem {
  name: string;
  rollout_percent: number;
}

const ProjectPage: React.FC = () => {
  const { projectId = '' } = useParams();
  const queryClient = useQueryClient();

  const { data: projectResp, isLoading: loadingProject, error: projectError } = useQuery({
    queryKey: ['project', projectId],
    queryFn: async () => {
      const res = await apiClient.getProject(projectId);
      return res.data as ProjectResponse;
    },
    enabled: !!projectId,
  });

  const { data: features, isLoading: loadingFeatures, error: featuresError } = useQuery<Feature[]>({
    queryKey: ['project-features', projectId],
    queryFn: async () => {
      const res = await apiClient.listProjectFeatures(projectId);
      return res.data;
    },
    enabled: !!projectId,
  });

  // Create Feature Dialog state
  const [open, setOpen] = useState(false);
  const [keyValue, setKeyValue] = useState('');
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [kind, setKind] = useState<FeatureKind>('boolean');
  const [defaultVariant, setDefaultVariant] = useState('off');
  const [enabled, setEnabled] = useState(true);
  const [variants, setVariants] = useState<VariantFormItem[]>([{ name: 'control', rollout_percent: 100 }]);
  const [ruleJson, setRuleJson] = useState('');
  const [ruleTargetVariantIndex, setRuleTargetVariantIndex] = useState<number | ''>('');
  const [rulePriority, setRulePriority] = useState<number | ''>(0);
  const [formError, setFormError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  const resetForm = () => {
    setKeyValue('');
    setName('');
    setDescription('');
    setKind('boolean');
    setDefaultVariant('off');
    setEnabled(true);
    setVariants([{ name: 'control', rollout_percent: 100 }]);
    setRuleJson('');
    setRuleTargetVariantIndex('');
    setRulePriority(0);
    setFormError(null);
  };


  const createFeatureMutation = useMutation({
    mutationFn: async () => {
      setFormError(null);

      // Basic validation
      if (!keyValue.trim() || !name.trim()) {
        throw new Error('Key and Name are required');
      }
      if (kind === 'multivariant') {
        if (!variants.length) throw new Error('At least one variant is required for multivariant features');
      }

      // Create feature
      const dv = kind === 'boolean' ? (defaultVariant === 'on' ? 'on' : 'off') : (defaultVariant.trim() || variants[0]?.name || 'control');
      const featureRes = await apiClient.createProjectFeature(projectId, {
        key: keyValue.trim(),
        name: name.trim(),
        description: description.trim() || undefined,
        kind,
        default_variant: dv,
        enabled,
      });
      const feature = (featureRes.data as { feature: Feature }).feature;

      // Create variants for multivariant
      const createdVariantIds: string[] = [];
      if (kind === 'multivariant') {
        for (const v of variants) {
          const vRes = await apiClient.createFeatureFlagVariant(feature.id, {
            name: v.name,
            rollout_percent: Number(v.rollout_percent) || 0,
          });
          const flagVariant = (vRes.data as { flag_variant: { id: string } }).flag_variant;
          createdVariantIds.push(flagVariant.id);
        }
      }

      // Create rule if provided (only when multivariant and a target variant selected)
      if (kind === 'multivariant' && ruleJson.trim() && ruleTargetVariantIndex !== '') {
        let conditionObj: Record<string, unknown>;
        try {
          conditionObj = JSON.parse(ruleJson) as Record<string, unknown>;
        } catch {
          throw new Error('Invalid rule JSON');
        }
        const idx = ruleTargetVariantIndex as number;
        const variantId = createdVariantIds[idx];
        if (!variantId) throw new Error('Selected rule target variant is invalid');
        await apiClient.createFeatureRule(feature.id, {
          condition: conditionObj,
          flag_variant_id: variantId,
          priority: rulePriority === '' ? 0 : Number(rulePriority),
        });
      }

      return feature;
    },
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['project-features', projectId] });
      setOpen(false);
      resetForm();
    },
    onError: (err: unknown) => {
      let msg = 'Failed to create feature';
      if (typeof err === 'object' && err !== null) {
        const e = err as { response?: { data?: { error?: { message?: string } } }; message?: string };
        msg = e.response?.data?.error?.message || e.message || msg;
      }
      setFormError(msg);
    },
    onSettled: () => setSubmitting(false),
  });

  const handleAddVariant = () => setVariants((prev) => [...prev, { name: '', rollout_percent: 0 }]);
  const handleRemoveVariant = (index: number) => setVariants((prev) => prev.filter((_, i) => i !== index));
  const handleVariantChange = (index: number, field: keyof VariantFormItem, value: string) => {
    setVariants((prev) => prev.map((v, i) => (i === index ? { ...v, [field]: field === 'rollout_percent' ? Number(value) : value } : v)));
  };

  const project = projectResp?.project;

  return (
    <AuthenticatedLayout showBackButton backTo="/dashboard">
      <PageHeader
        title={project ? project.name : 'Project'}
        subtitle={project ? `ID: ${project.id}` : 'Project details'}
        icon={<FlagIcon />}
        gradientVariant="default"
        subtitleGradientVariant="default"
      />

      <Paper sx={{ p: 3 }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
          <Typography variant="h6" className="gradient-subtitle">Features</Typography>
          <Button variant="contained" startIcon={<AddIcon />} onClick={() => setOpen(true)}>
            Add Feature
          </Button>
        </Box>

        {(loadingProject || loadingFeatures) && (
          <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
            <CircularProgress />
          </Box>
        )}
        {(projectError || featuresError) && (
          <Typography color="error">Failed to load project or features.</Typography>
        )}

        {!loadingFeatures && features && features.length > 0 ? (
          <Grid container spacing={2}>
            {features.map((f) => (
              <Grid item xs={12} md={6} key={f.id}>
                <Paper sx={{ p: 2, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                  <Box>
                    <Typography variant="subtitle1">{f.name}</Typography>
                    <Typography variant="body2" color="text.secondary">{f.key}</Typography>
                    <Box sx={{ mt: 1, display: 'flex', gap: 1, flexWrap: 'wrap' }}>
                      <Chip size="small" label={`kind: ${f.kind}`} />
                      <Chip size="small" label={`default: ${f.default_variant}`} />
                      <Chip size="small" label={f.enabled ? 'enabled' : 'disabled'} color={f.enabled ? 'success' : 'default'} />
                    </Box>
                  </Box>
                </Paper>
              </Grid>
            ))}
          </Grid>
        ) : !loadingFeatures ? (
          <Typography variant="body2">No features yet.</Typography>
        ) : null}
      </Paper>

      {/* Create Feature Dialog */}
      <Dialog open={open} onClose={() => setOpen(false)} fullWidth maxWidth="md">
        <DialogTitle className="gradient-text-purple">Create Feature</DialogTitle>
        <DialogContent>
          <Box sx={{ display: 'grid', gridTemplateColumns: { xs: '1fr', md: '1fr 1fr' }, gap: 2, mt: 1 }}>
            <TextField label="Key" value={keyValue} onChange={(e) => setKeyValue(e.target.value)} required fullWidth />
            <TextField label="Name" value={name} onChange={(e) => setName(e.target.value)} required fullWidth />
            <TextField label="Description" value={description} onChange={(e) => setDescription(e.target.value)} fullWidth multiline minRows={2} />
            <TextField select label="Kind" value={kind} onChange={(e) => { const v = e.target.value as FeatureKind; setKind(v); if (v === 'boolean') { setDefaultVariant(defaultVariant === 'on' ? 'on' : 'off'); } }} fullWidth>
              {kindOptions.map(k => (
                <MenuItem key={k} value={k}>{k}</MenuItem>
              ))}
            </TextField>
            {kind === 'boolean' ? (
              <Box>
                <Typography variant="caption">Default Value</Typography>
                <RadioGroup
                  row
                  value={defaultVariant === 'on' ? 'on' : 'off'}
                  onChange={(e) => setDefaultVariant(e.target.value)}
                >
                  <FormControlLabel value="on" control={<Radio />} label="on" />
                  <FormControlLabel value="off" control={<Radio />} label="off" />
                </RadioGroup>
              </Box>
            ) : (
              <TextField label="Default Variant" value={defaultVariant} onChange={(e) => setDefaultVariant(e.target.value)} fullWidth helperText="Should match one of variants" />
            )}
            <FormControlLabel
              control={<Switch checked={enabled} onChange={(e) => setEnabled(e.target.checked)} />}
              label="Enabled"
            />
          </Box>

          {kind === 'multivariant' && (
            <Box sx={{ mt: 3 }}>
              <Typography variant="subtitle1" sx={{ mb: 1 }}>Variants</Typography>
              {variants.map((v, i) => (
                <Box key={i} sx={{ display: 'grid', gridTemplateColumns: '2fr 1fr auto', gap: 1, alignItems: 'center', mb: 1 }}>
                  <TextField label="Name" value={v.name} onChange={(e) => handleVariantChange(i, 'name', e.target.value)} fullWidth />
                  <TextField label="Rollout %" type="number" value={v.rollout_percent} onChange={(e) => handleVariantChange(i, 'rollout_percent', e.target.value)} fullWidth />
                  <IconButton aria-label="delete" onClick={() => handleRemoveVariant(i)}>
                    <DeleteIcon />
                  </IconButton>
                </Box>
              ))}
              <Button size="small" startIcon={<AddIcon />} onClick={handleAddVariant}>Add Variant</Button>

              <Box sx={{ mt: 2 }}>
                <Typography variant="subtitle1" sx={{ mb: 1 }}>Rule (optional)</Typography>
                <TextField label="Rule JSON" value={ruleJson} onChange={(e) => setRuleJson(e.target.value)} fullWidth multiline minRows={3} placeholder='{"user_id": 123}' />
                <TextField label="Rule Priority" type="number" value={rulePriority} onChange={(e) => setRulePriority(e.target.value === '' ? '' : Number(e.target.value))} fullWidth sx={{ mt: 1 }} helperText="Lower numbers run first" />
                <TextField select fullWidth sx={{ mt: 1 }} label="Rule Target Variant" value={ruleTargetVariantIndex} onChange={(e) => setRuleTargetVariantIndex(e.target.value === '' ? '' : Number(e.target.value))} helperText="Which variant to return when rule matches">
                  <MenuItem value="">None</MenuItem>
                  {variants.map((v, i) => (
                    <MenuItem key={i} value={i}>{v.name || `Variant #${i+1}`}</MenuItem>
                  ))}
                </TextField>
              </Box>
            </Box>
          )}

          {formError && (
            <Typography color="error" sx={{ mt: 2 }}>{formError}</Typography>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => { setOpen(false); resetForm(); }}>
            Cancel
          </Button>
          <Button variant="contained" onClick={() => { setSubmitting(true); createFeatureMutation.mutate(); }} disabled={submitting}>
            {submitting ? 'Creating...' : 'Create'}
          </Button>
        </DialogActions>
      </Dialog>
    </AuthenticatedLayout>
  );
};

export default ProjectPage;
