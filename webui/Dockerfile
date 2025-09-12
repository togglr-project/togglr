FROM node:24.3 AS build

ARG VITE_VERSION
ARG VITE_BUILD_TIME
ARG VITE_IS_DEMO

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY . .

# Build with default values - will be overridden at runtime
RUN VITE_VERSION=${VITE_VERSION:-dev} \
    VITE_BUILD_TIME=${VITE_BUILD_TIME:-$(date -u +%Y-%m-%dT%H:%M:%SZ)} \
    VITE_IS_DEMO=${VITE_IS_DEMO:-false} \
    npm run build

FROM nginx:stable-alpine

COPY --from=build /app/dist /usr/share/nginx/html

# Copy user agreement files explicitly to ensure they're available
COPY ./public/USER_AGREEMENT_EN.txt /usr/share/nginx/html/USER_AGREEMENT_EN.txt
COPY ./public/USER_AGREEMENT_RU.txt /usr/share/nginx/html/USER_AGREEMENT_RU.txt

COPY ./nginx/nginx.conf /etc/nginx/conf.d/default.conf

# Copy startup script
COPY ./docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh

EXPOSE 80

ENTRYPOINT ["/docker-entrypoint.sh"]
