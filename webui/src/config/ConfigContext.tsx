import React, { createContext, useContext, type ReactNode } from 'react';

interface ConfigContextType {
  isDemo: boolean;
}

// Create the context with a default value
const ConfigContext = createContext<ConfigContextType>({
  isDemo: false,
});

// Custom hook to use the config context
export const useConfig = (): ConfigContextType => {
  return useContext(ConfigContext);
};

interface ConfigProviderProps {
  children: ReactNode;
}

export const ConfigProvider: React.FC<ConfigProviderProps> = ({ children }) => {
  // Read the IS_DEMO environment variable
  // In Vite, environment variables are prefixed with VITE_
  const isDemo = import.meta.env.VITE_IS_DEMO === 'true';

  return (
    <ConfigContext.Provider value={{ isDemo }}>
      {children}
    </ConfigContext.Provider>
  );
};