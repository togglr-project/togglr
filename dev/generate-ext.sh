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

# Add LDAP and Keycloak dependencies to togglr-backend
# Look for the togglr-mailhog dependency and add LDAP and Keycloak after it
mailhog_pattern = r'(togglr-mailhog:\s*\n\s*condition: service_healthy)'
dependencies = r'\1\n      togglr-openldap:\n        condition: service_healthy\n      togglr-keycloak:\n        condition: service_healthy'

content = re.sub(mailhog_pattern, dependencies, content)

# Add volumes to togglr-backend
# Find togglr-backend and add volumes section
backend_pattern = r'(togglr-backend:.*?command: \["/bin/app", "server"\])'
volumes_addition = r'\1\n    volumes:\n      - "./secrets:/opt/togglr/secrets"'

content = re.sub(backend_pattern, volumes_addition, content, flags=re.DOTALL)

# Add LDAP and Keycloak services before volumes section
additional_services = '''
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

  togglr-keycloak-db:
    image: postgres:17
    container_name: togglr-keycloak-db
    restart: always
    environment:
      POSTGRES_USER: keycloak
      POSTGRES_PASSWORD: keycloak123
      POSTGRES_DB: keycloak
      PGDATA: /var/lib/postgresql/data/keycloak
    volumes:
      - "keycloak_db_togglr:/var/lib/postgresql/data/keycloak"
    ports:
      - "5433:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U keycloak -d keycloak"]
      interval: 10s
      timeout: 5s
      retries: 5
      start_period: 30s

  togglr-keycloak:
    container_name: togglr-keycloak.local
    image: quay.io/keycloak/keycloak:26.4
    environment:
      KEYCLOAK_FEATURES: saml
      KEYCLOAK_ADMIN: admin
      KEYCLOAK_ADMIN_PASSWORD: Dev123456
      KC_DB: postgres
      KC_DB_URL: jdbc:postgresql://togglr-keycloak-db:5432/keycloak
      KC_DB_USERNAME: keycloak
      KC_DB_PASSWORD: keycloak123
      KC_HOSTNAME_STRICT: false
      KC_HOSTNAME_STRICT_HTTPS: false
      KC_HTTP_ENABLED: true
      KC_HTTPS_ENABLED: false
      KC_PROXY_HEADERS: xforwarded
      KC_PROXY: edge
      KC_HEALTH_ENABLED: true
      KC_LOG_LEVEL: org.keycloak.saml.SP:debug,org.keycloak.protocol.saml:debug,org.keycloak.events:debug
      KC_HOSTNAME: togglr-keycloak.local
      KC_HOSTNAME_PORT: 8080
      KC_FRONTEND_URL: http://togglr-keycloak.local:9090
      KC_SPI_X_FRAME_OPTIONS_ENABLED: "false"
      KC_HTTP_COOKIE_SAME_SITE: "None"
    ports:
      - "9090:8080"
    extra_hosts:
      - "host.docker.internal:host-gateway"
    depends_on:
      togglr-keycloak-db:
        condition: service_healthy
      togglr-openldap:
        condition: service_healthy
    restart: always
    command: start-dev
    healthcheck:
      test: ["CMD-SHELL", "bash -c 'cat < /dev/null > /dev/tcp/127.0.0.1/8080'"]
      interval: 30s
      timeout: 5s
      retries: 5
      start_period: 60s
'''

# Add additional services before volumes section
volumes_pattern = r'^volumes:'
content = re.sub(volumes_pattern, additional_services + '\nvolumes:', content, flags=re.MULTILINE)

# Add additional volumes to the volumes section at the end
# Find the volumes section and add LDAP and Keycloak volumes after the last volume
volumes_end_pattern = r'(nats_togglr:\s*)$'
additional_volumes = r'\1\n  ldap_data_togglr:\n  ldap_config_togglr:\n  keycloak_db_togglr:'
content = re.sub(volumes_end_pattern, additional_volumes, content, flags=re.MULTILINE)

# Clean up extra empty lines
content = re.sub(r'\n\n\n+', '\n\n', content)

# Write the modified content back
with open('dev/docker-compose.ext.yml', 'w') as f:
    f.write(content)

print("✅ Generated docker-compose.ext.yml successfully!")
print("📝 Added LDAP services: togglr-openldap, togglr-phpldapadmin")
print("📝 Added Keycloak services: togglr-keycloak-db, togglr-keycloak")
print("🔗 Added LDAP and Keycloak dependencies to togglr-backend")
print("📦 Added volumes to togglr-backend: secrets")
print("📦 Added volumes: ldap_data_togglr, ldap_config_togglr, keycloak_db_togglr")
EOF