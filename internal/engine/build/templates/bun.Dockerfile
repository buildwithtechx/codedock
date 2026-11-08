FROM oven/bun:1-alpine

WORKDIR /app

COPY package*.json bun.lock* ./

RUN bun install --frozen-lockfile

COPY . .

RUN bun run build --if-present

EXPOSE 3000

CMD ["bun", "start"]

