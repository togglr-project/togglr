#!/bin/bash

# Script to generate docker-compose.ext.yml from docker-compose.yml
# This ensures both files stay in sync

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BASE_FILE="$SCRIPT_DIR/docker-compose.yml"
EXT_FILE="$SCRIPT_DIR/docker-compose.ext.yml"

echo "Generating docker-compose.ext.yml from docker-compose.yml..."

# Create the extended version by copying the base file
cp "$BASE_FILE" "$EXT_FILE"

# Add LDAP dependency to togglr-backend using Python for precise formatting
python3 << 'EOF'
import re

# Read the file
with open('dev/docker-compose.ext.yml', 'r') as f:
    content = f.read()

# Add LDAP dependency to togglr-backend
# Look for the togglr-mailhog dependency and add LDAP after it
mailhog_pattern = r'(togglr-mailhog:\s*\n\s*condition: service_healthy)'
ldap_dependency = r'\1\n      togglr-openldap:\n        condition: service_healthy'

content = re.sub(mailhog_pattern, ldap_dependency, content)

# Add LDAP services before volumes section
ldap_services = '''
  togglr-openldap:
    container_name: togglr-openldap
    image: osixia/openldap:1.5.0
    environment:
      LDAP_ORGANISATION: "Togglr"
      LDAP_DOMAIN: "togglr.local"
      LDAP_ADMIN_PASSWORD: "Dev123456"
      LDAP_CONFIG_PASSWORD: "Dev123456"
      LDAP_READONLY_USER: "false"
      LDAP_TLS_VERIFY_CLIENT: "never"
    volumes:
      - ldap_data_togglr:/var/lib/ldap
      - ldap_config_togglr:/etc/ldap/slapd.d
      - ./ldap/ldif:/container/service/slapd/assets/config/bootstrap/ldif/custom
    ports:
      - "389:389"
      - "636:636"
    restart: always
    command: ["--copy-service", "--loglevel", "debug"]
    healthcheck:
      test: ["CMD", "ldapwhoami", "-x", "-H", "ldap://localhost:389"]
      interval: 30s
      timeout: 10s
      retries: 5
      start_period: 30s

  togglr-phpldapadmin:
    container_name: togglr-phpldapadmin
    image: osixia/phpldapadmin:0.9.0
    ports:
      - "9080:80"
    environment:
      PHPLDAPADMIN_LDAP_HOSTS: "togglr-openldap"
      PHPLDAPADMIN_HTTPS: "false"
    depends_on:
      togglr-openldap:
        condition: service_healthy
    restart: always
    healthcheck:
      test: ["CMD-SHELL", "bash -c 'cat < /dev/null > /dev/tcp/127.0.0.1/80'"]
      interval: 30s
      timeout: 5s
      retries: 5
      start_period: 30s
'''

# Add LDAP services before volumes section
volumes_pattern = r'^volumes:'
content = re.sub(volumes_pattern, ldap_services + '\nvolumes:', content, flags=re.MULTILINE)

# Add LDAP volumes to the volumes section at the end
# Find the volumes section and add LDAP volumes after the last volume
# Look for the end of volumes section and add LDAP volumes
volumes_end_pattern = r'(nats_togglr:\s*)$'
ldap_volumes = r'\1\n  ldap_data_togglr:\n  ldap_config_togglr:'
content = re.sub(volumes_end_pattern, ldap_volumes, content, flags=re.MULTILINE)

# Clean up extra empty lines
content = re.sub(r'\n\n\n+', '\n\n', content)

# Write the modified content back
with open('dev/docker-compose.ext.yml', 'w') as f:
    f.write(content)

print("✅ Generated docker-compose.ext.yml successfully!")
print("📝 Added LDAP services: togglr-openldap, togglr-phpldapadmin")
print("🔗 Added LDAP dependency to togglr-backend")
print("📦 Added LDAP volumes: ldap_data_togglr, ldap_config_togglr")
EOF
