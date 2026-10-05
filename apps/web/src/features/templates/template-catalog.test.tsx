import { cleanup, fireEvent, render, screen } from '@testing-library/react';
import { afterEach, describe, expect, it } from 'vitest';
import { TemplateCatalog } from './template-catalog';

afterEach(cleanup);
describe('template discovery', () => {
  it('combines a normalized search with a category filter', () => {
    render(<TemplateCatalog />);
    fireEvent.change(screen.getByRole('searchbox'), { target: { value: '  REDIS  ' } });
    expect(screen.getByRole('heading', { name: 'Redis' })).toBeTruthy();
    fireEvent.click(screen.getByRole('button', { name: 'AI' }));
    expect(screen.queryByRole('heading', { name: 'Redis' })).toBeNull();
    expect(screen.getByText('No recipes match your search.')).toBeTruthy();
    fireEvent.click(screen.getByRole('button', { name: 'Reset filters' }));
    expect(screen.getByRole('heading', { name: 'Redis' })).toBeTruthy();
    expect(screen.getByRole('heading', { name: 'Ollama' })).toBeTruthy();
  });
  it('takes a selected recipe to its setup instructions', () => {
    render(<TemplateCatalog />);
    fireEvent.click(screen.getByRole('button', { name: 'Storage' }));
    const link = screen.getByRole('link', { name: /MinIO/ });
    expect(link.getAttribute('href')).toBe(
      'https://docs.codedock.run/deployments/templates/#minio'
    );
    expect(screen.queryByRole('heading', { name: 'PostgreSQL' })).toBeNull();
  });
  it('finds recipes by their deployment workflow', () => {
    render(<TemplateCatalog />);
    fireEvent.change(screen.getByRole('searchbox'), { target: { value: 'Compose' } });
    expect(screen.getByRole('heading', { name: 'Supabase' })).toBeTruthy();
    expect(screen.queryByRole('heading', { name: 'Redis' })).toBeNull();
  });
});
