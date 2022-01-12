package gorobohash

import (
	"crypto/sha512"
	"encoding/binary"
	"encoding/hex"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"strconv"

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

//listDirs
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

func getListOfFiles(path string) {

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
		if isContain(colors, r.Colors) {
			roboset = "set1/" + colors
		} else {
			radomColor := r.Colors[r.Hasharray[0]%int16(len(r.Colors))]
			roboset = "set1/" + radomColor
		}
	}

	if isContain(bgset, r.BGSets) {
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

	roboparts := listDirs("material/sets/" + roboset)

	if bgset != "" {

	}

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
	imgFile.Close()

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

func isContain(a string, b []string) bool {
	for _, v := range b {
		if a == v {
			return true
		}
	}
	return false
}

func Hex2Dec(val string) int {
	n, err := strconv.ParseUint(val, 16, 32)
	if err != nil {
		log.Fatalln(err)
	}
	return int(n)
}
