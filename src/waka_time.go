package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type Language struct {
	Name         string  `json:"name"`
	TotalSeconds float64 `json:"total_seconds"`
	Percent      float64 `json:"percent"`
	Digital      string  `json:"digital"`
	Text         string  `json:"text"`
	Hours        int     `json:"hours"`
	Minutes      int     `json:"minutes"`
	Seconds      int     `json:"seconds"`
}

type Data struct {
	Languages []Language `json:"languages"`
}

type WakaTimeStats struct {
	Data `json:"data"`
}

func getWakaTimeStats(key string) *WakaTimeStats {
	url := fmt.Sprintf("https://wakatime.com/api/v1/users/current/stats/last_7_days?api_key=%s", key)

	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	statistics := new(WakaTimeStats)
	err = json.Unmarshal(responseData, statistics)
	if err != nil {
		log.Fatal(err)
	}

	return statistics
}

func getRepresentData(languages []Language,
	showTitle bool, tableHeight int, cellPadding int,
	blockLength int, blockRunes []rune,
	startSection string, endSection string) string {

	namePadding := getPad(languages, tableHeight)

	var buffer strings.Builder

	buffer.WriteString(startSection)

	buffer.WriteString("\n```text\n")

	if showTitle {
		buffer.WriteString(getHeader())
	}

	for i := range languages {
		if i == tableHeight {
			break
		}
		lang := &languages[i]

		if lang.Hours == 0 && lang.Minutes == 0 {
			continue
		}

		buffer.WriteString(fmt.Sprintf("%*[2]s%-16[3]s%*[5]s%05.2[6]f%%\n",
			-1*namePadding+cellPadding,
			lang.Name,
			lang.Text,
			-1*blockLength+cellPadding,
			makeProgressGraph(lang.Percent, blockLength, blockRunes),
			lang.Percent))
	}

	buffer.WriteString("```\n")
	buffer.WriteString(endSection)

	return buffer.String()
}

func getPad(languages []Language, tableHeight int) int {
	var pad int
	for i := range languages {
		if i == tableHeight {
			break
		}
		if length := len(languages[i].Name); length > pad {
			pad = length
		}
	}
	return pad
}

func getHeader() string {
	weekEnd := time.Now().AddDate(0, 0, -1)
	weekStart := weekEnd.AddDate(0, 0, -7)
	_, week := weekStart.ISOWeek()

	return fmt.Sprintf("Week #%v: %s - %s\n", week, weekStart.Format("02 January, 2006"), weekEnd.Format("02 January, 2006"))
}

func makeProgressGraph(percent float64, blockLength int, blockRunes []rune) string {
	n := len(blockRunes) - 1

	doneBlock, remainderBlock, emptyBlock := getBlockSizes(percent, float64(blockLength), blockRunes)

	var graph strings.Builder
	graph.WriteString(strings.Repeat(string(blockRunes[n]), doneBlock))
	if remainderBlock > 0 {
		graph.WriteString(string(blockRunes[remainderBlock]))
		emptyBlock--
	}
	graph.WriteString(strings.Repeat(string(blockRunes[0]), emptyBlock))

	return graph.String()
}

func getBlockSizes(percent float64, blockLength float64, blockRunes []rune) (int, int, int) {
	n := float64(len(blockRunes)) - 1

	doneBlock := math.Floor(blockLength*percent/100.0 + 0.5/n)
	remainderBlock := math.Floor((blockLength*percent/100.0-doneBlock)*n + 0.5)
	emptyBlock := blockLength - doneBlock

	return int(doneBlock), int(remainderBlock), int(emptyBlock)
}

func replaceContent(oldReadme, wakaData, startComment, endComment string) string {
	var regex = fmt.Sprintf("%v(\\w|\\W)+%v", startComment, endComment)

	var re = regexp.MustCompile(regex)
	newReadme := re.ReplaceAllString(oldReadme, wakaData)

	return newReadme
}
