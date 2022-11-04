package controllers

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"server/config"
	"sync"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	once             sync.Once
	exitCtrlInstance *ExitController
)

// ExitController type to handle exit and cancellation
type ExitController struct {
	ctx           context.Context
	ctxCancelFunc func()
	ctxCancelDone chan bool
	handlers      map[string]func()
	//log           logging.Log
	mutex sync.RWMutex
}

// NewExitController returns a singleton instance of the controller, repeated calls return same instance
func NewExitController() *ExitController {
	once.Do(func() {
		timeout := time.Duration(config.GetEnvironment().ENGINE_MAX_RUNTIME) * time.Second
		log.Tracef("Setting up context with timeout of %s", timeout)
		ctx, ctxCancelFunc := context.WithTimeout(context.Background(), timeout)
		exitCtrlInstance = &ExitController{
			ctx:           ctx,
			ctxCancelFunc: ctxCancelFunc,
			ctxCancelDone: make(chan bool),
			handlers:      make(map[string]func(), 10),
		}
	})
	return exitCtrlInstance
}

// AddCancelHandler adds a function to be invoked when the context is cancelled
// Provided function will be invoked in its own go routine therefore its implementation
// need to wait for what ever needs to be done upon cancelling.
func AddCancelHandler(name string, handler func()) error {
	if exitCtrlInstance == nil {
		return errors.New("not initialized yet")
	}
	return exitCtrlInstance.addCancelHandler(name, handler)
}

// Done returns channel that's closed when all cancel handlers are done
func (c *ExitController) Done() <-chan bool {
	return c.ctxCancelDone
}

// Context returns cancellable context used by this exit controller
func (c *ExitController) Context() context.Context {
	return c.ctx
}

// Start starts the controller's routines waiting for either OS signal or context cancellation (cancellation only useful in testing)
func (c *ExitController) Start() {
	go c.waitForCancel() // this will start a method that invokes handlers when context gets cancelled
	go c.waitForSignal() // this will start a method that cancels the context on system signal
}

// Stop cancels the context - only useful for testing
func (c *ExitController) Stop() error {
	// check that the context hasn't been cancelled by getting error from the context
	if err := c.ctx.Err(); err != nil {
		return err
	}
	c.ctxCancelFunc()
	<-c.ctxCancelDone // this will block here until the channel is closed which maeans all cancel handlers are done
	return nil
}

func (c *ExitController) addCancelHandler(name string, handler func()) error {
	// check that the context hasn't been cancelled by getting error from the context
	if err := c.ctx.Err(); err != nil {
		return err
	}
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.handlers[name] = handler
	return nil
}

func (c *ExitController) waitForCancel() {
	<-c.ctx.Done() // this will block until the context is cancelled
	c.mutex.Lock()
	defer c.mutex.Unlock()
	log.Println("Cancellation signal received, will invoke cancel handlers...")
	var cancelWaitGroup sync.WaitGroup
	for handlerName, handlerFunc := range c.handlers {
		cancelWaitGroup.Add(1)
		log.Printf("Invoking cancel handlers: %s", handlerName)
		go invokeCancelHandler(&cancelWaitGroup, handlerFunc) // this will invoke each cancel handler in its own go routine
	}
	cancelWaitGroup.Wait() // this will block until all the cancel handlers are done
	close(c.ctxCancelDone) // this will indicate that all cancel handlers are done
}

func invokeCancelHandler(wg *sync.WaitGroup, handler func()) {
	handler()
	wg.Done()
}

func (c *ExitController) waitForSignal() { //wg *sync.WaitGroup) {
	signalChannel := signalChannel()
	for {
		select {
		case <-c.ctx.Done():
			log.Println("Not waiting for signal, context has been cancelled.")
			return
		case s, ok := <-signalChannel:
			if !ok {
				s = "unknown"
			}
			log.Printf("Termination signal %s received, cancelling context...", s)
			c.ctxCancelFunc() // this will cancel the context
			return
		}
	}
}

// Handle following signals: SIGHUP, SIGINT, SIGQUIT, SIGABRT, and SIGTERM
// Signal SIGKILL can not be handled by a program.

// signalChannel returns a channel that gets true value passed when terminating signal was received
func signalChannel() chan string {

	done := make(chan string)
	// create a channel and configure it to be notified of the signals
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM)
	// start a go routine that will wait for signal (any of the above)
	go func() {
		s := <-sc
		signal.Stop(sc)
		done <- s.String()
		close(done)
	}()

	return done
}
