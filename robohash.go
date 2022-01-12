package gorobohash

import (
	"crypto/sha512"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/nfnt/resize"
)

//resource define resource path
type resource struct {
	// ResourceDir string
	BGSets    []string
	Colors    []string
	Format    string
	Hasharray []int16
	HexDigest string
	Iter      int
	HashCount int
	Sets      []string
}

//newResource constructor
func NewResource(name string) *resource {
	r := &resource{}
	r.HashCount = 11
	h := sha512.New()
	if _, err := h.Write([]byte(name)); err != nil {
		log.Fatalln(err)
	}
	r.HexDigest = hex.EncodeToString(h.Sum(nil))
	r.createHahes(r.HashCount)
	r.Sets = listDirs("material/sets")
	r.BGSets = listDirs("material/backgrounds")
	r.Colors = listDirs("material/sets/set1")
	r.Format = "png"
	r.Iter = 4
	log.Printf("%+v", r)
	return r
}

//createHahes createHahes
func (r *resource) createHahes(count int) {
	for i := 0; i < count; i++ {
		blockSize := len(r.HexDigest) / count
		currentStart := (1+i)*blockSize - blockSize
		currentEnd := (1 + i) * blockSize
		r.Hasharray = append(r.Hasharray, int16(binary.BigEndian.Uint16([]byte(r.HexDigest[currentStart:currentEnd]))))
	}
}

//listDirs dir list
func listDirs(path string) []string {
	rtn := []string{}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatalln(err)
		return rtn
	}
	for _, f := range files {
		rtn = append(rtn, f.Name())
	}
	return rtn
}

//getListOfFiles dir file list
func getListOfFiles(path string) []string {
	// chosenFiles := []string{}
	// directories := []string{}
	// fileList := []string{}
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		fmt.Println(info.Name())
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return []string{}
}

//assemble Build our Robot! Returns the robot image itself.
func (r *resource) assemble(roboset, colors, bgset, format string, x, y int) {
	roboset = r.Sets[0]
	if roboset == "any" {
		roboset = r.Sets[r.Hasharray[1]%int16(len(r.Sets))]
	}
	// if isContain(roboset, r.Sets) {
	// 	roboset = roboset
	// }

	if roboset == "set1" {
		if exist, _ := isContain(colors, r.Colors); exist {
			roboset = "set1/" + colors
		} else {
			roboset = "set1/" + r.Colors[r.Hasharray[0]%int16(len(r.Colors))]
		}
	}

	if exist, _ := isContain(bgset, r.BGSets); exist {
		bgset = bgset
	} else if bgset == "any" {
		bgset = r.BGSets[r.Hasharray[2]%int16(len(r.BGSets))]
	}

	if format == "" {
		format = r.Format
	}

	if x == 0 {
		x = 300
	}
	if y == 0 {
		y = 300
	}
	setsDir := "material/sets/" + roboset + "/"
	roboparts := getListOfFiles(setsDir)
	roboparts = sortSets(roboparts)
	if bgset != "" {
		bgList := []string{}
		backgroud := listDirs("material/sets/" + bgset)
		for _, v := range backgroud {
			if !strings.HasPrefix(v, ".") {
				bgList = append(bgList, "material/sets/"+bgset+v)
			}
		}
		backgroud = []string{bgList[r.Hasharray[3]%int16(len(bgList))]}
	}
	fmt.Println(roboparts[0])
	imgFile, err := os.Open(roboparts[0])
	if err != nil {
		log.Fatalln(err)
		return
	}
	imgDec, _, err := image.Decode(imgFile)
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer imgFile.Close()

	resizeImg := resize.Resize(1024, 1024, imgDec, resize.Lanczos3)

	out, err := os.Create("test." + format)
	if err != nil {
		log.Fatalln(err)
	}
	defer out.Close()
	if err := jpeg.Encode(out, resizeImg, nil); err != nil {
		log.Fatalln(err)
	}
}

//isContain is slice contain obj
func isContain(obj string, slice []string) (bool, int) {
	for k, v := range slice {
		if obj == v {
			return true, k
		}
	}
	return false, -1
}

//isContain is slice contain like obj
func isContainLike(obj string, slice []string) (bool, int) {
	for k, v := range slice {
		if strings.Contains(v, obj) {
			return true, k
		}
	}
	return false, -1
}

func Hex2Dec(val string) int {
	n, err := strconv.ParseUint(val, 16, 32)
	if err != nil {
		log.Fatalln(err)
	}
	return int(n)
}

//sortSets origin python: roboparts.sort(key=lambda x: x.split("#")[1])
func sortSets(sets []string) []string {
	var temp, rtn []string
	for _, v := range sets {
		temp = append(temp, strings.Split(v, "#")[1])
	}
	sort.Strings(temp)
	//E.g. [01Body 02Face Accessory Eyes Mouth]
	for _, v := range temp {
		if exist, key := isContainLike(v, sets); exist {
			rtn = append(rtn, sets[key])
		}
	}
	if len(rtn) == 0 {
		log.Fatal("sortSets length 0")
	}
	//[003#01Body 004#02Face 002#Accessory 001#Eyes 000#Mouth]
	return rtn
}
