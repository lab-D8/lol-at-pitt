package main

import (
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"

	"github.com/TrevorSStone/goriot"
	"github.com/lab-D8/lol-at-pitt/ols"
	"gopkg.in/mgo.v2"
	"log"
)

func dumpDb(filename string) {
	session, _ := mgo.Dial(MongoLocation)
	db := session.DB(DatabaseName)
	json_blob := map[string]interface{}{}
	var teams ols.Teams
	db.C("teams").Find(map[string]string{}).All(&teams)
	json_blob["Teams"] = teams

	var players ols.Players
	db.C("players").Find(map[string]string{}).All(&players)
	db.CollectionNames()
	json_blob["Players"] = players

	data, _ := json.MarshalIndent(json_blob, "", "  ")
	ioutil.WriteFile(filename, data, 0644)
}

func initDbPlayers(players ols.Players) {
	session, err := mgo.Dial(MongoLocation)
	if err != nil {
		panic(err)
	}
	db := session.DB(DatabaseName)
	db.C("players").DropCollection()
	for _, player := range players {
		db.C("players").Insert(player)
	}
	session.Close()

}

func initDbTeams(teams ols.Teams) {
	session, err := mgo.Dial(MongoLocation)
	defer session.Close()

	if err != nil {
		panic(err)
	}
	db := session.DB(DatabaseName)
	db.C("teams").DropCollection()

	for _, team := range teams {
		db.C("teams").Insert(team)
	}
}

func upload(json_file string) {
	var db_blob DB
	file, _ := os.Open(json_file)
	defer file.Close()

	data, _ := ioutil.ReadAll(file)
	json.Unmarshal(data, &db_blob)
	players := db_blob.Players
	initDbPlayers(players)

	teams := db_blob.Teams
	initDbTeams(teams)
}

func UploadPlayers(filename string) {
	r, _ := os.Open(filename)
	csvReader := csv.NewReader(r)
	allData, _ := csvReader.ReadAll()

	for _, record := range allData {
		normalizedSummonerName := goriot.NormalizeSummonerName(record[1])[0]
		amt, err := strconv.Atoi(record[2])

		if err != nil {
			log.Println("Skipping ", record[0])
			continue
		}

		proficiencies := []ols.Proficiency{}
		for _, i := range []int{3, 4, 5} {
			if record[i] == "None" {
				break
			} else {
				score, _ := strconv.Atoi(record[i+3])
				proficiencies = append(proficiencies, ols.Proficiency{
					Position: record[i],
					Score:    score,
				})
			}
		}

		totalProficiency, _ := strconv.ParseFloat(record[9], 64)

		player := &ols.Player{
			Name:             record[0],
			Ign:              record[1],
			NormalizedIgn:    normalizedSummonerName,
			Score:            amt,
			Proficiencies:    proficiencies,
			Id:               rand.Int63(),
			ProficiencyTotal: totalProficiency,
		}

		ols.GetPlayersDAO().Save(*player)
	}
}

func UploadCaptains(filename string) {
	r, _ := os.Open(filename)
	csvReader := csv.NewReader(r)
	allData, _ := csvReader.ReadAll()
	rand.Seed(51)
	for _, record := range allData[1:] {
		normalizedSummonerName := goriot.NormalizeSummonerName(record[1])[0]
		captain := &ols.Player{Name: record[0], Ign: record[1], NormalizedIgn: normalizedSummonerName, Id: rand.Int63()}
		if captain != nil {
			team := ols.Team{Name: captain.Ign + "'s team", Captain: captain.Id}
			team.Points, _ = strconv.Atoi(record[3])
			ols.GetPlayersDAO().Save(*captain)
			ols.GetTeamsDAO().Save(team)

		}
	}
}

func deleteDb() {
	session, err := mgo.Dial(MongoLocation)
	defer session.Close()

	if err != nil {
		panic(err)
	}
	db := session.DB(DatabaseName)
	db.DropDatabase()
}
