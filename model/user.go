package model

import (
	"log"
	"strings"
	"time"

	"github.com/Bit0r/online-store/model/perm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           uint64
	Name         string
	CreationTime time.Time
	Privileges   []string
}

type Users = []User

func VerifyUser(name, passwd string) (id uint64) {
	var hash string
	query := `select id, passwd
	from user
	where name = ?`

	if db.QueryRow(query, name).Scan(&id, &hash) != nil {
		// 验证用户名是否正确
		return 0
	}

	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(passwd)) != nil {
		// 验证密码是否正确
		return 0
	}

	return
}

func AddUser(name, passwd string) (id uint64) {
	query := `insert ignore into user(name, passwd) values(?, ?)`
	hash, _ := bcrypt.GenerateFromPassword([]byte(passwd), 0)
	result, err := db.Exec(query, name, hash)
	if count, _ := result.RowsAffected(); err != nil || count == 0 {
		return 0
	} else {
		id, _ := result.LastInsertId()
		return uint64(id)
	}
}

func GetPrivileges(userID uint64) ([]string, error) {
	query := `select privilege
	from user_privilege
	where user_id = ?
	order by privilege`
	rs, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}

	var privilege string
	var privileges []string
	for rs.Next() {
		err := rs.Scan(&privilege)
		if err != nil {
			log.Println(err)
		} else {
			privileges = append(privileges, privilege)
		}
	}

	return privileges, nil
}

func GetPrivilegeSet(userID uint64) (perm.PrivilegeSet, error) {
	privileges, err := GetPrivileges(userID)
	if err != nil {
		return 0, err
	}
	return perm.NewByStr(privileges...), nil
}

func HasPrivilege(userID uint64, privilege string) bool {
	query := `select true
	from user_privilege
	where user_id = ? and (privilege = ? or privilege = 'all')`

	flag := false
	db.QueryRow(query, userID, privilege).Scan(&flag)
	return flag
}

func HasPrivileges(userID uint64, privileges []string) bool {
	holders := strings.Repeat(",?", len(privileges))[1:]
	query := `select count(*) = ?
	from user_privilege
	where user_id = ? and privilege in (` + holders + ")"

	args := make([]interface{}, len(privileges)+2)
	args[0], args[1] = len(privileges), userID
	for i, v := range privileges {
		args[i+2] = v
	}

	flag := false
	db.QueryRow(query, args...).Scan(&flag)
	return flag
}

func CountUsers(isAdmin bool) (count uint64) {
	query := "select count(distinct id)"
	if isAdmin {
		query += ` from user u, user_privilege up
		where u.id = up.user_id`
	} else {
		query += ` from user u
		left join user_privilege up
		on u.id = up.user_id
		where up.user_id is null`
	}
	db.QueryRow(query).Scan(&count)
	return
}

func GetUsers(isAdmin bool, limit Limit) (users Users, err error) {
	query := `select distinct id, name, creation_time`

	if isAdmin {
		query += ` from user u, user_privilege up
		where u.id = up.user_id`
	} else {
		query += ` from user u
		left join user_privilege up
		on u.id = up.user_id
		where up.user_id is null`
	}

	query += " limit ?, ?"

	rows, err := db.Query(query, limit.Offset, limit.Count)
	if err != nil {
		return
	}
	defer rows.Close()

	query = `select privilege
	from user_privilege
	where user_id = ?`
	var user User
	for rows.Next() {
		err = rows.Scan(&user.ID, &user.Name, &user.CreationTime)
		if err != nil {
			log.Println(err)
			continue
		}

		if !isAdmin {
			users = append(users, user)
			continue
		}

		// 如果是管理员，则获取所有权限
		user.Privileges, err = GetPrivileges(user.ID)
		if err != nil {
			log.Println(err)
			continue
		}
		users = append(users, user)
	}
	return users, nil
}

func GetUser(id uint64) (user User, err error) {
	query := `select id, name, creation_time
	from user
	where id = ?`

	err = db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.CreationTime)
	if err != nil {
		return
	}

	user.Privileges, err = GetPrivileges(id)
	if err != nil {
		return
	}
	return
}

func SetPrivileges(userID uint64, privileges []string) (err error) {
	txn, err := db.Begin()
	if err != nil {
		return
	}
	query := `delete from user_privilege
	where user_id = ?`
	txn.Exec(query, userID)

	if len(privileges) == 0 {
		// 如果没有权限，则不需要插入
		txn.Commit()
		return
	}

	holders := strings.Repeat(",(?,?)", len(privileges))[1:]
	query = "insert into user_privilege(user_id, privilege) values" + holders
	args := make([]any, len(privileges)*2)
	for i, priv := range privileges {
		args[i*2] = userID
		args[i*2+1] = priv
	}
	_, err = txn.Exec(query, args...)
	if err != nil {
		// 如果插入失败，则回滚
		txn.Rollback()
		return
	}
	txn.Commit()
	return
}
