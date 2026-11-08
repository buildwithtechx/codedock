import { beforeEach, describe, expect, it } from 'vitest';
import { AI_PERSONAS } from '#/interfaces/ai';
import { useAiSidebarStore } from '#/stores/ai-sidebar-store';

describe('AI Specialist Personas', () => {
  beforeEach(() => {
    useAiSidebarStore.setState({
      isOpen: true,
      activePersona: 'general',
      threads: [
        {
          id: 'thread_test_1',
          title: 'Test conversation',
          updatedAt: new Date().toISOString(),
          messages: [],
          persona: 'general',
        },
      ],
      activeThreadId: 'thread_test_1',
    });
  });

  it('defines the 4 required specialist personas with proper readOnly security flags', () => {
    const ids = AI_PERSONAS.map((p) => p.id);
    expect(ids).toContain('general');
    expect(ids).toContain('logs');
    expect(ids).toContain('security');
    expect(ids).toContain('devops');

    const logs = AI_PERSONAS.find((p) => p.id === 'logs');
    expect(logs?.readOnly).toBe(true);
    expect(logs?.icon).toBe('terminal');

    const security = AI_PERSONAS.find((p) => p.id === 'security');
    expect(security?.readOnly).toBe(true);
    expect(security?.icon).toBe('shield-check');

    const devops = AI_PERSONAS.find((p) => p.id === 'devops');
    expect(devops?.readOnly).toBe(false);
    expect(devops?.icon).toBe('rocket');

    const general = AI_PERSONAS.find((p) => p.id === 'general');
    expect(general?.readOnly).toBe(false);
  });

  it('allows switching persona and updates the active thread state in store', () => {
    const store = useAiSidebarStore.getState();
    expect(store.activePersona).toBe('general');

    store.setPersona('logs');
    expect(useAiSidebarStore.getState().activePersona).toBe('logs');
    expect(useAiSidebarStore.getState().threads[0].persona).toBe('logs');

    store.setPersona('security');
    expect(useAiSidebarStore.getState().activePersona).toBe('security');
    expect(useAiSidebarStore.getState().threads[0].persona).toBe('security');
  });

  it('creates new thread preserving or overriding persona', () => {
    const store = useAiSidebarStore.getState();
    store.setPersona('logs');

    const newId = store.startNewThread('security');
    const updated = useAiSidebarStore.getState();
    expect(updated.activeThreadId).toBe(newId);
    expect(updated.activePersona).toBe('security');
    const newThread = updated.threads.find((t) => t.id === newId);
    expect(newThread?.persona).toBe('security');
  });

  it('restores thread persona when selecting an existing thread', () => {
    useAiSidebarStore.setState({
      threads: [
        {
          id: 't1',
          title: 'Logs thread',
          updatedAt: new Date().toISOString(),
          messages: [],
          persona: 'logs',
        },
        {
          id: 't2',
          title: 'Security thread',
          updatedAt: new Date().toISOString(),
          messages: [],
          persona: 'security',
        },
      ],
      activeThreadId: 't1',
      activePersona: 'logs',
    });

    const store = useAiSidebarStore.getState();
    store.selectThread('t2');

    expect(useAiSidebarStore.getState().activeThreadId).toBe('t2');
    expect(useAiSidebarStore.getState().activePersona).toBe('security');
  });
});
