package psh

import (
	"github.com/coocood/assrt"
	"testing"
)

func TestShConstruction(t *testing.T) {
	assert := assrt.NewAssert(t)

	echo := Sh("echo")

	assert.Equal(
		"echo",
		echo.expose().cmd,
	)
}

func TestShBakeArgs(t *testing.T) {
	assert := assrt.NewAssert(t)

	echo := Sh("echo")
	echo = echo.BakeArgs("a", "b")
	echo = echo.BakeArgs("c")

	assert.Equal(
		[]string{"a", "b", "c"},
		echo.expose().args,
	)
}

func TestShBakeArgsMagic(t *testing.T) {
	assert := assrt.NewAssert(t)

	echo := Sh("echo")("a", "b")("c")

	assert.Equal(
		[]string{"a", "b", "c"},
		echo.expose().args,
	)
}

func TestShBakeArgsForked(t *testing.T) {
	assert := assrt.NewAssert(t)

	echo := Sh("echo")
	echo1 := echo.BakeArgs("a", "b")
	echo2 := echo.BakeArgs("c")

	assert.Equal(
		0,
		len(echo.expose().args),
	)
	assert.Equal(
		[]string{"a", "b"},
		echo1.expose().args,
	)
	assert.Equal(
		[]string{"c"},
		echo2.expose().args,
	)
}

func TestShBakeArgsMagicForked(t *testing.T) {
	assert := assrt.NewAssert(t)

	echo := Sh("echo")
	echo1 := echo("a", "b")
	echo2 := echo("c")

	assert.Equal(
		0,
		len(echo.expose().args),
	)
	assert.Equal(
		[]string{"a", "b"},
		echo1.expose().args,
	)
	assert.Equal(
		[]string{"c"},
		echo2.expose().args,
	)
}

func TestShBakeArgsForkedDeeper(t *testing.T) {
	assert := assrt.NewAssert(t)

	echo := Sh("echo").BakeArgs("")
	echo1 := echo.BakeArgs("a", "b")
	echo2 := echo.BakeArgs("c")

	assert.Equal(
		[]string{""},
		echo.expose().args,
	)
	assert.Equal(
		[]string{"", "a", "b"},
		echo1.expose().args,
	)
	assert.Equal(
		[]string{"", "c"},
		echo2.expose().args,
	)
}

func TestShBakeArgsMagicForkedDeeper(t *testing.T) {
	assert := assrt.NewAssert(t)

	echo := Sh("echo")("")
	echo1 := echo("a", "b")
	echo2 := echo("c")

	assert.Equal(
		[]string{""},
		echo.expose().args,
	)
	assert.Equal(
		[]string{"", "a", "b"},
		echo1.expose().args,
	)
	assert.Equal(
		[]string{"", "c"},
		echo2.expose().args,
	)
}
