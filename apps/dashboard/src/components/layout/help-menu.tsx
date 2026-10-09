import {
  BookOpen,
  CircleHelp,
  Flag,
  LifeBuoy,
  type LucideIcon,
  MessagesSquare,
} from 'lucide-react';
import { Button } from '#/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '#/components/ui/dropdown-menu';

const DOCS_URL = 'https://docs.codedock.run';
const SUPPORT_URL = 'https://codedock.run/support';
const ISSUES_URL = 'https://github.com/buildwithtechx/codedock/issues';
const COMMUNITY_URL = 'https://codedock.run/community';

const actions: { id: string; label: string; icon: LucideIcon; url: string }[] = [
  { id: 'support', label: 'Contact support', icon: LifeBuoy, url: SUPPORT_URL },
  { id: 'report-issue', label: 'Report an issue', icon: Flag, url: ISSUES_URL },
  { id: 'feedback', label: 'Share feedback', icon: MessagesSquare, url: SUPPORT_URL },
  { id: 'docs', label: 'Documentation', icon: BookOpen, url: DOCS_URL },
  { id: 'community', label: 'Community', icon: MessagesSquare, url: COMMUNITY_URL },
];

export function HelpMenu() {
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="outline" size="icon" aria-label="Help and support">
          <CircleHelp className="h-4 w-4" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-52">
        {actions.map((action) => (
          <DropdownMenuItem key={action.id} asChild>
            <a href={action.url} target="_blank" rel="noopener noreferrer">
              <action.icon className="h-4 w-4" />
              {action.label}
            </a>
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
