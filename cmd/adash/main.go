// Management Console
package main

import (
	"context"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
	"time"

	adash "bitbucket.org/peregrinetraders/mc/apps/adash/service"
	"bitbucket.org/peregrinetraders/mc/apps/adash/service/charter"
	"bitbucket.org/peregrinetraders/mc/apps/adash/service/holidayer"
	"bitbucket.org/peregrinetraders/mc/apps/adash/service/redishashtracker"
	"bitbucket.org/peregrinetraders/mc/apps/adash/service/redislisttracker"
	"bitbucket.org/peregrinetraders/mc/pkg/aaa"
	"bitbucket.org/peregrinetraders/mc/pkg/redistracker"

	"bitbucket.org/peregrinetraders/mc/pkg/store"
	//store "github.com/PeregrineTradersDevTeam/redis-farm"
	"github.com/axkit/hms"
	"github.com/mediocregopher/radix/v3"
	"github.com/regorov/go-chart/drawing"
	"github.com/rs/zerolog"
	"github.com/urfave/cli"
)

// EnvVarPrefix holds environment variables prefix related to applications.
const EnvVarPrefix = "PT_AD_"

func main() {

	f, err := os.Create("./cpu.pprof")
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	app := cli.NewApp()
	app.Name = "Peregrine Traders AD"
	app.Usage = "Peregrine Traders Algo Dashboard"
	app.Version = ReleaseNumber + " (" + BuildTime + ")"

	app.Commands = []cli.Command{
		{
			Name:    "start",
			Aliases: []string{"s"},
			Usage:   "start microservice",
			Action:  start,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "listen, l",
					Value:  ":8080",
					Usage:  "HTTP listener",
					EnvVar: EnvVarPrefix + "HTTP_LISTENER",
				},

				cli.BoolFlag{
					Name:   "debug, d",
					Usage:  "debug mode activation flag",
					EnvVar: EnvVarPrefix + "DEBUG",
				},

				cli.StringFlag{
					Name:   "runtime-listener, rl",
					Usage:  "golang runtime HTTP listener",
					EnvVar: EnvVarPrefix + "RUNTIME_HTTP_LISTENER",
				},
				cli.IntFlag{
					Name:   "db",
					Value:  0,
					Usage:  "redis database number",
					EnvVar: EnvVarPrefix + "DB",
				},
				cli.StringFlag{
					Name:   "redis",
					Value:  "redis:6379",
					Usage:  "redis server listener as ip:port",
					EnvVar: EnvVarPrefix + "REDIS_LISTENER",
				},
				cli.StringFlag{
					Name:   "frontfolder, ff",
					Value:  "./frontend",
					Usage:  "html, css and js location",
					EnvVar: EnvVarPrefix + "FRONTEND_FOLDER",
				},
				cli.StringFlag{
					Name:   "lcfolder, lcf",
					Value:  "",
					Usage:  "launcher config files location",
					EnvVar: EnvVarPrefix + "LAUNCHER_CONFIG_FOLDER",
				},
				cli.StringFlag{
					Name:   "restart-allowed-before, rab",
					Value:  "",
					Usage:  "platform restart allowed before time in UTC",
					EnvVar: EnvVarPrefix + "RESTART_ALLOWED_BEFORE",
				},
				cli.StringFlag{
					Name:   "platform-restart-command, prc",
					Value:  "",
					Usage:  "platform restart shell command",
					EnvVar: EnvVarPrefix + "PLATFORM_RESTART_COMMMAND",
				},
				cli.StringFlag{
					Name:   "market-holidays, mh",
					Value:  "market-holidays.csv",
					Usage:  "market holidays calendar",
					EnvVar: EnvVarPrefix + "MARKET_HOLIDAY_FILE",
				},
			},
		},
		{
			Name:   "stop",
			Usage:  "stop microservice",
			Action: stop,
		},
	}

	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}

func stop(c *cli.Context) error {

	return nil
}

