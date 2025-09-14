import React, { useEffect, useMemo, useState } from 'react';
import {
  Box,
  Button,
  Chip,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  FormControlLabel,
  Grid,
  IconButton,
  MenuItem,
  Radio,
  RadioGroup,
  Switch,
  TextField,
  Typography,
} from '@mui/material';
import { Add as AddIcon, Delete as DeleteIcon } from '@mui/icons-material';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import apiClient from '../../api/apiClient';
import type { Feature, FeatureKind } from '../../generated/api/client';

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

const kindOptions: FeatureKind[] = ['boolean', 'multivariant'];

type OperatorOption = 'eq' | 'neq' | 'in' | 'not_in' | 'gt' | 'gte' | 'lt' | 'lte' | 'regex' | 'percentage';
interface RuleConditionItem { attribute: string; operator: OperatorOption; value: string }
interface RuleFormItem { id: string; flag_variant_id: string; priority: number | ''; conditions: RuleConditionItem[] }
interface VariantFormItem { id: string; name: string; rollout_percent: number }

export interface CreateFeatureDialogProps {
  open: boolean;
  onClose: () => void;
  projectId: string;
}

const CreateFeatureDialog: React.FC<CreateFeatureDialogProps> = ({ open, onClose, projectId }) => {
  const queryClient = useQueryClient();

  // Form state
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

  // Derived validation helpers
  const totalRollout = useMemo(() => variants.reduce((sum, v) => sum + (Number.isFinite(Number(v.rollout_percent)) ? Number(v.rollout_percent) : 0), 0), [variants]);
  const rolloutSumValid = Math.round(totalRollout) === 100;
  const variantsValuesValid = variants.every((v) => v.name.trim().length > 0 && Number.isFinite(Number(v.rollout_percent)) && v.rollout_percent > 0 && v.rollout_percent <= 100);
  const variantsValid = variantsValuesValid && rolloutSumValid;
  const variantNames = variants.map(v => v.name.trim()).filter(n => n.length > 0);
  const hasAtLeastTwoVariants = variantNames.length >= 2;

  const priorityCounts = rules.reduce((acc, r) => {
    const p = r.priority;
    if (typeof p === 'number') {
      acc[p] = (acc[p] || 0) + 1;
    }
    return acc;
  }, {} as Record<number, number>);
  const hasDuplicatePriorities = Object.values(priorityCounts).some((c) => c > 1);

  const rulesValid = rules.every((r) =>
    r.flag_variant_id && variants.some((v) => v.id === r.flag_variant_id) &&
    r.conditions.length > 0 &&
    r.conditions.every((c) => c.attribute.trim().length > 0) &&
    typeof r.priority === 'number' && Number.isInteger(r.priority) && r.priority >= 0 && r.priority <= 255
  ) && !hasDuplicatePriorities;

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

      if (!keyValue.trim() || !name.trim()) {
        throw new Error('Key and Name are required');
      }
      if (kind === 'multivariant') {
        if (!hasAtLeastTwoVariants) throw new Error('At least two variants are required for multivariant features');
        if (!variantsValid) throw new Error('Variants must have names, rollout between 1 and 100, and total rollout must equal 100');
        const names = variants.map((v) => v.name.trim()).filter((n) => n.length > 0);
        if (!names.includes(defaultVariant.trim())) throw new Error('Default Variant must be one of the variants');
      }

      const dv = kind === 'boolean' ? (defaultVariant === 'on' ? 'on' : 'off') : defaultVariant.trim();

      let inlineVariants: { id: string; name: string; rollout_percent: number }[] | undefined;
      let inlineRules: { id: string; conditions: any[]; flag_variant_id: string; priority?: number }[] | undefined;

      if (kind === 'multivariant') {
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
      onClose();
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

  const ruleOperatorOptions: OperatorOption[] = ['eq','neq','in','not_in','gt','gte','lt','lte','regex','percentage'];

  const addRule = () => setRules((prev) => {
    const nums = prev.map((r) => (typeof r.priority === 'number' ? r.priority : 0));
    const next = nums.length ? Math.min(255, Math.max(...nums) + 1) : 0;
    return [
      ...prev,
      { id: genId(), flag_variant_id: variants[0]?.id || '', priority: next, conditions: [{ attribute: '', operator: 'eq', value: '' }] }
    ];
  });
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

  const defaultVariantErrorText = useMemo(() => {
    if (kind !== 'multivariant') return '';
    if (!hasAtLeastTwoVariants) return 'Add at least two variants first';
    if (variantNames.length > 0 && !variantNames.includes(defaultVariant)) return 'Select one of the defined variants';
    return '';
  }, [kind, hasAtLeastTwoVariants, variantNames, defaultVariant]);

  const canCreate = useMemo(() => {
    if (!keyValue.trim() || !name.trim()) return false;
    if (kind === 'multivariant') {
      if (!hasAtLeastTwoVariants) return false;
      if (!variantsValid) return false;
      if (!variantNames.includes(defaultVariant)) return false;
      if (!rulesValid && rules.length > 0) return false;
    }
    return !submitting;
  }, [keyValue, name, kind, hasAtLeastTwoVariants, variantsValid, variantNames, defaultVariant, rulesValid, rules.length, submitting]);

  return (
    <Dialog open={open} onClose={onClose} fullWidth maxWidth="md">
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
              error={!!defaultVariantErrorText}
              helperText={defaultVariantErrorText}
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
            <Typography variant="subtitle1">Variants</Typography>
            <Box sx={{ mt: 1, display: 'flex', flexDirection: 'column', gap: 1 }}>
              {variants.map((v, index) => (
                <Box key={v.id} sx={{ display: 'grid', gridTemplateColumns: { xs: '1fr', md: '1fr 120px 40px' }, gap: 1, alignItems: 'center' }}>
                  <TextField
                    label="Name"
                    value={v.name}
                    onChange={(e) => handleVariantChange(index, 'name', e.target.value)}
                    required
                    fullWidth
                  />
                  <TextField
                    label="Rollout %"
                    type="number"
                    inputProps={{ min: 1, max: 100, step: 1 }}
                    value={v.rollout_percent}
                    onChange={(e) => handleVariantChange(index, 'rollout_percent', e.target.value)}
                    required
                  />
                  <IconButton aria-label="delete-variant" onClick={() => handleRemoveVariant(index)} disabled={variants.length <= 1}>
                    <DeleteIcon />
                  </IconButton>
                </Box>
              ))}
              <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <Button startIcon={<AddIcon />} onClick={handleAddVariant}>Add Variant</Button>
                <Typography variant="body2" color={rolloutSumValid ? 'text.secondary' : 'error'}>
                  Total rollout: {totalRollout}% (must equal 100)
                </Typography>
              </Box>
              {!variantsValuesValid && (
                <Typography variant="body2" color="error">Each variant must have a name and rollout in 1..100</Typography>
              )}
              {!hasAtLeastTwoVariants && (
                <Typography variant="body2" color="error">At least two variants are required</Typography>
              )}
            </Box>

            <Box sx={{ mt: 3 }}>
              <Typography variant="subtitle1">Rules</Typography>
              <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
                Define rules to route users to specific variants. Priorities must be unique (0-255).
              </Typography>
              {rules.map((r, ri) => {
                const priorityError = !(typeof r.priority === 'number' && Number.isInteger(r.priority) && r.priority >= 0 && r.priority <= 255);
                const isDup = typeof r.priority === 'number' && (priorityCounts[r.priority] || 0) > 1;
                return (
                  <Box key={r.id} sx={{ border: '1px solid', borderColor: 'divider', borderRadius: 1, p: 1.5, mb: 1 }}>
                    <Box sx={{ display: 'grid', gridTemplateColumns: { xs: '1fr', md: '1fr 200px 120px 40px' }, gap: 1, alignItems: 'center' }}>
                      <TextField
                        label="Priority"
                        type="number"
                        inputProps={{ min: 0, max: 255, step: 1 }}
                        value={r.priority}
                        onChange={(e) => updateRule(ri, { priority: e.target.value === '' ? '' : Number(e.target.value) })}
                        error={priorityError || isDup}
                        helperText={priorityError ? '0..255 integer' : (isDup ? 'Priority must be unique' : '')}
                      />
                      <TextField
                        select
                        label="Target Variant"
                        value={r.flag_variant_id}
                        onChange={(e) => updateRule(ri, { flag_variant_id: e.target.value })}
                        required
                      >
                        {variants.map(v => (
                          <MenuItem key={v.id} value={v.id}>{v.name || v.id}</MenuItem>
                        ))}
                      </TextField>
                      <Box />
                      <IconButton aria-label="delete-rule" onClick={() => removeRule(ri)}>
                        <DeleteIcon />
                      </IconButton>
                    </Box>

                    <Box sx={{ mt: 1 }}>
                      <Typography variant="body2" color="text.secondary">Conditions</Typography>
                      {r.conditions.map((c, ci) => (
                        <Grid container spacing={1} alignItems="center" key={ci} sx={{ mt: 0.5 }}>
                          <Grid item xs={12} md={4}>
                            <TextField
                              label="Attribute"
                              value={c.attribute}
                              onChange={(e) => changeRuleConditionInline(ri, ci, 'attribute', e.target.value)}
                              required
                              fullWidth
                            />
                          </Grid>
                          <Grid item xs={12} md={3}>
                            <TextField
                              select
                              label="Operator"
                              value={c.operator}
                              onChange={(e) => changeRuleConditionInline(ri, ci, 'operator', e.target.value)}
                              fullWidth
                            >
                              {ruleOperatorOptions.map(op => (
                                <MenuItem key={op} value={op}>{op}</MenuItem>
                              ))}
                            </TextField>
                          </Grid>
                          <Grid item xs={12} md={4}>
                            <TextField
                              label="Value (JSON or plain)"
                              value={c.value}
                              onChange={(e) => changeRuleConditionInline(ri, ci, 'value', e.target.value)}
                              fullWidth
                            />
                          </Grid>
                          <Grid item xs={12} md={1}>
                            <IconButton aria-label="delete-condition" onClick={() => removeRuleConditionInline(ri, ci)}>
                              <DeleteIcon />
                            </IconButton>
                          </Grid>
                        </Grid>
                      ))}
                      <Box sx={{ mt: 1 }}>
                        <Button size="small" startIcon={<AddIcon />} onClick={() => addRuleConditionInline(ri)}>Add Condition</Button>
                      </Box>
                    </Box>
                  </Box>
                );
              })}
              <Button variant="outlined" startIcon={<AddIcon />} onClick={addRule}>Add Rule</Button>
              {hasDuplicatePriorities && (
                <Box sx={{ mt: 1 }}>
                  <Chip size="small" color="warning" label="Duplicate priorities detected" />
                </Box>
              )}
            </Box>
          </Box>
        )}

        {formError && (
          <Typography color="error" sx={{ mt: 2 }}>{formError}</Typography>
        )}
      </DialogContent>
      <DialogActions>
        <Button onClick={() => { resetForm(); onClose(); }}>Cancel</Button>
        <Button
          onClick={() => { setSubmitting(true); createFeatureMutation.mutate(); }}
          variant="contained"
          disabled={!canCreate}
        >
          Create
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default CreateFeatureDialog;
