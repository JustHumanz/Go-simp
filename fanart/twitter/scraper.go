package twitter

import (
	"encoding/json"
	"net/url"
	"strings"
	"sync"

	"github.com/JustHumanz/Go-simp/config"
	database "github.com/JustHumanz/Go-simp/database"
	engine "github.com/JustHumanz/Go-simp/engine"
	log "github.com/sirupsen/logrus"
)

func Changebearer() {
	twbearerold := twbearer
	for _, Token := range config.TwitterToken {
		if Token != twbearer {
			twbearer = Token
		}
	}
	log.WithFields(log.Fields{
		"Old": twbearerold,
		"New": twbearer,
	}).Info("Change Twitter bearer")
}

func CurlTwitter(TwitterPayloads []string, Group int64) {
	var (
		TWdata TwitterStruct
	)
	for _, Payload := range TwitterPayloads {
		var (
			body []byte
			err  error
		)
		url := "https://api.twitter.com/1.1/search/tweets.json?" + Payload
		body, err = engine.Curl(url, []string{"Authorization", "Bearer " + twbearer})
		if err != nil {
			if strings.HasPrefix(err.Error(), "401 Unauthorized") {
				Changebearer()
				body, err = engine.CoolerCurl(url, []string{"Authorization", "Bearer " + twbearer})
				if err != nil {
					log.Error(err, string(body))
					return
				}
			} else {
				log.Error(err, string(body))
				return
			}
		}
		err = json.Unmarshal(body, &TWdata)
		if err != nil {
			log.Error(err)
		}

		wg := new(sync.WaitGroup)
		for _, Statuses := range TWdata.CheckNew() {
			wg.Add(1)
			go Statuses.CheckHashTag(database.GetHashtag(Group), wg)
		}
		wg.Wait()
	}
}

//Start create payload from twitter by group name
func CreatePayload(Group int64) []string {
	var (
		Hashtags      []string
		FinalTags     []string
		GroupsHashtag = database.GetHashtag(Group)
	)
	if len(GroupsHashtag) < 7 {
		for _, hashtag := range database.GetHashtag(Group) {
			if hashtag.TwitterHashtags != "" {
				Hashtags = append(Hashtags, hashtag.TwitterHashtags)
			}
		}
		params := url.Values{}
		params.Add("q", strings.Join(Hashtags, " OR ")+" "+"filter:links -filter:replies filter:media -filter:retweets")

		log.WithFields(log.Fields{
			"Hashtags": strings.Join(Hashtags, " OR "),
		}).Info("Start Curl new art")
		return []string{params.Encode()}
	} else {
		for z, hashtag := range GroupsHashtag {
			if hashtag.TwitterHashtags != "" {
				Hashtags = append(Hashtags, hashtag.TwitterHashtags)
				if z == len(GroupsHashtag)/3 || z == len(GroupsHashtag)-1 {
					params := url.Values{}
					params.Add("q", strings.Join(Hashtags, " OR ")+" "+"filter:links -filter:replies filter:media -filter:retweets")
					FinalTags = append(FinalTags, params.Encode())
				}
			}
		}
		log.WithFields(log.Fields{
			"Hashtags": strings.Join(Hashtags, " OR "),
		}).Info("Start Curl new art")
		return FinalTags
	}
}
