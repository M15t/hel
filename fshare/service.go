package fshare

// Configuration struct
type Configuration struct {
	UserEmail string
	UserPwd   string
	UserAgent string
	AppKey    string
	Credential
}

// Credential struct
type Credential struct {
	Token     *string
	SessionID *string
}

// New initializes Fshare service with default config
func New() *Fshare {
	return &Fshare{
		Configuration{
			UserEmail: "trees1411@yahoo.com",
			UserPwd:   "pVPdzJEJaRhjG@8",
			UserAgent: "checkAPI-UT42RX",
			AppKey:    "dMnqMMZMUnN5YpvKENaEhdQQ5jxDqddt",
		},
	}
}

// Fshare represents the fshare service
type Fshare struct {
	cfg Configuration
}
