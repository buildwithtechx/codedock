import { ArrowUpRight, Boxes, Search } from 'lucide-react';
import { useState } from 'react';
import { SiteLink } from '../../components/site-link';
import { productLinks } from '../../lib/product-links';
import { templateRecipes } from './template-recipes';

export function TemplateCatalog() {
  const [search, setSearch] = useState('');
  const [category, setCategory] = useState('All');
  const categories = ['All', ...new Set(templateRecipes.map((recipe) => recipe.category))];
  const visible = templateRecipes.filter(
    (recipe) =>
      (category === 'All' || category === recipe.category) &&
      `${recipe.name} ${recipe.description} ${recipe.category} ${recipe.workflow}`
        .toLowerCase()
        .includes(search.toLowerCase().trim())
  );
  return (
    <div className="mx-auto max-w-6xl px-6 py-16 sm:py-24">
      <p className="font-semibold text-primary text-xs uppercase tracking-widest">
        Template library
      </p>
      <h1 className="mt-5 max-w-4xl text-balance font-extrabold text-4xl tracking-tight sm:text-6xl">
        Start with the services your application needs.
      </h1>
      <p className="mt-6 max-w-3xl text-lg text-muted-foreground leading-relaxed">
        Explore database workflows and container recipes. Each guide explains the setup,
        persistence, and networking you need before deploying.
      </p>
      <p className="mt-4 max-w-3xl text-muted-foreground text-sm leading-relaxed">
        Database workflows use Codedock's provisioning tools. Docker and Compose recipes require
        configuration; they are not automatic one-click installs.
      </p>
      <div className="mt-10 flex flex-col gap-5">
        <label className="flex max-w-md items-center gap-3 rounded-xl border border-border bg-card/70 px-4 py-3 focus-within:border-primary focus-within:ring-2 focus-within:ring-primary/30">
          <Search aria-hidden="true" className="size-4 text-muted-foreground" />
          <span className="sr-only">Search template recipes</span>
          <input
            value={search}
            onChange={(event) => setSearch(event.target.value)}
            placeholder="Search services, databases, or workflows..."
            type="search"
            className="min-w-0 flex-1 bg-transparent text-sm outline-none"
          />
        </label>
        <fieldset className="flex flex-wrap gap-2">
          <legend className="sr-only">Filter by category</legend>
          {categories.map((item) => (
            <button
              key={item}
              type="button"
              aria-pressed={category === item}
              onClick={() => setCategory(item)}
              className={`rounded-full border px-4 py-2 text-xs transition-colors ${category === item ? 'border-primary bg-primary text-white' : 'border-border text-muted-foreground hover:border-primary/50'}`}
            >
              {item}
            </button>
          ))}
        </fieldset>
      </div>
      <p aria-live="polite" className="mt-6 text-muted-foreground text-xs">
        {visible.length} {visible.length === 1 ? 'recipe' : 'recipes'}
      </p>
      <div className="mt-5 grid gap-5 sm:grid-cols-2 lg:grid-cols-3">
        {visible.map((recipe) => (
          <SiteLink
            key={recipe.id}
            href={productLinks.docs + recipe.docsPath}
            className="group flex flex-col rounded-2xl border border-border bg-card/70 p-6 transition-colors hover:border-primary/50"
          >
            <div className="mb-6 flex items-center justify-between">
              <span className="flex size-11 items-center justify-center rounded-xl border border-primary/20 bg-primary/10">
                <Boxes className="size-5 text-primary" />
              </span>
              <span className="rounded-full border border-border px-2 py-1 text-[10px] text-muted-foreground">
                {recipe.workflow}
              </span>
            </div>
            <p className="text-primary text-xs">{recipe.category}</p>
            <h2 className="mt-2 font-bold text-xl">{recipe.name}</h2>
            <p className="mt-3 flex-1 text-muted-foreground text-sm leading-relaxed">
              {recipe.description}
            </p>
            <span className="mt-6 inline-flex items-center gap-2 font-semibold text-primary text-xs">
              Open setup guide
              <ArrowUpRight className="size-3.5" />
            </span>
          </SiteLink>
        ))}
      </div>
      {visible.length === 0 && (
        <div className="mt-5 rounded-2xl border border-border p-10 text-center">
          <h2 className="font-semibold">No recipes match your search.</h2>
          <p className="mt-2 text-muted-foreground text-sm">
            Try another name or reset the filters.
          </p>
          <button
            type="button"
            onClick={() => {
              setSearch('');
              setCategory('All');
            }}
            className="mt-5 rounded-lg bg-primary px-4 py-2 text-sm text-white"
          >
            Reset filters
          </button>
        </div>
      )}
    </div>
  );
}
