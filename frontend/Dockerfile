# ─────────────────────────────────────────────
# frontend/Dockerfile  –  Next.js + pnpm v10
# ─────────────────────────────────────────────
########################
# 1. “deps”  – install *all* deps (dev + prod)
########################
FROM node:20-alpine AS deps
WORKDIR /app

# Corepack ships with Node ≥16; this activates the exact pnpm used by the repo
RUN corepack enable && corepack prepare pnpm@10.13.1 --activate \
  && apk add --no-cache libc6-compat          # fixes native deps like sharp/oxide

COPY package.json pnpm-lock.yaml ./

# Full install so dev tools (typescript, tailwind) exist for the build step
# --frozen-lockfile  ➜ strict, CI‑friendly
RUN pnpm install --frozen-lockfile            # runs build scripts once ✔

########################
# 2. “builder” – copy source & build bundle
########################
FROM node:20-alpine AS builder
WORKDIR /app
RUN corepack enable                            # pnpm already in PATH

# Re‑use deps from previous stage (fast!)
COPY --from=deps /app/node_modules ./node_modules
COPY --from=deps /app/pnpm-lock.yaml ./pnpm-lock.yaml
COPY --from=deps /app/package.json ./package.json

# Bring in the actual source (thanks to .dockerignore, no local node_modules!)
COPY . .
ARG NEXT_PUBLIC_API
ARG NEXT_PUBLIC_WS
ENV NEXT_PUBLIC_API=$NEXT_PUBLIC_API
ENV NEXT_PUBLIC_WS=$NEXT_PUBLIC_WS

RUN pnpm run build      # generates .next/, etc.

########################
# 3. “runner” – tiny runtime image
########################
FROM node:20-alpine AS runner
WORKDIR /app
ENV NODE_ENV=production

# Copy only the artefacts we need
COPY --from=builder /app/.next       ./.next
COPY --from=builder /app/public      ./public
COPY --from=builder /app/package.json ./package.json
COPY --from=builder /app/pnpm-lock.yaml ./pnpm-lock.yaml

# Install *production‑only* deps, then prune store
RUN corepack enable \
  && pnpm install --frozen-lockfile --prod \
  && pnpm store prune

EXPOSE 3010
CMD ["pnpm", "start"]                 # package.json: "start": "next start -p 3000"
