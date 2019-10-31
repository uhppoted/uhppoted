package commands

import (
	"flag"
	"fmt"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/eventlog"
	"log"
	"os"
	"path/filepath"
	"sync"
	"syscall"
	"uhppoted/config"
	filelogger "uhppoted/eventlog"
)

type service struct {
	name   string
	conf   *config.Config
	logger *log.Logger
}

type EventLog struct {
	log *eventlog.Log
}

var wd = workdir()
var configuration = flag.String("config", filepath.Join(wd, "uhppoted.conf"), "Path for the configuration file")
var dir = flag.String("dir", wd, "Working directory")
var logfile = flag.String("logfile", filepath.Join(wd, "logs", "uhppoted.log"), "uhppoted log file")
var logfilesize = flag.Int("logfilesize", 10, "uhppoted log file size")
var pidFile = flag.String("pid", filepath.Join(wd, "uhppoted.pid"), "uhppoted PID file")
var console = flag.Bool("console", false, "Run as command-line application")

func (c *Run) Execute(ctx Context) error {
	log.Printf("uhppoted service - %s (PID %d)\n", "Microsoft Windows", os.Getpid())

	return execute(ctx)
}

func (s *service) Execute(args []string, r <-chan svc.ChangeRequest, status chan<- svc.Status) (ssec bool, errno uint32) {
	s.logger.Printf("uhppoted service - Execute\n")

	const commands = svc.AcceptStop | svc.AcceptShutdown

	status <- svc.Status{State: svc.StartPending}

	interrupt := make(chan os.Signal, 1)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			err := listen(s.conf, s.logger, interrupt)

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
			s.logger.Printf("uhppoted service - select: %v  %v\n", c.Cmd, c.CurrentStatus)
			switch c.Cmd {
			case svc.Interrogate:
				s.logger.Printf("uhppoted service - svc.Interrogate %v\n", c.CurrentStatus)
				status <- c.CurrentStatus

			case svc.Stop:
				interrupt <- syscall.SIGINT
				s.logger.Printf("uhppoted service- svc.Stop\n")
				break loop

			case svc.Shutdown:
				interrupt <- syscall.SIGTERM
				s.logger.Printf("uhppoted service - svc.Shutdown\n")
				break loop

			default:
				s.logger.Printf("uhppoted service - svc.????? (%v)\n", c.Cmd)
			}
		}
	}

	s.logger.Printf("uhppoted service - stopping\n")
	status <- svc.Status{State: svc.StopPending}
	wg.Wait()
	status <- svc.Status{State: svc.Stopped}
	s.logger.Printf("uhppoted service - stopped\n")

	return false, 0
}

func start(c *config.Config, logfile string, logfilesize int) {
	var logger *log.Logger

	eventlogger, err := eventlog.Open("uhppoted")
	if err != nil {
		events := filelogger.Ticker{Filename: logfile, MaxSize: logfilesize}
		logger = log.New(&events, "", log.Ldate|log.Ltime|log.LUTC)
	} else {
		defer eventlogger.Close()

		events := EventLog{eventlogger}
		logger = log.New(&events, "uhppoted", log.Ldate|log.Ltime|log.LUTC)
	}

	logger.Printf("uhppoted service - start\n")

	if *console {
		run(c, logger)
		return
	}

	uhppoted := service{
		name:   "uhppoted",
		conf:   c,
		logger: logger,
	}

	logger.Printf("uhppoted service - starting\n")
	err = svc.Run("uhppoted", &uhppoted)

	if err != nil {
		fmt.Printf("   Unable to execute ServiceManager.Run request (%v)\n", err)
		fmt.Println()
		fmt.Println("   To run uhppoted as a command line application, type:")
		fmt.Println()
		fmt.Println("     > uhppoted --console")
		fmt.Println()

		logger.Fatalf("Error executing ServiceManager.Run request: %v", err)
		return
	}

	logger.Printf("uhppoted daemon - started\n")
}

func (e *EventLog) Write(p []byte) (int, error) {
	err := e.log.Info(1, string(p))
	if err != nil {
		return 0, err
	}

	return len(p), nil
}
