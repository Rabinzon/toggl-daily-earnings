package main

import (
	"github.com/jason0x43/go-toggl"
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {
	var earling int
	file, err := ioutil.ReadFile(".config")
	if err != nil {
		log.Fatal(err)
	}

	config := strings.Split(string(file), "\n")
	rate, _:= strconv.ParseFloat(config[2], 64)

	account, err := toggl.NewSession(config[0], config[1])

	if err != nil {
		log.Fatal(err)
	}

	s, err := account.GetAccount()
	workSpaceId := s.Data.Workspaces[0].ID
	total, err := getTotalGrand(account, workSpaceId, rate)

	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("bash", "-c", "$(pwd)/indicator")
	grepIn, _ := cmd.StdinPipe()

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	for {
		s, err = account.GetAccount()
		workSpaceId = s.Data.Workspaces[0].ID

		if err != nil {
			log.Fatal(err)
		}

		TimeEntries := s.Data.TimeEntries
		last := TimeEntries[len(TimeEntries)-1]

		if last.Stop == nil {
			current := toHours(int((time.Now().Unix() + last.Duration) * 1000))
			earling = int(current * rate) + total
		} else {
			earling, _ = getTotalGrand(account, workSpaceId, rate)
			total = earling
		}

		grepIn.Write([]byte(strconv.Itoa(earling) + "\n"))
		time.Sleep(20000 * time.Millisecond)
	}
}

func getTotalGrand(account toggl.Session, workSpaceId int, rate float64) (int, error) {
	today := time.Now().Local().Format("2006-01-02")
	report, err := account.GetSummaryReport(workSpaceId, today, today)

	if err != nil {
		return 0, err
	}

	return int(toHours(report.TotalGrand) * rate), nil
}

func toHours(m int) float64 {
	return float64(m) / float64(3600000)
}
