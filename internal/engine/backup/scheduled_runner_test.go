package backup
import (
 "context"
 "testing"
 "time"
 "codedock.run/codedock/internal/models"
)
func TestCronDispatchesToDurableRunnerWithPolicyTimeout(t *testing.T) {
 manager:=NewBackupManager(nil,nil,t.TempDir());called:=false
 manager.SetScheduledRunner(func(ctx context.Context,id string)error {called=true;if id!="policy" {t.Fatal("wrong scheduled policy")};deadline,ok:=ctx.Deadline();if !ok || time.Until(deadline)>5*time.Second {t.Fatal("policy timeout not applied")};return nil})
 cfg:=&models.BackupConfig{ID:"policy",Schedule:"0 2 * * *",Timeout:5,Status:models.BackupConfigStatusActive,BackupEnabled:true}
 if err:=manager.RegisterBackup(cfg);err!=nil {t.Fatal(err)}
 entry:=manager.cronEngine.Entry(manager.entries[cfg.ID]);if entry.Job==nil {t.Fatal("missing cron job")};entry.Job.Run();if !called {t.Fatal("scheduled backup bypassed durable runner")}
}
