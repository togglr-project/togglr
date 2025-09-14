import React, { useEffect, useMemo, useState } from 'react';
import {
  Box,
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Alert,
  CircularProgress,
  Typography,
  TextField,
  Grid,
  Switch,
  FormControlLabel,
  Select,
  MenuItem,
  InputLabel,
  FormControl,
  IconButton,
  Divider,
  Tooltip,
  Chip,
  Radio,
  RadioGroup,
} from '@mui/material';
import { Add, Delete } from '@mui/icons-material';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import apiClient from '../../api/apiClient';
import type { CreateFeatureRequest, FeatureDetailsResponse, RuleCondition } from '../../generated/api/client';
import { FeatureKind as FeatureKindEnum, RuleOperator as RuleOperatorEnum } from '../../generated/api/client';

export interface EditFeatureDialogProps {
  open: boolean;
  onClose: () => void;
  featureDetails: FeatureDetailsResponse | null;
}

const uuid = () => (typeof crypto !== 'undefined' && 'randomUUID' in crypto ? crypto.randomUUID() : Math.random().toString(36).slice(2));

const EditFeatureDialog: React.FC<EditFeatureDialogProps> = ({ open, onClose, featureDetails }) => {
  const queryClient = useQueryClient();

  type VariantForm = { id: string; name: string; rollout_percent: number };
  type RuleForm = { id: string; priority: number; flag_variant_id: string; conditions: RuleCondition[] };

  const [keyVal, setKeyVal] = useState('');
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [kind, setKind] = useState<string>(FeatureKindEnum.Boolean);
  const [enabled, setEnabled] = useState<boolean>(false);
  const [variants, setVariants] = useState<VariantForm[]>([]);
  const [defaultVariantId, setDefaultVariantId] = useState<string>('');
  const [rules, setRules] = useState<RuleForm[]>([]);

  const [error, setError] = useState<string | null>(null);

  // Initialize form from featureDetails
  useEffect(() => {
    if (!open || !featureDetails) return;
    const f = featureDetails.feature;
    setKeyVal(f.key);
    setName(f.name);
    setDescription(f.description || '');
    setKind(f.kind);
    setEnabled(f.enabled);
    const vars = (featureDetails.variants || []).map(v => ({ id: v.id, name: v.name, rollout_percent: v.rollout_percent }));
    setVariants(vars);
    const defId = vars.find(v => v.id === f.default_variant || v.name === f.default_variant)?.id || vars[0]?.id || '';
    setDefaultVariantId(defId);
    const rls = (featureDetails.rules || []).map(r => ({ id: r.id, priority: r.priority, flag_variant_id: r.flag_variant_id, conditions: r.conditions }));
    setRules(rls);
    setError(null);
  }, [open, featureDetails]);

  const featureId = featureDetails?.feature.id;
  const projectId = featureDetails?.feature.project_id;

  const updateMutation = useMutation({
    mutationFn: async (body: CreateFeatureRequest) => {
      if (!featureId) throw new Error('No feature id');
      const res = await apiClient.updateFeature(featureId, body);
      return res.data;
    },
    onSuccess: () => {
      if (featureId) {
        queryClient.invalidateQueries({ queryKey: ['feature-details', featureId] });
      }
      if (projectId) {
        queryClient.invalidateQueries({ queryKey: ['project-features', projectId] });
      }
      onClose();
    },
  });

  const variantNameById = useMemo(() => Object.fromEntries(variants.map(v => [v.id, v.name || v.id])), [variants]);

  const addVariant = () => {
    const id = uuid();
    setVariants(prev => [...prev, { id, name: `variant_${prev.length + 1}`, rollout_percent: 0 }]);
    if (!defaultVariantId) setDefaultVariantId(id);
  };
  const removeVariant = (id: string) => {
    // prevent removing default or used by rules
    if (id === defaultVariantId) return;
    if (rules.some(r => r.flag_variant_id === id)) return;
    setVariants(prev => prev.filter(v => v.id !== id));
  };

  const addRule = () => {
    const id = uuid();
    const priority = (rules.length ? Math.max(...rules.map(r => r.priority)) + 1 : 1);
    const firstVariant = variants[0]?.id || '';
    setRules(prev => [...prev, { id, priority, flag_variant_id: firstVariant, conditions: [] }]);
  };
  const removeRule = (id: string) => setRules(prev => prev.filter(r => r.id !== id));

  const addCondition = (ruleId: string) => {
    setRules(prev => prev.map(r => r.id === ruleId ? { ...r, conditions: [...r.conditions, { attribute: '', operator: RuleOperatorEnum.Eq, value: '' }] } : r));
  };
  const removeCondition = (ruleId: string, idx: number) => {
    setRules(prev => prev.map(r => r.id === ruleId ? { ...r, conditions: r.conditions.filter((_, i) => i !== idx) } : r));
  };

  const parseValueSmart = (input: string): any => {
    const t = input.trim();
    if (t === '') return '';
    if (t === 'true') return true;
    if (t === 'false') return false;
    if (!isNaN(Number(t))) return Number(t);
    try { return JSON.parse(t); } catch { return input; }
  };

  const validate = (): string | null => {
    if (!keyVal.trim()) return 'Key is required';
    if (!name.trim()) return 'Name is required';
    if (!kind) return 'Kind is required';
    if (!defaultVariantId) return 'Default variant is required';
    if (!variants.find(v => v.id === defaultVariantId)) return 'Default variant must be one of the variants';
    for (const v of variants) {
      if (!v.name.trim()) return 'Variant name cannot be empty';
      if (v.rollout_percent < 0 || v.rollout_percent > 100) return 'Variant rollout percent must be between 0 and 100';
    }
    for (const r of rules) {
      if (!variants.find(v => v.id === r.flag_variant_id)) return 'Each rule must target an existing variant';
      for (const c of r.conditions) {
        if (!c.attribute.trim()) return 'Condition attribute cannot be empty';
        if (!c.operator) return 'Condition operator is required';
      }
    }
    return null;
  };

  const onSave = () => {
    setError(null);
    const err = validate();
    if (err) { setError(err); return; }

    const payload: CreateFeatureRequest = {
      key: keyVal,
      name,
      description: description || undefined,
      kind: kind as any,
      default_variant: (variants.find(v => v.id === defaultVariantId)?.name ?? defaultVariantId),
      enabled,
      variants: variants.map(v => ({ id: v.id, name: v.name, rollout_percent: Number(v.rollout_percent) })),
      rules: rules.map(r => ({
        id: r.id,
        priority: r.priority,
        flag_variant_id: r.flag_variant_id,
        conditions: r.conditions.map(c => ({ attribute: c.attribute, operator: c.operator, value: c.value })),
      })),
    };

    updateMutation.mutate(payload);
  };

  const disabled = !featureDetails || updateMutation.isPending;

  const ruleOperatorOptions = Object.values(RuleOperatorEnum);
  const featureKindOptions = Object.values(FeatureKindEnum);

  const canDeleteVariant = (id: string) => id !== defaultVariantId && !rules.some(r => r.flag_variant_id === id);
  const onId = useMemo(() => variants.find(v => (v.name || v.id).toLowerCase() === 'on')?.id, [variants]);
  const offId = useMemo(() => variants.find(v => (v.name || v.id).toLowerCase() === 'off')?.id, [variants]);

  return (
    <Dialog open={open} onClose={disabled ? undefined : onClose} fullWidth maxWidth="md">
      <DialogTitle className="gradient-text-purple">Edit Feature</DialogTitle>
      <DialogContent>
        {!featureDetails ? (
          <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
            <CircularProgress />
          </Box>
        ) : (
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, mt: 1 }}>
            {error && <Alert severity="error">{error}</Alert>}
            {updateMutation.isError && (
              <Alert severity="error">{(updateMutation.error as any)?.message || 'Failed to update feature'}</Alert>
            )}

            {/* Basic info */}
            <Typography variant="subtitle1">Basic</Typography>
            <Grid container spacing={2}>
              <Grid item xs={12} md={6}>
                <TextField label="Key" value={keyVal} onChange={(e) => setKeyVal(e.target.value)} fullWidth disabled />
              </Grid>
              <Grid item xs={12} md={6}>
                <TextField label="Name" value={name} onChange={(e) => setName(e.target.value)} fullWidth />
              </Grid>
              <Grid item xs={12}>
                <TextField label="Description" value={description} onChange={(e) => setDescription(e.target.value)} fullWidth multiline minRows={2} />
              </Grid>
              <Grid item xs={12} md={6}>
                <FormControl fullWidth>
                  <InputLabel id="kind-label">Kind</InputLabel>
                  <Select labelId="kind-label" label="Kind" value={kind} onChange={(e) => setKind(e.target.value)} disabled>
                    {featureKindOptions.map(k => (
                      <MenuItem key={k} value={k}>{k}</MenuItem>
                    ))}
                  </Select>
                </FormControl>
              </Grid>
              <Grid item xs={12} md={6} sx={{ display: 'flex', alignItems: 'center' }}>
                <FormControlLabel control={<Switch checked={enabled} onChange={(e) => setEnabled(e.target.checked)} />} label="Enabled" />
              </Grid>
            </Grid>

            <Divider />

            {/* Variants / Default */}
            {kind === FeatureKindEnum.Multivariant ? (
              <>
                <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                  <Typography variant="subtitle1">Variants</Typography>
                  <Button startIcon={<Add />} onClick={addVariant} size="small">Add variant</Button>
                </Box>
                <Grid container spacing={1}>
                  {variants.map((v, idx) => (
                    <React.Fragment key={v.id}>
                      <Grid item xs={12} md={4}>
                        <TextField label={`Variant name #${idx + 1}`} value={v.name} onChange={(e) => setVariants(prev => prev.map(x => x.id === v.id ? { ...x, name: e.target.value } : x))} fullWidth />
                      </Grid>
                      <Grid item xs={12} md={3}>
                        <TextField type="number" label="Rollout %" inputProps={{ min: 0, max: 100 }} value={v.rollout_percent}
                          onChange={(e) => setVariants(prev => prev.map(x => x.id === v.id ? { ...x, rollout_percent: Number(e.target.value) } : x))} fullWidth />
                      </Grid>
                      <Grid item xs={12} md={4}>
                        <TextField label="ID" value={v.id} fullWidth disabled />
                      </Grid>
                      <Grid item xs={12} md={1} sx={{ display: 'flex', alignItems: 'center' }}>
                        <Tooltip title={canDeleteVariant(v.id) ? 'Delete variant' : 'Cannot delete: used in default or rules'}>
                          <span>
                            <IconButton onClick={() => removeVariant(v.id)} disabled={!canDeleteVariant(v.id)} aria-label="delete variant">
                              <Delete />
                            </IconButton>
                          </span>
                        </Tooltip>
                      </Grid>
                    </React.Fragment>
                  ))}
                </Grid>
                <Box sx={{ display: 'flex', gap: 2, alignItems: 'center' }}>
                  <FormControl sx={{ minWidth: 240 }}>
                    <InputLabel id="default-var-label">Default variant</InputLabel>
                    <Select labelId="default-var-label" label="Default variant" value={defaultVariantId} onChange={(e) => setDefaultVariantId(e.target.value)}>
                      {variants.map(v => (
                        <MenuItem key={v.id} value={v.id}>{v.name || v.id}</MenuItem>
                      ))}
                    </Select>
                  </FormControl>
                  <Typography variant="body2" color="text.secondary">Variants total: {variants.length}</Typography>
                </Box>
              </>
            ) : (
              <Box sx={{ display: 'flex', gap: 2, alignItems: 'center' }}>
                <Typography variant="subtitle2">Default</Typography>
                <Chip label={featureDetails?.feature.default_variant ?? ''} />
              </Box>
            )}

            <Divider />

            {/* Rules */}
            <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
              <Typography variant="subtitle1">Rules</Typography>
              <Button startIcon={<Add />} onClick={addRule} size="small">Add rule</Button>
            </Box>
            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
              {rules.sort((a,b) => a.priority - b.priority).map((r, rIndex) => (
                <Box key={r.id} sx={{ border: '1px solid', borderColor: 'divider', borderRadius: 1, p: 1.5 }}>
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, flexWrap: 'wrap' }}>
                    <TextField type="number" label="Priority" value={r.priority} sx={{ width: 140 }}
                      onChange={(e) => setRules(prev => prev.map(x => x.id === r.id ? { ...x, priority: Number(e.target.value) } : x))} />
                    <FormControl sx={{ minWidth: 240 }}>
                      <InputLabel id={`rule-var-${r.id}`}>Target variant</InputLabel>
                      <Select labelId={`rule-var-${r.id}`} label="Target variant" value={r.flag_variant_id}
                        onChange={(e) => setRules(prev => prev.map(x => x.id === r.id ? { ...x, flag_variant_id: e.target.value } : x))}>
                        {variants.map(v => (
                          <MenuItem key={v.id} value={v.id}>{variantNameById[v.id]}</MenuItem>
                        ))}
                      </Select>
                    </FormControl>
                    <IconButton onClick={() => removeRule(r.id)} aria-label="delete rule"><Delete /></IconButton>
                  </Box>

                  <Box sx={{ mt: 1 }}>
                    <Typography variant="body2" color="text.secondary">Conditions</Typography>
                    <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1 }}>
                      {r.conditions.map((c, cIdx) => (
                        <Grid key={`${r.id}-${cIdx}`} container spacing={1} alignItems="center">
                          <Grid item xs={12} md={4}>
                            <TextField label="Attribute" value={c.attribute}
                              onChange={(e) => setRules(prev => prev.map(x => x.id === r.id ? { ...x, conditions: x.conditions.map((cc, i) => i === cIdx ? { ...cc, attribute: e.target.value } : cc) } : x))} fullWidth />
                          </Grid>
                          <Grid item xs={12} md={3}>
                            <FormControl fullWidth>
                              <InputLabel id={`op-${r.id}-${cIdx}`}>Operator</InputLabel>
                              <Select labelId={`op-${r.id}-${cIdx}`} label="Operator" value={c.operator}
                                onChange={(e) => setRules(prev => prev.map(x => x.id === r.id ? { ...x, conditions: x.conditions.map((cc, i) => i === cIdx ? { ...cc, operator: e.target.value as any } : cc) } : x))}>
                                {ruleOperatorOptions.map(op => (
                                  <MenuItem key={op} value={op}>{op}</MenuItem>
                                ))}
                              </Select>
                            </FormControl>
                          </Grid>
                          <Grid item xs={12} md={4}>
                            <TextField label="Value" value={String(c.value ?? '')}
                              onChange={(e) => setRules(prev => prev.map(x => x.id === r.id ? { ...x, conditions: x.conditions.map((cc, i) => i === cIdx ? { ...cc, value: parseValueSmart(e.target.value) } : cc) } : x))} fullWidth />
                          </Grid>
                          <Grid item xs={12} md={1} sx={{ display: 'flex', alignItems: 'center' }}>
                            <IconButton onClick={() => removeCondition(r.id, cIdx)} aria-label="delete condition"><Delete /></IconButton>
                          </Grid>
                        </Grid>
                      ))}
                      <Button startIcon={<Add />} onClick={() => addCondition(r.id)} size="small" sx={{ alignSelf: 'flex-start' }}>Add condition</Button>
                    </Box>
                  </Box>
                </Box>
              ))}
              {rules.length === 0 && <Typography variant="body2" color="text.secondary">No rules</Typography>}
            </Box>
          </Box>
        )}
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose} disabled={disabled}>Cancel</Button>
        <Button onClick={onSave} variant="contained" disabled={disabled}>{updateMutation.isPending ? 'Saving…' : 'Save'}</Button>
      </DialogActions>
    </Dialog>
  );
};

export default EditFeatureDialog;
