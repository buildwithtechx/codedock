/**
 * Registers the Codedock PWA service worker for app-shell caching,
 * offline fallback, and Web Push notifications.
 */
export function registerServiceWorker(): void {
  if (typeof window === 'undefined' || !('serviceWorker' in navigator)) {
    return;
  }

  // Register on production builds or when explicitly enabled in dev
  const isProd = import.meta.env.PROD;
  const isExplicitlyEnabled =
    typeof localStorage !== 'undefined' && localStorage.getItem('codedock:enable-sw') === 'true';

  if (!isProd && !isExplicitlyEnabled) {
    return;
  }

  window.addEventListener('load', () => {
    navigator.serviceWorker
      .register('/sw.js', { scope: '/' })
      .then((registration) => {
        registration.addEventListener('updatefound', () => {
          const installingWorker = registration.installing;
          if (installingWorker) {
            installingWorker.addEventListener('statechange', () => {
              if (installingWorker.state === 'installed' && navigator.serviceWorker.controller) {
                // New content is available; send skip waiting message if needed
                installingWorker.postMessage({ action: 'SKIP_WAITING' });
              }
            });
          }
        });
      })
      .catch((error) => {
        console.warn('Service worker registration failed:', error);
      });
  });
}
