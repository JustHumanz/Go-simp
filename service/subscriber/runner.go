package main

import (
	"context"
	"encoding/json"
	"flag"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/nicklaw5/helix"
	"github.com/robfig/cron/v3"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/JustHumanz/Go-Simp/pkg/database"
	network "github.com/JustHumanz/Go-Simp/pkg/network"
	pilot "github.com/JustHumanz/Go-Simp/service/pilot/grpc"
	"github.com/JustHumanz/Go-Simp/service/utility/runfunc"
	log "github.com/sirupsen/logrus"
)

var (
	Bot          *discordgo.Session
	configfile   config.ConfigFile
	Payload      *[]database.Group
	gRCPconn     pilot.PilotServiceClient
	TwitchClient *helix.Client
	Youtube      *bool
	BiliBili     *bool
	Twitter      *bool
	Twitch       *bool
)

const (
	ModuleState = "Twitter Fanart"
)

func init() {
	Youtube = flag.Bool("Youtube", false, "Enable youtube module")
	BiliBili = flag.Bool("BiliBili", false, "Enable bilibili module")
	Twitter = flag.Bool("Twitter", false, "Enable twitter module")
	Twitch = flag.Bool("Twitch", false, "Enable Twitch module")
	flag.Parse()

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: true})
	gRCPconn = pilot.NewPilotServiceClient(network.InitgRPC(config.Pilot))
}

func main() {
	GetPayload := func() {
		res, err := gRCPconn.ReqData(context.Background(), &pilot.ServiceMessage{
			Message: "Send me nude",
			Service: ModuleState,
		})
		if err != nil {
			if configfile.Discord != "" {
				pilot.ReportDeadService(err.Error(), ModuleState)
			}
			log.Error("Error when request payload: %s", err)
		}
		err = json.Unmarshal(res.ConfigFile, &configfile)
		if err != nil {
			log.Error(err)
		}

		err = json.Unmarshal(res.VtuberPayload, &Payload)
		if err != nil {
			log.Error(err)
		}
	}

	GetPayload()
	configfile.InitConf()

	var err error
	Bot, err = discordgo.New("Bot " + configfile.Discord)
	if err != nil {
		log.Error(err)
	}

	database.Start(configfile)

	TwitchClient, err = helix.NewClient(&helix.Options{
		ClientID:     config.GoSimpConf.Twitch.ClientID,
		ClientSecret: config.GoSimpConf.Twitch.ClientSecret,
	})
	if err != nil {
		log.Error(err)
	}
	resp, err := TwitchClient.RequestAppAccessToken([]string{"user:read:email"})
	if err != nil {
		log.Error(err)
	}

	TwitchClient.SetAppAccessToken(resp.Data.AccessToken)

	c := cron.New()
	c.Start()

	c.AddFunc(config.CheckPayload, GetPayload)
	if *Youtube {
		c.AddFunc(config.YoutubeSubscriber, CheckYoutube)
		log.Info("Add youtube subscriber to cronjob")
	}

	if *BiliBili {
		c.AddFunc(config.BiliBiliFollowers, CheckBiliBili)
		log.Info("Add bilibili followers to cronjob")
	}

	if *Twitter {
		c.AddFunc(config.TwitterFollowers, CheckTwitter)
		log.Info("Add twitter followers to cronjob")
	}

	if *Twitch {
		c.AddFunc(config.TwitchFollowers, CheckTwitch)
		log.Info("Add twitch followers to cronjob")
	}

	_, err = gRCPconn.ModuleList(context.Background(), &pilot.ModuleData{
		Module:  "YoutubeSubscriber",
		Enabled: *Youtube,
	})
	if err != nil {
		log.Error(err)
	}
	_, err = gRCPconn.ModuleList(context.Background(), &pilot.ModuleData{
		Module:  "BiliBiliFollowers",
		Enabled: *BiliBili,
	})
	if err != nil {
		log.Error(err)
	}
	_, err = gRCPconn.ModuleList(context.Background(), &pilot.ModuleData{
		Module:  "TwitterFollowers",
		Enabled: *Twitter,
	})
	if err != nil {
		log.Error(err)
	}

	go pilot.RunHeartBeat(gRCPconn, "Subscriber")
	runfunc.Run(Bot)
}

func SendNude(Embed *discordgo.MessageEmbed, Group database.Group, Member database.Member) {
	if match, _ := regexp.MatchString("404.jpg", Group.IconURL); match {
		Embed.Author.IconURL = ""
	}
	ChannelData, err := Group.GetChannelByGroup(Member.Region)
	if err != nil {
		log.Error(err)
	}
	for i, Channel := range ChannelData {

		Channel.SetMember(Member)
		Tmp := &Channel
		ctx := context.Background()
		UserTagsList, err := Tmp.SetMember(Member).SetGroup(Group).GetUserList(ctx)
		if err != nil {
			log.Error(err)
		}
		msg, err := Bot.ChannelMessageSendEmbed(Channel.ChannelID, Embed)
		if err != nil {
			log.Error(msg, err)
		}
		if UserTagsList != nil {
			msg, err = Bot.ChannelMessageSend(Channel.ChannelID, "UserTags: "+strings.Join(UserTagsList, " "))
			if err != nil {
				log.Error(msg, err)
			}
		}
		if i%config.Waiting == 0 && configfile.LowResources {
			log.WithFields(log.Fields{
				"Func":  "Subscriber",
				"Value": config.Waiting,
			}).Warn("Waiting send message")
			time.Sleep(100 * time.Millisecond)
		}
	}
}

type Subs struct {
	Kind     string `json:"kind"`
	Etag     string `json:"etag"`
	PageInfo struct {
		ResultsPerPage int `json:"resultsPerPage"`
	} `json:"pageInfo"`
	Items []struct {
		Kind       string `json:"kind"`
		Etag       string `json:"etag"`
		ID         string `json:"id"`
		Statistics struct {
			ViewCount             string `json:"viewCount"`
			CommentCount          string `json:"commentCount"`
			SubscriberCount       string `json:"subscriberCount"`
			HiddenSubscriberCount bool   `json:"hiddenSubscriberCount"`
			VideoCount            string `json:"videoCount"`
		} `json:"statistics"`
	} `json:"items"`
}

type BiliBiliStat struct {
	LikeView LikeView
	Follow   BiliFollow
	Videos   int
}

type LikeView struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    struct {
		Archive struct {
			View int `json:"view"`
		} `json:"archive"`
		Article struct {
			View int `json:"view"`
		} `json:"article"`
		Likes int `json:"likes"`
	} `json:"data"`
}

type BiliFollow struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    struct {
		Mid       int `json:"mid"`
		Following int `json:"following"`
		Whisper   int `json:"whisper"`
		Black     int `json:"black"`
		Follower  int `json:"follower"`
	} `json:"data"`
}
