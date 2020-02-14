package models

import (
	"bookzone/common"
	"bookzone/sysinit"
	"bookzone/util"
	"bookzone/util/log"
	"regexp"
	"time"
	"errors"
	"ziyoubiancheng/mbook/utils"
)

type Member struct {
	MemberId      int       `xorm:"pk autoincr" json:"member_id"`
	Account       string    `json:"account"`
	Nickname      string    `json:"nickname"`
	Password      string    `json:"-"`
	Description   string    `json:"description"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
	Avatar        string    `json:"avatar"`
	Role          int       `json:"role"`
	RoleName      string    `xorm:"-" json:"role_name"`
	Status        int       `json:"status"`
	CreateTime    time.Time `json:"create_time"`
	CreateAt      int       `json:"create_at"`
	LastLoginTime time.Time `json:"last_login_time"`
}

func (this *Member) TableName() string {
	return "md_members"
}

func NewMember() *Member {
	return &Member{}
}

func (this *Member) Add() error {
	var err error
	if this.Email == "" {
		return errors.New("please input email")
	}

	if ok, err := regexp.MatchString(common.RegexpEmail, this.Email); !ok || err != nil {
		return errors.New("invalid email format")
	}

	if l := len(this.Password); l < 6 || l > 20 {
		return errors.New("please limit password to 6 ~ 20 characters")
	}

	has, err := sysinit.DatabaseEngine.SQL("select * from md_members where nickname = ? or account = ? or email = ?",
		this.Nickname, this.Account, this.Email).Exist()
	if err != nil {
		return errors.New("fail to query database")
	}
	if has {
		return errors.New("already exist")
	}

	hash, err := utils.PasswordHash(this.Password)
	if err != nil {
		return err
	}
	this.Password = hash
	_, err = sysinit.DatabaseEngine.Insert(this)
	if err != nil {
		return err
	}

	this.RoleName = common.Role(this.Role)
	return nil
}

func (this *Member) Update(cols ...string) error {
	if this.Email == "" {
		return errors.New("email empty")
	}

	var err error
	condiBean := &Member{MemberId: this.MemberId}
	if len(cols) == 0 {
		_, err = sysinit.DatabaseEngine.Update(this, condiBean)
	} else {
		_, err = sysinit.DatabaseEngine.Cols(cols...).Update(this, condiBean)
	}

	return err
}

func (this *Member) Find(id int) (*Member, error) {
	var err error
	this.MemberId = id
	if _, err = sysinit.DatabaseEngine.Get(this); err != nil {
		return nil, err
	}

	this.RoleName = common.Role(this.Role)
	return this, nil
}

func (this *Member) Login(account string, password string) (*Member, error) {
	member := &Member{Account: account, Status: 0}
	has, err := sysinit.DatabaseEngine.Get(member)

	if err != nil || !has {
		return nil, errors.New("fail to get user info")
	}

	ok, err := util.PasswordVerify(member.Password, password)
	if ok && err == nil {
		member.RoleName = common.Role(member.Role)
		return member, nil
	}

	return member, errors.New("password not match")
}

func (this *Member) IsAdministrator() bool {
	if this == nil || this.MemberId <= 0 {
		return false
	}

	return this.Role == 0 || this.Role == 1
}

func (this *Member) GetUsernameByUid(id int) string {
	member := &Member{MemberId: id}
	_, err := sysinit.DatabaseEngine.Get(member)
	if err != nil {
		log.Errorf(err.Error())
		return ""
	}
	return member.Account
}

func (this *Member) GetNicknameByUid(id int) string {
	member := &Member{MemberId: id}
	_, err := sysinit.DatabaseEngine.Get(member)
	if err != nil {
		log.Errorf(err.Error())
		return ""
	}
	return member.Nickname
}

func (this *Member) GetByUsername(name string) (*Member, error) {
	member := &Member{Account: name}
	_, err := sysinit.DatabaseEngine.Get(member)
	if err != nil {
		log.Errorf(err.Error())
		return nil, err
	}
	return member, nil
}