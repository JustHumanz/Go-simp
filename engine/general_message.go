package engine

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"sync"

	config "github.com/JustHumanz/Go-simp/config"
	database "github.com/JustHumanz/Go-simp/database"
	"github.com/bwmarrin/discordgo"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
)

//Fanart discord message handler
func Fanart(s *discordgo.Session, m *discordgo.MessageCreate) {
	m.Content = strings.ToLower(m.Content)
	Prefix := config.PFanart
	var (
		Member      bool
		Group       bool
		Pic         = config.NotFound
		Msg         string
		wg          sync.WaitGroup
		embed       *discordgo.MessageEmbed
		DynamicData Dynamic_svr
	)

	Color, err := GetColor("/tmp/mem.tmp", m.Author.AvatarURL("80"))

	if strings.HasPrefix(m.Content, Prefix) {
		SendNude := func(Title, Author, Text, URL, Pic, Msg string, Color int, State, Dynamic string) bool {
			Msg = Msg + " *sometimes image not showing,because image oversize*"
			if State == "TBiliBili" {
				var (
					body    []byte
					errcurl error
					urls    = "https://api.vc.bilibili.com/dynamic_svr/v1/dynamic_svr/get_dynamic_detail?dynamic_id=" + Dynamic
				)
				body, errcurl = Curl(urls, nil)
				if errcurl != nil {
					log.Error(errcurl, string(body))
					log.Info("Trying use tor")

					body, errcurl = CoolerCurl(urls)
					if errcurl != nil {
						log.Error(errcurl)
					}
				}
				json.Unmarshal(body, &DynamicData)
				embed = NewEmbed().
					SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
					SetTitle(Author).
					SetThumbnail(DynamicData.GetUserAvatar()).
					SetDescription(Text).
					SetURL(URL).
					SetImage(Pic).
					SetColor(Color).
					InlineAllFields().
					SetFooter(Msg, config.TwitterIMG).MessageEmbed
			} else {
				embed = NewEmbed().
					SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
					SetTitle(Author).
					SetThumbnail(GetUserAvatar(Author)).
					SetDescription(Text).
					SetURL(URL).
					SetImage(Pic).
					SetColor(Color).
					InlineAllFields().
					SetFooter(Msg, config.TwitterIMG).MessageEmbed
			}
			msg, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
			if err != nil {
				log.Error(err, msg)
			}
			err = Reacting(map[string]string{
				"ChannelID": m.ChannelID,
				"Content":   m.Content,
				"Prefix":    Prefix,
			})
			if err != nil {
				log.Error(err)
			}
			return true
		}
		for i := 0; i < len(GroupData); i++ {
			Data2 := database.GetName(GroupData[i].ID)
			GroupData[i].NameGroup = strings.ToLower(GroupData[i].NameGroup)
			wg.Add(2)

			go func() {
				for ii := 0; ii < len(Data2); ii++ {
					Data2[ii].Name = strings.ToLower(Data2[ii].Name)
					if m.Content == Prefix+Data2[ii].Name || m.Content == Prefix+Data2[ii].JpName {
						DataFix := Data2[ii].GetMemberURL()
						if DataFix.Videos != "" {
							Msg = "Video type,check original post"
						} else if DataFix.Photos != nil {
							Pic = DataFix.Photos[0]
							Color, err = GetColor("/tmp/mem.tmp", DataFix.Photos[0])
							if err != nil {
								log.Error(err)
							}
						}
						Member = SendNude(FixName(Data2[ii].EnName, Data2[ii].JpName),
							DataFix.Author, RemovePic(DataFix.Text),
							DataFix.PermanentURL,
							Pic, Msg, Color,
							DataFix.State, DataFix.Dynamic_id)
						break
					} else {
						Member = false
					}
				}
				wg.Done()
			}()
			go func() {
				if m.Content == Prefix+GroupData[i].NameGroup {
					DataFix := GroupData[i].GetGroupURL()
					if DataFix.Videos != "" {
						Pic = DataFix.Videos
						Msg = "Video type,check original post"
					} else if DataFix.Photos != nil {
						Pic = DataFix.Photos[0]
						Color, err = GetColor("/tmp/mem.tmp", DataFix.Photos[0])
						if err != nil {
							log.Error(err)
						}
					}

					Group = SendNude(FixName(DataFix.EnName, DataFix.JpName),
						DataFix.Author, RemovePic(DataFix.Text),
						DataFix.PermanentURL,
						Pic, Msg, Color,
						DataFix.State, DataFix.Dynamic_id)
				} else {
					Group = false
				}
				wg.Done()
			}()
			wg.Wait()

			if Member || Group {
				return
			}
		}
		if !Group && !Member {
			s.ChannelMessageSend(m.ChannelID, "`"+m.Content[len(Prefix):]+"` was invalid name")
		}
	}
}

