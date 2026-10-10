import { useQueryClient } from '@tanstack/react-query';
import { Box, Search } from 'lucide-react';
import type React from 'react';
import { useMemo, useState } from 'react';
import { toast } from 'sonner';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '#/components/ui/dialog';
import { Input } from '#/components/ui/input';
import { appsService } from '#/services/apps';
import { deploymentsService } from '#/services/deployments';
import { COMPANION_CATALOG, type CompanionCatalogItem } from './companion-catalog';
import { ServiceConfigForm } from './service-config-form';

interface AddServiceModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  projectId: string;
  environmentId: string;
  onCreated?: () => void;
}

export function AddServiceModal({
  open,
  onOpenChange,
  projectId,
  environmentId,
  onCreated,
}: AddServiceModalProps) {
  const queryClient = useQueryClient();
  const [mode, setMode] = useState<'catalog' | 'custom'>('catalog');
  const [selectedCompanion, setSelectedCompanion] = useState<CompanionCatalogItem | null>(null);
  const [category, setCategory] = useState<string>('all');
  const [search, setSearch] = useState('');
  const [name, setName] = useState('');
  const [image, setImage] = useState('');
  const [port, setPort] = useState(8080);
  const [volume, setVolume] = useState('');
  const [envRows, setEnvRows] = useState<{ key: string; value: string }[]>([]);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const resetForm = () => {
    setSelectedCompanion(null);
    setName('');
    setImage('');
    setPort(8080);
    setVolume('');
    setEnvRows([]);
    setSearch('');
    setCategory('all');
    setMode('catalog');
  };

  const handleSelectCompanion = (item: CompanionCatalogItem) => {
    setSelectedCompanion(item);
    setName(item.id);
    setImage(item.image);
    setPort(item.port);
    setVolume(item.defaultVolume || '');
    setEnvRows(item.defaultEnv.map((e) => ({ ...e })));
  };

  const filteredCatalog = useMemo(() => {
    return COMPANION_CATALOG.filter((item) => {
      if (category !== 'all' && item.category !== category) return false;
      if (!search.trim()) return true;
      const q = search.toLowerCase();
      return (
        item.name.toLowerCase().includes(q) ||
        item.description.toLowerCase().includes(q) ||
        item.image.toLowerCase().includes(q)
      );
    });
  }, [category, search]);

  const handleAddEnvRow = () => {
    setEnvRows((prev) => [...prev, { key: '', value: '' }]);
  };

  const handleUpdateEnvRow = (index: number, field: 'key' | 'value', val: string) => {
    setEnvRows((prev) => {
      const copy = [...prev];
      copy[index] = { ...copy[index], [field]: val };
      return copy;
    });
  };

  const handleRemoveEnvRow = (index: number) => {
    setEnvRows((prev) => prev.filter((_, i) => i !== index));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim() || !image.trim() || !environmentId || !projectId) {
      toast.error('Service name and Docker image are required');
      return;
    }

    setIsSubmitting(true);
    try {
      const appRes = await appsService.createApp(environmentId, {
        projectId,
        name: name.trim(),
        repositoryUrl: '',
        imageRef: image.trim(),
        branch: 'main',
        rootDirectory: '',
        runtimeMode: 'web',
        installCommand: '',
        buildCommand: '',
        startCommand: '',
        dockerfilePath: '',
        buildEngine: 'dockerfile',
        internalPort: Number(port) || 80,
        domain: '',
        staticOutput: '',
        healthCheckPath: '',
      });

      const serviceId = appRes?.data?.id;

      if (serviceId) {
        for (const row of envRows) {
          if (row.key.trim()) {
            try {
              await appsService.createVariable(serviceId, {
                key: row.key.trim(),
                value: row.value.trim(),
                isSecret: false,
              });
            } catch {}
          }
        }

        try {
          await deploymentsService.trigger(serviceId);
        } catch {}
      }

      await queryClient.invalidateQueries({ queryKey: ['apps'] });
      await queryClient.invalidateQueries({ queryKey: ['canvas'] });
      toast.success(`Service "${name}" added successfully`);
      onOpenChange(false);
      resetForm();
      onCreated?.();
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Failed to add service';
      toast.error(msg);
    } finally {
      setIsSubmitting(false);
    }
  };

  const isConfiguring = selectedCompanion !== null || mode === 'custom';

  return (
    <Dialog
      open={open}
      onOpenChange={(next) => {
        if (!next) resetForm();
        onOpenChange(next);
      }}
    >
      <DialogContent className="max-h-[90vh] overflow-y-auto sm:max-w-2xl">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <Box className="size-5 text-primary" />
            Add Service
          </DialogTitle>
          <DialogDescription>
            Add a companion database, cache, or custom Docker container to this project.
          </DialogDescription>
        </DialogHeader>

        {!isConfiguring ? (
          <div className="space-y-4">
            <div className="flex items-center justify-between gap-2 border-border/50 border-b pb-3">
              <div className="flex rounded-lg bg-muted/50 p-1 text-xs">
                <button
                  type="button"
                  onClick={() => setMode('catalog')}
                  className={`rounded-md px-3 py-1 font-medium transition-colors ${
                    mode === 'catalog'
                      ? 'bg-background text-foreground shadow-sm'
                      : 'text-muted-foreground hover:text-foreground'
                  }`}
                >
                  Curated Catalog
                </button>
                <button
                  type="button"
                  onClick={() => {
                    setMode('custom');
                    setName('');
                    setImage('');
                    setPort(8080);
                  }}
                  className={`rounded-md px-3 py-1 font-medium transition-colors ${
                    mode === 'custom'
                      ? 'bg-background text-foreground shadow-sm'
                      : 'text-muted-foreground hover:text-foreground'
                  }`}
                >
                  Custom Docker Image
                </button>
              </div>

              <div className="relative w-48">
                <Search className="pointer-events-none absolute top-1/2 left-2.5 size-3.5 -translate-y-1/2 text-muted-foreground" />
                <Input
                  value={search}
                  onChange={(e) => setSearch(e.target.value)}
                  placeholder="Search..."
                  className="h-8 pl-8 text-xs"
                />
              </div>
            </div>

            <div className="flex flex-wrap gap-1">
              {(['all', 'database', 'cache', 'queue', 'search', 'tools'] as const).map((cat) => (
                <button
                  key={cat}
                  type="button"
                  onClick={() => setCategory(cat)}
                  className={`rounded-full px-2.5 py-0.5 font-medium text-xs capitalize transition-colors ${
                    category === cat
                      ? 'bg-primary text-primary-foreground'
                      : 'bg-muted/60 text-muted-foreground hover:text-foreground'
                  }`}
                >
                  {cat}
                </button>
              ))}
            </div>

            <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
              {filteredCatalog.map((item) => (
                <button
                  key={item.id}
                  type="button"
                  onClick={() => handleSelectCompanion(item)}
                  className="group flex flex-col justify-between rounded-xl border border-border/60 bg-card p-4 text-left transition-all hover:border-primary/50 hover:shadow-sm"
                >
                  <div>
                    <div className="flex items-center justify-between">
                      <p className="font-semibold text-foreground text-sm">{item.name}</p>
                      <span className="rounded-md bg-muted px-1.5 py-0.5 font-mono text-[10px] text-muted-foreground">
                        :{item.port}
                      </span>
                    </div>
                    <p className="mt-1 line-clamp-2 text-muted-foreground text-xs">
                      {item.description}
                    </p>
                  </div>
                  <div className="mt-3 flex items-center justify-between border-border/40 border-t pt-2 font-mono text-[11px] text-muted-foreground">
                    <span className="max-w-[160px] truncate">{item.image}</span>
                    <span className="text-primary group-hover:underline">Configure →</span>
                  </div>
                </button>
              ))}
            </div>
          </div>
        ) : (
          <ServiceConfigForm
            name={name}
            onNameChange={setName}
            image={image}
            onImageChange={setImage}
            port={port}
            onPortChange={setPort}
            volume={volume}
            onVolumeChange={setVolume}
            envRows={envRows}
            onAddEnvRow={handleAddEnvRow}
            onUpdateEnvRow={handleUpdateEnvRow}
            onRemoveEnvRow={handleRemoveEnvRow}
            onBack={() => {
              setSelectedCompanion(null);
              if (mode === 'custom') setMode('catalog');
            }}
            onCancel={() => {
              onOpenChange(false);
              resetForm();
            }}
            onSubmit={handleSubmit}
            isSubmitting={isSubmitting}
          />
        )}
      </DialogContent>
    </Dialog>
  );
}
