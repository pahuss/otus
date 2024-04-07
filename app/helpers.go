package app

import (
	"encoding/csv"
	"fmt"
	"github.com/pahuss/otus/models"
	"github.com/pahuss/otus/repository"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func CopyRemoteCsvFile(fileUrl string, filePath string) {
	// if file already exists, do nothing
	if _, err := os.Stat(filePath); err == nil {
		return
	}

	resp, err := http.Get(fileUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	out, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
}

func ReadCsvPeopleToDatabase(filePath string, app *App) {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)

	rf := models.User{}
	r := &repository.UserRepository{
		Db: app.Db,
	}

	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		// do something with read line
		names := strings.Fields(rec[0])
		rf.FirstName = names[0]
		rf.LastName = names[1]
		rf.Age = stringToInt(rec[1])
		rf.Email = randomString(10) + "@social.net"

		i, err := r.InsertProfile(&rf)

		fmt.Printf("%+v %+v %+v\n", rec[0], i, err)
	}
}

func stringToInt(val string) int {
	i, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}
	return i
}

func randomString(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[randInt(len(letters))]
	}
	return string(b)
}

func randInt(n int) int {
	return int(rand.Int63n(int64(n)))
}
