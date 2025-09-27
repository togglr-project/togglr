import React from 'react';
import { Box, Button, Grid, IconButton, MenuItem, Select, TextField, Typography, Tooltip } from '@mui/material';
import Autocomplete from '@mui/material/Autocomplete';
import type { SelectChangeEvent } from '@mui/material/Select';
import { Add as AddIcon, Delete as DeleteIcon, PlaylistAdd as AddGroupIcon } from '@mui/icons-material';
import { useQuery } from '@tanstack/react-query';
import apiClient from '../../api/apiClient';
import type { RuleConditionExpression, RuleConditionGroup, RuleCondition, LogicalOperator as LogicalOperatorType, RuleOperator as RuleOperatorType } from '../../generated/api/client';
import { LogicalOperator, RuleOperator } from '../../generated/api/client';

// UI models keep value as string for editing; convert on output
export interface UiCondition { attribute: string; operator: RuleOperatorType; value: string }
export interface UiGroup { operator: LogicalOperatorType; children: UiExpression[] }
export interface UiExpression { condition?: UiCondition; group?: UiGroup }

const ruleOperatorOptions: RuleOperatorType[] = Object.values(RuleOperator);
const logicalOperatorOptions: LogicalOperatorType[] = Object.values(LogicalOperator);

const parseValueSmart = (input: string): string | number | boolean => {
  const t = (input ?? '').trim();
  if (t === '') return '';
  if (t === 'true') return true;
  if (t === 'false') return false;
  if (!Number.isNaN(Number(t)) && /^-?\d+(\.\d+)?$/.test(t)) return Number(t);
  try { return JSON.parse(t); } catch { return input; }
};

const valueToString = (val: unknown): string => {
  if (val === undefined || val === null) return '';
  if (typeof val === 'string') return val;
  try { return JSON.stringify(val); } catch { return String(val); }
};

export const toUiExpression = (expr?: RuleConditionExpression): UiExpression => {
  if (!expr || (!expr.condition && !expr.group)) {
    return { group: { operator: LogicalOperator.And, children: [{ condition: { attribute: '', operator: RuleOperator.Eq, value: '' } }] } };
  }
  if (expr.condition) {
    const c = expr.condition as RuleCondition;
    return { condition: { attribute: c.attribute || '', operator: c.operator || RuleOperator.Eq, value: valueToString(c.value) } };
  }
  if (expr.group) {
    const g = expr.group as RuleConditionGroup;
    return {
      group: {
        operator: g.operator || LogicalOperator.And,
        children: (g.children || []).map(ch => toUiExpression(ch))
      }
    };
  }
  return { group: { operator: LogicalOperator.And, children: [{ condition: { attribute: '', operator: RuleOperator.Eq, value: '' } }] } };
};

export const fromUiExpression = (ui: UiExpression): RuleConditionExpression => {
  if (ui.condition) {
    return { condition: { attribute: ui.condition.attribute || '', operator: ui.condition.operator || RuleOperator.Eq, value: parseValueSmart(ui.condition.value) } };
  }
  if (ui.group) {
    return {
      group: {
        operator: ui.group.operator || LogicalOperator.And,
        children: (ui.group.children || []).map(fromUiExpression),
      }
    };
  }
  // Fallback to empty condition
  return { condition: { attribute: '', operator: RuleOperator.Eq, value: '' } };
};

const ensureAtLeastOneChild = (g: UiGroup): UiGroup => {
  if (!g.children || g.children.length === 0) {
    return { ...g, children: [{ condition: { attribute: '', operator: RuleOperator.Eq, value: '' } }] };
  }
  return g;
};

const ConditionRow: React.FC<{
  value: UiCondition;
  onChange: (next: UiCondition) => void;
  onRemove?: () => void;
  attributeOptions?: string[];
}> = ({ value, onChange, onRemove, attributeOptions }) => {
  return (
    <Grid container spacing={1} alignItems="center" sx={{ mt: 0.5 }}>
      <Grid item xs={12} md={4}>
        <Autocomplete
          freeSolo
          disableClearable
          options={attributeOptions || []}
          inputValue={value.attribute}
          onInputChange={(_, newInputValue) => onChange({ ...value, attribute: newInputValue })}
          onChange={(_, newValue) => onChange({ ...value, attribute: (newValue as string) || '' })}
          renderInput={(params) => (
            <TextField
              {...params}
              label="Attribute"
              required
              fullWidth
              size="small"
            />
          )}
        />
      </Grid>
      <Grid item xs={12} md={3}>
        <TextField
          select
          label="Operator"
          value={value.operator}
          onChange={(e) => onChange({ ...value, operator: e.target.value as RuleOperatorType })}
          fullWidth
          size="small"
        >
          {ruleOperatorOptions.map(op => (
            <MenuItem key={op} value={op}>{op}</MenuItem>
          ))}
        </TextField>
      </Grid>
      <Grid item xs={12} md={4}>
        <TextField
          label="Value"
          value={value.value}
          onChange={(e) => onChange({ ...value, value: e.target.value })}
          fullWidth
          size="small"
        />
      </Grid>
      <Grid item xs={12} md={1}>
        {onRemove && (
          <IconButton aria-label="delete-condition" onClick={onRemove}>
            <DeleteIcon />
          </IconButton>
        )}
      </Grid>
    </Grid>
  );
};

