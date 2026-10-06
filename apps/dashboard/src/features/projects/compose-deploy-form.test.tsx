import { act, fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { ComposeDeployForm } from './compose-deploy-form';

const mutations = vi.hoisted(() => ({ analyze: vi.fn(), deploy: vi.fn() }));
vi.mock('#/hooks/use-compose', () => ({
  useAnalyzeCompose: () => ({ mutate: mutations.analyze, isPending: false }),
  useDeployCompose: () => ({ mutate: mutations.deploy, isPending: false }),
}));
vi.mock('@monaco-editor/react', () => ({
  default: ({ value, onChange }: { value: string; onChange: (value: string) => void }) => (
    <textarea
      aria-label="Compose source"
      value={value}
      onChange={(event) => onChange(event.target.value)}
    />
  ),
}));
vi.mock('sonner', () => ({ toast: { success: vi.fn(), error: vi.fn() } }));

describe('ComposeDeployForm', () => {
  it('rejects an analysis response after its source changed', () => {
    render(<ComposeDeployForm projectId="project" />);
    fireEvent.click(screen.getByRole('button', { name: 'Analyze' }));
    const callback = mutations.analyze.mock.calls[0][1].onSuccess;
    fireEvent.change(screen.getByLabelText('Compose source'), {
      target: { value: 'services: {}' },
    });
    act(() => callback({ appServices: [{ name: 'old-web', imageRef: 'nginx' }], databases: [] }));
    expect(screen.queryByText('old-web')).toBeNull();
    expect(
      (screen.getByRole('button', { name: 'Import resources' }) as HTMLButtonElement).disabled
    ).toBe(true);
    expect(mutations.deploy).not.toHaveBeenCalled();
  });
});
