import { Bot, Check, GitBranch, MessageSquare, Rocket, Search, Workflow } from 'lucide-react';
import type React from 'react';
import { useId, useRef, useState } from 'react';

interface PreviewExample {
  label: string;
  icon: React.ComponentType<{ className?: string }>;
  agent: string;
  avatarBg: string;
  prompt: string;
  steps: string[];
  result: string;
  detail: string;
  caption: string;
  resultIcon: React.ComponentType<{ className?: string }>;
}

const EXAMPLES: PreviewExample[] = [
  {
    label: 'Deploy',
    icon: Rocket,
    agent: 'Codedock assistant',
    avatarBg: 'bg-emerald-500/20 text-emerald-400 border-emerald-500/30',
    prompt: 'Deploy my API from GitHub and connect its domain.',
    steps: [
      'Inspect the repository',
      'Build and deploy the app',
      'Connect the domain and verify health',
    ],
    result: 'Your API is live.',
    detail: 'Build passed, domain connected, health checks passing.',
    caption: 'From a request to a running app, in one conversation.',
    resultIcon: Check,
  },
  {
    label: 'Investigate',
    icon: Search,
    agent: 'Deployment investigator',
    avatarBg: 'bg-blue-500/20 text-blue-400 border-blue-500/30',
    prompt: 'The latest deploy failed. Find the cause and fix it.',
    steps: [
      'Read the build logs',
      'Trace the failure to a code change',
      'Prepare a fix in a pull request',
    ],
    result: 'Found the cause. Fix ready to review.',
    detail: 'A missing import broke the build. The patch is in a new branch.',
    caption: 'An agent that follows the evidence and takes the next step.',
    resultIcon: GitBranch,
  },
  {
    label: 'Automate',
    icon: Workflow,
    agent: 'Production watcher',
    avatarBg: 'bg-purple-500/20 text-purple-400 border-purple-500/30',
    prompt: 'When my app goes down, investigate and send me a report.',
    steps: [
      'Unhealthy service triggers an investigation',
      'Agent checks logs and recent changes',
      'Findings arrive in your chat',
    ],
    result: 'The next incident has an agent on it.',
    detail: 'Give it a scope, choose a trigger, and enable the automation.',
    caption: 'Around the clock, while Codedock and your agent host are online.',
    resultIcon: MessageSquare,
  },
];

export function AssistantIntroPreview() {
  const [selected, setSelected] = useState(0);
  const tabs = useRef<Array<HTMLButtonElement | null>>([]);
  const id = useId();
  const example = EXAMPLES[selected];
  const ResultIconComp = example.resultIcon;

  return (
    <div className="mx-auto w-full min-w-0 max-w-[540px] rounded-2xl border border-border/50 bg-card shadow-sm">
      <div className="px-5 pt-5 sm:px-6 sm:pt-6">
        <p className="mb-3 font-medium text-muted-foreground text-xs">Explore an example</p>
        <div
          role="tablist"
          aria-label="What agents can do"
          className="flex gap-1 rounded-xl bg-muted/60 p-1"
        >
          {EXAMPLES.map((item, index) => {
            const TabIcon = item.icon;
            const isTabActive = selected === index;
            return (
              <button
                key={item.label}
                ref={(el) => {
                  tabs.current[index] = el;
                }}
                type="button"
                role="tab"
                id={`${id}-tab-${index}`}
                aria-selected={isTabActive}
                aria-controls={`${id}-panel`}
                tabIndex={isTabActive ? 0 : -1}
                onClick={() => setSelected(index)}
                className={`flex min-w-0 flex-1 items-center justify-center gap-1.5 rounded-lg px-2 py-2 font-medium text-xs transition-colors focus-visible:outline-none sm:text-sm ${
                  isTabActive
                    ? 'bg-card text-foreground shadow-xs'
                    : 'text-muted-foreground hover:bg-card/50 hover:text-foreground'
                }`}
              >
                <TabIcon className="h-4 w-4 shrink-0" />
                {item.label}
              </button>
            );
          })}
        </div>
      </div>

      <div
        role="tabpanel"
        id={`${id}-panel`}
        aria-labelledby={`${id}-tab-${selected}`}
        className="rounded-2xl p-5 focus-visible:outline-none sm:p-6"
      >
        <div className="flex items-center gap-2.5">
          <div
            className={`flex h-8 w-8 items-center justify-center rounded-xl border ${example.avatarBg}`}
          >
            <Bot className="h-4 w-4" />
          </div>
          <p className="min-w-0 font-medium text-foreground text-sm">{example.agent}</p>
        </div>

        <div className="mt-5 space-y-5">
          <p className="ms-auto max-w-[92%] rounded-2xl rounded-ee-md bg-muted/80 px-4 py-3 text-foreground text-sm leading-6">
            {example.prompt}
          </p>

          <ol className="my-6 space-y-4" aria-label="Example actions">
            {example.steps.map((step, index) => (
              <li
                key={step}
                className="relative flex items-start gap-3 text-[13px] text-muted-foreground leading-5"
              >
                {index < example.steps.length - 1 && (
                  <span
                    aria-hidden="true"
                    className="absolute start-[9px] top-6 -bottom-3 border-border border-s border-dashed"
                  />
                )}
                <span className="flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-muted text-muted-foreground">
                  <Check className="h-3 w-3" />
                </span>
                <span>{step}</span>
              </li>
            ))}
          </ol>

          <div className="rounded-xl border border-border/40 bg-muted/40 p-4">
            <p className="flex items-start gap-2.5 font-medium text-foreground text-sm leading-6">
              <ResultIconComp className="mt-0.5 h-5 w-5 shrink-0 text-emerald-500" />
              {example.result}
            </p>
            <p className="mt-1 ps-[30px] text-muted-foreground text-xs leading-5">
              {example.detail}
            </p>
          </div>
        </div>
      </div>

      <p className="border-border/40 border-t px-5 py-4 text-muted-foreground text-xs leading-5 sm:px-6">
        {example.caption}
      </p>
    </div>
  );
}
