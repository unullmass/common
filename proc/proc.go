/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package proc

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var stop chan os.Signal
var QuitChan chan bool
var wg *sync.WaitGroup
var err error
var mux sync.Mutex

var ErrWaitTimeout = errors.New("Timed out waiting for tasks to complete")

func init(){
	stop = make(chan os.Signal)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	QuitChan = make(chan bool)
	wg = &sync.WaitGroup{}
}


func AddTask() <-chan bool{
	wg.Add (1)
	return QuitChan
}


func TaskDone(){
	wg.Done()
}

func WaitForQuitAndCleanup(timeout time.Duration) error{
	WaitForQuitAndSignalTasks()
	return WaitFinalCleanup(timeout)
}

func WaitForQuitAndSignalTasks() {
	// wait for the stop signal for the process
	<-stop
	log.Println("Received quit. Sending shutdown and waiting on goroutines...")
	// send stop to the shutdown channel. All the routines waiting on the terminate signal
	// will receive when we close the channel
	close(QuitChan)
}

func WaitFinalCleanup(timeout time.Duration) (error){
	// wait for all the processes that added themselves to 
    c := make(chan struct{})
    go func() {
        defer close(c)
        wg.Wait()
	}()
	
	var success bool
    select {
    case <-c:
        success = true // completed normally
    case <-time.After(timeout*time.Second):
        success = false
	}
	
	if err != nil && success == false {
		return ErrWaitTimeout
	}
	return err
}

func SetError(e error){
	mux.Lock()
	err = e 
	mux.Unlock()
}

func EndProcess(){
	stop <- syscall.SIGTERM
}


