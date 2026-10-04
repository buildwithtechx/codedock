import { appsCommand } from './commands/apps-command.js';
import { databasesCommand } from './commands/databases-command.js';
import { deployCommand } from './commands/deploy-command.js';
import { envCommand } from './commands/env-command.js';
import { loginCommand } from './commands/login-command.js';
import { logoutCommand } from './commands/logout-command.js';
import { projectsCommand } from './commands/projects-command.js';
import { serversCommand } from './commands/servers-command.js';
import { statusCommand } from './commands/status-command.js';
import { CLI_VERSION, versionCommand } from './commands/version-command.js';
import { whoamiCommand } from './commands/whoami-command.js';
import { loadConfig } from './config-store.js';
import { printError } from './output-format.js';
import type { CliContext } from './types.js';

function printHelp(): void {
  console.log(`
\x1b[1mcodedock\x1b[0m v${CLI_VERSION} - Codedock PaaS Command Line Interface

\x1b[1mUSAGE\x1b[0m
  codedock <command> [subcommand] [arguments...] [options...]

\x1b[1mCOMMANDS\x1b[0m
  login                      Authenticate with your Codedock server
  logout                     Clear saved authentication credentials
  whoami, me                 Display current authenticated user
  status <service-id>        Check application service status and deployments
  servers                    Manage server nodes (list, create, delete)
  projects                   Manage projects (list, create, delete)
  apps                       Manage applications/services (list, create, delete, logs, deployments)
  env                        Manage project environment variables (list, get, set)
  db                         Manage databases and trigger backups (list, create, backup)
  deploy [path|service-id]   Deploy local directory or trigger deployment for a service
  version                    Print CLI version

\x1b[1mGLOBAL OPTIONS\x1b[0m
  --server <url>             Override Codedock server URL
  --token <token>            Override authentication token
  --json                     Output raw JSON
  -h, --help                 Display this help message
  -v, --version              Print CLI version

\x1b[1mEXAMPLES\x1b[0m
  codedock login --server http://localhost:8080 --email admin@example.com
  codedock whoami
  codedock servers list
  codedock projects create "My New API"
  codedock env set <project-id> PORT=8080 NODE_ENV=production
  codedock deploy .
`);
}

function parseCliArgs(rawArgs: string[]): {
  command?: string;
  positional: string[];
  options: Record<string, string | boolean>;
} {
  const positional: string[] = [];
  const options: Record<string, string | boolean> = {};

  let i = 0;
  while (i < rawArgs.length) {
    const arg = rawArgs[i];
    if (arg === '--help' || arg === '-h') {
      options.help = true;
      i++;
    } else if (arg === '--version' || arg === '-v') {
      options.version = true;
      i++;
    } else if (arg === '--json') {
      options.json = true;
      i++;
    } else if (arg === '--is-local') {
      options.isLocal = true;
      i++;
    } else if (arg.startsWith('--')) {
      const key = arg.slice(2);
      const next = rawArgs[i + 1];
      if (next && !next.startsWith('-')) {
        options[key] = next;
        i += 2;
      } else {
        options[key] = true;
        i++;
      }
    } else {
      positional.push(arg);
      i++;
    }
  }

  const command = positional[0]?.toLowerCase();
  return {
    command,
    positional: positional.slice(1),
    options,
  };
}

export async function executeCli(rawArgs: string[]): Promise<void> {
  const { command, positional, options } = parseCliArgs(rawArgs);
  const cfg = loadConfig();

  const serverUrl =
    (typeof options.server === 'string' ? options.server : '') ||
    process.env.CODEDOCK_SERVER_URL ||
    cfg.serverUrl ||
    'http://localhost:8080';

  const token =
    (typeof options.token === 'string' ? options.token : '') ||
    process.env.CODEDOCK_TOKEN ||
    cfg.token;

  const ctx: CliContext = {
    config: cfg,
    serverUrl,
    token,
    json: Boolean(options.json),
  };

  if (options.version) {
    await versionCommand(ctx);
    return;
  }

  if (options.help || !command) {
    printHelp();
    return;
  }

  switch (command) {
    case 'login':
      await loginCommand(ctx, positional, {
        server: typeof options.server === 'string' ? options.server : undefined,
        email: typeof options.email === 'string' ? options.email : undefined,
        password: typeof options.password === 'string' ? options.password : undefined,
        token: typeof options.token === 'string' ? options.token : undefined,
        totp: typeof options.totp === 'string' ? options.totp : undefined,
      });
      break;

    case 'logout':
      await logoutCommand(ctx);
      break;

    case 'whoami':
    case 'me':
      await whoamiCommand(ctx);
      break;

    case 'status':
      await statusCommand(ctx, positional);
      break;

    case 'servers':
    case 'server':
      await serversCommand(ctx, positional, {
        name: typeof options.name === 'string' ? options.name : undefined,
        ip: typeof options.ip === 'string' ? options.ip : undefined,
        port: typeof options.port === 'string' ? options.port : undefined,
        user: typeof options.user === 'string' ? options.user : undefined,
        key: typeof options.key === 'string' ? options.key : undefined,
        isLocal: Boolean(options.isLocal),
      });
      break;

    case 'projects':
    case 'project':
      await projectsCommand(ctx, positional, {
        name: typeof options.name === 'string' ? options.name : undefined,
        description: typeof options.description === 'string' ? options.description : undefined,
      });
      break;

    case 'apps':
    case 'app':
      await appsCommand(ctx, positional, {
        name: typeof options.name === 'string' ? options.name : undefined,
        project: typeof options.project === 'string' ? options.project : undefined,
        repo: typeof options.repo === 'string' ? options.repo : undefined,
        branch: typeof options.branch === 'string' ? options.branch : undefined,
      });
      break;

    case 'env':
    case 'environments':
    case 'environment':
      await envCommand(ctx, positional);
      break;

    case 'db':
    case 'databases':
    case 'database':
      await databasesCommand(ctx, positional, {
        name: typeof options.name === 'string' ? options.name : undefined,
        type: typeof options.type === 'string' ? options.type : undefined,
        server: typeof options.server === 'string' ? options.server : undefined,
      });
      break;

    case 'deploy':
      await deployCommand(ctx, positional, {
        project: typeof options.project === 'string' ? options.project : undefined,
        branch: typeof options.branch === 'string' ? options.branch : undefined,
      });
      break;

    case 'version':
      await versionCommand(ctx);
      break;

    default:
      printError(`Unknown command: "${command}". Run "codedock --help" for available commands.`);
      process.exit(1);
  }
}