//Get Guild name *Owner only*
func Humanz(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == config.OwnerDiscordID {
		if m.Content == "!>list" {
			for _, Guild := range s.State.Guilds {
				s.ChannelMessageSend(m.ChannelID, Guild.Name)
			}
		}
	}
}

//Tags command message handler
func Tags(s *discordgo.Session, m *discordgo.MessageCreate) {
	Prefix := config.PGeneral
	m.Content = strings.ToLower(m.Content)
	if strings.HasPrefix(m.Content, Prefix) {
		var (
			counter bool
		)
		User := database.UserStruct{
			DiscordID:       m.Author.ID,
			DiscordUserName: m.Author.Username,
			Channel_ID:      m.ChannelID,
		}
		if strings.HasPrefix(m.Content, Prefix+"tag me") {
			VtuberName := strings.TrimSpace(strings.Replace(m.Content, Prefix+"tag me", "", -1))
			Tag := func(NameArry, GroupName string, GroupID, MemberID int64, Group bool) {
				if database.CheckChannelEnable(m.ChannelID, NameArry, GroupID) {
					User.GroupID = GroupID
					var Already []string
					if Group {
						for _, Member := range database.GetName(GroupID) {
							err := User.Adduser(Member.ID)
							if err != nil {
								Already = append(Already, "`"+Member.Name+"`")
							} else {
								counter = true
							}
						}
						if Already != nil {
							s.ChannelMessageSend(m.ChannelID, "<@"+m.Author.ID+"> Already Added "+strings.Join(Already, " ")+" **"+GroupName+"**")
						}
					} else {
						err := User.Adduser(MemberID)
						if err != nil {
							s.ChannelMessageSend(m.ChannelID, "<@"+m.Author.ID+"> Already Added `"+NameArry+"`")
						} else {
							counter = true
						}
					}
				} else {
					s.ChannelMessageSend(m.ChannelID, "look like this channel not enable `"+GroupName+"`")
				}
			}

			if (VtuberName) != "" {
				FindNameArry := strings.Split(VtuberName, ",")
				for i := 0; i < len(FindNameArry); i++ {
					Data, err := FindName(FindNameArry[i])
					if err != nil {
						VTuberGroup, err := FindGropName(FindNameArry[i])
						if err != nil {
							s.ChannelMessageSend(m.ChannelID, "`"+FindNameArry[i]+"` was invalid")
							return
						}
						Tag(FindNameArry[i], VTuberGroup.NameGroup, VTuberGroup.ID, 0, true)
					} else {
						Tag(FindNameArry[i], Data.GroupName, Data.GroupID, Data.MemberID, false)
					}
				}
				if counter {
					s.ChannelMessageSend(m.ChannelID, "done")
				}
			} else {
				s.ChannelMessageSend(m.ChannelID, "Incomplete `tag me` command")
			}
		} else if strings.HasPrefix(m.Content, Prefix+"del tag") {
			VtuberName := strings.TrimSpace(strings.Replace(m.Content, Prefix+"del tag", "", -1))
			Del := func(NameArry, GroupName string, GroupID, MemberID int64, Group bool) {
				if database.CheckChannelEnable(m.ChannelID, NameArry, GroupID) {
					User.GroupID = GroupID
					if Group {
						for _, Mem := range database.GetName(GroupID) {
							err := User.Deluser(Mem.ID)
							if err != nil {
								log.Error(err)
								s.ChannelMessageSend(m.ChannelID, "`"+Mem.Name+"`,Already removed or This channel not Enable "+GroupName)
							} else {
								counter = true
							}
						}
					} else {
						err := User.Deluser(MemberID)
						if err != nil {
							log.Error(err)
							s.ChannelMessageSend(m.ChannelID, "`"+NameArry+"`,Already removed or This channel not Enable "+GroupName)
						} else {
							counter = true
						}
					}
				} else {
					s.ChannelMessageSend(m.ChannelID, "look like this channel not enable `"+NameArry+"`")
				}
			}

			if (VtuberName) != "" {
				FindNameArry := strings.Split(VtuberName, ",")
				for i := 0; i < len(FindNameArry); i++ {
					Data, err := FindName(FindNameArry[i])
					if err != nil {
						VTuberGroup, err := FindGropName(FindNameArry[i])
						if err != nil {
							s.ChannelMessageSend(m.ChannelID, "`"+FindNameArry[i]+"` was invalid")
							return
						}
						Del(FindNameArry[i], VTuberGroup.NameGroup, VTuberGroup.ID, 0, true)
					} else {
						Del(FindNameArry[i], Data.GroupName, Data.GroupID, Data.MemberID, false)
					}
				}
				if counter {
					s.ChannelMessageSend(m.ChannelID, "done")
				}
			} else {
				s.ChannelMessageSend(m.ChannelID, "Incomplete del tag command")
			}
		}
	}
}

