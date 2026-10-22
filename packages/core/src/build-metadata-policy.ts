/** Remote-build policy only. Keep provider DNS reachable; deny metadata and
 * credential agents by destination regardless of the URL or DNS spelling. */
export const BUILD_METADATA_POLICY_VERSION = "1";

export function isProtectedBuildKitConfig(config: unknown): boolean {
  const value = config as { Labels?: Record<string, unknown>; Entrypoint?: unknown } | null;
  return (
    value?.Labels?.["codedock.build-metadata-policy"] === BUILD_METADATA_POLICY_VERSION &&
    JSON.stringify(value.Entrypoint) === JSON.stringify(["/usr/local/bin/codedock-buildkitd"])
  );
}

export const BUILD_METADATA_IPV4 = [
  "169.254.169.254/32", // EC2, Azure, GCE, OCI, DigitalOcean, Hetzner
  "169.254.170.2/32",   // ECS task credentials
  "169.254.170.23/32",  // EKS Pod Identity credentials
  "100.100.100.200/32", // Alibaba ECS
  "169.254.42.42/32",   // Scaleway
] as const;

export const BUILD_METADATA_IPV6 = [
  "fd00:ec2::254/128", // EC2 IMDS
  "fd00:ec2::23/128",  // EKS Pod Identity credentials
  "fd20:ce::254/128",  // GCE IPv6-only metadata and DNS
] as const;

export function buildMetadataFirewallScript(): string {
  const chain = "CODEDOCK_METADATA_V1";
  return (
    [
      "set -eu",
      ...(
        [
          ["iptables", BUILD_METADATA_IPV4],
          ["ip6tables", BUILD_METADATA_IPV6],
        ] as const
      ).flatMap(([tool, destinations]) => [
        `command -v ${tool} >/dev/null`,
        `${tool} -w 5 -t raw -N ${chain} 2>/dev/null || ${tool} -w 5 -t raw -S ${chain} >/dev/null`,
        `${tool} -w 5 -t raw -F ${chain}`,
        `${tool} -w 5 -t raw -A ${chain} -d ${tool === "iptables" ? "169.254.169.254/32" : "fd20:ce::254/128"} -p udp --dport 53 -j RETURN`,
        `${tool} -w 5 -t raw -A ${chain} -d ${tool === "iptables" ? "169.254.169.254/32" : "fd20:ce::254/128"} -p tcp --dport 53 -j RETURN`,
        ...destinations.map(
          (destination) => `${tool} -w 5 -t raw -A ${chain} -d ${destination} -j DROP`,
        ),
        ...["PREROUTING", "OUTPUT"].flatMap((hook) => [
          `${tool} -w 5 -t raw -I ${hook} 1 -j ${chain}`,
          `test "$(${tool} -w 5 -t raw -S ${hook} | sed -n '/^-A /{p;q;}')" = '-A ${hook} -j ${chain}'`,
        ]),
        ...destinations.map(
          (destination) => `${tool} -w 5 -t raw -C ${chain} -d ${destination} -j DROP`,
        ),
      ]),
    ].join("\n") + "\n"
  );
}

export function buildKitGuardEntrypoint(): string {
  return (
    "#!/bin/sh\n" +
    buildMetadataFirewallScript() +
    `printf '%s\\n' 'CODEDOCK_BUILD_METADATA_GUARD_READY:${BUILD_METADATA_POLICY_VERSION}' >&2\n` +
    'exec /usr/bin/buildkitd "$@"\n'
  );
}
