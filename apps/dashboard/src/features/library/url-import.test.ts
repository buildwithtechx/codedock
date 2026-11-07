import { describe, expect, it } from 'vitest';
import { parseGitUrl } from './url-import';

describe('parseGitUrl', () => {
  it('parses https github repository without .git', () => {
    expect(parseGitUrl('https://github.com/owner/repo')).toEqual({
      repositoryUrl: 'https://github.com/owner/repo',
      branch: 'main',
      name: 'repo',
    });
  });

  it('preserves http protocol and .git suffix', () => {
    expect(parseGitUrl('http://git.internal/org/project.git')).toEqual({
      repositoryUrl: 'http://git.internal/org/project.git',
      branch: 'main',
      name: 'project',
    });
  });

  it('preserves https protocol with .git suffix', () => {
    expect(parseGitUrl('https://gitlab.com/group/service.git')).toEqual({
      repositoryUrl: 'https://gitlab.com/group/service.git',
      branch: 'main',
      name: 'service',
    });
  });

  it('parses scp git syntax', () => {
    expect(parseGitUrl('git@github.com:owner/repo.git')).toEqual({
      repositoryUrl: 'git@github.com:owner/repo.git',
      branch: 'main',
      name: 'repo',
    });
  });

  it('rejects URLs with usernames or passwords', () => {
    expect(parseGitUrl('https://user:pass@github.com/owner/repo')).toBeNull();
  });

  it('rejects URLs with invalid depths', () => {
    expect(parseGitUrl('https://github.com/owner/repo/sub')).toBeNull();
    expect(parseGitUrl('https://github.com/owner')).toBeNull();
  });

  it('rejects unsupported protocols', () => {
    expect(parseGitUrl('ssh://git@github.com/owner/repo.git')).toBeNull();
    expect(parseGitUrl('ftp://github.com/owner/repo')).toBeNull();
  });
});
