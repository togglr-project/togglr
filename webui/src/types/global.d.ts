declare global {
  interface Window {
    TOGGLR_CONFIG: {
      API_BASE_URL: string;
      VERSION: string;
      BUILD_TIME: string;
    };
  }
}

export {}; 