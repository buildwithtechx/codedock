package repositories

import "context"

func (r *DeploymentRepo) RecoverInterrupted(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `UPDATE deployments SET status='FAILED', finished_at=CURRENT_TIMESTAMP, build_logs=COALESCE(build_logs,'')||chr(10)||'Daemon restarted during deployment. Retained rollout containers have been reconciled; inspect runtime before retrying.' WHERE status IN ('pending','PREPARING','CLONING','PULLING','BUILDING','STARTING','READINESS','ROUTING')`)
	return err
}
