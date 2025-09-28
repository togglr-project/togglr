import { createTheme } from '@mui/material/styles';
import type { PaletteMode } from '@mui/material';

// Theme type limited to MUI palette modes
export type ThemeMode = PaletteMode;

// Dark theme color constants
const DARK_BACKGROUND_COLOR = 'rgba(28, 30, 35, 0.95)';
const DARK_PAPER_COLOR = 'rgba(38, 40, 48, 0.9)';
const PRIMARY_COLOR = 'rgb(130, 82, 255)'; // Slightly more blue-purple for a more refined look
const PRIMARY_LIGHT = 'rgba(150, 110, 255, 0.85)';
const SECONDARY_COLOR = '#FF5A8C'; // Soft pink for better contrast with purple
const DARK_TEXT_COLOR = 'rgb(245, 245, 250)'; // Slightly off-white for better eye comfort

// Light theme color constants
const LIGHT_BACKGROUND_COLOR = 'rgba(248, 249, 252, 0.95)'; // Slightly blue tint for a fresher look
const LIGHT_PAPER_COLOR = 'rgba(255, 255, 255, 0.92)';
const LIGHT_TEXT_COLOR = 'rgb(40, 42, 54)'; // Slightly blue-black for better contrast


// Create theme based on mode
export const getTheme = (mode: ThemeMode) => {
  // Determine the actual palette mode for MUI
  const actualMode: PaletteMode = (mode === 'light') ? 'light' : 'dark';

  // Get primary and secondary colors based on theme
  let primaryMain = PRIMARY_COLOR;
  let primaryLight = PRIMARY_LIGHT;
  let secondaryMain = SECONDARY_COLOR;
  let backgroundColor = mode === 'dark' ? DARK_BACKGROUND_COLOR : LIGHT_BACKGROUND_COLOR;
  let paperColor = mode === 'dark' ? DARK_PAPER_COLOR : LIGHT_PAPER_COLOR;
  let textPrimary = mode === 'dark' ? DARK_TEXT_COLOR : LIGHT_TEXT_COLOR;
  let textSecondary = mode === 'dark' ? 'rgba(220, 220, 230, 0.7)' : 'rgba(60, 65, 75, 0.75)';


  return createTheme({
    palette: {
      mode: actualMode,
      primary: {
        main: primaryMain,
        light: primaryLight,
      },
      secondary: {
        main: secondaryMain,
      },
      background: {
        default: backgroundColor,
        paper: paperColor,
      },
      text: {
        primary: textPrimary,
        secondary: textSecondary,
      },
      error: {
        main: mode === 'light' ? '#E53935' : '#FF5A5A',
      },
      warning: {
        main: mode === 'light' ? '#F9A825' : '#FFAA5A',
      },
      info: {
        main: mode === 'light' ? '#2196F3' : '#5AC8FF',
      },
      success: {
        main: mode === 'light' ? '#4CAF50' : '#5AFF8F',
      },
    },
    typography: {
      fontFamily: '"Inter", "Roboto", "Helvetica", "Arial", sans-serif',
      h1: {
        fontWeight: 700,
        letterSpacing: '-0.01em',
      },
      h2: {
        fontWeight: 700,
        letterSpacing: '-0.01em',
      },
      h3: {
        fontWeight: 600,
        letterSpacing: '-0.01em',
      },
      h4: {
        fontWeight: 600,
        letterSpacing: '-0.01em',
      },
      h5: {
        fontWeight: 600,
      },
      h6: {
        fontWeight: 600,
      },
      subtitle1: {
        fontWeight: 500,
      },
      subtitle2: {
        fontWeight: 500,
      },
      button: {
        textTransform: 'none',
        fontWeight: 500,
      },
    },
    shape: {
      borderRadius: 8,
    },
    components: {
      MuiButton: {
        styleOverrides: {
          root: {
            borderRadius: 8,
            padding: '6px 12px',
            fontWeight: 500,
            color: '#ffffff',
            boxShadow: '0 2px 8px 0 rgba(0, 0, 0, 0.1)',
            background: 'linear-gradient(45deg, rgba(130, 82, 255, 0.9) 30%, rgba(150, 110, 255, 0.9) 90%)',
            transition: 'all 0.2s ease-in-out',
            '&:hover': {
              background: 'linear-gradient(45deg, rgba(130, 82, 255, 1) 30%, rgba(150, 110, 255, 1) 90%)',
              boxShadow: '0 4px 12px 0 rgba(0, 0, 0, 0.2)',
              transform: 'translateY(-1px)',
            },
          },
          outlined: {
            borderWidth: 2,
            '&:hover': {
              borderWidth: 2,
            },
          },
        },
      },
      MuiPaper: {
        styleOverrides: {
          root: {
            borderRadius: 8,
            backgroundImage: mode === 'dark'
              ? 'linear-gradient(to bottom, rgba(45, 48, 56, 0.7), rgba(35, 38, 46, 0.85))'
              : 'linear-gradient(to bottom, rgba(255, 255, 255, 0.9), rgba(248, 249, 252, 0.85))',
            backdropFilter: 'blur(12px)',
            boxShadow: mode !== 'light'
              ? '0 8px 24px 0 rgba(0, 0, 0, 0.2)'
              : '0 8px 24px 0 rgba(0, 0, 0, 0.08)',
          },
        },
      },
      MuiAppBar: {
        styleOverrides: {
          root: {
            background: mode === 'dark'
              ? 'linear-gradient(90deg, rgba(28, 30, 35, 0.95) 0%, rgba(35, 38, 46, 0.95) 100%)'
              : 'linear-gradient(90deg, rgba(248, 249, 252, 0.95) 0%, rgba(255, 255, 255, 0.95) 100%)',
            backdropFilter: 'blur(10px)',
            boxShadow: mode !== 'light'
              ? '0 2px 12px 0 rgba(0, 0, 0, 0.2)'
              : '0 2px 12px 0 rgba(0, 0, 0, 0.06)',
          },
        },
      },
      MuiCard: {
        styleOverrides: {
          root: {
            borderRadius: 8,
            background: mode === 'dark'
              ? 'linear-gradient(135deg, rgba(45, 48, 56, 0.85) 0%, rgba(35, 38, 46, 0.85) 100%)'
              : 'linear-gradient(135deg, rgba(255, 255, 255, 0.9) 0%, rgba(248, 249, 252, 0.9) 100%)',
            backdropFilter: 'blur(10px)',
            boxShadow: mode !== 'light'
              ? '0 4px 20px 0 rgba(0, 0, 0, 0.2)'
              : '0 4px 20px 0 rgba(0, 0, 0, 0.05)',
            transition: 'all 0.3s ease-in-out',
            '&:hover': {
              boxShadow: mode !== 'light'
                ? '0 8px 30px 0 rgba(0, 0, 0, 0.3)'
                : '0 8px 30px 0 rgba(0, 0, 0, 0.1)',
            },
          },
        },
      },
      MuiChip: {
        styleOverrides: {
          root: {
            borderRadius: 6,
            fontWeight: 500,
          },
          outlined: {
            borderWidth: 1.5,
          },
        },
      },
      MuiTextField: {
        styleOverrides: {
          root: {
            '& .MuiOutlinedInput-root': {
              borderRadius: 8,
              // Override browser autofill styling
              '& input:-webkit-autofill': {
                WebkitBoxShadow: `0 0 0 100px ${paperColor} inset`,
                WebkitTextFillColor: textPrimary,
                caretColor: textPrimary,
                borderRadius: 'inherit',
                '&:focus': {
                  WebkitBoxShadow: `0 0 0 100px ${paperColor} inset`,
                },
              },
              // Firefox autofill
              '& input:autofill': {
                background: paperColor,
                color: textPrimary,
              },
            },
          },
        },
        variants: [
          {
            props: { size: 'small' },
            style: {
              '& .MuiOutlinedInput-root': {
                height: 36,
                '& .MuiInputBase-input': {
                  padding: '8px 12px',
                  fontSize: '0.875rem',
                },
              },
              '& .MuiInputLabel-root': {
                fontSize: '0.875rem',
                transform: 'translate(14px, 9px) scale(1)',
                '&.MuiInputLabel-shrink': {
                  transform: 'translate(14px, -9px) scale(0.75)',
                },
              },
            },
          },
        ],
      },
      MuiSelect: {
        styleOverrides: {
          outlined: {
            borderRadius: 8,
          },
        },
        variants: [
          {
            props: { size: 'small' },
            style: {
              height: 36,
              '& .MuiSelect-select': {
                padding: '8px 12px',
                fontSize: '0.875rem',
              },
            },
          },
        ],
      },
      MuiMenuItem: {
        styleOverrides: {
          root: {
            borderRadius: 4,
            margin: '1px 2px',
            '&:hover': {
              backgroundColor: mode === 'dark' ? 'rgba(130, 82, 255, 0.1)' : 'rgba(130, 82, 255, 0.05)',
            },
          },
        },
      },
      MuiListItemButton: {
        styleOverrides: {
          root: {
            borderRadius: 6,
            margin: '1px 2px',
            '&:hover': {
              backgroundColor: mode === 'dark' ? 'rgba(130, 82, 255, 0.1)' : 'rgba(130, 82, 255, 0.05)',
            },
          },
        },
      },
      MuiDivider: {
        styleOverrides: {
          root: {
            borderColor: mode === 'dark' ? 'rgba(255, 255, 255, 0.1)' : 'rgba(0, 0, 0, 0.08)',
          },
        },
      },
      MuiInput: {
        styleOverrides: {
          root: {
            // Override browser autofill styling
            '& input:-webkit-autofill': {
              WebkitBoxShadow: `0 0 0 100px ${paperColor} inset`,
              WebkitTextFillColor: textPrimary,
              caretColor: textPrimary,
              '&:focus': {
                WebkitBoxShadow: `0 0 0 100px ${paperColor} inset`,
              },
            },
            // Firefox autofill
            '& input:autofill': {
              background: paperColor,
              color: textPrimary,
            },
          },
        },
      },
      MuiFormControl: {
        variants: [
          {
            props: { size: 'small' },
            style: {
              '& .MuiInputLabel-root': {
                fontSize: '0.875rem',
                transform: 'translate(14px, 9px) scale(1)',
                '&.MuiInputLabel-shrink': {
                  transform: 'translate(14px, -9px) scale(0.75)',
                },
              },
            },
          },
        ],
      },
      MuiFilledInput: {
        styleOverrides: {
          root: {
            // Override browser autofill styling
            '& input:-webkit-autofill': {
              WebkitBoxShadow: `0 0 0 100px ${mode === 'light' ? 'rgba(0, 0, 0, 0.06)' : 'rgba(255, 255, 255, 0.09)'} inset`,
              WebkitTextFillColor: textPrimary,
              caretColor: textPrimary,
              '&:focus': {
                WebkitBoxShadow: `0 0 0 100px ${mode === 'light' ? 'rgba(0, 0, 0, 0.06)' : 'rgba(255, 255, 255, 0.09)'} inset`,
              },
            },
            // Firefox autofill
            '& input:autofill': {
              background: mode === 'light' ? 'rgba(0, 0, 0, 0.06)' : 'rgba(255, 255, 255, 0.09)',
              color: textPrimary,
            },
          },
        },
      },
    },
  });
};

// Default theme is dark
const theme = getTheme('dark');

export default theme;
