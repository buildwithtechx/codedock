import { useState } from 'react';
import { Button } from '#/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '#/components/ui/card';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { useCreateMailServer } from '#/hooks/use-mail';

export function MailServerForm({ onDone }: { onDone: () => void }) {
  const [domain, setDomain] = useState('');
  const [maildUrl, setMaildUrl] = useState('http://localhost:8081');
  const [webmailUrl, setWebmailUrl] = useState('');
  const [apiToken, setApiToken] = useState('');
  const createServer = useCreateMailServer();

  const submit = () => {
    createServer.mutate(
      { domain, maildUrl, webmailUrl: webmailUrl || undefined, apiToken },
      { onSuccess: () => onDone() }
    );
  };

  const valid = domain.trim() !== '' && maildUrl.trim() !== '' && apiToken.trim() !== '';

  return (
    <Card className="border-border/60 bg-card">
      <CardHeader>
        <CardTitle>Register mail server</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <div className="space-y-2">
            <Label htmlFor="mail-domain">Primary domain</Label>
            <Input
              id="mail-domain"
              placeholder="example.com"
              value={domain}
              onChange={(e) => setDomain(e.target.value)}
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="mail-maild">maild URL</Label>
            <Input
              id="mail-maild"
              placeholder="http://mail-host:8081"
              value={maildUrl}
              onChange={(e) => setMaildUrl(e.target.value)}
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="mail-webmail">Webmail URL (optional)</Label>
            <Input
              id="mail-webmail"
              placeholder="https://mail.example.com"
              value={webmailUrl}
              onChange={(e) => setWebmailUrl(e.target.value)}
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="mail-token">maild API token</Label>
            <Input
              id="mail-token"
              type="password"
              placeholder="MAILD_API_TOKEN value"
              value={apiToken}
              onChange={(e) => setApiToken(e.target.value)}
            />
          </div>
        </div>
        <div className="flex justify-end gap-2">
          <Button variant="outline" onClick={onDone}>
            Cancel
          </Button>
          <Button disabled={!valid || createServer.isPending} onClick={submit}>
            {createServer.isPending ? 'Registering…' : 'Register'}
          </Button>
        </div>
      </CardContent>
    </Card>
  );
}
