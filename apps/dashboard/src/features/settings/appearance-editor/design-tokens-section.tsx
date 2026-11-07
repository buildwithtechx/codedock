import { Button } from '@codedock/ui/button';
import { Input } from '@codedock/ui/input';
import { RotateCcw, Search } from 'lucide-react';
import { useMemo, useState } from 'react';
import { TOKEN_CATEGORIES, TOKEN_DEFINITIONS, type TokenCategory } from './appearance-tokens-data';

interface DesignTokensSectionProps {
  tokens: Record<string, string>;
  onSetToken: (key: string, value: string) => void;
  onResetToken: (key: string) => void;
}

function hexFromValue(val: string): string {
  if (val.startsWith('#') && (val.length === 7 || val.length === 4)) {
    return val;
  }
  return '#7c3aed';
}

export function DesignTokensSection({
  tokens,
  onSetToken,
  onResetToken,
}: DesignTokensSectionProps) {
  const [search, setSearch] = useState('');
  const [activeCategory, setActiveCategory] = useState<TokenCategory>('surfaces');

  const filteredTokens = useMemo(() => {
    return TOKEN_DEFINITIONS.filter((def) => {
      const matchesSearch =
        !search ||
        def.key.toLowerCase().includes(search.toLowerCase()) ||
        def.label.toLowerCase().includes(search.toLowerCase());
      const matchesCat = search ? true : def.category === activeCategory;
      return matchesSearch && matchesCat;
    });
  }, [search, activeCategory]);

  return (
    <section className="rounded-2xl border border-border/60 bg-card p-5 shadow-xs">
      <h3 className="font-semibold text-foreground text-sm">Design tokens</h3>

      {/* Search Input */}
      <div className="relative mt-3">
        <Search className="absolute top-1/2 left-3 size-4 -translate-y-1/2 text-muted-foreground" />
        <Input
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          placeholder="Find a token..."
          className="h-9 border-border/60 bg-background/50 pl-9 text-xs focus-visible:ring-1"
        />
      </div>

      {/* Categories Bar */}
      <div className="mt-3 flex flex-wrap gap-1.5 border-border/40 border-b pb-3">
        {TOKEN_CATEGORIES.map((cat) => {
          const isSelected = activeCategory === cat.id && !search;
          return (
            <button
              key={cat.id}
              type="button"
              onClick={() => {
                setSearch('');
                setActiveCategory(cat.id);
              }}
              className={`rounded-lg px-2.5 py-1.5 font-medium text-xs transition-colors ${
                isSelected
                  ? 'bg-foreground/[0.08] text-foreground'
                  : 'text-muted-foreground hover:bg-muted/50 hover:text-foreground'
              }`}
            >
              {cat.label}
            </button>
          );
        })}
      </div>

      {/* Token Rows List */}
      <div className="mt-3 divide-y divide-border/40">
        {filteredTokens.map((def) => {
          const currentValue = tokens[def.key] ?? def.defaultValue;
          const isOverridden = Boolean(tokens[def.key]);

          return (
            <div
              key={def.key}
              className="group flex flex-col gap-2 py-2.5 sm:flex-row sm:items-center sm:justify-between"
            >
              <div className="min-w-0 flex-1">
                <span className="font-mono text-foreground/90 text-xs">{def.key}</span>
              </div>

              <div className="flex items-center gap-2">
                {/* Color Swatch Button with Hidden Native Color Picker */}
                <label className="relative flex size-8 shrink-0 cursor-pointer items-center justify-center overflow-hidden rounded-lg border border-border/80 shadow-inner">
                  <div
                    className="size-full"
                    style={{
                      backgroundColor: currentValue.startsWith('#')
                        ? currentValue
                        : currentValue.startsWith('rgba')
                          ? currentValue
                          : '#8b5cf6',
                    }}
                  />
                  <input
                    type="color"
                    value={hexFromValue(currentValue)}
                    onChange={(e) => onSetToken(def.key, e.target.value)}
                    className="sr-only"
                    aria-label={`Pick color for ${def.key}`}
                  />
                </label>

                {/* Direct Text Input */}
                <Input
                  value={currentValue}
                  onChange={(e) => onSetToken(def.key, e.target.value)}
                  className="h-8 w-48 border-border/60 bg-background/50 font-mono text-xs"
                  aria-label={`Token value for ${def.key}`}
                />

                {isOverridden && (
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={() => onResetToken(def.key)}
                    className="size-8 text-muted-foreground hover:text-foreground"
                    title="Reset to default"
                  >
                    <RotateCcw className="size-3.5" />
                  </Button>
                )}
              </div>
            </div>
          );
        })}

        {filteredTokens.length === 0 && (
          <p className="py-6 text-center text-muted-foreground text-xs">
            No design tokens matched your search.
          </p>
        )}
      </div>
    </section>
  );
}
