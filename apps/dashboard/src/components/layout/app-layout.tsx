import { useMutation } from '@tanstack/react-query';
import { useRouterState } from '@tanstack/react-router';
import { AlertCircle, Menu, X } from 'lucide-react';
import type * as React from 'react';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { OrganizationOnboarding } from '#/features/organizations';
import { apiClient } from '#/lib/api-client';
import { useAuthStore } from '#/stores/auth-store';
import { Button } from '@/components/ui/button';
import { AppSidebar } from './app-sidebar';
import { CommandPalette } from './command-palette';

export function AppLayout({ children }: { children: React.ReactNode }) {
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const [commandOpen, setCommandOpen] = useState(false);
  const pathname = useRouterState({
    select: (state) => state.location.pathname,
  });

  useEffect(() => {
    setMobileMenuOpen(false);
  }, [pathname]);

  const { user } = useAuthStore();
  const [isResending, setIsResending] = useState(false);

  const resendMutation = useMutation({
    mutationFn: async () => {
      return apiClient.post('/auth/email/resend', { email: user?.email });
    },
    onSuccess: () => {
      toast.success('Verification email sent! Please check your inbox.');
    },
    onError: () => {
      toast.error('Failed to send verification email. Please try again later.');
    },
    onSettled: () => {
      setIsResending(false);
    },
  });

  const handleResend = () => {
    setIsResending(true);
    resendMutation.mutate();
  };

  return (
    <div className="relative min-h-screen bg-background md:p-3">
      <AppSidebar
        collapsed={sidebarCollapsed}
        onToggle={() => setSidebarCollapsed((p) => !p)}
        mobileOpen={mobileMenuOpen}
        onMobileClose={() => setMobileMenuOpen(false)}
      />
      <div
        className={`relative flex min-h-screen flex-1 flex-col ${sidebarCollapsed ? 'md:pl-24' : 'md:pl-[296px]'}`}
      >
        <div className="flex h-14 items-center px-4 md:hidden">
          <button
            type="button"
            onClick={() => setMobileMenuOpen((p) => !p)}
            className="flex h-9 w-9 items-center justify-center rounded-xl border border-border/60 text-muted-foreground transition-colors hover:bg-muted"
          >
            {mobileMenuOpen ? <X className="h-4 w-4" /> : <Menu className="h-4 w-4" />}
          </button>
        </div>
        <main className="dashboard-content flex-1 overflow-auto px-4 py-6 md:px-8 md:py-8">
          <div key={pathname} className="page-transition mx-auto w-full max-w-[90rem]">
            {user && user.emailVerified === false && (
              <div className="mb-6 flex flex-col items-center justify-between gap-4 rounded-lg border border-warning/50 bg-warning/20 p-4 sm:flex-row">
                <div className="flex items-center gap-3">
                  <AlertCircle className="h-5 w-5 text-warning" />
                  <p className="font-medium text-sm">
                    Please verify your email address to unlock all features.
                  </p>
                </div>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={handleResend}
                  disabled={isResending}
                  className="whitespace-nowrap"
                >
                  {isResending ? 'Sending...' : 'Resend Email'}
                </Button>
              </div>
            )}
            {children}
          </div>
        </main>
      </div>
      <CommandPalette open={commandOpen} onOpenChange={setCommandOpen} />
      <OrganizationOnboarding />
    </div>
  );
}
