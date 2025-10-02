import React from 'react';
import { Box, Paper, Typography } from '@mui/material';

interface JsonDiffViewerProps {
  oldValue: any;
  newValue: any;
}

interface DiffLine {
  type: 'added' | 'removed' | 'unchanged';
  content: string;
  lineNumber?: number;
  key?: string;
}

const JsonDiffViewer: React.FC<JsonDiffViewerProps> = ({ oldValue, newValue }) => {
  // Sort object keys for consistent comparison
  const sortObjectKeys = (obj: any): any => {
    if (obj === null || typeof obj !== 'object') {
      return obj;
    }
    
    if (Array.isArray(obj)) {
      return obj.map(sortObjectKeys);
    }
    
    const sortedKeys = Object.keys(obj).sort();
    const sortedObj: any = {};
    
    for (const key of sortedKeys) {
      sortedObj[key] = sortObjectKeys(obj[key]);
    }
    
    return sortedObj;
  };

  const formatJson = (obj: any): string => {
    const sortedObj = sortObjectKeys(obj);
    return JSON.stringify(sortedObj, null, 2);
  };

  // Create a proper JSON diff that compares by keys
  const createJsonDiff = (oldObj: any, newObj: any): DiffLine[] => {
    const diff: DiffLine[] = [];
    
    // Get all unique keys from both objects
    const allKeys = new Set([
      ...Object.keys(oldObj || {}),
      ...Object.keys(newObj || {})
    ]);
    
    const sortedKeys = Array.from(allKeys).sort();
    
    // Add opening brace
    diff.push({ type: 'unchanged', content: '{' });
    
    for (let i = 0; i < sortedKeys.length; i++) {
      const key = sortedKeys[i];
      const oldValue = oldObj?.[key];
      const newValue = newObj?.[key];
      const isLast = i === sortedKeys.length - 1;
      
      const oldExists = oldValue !== undefined;
      const newExists = newValue !== undefined;
      
      if (!oldExists && newExists) {
        // Key was added
        const valueStr = JSON.stringify(newValue, null, 2).split('\n').map((line, idx) => 
          idx === 0 ? line : '  ' + line
        ).join('\n');
        diff.push({ 
          type: 'added', 
          content: `  "${key}": ${valueStr}${isLast ? '' : ','}`,
          key 
        });
      } else if (oldExists && !newExists) {
        // Key was removed
        const valueStr = JSON.stringify(oldValue, null, 2).split('\n').map((line, idx) => 
          idx === 0 ? line : '  ' + line
        ).join('\n');
        diff.push({ 
          type: 'removed', 
          content: `  "${key}": ${valueStr}${isLast ? '' : ','}`,
          key 
        });
      } else if (oldExists && newExists) {
        // Key exists in both - check if values are different
        const oldStr = JSON.stringify(oldValue);
        const newStr = JSON.stringify(newValue);
        
        if (oldStr === newStr) {
          // Values are identical
          const valueStr = JSON.stringify(oldValue, null, 2).split('\n').map((line, idx) => 
            idx === 0 ? line : '  ' + line
          ).join('\n');
          diff.push({ 
            type: 'unchanged', 
            content: `  "${key}": ${valueStr}${isLast ? '' : ','}`,
            key 
          });
        } else {
          // Values are different - show both
          const oldValueStr = JSON.stringify(oldValue, null, 2).split('\n').map((line, idx) => 
            idx === 0 ? line : '  ' + line
          ).join('\n');
          const newValueStr = JSON.stringify(newValue, null, 2).split('\n').map((line, idx) => 
            idx === 0 ? line : '  ' + line
          ).join('\n');
          
          diff.push({ 
            type: 'removed', 
            content: `  "${key}": ${oldValueStr}${isLast ? '' : ','}`,
            key 
          });
          diff.push({ 
            type: 'added', 
            content: `  "${key}": ${newValueStr}${isLast ? '' : ','}`,
            key 
          });
        }
      }
    }
    
    // Add closing brace
    diff.push({ type: 'unchanged', content: '}' });
    
    return diff;
  };

  const diff = createJsonDiff(oldValue, newValue);

  const renderLine = (line: DiffLine, index: number) => {
    const getLineStyle = () => {
      switch (line.type) {
        case 'added':
          return {
            bgcolor: 'success.light',
            color: 'success.contrastText',
            borderLeft: '3px solid',
            borderLeftColor: 'success.main',
            pl: 1
          };
        case 'removed':
          return {
            bgcolor: 'error.light',
            color: 'error.contrastText',
            borderLeft: '3px solid',
            borderLeftColor: 'error.main',
            pl: 1
          };
        default:
          return {
            bgcolor: 'background.paper',
            color: 'text.primary',
            pl: 1
          };
      }
    };

    return (
      <Box
        key={index}
        sx={{
          display: 'flex',
          alignItems: 'flex-start',
          minHeight: '20px',
          fontFamily: 'monospace',
          fontSize: '0.75rem',
          ...getLineStyle()
        }}
      >
        <Box sx={{ minWidth: '20px', textAlign: 'center', opacity: 0.8, mr: 1 }}>
          {line.type === 'added' ? '+' : line.type === 'removed' ? '-' : ' '}
        </Box>
        <Box sx={{ flex: 1, whiteSpace: 'pre-wrap' }}>
          {line.content || ' '}
        </Box>
      </Box>
    );
  };

  return (
    <Box sx={{ 
      border: 1, 
      borderColor: 'divider', 
      borderRadius: 1,
      overflow: 'hidden',
      maxHeight: '400px',
      overflowY: 'auto'
    }}>
      <Box sx={{ 
        bgcolor: 'background.paper', 
        p: 1, 
        borderBottom: 1, 
        borderColor: 'divider',
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center'
      }}>
        <Typography variant="caption" sx={{ fontWeight: 'bold', color: 'text.primary' }}>
          JSON Diff
        </Typography>
        <Box sx={{ display: 'flex', gap: 1 }}>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
            <Box sx={{ width: 12, height: 12, bgcolor: 'error.main', borderRadius: 0.5 }} />
            <Typography variant="caption" sx={{ color: 'text.primary' }}>Removed</Typography>
          </Box>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
            <Box sx={{ width: 12, height: 12, bgcolor: 'success.main', borderRadius: 0.5 }} />
            <Typography variant="caption" sx={{ color: 'text.primary' }}>Added</Typography>
          </Box>
        </Box>
      </Box>
      <Box>
        {diff.map((line, index) => renderLine(line, index))}
      </Box>
    </Box>
  );
};

export default JsonDiffViewer;
