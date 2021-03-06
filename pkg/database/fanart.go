package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"math/rand"
	"regexp"
	"strings"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	log "github.com/sirupsen/logrus"
)

//GetFanart Get Member fanart URL from TBiliBili and Twitter
func GetFanart(GroupID, MemberID int64) (*DataFanart, error) {
	var (
		Data     DataFanart
		PhotoTmp sql.NullString
		Video    sql.NullString
	)

	Twitter := func() error {
		rows, err := DB.Query(`SELECT Twitter.* FROM Vtuber.Twitter Inner Join Vtuber.VtuberMember on VtuberMember.id = Twitter.VtuberMember_id Inner Join Vtuber.VtuberGroup on VtuberGroup.id = VtuberMember.VtuberGroup_id where (VtuberGroup.id=? OR VtuberMember.id=?)  ORDER by RAND() LIMIT 1`, GroupID, MemberID)
		if err != nil {
			return err
		}

		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&Data.ID, &Data.PermanentURL, &Data.Author, &Data.Likes, &PhotoTmp, &Video, &Data.Text, &Data.TweetID, &Data.Member.ID)
			if err != nil {
				return err
			}
		}

		if Data.ID == 0 {
			return errors.New("Vtuber don't have any fanart in Twitter")
		}
		Data.State = config.TwitterArt
		return nil
	}
	Tbilibili := func() error {
		rows, err := DB.Query(`SELECT TBiliBili.* FROM Vtuber.TBiliBili Inner Join Vtuber.VtuberMember on VtuberMember.id = TBiliBili.VtuberMember_id Inner Join Vtuber.VtuberGroup on VtuberGroup.id = VtuberMember.VtuberGroup_id where (VtuberGroup.id=? OR VtuberMember.id=?)  ORDER by RAND() LIMIT 1`, GroupID, MemberID)
		if err != nil {
			return err
		}

		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&Data.ID, &Data.PermanentURL, &Data.Author, &Data.Likes, &PhotoTmp, &Video, &Data.Text, &Data.Dynamic_id, &Data.Member.ID)
			if err != nil {
				return err
			}
		}
		if Data.ID == 0 {
			return errors.New("Vtuber don't have any fanart in BiliBili")
		}

		Data.State = config.BiliBiliArt
		return nil
	}

	Pixiv := func() error {
		rows, err := DB.Query(`SELECT Pixiv.* FROM Vtuber.Pixiv Inner Join Vtuber.VtuberMember on VtuberMember.id = Pixiv.VtuberMember_id Inner Join Vtuber.VtuberGroup on VtuberGroup.id = VtuberMember.VtuberGroup_id where (VtuberGroup.id=? OR VtuberMember.id=?)  ORDER by RAND() LIMIT 1`, GroupID, MemberID)
		if err != nil {
			return err
		}

		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&Data.ID, &Data.PermanentURL, &Data.Author, &PhotoTmp, &Data.Text, &Data.PixivID, &Data.Member.ID)
			if err != nil {
				return err
			}
		}
		if Data.ID == 0 {
			return errors.New("Vtuber don't have any fanart in Pixiv")
		}

		Data.State = config.PixivArt
		return nil
	}

	gachaint := rand.Intn(3-1) + 1
	if gachaint == 1 {
		err := Pixiv()
		if err != nil {
			return nil, err
		}
	} else if gachaint == 2 {
		err := Tbilibili()
		if err != nil {
			return nil, err
		}
	} else {
		err := Twitter()
		if err != nil {
			return nil, err
		}
	}

	Data.Videos = Video.String
	Data.Photos = strings.Fields(PhotoTmp.String)
	return &Data, nil

}

//GetFanart Get Member fanart URL from TBiliBili and Twitter
func GetLewd(GroupID, MemberID int64) (*DataFanart, error) {
	var (
		Data     DataFanart
		PhotoTmp sql.NullString
		Video    sql.NullString
	)

	rows, err := DB.Query(`SELECT Lewd.* FROM Vtuber.Lewd Inner Join Vtuber.VtuberMember on VtuberMember.id = Lewd.VtuberMember_id Inner Join Vtuber.VtuberGroup on VtuberGroup.id = VtuberMember.VtuberGroup_id where (VtuberGroup.id=? OR VtuberMember.id=?)  ORDER by RAND() LIMIT 1`, GroupID, MemberID)
	if err != nil {
		return nil, err
	} else if err == sql.ErrNoRows {
		return nil, errors.New("Vtuber don't have any fanart in Twitter")
	}

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&Data.ID, &Data.PermanentURL, &Data.Author, &PhotoTmp, &Video, &Data.Text, &Data.TweetID, &Data.PixivID, &Data.Member.ID)
		if err != nil {
			return nil, err
		}
	}
	if Data.PixivID != "" {
		Data.State = config.PixivArt
	} else {
		Data.State = config.TwitterArt
	}
	Data.Photos = strings.Fields(PhotoTmp.String)
	Data.Videos = Video.String
	return &Data, nil

}

