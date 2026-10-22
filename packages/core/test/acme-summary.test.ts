import { describe, expect, it } from 'vitest';
import { summarizeCertbotFailure } from '../src/acme-summary';

const OPENER = 'Certbot failed to authenticate some domains (authenticator: webroot).';

describe('summarizeCertbotFailure', () => {
  it.each([502, 504])(
    'HTTP %s -> inspect challenge forwarding without inferring a takeover',
    (status) => {
      const out = `${OPENER}
  Type:   unauthorized
  Detail: Invalid response from http://app.example.test/.well-known/acme-challenge/abc: ${status}${status === 504 ? ' Gateway Timeout' : ' Bad Gateway'}`;
      const summary = summarizeCertbotFailure(out, 'app.example.test');
      expect(summary).toContain(`HTTP ${status}`);
      expect(summary).toContain('49180');
      expect(summary).toMatch(/challenge.*responder|forwarding.*acme-challenge/i);
      expect(summary).not.toMatch(/Port 80.*occupied|another web server/i);
    }
  );

  it('does not diagnose a gateway failure from digits in the challenge URL', () => {
    const out = `${OPENER}
  Detail: Invalid response from http://502.example.test/.well-known/acme-challenge/token-504: 404`;
    const summary = summarizeCertbotFailure(out, '502.example.test');
    expect(summary).not.toMatch(/HTTP 502|HTTP 504|gateway could not/i);
    expect(summary).toContain('HTTP 404');
  });

  it('names an HTTPS gateway failure without claiming port 80 answered', () => {
    const out = `${OPENER}
  Detail: Invalid response from https://app.example.test/.well-known/acme-challenge/abc: 502`;
    const summary = summarizeCertbotFailure(out, 'app.example.test');
    expect(summary).toMatch(/redirected to HTTPS/i);
    expect(summary).not.toContain('Port 80 is occupied');
  });

  it('diagnoses actual port bind collision', () => {
    const out = `Problem binding to port 80: Could not bind to IPv4 or IPv6.`;
    const summary = summarizeCertbotFailure(out, 'app.example.test');
    expect(summary).toContain('Port 80 is occupied by another process');
  });

  it('diagnoses Cloudflare 525 error', () => {
    const out = `Invalid response from http://app.example.test/.well-known/acme-challenge/abc: 525 Cloudflare`;
    const summary = summarizeCertbotFailure(out, 'app.example.test');
    expect(summary).toContain('Cloudflare answered the challenge');
  });
});
