module github.com/JustHumanz/Go-simp

go 1.15

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/bwmarrin/discordgo v0.22.0
	github.com/cenkalti/dominantcolor v0.0.0-20171020061837-df772e8dd39e
	github.com/go-sql-driver/mysql v1.5.0
	github.com/hako/durafmt v0.0.0-20200710122514-c0fb7b4da026
	github.com/lib/pq v1.9.0
	github.com/n0madic/twitter-scraper v0.0.0-20201109091442-a5bcf09531d2
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646 // indirect
	github.com/olekukonko/tablewriter v0.0.4
	github.com/sirupsen/logrus v1.7.0
	gopkg.in/robfig/cron.v2 v2.0.0-20150107220207-be2e0b0deed5
)

replace github.com/n0madic/twitter-scraper v0.0.0-20201109091442-a5bcf09531d2 => github.com/JustHumanz/twitter-scraper v0.0.0-20201203144216-f27c756e93d2
