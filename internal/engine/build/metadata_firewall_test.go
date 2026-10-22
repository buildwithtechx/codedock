package build

import (
	"strings"
	"testing"
)

func TestBuildMetadataFirewallScript(t *testing.T) {
	script := BuildMetadataFirewallScript()

	// Verify key iptables / ip6tables directives
	if !strings.Contains(script, "CODEDOCK_METADATA_V1") {
		t.Fatal("expected CODEDOCK_METADATA_V1 chain in script")
	}
	if !strings.Contains(script, "169.254.169.254/32 -j DROP") {
		t.Fatal("expected DROP rule for 169.254.169.254/32")
	}
	if !strings.Contains(script, "100.100.100.200/32 -j DROP") {
		t.Fatal("expected DROP rule for Alibaba 100.100.100.200/32")
	}
	if !strings.Contains(script, "fd00:ec2::254/128 -j DROP") {
		t.Fatal("expected DROP rule for IPv6 EC2 IMDS")
	}
	// Verify DNS exception on port 53
	if !strings.Contains(script, "-p udp --dport 53 -j RETURN") {
		t.Fatal("expected DNS UDP port 53 exception")
	}
	if !strings.Contains(script, "-p tcp --dport 53 -j RETURN") {
		t.Fatal("expected DNS TCP port 53 exception")
	}
}

func TestBuildKitGuardEntrypoint(t *testing.T) {
	entrypoint := BuildKitGuardEntrypoint()
	if !strings.HasPrefix(entrypoint, "#!/bin/sh\n") {
		t.Fatal("expected shebang in entrypoint")
	}
	if !strings.Contains(entrypoint, "CODEDOCK_BUILD_METADATA_GUARD_READY:1") {
		t.Fatal("expected guard ready notice in entrypoint")
	}
	if !strings.Contains(entrypoint, "exec /usr/bin/buildkitd") {
		t.Fatal("expected exec buildkitd in entrypoint")
	}
}

func TestIsProtectedBuildKitConfig(t *testing.T) {
	validLabels := map[string]string{
		BuildMetadataPolicyLabel: BuildMetadataPolicyVersion,
	}
	validEntrypoint := []string{"/usr/local/bin/codedock-buildkitd"}

	if !IsProtectedBuildKitConfig(validLabels, validEntrypoint) {
		t.Fatal("expected true for valid labels and entrypoint")
	}

	invalidLabels := map[string]string{
		BuildMetadataPolicyLabel: "old",
	}
	if IsProtectedBuildKitConfig(invalidLabels, validEntrypoint) {
		t.Fatal("expected false for invalid label version")
	}

	if IsProtectedBuildKitConfig(validLabels, []string{"/bin/bash"}) {
		t.Fatal("expected false for invalid entrypoint")
	}
}

func TestValidateOutboundTarget(t *testing.T) {
	// Metadata IPs must be rejected
	err := ValidateOutboundTarget("169.254.169.254")
	if err == nil {
		t.Fatal("expected error for 169.254.169.254")
	}

	err = ValidateOutboundTarget("169.254.169.254:80")
	if err == nil {
		t.Fatal("expected error for 169.254.169.254:80")
	}

	err = ValidateOutboundTarget("169.254.170.2")
	if err == nil {
		t.Fatal("expected error for ECS metadata IP")
	}

	err = ValidateOutboundTarget("100.100.100.200")
	if err == nil {
		t.Fatal("expected error for Alibaba metadata IP")
	}

	// Normal public IPs or domains must pass
	if err := ValidateOutboundTarget("8.8.8.8"); err != nil {
		t.Fatalf("unexpected error for 8.8.8.8: %v", err)
	}

	if err := ValidateOutboundTarget("example.com"); err != nil {
		t.Fatalf("unexpected error for example.com: %v", err)
	}
}