//Check user permission
func CheckPermission(User, Channel string) bool {
	Debugging(GetFunctionName(CheckPermission), "In", fmt.Sprint(User, Channel))
	a, err := BotSession.UserChannelPermissions(User, Channel)
	BruhMoment(err, "", false)
	Permission := 16
	if a&Permission != 0 {
		Debugging(GetFunctionName(CheckPermission), "Out", true)
		return true
	} else {
		Debugging(GetFunctionName(CheckPermission), "Out", false)
		return false
	}
}

//Enable command message handler
func Enable(s *discordgo.Session, m *discordgo.MessageCreate) {
	m.Content = strings.ToLower(m.Content)
	Prefix := config.PGeneral
	if strings.HasPrefix(m.Content, Prefix) {
		var (
			counter bool
			tagtype int
		)
		CommandArray := strings.Split(m.Content, " ")
		if CommandArray[0] == Prefix+"enable" {
			if len(CommandArray) > 1 {
				if CommandArray[1] == "art" {
					tagtype = 1
				} else if CommandArray[1] == "live" {
					tagtype = 2
				} else {
					tagtype = 3
				}
				FindGroupArry := strings.Split(strings.TrimSpace(CommandArray[len(CommandArray)-1]), ",")

				for i := 0; i < len(FindGroupArry); i++ {
					VTuberGroup, err := FindGropName(FindGroupArry[i])
					if err != nil {
						s.ChannelMessageSend(m.ChannelID, "`"+FindGroupArry[i]+"`,Name of Vtuber Group was invalid")
						return
					}
					if CheckPermission(m.Author.ID, m.ChannelID) {
						if database.ChannelCheck(VTuberGroup.ID, m.ChannelID) {
							counter = false
						} else {
							err := database.AddChannel(m.ChannelID, tagtype, VTuberGroup.ID)
							if err != nil {
								log.Error(err)
								s.ChannelMessageSend(m.ChannelID, "Something error XD")
							}
							counter = true
						}
					} else {
						s.ChannelMessageSend(m.ChannelID, "You don't have permission to enable/disable/update")
						return
					}
				}
				if counter {
					s.ChannelMessageSend(m.ChannelID, "done")
				} else {
					s.ChannelMessageSend(m.ChannelID, "already added")
				}
			} else {
				s.ChannelMessageSend(m.ChannelID, "Incomplete enable command")
			}
		} else if CommandArray[0] == Prefix+"disable" {
			if len(CommandArray) > 1 {
				FindGroupArry := strings.Split(strings.TrimSpace(CommandArray[1]), ",")
				for i := 0; i < len(FindGroupArry); i++ {
					VTuberGroup, err := FindGropName(FindGroupArry[i])
					if err != nil {
						s.ChannelMessageSend(m.ChannelID, "`"+FindGroupArry[i]+"`,Name of Vtuber Group was not valid")
						return
					}
					if CheckPermission(m.Author.ID, m.ChannelID) {
						if database.ChannelCheck(VTuberGroup.ID, m.ChannelID) {
							err := database.DelChannel(m.ChannelID, VTuberGroup.ID)
							if err != nil {
								log.Error(err)
								s.ChannelMessageSend(m.ChannelID, "Something error XD")
								return
							}
							counter = true
						} else {
							counter = false
						}
					} else {
						s.ChannelMessageSend(m.ChannelID, "You don't have permission to enable/disable/update")
						return
					}
				}
				if counter {
					s.ChannelMessageSend(m.ChannelID, "done")
				} else {
					s.ChannelMessageSend(m.ChannelID, "already removed")
				}
			}
		} else if CommandArray[0] == Prefix+"update" {
			if len(CommandArray) > 1 {
				if CommandArray[1] == "art" {
					tagtype = 1
				} else if CommandArray[1] == "live" {
					tagtype = 2
				} else {
					tagtype = 3
				}
				FindGroupArry := strings.Split(strings.TrimSpace(CommandArray[len(CommandArray)-1]), ",")

				for i := 0; i < len(FindGroupArry); i++ {
					VTuberGroup, err := FindGropName(FindGroupArry[i])
					if err != nil {
						s.ChannelMessageSend(m.ChannelID, "`"+FindGroupArry[i]+"`,Name of Vtuber Group was invalid")
						return
					}
					if CheckPermission(m.Author.ID, m.ChannelID) {
						err := database.UpdateChannel(m.ChannelID, tagtype, VTuberGroup.ID)
						if err != nil {
							counter = false
						} else {
							counter = true
						}
					} else {
						s.ChannelMessageSend(m.ChannelID, "You don't have permission to enable/disable/update")
						return
					}
				}
				if counter {
					s.ChannelMessageSend(m.ChannelID, "done")
				} else {
					s.ChannelMessageSend(m.ChannelID, "already added")
				}
			} else {
				s.ChannelMessageSend(m.ChannelID, "Incomplete update command")
			}
		}
	}
}

