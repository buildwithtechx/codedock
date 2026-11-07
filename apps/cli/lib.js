export function mapPlatform(platform, arch) {
  const os = { darwin: 'darwin', linux: 'linux', win32: 'windows' }[platform];
  const cpu = { x64: 'amd64', arm64: 'arm64' }[arch];
  if (!os || !cpu) {
    throw new Error(`Unsupported platform: ${platform}/${arch}`);
  }
  return { os, arch: cpu, binary: os === 'windows' ? 'codedock.exe' : 'codedock' };
}

export function artifactName(os, arch) {
  return `codedock_${os}_${arch}.tar.gz`;
}

export function releaseTag(version) {
  return `v${version}`;
}

export function downloadUrl(version, os, arch) {
  return `https://github.com/buildwithtechx/codedock/releases/download/${releaseTag(version)}/${artifactName(os, arch)}`;
}
