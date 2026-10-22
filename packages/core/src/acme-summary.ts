export interface CertbotDiagnosis {
  domain: string;
  category:
    | 'gateway_error'
    | 'port_occupied'
    | 'firewall_timeout'
    | 'cloudflare_525'
    | 'dns_or_not_found'
    | 'generic';
  summary: string;
}

export const ACME_HTTP01_PORT = 49180;

export function summarizeCertbotFailure(output: string, domain: string): string {
  const text = output || '';

  const challengeUrlMatch = text.match(/https?:\/\/\S*\.well-known\/acme-challenge\/\S*/i);
  const challengeUrl = challengeUrlMatch ? challengeUrlMatch[0] : null;
  const redirectedToHttps = challengeUrl ? /^https:/i.test(challengeUrl) : false;

  const statusMatch = text.match(/invalid response from https?:\/\/\S+:\s*(\d{3})\b/i);
  const challengeStatus = statusMatch ? statusMatch[1] : null;

  if (
    /bind: address already in use|could not bind to IPv4 or IPv6|failed to bind to :80\b/i.test(
      text
    )
  ) {
    return (
      `Port 80 is occupied by another process on this server. Another web server or container ` +
      `is actively bound to :80. Stop the process holding port 80 before requesting a certificate.`
    );
  }

  if (challengeStatus === '502' || challengeStatus === '504') {
    return (
      `The certificate challenge for ${domain} returned HTTP ${challengeStatus}` +
      (redirectedToHttps ? ' after being redirected to HTTPS. ' : '. ') +
      `A gateway could not reach its challenge responder. Check forwarding for ` +
      `/.well-known/acme-challenge/ and that Codedock's temporary listener on port ${ACME_HTTP01_PORT} ` +
      `is reachable from the edge, then retry. This response alone does not establish a port-80 conflict.`
    );
  }

  if (challengeStatus === '525' || /\b525\b.*cloudflare/i.test(text)) {
    return (
      `Cloudflare answered the challenge itself (HTTP 525). To fix: ` +
      `1) Set SSL/TLS in Cloudflare to "Full" (not Full Strict) until certificate issuance completes, ` +
      `2) Or temporarily pause Cloudflare proxy (grey cloud), or 3) Upload an Origin CA certificate.`
    );
  }

  if (/timeout during connect|connection refused|failed to connect|timed? ?out/i.test(text)) {
    return (
      `Port 80 for ${domain} isn't reachable from the internet — a cloud security group or firewall ` +
      `is likely blocking inbound traffic on port 80.`
    );
  }

  if (/404|unauthorized|not found/i.test(text)) {
    return (
      `The ACME challenge token was not found (HTTP 404). Verify that ${domain} points to this server's ` +
      `public IP address and that reverse-proxy path routing is not intercepting /.well-known/acme-challenge/.`
    );
  }

  return `Let's Encrypt could not verify the challenge for ${domain}. Check your DNS records and server connectivity.`;
}
