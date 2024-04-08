package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/yeeckfy/gomp"
)

const (
	ColorWhite gomp.Color = 0xFFFFFFFF
)

type City int

const (
	CityLosSantos City = iota
	CitySanFierro
	CityLasVenturas
)

type Character struct {
	*gomp.Player
	citySelection      City
	hasCitySelected    bool
	lastCitySelectedAt time.Time
}

var chars = make(map[int]*Character, 1000)

var lsTd *gomp.Textdraw
var sfTd *gomp.Textdraw
var lvTd *gomp.Textdraw
var classSelHelperTd *gomp.Textdraw

var vehFiles = []string{
	"bone.txt", "flint.txt", "ls_airport.txt", "ls_gen_inner.txt", "ls_gen_outer.txt",
	"ls_law.txt", "lv_airport.txt", "lv_gen.txt", "lv_law.txt", "pilots.txt",
	"red_county.txt", "sf_airport.txt", "sf_gen.txt", "sf_law.txt", "sf_train.txt",
	"tierra.txt", "trains_platform.txt", "trains.txt", "whetstone.txt",
}

func onGameModeInit(evt *gomp.GameModeInitEvent) {
	gomp.SetGameModeText("Grand Larceny")
	gomp.SetPlayerMarkerMode(gomp.PlayerMarkerModeGlobal)
	gomp.EnableNametags()
	gomp.SetNametagDrawRadius(40.0)
	gomp.EnableStuntBonuses()
	gomp.DisableEntryExitMarkers()
	gomp.SetWeather(2)
	gomp.SetWorldTime(11)

	lsTd, _ = NewCityNameTextdraw("Los Santos")
	sfTd, _ = NewCityNameTextdraw("San Fierro")
	lvTd, _ = NewCityNameTextdraw("Las Venturas")

	classSelHelperTd, _ = gomp.NewTextdraw(gomp.Vector2{X: 10.0, Y: 415.0}, "Press ~b~~k~~GO_LEFT~ ~w~or ~b~~k~~GO_RIGHT~ ~w~to switch cities.~n~ Press ~r~~k~~PED_FIREWEAPON~ ~w~to select.", nil)
	classSelHelperTd.EnableBox()
	classSelHelperTd.SetBoxColor(0x222222BB)
	classSelHelperTd.SetLetterSize(gomp.Vector2{X: 0.3, Y: 1.0})
	classSelHelperTd.SetTextSize(gomp.Vector2{X: 400.0, Y: 40.0})
	classSelHelperTd.SetStyle(gomp.TextdrawStyle2)
	classSelHelperTd.SetShadow(0)
	classSelHelperTd.SetOutline(1)
	classSelHelperTd.SetBackgroundColor(0x000000FF)
	classSelHelperTd.SetColor(0xFFFFFFFF)

	gomp.NewClass(0, 298, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 298, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 299, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 300, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 301, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 302, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 303, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 304, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 305, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 280, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 281, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 282, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 283, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 284, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 285, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 286, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 287, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 288, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 289, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 265, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 266, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 267, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 268, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 269, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 270, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 1, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 2, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 3, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 4, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 5, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 6, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 8, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 42, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 65, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	//gomp.NewClass(0, 74, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 86, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 119, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 149, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 208, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 273, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 289, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)

	gomp.NewClass(0, 47, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 48, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 49, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 50, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 51, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 52, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 53, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 54, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 55, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 56, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 57, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 58, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 68, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 69, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 70, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 71, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 72, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 73, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 75, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 76, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 78, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 79, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 80, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 81, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 82, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 83, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 84, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 85, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 87, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 88, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 89, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 91, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 92, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 93, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 95, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 96, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 97, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 98, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)
	gomp.NewClass(0, 99, gomp.Vector3{X: 1759.0189, Y: -1898.1260, Z: 13.5622}, 266.4503, -1, -1, -1, -1, -1, -1)

	var vehCnt int
	for _, vehFile := range vehFiles {
		filename := filepath.Join("scriptfiles/vehicles", vehFile)

		cnt, err := LoadStaticVehiclesFromFile(filename)
		if err != nil {
			fmt.Printf("Failed to load vehicles from %s: %v\n", filename, err)
		}

		vehCnt += cnt
	}

	fmt.Printf("Total vehicles from files: %d\n", vehCnt)
}

func onPlayerConnect(evt *gomp.PlayerConnectEvent) {
	char := &Character{
		Player:             evt.Player,
		citySelection:      -1,
		hasCitySelected:    false,
		lastCitySelectedAt: time.Now(),
	}

	chars[char.ID()] = char

	char.ShowGameText("~w~Grand Larceny", 3*time.Second, 4)
	char.SendMessage("Welcome to {88AA88}G{FFFFFF}rand {88AA88}L{FFFFFF}arceny", ColorWhite)
}

func onPlayerSpawn(evt *gomp.PlayerSpawnEvent) {
	char := chars[evt.Player.ID()]

	if char.IsBot() {
		return
	}

	char.SetInterior(0)
	char.ShowClock()
	char.ResetMoney()
	char.GiveMoney(30000)
}

func onPlayerRequestClass(evt *gomp.PlayerRequestClassEvent) {
	char := chars[evt.Player.ID()]

	if char.IsBot() {
		return
	}

	if char.hasCitySelected {
		setupCharSelection(char)
		return
	}

	if char.State() != gomp.PlayerStateSpectating {
		char.EnableSpectating()
		classSelHelperTd.ShowFor(char.Player)
		char.citySelection = -1
	}
}

func onPlayerUpdate(evt *gomp.PlayerUpdateEvent) {
	char := chars[evt.Player.ID()]

	if char.IsBot() {
		return
	}

	if !char.hasCitySelected && char.State() == gomp.PlayerStateSpectating {
		handleCitySelection(char)
		return
	}

	if char.ArmedWeapon() == gomp.WeaponMinigun {
		char.Kick()
		return
	}
}

func onPlayerDeath(evt *gomp.PlayerDeathEvent) {
	char := chars[evt.Player.ID()]

	char.hasCitySelected = false

	var killer *Character
	if evt.Killer != nil {
		killer = chars[evt.Killer.ID()]
	}

	if killer == nil {
		char.ResetMoney()
		return
	}

	if char.Money() > 0 {
		killer.GiveMoney(char.Money())
		char.ResetMoney()
	}
}

func NewCityNameTextdraw(cityName string) (*gomp.Textdraw, error) {
	td, err := gomp.NewTextdraw(gomp.Vector2{X: 10.0, Y: 380.0}, cityName, nil)
	if err != nil {
		return nil, err
	}

	td.DisableBox()
	td.SetLetterSize(gomp.Vector2{X: 1.25, Y: 3.0})
	td.SetStyle(gomp.TextdrawStyle0)
	td.SetShadow(0)
	td.SetOutline(1)
	td.SetColor(0xEEEEEEFF)
	td.SetBackgroundColor(0x000000FF)

	return td, nil
}

func LoadStaticVehiclesFromFile(filename string) (int, error) {
	inf, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer inf.Close()

	r := bufio.NewReader(inf)

	const delim = '\n'
	var cnt int
	var eof bool

	for !eof {
		line, err := r.ReadString(delim)
		if err != nil {
			if err == io.EOF {
				eof = true
			} else {
				return cnt, err
			}
		}

		line = strings.TrimSuffix(line, string(delim))

		split := strings.Split(line, ",")

		model, err := strconv.ParseInt(split[0], 10, 0)
		if err != nil {
			continue
		}

		spawnX, err := strconv.ParseFloat(split[1], 32)
		if err != nil {
			continue
		}

		spawnY, err := strconv.ParseFloat(split[2], 32)
		if err != nil {
			continue
		}

		spawnZ, err := strconv.ParseFloat(split[3], 32)
		if err != nil {
			continue
		}

		rot, err := strconv.ParseFloat(split[4], 32)
		if err != nil {
			continue
		}

		primaryColor, err := strconv.ParseInt(split[5], 10, 0)
		if err != nil {
			continue
		}
		_ = primaryColor

		secColAndName := strings.Split(split[6], ";")

		secondaryColor, err := strconv.ParseInt(strings.TrimSpace(secColAndName[0]), 10, 0)
		if err != nil {
			continue
		}
		_ = secondaryColor

		veh, err := gomp.NewStaticVehicle(gomp.VehicleModel(model), gomp.Vector3{X: float32(spawnX), Y: float32(spawnY), Z: float32(spawnZ)}, float32(rot))
		if err != nil {
			continue
		}

		veh.SetColor(gomp.VehicleColor{
			Primary:   int(primaryColor),
			Secondary: int(secondaryColor),
		})

		cnt++
	}

	fmt.Printf("Loaded %d vehicles from: %s\n", cnt, filename)

	return cnt, nil
}

func handleCitySelection(char *Character) {
	if char.citySelection == -1 {
		switchToNextCity(char)
		return
	}

	if time.Since(char.lastCitySelectedAt) < 500*time.Millisecond {
		return
	}

	keyData := char.KeyData()

	if keyData.Keys&gomp.PlayerKeyFire != 0 {
		char.hasCitySelected = true
		lsTd.HideFor(char.Player)
		sfTd.HideFor(char.Player)
		lvTd.HideFor(char.Player)
		classSelHelperTd.HideFor(char.Player)
		char.DisableSpectating()
		return
	}

	if keyData.LeftRight > 0 {
		switchToNextCity(char)
	} else if keyData.LeftRight < 0 {
		switchToPrevCity(char)
	}
}

func switchToNextCity(char *Character) {
	char.citySelection++
	if char.citySelection > CityLasVenturas {
		char.citySelection = CityLosSantos
	}

	char.PlaySound(1052, gomp.Vector3{})
	char.lastCitySelectedAt = time.Now()
	setupSelectedCity(char)
}

func switchToPrevCity(char *Character) {
	char.citySelection--
	if char.citySelection < CityLosSantos {
		char.citySelection = CityLasVenturas
	}

	char.PlaySound(1053, gomp.Vector3{})
	char.lastCitySelectedAt = time.Now()
	setupSelectedCity(char)
}

func setupSelectedCity(char *Character) {
	if char.citySelection == -1 {
		char.citySelection = CityLosSantos
	}

	char.SetInterior(0)

	if char.citySelection == CityLosSantos {
		char.SetCameraPosition(gomp.Vector3{X: 1630.6136, Y: -2286.0298, Z: 110.0})
		char.SetCameraLookAt(gomp.Vector3{X: 1887.6034, Y: -1682.1442, Z: 47.6167}, gomp.PlayerCameraCutTypeCut)

		lsTd.ShowFor(char.Player)
		sfTd.HideFor(char.Player)
		lvTd.HideFor(char.Player)
	} else if char.citySelection == CitySanFierro {
		char.SetCameraPosition(gomp.Vector3{X: -1300.8754, Y: 68.0546, Z: 129.4823})
		char.SetCameraLookAt(gomp.Vector3{X: -1817.9412, Y: 769.3878, Z: 132.6589}, gomp.PlayerCameraCutTypeCut)

		lsTd.HideFor(char.Player)
		sfTd.ShowFor(char.Player)
		lvTd.HideFor(char.Player)
	} else if char.citySelection == CityLasVenturas {
		char.SetCameraPosition(gomp.Vector3{X: 1310.6155, Y: 1675.9182, Z: 110.7390})
		char.SetCameraLookAt(gomp.Vector3{X: 2285.2944, Y: 1919.3756, Z: 68.2275}, gomp.PlayerCameraCutTypeCut)

		lsTd.HideFor(char.Player)
		sfTd.HideFor(char.Player)
		lvTd.ShowFor(char.Player)
	}
}

func setupCharSelection(char *Character) {
	if char.citySelection == CityLosSantos {
		char.SetInterior(11)
		char.SetPosition(gomp.Vector3{X: 508.7362, Y: -87.4335, Z: 998.9609})
		char.SetFacingAngle(0.0)
		char.SetCameraPosition(gomp.Vector3{X: 508.7362, Y: -83.4335, Z: 998.9609})
		char.SetCameraLookAt(gomp.Vector3{X: 508.7362, Y: -87.4335, Z: 998.9609}, gomp.PlayerCameraCutTypeCut)
	} else if char.citySelection == CitySanFierro {
		char.SetInterior(3)
		char.SetPosition(gomp.Vector3{X: -2673.8381, Y: 1399.7424, Z: 918.3516})
		char.SetFacingAngle(181.0)
		char.SetCameraPosition(gomp.Vector3{X: -2673.2776, Y: 1394.3859, Z: 918.3516})
		char.SetCameraLookAt(gomp.Vector3{X: -2673.8381, Y: 1399.7424, Z: 918.3516}, gomp.PlayerCameraCutTypeCut)
	} else if char.citySelection == CityLasVenturas {
		char.SetInterior(3)
		char.SetPosition(gomp.Vector3{X: 349.0453, Y: 193.2271, Z: 1014.1797})
		char.SetFacingAngle(286.25)
		char.SetCameraPosition(gomp.Vector3{X: 352.9164, Y: 194.5702, Z: 1014.1875})
		char.SetCameraLookAt(gomp.Vector3{X: 349.0453, Y: 193.2271, Z: 1014.1797}, gomp.PlayerCameraCutTypeCut)
	}
}

func main() {}

func init() {
	gomp.On(gomp.EventTypeGameModeInit, onGameModeInit)
	gomp.On(gomp.EventTypePlayerConnect, onPlayerConnect)
	gomp.On(gomp.EventTypePlayerSpawn, onPlayerSpawn)
	gomp.On(gomp.EventTypePlayerRequestClass, onPlayerRequestClass)
	gomp.On(gomp.EventTypePlayerUpdate, onPlayerUpdate)
	gomp.On(gomp.EventTypePlayerDeath, onPlayerDeath)
}