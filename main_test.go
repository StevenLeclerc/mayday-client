package main

import (
	"os"
	"testing"
	"time"
)

func Test_run(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		//timer := time.NewTimer(time.Second * 2)
		//timer.C
		go main()
		time.Sleep(time.Second * 5)
		os.Exit(0)
	})
}