const GroupEditor: React.FC<{
  value: UiGroup;
  onChange: (next: UiGroup) => void;
  canRemove?: boolean;
  onRemove?: () => void;
  depth?: number;
  attributeOptions?: string[];
}> = ({ value, onChange, canRemove, onRemove, depth = 0, attributeOptions }) => {
  const g = ensureAtLeastOneChild(value);

  const setOperator = (op: LogicalOperatorType) => onChange({ ...g, operator: op });
  const updateChild = (idx: number, child: UiExpression) => {
    const next = g.children.slice();
    next[idx] = child;
    onChange({ ...g, children: next });
  };
  const addCondition = () => onChange({ ...g, children: [...g.children, { condition: { attribute: '', operator: RuleOperator.Eq, value: '' } }] });
  const addGroup = () => onChange({ ...g, children: [...g.children, { group: { operator: LogicalOperator.And, children: [{ condition: { attribute: '', operator: RuleOperator.Eq, value: '' } }] } }] });
  const removeChild = (idx: number) => onChange({ ...g, children: g.children.filter((_, i) => i !== idx) });

  return (
    <Box sx={{ borderLeft: depth === 0 ? 'none' : '2px solid', borderColor: 'divider', pl: depth === 0 ? 0 : 2, mt: depth === 0 ? 0 : 1 }}>
      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 1 }}>
        <Typography variant="subtitle2">Group</Typography>
        <Select size="small" value={g.operator} onChange={(e: SelectChangeEvent) => setOperator(e.target.value as LogicalOperatorType)}>
          {logicalOperatorOptions.map(op => (
            <MenuItem key={op} value={op}>{op}</MenuItem>
          ))}
        </Select>
        <Box sx={{ flexGrow: 1 }} />
        {canRemove && onRemove && (
          <Tooltip title="Delete group">
            <IconButton size="small" onClick={onRemove}><DeleteIcon fontSize="small" /></IconButton>
          </Tooltip>
        )}
      </Box>

      <Box>
        {g.children.map((ch, idx) => (
          <Box key={idx} sx={{ mb: 1 }}>
            {ch.condition ? (
              <ConditionRow
                value={ch.condition}
                onChange={(next) => updateChild(idx, { condition: next })}
                onRemove={() => removeChild(idx)}
                attributeOptions={attributeOptions}
              />
            ) : ch.group ? (
              <GroupEditor
                value={ch.group}
                onChange={(next) => updateChild(idx, { group: next })}
                canRemove
                onRemove={() => removeChild(idx)}
                depth={depth + 1}
                attributeOptions={attributeOptions}
              />
            ) : null}
          </Box>
        ))}
      </Box>

      <Box sx={{ display: 'flex', gap: 1, mt: 1 }}>
        <Button size="small" startIcon={<AddIcon />} onClick={addCondition}>Add Condition</Button>
        <Button size="small" startIcon={<AddGroupIcon />} onClick={addGroup}>Add Group</Button>
      </Box>
    </Box>
  );
};

const ConditionExpressionBuilder: React.FC<{
  value?: RuleConditionExpression;
  onChange: (expr: RuleConditionExpression) => void;
  attributeOptions?: string[];
}> = ({ value, onChange, attributeOptions }) => {
  const [ui, setUi] = React.useState<UiExpression>(() => toUiExpression(value));

  React.useEffect(() => {
    setUi(toUiExpression(value));
  }, [value]);

  const updateUi = (nextUi: UiExpression) => {
    setUi(nextUi);
    onChange(fromUiExpression(nextUi));
  };

  // Ensure root is always a group for easier editing
  const ensureRootGroup = (u: UiExpression): UiGroup => {
    if (u.group) return ensureAtLeastOneChild(u.group);
    if (u.condition) return { operator: LogicalOperator.And, children: [u] };
    return { operator: LogicalOperator.And, children: [{ condition: { attribute: '', operator: RuleOperator.Eq, value: '' } }] };
  };

  const rootGroup = ensureRootGroup(ui);

  const { data: attrsData } = useQuery<{ name: string }[]>({
    queryKey: ['rule-attributes'],
    queryFn: async () => {
      const res = await apiClient.listRuleAttributes();
      return res.data as { name: string }[];
    },
  });
  const fetchedAttrNames = React.useMemo(() => (Array.isArray(attrsData) ? attrsData.map((a: { name: string }) => a.name).filter(Boolean) : []), [attrsData]);
  const effectiveOptions = attributeOptions ?? fetchedAttrNames;

  return (
    <Box>
      <GroupEditor value={rootGroup} onChange={(g) => updateUi({ group: g })} attributeOptions={effectiveOptions} />
    </Box>
  );
};

export default ConditionExpressionBuilder;
