import { Eye, EyeOff, Plus, Save, Trash2 } from 'lucide-react';
import { useState } from 'react';
import { Button } from '#/components/ui/button';
import { Card } from '#/components/ui/card';
import { Input } from '#/components/ui/input';
import { useGetVars, useSetVars } from '#/features/projects';

interface ProjectVariablesTabProps {
  projectId: string;
}

export function ProjectVariablesTab({ projectId }: ProjectVariablesTabProps) {
  const { data: varsRes, isLoading, refetch } = useGetVars(projectId);
  const setVarsMutation = useSetVars();

  const [entries, setEntries] = useState<Array<{ key: string; value: string; show?: boolean }>>([]);
  const [initialized, setInitialized] = useState(false);
  const [newKey, setNewKey] = useState('');
  const [newValue, setNewValue] = useState('');

  if (varsRes?.data && !initialized && !isLoading) {
    const rawVars = varsRes.data || {};
    const parsed = Object.entries(rawVars).map(([key, value]) => ({
      key,
      value: String(value),
      show: false,
    }));
    setEntries(parsed);
    setInitialized(true);
  }

  const handleSave = async (updated: Array<{ key: string; value: string }>) => {
    const variables: Record<string, string> = {};
    for (const item of updated) {
      if (item.key.trim()) {
        variables[item.key.trim()] = item.value;
      }
    }
    try {
      await setVarsMutation.mutateAsync({
        id: projectId,
        payload: { variables },
      });
      refetch();
    } catch {
      // Handled
    }
  };

  const handleAdd = () => {
    if (!newKey.trim()) return;
    const updated = [...entries, { key: newKey.trim(), value: newValue, show: false }];
    setEntries(updated);
    setNewKey('');
    setNewValue('');
    void handleSave(updated);
  };

  const handleDelete = (index: number) => {
    const updated = entries.filter((_, i) => i !== index);
    setEntries(updated);
    void handleSave(updated);
  };

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <div>
          <h3 className="font-semibold text-foreground/90 text-sm">
            Project Environment Variables
          </h3>
          <p className="mt-0.5 text-muted-foreground text-xs">
            Variables configured here are inherited by all workloads and services within this
            project.
          </p>
        </div>
      </div>

      <Card className="space-y-4 p-4">
        <div className="space-y-2.5">
          {entries.length === 0 ? (
            <p className="py-2 text-center text-muted-foreground text-xs">
              No project environment variables configured yet.
            </p>
          ) : (
            entries.map((item, index) => (
              <div key={item.key} className="flex items-center gap-2">
                <Input className="w-1/3 font-mono text-xs" value={item.key} disabled />
                <div className="relative flex-1">
                  <Input
                    className="pr-9 font-mono text-xs"
                    type={item.show ? 'text' : 'password'}
                    value={item.value}
                    onChange={(e) => {
                      const updated = [...entries];
                      updated[index].value = e.target.value;
                      setEntries(updated);
                    }}
                  />
                  <button
                    type="button"
                    onClick={() => {
                      const updated = [...entries];
                      updated[index].show = !updated[index].show;
                      setEntries(updated);
                    }}
                    className="absolute top-1/2 right-2.5 -translate-y-1/2 text-muted-foreground hover:text-foreground"
                  >
                    {item.show ? (
                      <EyeOff className="h-3.5 w-3.5" />
                    ) : (
                      <Eye className="h-3.5 w-3.5" />
                    )}
                  </button>
                </div>
                <Button
                  size="icon"
                  variant="ghost"
                  className="h-9 w-9 shrink-0 text-muted-foreground hover:text-destructive"
                  onClick={() => handleDelete(index)}
                >
                  <Trash2 className="h-4 w-4" />
                </Button>
              </div>
            ))
          )}
        </div>

        <div className="border-border/60 border-t pt-3">
          <div className="flex items-center gap-2">
            <Input
              placeholder="VARIABLE_NAME"
              className="w-1/3 font-mono text-xs uppercase"
              value={newKey}
              onChange={(e) => setNewKey(e.target.value.toUpperCase())}
            />
            <Input
              placeholder="Value"
              className="flex-1 font-mono text-xs"
              value={newValue}
              onChange={(e) => setNewValue(e.target.value)}
            />
            <Button
              size="sm"
              className="h-9 shrink-0 gap-1.5 text-xs"
              onClick={handleAdd}
              disabled={!newKey.trim() || setVarsMutation.isPending}
            >
              <Plus className="h-3.5 w-3.5" />
              Add
            </Button>
          </div>
        </div>

        {entries.length > 0 && (
          <div className="flex justify-end border-border/60 border-t pt-2">
            <Button
              size="sm"
              className="gap-1.5 text-xs"
              onClick={() => void handleSave(entries)}
              disabled={setVarsMutation.isPending}
            >
              <Save className="h-3.5 w-3.5" />
              {setVarsMutation.isPending ? 'Saving...' : 'Save Changes'}
            </Button>
          </div>
        )}
      </Card>
    </div>
  );
}
