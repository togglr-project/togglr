import React, { useState } from 'react';
import {
  Card,
  CardContent,
  CardActions,
  Typography,
  IconButton,
  Tooltip,
  Box,
  // Chip,
  Collapse,
  Button,
} from '@mui/material';
import {
  Edit as EditIcon,
  // Delete as DeleteIcon,
  ExpandMore as ExpandMoreIcon,
  ExpandLess as ExpandLessIcon,
  // Code as CodeIcon,
  Settings as SettingsIcon,
  // ToggleOn as ToggleOnIcon,
  // Numbers as NumbersIcon,
} from '@mui/icons-material';
import type { ProjectSetting } from '../../generated/api/client';
import { getSettingDefinition } from '../../constants/projectSettings';

interface ProjectSettingCardProps {
  setting: ProjectSetting;
  onEdit: () => void;
  onDelete: () => void;
}

const ProjectSettingCard: React.FC<ProjectSettingCardProps> = ({
  setting,
  onEdit,
  // onDelete,
}) => {
  const [expanded, setExpanded] = useState(false);

  const toggleExpanded = () => {
    setExpanded(!expanded);
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const isValidJson = (value: string) => {
    try {
      JSON.parse(value);
      return true;
    } catch {
      return false;
    }
  };

  const formatJsonValue = (value: string) => {
    if (isValidJson(value)) {
      try {
        return JSON.stringify(JSON.parse(value), null, 2);
      } catch {
        return value;
      }
    }
    return value;
  };

  const getValueType = () => {
    const settingDef = getSettingDefinition(setting.name);
    if (settingDef) {
      return settingDef.type;
    }
    return isValidJson(setting.value) ? 'json' : 'text';
  };

  const formatValue = (value: string, type: string) => {
    switch (type) {
      case 'json':
        return formatJsonValue(value);
      case 'boolean':
        return value === 'true' ? 'True' : 'False';
      case 'integer':
        return parseInt(value, 10).toLocaleString();
      case 'double':
        return parseFloat(value).toLocaleString();
      default:
        return value;
    }
  };

  const valueType = getValueType();
  // const isPredefined = isPredefinedSetting(setting.name);

  return (
    <Card sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
      <CardContent sx={{ flexGrow: 1, p: 2 }}>
        <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', mb: 1 }}>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, flexGrow: 1 }}>
            <SettingsIcon color="primary" sx={{ fontSize: 18 }} />
            <Typography variant="subtitle1" component="h3" sx={{ fontWeight: 500 }}>
              {setting.name}
            </Typography>
          </Box>
          <IconButton
            size="small"
            onClick={toggleExpanded}
            aria-label={expanded ? 'Collapse' : 'Expand'}
            sx={{ p: 0.5 }}
          >
            {expanded ? <ExpandLessIcon fontSize="small" /> : <ExpandMoreIcon fontSize="small" />}
          </IconButton>
        </Box>

        <Typography variant="caption" color="text.secondary" sx={{ mb: 1, display: 'block' }}>
          Updated: {formatDate(setting.updated_at)}
        </Typography>

        {!expanded && (
          <Typography
            variant="body2"
            color="text.secondary"
            sx={{
              overflow: 'hidden',
              textOverflow: 'ellipsis',
              display: '-webkit-box',
              WebkitLineClamp: 1,
              WebkitBoxOrient: 'vertical',
              fontFamily: valueType === 'json' ? 'monospace' : 'inherit',
              fontSize: '0.875rem',
            }}
          >
            {formatValue(setting.value, valueType)}
          </Typography>
        )}

        <Collapse in={expanded}>
          <Box sx={{ mt: 1 }}>
            <Typography variant="subtitle2" gutterBottom sx={{ fontSize: '0.875rem' }}>
              Value:
            </Typography>
            <Box
              sx={{
                bgcolor: 'grey.50',
                border: '1px solid',
                borderColor: 'grey.200',
                borderRadius: 1,
                p: 1.5,
                fontFamily: valueType === 'json' ? 'monospace' : 'inherit',
                fontSize: '0.8rem',
                overflow: 'auto',
                maxHeight: 150,
                whiteSpace: valueType === 'json' ? 'pre-wrap' : 'normal',
                wordBreak: 'break-word',
              }}
            >
              {formatValue(setting.value, valueType)}
            </Box>
          </Box>
        </Collapse>
      </CardContent>

      <CardActions sx={{ justifyContent: 'space-between', px: 1.5, py: 1 }}>
        <Box>
          <Tooltip title="Edit">
            <IconButton
              size="small"
              onClick={onEdit}
              color="primary"
              sx={{ p: 0.5 }}
            >
              <EditIcon fontSize="small" />
            </IconButton>
          </Tooltip>
          {/*<Tooltip title={isPredefined ? "Cannot delete predefined settings" : "Delete"}>*/}
          {/*  <span>*/}
          {/*    <IconButton*/}
          {/*      size="small"*/}
          {/*      onClick={onDelete}*/}
          {/*      color="error"*/}
          {/*      disabled={isPredefined}*/}
          {/*      sx={{ p: 0.5 }}*/}
          {/*    >*/}
          {/*      <DeleteIcon fontSize="small" />*/}
          {/*    </IconButton>*/}
          {/*  </span>*/}
          {/*</Tooltip>*/}
        </Box>
        <Button
          size="small"
          onClick={toggleExpanded}
          endIcon={expanded ? <ExpandLessIcon fontSize="small" /> : <ExpandMoreIcon fontSize="small" />}
          sx={{ fontSize: '0.75rem', minWidth: 'auto', px: 1 }}
        >
          {expanded ? 'Less' : 'More'}
        </Button>
      </CardActions>
    </Card>
  );
};

export default ProjectSettingCard;
