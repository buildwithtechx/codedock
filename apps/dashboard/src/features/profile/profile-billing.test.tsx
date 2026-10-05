import { cleanup, render } from '@testing-library/react';
import type { ComponentType } from 'react';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { Route } from '#/routes/_dashboard.profile';

const state = vi.hoisted(() => ({ cloudMode: undefined as boolean | undefined, billing: vi.fn() }));
vi.mock('#/components/layout/page-header', () => ({ PageHeader: () => null }));
vi.mock('#/features/settings', () => ({
  useGetPublicSettings: () => ({ data: { data: { cloudMode: state.cloudMode } } }),
}));
vi.mock('#/features/profile/security-2fa-setup', () => ({ Security2FASetup: () => null }));
vi.mock('#/features/profile/user-profile-form', () => ({
  ProfileEmailForm: () => null,
  ProfileNameForm: () => null,
  ProfilePasswordForm: () => null,
}));
vi.mock('#/features/profile/billing-section', () => ({
  BillingSection: () => {
    state.billing();
    return <div>Cloud billing</div>;
  },
}));

afterEach(() => {
  cleanup();
  vi.clearAllMocks();
});
describe('profile billing visibility', () => {
  for (const cloudMode of [undefined, false, true]) {
    it(`only mounts billing when the server reports cloud mode (${cloudMode})`, () => {
      state.cloudMode = cloudMode;
      const Profile = Route.options.component as ComponentType;
      const view = render(<Profile />);
      expect(view.queryByText('Cloud billing') !== null).toBe(cloudMode === true);
      expect(state.billing).toHaveBeenCalledTimes(cloudMode === true ? 1 : 0);
    });
  }
});
