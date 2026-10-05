export function StatusBadge({ status }: { status: string }) {
  const s = (status || '').toLowerCase();
  let color = 'bg-gray-500/10 text-gray-500 border-gray-500/20';
  let dot = 'bg-gray-500';

  if (s === 'running' || s === 'online' || s === 'healthy') {
    color = 'bg-emerald-500/10 text-emerald-500 border-emerald-500/20';
    dot = 'bg-emerald-500 shadow-[0_0_8px_rgba(16,185,129,0.4)]';
  } else if (s === 'failed' || s === 'error' || s === 'stopped') {
    color = 'bg-red-500/10 text-red-500 border-red-500/20';
    dot = 'bg-red-500';
  } else if (s === 'deploying' || s === 'pending' || s === 'building') {
    color = 'bg-amber-500/10 text-amber-500 border-amber-500/20';
    dot = 'bg-amber-500 animate-pulse';
  }

  return (
    <div
      className={`flex items-center gap-1.5 rounded-full border px-2.5 py-0.5 font-medium text-xs ${color}`}
    >
      <div className={`h-1.5 w-1.5 rounded-full ${dot}`} />
      <span className="capitalize">{status || 'Unknown'}</span>
    </div>
  );
}
