package main

import (
	"fmt"
	"os"
	"path"

	"bitbucket.org/mischief/libauth"
	"github.com/altid/libs/fs"
	"github.com/mischief/ndb"
)

type config struct {
	log  string
	pass string
	user string
}

func newConfig() (*config, error) {
	confdir, err := fs.UserConfDir()
	if err != nil {
		return nil, err
	}
	filePath := path.Join(confdir, "altid", "config")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, err
	}
	conf, err := ndb.Open(filePath)
	if err != nil {
		return nil, err
	}
	recs := conf.Search("service", *srv)
	switch len(recs) {
	case 0:
		return nil, fmt.Errorf("unable to find entry for %s", *srv)
	case 1:
		return readRecord(recs[0])
	}
	return nil, fmt.Errorf("ound multiple entries for %s, unable to continue", *srv)
}

func readRecord(rec ndb.Record) (*config, error) {
	datadir, err := fs.UserShareDir()
	if err != nil {
		datadir = "/tmp/altid"
	}
	conf := &config{
		log: path.Join(datadir, *srv),
	}
	for _, tup := range rec {
		switch tup.Attr {
		case "log":
			conf.log = path.Join(tup.Val, *srv)
		case "auth":
			conf.pass = tup.Val
		case "user":
			conf.user = tup.Val
		}
	}
	if len(conf.pass) > 5 && conf.pass[:5] == "pass=" {
		conf.pass = conf.pass[5:]
	}
	if conf.pass == "factotum" {
		userPwd, err := libauth.Getuserpasswd(
			"proto=pass service=discord server=discordapp.com user=%s",
			conf.user,
		)
		if err != nil {
			return nil, err
		}
		conf.pass = userPwd.Password
	}
	return conf, nil
}
