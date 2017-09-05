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
	"time"

	"github.com/minero/minero/proto/nbt"
)

var (
	APIURL              = "https://sessionserver.mojang.com/session/minecraft/profile/%s"
	SKINURL             = "https://visage.surgeplay.com/frontfull/50/%s"
	SKINDIR             = "minecraft-map/static/markers/"
	JSPATH              = "minecraft-map/player-markers.js"
	JSTMPLPATH          = "player-markers-tmpl.js"
	DATDIRS             = []string{"world/playerdata/"}
	PLAYERTIMEOUT int64 = 0
	CACHETIME     int64 = 24 * 60 * 60
	CACHEDIR            = ".player-marker-cache/"
	NOW                 = time.Now().Unix()
	OLDPLAYER           = errors.New("OLD PLAYER")
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
		log.Printf("Could not parse config file. Using defaults\n\n")
		fmt.Printf("\n%v\n", err)
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
		if SKINDIR[len(SKINDIR)-1] != '/' {
			SKINDIR += "/"
		}
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
	if config.PLAYERTIMEOUT > 0 {
		PLAYERTIMEOUT = config.PLAYERTIMEOUT
	}
	if config.CACHETIME > -1 {
		CACHETIME = config.CACHETIME * 60 * 60
	}
	if config.CACHEDIR != "" {
		CACHEDIR = config.CACHEDIR
		if CACHEDIR[len(CACHEDIR)-1] != '/' {
			CACHEDIR += "/"
		}
	}
	if CACHETIME > 0 {
		if err := os.MkdirAll(CACHEDIR, 0755); err != nil {
			log.Fatalf("\nError creating cache directory %s\nmain.go >> init() >> os.MkDirAll() >> %v\n", CACHEDIR, err)
		}
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

	playerC, errC := make(chan *Player), make(chan error)
	i := 0

	for _, DATDIR := range DATDIRS {

		if DATDIR[len(DATDIR)-1] != '/' {
			DATDIR += "/"
		}

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
				player, err := NewPlayer(path)
				if err != nil {
					if err == OLDPLAYER {
						errC <- err
					} else {
						errC <- fmt.Errorf("error creating new player from %s\nmain.go >> NewPlayer() >> %v\n", path, err)
					}
				} else {
					playerC <- player
				}
			}()
			i++
		}
	}

	var overworldPlayers []*Player
	var netherPlayers []*Player
	var endPlayers []*Player

	for j := 0; j < i; j++ {
		select {
		case player := <-playerC:
			switch player.Dimension {
			case -1:
				netherPlayers = append(netherPlayers, player)
			case 1:
				endPlayers = append(endPlayers, player)
			default:
				overworldPlayers = append(overworldPlayers, player)
			}
		case err := <-errC:
			if err != OLDPLAYER {
				log.Println(err)
			}
		}
	}

	if err := tmpl.Execute(jsFile, map[string]interface{}{
		"overworldPlayers": overworldPlayers,
		"endPlayers":       endPlayers,
		"netherPlayers":    netherPlayers,
	}); err != nil {
		log.Fatalf("error generating jsFile\nmain.go >> tmpl.Execute() >> %v\n", err)
	}

}

type Config struct {
	APIURL        string   `json:"api-url, omitempty"`
	SKINURL       string   `json:"skin-url, omitempty"`
	SKINDIR       string   `json:"skin-dir, omitempty"`
	JSPATH        string   `json:"js-path, omitempty"`
	JSTMPLPATH    string   `json:"js-tmpl-path, omitempty"`
	DATDIRS       []string `json:"dat-dirs, omitempty"`
	PLAYERTIMEOUT int64    `json:"player-timeout, omitempty"`
	CACHETIME     int64    `json:"cache-time, omitempty"`
	CACHEDIR      string   `json:"cache-dir, omitempty"`
}

type Player struct {
	X         int
	Y         int
	Z         int
	Dimension int
	Uuid      string
	Username  string
	ModTime   int64
	DifTime   int64
}

type MinecraftProfile struct {
	Name string `json:"name"`
}

