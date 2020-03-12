package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "rates"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "database",
			Value:  "root@/rates?parseTime=true",
			Usage:  "MySQL Connection string",
			EnvVar: "RATES_DATABASE",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:   "collect",
			Usage:  "collect rates from source",
			Action: collect,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "source",
					Value:  "https://www.bank.lv/vk/ecb_rss.xml",
					Usage:  "URL to RSS Feed for collector to process",
					EnvVar: "RATES_SOURCE",
				},
			},
		},
		{
			Name:   "serve",
			Usage:  "start serving API",
			Action: serve,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "listen",
					Value:  ":3333",
					Usage:  "Address to bind/listen to",
					EnvVar: "RATES_LISTEN",
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
