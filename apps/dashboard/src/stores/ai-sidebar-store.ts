import { create } from 'zustand';
import type {
  AiConversationThread,
  AiPendingAction,
  AiResourceTag,
  AiToolStep,
  AiTurnMessage,
} from '#/interfaces/ai';

export type { AiConversationThread, AiPendingAction, AiResourceTag, AiToolStep, AiTurnMessage };

interface AiSidebarState {
  isOpen: boolean;
  isExpanded: boolean;
  mode: 'ask' | 'agent';
  threads: AiConversationThread[];
  activeThreadId: string;
  openSidebar: () => void;
  closeSidebar: () => void;
  toggleSidebar: () => void;
  toggleExpanded: () => void;
  setMode: (mode: 'ask' | 'agent') => void;
  startNewThread: () => string;
  selectThread: (id: string) => void;
  deleteThread: (id: string) => void;
  appendMessage: (msg: AiTurnMessage) => void;
  updateActionStatus: (
    messageId: string,
    actionId: string,
    status: 'approved' | 'cancelled' | 'failed'
  ) => void;
}

const STORAGE_KEY = 'codedock_ai_threads_v1';

function createEmptyThread(): AiConversationThread {
  return {
    id: `thread_${Date.now().toString(36)}`,
    title: 'New conversation',
    updatedAt: new Date().toISOString(),
    messages: [],
  };
}

function loadThreads(): AiConversationThread[] {
  if (typeof window === 'undefined') return [createEmptyThread()];
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (!raw) return [createEmptyThread()];
    const parsed = JSON.parse(raw) as AiConversationThread[];
    return Array.isArray(parsed) && parsed.length > 0 ? parsed : [createEmptyThread()];
  } catch {
    return [createEmptyThread()];
  }
}

function saveThreads(threads: AiConversationThread[]) {
  if (typeof window === 'undefined') return;
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(threads.slice(0, 15)));
  } catch {}
}

const initialThreads = loadThreads();

export const useAiSidebarStore = create<AiSidebarState>((set) => ({
  isOpen: false,
  isExpanded: false,
  mode: 'ask',
  threads: initialThreads,
  activeThreadId: initialThreads[0].id,
  openSidebar: () => set({ isOpen: true }),
  closeSidebar: () => set({ isOpen: false }),
  toggleSidebar: () => set((s) => ({ isOpen: !s.isOpen })),
  toggleExpanded: () => set((s) => ({ isExpanded: !s.isExpanded })),
  setMode: (mode) => set({ mode }),
  startNewThread: () => {
    const fresh = createEmptyThread();
    set((s) => {
      const kept = s.threads.filter((t) => t.messages.length > 0);
      const next = [fresh, ...kept];
      saveThreads(next);
      return { threads: next, activeThreadId: fresh.id };
    });
    return fresh.id;
  },
  selectThread: (id) => set({ activeThreadId: id }),
  deleteThread: (id) =>
    set((s) => {
      const remaining = s.threads.filter((t) => t.id !== id);
      const next = remaining.length > 0 ? remaining : [createEmptyThread()];
      saveThreads(next);
      return {
        threads: next,
        activeThreadId: s.activeThreadId === id ? next[0].id : s.activeThreadId,
      };
    }),
  appendMessage: (msg) =>
    set((s) => {
      const next = s.threads.map((t) => {
        if (t.id !== s.activeThreadId) return t;
        const messages = [...t.messages, msg];
        const title =
          t.title === 'New conversation' && msg.role === 'user'
            ? msg.content.slice(0, 42)
            : t.title;
        return { ...t, title, updatedAt: new Date().toISOString(), messages };
      });
      saveThreads(next);
      return { threads: next };
    }),
  updateActionStatus: (messageId, actionId, status) =>
    set((s) => {
      const next = s.threads.map((t) => {
        if (t.id !== s.activeThreadId) return t;
        return {
          ...t,
          messages: t.messages.map((m) =>
            m.id !== messageId
              ? m
              : {
                  ...m,
                  pendingActions: m.pendingActions?.map((a) =>
                    a.actionId === actionId ? { ...a, resolvedStatus: status } : a
                  ),
                }
          ),
        };
      });
      saveThreads(next);
      return { threads: next };
    }),
}));
