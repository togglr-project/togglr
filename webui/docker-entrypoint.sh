#!/bin/sh

# Default values
API_BASE_URL=${ETOGGLE_API_BASE_URL:-/}
VERSION=${ETOGGLE_VERSION:-dev}
BUILD_TIME=${ETOGGLE_BUILD_TIME:-$(date -u +%Y-%m-%dT%H:%M:%SZ)}

# Create runtime configuration
cat > /usr/share/nginx/html/config.js << EOF
// Runtime configuration
window.ETOGGLE_CONFIG = {
  API_BASE_URL: '${API_BASE_URL}',
  VERSION: '${VERSION}',
  BUILD_TIME: '${BUILD_TIME}'
};
EOF

echo "ETOGGLE UI configured with:"
echo "  API_BASE_URL: ${API_BASE_URL}"
echo "  VERSION: ${VERSION}"
echo "  BUILD_TIME: ${BUILD_TIME}"

# Start nginx
exec nginx -g "daemon off;" 