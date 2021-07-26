package pp

import (
	"fmt"
	"math/big"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"
	"unsafe"

	// Use fork until following PR is merged
	// https://github.com/mitchellh/colorstring/pull/3
	"github.com/k0kubun/colorstring"
)

type testCase struct {
	object interface{}
	expect string
}

type Foo struct {
	Bar       int
	Hoge      string
	Hello     map[string]string
	HogeHoges []HogeHoge
}

type FooPri struct {
	Public  string
	private string
}

type Piyo struct {
	Field1 map[string]string
	F2     *Foo
	Fie3   int
}

type HogeHoge struct {
	Hell  string
	World int
	A     interface{}
}

type EmptyStruct struct {
}

type User struct {
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	deletedAt time.Time
}

type LargeBuffer struct {
	Buf [1025]byte
}

type Private struct {
	b bool
	i int
	u uint
	f float32
	c complex128
}

type Circular struct {
	C *Circular
}

var c Circular = Circular{}
var nilSlice []int

func init() {
	c.C = &c
}

var (
	testCases = []testCase{
		{nil, "[cyan][bold]nil"},
		{nilSlice, "[][green]int[reset]([cyan][bold]nil[reset])"},
		{true, "[cyan][bold]true"},
		{false, "[cyan][bold]false"},
		{int(4), "[blue][bold]4"},
		{int8(8), "[blue][bold]8"},
		{int16(16), "[blue][bold]16"},
		{int32(32), "[blue][bold]32"},
		{int64(64), "[blue][bold]64"},
		{uint(4), "[blue][bold]4"},
		{uint8(8), "[blue][bold]8"},
		{uint16(16), "[blue][bold]16"},
		{uint32(32), "[blue][bold]32"},
		{uint64(64), "[blue][bold]64"},
		{uintptr(128), "[blue][bold]0x80"},
		{float32(2.23), "[magenta][bold]2.230000"},
		{float64(3.14), "[magenta][bold]3.140000"},
		{complex64(complex(3, -4)), "[blue][bold](3-4i)"},
		{complex128(complex(5, 6)), "[blue][bold](5+6i)"},
		{"string", `[red][bold]"[reset][red]string[reset][red][bold]"`},
		{[]string{}, "[][green]string[reset]{}"},
		{EmptyStruct{}, "pp.[green]EmptyStruct[reset]{}"},
		{
			[]*Piyo{nil, nil}, `
			[]*pp.[green]Piyo[reset]{
			  (*pp.[green]Piyo[reset])([cyan][bold]nil[reset]),
			  (*pp.[green]Piyo[reset])([cyan][bold]nil[reset]),
			}
			`,
		},
		{
			&c, `
				&pp.[green]Circular[reset]{
				  [yellow]C[reset]: &pp.[green]Circular[reset]{...},
				}
			`,
		},
		{"日本\t語\x00", `[red][bold]"[reset][red]日本[reset][magenta][bold]\t[reset][red]語[reset][magenta][bold]\x00[reset][red][bold]"`},
		{
			time.Date(2015, time.February, 14, 22, 15, 0, 0, time.UTC),
			"[blue][bold]2015[reset]-[blue][bold]02[reset]-[blue][bold]14[reset] [blue][bold]22[reset]:[blue][bold]15[reset]:[blue][bold]00[reset] [blue][bold]UTC[reset]",
		},
		{
			LargeBuffer{}, `
			pp.[green]LargeBuffer[reset]{
			  [yellow]Buf[reset]: [[blue]1025[reset]][green]uint8[reset]{...},
			}
			`,
		},
		{
			[]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, `
			[][green]uint8[reset]{
			  [blue][bold]0[reset], [blue][bold]1[reset], [blue][bold]2[reset], [blue][bold]3[reset], [blue][bold]4[reset], [blue][bold]5[reset], [blue][bold]6[reset], [blue][bold]7[reset], [blue][bold]8[reset], [blue][bold]9[reset], [blue][bold]0[reset], [blue][bold]1[reset], [blue][bold]2[reset], [blue][bold]3[reset], [blue][bold]4[reset], [blue][bold]5[reset],
			  [blue][bold]6[reset], [blue][bold]7[reset], [blue][bold]8[reset], [blue][bold]9[reset],
			}
			`,
		},
		{
			[]uint16{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, `
			[][green]uint16[reset]{
			  [blue][bold]0[reset], [blue][bold]1[reset], [blue][bold]2[reset], [blue][bold]3[reset], [blue][bold]4[reset], [blue][bold]5[reset], [blue][bold]6[reset], [blue][bold]7[reset],
			  [blue][bold]8[reset], [blue][bold]9[reset], [blue][bold]0[reset], [blue][bold]1[reset], [blue][bold]2[reset], [blue][bold]3[reset], [blue][bold]4[reset], [blue][bold]5[reset],
			  [blue][bold]6[reset], [blue][bold]7[reset], [blue][bold]8[reset], [blue][bold]9[reset],
			}
			`,
		},
		{
			[]uint32{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, `
			[][green]uint32[reset]{
			  [blue][bold]0[reset], [blue][bold]1[reset], [blue][bold]2[reset], [blue][bold]3[reset], [blue][bold]4[reset], [blue][bold]5[reset], [blue][bold]6[reset], [blue][bold]7[reset],
			  [blue][bold]8[reset], [blue][bold]9[reset], [blue][bold]0[reset], [blue][bold]1[reset], [blue][bold]2[reset], [blue][bold]3[reset], [blue][bold]4[reset], [blue][bold]5[reset],
			  [blue][bold]6[reset], [blue][bold]7[reset], [blue][bold]8[reset], [blue][bold]9[reset],
			}
			`,
		},
		{
			[]uint64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0}, `
			[][green]uint64[reset]{
			  [blue][bold]0[reset], [blue][bold]1[reset], [blue][bold]2[reset], [blue][bold]3[reset],
			  [blue][bold]4[reset], [blue][bold]5[reset], [blue][bold]6[reset], [blue][bold]7[reset],
			  [blue][bold]8[reset], [blue][bold]9[reset], [blue][bold]0[reset],
			}
			`,
		},
		{
			[][]byte{{0, 1, 2}, {3, 4}, {255}}, `
			[][green][]uint8[reset]{
			  [][green]uint8[reset]{
			    [blue][bold]0[reset], [blue][bold]1[reset], [blue][bold]2[reset],
			  },
			  [][green]uint8[reset]{
			    [blue][bold]3[reset], [blue][bold]4[reset],
			  },
			  [][green]uint8[reset]{
			    [blue][bold]255[reset],
			  },
			}
			`,
		},
		{
			map[string]interface{}{"foo": 10, "bar": map[int]int{20: 30}}, `
			[green]map[string]interface {}[reset]{
			  [red][bold]"[reset][red]bar[reset][red][bold]"[reset]: [green]map[int]int[reset]{
			    [blue][bold]20[reset]: [blue][bold]30[reset],
			  },
			  [red][bold]"[reset][red]foo[reset][red][bold]"[reset]: [blue][bold]10[reset],
			}
			`,
		},
	}

	arr [3]int
	tm  = time.Date(2015, time.January, 2, 0, 0, 0, 0, time.UTC)

	bigInt, _      = new(big.Int).SetString("-908f8474ea971baf", 16)
	bigFloat, _, _ = big.ParseFloat("3.1415926535897932384626433832795028", 10, 10, big.ToZero)

	checkCases = []interface{}{
		Private{b: false, i: 1, u: 2, f: 2.22, c: complex(5, 6)},
		map[string]int{"hell": 23, "world": 34},
		map[string]map[string]string{"s1": map[string]string{"v1": "m1", "va1": "me1"}, "si2": map[string]string{"v2": "m2"}},
		Foo{Bar: 1, Hoge: "a", Hello: map[string]string{"hel": "world", "a": "b"}, HogeHoges: []HogeHoge{HogeHoge{Hell: "a", World: 1}, HogeHoge{Hell: "bbb", World: 100}}},
		arr,
		[]string{"aaa", "bbb", "ccc"},
		make(chan bool, 10),
		func(a string, b float32) int { return 0 },
		&HogeHoge{},
		&Piyo{Field1: map[string]string{"a": "b", "cc": "dd"}, F2: &Foo{}, Fie3: 128},
		[]interface{}{1, 3},
		interface{}(1),
		HogeHoge{A: "test"},
		FooPri{Public: "hello", private: "world"},
		new(regexp.Regexp),
		unsafe.Pointer(new(regexp.Regexp)),
		"日本\t語\n\000\U00101234a",
		bigInt,
		bigFloat,
		&tm,
		&User{Name: "k0kubun", CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(), deletedAt: time.Now().UTC()},
	}
)

func TestFormat(t *testing.T) {
	for _, test := range testCases {
		actual := fmt.Sprintf("%s", defaultPrettyPrinter.format(test.object))

		trimmed := strings.Replace(test.expect, "\t", "", -1)
		trimmed = strings.TrimPrefix(trimmed, "\n")
		trimmed = strings.TrimSuffix(trimmed, "\n")
		expect := colorstring.Color(trimmed)
		if expect != actual {
			v := reflect.ValueOf(test.object)
			t.Errorf("\nTestCase: %#v\nType: %s\nExpect: %# v\nActual: %# v\n", test.object, v.Kind(), expect, actual)
			return
		}
		logResult(t, test.object, actual)
	}

	for _, object := range checkCases {
		actual := fmt.Sprintf("%s", defaultPrettyPrinter.format(object))
		logResult(t, object, actual)
	}
}

func logResult(t *testing.T, object interface{}, actual string) {
	if isMultiLine(actual) {
		t.Logf("%#v =>\n%s\n", object, actual)
	} else {
		t.Logf("%#v => %s\n", object, actual)
	}
}

func isMultiLine(text string) bool {
	return strings.Contains(text, "\n")
}
