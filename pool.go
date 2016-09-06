package gost

import (
	"sync"

	"os/exec"
	"sync/atomic"

	"strconv"

	"time"

	"github.com/mrsinham/go-selenium"
	"github.com/phayes/freeport"
)

type Pool interface {
	Get() (Worker, error)
	Put(w Worker)
	GetCurrentNb() int32
	GetMaxNb() int32
}

type Worker interface {
	GetDriver() selenium.WebDriver
	GetPort() int
	GetError() error
}

type pool struct {
	sync.Pool
	currentNb int32
	maxNb     int32
}

func (p *pool) Get() (Worker, error) {
	wd := p.Pool.Get().(Worker)
	if wd.GetError() != nil {
		return nil, wd.GetError()
	}

	// decrement currentNb
	atomic.AddInt32(&p.currentNb, ^int32(0))
	return wd, nil
}

func (p *pool) Put(w Worker) {
	p.Pool.Put(w)
	atomic.AddInt32(&p.currentNb, 1)
}

func (p *pool) GetMaxNb() int32 {
	return p.maxNb
}

func (p *pool) GetCurrentNb() int32 {
	return atomic.LoadInt32(&p.currentNb)
}

type worker struct {
	cmd       *exec.Cmd
	port      int
	err       error
	webdriver selenium.WebDriver
}

func (w *worker) GetDriver() selenium.WebDriver {
	return w.webdriver
}

func (w *worker) GetPort() int {
	return w.port
}

func (w *worker) GetError() error {
	return w.err
}

func newPool(maxInstance int32, pathToPJS string) *pool {
	p := &pool{maxNb: maxInstance}
	p.Pool = sync.Pool{
		New: func() interface{} {
			for {
				if atomic.LoadInt32(&p.currentNb) < maxInstance {
					break
				}
				// TODO: wait for loop time
			}

			// TODO: lock to avoid port collision

			fp := freeport.GetPort()
			fps := strconv.Itoa(fp)
			c := exec.Command(pathToPJS, []string{
				"--webdriver=" + fps,
			}...)

			err := c.Start()
			if err != nil {
				// TODO: insert logger
				return &worker{err: err}
			}

			<-time.After(500 * time.Millisecond)
			var wd selenium.WebDriver
			wd, err = selenium.NewRemote(make(selenium.Capabilities), "http://127.0.0.1:"+fps)
			if err != nil {
				return &worker{err: err}
			}

			// TODO: increase number

			return &worker{
				cmd:       c,
				port:      fp,
				webdriver: wd,
			}
		},
	}
	return p
}
