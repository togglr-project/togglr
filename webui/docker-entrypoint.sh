#!/bin/sh

# Default values
API_BASE_URL=${TOGGLR_API_BASE_URL:-/}
WS_BASE_URL=${TOGGLR_WS_BASE_URL:-}
VERSION=${TOGGLR_VERSION:-dev}
BUILD_TIME=${TOGGLR_BUILD_TIME:-$(date -u +%Y-%m-%dT%H:%M:%SZ)}

# Create runtime configuration
cat > /usr/share/nginx/html/config.js << EOF
// Runtime configuration
window.TOGGLR_CONFIG = {
  API_BASE_URL: '${API_BASE_URL}',
  WS_BASE_URL: '${WS_BASE_URL}',
  VERSION: '${VERSION}',
  BUILD_TIME: '${BUILD_TIME}'
};
EOF

echo "TOGGLR UI configured with:"
echo "  API_BASE_URL: ${API_BASE_URL}"
echo "  WS_BASE_URL: ${WS_BASE_URL}"
echo "  VERSION: ${VERSION}"
echo "  BUILD_TIME: ${BUILD_TIME}"

# Start nginx
exec nginx -g "daemon off;" 