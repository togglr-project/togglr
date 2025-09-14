import React, { useState, useEffect } from 'react';
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

// UUID generator (uses crypto.randomUUID when available)
const genId = (): string => {
  const g = (typeof crypto !== 'undefined' && typeof (crypto as any).randomUUID === 'function')
    ? (crypto as any).randomUUID()
    : 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
        const r = Math.random() * 16 | 0;
        const v = c === 'x' ? r : (r & 0x3) | 0x8;
        return v.toString(16);
      });
  return g;
};

type OperatorOption = 'eq' | 'neq' | 'in' | 'not_in' | 'gt' | 'gte' | 'lt' | 'lte' | 'regex' | 'percentage';
interface RuleConditionItem { attribute: string; operator: OperatorOption; value: string }
interface RuleFormItem { id: string; flag_variant_id: string; priority: number | ''; conditions: RuleConditionItem[] }

interface VariantFormItem {
  id: string;
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
  const [variants, setVariants] = useState<VariantFormItem[]>([{ id: genId(), name: 'control', rollout_percent: 100 }]);
  const [rules, setRules] = useState<RuleFormItem[]>([]);
  const [formError, setFormError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  // Variants validation helpers
  const totalRollout = variants.reduce((sum, v) => sum + (Number.isFinite(Number(v.rollout_percent)) ? Number(v.rollout_percent) : 0), 0);
  const rolloutSumValid = Math.round(totalRollout) === 100;
  const variantsValuesValid = variants.every((v) => v.name.trim().length > 0 && Number.isFinite(Number(v.rollout_percent)) && v.rollout_percent > 0 && v.rollout_percent <= 100);
  const variantsValid = variantsValuesValid && rolloutSumValid;
  const variantNames = variants.map(v => v.name.trim()).filter(n => n.length > 0);
  const hasAtLeastTwoVariants = variantNames.length >= 2;
  const defaultVariantValid = kind !== 'multivariant' || (variantNames.includes(defaultVariant));

  useEffect(() => {
    if (kind === 'multivariant') {
      const names = variantNames;
      if (!names.includes(defaultVariant)) {
        setDefaultVariant(names[0] || '');
      }
    }
  }, [kind, variants]);

  // Ensure rules reference existing variants when variants change
  useEffect(() => {
    setRules((prev) => prev.map((r) => (
      variants.some((v) => v.id === r.flag_variant_id)
        ? r
        : { ...r, flag_variant_id: variants[0]?.id || '' }
    )));
  }, [variants]);

  const resetForm = () => {
    setKeyValue('');
    setName('');
    setDescription('');
    setKind('boolean');
    setDefaultVariant('off');
    setEnabled(true);
    setVariants([{ id: genId(), name: 'control', rollout_percent: 100 }]);
    setRules([]);
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
        if (!hasAtLeastTwoVariants) throw new Error('At least two variants are required for multivariant features');
        if (!variantsValid) throw new Error('Variants must have names, rollout between 1 and 100, and total rollout must equal 100');
        const names = variants.map((v) => v.name.trim()).filter((n) => n.length > 0);
        if (!names.includes(defaultVariant.trim())) throw new Error('Default Variant must be one of the variants');
      }

      // Create feature with inline variants and rules (if multivariant)
      const dv = kind === 'boolean' ? (defaultVariant === 'on' ? 'on' : 'off') : defaultVariant.trim();

      let inlineVariants: { id: string; name: string; rollout_percent: number }[] | undefined;
      let inlineRules: { id: string; conditions: any[]; flag_variant_id: string; priority?: number }[] | undefined;

      if (kind === 'multivariant') {
        // Validate rules if present
        if (rules.length > 0 && !rulesValid) {
          throw new Error('Please fix rules: select a target variant and add at least one condition with attribute');
        }
        inlineVariants = variants.map((v) => ({ id: v.id, name: v.name.trim(), rollout_percent: Number(v.rollout_percent) || 0 }));
        if (rules.length > 0) {
          inlineRules = rules.map((r) => ({
            id: r.id,
            flag_variant_id: r.flag_variant_id,
            priority: r.priority === '' ? 0 : Number(r.priority),
            conditions: r.conditions.filter((c) => c.attribute.trim()).map((c) => {
              let val: any = c.value;
              const trimmed = (c.value || '').trim();
              if (trimmed.length > 0) {
                try { val = JSON.parse(trimmed); } catch { val = c.value; }
              } else {
                val = '';
              }
              return { attribute: c.attribute.trim(), operator: c.operator, value: val };
            }),
          }));
        }
      }

      const featureRes = await apiClient.createProjectFeature(projectId, {
        key: keyValue.trim(),
        name: name.trim(),
        description: description.trim() || undefined,
        kind,
        default_variant: dv,
        enabled,
        variants: inlineVariants,
        rules: inlineRules,
      } as any);
      const feature = (featureRes.data as { feature: Feature }).feature;
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

  const handleAddVariant = () => setVariants((prev) => [...prev, { id: genId(), name: '', rollout_percent: 1 }]);
  const handleRemoveVariant = (index: number) => setVariants((prev) => prev.filter((_, i) => i !== index));
  const handleVariantChange = (index: number, field: keyof VariantFormItem, value: string) => {
    setVariants((prev) => prev.map((v, i) => (i === index ? { ...v, [field]: field === 'rollout_percent' ? Number(value) : value } : v)));
  };

  // Rules helpers/state for Create Feature dialog
  const ruleOperatorOptions: OperatorOption[] = ['eq','neq','in','not_in','gt','gte','lt','lte','regex','percentage'];

  const addRule = () => setRules((prev) => [
    ...prev,
    { id: genId(), flag_variant_id: variants[0]?.id || '', priority: 0, conditions: [{ attribute: '', operator: 'eq', value: '' }] }
  ]);
  const removeRule = (index: number) => setRules((prev) => prev.filter((_, i) => i !== index));
  const updateRule = (index: number, patch: Partial<RuleFormItem>) => setRules((prev) => prev.map((r, i) => (i === index ? { ...r, ...patch } : r)));
  const addRuleConditionInline = (ruleIndex: number) => setRules((prev) => prev.map((r, i) => (
    i === ruleIndex ? { ...r, conditions: [...r.conditions, { attribute: '', operator: 'eq', value: '' }] } : r
  )));
  const removeRuleConditionInline = (ruleIndex: number, condIndex: number) => setRules((prev) => prev.map((r, i) => (
    i === ruleIndex ? { ...r, conditions: r.conditions.filter((_, ci) => ci !== condIndex) } : r
  )));
  const changeRuleConditionInline = (ruleIndex: number, condIndex: number, field: keyof RuleConditionItem, value: string) => setRules((prev) => prev.map((r, i) => (
    i === ruleIndex
      ? { ...r, conditions: r.conditions.map((c, ci) => (ci === condIndex ? { ...c, [field]: value as any } : c)) }
      : r
  )));

  const rulesValid = rules.every((r) =>
    r.flag_variant_id && variants.some((v) => v.id === r.flag_variant_id) &&
    r.conditions.length > 0 &&
    r.conditions.every((c) => c.attribute.trim().length > 0)
  );

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
            <TextField select label="Kind" value={kind} onChange={(e) => { const v = e.target.value as FeatureKind; setKind(v); if (v === 'boolean') { setDefaultVariant(defaultVariant === 'on' ? 'on' : 'off'); } else { const names = variants.map(vv => vv.name.trim()).filter(n => n.length > 0); setDefaultVariant(names[0] || ''); } }} fullWidth>
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
              <TextField
                select
                label="Default Variant"
                value={defaultVariant}
                onChange={(e) => setDefaultVariant(e.target.value)}
                fullWidth
                required
                disabled={variantNames.length === 0}
                error={!hasAtLeastTwoVariants || (variantNames.length > 0 && !variantNames.includes(defaultVariant))}
                helperText={!hasAtLeastTwoVariants ? 'Add at least two variants first' : (!variantNames.includes(defaultVariant) ? 'Select one of the defined variants' : '')}
              >
                {variantNames.map((n) => (
                  <MenuItem key={n} value={n}>{n}</MenuItem>
                ))}
              </TextField>
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
                  <TextField label="Rollout %" type="number" value={v.rollout_percent} onChange={(e) => handleVariantChange(i, 'rollout_percent', e.target.value)} fullWidth inputProps={{ min: 1, max: 100, step: 1 }} />
                  <IconButton aria-label="delete" onClick={() => handleRemoveVariant(i)}>
                    <DeleteIcon />
                  </IconButton>
                </Box>
              ))}
              <Button size="small" startIcon={<AddIcon />} onClick={handleAddVariant}>Add Variant</Button>

              <Box sx={{ mt: 2 }}>
                <Typography variant="body2" sx={{ mb: 0.5 }}>
                  Total rollout: <strong>{Number.isFinite(totalRollout) ? totalRollout : 0}%</strong>
                </Typography>
                {!rolloutSumValid && (
                  <Typography variant="body2" color="error" sx={{ mb: 1 }}>
                    Sum of rollout percentages must equal 100.
                  </Typography>
                )}
                {!variantsValuesValid && (
                  <Typography variant="body2" color="error" sx={{ mb: 1 }}>
                    Each variant must have a name and rollout between 1 and 100.
                  </Typography>
                )}
                <Typography variant="subtitle1" sx={{ mt: 2 }}>Rules</Typography>
                <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
                  Define optional rules that map conditions to a target variant.
                </Typography>
                {rules.map((r, ri) => (
                  <Box key={r.id} sx={{ border: '1px solid', borderColor: 'divider', borderRadius: 1, p: 1.5, mb: 1.5 }}>
                    <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 1 }}>
                      <Typography variant="subtitle2">Rule #{ri + 1}</Typography>
                      <IconButton aria-label="delete" onClick={() => removeRule(ri)} size="small"><DeleteIcon /></IconButton>
                    </Box>
                    <Box sx={{ display: 'grid', gridTemplateColumns: { xs: '1fr', md: '1fr 1fr' }, gap: 1, mb: 1 }}>
                      <TextField
                        select
                        label="Target Variant"
                        value={r.flag_variant_id}
                        onChange={(e) => updateRule(ri, { flag_variant_id: e.target.value })}
                        required
                        error={!!r && (!r.flag_variant_id || !variants.some(v => v.id === r.flag_variant_id))}
                        helperText="Select which variant this rule serves"
                      >
                        {variants.map(v => (
                          <MenuItem key={v.id} value={v.id}>{v.name || '(unnamed)'}</MenuItem>
                        ))}
                      </TextField>
                      <TextField
                        label="Rule Priority"
                        type="number"
                        value={r.priority}
                        onChange={(e) => updateRule(ri, { priority: e.target.value === '' ? '' : Number(e.target.value) })}
                        helperText="Lower numbers run first"
                      />
                    </Box>
                    <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
                      Conditions
                    </Typography>
                    {r.conditions.map((c, ci) => (
                      <Box key={ci} sx={{ display: 'grid', gridTemplateColumns: { xs: '1fr', md: '1.5fr 1fr 1.5fr auto' }, gap: 1, alignItems: 'center', mb: 1 }}>
                        <TextField label="Attribute" value={c.attribute} onChange={(e) => changeRuleConditionInline(ri, ci, 'attribute', e.target.value)} />
                        <TextField select label="Operator" value={c.operator} onChange={(e) => changeRuleConditionInline(ri, ci, 'operator', e.target.value)}>
                          {ruleOperatorOptions.map(op => (
                            <MenuItem key={op} value={op}>{op}</MenuItem>
                          ))}
                        </TextField>
                        <TextField label="Value" value={c.value} onChange={(e) => changeRuleConditionInline(ri, ci, 'value', e.target.value)} helperText="JSON or text" />
                        <IconButton aria-label="delete" onClick={() => removeRuleConditionInline(ri, ci)}>
                          <DeleteIcon />
                        </IconButton>
                      </Box>
                    ))}
                    <Button size="small" startIcon={<AddIcon />} onClick={() => addRuleConditionInline(ri)}>Add Condition</Button>
                  </Box>
                ))}
                <Button size="small" variant="outlined" startIcon={<AddIcon />} onClick={addRule} disabled={variants.length === 0}>
                  Add Rule
                </Button>
                {rules.length > 0 && !rulesValid && (
                  <Typography variant="body2" color="error" sx={{ mt: 1 }}>
                    Please fix rules: each must have a target variant and at least one condition with attribute.
                  </Typography>
                )}
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
          <Button variant="contained" onClick={() => { setSubmitting(true); createFeatureMutation.mutate(); }} disabled={submitting || !keyValue.trim() || !name.trim() || (kind === 'multivariant' && (!hasAtLeastTwoVariants || !variantsValid || !variantNames.includes(defaultVariant) || (rules.length > 0 && !rulesValid)))}>
            {submitting ? 'Creating...' : 'Create'}
          </Button>
        </DialogActions>
      </Dialog>

    </AuthenticatedLayout>
  );
};

export default ProjectPage;
