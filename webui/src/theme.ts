import { createTheme } from '@mui/material/styles';
import type { PaletteMode } from '@mui/material';

// Theme type limited to MUI palette modes
export type ThemeMode = PaletteMode;

// Dark theme color constants
const DARK_BACKGROUND_COLOR = 'rgba(28, 30, 35, 0.95)';
const DARK_PAPER_COLOR = 'rgb(38, 40, 48)';
const PRIMARY_COLOR = 'rgb(130, 82, 255)'; // Slightly more blue-purple for a more refined look
const PRIMARY_LIGHT = 'rgba(150, 110, 255, 0.85)';
const SECONDARY_COLOR = '#FF5A8C'; // Soft pink for better contrast with purple
const DARK_TEXT_COLOR = 'rgb(245, 245, 250)'; // Slightly off-white for better eye comfort

// Light theme color constants
const LIGHT_BACKGROUND_COLOR = 'rgb(248, 249, 252)'; // Slightly blue tint for a fresher look
const LIGHT_PAPER_COLOR = 'rgb(255, 255, 255)';
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
      body1: {
        lineHeight: 1.7,
        letterSpacing: '0.005em',
      },
      body2: {
        lineHeight: 1.7,
        letterSpacing: '0.005em',
      },
      button: {
        textTransform: 'none',
        fontWeight: 500,
      },
    },
    shape: {
      borderRadius: 10,
    },
    components: {
      MuiButton: {
        styleOverrides: {
          root: {
            borderRadius: 10,
            padding: '6px 12px',
            fontWeight: 500,
            color: '#ffffff',
            boxShadow: 'none',
            backgroundColor: primaryMain,
            transition: 'background-color 0.15s ease, color 0.15s ease, box-shadow 0.15s ease, transform 0.05s ease',
            '&:hover': {
              backgroundColor: mode === 'light' ? 'rgba(130, 82, 255, 0.9)' : 'rgba(130, 82, 255, 0.85)',
              boxShadow: '0 2px 8px rgba(0, 0, 0, 0.12)'
            },
            '&:active': {
              backgroundColor: mode === 'light' ? 'rgba(130, 82, 255, 0.85)' : 'rgba(130, 82, 255, 0.8)',
              transform: 'translateY(0)'
            },
            '&:focus-visible': {
              outline: 'none',
              boxShadow: '0 0 0 3px rgba(130, 82, 255, 0.28)'
            },
          },
          outlined: {
            borderWidth: 1,
            '&:hover': {
              borderWidth: 1,
            },
          },
        },
      },
      MuiPaper: {
        styleOverrides: {
          root: {
            borderRadius: 10,
            backgroundColor: paperColor,
            boxShadow: '0 1px 3px rgba(0, 0, 0, 0.04), 0 8px 24px rgba(0, 0, 0, 0.04)',
            border: `1px solid ${mode === 'dark' ? 'rgba(255, 255, 255, 0.07)' : 'rgba(0, 0, 0, 0.06)'}`,
          },
        },
      },
      MuiAppBar: {
        styleOverrides: {
          root: {
            backgroundColor: mode === 'dark' ? 'rgba(28, 30, 35, 0.75)' : 'rgba(255, 255, 255, 0.8)',
            boxShadow: '0 1px 1px rgba(0, 0, 0, 0.06)',
            borderBottom: `1px solid ${mode === 'dark' ? 'rgba(255, 255, 255, 0.06)' : 'rgba(0, 0, 0, 0.06)'}`,
            backdropFilter: 'saturate(180%) blur(8px)',
            WebkitBackdropFilter: 'saturate(180%) blur(8px)',
            overflow: 'hidden',
            backgroundClip: 'padding-box',
            WebkitMaskImage: '-webkit-radial-gradient(white, black)'
          },
        },
      },
      MuiCard: {
        styleOverrides: {
          root: {
            borderRadius: 10,
            backgroundColor: paperColor,
            boxShadow: '0 1px 3px rgba(0, 0, 0, 0.04), 0 8px 24px rgba(0, 0, 0, 0.04)',
            border: `1px solid ${mode === 'dark' ? 'rgba(255, 255, 255, 0.07)' : 'rgba(0, 0, 0, 0.06)'}`,
            transition: 'box-shadow 0.2s ease-in-out, transform 0.2s ease-in-out',
            '&:hover': {
              transform: 'translateY(-1px)',
              boxShadow: '0 4px 16px rgba(0, 0, 0, 0.10)'
            },
          },
        },
      },
      MuiDialog: {
        styleOverrides: {
          paper: {
            backgroundColor: mode === 'dark' ? 'rgba(16, 18, 22, 0.98)' : paperColor,
            border: `1px solid ${mode === 'dark' ? 'rgba(255, 255, 255, 0.06)' : 'rgba(0, 0, 0, 0.06)'}`,
            boxShadow: mode === 'dark' ? '0 10px 24px rgba(0, 0, 0, 0.32), 0 2px 8px rgba(0, 0, 0, 0.18)' : '0 6px 16px rgba(0, 0, 0, 0.08)',
          },
        },
      },
      MuiMenu: {
        styleOverrides: {
          paper: {
            backgroundColor: mode === 'dark' ? 'rgb(16, 18, 22)' : 'rgb(255, 255, 255)',
            borderRadius: 10,
            border: `1px solid ${mode === 'dark' ? 'rgba(255, 255, 255, 0.06)' : 'rgba(0, 0, 0, 0.06)'}`,
            boxShadow: mode === 'dark'
              ? '0 8px 20px rgba(0, 0, 0, 0.28), 0 2px 6px rgba(0, 0, 0, 0.16)'
              : '0 6px 16px rgba(0, 0, 0, 0.08)'
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
              borderRadius: 10,
              transition: 'box-shadow 0.15s ease, border-color 0.15s ease',
              '&.Mui-focused .MuiOutlinedInput-notchedOutline': {
                borderColor: primaryMain,
              },
              '&.Mui-focused': {
                boxShadow: '0 0 0 3px rgba(130, 82, 255, 0.20)'
              },
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
            borderRadius: 10,
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
              backgroundColor: mode === 'dark' ? 'rgba(130, 82, 255, 0.08)' : 'rgba(130, 82, 255, 0.04)',
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
              backgroundColor: mode === 'dark' ? 'rgba(130, 82, 255, 0.08)' : 'rgba(130, 82, 255, 0.04)',
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
            '&.Mui-focused': {
              boxShadow: '0 0 0 3px rgba(130, 82, 255, 0.20)'
            },
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
