import { Button } from '@codedock/ui/button';
import { Input } from '@codedock/ui/input';
import { Check, ChevronDown, ChevronUp, HelpCircle } from 'lucide-react';
import { useId, useState } from 'react';
import type { AiQuestionAnswer, AiQuestionGroup } from '#/interfaces/ai';

interface AskAiQuestionsProps {
  messageId: string;
  questions: AiQuestionGroup[];
  onAnswer: (answerText: string) => void;
  onDismiss?: (groupId: string) => void;
  disabled?: boolean;
}

export function AskAiQuestions({
  messageId: _messageId,
  questions,
  onAnswer,
  onDismiss,
  disabled = false,
}: AskAiQuestionsProps) {
  const [opened, setOpened] = useState(true);
  const [selectedGroup, setSelectedGroup] = useState<string | null>(null);
  const [drafts, setDrafts] = useState<Record<string, Record<string, AiQuestionAnswer>>>({});
  const panelId = useId();

  if (!questions || questions.length === 0) return null;

  const currentGroup = questions.find((g) => g.id === selectedGroup) || questions[0];
  const groupKey = currentGroup.id;
  const groupAnswers = drafts[groupKey] || {};

  const isReady = currentGroup.questions.every((q) => {
    if (q.optional) return true;
    const ans = groupAnswers[q.id];
    return Boolean(
      (ans?.options && ans.options.length > 0) || (ans?.text && ans.text.trim().length > 0)
    );
  });

  const handleUpdateAnswer = (questionId: string, val: AiQuestionAnswer) => {
    setDrafts((prev) => ({
      ...prev,
      [groupKey]: {
        ...prev[groupKey],
        [questionId]: val,
      },
    }));
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!isReady || disabled) return;

    // Format human-readable summary to send back as user answer
    const parts: string[] = [];
    for (const q of currentGroup.questions) {
      const ans = groupAnswers[q.id];
      if (!ans) continue;
      const selectedLabels = (ans.options || []).map((optId) => {
        const found = q.options?.find((o) => o.id === optId);
        return found ? found.label : optId;
      });

      if (selectedLabels.length > 0 && ans.text?.trim()) {
        parts.push(`${q.title}: ${selectedLabels.join(', ')} (${ans.text.trim()})`);
      } else if (selectedLabels.length > 0) {
        parts.push(`${q.title}: ${selectedLabels.join(', ')}`);
      } else if (ans.text?.trim()) {
        parts.push(`${q.title}: ${ans.text.trim()}`);
      }
    }

    const answerSummary = parts.length > 0 ? parts.join('\n') : 'Confirmed';
    onAnswer(answerSummary);
    if (onDismiss) {
      onDismiss(currentGroup.id);
    }
  };

  return (
    <section
      className="mb-2 overflow-hidden rounded-xl border border-primary/20 bg-card/95 shadow-lg backdrop-blur-xl transition-all"
      aria-label="Agent questions"
      data-assistant-questions-open={opened}
    >
      <button
        type="button"
        className="flex h-auto w-full items-center justify-between gap-3 px-3.5 py-2.5 text-left text-foreground transition-colors hover:bg-muted/40"
        aria-label={opened ? 'Hide agent questions' : 'Show agent questions'}
        aria-expanded={opened}
        aria-controls={panelId}
        disabled={disabled}
        onClick={() => setOpened((v) => !v)}
      >
        <div className="flex min-w-0 items-center gap-2">
          <HelpCircle className="size-4 shrink-0 text-primary" />
          <span role="status" className="truncate font-medium text-foreground text-xs">
            {questions.length === 1
              ? 'Waiting for your answer'
              : `${questions.length} requests need your answer`}
          </span>
        </div>
        <div className="flex items-center gap-1.5 text-[11px] text-muted-foreground">
          <span>{opened ? 'Hide' : 'Answer'}</span>
          {opened ? <ChevronUp className="size-3.5" /> : <ChevronDown className="size-3.5" />}
        </div>
      </button>

      <form
        id={panelId}
        hidden={!opened}
        aria-label="Answer agent questions"
        onSubmit={handleSubmit}
        className={opened ? 'border-border/50 border-t p-3 pt-2' : 'hidden'}
      >
        {questions.length > 1 && (
          <div className="mb-2.5 flex flex-wrap gap-1.5">
            {questions.map((g, idx) => (
              <Button
                key={g.id}
                type="button"
                variant={g.id === currentGroup.id ? 'secondary' : 'ghost'}
                size="sm"
                className="h-6 px-2 text-[10px]"
                aria-pressed={g.id === currentGroup.id}
                disabled={disabled}
                onClick={() => setSelectedGroup(g.id)}
              >
                Question {idx + 1}
              </Button>
            ))}
          </div>
        )}

        <div className="max-h-[min(45dvh,20rem)] space-y-3.5 overflow-y-auto pr-1">
          <h4 className="font-semibold text-foreground text-xs">
            {currentGroup.title || 'A question for you'}
          </h4>

          {currentGroup.questions.map((question) => {
            const currentVal = groupAnswers[question.id] || {};
            const selectedOptions = currentVal.options || [];

            return (
              <fieldset key={question.id} disabled={disabled} className="space-y-2">
                <legend className="font-medium text-[11px] text-muted-foreground">
                  {question.title}
                  {question.optional && ' (optional)'}
                </legend>

                {question.options && question.options.length > 0 && (
                  <div className="space-y-1.5">
                    {question.options.map((opt) => {
                      const isSelected = selectedOptions.includes(opt.id);
                      return (
                        <button
                          key={opt.id}
                          type="button"
                          onClick={() => {
                            let nextOpts: string[];
                            if (question.multiSelect) {
                              nextOpts = isSelected
                                ? selectedOptions.filter((id) => id !== opt.id)
                                : [...selectedOptions, opt.id];
                            } else {
                              nextOpts = [opt.id];
                            }
                            handleUpdateAnswer(question.id, {
                              ...currentVal,
                              options: nextOpts,
                            });
                          }}
                          className={`flex w-full cursor-pointer items-start gap-2.5 rounded-lg border p-2 text-left transition-colors ${
                            isSelected
                              ? 'border-primary bg-primary/10 text-foreground'
                              : 'border-border/60 bg-background/70 text-muted-foreground hover:bg-muted/50'
                          }`}
                        >
                          <div
                            className={`mt-0.5 flex size-3.5 shrink-0 items-center justify-center rounded ${
                              question.multiSelect ? 'rounded-xs' : 'rounded-full'
                            } border ${
                              isSelected
                                ? 'border-primary bg-primary text-primary-foreground'
                                : 'border-muted-foreground/40'
                            }`}
                          >
                            {isSelected && <Check className="size-2.5 stroke-[3]" />}
                          </div>
                          <div className="min-w-0 flex-1">
                            <span className="block font-medium text-foreground text-xs">
                              {opt.label}
                            </span>
                            {opt.description && (
                              <span className="mt-0.5 block text-[10px] text-muted-foreground leading-normal">
                                {opt.description}
                              </span>
                            )}
                          </div>
                        </button>
                      );
                    })}
                  </div>
                )}

                {(!question.options?.length || question.allowOther !== false) && (
                  <Input
                    type={question.isSecret ? 'password' : 'text'}
                    aria-label={`Answer: ${question.title}`}
                    value={currentVal.text || ''}
                    onChange={(e) =>
                      handleUpdateAnswer(question.id, {
                        ...currentVal,
                        text: e.target.value,
                      })
                    }
                    placeholder="Your answer or more context..."
                    className="h-8 bg-background/80 text-xs"
                    maxLength={8000}
                  />
                )}
              </fieldset>
            );
          })}
        </div>

        <div className="mt-3 flex items-center justify-end gap-2 border-border/40 border-t pt-2.5">
          <Button
            type="submit"
            size="sm"
            disabled={!isReady || disabled}
            className="h-7 px-3 text-xs"
          >
            Send answer
          </Button>
        </div>
      </form>
    </section>
  );
}
