package main

import (
	"strconv"
	"strings"
	"time"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	database "github.com/JustHumanz/Go-Simp/pkg/database"
	engine "github.com/JustHumanz/Go-Simp/pkg/engine"
	"github.com/bwmarrin/discordgo"
	"github.com/hako/durafmt"
	log "github.com/sirupsen/logrus"
)

func TwitchMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	Prefix := configfile.BotPrefix.Twitch
	m.Content = strings.ToLower(m.Content)
	if strings.HasPrefix(m.Content, Prefix) {
		CommandArray := strings.Split(m.Content, " ")
		if len(CommandArray) > 1 {
			Payload := strings.Split(strings.TrimSpace(CommandArray[1]), ",")
			if CommandArray[0] == Prefix+Live {
				for _, FindGroupArry := range Payload {
					VTuberGroup, err := FindGropName(FindGroupArry)
					if err != nil {
						Member := FindVtuber(FindGroupArry)
						if Member == (database.Member{}) {
							s.ChannelMessageSend(m.ChannelID, "`"+FindGroupArry+"`,Name of Vtuber Group or Vtuber Name was not found")
							return
						} else {
							LiveTwitch := database.TwitchGet(0, Member.ID, config.LiveStatus)
							FixName := engine.FixName(Member.EnName, Member.JpName)
							if LiveTwitch != nil {
								Color, err := engine.GetColor(config.TmpDir, m.Author.AvatarURL("128"))
								if err != nil {
									log.Error(err)
								}
								for _, LiveData := range LiveTwitch {
									loc := engine.Zawarudo(Member.Region)
									diff := time.Now().In(loc).Sub(LiveData.Schedul.In(loc))
									view, err := strconv.Atoi(LiveData.Viewers)
									if err != nil {
										log.Error(err)
									}
									_, err = s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
										SetTitle(FixName).
										SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
										SetThumbnail(Member.TwitchAvatar).
										SetImage(LiveData.Thumb).
										SetURL("https://twitch.tv/"+Member.TwitchName).
										AddField("Start live", durafmt.Parse(diff).LimitFirstN(2).String()+" Ago").
										AddField("Viewers", engine.NearestThousandFormat(float64(view))).
										InlineAllFields().
										AddField("Game", LiveData.Game).
										SetColor(Color).
										SetFooter(LiveData.Schedul.In(loc).Format(time.RFC822)).MessageEmbed)
									if err != nil {
										log.Error(err)
									}
								}
							} else {
								_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
									SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
									SetDescription("It looks like `"+FixName+"` doesn't have a livestream right now").
									SetImage(config.WorryIMG).MessageEmbed)
								if err != nil {
									log.Error(err)
								}
							}
						}
					} else {
						TwitchLive := database.TwitchGet(VTuberGroup.ID, 0, config.LiveStatus)
						if TwitchLive != nil {
							Color, err := engine.GetColor(config.TmpDir, m.Author.AvatarURL("128"))
							if err != nil {
								log.Error(err)
							}

							for _, LiveData := range TwitchLive {
								LiveData.AddMember(FindVtuber(LiveData.Member.ID))
								loc := engine.Zawarudo(LiveData.Member.Region)
								FixName := engine.FixName(LiveData.Member.EnName, LiveData.Member.JpName)
								diff := time.Now().In(loc).Sub(LiveData.Schedul.In(loc))
								view, err := strconv.Atoi(LiveData.Viewers)
								if err != nil {
									log.Error(err)
								}
								_, err = s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
									SetTitle(FixName).
									SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
									SetThumbnail(LiveData.Member.TwitchAvatar).
									SetImage(LiveData.Thumb).
									SetURL("https://twitch.tv/"+LiveData.Member.TwitchName).
									AddField("Start live", durafmt.Parse(diff).LimitFirstN(2).String()+" Ago").
									AddField("Viewers", engine.NearestThousandFormat(float64(view))).
									InlineAllFields().
									AddField("Game", LiveData.Game).
									SetColor(Color).
									SetFooter(LiveData.Schedul.In(loc).Format(time.RFC822)).MessageEmbed)
								if err != nil {
									log.Error(err)
								}
							}
						} else {
							_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
								SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
								SetDescription("It looks like `"+VTuberGroup.GroupName+"` doesn't have a livestream right now").
								SetImage(config.WorryIMG).MessageEmbed)
							if err != nil {
								log.Error(err)
							}
						}
					}
				}

			} else if CommandArray[0] == Prefix+Past || CommandArray[0] == Prefix+"last" {
				for _, FindGroupArry := range Payload {
					VTuberGroup, err := FindGropName(FindGroupArry)
					if err != nil {
						Member := FindVtuber(FindGroupArry)
						if Member == (database.Member{}) {
							s.ChannelMessageSend(m.ChannelID, "`"+FindGroupArry+"`,Name of Vtuber Group or Vtuber Name was not found")
							return
						} else {
							TwitchLive := database.TwitchGet(0, Member.ID, config.PastStatus)
							FixName := engine.FixName(Member.JpName, Member.EnName)
							if TwitchLive != nil {
								Color, err := engine.GetColor(config.TmpDir, m.Author.AvatarURL("128"))
								if err != nil {
									log.Error(err)
								}

								for _, LiveData := range TwitchLive {
									loc := engine.Zawarudo(Member.Region)
									diff := LiveData.Schedul.In(loc).Sub(time.Now().In(loc))
									view, err := strconv.Atoi(LiveData.Viewers)
									if err != nil {
										log.Error(err)
									}
									_, err = s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
										SetTitle(FixName).
										SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
										SetThumbnail(Member.TwitchAvatar).
										SetImage(LiveData.Thumb).
										SetURL("https://twitch.tv/"+Member.TwitchName).
										AddField("Start live", durafmt.Parse(diff).LimitFirstN(2).String()+" Ago").
										AddField("Viewers", engine.NearestThousandFormat(float64(view))).
										InlineAllFields().
										AddField("Game", LiveData.Game).
										SetColor(Color).
										SetFooter(LiveData.Schedul.In(loc).Format(time.RFC822)).MessageEmbed)
									if err != nil {
										log.Error(err)
									}
								}
							} else {
								_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
									SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
									SetDescription("It looks like `"+FixName+"` doesn't have a Past livestream right now").
									SetImage(config.WorryIMG).MessageEmbed)
								if err != nil {
									log.Error(err)
								}
							}
						}
					} else {
						TwitchLive := database.TwitchGet(VTuberGroup.ID, 0, config.PastStatus)
						if TwitchLive != nil {
							Color, err := engine.GetColor(config.TmpDir, m.Author.AvatarURL("128"))
							if err != nil {
								log.Error(err)
							}

							for _, LiveData := range TwitchLive {
								LiveData.AddMember(FindVtuber(LiveData.Member.ID))
								loc := engine.Zawarudo(LiveData.Member.Region)

								FixName := engine.FixName(LiveData.Member.EnName, LiveData.Member.JpName)
								view, err := strconv.Atoi(LiveData.Viewers)
								if err != nil {
									log.Error(err)
								}

								diff := time.Now().In(loc).Sub(LiveData.Schedul.In(loc))
								_, err = s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
									SetTitle(FixName).
									SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
									SetThumbnail(LiveData.Member.TwitchAvatar).
									SetImage(LiveData.Thumb).
									SetURL("https://twitch.tv/"+LiveData.Member.TwitchName).
									AddField("Start live", durafmt.Parse(diff).LimitFirstN(2).String()+" Ago").
									AddField("Viewers", engine.NearestThousandFormat(float64(view))).
									InlineAllFields().
									AddField("Game", LiveData.Game).
									SetColor(Color).
									SetFooter(LiveData.Schedul.In(loc).Format(time.RFC822)).MessageEmbed)
								if err != nil {
									log.Error(err)
								}
							}
						} else {
							_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
								SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
								SetDescription("It looks like `"+VTuberGroup.GroupName+"` doesn't have a Past livestream right now").
								SetImage(config.WorryIMG).MessageEmbed)
							if err != nil {
								log.Error(err)
							}
						}
					}
				}
			}
		} else {
			_, err := s.ChannelMessageSend(m.ChannelID, "Incomplete `"+Prefix+"` command")
			if err != nil {
				log.Error(err)
			}
		}
	}
}
