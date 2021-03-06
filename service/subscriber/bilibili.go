package main

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"sync"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	engine "github.com/JustHumanz/Go-Simp/pkg/engine"
	network "github.com/JustHumanz/Go-Simp/pkg/network"
	pilot "github.com/JustHumanz/Go-Simp/service/pilot/grpc"
	log "github.com/sirupsen/logrus"
)

func CheckBiliBili() {
	BiliBiliSession := []string{"Cookie", "SESSDATA=" + configfile.BiliSess}
	for _, Group := range *Payload {
		Names := Group.Members
		for _, Name := range Names {
			if Name.BiliBiliID != 0 && Name.Active() {
				var (
					wg        sync.WaitGroup
					bilistate BiliBiliStat
				)
				wg.Add(3)
				go func() {
					body, curlerr := network.CoolerCurl("https://api.bilibili.com/x/relation/stat?vmid="+strconv.Itoa(Name.BiliBiliID), BiliBiliSession)
					if curlerr != nil {
						log.Error(curlerr)
						gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
							Message: curlerr.Error(),
							Service: ModuleState,
						})
					}
					err := json.Unmarshal(body, &bilistate.Follow)
					if err != nil {
						log.Error(err)
					}
					defer wg.Done()
				}()

				go func() {
					body, curlerr := network.CoolerCurl("https://api.bilibili.com/x/space/upstat?mid="+strconv.Itoa(Name.BiliBiliID), BiliBiliSession)
					if curlerr != nil {
						log.Error(curlerr)
						gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
							Message: curlerr.Error(),
							Service: ModuleState,
						})
					}
					err := json.Unmarshal(body, &bilistate.LikeView)
					if err != nil {
						log.Error(err)
					}
					defer wg.Done()
				}()

				go func() {
					baseurl := "https://api.bilibili.com/x/space/arc/search?mid=" + strconv.Itoa(Name.BiliBiliID) + "&ps=100"
					url := []string{baseurl + "&tid=1", baseurl + "&tid=3", baseurl + "&tid=4"}
					for f := 0; f < len(url); f++ {
						body, curlerr := network.CoolerCurl(url[f], BiliBiliSession)
						if curlerr != nil {
							log.Error(curlerr)
							gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
								Message: curlerr.Error(),
								Service: ModuleState,
							})
						}
						var video engine.SpaceVideo
						err := json.Unmarshal(body, &video)
						if err != nil {
							log.Error(err)
						}
						bilistate.Videos += video.Data.Page.Count
					}
					defer wg.Done()
				}()
				wg.Wait()

				BiliFollowDB, err := Name.GetSubsCount()
				if err != nil {
					log.Error(err)
				}
				if bilistate.Follow.Data.Follower != 0 {
					if BiliFollowDB.BiliFollow != bilistate.Follow.Data.Follower {
						if bilistate.Follow.Data.Follower <= 10000 {
							for i := 0; i < 1000001; i += 100000 {
								if i == bilistate.Follow.Data.Follower {
									Avatar := Name.BiliBiliAvatar
									Color, err := engine.GetColor(config.TmpDir, Avatar)
									if err != nil {
										log.Error(err)
									}

									err = Name.RemoveSubsCache()
									if err != nil {
										log.Error(err)
									}
									Graph := "[View as Graph](" + os.Getenv("PrometheusURL") + "/graph?g0.expr=get_subscriber%7Bstate%3D%22BiliBili%22%2C%20vtuber%3D%22" + Name.Name + "%22%7D&g0.tab=0&g0.stacked=0&g0.range_input=4w)"
									SendNude(engine.NewEmbed().
										SetAuthor(Group.GroupName, Group.IconURL, "https://space.bilibili.com/"+strconv.Itoa(Name.BiliBiliID)).
										SetTitle(engine.FixName(Name.EnName, Name.JpName)).
										SetThumbnail(config.BiliBiliIMG).
										SetDescription("Congratulation for "+engine.NearestThousandFormat(float64(i))+" followers").
										SetImage(Avatar).
										AddField("Viewers", strconv.Itoa(bilistate.LikeView.Data.Archive.View)).
										AddField("Videos", strconv.Itoa(bilistate.Videos)).
										SetURL("https://space.bilibili.com/"+strconv.Itoa(Name.BiliBiliID)).
										InlineAllFields().
										AddField("Graph", Graph).
										SetColor(Color).MessageEmbed, Group, Name)
								}
							}
						} else {
							for i := 0; i < 10001; i += 1000 {
								if i == bilistate.Follow.Data.Follower {
									Avatar := Name.BiliBiliAvatar
									Color, err := engine.GetColor(config.TmpDir, Avatar)
									if err != nil {
										log.Error(err)
									}
									SendNude(engine.NewEmbed().
										SetAuthor(Group.GroupName, Group.IconURL, "https://space.bilibili.com/"+strconv.Itoa(Name.BiliBiliID)).
										SetTitle(engine.FixName(Name.EnName, Name.JpName)).
										SetThumbnail(config.BiliBiliIMG).
										SetDescription("Congratulation for "+engine.NearestThousandFormat(float64(i))+" followers").
										SetImage(Avatar).
										AddField("Views", engine.NearestThousandFormat(float64(bilistate.LikeView.Data.Archive.View))).
										AddField("Videos", engine.NearestThousandFormat(float64(bilistate.Videos))).
										SetURL("https://space.bilibili.com/"+strconv.Itoa(Name.BiliBiliID)).
										InlineAllFields().
										SetColor(Color).MessageEmbed, Group, Name)
								}
							}
						}
					}
					log.WithFields(log.Fields{
						"Past BiliBili Follower":    BiliFollowDB.BiliFollow,
						"Current BiliBili Follower": bilistate.Follow.Data.Follower,
						"Vtuber":                    Name.EnName,
					}).Info("Update BiliBili Follower")

					BiliFollowDB.SetMember(Name).SetGroup(Group).
						UpBiliFollow(bilistate.Follow.Data.Follower).
						UpBiliVideo(bilistate.Videos).
						UpBiliViews(bilistate.LikeView.Data.Archive.View).
						UpdateState(config.BiliLive).UpdateSubs()

					bin, err := BiliFollowDB.MarshalBinary()
					if err != nil {
						log.Error(err)
					}
					if config.GoSimpConf.Metric {
						gRCPconn.MetricReport(context.Background(), &pilot.Metric{
							MetricData: bin,
							State:      config.SubsState,
						})
					}
				}
			}
		}
	}
}
