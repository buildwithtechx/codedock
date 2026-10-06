package kubernetes

import (
	"codedock.run/codedock/internal/models"
	"encoding/base64"
	"testing"
)

func TestValidateNodeRequiresPrivateNonOverlappingNetworkAndRealFingerprint(t *testing.T) {
	valid := models.ClusterNode{ServerID: "96572286-6660-4dd6-b22d-40952f879001", PrivateIP: "192.168.3.4", Interface: "ens3", Fingerprint: "SHA256:" + base64.RawStdEncoding.EncodeToString(make([]byte, 32))}
	if err := ValidateNode(valid); err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct{ ip, network, key string }{
		{"8.8.8.8", "ens3", valid.Fingerprint},
		{"10.42.3.4", "ens3", valid.Fingerprint},
		{"10.43.3.4", "ens3", valid.Fingerprint},
		{valid.PrivateIP, "ens3;touch /tmp/x", valid.Fingerprint},
		{valid.PrivateIP, "ens3", "SHA256:bogus"},
	} {
		node := valid
		node.PrivateIP = test.ip
		node.Interface = test.network
		node.Fingerprint = test.key
		if err := ValidateNode(node); err == nil {
			t.Fatal("invalid node accepted", test)
		}
	}
}
