FROM oven/bun:latest
WORKDIR /ui
COPY . .
RUN bun install
RUN bun run build
EXPOSE 3000
CMD bun .output/server/index.mjs
