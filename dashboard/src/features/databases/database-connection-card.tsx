import { Copy, Eye, EyeOff, Loader2 } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '#/components/ui/card';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { type Database, databasesService } from '#/features/databases';

interface Props {
  database: Database;
}

export function DatabaseConnectionCard({ database }: Props) {
  const [showPassword, setShowPassword] = useState(false);
  const [revealedPassword, setRevealedPassword] = useState<string | null>(
    database.password || null
  );
  const [isRevealing, setIsRevealing] = useState(false);

  const handleRevealPassword = async (): Promise<string | null> => {
    if (revealedPassword) {
      return revealedPassword;
    }
    try {
      setIsRevealing(true);
      const res = await databasesService.revealCredentials(database.id);
      if (res.password) {
        setRevealedPassword(res.password);
        return res.password;
      }
      return null;
    } catch (_error) {
      toast.error('Failed to reveal credentials');
      return null;
    } finally {
      setIsRevealing(false);
    }
  };

  const toggleShowPassword = async () => {
    if (!showPassword && !revealedPassword) {
      const pwd = await handleRevealPassword();
      if (pwd) setShowPassword(true);
    } else {
      setShowPassword(!showPassword);
    }
  };

  const copyToClipboard = async (text: string, label: string) => {
    let copyText = text;
    if (copyText.includes('<password>') || label === 'Password') {
      const pwd = await handleRevealPassword();
      if (!pwd) {
        return;
      }
      if (label === 'Password') {
        copyText = pwd;
      } else {
        copyText = copyText.replace('<password>', pwd);
      }
    }
    if (copyText.includes('<password>')) {
      return;
    }
    try {
      await navigator.clipboard.writeText(copyText);
      toast.success(`${label} copied to clipboard`);
    } catch (_error) {
      toast.error('Failed to copy to clipboard');
    }
  };

  const getFormattedUrl = (includeMasking: boolean) => {
    const { engine, username, internalDns, port, databaseName } = database;
    let scheme = engine;
    if (engine === 'postgresql') scheme = 'postgres';
    if (engine === 'mongodb') scheme = 'mongodb';
    if (engine === 'redis' || engine === 'dragonfly' || engine === 'keydb') scheme = 'redis';

    const activePassword = revealedPassword;
    const pwdSegment = activePassword
      ? includeMasking && !showPassword
        ? '••••••••'
        : activePassword
      : '<password>';

    const hasPassword = activePassword === null || activePassword.length > 0;
    if (engine === 'redis' || engine === 'dragonfly' || engine === 'keydb') {
      return `${scheme}://${hasPassword ? `${pwdSegment}@` : ''}${internalDns}:${port}`;
    }

    return `${scheme}://${username}:${pwdSegment}@${internalDns}:${port}/${databaseName}`;
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Connection Details</CardTitle>
        <CardDescription>
          Use these credentials to connect to your database instance internally.
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="space-y-2">
          <Label>Connection URL</Label>
          <div className="flex space-x-2">
            <Input
              type="text"
              value={getFormattedUrl(true)}
              readOnly
              className="bg-muted/50 font-mono text-sm"
            />
            <Button
              variant="outline"
              size="icon"
              disabled={isRevealing}
              onClick={toggleShowPassword}
              title={showPassword ? 'Hide Password' : 'Reveal Password'}
            >
              {isRevealing ? (
                <Loader2 className="h-4 w-4 animate-spin" />
              ) : showPassword ? (
                <EyeOff className="h-4 w-4" />
              ) : (
                <Eye className="h-4 w-4" />
              )}
            </Button>
            <Button
              variant="outline"
              size="icon"
              onClick={() => copyToClipboard(getFormattedUrl(false), 'Connection URL')}
              title="Copy Connection URL"
            >
              <Copy className="h-4 w-4" />
            </Button>
          </div>
        </div>

        <div className="grid grid-cols-1 gap-4 pt-2 md:grid-cols-2">
          <div className="space-y-1">
            <Label className="text-muted-foreground text-xs">Host (Internal)</Label>
            <div className="flex items-center space-x-2">
              <span className="flex-1 font-mono text-sm">{database.internalDns}</span>
              <Button
                variant="ghost"
                size="icon"
                className="h-6 w-6"
                onClick={() => copyToClipboard(database.internalDns || '', 'Host')}
              >
                <Copy className="h-3 w-3" />
              </Button>
            </div>
          </div>
          <div className="space-y-1">
            <Label className="text-muted-foreground text-xs">Port</Label>
            <div className="flex items-center space-x-2">
              <span className="flex-1 font-mono text-sm">{database.port}</span>
              <Button
                variant="ghost"
                size="icon"
                className="h-6 w-6"
                onClick={() => copyToClipboard(database.port.toString(), 'Port')}
              >
                <Copy className="h-3 w-3" />
              </Button>
            </div>
          </div>
          <div className="space-y-1">
            <Label className="text-muted-foreground text-xs">Username</Label>
            <div className="flex items-center space-x-2">
              <span className="flex-1 font-mono text-sm">{database.username}</span>
              <Button
                variant="ghost"
                size="icon"
                className="h-6 w-6"
                onClick={() => copyToClipboard(database.username, 'Username')}
              >
                <Copy className="h-3 w-3" />
              </Button>
            </div>
          </div>
          <div className="space-y-1">
            <Label className="text-muted-foreground text-xs">Password</Label>
            <div className="flex items-center space-x-2">
              <span className="flex-1 font-mono text-muted-foreground text-sm italic">
                {revealedPassword
                  ? showPassword
                    ? revealedPassword
                    : '••••••••'
                  : '•••••••• (Click Eye to Reveal)'}
              </span>
              <Button
                variant="ghost"
                size="icon"
                className="h-6 w-6"
                disabled={isRevealing}
                onClick={() => copyToClipboard(revealedPassword || '', 'Password')}
                title="Copy Password"
              >
                <Copy className="h-3 w-3" />
              </Button>
            </div>
          </div>
          <div className="space-y-1">
            <Label className="text-muted-foreground text-xs">Database Name</Label>
            <div className="flex items-center space-x-2">
              <span className="flex-1 font-mono text-sm">{database.databaseName}</span>
              <Button
                variant="ghost"
                size="icon"
                className="h-6 w-6"
                onClick={() => copyToClipboard(database.databaseName, 'Database Name')}
              >
                <Copy className="h-3 w-3" />
              </Button>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
