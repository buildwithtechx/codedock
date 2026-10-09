import { useMemo } from 'react';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '#/components/ui/tabs';
import type { OneClickEnvVar } from '#/interfaces/templates';

export interface InstallVariableGroup {
  id: string;
  label: string;
  variables: OneClickEnvVar[];
}

export function groupInstallVariables(variables: OneClickEnvVar[]): InstallVariableGroup[] {
  const settings = variables.filter((variable) => !variable.secret && variable.input);
  const secrets = variables.filter((variable) => variable.secret);
  const advanced = variables.filter((variable) => !variable.secret && !variable.input);
  return [
    { id: 'settings', label: 'Settings', variables: settings },
    { id: 'secrets', label: 'Secrets', variables: secrets },
    { id: 'advanced', label: 'Environment overrides', variables: advanced },
  ].filter((group) => group.variables.length > 0);
}

export function findMissingRequired(
  variables: OneClickEnvVar[],
  values: Record<string, string>
): OneClickEnvVar[] {
  return variables.filter(
    (variable) =>
      variable.required && !(values[variable.key] ?? '').trim() && !variable.defaultValue
  );
}

export function InstallSettingsTabs({
  variables,
  values,
  onChange,
  disabled = false,
}: {
  variables: OneClickEnvVar[];
  values: Record<string, string>;
  onChange: (key: string, value: string) => void;
  disabled?: boolean;
}) {
  const groups = useMemo(() => groupInstallVariables(variables), [variables]);
  if (groups.length === 0) return null;
  if (groups.length === 1) {
    return (
      <div className="space-y-4">
        {groups[0].variables.map((variable) => (
          <VariableField
            key={variable.key}
            variable={variable}
            value={values[variable.key] ?? ''}
            onChange={onChange}
            disabled={disabled}
          />
        ))}
      </div>
    );
  }
  return (
    <Tabs defaultValue={groups[0].id} className="w-full">
      <TabsList>
        {groups.map((group) => (
          <TabsTrigger key={group.id} value={group.id}>
            {group.label}
          </TabsTrigger>
        ))}
      </TabsList>
      {groups.map((group) => (
        <TabsContent key={group.id} value={group.id} className="mt-4 space-y-4">
          {group.variables.map((variable) => (
            <VariableField
              key={variable.key}
              variable={variable}
              value={values[variable.key] ?? ''}
              onChange={onChange}
              disabled={disabled}
            />
          ))}
        </TabsContent>
      ))}
    </Tabs>
  );
}

function VariableField({
  variable,
  value,
  onChange,
  disabled,
}: {
  variable: OneClickEnvVar;
  value: string;
  onChange: (key: string, value: string) => void;
  disabled: boolean;
}) {
  const placeholder = variable.secret ? 'Leave blank to generate' : (variable.defaultValue ?? '');
  return (
    <div className="space-y-2">
      <Label htmlFor={`install-${variable.key}`}>
        {variable.label || variable.key}
        {variable.required && !variable.defaultValue ? ' *' : ''}
      </Label>
      <Input
        id={`install-${variable.key}`}
        type={variable.secret ? 'password' : 'text'}
        autoComplete="off"
        placeholder={placeholder}
        value={value}
        disabled={disabled}
        onChange={(event) => onChange(variable.key, event.target.value)}
      />
    </div>
  );
}
