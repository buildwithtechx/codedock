package compose

import (
	"testing"
)

func TestComposeRejectsDiscardedAndUnsafeFields(t *testing.T) {
	for name, content := range map[string]string{
		"privileged":       "services:\n  app:\n    image: alpine\n    privileged: true\n",
		"file variables":   "services:\n  app:\n    image: alpine\n    env_file: .env\n",
		"host mount":       "services:\n  app:\n    image: alpine\n    volumes: ['/var/run/docker.sock:/socket']\n",
		"anonymous volume": "services:\n  app:\n    image: alpine\n    volumes: ['/data']\n",
		"local build":      "services:\n  app:\n    build: .\n",
		"cycle":            "services:\n  a:\n    image: alpine\n    depends_on: [b]\n  b:\n    image: alpine\n    depends_on: [a]\n",
		"secret file":      "services:\n  app:\n    image: alpine\nsecrets:\n  password:\n    file: /etc/passwd\n",
	} {
		t.Run(name, func(t *testing.T) {
			if err := ValidateStackInput(content); err == nil {
				t.Fatal("unsafe or unsupported Compose fields were silently accepted")
			}
		})
	}
}

func TestComposeAcceptsPreservedCoreRuntimeFields(t *testing.T) {
	content := `services:
  web:
    build:
      context: https://github.com/example/app.git#main:frontend
      dockerfile: Dockerfile
      args: {PUBLIC_URL: "${PUBLIC_URL:?required}"}
    environment: [DATABASE_URL=postgres://database/app]
    command: [node, server.js]
    ports: ['8080:3000']
    volumes: [cache:/cache]
    networks: [private]
    depends_on:
      database: {condition: service_healthy}
    healthcheck:
      test: [CMD, curl, -f, 'http://localhost:3000/health']
  database:
    image: postgres:16
    environment: {POSTGRES_PASSWORD: 'example-only'}
    volumes: [data:/var/lib/postgresql/data]
    networks: [private]
    healthcheck: {test: [CMD, pg_isready]}
volumes: {data: {}, cache: {}}
networks: {private: {internal: true}}
`
	if err := ValidateStackInput(content); err != nil {
		t.Fatal(err)
	}
}

func TestComposeRejectsCrossStackOwnershipAndIgnoredClusterPlacement(t *testing.T) {
	id := "id"
	for _, content := range []string{
		`{"services":{"app":{"image":"alpine"}},"volumes":{"data":{"name":"other-data"}}}`,
		`{"services":{"app":{"image":"alpine","deploy":{"placement":{"constraints":["node.role == manager"]}}}}}`,
		`{"services":{"app":{"image":"alpine","labels":{"traefik.http.routers.other.rule":"Host(` + "`other.example`" + `)"}}}}`,
	} {
		if err := validateStackOwnership(content, id); err == nil {
			t.Fatal("unsupported ownership or target capability accepted")
		}
	}
}
