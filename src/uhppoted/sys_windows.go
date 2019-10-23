package main

import (
	"flag"
	"fmt"
	"golang.org/x/sys/windows/svc"
	"log"
	"os"
	"path/filepath"
	"sync"
	"syscall"
	"uhppoted/config"
	"uhppoted/eventlog"
)

// WINDOWS (PROVISIONAL - NOT TESTED)

type service struct {
	conf   *config.Config
	logger *log.Logger
}

// var pwd, _ = filepath.Abs(filepath.Dir(os.Args[0]))
var pwd = `C:\uhppoted`
var configuration = flag.String("config", filepath.Join(pwd, "uhppoted.cfg"), "Path for the configuration file")
var dir = flag.String("dir", pwd, "Working directory")
var logfile = flag.String("logfile", filepath.Join(pwd, "logs", "uhppoted.log"), "uhppoted log file")
var logfilesize = flag.Int("logfilesize", 10, "uhppoted log file size")
var pidFile = flag.String("pid", filepath.Join(pwd, "uhppoted.pid"), "uhppoted PID file")

func sysinit() {
	log.Printf("uhppoted daemon - %s (PID %d)\n", "Microsoft Windows", os.Getpid())
	fmt.Printf("uhppoted daemon - %s (PID %d)\n", "Microsoft Windows", os.Getpid())
}

func start(c *config.Config, logfile string, logfilesize int) {
	events := eventlog.Ticker{Filename: logfile, MaxSize: logfilesize}
	logger := log.New(&events, "", log.Ldate|log.Ltime|log.LUTC)

	logger.Printf("uhppoted daemon - start\n")

	interactive, err := svc.IsAnInteractiveSession()
	logger.Printf("uhppoted daemon - interactive: %v %v\n", interactive, err)
	// if err != nil {
	// 	logger.Printf("Error querying ServiceManager for interactive session status: %v", err)
	// 	log.Fatalf("Error querying ServiceManager for interactive session status: %v", err)
	// }

	// if interactive {
	// 	run(c, logfile, logfilesize)
	// 	return
	// }

	logger.Printf("uhppoted daemon - starting\n")
	err = svc.Run("uhppoted", &service{
		conf:   c,
		logger: logger,
	})
	logger.Printf("uhppoted daemon - started %v\n", err)

	if err != nil {
		logger.Printf("Error executing ServiceManager.Run request: %v", err)
		log.Fatalf("Error executing ServiceManager.Run request: %v", err)
	}
}

func (s *service) Execute(args []string, r <-chan svc.ChangeRequest, status chan<- svc.Status) (ssec bool, errno uint32) {
	s.logger.Printf("uhppoted daemon - Execute\n")

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
			s.logger.Printf("uhppoted daemon - select: %v  %v\n", c.Cmd, c.CurrentStatus)
			switch c.Cmd {
			case svc.Interrogate:
				s.logger.Printf("uhppoted daemon - svc.Interrogate %v\n", c.CurrentStatus)
				status <- c.CurrentStatus

			case svc.Stop:
				interrupt <- syscall.SIGINT
				s.logger.Printf("uhppoted daemon - svc.Stop\n")
				break loop

			case svc.Shutdown:
				interrupt <- syscall.SIGTERM
				s.logger.Printf("uhppoted daemon - svc.Shutdown\n")
				break loop

			default:
				s.logger.Printf("uhppoted daemon - svc.????? (%v)\n", c.Cmd)
			}
		}
	}

	s.logger.Printf("uhppoted daemon - stopping\n")
	status <- svc.Status{State: svc.StopPending}
	wg.Wait()
	status <- svc.Status{State: svc.Stopped}
	s.logger.Printf("uhppoted daemon - stopped\n")

	return false, 0
}

// func runService(name string, isDebug bool) {
// 	// var err error
// 	// if isDebug {
// 	// 	elog = debug.New(name)
// 	// } else {
// 	// 	elog, err = eventlog.Open(name)
// 	// 	if err != nil {
// 	// 		return
// 	// 	}
// 	// }
// 	// defer elog.Close()
//     //
// 	// elog.Info(1, fmt.Sprintf("starting %s service", name))
// 	// run := svc.Run
// 	// if isDebug {
// 	// 	run = debug.Run
// 	// }
// 	// err = run(name, &myservice{})
// 	// if err != nil {
// 	// 	elog.Error(1, fmt.Sprintf("%s service failed: %v", name, err))
// 	// 	return
// 	// }
// 	// elog.Info(1, fmt.Sprintf("%s service stopped", name))
// }
