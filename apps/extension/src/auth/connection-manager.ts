import * as vscode from 'vscode';
import type { CodedockConnection } from '../interfaces';

const SECRET_TOKEN_PREFIX = 'codedock:token:';
const STATE_PROFILES_KEY = 'codedock:profiles';
const STATE_ACTIVE_PROFILE_KEY = 'codedock:active_profile';

export class ConnectionManager {
  private readonly _onDidChangeConnection = new vscode.EventEmitter<CodedockConnection | null>();
  readonly onDidChangeConnection = this._onDidChangeConnection.event;

  constructor(private readonly context: vscode.ExtensionContext) {}

  async getActiveConnection(): Promise<CodedockConnection | null> {
    const activeId = this.context.globalState.get<string>(STATE_ACTIVE_PROFILE_KEY);
    if (!activeId) {
      return null;
    }

    const profiles = this.context.globalState.get<
      Record<string, Omit<CodedockConnection, 'token'>>
    >(STATE_PROFILES_KEY, {});
    const profile = profiles[activeId];
    if (!profile) {
      return null;
    }

    const token = await this.context.secrets.get(`${SECRET_TOKEN_PREFIX}${activeId}`);
    if (!token) {
      return null;
    }

    return {
      ...profile,
      token,
    };
  }

  async listProfiles(): Promise<
    Array<{ id: string; name: string; url: string; isActive: boolean }>
  > {
    const activeId = this.context.globalState.get<string>(STATE_ACTIVE_PROFILE_KEY);
    const profiles = this.context.globalState.get<
      Record<string, Omit<CodedockConnection, 'token'>>
    >(STATE_PROFILES_KEY, {});

    return Object.entries(profiles).map(([id, p]) => ({
      id,
      name: p.instanceName || p.url,
      url: p.url,
      isActive: id === activeId,
    }));
  }

  async saveConnection(conn: CodedockConnection, profileId: string = 'default'): Promise<void> {
    const profiles = this.context.globalState.get<
      Record<string, Omit<CodedockConnection, 'token'>>
    >(STATE_PROFILES_KEY, {});

    const { token, ...meta } = conn;
    profiles[profileId] = meta;

    await this.context.globalState.update(STATE_PROFILES_KEY, profiles);
    await this.context.secrets.store(`${SECRET_TOKEN_PREFIX}${profileId}`, token);
    await this.context.globalState.update(STATE_ACTIVE_PROFILE_KEY, profileId);

    this._onDidChangeConnection.fire(conn);
  }

  async switchProfile(profileId: string): Promise<boolean> {
    const profiles = this.context.globalState.get<
      Record<string, Omit<CodedockConnection, 'token'>>
    >(STATE_PROFILES_KEY, {});
    if (!profiles[profileId]) {
      return false;
    }

    await this.context.globalState.update(STATE_ACTIVE_PROFILE_KEY, profileId);
    const conn = await this.getActiveConnection();
    this._onDidChangeConnection.fire(conn);
    return true;
  }

  async clearActiveConnection(): Promise<void> {
    const activeId = this.context.globalState.get<string>(STATE_ACTIVE_PROFILE_KEY);
    if (activeId) {
      const profiles = this.context.globalState.get<
        Record<string, Omit<CodedockConnection, 'token'>>
      >(STATE_PROFILES_KEY, {});
      delete profiles[activeId];
      await this.context.globalState.update(STATE_PROFILES_KEY, profiles);
      await this.context.secrets.delete(`${SECRET_TOKEN_PREFIX}${activeId}`);
      await this.context.globalState.update(STATE_ACTIVE_PROFILE_KEY, undefined);
    }

    this._onDidChangeConnection.fire(null);
  }
}
