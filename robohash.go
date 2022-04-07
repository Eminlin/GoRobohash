package gorobohash

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"

	_ "image/png"

	"github.com/nfnt/resize"
)

//resource define resource path
type resource struct {
	bgSets       []string
	colors       []string
	format       string
	hashArray    []int64
	hexDigest    string
	iter         int
	hashCount    int
	sets         []string
	name         string
	materialPath string
	AssembleOptions
}

type AssembleOptions struct {
	RoboSet, Colors, BgSet string //optional
	Format                 string //optional default png
	OutputPath             string //optional default current dir
	X                      int    //optional default 300
	Y                      int    //optional default 300
}

//newResource new resource
func NewResource(text string, options *AssembleOptions) *resource {
	if text == "" {
		log.Fatalln("text is empty, you must set text")
	}
	r := &resource{}
	r.hashCount = 11
	h := sha512.New()
	if _, err := h.Write([]byte(text)); err != nil {
		log.Fatalln(err)
	}
	r.hexDigest = hex.EncodeToString(h.Sum(nil))
	r.createHahes(r.hashCount)
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatalln("runtime.Caller(0) failed")
	}
	r.materialPath = filepath.Dir(filename) + "/" + "material"
	r.sets = listDirs(r.materialPath + "/sets")
	r.bgSets = listDirs(r.materialPath + "/backgrounds")
	r.colors = listDirs(r.materialPath + "/sets/set1")
	r.Format = "png"
	r.iter = 4
	r.name = text
	if options.Format != "" {
		r.Format = "png"
	}
	if options.X == 0 {
		r.X = 300
	}
	if options.Y == 0 {
		r.Y = 300
	}
	if options.BgSet == "bmp" {
		log.Fatalln("bmp is not supported yet")
	}
	if options.OutputPath == "" {
		r.OutputPath = "./"
	}
	r.OutputPath = filepath.Dir(options.OutputPath)
	return r
}

//createHahes createHahes
func (r *resource) createHahes(count int) {
	for i := 0; i < count; i++ {
		blockSize := len(r.hexDigest) / count
		currentStart := (1+i)*blockSize - blockSize
		currentEnd := (1 + i) * blockSize
		temp, err := strconv.ParseInt(r.hexDigest[currentStart:currentEnd], 16, 64)
		if err != nil {
			log.Fatalln(err)
		}
		r.hashArray = append(r.hashArray, temp)
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
		if info.IsDir() && strings.Contains(path, "#") {
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
		chosenFiles = append(chosenFiles, filesInDir[r.hashArray[r.iter]%int64(len(filesInDir))])
		r.iter += 1
	}
	//chosenFiles E.g. [material\sets\set1\blue\003#01Body\000#blue_body-10.png material\sets\set1\blue\003#01Body\007#blue_body-06.png material\sets\set1\blue\004#02Face\000#blue_face-07.png material\sets\set1\blue\002#Accessory\002#blue_accessory-07.png material\sets\set1\blue\001#Eyes\009#blue_eyes-04.png material\sets\set1\blue\000#Mouth\009#blue_mouth-03.png]
	return chosenFiles
}

//Assemble Build our Robot! Return the robot image name itself.
func (r *resource) Assemble() (string, error) {
	if r.RoboSet == "any" {
		r.RoboSet = r.sets[r.hashArray[1]%int64(len(r.sets))]
	} else if status, _ := isContain(r.RoboSet, r.sets); !status {
		r.RoboSet = r.sets[0]
	}
	if r.RoboSet == "set1" {
		if exist, _ := isContain(r.Colors, r.colors); exist {
			r.RoboSet = "set1/" + r.Colors
		} else {
			r.RoboSet = "set1/" + r.colors[r.hashArray[0]%int64(len(r.colors))]
		}
	}
	if r.BgSet == "any" {
		r.BgSet = r.bgSets[r.hashArray[2]%int64(len(r.bgSets))]
	}
	roboparts := r.getListOfFiles(r.materialPath + "/sets/" + r.RoboSet + "/")
	roboparts = sortSets(roboparts)
	var background string
	if r.BgSet != "" {
		bgList := []string{}
		temp := listDirs(r.materialPath + "/backgrounds/" + r.BgSet)
		for _, v := range temp {
			if !strings.HasPrefix(v, ".") {
				bgList = append(bgList, r.materialPath+"/backgrounds/"+r.BgSet+"/"+v)
			}
		}
		background = bgList[r.hashArray[3]%int64(len(bgList))]
	}
	//roboparts[0]: material\sets\set1\white\003#01Body\006#white_body-05.png
	imgFile, err := os.Open(roboparts[0])
	if err != nil {
		return "", err
	}
	imgDec, _, err := image.Decode(imgFile)
	if err != nil {
		return "", err
	}
	defer imgFile.Close()
	resizeImg := resize.Resize(1024, 1024, imgDec, resize.Lanczos3)
	newImg := image.NewRGBA(image.Rect(0, 0, 1024, 1024))
	draw.Draw(newImg, newImg.Bounds(), resizeImg, resizeImg.Bounds().Min, draw.Over)
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
		draw.Draw(newImg, newImg.Bounds(), newImg, newImg.Bounds().Min, draw.Over)
		draw.Draw(newImg, newImg.Bounds(), tempResizeImg, newImg.Bounds().Min, draw.Over)
	}
	if r.BgSet != "" {
		imgFile, err := os.Open(background)
		if err != nil {
			return "", err
		}
		tempImg, _, err := image.Decode(imgFile)
		if err != nil {
			return "", err
		}
		imgFile.Close()
		tempResizeImg := resize.Resize(1024, 1024, tempImg, resize.Lanczos3)
		draw.Draw(newImg, newImg.Bounds(), tempResizeImg, tempImg.Bounds().Min, draw.Over)
	}
	output := fmt.Sprintf("%s/%s.%s", r.OutputPath, r.name, r.Format)
	out, err := os.Create(output)
	if err != nil {
		return "", err
	}
	defer out.Close()
	resizeImg = resize.Resize(uint(r.X), uint(r.Y), newImg, resize.Lanczos3)
	if err := jpeg.Encode(out, resizeImg, nil); err != nil {
		return "", err
	}
	return output, nil
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
