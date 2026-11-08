import { create } from 'zustand';
import type {
  AiConversationThread,
  AiPendingAction,
  AiPersonaId,
  AiResourceTag,
  AiToolStep,
  AiTurnMessage,
} from '#/interfaces/ai';

export type {
  AiConversationThread,
  AiPendingAction,
  AiPersonaId,
  AiResourceTag,
  AiToolStep,
  AiTurnMessage,
};

interface AiSidebarState {
  isOpen: boolean;
  isExpanded: boolean;
  mode: 'ask' | 'agent';
  activePersona: AiPersonaId;
  threads: AiConversationThread[];
  activeThreadId: string;
  pendingPrompt: string | null;
  autoSendPrompt: boolean;
  openSidebar: () => void;
  closeSidebar: () => void;
  toggleSidebar: () => void;
  toggleExpanded: () => void;
  setMode: (mode: 'ask' | 'agent') => void;
  setPersona: (persona: AiPersonaId) => void;
  openWithPrompt: (prompt: string, autoSend?: boolean) => void;
  clearPendingPrompt: () => void;
  startNewThread: (persona?: AiPersonaId) => string;
  selectThread: (id: string) => void;
  deleteThread: (id: string) => void;
  appendMessage: (msg: AiTurnMessage) => void;
  dismissQuestions: (messageId: string, questionGroupId: string) => void;
  updateActionStatus: (
    messageId: string,
    actionId: string,
    status: 'approved' | 'cancelled' | 'failed'
  ) => void;
}

const STORAGE_KEY = 'codedock_ai_threads_v1';

function createEmptyThread(persona: AiPersonaId = 'general'): AiConversationThread {
  return {
    id: `thread_${Date.now().toString(36)}`,
    title: 'New conversation',
    updatedAt: new Date().toISOString(),
    messages: [],
    persona,
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
  activePersona: initialThreads[0]?.persona || 'general',
  threads: initialThreads,
  activeThreadId: initialThreads[0].id,
  pendingPrompt: null,
  autoSendPrompt: false,
  openSidebar: () => set({ isOpen: true }),
  closeSidebar: () => set({ isOpen: false }),
  toggleSidebar: () => set((s) => ({ isOpen: !s.isOpen })),
  toggleExpanded: () => set((s) => ({ isExpanded: !s.isExpanded })),
  setMode: (mode) => set({ mode }),
  setPersona: (persona) =>
    set((s) => {
      const next = s.threads.map((t) => (t.id === s.activeThreadId ? { ...t, persona } : t));
      saveThreads(next);
      return { activePersona: persona, threads: next };
    }),
  openWithPrompt: (prompt, autoSend = true) =>
    set({
      isOpen: true,
      pendingPrompt: prompt,
      autoSendPrompt: autoSend,
    }),
  clearPendingPrompt: () => set({ pendingPrompt: null, autoSendPrompt: false }),
  startNewThread: (persona) => {
    let threadId = '';
    set((s) => {
      const targetPersona = persona || s.activePersona || 'general';
      const fresh = createEmptyThread(targetPersona);
      threadId = fresh.id;
      const kept = s.threads.filter((t) => t.messages.length > 0);
      const next = [fresh, ...kept];
      saveThreads(next);
      return { threads: next, activeThreadId: fresh.id, activePersona: targetPersona };
    });
    return threadId;
  },
  selectThread: (id) =>
    set((s) => {
      const target = s.threads.find((t) => t.id === id);
      return {
        activeThreadId: id,
        activePersona: target?.persona || 'general',
      };
    }),
  deleteThread: (id) =>
    set((s) => {
      const remaining = s.threads.filter((t) => t.id !== id);
      const next = remaining.length > 0 ? remaining : [createEmptyThread()];
      saveThreads(next);
      const nextActive =
        s.activeThreadId === id
          ? next[0]
          : s.threads.find((t) => t.id === s.activeThreadId) || next[0];
      return {
        threads: next,
        activeThreadId: nextActive.id,
        activePersona: nextActive.persona || 'general',
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
  dismissQuestions: (messageId, questionGroupId) =>
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
                  questions: m.questions?.filter((q) => q.id !== questionGroupId),
                }
          ),
        };
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
