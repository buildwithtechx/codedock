/** Fine-grained token capabilities for granular access control. */

export const CAPABILITY_AREAS = {
  domains: {
    root: 'project',
    resource: 'domain',
    verbs: { read: 'read', write: 'write', delete: 'admin' },
  },
  deployments: {
    root: 'project',
    resource: 'deployment',
    verbs: { read: 'read', trigger: 'write' },
  },
} as const;

type Areas = typeof CAPABILITY_AREAS;
export type CapabilityArea = keyof Areas;
export type Capability = {
  [A in CapabilityArea]: `${A}:${Extract<keyof Areas[A]['verbs'], string>}`;
}[CapabilityArea];
export type CoarseAction = 'read' | 'write' | 'admin';

interface CapabilityEntry {
  capability: Capability;
  area: CapabilityArea;
  root: string;
  resource: string;
  action: CoarseAction;
}

const ENTRIES: CapabilityEntry[] = (Object.keys(CAPABILITY_AREAS) as CapabilityArea[]).flatMap(
  (area) =>
    Object.entries(CAPABILITY_AREAS[area].verbs).map(([verb, action]) => ({
      capability: `${area}:${verb}` as Capability,
      area,
      root: CAPABILITY_AREAS[area].root,
      resource: CAPABILITY_AREAS[area].resource,
      action: action as CoarseAction,
    }))
);

const BY_CAPABILITY = new Map<string, CapabilityEntry>(ENTRIES.map((e) => [e.capability, e]));
const RANK: Record<CoarseAction, number> = { read: 1, write: 2, admin: 3 };

/** Every capability, in display order. */
export const CAPABILITIES: readonly Capability[] = ENTRIES.map((e) => e.capability);

/** Returns true if the string is a valid Capability token. */
export function isCapability(value: unknown): value is Capability {
  return typeof value === 'string' && BY_CAPABILITY.has(value);
}

/** Resolves capability for resource type and action. */
export function capabilityFor(resourceType: string, action: string): Capability | null {
  return (
    ENTRIES.find((e) => e.resource === resourceType && e.action === action)?.capability ?? null
  );
}

/** Returns the coarse action a capability maps to. */
export function capabilityAction(capability: Capability): CoarseAction {
  return BY_CAPABILITY.get(capability)!.action;
}

/** Returns capabilities for a root grant type. */
export function capabilitiesForRoot(rootType: string): Capability[] {
  return ENTRIES.filter((e) => e.root === rootType).map((e) => e.capability);
}

/** Returns capabilities implied by coarse permissions. */
export function capabilitiesImpliedBy(
  permissions: readonly string[],
  rootType = 'project'
): Capability[] {
  const level = Math.max(0, ...permissions.map((p) => RANK[p as CoarseAction] ?? 0));
  return ENTRIES.filter((e) => e.root === rootType && RANK[e.action] <= level).map(
    (e) => e.capability
  );
}