func start(c *cli.Context) error {

	// 1. logger configuration
	zerolog.TimeFieldFormat = "20060102T150405.999Z07:00"
	zerolog.TimestampFieldName = "t"
	zerolog.MessageFieldName = "msg"
	zerolog.LevelFieldName = "lvl"

	debug := c.IsSet("debug")

	lvl := zerolog.InfoLevel
	if debug {
		lvl = zerolog.DebugLevel
	}

	logger := zerolog.New(os.Stderr).Level(lvl).With().Timestamp().Logger()
	logger.Info().Str("version", BuildNumber).Msg("application started")

	// 2. launch runtime monitor.
	if c.IsSet("runtime-listener") {
		go func() {
			logger.Info().Str("runtime-listener", c.String("runtime-listener")).Msg("start runtime http listener")
			if err := http.ListenAndServe(c.String("runtime-listener"), nil); err != nil {
				logger.Info().Str("runtime-listener", c.String("runtime-listener")).Msg("start runtime http listener")
				logger.Error().Str("errmsg", err.Error()).Msg("runtime listener starting failed")
			}
		}()
	}

	// 3. create SIGINT capture.
	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGINT)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-stop
		logger.Info().Msg("Signal SIGINT received")
		cancel()
	}()

	// 4. write launching params to the log.
	info := func() {
		logger.Info().
			Str("version", ReleaseNumber).
			Str("build-time", BuildTime).
			Bool("debug", debug).
			Str("redis-address", c.String("redis")).
			Int("redis-db-num", c.Int("db")).
			Str("http-listener", c.String("listen")).
			Str("restart-allowed-before", c.String("restart-allowed-before")).
			Str("platform-restart-command", c.String("platform-restart-command")).
			Msg("launching with params")
	}

	info()

	// 5. connecting to redis.
	con := store.NewRedisStore(c.String("redis"), c.Int("db"), &logger)
	if err := con.Connect(); err != nil {
		logger.Error().Str("errmsg", err.Error()).Msg("connection to REDIS failed")
		return err
	}

	// 6. create charter engine.
	chr := charter.NewCharter(con, &logger)

	// 7.
	redisTracker := redistracker.New(&logger, con)

	astore := aaa.FileStore{}
	err := astore.Load("./auth.json")
	if err != nil {
		logger.Error().Str("errmsg", err.Error()).Msg("security params reading failed")
		return err
	}

	customConnFunc := func(network, addr string) (radix.Conn, error) {
		return radix.Dial(network, addr, radix.DialSelectDB(c.Int("db")))
	}

	// this pool will use our ConnFunc for all connections it creates
	pool, err := radix.NewPool("tcp", c.String("redis"), 10, radix.PoolConnFunc(customConnFunc))
	if err != nil {
		logger.Error().Str("errmsg", err.Error()).Msg("connection to REDIS failed")
		return err
	}
	defer pool.Close()

	logger.Info().Msg("rest api server starting")

	authservice := aaa.New(&astore)

	pullstore := adash.NewRedisStore(pool, c.Int("db"))

	service := adash.NewService(pullstore, &logger, c.String("lcfolder"))

	redisTracker.AddKey("CFG:SESSION", 5*time.Second, false, chr.Reset, service.SetSession)
	redisTracker.Start(ctx)

	err = service.Start(ctx)
	if err != nil {
		logger.Error().Str("errmsg", err.Error()).Msg("service core starting failed")
		return err
	}

	av := adash.AppVersion{
		Version:  ReleaseNumber,
		BuildAt:  BuildTime,
		BuildFor: BuildPlatform,
	}

	if av.Version == "" {
		av.Version = "N/A"
	}

	buf, err := ioutil.ReadFile("release.md")
	if err == nil {
		av.ReleaseNotes = string(buf)
	} else {
		logger.Warn().Str("msg", err.Error()).Msg("could not read file release.md")
	}

	h := holidayer.New(logger, c.String("market-holidays"))
	if err := h.Start(ctx); err != nil {
		logger.Error().Str("errmsg", err.Error()).Msg("service core starting failed")
	}

	api := adash.NewAPI(service,
		authservice,
		logger,
		c.String("frontfolder"),
		hms.HHMM(c.String("restart-allowed-before")),
		c.String("platform-restart-command"),
		av,
		h,
	)

	//ss.SetWsProxy(api.WsProxy())
	chr.RegisterCurve(charter.NewCurve("PNL:BUOYANCY",
		"",
		charter.RedisPlatformPnLFiller(&logger, con, "PNL:BUOYANCY", chr),
		charter.StyleMixedPnL,
		charter.PointsAsIs,
		charter.LabelMaxMin),
	)

	lt := redislisttracker.New(con, &logger, chr)
	lt.TrackList("PNL:BU",
		charter.StyleMixedPnL,
		charter.PointsAsIs,
		charter.LabelMaxMin)

	rht := redishashtracker.New(con, &logger, chr)

	rht.TrackHashAttribute("M", "pd", "LASTPD",
		charter.StyleInstancePnL,
		nil,
		charter.Label(-2, -1, charter.LabelMinMaxInterval, drawing.ColorWhite))

	rht.TrackHashAttribute("M", "pl", "LASTPL",
		charter.StyleInstancePnL,
		nil,
		charter.Label(-2, -1, charter.LabelMinMaxInterval, drawing.ColorWhite))

	rht.TrackHashAttribute("M", "pnl", "LASTPNL",
		charter.StyleInstancePnL,
		nil,
		charter.Label(-2, -1, charter.LabelMinMaxInterval, drawing.ColorWhite))

	rht.TrackHashAttribute("M", "l", "LASTPRICE",
		charter.StyleLastPrice,
		charter.PointsAsIs,
		charter.Label(-2, 10, charter.Last, drawing.ColorWhite)) // drawing.Color{R: 25, G: 0, B: 255, A: 255}))

	rht.TrackHashAttribute("M", "sahp", "SAHPRICE",
		charter.StyleLastPrice,
		charter.PointsAsIs,
		charter.Label(-2, 10, charter.Last, drawing.ColorWhite)) // drawing.Color{R: 25, G: 0, B: 255, A: 255}))

	rht.TrackHashAttribute("M", "hi", "HIPRICE",
		charter.StyleThreshold(drawing.ColorFromAlphaMixedRGBA(0, 0, 255, 255)),
		charter.PointsSquared,
		charter.Label(4, 10, charter.Last, drawing.ColorFromAlphaMixedRGBA(0, 249, 255, 255)))

	rht.TrackHashAttribute("M", "lo", "LOPRICE",
		charter.StyleThreshold(drawing.ColorRed), // drawing.ColorFromAlphaMixedRGBA(255, 231, 0, 255)),
		charter.PointsSquared,
		charter.Label(4, -1, charter.Last, drawing.ColorFromAlphaMixedRGBA(0, 249, 255, 255))) //drawing.ColorFromAlphaMixedRGBA(255, 231, 0, 255)))

	rht.TrackHashAttribute("M", "op", "OPPRICE",
		charter.StyleOpeningPrice(drawing.ColorFromAlphaMixedRGBA(253, 254, 2, 255)),
		charter.PointsSquared,
		charter.Label(4, 37, charter.Last, drawing.ColorFromAlphaMixedRGBA(253, 254, 2, 255)))

	rht.Start(ctx)
	lt.Start(ctx)

	//time.Sleep(time.Second * 5)

	// fo, err := os.Create("twolines.png")
	// if err != nil {
	// 	return err
	// }
	// defer fo.Close()

	// buf := bytes.NewBuffer(nil)

	// if err := charter.RenderMulti([]*charter.Curve{chr.Curve("M:A002.l"), chr.Curve("M:A002.hi"), chr.Curve("M:A002.lo"), chr.Curve("M:A002.op")}, 184, 67, buf); err != nil { //,
	// 	return err
	// }

	// if _, err := fo.Write(buf.Bytes()); err != nil {
	// 	return err
	// }
	//return nil

	if err := api.Start(ctx, c.String("listen"), rht, chr); err != nil {
		logger.Error().Str("errmsg", err.Error()).Msg("api server starting failed")
		return err
	}
	logger.Debug().Msg("API server started!")

	<-ctx.Done()

	return nil
}
