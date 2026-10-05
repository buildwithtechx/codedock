import {
  createMemoryHistory,
  createRootRoute,
  createRoute,
  createRouter,
  RouterProvider,
} from '@tanstack/react-router';
import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, describe, expect, it } from 'vitest';
import { MobileNavigation } from './mobile-navigation';

async function renderMenu() {
  const root = createRootRoute({ component: MobileNavigation });
  const index = createRoute({ getParentRoute: () => root, path: '/' });
  const agencies = createRoute({ getParentRoute: () => root, path: '/solutions/agencies' });
  const router = createRouter({
    routeTree: root.addChildren([index, agencies]),
    history: createMemoryHistory({ initialEntries: ['/'] }),
  });
  await router.load();
  render(<RouterProvider router={router} />);
  return router;
}
afterEach(cleanup);
describe('mobile navigation', () => {
  it('opens the menu and restores trigger focus after Escape', async () => {
    await renderMenu();
    const trigger = screen.getByRole('button', { name: 'Open menu' });
    fireEvent.click(trigger);
    expect(screen.getByRole('dialog', { name: 'Explore Codedock' })).toBeTruthy();
    expect(screen.getByRole('link', { name: 'Application deployment' })).toBeTruthy();
    fireEvent.keyDown(document, { key: 'Escape' });
    await waitFor(() => expect(screen.queryByRole('dialog')).toBeNull());
    expect(document.activeElement).toBe(trigger);
  });
  it('closes when a destination is selected', async () => {
    const router = await renderMenu();
    fireEvent.click(screen.getByRole('button', { name: 'Open menu' }));
    const link = screen.getByRole('link', { name: 'Agencies' });
    expect(link.getAttribute('href')).toBe('/solutions/agencies');
    fireEvent.click(link);
    await waitFor(() => expect(screen.queryByRole('dialog')).toBeNull());
    await waitFor(() => expect(router.state.location.pathname).toBe('/solutions/agencies'));
  });
});
