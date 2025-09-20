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
  Switch,
  TextField,
  Typography,
} from '@mui/material';
import Autocomplete from '@mui/material/Autocomplete';
import ConditionExpressionBuilder from '../conditions/ConditionExpressionBuilder';
import { Add as AddIcon, Delete as DeleteIcon } from '@mui/icons-material';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import apiClient from '../../api/apiClient';
import type { Feature, FeatureKind, RuleAction as RuleActionType, RuleConditionExpression, Segment } from '../../generated/api/client';
import { RuleAction as RuleActionEnum } from '../../generated/api/client';

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

const kindOptions: FeatureKind[] = ['simple', 'multivariant'];
const rolloutKeyOptions = ['user.id', 'user.email'];

type OperatorOption = 'eq' | 'neq' | 'in' | 'not_in' | 'gt' | 'gte' | 'lt' | 'lte' | 'regex' | 'percentage';
interface RuleConditionItem { attribute: string; operator: OperatorOption; value: string }
interface RuleFormItem { id: string; action: RuleActionType; flag_variant_id?: string; priority: number | ''; expression: RuleConditionExpression; segment_id?: string; is_customized: boolean; baseExpressionJson?: string }
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
  const [rolloutKey, setRolloutKey] = useState('');
  const [kind, setKind] = useState<FeatureKind>('simple');
  const [defaultVariant, setDefaultVariant] = useState('');
  const [enabled, setEnabled] = useState(true);
  const [variants, setVariants] = useState<VariantFormItem[]>([{ id: genId(), name: 'control', rollout_percent: 100 }]);
  const [rules, setRules] = useState<RuleFormItem[]>([]);

  // Load project segments for reuse in rules
  const { data: segments } = useQuery<Segment[]>({
    queryKey: ['project-segments', projectId],
    queryFn: async () => {
      const res = await apiClient.listProjectSegments(projectId);
      const resp = res.data as any;
      return Array.isArray(resp?.items) ? (resp.items as Segment[]) : (resp as Segment[]);
    },
    enabled: !!projectId,
  });
  const [formError, setFormError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  // Derived validation helpers
  // Feature key: allowed characters a-z, A-Z, 0-9, hyphen (-), underscore (_), colon (:), @, !, #, $, dot (.)
  const keyRegex = /^[A-Za-z0-9_:@!#$.-]+$/; // hyphen and dot are safe here; no spaces allowed
  const keyValid = useMemo(() => {
    const v = keyValue.trim();
    return v.length > 0 && keyRegex.test(v);
  }, [keyValue]);

  const totalRollout = useMemo(() => variants.reduce((sum, v) => sum + (Number.isFinite(Number(v.rollout_percent)) ? Number(v.rollout_percent) : 0), 0), [variants]);
  const rolloutSumValid = Math.round(totalRollout) === 100;
  const variantsValuesValid = variants.every((v) => v.name.trim().length > 0 && Number.isFinite(Number(v.rollout_percent)) && v.rollout_percent > 0 && v.rollout_percent <= 100);
  const variantsValid = variantsValuesValid && rolloutSumValid;
  const variantNames = variants.map(v => v.name.trim()).filter(n => n.length > 0);
  const hasAtLeastTwoVariants = variantNames.length >= 2;

  const hasDuplicatePriorities = [RuleActionEnum.Assign, RuleActionEnum.Include, RuleActionEnum.Exclude].some((action) => {
    const counts = rules.filter(r => r.action === action).reduce((acc, r) => {
      const p = r.priority;
      if (typeof p === 'number') {
        acc[p] = (acc[p] || 0) + 1;
      }
      return acc;
    }, {} as Record<number, number>);
    return Object.values(counts).some((c) => c > 1);
  });

  const hasValidLeaf = (e?: RuleConditionExpression): boolean => {
    if (!e) return false;
    if ((e as any).condition) {
      const c = (e as any).condition as { attribute?: string };
      return Boolean(c.attribute && c.attribute.trim().length > 0);
    }
    if ((e as any).group) {
      const g = (e as any).group as { children?: RuleConditionExpression[] };
      return Array.isArray(g.children) && g.children.some(ch => hasValidLeaf(ch));
    }
    return false;
  };

  const rulesValid = rules.every((r) =>
    (r.action === RuleActionEnum.Assign ? (r.flag_variant_id && variants.some((v) => v.id === r.flag_variant_id)) : true) &&
    hasValidLeaf(r.expression) &&
    typeof r.priority === 'number' && Number.isInteger(r.priority) && r.priority >= 0 && r.priority <= 255
  ) && !hasDuplicatePriorities;


  // Ensure assign rules reference existing variants when variants change
  useEffect(() => {
    setRules((prev) => prev.map((r) => {
      if (r.action !== RuleActionEnum.Assign) return r;
      if (r.flag_variant_id && variants.some((v) => v.id === r.flag_variant_id)) return r;
      const next = variants[0]?.id;
      return { ...r, flag_variant_id: next };
    }));
  }, [variants]);

  // When switching to simple, drop any assign rules (not allowed for simple)
  useEffect(() => {
    if (kind === 'simple') {
      setRules((prev) => prev.filter((r) => r.action !== RuleActionEnum.Assign));
    }
  }, [kind]);

  const resetForm = () => {
    setKeyValue('');
    setName('');
    setDescription('');
    setRolloutKey('');
    setKind('simple');
    setDefaultVariant('');
    setEnabled(true);
    setVariants([{ id: genId(), name: 'control', rollout_percent: 100 }]);
    setRules([]);
    setFormError(null);
  };

  const createFeatureMutation = useMutation({
    mutationFn: async () => {
      setFormError(null);

      if (!keyValid || !name.trim()) {
        if (!keyValid) {
          throw new Error('Key contains invalid characters. Allowed: a-z, A-Z, 0-9, -, _, :, @, !, #, $, .');
        }
        throw new Error('Key and Name are required');
      }
      if (kind === 'multivariant') {
        if (!rolloutKey.trim()) throw new Error('Rollout Key is required for multivariant features');
        if (!hasAtLeastTwoVariants) throw new Error('At least two variants are required for multivariant features');
        if (!variantsValid) throw new Error('Variants must have names, rollout between 1 and 100, and total rollout must equal 100');
      }

      const dv = defaultVariant.trim();

      let inlineVariants: { id: string; name: string; rollout_percent: number }[] | undefined;
      let inlineRules: { id: string; conditions: RuleConditionExpression; action: RuleActionType; flag_variant_id?: string; priority?: number }[] | undefined;

      if (kind === 'multivariant') {
        inlineVariants = variants.map((v) => ({ id: v.id, name: v.name.trim(), rollout_percent: Number(v.rollout_percent) || 0 }));
      }

      if (rules.length > 0) {
        if (!rulesValid) {
          throw new Error('Please fix rules');
        }
        inlineRules = rules.map((r) => ({
          id: r.id,
          action: r.action,
          flag_variant_id: r.action === RuleActionEnum.Assign ? r.flag_variant_id : undefined,
          priority: r.priority === '' ? 0 : Number(r.priority),
          conditions: r.expression,
          segment_id: r.segment_id || undefined,
          is_customized: r.segment_id ? Boolean(r.is_customized) : true,
        }));
      }

      const featureRes = await apiClient.createProjectFeature(projectId, {
        key: keyValue.trim(),
        name: name.trim(),
        description: description.trim() || undefined,
        kind,
        default_variant: dv,
        enabled,
        rollout_key: kind === 'multivariant' ? (rolloutKey.trim() || undefined) : undefined,
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

  const addRule = (action: RuleActionType) => setRules((prev) => {
    const nums = prev.map((r) => (typeof r.priority === 'number' ? r.priority : 0));
    const next = nums.length ? Math.min(255, Math.max(...nums) + 1) : 0;
    return [
      ...prev,
      { id: genId(), action, flag_variant_id: action === RuleActionEnum.Assign ? (variants[0]?.id || '') : undefined, priority: next, expression: { group: { operator: 'and', children: [{ condition: { attribute: '', operator: 'eq', value: '' } }] } as any }, segment_id: undefined, is_customized: true, baseExpressionJson: undefined }
    ];
  });
  const removeRuleById = (id: string) => setRules((prev) => prev.filter((r) => r.id !== id));
  const updateRuleById = (id: string, patch: Partial<RuleFormItem>) => setRules((prev) => prev.map((r) => {
    if (r.id !== id) return r;
    const next: RuleFormItem = { ...r, ...patch } as RuleFormItem;
    if (Object.prototype.hasOwnProperty.call(patch, 'expression')) {
      if (next.segment_id) {
        const base = next.baseExpressionJson;
        try {
          const exprJson = JSON.stringify(next.expression || {});
          next.is_customized = !base || base !== exprJson;
        } catch {
          next.is_customized = true;
        }
      } else {
        next.is_customized = true;
      }
    }
    return next;
  }));
  const setRuleExpression = (id: string, expr: RuleConditionExpression) => {
    setRules(prev => prev.map(r => {
      if (r.id !== id) return r;
      const next: RuleFormItem = { ...r, expression: expr };
      if (next.segment_id) {
        const base = next.baseExpressionJson;
        try {
          const exprJson = JSON.stringify(expr || {});
          next.is_customized = !base || base !== exprJson;
        } catch {
          next.is_customized = true;
        }
      } else {
        next.is_customized = true;
      }
      return next;
    }));
  };
  const handleSelectSegment = (id: string, segId: string) => {
    const seg = (segments || []).find(s => s.id === segId);
    setRules(prev => prev.map(r => {
      if (r.id !== id) return r;
      if (!seg) {
        return { ...r, segment_id: undefined, baseExpressionJson: undefined, is_customized: true };
      }
      const baseJson = JSON.stringify(seg.conditions || {});
      return { ...r, segment_id: seg.id, expression: seg.conditions as any, baseExpressionJson: baseJson, is_customized: false };
    }));
  };
  const addRuleCondition = (ruleId: string) => setRules((prev) => prev.map((r) => (
    r.id === ruleId ? { ...r, expression: { condition: { attribute: '', operator: 'eq', value: '' } } } : r
  )));
  const removeRuleCondition = (ruleId: string, condIndex: number) => setRules((prev) => prev.map((r) => (
    r.id === ruleId ? { ...r, expression: { condition: { attribute: '', operator: 'eq', value: '' } } } : r
  )));
  const changeRuleCondition = (ruleId: string, condIndex: number, field: keyof RuleConditionItem, value: string) => setRules((prev) => prev.map((r) => (
    r.id === ruleId
      ? { ...r, expression: { condition: { attribute: '', operator: 'eq', value: '' } } }
      : r
  )));


  const canCreate = useMemo(() => {
    if (!keyValid || !name.trim()) return false;
    if (kind === 'multivariant') {
      if (!rolloutKey.trim()) return false;
      if (!hasAtLeastTwoVariants) return false;
      if (!variantsValid) return false;
    }
    if (rules.length > 0 && !rulesValid) return false;
    return !submitting;
  }, [keyValid, name, kind, rolloutKey, hasAtLeastTwoVariants, variantsValid, rulesValid, rules.length, submitting]);

  return (
    <Dialog open={open} onClose={onClose} fullWidth maxWidth="md">
      <DialogTitle sx={{ color: 'primary.main' }}>Create Feature</DialogTitle>
      <DialogContent>
        <Box sx={{ display: 'grid', gridTemplateColumns: { xs: '1fr', md: '1fr 1fr' }, gap: 2, mt: 1 }}>
          <TextField 
            label="Key" 
            value={keyValue} 
            onChange={(e) => setKeyValue(e.target.value)} 
            required 
            fullWidth 
            error={keyValue.trim().length > 0 && !keyValid}
            helperText={!keyValid && keyValue.trim().length > 0 ? 'Allowed: letters (a-z, A-Z), digits (0-9), hyphen (-), underscore (_), colon (:), @, !, #, $, dot (.)' : undefined}
          />
          <TextField label="Name" value={name} onChange={(e) => setName(e.target.value)} required fullWidth />
          <TextField label="Description" value={description} onChange={(e) => setDescription(e.target.value)} fullWidth multiline minRows={2} />
          <TextField select label="Kind" value={kind} onChange={(e) => { const v = e.target.value as FeatureKind; setKind(v); }} fullWidth>
            {kindOptions.map(k => (
              <MenuItem key={k} value={k}>{k}</MenuItem>
            ))}
          </TextField>
          {kind === 'multivariant' && (
            <Autocomplete
              freeSolo
              options={rolloutKeyOptions}
              value={rolloutKey}
              onChange={(_, val) => setRolloutKey(val || '')}
              onInputChange={(_, val) => setRolloutKey(val)}
              renderInput={(params) => (
                <TextField {...params} label="Rollout Key" required fullWidth helperText="Select from suggestions or type any attribute name" />
              )}
            />
          )}
          {kind === 'simple' ? (
            <TextField
              label="Default value (returned when enabled)"
              value={defaultVariant}
              onChange={(e) => setDefaultVariant(e.target.value)}
              fullWidth
              helperText="Any value; when feature is enabled, this exact value is returned"
            />
          ) : (
            <TextField
              label="Default variant (free text)"
              value={defaultVariant}
              onChange={(e) => setDefaultVariant(e.target.value)}
              fullWidth
              helperText="May be any string; not required to match defined variants"
            />
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
          </Box>
        )}

        {/* Rules (shown for both boolean and multivariant) */}
        <Box sx={{ mt: 3 }}>
          <Typography variant="subtitle1" sx={{ mb: 1 }}>Rules</Typography>

          {/* Assign section (multivariant only) */}
          {kind === 'multivariant' && (
            <Box sx={{ border: '2px solid', borderColor: 'success.light', borderRadius: 1, p: 1.5, mb: 2 }}>
              <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 1 }}>
                <Typography variant="subtitle2">Assign request to variant rules</Typography>
                <Button startIcon={<AddIcon />} onClick={() => addRule(RuleActionEnum.Assign)} size="small" disabled={variants.length === 0}>Add Rule</Button>
              </Box>
              <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>Route matching requests directly to a specific variant.</Typography>
              <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1.5 }}>
                {rules.filter(r => r.action === RuleActionEnum.Assign).sort((a,b) => (Number(a.priority || 0) - Number(b.priority || 0))).map((r) => (
                  <Box key={r.id} sx={{ border: '1px solid', borderColor: 'divider', borderRadius: 1, p: 1 }}>
                    <Box sx={{ display: 'grid', gridTemplateColumns: { xs: '1fr', md: '140px 1fr 40px' }, gap: 1, alignItems: 'center' }}>
                      <TextField
                        label="Priority"
                        type="number"
                        inputProps={{ min: 0, max: 255, step: 1 }}
                        value={r.priority}
                        onChange={(e) => updateRuleById(r.id, { priority: e.target.value === '' ? '' : Number(e.target.value) })}
                      />
                      <TextField
                        select
                        label="Target variant"
                        value={r.flag_variant_id || ''}
                        onChange={(e) => updateRuleById(r.id, { flag_variant_id: String(e.target.value) })}
                      >
                        {variants.map(v => (
                          <MenuItem key={v.id} value={v.id}>{v.name || v.id}</MenuItem>
                        ))}
                      </TextField>
                      <IconButton aria-label="delete-rule" onClick={() => removeRuleById(r.id)}>
                        <DeleteIcon />
                      </IconButton>
                    </Box>

                    <Box sx={{ mt: 1, display: 'flex', alignItems: 'center', gap: 1 }}>
                      <TextField
                        select
                        label="Segment template"
                        value={r.segment_id || ''}
                        onChange={(e) => handleSelectSegment(r.id, String(e.target.value))}
                        sx={{ minWidth: 240 }}
                     >
                        <MenuItem value="">Custom (no segment)</MenuItem>
                        {(segments || []).map((s) => (
                          <MenuItem key={s.id} value={s.id}>{s.name}</MenuItem>
                        ))}
                      </TextField>
                      {r.segment_id && (
                        <Chip size="small" color={r.is_customized ? 'warning' : 'success'} label={r.is_customized ? 'Customized' : 'From segment'} />
                      )}
                    </Box>

                    <Box sx={{ mt: 1 }}>
                      <Typography variant="body2" color="text.secondary">Conditions</Typography>
                      <ConditionExpressionBuilder value={r.expression} onChange={(expr) => updateRuleById(r.id, { expression: expr })} />
                    </Box>
                  </Box>
                ))}
                {rules.filter(r => r.action === RuleActionEnum.Assign).length === 0 && (
                  <Typography variant="body2" color="text.secondary">No assign rules</Typography>
                )}
              </Box>
            </Box>
          )}

          {/* Include section */}
          <Box sx={{ border: '2px solid', borderColor: 'info.light', borderRadius: 1, p: 1.5, mb: 2 }}>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 1 }}>
              <Typography variant="subtitle2">Include rules</Typography>
              <Button startIcon={<AddIcon />} onClick={() => addRule(RuleActionEnum.Include)} size="small">Add Rule</Button>
            </Box>
            <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>Define experiment segment. Matching requests participate in rollout between variants (for multivariant) or are subjected to flag action (for simple).</Typography>
            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1.5 }}>
              {rules.filter(r => r.action === RuleActionEnum.Include).sort((a,b) => (Number(a.priority || 0) - Number(b.priority || 0))).map((r) => (
                <Box key={r.id} sx={{ border: '1px solid', borderColor: 'divider', borderRadius: 1, p: 1 }}>
                  <Box sx={{ display: 'grid', gridTemplateColumns: { xs: '1fr', md: '140px 40px' }, gap: 1, alignItems: 'center' }}>
                    <TextField
                      label="Priority"
                      type="number"
                      inputProps={{ min: 0, max: 255, step: 1 }}
                      value={r.priority}
                      onChange={(e) => updateRuleById(r.id, { priority: e.target.value === '' ? '' : Number(e.target.value) })}
                    />
                    <IconButton aria-label="delete-rule" onClick={() => removeRuleById(r.id)}>
                      <DeleteIcon />
                    </IconButton>
                  </Box>
                  <Box sx={{ mt: 1, display: 'flex', alignItems: 'center', gap: 1 }}>
                    <TextField
                      select
                      label="Segment template"
                      value={r.segment_id || ''}
                      onChange={(e) => handleSelectSegment(r.id, String(e.target.value))}
                      sx={{ minWidth: 240 }}
                    >
                      <MenuItem value="">Custom (no segment)</MenuItem>
                      {(segments || []).map((s) => (
                        <MenuItem key={s.id} value={s.id}>{s.name}</MenuItem>
                      ))}
                    </TextField>
                    {r.segment_id && (
                      <Chip size="small" color={r.is_customized ? 'warning' : 'success'} label={r.is_customized ? 'Customized' : 'From segment'} />
                    )}
                  </Box>
                  <Box sx={{ mt: 1 }}>
                    <Typography variant="body2" color="text.secondary">Conditions</Typography>
                    <ConditionExpressionBuilder value={r.expression} onChange={(expr) => updateRuleById(r.id, { expression: expr })} />
                  </Box>
                </Box>
              ))}
              {rules.filter(r => r.action === RuleActionEnum.Include).length === 0 && (
                <Typography variant="body2" color="text.secondary">No include rules</Typography>
              )}
            </Box>
          </Box>

          {/* Exclude section */}
          <Box sx={{ border: '2px solid', borderColor: 'error.light', borderRadius: 1, p: 1.5 }}>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 1 }}>
              <Typography variant="subtitle2">Exclude rules</Typography>
              <Button startIcon={<AddIcon />} onClick={() => addRule(RuleActionEnum.Exclude)} size="small">Add Rule</Button>
            </Box>
            <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>Matching requests are excluded from experiment and do not participate in variant rollout.</Typography>
            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1.5 }}>
              {rules.filter(r => r.action === RuleActionEnum.Exclude).sort((a,b) => (Number(a.priority || 0) - Number(b.priority || 0))).map((r) => (
                <Box key={r.id} sx={{ border: '1px solid', borderColor: 'divider', borderRadius: 1, p: 1 }}>
                  <Box sx={{ display: 'grid', gridTemplateColumns: { xs: '1fr', md: '140px 40px' }, gap: 1, alignItems: 'center' }}>
                    <TextField
                      label="Priority"
                      type="number"
                      inputProps={{ min: 0, max: 255, step: 1 }}
                      value={r.priority}
                      onChange={(e) => updateRuleById(r.id, { priority: e.target.value === '' ? '' : Number(e.target.value) })}
                    />
                    <IconButton aria-label="delete-rule" onClick={() => removeRuleById(r.id)}>
                      <DeleteIcon />
                    </IconButton>
                  </Box>
                  <Box sx={{ mt: 1, display: 'flex', alignItems: 'center', gap: 1 }}>
                    <TextField
                      select
                      label="Segment template"
                      value={r.segment_id || ''}
                      onChange={(e) => handleSelectSegment(r.id, String(e.target.value))}
                      sx={{ minWidth: 240 }}
                    >
                      <MenuItem value="">Custom (no segment)</MenuItem>
                      {(segments || []).map((s) => (
                        <MenuItem key={s.id} value={s.id}>{s.name}</MenuItem>
                      ))}
                    </TextField>
                    {r.segment_id && (
                      <Chip size="small" color={r.is_customized ? 'warning' : 'success'} label={r.is_customized ? 'Customized' : 'From segment'} />
                    )}
                  </Box>
                  <Box sx={{ mt: 1 }}>
                    <Typography variant="body2" color="text.secondary">Conditions</Typography>
                    <ConditionExpressionBuilder value={r.expression} onChange={(expr) => updateRuleById(r.id, { expression: expr })} />
                  </Box>
                </Box>
              ))}
              {rules.filter(r => r.action === RuleActionEnum.Exclude).length === 0 && (
                <Typography variant="body2" color="text.secondary">No exclude rules</Typography>
              )}
            </Box>
          </Box>

          {hasDuplicatePriorities && (
            <Box sx={{ mt: 1 }}>
              <Chip size="small" color="warning" label="Duplicate priorities detected" />
            </Box>
          )}
        </Box>

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
