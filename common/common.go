package common

const SessionName = "__bookzone_session__"

const RegexpEmail = `^(\w)+(\.\w+)*@(\w)+((\.\w+)+)$`

const PageSize = 20
const RollPage = 4

const WorkingDirectory = "./"

const (
	MemberSuperRole   = 0
	MemberAdminRole   = 1
	MemberGeneralRole = 2
)

func Role(role int) string {
	switch role {
	case MemberGeneralRole:
		return "普通用户"
	case MemberAdminRole:
		return "管理员"
	case MemberSuperRole:
		return "超级管理员"
	default:
		return ""
	}
}

const (
	BookFounder = 0
	BookAdmin   = 1
	BookEditor  = 2
	BookGeneral = 3
)

func BookRole(role int) string {
	switch role {
	case BookFounder:
		return "创建人"
	case BookAdmin:
		return "管理员"
	case BookEditor:
		return "编辑"
	case BookGeneral:
		return "普通用户"
	default:
		return ""
	}
}