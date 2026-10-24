package deployments

import (
	"os"
	"path/filepath"
	"strings"
)

func detectAppIcon(sourceDir string) string {
	if b, err := os.ReadFile(filepath.Join(sourceDir, "package.json")); err == nil {
		s := strings.ToLower(string(b))
		if strings.Contains(s, "\"next\"") {
			return "nextjs"
		}
		if strings.Contains(s, "\"nuxt\"") {
			return "nuxt"
		}
		if strings.Contains(s, "\"@sveltejs/kit\"") {
			return "sveltekit"
		}
		if strings.Contains(s, "\"@solidjs/start\"") {
			return "solidstart"
		}
		if strings.Contains(s, "\"@remix-run") {
			return "remix"
		}
		if strings.Contains(s, "\"astro\"") {
			return "astro"
		}
		if strings.Contains(s, "\"@tanstack/start\"") {
			return "tanstack"
		}
		if strings.Contains(s, "\"expo\"") {
			return "expo"
		}
		if strings.Contains(s, "\"react-native\"") {
			return "react-native"
		}
		if strings.Contains(s, "\"vite\"") {
			return "vite"
		}
		if strings.Contains(s, "\"vue\"") {
			return "vue"
		}
		if strings.Contains(s, "\"svelte\"") {
			return "svelte"
		}
		if strings.Contains(s, "\"@angular/core\"") {
			return "angular"
		}
		if strings.Contains(s, "\"@builder.io/qwik\"") {
			return "qwik"
		}
		if strings.Contains(s, "\"gatsby\"") {
			return "gatsby"
		}
		if strings.Contains(s, "\"@redwoodjs/") {
			return "redwoodjs"
		}
		if strings.Contains(s, "\"electron\"") {
			return "electron"
		}
		if strings.Contains(s, "\"@tauri-apps/") {
			return "tauri"
		}
		if strings.Contains(s, "\"hono\"") {
			return "hono"
		}
		if strings.Contains(s, "\"elysia\"") {
			return "elysia"
		}
		if strings.Contains(s, "\"@nestjs/core\"") {
			return "nestjs"
		}
		if strings.Contains(s, "\"fastify\"") {
			return "fastify"
		}
		if strings.Contains(s, "\"express\"") {
			return "express"
		}
		if strings.Contains(s, "\"koa\"") {
			return "koa"
		}
		if strings.Contains(s, "\"@adonisjs/") {
			return "adonisjs"
		}
		if strings.Contains(s, "\"@strapi/") {
			return "strapi"
		}
		if strings.Contains(s, "\"payload\"") {
			return "payload"
		}
		if strings.Contains(s, "\"@trpc/") {
			return "trpc"
		}
		if strings.Contains(s, "\"graphql\"") {
			return "graphql"
		}
		if strings.Contains(s, "\"react\"") {
			return "react"
		}
	}

	files := map[string]string{
		"next.config.js": "nextjs", "next.config.mjs": "nextjs",
		"vite.config.ts": "vite", "vite.config.js": "vite",
		"astro.config.mjs": "astro", "nuxt.config.ts": "nuxt",
		"nest-cli.json": "nestjs", "svelte.config.js": "sveltekit",
		"remix.config.js": "remix", "angular.json": "angular",
		"pom.xml": "java", "build.gradle": "java",
		"manage.py": "django", "Gemfile": "ruby",
		"composer.json": "php", "Cargo.toml": "rust",
		"go.mod": "golang", "artisan": "laravel",
	}

	for file, icon := range files {
		if _, err := os.Stat(filepath.Join(sourceDir, file)); err == nil {
			return icon
		}
	}

	if b, err := os.ReadFile(filepath.Join(sourceDir, "requirements.txt")); err == nil {
		s := strings.ToLower(string(b))
		if strings.Contains(s, "fastapi") {
			return "fastapi"
		}
		if strings.Contains(s, "flask") {
			return "flask"
		}
		if strings.Contains(s, "django") {
			return "django"
		}
		return "python"
	}
	if _, err := os.Stat(filepath.Join(sourceDir, "package.json")); err == nil {
		return "nodejs"
	}

	return "git"
}
