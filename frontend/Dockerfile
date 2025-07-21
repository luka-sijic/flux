########################
# 1. Install deps with pnpm
########################
FROM node:20-alpine AS deps
WORKDIR /app

# Enable corepack (ships with Node 16+) so we get pnpm without a global npm install
RUN corepack enable && corepack prepare pnpm@latest --activate \
    && apk add --no-cache libc6-compat      # for some node-gyp builds

# Copy only manifests first so layer caching works
COPY package.json pnpm-lock.yaml ./

# Install deps *without* devDependencies
RUN pnpm install --frozen-lockfile --prod

########################
# 2. Build production Next.js bundle
########################
FROM node:20-alpine AS builder
WORKDIR /app

# Re‑use dependencies from previous stage
COPY --from=deps /app/node_modules ./node_modules
COPY --from=deps /app/pnpm-lock.yaml ./pnpm-lock.yaml
COPY --from=deps /app/package.json ./package.json

# Bring in the rest of the source
COPY . .

RUN corepack enable && pnpm run build     # outputs .next/

########################
# 3. Lightweight runtime image
########################
FROM node:20-alpine AS runner
WORKDIR /app
ENV NODE_ENV=production

# Copy built artefacts and deps
COPY --from=builder /app ./

EXPOSE 3010
CMD ["pnpm", "start"]          # "start": "next start -p 3000" in package.json
