package models

import (
	"bookzone/sysinit"
	"bookzone/util"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Comments struct {
	Id 				int
	Uid 			int 	`orm:"index"`
	BookId 			int		`orm:"index"`
	Content 		string
	TimeCreate		time.Time
}

type BookCommentsResult struct {
	Uid 			int			`json:"uid"`
	Score 			int			`json:"score"`
	Avatar 			string		`json:"avatar"`
	Nickname 		string		`json:"nickname"`
	Content 		string		`json:"content"`
	TimeCreate  	time.Time	`json:"time_create"`
}

func NewComments() *Comments {
	return &Comments{}
}

func (this *Comments) AddComments(uid, bookId int, content string) error {
	var err error
	second := 10
	sql := "select id from " + fmt.Sprintf("md_comments_%04d", bookId % 2) + " where uid = ? and time_create > ? order by id desc"
	ret, err := sysinit.DatabaseEngine.QueryString(sql, uid, time.Now().Add(-time.Duration(second) * time.Second))
	if err != nil {
		return err
	}
	if len(ret) > 0 {
		return errors.New("frequency limit")
	}

	sql = "insert into " + fmt.Sprintf("md_comments_%04d", bookId % 2) + "(uid, book_id, content, time_create) values (?, ?, ?, ?)"
	_, err = sysinit.DatabaseEngine.Exec(sql, uid, bookId, content, time.Now())
	if err != nil {
		return err
	}

	sql = "update md_books set cnt_comment = cnt_comment + 1 where book_id = ?"
	_, err = sysinit.DatabaseEngine.Exec(sql, bookId)
	if err != nil {
		return err
	}
	return nil
}

func (this *Comments) BookComments(page, size, bookId int) ([]*BookCommentsResult, error) {
	sql := "select book_id, uid, content, time_create from " + fmt.Sprintf("md_comments_%04d", bookId % 2) + " where book_id = ? limit %v offset %v"
	sql = fmt.Sprintf(sql, size, (page - 1) * size)
	retSlice, err := sysinit.DatabaseEngine.QueryString(sql, bookId)

	if err != nil {
		return nil, err
	}

	var comments []*BookCommentsResult
	for _, data := range retSlice {
		var cmtRet BookCommentsResult
		util.Map2struct(data, &cmtRet)
		comments = append(comments, &cmtRet)
	}

	uids := []string{}
	for _, v := range comments {
		uids = append(uids, strconv.Itoa(v.Uid))
	}
	uidstr := strings.Join(uids, ",")
	sql = "select member_id, avatar, nickname from md_members where member_id in(" + uidstr + ")"

	retSlice, err = sysinit.DatabaseEngine.QueryString(sql)
	if err != nil {
		return nil, err
	}

	var members []*Member
	for _, data := range retSlice {
		var cmtRet Member
		util.Map2struct(data, &cmtRet)
		members = append(members, &cmtRet)
	}

	memberMap := make(map[int]*Member)
	for _, member := range members {
		memberMap[member.MemberId] = member
	}
	for k, v := range comments {
		comments[k].Avatar = memberMap[v.Uid].Avatar
		comments[k].Nickname = memberMap[v.Uid].Nickname
	}

	sql = "select uid, score from md_score where book_id = ? and uid in(" + uidstr + ")"
	scores := []*Score{}
	retSlice, err = sysinit.DatabaseEngine.QueryString(sql, bookId)
	if err != nil {
		return nil, err
	}
	for _, data := range retSlice {
		var score Score
		util.Map2struct(data, &score)
		scores = append(scores, &score)
	}
	scoreMap := make(map[int]*Score)
	for _, score := range scores {
		scoreMap[score.Uid] = score
	}
	for k, v := range comments {
		comments[k].Score = scoreMap[v.Uid].Score
	}

	return comments, nil
}

type Score struct {
	Id 				int				`xorm:"pk autoincr"`
	BookId 			int				`json:"book_id"`
	Uid 			int				`json:"uid"`
	Score 			int				`json:"score"`
	TimeCreate 		time.Time		`json:"time_create"`
}

func NewScore() *Score {
	return &Score{}
}

func (this *Score) TableName() string {
	return "md_score"
}

type BookScoreResult struct {
	Avatar 			string		`json:"avatar"`
	Nickname 		string		`json:"nickname"`
	Score 			string		`json:"score"`
	TimeCreate 		time.Time	`json:"time_create"`
}

func (this *Score) BookScores(p, listRows, bookId int) ([]*BookScoreResult, error) {
	sql := "select s.score,s.time_create,m.avatar,m.nickname from md_score s left join md_members m on m.member_id=s.uid where s.book_id=? order by s.id desc limit %v offset %v"
	sql = fmt.Sprintf(sql, listRows, (p - 1) * listRows)

	retSlice, err := sysinit.DatabaseEngine.QueryString(sql)
	if err != nil {
		return nil, err
	}

	var scoreRets []*BookScoreResult
	for _, data := range retSlice {
		var scoreResult BookScoreResult
		util.Map2struct(data, &scoreResult)
		scoreRets = append(scoreRets, &scoreResult)
	}

	return scoreRets, nil
}

func (this *Score) BookScoreByUid(uid, bookId int) int {
	score := &Score{Uid:uid, BookId: bookId}
	has, err := sysinit.DatabaseEngine.Get(score)
	if err != nil || !has {
		return -1
	} else {
		return score.Score
	}
}

func (this *Score) AddScore(uid, bookId, score int) error {
	scoreObj := Score{Uid: uid, BookId: bookId}
	_, err := sysinit.DatabaseEngine.Get(&scoreObj)
	if err != nil {
		return err
	}

	if scoreObj.Id > 0 {
		return errors.New("score object already exists")
	}

	score = score * 10
	scoreObj.Score = score
	scoreObj.TimeCreate = time.Now()
	_, err = sysinit.DatabaseEngine.Insert(&scoreObj)
	if err != nil {
		return err
	}

	if scoreObj.Id > 0 {
		var book = Book{BookId: bookId}
		_, err = sysinit.DatabaseEngine.Get(&book)
		if err != nil {
			return err
		}

		if book.CntScore == 0 {
			book.CntScore = 1
			book.Score = 0
		} else {
			book.CntScore = book.CntScore + 1
		}
		book.Score = (book.Score * (book.CntScore - 1) + score) / book.CntScore
		_, err = sysinit.DatabaseEngine.Update(&book, &Book{BookId:bookId})
		if err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("database insert error")
	}
}