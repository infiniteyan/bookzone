package common

type HttpCode int

const SessionId = "sid"

const MemberSessionName = "bookzone_session"
const UidSessionName   = "uid"

const (
	HttpCodeSuccess				HttpCode = iota
	HttpCodeErrorParameter
	HttpCodeErrorPassword
	HttpCodeErrorEmail
	HttpCodeErrorNickname
	HttpCodeErrorDatabase
	HttpCodeErrorLoginFirst
	HttpCodeErrorRegisterFail
	HttpCodeErrorLoginFail
)