//
// Copyright (c) 2019-2020 Krygon, Ltd.
//
// ALL RIGHTS RESERVED, DO NOT DISTRIBUTE!
//

package server

import (
	"net"
	"strconv"
	"time"

	"github.com/jaredfolkins/badactor"

	"github.com/pterodactyl/sftp-server/src/logger"
)

type Jailer struct {
	Studio *badactor.Studio
}

func NewJailer() *Jailer {
	jailer := &Jailer{
		Studio: badactor.NewStudio(1024),
	}

	jailer.Studio.AddRule(&badactor.Rule{
		Name:        "Authentication",
		Message:     "You are being rate-limited for being a retard.",
		StrikeLimit: 4,
		ExpireBase:  time.Second * 15,
		Sentence:    time.Minute * 5,
		Action:      &JailerAction{},
	})

	err := jailer.Studio.CreateDirectors(1024)
	if err != nil {
		panic(err)
		return nil
	}

	jailer.Studio.StartReaper(time.Second)

	return jailer
}

func (j *Jailer) Acquire(ip string) error {
	addr, _, err := net.SplitHostPort(ip)
	if err != nil {
		return err
	}

	if err := j.Studio.Infraction(addr, "Authentication"); err != nil {
		return err
	}

	strikes, err := j.Studio.Strikes(addr, "Authentication")
	if err != nil {
		return err
	}

	logger.Get().Info(addr + " has " + strconv.Itoa(strikes) + " strikes")
	return nil
}

func (j *Jailer) IsLimited(ip string) (bool, error) {
	addr, _, err := net.SplitHostPort(ip)
	if err != nil {
		return true, err
	}

	if j.Studio.IsJailed(addr) {
		logger.Get().Info(addr + " is rate-limited")
		return true, nil
	}

	return false, nil
}

type JailerAction struct {
}

func (j *JailerAction) WhenJailed(a *badactor.Actor, r *badactor.Rule) error {
	logger.Get().Info("someone has been jailed")
	return nil
}

func (j *JailerAction) WhenTimeServed(a *badactor.Actor, r *badactor.Rule) error {
	logger.Get().Info("someone has served their jail time")
	return nil
}
