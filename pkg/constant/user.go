package constant

const (
	EmailRegexPattern  = `^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`
	PasswdRegexPattern = `^^(?=.*[0-9])(?=.*[a-zA-Z])[0-9A-Za-z~!@#$%^&*._?]{8,15}$`
	PhoneRegexPattern  = `^(13[0-9]|14[01456879]|15[0-35-9]|16[2567]|17[0-8]|18[0-9]|19[0-35-9])\d{8}$`
)
