import {
  createMemoryHistory,
  createRootRoute,
  createRoute,
  createRouter,
  RouterProvider,
} from '@tanstack/react-router';
import { cleanup, fireEvent, render, screen } from '@testing-library/react';
import { afterEach, describe, expect, it } from 'vitest';
import { TemplateCatalog } from './template-catalog';

afterEach(cleanup);

async function renderCatalog() {
  const rootRoute = createRootRoute();
  const templatesRoute = createRoute({
    getParentRoute: () => rootRoute,
    path: '/templates',
    component: TemplateCatalog,
  });
  const detailRoute = createRoute({
    getParentRoute: () => rootRoute,
    path: '/templates/$templateId',
    component: () => <div />,
  });
  const routeTree = rootRoute.addChildren([templatesRoute, detailRoute]);
  const router = createRouter({
    routeTree,
    history: createMemoryHistory({ initialEntries: ['/templates'] }),
  });
  await router.load();
  render(<RouterProvider router={router} />);
}

describe('template discovery', () => {
  it('combines a normalized search with a category filter', async () => {
    await renderCatalog();
    fireEvent.change(screen.getByRole('searchbox'), { target: { value: '  REDIS  ' } });
    expect(screen.getByRole('heading', { name: 'Redis' })).toBeTruthy();
    fireEvent.click(screen.getByRole('button', { name: 'AI' }));
    expect(screen.queryByRole('heading', { name: 'Redis' })).toBeNull();
    expect(screen.getByText('No recipes match your search.')).toBeTruthy();
    fireEvent.click(screen.getByRole('button', { name: 'Reset filters' }));
    expect(screen.getByRole('heading', { name: 'Redis' })).toBeTruthy();
    expect(screen.getByRole('heading', { name: 'Ollama' })).toBeTruthy();
  });
  it('takes a catalogued recipe to its detail page', async () => {
    await renderCatalog();
    fireEvent.click(screen.getByRole('button', { name: 'Storage' }));
    const link = screen.getByRole('link', { name: /MinIO/ });
    expect(link.getAttribute('href')).toBe('/templates/minio');
    expect(screen.queryByRole('heading', { name: 'PostgreSQL' })).toBeNull();
  });
  it('keeps uncatalogued recipes on their setup guides', async () => {
    await renderCatalog();
    fireEvent.change(screen.getByRole('searchbox'), { target: { value: 'Redis' } });
    const link = screen.getByRole('link', { name: /Redis/ });
    expect(link.getAttribute('href')).toBe(
      'https://docs.codedock.run/deployments/templates/#redis'
    );
  });
  it('finds recipes by their deployment workflow', async () => {
    await renderCatalog();
    fireEvent.change(screen.getByRole('searchbox'), { target: { value: 'Compose' } });
    expect(screen.getByRole('heading', { name: 'Supabase' })).toBeTruthy();
    expect(screen.queryByRole('heading', { name: 'Redis' })).toBeNull();
  });
});
