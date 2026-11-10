package bare

import (
	"bytes"
	"context"
	"io"

	appsmod "codedock/internal/modules/apps"
)

func (r *Runtime) BackupData(ctx context.Context, app *appsmod.AppService, target appsmod.RuntimeTarget, output io.Writer) error {
	script, err := ownedScript(app)
	if err != nil {
		return err
	}
	script += "test -d \"$base/data\"\ntar -czf - -C \"$base/data\" .\n"
	return r.runner.StreamHost(ctx, toClusterNode(target.BareNode), script, nil, output)
}

func (r *Runtime) RestoreData(ctx context.Context, app *appsmod.AppService, target appsmod.RuntimeTarget, archive io.Reader) error {
	script, err := ownedScript(app)
	if err != nil {
		return err
	}
	script += "test ! -e \"$base/pending\"\ntest -d \"$base/data\"\ntest \"$(systemctl show \"$unit\" --property=ActiveState --value)\" != active\ntar -xzf - -C \"$base/data\"\nchown -R \"$user\":\"$user\" \"$base/data\"\nchmod 750 \"$base/data\"\n"
	var discard bytes.Buffer
	return r.runner.StreamHost(ctx, toClusterNode(target.BareNode), script, archive, &discard)
}
