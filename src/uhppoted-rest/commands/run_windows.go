package commands

import (
	"context"
	"flag"
	"fmt"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/eventlog"
	"log"
	"os"
	"path/filepath"
	"sync"
	"syscall"
	"uhppote"
	"uhppoted-rest/config"
	filelogger "uhppoted-rest/eventlog"
)

type Run struct {
	configuration string
	dir           string
	pidFile       string
	logFile       string
	logFileSize   int
	console       bool
	debug         bool
}

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
	pidFile:       filepath.Join(workdir(), "uhppoted-rest.pid"),
	logFile:       filepath.Join(workdir(), "logs", "uhppoted-rest.log"),
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
	log.Printf("uhppoted-rest daemon %s - %s (PID %d)\n", uhppote.VERSION, "Microsoft Windows", os.Getpid())

	f := func(c *config.Config) error {
		return r.start(c)
	}

	return r.execute(ctx, f)
}

func (r *Run) start(c *config.Config) error {
	var logger *log.Logger

	eventlogger, err := eventlog.Open("uhppoted-rest")
	if err != nil {
		events := filelogger.Ticker{Filename: r.logFile, MaxSize: r.logFileSize}
		logger = log.New(&events, "", log.Ldate|log.Ltime|log.LUTC)
	} else {
		defer eventlogger.Close()

		events := EventLog{eventlogger}
		logger = log.New(&events, "uhppoted-rest", log.Ldate|log.Ltime|log.LUTC)
	}

	logger.Printf("uhppoted-rest service - start\n")

	if r.console {
		r.run(c, logger)
		return nil
	}

	uhppoted := service{
		name:   "uhppoted-rest",
		conf:   c,
		logger: logger,
		cmd:    r,
	}

	logger.Printf("uhppoted-rest service - starting\n")
	err = svc.Run("uhppoted-rest", &uhppoted)

	if err != nil {
		fmt.Printf("   Unable to execute ServiceManager.Run request (%v)\n", err)
		fmt.Println()
		fmt.Println("   To run uhppoted-rest as a command line application, type:")
		fmt.Println()
		fmt.Println("     > uhppoted-rest --console")
		fmt.Println()

		logger.Fatalf("Error executing ServiceManager.Run request: %v", err)
		return err
	}

	logger.Printf("uhppoted-rest daemon - started\n")
	return nil
}

func (s *service) Execute(args []string, r <-chan svc.ChangeRequest, status chan<- svc.Status) (ssec bool, errno uint32) {
	s.logger.Printf("uhppoted-rest service - Execute\n")

	const commands = svc.AcceptStop | svc.AcceptShutdown

	status <- svc.Status{State: svc.StartPending}

	interrupt := make(chan os.Signal, 1)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			err := s.cmd.listen(s.conf, s.logger, interrupt)

			if err != nil {
				s.logger.Printf("ERROR: %v", err)
				continue
			}

			s.logger.Printf("exit\n")
			break
		}
	}()

	status <- svc.Status{State: svc.Running, Accepts: commands}

loop:
	for {
		select {
		case c := <-r:
			s.logger.Printf("uhppoted-rest service - select: %v  %v\n", c.Cmd, c.CurrentStatus)
			switch c.Cmd {
			case svc.Interrogate:
				s.logger.Printf("uhppoted-rest service - svc.Interrogate %v\n", c.CurrentStatus)
				status <- c.CurrentStatus

			case svc.Stop:
				interrupt <- syscall.SIGINT
				s.logger.Printf("uhppoted-rest service- svc.Stop\n")
				break loop

			case svc.Shutdown:
				interrupt <- syscall.SIGTERM
				s.logger.Printf("uhppoted-rest service - svc.Shutdown\n")
				break loop

			default:
				s.logger.Printf("uhppoted-rest service - svc.????? (%v)\n", c.Cmd)
			}
		}
	}

	s.logger.Printf("uhppoted-rest service - stopping\n")
	status <- svc.Status{State: svc.StopPending}
	wg.Wait()
	status <- svc.Status{State: svc.Stopped}
	s.logger.Printf("uhppoted-rest service - stopped\n")

	return false, 0
}

func (e *EventLog) Write(p []byte) (int, error) {
	err := e.log.Info(1, string(p))
	if err != nil {
		return 0, err
	}

	return len(p), nil
}
