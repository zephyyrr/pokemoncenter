package main

import (
	"flag"
	"fmt"
	"github.com/Sirupsen/logrus"
	"io"
	"os"
)

const (
	PlayerNameOffset  = 0x2598
	PlayerNameSize    = 11
	PokédexOwnOffset  = 0x25A3
	PokédexOwnSize    = 19
	PokédexSeenOffset = 0x25B6
	PokédexSeenSize   = 19
	CheckSumOffset    = 0x3523
	CheckSumStart     = PlayerNameOffset
	CheckSumEnd       = 0x3522

	Terminator = 0x50
)

var (
	FullPokédex = []byte{
		0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0x7F}
)

func main() {
	flag.Parse()
	file := flag.Arg(0)
	savefile, err := os.OpenFile(file, os.O_RDWR, os.ModePerm)
	if err != nil {
		logrus.WithField("file", file).WithField("error", err).Error("Error encountered when opening the file.")
	}

	fmt.Println("Checksum Before:", RecalculateChecksum(savefile))

	savefile.Seek(PokédexOwnOffset, 0)
	savefile.Write(FullPokédex)

	savefile.Seek(PokédexSeenOffset, 0)
	savefile.Write(FullPokédex)

	checksum := RecalculateChecksum(savefile)
	savefile.WriteAt([]byte{checksum}, CheckSumOffset)
	fmt.Println("Checksum After:", checksum)

}

func RecalculateChecksum(sav io.ReadSeeker) (sum byte) {
	sav.Seek(CheckSumStart, 0)
	toRead := CheckSumEnd - CheckSumStart + 1 //Might be a off-by-one here.
	sum = 255
	buf := make([]byte, toRead)

	for toRead > 0 {
		buf = buf[0:toRead] //Reslice to only hold what we need.
		n, err := sav.Read(buf)
		toRead -= n

		for _, b := range buf {
			sum -= b
		}

		if err != nil {
			logrus.Error(err)
			return
		}
	}

	return
}
