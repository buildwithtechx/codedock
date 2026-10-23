import { ExternalLink, Sparkles, X } from 'lucide-react';
import { useEffect, useState } from 'react';
import {
  InstanceUpgradeModal,
  OPEN_UPGRADE_MODAL_EVENT,
} from '#/features/settings/instance-upgrade-modal';
import { useGetUpdateStatus } from '#/hooks/use-settings';

const LS_LAST_SEEN = 'codedock_last_seen_version';
const LS_WHATS_NEW = 'codedock_whats_new_version';
const CHANGELOG_URL = 'https://github.com/buildwithtechx/codedock/releases';

export function UpdateCenter() {
  const { data: statusRes } = useGetUpdateStatus();
  const info = statusRes?.data;
  const currentVersion = info?.currentVersion ? info.currentVersion.replace(/^v/, '') : null;

  const [upgradeModalOpen, setUpgradeModalOpen] = useState(false);
  const [whatsNewVersion, setWhatsNewVersion] = useState<string | null>(null);
  const [notesOpen, setNotesOpen] = useState(false);

  useEffect(() => {
    const handleOpen = () => setUpgradeModalOpen(true);
    window.addEventListener(OPEN_UPGRADE_MODAL_EVENT, handleOpen);
    return () => window.removeEventListener(OPEN_UPGRADE_MODAL_EVENT, handleOpen);
  }, []);

  useEffect(() => {
    if (!currentVersion) return;
    try {
      const pendingWhatsNew = localStorage.getItem(LS_WHATS_NEW);
      if (pendingWhatsNew) {
        setWhatsNewVersion(pendingWhatsNew);
        return;
      }
      const lastSeen = localStorage.getItem(LS_LAST_SEEN);
      if (!lastSeen) {
        localStorage.setItem(LS_LAST_SEEN, currentVersion);
      } else if (lastSeen !== currentVersion) {
        setWhatsNewVersion(currentVersion);
      }
    } catch {
      /* ignore storage errors */
    }
  }, [currentVersion]);

  const dismissWhatsNew = () => {
    setWhatsNewVersion(null);
    try {
      localStorage.removeItem(LS_WHATS_NEW);
      if (currentVersion) {
        localStorage.setItem(LS_LAST_SEEN, currentVersion);
      }
    } catch {
      /* ignore */
    }
  };

  return (
    <>
      <InstanceUpgradeModal
        open={upgradeModalOpen}
        onOpenChange={(next) => {
          setUpgradeModalOpen(next);
          if (!next) {
            try {
              const pending = localStorage.getItem(LS_WHATS_NEW);
              if (pending) setWhatsNewVersion(pending);
            } catch {
              /* ignore */
            }
          }
        }}
      />

      {whatsNewVersion && (
        <div className="fixed inset-e-4 bottom-4 z-50 w-90 max-w-[calc(100vw-2rem)] overflow-hidden rounded-2xl border border-border/60 bg-popover/95 shadow-2xl backdrop-blur-xl">
          <div className="flex items-start gap-3 px-4 py-4">
            <div className="mt-0.5 flex size-9 shrink-0 items-center justify-center rounded-xl bg-primary/10 text-primary">
              <Sparkles className="size-4.5" />
            </div>
            <div className="min-w-0 flex-1">
              <p className="font-semibold text-[14px] text-foreground">
                Updated to v{whatsNewVersion}
              </p>
              <p className="mt-0.5 text-[12.5px] text-muted-foreground leading-relaxed">
                {info?.releaseNotes
                  ? 'See what changed in this release.'
                  : 'You are running the latest Codedock release.'}
              </p>
              <div className="mt-2.5 flex items-center gap-3">
                {info?.releaseNotes && (
                  <button
                    type="button"
                    onClick={() => setNotesOpen(true)}
                    className="font-medium text-[13px] text-foreground underline-offset-4 hover:underline"
                  >
                    What&apos;s new
                  </button>
                )}
                <a
                  href={info?.downloadUrl || CHANGELOG_URL}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="inline-flex items-center gap-1 font-medium text-[13px] text-foreground underline-offset-4 hover:underline"
                >
                  Changelog
                  <ExternalLink className="size-3.5" />
                </a>
              </div>
            </div>
            <button
              type="button"
              onClick={dismissWhatsNew}
              aria-label="Dismiss update notification"
              className="shrink-0 rounded-lg p-1 text-muted-foreground transition-colors hover:bg-foreground/6 hover:text-foreground"
            >
              <X className="size-4" />
            </button>
          </div>
        </div>
      )}

      {notesOpen && (
        <div
          className="fixed inset-0 z-60 flex items-center justify-center bg-black/50 p-4"
          onClick={() => setNotesOpen(false)}
        >
          <div
            className="flex max-h-[80vh] w-full max-w-lg flex-col overflow-hidden rounded-2xl border border-border/60 bg-card shadow-2xl"
            onClick={(e) => e.stopPropagation()}
          >
            <div className="flex items-center justify-between border-border/50 border-b px-5 py-4">
              <h3 className="font-semibold text-[15px] text-foreground">
                What&apos;s new in v{whatsNewVersion}
              </h3>
              <button
                type="button"
                onClick={() => setNotesOpen(false)}
                aria-label="Close release notes"
                className="rounded-lg p-1 text-muted-foreground transition-colors hover:bg-foreground/6 hover:text-foreground"
              >
                <X className="size-4" />
              </button>
            </div>
            <pre className="flex-1 overflow-auto whitespace-pre-wrap wrap-break-word px-5 py-4 font-mono text-[12.5px] text-muted-foreground leading-relaxed">
              {info?.releaseNotes?.trim() || 'No release notes provided for this build.'}
            </pre>
            <div className="border-border/50 border-t px-5 py-3">
              <a
                href={info?.downloadUrl || CHANGELOG_URL}
                target="_blank"
                rel="noopener noreferrer"
                className="inline-flex items-center gap-1 font-medium text-[13px] text-foreground underline-offset-4 hover:underline"
              >
                View full changelog
                <ExternalLink className="size-3.5" />
              </a>
            </div>
          </div>
        </div>
      )}
    </>
  );
}
