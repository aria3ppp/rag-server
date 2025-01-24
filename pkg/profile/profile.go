package profile

const (
	profileDebug int8 = iota
	profileRelease
)

const IsDebug = profileValue == profileDebug
const IsRelease = profileValue == profileRelease
