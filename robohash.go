package gorobohash

import (
	"crypto/sha512"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	_ "image/jpeg"
	_ "image/png"

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
	Name      string
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
	ex, err := os.Executable()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(ex)
	expath := filepath.Dir(ex)
	fmt.Println(expath)

	dir, _ := os.Getwd()
	fmt.Println(dir)
	// r.BGSets = listDirs(expath + "\\material\\sets")
	r.Sets = listDirs("material/sets")
	r.BGSets = listDirs("material/backgrounds")
	r.Colors = listDirs("material/sets/set1")
	r.Format = "png"
	r.Iter = 4
	r.Name = name
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
func (r *resource) getListOfFiles(path string) []string {
	chosenFiles := []string{}
	directories := []string{}
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			directories = append(directories, path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	sort.Strings(directories)
	for _, v := range directories {
		filesInDir := []string{}
		err := filepath.Walk(v, func(v string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				filesInDir = append(filesInDir, v)
			}
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
		chosenFiles = append(chosenFiles, filesInDir[r.Hasharray[r.Iter]%int16(len(filesInDir))])
		r.Iter += 1
	}
	//chosenFiles E.g. [material\sets\set1\blue\003#01Body\000#blue_body-10.png material\sets\set1\blue\003#01Body\007#blue_body-06.png material\sets\set1\blue\004#02Face\000#blue_face-07.png material\sets\set1\blue\002#Accessory\002#blue_accessory-07.png material\sets\set1\blue\001#Eyes\009#blue_eyes-04.png material\sets\set1\blue\000#Mouth\009#blue_mouth-03.png]
	return chosenFiles
}

//assemble Build our Robot! Returns the robot image itself.
func (r *resource) Assemble(roboset, colors, bgset, format string, x, y int) {
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
	roboparts := r.getListOfFiles("material/sets/" + roboset + "/")
	roboparts = sortSets(roboparts)
	var background string
	if bgset != "" {
		bgList := []string{}
		temp := listDirs("material/sets/" + bgset)
		for _, v := range temp {
			if !strings.HasPrefix(v, ".") {
				bgList = append(bgList, "material/sets/"+bgset+v)
			}
		}
		background = bgList[r.Hasharray[3]%int16(len(bgList))]
	}
	// fmt.Println(roboparts)
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
	newImg := image.NewRGBA(image.Rect(0, 0, 1024, 1024))
	for _, v := range roboparts {
		tempFile, err := os.Open(v)
		if err != nil {
			log.Fatalln(err)
		}
		tempImg, _, err := image.Decode(tempFile)
		if err != nil {
			log.Fatal(err)
		}
		tempFile.Close()
		tempResizeImg := resize.Resize(1024, 1024, tempImg, resize.Lanczos3)

		draw.Draw(newImg, newImg.Bounds(), resizeImg, resizeImg.Bounds().Min, draw.Over)
		draw.Draw(newImg, newImg.Bounds(), tempResizeImg, tempImg.Bounds().Min, draw.Over)
	}
	if bgset != "" {
		imgFile, err := os.Open(background)
		if err != nil {
			log.Fatalln(err)
			return
		}
		tempImg, _, err := image.Decode(imgFile)
		if err != nil {
			log.Fatal(err)
		}
		imgFile.Close()
		tempResizeImg := resize.Resize(1024, 1024, tempImg, resize.Lanczos3)
		draw.Draw(newImg, newImg.Bounds(), tempResizeImg, tempImg.Bounds().Min, draw.Over)
	}
	out, err := os.Create(fmt.Sprintf("output/%s.%s", r.Name, format))
	if err != nil {
		log.Fatalln(err)
	}
	defer out.Close()
	if err := jpeg.Encode(out, newImg, nil); err != nil {
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
