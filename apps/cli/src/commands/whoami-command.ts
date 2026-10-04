import { ApiClient } from '../api-client.js';
import { printError, printJson } from '../output-format.js';
import type { CliContext, UserProfile } from '../types.js';

export async function whoamiCommand(ctx: CliContext): Promise<void> {
  if (!ctx.token) {
    printError('Not logged in. Please run "codedock login" first.');
    process.exit(1);
  }

  const client = new ApiClient(ctx);
  try {
    const res = await client.get<UserProfile>('/api/auth/me');
    const user = res.data;
    if (!user) {
      throw new Error('User profile empty');
    }

    if (ctx.json) {
      printJson(user);
      return;
    }

    console.log('\x1b[1mLogged in as:\x1b[0m');
    console.log(`  Name:  ${user.name || '(unnamed)'}`);
    console.log(`  Email: ${user.email}`);
    console.log(`  Role:  ${user.role}`);
    if (user.planType) {
      console.log(`  Plan:  ${user.planType}`);
    }
    console.log(`  Host:  ${ctx.serverUrl}`);
  } catch (err: unknown) {
    const msg = err instanceof Error ? err.message : String(err);
    printError(`Failed to fetch user profile: ${msg}`);
    process.exit(1);
  }
}
