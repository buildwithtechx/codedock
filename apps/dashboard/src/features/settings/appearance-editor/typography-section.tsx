import { Button } from '@codedock/ui/button';
import { Input } from '@codedock/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@codedock/ui/select';
import { Upload } from 'lucide-react';
import { useRef, useState } from 'react';
import { toast } from 'sonner';

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

const DEFAULT_INTERFACE_FONTS = ['Gellix', 'Manrope', 'Inter', 'Geist', 'System', 'Custom'];
const DEFAULT_MONO_FONTS = ['JetBrains Mono', 'Fira Code', 'Geist Mono', 'System Mono', 'Custom'];

export function TypographySection({
  interfaceFont = 'Manrope',
  monoFont = 'JetBrains Mono',
  interfaceSize = 16,
  monoSize = 13,
  lineHeight = 1.5,
  fontWeight = 400,
  onChange,
}: TypographySectionProps) {
  const [interfaceFonts, setInterfaceFonts] = useState<string[]>(() => {
    return DEFAULT_INTERFACE_FONTS.includes(interfaceFont)
      ? DEFAULT_INTERFACE_FONTS
      : [interfaceFont, ...DEFAULT_INTERFACE_FONTS];
  });

  const [monoFonts, setMonoFonts] = useState<string[]>(() => {
    return DEFAULT_MONO_FONTS.includes(monoFont)
      ? DEFAULT_MONO_FONTS
      : [monoFont, ...DEFAULT_MONO_FONTS];
  });

  const interfaceInputRef = useRef<HTMLInputElement>(null);
  const monoInputRef = useRef<HTMLInputElement>(null);

  const handleFontUpload = async (role: 'sans' | 'mono', file?: File) => {
    if (!file) return;
    if (file.size > 256 * 1024) {
      toast.error('Font file size must be 256 KiB or less');
      return;
    }

    try {
      const buffer = await file.arrayBuffer();
      const header = new TextDecoder().decode(new Uint8Array(buffer.slice(0, 4)));
      if (header !== 'wOF2') {
        toast.error('Please upload a valid WOFF2 font file');
        return;
      }

      const cleanName = file.name.replace(/\.[^/.]+$/, '').slice(0, 32);
      const fontName = `Custom-${cleanName}`;

      if (typeof FontFace !== 'undefined') {
        const fontFace = new FontFace(fontName, buffer);
        await fontFace.load();
        document.fonts.add(fontFace);
      }

      if (role === 'sans') {
        setInterfaceFonts((prev) => (prev.includes(fontName) ? prev : [fontName, ...prev]));
        onChange({ interfaceFont: fontName });
      } else {
        setMonoFonts((prev) => (prev.includes(fontName) ? prev : [fontName, ...prev]));
        onChange({ monoFont: fontName });
      }

      toast.success(`Loaded custom ${role === 'sans' ? 'interface' : 'mono'} font "${cleanName}"`);
    } catch {
      toast.error('Failed to parse or load WOFF2 font file');
    }
  };

  return (
    <section className="rounded-2xl border border-border/60 bg-card p-5 shadow-xs">
      <h3 className="font-semibold text-foreground text-sm">Typography</h3>

      {/* Hidden file inputs for real WOFF2 font loading */}
      <input
        ref={interfaceInputRef}
        type="file"
        accept=".woff2,font/woff2"
        className="hidden"
        onChange={(e) => {
          void handleFontUpload('sans', e.target.files?.[0]);
          e.target.value = '';
        }}
      />
      <input
        ref={monoInputRef}
        type="file"
        accept=".woff2,font/woff2"
        className="hidden"
        onChange={(e) => {
          void handleFontUpload('mono', e.target.files?.[0]);
          e.target.value = '';
        }}
      />

      <div className="mt-4 space-y-4">
        {/* Row 1: Fonts */}
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <div className="space-y-1.5">
            <span className="font-medium text-muted-foreground text-xs">Interface font</span>
            <Select value={interfaceFont} onValueChange={(val) => onChange({ interfaceFont: val })}>
              <SelectTrigger className="h-9 bg-background/50 text-xs">
                <SelectValue placeholder="Select interface font" />
              </SelectTrigger>
              <SelectContent>
                {interfaceFonts.map((font) => (
                  <SelectItem key={font} value={font}>
                    {font}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-1.5">
            <span className="font-medium text-muted-foreground text-xs">Monospace font</span>
            <Select value={monoFont} onValueChange={(val) => onChange({ monoFont: val })}>
              <SelectTrigger className="h-9 bg-background/50 font-mono text-xs">
                <SelectValue placeholder="Select mono font" />
              </SelectTrigger>
              <SelectContent>
                {monoFonts.map((font) => (
                  <SelectItem key={font} value={font} className="font-mono">
                    {font}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        </div>

        {/* Row 2: Upload Buttons */}
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <Button
            type="button"
            variant="outline"
            size="sm"
            className="w-full justify-center gap-2 border-border/60 text-xs"
            onClick={() => interfaceInputRef.current?.click()}
          >
            <Upload className="size-3.5" />
            <span>Upload Interface WOFF2</span>
          </Button>

          <Button
            type="button"
            variant="outline"
            size="sm"
            className="w-full justify-center gap-2 border-border/60 text-xs"
            onClick={() => monoInputRef.current?.click()}
          >
            <Upload className="size-3.5" />
            <span>Upload Mono WOFF2</span>
          </Button>
        </div>

        {/* Row 3: Sizes */}
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <div className="space-y-1.5">
            <span className="font-medium text-muted-foreground text-xs">Interface size</span>
            <Input
              type="number"
              value={interfaceSize}
              onChange={(e) => onChange({ interfaceSize: e.target.value })}
              className="h-9 border-border/60 bg-background/50 text-xs"
            />
          </div>

          <div className="space-y-1.5">
            <span className="font-medium text-muted-foreground text-xs">Monospace size</span>
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
            <span className="font-medium text-muted-foreground text-xs">Line height</span>
            <Input
              type="number"
              step="0.1"
              value={lineHeight}
              onChange={(e) => onChange({ lineHeight: e.target.value })}
              className="h-9 border-border/60 bg-background/50 text-xs"
            />
          </div>

          <div className="space-y-1.5">
            <span className="font-medium text-muted-foreground text-xs">Font weight</span>
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
