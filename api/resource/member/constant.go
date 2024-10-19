package member

const (
	Success = "G0000"

	ErrSignUpErrBindRequest      = "M4001"
	ErrSignUpErrHashPassword     = "M4002"
	ErrSignUpErrRepoCreateMember = "M4003"

	ErrSignInErrBindRequest          = "M4101"
	ErrSignInErrGetMemberInformation = "M4102"
	ErrSignInErrCompareHashPassword  = "M4103"

	ErrGetProfileErrGetCachedMemberProfile  = "M4201"
	ErrGetProfileErrUnmarshalMemberProfile  = "M4202"
	ErrGetProfileErrGetMemberProfile        = "M4203"
	ErrGetProfileErrMarshalMemberProfile    = "M4204"
	ErrGetProfileErrSetMemberProfileToRedis = "M4205"
)
