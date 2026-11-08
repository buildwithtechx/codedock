import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import type { AiQuestionGroup } from '#/interfaces/ai';
import { AskAiQuestions } from './ask-ai-questions';

const sampleQuestions: AiQuestionGroup[] = [
  {
    id: 'grp-rollback',
    title: 'Choose Rollback Target',
    questions: [
      {
        id: 'deployment_choice',
        title: 'Which deployment would you like to restore?',
        allowOther: true,
        options: [
          { id: 'dep-1', label: 'v1.2.0 (Stable)', description: 'Deployed 2 hours ago' },
          { id: 'dep-2', label: 'v1.1.9 (Previous)', description: 'Deployed yesterday' },
        ],
      },
    ],
  },
];

describe('AskAiQuestions', () => {
  it('renders question header and options', () => {
    const handleAnswer = vi.fn();
    const handleDismiss = vi.fn();

    render(
      <AskAiQuestions
        messageId="msg-1"
        questions={sampleQuestions}
        onAnswer={handleAnswer}
        onDismiss={handleDismiss}
      />
    );

    expect(screen.getByText('Waiting for your answer')).toBeDefined();
    expect(screen.getByText('Choose Rollback Target')).toBeDefined();
    expect(screen.getByText('v1.2.0 (Stable)')).toBeDefined();
    expect(screen.getByText('v1.1.9 (Previous)')).toBeDefined();
  });

  it('allows selecting an option and submitting the answer', () => {
    const handleAnswer = vi.fn();
    const handleDismiss = vi.fn();

    render(
      <AskAiQuestions
        messageId="msg-1"
        questions={sampleQuestions}
        onAnswer={handleAnswer}
        onDismiss={handleDismiss}
      />
    );

    // Click option
    const optionBtn = screen.getByText('v1.2.0 (Stable)');
    fireEvent.click(optionBtn);

    // Click submit
    const submitBtn = screen.getByRole('button', { name: 'Send answer' });
    fireEvent.click(submitBtn);

    expect(handleAnswer).toHaveBeenCalledWith(
      'Which deployment would you like to restore?: v1.2.0 (Stable)'
    );
    expect(handleDismiss).toHaveBeenCalledWith('grp-rollback');
  });

  it('can toggle collapsed state', () => {
    const handleAnswer = vi.fn();

    const { container } = render(
      <AskAiQuestions messageId="msg-1" questions={sampleQuestions} onAnswer={handleAnswer} />
    );

    const toggleBtn = screen.getByRole('button', { name: 'Hide agent questions' });
    const section = container.querySelector('section');
    expect(section?.getAttribute('data-assistant-questions-open')).toBe('true');

    fireEvent.click(toggleBtn);
    expect(section?.getAttribute('data-assistant-questions-open')).toBe('false');
  });
});
