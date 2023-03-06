package audit

type AuditLog struct {
	UserId uint
	Method string
	Url    string
}
