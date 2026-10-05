import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, describe, expect, it } from 'vitest';
import { MobileNavigation } from './mobile-navigation';

afterEach(cleanup);
describe('mobile navigation', () => {
  it('opens the menu and restores trigger focus after Escape', async () => {
    render(<MobileNavigation />);
    const trigger = screen.getByRole('button', { name: 'Open menu' });
    fireEvent.click(trigger);
    expect(screen.getByRole('dialog', { name: 'Explore Codedock' })).toBeTruthy();
    expect(screen.getByRole('link', { name: 'Application deployment' })).toBeTruthy();
    fireEvent.keyDown(document, { key: 'Escape' });
    await waitFor(() => expect(screen.queryByRole('dialog')).toBeNull());
    expect(document.activeElement).toBe(trigger);
  });
  it('closes when a destination is selected', async () => {
    render(<MobileNavigation />);
    fireEvent.click(screen.getByRole('button', { name: 'Open menu' }));
    const link = screen.getByRole('link', { name: 'Agencies' });
    expect(link.getAttribute('href')).toBe('/solutions/agencies');
    link.setAttribute('href', '#selected');
    fireEvent.click(link);
    await waitFor(() => expect(screen.queryByRole('dialog')).toBeNull());
  });
});
