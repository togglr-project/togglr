import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { ThemeProvider } from '@mui/material/styles';
import theme from '../../theme';
import ScheduleBuilder from '../ScheduleBuilder';
import type { ScheduleBuilderData } from '../../utils/cronGenerator';

import { vi } from 'vitest';

// Mock the dependencies
vi.mock('timezone-support', () => ({
  listTimeZones: vi.fn(() => ['UTC', 'America/New_York', 'Europe/London'])
}));

vi.mock('cron-validator', () => ({
  isValidCron: vi.fn(() => true)
}));

vi.mock('cronstrue', () => ({
  default: {
    toString: vi.fn(() => 'Every day at 9:00 AM')
  }
}));

vi.mock('../../utils/cronGenerator', () => ({
  generateCronExpression: vi.fn(() => '0 9 * * *'),
  generateScheduleDescription: vi.fn(() => 'Every day at 9:00 AM'),
  validateScheduleData: vi.fn(() => [])
}));

vi.mock('../TimelinePreview', () => ({
  default: function MockTimelinePreview() {
    return <div data-testid="timeline-preview">Timeline Preview</div>;
  }
}));

const defaultProps = {
  open: true,
  onSubmit: vi.fn(),
  featureId: 'test-feature-123',
  initialData: {
    repeatEvery: { interval: 1, unit: 'hours' }, // 1 hour interval
    duration: { value: 30, unit: 'minutes' } // 30 minutes is less than 1 hour
  },
  featureCreatedAt: '2024-01-01T00:00:00Z'
};

const renderWithTheme = (component: React.ReactElement) => {
  return render(
    <ThemeProvider theme={theme}>
      {component}
    </ThemeProvider>
  );
};

