package gorobohash

import (
	"crypto/sha512"
	"encoding/hex"
)

//resource define resource path
type resource struct {
	Sets      string
	BGSets    string
	Colors    string
	Format    string
	Hasharray []string
	HexDigest string
	Iter      int
	HashCount int
}

//newResource constructor
func newResource(name string) *resource {
	h := sha512.New()
	h.Write([]byte(name))
	return &resource{
		Sets:      "material/sets",
		BGSets:    "material/backgrounds",
		Colors:    "material/sets/set1",
		Format:    "png",
		Iter:      4,
		HashCount: 11,
		HexDigest: hex.EncodeToString(h.Sum(nil)),
	}
}

func (r *resource) createHahes(count int) {
	for i := 0; i < count; i++ {
		blockSize := len(r.HexDigest) / count
		currentStart := (1+i)*blockSize - blockSize
		currentEnd := (1 + i) * blockSize
		r.Hasharray = append(r.Hasharray, r.HexDigest[currentStart:currentEnd])
	}
}

func listDirs() {

}
