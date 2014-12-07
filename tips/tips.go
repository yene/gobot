package tips

/*
 * load Loading Screen Tips
 * copyright valve
 */

import (
	"io/ioutil"
	"math/rand"
	"strings"
	"time"
)

func RandomTip() string {
	tips := tips()
	rand.Seed(time.Now().UTC().UnixNano())
	i := rand.Intn(len(tips))
	return strings.TrimSpace(tips[i])
}

func tips() []string {
	file, e := ioutil.ReadFile("./tips.txt")
	if e != nil {
		panic(e)
	}
	return strings.Split(string(file), "---")
}
