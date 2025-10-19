export interface ProjectSettingDefinition {
  name: string;
  type: 'boolean' | 'integer' | 'double' | 'text' | 'json';
  description: string;
  defaultValue?: any;
  validation?: {
    min?: number;
    max?: number;
    required?: boolean;
  };
}

export const PREDEFINED_PROJECT_SETTINGS: Record<string, ProjectSettingDefinition> = {
  auto_disable_requires_approval: {
    name: 'auto_disable_requires_approval',
    type: 'boolean',
    description: 'Requires approval before auto-disabling features',
    defaultValue: false,
  },
  auto_disable_enabled: {
    name: 'auto_disable_enabled',
    type: 'boolean',
    description: 'Enable automatic feature disabling',
    defaultValue: false,
  },
  auto_disable_error_threshold: {
    name: 'auto_disable_error_threshold',
    type: 'double',
    description: 'Error threshold in percents for auto-disabling features (0.0-100.0)',
    defaultValue: 10.0,
    validation: {
      min: 0.0,
      max: 100.0,
      required: true,
    },
  },
  auto_disable_time_window_sec: {
    name: 'auto_disable_time_window_sec',
    type: 'integer',
    description: 'Time window in seconds for error monitoring',
    defaultValue: 300,
    validation: {
      min: 1,
      max: 86400,
      required: true,
    },
  },
  audit_log_retention_days: {
    name: 'audit_log_retention_days',
    type: 'integer',
    description: 'Number of days to retain audit logs',
    defaultValue: 90,
    validation: {
      min: 1,
      max: 3650,
      required: true,
    },
  },
};

export const getSettingDefinition = (settingName: string): ProjectSettingDefinition | null => {
  return PREDEFINED_PROJECT_SETTINGS[settingName] || null;
};

export const isPredefinedSetting = (settingName: string): boolean => {
  return settingName in PREDEFINED_PROJECT_SETTINGS;
};
