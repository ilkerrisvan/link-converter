package link

import (
	"database/sql"
	"errors"
	"time"
	"trendyolcase/pkg/model"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

/*
If the webURL in the db, return the corresponding deeplink.
*/

func (l *Repository) GetDeepLinkIfWebURLExist(webURL string) (string, error) {
	link := new(model.Link)
	q, err := l.db.Prepare("select short_url from links where long_url = $1")
	if err != nil{
		panic("There is an error.")
	}
	err = q.QueryRow(webURL).Scan(&link.Deeplink)
	if err != nil {
		return "", errors.New("Database connection or query has problem.")
	} else {
		return link.Deeplink, nil
	}
}

/*
If the deeplink in the db, return the corresponding webURL.
*/

func (l *Repository) GetWebURLIfDeepLinkExist(deepLink string) (string, error) {
	link := new(model.Link)
	q, _ := l.db.Prepare("select long_url from links where short_url = $1")
	err := q.QueryRow(deepLink).Scan(&link.WebUrl)
	if err != nil {
		return "", nil
	} else {
		return link.WebUrl, nil
	}
}

/*
Adds weblink - webURL pairs to db.
*/

func (l *Repository) Insert(webURL string, deepLink string) bool {
	q, _ := l.db.Prepare("insert into links(long_url,short_url) values($1,$2)")
	_, err := q.Exec(webURL, deepLink)
	if err != nil {
		return false
	} else {
		return true
	}
}

/*
Adds logs about requests to db.
*/

func (l *Repository) InsertLog(logInformation string) bool {
	timeNow := time.Now()
	q, _ := l.db.Prepare("insert into logs(created_at,info) values($1,$2)")
	_, err := q.Exec(timeNow, logInformation)
	if err != nil {
		return false
	} else {
		return true
	}
}
