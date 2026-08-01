import { createFileRoute, Outlet } from '@tanstack/react-router';
import { FaDiscord, FaGithub, FaXTwitter } from 'react-icons/fa6';
import { OnboardingImport } from '#/features/auth';
import { useSystemStore } from '#/stores/system-store';

export const Route = createFileRoute('/_auth')({ component: AuthLayout });

function AuthLayout() {
  const siteName = useSystemStore((state) => state.siteName);

  return (
    <main className="grid min-h-dvh bg-background text-foreground lg:grid-cols-[minmax(0,0.94fr)_minmax(34rem,1.06fr)]">
      <section className="relative hidden overflow-hidden bg-[#36156e] px-10 py-9 text-white lg:flex lg:flex-col xl:px-16 xl:py-12">
        <div
          aria-hidden
          className="absolute inset-0 bg-[radial-gradient(circle_at_86%_16%,rgba(196,161,255,0.3),transparent_25%),radial-gradient(circle_at_15%_82%,rgba(111,45,214,0.6),transparent_34%)]"
        />
        <div
          aria-hidden
          className="absolute inset-y-0 right-0 w-[48%] bg-[linear-gradient(90deg,transparent_0%,rgba(255,255,255,0.04)_1px,transparent_1px),linear-gradient(0deg,transparent_0%,rgba(255,255,255,0.04)_1px,transparent_1px)] bg-size-[2rem_2rem]"
        />
        <div className="relative z-10 flex items-center gap-3">
          <img src="/apple-touch-icon.png" alt="" className="size-8 rounded-lg" />
          <div className="flex h-8 flex-col justify-center">
            <p className="font-bold text-base leading-none tracking-[-0.04em]">{siteName}</p>
            <p className="mt-1 text-[9px] text-white/55 uppercase leading-none tracking-[0.16em]">
              Self-hosted workspace
            </p>
          </div>
        </div>

        <div className="relative z-10 my-auto max-w-md py-16">
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
          <h2 className="relative mt-5 max-w-sm font-bold text-4xl leading-[0.98] tracking-[-0.045em] xl:text-5xl">
            Build it.
            <br />
            Run it your way.
          </h2>
          <p className="relative mt-6 max-w-sm text-sm text-white/70 leading-6">
            Codedock keeps services, releases, and infrastructure in the workspace your team
            operates.
          </p>
          <div className="relative mt-10 flex max-w-sm items-center gap-3">
            <span className="h-px flex-1 bg-white/25" />
            <span className="size-2 rounded-full bg-[#d8c7ff] shadow-[0_0_18px_rgba(216,199,255,0.9)]" />
            <span className="h-px flex-[2] bg-white/15" />
          </div>
        </div>

        <div className="relative z-10">
          <OnboardingImport />
        </div>
      </section>

      <section className="flex min-h-dvh flex-col px-5 py-6 sm:px-10 sm:py-9 lg:px-14 xl:px-20">
        <div className="flex items-center gap-3 lg:hidden">
          <img src="/apple-touch-icon.png" alt="" className="size-7 rounded-md" />
          <p className="font-bold tracking-[-0.04em]">{siteName}</p>
        </div>
        <div className="flex flex-1 items-center justify-center py-12 lg:py-16">
          <Outlet />
        </div>
        <div className="flex justify-end gap-2 text-muted-foreground">
          <a
            href="https://github.com/buildwithtechx/codedock"
            target="_blank"
            rel="noreferrer"
            aria-label="Codedock on GitHub"
            className="grid size-8 place-items-center rounded-lg transition-colors hover:bg-primary-soft hover:text-primary"
          >
            <FaGithub className="size-4" aria-hidden="true" />
          </a>
          <a
            href="https://x.com/codedock"
            target="_blank"
            rel="noreferrer"
            aria-label="Codedock on X"
            className="grid size-8 place-items-center rounded-lg transition-colors hover:bg-primary-soft hover:text-primary"
          >
            <FaXTwitter className="size-3.5" aria-hidden="true" />
          </a>
          <a
            href="https://discord.gg/codedock"
            target="_blank"
            rel="noreferrer"
            aria-label="Codedock on Discord"
            className="grid size-8 place-items-center rounded-lg transition-colors hover:bg-primary-soft hover:text-primary"
          >
            <FaDiscord className="size-4" aria-hidden="true" />
          </a>
        </div>
      </section>
    </main>
  );
}
