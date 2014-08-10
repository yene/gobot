package wisdom

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"
)

func main() {
	fmt.Println(RandomWisdom())
}

func RandomWisdom() string {
	wisdoms := wisdoms()
	rand.Seed(time.Now().UTC().UnixNano())
	i := rand.Intn(len(wisdoms))
	return wisdoms[i]
}

func wisdoms() []string {
	file, e := ioutil.ReadFile("./wisdom.txt")
	if e != nil {
		panic(e)
	}
	return strings.Split(string(file), "\n")
}
