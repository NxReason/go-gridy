package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testArgs = [][]string{
	{},
	{"-i", "input.jpg"},
	{"-i", "input1.jpg", "input2.jpg", "-g", "25*25"},
	{"-i", "input1.jpg", "-g", "8x8", "-r", "1024*768"},
}

func TestFindAllFlagsEmptyArgs(t *testing.T) {
	args := testArgs[0]

	get := FindAllFlags(args)
	want := make(map[int]string)

	assert.Equal(t, want, get)
}

func TestFindAllFlagsInputOnly(t *testing.T) {
	args := testArgs[1]

	get := FindAllFlags(args)
	want := map[int]string { 0: "-i" }

	assert.Equal(t, want, get)
}

func TestFindAllFlagsMulti(t *testing.T) {
	args := testArgs[2]

	get := FindAllFlags(args)
	want := map[int]string { 0: "-i", 3: "-g" }

	assert.Equal(t, want, get)
}

func TestMakeConfigEmpty(t *testing.T) {
	flagsPos := FindAllFlags(testArgs[0])

	get, errs := MakeConfig(testArgs[0], flagsPos)
	want := Config{}

	assert.Equal(t, want, get)
	assert.Equal(t, 0, len(errs))
}

func TestMakeConfigSingleInput(t *testing.T) {
	args := testArgs[1]
	flagsPos := FindAllFlags(args)

	get, errs := MakeConfig(args, flagsPos)
	want := Config {
		InputFiles: []string { "input.jpg" },
	}

	assert.Equal(t, want, get)
	assert.Equal(t, 0, len(errs))
}

func TestMakeConfigMultiInput(t *testing.T) {
	assert := assert.New(t)
	args := testArgs[2]
	flagsPos := FindAllFlags(args)

	get, errs := MakeConfig(args, flagsPos)
	want := Config {
		InputFiles: []string {"input1.jpg", "input2.jpg"},
		GridRows: 25,
		GridCols: 25,
		gridSet: true,
	}

	assert.Equal(want, get)
	assert.Equal(0, len(errs))
}

func TestMakeConfig_InvalidArgs(t *testing.T) {
	assert := assert.New(t)
	args := []string{ "-i", "-g", "25i25" }
	flagsPos := FindAllFlags(args)

	got, errs := MakeConfig(args, flagsPos)
	want := Config{}

	assert.Equal(want, got)
	assert.Len(errs, 1)
	assert.ErrorContains(errs[0], "wrong grid format")
}

func TestMakeConfig_GridAndReso(t *testing.T) {
	assert := assert.New(t)
	args := testArgs[3]
	flagsPos := FindAllFlags(args)

	got, errs := MakeConfig(args, flagsPos)
	want := Config {
		InputFiles: []string{"input1.jpg"},
		GridRows: 8,
		GridCols: 8,
		gridSet: true,
		OutputWidth: 1024,
		OutputHeight: 768,
	}

	assert.Len(errs, 0)
	assert.Equal(want, got)
}

func TestMakeConfig_InputsFromFolder(t *testing.T) {
	assert := assert.New(t)
	args := []string{ "-f", "../fsmall", "-g", "10x10" }
	flagsPos := FindAllFlags(args)

	got, errs := MakeConfig(args, flagsPos)
	
	assert.Len(got.InputFiles, 4)
	assert.Len(errs, 0)
}