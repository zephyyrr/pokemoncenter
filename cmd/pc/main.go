package main

import (
	"encoding/binary"
	"encoding/json"
	"flag"
	"github.com/Sirupsen/logrus"
	"github.com/zephyyrr/pokemoncenter"
	"io"
	"os"
)

var (
	output    = flag.String("o", "a.pkm", "Output location")
	decompile = flag.Bool("dec", false, "Decompile instead")
)

func main() {
	flag.Parse()
	input := flag.Arg(0)

	r, err := os.Open(input)
	if err != nil {
		logrus.WithField("file", input).Fatal(err)
		return
	}
	w, err := os.OpenFile(*output, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		logrus.WithField("file", *output).Fatal(err)
		return
	}

	if *decompile {
		if *output == "a.pkm" {
			//Transform to a .json output instead.
			*output = "a.json"
		}
		err = Decompile(r, w)
	} else {
		//Just straight compiling
		err = Compile(r, w)
	}

	if err != nil {
		logrus.Fatal(err)
		return
	}
}

func Compile(r io.Reader, w io.Writer) (err error) {
	var pkm pokemoncenter.Pokémon
	dec := json.NewDecoder(r)
	err = dec.Decode(&pkm)
	if err != nil {
		return
	}
	err = WriteBinaryPokémon(pkm, w)
	return
}

func Decompile(r io.Reader, w io.Writer) (err error) {
	var pkm pokemoncenter.Pokémon
	err = binary.Read(r, binary.BigEndian, &pkm)
	if err != nil {
		return
	}
	err = WriteJsonDecompile(pkm, w)
	return
}

func WriteBinaryPokémon(pkm pokemoncenter.Pokémon, w io.Writer) error {
	return binary.Write(w, binary.BigEndian, pkm)
}

func WriteJsonDecompile(pkm pokemoncenter.Pokémon, w io.Writer) error {
	jsondata, err := json.MarshalIndent(pkm, "", "  ")
	_, err = w.Write(jsondata)
	return err
}
