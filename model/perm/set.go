package perm

type PrivilegeSet uint64

const (
	Order = 1 << iota
	Book
	User
	All = 1<<iota - 1
)

var str2bit = map[string]PrivilegeSet{
	"order": Order,
	"book":  Book,
	"user":  User,
	"all":   All,
}

func NewByStr(privileges ...string) PrivilegeSet {
	var bits PrivilegeSet
	for _, privilege := range privileges {
		bits |= str2bit[privilege]
	}
	return bits
}

func (privileges PrivilegeSet) Has(privilege PrivilegeSet) bool {
	return privileges&privilege == privilege
}

func (privileges PrivilegeSet) HasPrivilege(privilege string) bool {
	return privileges.Has(str2bit[privilege])
}