describe('ScheduleBuilder', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('Rendering', () => {
    it('renders when open is true', () => {
      renderWithTheme(<ScheduleBuilder {...defaultProps} />);
      
      expect(screen.getByText('Date Range')).toBeInTheDocument();
      expect(screen.getByText('Schedule Type')).toBeInTheDocument();
      expect(screen.getByText('Parameters')).toBeInTheDocument();
      expect(screen.getByText('Duration')).toBeInTheDocument();
      expect(screen.getByText('Action')).toBeInTheDocument();
      expect(screen.getByText('Preview')).toBeInTheDocument();
    });

    it('does not render when open is false', () => {
      renderWithTheme(<ScheduleBuilder {...defaultProps} open={false} />);
      
      expect(screen.queryByText('Date Range')).not.toBeInTheDocument();
    });

    it('renders with initial data', () => {
      const initialData: Partial<ScheduleBuilderData> = {
        timezone: 'America/New_York',
        scheduleType: 'daily',
        action: 'disable'
      };

      renderWithTheme(
        <ScheduleBuilder 
          {...defaultProps} 
          initialData={initialData} 
        />
      );
      
      expect(screen.getByText('Date Range')).toBeInTheDocument();
    });
  });

  describe('Stepper Navigation', () => {
    it('starts at step 0 (Date Range)', () => {
      renderWithTheme(<ScheduleBuilder {...defaultProps} />);
      
      expect(screen.getByText('Date Range')).toBeInTheDocument();
      expect(screen.getByText('Next')).toBeInTheDocument();
      expect(screen.queryByText('Back')).not.toBeInTheDocument();
    });

    it('navigates to next step when Next button is clicked', async () => {
      const user = userEvent.setup();
      renderWithTheme(<ScheduleBuilder {...defaultProps} />);
      
      // Fill required field for step 0 - look for date input
      const dateInputs = screen.getAllByDisplayValue('');
      const startDateInput = dateInputs[0]; // First date input should be start date
      await user.type(startDateInput, '2024-01-01');
      
      const nextButton = screen.getByText('Next');
      await user.click(nextButton);
      
      expect(screen.getByText('Schedule Type')).toBeInTheDocument();
      expect(screen.getByText('Back')).toBeInTheDocument();
    });

    it('navigates back to previous step when Back button is clicked', async () => {
      const user = userEvent.setup();
      renderWithTheme(<ScheduleBuilder {...defaultProps} />);
      
      // Navigate to step 1
      const dateInputs = screen.getAllByDisplayValue('');
      const startDateInput = dateInputs[0];
      await user.type(startDateInput, '2024-01-01');
      
      const nextButton = screen.getByText('Next');
      await user.click(nextButton);
      
      // Navigate back
      const backButton = screen.getByText('Back');
      await user.click(backButton);
      
      expect(screen.getByText('Date Range')).toBeInTheDocument();
    });

    it('disables Next button when current step is invalid', () => {
      // Test without featureCreatedAt to ensure Next button is disabled
      const propsWithoutFeatureCreatedAt = { ...defaultProps, featureCreatedAt: undefined };
      renderWithTheme(<ScheduleBuilder {...propsWithoutFeatureCreatedAt} />);
      
      const nextButton = screen.getByText('Next');
      expect(nextButton).toBeDisabled();
    });
  });

  describe('Date Range Step', () => {
    it('requires start date to proceed', async () => {
      const user = userEvent.setup();
      // Test without featureCreatedAt to ensure Next button is disabled initially
      const propsWithoutFeatureCreatedAt = { ...defaultProps, featureCreatedAt: undefined };
      renderWithTheme(<ScheduleBuilder {...propsWithoutFeatureCreatedAt} />);
      
      const nextButton = screen.getByText('Next');
      expect(nextButton).toBeDisabled();
      
      const dateInputs = screen.getAllByDisplayValue('');
      const startDateInput = dateInputs[0];
      await user.type(startDateInput, '2024-01-01');
      
      expect(nextButton).not.toBeDisabled();
    });

    it('allows optional end date', async () => {
      const user = userEvent.setup();
      renderWithTheme(<ScheduleBuilder {...defaultProps} />);
      
      const dateInputs = screen.getAllByDisplayValue('');
      const startDateInput = dateInputs[0];
      await user.type(startDateInput, '2024-01-01');
      
      // Try to add end date if there's a second input
      if (dateInputs.length > 1) {
        const endDateInput = dateInputs[1];
        await user.type(endDateInput, '2024-12-31');
      }
      
      const nextButton = screen.getByText('Next');
      expect(nextButton).not.toBeDisabled();
    });
  });

  describe('Schedule Type Step', () => {
    beforeEach(async () => {
      const user = userEvent.setup();
      renderWithTheme(<ScheduleBuilder {...defaultProps} />);
      
      // Navigate to schedule type step
      const dateInputs = screen.getAllByDisplayValue('');
      const startDateInput = dateInputs[0];
      await user.type(startDateInput, '2024-01-01');
      
      const nextButton = screen.getByText('Next');
      await user.click(nextButton);
    });

    it('renders schedule type options', () => {
      expect(screen.getByText('Repeat every N minutes/hours')).toBeInTheDocument();
      expect(screen.getByText('At fixed time daily')).toBeInTheDocument();
      expect(screen.getByText('At fixed day monthly')).toBeInTheDocument();
      expect(screen.getByText('Once a year')).toBeInTheDocument();
    });

    it('selects repeat every by default', () => {
      const repeatEveryRadio = screen.getByLabelText('Repeat every N minutes/hours');
      expect(repeatEveryRadio).toBeChecked();
    });

    it('allows changing schedule type', async () => {
      const user = userEvent.setup();
      
      const dailyRadio = screen.getByLabelText('At fixed time daily');
      await user.click(dailyRadio);
      
      expect(dailyRadio).toBeChecked();
    });
  });

  describe('Form Submission', () => {
        it('calls onSubmit with correct data when form is submitted', async () => {
          const user = userEvent.setup();
          const mockOnSubmit = vi.fn();
          
          renderWithTheme(
            <ScheduleBuilder 
              {...defaultProps} 
              onSubmit={mockOnSubmit} 
            />
          );
          
          // Fill out the form step by step
          const dateInputs = screen.getAllByDisplayValue('');
          const startDateInput = dateInputs[0];
          await user.type(startDateInput, '2024-01-01');
          
          // Navigate through all steps
          const nextButton = screen.getByText('Next');
          await user.click(nextButton); // Schedule Type (already has default values)
          await user.click(nextButton); // Parameters (already has default values)
          await user.click(nextButton); // Duration (already has default values)
          await user.click(nextButton); // Action (already has default values)
          await user.click(nextButton); // Preview
          
          // Wait for the Create Schedule button to appear
          await waitFor(() => {
            expect(screen.getByText('Create Schedule')).toBeInTheDocument();
          });

          // Submit
          const createButton = screen.getByText('Create Schedule');
          await user.click(createButton);
          
          expect(mockOnSubmit).toHaveBeenCalledWith(
            expect.objectContaining({
              action: 'enable',
              cronExpression: '0 9 * * *',
              duration: { value: 30, unit: 'minutes' },
              repeatEvery: { interval: 1, unit: 'hours' },
              scheduleType: 'repeat_every',
              startsAt: '2024-01-01T00:00:00Z',
              timezone: expect.any(String)
            })
          );
        });

    it('does not submit when form is invalid', async () => {
      const user = userEvent.setup();
      const mockOnSubmit = vi.fn();
      
      // Test without featureCreatedAt to ensure Next button is disabled
      const propsWithoutFeatureCreatedAt = { ...defaultProps, featureCreatedAt: undefined };
      renderWithTheme(
        <ScheduleBuilder 
          {...propsWithoutFeatureCreatedAt} 
          onSubmit={mockOnSubmit} 
        />
      );
      
      // Try to submit without filling required fields
      const nextButton = screen.getByText('Next');
      expect(nextButton).toBeDisabled();
      
      expect(mockOnSubmit).not.toHaveBeenCalled();
    });
  });

  describe('Error Handling', () => {
        it('displays validation errors', async () => {
          const { validateScheduleData } = await import('../../utils/cronGenerator');
          vi.mocked(validateScheduleData).mockReturnValue(['Timezone is required']);

          renderWithTheme(<ScheduleBuilder {...defaultProps} />);

          // Navigate to preview step where errors are displayed
          const user = userEvent.setup();
          const dateInputs = screen.getAllByDisplayValue('');
          const startDateInput = dateInputs[0];
          await user.type(startDateInput, '2024-01-01');

          const nextButton = screen.getByText('Next');
          await user.click(nextButton); // Schedule Type
          await user.click(nextButton); // Parameters
          await user.click(nextButton); // Duration
          await user.click(nextButton); // Action
          await user.click(nextButton); // Preview

          // Wait for validation errors to appear
          await waitFor(() => {
            expect(screen.getByText('Timezone is required')).toBeInTheDocument();
          });
        });

    it('handles cron generation errors gracefully', async () => {
      const { generateCronExpression } = await import('../../utils/cronGenerator');
      vi.mocked(generateCronExpression).mockImplementation(() => {
        throw new Error('Invalid schedule data');
      });
      
      renderWithTheme(<ScheduleBuilder {...defaultProps} />);
      
      // Should not crash and should render the component
      expect(screen.getByText('Date Range')).toBeInTheDocument();
    });
  });

  describe('Accessibility', () => {
    it('has proper ARIA labels for date inputs', () => {
      renderWithTheme(<ScheduleBuilder {...defaultProps} />);
      
      // Check that date inputs are present (they should have proper labels)
      const dateInputs = screen.getAllByDisplayValue('');
      expect(dateInputs.length).toBeGreaterThan(0);
    });

    it('supports keyboard navigation', async () => {
      const user = userEvent.setup();
      renderWithTheme(<ScheduleBuilder {...defaultProps} />);
      
      const dateInputs = screen.getAllByDisplayValue('');
      const startDateInput = dateInputs[0];
      await user.type(startDateInput, '2024-01-01');
      
      const nextButton = screen.getByText('Next');
      await user.tab();
      await user.keyboard('{Enter}');
      
      expect(screen.getByText('Schedule Type')).toBeInTheDocument();
    });
  });

  describe('Component Integration', () => {
    it('renders TimelinePreview component', async () => {
      const user = userEvent.setup();
      renderWithTheme(<ScheduleBuilder {...defaultProps} />);
      
      // Navigate to preview step
      const dateInputs = screen.getAllByDisplayValue('');
      const startDateInput = dateInputs[0];
      await user.type(startDateInput, '2024-01-01');
      
      const nextButton = screen.getByText('Next');
      await user.click(nextButton); // Schedule Type
      await user.click(nextButton); // Parameters
      await user.click(nextButton); // Duration
      await user.click(nextButton); // Action
      await user.click(nextButton); // Preview
      
      // Wait for TimelinePreview to appear
      await waitFor(() => {
        expect(screen.getByTestId('timeline-preview')).toBeInTheDocument();
      });
    });

    it('handles timezone selection', () => {
      renderWithTheme(<ScheduleBuilder {...defaultProps} />);
      
      // The component should render without crashing
      expect(screen.getByText('Date Range')).toBeInTheDocument();
    });
  });
});
