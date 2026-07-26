import { Database, Eye, EyeOff, Info } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogTitle,
  DialogTrigger,
} from '#/components/ui/dialog';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '#/components/ui/select';
import { useCreateS3Destination } from '#/features/backups/hooks';

type Props = {
  isOpen: boolean;
  setIsOpen: (open: boolean) => void;
  trigger: React.ReactNode;
};

export function CreateS3DestinationDialog({ isOpen, setIsOpen, trigger }: Props) {
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [provider, setProvider] = useState('r2');
  const [endpoint, setEndpoint] = useState('');
  const [bucket, setBucket] = useState('');
  const [region, setRegion] = useState('us-east-1');
  const [accessKeyId, setAccessKeyId] = useState('');
  const [secretAccessKey, setSecretAccessKey] = useState('');
  const [showSecret, setShowSecret] = useState(false);

  const createS3Dest = useCreateS3Destination();

  const handleProviderChange = (value: string) => {
    setProvider(value);
    if (value === 'r2') {
      setEndpoint('https://<account_id>.r2.cloudflarestorage.com');
      setRegion('auto');
    } else if (value === 's3') {
      setRegion('us-east-1');
      setEndpoint('https://s3.us-east-1.amazonaws.com');
    } else {
      setEndpoint('');
      setRegion('');
    }
  };

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await createS3Dest.mutateAsync({
        payload: {
          projectId: 'global',
          name,
          description,
          provider,
          endpoint,
          bucket,
          region,
          accessKeyId,
          secretAccessKey,
        },
      });
      toast.success('S3 destination added successfully');
      setIsOpen(false);
      setName('');
      setDescription('');
      setProvider('r2');
      setEndpoint('');
      setBucket('');
      setRegion('us-east-1');
      setAccessKeyId('');
      setSecretAccessKey('');
    } catch {
      toast.error('Failed to add S3 destination');
    }
  };

  return (
    <Dialog open={isOpen} onOpenChange={setIsOpen}>
      <DialogTrigger asChild>{trigger}</DialogTrigger>
      <DialogContent className="border/50 gap-0 bg-card/95 p-0 backdrop-blur-xl sm:max-w-150 [&>button]:hidden">
        <div className="px-5 pt-5 pb-4">
          <div className="flex items-start justify-between">
            <div className="flex flex-col">
              <DialogTitle className="flex items-center gap-2 font-bold text-foreground text-xl tracking-tight">
                <Database className="h-5 w-5 text-primary" />
                New S3 Storage
              </DialogTitle>
              <DialogDescription>Connect compatible S3 storage</DialogDescription>
            </div>
            <DialogClose asChild>
              <Button
                variant="ghost"
                className="font-medium text-foreground/80 text-sm hover:bg-transparent hover:text-foreground"
              >
                CLOSE
              </Button>
            </DialogClose>
          </div>
        </div>

        <form onSubmit={handleSave} className="space-y-4 px-5 pt-2 pb-5">
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-1">
              <Label className="font-semibold text-foreground text-xs uppercase tracking-wider">
                Name *
              </Label>
              <Input
                placeholder="My R2 Backup"
                value={name}
                onChange={(e) => setName(e.target.value)}
                required
                className="border/50 h-10 rounded-lg bg-background/50 text-sm focus:border-primary focus:ring-0"
              />
            </div>
            <div className="space-y-1">
              <Label className="font-semibold text-foreground text-xs uppercase tracking-wider">
                Provider
              </Label>
              <Select value={provider} onValueChange={handleProviderChange}>
                <SelectTrigger className="border/50 h-10 rounded-lg bg-background/50 text-sm focus:ring-0">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="r2">Cloudflare R2</SelectItem>
                  <SelectItem value="s3">AWS S3</SelectItem>
                  <SelectItem value="minio">MinIO</SelectItem>
                  <SelectItem value="wasabi">Wasabi</SelectItem>
                  <SelectItem value="b2">Backblaze B2</SelectItem>
                  <SelectItem value="other">Other S3 Compatible</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>

          <div className="space-y-1">
            <Label className="font-semibold text-foreground text-xs uppercase tracking-wider">
              Description
            </Label>
            <Input
              placeholder="Production backups bucket"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              className="border/50 h-10 rounded-lg bg-background/50 text-sm focus:border-primary focus:ring-0"
            />
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-1">
              <Label className="font-semibold text-foreground text-xs uppercase tracking-wider">
                Endpoint URL *
              </Label>
              <Input
                placeholder="https://..."
                value={endpoint}
                onChange={(e) => setEndpoint(e.target.value)}
                required
                className="border/50 h-10 rounded-lg bg-background/50 font-mono text-sm focus:border-primary focus:ring-0"
              />
            </div>
            <div className="space-y-1">
              <Label className="font-semibold text-foreground text-xs uppercase tracking-wider">
                Bucket Name *
              </Label>
              <Input
                placeholder="my-backups"
                value={bucket}
                onChange={(e) => setBucket(e.target.value)}
                required
                className="border/50 h-10 rounded-lg bg-background/50 font-mono text-sm focus:border-primary focus:ring-0"
              />
            </div>
          </div>

          <div className="space-y-1">
            <Label className="font-semibold text-foreground text-xs uppercase tracking-wider">
              Region
            </Label>
            <Input
              placeholder="us-east-1 or auto"
              value={region}
              onChange={(e) => setRegion(e.target.value)}
              className="border/50 h-10 rounded-lg bg-background/50 font-mono text-sm focus:border-primary focus:ring-0"
            />
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-1">
              <Label className="font-semibold text-foreground text-xs uppercase tracking-wider">
                Access Key ID *
              </Label>
              <Input
                placeholder="AKIA..."
                value={accessKeyId}
                onChange={(e) => setAccessKeyId(e.target.value)}
                required
                className="border/50 h-10 rounded-lg bg-background/50 font-mono text-sm focus:border-primary focus:ring-0"
              />
            </div>
            <div className="space-y-1">
              <Label className="font-semibold text-foreground text-xs uppercase tracking-wider">
                Secret Access Key *
              </Label>
              <div className="relative">
                <Input
                  type={showSecret ? 'text' : 'password'}
                  placeholder="••••••••"
                  value={secretAccessKey}
                  onChange={(e) => setSecretAccessKey(e.target.value)}
                  required
                  className="border/50 h-10 rounded-lg bg-background/50 pr-10 font-mono text-sm focus:border-primary focus:ring-0"
                />
                <Button
                  type="button"
                  variant="ghost"
                  size="sm"
                  onClick={() => setShowSecret(!showSecret)}
                  className="absolute top-1/2 right-1 -translate-y-1/2 text-muted-foreground hover:text-foreground"
                >
                  {showSecret ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                </Button>
              </div>
            </div>
          </div>

          <div className="border/50 flex items-center gap-2 rounded-lg bg-muted/40 p-3 text-muted-foreground text-xs">
            <Info className="h-4 w-4 shrink-0 text-primary" />
            <span>Credentials are safely encrypted before being saved into your database.</span>
          </div>

          <div className="flex justify-end gap-3 pt-2">
            <DialogClose asChild>
              <Button type="button" variant="outline">
                Cancel
              </Button>
            </DialogClose>
            <Button type="submit" disabled={createS3Dest.isPending}>
              {createS3Dest.isPending ? 'Saving…' : 'Save S3 Destination'}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}