//helmp command message handler
func Help(s *discordgo.Session, m *discordgo.MessageCreate) {
	m.Content = strings.ToLower(m.Content)
	Prefix := config.PGeneral
	Color, err := GetColor("/tmp/discordpp.tmp", m.Author.AvatarURL("128"))
	if err != nil {
		log.Error(err)
	}
	if strings.HasPrefix(m.Content, Prefix) {
		if m.Content == Prefix+"help en" || m.Content == Prefix+"help" {
			s.ChannelMessageSendEmbed(m.ChannelID, NewEmbed().
				SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
				SetTitle("Help").
				AddField(Prefix+"Enable [art/live/all] [Vtuber Group]", "This command will declare if [Vtuber Group] enable in this channel\nExample:\n`"+config.PGeneral+"enable hanayori` so other users can use `"+config.PGeneral+"tag me kanochi`").
				AddField(Prefix+"Update [art/live/all] [Vtuber Group]", "Use this command if you want to change enable state").
				AddField(Prefix+"Disable [Vtuber Group]", "Just like enable but this disable command :3 ").
				AddField(config.PFanart+"[Group/Member name]", "Show fanart with randomly with their fanart hashtag\nExample: \n`"+config.PFanart+"Kanochi` or `"+config.PFanart+"hololive`").
				AddField(Prefix+"Tag me [Group/Member name]", "This command will add you to tags list if any new fanart\nExample: \n`"+config.PFanart+"tag me Kanochi`,then you will get tagged when there is a new fanart of kano").
				AddField(Prefix+"Del tag [Group/Member name]", "This command will remove you from tags list").
				AddField(Prefix+"My tags", "Show all lists that you are subscribed").
				AddField(Prefix+"Channel tags", "Show what is enable in this channel").
				AddField(Prefix+"Vtuber data", "Show available Vtuber data ").
				AddField(config.PYoutube+"Upcoming [Vtuber Group/Member]", "This command will show all Upcoming live streams on Youtube").
				AddField(config.PYoutube+"Live [Vtuber Group/Member]", "This command will show all live streams right now on Youtube").
				AddField(config.PYoutube+"Last [Vtuber Group/Member]", "This command will show all past streams on Youtube [only 5]").
				AddField(config.PYoutube+"[Upcoming/Live/Last] [Member name]", "This command will show all Vtuber member Upcoming/Live/Past streams on Youtube").
				AddField("~~"+config.PBilibili+"Upcoming [Vtuber Group/Member]~~", "~~This command will show all Upcoming live streams on BiliBili~~").
				AddField(config.PBilibili+"Live [Vtuber Group/Member]", "This command will show all live streams right now on BiliBili").
				AddField(config.PBilibili+"Last [Vtuber Group/Member]", "This command will show all past streams on BiliBili").
				AddField("sp_"+config.PBilibili+"[Vtuber Group/Member]", "This command will show latest video upload on BiliBili").
				AddField(config.PBilibili+"[Upcoming/Live/Last] [Member name]", "This command will show all Vtuber member Upcoming/Live/Past streams on BiliBili").
				AddField(Prefix+"Help EN", "Well,you using it right now").
				AddField(Prefix+"Help JP", "Like this but in Japanese").
				SetThumbnail("https://justhumanz.me/bsd.png").
				SetFooter("Only user with permission \"Manage Channel or higher\" can Enable/Disable/Update Vtuber Group").
				SetColor(Color).MessageEmbed)
			return
		} else if m.Content == Prefix+"help jp" { //i'm just joking lol
			s.ChannelMessageSendEmbed(m.ChannelID, NewEmbed().
				SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
				SetTitle("Help").
				SetDescription("日本語が話せるようになってヘルプメニューを作りたい\n~Dev").
				SetImage("https://i.pinimg.com/736x/c8/1f/55/c81f554101a36e4b853318b7f937fcbc.jpg").
				SetFooter("More like,help me").
				SetColor(Color).MessageEmbed)
			return
		}
	}
}

