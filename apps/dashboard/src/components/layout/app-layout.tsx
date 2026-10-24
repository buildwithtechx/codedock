import { useMutation } from '@tanstack/react-query';
import { useRouterState } from '@tanstack/react-router';
import { AlertCircle, ArrowUp } from 'lucide-react';
import type * as React from 'react';
import { useEffect, useRef, useState } from 'react';
import { toast } from 'sonner';
import { Button } from '@codedock/ui/button';
import { OrganizationOnboarding } from '#/features/organizations';
import { apiClient } from '#/lib/api-client';
import { useAiSidebarStore } from '#/stores/ai-sidebar-store';
import { useAuthStore } from '#/stores/auth-store';
import { AskAiSidebar } from '../ai/ask-ai-sidebar';
import { AppSidebar } from './app-sidebar';
import { BottomDock } from './bottom-dock';
import { CommandPalette } from './command-palette';
import { TopHeader } from './top-header';
import { UpdateCenter } from './update-center';

export function AppLayout({ children }: { children: React.ReactNode }) {
  const pathname = useRouterState({
    select: (state) => state.location.pathname,
  });
  const isAiSidebarOpen = useAiSidebarStore((s) => s.isOpen);
  const [sidebarCollapsed, setSidebarCollapsed] = useState(() => {
    if (typeof window === 'undefined') return false;
    return localStorage.getItem('codedock_auto_collapse_sidebar') === 'true';
  });
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const [commandOpen, setCommandOpen] = useState(false);
  const mainRef = useRef<HTMLElement>(null);

  useEffect(() => {
    setMobileMenuOpen(false);
    mainRef.current?.scrollTo({ top: 0, behavior: 'instant' });
  }, [pathname]);

  useEffect(() => {
    const prefersAlwaysCollapsed =
      typeof window !== 'undefined' &&
      localStorage.getItem('codedock_auto_collapse_sidebar') === 'true';
    if (isAiSidebarOpen) {
      setSidebarCollapsed(true);
    } else if (!prefersAlwaysCollapsed) {
      setSidebarCollapsed(false);
    }
  }, [isAiSidebarOpen]);

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
    <div className="relative flex h-dvh min-h-0 flex-col overflow-hidden bg-background md:flex-row">
      <AppSidebar
        collapsed={sidebarCollapsed}
        onToggle={() => setSidebarCollapsed((p) => !p)}
        mobileOpen={mobileMenuOpen}
        onMobileClose={() => setMobileMenuOpen(false)}
        onOpenCommand={() => setCommandOpen(true)}
      />
      <div className="relative flex min-h-0 min-w-0 flex-1 flex-col bg-background text-foreground">
        <TopHeader
          mobileOpen={mobileMenuOpen}
          onToggleMobile={() => setMobileMenuOpen((p) => !p)}
          onOpenCommand={() => setCommandOpen(true)}
        />
        <main
          ref={mainRef}
          className="dashboard-content relative min-h-0 min-w-0 flex-1 overflow-y-auto overscroll-y-contain bg-background p-4 text-foreground sm:p-5 lg:p-6"
        >
          <div id="top" aria-hidden="true" className="h-0 w-0 pointer-events-none" />
          <div key={pathname} className="page-transition mx-auto w-full max-w-[1600px]">
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
            <div className="mt-8 flex justify-center pb-8">
              <button
                type="button"
                onClick={() => mainRef.current?.scrollTo({ top: 0, behavior: 'smooth' })}
                className="sr-only focus:not-sr-only focus:inline-flex focus:items-center focus:gap-1.5 focus:rounded-lg focus:border focus:border-border/50 focus:bg-card focus:px-3 focus:py-1.5 focus:font-medium focus:text-foreground focus:text-xs focus:shadow-sm"
                aria-label="Back to top"
              >
                <ArrowUp className="h-3.5 w-3.5" />
                Back to top
              </button>
            </div>
          </div>
        </main>
      </div>
      <AskAiSidebar />
      <BottomDock />
      <CommandPalette open={commandOpen} onOpenChange={setCommandOpen} />
      <OrganizationOnboarding />
      <UpdateCenter />
    </div>
  );
}