func NewPlayer(path string) (*Player, error) {

	fs, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if PLAYERTIMEOUT > 0 && NOW-fs.ModTime().Unix() > PLAYERTIMEOUT {
		return nil, OLDPLAYER
	}

	player := &Player{}
	filename := filepath.Base(path)
	player.Uuid = strings.Replace(strings.TrimSuffix(filename, filepath.Ext(filename)), "-", "", -1)
	player.ModTime = fs.ModTime().Unix()
	player.DifTime = NOW - fs.ModTime().Unix()

	if err := player.SetLocation(path); err != nil {
		return nil, err
	}
	if err := player.SetUsername(); err != nil {
		return nil, err
	}

	if err := GetSkin(player.Uuid); err != nil {
		return nil, err
	}

	return player, nil
}

func (p *Player) SetLocation(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("SetLocation() >> os.Open() >> %v\n", err)
	}
	defer f.Close()

	r, err := gzip.NewReader(f)
	if err != nil {
		return fmt.Errorf("SetLocation() >> gzip.NewReader() >> %v\n", err)
	}
	defer r.Close()

	c, err := nbt.Read(r)
	if err != nil {
		return fmt.Errorf("SetLocation() >> nbt.Read() >> %v\n", err)
	}

	i := c.Value["Dimension"].(*nbt.Int32)

	p.Dimension = int(i.Int32)

	ps, ok := c.Value["Pos"].(*nbt.List)
	if !ok {
		return err
	}
	pos := ps.Value
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

	p.X = int(x.Float64)
	p.Y = int(y.Float64)
	p.Z = int(z.Float64)

	return nil
}

func (p *Player) SetUsername() error {

	fs, err := os.Stat(CACHEDIR + p.Uuid + ".txt")
	if err != nil {
		username, err := RequestUsername(p.Uuid)
		if err != nil {
			return fmt.Errorf("SetUsername() >> RequestUsername() >> %v\n", err)
		}
		p.Username = username
		return nil
	}

	if NOW-fs.ModTime().Unix() > CACHETIME {
		username, err := RequestUsername(p.Uuid)
		if err != nil {
			return fmt.Errorf("SetUsername() >> RequestUsername() >> %v\n", err)
		}
		p.Username = username
		return nil
	}

	b, err := ioutil.ReadFile(CACHEDIR + p.Uuid + ".txt")
	if err != nil {
		username, err := RequestUsername(p.Uuid)
		if err != nil {
			return fmt.Errorf("SetUsername() >> RequestUsername() >> %v\n", err)
		}
		p.Username = username
		return nil
	}
	p.Username = strings.Replace(string(b), "\n", "", -1)
	return nil
}

func RequestUsername(uuid string) (string, error) {
	resp, err := http.Get(fmt.Sprintf(APIURL, uuid))
	if err != nil {
		return "", fmt.Errorf("http.Get() >> %v\n", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ioutil.ReadAll() >> %v\n", err)
	}

	var mp MinecraftProfile
	if err := json.Unmarshal(body, &mp); err != nil {
		return "", fmt.Errorf("json.Unmarshal() >> %v\n", err)
	}

	if CACHETIME > 0 {
		ioutil.WriteFile(CACHEDIR+uuid+".txt", []byte(mp.Name), 0666)
	}

	return mp.Name, nil
}

func GetSkin(uuid string) error {

	fs, err := os.Stat(SKINDIR + uuid + ".png")
	if err != nil {
		if err := RequestSkin(uuid); err != nil {
			return fmt.Errorf("GetSkin() >> RequestSkin() >> %v\n", err)
		}
		return nil
	}

	if NOW-fs.ModTime().Unix() > CACHETIME {
		if err := RequestSkin(uuid); err != nil {
			return fmt.Errorf("GetSkin() >> RequestSkin() >> %v\n", err)
		}
		return nil
	}
	return nil
}

func RequestSkin(uuid string) error {
	resp, err := http.Get(fmt.Sprintf(SKINURL, uuid))
	if err != nil {
		return fmt.Errorf("http.Get() >> %v\n", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ioutil.ReadAll() >> %v\n", err)
	}

	if err := ioutil.WriteFile(SKINDIR+uuid+".png", body, 0666); err != nil {
		return fmt.Errorf("ioutil.WriteFile() >> %v\n", err)
	}

	return nil

}
