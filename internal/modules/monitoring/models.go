package monitoring

type Vault interface {
	Encrypt(string) (string, error)
	Decrypt(string) (string, error)
}
