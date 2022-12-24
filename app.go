package main

import (
	"fmt"
	"net/url"
	"strings"
)

type App struct {
	hosts   *Hosts
	storage *Storage
}

func NewApp(
	hosts *Hosts,
	storage *Storage,
) *App {
	return &App{
		hosts:   hosts,
		storage: storage,
	}
}

type Cmd struct {
	Unblock bool
	Block   bool
	Website string
}

func (app *App) Handle(cmd Cmd) error {
	switch {
	case cmd.Website != "": // Banning a website
		return app.banWebsite(cmd.Website)
	case cmd.Unblock:
		return app.unblock()
	case cmd.Block:
		return app.block()
	default:
		if isPaused, err := app.storage.IsStatus(StatusPaused); err != nil {
			return err
		} else if isPaused {
			fmt.Println("Application in paused state. Doing nothing")
			return nil // Do nothing
		}

		return app.block()
	}
}

func (app *App) banWebsite(website string) error {
	if !strings.HasPrefix(website, "http") {
		website = "http://" + website
	}

	u, err := url.ParseRequestURI(website)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	web := strings.ToLower(u.Hostname())

	fmt.Printf("Banning %q...\n", web)
	return app.storage.BanWebsite(web)
}

func (app *App) block() error {
	if err := app.storage.Resume(); err != nil {
		return err
	}

	webs, err := app.storage.Websites()
	if err != nil {
		return err
	}

	return app.hosts.Block(webs)

}

func (app *App) unblock() error {
	webs, err := app.storage.Websites()
	if err != nil {
		return err
	}

	return app.hosts.Block(webs)
}
