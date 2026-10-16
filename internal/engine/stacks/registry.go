package stacks

type StackDefinition struct {
	ID              string
	Name            string
	RootMarkers     []string
	Deps            []string
	DefaultRoles    []StackRole
	OptimizeCommand string
}

var Registry = map[string]StackDefinition{
	"laravel": {
		ID:              "laravel",
		Name:            "Laravel",
		RootMarkers:     []string{"artisan", "composer.json"},
		Deps:            []string{"laravel/framework"},
		OptimizeCommand: "php artisan optimize",
		DefaultRoles: []StackRole{
			{Name: "queue", Kind: RoleKindWorker, Command: "exec php artisan queue:work", StopGracePeriod: "90s", Unless: &RoleGates{Deps: []string{"laravel/horizon"}}},
			{Name: "horizon", Kind: RoleKindWorker, Command: "exec php artisan horizon", StopGracePeriod: "90s", When: &RoleGates{Deps: []string{"laravel/horizon"}}},
			{Name: "scheduler", Kind: RoleKindScheduler, Command: "exec php artisan schedule:work", StopGracePeriod: "10s", Singleton: true},
		},
	},
	"symfony": {
		ID:          "symfony",
		Name:        "Symfony",
		RootMarkers: []string{"composer.json", "symfony.lock"},
		Deps:        []string{"symfony/framework-bundle"},
	},
	"rails": {
		ID:          "rails",
		Name:        "Ruby on Rails",
		RootMarkers: []string{"Gemfile", "bin/rails"},
		DefaultRoles: []StackRole{
			{Name: "jobs", Kind: RoleKindWorker, Command: "exec bin/jobs", StopGracePeriod: "10s", When: &RoleGates{Deps: []string{"solid_queue"}, Files: []string{"bin/jobs"}}},
			{Name: "sidekiq", Kind: RoleKindWorker, Command: "exec bundle exec sidekiq", StopGracePeriod: "45s", When: &RoleGates{Deps: []string{"sidekiq"}}},
		},
	},
	"sinatra": {
		ID:          "sinatra",
		Name:        "Sinatra",
		RootMarkers: []string{"Gemfile"},
		Deps:        []string{"sinatra"},
	},
	"django": {
		ID:          "django",
		Name:        "Django",
		RootMarkers: []string{"manage.py"},
	},
	"flask": {
		ID:          "flask",
		Name:        "Flask",
		RootMarkers: []string{"requirements.txt"},
		Deps:        []string{"flask"},
	},
	"nextjs": {
		ID:          "nextjs",
		Name:        "Next.js",
		RootMarkers: []string{"package.json"},
		Deps:        []string{"next"},
	},
	"node": {
		ID:          "node",
		Name:        "Node.js",
		RootMarkers: []string{"package.json"},
	},
	"go": {
		ID:          "go",
		Name:        "Go",
		RootMarkers: []string{"go.mod"},
	},
	"rust": {
		ID:          "rust",
		Name:        "Rust",
		RootMarkers: []string{"Cargo.toml"},
	},
	"python": {
		ID:          "python",
		Name:        "Python",
		RootMarkers: []string{"requirements.txt"},
	},
	"php": {
		ID:          "php",
		Name:        "PHP",
		RootMarkers: []string{"composer.json"},
	},
	"ruby": {
		ID:          "ruby",
		Name:        "Ruby",
		RootMarkers: []string{"Gemfile"},
	},
	"static": {
		ID:          "static",
		Name:        "Static Site",
		RootMarkers: []string{"index.html"},
	},
	"docker": {
		ID:          "docker",
		Name:        "Dockerfile",
		RootMarkers: []string{"Dockerfile"},
	},
}

var DetectOrder = []string{
	"laravel", "symfony", "rails", "sinatra", "django", "flask",
	"nextjs", "node", "go", "rust", "python", "php", "ruby", "static", "docker",
}
