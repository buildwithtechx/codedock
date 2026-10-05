export function trackInstallCopy(source: string, command: string) {
  if (!window.posthog?.capture) {
    window.__codedock_install_copies ??= [];
    if (window.__codedock_install_copies.length < 100) {
      window.__codedock_install_copies.push({ source, command });
    }
    return;
  }
  try {
    window.posthog.capture('install_command_copied', { source, command });
  } catch {
    return;
  }
}

export function flushInstallCopies() {
  if (!window.posthog?.capture) return;
  const copies = window.__codedock_install_copies ?? [];
  window.__codedock_install_copies = [];
  for (const { source, command } of copies) {
    trackInstallCopy(source, command);
  }
}
