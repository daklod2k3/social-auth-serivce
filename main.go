package main

import (
	"auth/cmd"
	"auth/internal/global"
	"shared/initialize"
	"sync"
)

func main() {

	wg := sync.WaitGroup{}
	wg.Add(2)

	global.InitGlobal()

	initialize.InitGlobal(&initialize.Type{
		Config: global.Config,
		Logger: global.Logger,
	})

	go func() {
		defer wg.Done()
		cmd.StartGRPCServer()

	}()

	go func() {
		defer wg.Done()
		go cmd.StartAPI()

	}()

	wg.Wait()
}
