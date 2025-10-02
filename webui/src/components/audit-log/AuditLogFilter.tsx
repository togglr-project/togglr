import React from 'react';
import { Box, TextField, MenuItem, Select, InputLabel, FormControl } from '@mui/material';

export interface AuditLogFilterValue {
  environmentKey?: string;
  entity?: string;
  entityId?: string;
  actor?: string;
  from?: string; // ISO datetime-local
  to?: string;   // ISO datetime-local
  sortBy?: string;
  sortOrder?: 'asc' | 'desc';
}

interface Props {
  value: AuditLogFilterValue;
  environments: { id: number; key: string; name: string }[];
  onChange: (next: AuditLogFilterValue) => void;
}

const AuditLogFilter: React.FC<Props> = ({ value, environments, onChange }) => {
  const handle = (patch: Partial<AuditLogFilterValue>) => onChange({ ...value, ...patch });

  return (
    <Box sx={{ display: 'flex', gap: 2, flexWrap: 'wrap', alignItems: 'center' }}>
      <FormControl sx={{ minWidth: 160 }}>
        <InputLabel id="env-key-label">Environment</InputLabel>
        <Select
          labelId="env-key-label"
          label="Environment"
          size="small"
          value={value.environmentKey ?? ''}
          onChange={(e) => handle({ environmentKey: e.target.value || undefined })}
        >
          <MenuItem value="">
            <em>All</em>
          </MenuItem>
          {environments.map((e) => (
            <MenuItem key={e.id} value={e.key}>
              {e.key}
            </MenuItem>
          ))}
        </Select>
      </FormControl>

      <TextField
        label="Entity"
        size="small"
        placeholder="feature, rule, ..."
        value={value.entity ?? ''}
        onChange={(e) => handle({ entity: e.target.value || undefined })}
      />

      <TextField
        label="Entity ID"
        size="small"
        placeholder="UUID"
        value={value.entityId ?? ''}
        onChange={(e) => handle({ entityId: e.target.value || undefined })}
      />

      <TextField
        label="Actor"
        size="small"
        placeholder="username or source"
        value={value.actor ?? ''}
        onChange={(e) => handle({ actor: e.target.value || undefined })}
      />

      <TextField
        label="From"
        type="datetime-local"
        size="small"
        InputLabelProps={{ shrink: true }}
        value={value.from ?? ''}
        onChange={(e) => handle({ from: e.target.value || undefined })}
      />

      <TextField
        label="To"
        type="datetime-local"
        size="small"
        InputLabelProps={{ shrink: true }}
        value={value.to ?? ''}
        onChange={(e) => handle({ to: e.target.value || undefined })}
      />

      <FormControl sx={{ minWidth: 140 }}>
        <InputLabel id="sort-by-label">Sort by</InputLabel>
        <Select
          labelId="sort-by-label"
          label="Sort by"
          size="small"
          value={value.sortBy ?? 'created_at'}
          onChange={(e) => handle({ sortBy: (e.target.value as string) || 'created_at' })}
        >
          <MenuItem value="created_at">Created At</MenuItem>
          <MenuItem value="actor">Actor</MenuItem>
          <MenuItem value="entity">Entity</MenuItem>
          <MenuItem value="action">Action</MenuItem>
          <MenuItem value="environment_key">Environment</MenuItem>
          <MenuItem value="username">Username</MenuItem>
        </Select>
      </FormControl>

      <FormControl sx={{ minWidth: 120 }}>
        <InputLabel id="sort-order-label">Order</InputLabel>
        <Select
          labelId="sort-order-label"
          label="Order"
          size="small"
          value={value.sortOrder ?? 'desc'}
          onChange={(e) => handle({ sortOrder: (e.target.value as 'asc' | 'desc') || 'desc' })}
        >
          <MenuItem value="asc">asc</MenuItem>
          <MenuItem value="desc">desc</MenuItem>
        </Select>
      </FormControl>
    </Box>
  );
};

export default AuditLogFilter;
