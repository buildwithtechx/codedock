interface ThemePreviewProps {
  theme: 'light' | 'dim' | 'dark';
}

export function ThemePreview({ theme }: ThemePreviewProps) {
  const isLight = theme === 'light';
  const isDim = theme === 'dim';

  const bgPage = isLight ? 'bg-zinc-100' : isDim ? 'bg-[#181a1f]' : 'bg-zinc-950';
  const bgCard = isLight ? 'bg-white' : isDim ? 'bg-[#22262e]' : 'bg-zinc-900';
  const borderCol = isLight ? 'border-zinc-200' : isDim ? 'border-[#2f3542]' : 'border-zinc-800';
  const textTitle = isLight ? 'bg-zinc-900' : 'bg-zinc-100';
  const textMuted = isLight ? 'bg-zinc-300' : isDim ? 'bg-zinc-600' : 'bg-zinc-700';

  return (
    <div
      aria-hidden="true"
      className={`flex h-24 gap-2 overflow-hidden rounded-lg border ${borderCol} ${bgPage} p-2.5 transition-colors`}
    >
      <div className={`flex w-1/4 shrink-0 flex-col gap-1.5 rounded-md ${bgCard} p-1.5`}>
        <span className={`mb-1 size-2 shrink-0 rounded-full ${textTitle}`} />
        <span className={`h-1.5 rounded-full ${textMuted}`} />
        <span className={`h-1 w-3/4 rounded-full ${textMuted} opacity-60`} />
        <span className={`h-1 w-3/4 rounded-full ${textMuted} opacity-60`} />
      </div>
      <div className="flex min-w-0 flex-1 flex-col gap-2 py-1">
        <span className={`h-1.5 w-3/5 shrink-0 rounded-full ${textMuted}`} />
        <div className="grid flex-1 grid-cols-2 gap-1.5">
          <span className={`rounded-md ${bgCard}`} />
          <span className={`rounded-md ${bgCard}`} />
        </div>
        <span className={`h-4 shrink-0 rounded-md ${bgCard}`} />
      </div>
    </div>
  );
}