//Status command message handler
func Status(s *discordgo.Session, m *discordgo.MessageCreate) {
	m.Content = strings.ToLower(m.Content)
	Prefix := config.PGeneral
	Color, err := GetColor("/tmp/discordpp", m.Author.AvatarURL("128"))
	if err != nil {
		log.Error(err)
	}

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)

	if strings.HasPrefix(m.Content, Prefix) {
		if m.Content == Prefix+"my tags" {
			list := database.UserStatus(m.Author.ID, m.ChannelID)
			if list != nil {
				s.ChannelMessageSendEmbed(m.ChannelID, NewEmbed().
					SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
					SetThumbnail(m.Author.AvatarURL("128")).
					AddField("Username", m.Author.Username).
					AddField("Vtuber Tag", strings.Join(list, ",")).
					SetColor(Color).MessageEmbed)
			} else {
				s.ChannelMessageSendEmbed(m.ChannelID, NewEmbed().
					SetTitle("404 Not found").
					SetImage(config.NotFound).
					SetColor(Color).MessageEmbed)
			}
		} else if m.Content == Prefix+"channel tags" {
			list, Type := database.ChannelStatus(m.ChannelID)
			if list != nil {
				var (
					tbl     [][]string
					Typestr string
				)
				for i := 0; i < len(list); i++ {
					if Type[i] == 1 {
						Typestr = "Art"
					} else if Type[i] == 2 {
						Typestr = "Live"
					} else {
						Typestr = "All"
					}
					tmp := [][]string{[]string{list[i], Typestr}}
					tbl = append(tbl, tmp...)
				}
				table.SetHeader([]string{"Vtuber Group", "Type"})
				for _, v := range tbl {
					table.Append(v)
				}
				table.Render()

				s.ChannelMessageSendEmbed(m.ChannelID, NewEmbed().
					SetAuthor(m.Author.Username, m.Author.AvatarURL("80")).
					SetDescription("```\r"+tableString.String()+"```").
					SetThumbnail(config.Go_Simp).
					SetColor(Color).MessageEmbed)
			} else {
				s.ChannelMessageSendEmbed(m.ChannelID, NewEmbed().
					SetTitle("404 Not found").
					SetImage(config.NotFound).
					SetColor(Color).MessageEmbed)
			}
		} else if m.Content == Prefix+"vtuber data" {
			var (
				tbl   [][]string
				Data2 []database.Name
			)
			Data := database.GetGroup()
			for i := 0; i < len(Data); i++ {
				Data2 = database.GetName(Data[i].ID)
				for ii := 0; ii < len(Data2); ii++ {
					tmp := [][]string{[]string{Data2[ii].Name, Data[i].NameGroup}}
					tbl = append(tbl, tmp...)
				}
			}

			table.SetHeader([]string{"Nickname", "Group"})
			for _, v := range tbl {
				table.Append(v)
			}
			table.Render()
			embed := NewEmbed().
				SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
				SetThumbnail(config.Go_Simp).
				SetDescription("```" + tableString.String() + "```").
				SetColor(Color).
				SetFooter("use `Nickname` as parameter").MessageEmbed
			s.ChannelMessageSendEmbed(m.ChannelID, embed)
		}
	}
}

