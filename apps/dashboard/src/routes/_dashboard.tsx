import { createFileRoute, Outlet, redirect } from '@tanstack/react-router';
import { AppLayout } from '#/components/layout/app-layout';
import { useAuthStore } from '#/stores/auth-store';

export const Route = createFileRoute('/_dashboard')({
  beforeLoad: () => {
    const { isAuthenticated } = useAuthStore.getState();
    if (!isAuthenticated) {
      throw redirect({
        to: '/signin',
      });
    }
  },
  component: DashboardLayout,
});

function DashboardLayout() {
  return (
    <AppLayout>
      <Outlet />
    </AppLayout>
  );
}
