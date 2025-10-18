#!/bin/bash

# Togglr Platform Installer
# Bash version of the Go installer

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration variables
INSTALL_DIR="/opt/togglr"
DOCKER_REGISTRY=""
PLATFORM_VERSION="latest"

# Function to print colored output
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if running as root
check_root() {
    if [[ $EUID -ne 0 ]]; then
        print_error "This installer must be run as root (sudo)"
        exit 1
    fi
}

# Function to check required commands
check_required_commands() {
    local missing_commands=()
    
    # Check for curl
    if ! command -v curl &> /dev/null; then
        missing_commands+=("curl")
    fi
    
    # Check for openssl
    if ! command -v openssl &> /dev/null; then
        missing_commands+=("openssl")
    fi
    
    if [[ ${#missing_commands[@]} -gt 0 ]]; then
        print_error "Missing required commands: ${missing_commands[*]}"
        print_info "Please install the missing commands and try again."
        print_info "On Ubuntu/Debian: sudo apt-get install curl openssl"
        print_info "On CentOS/RHEL: sudo yum install curl openssl"
        print_info "On macOS: brew install curl openssl"
        exit 1
    fi
}

# Function to print welcome message
print_welcome() {
    echo "================================================="
    echo "       Welcome to Togglr Platform Installer      "
    echo "================================================="
    echo "This installer will set up the Togglr platform on your system."
    echo "It will create necessary directories and configuration files."
    echo
}

# Function to read user input with validation
read_input() {
    local prompt="$1"
    local validation_func="$2"
    local input=""

    while true; do
        echo -n "$prompt: " >&2
        read -r input
        input=$(echo "$input" | xargs) # trim whitespace

        if [[ -n "$input" ]]; then
            if [[ -n "$validation_func" ]]; then
                if $validation_func "$input"; then
                    break
                fi
            else
                break
            fi
        else
            print_error "Input cannot be empty. Please try again."
        fi
    done

    echo "$input" >&2
}

# Function to read yes/no input
read_yes_no() {
    local prompt="$1"
    local input=""

    while true; do
        echo -n "$prompt (y/n): " >&2
        read -r input
        input=$(echo "$input" | tr '[:upper:]' '[:lower:]')

        if [[ "$input" == "y" || "$input" == "yes" ]]; then
            return 0
        elif [[ "$input" == "n" || "$input" == "no" ]]; then
            return 1
        else
            print_error "Please enter 'y' or 'n'"
        fi
    done
}

# Function to validate email
validate_email() {
    local email="$1"
    if [[ "$email" =~ ^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$ ]] && [[ ! "$email" =~ \.\..*@ ]]; then
        return 0
    else
        print_error "Please enter a valid email address"
        return 1
    fi
}

# Function to generate random string
generate_random_string() {
    local length="$1"
    LC_ALL=C tr -dc 'A-Za-z0-9' < /dev/urandom | head -c "$length"
}

# Function to create directories
create_directories() {
    print_info "Creating installation directories..."
    
    local directories=(
        "$INSTALL_DIR"
        "$INSTALL_DIR/nginx/ssl"
        "$INSTALL_DIR/secrets"
        "$INSTALL_DIR/nats"
        "$INSTALL_DIR/postgresql"
    )
    
    for dir in "${directories[@]}"; do
        print_info "Creating directory: $dir"
        mkdir -p "$dir"
        chmod 600 "$dir"
    done
}

# Function to generate SSL certificate
generate_ssl_certificate() {
    local domain="$1"
    local key_path="$INSTALL_DIR/nginx/ssl/nginx_key.pem"
    local cert_path="$INSTALL_DIR/nginx/ssl/nginx_cert.pem"
    
    print_info "Generating self-signed SSL certificate for $domain..."
    
    # Generate private key
    openssl genrsa -out "$key_path" 2048 2>/dev/null
    
    # Generate certificate
    openssl req -new -x509 -key "$key_path" -out "$cert_path" -days 3650 -subj "/C=US/ST=State/L=City/O=Togglr/CN=$domain" 2>/dev/null
    
    # Set proper permissions
    chmod 600 "$key_path"
    chmod 644 "$cert_path"
    
    print_success "SSL certificate generated successfully"
}

# Function to render template
render_template() {
    local template_file="$1"
    local output_file="$2"
    local temp_file="$3"
    
    # Create temporary file with template content
    cat > "$temp_file" << 'EOF'
#!/bin/bash
# Template renderer
EOF
    
    # Add template content
    cat "$template_file" >> "$temp_file"
    
    # Make it executable and run
    chmod +x "$temp_file"
    "$temp_file" > "$output_file"
    
    # Clean up
    rm -f "$temp_file"
}

# Function to collect user input from environment or prompt
collect_user_input() {
    print_info "=== Platform Configuration ==="
    
    # Check if running in non-interactive mode with environment variables
    if [[ -n "$TOGGLR_ADMIN_EMAIL" && -n "$TOGGLR_DOMAIN" && -n "$TOGGLR_MAILER_ADDR" ]]; then
        print_info "Using environment variables for configuration..."
        ADMIN_EMAIL="$TOGGLR_ADMIN_EMAIL"
        DOMAIN="$TOGGLR_DOMAIN"
        FRONTEND_URL="https://$DOMAIN"
        MAILER_ADDR="$TOGGLR_MAILER_ADDR"
        MAILER_USER="${TOGGLR_MAILER_USER:-admin@$DOMAIN}"
        MAILER_PASSWORD="${TOGGLR_MAILER_PASSWORD:-password}"
        MAILER_FROM="${TOGGLR_MAILER_FROM:-noreply@$DOMAIN}"
        HAS_EXISTING_SSL_CERT="${TOGGLR_HAS_SSL_CERT:-false}"
        return
    fi
    
    # Get admin email
    ADMIN_EMAIL=$(read_input "Enter administrator email" "validate_email")
    
    # Get domain
    DOMAIN=$(read_input "Enter domain for the platform")
    FRONTEND_URL="https://$DOMAIN"
    
    # Ask about SSL certificate
    print_info "Asking about SSL certificate..."
    if read_yes_no "Do you have an existing SSL certificate for this domain?"; then
        HAS_EXISTING_SSL_CERT=true
        print_info "You will need to place your SSL certificate and key files at:"
        print_info "  - Certificate: $INSTALL_DIR/nginx/ssl/nginx_cert.pem"
        print_info "  - Key: $INSTALL_DIR/nginx/ssl/nginx_key.pem"
        print_info "You will be reminded about this at the end of installation."
    else
        HAS_EXISTING_SSL_CERT=false
        print_info "A self-signed SSL certificate will be generated for you at the end of installation."
    fi
    
    print_info "=== SMTP Server Configuration ==="
    
    # Get SMTP server details
    MAILER_ADDR=$(read_input "Enter SMTP server address (including port)")
    MAILER_USER=$(read_input "Enter SMTP user")
    MAILER_PASSWORD=$(read_input "Enter SMTP password")
    MAILER_FROM=$(read_input "Enter email address for sending emails (from)")
    
}

# Function to generate secrets
generate_secrets() {
    print_info "Generating secure passwords and keys..."
    
    PG_PASSWORD=$(generate_random_string 12)
    SECRET_KEY=$(generate_random_string 32)
    JWT_SECRET_KEY=$(generate_random_string 32)
    ADMIN_TMP_PASSWORD=$(generate_random_string 12)
    
    print_success "Generated secure passwords and keys"
}

# Function to create platform.env
create_platform_env() {
    local platform_env_file="$INSTALL_DIR/platform.env"
    
    cat > "$platform_env_file" << EOF
DOCKER_REGISTRY=$DOCKER_REGISTRY
DOMAIN=$DOMAIN
DB_PASSWORD=$PG_PASSWORD
SSL_CERT=nginx_cert.pem
SSL_CERT_KEY=nginx_key.pem
PLATFORM_VERSION=$PLATFORM_VERSION
EOF
    
    print_success "Created $platform_env_file"
}

# Function to create config.env
create_config_env() {
    local config_env_file="$INSTALL_DIR/config.env"
    
    cat > "$config_env_file" << EOF
LOGGER_LEVEL=info
FRONTEND_URL=$FRONTEND_URL
SECRET_KEY=$SECRET_KEY

ADMIN_EMAIL=$ADMIN_EMAIL
ADMIN_TMP_PASSWORD=$ADMIN_TMP_PASSWORD

API_SERVER_ADDR=:8080
API_SERVER_READ_TIMEOUT=15s
API_SERVER_WRITE_TIMEOUT=30s
API_SERVER_IDLE_TIMEOUT=60s

TECH_SERVER_ADDR=:8081
TECH_SERVER_READ_TIMEOUT=15s
TECH_SERVER_WRITE_TIMEOUT=30s
TECH_SERVER_IDLE_TIMEOUT=60s

SDK_SERVER_ADDR=:8090
SDK_SERVER_READ_TIMEOUT=15s
SDK_SERVER_WRITE_TIMEOUT=30s
SDK_SERVER_IDLE_TIMEOUT=60s

WS_SERVER_ADDR=:8082
WS_SERVER_READ_TIMEOUT=15s
WS_SERVER_WRITE_TIMEOUT=30s
WS_SERVER_IDLE_TIMEOUT=60s

# PostgreSQL
POSTGRES_HOST=togglr-postgresql
POSTGRES_DATABASE=togglr
POSTGRES_PASSWORD=$PG_PASSWORD
POSTGRES_PORT=5432
POSTGRES_USER=togglr
MIGRATIONS_DIR=/migrations

# NATS
NATS_URL=nats://togglr-nats:4222

# JWT
JWT_SECRET_KEY=$JWT_SECRET_KEY
ACCESS_TOKEN_TTL=3h
REFRESH_TOKEN_TTL=168h
RESET_PASSWORD_TTL=8h

# SAML
SAML_ENABLED=false
SAML_CREATE_CERTS=false
SAML_ENTITY_ID=https://$DOMAIN/api/v1/saml/metadata
SAML_CALLBACK_URL="https://$DOMAIN/api/v1/auth/sso/callback"
SAML_CERTIFICATE_PATH="/opt/togglr/secrets/saml_cert.pem"
SAML_PRIVATE_KEY_PATH="/opt/togglr/secrets/saml_key.pem"
SAML_ALLOWED_ISSUERS=http://idp.$DOMAIN/realms/togglr,https://$DOMAIN/api/v1/saml/acs
SAML_ATTRIBUTE_MAPPING=uid:username,email:email
SAML_PUBLIC_ROOT_URL="https://$DOMAIN"
SAML_IDP_METADATA_URL="http://idp.$DOMAIN/realms/togglr/protocol/saml/descriptor"

# Mailer
MAILER_ADDR=$MAILER_ADDR
MAILER_USER=$MAILER_USER
MAILER_PASSWORD=$MAILER_PASSWORD
MAILER_FROM=$MAILER_FROM
MAILER_ALLOW_INSECURE=false
MAILER_BASE_URL="https://$DOMAIN"
EOF
    
    print_success "Created $config_env_file"
}

# Function to download docker-compose.yml
download_docker_compose() {
    local docker_compose_file="$INSTALL_DIR/docker-compose.yml"
    local docker_compose_url="https://raw.githubusercontent.com/togglr-project/togglr/main/prod/docker-compose.yml"
    
    print_info "Downloading docker-compose.yml from GitHub..."
    
    if curl -fsSL "$docker_compose_url" -o "$docker_compose_file"; then
        print_success "Downloaded $docker_compose_file"
    else
        print_error "Failed to download docker-compose.yml from $docker_compose_url"
        exit 1
    fi
}

# Function to download NATS configuration
download_nats_config() {
    local nats_conf_file="$INSTALL_DIR/nats/nats.conf"
    local nats_conf_url="https://raw.githubusercontent.com/togglr-project/togglr/main/prod/nats.conf"
    
    print_info "Downloading nats.conf from GitHub..."
    
    if curl -fsSL "$nats_conf_url" -o "$nats_conf_file"; then
        print_success "Downloaded $nats_conf_file"
    else
        print_error "Failed to download nats.conf from $nats_conf_url"
        exit 1
    fi
}

# Function to create Makefile
create_makefile() {
    local makefile="$INSTALL_DIR/Makefile"
    
    # Create Makefile with proper tabs
    cat > "$makefile" << 'EOF'
_COMPOSE=docker compose -f docker-compose.yml --project-name togglr --env-file platform.env

.DEFAULT_GOAL := help

.PHONY: help
help: ## Print this message
	@echo "$$(grep -hE '^\S+:.*##' $(MAKEFILE_LIST) | sed -e 's/:.*##\s*/:/' -e 's/^\(.\+\):\(.*\)/\\x1b[36m\1\\x1b[m:\2/' | column -c2 -t -s :)"

.PHONY: up
up: ## Up the environment in docker compose
	${_COMPOSE} up -d

.PHONY: down
down: ## Down the environment in docker compose
	${_COMPOSE} down --remove-orphans

.PHONY: pull
pull: ## Pull images from remote Docker registry
	${_COMPOSE} pull
EOF
    
    # Ensure tabs are used instead of spaces (compatible with both Linux and macOS)
    if [[ "$OSTYPE" == "darwin"* ]]; then
        sed -i '' 's/^    /\t/g' "$makefile"
    else
        sed -i 's/^    /\t/g' "$makefile"
    fi
    
    print_success "Created $makefile"
}

# Function to print final information
print_final_info() {
    echo
    print_success "Installation completed successfully!"
    print_info "The platform has been installed in $INSTALL_DIR"
    print_info "You can check and modify settings in $INSTALL_DIR/platform.env and $INSTALL_DIR/config.env"
    print_info "A Makefile has been created in $INSTALL_DIR with commands for starting and stopping the platform"
    
    echo
    print_info "ADMIN LOGIN INFORMATION:"
    print_info "  - Email: $ADMIN_EMAIL"
    print_info "  - Temporary Password: $ADMIN_TMP_PASSWORD"
    print_info "Please use these credentials to log in to the platform. You will be prompted to change the password on first login."
    
    if [[ "$HAS_EXISTING_SSL_CERT" == true ]]; then
        echo
        print_warning "REMINDER: Don't forget to place your SSL certificate and key files at:"
        print_info "  - Certificate: $INSTALL_DIR/nginx/ssl/nginx_cert.pem"
        print_info "  - Key: $INSTALL_DIR/nginx/ssl/nginx_key.pem"
    fi
    
}

# Main installation function
main() {
    # Check if running as root
    check_root
    
    # Check required commands
    check_required_commands
    
    # Check if running in interactive terminal or with environment variables
    if [[ ! -t 0 && -z "$TOGGLR_ADMIN_EMAIL" ]]; then
        print_error "This installer requires an interactive terminal or environment variables."
        print_info "Please run the installer in an interactive terminal session or set environment variables:"
        print_info "  TOGGLR_ADMIN_EMAIL=admin@example.com"
        print_info "  TOGGLR_DOMAIN=example.com"
        print_info "  TOGGLR_MAILER_ADDR=smtp.example.com:587"
        print_info "Example: sudo ./install.sh"
        exit 1
    fi
    
    # Print welcome message
    print_welcome
    
    # Inform about installation directory
    print_info "The platform will be installed in the $INSTALL_DIR directory."
    echo
    
    # Ask for confirmation to proceed
    if ! read_yes_no "Do you want to continue with the installation?"; then
        print_info "Installation cancelled."
        exit 0
    fi
    
    # Collect user input
    collect_user_input
    
    # Generate passwords and other required values
    generate_secrets
    
    # Create required directories
    create_directories
    
    # Create configuration files
    create_platform_env
    create_config_env
    download_docker_compose
    download_nats_config
    create_makefile
    
    # Handle SSL certificate based on user's choice
    if [[ "$HAS_EXISTING_SSL_CERT" == false ]]; then
        generate_ssl_certificate "$DOMAIN"
    fi
    
    # Print final information
    print_final_info
}

# Run main function
main "$@"
