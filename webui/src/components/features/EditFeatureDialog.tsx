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
} from '@mui/material';
import Autocomplete from '@mui/material/Autocomplete';
import ConditionExpressionBuilder from '../conditions/ConditionExpressionBuilder';
import { Add, Delete, Sync } from '@mui/icons-material';
import { useMutation, useQueryClient, useQuery } from '@tanstack/react-query';
import apiClient from '../../api/apiClient';
import type { CreateFeatureRequest, FeatureDetailsResponse, RuleConditionExpression, RuleAction as RuleActionType, Segment, ProjectTag } from '../../generated/api/client';
import { FeatureKind as FeatureKindEnum, RuleAction as RuleActionEnum } from '../../generated/api/client';
import TagSelector from './TagSelector';

export interface EditFeatureDialogProps {
  open: boolean;
  onClose: () => void;
  featureDetails: FeatureDetailsResponse | null;
}

const uuid = () => (typeof crypto !== 'undefined' && 'randomUUID' in crypto ? crypto.randomUUID() : Math.random().toString(36).slice(2));

const EditFeatureDialog: React.FC<EditFeatureDialogProps> = ({ open, onClose, featureDetails }) => {
  const queryClient = useQueryClient();

  type VariantForm = { id: string; name: string; rollout_percent: number };
  type RuleForm = { id: string; priority: number; action: RuleActionType; flag_variant_id?: string; expression: RuleConditionExpression; segment_id?: string; is_customized?: boolean; baseExpressionJson?: string };

  const [keyVal, setKeyVal] = useState('');
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [rolloutKey, setRolloutKey] = useState('');
  const [kind, setKind] = useState<string>(FeatureKindEnum.Simple);
  const [enabled, setEnabled] = useState<boolean>(false);
  const [variants, setVariants] = useState<VariantForm[]>([]);
  const [defaultVariant, setDefaultVariant] = useState<string>('');
  const [rules, setRules] = useState<RuleForm[]>([]);
  const [selectedTags, setSelectedTags] = useState<ProjectTag[]>([]);

  const [syncing, setSyncing] = useState<Record<string, boolean>>({});
  const [syncErrors, setSyncErrors] = useState<Record<string, string | undefined>>({});

  const [error, setError] = useState<string | null>(null);

  // Initialize form from featureDetails
  useEffect(() => {
    if (!open || !featureDetails) return;
    const f = featureDetails.feature;
    setKeyVal(f.key);
    setName(f.name);
    setDescription(f.description || '');
    setKind(f.kind);
    setRolloutKey(f.rollout_key || '');
    setEnabled(f.enabled);
    const vars = (featureDetails.variants || []).map(v => ({ id: v.id, name: v.name, rollout_percent: v.rollout_percent }));
    setVariants(vars);
    setDefaultVariant(f.default_variant || '');
    const rls = (featureDetails.rules || []).map(r => ({ id: r.id, priority: r.priority, action: r.action as RuleActionType, flag_variant_id: r.flag_variant_id, expression: r.conditions as any, segment_id: (r as any).segment_id, is_customized: (r as any).is_customized }));
    setRules(rls);
    setSelectedTags(featureDetails.tags || []);
    setError(null);
  }, [open, featureDetails]);

  const featureId = featureDetails?.feature.id;
  const projectId = featureDetails?.feature.project_id;

  // Load project segments for segment templates
  const { data: segments } = useQuery<Segment[]>({
    queryKey: ['project-segments', projectId],
    queryFn: async () => {
      const res = await apiClient.listProjectSegments(projectId!);
      const resp = res.data as any;
      return Array.isArray(resp?.items) ? (resp.items as Segment[]) : (resp as Segment[]);
    },
    enabled: Boolean(projectId),
  });

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
  };
  const removeVariant = (id: string) => {
    // prevent removing variants referenced by assign rules
    if (rules.some(r => r.action === RuleActionEnum.Assign && r.flag_variant_id === id)) return;
    setVariants(prev => prev.filter(v => v.id !== id));
  };

  const addRule = (action: RuleActionType) => {
    const id = uuid();
    const priority = (rules.length ? Math.max(...rules.map(r => r.priority)) + 1 : 1);
    const firstVariant = variants[0]?.id;
    setRules(prev => [...prev, { id, priority, action, flag_variant_id: action === RuleActionEnum.Assign ? firstVariant : undefined, expression: { group: { operator: 'and', children: [{ condition: { attribute: '', operator: 'eq', value: '' } }] } as any } }]);
  };
  const removeRule = (id: string) => setRules(prev => prev.filter(r => r.id !== id));

  // Update rule expression and mark customized when segment is attached
  const setRuleExpression = (id: string, expr: RuleConditionExpression) => {
    setRules(prev => prev.map(r => {
      if (r.id !== id) return r;
      const next: RuleForm = { ...r, expression: expr };
      if (next.segment_id) {
        try {
          const exprJson = JSON.stringify(expr || {});
          next.is_customized = !next.baseExpressionJson || next.baseExpressionJson !== exprJson;
        } catch {
          next.is_customized = true;
        }
      }
      return next;
    }));
  };

  // Attach/detach a segment to a rule and copy its conditions
  const handleSelectSegment = (id: string, segId: string) => {
    const seg = (segments || []).find(s => s.id === segId);
    setRules(prev => prev.map(r => {
      if (r.id !== id) return r;
      if (!seg) {
        return { ...r, segment_id: undefined, is_customized: true, baseExpressionJson: undefined } as RuleForm;
      }
      const baseJson = JSON.stringify(seg.conditions || {});
      return { ...r, segment_id: seg.id, expression: seg.conditions as any, is_customized: false, baseExpressionJson: baseJson } as RuleForm;
    }));
  };

  const handleSyncRule = async (ruleId: string) => {
    if (!featureId) return;
    setSyncErrors(prev => ({ ...prev, [ruleId]: undefined }));
    setSyncing(prev => ({ ...prev, [ruleId]: true }));
    try {
      const res = await apiClient.syncCustomizedFeatureRule(featureId, ruleId);
      const ru: any = res.data as any;
      setRules(prev => prev.map(r => {
        if (r.id !== ruleId) return r;
        const updatedExpr = (ru && ru.conditions) ? (ru.conditions as any) : r.expression;
        const updated: RuleForm = {
          ...r,
          priority: ru?.priority ?? r.priority,
          action: ru?.action ?? r.action,
          flag_variant_id: ru?.flag_variant_id ?? r.flag_variant_id,
          expression: updatedExpr,
          segment_id: ru?.segment_id ?? r.segment_id,
          is_customized: Boolean(ru?.is_customized),
          baseExpressionJson: updatedExpr ? JSON.stringify(updatedExpr) : r.baseExpressionJson,
        };
        return updated;
      }));
      // refresh caches
      if (featureId) {
        await queryClient.invalidateQueries({ queryKey: ['feature-details', featureId] });
      }
      if (projectId) {
        await queryClient.invalidateQueries({ queryKey: ['project-features', projectId] });
      }
    } catch (e: any) {
      setSyncErrors(prev => ({ ...prev, [ruleId]: e?.message || 'Failed to sync rule' }));
    } finally {
      setSyncing(prev => ({ ...prev, [ruleId]: false }));
    }
  };

  const parseValueSmart = (input: string): any => {
    const t = input.trim();
    if (t === '') return '';
    if (t === 'true') return true;
    if (t === 'false') return false;
    if (!isNaN(Number(t))) return Number(t);
    try { return JSON.parse(t); } catch { return input; }
  };

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

  const validate = (): string | null => {
    if (!keyVal.trim()) return 'Key is required';
    if (!name.trim()) return 'Name is required';
    if (!kind) return 'Kind is required';
    if (kind === FeatureKindEnum.Multivariant && !rolloutKey.trim()) return 'Rollout Key is required for multivariant features';
    for (const v of variants) {
      if (!v.name.trim()) return 'Variant name cannot be empty';
      if (v.rollout_percent < 0 || v.rollout_percent > 100) return 'Variant rollout percent must be between 0 and 100';
    }
    for (const r of rules) {
      if (r.action === RuleActionEnum.Assign) {
        if (!r.flag_variant_id || !variants.find(v => v.id === r.flag_variant_id)) return 'Assign rules must target an existing variant';
      }
      if (!hasValidLeaf(r.expression)) return 'Each rule must have at least one valid condition';
    }
    if (hasDuplicatePriorities) return 'Rule priorities must be unique within each rules section (Assign/Include/Exclude)';
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
      default_variant: defaultVariant,
      enabled,
      rollout_key: kind === FeatureKindEnum.Multivariant ? (rolloutKey.trim() || undefined) : undefined,
      variants: variants.map(v => ({ id: v.id, name: v.name, rollout_percent: Number(v.rollout_percent) })),
      rules: rules.map(r => ({
        id: r.id,
        priority: r.priority,
        action: r.action as any,
        flag_variant_id: r.action === RuleActionEnum.Assign ? r.flag_variant_id : undefined,
        conditions: r.expression,
        segment_id: r.segment_id || undefined,
        is_customized: r.segment_id ? Boolean(r.is_customized) : true,
      })),
    };

    updateMutation.mutate(payload);
    
    // Update tags after feature update
    if (featureId) {
      updateFeatureTags(featureId);
    }
  };

  const updateFeatureTags = async (featureId: string) => {
    if (!featureDetails) return;
    
    const currentTags = featureDetails.tags || [];
    const newTags = selectedTags;
    
    // Find tags to add
    const tagsToAdd = newTags.filter(newTag => 
      !currentTags.some(currentTag => currentTag.id === newTag.id)
    );
    
    // Find tags to remove
    const tagsToRemove = currentTags.filter(currentTag => 
      !newTags.some(newTag => newTag.id === currentTag.id)
    );
    
    // Add new tags
    for (const tag of tagsToAdd) {
      try {
        await apiClient.addFeatureTag(featureId, { tag_id: tag.id });
      } catch (err) {
        console.warn('Failed to add tag to feature:', err);
      }
    }
    
    // Remove old tags
    for (const tag of tagsToRemove) {
      try {
        await apiClient.removeFeatureTag(featureId, tag.id);
      } catch (err) {
        console.warn('Failed to remove tag from feature:', err);
      }
    }
  };

  const disabled = !featureDetails || updateMutation.isPending;

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

  const ruleOperatorOptions: any[] = [];
  const featureKindOptions = Object.values(FeatureKindEnum);
  const rolloutKeyOptions = ['user.id', 'user.email'];

  // Default variant is free text now; allow deletion unless referenced by an Assign rule
  const canDeleteVariant = (id: string) => !rules.some(r => r.action === RuleActionEnum.Assign && r.flag_variant_id === id);
  const onId = useMemo(() => variants.find(v => (v.name || v.id).toLowerCase() === 'on')?.id, [variants]);
  const offId = useMemo(() => variants.find(v => (v.name || v.id).toLowerCase() === 'off')?.id, [variants]);

  return (
    <Dialog open={open} onClose={disabled ? undefined : onClose} fullWidth maxWidth="md">
      <DialogTitle sx={{ color: 'primary.main' }}>Edit Feature</DialogTitle>
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
            <Grid container spacing={1.5}>
              <Grid item xs={12} md={6}>
                <TextField label="Key" value={keyVal} onChange={(e) => setKeyVal(e.target.value)} fullWidth disabled size="small" />
              </Grid>
              <Grid item xs={12} md={6}>
                <TextField label="Name" value={name} onChange={(e) => setName(e.target.value)} fullWidth size="small" />
              </Grid>
              <Grid item xs={12}>
                <TextField label="Description" value={description} onChange={(e) => setDescription(e.target.value)} fullWidth multiline minRows={2} />
              </Grid>
              <Grid item xs={12} md={6}>
                <FormControl fullWidth size="small">
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
              {kind === FeatureKindEnum.Multivariant && (
                <Grid item xs={12} md={6}>
                  <Autocomplete
                    freeSolo
                    options={rolloutKeyOptions}
                    value={rolloutKey}
                    onChange={(_, val) => setRolloutKey(val || '')}
                    onInputChange={(_, val) => setRolloutKey(val)}
                    renderInput={(params) => (
                      <TextField {...params} label="Rollout Key" required fullWidth size="small" helperText="Select from suggestions or type any attribute name" />
                    )}
                  />
                </Grid>
              )}
            </Grid>

            {/* Tags */}
            {projectId && (
              <Box>
                <TagSelector
                  projectId={projectId}
                  selectedTags={selectedTags}
                  onChange={setSelectedTags}
                  disabled={disabled}
                />
              </Box>
            )}

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
                        <TextField label={`Variant name #${idx + 1}`} value={v.name} onChange={(e) => setVariants(prev => prev.map(x => x.id === v.id ? { ...x, name: e.target.value } : x))} fullWidth size="small" />
                      </Grid>
                      <Grid item xs={12} md={3}>
                        <TextField type="number" label="Rollout %" inputProps={{ min: 0, max: 100 }} value={v.rollout_percent}
                          onChange={(e) => setVariants(prev => prev.map(x => x.id === v.id ? { ...x, rollout_percent: Number(e.target.value) } : x))} fullWidth size="small" />
                      </Grid>
                      <Grid item xs={12} md={4}>
                        <TextField label="ID" value={v.id} fullWidth disabled size="small" />
                      </Grid>
                      <Grid item xs={12} md={1} sx={{ display: 'flex', alignItems: 'center' }}>
                        <Tooltip title={canDeleteVariant(v.id) ? 'Delete variant' : 'Cannot delete: used in assign rules'}>
                          <span>
                            <IconButton onClick={() => removeVariant(v.id)} disabled={!canDeleteVariant(v.id)} aria-label="delete variant" size="small">
                              <Delete />
                            </IconButton>
                          </span>
                        </Tooltip>
                      </Grid>
                    </React.Fragment>
                  ))}
                </Grid>
                <Box sx={{ display: 'flex', gap: 2, alignItems: 'center' }}>
                  <TextField
                    label="Default variant (free text)"
                    value={defaultVariant}
                    onChange={(e) => setDefaultVariant(e.target.value)}
                    fullWidth
                    size="small"
                    helperText="May be any string; not required to match defined variants"
                  />
                  <Typography variant="body2" color="text.secondary">Variants total: {variants.length}</Typography>
                </Box>
              </>
            ) : (
              <Box sx={{ display: 'flex', gap: 2, alignItems: 'center' }}>
                <TextField
                  label="Default value (returned when enabled)"
                  value={defaultVariant}
                  onChange={(e) => setDefaultVariant(e.target.value)}
                  fullWidth
                  size="small"
                  helperText="Any value; when feature is enabled, this exact value is returned"
                />
              </Box>
            )}

            <Divider />

            {/* Rules */}
            <Box>
              <Typography variant="subtitle1" sx={{ mb: 1 }}>Rules</Typography>

              {/* Assign section */}
              <Box sx={{ border: '2px solid', borderColor: 'success.light', borderRadius: 1, p: 1.5, mb: 2 }}>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 1 }}>
                  <Typography variant="subtitle2">Assign request to variant rules</Typography>
                  <Button startIcon={<Add />} onClick={() => addRule(RuleActionEnum.Assign)} size="small" disabled={variants.length === 0}>Add Rule</Button>
                </Box>
                <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>Route matching requests directly to a specific variant.</Typography>
                <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1.5 }}>
                  {rules.filter(r => r.action === RuleActionEnum.Assign).sort((a,b) => a.priority - b.priority).map((r) => (
                    <Box key={r.id} sx={{ border: '1px solid', borderColor: 'divider', borderRadius: 1, p: 1 }}>
                      <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, flexWrap: 'wrap' }}>
                        <TextField type="number" label="Priority" value={r.priority} sx={{ width: 140 }}
                          onChange={(e) => setRules(prev => prev.map(x => x.id === r.id ? { ...x, priority: Number(e.target.value) } : x))} size="small" />
                        <FormControl sx={{ minWidth: 240 }} size="small">
                          <InputLabel id={`rule-var-${r.id}`}>Target variant</InputLabel>
                          <Select labelId={`rule-var-${r.id}`} label="Target variant" value={r.flag_variant_id || ''}
                            onChange={(e) => setRules(prev => prev.map(x => x.id === r.id ? { ...x, flag_variant_id: String(e.target.value) } : x))}>
                            {variants.map(v => (
                              <MenuItem key={v.id} value={v.id}>{variantNameById[v.id]}</MenuItem>
                            ))}
                          </Select>
                        </FormControl>
                        <IconButton onClick={() => removeRule(r.id)} aria-label="delete rule" size="small"><Delete /></IconButton>
                      </Box>

                      <Box sx={{ mt: 1, display: 'flex', alignItems: 'center', gap: 1 }}>
                        <TextField
                          select
                          label="Segment template"
                          size="small"
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
                          <>
                            <Chip size="small" color={r.is_customized ? 'warning' : 'success'} label={r.is_customized ? 'Customized' : 'From segment'} />
                            {r.is_customized && (
                              <Tooltip title="Sync with segment">
                                <span>
                                  <IconButton onClick={() => handleSyncRule(r.id)} disabled={!!syncing[r.id]} aria-label="sync rule" size="small">
                                    <Sync />
                                  </IconButton>
                                </span>
                              </Tooltip>
                            )}
                          </>
                        )}
                      </Box>

                      <Box sx={{ mt: 1 }}>
                        <Typography variant="body2" color="text.secondary">Conditions</Typography>
                        <ConditionExpressionBuilder value={r.expression} onChange={(expr) => setRuleExpression(r.id, expr)} />
                      </Box>
                    </Box>
                  ))}
                  {rules.filter(r => r.action === RuleActionEnum.Assign).length === 0 && (
                    <Typography variant="body2" color="text.secondary">No assign rules</Typography>
                  )}
                </Box>
              </Box>

              {/* Include section */}
              <Box sx={{ border: '2px solid', borderColor: 'info.light', borderRadius: 1, p: 1.5, mb: 2 }}>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 1 }}>
                  <Typography variant="subtitle2">Include rules</Typography>
                  <Button startIcon={<Add />} onClick={() => addRule(RuleActionEnum.Include)} size="small">Add Rule</Button>
                </Box>
                <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>Define experiment segment. Matching requests participate in rollout between variants (for multivariant) or are subjected to flag action (for simple).</Typography>
                <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1.5 }}>
                  {rules.filter(r => r.action === RuleActionEnum.Include).sort((a,b) => a.priority - b.priority).map((r) => (
                    <Box key={r.id} sx={{ border: '1px solid', borderColor: 'divider', borderRadius: 1, p: 1 }}>
                      <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, flexWrap: 'wrap' }}>
                        <TextField type="number" label="Priority" value={r.priority} sx={{ width: 140 }}
                          onChange={(e) => setRules(prev => prev.map(x => x.id === r.id ? { ...x, priority: Number(e.target.value) } : x))} size="small" />
                        <IconButton onClick={() => removeRule(r.id)} aria-label="delete rule" size="small"><Delete /></IconButton>
                      </Box>

                      <Box sx={{ mt: 1, display: 'flex', alignItems: 'center', gap: 1 }}>
                        <TextField
                          select
                          label="Segment template"
                          size="small"
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
                          <>
                            <Chip size="small" color={r.is_customized ? 'warning' : 'success'} label={r.is_customized ? 'Customized' : 'From segment'} />
                            {r.is_customized && (
                              <Tooltip title="Sync with segment">
                                <span>
                                  <IconButton onClick={() => handleSyncRule(r.id)} disabled={!!syncing[r.id]} aria-label="sync rule" size="small">
                                    <Sync />
                                  </IconButton>
                                </span>
                              </Tooltip>
                            )}
                          </>
                        )}
                      </Box>

                      <Box sx={{ mt: 1 }}>
                        <Typography variant="body2" color="text.secondary">Conditions</Typography>
                        <ConditionExpressionBuilder value={r.expression} onChange={(expr) => setRuleExpression(r.id, expr)} />
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
                  <Button startIcon={<Add />} onClick={() => addRule(RuleActionEnum.Exclude)} size="small">Add Rule</Button>
                </Box>
                <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>Matching requests are excluded from experiment and do not participate in variant rollout.</Typography>
                <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1.5 }}>
                  {rules.filter(r => r.action === RuleActionEnum.Exclude).sort((a,b) => a.priority - b.priority).map((r) => (
                    <Box key={r.id} sx={{ border: '1px solid', borderColor: 'divider', borderRadius: 1, p: 1 }}>
                      <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, flexWrap: 'wrap' }}>
                        <TextField type="number" label="Priority" value={r.priority} sx={{ width: 140 }}
                          onChange={(e) => setRules(prev => prev.map(x => x.id === r.id ? { ...x, priority: Number(e.target.value) } : x))} size="small" />
                        <IconButton onClick={() => removeRule(r.id)} aria-label="delete rule" size="small"><Delete /></IconButton>
                      </Box>

                      <Box sx={{ mt: 1, display: 'flex', alignItems: 'center', gap: 1 }}>
                        <TextField
                          select
                          label="Segment template"
                          size="small"
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
                          <>
                            <Chip size="small" color={r.is_customized ? 'warning' : 'success'} label={r.is_customized ? 'Customized' : 'From segment'} />
                            {r.is_customized && (
                              <Tooltip title="Sync with segment">
                                <span>
                                  <IconButton onClick={() => handleSyncRule(r.id)} disabled={!!syncing[r.id]} aria-label="sync rule" size="small">
                                    <Sync />
                                  </IconButton>
                                </span>
                              </Tooltip>
                            )}
                          </>
                        )}
                      </Box>

                      <Box sx={{ mt: 1 }}>
                        <Typography variant="body2" color="text.secondary">Conditions</Typography>
                        <ConditionExpressionBuilder value={r.expression} onChange={(expr) => setRuleExpression(r.id, expr)} />
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
                  <Chip size="small" color="warning" label="Duplicate priorities detected within a section" />
                </Box>
              )}
            </Box>
          </Box>
        )}
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose} disabled={disabled} size="small">Cancel</Button>
        <Button onClick={onSave} variant="contained" disabled={disabled} size="small">{updateMutation.isPending ? 'Saving…' : 'Save'}</Button>
      </DialogActions>
    </Dialog>
  );
};

export default EditFeatureDialog;
