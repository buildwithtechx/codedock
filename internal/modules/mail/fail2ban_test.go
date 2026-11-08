package mail

import (
	"testing"
)

func TestParseFail2banStatus(t *testing.T) {
	t.Run("valid multi-jail status", func(t *testing.T) {
		output := "Status\n|- Number of jail:\t2\n`- Jail list:\tpostfix, dovecot\n"
		jails, err := ParseFail2banStatus(output)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(jails) != 2 || jails[0] != "postfix" || jails[1] != "dovecot" {
			t.Fatalf("expected [postfix, dovecot], got %v", jails)
		}
	})

	t.Run("zero jails", func(t *testing.T) {
		output := "Status\n|- Number of jail:\t0\n`- Jail list:\t\n"
		jails, err := ParseFail2banStatus(output)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(jails) != 0 {
			t.Fatalf("expected 0 jails, got %v", jails)
		}
	})

	t.Run("count mismatch rejects", func(t *testing.T) {
		output := "Status\n|- Number of jail:\t3\n`- Jail list:\tpostfix, dovecot\n"
		_, err := ParseFail2banStatus(output)
		if err == nil {
			t.Fatal("expected error on count mismatch")
		}
	})

	t.Run("malformed output rejects", func(t *testing.T) {
		_, err := ParseFail2banStatus("completely invalid gibberish")
		if err == nil {
			t.Fatal("expected error on invalid output")
		}
	})
}

func TestParseFail2banJailStatus(t *testing.T) {
	t.Run("valid jail with IPv4 and IPv6 banned IPs", func(t *testing.T) {
		output := "Status for the jail: postfix\n" +
			"|- Filter\n" +
			"|  |- Currently failed:\t3\n" +
			"|  |- Total failed:\t42\n" +
			"|  `- File list:\t/var/log/mail.log\n" +
			"`- Actions\n" +
			"   |- Currently banned:\t2\n" +
			"   |- Total banned:\t15\n" +
			"   `- Banned IP list:\t198.51.100.1 2001:db8::1\n"

		jail, err := ParseFail2banJailStatus("postfix", output)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if jail.Jail != "postfix" {
			t.Errorf("expected postfix, got %q", jail.Jail)
		}
		if jail.CurrentlyFailed != 3 || jail.TotalFailed != 42 {
			t.Errorf("unexpected failed counts: %d, %d", jail.CurrentlyFailed, jail.TotalFailed)
		}
		if jail.CurrentlyBanned != 2 || jail.TotalBanned != 15 {
			t.Errorf("unexpected banned counts: %d, %d", jail.CurrentlyBanned, jail.TotalBanned)
		}
		if len(jail.BannedIPs) != 2 || jail.BannedIPs[0] != "198.51.100.1" || jail.BannedIPs[1] != "2001:db8::1" {
			t.Errorf("unexpected banned IPs: %v", jail.BannedIPs)
		}
	})

	t.Run("valid jail with no banned IPs", func(t *testing.T) {
		output := "Status for the jail: dovecot\n" +
			"|- Filter\n" +
			"|  |- Currently failed:\t0\n" +
			"|  |- Total failed:\t0\n" +
			"|  `- File list:\t/var/log/mail.log\n" +
			"`- Actions\n" +
			"   |- Currently banned:\t0\n" +
			"   |- Total banned:\t0\n" +
			"   `- Banned IP list:\t\n"

		jail, err := ParseFail2banJailStatus("dovecot", output)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(jail.BannedIPs) != 0 {
			t.Errorf("expected empty banned IPs, got %v", jail.BannedIPs)
		}
	})

	t.Run("mismatched jail name in output rejects", func(t *testing.T) {
		output := "Status for the jail: other_jail\n" +
			"|- Filter\n" +
			"|  |- Currently failed:\t0\n" +
			"|  |- Total failed:\t0\n" +
			"|  `- File list:\t/var/log/mail.log\n" +
			"`- Actions\n" +
			"   |- Currently banned:\t0\n" +
			"   |- Total banned:\t0\n" +
			"   `- Banned IP list:\t\n"

		_, err := ParseFail2banJailStatus("dovecot", output)
		if err == nil {
			t.Fatal("expected error on mismatched jail name")
		}
	})

	t.Run("invalid IP rejects", func(t *testing.T) {
		output := "Status for the jail: postfix\n" +
			"|- Filter\n" +
			"|  |- Currently failed:\t1\n" +
			"|  |- Total failed:\t1\n" +
			"|  `- File list:\t/var/log/mail.log\n" +
			"`- Actions\n" +
			"   |- Currently banned:\t1\n" +
			"   |- Total banned:\t1\n" +
			"   `- Banned IP list:\t999.999.999.999\n"

		_, err := ParseFail2banJailStatus("postfix", output)
		if err == nil {
			t.Fatal("expected error on invalid IP format")
		}
	})
}

func TestValidateIPAndJail(t *testing.T) {
	if err := ValidateIP("1.2.3.4"); err != nil {
		t.Errorf("expected valid IPv4: %v", err)
	}
	if err := ValidateIP("2001:0db8:85a3:0000:0000:8a2e:0370:7334"); err != nil {
		t.Errorf("expected valid IPv6: %v", err)
	}
	if err := ValidateIP("not-an-ip"); err == nil {
		t.Errorf("expected error for invalid IP")
	}

	if err := ValidateJailName("postfix"); err != nil {
		t.Errorf("expected valid jail: %v", err)
	}
	if err := ValidateJailName("dovecot-ssl"); err != nil {
		t.Errorf("expected valid jail: %v", err)
	}
	if err := ValidateJailName("bad;command"); err == nil {
		t.Errorf("expected error for jail with semicolon")
	}
}
