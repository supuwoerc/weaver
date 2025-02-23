package constant

const (
	EmailRegexPattern  = `^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`
	PasswdRegexPattern = `^[a-fA-F0-9]{32}$`
	PhoneRegexPattern  = `^(13[0-9]|14[01456879]|15[0-35-9]|16[2567]|17[0-8]|18[0-9]|19[0-35-9])\d{8}$`
	ClaimsContextKey   = "claims"
	UserTokenPairKey   = "user_cache:token"
)

//go:generate stringer -type=UserStatus -linecomment -output user_status_string.go
type UserStatus int

const (
	Inactive UserStatus = iota + 1 // inactive
	Normal                         // normal
	Disabled                       // disabled
)

//go:generate stringer -type=UserGender -linecomment -output user_gender_string.go
type UserGender int

const (
	Male    UserGender = iota + 1 // male
	Females                       // females
)
