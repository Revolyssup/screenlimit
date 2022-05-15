package sysstats

import (
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/Revolyssup/screenlimit/db"
	"github.com/Revolyssup/screenlimit/db/events"
	"github.com/Revolyssup/screenlimit/policy"
)

type StatCollector struct {
	programs    *[]string //eg- brave, chrome, firefox or any other app we want to monitor.
	actionOnApp map[string]policy.ActionType
	store       *db.EventStore
	buf         string
	mx          sync.Mutex
}

func New(prog *[]string, store *db.EventStore) *StatCollector {
	return &StatCollector{
		store:       store,
		programs:    prog,
		actionOnApp: make(map[string]policy.ActionType),
	}
}
func (s *StatCollector) AddActionApp(app string, action string) {
	s.mx.Lock()
	defer s.mx.Unlock()
	fmt.Println("adding action on app ", app)
	(s.actionOnApp)[app] = policy.ActionType(action)
}
func (s *StatCollector) AddApp(app string) {
	s.mx.Lock()
	defer s.mx.Unlock()
	fmt.Println("adding app ", app)
	*(s.programs) = append(*(s.programs), app)
}
func (s *StatCollector) GetApp() []string {
	s.mx.Lock()
	defer s.mx.Unlock()
	return *s.programs
}

// func (s *StatCollector) Log(ts string, ev string, eventdesciption string) error {
// 	return s.DB.Model(&db.Event{}).Create(&db.Event{
// 		Timestamp:  ts,
// 		Action:     eventdesciption,
// 		ActionType: actions.Type(ev),
// 	}).Error
// }

//TODO: It will run in the background extracting statistics and pushing them into the database
func (s *StatCollector) Run() {
	cmd := exec.Command("bpftrace", "-e", "t:syscalls:sys_enter_execve {printf(\"pid: %d,--comm: %s\\n\",pid,comm);}")
	var stdout, stderr []byte
	var errStdout, errStderr error
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()
	err := cmd.Start()
	if err != nil {
		log.Fatalf("cmd.Start() failed with '%s'\n", err)
	}
	fmt.Println("started bpftrace")

	// cmd.Wait() should be called only after we finish reading
	// from stdoutIn and stderrIn.
	// wg ensures that we finish
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		stdout, errStdout = copyAndCapture(s, stdoutIn)
		wg.Done()
	}()

	stderr, errStderr = copyAndCapture(s, stderrIn)

	wg.Wait()

	err = cmd.Wait()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	if errStdout != nil || errStderr != nil {
		log.Fatal("failed to capture stdout or stderr\n")
	}
	outStr, errStr := string(stdout), string(stderr)
	fmt.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)

}
func copyAndCapture(w io.Writer, r io.Reader) ([]byte, error) {
	var out []byte
	buf := make([]byte, 1024, 1024)
	for {
		n, err := r.Read(buf[:])
		if n > 0 {
			d := buf[:n]
			out = append(out, d...)
			_, err := w.Write(d)
			if err != nil {
				return out, err
			}
		}
		if err != nil {
			// Read returns io.EOF at the end of file, which is not an error for us
			if err == io.EOF {
				err = nil
			}
			return out, err
		}
	}
}

func (m *StatCollector) Write(p []byte) (n int, err error) {
	buf := string(p)
	for _, line := range strings.Split(buf, "\n") {
		for _, app := range *m.programs {
			if strings.Contains(line, app) {
				pid := strings.TrimSuffix(strings.TrimPrefix(line, "pid: "), ",")
				fmt.Println("pid is ", pid)
				m.store.Add(time.Now().GoString(), "Child opened "+app, events.Child, pid)
				fmt.Println("LOFE ", (m.actionOnApp)[app])
				(m.actionOnApp)[app].Exec()
			}
		}
	}

	return len(p), nil
}
