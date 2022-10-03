package auth

type Authenticator interface {
	VerifyBearer(token string) (string, error)
}
