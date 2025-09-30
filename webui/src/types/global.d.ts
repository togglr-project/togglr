declare global {
  interface Window {
    TOGGLR_CONFIG: {
      API_BASE_URL: string;
      WS_BASE_URL?: string;
      VERSION: string;
      BUILD_TIME: string;
    };
    __RQ?: import('@tanstack/react-query').QueryClient;
  }
}

export {}; 