//Find a valid Vtuber name from message handler
func FindName(MemberName string) (NameStruct, error) {
	Data := database.GetGroup()
	for i := 0; i < len(Data); i++ {
		Names := database.GetName(Data[i].ID)
		for _, Name := range Names {
			if strings.ToLower(Name.Name) == MemberName || strings.ToLower(Name.JpName) == MemberName {
				return NameStruct{
					GroupName: Data[i].NameGroup,
					GroupID:   Data[i].ID,
					MemberID:  Name.ID,
				}, nil
			}
		}
	}
	return NameStruct{}, errors.New("Name Vtuber not valid")
}

type NameStruct struct {
	GroupName string
	GroupID   int64
	MemberID  int64
}

//Find a valid Vtuber Group from message handler
func FindGropName(GroupName string) (database.GroupName, error) {
	funcvar := GetFunctionName(FindGropName)
	Debugging(funcvar, "In", GroupName)
	Data := database.GetGroup()
	for i := 0; i < len(Data); i++ {
		if strings.ToLower(Data[i].NameGroup) == strings.ToLower(GroupName) {
			Debugging(funcvar, "Out", fmt.Sprint(Data[i].ID, nil))
			return Data[i], nil
		}
	}
	Debugging(funcvar, "Out", 0)
	return database.GroupName{}, errors.New(GroupName + " Name Vtuber not valid")
}

//Remove twitter pic
func RemovePic(text string) string {
	return regexp.MustCompile(`(?m)^(.*?)pic\.twitter.com\/.+`).ReplaceAllString(text, "${1}$2")
}

func Reacting(Data map[string]string) error {
	EmojiList := config.EmojiFanart
	ChannelID := Data["ChannelID"]
	MessID, err := BotSession.Channel(ChannelID)
	if err != nil {
		return errors.New(err.Error() + " ChannelID: " + ChannelID)
	}
	for l := 0; l < len(EmojiList); l++ {
		if Data["Content"][len(Data["Prefix"]):] == "kanochi" { //don't change this ("kanochi") *kalau di rubah w tandain lo a*g >:'( *
			err := BotSession.MessageReactionAdd(ChannelID, MessID.LastMessageID, EmojiList[0])
			if err != nil {
				return errors.New(err.Error() + " ChannelID: " + ChannelID)
				//log.Error(err, ChannelID)
			}
			break
		} else {
			err := BotSession.MessageReactionAdd(ChannelID, MessID.LastMessageID, EmojiList[l])
			if err != nil {
				return errors.New(err.Error() + " ChannelID: " + ChannelID)
				//log.Error(err, ChannelID)
				//break
			}
		}
	}
	return nil
}

//Get twitter avatar
func GetUserAvatar(username string) string {
	funcvar := GetFunctionName(GetUserAvatar)
	Debugging(funcvar, "In", username)

	t := regexp.MustCompile("[[:^ascii:]]").ReplaceAllLiteralString(username, "")
	resp, err := http.Get("https://mobile.twitter.com/" + t)
	if err != nil {
		log.Error(err)
	}

	defer resp.Body.Close()
	bit, err := ioutil.ReadAll(resp.Body)
	BruhMoment(err, "", false)

	var avatar string
	re := regexp.MustCompile(`(?ms)avatar.*?(http.*?)"`)
	if len(re.FindStringIndex(string(bit))) > 0 {
		re2 := regexp.MustCompile(`<img[^>]+\bsrc=["']([^"']+)["']`)
		submatchall := re2.FindAllStringSubmatch(re.FindString(string(bit)), -1)
		for _, element := range submatchall {
			avatar = strings.Replace(element[1], "normal.jpg", "400x400.jpg", -1)
		}
	}
	Debugging(funcvar, "In", avatar)
	return avatar
}

//Get bilibili user avatar
func (Data Dynamic_svr) GetUserAvatar() string {
	return Data.Data.Card.Desc.UserProfile.Info.Face
}

//Guild join handler
/*
func GuildJoin(s *discordgo.Session, g *discordgo.GuildCreate) {
	for _, Channel := range g.Guild.Channels {
		fmt.Println(Channel.Name
			BotPermission, err := s.UserChannelPermissions(BotID, Channel.ID)
			if err != nil {
				log.Error(err)
			}
			if Channel.Type == 0 && BotPermission&2048 != 0 {
				s.ChannelMessageSend(Channel.ID, "Thx for invite me to this channel <3 ")
				s.ChannelMessageSend(Channel.ID, "Type `"+config.PGeneral+"help` to show options")
				return
			}
	}
}

*/