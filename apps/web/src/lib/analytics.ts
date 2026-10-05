export function trackInstallCopy(source: string, command: string) {
  try {
    window.posthog?.capture?.('install_command_copied', { source, command });
  } catch {
    return;
  }
}
