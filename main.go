package main

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/minero/minero/proto/nbt"
)

var (
	APIURL     = "https://sessionserver.mojang.com/session/minecraft/profile/"
	SKINURL    = "https://visage.surgeplay.com/frontfull/50/"
	SKINDIR    = "minecraft-map/static/markers/"
	JSPATH     = "minecraft-map/player-markers.js"
	JSTMPLPATH = "player-markers-tmpl.js"
	DATDIRS    = []string{"world/playerdata/"}
)

func init() {

	configPath := flag.String("c", "config.json", "full path to config file in JSON format")
	flag.Parse()
	if *configPath == "" {
		*configPath = "config.json"
	}

	b, err := ioutil.ReadFile(*configPath)
	if err != nil {
		log.Println("Could not read config file. Using defaults")
		return
	}

	var config Config
	if err := json.Unmarshal(b, &config); err != nil {
		log.Println("Could not parse config file. Using defaults")
		return
	}

	if config.APIURL != "" {
		APIURL = config.APIURL
	}
	if config.SKINURL != "" {
		SKINURL = config.SKINURL
	}
	if config.SKINDIR != "" {
		SKINDIR = config.SKINDIR
	}
	if config.JSPATH != "" {
		JSPATH = config.JSPATH
	}
	if config.JSTMPLPATH != "" {
		JSTMPLPATH = config.JSTMPLPATH
	}
	if config.DATDIRS != nil && len(config.DATDIRS) > 0 {
		DATDIRS = config.DATDIRS
	}
}

func main() {
	if err := os.MkdirAll(SKINDIR, 0777); err != nil {
		log.Fatalf("error creating skins directory %s\nmain.go >> os.MkdirAll() >> %v\n", SKINDIR, err)
	}

	tmpl, err := template.ParseFiles(JSTMPLPATH)
	if err != nil {
		log.Fatalf("error parsing template %q\nmain.go >> template.ParseFiles() >> %v\v", JSTMPLPATH, err)
	}

	jsFile, err := os.Create(JSPATH)
	if err != nil {
		log.Fatalf("error creating js file %s\nmain.go >> os.Create() >> %v\n", JSPATH, err)
	}

	userC, errC := make(chan *User), make(chan error)
	i := 0

	for _, DATDIR := range DATDIRS {

		files, err := ioutil.ReadDir(DATDIR)
		if err != nil {
			log.Printf("error reading directory %s\nmain.go >> ioutil.ReadDir() >> %v\n", DATDIR, err)
			continue
		}

		for _, fi := range files {
			if filepath.Ext(fi.Name()) != ".dat" {
				continue
			}

			path := DATDIR + fi.Name()

			go func() {
				user, err := NewUser(path)
				if err != nil {
					errC <- fmt.Errorf("error creating new user from %s\nmain.go >> NewUser() >> %v\n", path, err)
				}
				userC <- user
			}()
			i++
		}
	}

	var users []*User

	for j := 0; j < i; j++ {
		select {
		case user := <-userC:
			users = append(users, user)
		case err := <-errC:
			log.Println(err)
		}
	}

	if err := tmpl.Execute(jsFile, users); err != nil {
		log.Fatalf("error generating jsFile\nmain.go >> tmpl.Execute() >> %v\n", err)
	}

}

func CreateUser(path string, users chan *User) {

	user, err := NewUser(path)
	if err != nil {
		log.Printf("error creating new user from %s\nmain.go >> NewUser() >> %v\n", path, err)
		return
	}

	users <- user

}

type Config struct {
	APIURL     string   `json:"api-url, omitempty"`
	SKINURL    string   `json:"skin-url, omitempty"`
	SKINDIR    string   `json:"skin-dir, omitempty"`
	JSPATH     string   `json:"js-path, omitempty"`
	JSTMPLPATH string   `json:"js-tmpl-path,omitempty"`
	DATDIRS    []string `json:"dat-dirs, omitempty"`
}

type User struct {
	X         int
	Y         int
	Z         int
	Dimension int
	Uuid      string
	Username  string
}

type MinecraftProfile struct {
	Name string `json:"name"`
}

func NewUser(path string) (*User, error) {
	user := &User{}
	filename := filepath.Base(path)
	user.Uuid = strings.Replace(strings.TrimSuffix(filename, filepath.Ext(filename)), "-", "", -1)

	if err := user.SetLocation(path); err != nil {
		return nil, err
	}

	if err := user.SetUsername(); err != nil {
		return nil, err
	}

	if err := user.GetSkin(); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *User) SetLocation(path string) error {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	r, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer r.Close()

	c, err := nbt.Read(r)
	if err != nil {
		return err
	}

	i := c.Value["Dimension"].(*nbt.Int32)

	u.Dimension = int(i.Int32)

	p, ok := c.Value["Pos"].(*nbt.List)
	if !ok {
		return err
	}
	pos := p.Value
	if len(pos) != 3 {
		return fmt.Errorf("pos wrong length. expected 3 got %d", len(pos))
	}

	x, ok := pos[0].(*nbt.Float64)
	if !ok {
		return errors.New("Invalid \"x\" type")
	}

	y, ok := pos[1].(*nbt.Float64)
	if !ok {
		return errors.New("Invalid \"y\" type")
	}

	z, ok := pos[2].(*nbt.Float64)
	if !ok {
		return errors.New("Invalid \"z\" type")
	}

	u.X = int(x.Float64)
	u.Y = int(y.Float64)
	u.Z = int(z.Float64)

	return nil
}

func (u *User) SetUsername() error {
	resp, err := http.Get(APIURL + u.Uuid)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var p MinecraftProfile
	if err := json.Unmarshal(body, &p); err != nil {
		return err
	}

	u.Username = p.Name

	return nil
}

func (u *User) GetSkin() error {
	resp, err := http.Get(SKINURL + u.Uuid)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(SKINDIR+u.Uuid+".png", body, 0666); err != nil {
		return err
	}

	return nil
}
