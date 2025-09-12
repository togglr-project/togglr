declare global {
  interface Window {
    ETOGGL_CONFIG: {
      API_BASE_URL: string;
      VERSION: string;
      BUILD_TIME: string;
    };
  }
}

export {}; 