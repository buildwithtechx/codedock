const http = require('node:http');
const fs = require('node:fs');
const path = require('node:path');
const crypto = require('node:crypto');

const PORT = process.env.PORT || 80;
const SCRIPTS = new Set(['/install.sh', '/upgrade.sh', '/cli', '/base.sh', '/worker.sh']);
const POSTHOG_KEY = process.env.POSTHOG_API_KEY;
const POSTHOG_HOST = process.env.POSTHOG_HOST || 'https://us.i.posthog.com';

function getVersion() {
  if (process.env.CODEDOCK_VERSION) return process.env.CODEDOCK_VERSION;
  if (process.env.VERSION) return process.env.VERSION;
  try {
    const pkgPath = path.resolve(__dirname, '../package.json');
    if (fs.existsSync(pkgPath)) {
      const pkg = JSON.parse(fs.readFileSync(pkgPath, 'utf8'));
      if (pkg.version) return pkg.version;
    }
  } catch (e) {
    // ignore
  }
  return '0.1.0';
}

function hashValue(value) {
  const salt = process.env.POSTHOG_DISTINCT_ID_SALT || POSTHOG_KEY || 'codedock-install';
  return crypto.createHash('sha256').update(`${salt}:${value}`).digest('hex').slice(0, 32);
}

function trackEvent(eventName, req) {
  if (!POSTHOG_KEY) return;
  const rawIp = req.headers['x-forwarded-for'] || req.socket.remoteAddress || '';
  const ip = typeof rawIp === 'string' ? rawIp.split(',')[0].trim() : '';
  const userAgent = req.headers['user-agent'] || 'unknown';
  const distinctId = `installer:${hashValue(`${ip}:${userAgent}`)}`;

  fetch(`${POSTHOG_HOST}/capture/`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      api_key: POSTHOG_KEY,
      event: eventName,
      distinct_id: distinctId,
      properties: {
        $ip: ip,
        $user_agent: userAgent,
        $current_url: req.url,
        version: getVersion(),
        script_requested: req.url === '/' ? '/install.sh' : req.url,
      },
    }),
  }).catch((err) => console.error('PostHog error:', err));
}

const server = http.createServer((req, res) => {
  if (req.url === '/version') {
    res.writeHead(200, { 'Content-Type': 'application/json' });
    res.end(JSON.stringify({ version: getVersion(), latest: true }));
    return;
  }

  let filename = path.basename(req.url);
  if (req.url === '/') filename = 'install.sh';
  if (req.url === '/cli') filename = 'install-cli.sh';

  const file = path.join(__dirname, filename);
  const name = path.basename(req.url === '/' ? '/install.sh' : req.url);

  if (!SCRIPTS.has(req.url === '/' ? '/install.sh' : req.url) || !fs.existsSync(file)) {
    res.writeHead(404);
    res.end('Not found');
    return;
  }

  res.writeHead(200, { 'Content-Type': 'text/x-shellscript' });
  fs.createReadStream(file).pipe(res);
  trackEvent('script_downloaded', req);
});

server.listen(PORT, () => {
  console.log(`📦 install server running on port ${PORT}`);
});
