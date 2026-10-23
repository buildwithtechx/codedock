import { Upload } from 'lucide-react';
import { Button } from '@codedock/ui/button';
import { Input } from '@codedock/ui/input';

interface TypographySectionProps {
  interfaceFont?: string;
  monoFont?: string;
  interfaceSize?: string | number;
  monoSize?: string | number;
  lineHeight?: string | number;
  fontWeight?: string | number;
  onChange: (fields: {
    interfaceFont?: string;
    monoFont?: string;
    interfaceSize?: string | number;
    monoSize?: string | number;
    lineHeight?: string | number;
    fontWeight?: string | number;
  }) => void;
}

const INTERFACE_FONTS = ['Gellix', 'Manrope', 'Inter', 'Geist', 'System', 'Custom'];
const MONO_FONTS = ['JetBrains Mono', 'Fira Code', 'Geist Mono', 'System Mono', 'Custom'];

export function TypographySection({
  interfaceFont = 'Manrope',
  monoFont = 'JetBrains Mono',
  interfaceSize = 16,
  monoSize = 13,
  lineHeight = 1.5,
  fontWeight = 400,
  onChange,
}: TypographySectionProps) {
  return (
    <section className="rounded-2xl border border-border/60 bg-card p-5 shadow-xs">
      <h3 className="font-semibold text-foreground text-sm">Typography</h3>

      <div className="mt-4 space-y-4">
        {/* Row 1: Fonts */}
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <div className="space-y-1.5">
            <label className="font-medium text-muted-foreground text-xs">Interface font</label>
            <select
              value={interfaceFont}
              onChange={(e) => onChange({ interfaceFont: e.target.value })}
              className="flex h-9 w-full rounded-lg border border-border/60 bg-background/50 px-3 py-1 font-medium text-xs text-foreground shadow-xs outline-none focus:ring-1 focus:ring-primary"
            >
              {INTERFACE_FONTS.map((font) => (
                <option key={font} value={font}>
                  {font}
                </option>
              ))}
            </select>
          </div>

          <div className="space-y-1.5">
            <label className="font-medium text-muted-foreground text-xs">Monospace font</label>
            <select
              value={monoFont}
              onChange={(e) => onChange({ monoFont: e.target.value })}
              className="flex h-9 w-full rounded-lg border border-border/60 bg-background/50 px-3 py-1 font-mono text-xs text-foreground shadow-xs outline-none focus:ring-1 focus:ring-primary"
            >
              {MONO_FONTS.map((font) => (
                <option key={font} value={font}>
                  {font}
                </option>
              ))}
            </select>
          </div>
        </div>

        {/* Row 2: Upload Buttons */}
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <Button
            type="button"
            variant="outline"
            size="sm"
            className="w-full justify-center gap-2 border-border/60 text-xs"
            onClick={() => {}}
          >
            <Upload className="size-3.5" />
            <span>Upload WOFF2</span>
          </Button>

          <Button
            type="button"
            variant="outline"
            size="sm"
            className="w-full justify-center gap-2 border-border/60 text-xs"
            onClick={() => {}}
          >
            <Upload className="size-3.5" />
            <span>Upload WOFF2</span>
          </Button>
        </div>

        {/* Row 3: Sizes */}
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <div className="space-y-1.5">
            <label className="font-medium text-muted-foreground text-xs">Interface size</label>
            <Input
              type="number"
              value={interfaceSize}
              onChange={(e) => onChange({ interfaceSize: e.target.value })}
              className="h-9 border-border/60 bg-background/50 text-xs"
            />
          </div>

          <div className="space-y-1.5">
            <label className="font-medium text-muted-foreground text-xs">Monospace size</label>
            <Input
              type="number"
              value={monoSize}
              onChange={(e) => onChange({ monoSize: e.target.value })}
              className="h-9 border-border/60 bg-background/50 text-xs"
            />
          </div>
        </div>

        {/* Row 4: Line height & Font weight */}
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <div className="space-y-1.5">
            <label className="font-medium text-muted-foreground text-xs">Line height</label>
            <Input
              type="number"
              step="0.1"
              value={lineHeight}
              onChange={(e) => onChange({ lineHeight: e.target.value })}
              className="h-9 border-border/60 bg-background/50 text-xs"
            />
          </div>

          <div className="space-y-1.5">
            <label className="font-medium text-muted-foreground text-xs">Font weight</label>
            <Input
              type="number"
              step="50"
              value={fontWeight}
              onChange={(e) => onChange({ fontWeight: e.target.value })}
              className="h-9 border-border/60 bg-background/50 text-xs"
            />
          </div>
        </div>

        <p className="pt-1 text-[11px] text-muted-foreground/80">
          WOFF2 fonts travel with your theme. Up to 256 KiB per font.
        </p>
      </div>
    </section>
  );
}
