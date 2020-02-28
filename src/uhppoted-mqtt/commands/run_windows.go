package commands

import (
	"context"
	"flag"
	"fmt"
	"github.com/uhppoted/uhppote-core/uhppote"
	"github.com/uhppoted/uhppoted-api/config"
	filelogger "github.com/uhppoted/uhppoted-api/eventlog"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/eventlog"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
)

type service struct {
	name   string
	conf   *config.Config
	logger *log.Logger
	cmd    *Run
}

type EventLog struct {
	log *eventlog.Log
}

var RUN = Run{
	configuration: filepath.Join(workdir(), "uhppoted.conf"),
	dir:           workdir(),
	pidFile:       filepath.Join(workdir(), fmt.Sprintf("%s.pid", SERVICE)),
	logFile:       filepath.Join(workdir(), "logs", fmt.Sprintf("%s.log", SERVICE)),
	logFileSize:   10,
	console:       false,
	debug:         false,
}

func (r *Run) FlagSet() *flag.FlagSet {
	flagset := flag.NewFlagSet("", flag.ExitOnError)

	flagset.StringVar(&r.configuration, "config", r.configuration, "Sets the configuration file path")
	flagset.StringVar(&r.dir, "dir", r.dir, "Work directory")
	flagset.StringVar(&r.pidFile, "pid", r.pidFile, "Sets the service PID file path")
	flagset.StringVar(&r.logFile, "logfile", r.logFile, "Sets the log file path")
	flagset.IntVar(&r.logFileSize, "logfilesize", r.logFileSize, "Sets the log file size before forcing a log rotate")
	flagset.BoolVar(&r.console, "console", r.console, "Run as command-line application")
	flagset.BoolVar(&r.debug, "debug", r.debug, "Displays vaguely useful internal information")

	return flagset
}

func (r *Run) Execute(ctx context.Context) error {
	log.Printf("%s service %s - %s (PID %d)\n", SERVICE, uhppote.VERSION, "Microsoft Windows", os.Getpid())

	f := func(c *config.Config) error {
		return r.start(c)
	}

	return r.execute(ctx, f)
}

func (r *Run) start(c *config.Config) error {
	var logger *log.Logger

	eventlogger, err := eventlog.Open(SERVICE)
	if err != nil {
		events := filelogger.Ticker{Filename: r.logFile, MaxSize: r.logFileSize}
		logger = log.New(&events, "", log.Ldate|log.Ltime|log.LUTC)
	} else {
		defer eventlogger.Close()

		events := EventLog{eventlogger}
		logger = log.New(&events, SERVICE, log.Ldate|log.Ltime|log.LUTC)
	}

	logger.Printf("%s service - start\n", SERVICE)

	if r.console {
		interrupt := make(chan os.Signal, 1)

		signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
		r.run(c, logger, interrupt)
		return nil
	}

	uhppoted := service{
		name:   SERVICE,
		conf:   c,
		logger: logger,
		cmd:    r,
	}

	logger.Printf("%s service - starting\n", SERVICE)
	err = svc.Run(SERVICE, &uhppoted)
	if err != nil {
		fmt.Printf("   Unable to execute ServiceManager.Run request (%v)\n", err)
		fmt.Println()
		fmt.Printf("   To run %s as a command line application, type:\n", SERVICE)
		fmt.Println()
		fmt.Printf("     > %s --console\n", SERVICE)
		fmt.Println()

		logger.Fatalf("Error executing ServiceManager.Run request: %v", err)
		return err
	}

	logger.Printf("%s daemon - started\n", SERVICE)
	return nil
}

func (s *service) Execute(args []string, r <-chan svc.ChangeRequest, status chan<- svc.Status) (ssec bool, errno uint32) {
	s.logger.Printf("%s service - Execute\n", SERVICE)

	const commands = svc.AcceptStop | svc.AcceptShutdown

	status <- svc.Status{State: svc.StartPending}

	interrupt := make(chan os.Signal, 1)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		s.cmd.run(s.conf, s.logger, interrupt)

		s.logger.Printf("exit\n")
	}()

	status <- svc.Status{State: svc.Running, Accepts: commands}

loop:
	for {
		select {
		case c := <-r:
			s.logger.Printf("%s service - select: %v  %v\n", SERVICE, c.Cmd, c.CurrentStatus)
			switch c.Cmd {
			case svc.Interrogate:
				s.logger.Printf("%s service - svc.Interrogate %v\n", SERVICE, c.CurrentStatus)
				status <- c.CurrentStatus

			case svc.Stop:
				interrupt <- syscall.SIGINT
				s.logger.Printf("%s service- svc.Stop\n", SERVICE)
				break loop

			case svc.Shutdown:
				interrupt <- syscall.SIGTERM
				s.logger.Printf("%s service - svc.Shutdown\n", SERVICE)
				break loop

			default:
				s.logger.Printf("%s service - svc.????? (%v)\n", SERVICE, c.Cmd)
			}
		}
	}

	s.logger.Printf("%s service - stopping\n", SERVICE)
	status <- svc.Status{State: svc.StopPending}
	wg.Wait()
	status <- svc.Status{State: svc.Stopped}
	s.logger.Printf("%s service - stopped\n", SERVICE)

	return false, 0
}

func (e *EventLog) Write(p []byte) (int, error) {
	err := e.log.Info(1, string(p))
	if err != nil {
		return 0, err
	}

	return len(p), nil
}
