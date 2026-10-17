import { Eye, EyeOff, KeyRound, Lock, Upload } from 'lucide-react';
import type React from 'react';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';

export interface ServerSshAuthFieldsProps {
  sshAuthMethod: 'key' | 'password';
  setSshAuthMethod: (method: 'key' | 'password') => void;
  sshPrivateKey: string;
  setSshPrivateKey: (key: string) => void;
  sshPassword: string;
  setSshPassword: (password: string) => void;
  showPassword: boolean;
  setShowPassword: (show: boolean | ((prev: boolean) => boolean)) => void;
  fileInputRef: React.RefObject<HTMLInputElement | null>;
  handleFileUpload: (event: React.ChangeEvent<HTMLInputElement>) => void;
}

export function ServerSshAuthFields({
  sshAuthMethod,
  setSshAuthMethod,
  sshPrivateKey,
  setSshPrivateKey,
  sshPassword,
  setSshPassword,
  showPassword,
  setShowPassword,
  fileInputRef,
  handleFileUpload,
}: ServerSshAuthFieldsProps) {
  return (
    <div className="space-y-3 rounded-xl border border-border/70 bg-muted/20 p-4">
      <div className="flex items-center justify-between">
        <Label className="font-semibold text-xs">Authentication Method</Label>
        <div className="inline-flex rounded-lg border border-border/80 bg-background p-0.5 text-xs">
          <button
            type="button"
            onClick={() => setSshAuthMethod('key')}
            className={`inline-flex items-center gap-1.5 rounded-md px-3 py-1 font-medium transition-colors ${
              sshAuthMethod === 'key'
                ? 'bg-primary text-primary-foreground'
                : 'text-muted-foreground hover:text-foreground'
            }`}
          >
            <KeyRound className="h-3 w-3" />
            Private Key
          </button>
          <button
            type="button"
            onClick={() => setSshAuthMethod('password')}
            className={`inline-flex items-center gap-1.5 rounded-md px-3 py-1 font-medium transition-colors ${
              sshAuthMethod === 'password'
                ? 'bg-primary text-primary-foreground'
                : 'text-muted-foreground hover:text-foreground'
            }`}
          >
            <Lock className="h-3 w-3" />
            Password
          </button>
        </div>
      </div>

      {sshAuthMethod === 'key' ? (
        <div className="space-y-2">
          <div className="flex items-center justify-between">
            <span className="text-muted-foreground text-xs">
              Paste OpenSSH / PEM key or upload key file
            </span>
            <Button
              type="button"
              variant="outline"
              size="sm"
              className="h-7 gap-1 text-xs"
              onClick={() => fileInputRef.current?.click()}
            >
              <Upload className="h-3 w-3" />
              Upload file
            </Button>
            <input ref={fileInputRef} type="file" className="hidden" onChange={handleFileUpload} />
          </div>
          <textarea
            value={sshPrivateKey}
            onChange={(e) => setSshPrivateKey(e.target.value)}
            placeholder="-----BEGIN OPENSSH PRIVATE KEY-----&#10;...&#10;-----END OPENSSH PRIVATE KEY-----"
            rows={4}
            spellCheck={false}
            className="w-full rounded-lg border border-border bg-background p-2.5 font-mono text-xs outline-none focus:ring-2 focus:ring-primary/20"
            required
          />
        </div>
      ) : (
        <div className="space-y-1.5">
          <Label htmlFor="srv-pass" className="text-xs">
            SSH Password
          </Label>
          <div className="relative">
            <Input
              id="srv-pass"
              type={showPassword ? 'text' : 'password'}
              value={sshPassword}
              onChange={(e) => setSshPassword(e.target.value)}
              placeholder="••••••••••••"
              className="pr-10"
              required
            />
            <button
              type="button"
              onClick={() => setShowPassword((p) => !p)}
              className="absolute top-1/2 right-3 -translate-y-1/2 text-muted-foreground hover:text-foreground"
            >
              {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