//DeleteFanart Delete fanart when get 404 error status
func (Data DataFanart) DeleteFanart(e string) error {
	if notfound, _ := regexp.MatchString("404", e); notfound {
		log.Info("Delete fanart metadata ", Data.PermanentURL)
		if Data.State == "Twitter" {
			stmt, err := DB.Prepare(`DELETE From Twitter WHERE id=?`)
			if err != nil {
				return err
			}
			defer stmt.Close()

			stmt.Exec(Data.ID)
			return nil
		} else {
			stmt, err := DB.Prepare(`DELETE From TBiliBili WHERE id=?`)
			if err != nil {
				return err
			}
			defer stmt.Close()

			stmt.Exec(Data.ID)
			return nil
		}
	} else {
		return nil
	}
}

func AddLewd(Data DataFanart) error {
	stmt, err := DB.Prepare(`INSERT INTO Lewd (PermanentURL,Author,Photos,Videos,Text,TweetID,PixivID,VtuberMember_id) values(?,?,?,?,?,?,?,?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(Data.PermanentURL, Data.Author, strings.Join(Data.Photos, "\n"), Data.Videos, Data.Text, Data.TweetID, Data.PixivID, Data.Member.ID)
	if err != nil {
		return err
	}

	_, err = res.LastInsertId()
	if err != nil {
		return err
	}
	return nil
}

//CheckMemberFanart Check if `that` was a new fanart
func (FanArt DataFanart) CheckTweetFanArt() (bool, error) {
	if FanArt.Lewd {
		var (
			id int
		)
		err := DB.QueryRow(`SELECT id FROM Lewd WHERE TweetID=?`, FanArt.TweetID).Scan(&id)
		if err == sql.ErrNoRows {
			log.WithFields(log.Fields{
				"Name":    FanArt.Member.EnName,
				"Hashtag": FanArt.Member.TwitterHashtags,
				"Lewd":    FanArt.Lewd,
				"URL":     FanArt.PermanentURL,
			}).Info("New Fanart")
			err = AddLewd(FanArt)
			if err != nil {
				return false, err
			}
			return true, nil
		}
	} else {
		var (
			id int
		)
		err := DB.QueryRow(`SELECT id FROM Twitter WHERE TweetID=?`, FanArt.TweetID).Scan(&id)
		if err == sql.ErrNoRows {
			log.WithFields(log.Fields{
				"Name":    FanArt.Member.EnName,
				"Hashtag": FanArt.Member.TwitterHashtags,
				"Lewd":    FanArt.Lewd,
				"URL":     FanArt.PermanentURL,
			}).Info("New Fanart")

			stmt, err := DB.Prepare(`INSERT INTO Twitter (PermanentURL,Author,Likes,Photos,Videos,Text,TweetID,VtuberMember_id) values(?,?,?,?,?,?,?,?)`)
			if err != nil {
				return false, err
			}
			defer stmt.Close()

			res, err := stmt.Exec(FanArt.PermanentURL, FanArt.Author, FanArt.Likes, strings.Join(FanArt.Photos, "\n"), FanArt.Videos, FanArt.Text, FanArt.TweetID, FanArt.Member.ID)
			if err != nil {
				return false, err
			}

			_, err = res.LastInsertId()
			if err != nil {
				return false, err
			}
			return true, nil
		} else if err != nil {
			return false, err
		} else {
			if !config.GoSimpConf.LowResources {
				//update like
				log.WithFields(log.Fields{
					"Name":    FanArt.Member.EnName,
					"Hashtag": FanArt.Member.TwitterHashtags,
					"Likes":   FanArt.Likes,
				}).Info("Update like")
				_, err := DB.Exec(`Update Twitter set Likes=? Where id=? `, FanArt.Likes, id)
				if err != nil {
					return false, err
				}
			}
		}
	}
	return false, nil
}

func (FanArt DataFanart) CheckTBiliBiliFanArt() (bool, error) {
	var tmp int64
	row := DB.QueryRow("SELECT id FROM Vtuber.TBiliBili where Dynamic_id=?", FanArt.Dynamic_id)
	err := row.Scan(&tmp)
	if err == sql.ErrNoRows {
		log.WithFields(log.Fields{
			"Vtuber": FanArt.Member.EnName,
			"Img":    FanArt.Photos,
		}).Info("New Fanart")
		stmt, err := DB.Prepare(`INSERT INTO TBiliBili (PermanentURL,Author,Likes,Photos,Videos,Text,Dynamic_id,VtuberMember_id) values(?,?,?,?,?,?,?,?)`)
		if err != nil {
			return false, err
		}
		defer stmt.Close()

		res, err := stmt.Exec(FanArt.PermanentURL, FanArt.Author, FanArt.Likes, strings.Join(FanArt.Photos, "\n"), FanArt.Videos, FanArt.Text, FanArt.Dynamic_id, FanArt.Member.ID)
		if err != nil {
			return false, err
		}

		_, err = res.LastInsertId()
		if err != nil {
			return false, err
		}
		return true, nil
	} else if err != nil {
		log.Error(err)
	}
	return false, nil

}

func (FanArt DataFanart) CheckPixivFanArt() (bool, error) {
	if FanArt.Lewd {
		var tmp int64
		row := DB.QueryRow("SELECT id FROM Vtuber.Lewd where PixivID=?", FanArt.PixivID)
		err := row.Scan(&tmp)
		if err == sql.ErrNoRows {
			log.WithFields(log.Fields{
				"Vtuber": FanArt.Member.EnName,
				"Img":    FanArt.Photos,
				"URL":    FanArt.PermanentURL,
			}).Info("New Lewd Fanart")
			err = AddLewd(FanArt)
			if err != nil {
				return false, err
			}
			return true, nil
		} else if err != nil {
			return false, err
		} else {
			return false, errors.New(FanArt.Member.Name + " Still same")
		}
	} else {
		var tmp int64
		row := DB.QueryRow("SELECT id FROM Vtuber.Pixiv where PixivID=?", FanArt.PixivID)
		err := row.Scan(&tmp)
		if err == sql.ErrNoRows {
			log.WithFields(log.Fields{
				"Vtuber": FanArt.Member.EnName,
				"Img":    FanArt.Photos,
				"URL":    FanArt.PermanentURL,
			}).Info("New Fanart")
			stmt, err := DB.Prepare(`INSERT INTO Pixiv (PermanentURL,Author,Photos,Text,PixivID,VtuberMember_id) values(?,?,?,?,?,?)`)
			if err != nil {
				return false, err
			}
			defer stmt.Close()

			res, err := stmt.Exec(FanArt.PermanentURL, FanArt.Author, strings.Join(FanArt.Photos, "\n"), FanArt.Text, FanArt.PixivID, FanArt.Member.ID)
			if err != nil {
				return false, err
			}

			_, err = res.LastInsertId()
			if err != nil {
				return false, err
			}
			return true, nil

		} else if err != nil {
			return false, err
		} else {
			return false, errors.New(FanArt.Member.Name + " Still same")
		}
	}
}

//DataFanart fanart struct
type DataFanart struct {
	ID           int64
	Member       Member
	Group        Group
	PermanentURL string
	Author       string
	AuthorAvatar string
	Photos       []string
	Videos       string
	Text         string
	Likes        int
	Dynamic_id   string
	TweetID      string
	PixivID      string
	Lewd         bool
	State        string
}

func (Data *DataFanart) MarshallBin() []byte {
	bit, err := json.Marshal(Data)
	if err != nil {
		log.Error(err)
	}
	return bit
}

func (Data *DataFanart) AddMember(new Member) *DataFanart {
	Data.Member = new
	return Data
}

func (Data *DataFanart) AddGroup(new Group) *DataFanart {
	Data.Group = new
	return Data
}

func (Data *DataFanart) AddPermanentURL(new string) *DataFanart {
	Data.PermanentURL = new
	return Data
}

func (Data *DataFanart) AddAuthor(new string) *DataFanart {
	Data.Author = new
	return Data
}

func (Data *DataFanart) AddAuthorAvatar(new string) *DataFanart {
	Data.AuthorAvatar = new
	return Data
}

func (Data *DataFanart) AddPhotos(new []string) *DataFanart {
	Data.Photos = new
	return Data
}

func (Data *DataFanart) AddVideos(new string) *DataFanart {
	Data.Videos = new
	return Data
}

func (Data *DataFanart) AddText(new string) *DataFanart {
	Data.Text = new
	return Data
}

func (Data *DataFanart) AddDynamicID(new string) *DataFanart {
	Data.Dynamic_id = new
	return Data
}

func (Data *DataFanart) AddTweetID(new string) *DataFanart {
	Data.TweetID = new
	return Data
}

func (Data *DataFanart) AddPixivID(new string) *DataFanart {
	Data.PixivID = new
	return Data
}

func (Data *DataFanart) SetLewd(new bool) *DataFanart {
	Data.Lewd = new
	return Data
}

func (Data *DataFanart) SetState(new string) *DataFanart {
	Data.State = new
	return Data
}

func (p *Group) GetRandomFanart() (*DataFanart, error) {
	b, err := GetFanart(p.ID, 0)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (p *Member) GetRandomFanart() (*DataFanart, error) {
	b, err := GetFanart(0, p.ID)
	if err != nil {
		return nil, err
	}
	return b, nil
}
