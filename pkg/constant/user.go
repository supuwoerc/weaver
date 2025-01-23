package constant

const (
	EmailRegexPattern  = `^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`
	PasswdRegexPattern = `^[a-fA-F0-9]{32}$`
	PhoneRegexPattern  = `^(13[0-9]|14[01456879]|15[0-35-9]|16[2567]|17[0-8]|18[0-9]|19[0-35-9])\d{8}$`
	UserContextKey     = "user"
	ClaimsContextKey   = "claims"
)
