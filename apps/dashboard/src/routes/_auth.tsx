import { createFileRoute, Outlet } from '@tanstack/react-router';
import { useSystemStore } from '#/stores/system-store';

export const Route = createFileRoute('/_auth')({ component: AuthLayout });

function AuthLayout() {
  const siteName = useSystemStore((state) => state.siteName);

  return (
    <main className="grid min-h-dvh bg-background text-foreground lg:grid-cols-[minmax(0,0.94fr)_minmax(34rem,1.06fr)]">
      <section className="relative hidden overflow-hidden bg-[#36156e] px-10 py-9 text-white lg:flex lg:flex-col xl:px-16 xl:py-12">
        <div className="flex items-center gap-3">
          <span className="grid size-10 place-items-center rounded-xl border border-white/25 bg-white/10 font-extrabold text-xs tracking-[-0.1em]">
            CD
          </span>
          <div>
            <p className="font-bold text-base tracking-[-0.04em]">{siteName}</p>
            <p className="mt-0.5 text-[10px] text-white/65 uppercase tracking-[0.18em]">
              Deployment workspace
            </p>
          </div>
        </div>

        <div className="relative my-auto max-w-md py-16">
          <div
            aria-hidden
            className="absolute -top-8 -left-14 size-56 rounded-full border border-[#c6a8ff]/40"
          />
          <div
            aria-hidden
            className="absolute top-10 left-11 size-48 rounded-full border border-white/15"
          />
          <div
            aria-hidden
            className="absolute top-24 left-24 size-24 rounded-full bg-[#d8c7ff]/15"
          />
          <p className="relative font-semibold text-[#d8c7ff] text-[10px] uppercase tracking-[0.18em]">
            Built for your stack
          </p>
          <h2 className="relative mt-5 max-w-sm font-extrabold text-5xl leading-[0.95] tracking-[-0.06em] xl:text-6xl">
            Run what
            <br />
            you build.
            <br />
            Keep it yours.
          </h2>
          <p className="relative mt-6 max-w-sm text-sm text-white/70 leading-6">
            Codedock keeps services, releases, and infrastructure in the workspace your team
            operates.
          </p>
        </div>

        <div className="border-white/15 border-t pt-5 text-[10px] text-white/60 uppercase tracking-[0.14em]">
          Infrastructure on your terms
        </div>
      </section>

      <section className="flex min-h-dvh flex-col px-5 py-6 sm:px-10 sm:py-9 lg:px-14 xl:px-20">
        <div className="flex items-center gap-3 lg:hidden">
          <span className="grid size-9 place-items-center rounded-lg bg-primary font-extrabold text-[10px] text-primary-foreground tracking-[-0.1em]">
            CD
          </span>
          <p className="font-bold tracking-[-0.04em]">{siteName}</p>
        </div>
        <div className="flex flex-1 items-center justify-center py-12 lg:py-16">
          <Outlet />
        </div>
        <p className="text-[10px] text-muted-foreground uppercase tracking-[0.14em]">
          Codedock deployment workspace
        </p>
      </section>
    </main>
  );
}
