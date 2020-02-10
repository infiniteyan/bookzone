package models

import (
	"bookzone/sysinit"
)

type RelationShip struct {
	RelationshipId 	int		`orm:"pk:auto;" json:"relationship_id"`
	MemberId 		int  	`json:"member_id"`
	BookId			int     `json:"book_id"`
	RoleId			int		`json:"role_id"`
}

func (this *RelationShip) TableName() string {
	return "md_relationship"
}

func NewRelationship() *RelationShip {
	return &RelationShip{}
}

func (this *RelationShip) Select(bookId, memberId int) (*RelationShip, error) {
	relationship := &RelationShip{BookId: bookId, MemberId: memberId}
	_, err := sysinit.DatabaseEngine.Get(relationship)
	if err != nil {
		return nil, err
	} else {
		return relationship, nil
	}
}

func (this *RelationShip) SelectRoleId(bookId, memberId int) (int, error) {
	relationship := &RelationShip{BookId: bookId, MemberId: memberId}
	_, err := sysinit.DatabaseEngine.Get(relationship)
	if err != nil {
		return -1, err
	} else {
		return relationship.RoleId, nil
	}
}

func (this *RelationShip) Insert() error {
	_, err := sysinit.DatabaseEngine.Insert(this)
	return err
}

func (this *RelationShip) Update() error {
	_, err := sysinit.DatabaseEngine.Update(this, &RelationShip{RelationshipId:this.RelationshipId})
	return err
}