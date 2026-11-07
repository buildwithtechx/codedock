import { describe, expect, it } from 'vitest';
import { BillingGateErrorCodes, RollbackErrorCodes } from '../src/errors';
import type {
  Deployment,
  ProjectConfig,
  RollbackCapacityUI,
  UpdateRollbackSettingsRequest,
} from '../src/index';

describe('@codedock/contracts', () => {
  it('exports standard rollback error codes', () => {
    expect(RollbackErrorCodes.NOT_READY).toBe('ROLLBACK_NOT_READY');
    expect(RollbackErrorCodes.ALREADY_ACTIVE).toBe('ROLLBACK_ALREADY_ACTIVE');
    expect(RollbackErrorCodes.ARTIFACT_GONE).toBe('ROLLBACK_ARTIFACT_GONE');
    expect(RollbackErrorCodes.CONFIGURATION_REDACTED).toBe('ROLLBACK_CONFIGURATION_REDACTED');
  });

  it('exports standard billing gate error codes', () => {
    expect(BillingGateErrorCodes.GATE_REQUIRED).toBe('BILLING_GATE_REQUIRED');
    expect(BillingGateErrorCodes.PAYMENT_REQUIRED).toBe('PAYMENT_REQUIRED');
  });

  it('types compile and validate interface shapes', () => {
    const deployment: Partial<Deployment> = {
      id: 'dep-123',
      status: 'SUCCESS',
    };
    const config: Partial<ProjectConfig> = {
      id: 'proj-123',
      rollbackWindow: 5,
    };
    const capacity: Partial<RollbackCapacityUI> = {
      window: 5,
      maxWindow: 20,
      source: 'instance-default',
    };
    const req: UpdateRollbackSettingsRequest = {
      rollbackWindow: 7,
      defaultRollbackStrategy: 'snapshot',
    };

    expect(deployment.id).toBe('dep-123');
    expect(config.rollbackWindow).toBe(5);
    expect(capacity.maxWindow).toBe(20);
    expect(req.rollbackWindow).toBe(7);
  });
});
