package tkc

// build 202406180001

import (
	"archive/zip"
	"bufio"
	"bytes"
	"compress/gzip"
	"container/list"
	"context"
	"crypto"
	"crypto/aes"
	"crypto/md5"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha512"
	cryptorand "crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"encoding"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"encoding/xml"
	"errors"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"html/template"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	drawOri "image/draw"
	"io"
	"io/fs"
	"math"
	"math/big"
	"math/rand"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/smtp"
	"net/textproto"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
	"unicode"
	"unsafe"

	"golang.org/x/image/bmp"
	"golang.org/x/image/draw"
	"golang.org/x/term"

	"github.com/eiannone/keyboard"
	"github.com/topxeq/gods/lists/arraylist"
	"github.com/topxeq/gods/lists/doublylinkedlist"
	"github.com/topxeq/gods/sets/hashset"
	"github.com/topxeq/gods/sets/linkedhashset"
	"github.com/topxeq/gods/stacks/linkedliststack"
	"github.com/topxeq/gods/trees/btree"
	"github.com/topxeq/regexpx"
	"github.com/topxeq/xmlx"
	"github.com/topxeq/docxrepl"

	"github.com/aarzilli/sandblast"
	"github.com/topxeq/goph"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/net/html"
	
	"github.com/mholt/archives"

	"github.com/atotto/clipboard"
	"github.com/beevik/etree"
	jsoniter "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
	"github.com/tidwall/pretty"
	"github.com/topxeq/mahonia"
	"github.com/topxeq/socks"
	"github.com/topxeq/uuid"

	"github.com/mholt/archiver/v3"

	"github.com/jhillyerd/enmime"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"

	"github.com/yuin/goldmark"
	mdhighlighting "github.com/yuin/goldmark-highlighting"
	mdast "github.com/yuin/goldmark/ast"
	mdextension "github.com/yuin/goldmark/extension"
	mdparser "github.com/yuin/goldmark/parser"
	mdhtml "github.com/yuin/goldmark/renderer/html"

	zipx "github.com/yeka/zip"

	"github.com/mozillazg/go-pinyin"

	"github.com/davecgh/go-spew/spew"

	"github.com/expr-lang/expr"

	"unicode/utf8"

	"github.com/buger/jsonparser"
	"github.com/mailru/easyjson/jwriter"

	"github.com/pquerna/otp/totp"

	"github.com/ALTree/bigfloat"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/table"
	bubbletea "github.com/charmbracelet/bubbletea"
	
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/disk"
	
	"github.com/sergi/go-diff/diffmatchpatch"
	
	"github.com/adhocore/gronx"
	"github.com/adhocore/gronx/pkg/tasker"
	
	"github.com/fogleman/gg"
)

var VersionG = "v1.0.1"
var BuildG = "20230712"

type TK struct {
	Version string
}

var TKX = &TK{Version: VersionG}

func (pA *TK) NewTK() *TK {
	return &TK{Version: VersionG}
}

var NewTK = TKX.NewTK

func (pA *TK) GetVersion() string {
	return pA.Version
}

var GetVersion = TKX.GetVersion

// 类型 types structs

type ExitCallback func()

type Message struct {
	Type        string
	From        string
	To          string
	Value       interface{}
	BoolValue   bool
	IntValue    int
	FloatValue  float64
	StringValue string
	Resv1       interface{}
}

type UndefinedStruct struct {
	int
}

func (o UndefinedStruct) String() string {
	return "undefined"
}

var Undefined UndefinedStruct = UndefinedStruct{0}

func (pA *TK) IsUndefined(vA interface{}) bool {
	// Pl("--- %#v ----> %#v --- \n", vA, Undefined)
	return vA == Undefined
}

var IsUndefined = TKX.IsUndefined

// SimpleFlexObject
type SimpleFlexObject struct {
	IsValid  bool
	SimpleMode bool
	Items    []string
	ItemsMap map[string]int
	KeysMap  map[int]string
}

// example: |||name^^张三$$value1$$age^^16$$remark^^积分赛`+"\u2600\u2600"+`多凉快`+"\u2601\u2601"+`防hjkdjhsfdfsh`+"\u2602\u2602"+`kds1468273☀☀khfjh^*&68435静电☁☁连接☂☂方式$$value2$$回房间看到回复是大家课后反馈复合时空的|||
// will decode to object &tkc.SimpleFlexObject{IsValid:true, SimpleMode:false, Items:[]string{"张三", "value1", "16", "积分赛|||多凉快$$防hjkdjhsfdfsh^^kds1468273|||khfjh^*&68435静电$$连接^^方式", "value2", "回房间看到回复是大家课后反馈复合时空的"}
// ☀☀ or \u2600\u2600 for |||
// ☁☁ or \u2601\u2601 for $$
// ☂☂ or \u2602\u2602 for ^^

func (pA *TK) NewSimpleFlexObject(argsA ...interface{}) *SimpleFlexObject {
	objT := &SimpleFlexObject{Items: make([]string, 0, 10), ItemsMap: make(map[string]int, 0), KeysMap: make(map[int]string, 0), IsValid: true, SimpleMode: false}
	
	for _, v := range argsA {
		b, ok := v.(bool)
		if ok {
			objT.SimpleMode = b
			continue
		}
		
		s, ok := v.(string)
		if ok {
			objT.Decode(s, true)
			continue
		}
		
		o, ok := v.(*SimpleFlexObject)
		if ok {
			return o.Copy()
		}
		
		s1 := ToStr(v)
		objT.Decode(s1, true)
	}
	
	return objT
}

var NewSimpleFlexObject = TKX.NewSimpleFlexObject

func (pA *TK) NewSimpleFlexObjectSimple() *SimpleFlexObject {
	return &SimpleFlexObject{Items: make([]string, 0, 10), ItemsMap: make(map[string]int, 0), KeysMap: make(map[int]string, 0), IsValid: true, SimpleMode: true}
}

var NewSimpleFlexObjectSimple = TKX.NewSimpleFlexObjectSimple

func (pA *TK) NewSimpleFlexObjectWithString(strA string, mustA bool) *SimpleFlexObject {
	rs := &SimpleFlexObject{Items: make([]string, 0, 10), ItemsMap: make(map[string]int, 0), KeysMap: make(map[int]string, 0), IsValid: true, SimpleMode: false}

	rs.Decode(strA, mustA)

	return rs
}

var NewSimpleFlexObjectWithString = TKX.NewSimpleFlexObjectWithString

func (pA *TK) NewSimpleFlexObjectWithStringComplex(strA string, mustA bool) *SimpleFlexObject {
	rs := &SimpleFlexObject{Items: make([]string, 0, 10), ItemsMap: make(map[string]int, 0), KeysMap: make(map[int]string, 0), IsValid: true, SimpleMode: false}

	rs.Decode(strA, mustA)

	return rs
}

var NewSimpleFlexObjectWithStringComplex = TKX.NewSimpleFlexObjectWithStringComplex

func (p *SimpleFlexObject) Copy() *SimpleFlexObject {
	rs := &SimpleFlexObject{Items: make([]string, len(p.Items)), ItemsMap: make(map[string]int, len(p.ItemsMap)), KeysMap: make(map[int]string, len(p.KeysMap)), IsValid: p.IsValid, SimpleMode: p.SimpleMode}
	
	for i, v := range p.Items {
		rs.Items[i] = v
	}

	for k, v := range p.ItemsMap {
		rs.ItemsMap[k] = v
	}

	for k, v := range p.KeysMap {
		rs.KeysMap[k] = v
	}

	return rs
}

func (p *SimpleFlexObject) Decode(strA string, mustA bool) {
	listT := strings.Split(strA, "|||")

	var tmps0 string

	if len(listT) < 3 {
		if mustA {
			tmps0 = listT[0]
		} else {
			p.IsValid = false
			return
		}
	} else {
		tmps0 = listT[1]
	}

	list1T := strings.Split(tmps0, "$$")

	len1T := len(list1T)

	//	p.Items =  make([]string, len1T)
	//
	//	p.ItemsMap =

	baseLenT := len(p.Items)

	for i := 0; i < len1T; i++ {
		list2T := strings.Split(list1T[i], "^^")

		if len(list2T) > 1 {
			if p.SimpleMode {
				p.Items = append(p.Items, list2T[1])
				p.ItemsMap[list2T[0]] = baseLenT + i
				p.KeysMap[baseLenT+i] = list2T[0]
			} else {
				p.Items = append(p.Items, strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(list2T[1], "\u2600\u2600", "|||"), "\u2601\u2601", "$$"), "\u2602\u2602", "^^"))
				p.ItemsMap[list2T[0]] = baseLenT + i
				p.KeysMap[baseLenT+i] = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(list2T[0], "\u2600\u2600", "|||"), "\u2601\u2601", "$$"), "\u2602\u2602", "^^")
			}
		} else {
			p.Items = append(p.Items, list2T[0])
		}
	}

}

func (p *SimpleFlexObject) Encode(defaultA ...string) string {
	defaultT := ""

	if len(defaultA) > 0 {
		defaultT = defaultA[0]
	}

	if !p.IsValid {
		return defaultT
	}

	var sbufT strings.Builder

	len0T := len(p.Items)

	sbufT.WriteString("|||")

	for i := 0; i < len0T; i++ {
		if i != 0 {
			sbufT.WriteString("$$")
		}

		tmps1, b := p.KeysMap[i]

		if b {
			if p.SimpleMode {
				sbufT.WriteString(tmps1)
			} else {
				sbufT.WriteString(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(tmps1, "|||", "\u2600\u2600"), "$$", "\u2601\u2601"), "^^", "\u2602\u2602"))
			}
			
			sbufT.WriteString("^^")
		}

		if p.SimpleMode {
			sbufT.WriteString(p.Items[i])
		} else {
			sbufT.WriteString(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(p.Items[i], "|||", "\u2600\u2600"), "$$", "\u2601\u2601"), "^^", "\u2602\u2602"))
		}

	}

	sbufT.WriteString("|||")

	return sbufT.String()
}

func (p *SimpleFlexObject) GetItemWithKey(keyA string, defaultA ...string) string {
	defaultT := ""

	if len(defaultA) > 0 {
		defaultT = defaultA[0]
	}

	if !p.IsValid {
		return defaultT
	}

	tmpi, b := p.ItemsMap[keyA]

	if !b {
		return defaultT
	}

	if (tmpi < 0) || (tmpi >= len(p.Items)) {
		return defaultT
	}

	return p.Items[tmpi]

}

func (p *SimpleFlexObject) KeyExists(keyA string) bool {
	if !p.IsValid {
		return false
	}

	_, b := p.ItemsMap[keyA]

	return b
	
}

func (p *SimpleFlexObject) GetIndexByKey(keyA string) int {
	if !p.IsValid {
		return -1
	}

	idxT, b := p.ItemsMap[keyA]
	
	if !b {
		return -1
	}

	return idxT
	
}

func (p *SimpleFlexObject) GetItem(idxA int, defaultA ...string) string {
	defaultT := ""

	if len(defaultA) > 0 {
		defaultT = defaultA[0]
	}

	if !p.IsValid {
		return defaultT
	}

	lenT := len(p.Items)

	if (idxA < 0) || (idxA >= lenT) {
		return defaultT
	}
	
	return p.Items[idxA]

}

func (p *SimpleFlexObject) SetItemWithKey(keyA string, valueA string) string {
	if !p.IsValid {
		return "TXERROR:object is invalid"
	}

	tmpi, b := p.ItemsMap[keyA]

	if b {
		p.Items[tmpi] = valueA
		return ""
	}

	baseLenT := len(p.Items)

	p.Items = append(p.Items, valueA)
	p.ItemsMap[keyA] = baseLenT
	p.KeysMap[baseLenT] = keyA

	return ""

}

func (p *SimpleFlexObject) AddItemWithKey(keyA string, valueA string) string {
	if !p.IsValid {
		return "TXERROR:object is invalid"
	}

	_, b := p.ItemsMap[keyA]

	if b {
		return "TXERROR:key already exists"
	}

	baseLenT := len(p.Items)

	p.Items = append(p.Items, valueA)
	p.ItemsMap[keyA] = baseLenT
	p.KeysMap[baseLenT] = keyA

	return ""

}

func (p *SimpleFlexObject) AddItem(valueA string) string {
	if !p.IsValid {
		return "TXERROR:object is invalid"
	}

	p.Items = append(p.Items, valueA)

	return ""

}

func (p *SimpleFlexObject) InsertItem(valueA string, idxA int) string {
	if !p.IsValid {
		return "TXERROR:object is invalid"
	}
	
	lenT := len(p.Items)

	if idxA <= 0 {
		idxA = 0
		p.Items = append([]string{valueA}, p.Items...)
	} else if idxA >= lenT {
		idxA = lenT
		p.Items = append(p.Items, valueA)
	} else {
		p.Items = append(p.Items[0:idxA], append([]string{valueA}, p.Items[idxA:]...)...)
//		p.Items = append(tmpl, p.Items[idxA:]...)
	}
	
	ary1T := make([]int, 0, len(p.KeysMap))

	for idxT, _ := range p.KeysMap {
		if idxT >= idxA {
			ary1T = append(ary1T, idxT)
		}
	}
	
	sort.Sort(sort.Reverse(sort.IntSlice(ary1T)))
	
	for _, idxT := range ary1T {
		p.KeysMap[idxT + 1] = p.KeysMap[idxT]
		p.ItemsMap[p.KeysMap[idxT]] = idxT + 1
		delete(p.KeysMap, idxT)
	}

	return ""
}

func (p *SimpleFlexObject) SetItem(valueA string, idxA int) string {
	if !p.IsValid {
		return "TXERROR:object is invalid"
	}
	
	lenT := len(p.Items)

	if idxA < 0 || idxA >= lenT {
		return fmt.Sprintf("TXERROR:out of index (%v/%v)", idxA, lenT)
	}
	
	p.Items[idxA] = valueA

	return ""
}

func (p *SimpleFlexObject) SetItemByKey(valueA string, keyA string) string {
	if !p.IsValid {
		return "TXERROR:object is invalid"
	}
	
	idxT, keyOk := p.ItemsMap[keyA]
	
	if !keyOk {
		return "TXERROR:key not exists"
	}
	
	p.Items[idxT] = valueA

	return ""
}

func (p *SimpleFlexObject) InsertItemWithKey(keyA string, valueA string, idxA int) string {
	if !p.IsValid {
		return "TXERROR:object is invalid"
	}
	
	_, keyOk := p.ItemsMap[keyA]
	
	if keyOk {
		return "TXERROR:key already exists"
	}
	
	lenT := len(p.Items)

	if idxA <= 0 {
		idxA = 0
		p.Items = append([]string{valueA}, p.Items...)
	} else if idxA >= lenT {
		idxA = lenT
		p.Items = append(p.Items, valueA)
	} else {
		p.Items = append(p.Items[0:idxA], append([]string{valueA}, p.Items[idxA:]...)...)
//		p.Items = append(tmpl, p.Items[idxA:]...)
	}
	
	ary1T := make([]int, 0, len(p.KeysMap))

	for idxT, _ := range p.KeysMap {
		if idxT >= idxA {
			ary1T = append(ary1T, idxT)
		}
	}
	
	sort.Sort(sort.Reverse(sort.IntSlice(ary1T)))
	
	for _, idxT := range ary1T {
		p.KeysMap[idxT + 1] = p.KeysMap[idxT]
		p.ItemsMap[p.KeysMap[idxT]] = idxT + 1
		delete(p.KeysMap, idxT)
	}

	p.KeysMap[idxA] = keyA
	p.ItemsMap[keyA] = idxA

	return ""
}

func (p *SimpleFlexObject) RemoveItemByIndex(idxA int) string {
	if !p.IsValid {
		return "TXERROR:object is invalid"
	}
	
	lenT := len(p.Items)

	if (idxA < 0) || (idxA >= lenT) {
		return "TXERROR:index out of range"
	}
	
	if idxA == 0 {
		p.Items = p.Items[1:]
	} else if idxA == lenT - 1 {
		p.Items = p.Items[0:lenT - 1]
	} else {
		p.Items = append(p.Items[0:idxA], p.Items[idxA+1:]...)
	}

	keyT, ok := p.KeysMap[idxA]
	
	if ok {
		delete(p.KeysMap, idxA)
		delete(p.ItemsMap, keyT)
	}
	
	ary1T := make([]int, 0, len(p.KeysMap))
//	ary2T := make([]string, 0, len(p.KeysMap))

	for idxT, _ := range p.KeysMap {
		if idxT > idxA {
			ary1T = append(ary1T, idxT)
//			ary2T = append(ary2T, keyT)
		}
	}
	
	sort.Ints(ary1T)
//	Plv(ary1T)
//	Plv(ary2T)
	
	for _, idxT := range ary1T {
		p.KeysMap[idxT - 1] = p.KeysMap[idxT]
		p.ItemsMap[p.KeysMap[idxT]] = idxT - 1
		delete(p.KeysMap, idxT)
	}

	return ""
}

func (p *SimpleFlexObject) RemoveItemByKey(keyA string) string {
	if !p.IsValid {
		return "TXERROR:object is invalid"
	}
	
	idx0T, ok := p.ItemsMap[keyA]
	
	if !ok {
		return "TXERROR:key not exists"
	}

	lenT := len(p.Items)

	if (idx0T < 0) || (idx0T >= lenT) {
		return "TXERROR:index out of range"
	}
	
	if idx0T == 0 {
		p.Items = p.Items[1:]
	} else if idx0T == lenT - 1 {
		p.Items = p.Items[0:lenT - 1]
	} else {
		p.Items = append(p.Items[0:idx0T], p.Items[idx0T+1:]...)
	}

	keyT, ok := p.KeysMap[idx0T]
	
	if ok {
		delete(p.KeysMap, idx0T)
		delete(p.ItemsMap, keyT)
	}

	ary1T := make([]int, 0, len(p.KeysMap))

	for idxT, _ := range p.KeysMap {
		if idxT > idx0T {
			ary1T = append(ary1T, idxT)
		}
	}
	
	sort.Ints(ary1T)
	
	for _, idxT := range ary1T {
		p.KeysMap[idxT - 1] = p.KeysMap[idxT]
		p.ItemsMap[p.KeysMap[idxT]] = idxT - 1
		delete(p.KeysMap, idxT)
	}

	return ""
}

func (v SimpleFlexObject) Size() int {
	if !v.IsValid {
		return 0
	}

	return len(v.Items)
}

func (v SimpleFlexObject) KeySize() int {
	if !v.IsValid {
		return 0
	}

	return len(v.ItemsMap)
}

func (v SimpleFlexObject) GetMap() map[string]string {
	if !v.IsValid {
		return nil
	}

	mapT := make(map[string]string, len(v.ItemsMap))
	
	for k, kv := range v.ItemsMap {
		mapT[k] = v.Items[kv]
	}
	
	return mapT
}

func (v SimpleFlexObject) GetIndexMap() map[string]string {
	if !v.IsValid {
		return nil
	}

	mapT := make(map[string]string, len(v.Items))
	
	for i, iv := range v.Items {
		mapT[IntToStr(i)] = iv
	}
	
	return mapT
}

func (v SimpleFlexObject) GetItems() []string {
	if !v.IsValid {
		return nil
	}

	aryT := make([]string, 0, len(v.Items))
	
	for _, v := range v.Items {
		aryT = append(aryT, v)
	}
	
	return aryT
}

func (v SimpleFlexObject) GetItemsWithoutKey() []string {
	if !v.IsValid {
		return nil
	}

	aryT := make([]string, 0, len(v.Items))
	
	for i, iv := range v.Items {
		_, ok := v.KeysMap[i]
		
		if !ok {
			aryT = append(aryT, iv)
		}
	}
	
	return aryT
}

func (v SimpleFlexObject) GetItemsWithKey() []string {
	if !v.IsValid {
		return nil
	}

	aryT := make([]string, 0, len(v.Items))
	
	for i, iv := range v.Items {
		_, ok := v.KeysMap[i]
		
		if ok {
			aryT = append(aryT, iv)
		}
	}
	
	return aryT
}

func (v SimpleFlexObject) GetKeys() []string {
	if !v.IsValid {
		return nil
	}

	aryT := make([]string, 0, len(v.ItemsMap))
	
	for k, _ := range v.ItemsMap {
		aryT = append(aryT, k)
	}
	
	return aryT
}

func (p SimpleFlexObject) String() string {
	if !p.IsValid {
		return "simpleFlexObject(invalid)"
	}

	var sb strings.Builder

	sb.WriteString("{")

	i := 0
	for tmps, v := range p.ItemsMap {
		if i > 0 {
			sb.WriteString(", ")
		}

		sb.WriteString(tmps)
		sb.WriteString(": ")
		sb.WriteString(IntToStr(v))

		i++
	}

	sb.WriteString("}")

	sb.WriteString("{")

	i = 0
	for tmpi, v := range p.KeysMap {
		if i > 0 {
			sb.WriteString(", ")
		}

		sb.WriteString(IntToStr(tmpi))
		sb.WriteString(": ")
		sb.WriteString(v)

		i++
	}

	sb.WriteString("}")

	sb.WriteString("[")

	i = 0
	for _, v := range p.Items {
		if i > 0 {
			sb.WriteString(", ")
		}

		sb.WriteString(strconv.Quote(v))
		i++
	}

	sb.WriteString("]")

	return sb.String()
}

// IntSumBuffer
type IntSumBuffer struct {
	Size    int
	Data    []int
	Pointer int
}

func (pA *TK) NewIntSumBuffer(sizeA int) *IntSumBuffer {
	return &IntSumBuffer{Data: make([]int, sizeA), Size: sizeA, Pointer: 0}
}

var NewIntSumBuffer = TKX.NewIntSumBuffer

func (p *IntSumBuffer) Add(vA int) (int, bool) {
	p.Data[p.Pointer] = vA

	p.Pointer++

	if p.Pointer >= p.Size {
		var sumT int = 0

		for i := 0; i < p.Size; i++ {
			sumT += p.Data[i]
		}

		p.Pointer = 0
		return sumT, true
	}

	return vA, false
}

func (p *IntSumBuffer) AddCompact(vA int) interface{} {
	p.Data[p.Pointer] = vA

	p.Pointer++

	if p.Pointer >= p.Size {
		var sumT int = 0

		for i := 0; i < p.Size; i++ {
			sumT += p.Data[i]
		}

		p.Pointer = 0
		return sumT
	}

	return false
}

func (p *IntSumBuffer) AddAvg(vA int) (int, bool) {
	p.Data[p.Pointer] = vA

	p.Pointer++

	if p.Pointer >= p.Size {
		var sumT int = 0

		for i := 0; i < p.Size; i++ {
			sumT += p.Data[i]
		}

		p.Pointer = 0
		return sumT / p.Size, true
	}

	return vA, false
}

func (p *IntSumBuffer) AddAvgCompact(vA int) interface{} {
	p.Data[p.Pointer] = vA

	p.Pointer++

	if p.Pointer >= p.Size {
		var sumT int = 0

		for i := 0; i < p.Size; i++ {
			sumT += p.Data[i]
		}

		p.Pointer = 0
		return sumT / p.Size
	}

	return false
}

// MapRef
type MapRef struct {
	Data interface{}
	Key  string
}

func (pA *TK) GetMapRef(mapA interface{}, keyA string) *MapRef {
	valueT := reflect.ValueOf(mapA)

	kindT := valueT.Kind()

	if kindT != reflect.Map {
		return nil
	}

	return &MapRef{Data: mapA, Key: keyA}
}

var GetMapRef = TKX.GetMapRef

func (p *MapRef) GetValue() interface{} {
	valueT := reflect.ValueOf(p.Data)

	kindT := valueT.Kind()

	// if kindT != reflect.Pointer {
	// 	return fmt.Errorf("not pointer type")
	// }

	// elemT := valueT.Elem()

	if kindT != reflect.Map {
		return nil
	}

	rv := valueT.MapIndex(reflect.ValueOf(p.Key))

	if !rv.IsValid() {
		return nil
	}

	return rv.Interface()

}

func (p *MapRef) SetValue(vA interface{}) (result error) {
	defer func() {
		r := recover()

		if r != nil {
			result = fmt.Errorf("failed to set map value: %v", r)
			return
		}
	}()

	valueT := reflect.ValueOf(p.Data)

	kindT := valueT.Kind()

	// if kindT != reflect.Pointer {
	// 	return fmt.Errorf("not pointer type")
	// }

	// elemT := valueT.Elem()

	if kindT != reflect.Map {
		result = fmt.Errorf("not a map")
		return
	}

	valueT.SetMapIndex(reflect.ValueOf(p.Key), reflect.ValueOf(vA))

	result = nil
	return
}

type FlexRef struct {
	Type  string
	Data  interface{}
	Key   string
	Index int
}

func (pA *TK) GetFlexRef(dataA interface{}, typeA string, keyA string, indexA int) *FlexRef {
	if typeA == "" {
		return &FlexRef{Type: "", Data: dataA, Key: keyA, Index: indexA}
	}

	if typeA == "map" {
		valueT := reflect.ValueOf(dataA)

		kindT := valueT.Kind()

		if kindT != reflect.Map {
			return nil
		}

		return &FlexRef{Type: "map", Data: dataA, Key: keyA, Index: indexA}

	}

	if typeA == "array" || typeA == "slice" {
		valueT := reflect.ValueOf(dataA)

		kindT := valueT.Kind()

		if kindT != reflect.Slice && kindT != reflect.Slice {
			return nil
		}

		if kindT == reflect.Slice {
			return &FlexRef{Type: "array", Data: dataA, Key: keyA, Index: indexA}
		}

		if kindT == reflect.Array {
			return &FlexRef{Type: "array", Data: dataA, Key: keyA, Index: indexA}
		}

		return nil

	}

	return nil
}

var GetFlexRef = TKX.GetFlexRef

func (p *FlexRef) GetValue() (result interface{}) {
	defer func() {
		r := recover()

		if r != nil {
			result = nil
			return
		}
	}()

	valueT := reflect.ValueOf(p.Data)

	kindT := valueT.Kind()

	if p.Type == "map" {

		// if kindT != reflect.Pointer {
		// 	return fmt.Errorf("not pointer type")
		// }

		// elemT := valueT.Elem()

		if kindT != reflect.Map {
			return nil
		}

		rv := valueT.MapIndex(reflect.ValueOf(p.Key))

		if !rv.IsValid() {
			return nil
		}

		return rv.Interface()

	}

	if p.Type == "array" {
		// if kindT != reflect.Pointer {
		// 	return fmt.Errorf("not pointer type")
		// }

		// elemT := valueT.Elem()

		if kindT != reflect.Slice && kindT != reflect.Array {
			return nil
		}

		rv := valueT.Index(p.Index)

		if !rv.IsValid() {
			return nil
		}

		return rv.Interface()

	}

	if kindT != reflect.Pointer {
		return nil
	}

	rv := valueT.Elem()

	if !rv.IsValid() {
		return nil
	}

	return rv
}

func (p *FlexRef) SetValue(vA interface{}) (result error) {
	defer func() {
		r := recover()

		if r != nil {
			result = fmt.Errorf("failed to set map value: %v", r)
			return
		}
	}()

	valueT := reflect.ValueOf(p.Data)

	kindT := valueT.Kind()

	// Pl("SetValue: %#v %#v %v", p.Data, vA, kindT)

	// if kindT != reflect.Pointer {
	// 	return fmt.Errorf("not pointer type")
	// }

	// elemT := valueT.Elem()

	if kindT == reflect.Map {
		valueT.SetMapIndex(reflect.ValueOf(p.Key), reflect.ValueOf(vA))
		return nil
	}

	if kindT == reflect.Array || kindT == reflect.Slice {
		valueT.Index(p.Index).Set(reflect.ValueOf(vA))
		return nil
	}

	result = fmt.Errorf("type not supported")

	return
}

type List struct {
	Data []interface{}

	Type string
}

func NewList(vA ...interface{}) *List {
	p := &List{Data: make([]interface{}, 0, len(vA)), Type: ""}
	p.Data = append(p.Data, vA...)

	return p
}

func (p *List) Clear() {
	p.Data = make([]interface{}, 0)
}

func (p *List) AddItem(vA interface{}) {
	p.Data = append(p.Data, vA)
}

func (p *List) AddItems(vA ...interface{}) {
	p.Data = append(p.Data, vA...)
}

type QuickObject struct {
	Type  string
	Value interface{}
	Id    string
	Name  string
}

func (pA *TK) NewQuickObject(argsA ...interface{}) *QuickObject {
	p1 := &QuickObject{}

	for _, v := range argsA {
		s1, ok := v.(string)

		if ok {
			if strings.HasPrefix(s1, "-type=") {
				p1.Type = s1[6:]
			} else if strings.HasPrefix(s1, "-id=") {
				p1.Id = s1[4:]
			} else if strings.HasPrefix(s1, "-name=") {
				p1.Name = s1[6:]
			} else if strings.HasPrefix(s1, "-value=") {
				p1.Value = s1[7:]
			} else {
				p1.Value = s1
			}
		} else {
			p1.Value = v
		}
	}

	return p1
}

var NewQuickObject = TKX.NewQuickObject

// 自定义操作代理类型

type TXDelegate func(actionA string, objA interface{}, dataA interface{}, paramsA ...interface{}) interface{}

type CompactDelegate func(inputA interface{}) interface{}

type QuickDelegate func(strA string) string

type QuickVarDelegate func(argsA ...interface{}) interface{}

// 自定义集合类型

type TXCollection struct {
	Content map[string]int
}

func (pA *TK) CreateTXCollection(sizeA int) *TXCollection {
	rs := &TXCollection{}
	rs.Content = make(map[string]int, sizeA)

	return rs
}

var CreateTXCollection = TKX.CreateTXCollection

func (p *TXCollection) InitX(sizeA int) {
	p.Content = make(map[string]int, sizeA)
}

func (p *TXCollection) Length() int {
	if p.Content == nil {
		return 0
	}

	return len(p.Content)
}

func (p *TXCollection) Size() int {
	return p.Length()
}

func (p *TXCollection) AddOrSet(strA string) {
	p.Content[strA] = 1
}

func (p *TXCollection) AddOrSetExcludeEmpty(strA string) {
	if strA == "" {
		return
	}

	p.Content[strA] = 1
}

func (p *TXCollection) Contains(strA string) bool {
	_, ok := p.Content[strA]

	return ok
}

func (p *TXCollection) Remove(strA string) bool {
	_, ok := p.Content[strA]

	if ok {
		delete(p.Content, strA)
		return true
	}

	return false
}

func (p *TXCollection) GetList() []string {
	if (p == nil) || (p.Content == nil) {
		return make([]string, 0)
	}

	rs := make([]string, 0, len(p.Content))

	for k := range p.Content {
		rs = append(rs, k)
	}

	return rs
}

func (p *TXCollection) GetSortedList(ifUpA bool) []string {
	rs := p.GetList()

	if ifUpA {
		sort.Sort(sort.StringSlice(rs))
	} else {
		sort.Sort(sort.Reverse(sort.StringSlice(rs)))
	}

	return rs
}

func (p *TXCollection) GetListString(ifUpA bool, sepA string) string {
	rs := p.GetList()

	if ifUpA {
		sort.Sort(sort.StringSlice(rs))
	} else {
		sort.Sort(sort.Reverse(sort.StringSlice(rs)))
	}

	return strings.Join(rs, sepA)
}

// 存放TX格式的网络API返回结果
type TXResult map[string]string // struct {

// Status  string
// Value   string
// Debuf   string
// Reason  string
// Object  string
// Object2 string
// Object3 string
// Token   string
//}

var InvalidTXResultG = TXResult{"Status": "fail", "Value": "general error"}

func (vA TXResult) String() string {
	return ToJSONX(vA)
}

func (vA TXResult) Wrap() string {
	return "TX_RESULT: " + ToJSONX(vA) + "_XT"
}

func (vA TXResult) Get(keyA string) string {
	return vA[keyA]
}

func (pA *TXResult) Set(keyA string, valueA string) {
	(*pA)[keyA] = valueA
}

func (pA *TXResult) Remove(keyA string) {
	delete(*pA, keyA)
}

func (pA *TXResult) FromJSON(jsonA string) interface{} {
	vT := NewTXResultFromJSON(jsonA)

	if IsError(vT) {
		return vT
	}

	// nv1 := vT.(*TXResult)

	// for k, v = range nv1 {
	// 	pA[k] = v
	// }

	*pA = *(vT.(*TXResult))

	return pA
}

func (pA *TK) NewTXResult(statusA string, valueA ...string) *TXResult {
	var valueT = ""

	if len(valueA) > 0 {
		valueT = valueA[0]
	}

	pT := &TXResult{"Status": statusA, "Value": valueT}

	return pT
}

var NewTXResult = TKX.NewTXResult

func (pA *TK) NewTXResultFromJSON(jsonA string) interface{} {
	var rs = new(TXResult)

	errT := jsoniter.Unmarshal([]byte(jsonA), rs)

	if errT != nil {
		return errT
	}

	return rs
}

var NewTXResultFromJSON = TKX.NewTXResultFromJSON

func (pA *TK) NewTXResultFromJSONX(jsonA string) *TXResult {
	var rs = new(TXResult)

	errT := jsoniter.Unmarshal([]byte(jsonA), rs)

	if errT != nil {
		return NewTXResult("fail", "failed to parse TXResult: "+errT.Error())
	}

	return rs
}

var NewTXResultFromJSONX = TKX.NewTXResultFromJSONX

func (pA *TK) NewTXResultFromWrap(strA string) interface{} {
	strT := RegFindFirstX(strA, `(?sm)TX_RESULT: ({.*?"Status".*?)_XT`, 1)

	if IsErrorString(strT) {
		return ErrStrToErr(strT)
	}

	var rs = NewTXResultFromJSON(strT)

	return rs
}

var NewTXResultFromWrap = TKX.NewTXResultFromWrap

func (pA *TK) NewTXResultFromWrapX(strA string) *TXResult {
	strT := RegFindFirstX(strA, `(?sm)TX_RESULT: ({.*?"Status".*?)_XT`, 1)

	if IsErrorString(strT) {
		return NewTXResult("fail", "failed to parse TXResult wrapped: "+GetErrStr(strT))
	}

	var rs = NewTXResultFromJSONX(strT)

	return rs
}

var NewTXResultFromWrapX = TKX.NewTXResultFromWrapX

func (pA *TK) TXResultFromStringE(strA string) (*TXResult, error) {
	pT := new(TXResult)

	errT := json.Unmarshal([]byte(strA), pT)

	if errT != nil {
		return nil, errT
	}

	return pT, nil
}

var TXResultFromStringE = TKX.TXResultFromStringE

func (pA *TK) TXResultFromString(strA string) *TXResult {
	pT := new(TXResult)

	errT := json.Unmarshal([]byte(strA), pT)

	if errT != nil {
		return nil
	}

	return pT
}

var TXResultFromString = TKX.TXResultFromString

func (pA *TK) TXResultFromStringSafely(strA string) *TXResult {
	p := new(TXResult)

	errT := json.Unmarshal([]byte(strA), p)

	if errT != nil {
		return &InvalidTXResultG
	}

	return p
}

var TXResultFromStringSafely = TKX.TXResultFromStringSafely

func (pA *TK) Resultf(statusA string, formatA string, argsA ...interface{}) *TXResult {
	return NewTXResult(statusA, fmt.Sprintf(formatA, argsA...))
}

var Resultf = TKX.Resultf

// 全局变量 global variables

var logFileG = "c:\\tk.log"
var TimeFormat = "2006-01-02 15:04:05"
var TimeFormatMS = "2006-01-02 15:04:05.000"
var TimeFormatMSCompact = "20060102150405.000"
var TimeFormatCompact = "20060102150405"
var TimeFormatCompact2 = "2006/01/02 15:04:05"
var TimeFormatDateCompact = "20060102"
var GlobalEnvSetG *TXCollection = nil

var variableG = make(map[string]interface{})

var varMutexG sync.Mutex
var fileVarMutexG sync.Mutex

// global locks, 0 - general purpose, 1 - lock for file seq generator
var LocksG [10]sync.RWMutex

// 全局环境集合相关 global environment related

func (pA *TK) SetGlobalEnv(vA string) {
	if GlobalEnvSetG == nil {
		GlobalEnvSetG = CreateTXCollection(100)
	}

	GlobalEnvSetG.AddOrSetExcludeEmpty(vA)
}

var SetGlobalEnv = TKX.SetGlobalEnv

func (pA *TK) RemoveGlobalEnv(vA string) {
	if GlobalEnvSetG == nil {
		GlobalEnvSetG = CreateTXCollection(100)
	}

	GlobalEnvSetG.Remove(vA)
}

var RemoveGlobalEnv = TKX.RemoveGlobalEnv

func (pA *TK) GetGlobalEnvList() []string {
	if GlobalEnvSetG == nil {
		GlobalEnvSetG = CreateTXCollection(100)
	}

	return GlobalEnvSetG.GetList()
}

var GetGlobalEnvList = TKX.GetGlobalEnvList

func (pA *TK) GetGlobalEnvString() string {
	if GlobalEnvSetG == nil {
		GlobalEnvSetG = CreateTXCollection(100)
	}

	return GlobalEnvSetG.GetListString(true, ",")
}

var GetGlobalEnvString = TKX.GetGlobalEnvString

func (pA *TK) HasGlobalEnv(vA string) bool {
	if GlobalEnvSetG == nil {
		GlobalEnvSetG = CreateTXCollection(100)
	}

	return GlobalEnvSetG.Contains(vA)
}

var HasGlobalEnv = TKX.HasGlobalEnv

// Global variables related 全局变量相关

func (pA *TK) LockN(idxA ...int) {
	idxT := 0

	if len(idxA) > 0 {
		idxT = idxA[0]
	}

	if idxT < 0 || idxT > 9 {
		idxT = 0
	}

	LocksG[idxT].Lock()
}

var LockN = TKX.LockN

func (pA *TK) UnlockN(idxA ...int) {
	// defer func() {
	// 	Pl("defer: %v", "r")
	// 	r := recover()

	// 	if r != nil {
	// 		Pl("recover: %v", r)
	// 		return
	// 	}
	// }()

	idxT := 0

	if len(idxA) > 0 {
		idxT = idxA[0]
	}

	if idxT < 0 || idxT > 9 {
		idxT = 0
	}

	LocksG[idxT].Unlock()
}

var UnlockN = TKX.UnlockN

func (pA *TK) TryLockN(idxA ...int) bool {
	idxT := 0

	if len(idxA) > 0 {
		idxT = idxA[0]
	}

	if idxT < 0 || idxT > 9 {
		idxT = 0
	}

	return LocksG[idxT].TryLock()
}

var TryLockN = TKX.TryLockN

func (pA *TK) RLockN(idxA ...int) {
	idxT := 0

	if len(idxA) > 0 {
		idxT = idxA[0]
	}

	if idxT < 0 || idxT > 9 {
		idxT = 0
	}

	LocksG[idxT].RLock()
}

var RLockN = TKX.RLockN

func (pA *TK) RUnlockN(idxA ...int) {
	// defer func() {
	// 	r := recover()

	// 	if r != nil {
	// 		return
	// 	}
	// }()

	idxT := 0

	if len(idxA) > 0 {
		idxT = idxA[0]
	}

	if idxT < 0 || idxT > 9 {
		idxT = 0
	}

	LocksG[idxT].RUnlock()
}

var RUnlockN = TKX.RUnlockN

func (pA *TK) TryRLockN(idxA ...int) bool {
	idxT := 0

	if len(idxA) > 0 {
		idxT = idxA[0]
	}

	if idxT < 0 || idxT > 9 {
		idxT = 0
	}

	return LocksG[idxT].TryRLock()
}

var TryRLockN = TKX.TryRLockN

func (pA *TK) GetVar(nameA string) interface{} {
	varMutexG.Lock()
	rs, ok := variableG[nameA]
	varMutexG.Unlock()

	if !ok {
		GenerateErrorString("no key")
	}

	return rs
}

var GetVar = TKX.GetVar

func (pA *TK) GetVarEx(nameA string, argsA ...interface{}) interface{} {
	var defaultT interface{} = nil

	if len(argsA) > 0 {
		defaultT = argsA[0]
	}

	varMutexG.Lock()
	rs, ok := variableG[nameA]
	varMutexG.Unlock()

	if !ok {
		return defaultT
	}

	return rs
}

var GetVarEx = TKX.GetVarEx

func (pA *TK) SetVar(nameA string, valueA interface{}) {
	varMutexG.Lock()
	variableG[nameA] = valueA
	varMutexG.Unlock()
}

var SetVar = TKX.SetVar

func (pA *TK) DeleteVar(nameA string) {
	varMutexG.Lock()
	delete(variableG, nameA)
	varMutexG.Unlock()
}

var DeleteVar = TKX.DeleteVar

func (pA *TK) ClearVar() {
	varMutexG.Lock()
	variableG = make(map[string]interface{})
	varMutexG.Unlock()
}

var ClearVar = TKX.ClearVar

func (pA *TK) SizeVar() int {
	var rs int
	varMutexG.Lock()
	rs = len(variableG)
	varMutexG.Unlock()

	return rs
}

var SizeVar = TKX.SizeVar

func (pA *TK) GetFileVar(fileNameA string) interface{} {
	var rs interface{}
	fileVarMutexG.Lock()
	errT := LoadJSONFromFile(fileNameA, &rs)
	fileVarMutexG.Unlock()

	if errT != nil {
		return errT
	}

	return rs
}

var GetFileVar = TKX.GetFileVar

func (pA *TK) SetFileVar(fileNameA string, valueA interface{}) error {
	fileVarMutexG.Lock()
	errT := SaveJSONIndentToFile(valueA, fileNameA)
	fileVarMutexG.Unlock()

	return errT
}

var SetFileVar = TKX.SetFileVar

// 字符串相关函数 string related

// func (pA *TK) FindFirstDiffPosInStrs(str1A string, str2A string) int {
// 	len1 := len(str1A)
// 	len2 := len(str2A)

// 	var lenT int

// 	if len1 > len2 {
// 		lenT = len2
// 	} else {
// 		lenT = len1
// 	}

// 	i := 0

// 	for ; i <= lenT; i++ {
// 		if str1A[i] != str2A[i] {
// 			return i
// 		}
// 	}

// 	return i

// }

// var FindFirstDiffPosInStrs = TKX.FindFirstDiffPosInStrs

func (pA *TK) LimitString(strA string, lenA int, optsA ...string) string {
	if lenA < 0 {
		return strA
	}

	suffixT := GetSwitch(optsA, "-suffix=", "...")
	fromEndT := IfSwitchExists(optsA, "-end")

	runeT := IfSwitchExists(optsA, "-rune")

	if runeT {
		runeStrT := []rune(strA)

		var sb strings.Builder

		if !fromEndT {
			for i, v := range runeStrT {
				sb.WriteRune(v)

				if sb.Len() > lenA {
					return string(runeStrT[:i]) + suffixT
				}
			}

			return strA + suffixT
		} else {
			lenT := len(runeStrT)
			for i := lenT - 1; i >= 0; i-- {
				sb.WriteRune(runeStrT[i])

				if sb.Len() > lenA {
					return string(runeStrT[i+1:]) + suffixT
				}
			}

			return strA + suffixT
		}
	}

	lenT := len(strA)

	diffT := lenT - lenA

	if diffT <= 0 {
		return strA
	}

	if fromEndT {
		return strA[lenT-lenA:lenT] + suffixT
	}

	return strA[:lenA] + suffixT

}

var LimitString = TKX.LimitString

func (pA *TK) PadString(strA string, lenA int, optsA ...string) string {
	fillStrT := GetSwitch(optsA, "-fill=", "0")

	ifRightT := IfSwitchExists(optsA, "-right") || (GetSwitch(optsA, "-right=", "") == "true")

	lenT := len(strA)

	diffT := lenA - lenT

	if diffT <= 0 {
		return strA
	}

	len2T := len(fillStrT)

	if ifRightT {
		return strA + strings.Repeat(fillStrT, diffT/len2T)
	}

	return strings.Repeat(fillStrT, diffT/len2T) + strA

}

var PadString = TKX.PadString

func (pA *TK) FindSubStringAll(strA string, subStrA string) [][]int {
	bufT := make([][]int, 0, 100)

	subLenT := len(subStrA)

	for start := 0; start < len(strA); {
		remaining := strA[start:] // slicing the string is cheap
		loc := strings.Index(remaining, subStrA)
		if loc < 0 {
			return bufT
		}
		// loc[0] is the start index of the match,
		// loc[1] is the end index (exclusive)

		bufT = append(bufT, []int{loc + start, loc + subLenT + start})

		start += loc + subLenT
	}

	return bufT
}

var FindSubStringAll = TKX.FindSubStringAll

func (pA *TK) IsEmptyTrim(strA string) bool {
	return (Trim(strA) == "")
}

var IsEmptyTrim = TKX.IsEmptyTrim

func (pA *TK) RuneLen(strA string) int {
	return len([]rune(strA))
}

var RuneLen = TKX.RuneLen

func (pA *TK) JoinList(aryA interface{}, sepA string, defaultA ...string) string {
	var defaultT string = ""
	if (defaultA != nil) && (len(defaultA) > 0) {
		defaultT = defaultA[0]
	}

	if aryA == nil {
		return defaultT
	}

	switch v := aryA.(type) {
	case []string:
		return strings.Join(v, sepA)
	case []interface{}:
		var bufT strings.Builder
		for j, jv := range v {
			if j > 0 {
				bufT.WriteString(sepA)
			}

			bufT.WriteString(fmt.Sprintf("%v", jv))
		}

		return bufT.String()
	case []int:
		var bufT strings.Builder
		for j, jv := range v {
			if j > 0 {
				bufT.WriteString(sepA)
			}

			bufT.WriteString(fmt.Sprintf("%v", jv))
		}

		return bufT.String()
	case []int64:
		var bufT strings.Builder
		for j, jv := range v {
			if j > 0 {
				bufT.WriteString(sepA)
			}

			bufT.WriteString(fmt.Sprintf("%v", jv))
		}

		return bufT.String()
	case []float64:
		var bufT strings.Builder
		for j, jv := range v {
			if j > 0 {
				bufT.WriteString(sepA)
			}

			bufT.WriteString(fmt.Sprintf("%v", jv))
		}

		return bufT.String()
	case []byte:
		var bufT strings.Builder
		for j, jv := range v {
			if j > 0 {
				bufT.WriteString(sepA)
			}

			bufT.WriteString(fmt.Sprintf("%X", jv))
		}

		return bufT.String()
	case map[string]string:
		var bufT strings.Builder
		j := 0
		for k, kv := range v {
			if j > 0 {
				bufT.WriteString(sepA)
			}

			bufT.WriteString(fmt.Sprintf("%v: %v", k, kv))
		}

		return bufT.String()
	case map[string]interface{}:
		var bufT strings.Builder
		j := 0
		for k, kv := range v {
			if j > 0 {
				bufT.WriteString(sepA)
			}

			bufT.WriteString(fmt.Sprintf("%v: %v", k, kv))
		}

		return bufT.String()
	}

	return defaultT
}

var JoinList = TKX.JoinList

// StartsWith 检查字符串strA开始是否是subStrA
func (pA *TK) StartsWith(strA string, subStrA ...string) bool {
	for _, v := range subStrA {
		if strings.HasPrefix(strA, v) {
			return true
		}
	}

	return false
}

var StartsWith = TKX.StartsWith

func (pA *TK) StartsWithIgnoreCase(strA string, subStrA string) bool {

	return strings.HasPrefix(strings.ToLower(strA), strings.ToLower(subStrA))
}

var StartsWithIgnoreCase = TKX.StartsWithIgnoreCase

func (pA *TK) StartsWithUpper(wordA string) bool {
	if len(wordA) < 1 {
		return false
	}

	return (wordA[0] >= 'A') && (wordA[0] <= 'Z')
}

var StartsWithUpper = TKX.StartsWithUpper

func (pA *TK) StartsWithDigit(strA string) bool {
	if len(strA) < 1 {
		return false
	}

	ts := strA[0:1]
	switch ts {
	case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
		return true
	default:
		return false

	}
}

var StartsWithDigit = TKX.StartsWithDigit

func (pA *TK) Contains(strA string, subStrA string) bool {
	return strings.Contains(strA, subStrA)
}

var Contains = TKX.Contains

func (pA *TK) ContainsIgnoreCase(strA string, subStrA string) bool {
	return strings.Contains(strings.ToLower(strA), strings.ToLower(subStrA))
}

var ContainsIgnoreCase = TKX.ContainsIgnoreCase

func (pA *TK) ToLower(strA string) string {
	return strings.ToLower(strA)
}

var ToLower = TKX.ToLower

func (pA *TK) ToUpper(strA string) string {
	return strings.ToUpper(strA)
}

var ToUpper = TKX.ToUpper

// EndsWith 检查字符串strA结尾是否是subStrA
func (pA *TK) EndsWith(strA string, subStrA ...string) bool {
	for _, v := range subStrA {
		if strings.HasSuffix(strA, v) {
			return true
		}
	}

	return false
}

var EndsWith = TKX.EndsWith

// EndsWithIgnoreCase 检查字符串strA结尾是否是subStrA，不区分大小写
func (pA *TK) EndsWithIgnoreCase(strA string, subStrA string) bool {

	return strings.HasSuffix(strings.ToLower(strA), strings.ToLower(subStrA))
}

var EndsWithIgnoreCase = TKX.EndsWithIgnoreCase

// Trim 仅仅封装了strings.TrimSpace，如果带有参数，则调用strings.Trim
func (pA *TK) Trim(strA string, cutSetA ...string) string {
	// Pl("Trim: %#v %#v", strA, cutSetA)
	if len(cutSetA) < 1 {
		return strings.TrimSpace(strA)
	}

	return strings.Trim(strA, cutSetA[0])
}

var Trim = TKX.Trim

func (pA *TK) TrimCharSet(strA string, charSetA string) string {
	return strings.Trim(strA, charSetA)
}

var TrimCharSet = TKX.TrimCharSet

func (pA *TK) InterfaceToStringArray(aryT []interface{}) []string {
	bufT := make([]string, len(aryT))

	for _, v := range aryT {
		bufT = append(bufT, fmt.Sprintf("%v", v))
	}

	return bufT
}

var InterfaceToStringArray = TKX.InterfaceToStringArray

// InStrings 第一个可变参数如果以“-”开头，将表示参数开关，-it表示忽略大小写，并且trim再比较（strA并不trim）
func (pA *TK) InStrings(strA string, argsA ...string) bool {
	ignoreCaseT := false
	trimT := false

	for i, arg := range argsA {
		if i == 0 {
			if len(arg) > 0 {
				if arg[0] == '-' {
					if strings.Contains(arg, "i") {
						ignoreCaseT = true
					}

					if strings.Contains(arg, "t") {
						trimT = true
					}

					if ignoreCaseT {
						strA = strings.ToLower(strA)
					}

					continue
				}
			}
		}

		if ignoreCaseT {
			arg = strings.ToLower(arg)
		}

		if trimT {
			arg = Trim(arg)
		}

		if strA == arg {
			return true
		}
	}

	return false
}

var InStrings = TKX.InStrings

func (pA *TK) InStringsTrim(strA string, argsA ...string) bool {
	for _, arg := range argsA {
		if strA == Trim(arg) {
			return true
		}
	}

	return false
}

var InStringsTrim = TKX.InStringsTrim

func (pA *TK) IsValidEmail(strA string) bool {
	return RegMatch(strA, `[a-zA-Z0-9]+?[a-zA-Z0-9\.\-_]*?@[a-zA-Z0-9]+?(\.[a-zA-Z0-9\.\-_]*)+`)
}

var IsValidEmail = TKX.IsValidEmail

func (pA *TK) GetSliceMaxLen(strA string, maxBytesA int) string {
	lenT := len(strA)

	if lenT <= maxBytesA {
		return strA
	}

	return strA[:maxBytesA]
}

var GetSliceMaxLen = TKX.GetSliceMaxLen

func (pA *TK) FindFirstDiffIndex(strA string, str2A string) int {
	lent1 := len(strA)
	lent2 := len(str2A)

	lent := lent1

	if lent2 < lent {
		lent = lent2
	}

	i := 0

	for i = 0; i < lent; i++ {
		if strA[i] != str2A[i] {
			return i
		}
	}

	if lent1 == lent2 {
		return -1
	}

	return i

}

var FindFirstDiffIndex = TKX.FindFirstDiffIndex

func (pA *TK) FindSamePrefix(strA, str2A string) string {
	idxT := FindFirstDiffIndex(strA, str2A)

	if idxT < 0 {
		return strA
	}

	if idxT > len(strA) {
		return str2A[:idxT]
	}

	return strA[:idxT]
}

var FindSamePrefix = TKX.FindSamePrefix

// IsErrorString 判断是否表示出错的字符串
func (pA *TK) IsErrorString(errStrA string) bool {
	return StartsWith(errStrA, "TXERROR:")
}

var IsErrorString = TKX.IsErrorString

func (pA *TK) IsErrStr(errStrA string) bool {
	return StartsWith(errStrA, "TXERROR:")
}

var IsErrStr = TKX.IsErrStr

func (pA *TK) IsErrStrX(errStrA interface{}) bool {
	nv, ok := errStrA.(string)

	if !ok {
		return false
	}

	return StartsWith(nv, "TXERROR:")
}

var IsErrStrX = TKX.IsErrStrX

// GetErrorString 获取出错字符串中的出错原因部分
func (pA *TK) GetErrorString(errStrA string) string {
	if StartsWith(errStrA, "TXERROR:") {
		return errStrA[8:]
	} else {
		return errStrA
	}
}

var GetErrorString = TKX.GetErrorString

// GetErrorStringSafely 获取出错字符串中的出错原因部分，如果不是出错字符串则返回原串
func (pA *TK) GetErrorStringSafely(errStrA string) string {
	if StartsWith(errStrA, "TXERROR:") {
		return errStrA[8:]
	} else {
		return errStrA
	}
}

var GetErrorStringSafely = TKX.GetErrorStringSafely

func (pA *TK) GetErrStr(errStrA string) string {
	if StartsWith(errStrA, "TXERROR:") {
		return errStrA[8:]
	} else {
		return errStrA
	}
}

var GetErrStr = TKX.GetErrStr

// GenerateErrorString 生成一个出错字符串
func (pA *TK) GenerateErrorString(errStrA string) string {
	return "TXERROR:" + errStrA
}

var GenerateErrorString = TKX.GenerateErrorString

func (pA *TK) ErrStr(errStrA string) string {
	return "TXERROR:" + errStrA
}

var ErrStr = TKX.ErrStr

// GenerateErrorStringF 生成一个出错字符串，但可以加上格式，类似Printf
func (pA *TK) GenerateErrorStringF(formatA string, argsA ...interface{}) string {
	return fmt.Sprintf("TXERROR:"+formatA, argsA...)
}

var GenerateErrorStringF = TKX.GenerateErrorStringF

func (pA *TK) ErrStrF(formatA string, argsA ...interface{}) string {
	return fmt.Sprintf("TXERROR:"+formatA, argsA...)
}

var ErrStrF = TKX.ErrStrF

func (pA *TK) ErrStrf(formatA string, argsA ...interface{}) string {
	return fmt.Sprintf("TXERROR:"+formatA, argsA...)
}

var ErrStrf = TKX.ErrStrf

// ErrorStringToError convert errorstring to error, if not, return nil
func (pA *TK) ErrorStringToError(strA string) error {
	if IsErrorString(strA) {
		return fmt.Errorf("%v", GetErrorString(strA))
	}

	return nil
}

var ErrorStringToError = TKX.ErrorStringToError

func (pA *TK) ErrStrToErr(strA string) error {
	if IsErrorString(strA) {
		return fmt.Errorf("%v", GetErrorString(strA))
	}

	return nil
}

var ErrStrToErr = TKX.ErrStrToErr

func (pA *TK) ErrToStr(errA error) string {
	return fmt.Sprintf("TXERROR:%v", errA.Error())
}

var ErrToStr = TKX.ErrToStr

func (pA *TK) ErrToStrF(formatA string, errA error) string {
	return fmt.Sprintf("TXERROR:"+formatA, errA.Error())
}

var ErrToStrF = TKX.ErrToStrF

func (pA *TK) Replace(strA, findA, replaceA string) string {
	return strings.Replace(strA, findA, replaceA, -1)
}

var Replace = TKX.Replace

func (pA *TK) ReplaceHtmlByMap(strA string, mapA map[string]string) string {
	if mapA == nil {
		return strA
	}

	for k, v := range mapA {
		strA = Replace(strA, "TX_"+k+"_XT", v)
	}

	return strA
}

var ReplaceHtmlByMap = TKX.ReplaceHtmlByMap

func (pA *TK) CleanHtmlPlaceholders(strA string) string {
	strA = RegReplace(strA, `TX_.*?_XT`, "")

	return strA
}

var CleanHtmlPlaceholders = TKX.CleanHtmlPlaceholders

func (pA *TK) StringReplace(strA string, argsA ...string) string {
	if len(argsA) < 2 {
		return strA
	}

	lenT := len(argsA) / 2

	strT := strA
	for i := 0; i < lenT; i++ {
		strT = strings.Replace(strT, argsA[i*2], argsA[i*2+1], -1)
	}

	return strT
}

var StringReplace = TKX.StringReplace

func (pA *TK) SplitLines(strA string) []string {
	if !strings.Contains(strA, "\n") {
		if strings.Contains(strA, "\r") {
			return strings.Split(strA, "\r")
		}
	}
	strT := Replace(strA, "\r", "")
	return strings.Split(strT, "\n")
}

var SplitLines = TKX.SplitLines

func (pA *TK) SplitLinesRemoveEmpty(strA string) []string {
	if !strings.Contains(strA, "\n") {
		if strings.Contains(strA, "\r") {
			strT := RegReplace(strA, "\\r\\s*\\r", "\r")
			return strings.Split(strT, "\r")
		}
	}

	strT := Replace(strA, "\r", "")
	strT = RegReplace(strT, "\\n\\s*\\n", "\n")
	strT = RegReplace(strT, `^\s*\n`, "")
	strT = RegReplace(strT, `\n\s*$`, "")
	return strings.Split(strT, "\n")
}

var SplitLinesRemoveEmpty = TKX.SplitLinesRemoveEmpty

func (pA *TK) Split(strA string, sepA string) []string {
	return strings.Split(strA, sepA)
}

var Split = TKX.Split

func (pA *TK) SplitN(strA string, sepA string, countA int) []string {
	return strings.SplitN(strA, sepA, countA)
}

var SplitN = TKX.SplitN

// SplitByLen 按长度拆分一个字符串为数组，注意由于是rune，可能不是按字节长度，例： listT = strSplitByLen(strT, 10)，可以加第三个参数表示字节数不能超过多少，加第四个参数表示分隔符（遇上分隔符从分隔符后重新计算长度，也就是说分割长度可以超过指定的个数，一般用于有回车的情况）
func (pA *TK) SplitByLen(strA string, lenA int, byteLimitA ...interface{}) []string {
	if len(strA) == 0 {
		return nil
	}
	if lenA >= len(strA) {
		return []string{strA}
	}

	ifByteLimitT := len(byteLimitA) > 0
	ifSepT := len(byteLimitA) > 1

	var byteLimitT int
	var sepT string = ""

	if ifByteLimitT {
		byteLimitT = ToInt(byteLimitA[0])

		if byteLimitT <= 0 {
			byteLimitT = lenA
		}
	}

	if ifSepT {
		sepT = ToStr(byteLimitA[1])
	}

	var sepStartT int

	var byteLenT int
	var curStrT string

	var chunks []string = make([]string, 0, (len(strA)-1)/lenA+1)
	currentLen := 0
	sepStartT = 0
	currentStart := 0
	for i := range strA {
		if sepT != "" && i > 0 {
			if strings.ContainsAny(strA[i-1:i], sepT) {
				// chunks = append(chunks, curStrT)
				currentLen = 0
				sepStartT = i
				// currentStart = i

				// currentLen++
				// continue
			}
		}

		curStrT = strA[currentStart:i]
		if ifByteLimitT {
			byteLenT = len([]byte(strA[sepStartT:i]))
			if byteLenT >= byteLimitT {
				chunks = append(chunks, curStrT)
				currentLen = 0
				currentStart = i
				sepStartT = i
				currentLen++
				continue
			}
		}
		if currentLen == lenA {
			chunks = append(chunks, curStrT)
			currentLen = 0
			currentStart = i
			sepStartT = i
		}
		currentLen++
	}
	chunks = append(chunks, strA[currentStart:])
	return chunks
}

var SplitByLen = TKX.SplitByLen

func (pA *TK) SplitByLenRune(s string, chunkSize int) []string {
	if chunkSize >= len(s) {
		return []string{s}
	}
	var chunks []string
	chunk := make([]rune, chunkSize)
	len := 0
	for _, r := range s {
		chunk[len] = r
		len++
		if len == chunkSize {
			chunks = append(chunks, string(chunk))
			len = 0
		}
	}
	if len > 0 {
		chunks = append(chunks, string(chunk[:len]))
	}
	return chunks
}

var SplitByLenRune = TKX.SplitByLenRune

func (pA *TK) JoinLines(strListA []string, sepA ...string) string {
	if strListA == nil {
		return GenerateErrorString("nil list")
	}

	var sepT string

	if len(sepA) > 0 {
		sepT = sepA[0]
	} else {
		sepT = "\n"
	}

	return strings.Join(strListA, sepT)
}

var JoinLines = TKX.JoinLines

func (pA *TK) JoinLinesBySeparator(strListA []string, sepA ...string) string {
	if strListA == nil {
		return GenerateErrorString("nil list")
	}

	var sepT string

	if len(sepA) > 0 {
		sepT = sepA[0]
	} else {
		sepT = "\n"
	}

	return strings.Join(strListA, sepT)
}

var JoinLinesBySeparator = TKX.JoinLinesBySeparator

// StartsWithBOM if a UTF-8 string starts with BOM
func (pA *TK) StartsWithBOM(strA string) bool {
	bom := []byte{0xEF, 0xBB, 0xBF}

	if StartsWith(strA, string(bom)) {
		return true
	}

	return false

}

var StartsWithBOM = TKX.StartsWithBOM

// RemoveBOM if a UTF-8 string starts with BOM, remove it
func (pA *TK) RemoveBOM(strA string) string {
	bufT := []byte(strA)

	if len(bufT) < 3 {
		return strA
	}

	if bufT[0] == 0xEF && bufT[1] == 0xBB && bufT[2] == 0xBF {
		bufT = bufT[3:]
	}

	return string(bufT)

}

var RemoveBOM = TKX.RemoveBOM

// EnsureValidFileName 确保文件名合理，不保证长度
func (pA *TK) EnsureValidFileName(fileNameA string) string {
	s1 := []rune(fileNameA)

	o1 := []rune{}

	for _, v := range s1 {
		switch v {
		case '%':
			o1 = append(o1, '%')
			o1 = append(o1, '%')
		case ' ', '/', '?', '\\', '*', '"', '<', '>', '\'', '|', ':':
			o1 = append(o1, []rune("%"+hex.EncodeToString([]byte(string(v))))...)
		default:
			o1 = append(o1, v)
		}

	}

	oo := string(o1)

	if strings.HasPrefix(oo, ".") {
		oo = "%2E" + oo[1:]
	}

	return oo
}

var EnsureValidFileName = TKX.EnsureValidFileName

// EnsureValidFileNameX 确保文件名合理并且长度合适
func (pA *TK) EnsureValidFileNameX(fileNameA string) string {
	rs := EncodeStringSimple(fileNameA)

	var extT string
	if len(rs) > 180 {
		extT = filepath.Ext(rs)

		var tmpfn string
		if extT == "" {
			tmpfn = rs
		} else {
			tmpfn = rs[:len(rs)-len(extT)]
		}

		if len(extT) > 30 {
			extT = extT[:30]
			extT = RegReplace(extT, `(%[A-F0-9]?)$`, "")
		}

		if len(tmpfn) > 160 {
			tmpfn = tmpfn[:180]
			tmpfn = RegReplace(tmpfn, `(%[A-F0-9]?)$`, "")
		}

		return tmpfn + extT
	}

	tmps := RegReplace(rs, `(%[A-F0-9]?)$`, "")
	if tmps != rs {
		return tmps
	}

	tmps = RegReplace(tmps, `(%[A-F0-9]?)(\.[^\.]*?)$`, "$2")
	return tmps
}

var EnsureValidFileNameX = TKX.EnsureValidFileNameX

// SyncMap
type SyncMap struct {
	ValueM map[string]interface{}
	LockM  *sync.RWMutex
}

func (pA *TK) NewSyncMap() *SyncMap {
	mapT := &SyncMap{}

	mapT.ValueM = map[string]interface{}{}
	mapT.LockM = new(sync.RWMutex)

	return mapT
}

var NewSyncMap = TKX.NewSyncMap

func (p *SyncMap) Lock() {
	p.LockM.Lock()
}

func (p *SyncMap) TryLock() {
	p.LockM.TryLock()
}

func (p *SyncMap) Unlock() {
	p.LockM.Unlock()
}

func (p *SyncMap) RLock() {
	p.LockM.RLock()
}

func (p *SyncMap) TryRLock() {
	p.LockM.TryRLock()
}

func (p *SyncMap) RUnlock() {
	p.LockM.RUnlock()
}

func (p *SyncMap) QuickClear() {
	p.ValueM = map[string]interface{}{}
}

func (p *SyncMap) QuickSet(keyA string, vA interface{}) {
	p.ValueM[keyA] = vA
}

func (p *SyncMap) QuickDelete(keyA string) {
	delete(p.ValueM, keyA)
}

func (p *SyncMap) QuickGet(keyA string, defaultA ...interface{}) interface{} {
	vA, ok := p.ValueM[keyA]

	if !ok {
		if len(defaultA) > 0 {
			vA = defaultA[0]
		} else {
			vA = nil
		}
	}

	return vA
}

func (p *SyncMap) QuickSize() int {
	vA := len(p.ValueM)

	return vA
}

func (p *SyncMap) Reset() {
	p.ValueM = map[string]interface{}{}
	p.LockM = new(sync.RWMutex)
}

func (p *SyncMap) Clear() {
	p.LockM.Lock()
	p.ValueM = map[string]interface{}{}
	p.LockM.Unlock()
}

func (p *SyncMap) TryClear() bool {
	bT := p.LockM.TryLock()
	if !bT {
		return false
	}

	p.ValueM = map[string]interface{}{}
	p.LockM.Unlock()

	return true
}

func (p *SyncMap) Set(keyA string, vA interface{}) {
	p.LockM.Lock()
	p.ValueM[keyA] = vA
	p.LockM.Unlock()
}

func (p *SyncMap) TrySet(keyA string, vA interface{}) bool {
	bT := p.LockM.TryLock()
	if !bT {
		return false
	}

	p.ValueM[keyA] = vA
	p.LockM.Unlock()

	return true
}

func (p *SyncMap) Delete(keyA string) {
	p.LockM.Lock()
	delete(p.ValueM, keyA)
	p.LockM.Unlock()
}

func (p *SyncMap) TryDelete(keyA string) bool {
	bT := p.LockM.TryLock()
	if !bT {
		return false
	}

	delete(p.ValueM, keyA)
	p.LockM.Unlock()

	return true
}

func (p *SyncMap) Get(keyA string, defaultA ...interface{}) interface{} {
	p.LockM.RLock()

	vA, ok := p.ValueM[keyA]

	if !ok {
		if len(defaultA) > 0 {
			vA = defaultA[0]
		} else {
			vA = nil
		}
	}

	p.LockM.RUnlock()

	return vA
}

func (p *SyncMap) GetList() interface{} {
	rs := make(map[string]interface{})

	p.LockM.RLock()

	for k, v := range p.ValueM {
		rs[k] = v
	}

	p.LockM.RUnlock()

	return rs
}

func (p *SyncMap) TryGet(keyA string, defaultA ...interface{}) interface{} {
	bT := p.LockM.TryRLock()
	if !bT {
		if len(defaultA) > 0 {
			return defaultA[0]
		} else {
			return nil
		}
	}

	vA, ok := p.ValueM[keyA]

	if !ok {
		if len(defaultA) > 0 {
			vA = defaultA[0]
		} else {
			vA = nil
		}
	}

	p.LockM.RUnlock()

	return vA
}

func (p *SyncMap) Size() int {
	p.LockM.RLock()

	vA := len(p.ValueM)

	p.LockM.RUnlock()

	return vA
}

func (p *SyncMap) TrySize() int {
	bT := p.LockM.TryRLock()
	if !bT {
		return -1
	}

	vA := len(p.ValueM)

	p.LockM.RUnlock()

	return vA
}

// Seq
type Seq struct {
	ValueM int
	MutexM sync.Mutex
}

func (pA *TK) NewSeq() *Seq {
	seqT := &Seq{}

	return seqT
}

var NewSeq = TKX.NewSeq

func (p *Seq) Reset(valueA ...int) {
	var valueT int = 0

	if len(valueA) > 0 {
		valueT = valueA[0]
	}

	p.MutexM.Lock()
	p.ValueM = valueT
	p.MutexM.Unlock()
}

func (p *Seq) Get() int {
	var result int

	p.MutexM.Lock()
	p.ValueM++
	result = p.ValueM
	p.MutexM.Unlock()

	return result
}

func (p *Seq) GetCurrent() int {
	var result int

	p.MutexM.Lock()
	result = p.ValueM
	p.MutexM.Unlock()

	return result
}

func (p *Seq) String() string {
	var result int

	p.MutexM.Lock()
	result = p.ValueM
	p.MutexM.Unlock()

	return ToStr(result)
}

var AutoSeq = &Seq{}

func (pA *TK) GetSeq() int {
	return AutoSeq.Get()
}

var GetSeq = TKX.GetSeq

func (pA *TK) ResetSeq(valueA ...int) {
	AutoSeq.Reset(valueA...)
}

var ResetSeq = TKX.ResetSeq

// StringRing
type StringRing struct {
	Buf    []string
	Start  int
	End    int
	MutexM *sync.Mutex
}

func (pA *TK) NewStringRing(sizeA ...int) *StringRing {
	ringT := &StringRing{}

	ringT.MutexM = new(sync.Mutex)

	ringT.Reset(sizeA...)

	return ringT
}

var NewStringRing = TKX.NewStringRing

func (p *StringRing) Reset(sizeA ...int) {
	var sizeT int = 10

	if len(sizeA) > 0 {
		sizeT = sizeA[0]
	}

	p.MutexM.Lock()
	p.Buf = make([]string, sizeT+1)

	p.Start = 0
	p.End = 0
	p.MutexM.Unlock()
}

func (p *StringRing) Push(strA string) {
	p.MutexM.Lock()
	p.Buf[p.Start] = strA
	p.Start++

	if p.Start >= len(p.Buf) {
		p.Start = 0
	}

	if p.Start == p.End {
		p.End++

		if p.End >= len(p.Buf) {
			p.End = 0
		}
	}
	p.MutexM.Unlock()
}

func (p *StringRing) Get(indexA ...int) string {
	var idxT int
	var rs string

	p.MutexM.Lock()
	if len(indexA) > 0 {
		idxT = indexA[0]
	} else {
		idxT = p.Start
	}

	if idxT < 0 || idxT >= len(p.Buf) {
		return ErrStrf("out of range")
	}

	rs = p.Buf[idxT]
	p.MutexM.Unlock()

	return rs
}

func (v StringRing) GetList() []string {
	v.MutexM.Lock()
	bufT := make([]string, 0, len(v.Buf))

	i := v.End
	for i != v.Start {
		bufT = append(bufT, v.Buf[i])

		i++
		if i >= len(v.Buf) {
			i = 0
		}
	}

	v.MutexM.Unlock()
	return bufT
}

func (v StringRing) GetString(sepA ...string) string {
	var sepT string = "\n"

	if len(sepA) > 0 {
		sepT = sepA[0]
	}

	bufT := v.GetList()

	return strings.Join(bufT, sepT)
}

func (v StringRing) String() string {
	v.MutexM.Lock()
	rs := Spr("Start: %v, End: %v, Buf: %v", v.Start, v.End, v.Buf)
	v.MutexM.Unlock()

	return rs
}

// ByteQueue(default size 10, set to -1 for infinite)

type ByteQueueItem struct {
	Prev  *ByteQueueItem
	Next  *ByteQueueItem
	Value byte
}

func (v ByteQueueItem) String() string {
	return fmt.Sprintf("%v", v.Value)
}

type ByteQueue struct {
	Head  *ByteQueueItem
	Tail  *ByteQueueItem
	CapM  int
	SizeM int
}

func (pA *TK) NewByteQueue(capA ...int) *ByteQueue {
	ringT := &ByteQueue{}

	ringT.Reset(capA...)

	return ringT
}

var NewByteQueue = TKX.NewByteQueue

func (p *ByteQueue) Reset(capA ...int) {
	var sizeT int = 10

	if len(capA) > 0 {
		sizeT = capA[0]
	}

	p.CapM = sizeT
	p.SizeM = 0

	p.Head = nil
	p.Tail = nil
}

func (p *ByteQueue) Clear() {
	p.SizeM = 0

	p.Head = nil
	p.Tail = nil
}

func (p *ByteQueue) Size() int {
	return p.SizeM
}

func (p *ByteQueue) Push(byteA byte) {
	if p.CapM >= 0 {
		if p.SizeM >= p.CapM {
			if p.SizeM == 1 {
				p.Head = nil
				p.Tail = nil
			} else {
				p.Head = p.Head.Next
				p.Head.Prev = nil
			}

			p.SizeM--
		}
	}

	itemT := &ByteQueueItem{Value: byteA}

	if p.Tail == nil {
		p.Head = itemT
		p.Tail = itemT
	} else {
		itemT.Prev = p.Tail
		p.Tail.Next = itemT
		p.Tail = itemT
	}

	p.SizeM++
}

func (p *ByteQueue) Insert(idxA int, byteA byte) error {
	if idxA < 0 || idxA >= p.Size() {
		return fmt.Errorf("out of index: %v/%v", idxA, p.Size())
	}

	if p.CapM >= 0 {
		if p.SizeM >= p.CapM {
			if p.SizeM == 1 {
				p.Head = nil
				p.Tail = nil
			} else {
				p.Tail = p.Tail.Prev
				p.Tail.Next = nil
			}

			p.SizeM--
		}
	}

	itemT := &ByteQueueItem{Value: byteA}

	if idxA == 0 {
		if p.Tail == nil {
			p.Head = itemT
			p.Tail = itemT
		} else {
			p.Head.Prev = itemT
			itemT.Next = p.Head

			p.Head = itemT
		}

		p.SizeM++
		return nil
	}

	currentT := p.Head
	cntT := 0

	for cntT < idxA {
		currentT = currentT.Next
		cntT++
	}

	if currentT == nil {
		// return fmt.Errorf("out of index: %v/%v", cntT, p.Size())
		if p.Tail == nil {
			p.Head = itemT
			p.Tail = itemT
		} else {
			itemT.Prev = p.Tail
			p.Tail.Next = itemT
			p.Tail = itemT
			itemT.Next = currentT
		}

		p.SizeM++

		return nil
	}

	// Pl("currentT: %#v", currentT)

	itemT.Prev = currentT.Prev
	currentT.Prev.Next = itemT
	currentT.Prev = itemT
	itemT.Next = currentT

	p.SizeM++

	return nil
}

func (p *ByteQueue) Remove(idxA int) error {
	if idxA < 0 || idxA >= p.Size() {
		return fmt.Errorf("out of index: %v/%v", idxA, p.Size())
	}

	if idxA == 0 {
		if p.Size() == 1 {
			p.Head = nil
			p.Tail = nil
		} else {
			p.Head = p.Head.Next
			p.Head.Prev = nil
		}

		p.SizeM--

		return nil
	}

	if idxA == p.Size()-1 {
		if p.Size() == 1 {
			p.Head = nil
			p.Tail = nil
		} else {
			p.Tail = p.Tail.Prev
			p.Tail.Next = nil
		}

		p.SizeM--

		return nil

	}

	currentT := p.Head
	cntT := 0

	for cntT < idxA {
		currentT = currentT.Next
		cntT++
	}

	currentT.Prev.Next = currentT.Next
	currentT.Next.Prev = currentT.Prev

	p.SizeM--

	return nil
}

// no indexA to get first item, indexA < 0 to get the last item
func (p *ByteQueue) Get(indexA ...int) interface{} {
	var idxT int

	if len(indexA) > 0 {
		idxT = indexA[0]
	} else {
		idxT = 0
	}

	if idxT < 0 {
		idxT = p.Size() - 1
	}

	if idxT < 0 || idxT >= p.Size() {
		return fmt.Errorf("out of index")
	}

	currentT := p.Head
	cntT := 0

	for cntT < idxT {
		currentT = currentT.Next
		cntT++
	}

	if currentT == nil {
		return fmt.Errorf("out of index")
	}

	return currentT.Value
}

// pop the last item
func (p *ByteQueue) Pop() interface{} {
	sizeT := p.Size()
	if sizeT < 1 {
		return fmt.Errorf("no item")
	}

	if sizeT == 1 {
		rs := p.Head.Value

		p.Head = nil
		p.Tail = nil

		p.SizeM = 0

		return rs
	}

	rs := p.Tail.Value

	p.Tail = p.Tail.Prev
	p.Tail.Next = nil

	p.SizeM--

	return rs
}

// pick and pop the last item
func (p *ByteQueue) Pick() interface{} {
	sizeT := p.Size()
	if sizeT < 1 {
		return fmt.Errorf("no item")
	}

	if sizeT == 1 {
		rs := p.Head.Value

		p.Head = nil
		p.Tail = nil

		p.SizeM = 0

		return rs
	}

	rs := p.Head.Value

	p.Head = p.Head.Next
	p.Head.Prev = nil

	p.SizeM--

	return rs
}

func (v ByteQueue) String() string {
	var bufT strings.Builder

	currentT := v.Head

	for currentT != nil {
		bufT.WriteString(fmt.Sprintf(" %v", currentT))

		currentT = currentT.Next
	}

	// rs := Spr("Head: %v, Tail: %v, Buf: %#v", v.Head, v.Tail, v)

	return bufT.String()
}

func (v ByteQueue) GetInfo() string {
	var bufT strings.Builder

	currentT := v.Head

	bufT.WriteString("(ByteQueue)[")
	for currentT != nil {
		bufT.WriteString(fmt.Sprintf(" %v", currentT))

		currentT = currentT.Next
	}
	bufT.WriteString(" ]")

	// rs := Spr("Head: %v, Tail: %v, Buf: %#v", v.Head, v.Tail, v)

	return bufT.String()
}

func (v ByteQueue) GetList() []byte {
	var bufT []byte = make([]byte, v.Size())

	currentT := v.Head

	cntT := 0

	for currentT != nil {
		bufT[cntT] = currentT.Value

		currentT = currentT.Next
		cntT++
	}

	return bufT
}

// AnyQueue

type AnyQueueItem struct {
	Prev  *AnyQueueItem
	Next  *AnyQueueItem
	Value interface{}
}

func (v AnyQueueItem) String() string {
	return fmt.Sprintf("%v", v.Value)
}

// A queue with size limit(default 10, set to -1 for infinite)
type AnyQueue struct {
	Head  *AnyQueueItem
	Tail  *AnyQueueItem
	CapM  int
	SizeM int
}

func (pA *TK) NewAnyQueue(capA ...int) *AnyQueue {
	ringT := &AnyQueue{}

	ringT.Reset(capA...)

	return ringT
}

var NewAnyQueue = TKX.NewAnyQueue

func (p *AnyQueue) Reset(capA ...int) {
	var sizeT int = 10

	if len(capA) > 0 {
		sizeT = capA[0]
	}

	p.CapM = sizeT
	p.SizeM = 0

	p.Head = nil
	p.Tail = nil
}

func (p *AnyQueue) Clear() {
	p.SizeM = 0

	p.Head = nil
	p.Tail = nil
}

func (p *AnyQueue) Size() int {
	return p.SizeM
}

func (p *AnyQueue) Push(byteA interface{}) {
	if p.SizeM >= p.CapM {
		if p.SizeM == 1 {
			p.Head = nil
			p.Tail = nil
		} else {
			p.Head = p.Head.Next
			p.Head.Prev = nil
		}

		p.SizeM--
	}

	itemT := &AnyQueueItem{Value: byteA}

	if p.Tail == nil {
		p.Head = itemT
		p.Tail = itemT
	} else {
		itemT.Prev = p.Tail
		p.Tail.Next = itemT
		p.Tail = itemT
	}

	p.SizeM++
}

func (p *AnyQueue) Insert(idxA int, byteA interface{}) error {
	if idxA < 0 || idxA >= p.Size() {
		return fmt.Errorf("out of index: %v/%v", idxA, p.Size())
	}

	if p.CapM >= 0 {
		if p.SizeM >= p.CapM {
			if p.SizeM == 1 {
				p.Head = nil
				p.Tail = nil
			} else {
				p.Tail = p.Tail.Prev
				p.Tail.Next = nil
			}

			p.SizeM--
		}
	}

	itemT := &AnyQueueItem{Value: byteA}

	if idxA == 0 {
		if p.Tail == nil {
			p.Head = itemT
			p.Tail = itemT
		} else {
			p.Head.Prev = itemT
			itemT.Next = p.Head

			p.Head = itemT
		}

		p.SizeM++
		return nil
	}

	currentT := p.Head
	cntT := 0

	for cntT < idxA {
		currentT = currentT.Next
		cntT++
	}

	if currentT == nil {
		// return fmt.Errorf("out of index: %v/%v", cntT, p.Size())
		if p.Tail == nil {
			p.Head = itemT
			p.Tail = itemT
		} else {
			itemT.Prev = p.Tail
			p.Tail.Next = itemT
			p.Tail = itemT
			itemT.Next = currentT
		}

		p.SizeM++

		return nil
	}

	// Pl("currentT: %#v", currentT)

	itemT.Prev = currentT.Prev
	currentT.Prev.Next = itemT
	currentT.Prev = itemT
	itemT.Next = currentT

	p.SizeM++

	return nil
}

func (p *AnyQueue) Remove(idxA int) error {
	if idxA < 0 || idxA >= p.Size() {
		return fmt.Errorf("out of index: %v/%v", idxA, p.Size())
	}

	if idxA == 0 {
		if p.Size() == 1 {
			p.Head = nil
			p.Tail = nil
		} else {
			p.Head = p.Head.Next
			p.Head.Prev = nil
		}

		p.SizeM--

		return nil
	}

	if idxA == p.Size()-1 {
		if p.Size() == 1 {
			p.Head = nil
			p.Tail = nil
		} else {
			p.Tail = p.Tail.Prev
			p.Tail.Next = nil
		}

		p.SizeM--

		return nil

	}

	currentT := p.Head
	cntT := 0

	for cntT < idxA {
		currentT = currentT.Next
		cntT++
	}

	currentT.Prev.Next = currentT.Next
	currentT.Next.Prev = currentT.Prev

	p.SizeM--

	return nil
}

func (p *AnyQueue) SetByIndex(idxA int, valueA interface{}) error {
	if idxA < 0 || idxA >= p.Size() {
		return fmt.Errorf("out of index: %v/%v", idxA, p.Size())
	}

	currentT := p.Head
	cntT := 0

	for cntT < idxA {
		currentT = currentT.Next
		cntT++
	}

	currentT.Value = valueA

	return nil
}

// no indexA to get first item, indexA < 0 to get the last item
func (p *AnyQueue) Get(indexA ...int) interface{} {
	var idxT int

	if len(indexA) > 0 {
		idxT = indexA[0]
	} else {
		idxT = 0
	}

	if idxT < 0 {
		idxT = p.Size() - 1
	}

	if idxT < 0 || idxT >= p.Size() {
		return fmt.Errorf("out of index")
	}

	currentT := p.Head
	cntT := 0

	for cntT < idxT {
		currentT = currentT.Next
		cntT++
	}

	if currentT == nil {
		return fmt.Errorf("out of index")
	}

	return currentT.Value
}

// pop the last item
func (p *AnyQueue) Pop() interface{} {
	sizeT := p.Size()
	if sizeT < 1 {
		return fmt.Errorf("no item")
	}

	if sizeT == 1 {
		rs := p.Head.Value

		p.Head = nil
		p.Tail = nil

		p.SizeM = 0

		return rs
	}

	rs := p.Tail.Value

	p.Tail = p.Tail.Prev
	p.Tail.Next = nil

	p.SizeM--

	return rs
}

// pick and pop the first item
func (p *AnyQueue) Pick() interface{} {
	sizeT := p.Size()
	if sizeT < 1 {
		return fmt.Errorf("no item")
	}

	if sizeT == 1 {
		rs := p.Head.Value

		p.Head = nil
		p.Tail = nil

		p.SizeM = 0

		return rs
	}

	rs := p.Head.Value

	p.Head = p.Head.Next
	p.Head.Prev = nil

	p.SizeM--

	return rs
}

func (v AnyQueue) String() string {
	var bufT strings.Builder

	bufT.WriteString("[")

	currentT := v.Head

	firstT := true

	for currentT != nil {
		if firstT {
			firstT = false
			bufT.WriteString(fmt.Sprintf("%v", currentT))
		} else {
			bufT.WriteString(fmt.Sprintf(" %v", currentT))
		}

		currentT = currentT.Next
	}

	bufT.WriteString("]")
	// rs := Spr("Head: %v, Tail: %v, Buf: %#v", v.Head, v.Tail, v)

	return bufT.String()
}

func (v AnyQueue) GetInfo() string {
	var bufT strings.Builder

	currentT := v.Head

	for currentT != nil {
		bufT.WriteString(fmt.Sprintf(" %#v", currentT))

		currentT = currentT.Next
	}

	// rs := Spr("Head: %v, Tail: %v, Buf: %#v", v.Head, v.Tail, v)

	return bufT.String()
}

func (v AnyQueue) GetList() []interface{} {
	var bufT []interface{} = make([]interface{}, v.Size())

	currentT := v.Head

	cntT := 0

	for currentT != nil {
		bufT[cntT] = currentT.Value

		currentT = currentT.Next
		cntT++
	}

	return bufT
}

// TXString 相关

type TXString struct {
	string
	Err string
	Obj interface{}
}

func (pA *TK) CreateString(strA string, errA string) *TXString {
	strT := &TXString{}
	strT.string = strA
	strT.Err = errA
	return strT
}

var CreateString = TKX.CreateString

func (pA *TK) CreateStringSimple(strA string) *TXString {
	return &TXString{string: strA, Err: ""}
}

var CreateStringSimple = TKX.CreateStringSimple

func (pA *TK) CreateStringWithObject(strA string, objA interface{}) *TXString {
	return &TXString{string: strA, Err: "", Obj: objA}
}

var CreateStringWithObject = TKX.CreateStringWithObject

func (pA *TK) CreateStringEmpty() *TXString {
	return &TXString{string: "", Err: ""}
}

var CreateStringEmpty = TKX.CreateStringEmpty

func (pA *TK) CreateStringSuccess() *TXString {
	return &TXString{string: "", Err: ""}
}

var CreateStringSuccess = TKX.CreateStringSuccess

func (pA *TK) CreateStringError(errA string) *TXString {
	return &TXString{string: "", Err: errA}
}

var CreateStringError = TKX.CreateStringError

func (pA *TK) CreateStringErrorF(formatA string, argsA ...interface{}) *TXString {
	return &TXString{string: "", Err: fmt.Sprintf(formatA, argsA...)}
}

var CreateStringErrorF = TKX.CreateStringErrorF

func (pA *TK) CreateStringErrorFromTXError(errA string) *TXString {
	if IsErrorString(errA) {
		return &TXString{string: "", Err: GetErrorString(errA)}
	}

	return &TXString{string: errA, Err: ""}
}

var CreateStringErrorFromTXError = TKX.CreateStringErrorFromTXError

func (v TXString) String() string {
	return v.string
}

func (p *TXString) Length() int {
	return len(p.string)
}

func (p *TXString) Len() int {
	return len(p.string)
}

func (p *TXString) CutToLen(lenA int) string {
	lenT := len(p.string)

	if lenT > lenA {
		p.string = p.string[:lenA]
	}

	return p.string
}

func (p *TXString) Error() string {
	return p.Err
}

func (p *TXString) JSONString(optsA ...string) string {
	return ToJSONX(p, optsA...)
}

func (p *TXString) StringEmptyIfError() string {
	if p.IsError() {
		return ""
	}

	return p.string
}

func (p *TXString) InitWithString(strA string) *TXString {
	p.string = strA
	p.Err = ""
	return p
}

func (p *TXString) Set(strA string) *TXString {
	p.string = strA
	p.Err = ""
	return p
}

func (p *TXString) Trim(cutSetA ...string) *TXString {
	p.string = Trim(p.string, cutSetA...)
	return p
}

func (p *TXString) SplitLines() []string {
	return SplitLines(p.string)
}

func (p *TXString) IsError() bool {
	return (p.Err != "")
}

func (p *TXString) IsErrStr() bool {
	return IsErrStr(p.string)
}

func (p *TXString) Print() {
	Pl("string: %v, error: %v", p.string, p.Err)
}

func (p *TXString) PrintWithPrefixTime(prefixA string) {
	Pl("[%v] %vstring: %v, error: %v", GetNowTimeStringFormal(), prefixA, p.string, p.Err)
}

func (p *TXString) PrintWithPrefixTimeLast(prefixA string) {
	Pl("%vstring: %v, error: %v, [%v]", prefixA, p.string, p.Err, GetNowTimeStringFormal())
}

func (p *TXString) PrintWithTimeLast() {
	Pl("string: %v, error: %v, [%v]", p.string, p.Err, GetNowTimeStringFormal())
}

func (p *TXString) PrintResultWithTimeLast() {
	if p.IsError() {
		Pl("Error: %v [%v]", p.Err, GetNowTimeStringFormal())
	} else {
		Pl("Success: %v [%v]", p.string, GetNowTimeStringFormal())
	}
}

func (p *TXString) Replace(patternA string, replacementA string) *TXString {
	p.string = Replace(p.string, patternA, replacementA)

	return p
}

func (p *TXString) Contains(patternA string) bool {
	return strings.Contains(p.string, patternA)
}

func (p *TXString) RegReplace(patternA string, replacementA string) *TXString {
	p.string = RegReplaceX(p.string, patternA, replacementA)

	return p
}

func (p *TXString) RegReplaceX(patternA string, replacementA string) *TXString {
	p.string = RegReplaceX(p.string, patternA, replacementA)

	return p
}

func (p *TXString) RegFindAll(patternA string, groupA int) []string {
	return RegFindAllX(p.string, patternA, groupA)
}

func (p *TXString) RegFindFirst(patternA string, groupA int) string {
	return RegFindFirstX(p.string, patternA, groupA)
}

func (p *TXString) RegFindFirstX(patternA string, groupA int) string {
	return RegFindFirstX(p.string, patternA, groupA)
}

func (p *TXString) List() []string {
	return SplitLines(p.string)
}

func (p *TXString) ToStringList() []string {
	return SplitLines(p.string)
}

func (p *TXString) ToStringListRemoveEmpty() []string {
	return SplitLinesRemoveEmpty(p.string)
}

func (p *TXString) ErrorString() string {
	return p.Err
}

func (p *TXString) ErrorStringF(formatA string) string {
	return fmt.Sprintf(formatA, p.Err)
}

func (p *TXString) EQ(strA string) bool {
	return (p.string == strA)
}

func (p *TXString) Equals(strA string) bool {
	return (p.string == strA)
}

func (p *TXString) EqualsIgnoreCase(strA string) bool {
	return (strings.ToLower(p.string) == strings.ToLower(strA))
}

func (p *TXString) StartsWith(strA string) bool {
	return StartsWith(p.string, strA)
}

func (p *TXString) EndsWith(strA string) bool {
	return EndsWith(p.string, strA)
}

func (p *TXString) IsEmpty() bool {
	return (p.string == "")
}

func (p *TXString) IsEmptyTrim() bool {
	return (Trim(p.string) == "")
}

func (p *TXString) ContainsInHtmlAttributeString(substrA string) bool {
	tmpStr := p.Trim()
	if tmpStr.IsEmpty() {
		return false
	}

	if tmpStr.EQ(substrA) {
		return true
	} else if tmpStr.StartsWith(substrA + " ") {
		return true
	} else if tmpStr.EndsWith(" " + substrA) {
		return true
	}

	return false
}

func (p *TXString) PlErr(prefixA string) *TXString {
	if p.IsError() {
		Pl(prefixA + p.Err)
	}

	return p
}

func (p *TXString) PlSuccessOrErr(workA string) *TXString {
	if p.IsError() {
		Pl(workA + " failed: " + p.Err)
	} else {
		Pl(workA + " successfully")
	}

	return p
}

func (p *TXString) Save(fileA string) *TXString {
	if p.IsError() {
		return p
	}

	p.Err = SaveStringToFile(p.string, fileA)
	return p
}

type TXStringSlice struct {
	body []string
}

type TXStringArray []string

func (aryM TXStringArray) Contains(strA string) bool {
	if len(aryM) < 1 {
		return false
	}

	for _, v := range aryM {
		if v == strA {
			return true
		}
	}

	return false
}

func (aryM TXStringArray) ContainsIgnoreCase(strA string) bool {
	if len(aryM) < 1 {
		return false
	}

	for _, v := range aryM {
		if strings.ToLower(v) == strings.ToLower(strA) {
			return true
		}
	}

	return false
}

func (pA *TK) GenerateErrorStringTX(errStrA string) *TXString {
	return CreateString("", errStrA)
}

var GenerateErrorStringTX = TKX.GenerateErrorStringTX

func (pA *TK) GenerateErrorStringFTX(formatA string, argsA ...interface{}) *TXString {
	return CreateString("", fmt.Sprintf(formatA, argsA...))
}

var GenerateErrorStringFTX = TKX.GenerateErrorStringFTX

func (pA *TK) LoadStringTX(fileNameA string) *TXString {
	if !IfFileExists(fileNameA) {
		return GenerateErrorStringTX("file not exists")
	}

	fileT, err := os.Open(fileNameA)
	if err != nil {
		return GenerateErrorStringTX(err.Error())
	}

	defer fileT.Close()

	fileContentT, err := io.ReadAll(fileT)
	if err != nil {
		return GenerateErrorStringTX(err.Error())
	}

	return CreateStringSimple(string(fileContentT))
}

var LoadStringTX = TKX.LoadStringTX

func (pA *TK) RegContains(strA, patternA string) bool {
	regexpT, errT := regexp.Compile(patternA)

	if errT != nil {
		return false
	}

	return !(regexpT.FindStringIndex(strA) == nil)
}

var RegContains = TKX.RegContains

func (pA *TK) RegSplit(strA, patternA string, nA ...int) []string {
	regexpT, errT := regexp.Compile(patternA)

	if errT != nil {
		return nil
	}

	var nT int = -1

	if (nA != nil) && len(nA) > 0 {
		nT = nA[0]
	}

	return regexpT.Split(strA, nT)
}

var RegSplit = TKX.RegSplit

func (pA *TK) RegSplitX(strA, patternA string, nA ...int) []string {
	regexpT, errT := regexpx.Compile(patternA)

	if errT != nil {
		return nil
	}

	var nT int = -1

	if (nA != nil) && len(nA) > 0 {
		nT = nA[0]
	}

	return regexpT.Split(strA, nT)
}

var RegSplitX = TKX.RegSplitX

func (pA *TK) RegContainsX(strA, patternA string) bool {
	regexpT, errT := regexpx.Compile(patternA)

	if errT != nil {
		return false
	}

	return !(regexpT.FindStringIndex(strA) == nil)
}

var RegContainsX = TKX.RegContainsX

func (pA *TK) RegCount(strA, patternA string) int {
	regexpT, errT := regexpx.Compile(patternA)

	if errT != nil {
		return 0
	}

	rs := regexpT.FindAllStringSubmatchIndex(strA, -1)

	return len(rs)
}

var RegCount = TKX.RegCount
var RegCountX = TKX.RegCount

func (pA *TK) RegFindFirstTX(strA, patternA string, groupA int) *TXString {
	regexpT, errT := regexp.Compile(patternA)

	if errT != nil {
		return CreateStringError("invalid pattern")
	}

	rT := regexpT.FindStringSubmatch(strA)
	if rT == nil {
		return CreateStringError("no match")
	}

	if groupA < len(rT) {
		return CreateStringSimple(rT[groupA])
	}

	return CreateStringError("no group")
}

var RegFindFirstTX = TKX.RegFindFirstTX

func (pA *TK) LoadDualLineListFromString(strA string) [][]string {
	rs := SplitLines(strA)

	lenT := len(rs) / 2

	bufT := make([][]string, lenT)

	var bufP []string

	for i := 0; i < lenT; i++ {
		bufP = make([]string, 2)

		bufP[0] = rs[i*2]
		bufP[1] = rs[i*2+1]

		bufT[i] = bufP
	}

	return bufT
}

var LoadDualLineListFromString = TKX.LoadDualLineListFromString

// 表达式相关 expression related

// FlexEval 计算一个表达式，支持普通语法，第一个参数是表达式字符串，然后是0个或多个参数，在表达式中可以用v1、v2……来指代，表达式采用 github.com/antonmedv/expr 提供的表达式计算引擎，相关进一步文档也可以从这里获取
func (pA *TK) FlexEval(exprA string, varsA ...interface{}) interface{} {
	envT := map[string]interface{}{}

	for i, v := range varsA {
		envT[fmt.Sprintf("v%d", i+1)] = v
	}

	outT, errT := expr.Eval(exprA, envT)

	if errT != nil {
		return errT
	}

	return outT
}

var FlexEval = TKX.FlexEval

// 类似FlexEval，区别是：FlexEval从第二个参数开始可以接受多个参数，并在表达式中以v1、v2这样来指代，而FlexEvalMap则只允许有一个参数，需要是映射类型，这样可以直接用键名在表达式中引用这些变量
func (pA *TK) FlexEvalMap(exprA string, varsA ...interface{}) interface{} {
	var varsT map[string]interface{}
	var ok bool

	if len(varsA) < 1 {
		varsT = map[string]interface{}{}
	} else {
		varsT, ok = varsA[0].(map[string]interface{})

		if !ok {
			varsT = map[string]interface{}{}
		}
	}

	outT, errT := expr.Eval(exprA, varsT)

	if errT != nil {
		return errT
	}

	return outT
}

var FlexEvalMap = TKX.FlexEvalMap

// 正则表达式相关 regex related

func (pA *TK) RegReplace(strA, patternA, replaceA string) string {
	regexpT, errT := regexp.Compile(patternA)

	if errT != nil {
		return strA
	}

	return regexpT.ReplaceAllString(strA, replaceA)
}

var RegReplace = TKX.RegReplace

func (pA *TK) RegReplaceX(strA, patternA, replaceA string) string {
	regexpT, errT := regexpx.Compile(patternA)

	if errT != nil {
		return strA
	}

	return regexpT.ReplaceAllString(strA, replaceA)
}

var RegReplaceX = TKX.RegReplaceX

func (pA *TK) RegFindAll(strA, patternA string, groupA int) []string {
	regexpT, errT := regexp.Compile(patternA)
	if errT != nil {
		return nil
	}

	rT := regexpT.FindAllStringSubmatch(strA, -1)
	if rT == nil {
		return nil
	}

	if groupA < len(rT[0]) {
		rAryT := make([]string, 0, len(rT[0]))
		for i := range rT {
			rAryT = append(rAryT, rT[i][groupA])
		}
		return rAryT
	}

	return nil
}

var RegFindAll = TKX.RegFindAll

func (pA *TK) RegFindAllX(strA, patternA string, groupA int) []string {
	regexpT, errT := regexpx.Compile(patternA)
	if errT != nil {
		return nil
	}

	rT := regexpT.FindAllStringSubmatch(strA, -1)
	if rT == nil {
		return nil
	}

	if groupA < len(rT[0]) {
		rAryT := make([]string, 0, len(rT[0]))
		for i := range rT {
			rAryT = append(rAryT, rT[i][groupA])
		}
		return rAryT
	}

	return nil
}

var RegFindAllX = TKX.RegFindAllX

func (pA *TK) RegFindAllGroups(strA, patternA string) [][]string {
	regexpT, errT := regexp.Compile(patternA)
	if errT != nil {
		return nil
	}

	return regexpT.FindAllStringSubmatch(strA, -1)
}

var RegFindAllGroups = TKX.RegFindAllGroups

func (pA *TK) RegFindAllGroupsX(strA, patternA string) [][]string {
	regexpT, errT := regexpx.Compile(patternA)
	if errT != nil {
		return nil
	}

	return regexpT.FindAllStringSubmatch(strA, -1)
}

var RegFindAllGroupsX = TKX.RegFindAllGroupsX

func (pA *TK) RegFindFirstGroupsX(strA, patternA string) []string {
	regexpT, errT := regexpx.Compile(patternA)
	if errT != nil {
		return nil
	}

	return regexpT.FindStringSubmatch(strA)
}

var RegFindFirstGroupsX = TKX.RegFindFirstGroupsX

// RegFindFirst returns error string if no match or no matching group
func (pA *TK) RegFindFirst(strA, patternA string, groupA int) string {
	regexpT, errT := regexp.Compile(patternA)

	if errT != nil {
		return GenerateErrorStringF("invalid pattern")
	}

	rT := regexpT.FindStringSubmatch(strA)
	if rT == nil {
		return GenerateErrorString("no match")
	}

	if groupA < len(rT) {
		return rT[groupA]
	}

	return GenerateErrorString("no group")
}

var RegFindFirst = TKX.RegFindFirst

func (pA *TK) RegFindFirstX(strA, patternA string, groupA int) string {
	regexpT, errT := regexpx.Compile(patternA)

	if errT != nil {
		return GenerateErrorStringF("invalid pattern")
	}

	rT := regexpT.FindStringSubmatch(strA)
	if rT == nil {
		return GenerateErrorString("no match")
	}

	if groupA < len(rT) {
		return rT[groupA]
	}

	return GenerateErrorString("no group")
}

var RegFindFirstX = TKX.RegFindFirstX

// RegFindFirstIndex the first match location
func (pA *TK) RegFindFirstIndex(strA, patternA string) (int, int) {
	regexpT, errT := regexp.Compile(patternA)

	if errT != nil {
		return -1, -1
	}

	rT := regexpT.FindStringIndex(strA)
	if rT == nil {
		return -1, -1
	}

	return rT[0], rT[1]
}

var RegFindFirstIndex = TKX.RegFindFirstIndex

func (pA *TK) RegFindFirstIndexX(strA, patternA string) (int, int) {
	regexpT, errT := regexpx.Compile(patternA)

	if errT != nil {
		return -1, -1
	}

	rT := regexpT.FindStringIndex(strA)
	if rT == nil {
		return -1, -1
	}

	return rT[0], rT[1]
}

var RegFindFirstIndexX = TKX.RegFindFirstIndexX

func (pA *TK) RegFindAllIndexX(strA, patternA string) [][]int {
	regexpT, errT := regexpx.Compile(patternA)

	if errT != nil {
		return nil
	}

	rT := regexpT.FindAllStringSubmatchIndex(strA, -1)
	if rT == nil {
		return nil
	}

	return rT
}

var RegFindAllIndexX = TKX.RegFindAllIndexX

func (pA *TK) RegFindAllByGroupIndexX(strA, patternA string, groupA int) [][]int {
	regexpT, errT := regexpx.Compile(patternA)

	if errT != nil {
		return nil
	}

	rT := regexpT.FindAllStringSubmatchIndex(strA, -1)
	if rT == nil {
		return nil
	}

	findsT := make([][]int, 0, 10)

	for _, v := range rT {
		if len(v) < (groupA+1)*2 {
			continue
		}

		findsT = append(findsT, []int{v[groupA*2], v[groupA*2+1]})
	}

	if len(findsT) < 0 {
		return nil
	}

	return findsT
}

var RegFindAllByGroupIndexX = TKX.RegFindAllByGroupIndexX

func (pA *TK) RegFindFirstGroupsIndexX(strA, patternA string) []int {
	regexpT, errT := regexpx.Compile(patternA)

	if errT != nil {
		return nil
	}

	rT := regexpT.FindStringSubmatchIndex(strA)

	return rT
}

var RegFindFirstGroupsIndexX = TKX.RegFindFirstGroupsIndexX

func (pA *TK) RegFindFirstGroupIndexX(strA, patternA string, groupA int) []int {
	regexpT, errT := regexpx.Compile(patternA)

	if errT != nil {
		return nil
	}

	rT := regexpT.FindStringSubmatchIndex(strA)

	if len(rT) < (groupA+1)*2 {
		return nil
	}

	return []int{rT[groupA*2], rT[groupA*2+1]}
}

var RegFindFirstGroupIndexX = TKX.RegFindFirstGroupIndexX

func (pA *TK) RegStartsWith(strA, patternA string) bool {
	startT, _ := RegFindFirstIndex(strA, patternA)

	if startT == 0 {
		return true
	}

	return false
}

var RegStartsWith = TKX.RegStartsWith

func (pA *TK) RegStartsWithX(strA, patternA string) bool {
	startT, _ := RegFindFirstIndexX(strA, patternA)

	if startT == 0 {
		return true
	}

	return false
}

var RegStartsWithX = TKX.RegStartsWithX

func (pA *TK) RegMatch(strA, patternA string) bool {
	regexpT, errT := regexp.Compile(patternA)

	if errT != nil {
		return false
	}

	tmps := regexpT.FindString(strA)
	if tmps == strA {
		idxt := regexpT.FindStringIndex(strA)

		if idxt != nil {
			return true
		}

		return false
	}

	return false
}

var RegMatch = TKX.RegMatch

func (pA *TK) RegMatchX(strA, patternA string) bool {
	// Pl("str: %#v, reg: %#v", strA, patternA)
	regexpT, errT := regexpx.Compile(patternA)

	if errT != nil {
		return false
	}

	tmps := regexpT.FindString(strA)
	if tmps == strA {
		idxt := regexpT.FindStringIndex(strA)

		if idxt != nil {
			return true
		}

		return false
	}

	return false
}

var RegMatchX = TKX.RegMatchX

func (pA *TK) RegQuote(textA string) string {
	// textT := regexpx.QuoteMeta(textA)

	textT := strings.ReplaceAll(textA, "\\", "\\\\")

	textT = strings.ReplaceAll(textT, "-", "\\-")
	textT = strings.ReplaceAll(textT, ":", "\\:")
	textT = strings.ReplaceAll(textT, "\"", "\\\"")
	textT = strings.ReplaceAll(textT, "[", "\\[")
	textT = strings.ReplaceAll(textT, "]", "\\]")
	textT = strings.ReplaceAll(textT, "(", "\\(")
	textT = strings.ReplaceAll(textT, ")", "\\)")
	textT = strings.ReplaceAll(textT, "{", "\\{")
	textT = strings.ReplaceAll(textT, "}", "\\}")
	textT = strings.ReplaceAll(textT, "\t", "\\t")
	textT = strings.ReplaceAll(textT, "\r", "\\r")
	textT = strings.ReplaceAll(textT, "\n", "\\n")

	return textT
}

var RegQuote = TKX.RegQuote

// 随机数相关 random related

var ifRandomizedG = false

// Randomize 初始化随机数种子
func (pA *TK) Randomize(seedA ...int) {
	var seedT int64

	if len(seedA) > 0 {
		seedT = int64(seedA[0])
	} else {
		seedT = time.Now().Unix()
	}

	if !ifRandomizedG {
		rand.Seed(seedT)
		ifRandomizedG = true
	}
}

var Randomize = TKX.Randomize

// GetRandomIntLessThan 获取[0-maxA)之间的随机数
func (pA *TK) GetRandomIntLessThan(maxA int) int {
	Randomize()

	randT := rand.Intn(maxA)

	return randT
}

var GetRandomIntLessThan = TKX.GetRandomIntLessThan

func (pA *TK) GetRandomFloat() float64 {
	Randomize()

	randT := rand.Float64()

	return randT
}

var GetRandomFloat = TKX.GetRandomFloat

func (pA *TK) GetRandomInt64LessThan(maxA int64) int64 {
	Randomize()

	randT := rand.Int63n(maxA)

	return randT
}

var GetRandomInt64LessThan = TKX.GetRandomInt64LessThan

// GetRandomIntInRange 获取[minA-maxA]之间的随机数
func (pA *TK) GetRandomIntInRange(minA int, maxA int) int {
	Randomize()

	randT := rand.Intn(maxA+1-minA) + minA

	return randT
}

var GetRandomIntInRange = TKX.GetRandomIntInRange

func (pA *TK) GetRandomInt64InRange(minA int64, maxA int64) int64 {
	Randomize()

	randT := rand.Int63n(maxA+1-minA) + minA

	return randT
}

var GetRandomInt64InRange = TKX.GetRandomInt64InRange

func (pA *TK) GenerateRandomString(minCharA, maxCharA int, hasUpperA, hasLowerA, hasDigitA, hasSpecialCharA, hasSpaceA bool, hasInvalidChars bool) string {
	Randomize()

	if minCharA <= 0 {
		return ""
	}

	if maxCharA <= 0 {
		return ""
	}

	if minCharA > maxCharA {
		return ""
	}

	countT := minCharA + rand.Intn(maxCharA+1-minCharA)

	baseT := ""
	if hasUpperA {
		baseT += "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}

	if hasLowerA {
		baseT += "abcdefghijklmnopqrstuvwxyz"
	}

	if hasDigitA {
		baseT += "0123456789"
	}

	if hasSpecialCharA {
		baseT += "!@#$%^&*-=[]{}."
	}

	if hasSpaceA {
		baseT += " "
	}

	if hasInvalidChars {
		baseT += "/\\:*\"<>|(),+?;"
	}

	rStrT := ""
	var idxT int

	for i := 0; i < countT; i++ {
		idxT = rand.Intn(len(baseT))
		rStrT += baseT[idxT:(idxT + 1)]
	}

	return rStrT
}

var GenerateRandomString = TKX.GenerateRandomString

func (pA *TK) GenerateRandomStringX(argsA ...string) string {
	Randomize()

	minCharA := GetSwitchWithDefaultIntValue(argsA, "-min=", 6)

	if minCharA <= 0 {
		return ""
	}

	maxCharA := GetSwitchWithDefaultIntValue(argsA, "-max=", 6)

	if maxCharA <= 0 {
		return ""
	}

	if minCharA > maxCharA {
		return ""
	}

	countT := minCharA + rand.Intn(maxCharA+1-minCharA)

	baseT := ""

	hasUpperA := !IfSwitchExists(argsA, "-noUpper")

	if hasUpperA {
		baseT += "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}

	hasLowerA := !IfSwitchExists(argsA, "-noLower")
	if hasLowerA {
		baseT += "abcdefghijklmnopqrstuvwxyz"
	}

	hasDigitA := !IfSwitchExists(argsA, "-noDigit")
	if hasDigitA {
		baseT += "0123456789"
	}

	hasSpecialCharA := IfSwitchExists(argsA, "-special")
	if hasSpecialCharA {
		baseT += "!@#$%^&*-=[]{}."
	}

	hasSpaceA := IfSwitchExists(argsA, "-space")
	if hasSpaceA {
		baseT += " "
	}

	hasInvalidChars := IfSwitchExists(argsA, "-invalid")
	if hasInvalidChars {
		baseT += "/\\:*\"<>|(),+?;"
	}

	rStrT := ""
	var idxT int

	for i := 0; i < countT; i++ {
		idxT = rand.Intn(len(baseT))
		rStrT += baseT[idxT:(idxT + 1)]
	}

	return rStrT
}

var GenerateRandomStringX = TKX.GenerateRandomStringX

// RandomX 是一个线程不安全的随机数产生器
type RandomX struct {
	r uint64
}

func (p *RandomX) Randomize() {
	tmpc := time.Now().UnixNano()

	tmpc = (tmpc & 0x0000FFFF) ^ tmpc

	p.r = uint64(tmpc)
}

func (pA *TK) NewRandomGenerator() *RandomX {
	p := &RandomX{r: 0}

	p.Randomize()

	return p
}

var NewRandomGenerator = TKX.NewRandomGenerator

func (p *RandomX) Int64() int64 {
	tmpc := p.r

	tmpc ^= tmpc << 13
	tmpc ^= tmpc >> 17
	tmpc ^= tmpc << 5

	if tmpc < 0 {
		tmpc = -tmpc
	}

	p.r = tmpc
	return int64(tmpc)
}

func (p *RandomX) Float64() float64 {
	tmpc := p.Int64N(1000000000)

	tmpf := float64(tmpc) / 1000000000.0

	return tmpf
}

func (p *RandomX) Int64N(maxA int64) int64 {
	tmpc := p.Int64() % maxA
	if tmpc < 0 {
		tmpc = -tmpc
	}

	return tmpc
}

func (p *RandomX) Int() int {
	tmpc := p.Int64()

	return int(tmpc)
}

// ShuffleStringArray 把字符串数组随机化打乱timesA次
func (pA *TK) ShuffleStringArray(aryA []string, timesA int) {
	Randomize()

	var x, y int
	lent := len(aryA)
	for i := 0; i < timesA; i++ {
		x = GetRandomIntLessThan(lent)
		y = GetRandomIntLessThan(lent)
		if x == y {
			i--
			continue
		}

		aryA[x], aryA[y] = aryA[y], aryA[x]
	}
}

var ShuffleStringArray = TKX.ShuffleStringArray

// GetRandomizeStringArrayCopy 获得一个随机化后的字符串数组
func (pA *TK) GetRandomizeStringArrayCopy(aryA []string) []string {
	Randomize()

	lenT := len(aryA)

	aryT := aryA[0:lenT]

	rs := make([]string, lenT)

	var idxT int

	for i := 0; i < lenT; i++ {
		idxT = GetRandomIntLessThan(len(aryT))

		rs[i] = aryT[idxT]

		aryT = DeleteItemInStringArray(aryT, idxT)
	}

	return rs
}

var GetRandomizeStringArrayCopy = TKX.GetRandomizeStringArrayCopy

func (pA *TK) GetRandomizeSubStringArrayCopy(aryA []string, subCountA int) []string {
	Randomize()

	lenT := len(aryA)

	if subCountA > lenT {
		return nil
	}

	aryT := aryA[0:lenT]

	rs := make([]string, subCountA)

	for i := 0; i < subCountA; i++ {
		idxT := GetRandomIntLessThan(len(aryT))

		rs[i] = aryT[idxT]

		aryT = DeleteItemInStringArray(aryT, idxT)
	}

	return rs
}

var GetRandomizeSubStringArrayCopy = TKX.GetRandomizeSubStringArrayCopy

// GetRandomizeIntArrayCopy 获得一个随机化顺序后的int数组
func (pA *TK) GetRandomizeIntArrayCopy(aryA []int) []int {
	Randomize()

	lenT := len(aryA)

	aryT := aryA[0:lenT]

	rs := make([]int, lenT)

	var idxT int

	for i := 0; i < lenT; i++ {
		idxT = GetRandomIntLessThan((len(aryT)))

		rs[i] = aryT[idxT]

		aryT = DeleteItemInIntArray(aryT, idxT)
	}

	return rs
}

var GetRandomizeIntArrayCopy = TKX.GetRandomizeIntArrayCopy

func (pA *TK) GetRandomizeInt64ArrayCopy(aryA []int64) []int64 {
	Randomize()

	lenT := len(aryA)

	aryT := aryA[0:lenT]

	rs := make([]int64, lenT)

	var idxT int64

	for i := 0; i < lenT; i++ {
		idxT = GetRandomInt64LessThan((int64)(len(aryT)))

		rs[i] = aryT[idxT]

		aryT = DeleteItemInInt64Array(aryT, idxT)
	}

	return rs
}

var GetRandomizeInt64ArrayCopy = TKX.GetRandomizeInt64ArrayCopy

func (pA *TK) GetRandomSubDualList(listA [][]string, countA int) [][]string {
	if countA > len(listA) {
		countA = len(listA)
	}

	l := make([][]string, countA)

	if countA < 1 {
		return l
	}

	mapt := make(map[int]bool)

	lent := len(listA)

	for i := 0; i < countA; i++ {
		for true {
			idx := GetRandomIntLessThan(lent)
			if !mapt[idx] {
				mapt[idx] = true
				l[i] = listA[idx]
				break
			}
		}
	}

	return l
}

var GetRandomSubDualList = TKX.GetRandomSubDualList

func (pA *TK) JoinDualList(listA [][]string, sepItemA, sepInItemA string, withLineNumberA bool) string {
	if listA == nil {
		return ""
	}

	var bufT bytes.Buffer

	for i, v := range listA {
		if i != 0 {
			bufT.WriteString(sepItemA)

		}

		if withLineNumberA {
			bufT.WriteString(Spr("%v. ", i+1))
		}

		bufT.WriteString(v[0])
		bufT.WriteString(sepInItemA)
		bufT.WriteString(v[1])
	}

	return bufT.String()
}

var JoinDualList = TKX.JoinDualList

// 时间相关 time related

// GetNowDateString output likes 20150409
func (pA *TK) GetNowDateString() string {
	t := time.Now()
	return fmt.Sprintf("%04d%02d%02d", t.Year(), t.Month(), t.Day())
}

var GetNowDateString = TKX.GetNowDateString

func (pA *TK) GetNowDateStringFormal() string {
	t := time.Now()
	return fmt.Sprintf("%04d-%02d-%02d", t.Year(), t.Month(), t.Day())
}

var GetNowDateStringFormal = TKX.GetNowDateStringFormal

// GetNowTimeString GetNowTimeString
// "20060102150405"
func (pA *TK) GetNowTimeString() string {
	t := time.Now()
	return fmt.Sprintf("%04d%02d%02d%02d%02d%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}

var GetNowTimeString = TKX.GetNowTimeString

// GetNowTimeStringFormat GetNowTimeStringFormat
// "2006-01-02 15:04:05.000"
func (pA *TK) GetNowTimeStringFormat(formatA ...string) string {
	var formatT string

	if len(formatA) < 1 {
		formatT = TimeFormat
	} else if formatA[0] == "" || formatA[0] == "f" || formatA[0] == "formal" || formatA[0] == "-formal" {
		formatT = TimeFormat
	} else if formatA[0] == "c" || formatA[0] == "compact" || formatA[0] == "-compact" {
		formatT = TimeFormatCompact
	} else {
		formatT = formatA[0]
	}

	t := time.Now()
	return t.Format(formatT)
}

var GetNowTimeStringFormat = TKX.GetNowTimeStringFormat

// GetNowTimeStringFormal get the time string for now as "2020-02-02 08:09:15"
func (pA *TK) GetNowTimeStringFormal() string {
	t := time.Now()
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}

var GetNowTimeStringFormal = TKX.GetNowTimeStringFormal

func (pA *TK) GetNowTimeOnlyStringBeijing() string {
	t := time.Now()
	loc := time.FixedZone("GMT+08:00", 8*3600)
	t = t.In(loc)
	return fmt.Sprintf("%02d%02d%02d", t.Hour(), t.Minute(), t.Second())
}

var GetNowTimeOnlyStringBeijing = TKX.GetNowTimeOnlyStringBeijing

func (pA *TK) GetTimeFromUnixTimeStamp(timeStampA int64) time.Time {
	return time.Unix(timeStampA, 0)
}

var GetTimeFromUnixTimeStamp = TKX.GetTimeFromUnixTimeStamp

func (pA *TK) TimeStampNanoToTime(timeStampA int64) time.Time {
	return time.Unix(timeStampA/1000000000, timeStampA%1000000000)
}

var TimeStampNanoToTime = TKX.TimeStampNanoToTime

func (pA *TK) TimeStampMilliToTime(timeStampA int64) time.Time {
	return time.UnixMilli(timeStampA)
}

var TimeStampMilliToTime = TKX.TimeStampMilliToTime

// in millisecs
func (pA *TK) GetNowTick() int {
	return int(time.Now().UnixNano() / 1000000)
}

var GetNowTick = TKX.GetNowTick
var GetNowTimeStamp = TKX.GetNowTick

func (pA *TK) GetNowTickNano() int {
	return int(time.Now().UnixNano())
}

var GetNowTickNano = TKX.GetNowTickNano

func (pA *TK) GetTimeFromUnixTimeStampMid(timeStampStrA string) time.Time {
	if len(timeStampStrA) < 13 {
		if len(timeStampStrA) == 10 {
			timeStampStrA = timeStampStrA + "000"
		} else {
			return time.Time{}
		}
	}

	return time.Unix(StrToInt64WithDefaultValue(timeStampStrA[:10], 0), StrToInt64WithDefaultValue(timeStampStrA[10:], 0))
}

var GetTimeFromUnixTimeStampMid = TKX.GetTimeFromUnixTimeStampMid

// TimeStampToTime 时间戳转换为时间，如果参数是nil则返回当前时间，如果参数是整数，则按纳秒转换，如果是字符串，则可转换13位（毫秒）或10位（秒）的时间戳，此时如果转换失败则返回时间的零值（1970年...）
func (pA *TK) TimeStampToTime(timeStampA interface{}) time.Time {
	if timeStampA == nil {
		return time.Now()
	}

	nv2, ok := timeStampA.(int64)

	if ok {
		return time.Unix(nv2/1000000000, nv2%1000000000)
	}

	nv1, ok := timeStampA.(int)

	if ok {
		return time.Unix(int64(nv1)/1000000000, int64(nv1)%1000000000)
	}

	nv3, ok := timeStampA.(string)

	if !ok {
		return time.Now()
	}

	if len(nv3) < 13 {
		if len(nv3) == 10 {
			nv3 = nv3 + "000"
		} else {
			return time.Time{}
		}
	}

	return time.Unix(StrToInt64WithDefaultValue(nv3[:10], 0), StrToInt64WithDefaultValue(nv3[10:], 0))
}

var TimeStampToTime = TKX.TimeStampToTime
var TickToTime = TKX.TimeStampToTime

func (pA *TK) GetTimeStamp(timeA time.Time) string {
	return PadString(Int64ToStr(timeA.Unix()), 10)
}

var GetTimeStamp = TKX.GetTimeStamp

func (pA *TK) GetTimeStampMid(timeA time.Time) string {
	tmps := PadString(Int64ToStr(timeA.UnixNano()), 13)
	tmps = tmps[:len(tmps)-6]
	tmps = PadString(tmps, 13)
	return tmps[:13]
}

var GetTimeStampMid = TKX.GetTimeStampMid

func (pA *TK) GetTimeStampNano(timeA time.Time) string {
	return Int64ToStr(timeA.UnixNano())
}

var GetTimeStampNano = TKX.GetTimeStampNano

func (pA *TK) GetTimeStampNanoInt(timeA time.Time) int {
	return int(timeA.UnixNano())
}

var GetTimeStampNanoInt = TKX.GetTimeStampNanoInt

func (pA *TK) NowToFileName() string {
	return StringReplace(time.Now().String(), "-", "_", " ", "_", ":", "_", ".", "_", "+", "_", "=", "_")
}

var NowToFileName = TKX.NowToFileName

func (pA *TK) GetNowTimeStringHourMinute() string {
	t := time.Now()
	return fmt.Sprintf("%02d:%02d", t.Hour(), t.Minute())
}

var GetNowTimeStringHourMinute = TKX.GetNowTimeStringHourMinute

func (pA *TK) GetNowMinutesInDay() int {
	t := time.Now()

	rs := int(t.Hour())*60 + int(t.Minute())

	return rs
}

var GetNowMinutesInDay = TKX.GetNowMinutesInDay

func (pA *TK) NowToStrUTC(formatA string) string {
	n := time.Now().UTC()
	if formatA == "" {
		return (n.Format(TimeFormat))
	}

	return n.Format(formatA)
}

var NowToStrUTC = TKX.NowToStrUTC

func (pA *TK) GetTimeStringDiffMS(str1A, str2A, formatA string, defaultA int64) int64 {
	formatT := Trim(formatA)
	if formatT == "" {
		formatT = TimeFormat
	}

	t1, err := time.Parse(formatT, str1A)
	if err != nil {
		return defaultA
	}

	t2, err := time.Parse(formatT, str2A)
	if err != nil {
		return defaultA
	}

	diffT := t2.Sub(t1)

	return int64(diffT / 1000000)
}

var GetTimeStringDiffMS = TKX.GetTimeStringDiffMS

// return: 1 if str1A > str2A, -1 if str1A < str2A, 0: equal, error if invalid format
func (pA *TK) CompareTimeString(str1A, str2A, formatA string) (int, error) {
	formatT := Trim(formatA)
	if formatT == "" {
		formatT = TimeFormat
	}

	t1, err := time.Parse(formatT, str1A)
	if err != nil {
		return 0, Errf("invalid format for time1")
	}

	t2, err := time.Parse(formatT, str2A)
	if err != nil {
		return 0, Errf("invalid format for time2")
	}

	diffT := t2.Sub(t1)

	if diffT > 0 {
		return -1, nil
	} else if diffT < 0 {
		return 1, nil
	}

	return 0, nil
}

var CompareTimeString = TKX.CompareTimeString

// input: 202301, output 20230201000000, return TXERROR string for error
func (pA *TK) GetNextMonthStart(timeStrA string) string {
	t1, err := time.Parse("200601", timeStrA)
	if err != nil {
		return ErrStrf("invalid format for time")
	}

	t2 := t1.AddDate(0, 1, 0)

	return t2.Format("200601") + "01000000"
}

var GetNextMonthStart = TKX.GetNextMonthStart

func (pA *TK) ToTime(timeA interface{}, defaultA ...interface{}) interface{} {
	// Pl("ToTime: %#v, %#v", timeA, defaultA)
	timeT, ok := timeA.(time.Time)
	if ok {
		return timeT
	}

	var defaultT interface{} = Errf("failed to convert to time: (%T)%#v", timeA, timeA)
	ifLocalT := true

	formatT := ""

	for _, v := range defaultA {
		nv, ok := v.(time.Time)
		if ok {
			defaultT = nv
			continue
		}

		argT, ok := v.(string)
		if !ok {
			continue
		}

		if argT == "-global" {
			ifLocalT = false
			continue
		}

		if argT == "-defaultNow" {
			defaultT = time.Now()
			continue
		}

		if argT == "-defaultErr" {
			defaultT = fmt.Errorf("failed to convert to time: %#v", timeA)
			continue
		}

		if argT == "-defaultErrStr" {
			defaultT = ErrStrf("failed to convert to time: %#v", timeA)
			continue
		}

		if StartsWith(argT, "-format=") {
			tmpStrT := argT[len("-format="):]
			if StartsWith(tmpStrT, "\"") && EndsWith(tmpStrT, "\"") {
				formatT = tmpStrT[1 : len(tmpStrT)-1]
			} else {
				formatT = tmpStrT
			}

			continue
		}

		t := TKX.ToTime(argT)

		if !IsErrX(t) {
			defaultT = t
			continue
		}

	}

	strT, ok := timeA.(string)
	if !ok {
		return defaultT
	}

	if strT == "now" {
		if ifLocalT {
			return time.Now().Local()
		} else {
			return time.Now().UTC()
		}
	}

	var t time.Time
	var err error

	if formatT != "" {
		if ifLocalT {
			t, err = time.ParseInLocation(formatT, strT, time.Local)
			if err == nil {
				return t
			}
		} else {
			t, err = time.Parse(formatT, strT)
			if err == nil {
				return t
			}
		}

		return defaultT
	}

	if ifLocalT {
		t, err = time.ParseInLocation(TimeFormat, strT, time.Local)
	} else {
		t, err = time.Parse(TimeFormat, strT)
	}
	if err == nil {
		return t
	}

	if ifLocalT {
		t, err = time.ParseInLocation(TimeFormatCompact, strT, time.Local)
	} else {
		t, err = time.Parse(TimeFormatCompact, strT)
	}
	if err == nil {
		return t
	}

	if ifLocalT {
		t, err = time.ParseInLocation(TimeFormatMS, strT, time.Local)
	} else {
		t, err = time.Parse(TimeFormatMS, strT)
	}
	if err == nil {
		return t
	}

	if ifLocalT {
		t, err = time.ParseInLocation(TimeFormatDateCompact, strT, time.Local)
	} else {
		t, err = time.Parse(TimeFormatDateCompact, strT)
	}
	if err == nil {
		return t
	}

	// lenT := len(strT)

	// if lenT == 10 {
	// 	strT = strT + "000"
	// }

	// if lenT == 17 {
	// 	strT = strT[:13]
	// }

	// if len(strT) != 13 {
	// 	return defaultT
	// }

	listT := strings.SplitN(strT, ",", 2)

	str1T := listT[0]

	str2T := ""

	if len(listT) > 1 {
		str2T = strings.TrimSpace(listT[1])
	}

	maxT := int64(MAX_INT)

	if len(str1T) > 10 {
		str1T = str1T[0:10]
		str2T = strings.TrimSpace(str1T[10:])
	}

	tickT := StrToInt64(str1T, maxT)

	if tickT == maxT {
		return defaultT
	}

	var tick2T int64

	if str2T == "" {
		tick2T = 0
	} else {
		str2T = PadString(str2T, 9, "-right")
		tick2T = StrToInt64(str2T, maxT)

		if tick2T == maxT {
			return defaultT
		}
	}

	return time.Unix(tickT, tick2T)

	// return defaultT
}

var ToTime = TKX.ToTime

func (pA *TK) StrToTime(strA string, defaultA time.Time) time.Time {
	t, err := time.Parse(TimeFormat, strA)
	if err != nil {
		return defaultA
	}

	return t
}

var StrToTime = TKX.StrToTime

// StrToTimeByFormat default "2006-01-02 15:04:05"
func (pA *TK) StrToTimeByFormat(strA string, formatA string) (time.Time, error) {

	if formatA == "" {
		formatA = "2006-01-02 15:04:05"
	}

	return time.ParseInLocation(formatA, strA, time.Local)
}

var StrToTimeByFormat = TKX.StrToTimeByFormat

// FormatTime default format "2006-01-02 15:04:05"
func (pA *TK) FormatTime(timeA time.Time, formatA ...string) string {
	formatT := ""

	if (formatA != nil) && (len(formatA) > 0) {
		formatT = formatA[0]
	}

	if formatT == "" {
		formatT = "2006-01-02 15:04:05"
	} else if formatT == "compact" || formatT == "-compact" {
		formatT = TimeFormatCompact
	} else if formatT == "full" || formatT == "-full" {
		formatT = TimeFormatMS
	}

	return timeA.Format(formatT)
}

var FormatTime = TKX.FormatTime

func (pA *TK) TimeAddSecs(timeA time.Time, secsA float64) time.Time {
	return timeA.Add(time.Duration(secsA * float64(time.Second)))
}

var TimeAddSecs = TKX.TimeAddSecs

func (pA *TK) TimeAddDate(timeA time.Time, yearA int, monthA int, dayA int) time.Time {
	return timeA.AddDate(yearA, monthA, dayA)
}

var TimeAddDate = TKX.TimeAddDate

// IsYesterday 判断字符串是否是昨天，formatA默认为"20060102"格式
func (pA *TK) IsYesterday(dateStrA string, formatA string) bool {
	if formatA == "" {
		formatA = "20060102"
	}

	// dateT, errT := time.Parse(formatA, dateStrA)
	// if errT != nil {
	// 	return false
	// }

	if time.Now().Add(time.Hour*24*time.Duration(-1)).Format(formatA) == dateStrA {
		return true
	}

	return false
}

var IsYesterday = TKX.IsYesterday

// 切片、数组相关 slice related and array related

func (pA *TK) ReverseStringSlice(ss []string) {
	last := len(ss) - 1
	for i := 0; i < len(ss)/2; i++ {
		ss[i], ss[last-i] = ss[last-i], ss[i]
	}
}

var ReverseStringSlice = TKX.ReverseStringSlice

func (pA *TK) GetArrayItem(aryA interface{}, idxA int, defaultA ...interface{}) interface{} {
	// Pl("GetArrayItem: %#v %v %#v", aryA, idxA, defaultA)
	var defaultT interface{} = Undefined

	if len(defaultA) > 0 {
		defaultT = defaultA[0]
	}

	if aryA == nil {
		return defaultT
	}

	switch nv := aryA.(type) {
	case []interface{}:
		if (idxA < 0) || (idxA >= len(nv)) {
			return defaultT
		}

		return nv[idxA]
	case []int:
		if (idxA < 0) || (idxA >= len(nv)) {
			return defaultT
		}

		return nv[idxA]
	case []byte:
		if (idxA < 0) || (idxA >= len(nv)) {
			return defaultT
		}

		return nv[idxA]
	case []rune:
		if (idxA < 0) || (idxA >= len(nv)) {
			return defaultT
		}

		return nv[idxA]
	case []string:
		if (idxA < 0) || (idxA >= len(nv)) {
			return defaultT
		}

		return nv[idxA]
	case []float64:
		if (idxA < 0) || (idxA >= len(nv)) {
			return defaultT
		}

		return nv[idxA]
	case []map[string]string:
		if (idxA < 0) || (idxA >= len(nv)) {
			return defaultT
		}

		return nv[idxA]
	case []map[string]interface{}:
		if (idxA < 0) || (idxA >= len(nv)) {
			return defaultT
		}

		return nv[idxA]
	case []bool:
		if (idxA < 0) || (idxA >= len(nv)) {
			return defaultT
		}

		return nv[idxA]
	case []uint:
		if (idxA < 0) || (idxA >= len(nv)) {
			return defaultT
		}

		return nv[idxA]
	case []int64:
		if (idxA < 0) || (idxA >= len(nv)) {
			return defaultT
		}

		return nv[idxA]
	case []uint64:
		if (idxA < 0) || (idxA >= len(nv)) {
			return defaultT
		}

		return nv[idxA]
	case []float32:
		if (idxA < 0) || (idxA >= len(nv)) {
			return defaultT
		}

		return nv[idxA]
	case []map[string]int:
		if (idxA < 0) || (idxA >= len(nv)) {
			return defaultT
		}

		return nv[idxA]
	default:
		valueT := reflect.ValueOf(aryA)

		kindT := valueT.Kind()

		if kindT == reflect.Array || kindT == reflect.Slice || kindT == reflect.String {
			lenT := valueT.Len()

			if (idxA < 0) || (idxA >= lenT) {
				return defaultT
			}

			return valueT.Index(idxA).Interface()
		}

		return defaultT
	}

	// return defaultT
}

var GetArrayItem = TKX.GetArrayItem

func (pA *TK) GetArraySlice(aryA interface{}, startA int, endA int, defaultA ...interface{}) interface{} {
	// Pl("GetArraySlice: %#v %v %#v", aryA, idxA, defaultA)
	var defaultT interface{} = Undefined

	if len(defaultA) > 0 {
		defaultT = defaultA[0]
	}

	if aryA == nil {
		return defaultT
	}

	switch nv := aryA.(type) {
	case []interface{}:
		lenT := len(nv)
		if (startA < 0) || (startA >= lenT) {
			return defaultT
		}

		if endA < 0 {
			endA = lenT
		} else if endA > lenT {
			return defaultT
		}

		return nv[startA:endA]
	case []int:
		lenT := len(nv)
		if (startA < 0) || (startA >= lenT) {
			return defaultT
		}

		if endA < 0 {
			endA = lenT
		} else if endA > lenT {
			return defaultT
		}

		return nv[startA:endA]
	case []byte:
		lenT := len(nv)
		if (startA < 0) || (startA >= lenT) {
			return defaultT
		}

		if endA < 0 {
			endA = lenT
		} else if endA > lenT {
			return defaultT
		}

		return nv[startA:endA]
	case []rune:
		lenT := len(nv)
		if (startA < 0) || (startA >= lenT) {
			return defaultT
		}

		if endA < 0 {
			endA = lenT
		} else if endA > lenT {
			return defaultT
		}

		return nv[startA:endA]
	case []string:
		lenT := len(nv)
		if (startA < 0) || (startA >= lenT) {
			return defaultT
		}

		if endA < 0 {
			endA = lenT
		} else if endA > lenT {
			return defaultT
		}

		return nv[startA:endA]
	case []float64:
		lenT := len(nv)
		if (startA < 0) || (startA >= lenT) {
			return defaultT
		}

		if endA < 0 {
			endA = lenT
		} else if endA > lenT {
			return defaultT
		}

		return nv[startA:endA]
	case []map[string]string:
		lenT := len(nv)
		if (startA < 0) || (startA >= lenT) {
			return defaultT
		}

		if endA < 0 {
			endA = lenT
		} else if endA > lenT {
			return defaultT
		}

		return nv[startA:endA]
	case []map[string]interface{}:
		lenT := len(nv)
		if (startA < 0) || (startA >= lenT) {
			return defaultT
		}

		if endA < 0 {
			endA = lenT
		} else if endA > lenT {
			return defaultT
		}

		return nv[startA:endA]
	case []bool:
		lenT := len(nv)
		if (startA < 0) || (startA >= lenT) {
			return defaultT
		}

		if endA < 0 {
			endA = lenT
		} else if endA > lenT {
			return defaultT
		}

		return nv[startA:endA]
	case []uint:
		lenT := len(nv)
		if (startA < 0) || (startA >= lenT) {
			return defaultT
		}

		if endA < 0 {
			endA = lenT
		} else if endA > lenT {
			return defaultT
		}

		return nv[startA:endA]
	case []int64:
		lenT := len(nv)
		if (startA < 0) || (startA >= lenT) {
			return defaultT
		}

		if endA < 0 {
			endA = lenT
		} else if endA > lenT {
			return defaultT
		}

		return nv[startA:endA]
	case []uint64:
		lenT := len(nv)
		if (startA < 0) || (startA >= lenT) {
			return defaultT
		}

		if endA < 0 {
			endA = lenT
		} else if endA > lenT {
			return defaultT
		}

		return nv[startA:endA]
	case []float32:
		lenT := len(nv)
		if (startA < 0) || (startA >= lenT) {
			return defaultT
		}

		if endA < 0 {
			endA = lenT
		} else if endA > lenT {
			return defaultT
		}

		return nv[startA:endA]
	case []map[string]int:
		lenT := len(nv)
		if (startA < 0) || (startA >= lenT) {
			return defaultT
		}

		if endA < 0 {
			endA = lenT
		} else if endA > lenT {
			return defaultT
		}

		return nv[startA:endA]
	default:
		valueT := reflect.ValueOf(aryA)

		kindT := valueT.Kind()

		if kindT == reflect.Array || kindT == reflect.Slice || kindT == reflect.String {
			lenT := valueT.Len()

			if (startA < 0) || (startA >= lenT) {
				return defaultT
			}

			if endA < 0 {
				endA = lenT
			} else if endA > lenT {
				return defaultT
			}

			return valueT.Slice(startA, endA).Interface()
		}

		return defaultT
	}

	// return defaultT
}

var GetArraySlice = TKX.GetArraySlice

func (pA *TK) ArrayContains(aryA interface{}, vA interface{}) bool {
	if aryA == nil {
		return false
	}

	ary1, ok := aryA.([]string)

	if ok {
		for _, v := range ary1 {
			if v == vA {
				return true
			}
		}

		return false
	}

	ary2, ok := aryA.([]int)

	if ok {
		for _, v := range ary2 {
			if v == vA {
				return true
			}
		}

		return false
	}

	ary3, ok := aryA.([]byte)

	if ok {
		for _, v := range ary3 {
			if v == vA {
				return true
			}
		}

		return false
	}

	ary4, ok := aryA.([]float64)

	if ok {
		for _, v := range ary4 {
			if v == vA {
				return true
			}
		}

		return false
	}

	ary5, ok := aryA.([]rune)

	if ok {
		for _, v := range ary5 {
			if v == vA {
				return true
			}
		}

		return false
	}

	aryz, ok := aryA.([]interface{})

	if ok {
		for _, v := range aryz {
			if v == vA {
				return true
			}
		}

		return false
	}

	valueT := reflect.ValueOf(aryA)

	kindT := valueT.Kind()

	if kindT == reflect.Array || kindT == reflect.Slice {
		lenT := valueT.Len()

		for i := 0; i < lenT; i++ {
			if vA == valueT.Index(i).Interface() {
				return true
			}
		}
	}

	return false
}

var ArrayContains = TKX.ArrayContains

func (pA *TK) SafelyGetStringInArray(aryA interface{}, idxA int, optionsA ...string) string {
	defaultT := GetSwitch(optionsA, "-default=", "")

	if aryA == nil {
		return defaultT
	}

	switch nv := aryA.(type) {
	case []string:
		lenT := len(nv)

		if (idxA >= lenT) || (idxA < 0) {
			return defaultT
		}

		if IfSwitchExists(optionsA, "-nofloat") {
			return NilEmptyStringNoFloat(nv[idxA])
		}

		return NilToEmptyStr(nv[idxA])
	case []interface{}:
		lenT := len(nv)

		if (idxA >= lenT) || (idxA < 0) {
			return defaultT
		}

		if IfSwitchExists(optionsA, "-nofloat") {
			return NilEmptyStringNoFloat(nv[idxA])
		}

		return NilToEmptyStr(nv[idxA])
	}

	return defaultT
}

var SafelyGetStringInArray = TKX.SafelyGetStringInArray

func (pA *TK) GetRandomItem(aryA []interface{}) interface{} {
	if aryA == nil {
		return nil
	}

	lenT := len(aryA)

	return aryA[GetRandomIntLessThan(lenT)]
}

var GetRandomItem = TKX.GetRandomItem

func (pA *TK) PickRandomItem(aryA []interface{}) interface{} {
	if aryA == nil {
		return nil
	}

	lenT := len(aryA)

	idxT := GetRandomIntLessThan(lenT)

	itemT := aryA[idxT]

	DeleteItemInArray(aryA, idxT)

	return itemT
}

var PickRandomItem = TKX.PickRandomItem

func (pA *TK) GetRandomStringItem(aryA []string) string {
	if aryA == nil {
		return ErrStrF("nil input")
	}

	lenT := len(aryA)

	return aryA[GetRandomIntLessThan(lenT)]
}

var GetRandomStringItem = TKX.GetRandomStringItem

// DeleteItemInStringArray 删除字符串切片中的某一项
func (pA *TK) DeleteItemInStringArray(aryA []string, idxA int) []string {
	rs := make([]string, 0, len(aryA)-1)
	rs = append(rs, aryA[:idxA]...)
	rs = append(rs, aryA[idxA+1:]...)
	return rs
}

var DeleteItemInStringArray = TKX.DeleteItemInStringArray

// DeleteItemInArray 删除切片中的某一项
func (pA *TK) DeleteItemInArray(aryA []interface{}, idxA int) []interface{} {
	rs := make([]interface{}, 0, len(aryA)-1)
	rs = append(rs, aryA[:idxA]...)
	rs = append(rs, aryA[idxA+1:]...)
	return rs
}

var DeleteItemInArray = TKX.DeleteItemInArray

// DeleteItemInIntArray 删除字符串切片中的某一项
func (pA *TK) DeleteItemInIntArray(aryA []int, idxA int) []int {
	rs := make([]int, 0, len(aryA)-1)
	rs = append(rs, aryA[:idxA]...)
	rs = append(rs, aryA[idxA+1:]...)
	return rs
}

var DeleteItemInIntArray = TKX.DeleteItemInIntArray

func (pA *TK) DeleteItemInInt64Array(aryA []int64, idxA int64) []int64 {
	rs := make([]int64, 0, len(aryA)-1)
	rs = append(rs, aryA[:idxA]...)
	rs = append(rs, aryA[idxA+1:]...)
	return rs
}

var DeleteItemInInt64Array = TKX.DeleteItemInInt64Array

func (pA *TK) DeleteItemInFloat64Array(aryA []float64, idxA int64) []float64 {
	rs := make([]float64, 0, len(aryA)-1)
	rs = append(rs, aryA[:idxA]...)
	rs = append(rs, aryA[idxA+1:]...)
	return rs
}

var DeleteItemInFloat64Array = TKX.DeleteItemInFloat64Array

func (pA *TK) ContainsIn(strA string, subStrsA ...string) bool {
	for _, v := range subStrsA {
		if strings.Contains(strA, v) {
			return true
		}
	}

	return false
}

var ContainsIn = TKX.ContainsIn

func (pA *TK) RegContainsIn(strA string, subStrsA ...string) bool {
	for _, v := range subStrsA {
		if RegContainsX(strA, v) {
			return true
		}
	}

	return false
}

var RegContainsIn = TKX.RegContainsIn

func (pA *TK) ContainsInStringList(strA string, strListA []string) bool {
	if strListA == nil {
		return false
	}

	for _, v := range strListA {
		if strA == v {
			return true
		}
	}

	return false
}

var ContainsInStringList = TKX.ContainsInStringList

func (pA *TK) IndexInStringList(strA string, strListA []string) int {
	if strListA == nil {
		return -1
	}

	for i, v := range strListA {
		if strA == v {
			return i
		}
	}

	return -1
}

var IndexInStringList = TKX.IndexInStringList

func (pA *TK) IndexInStringListFromEnd(strA string, strListA []string) int {
	if strListA == nil {
		return -1
	}

	lent := len(strListA)

	for i := lent - 1; i >= 0; i-- {
		if strA == strListA[i] {
			return i
		}
	}

	return -1
}

var IndexInStringListFromEnd = TKX.IndexInStringListFromEnd

func (pA *TK) GetStringSliceFilled(filledString string, countA int) []string {
	if countA < 0 {
		countA = 0
	}

	bufT := make([]string, countA)

	for i := 0; i < countA; i++ {
		bufT[i] = filledString
	}

	return bufT
}

var GetStringSliceFilled = TKX.GetStringSliceFilled

// Len64 返回int64结果的len
func (pA *TK) Len64(aryA []string) int64 {
	return (int64)(len(aryA))
}

var Len64 = TKX.Len64

func (pA *TK) Int64ArrayToFloat64Array(aryA []int64) []float64 {
	if aryA == nil {
		return nil
	}

	lent := len(aryA)
	rs := make([]float64, lent)

	for i := 0; i < lent; i++ {
		rs[i] = float64(aryA[i])

	}

	return rs

}

var Int64ArrayToFloat64Array = TKX.Int64ArrayToFloat64Array

func (pA *TK) ByteSliceToStringDec(bufA []byte, sepA string) string {
	if bufA == nil {
		return ""
	}

	var outBufT strings.Builder

	lenT := len(bufA)

	for i := 0; i < lenT; i++ {
		if i != 0 {
			outBufT.WriteString(sepA)
		}

		outBufT.WriteString(fmt.Sprintf("%d", bufA[i]))
	}

	return outBufT.String()
}

var ByteSliceToStringDec = TKX.ByteSliceToStringDec

// 映射相关 map related

func (pA *TK) GetMapItem(mapA interface{}, keyA interface{}, defaultA ...interface{}) interface{} {
	var defaultT interface{} = Undefined

	if len(defaultA) > 0 {
		defaultT = defaultA[0]
	}

	if mapA == nil {
		return defaultT
	}

	var rv interface{}
	var ok bool

	switch nv := mapA.(type) {
	case map[string]interface{}:
		rv, ok = nv[ToStr(keyA)]
	case map[string]int:
		rv, ok = nv[ToStr(keyA)]
	case map[string]byte:
		rv, ok = nv[ToStr(keyA)]
	case map[string]rune:
		rv, ok = nv[ToStr(keyA)]
	case map[string]float64:
		rv, ok = nv[ToStr(keyA)]
	case map[string]string:
		rv, ok = nv[ToStr(keyA)]
	case map[string]map[string]string:
		rv, ok = nv[ToStr(keyA)]
	case map[string]map[string]interface{}:
		rv, ok = nv[ToStr(keyA)]
	case *OrderedMap:
		rv, ok = nv.Get(ToStr(keyA))
	default:
		// tk.Plo("here1", v1, v2o)
		valueT := reflect.ValueOf(mapA)

		kindT := valueT.Kind()

		if kindT == reflect.Map {
			rv := valueT.MapIndex(reflect.ValueOf(keyA))

			if !rv.IsValid() {
				return defaultT
			}

			return rv.Interface()
		}

		rv := ReflectGetMember(mapA, keyA)

		if IsError(rv) {
			return defaultT
		}

		return rv
	}

	if !ok {
		return defaultT
	}

	return rv
}

var GetMapItem = TKX.GetMapItem

func (pA *TK) SetMapItem(vA interface{}, keyA interface{}, valueA interface{}) error {
	switch nv := vA.(type) {
	case map[int]int:
		nv[ToInt(keyA)] = ToInt(valueA)
	case map[int]float64:
		nv[ToInt(keyA)] = ToFloat(valueA)
	case map[string]int:
		nv[ToStr(keyA)] = ToInt(valueA)
	case map[string]bool:
		nv[ToStr(keyA)] = ToBool(valueA)
	case map[string]byte:
		nv[ToStr(keyA)] = ToByte(valueA)
	case map[string]rune:
		nv[ToStr(keyA)] = ToRune(valueA)
	case map[string]float64:
		nv[ToStr(keyA)] = ToFloat(valueA)
	case map[string]string:
		nv[ToStr(keyA)] = ToStr(valueA)
	case map[string]interface{}:
		nv[ToStr(keyA)] = valueA
	case *OrderedMap:
		nv.Set(keyA, valueA)
	case url.Values:
		nv.Set(ToStr(keyA), ToStr(valueA))
	case *url.Values:
		nv.Set(ToStr(keyA), ToStr(valueA))
	default:
		valueT := reflect.ValueOf(vA)

		kindT := valueT.Kind()

		if kindT == reflect.Map {
			valueT.SetMapIndex(reflect.ValueOf(keyA), reflect.ValueOf(valueA))
			return nil
		}

		return fmt.Errorf("unsupported type: %T(%v)", vA, vA)
	}

	return nil
}

var SetMapItem = TKX.SetMapItem

// GetValueOfMSS get the value for key in map[string]string
// returns default value if not ok
func (pA *TK) GetValueOfMSS(mapA map[string]string, keyA string, defaultA string) string {
	v, ok := mapA[keyA]

	if ok {
		return v
	}

	return defaultA
}

var GetValueOfMSS = TKX.GetValueOfMSS

func (pA *TK) GetMapKeys(vA interface{}, argsA ...string) interface{} {
	if vA == nil {
		return fmt.Errorf("nil input")
	}

	var rs []interface{}

	switch nv := vA.(type) {
	case map[string]string:
		rs = make([]interface{}, len(nv))

		countT := 0
		for k, _ := range nv {
			rs[countT] = k
			countT++
		}

		return rs
	case map[string]bool:
		rs = make([]interface{}, len(nv))

		countT := 0
		for k, _ := range nv {
			rs[countT] = k
			countT++
		}

		return rs
	case map[string]byte:
		rs = make([]interface{}, len(nv))

		countT := 0
		for k, _ := range nv {
			rs[countT] = k
			countT++
		}

		return rs
	case map[string]rune:
		rs = make([]interface{}, len(nv))

		countT := 0
		for k, _ := range nv {
			rs[countT] = k
			countT++
		}

		return rs
	case map[string]int:
		rs = make([]interface{}, len(nv))

		countT := 0
		for k, _ := range nv {
			rs[countT] = k
			countT++
		}

		return rs
	case map[string]int64:
		rs = make([]interface{}, len(nv))

		countT := 0
		for k, _ := range nv {
			rs[countT] = k
			countT++
		}

		return rs
	case map[string]float64:
		rs = make([]interface{}, len(nv))

		countT := 0
		for k, _ := range nv {
			rs[countT] = k
			countT++
		}

		return rs
	case map[int]string:
		rs = make([]interface{}, len(nv))

		countT := 0
		for k, _ := range nv {
			rs[countT] = k
			countT++
		}

		return rs
	case map[int]int:
		rs = make([]interface{}, len(nv))

		countT := 0
		for k, _ := range nv {
			rs[countT] = k
			countT++
		}

		return rs
	case map[int]float64:
		rs = make([]interface{}, len(nv))

		countT := 0
		for k, _ := range nv {
			rs[countT] = k
			countT++
		}

		return rs
	case map[string]interface{}:
		rs = make([]interface{}, len(nv))

		countT := 0
		for k, _ := range nv {
			rs[countT] = k
			countT++
		}

		return rs
	case *OrderedMap:
		return nv.GetKeys()
	default:
		valueT := reflect.ValueOf(vA)

		kindT := valueT.Kind()

		if kindT == reflect.Map {
			keysT := valueT.MapKeys()

			rvo := reflect.MakeSlice(valueT.Type().Key(), 0, len(keysT))

			for _, v := range keysT {
				reflect.Append(rvo, v)
			}

			rvi := rvo.Interface()

			return rvi
		}

	}

	return fmt.Errorf("unsupported type: %T(%v)", vA, vA)
}

var GetMapKeys = TKX.GetMapKeys

// 系统相关函数 system related

func (pA *TK) GetSystemInfo(optsA ...string) interface{} {
	diskT := strings.TrimSpace(GetSwitch(optsA, "-disk=", ""))
	
	cpuTimeT := ToFloat(GetSwitch(optsA, "-cpuTime=", "0.1"), 0.1)
	
	rs := map[string]interface{} {
		"Cpu": map[string]interface{} {
			"Total": 0,
			"Free": 0,
			"Used": 0,
			"Percent": 0,
		},
		"Mem": map[string]interface{} {
			"Total": 0,
			"Free": 0,
			"Used": 0,
			"Percent": 0,
		},
		"Disk": map[string]interface{} {
			"Total": 0,
			"Free": 0,
			"Used": 0,
			"Percent": 0,
		},
	}
	
	v0, errT := cpu.Percent(time.Duration(float64(time.Second) * cpuTimeT), false)
	
	if errT != nil {
		return errT
	}
	
//	rs["Cpu"].(map[string]interface{})["Total"] = v0.Total
//	rs["Cpu"]["Free"] = v0.Free
//	rs["Cpu"]["Used"] = v0.Total - v0.Free
	rs["Cpu"].(map[string]interface{})["Percent"] = v0[0]

	v1, errT := mem.VirtualMemory()

	if errT != nil {
		return errT
	}
	
	rs["Mem"].(map[string]interface{})["Total"] = v1.Total
	rs["Mem"].(map[string]interface{})["Free"] = v1.Free
	rs["Mem"].(map[string]interface{})["Used"] = v1.Used
	rs["Mem"].(map[string]interface{})["Percent"] = v1.UsedPercent

	if diskT != "" {
		v3, errT := disk.Usage(diskT)
		
		if errT != nil {
			return errT
		}
		
		rs["Disk"].(map[string]interface{})["Total"] = v3.Total
		rs["Disk"].(map[string]interface{})["Free"] = v3.Free
		rs["Disk"].(map[string]interface{})["Used"] = v3.Used
		rs["Disk"].(map[string]interface{})["Percent"] = v3.UsedPercent

	}
	return rs
}

var GetSystemInfo = TKX.GetSystemInfo

func (pA *TK) GetChar() interface{} {
	char, key, err := keyboard.GetSingleKey()
	if err != nil {
		return err
	}
	// fmt.Printf("You pressed: %q\r\n", char)

	return fmt.Sprintf("%v+%v", key, char)
}

var GetChar = TKX.GetChar

func (pA *TK) GetChar2() interface{} {
	// switch stdin into 'raw' mode
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	b := make([]byte, 1)
	_, err = os.Stdin.Read(b)
	if err != nil {
		return err
	}

	return ToStr(int(b[0]))
}

var GetChar2 = TKX.GetChar2

func (pA *TK) Sscanf(strA string, formatA string, argsA ...interface{}) interface{} {
	rs, errT := fmt.Sscanf(strA, formatA, argsA...)

	if errT != nil {
		return errT
	}

	return rs
}

var Sscanf = TKX.Sscanf

// GetOSArgs return os.Args
func (pA *TK) GetOSArgs() []string {
	return os.Args
}

var GetOSArgs = TKX.GetOSArgs

func (pA *TK) GetOSArgsShort() []string {
	return os.Args[1:]
}

var GetOSArgsShort = TKX.GetOSArgsShort

// EnsureBasePath make sure a base path for application is exists, otherwise create it
// first look for c:\nameA(Windows) or /nameA(Mac&Linux), then the application path
// if nameA contains ".", "/", "\\", will use it as basepath instead
func (pA *TK) EnsureBasePath(nameA string) (string, error) {
	var basePathT string

	if ContainsIn(nameA, ".", "/", "\\") {
		baseT, errT := filepath.Abs(nameA)
		if errT != nil {
			return "", errT
		}

		basePathT = baseT
	} else {
		if strings.HasPrefix(runtime.GOOS, "win") {
			basePathT = "c:\\" + nameA
		} else {
			basePathT = "/" + nameA
		}

	}

	errT := EnsureMakeDirsE(basePathT)

	if (errT != nil) || (!IfFileExists(basePathT)) {
		basePathT = GetApplicationPath()
	}

	// filePathT := filepath.Join(basePathT, "basePathT.txt")

	// errT = SaveStringToFileE("test", filePathT)

	// if (errT != nil) || (!IfFileExists(filePathT)) {
	// 	return "", Errf("init base path failed")
	// }

	// errT = RemoveFile(filePathT)

	// if (errT != nil) || (IfFileExists(filePathT)) {
	// 	return "", Errf("init base path failed: failed to remove tmp file")
	// }

	return basePathT, nil
}

var EnsureBasePath = TKX.EnsureBasePath

// EnsureBasePath make sure a base path in user's home directory for application is exists, otherwise create it
// if failed, will try to create in the application path
func (pA *TK) EnsureBasePathInHome(nameA string) string {
	var basePathT string

	nameT := strings.TrimSpace(nameA)

	if nameT == "" {
		return ErrStrf("empty path name")
	}

	homePathT, errT := os.UserHomeDir()

	if errT != nil {
		return ErrStrf("failed to get home path: %v", errT)
	}

	basePathT = filepath.Join(homePathT, nameT)

	errT = EnsureMakeDirsE(basePathT)

	if errT != nil {
		return ErrStrf("failed to create path: %v", errT)
	}

	if !IfFileExists(basePathT) {
		return ErrStrf("failed to create path: %v", basePathT)
	}

	return basePathT
}

var EnsureBasePathInHome = TKX.EnsureBasePathInHome

// CreateTempFile dirA如果为空，则在系统临时目录下。patternA 可以是example或example*.txt这样，会自动替换*为随机字符串，创建者负责删除
func (pA *TK) CreateTempFile(dirA string, patternA string, optsA ...string) (string, error) {
	content := []byte("")
	tmpfile, err := os.CreateTemp(dirA, patternA)
	if err != nil {
		return "", err
	}

	defer tmpfile.Close()

	rs := tmpfile.Name()

	// defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write(content); err != nil {
		return rs, err
	}

	if err := tmpfile.Close(); err != nil {
		return rs, err
	}

	return rs, nil
}

var CreateTempFile = TKX.CreateTempFile

// CreateTempDir dirA如果为空，则在系统临时目录下。patternA 可以是example或example*这样，会自动替换*为随机字符串，创建者负责删除
func (pA *TK) CreateTempDir(dirA string, patternA string, optsA ...string) (string, error) {
	tmpDir, err := os.MkdirTemp(dirA, patternA)
	if err != nil {
		return "", err
	}

	return tmpDir, nil
}

var CreateTempDir = TKX.CreateTempDir

func (pA *TK) CopyFile(src, dst string, optsA ...string) error {

	srcFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	mode := srcFileStat.Mode()

	if !mode.IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	if mode.IsDir() {
		return fmt.Errorf("%s is a folder", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}

	defer source.Close()

	forceT := IfSwitchExists(optsA, "-force")

	if !forceT {
		if IfFileExists(dst) {
			return fmt.Errorf("file %s already exists", dst)
		}
	}

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}

	defer destination.Close()

	bufferSizeT := GetSwitchWithDefaultIntValue(optsA, "-bufferSize=", 1000000)

	if bufferSizeT <= 0 {
		bufferSizeT = 1000000
	}

	buf := make([]byte, bufferSizeT)
	for {
		n, err := source.Read(buf)

		if err != nil && err != io.EOF {
			return err
		}

		if n == 0 {
			break
		}

		if _, err := destination.Write(buf[:n]); err != nil {
			return err
		}
	}

	return err
}

var CopyFile = TKX.CopyFile

// GetCurrentThreadID get goroutineid
func (pA *TK) GetCurrentThreadID() string {
	var buf [64]byte

	n := runtime.Stack(buf[:], false)

	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]

	id, errT := strconv.Atoi(idField)
	if errT != nil {
		return GenerateErrorStringF("failed to get goroutine id: %v", errT)
	}

	return IntToStr(id)
}

var GetCurrentThreadID = TKX.GetCurrentThreadID

// Exit usage: Exit() or Exit(number)
func (pA *TK) Exit(c ...int) {
	if c == nil || len(c) < 1 {
		os.Exit(0)
	}

	os.Exit(c[0])
}

var Exit = TKX.Exit

// RunWinFileWithSystemDefault run a program or open a file with default program in Windows
func (pA *TK) RunWinFileWithSystemDefault(fileA string) string {
	if StartsWith(runtime.GOOS, "darwin") {
		return SystemCmd("open", fileA)
	} else if EndsWith(runtime.GOOS, "nux") {
		return SystemCmd("xdg-open", fileA)
	}

	cmd := exec.Command("cmd", "/C", "start", "", fileA)
	err := cmd.Start()
	if err != nil {
		return err.Error()
	}

	return ""
}

var RunWinFileWithSystemDefault = TKX.RunWinFileWithSystemDefault

// SystemCmd runs a system command, such as "cmd /c dir", "cmd /k copy a.txt b.txt".
func (pA *TK) SystemCmd(cmdA string, argsA ...string) string {
	var out, errBuf bytes.Buffer

	cmd := exec.Command(cmdA, argsA...)

	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	
	errT := cmd.Run()
	if errT != nil {
		return GenerateErrorStringF("failed: %v (%v/%v)", errT, out.String(), errBuf.String())
	}

	return out.String()
}

var SystemCmd = TKX.SystemCmd

// SystemCmdDetached runs a system command and detaches it, such as "cmd /c dir", "cmd /k copy a.txt b.txt".
func (pA *TK) SystemCmdDetached(cmdA string, argsA ...string) string {
	cmd := exec.Command(cmdA, argsA...)

	errT := cmd.Start()
	if errT != nil {
		return GenerateErrorStringF("failed to start new process and detach: %v", errT)
	}
	
	cmd.Process.Release()

	return ""
}

var SystemCmdDetached = TKX.SystemCmdDetached

// NewSSHClient create SSH client with fewer settings
func (pA *TK) NewSSHClient(hostA string, portA interface{}, userA string, passA string) (*goph.Client, error) {
	if strings.HasPrefix(passA, "//TXDEF#") {
		passA = DecryptStringByTXDEF(passA)
	}

	authT := goph.Password(passA)

	portT := ToStr(portA)

	clientT, errT := goph.NewConn(&goph.Config{
		User:     userA,
		Addr:     hostA,
		Port:     uint(ToInt(portT)),
		Auth:     authT,
		Timeout:  goph.DefaultTimeout,
		Callback: ssh.InsecureIgnoreHostKey(),
	})

	return clientT, errT

	// clientT := &goph.Client{
	// 	Addr: hostA,
	// 	Port: portA,
	// 	User: userA,
	// 	Auth: authT,
	// }

	// errT := goph.Conn(clientT, &ssh.ClientConfig{
	// 	User:    clientT.User,
	// 	Auth:    clientT.Auth,
	// 	Timeout: 20 * time.Second,
	// 	HostKeyCallback: func(host string, remote net.Addr, key ssh.PublicKey) error {
	// 		return nil
	// 		// hostFound, err := goph.CheckKnownHost(host, remote, key, "")

	// 		// if hostFound && err != nil {
	// 		// 	return err
	// 		// }

	// 		// if hostFound && err == nil {
	// 		// 	return nil
	// 		// }

	// 		// return goph.AddKnownHost(host, remote, key, "")
	// 	},
	// })

	// clientT, errT := goph.NewConn(userA, hostA, authT, func(host string, remote net.Addr, key ssh.PublicKey) error {

	// 	hostFound, err := goph.CheckKnownHost(host, remote, key, "")

	// 	if hostFound && err != nil {
	// 		return err
	// 	}

	// 	if hostFound && err == nil {
	// 		return nil
	// 	}

	// 	return goph.AddKnownHost(host, remote, key, "")
	// })

	return clientT, errT
}

var NewSSHClient = TKX.NewSSHClient

// Prf 仅仅是封装了fmt.Printf函数，但会返回format字符串
func (pA *TK) Prf(formatA string, argsA ...interface{}) string {
	fmt.Printf(formatA, argsA...)

	return formatA
}

var Prf = TKX.Prf

// Prl 仅仅封装了fmt.Println函数
func (pA *TK) Prl(a ...interface{}) {
	fmt.Println(a...)
}

var Prl = TKX.Prl

// Pln 仅仅封装了fmt.Println函数
func (pA *TK) Pln(a ...interface{}) {
	fmt.Println(a...)
}

var Pln = TKX.Pln

// Printf 仅仅封装了fmt.Printf函数，与其完全一致
func (pA *TK) Printf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}

var Printf = TKX.Printf

// Printfln 仅仅封装了fmt.Printf函数，但结尾会多输出一个换行符
func (pA *TK) Printfln(format string, a ...interface{}) {
	fmt.Printf(format+"\n", a...)
}

var Printfln = TKX.Printfln

// Spr 仅仅是封装了fmt.Sprintf函数
func (pA *TK) Spr(formatA string, argsA ...interface{}) string {
	return fmt.Sprintf(formatA, argsA...)
}

var Spr = TKX.Spr

// Pr 即fmt.Print
func (pA *TK) Pr(argsA ...interface{}) {
	fmt.Print(argsA...)
}

var Pr = TKX.Pr

// Pl 类似Pr，但结尾会加有一个回车
func (pA *TK) Pl(formatA string, argsA ...interface{}) {
	fmt.Printf(formatA+"\n", argsA...)
}

var Pl = TKX.Pl

func (pA *TK) Plo(vA ...interface{}) {
	lenT := len(vA)

	if lenT < 1 {
		fmt.Println()
		return
	}

	// Pl("Plo (%T)%#v", vA[0], vA[0])
	if lenT == 1 {
		fmt.Printf("(%T)%#v\n", vA[0], vA[0])
		return
	}

	for i, v := range vA {
		fmt.Printf("[%v] (%T)%#v\n", i, v, v)
	}
}

var Plo = TKX.Plo

func (pA *TK) Plt(vA ...interface{}) {
	lenT := len(vA)

	if lenT < 1 {
		fmt.Println()
		return
	}

	if lenT == 1 {
		fmt.Printf("(%T)%v\n", vA[0], vA[0])
		return
	}

	for i, v := range vA {
		fmt.Printf("[%v] (%T)%v\n", i, v, v)
	}
}

var Plt = TKX.Plt

func (pA *TK) Plos(vA ...interface{}) {
	for i, v := range vA {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Printf("(%T)%#v", v, v)
	}

	fmt.Print("\n") // test1
}

var Plos = TKX.Plos

func (pA *TK) Plosr(vA ...interface{}) {
	for i, v := range vA {
		fmt.Printf("[%v] (%T)%v\n", i, v, v)
	}
}

var Plosr = TKX.Plosr

// PlNow 类似Pl，但前面会加有当前时间标记
func (pA *TK) PlNow(formatA string, argsA ...interface{}) {
	fmt.Printf(fmt.Sprintf("[%v] ", time.Now().Format(TimeFormatCompact2))+formatA+"\n", argsA...)
}

var PlNow = TKX.PlNow

// PlVerbose 类似Pl，但仅在verboseA为true时才输出
func (pA *TK) PlVerbose(verboseA bool, formatA string, argsA ...interface{}) {
	if verboseA {
		fmt.Printf(formatA+"\n", argsA...)
	}
}

var PlVerbose = TKX.PlVerbose

// Fpl 类似Pl，但向流中写入(Fprintf)
func (pA *TK) Fpl(wA io.Writer, formatA string, argsA ...interface{}) {
	fmt.Fprintf(wA, formatA+"\n", argsA...)
}

var Fpl = TKX.Fpl

// Fpr 类似Pr，但向流中写入(Fprintf)
func (pA *TK) Fpr(wA io.Writer, formatA string, argsA ...interface{}) {
	fmt.Fprintf(wA, formatA, argsA...)
}

var Fpr = TKX.Fpr

func (pA *TK) PlvWithError(vA interface{}, errStrA string) {
	if errStrA == "" {
		fmt.Printf("%v\n", vA)
	} else {
		fmt.Printf("Error: %v\n", vA)
	}
}

var PlvWithError = TKX.PlvWithError

func (pA *TK) PlAndExit(formatA string, argsA ...interface{}) {
	fmt.Printf(formatA+"\n", argsA...)
	os.Exit(0)
}

var PlAndExit = TKX.PlAndExit

// PrlErrSimple 输出错误信息，结尾加一个回车
func (pA *TK) PlErrSimple(formatA string, argsA ...interface{}) {
	fmt.Printf("Error: "+formatA+"\n", argsA...)
}

var PlErrSimple = TKX.PlErrSimple

func (pA *TK) PlErrSimpleAndExit(formatA string, argsA ...interface{}) {
	fmt.Printf("Error: "+formatA+"\n", argsA...)
	os.Exit(0)
}

var PlErrSimpleAndExit = TKX.PlErrSimpleAndExit

func (pA *TK) PlErrAndExit(errA error) {
	fmt.Printf("Error: " + errA.Error() + "\n")
	os.Exit(0)
}

var PlErrAndExit = TKX.PlErrAndExit

func (pA *TK) PlTXErr(strA string) {
	fmt.Printf("Error: " + GetErrorString(strA) + "\n")
}

var PlTXErr = TKX.PlTXErr

func (pA *TK) PlSimpleErrorString(strA string) {
	fmt.Printf("Error: " + strA + "\n")
}

var PlSimpleErrorString = TKX.PlSimpleErrorString

func (pA *TK) PlErr(errA error) {
	if errA == nil {
		return
	}

	Pl("Error: %v", errA.Error())
}

var PlErr = TKX.PlErr

func (pA *TK) PlErrX(errA interface{}) {
	if !IsErrX(errA) {
		return
	}

	Pl("Error: %v", GetErrStrX(errA))
}

var PlErrX = TKX.PlErrX

func (pA *TK) PlErrString(strA string) {
	if !IsErrorString(strA) {
		return
	}

	Pl("Error: %v", GetErrorString(strA))
}

var PlErrString = TKX.PlErrString

func (pA *TK) PlErrWithPrefix(prefixA string, errA error) {
	if errA == nil {
		return
	}

	Pl("%v%v", prefixA, errA.Error())
}

var PlErrWithPrefix = TKX.PlErrWithPrefix

// Plv output one variable
func (pA *TK) Plv(argsA ...interface{}) {
	for i := 0; i < len(argsA); i++ {
		if i > 0 {
			fmt.Printf(" ")
		}

		fmt.Printf("%#v", argsA[i])
	}

	fmt.Println()
}

var Plv = TKX.Plv

func (pA *TK) Plvx(argsA interface{}) {
	fmt.Printf("[TYPE] %T [VALUE] %v [ITYPE] %#v\n", argsA, argsA, argsA)
}

var Plvx = TKX.Plvx

// Plvs output several variables, seperated by sepA
func (pA *TK) Plvs(sepA string, argsA ...interface{}) {
	lenT := len(argsA)

	strListA := GetStringSliceFilled("%#v", lenT)

	formatT := strings.Join(strListA, sepA)

	fmt.Printf(formatT+"\n", argsA...)
}

var Plvs = TKX.Plvs

// Plvsr output several variables, seperated by \n (new line character)
func (pA *TK) Plvsr(argsA ...interface{}) {
	Plvs("\n", argsA...)
}

var Plvsr = TKX.Plvsr

// Errf wrap fmt.Errorf function
func (pA *TK) Errf(formatA string, argsA ...interface{}) error {
	return fmt.Errorf(formatA, argsA...)
}

var Errf = TKX.Errf

func (pA *TK) FatalErr(prefixA string, errA error) {
	Pl("%v%v", prefixA, errA.Error())

	os.Exit(1)
}

var FatalErr = TKX.FatalErr

func (pA *TK) FatalErrf(formatA string, errA error) {
	Pl(formatA, errA.Error())

	os.Exit(1)
}

var FatalErrf = TKX.FatalErrf

func (pA *TK) Fatalf(formatA string, argsA ...interface{}) {
	Pl(formatA, argsA...)

	os.Exit(1)
}

var Fatalf = TKX.Fatalf

func (pA *TK) FatalfByCondition(ignoreA bool, formatA string, argsA ...interface{}) {
	Pl(formatA, argsA...)

	if !ignoreA {
		os.Exit(1)
	}
}

var FatalfByCondition = TKX.FatalfByCondition

func (pA *TK) CheckErr(prefixA string, errA error) {
	if errA == nil {
		return
	}

	Pl("%v%v", prefixA, errA.Error())

	os.Exit(1)
}

var CheckErr = TKX.CheckErr

func (pA *TK) CheckErrf(formatA string, argsA ...interface{}) {
	var errT error = nil

	if argsA == nil {
		return
	}

	for _, v := range argsA {
		tmpV, ok := v.(error)
		if !ok {
			continue
		}

		errT = tmpV
	}

	if errT == nil {
		return
	}

	Pl(formatA, argsA...)

	os.Exit(1)
}

var CheckErrf = TKX.CheckErrf

func (pA *TK) CheckErrStrf(formatA string, errStrA string, argsA ...interface{}) {
	if !IsErrStr(errStrA) {
		return
	}

	Pl(formatA, append([]interface{}{GetErrStr(errStrA)}, argsA...)...)

	os.Exit(1)
}

var CheckErrStrf = TKX.CheckErrStrf

func (pA *TK) CheckErrStr(errStrA string, argsA ...interface{}) {
	if !IsErrStr(errStrA) {
		return
	}

	Pl("error: %v", append([]interface{}{GetErrStr(errStrA)}, argsA...)...)

	os.Exit(1)
}

var CheckErrStr = TKX.CheckErrStr

func (pA *TK) CheckErrCompact(errA error) {
	if errA == nil {
		return
	}

	Prl(errA.Error())

	os.Exit(1)
}

var CheckErrCompact = TKX.CheckErrCompact

// GetEnv same as os.Getenv
func (pA *TK) GetEnv(keyA string) string {
	return os.Getenv(keyA)
}

var GetEnv = TKX.GetEnv

func (pA *TK) SetEnv(keyA string, valueA string) string {
	return ErrToErrStr(os.Setenv(keyA, valueA))
}

var SetEnv = TKX.SetEnv

// JoinPath same as filepath.Join
func (pA *TK) JoinPath(elemA ...string) string {
	return filepath.Join(elemA...)
}

var JoinPath = TKX.JoinPath

// GetUserInput 获取键盘输入，不太可靠
func (pA *TK) GetUserInput(promptA string) string {
	if promptA != "" {
		fmt.Print(promptA)
	}

	var textT string
	_, errT := fmt.Scanln(&textT)
	if errT != nil {
		return GenerateErrorString(errT.Error())
	}

	return textT
}

var GetUserInput = TKX.GetUserInput

// GetInputf like GetInput, but allows using printf for prompt string
func (pA *TK) GetInputf(formatA string, aA ...interface{}) string {
	fmt.Printf(formatA, aA...)

	// var stdinBufferedReaderT *bufio.Reader
	var stdinBufferedScannerT *bufio.Scanner

	stdinBufferedScannerT = bufio.NewScanner(os.Stdin)

	if stdinBufferedScannerT.Scan() {
		rStrT := stdinBufferedScannerT.Text()

		errT := stdinBufferedScannerT.Err()
		if errT != nil {
			if errT == io.EOF {
				return GenerateErrorStringF("EOF", rStrT)
			}

			return GenerateErrorStringF(errT.Error())
		}

		return rStrT
	}

	errT := stdinBufferedScannerT.Err()
	if errT != nil {
		if errT == io.EOF {
			return GenerateErrorStringF("EOF", stdinBufferedScanner.Text())
		}

		return GenerateErrorStringF(errT.Error())
	}

	return GenerateErrorStringF("EOF")
}

var GetInputf = TKX.GetInputf

func (pA *TK) GetInput(argsA ...string) string {
	prefixT := GetSwitch(argsA, "-prompt=", "")

	if prefixT != "" {
		fmt.Print(prefixT)
	}

	// var stdinBufferedReaderT *bufio.Reader
	var stdinBufferedScannerT *bufio.Scanner

	stdinBufferedScannerT = bufio.NewScanner(os.Stdin)

	// defaultT := GetSwitch(argsA, "-default=", "")

	// if defaultT != "" {
	// 	stdinBufferedScannerT.Buffer([]byte(defaultT), bufio.MaxScanTokenSize)
	// }

	if stdinBufferedScannerT.Scan() {
		rStrT := stdinBufferedScannerT.Text()

		errT := stdinBufferedScannerT.Err()
		if errT != nil {
			if errT == io.EOF {
				return GenerateErrorStringF("EOF", rStrT)
			}

			return GenerateErrorStringF(errT.Error())
		}

		return rStrT
	}

	errT := stdinBufferedScannerT.Err()
	if errT != nil {
		if errT == io.EOF {
			return GenerateErrorStringF("EOF", stdinBufferedScanner.Text())
		}

		return GenerateErrorStringF(errT.Error())
	}

	return GenerateErrorStringF("EOF")
}

var GetInput = TKX.GetInput

var stdinBufferedReader *bufio.Reader
var stdinBufferedScanner *bufio.Scanner

// GetInputBufferedScan 获取键盘输入
func (pA *TK) GetInputBufferedScan() string {
	if stdinBufferedScanner == nil {
		stdinBufferedScanner = bufio.NewScanner(os.Stdin)
	}

	if stdinBufferedScanner.Scan() {
		rStrT := stdinBufferedScanner.Text()

		errT := stdinBufferedScanner.Err()
		if errT != nil {
			if errT == io.EOF {
				return GenerateErrorStringF("EOF", rStrT)
			}

			return GenerateErrorStringF(errT.Error())
		}

		return rStrT
	}

	errT := stdinBufferedScanner.Err()
	if errT != nil {
		if errT == io.EOF {
			return GenerateErrorStringF("EOF", stdinBufferedScanner.Text())
		}

		return GenerateErrorStringF(errT.Error())
	}

	return GenerateErrorStringF("EOF")
}

var GetInputBufferedScan = TKX.GetInputBufferedScan

func (pA *TK) GetInputPasswordf(formatA string, aA ...interface{}) string {
	fmt.Printf(formatA, aA...)

	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return ErrStrF("failed to get password: %v", err)
	}

	return string(bytePassword)
}

var GetInputPasswordf = TKX.GetInputPasswordf

func (pA *TK) Sleep(secA float64) {
	time.Sleep(time.Duration(secA*1000) * time.Millisecond)
}

var Sleep = TKX.Sleep

func (pA *TK) SleepSeconds(secA int) {
	time.Sleep(time.Duration(secA) * time.Second)
}

var SleepSeconds = TKX.SleepSeconds

func (pA *TK) SleepMilliSeconds(msA int) {
	time.Sleep(time.Duration(msA) * time.Millisecond)
}

var SleepMilliSeconds = TKX.SleepMilliSeconds

func (pA *TK) GetRuntimeStack() string {
	return string(debug.Stack())
}

var GetRuntimeStack = TKX.GetRuntimeStack

func (pA *TK) GetOSName() string {
	return runtime.GOOS
}

var GetOSName = TKX.GetOSName

func (pA *TK) GetOSArch() string {
	return runtime.GOARCH
}

var GetOSArch = TKX.GetOSArch

func (pA *TK) GetCurrentDir() string {
	strT, errT := os.Getwd()
	if errT != nil {
		strT, errT = filepath.Abs(".")
		if errT != nil {
			return "."
		}
	}

	return strT
}

var GetCurrentDir = TKX.GetCurrentDir

func (pA *TK) SetCurrentDir(dirA string) string {
	errT := os.Chdir(dirA)
	if errT != nil {
		return ErrToStr(errT)
	}

	return ""
}

var SetCurrentDir = TKX.SetCurrentDir

func (pA *TK) GetHomeDir() string {
	strT, errT := os.UserHomeDir()
	if errT != nil {
		strT, errT = filepath.Abs(".")
		if errT != nil {
			return "."
		}
	}

	return strT
}

var GetHomeDir = TKX.GetHomeDir

func (pA *TK) GetTempDir() string {
	return os.TempDir()
}

var GetTempDir = TKX.GetTempDir

func (pA *TK) GetApplicationPath() string {
	file, _ := exec.LookPath(os.Args[0])
	pathT, _ := filepath.Abs(file)

	dirT, _ := filepath.Split(pathT)

	return dirT
}

var GetApplicationPath = TKX.GetApplicationPath

func (pA *TK) GetExecutablePath() string {
	pathT, errT := os.Executable()

	if errT != nil {
		return ErrorToString(errT)
	}

	return pathT
}

var GetExecutablePath = TKX.GetExecutablePath

func (pA *TK) GetFileAbs(fileA string) string {
	pathT, errT := filepath.Abs(fileA)

	if errT != nil {
		return ErrToErrStr(errT)
	}

	return pathT
}

var GetFileAbs = TKX.GetFileAbs

func (pA *TK) EnsureMakeDirs(pathA string) string {
	if !IfFileExists(pathA) {
		os.MkdirAll(pathA, 0777)
		return ""
	} else {
		if IsDirectory(pathA) {
			return ""
		} else {
			return GenerateErrorString("a file with same name exists")
		}
	}
}

var EnsureMakeDirs = TKX.EnsureMakeDirs

func (pA *TK) EnsureMakeDirsE(pathA string) error {
	if !IfFileExists(pathA) {
		os.MkdirAll(pathA, 0777)

		if !IfFileExists(pathA) {
			return fmt.Errorf("failed to make directory")
		}
		return nil
	} else {
		if IsDirectory(pathA) {
			return nil
		} else {
			return fmt.Errorf("a file with same name exists")
		}
	}
}

var EnsureMakeDirsE = TKX.EnsureMakeDirsE

// func GetCurrentThreadID() int {
// 	var user32 *syscall.DLL
// 	var GetCurrentThreadId *syscall.Proc
// 	var err error

// 	user32, err = syscall.LoadDLL("Kernel32.dll")
// 	if err != nil {
// 		fmt.Printf("syscall.LoadDLL fail: %v\n", err.Error())
// 		return 0
// 	}
// 	GetCurrentThreadId, err = user32.FindProc("GetCurrentThreadId")
// 	if err != nil {
// 		fmt.Printf("user32.FindProc fail: %v\n", err.Error())
// 		return 0
// 	}

// 	var pid uintptr
// 	pid, _, err = GetCurrentThreadId.Call()

// 	return int(pid)
// }

// 命令行分析

// AnalyzeCommandLineParamter 分解命令行参数，注意如果要带双引号，需要从开始到结束都括上，例如save "-fileName=abc.txt"，而不是save -fileName="abc.txt"
func (pA *TK) AnalyzeCommandLineParamter(cmdLineA string) []string {
	return regexp.MustCompile("( |\\\".*?\\\"|'.*?')").Split(cmdLineA, -1)
}

var AnalyzeCommandLineParamter = TKX.AnalyzeCommandLineParamter

// GetParameterByIndexWithDefaultValue 按顺序序号获取命令行参数，其中0代表第一个参数，也就是软件名称或者命令名称，1开始才是第一个参数，注意参数不包括开关，即类似-verbose=true这样的
func (pA *TK) GetParameterByIndexWithDefaultValue(argsA []string, idxA int, defaultA string) string {
	if idxA == -1 {
		idxA = 1
	}

	if (idxA >= len(argsA)) || (idxA < 0) {
		return defaultA
	}

	var cnt int
	for _, argT := range argsA {
		if StartsWith(argT, "-") {
			continue
		}

		if cnt == idxA {
			if StartsWith(argT, "\"") && EndsWith(argT, "\"") {
				return argT[1 : len(argT)-1]
			}

			return argT
		}

		cnt++

	}

	return defaultA
}

var GetParameterByIndexWithDefaultValue = TKX.GetParameterByIndexWithDefaultValue

func (pA *TK) GetParamI(argsA []interface{}, idxA int, defaultA ...string) string {
	defaultT := ""

	if len(defaultA) > 0 {
		defaultT = defaultA[0]
	}

	if idxA == -1 {
		idxA = 1
	}

	if (idxA >= len(argsA)) || (idxA < 0) {
		return defaultT
	}

	var cnt int
	for _, argT := range argsA {
		strT := ToStr(argT)
		if StartsWith(strT, "-") {
			continue
		}

		if cnt == idxA {
			if StartsWith(strT, "\"") && EndsWith(strT, "\"") {
				return strT[1 : len(strT)-1]
			}

			return strT
		}

		cnt++

	}

	return defaultT
}

var GetParamI = TKX.GetParamI

func (pA *TK) GetParameter(argsA []string, idxA int) string {
	return GetParameterByIndexWithDefaultValue(argsA, idxA, ErrStrF("failed"))
}

var GetParameter = TKX.GetParameter

func (pA *TK) GetParam(argsA []string, optsA ...interface{}) string {
	idxT := 1
	defaultT := ""

	for _, v := range optsA {
		if _, ok := v.(int); ok {
			idxT = v.(int)
		}

		if _, ok := v.(string); ok {
			defaultT = v.(string)
		}
	}

	return GetParameterByIndexWithDefaultValue(argsA, idxT, defaultT)
}

var GetParam = TKX.GetParam

// GetAllParameters 获取命令行参数中所有非开关参数
func (pA *TK) GetAllParameters(argsA []string) []string {
	aryT := make([]string, 0, len(argsA))

	for _, argT := range argsA {
		if StartsWith(argT, "-") {
			continue
		}

		aryT = append(aryT, argT)
	}

	return aryT
}

var GetAllParameters = TKX.GetAllParameters

func (pA *TK) GetAllOSParameters() []string {
	return GetAllParameters(os.Args)
}

var GetAllOSParameters = TKX.GetAllOSParameters

// GetAllSwitches 获取命令行参数中所有开关参数
func (pA *TK) GetAllSwitches(argsA []string) []string {
	aryT := make([]string, 0, len(argsA))

	for _, argT := range argsA {
		if !StartsWith(argT, "-") {
			continue
		}

		aryT = append(aryT, argT)
	}

	return aryT
}

var GetAllSwitches = TKX.GetAllSwitches

// ParseCommandLine 分析命令行字符串，类似os.Args的获取过程
func (pA *TK) ParseCommandLine(commandA string) ([]string, error) {
	var args []string

	state := "start"
	current := ""
	quote := "\""
	escapeNext := false

	command := []rune(commandA)

	for i := 0; i < len(command); i++ {
		c := command[i]

		if escapeNext {
			current += string(c)
			escapeNext = false
			continue
		}

		if c == '\\' {
			escapeNext = true
			continue
		}

		if state == "quotes" {
			if string(c) != quote {
				current += string(c)
			} else {
				args = append(args, current)
				current = ""
				state = "start"
			}
			continue
		}

		if c == '"' || c == '\'' || c == '`' {
			state = "quotes"
			quote = string(c)
			continue
		}

		if state == "arg" {
			if c == ' ' || c == '\t' {
				args = append(args, current)
				current = ""
				state = "start"
			} else {
				current += string(c)
			}
			// Pl("state: %v, current: %v, args: %v", state, current, args)
			continue
		}

		if c != ' ' && c != '\t' {
			state = "arg"
			current += string(c)
		}
	}

	if state == "quotes" {
		return []string{}, fmt.Errorf("Unclosed quote in command line: %v", command)
	}

	if current != "" {
		args = append(args, current)
	}

	return args, nil
}

var ParseCommandLine = TKX.ParseCommandLine

func (pA *TK) ParseCommandLineCompact(commandA string) []string {
	rsT, errT := ParseCommandLine(commandA)

	if errT != nil {
		return []string{commandA}
	}

	return rsT
}

var ParseCommandLineCompact = TKX.ParseCommandLineCompact

// GetSwitchWithDefaultValue 获取命令行参数中的开关，用法：tmps := tk.GetSwitchWithDefaultValue(args, "-verbose=", "false")
func (pA *TK) GetSwitchWithDefaultValue(argsA []string, switchStrA string, defaultA string) string {
	if argsA == nil {
		return defaultA
	}

	if len(argsA) < 1 {
		return defaultA
	}

	tmpStrT := ""
	for _, argT := range argsA {
		if StartsWith(argT, switchStrA) {
			tmpStrT = argT[len(switchStrA):]
			if StartsWith(tmpStrT, "\"") && EndsWith(tmpStrT, "\"") {
				return tmpStrT[1 : len(tmpStrT)-1]
			}

			return tmpStrT
		}

	}

	return defaultA

}

var GetSwitchWithDefaultValue = TKX.GetSwitchWithDefaultValue

func (pA *TK) GetSwitch(argsA []string, switchStrA string, defaultA ...string) string {

	ifDefaultT := true
	var defaultT string

	if defaultA == nil || len(defaultA) < 1 {
		ifDefaultT = false
	}

	if ifDefaultT {
		defaultT = defaultA[0]
	}

	if argsA == nil {
		if ifDefaultT {
			return defaultT
		}
		return ErrStr("nil input")
	}

	if len(argsA) < 1 {
		if ifDefaultT {
			return defaultT
		}
		return ErrStr("not exists")
	}

	tmpStrT := ""
	for _, argT := range argsA {
		if StartsWith(argT, switchStrA) {
			tmpStrT = argT[len(switchStrA):]
			if StartsWith(tmpStrT, "\"") && EndsWith(tmpStrT, "\"") {
				return tmpStrT[1 : len(tmpStrT)-1]
			}

			return tmpStrT
		}

	}

	if ifDefaultT {
		return defaultT
	}
	return ErrStr("not exists")
}

var GetSwitch = TKX.GetSwitch

func (pA *TK) GetSwitchI(argsA []interface{}, switchStrA string, defaultA ...string) string {
	var defaultT string

	if len(defaultA) > 0 {
		defaultT = defaultA[0]
	}

	if argsA == nil {
		return defaultT
	}

	if len(argsA) < 1 {
		return defaultT
	}

	tmpStrT := ""
	for _, argIT := range argsA {
		argT, ok := argIT.(string)
		if !ok {
			continue
		}

		if StartsWith(argT, switchStrA) {
			tmpStrT = argT[len(switchStrA):]
			if StartsWith(tmpStrT, "\"") && EndsWith(tmpStrT, "\"") {
				return tmpStrT[1 : len(tmpStrT)-1]
			}

			return tmpStrT
		}

	}

	return defaultT

}

var GetSwitchI = TKX.GetSwitchI

// GetSwitchWithDefaultIntValue 与GetSwitchWithDefaultValue类似，返回一个整数
func (pA *TK) GetSwitchWithDefaultIntValue(argsA []string, switchStrA string, defaultA ...int) int {
	defaultT := 0

	if len(defaultA) > 0 {
		defaultT = defaultA[0]
	}

	tmpStrT := GetSwitchWithDefaultValue(argsA, switchStrA, IntToStr(defaultT))

	return StrToIntWithDefaultValue(tmpStrT, defaultT)
}

var GetSwitchWithDefaultIntValue = TKX.GetSwitchWithDefaultIntValue

func (pA *TK) GetSwitchWithDefaultFloatValue(argsA []string, switchStrA string, defaultA ...float64) float64 {
	defaultT := 0.0

	if len(defaultA) > 0 {
		defaultT = defaultA[0]
	}

	tmpStrT := GetSwitchWithDefaultValue(argsA, switchStrA, Float64ToStr(defaultT))

	return StrToFloat64WithDefaultValue(tmpStrT, defaultT)
}

var GetSwitchWithDefaultFloatValue = TKX.GetSwitchWithDefaultFloatValue

func (pA *TK) GetSwitchWithDefaultInt64Value(argsA []string, switchStrA string, defaultA int64) int64 {
	tmpStrT := GetSwitchWithDefaultValue(argsA, switchStrA, string(defaultA))

	return StrToInt64WithDefaultValue(tmpStrT, defaultA)
}

var GetSwitchWithDefaultInt64Value = TKX.GetSwitchWithDefaultInt64Value

// IfSwitchExists 判断命令行参数中是否存在开关，用法：flag := IfSwitchExists(args, "-restart")
// func (pA *TK) IfSwitchExists(argsA []string, switchStrA string) bool {
// 	if argsA == nil {
// 		return false
// 	}

// 	if len(argsA) < 1 {
// 		return false
// 	}

// 	for _, argT := range argsA {
// 		if StartsWith(argT, switchStrA) {
// 			return true
// 		}

// 	}

// 	return false
// }

// IfSwitchExistsWhole 判断命令行参数中是否存在开关（完整的，），用法：flag := IfSwitchExistsWhole(args, "-restart")
func (pA *TK) IfSwitchExists(argsA []string, switchStrA string) bool {
	if argsA == nil {
		return false
	}

	if len(argsA) < 1 {
		return false
	}

	for _, argT := range argsA {
		if argT == switchStrA {
			return true
		}
	}

	return false
}

var IfSwitchExistsWhole = TKX.IfSwitchExists

var IfSwitchExists = TKX.IfSwitchExists

func (pA *TK) IfSwitchExistsWholeI(argsA []interface{}, switchStrA string) bool {
	if argsA == nil {
		return false
	}

	if len(argsA) < 1 {
		return false
	}

	for _, argT := range argsA {
		s, ok := argT.(string)
		if !ok {
			continue
		}

		if s == switchStrA {
			return true
		}

	}

	return false
}

var IfSwitchExistsWholeI = TKX.IfSwitchExistsWholeI

// 各种转换 conversion related

func (pA *TK) BoolToInt(b bool) int {
	if b {
		return 1
	}

	return 0
}

var BoolToInt = TKX.BoolToInt

func (pA *TK) AnyArrayToStringArray(vA []interface{}) []string {
	mssT := make([]string, len(vA))
	for i, v := range vA {
		vv := ToStr(v)
		mssT[i] = vv
	}
	return mssT
}

var AnyArrayToStringArray = TKX.AnyArrayToStringArray

func (pA *TK) StringArrayToAnyArray(vA []string) []interface{} {
	mssT := make([]interface{}, len(vA))
	for i, v := range vA {
		mssT[i] = v
	}
	return mssT
}

var StringArrayToAnyArray = TKX.StringArrayToAnyArray

func (pA *TK) MSI2MSS(vA map[string]interface{}) map[string]string {
	mssT := make(map[string]string, len(vA))
	for k, v := range vA {
		vv, ok := v.(string)
		if ok {
			mssT[k] = vv
		}
	}
	return mssT
}

var MSI2MSS = TKX.MSI2MSS

func (pA *TK) MSS2MSI(vA map[string]string) map[string]interface{} {
	msiT := make(map[string]interface{}, len(vA))
	for k, v := range vA {
		msiT[k] = v
	}
	return msiT
}

var MSS2MSI = TKX.MSS2MSI

func (pA *TK) ToInterface(vA interface{}) interface{} {
	return vA
}

var ToInterface = TKX.ToInterface

func (pA *TK) ToVar(pointerA *interface{}) interface{} {
	return *pointerA
}

var ToVar = TKX.ToVar

func (pA *TK) ToPointer(pointerA *interface{}, typeA ...string) interface{} {

	typeT := ""
	if len(typeA) > 0 {
		typeT = Trim(typeA[0])
	}

	if typeT != "" {
		switch typeT {
		case "int":
			return (*int)(unsafe.Pointer(pointerA))
		case "[]int":
			return (*[]int)(unsafe.Pointer(pointerA))
		case "uint":
			return (*uint)(unsafe.Pointer(pointerA))
		case "uint8": // byte
			return (*uint8)(unsafe.Pointer(pointerA))
		case "[]uint8": // []byte
			return (*[]uint8)(unsafe.Pointer(pointerA))
		case "uint16":
			return (*uint16)(unsafe.Pointer(pointerA))
		case "uint32":
			return (*uint32)(unsafe.Pointer(pointerA))
		case "uint64":
			return (*uint64)(unsafe.Pointer(pointerA))
		case "int8":
			return (*int8)(unsafe.Pointer(pointerA))
		case "int16":
			return (*int16)(unsafe.Pointer(pointerA))
		case "int32": // rune
			return (*int32)(unsafe.Pointer(pointerA))
		case "[]int32": // []rune
			return (*[]int32)(unsafe.Pointer(pointerA))
		case "int64":
			return (*int64)(unsafe.Pointer(pointerA))
		case "[]int64":
			return (*[]int64)(unsafe.Pointer(pointerA))
		case "complex64":
			return (*complex64)(unsafe.Pointer(pointerA))
		case "[]complex64":
			return (*[]complex64)(unsafe.Pointer(pointerA))
		case "complex128":
			return (*complex128)(unsafe.Pointer(pointerA))
		case "[]complex128":
			return (*[]complex128)(unsafe.Pointer(pointerA))
		case "float32":
			return (*float32)(unsafe.Pointer(pointerA))
		case "[]float32":
			return (*[]float32)(unsafe.Pointer(pointerA))
		case "float64":
			return (*float64)(unsafe.Pointer(pointerA))
		case "[]float64":
			return (*[]float64)(unsafe.Pointer(pointerA))
		case "bool":
			return (*bool)(unsafe.Pointer(pointerA))
		case "string":
			return (*string)(unsafe.Pointer(pointerA))
		case "[]string":
			return (*[]string)(unsafe.Pointer(pointerA))
		case "map[string]string":
			return (*map[string]string)(unsafe.Pointer(pointerA))
		case "[]map[string]string":
			return (*[]map[string]string)(unsafe.Pointer(pointerA))
		case "time.Time", "time":
			return (*time.Time)(unsafe.Pointer(pointerA))
		case "[]interface{}", "[]var", "[]":
			return (*[]interface{})(unsafe.Pointer(pointerA))
		case "map[string]interface{}", "msi":
			return (*map[string]interface{})(unsafe.Pointer(pointerA))
		case "[]map[string]interface{}", "[]msi":
			return (*[]map[string]interface{})(unsafe.Pointer(pointerA))
		case "*interface{}", "*var":
			return (*interface{})(unsafe.Pointer(pointerA))
		case "interface{}", "var", "":
			return pointerA
		}

		return pointerA
	}

	v := *pointerA

	Pl("type: %T, v: %v", v, v)

	switch v.(type) {
	case int:
		return (*int)(unsafe.Pointer(pointerA))
	case []int:
		return (*[]int)(unsafe.Pointer(pointerA))
	case uint:
		return (*uint)(unsafe.Pointer(pointerA))
	case uint8: // byte
		return (*uint8)(unsafe.Pointer(pointerA))
	case []uint8: // []byte
		return (*[]uint8)(unsafe.Pointer(pointerA))
	case uint16:
		return (*uint16)(unsafe.Pointer(pointerA))
	case uint32:
		return (*uint32)(unsafe.Pointer(pointerA))
	case uint64:
		return (*uint64)(unsafe.Pointer(pointerA))
	case int8:
		return (*int8)(unsafe.Pointer(pointerA))
	case int16:
		return (*int16)(unsafe.Pointer(pointerA))
	case int32: // rune
		return (*int32)(unsafe.Pointer(pointerA))
	case []int32: // []rune
		return (*[]int32)(unsafe.Pointer(pointerA))
	case int64:
		return (*int64)(unsafe.Pointer(pointerA))
	case []int64:
		return (*[]int64)(unsafe.Pointer(pointerA))
	case complex64:
		return (*complex64)(unsafe.Pointer(pointerA))
	case []complex64:
		return (*[]complex64)(unsafe.Pointer(pointerA))
	case complex128:
		return (*complex128)(unsafe.Pointer(pointerA))
	case []complex128:
		return (*[]complex128)(unsafe.Pointer(pointerA))
	case float32:
		return (*float32)(unsafe.Pointer(pointerA))
	case []float32:
		return (*[]float32)(unsafe.Pointer(pointerA))
	case float64:
		return (*float64)(unsafe.Pointer(pointerA))
	case []float64:
		return (*[]float64)(unsafe.Pointer(pointerA))
	case bool:
		return (*bool)(unsafe.Pointer(pointerA))
	case string:
		return (*string)(unsafe.Pointer(pointerA))
	case []string:
		Pl("%T %v %#v %#v", pointerA, pointerA, pointerA, *pointerA)
		tmp := unsafe.Pointer(pointerA)
		Pl("%T %v %#v", tmp, tmp, *(*[]string)(tmp))
		return (*[]string)(unsafe.Pointer(pointerA))
	case map[string]string:
		return (*map[string]string)(unsafe.Pointer(pointerA))
	case []map[string]string:
		return (*[]map[string]string)(unsafe.Pointer(pointerA))
	case time.Time:
		return (*time.Time)(unsafe.Pointer(pointerA))
	case []interface{}:
		return (*[]interface{})(unsafe.Pointer(pointerA))
	case map[string]interface{}:
		return (*map[string]interface{})(unsafe.Pointer(pointerA))
	case []map[string]interface{}:
		return (*[]map[string]interface{})(unsafe.Pointer(pointerA))
	case *interface{}:
		return (*interface{})(unsafe.Pointer(pointerA))
	case interface{}:
		return pointerA
	}

	return pointerA
}

var ToPointer = TKX.ToPointer

func (pA *TK) GetSystemEndian() binary.ByteOrder {
	var n uint64 = 1

	buf := make([]byte, binary.MaxVarintLen64)

	binary.PutUvarint(buf, n)

	if buf[0] == 1 {
		return binary.LittleEndian
	}

	return binary.BigEndian
}

var GetSystemEndian = TKX.GetSystemEndian

// GetSystemEndianInt return 1 indicate BigEndian, 0 indicate LitteEndian
func (pA *TK) GetSystemEndianInt() int {
	var n uint64 = 1

	buf := make([]byte, binary.MaxVarintLen64)

	binary.PutUvarint(buf, n)

	if buf[0] == 1 {
		return 0
	}

	return 1
}

var GetSystemEndianInt = TKX.GetSystemEndianInt

func (pA *TK) CompareBytes(buf1 []byte, buf2 []byte, limitA ...int) [][]int {
	limitT := 0

	if len(limitA) > 0 {
		limitT = limitA[0]
	}

	len1 := len(buf1)

	len2 := len(buf2)

	lenT := len1

	if lenT < len2 {
		lenT = len2
	}

	var c1 int
	var c2 int

	diffBufT := make([][]int, 0, 100)

	countT := 0

	for i := 0; i < lenT; i++ {
		if (limitT > 0) && (countT >= limitT) {
			break
		}

		if i >= len1 {
			c1 = 256
		} else {
			c1 = int(buf1[i])
		}

		if i >= len2 {
			c2 = 256
		} else {
			c2 = int(buf2[i])
		}

		// if i%1000000 == 0 {
		// 	pl("%v(%v) - %v %v", i, i, c1, c2)
		// }

		if c1 != c2 {
			diffBufT = append(diffBufT, []int{i, c1, c2})
			// pl("%v(%v) - %v %v", i, i, c1, c2)
			countT++
		}

	}

	return diffBufT
}

var CompareBytes = TKX.CompareBytes

func (pA *TK) BytesToData(bytesA []byte, dataA interface{}, optsA ...string) error {
	defaultEndianT := GetSystemEndian()

	if len(optsA) > 0 {
		endianStrA := GetSwitch(optsA, "-endian=", "")

		if endianStrA != "" {
			if StartsWith(endianStrA, "B") {
				defaultEndianT = binary.BigEndian
			} else if StartsWith(endianStrA, "L") {
				defaultEndianT = binary.LittleEndian
			}
		}
	}

	bufT := bytes.NewBuffer(bytesA)

	errT := binary.Read(bufT, defaultEndianT, dataA)

	return errT
}

var BytesToData = TKX.BytesToData

func (pA *TK) DataToBytes(dataA interface{}, optsA ...string) interface{} {
	nv1, ok := dataA.([]byte)

	if ok {
		return nv1
	}

	nv2, ok := dataA.(string)

	if ok {
		return []byte(nv2)
	}

	nv3, ok := dataA.([]rune)

	if ok {
		return []byte(string(nv3))
	}

	defaultEndianT := GetSystemEndian()

	if len(optsA) > 0 {
		endianStrA := GetSwitch(optsA, "-endian=", "")

		if endianStrA != "" {
			if StartsWith(endianStrA, "B") {
				defaultEndianT = binary.BigEndian
			} else if StartsWith(endianStrA, "L") {
				defaultEndianT = binary.LittleEndian
			}
		}
	}

	bufT := new(bytes.Buffer)

	var errT error

	nv4, ok := dataA.(int)

	if ok {
		errT = binary.Write(bufT, defaultEndianT, uint64(nv4))

		if errT != nil {
			return errT
		}

		return bufT.Bytes()

	}

	errT = binary.Write(bufT, defaultEndianT, dataA)

	if errT != nil {
		return errT
	}

	return bufT.Bytes()
}

var DataToBytes = TKX.DataToBytes

func (pA *TK) ToByte(vA interface{}, defaultA ...byte) byte {
	var defaultT byte = 0

	if len(defaultA) > 0 {
		defaultT = defaultA[0]
	}

	switch vT := vA.(type) {
	case bool:
		if vT {
			return 1
		}

		return 0
	case byte:
		return vT
	case int8:
		return byte(vT)
	case int16:
		return byte(vT)
	case int32:
		return byte(vT)
	case int:
		return byte(vT)
	case int64:
		return byte(vT)
	case uint16:
		return byte(vT)
	case uint32:
		return byte(vT)
	case uint:
		return byte(vT)
	case uint64:
		return byte(vT)
	case float32:
		return byte(vT)
	case float64:
		return byte(vT)
	case string:
		return byte(StrToInt(vT, int(defaultT)))
	case []byte:
		return vT[0]
	case []rune:
		return byte(vT[0])
	}

	var bufT bytes.Buffer

	errT := binary.Write(&bufT, GetSystemEndian(), vA)

	if errT != nil {
		return defaultT
	}

	bytesT := bufT.Bytes()

	if len(bytesT) < 1 {
		return defaultT
	}

	return bytesT[0]
}

var ToByte = TKX.ToByte

func (pA *TK) NilToEmptyStr(vA interface{}) string {
	if vA == nil {
		return ""
	}

	switch vA.(type) {
	case string:
		return vA.(string)
	case error:
		return ""
	case float64:
		return Float64ToStr(vA.(float64))
	default:
		return fmt.Sprintf("%v", vA)
	}
}

var NilToEmptyStr = TKX.NilToEmptyStr

func (pA *TK) NilEmptyStringNoFloat(vA interface{}) string {
	if vA == nil {
		return ""
	}

	vf, ok := vA.(float64)
	if ok {
		return ToStr(int(vf))
	}

	vs, ok := vA.(string)
	if ok {
		return vs
	}

	_, ok = vA.(error)
	if ok {
		return ""
	}

	return fmt.Sprintf("%v", vA)
}

var NilEmptyStringNoFloat = TKX.NilEmptyStringNoFloat

func (pA *TK) StrToBool(strA string) bool {
	lowerStr := strings.ToLower(strA)
	if lowerStr == "yes" || lowerStr == "true" {
		return true
	}

	if lowerStr == "no" || lowerStr == "false" {
		return false
	}

	return false
}

var StrToBool = TKX.StrToBool

func (pA *TK) StrToBytes(strA string) []byte {
	return []byte(strA)
}

var StrToBytes = TKX.StrToBytes

// not completed, like javascript, no value -> false, other true
func (pA *TK) ToBool(vA interface{}) bool {
	if IsNil(vA) {
		return false
	}

	boolT, ok := vA.(bool)

	if ok {
		return boolT
	}

	strT, ok := vA.(string)

	if ok {
		lowerStr := strings.ToLower(strT)
		if lowerStr == "yes" || lowerStr == "true" {
			return true
		}

		if lowerStr == "no" || lowerStr == "false" {
			return false
		}

		return false
	}

	// Z...
	return true
}

var ToBool = TKX.ToBool

func (pA *TK) ToBoolWithDefaultValue(vA interface{}, defaultA ...bool) bool {
	defaultT := false

	if (defaultA != nil) && (len(defaultA) > 0) {
		defaultT = defaultA[0]
	}

	if IsNil(vA) {
		return defaultT
	}

	boolT, ok := vA.(bool)

	if ok {
		return boolT
	}

	strT, ok := vA.(string)

	if ok {
		lowerStr := strings.ToLower(strT)
		if lowerStr == "yes" || lowerStr == "true" {
			return true
		}

		if lowerStr == "no" || lowerStr == "false" {
			return false
		}

		return defaultT
	}

	// Z...
	return defaultT
}

var ToBoolWithDefaultValue = TKX.ToBoolWithDefaultValue

func (pA *TK) BoolToStr(vA bool) string {
	if vA {
		return "true"
	} else {
		return "false"
	}
}

var BoolToStr = TKX.BoolToStr

func (pA *TK) ByteToHex(byteA byte) string {
	return Spr("%X", byteA)
}

var ByteToHex = TKX.ByteToHex

// IntToStr 整形转字符串
func (pA *TK) IntToStr(intA int) string {
	return strconv.Itoa(intA)
}

var IntToStr = TKX.IntToStr

func (pA *TK) IntToStrX(n interface{}, defaultA ...string) string {
	var defaultT string = ""
	if (defaultA != nil) && (len(defaultA) > 0) {
		defaultT = defaultA[0]
	}

	switch nv := n.(type) {
	case int:
		return fmt.Sprintf("%v", nv)
	case int8:
		return fmt.Sprintf("%v", nv)
	case int16:
		return fmt.Sprintf("%v", nv)
	case int32:
		return fmt.Sprintf("%v", nv)
	case int64:
		return strconv.FormatInt(nv, 10)
	case float64:
		return Float64ToStr(nv)
	case float32:
		tmps := fmt.Sprintf("%f", nv)
		if Contains(tmps, ".") {
			tmps = strings.TrimRight(tmps, "0")
			tmps = strings.TrimRight(tmps, ".")
		}

		return tmps
	case string:
		nT, errT := strconv.ParseInt(nv, 10, 0)
		if errT != nil {
			return defaultT
		}

		return fmt.Sprintf("%v", nT)
	default:
		nT, errT := strconv.ParseInt(fmt.Sprintf("%v", nv), 10, 0)
		if errT != nil {
			return defaultT
		}

		return fmt.Sprintf("%v", nT)
	}

}

var IntToStrX = TKX.IntToStrX

func (pA *TK) Int64ToStr(intA int64) string {
	return strconv.FormatInt(intA, 10)
}

var Int64ToStr = TKX.Int64ToStr

func (pA *TK) ToStr(v interface{}) string {
	// Pl("ToStr: (%T) %#v", v, v)
	switch nv := v.(type) {
	case float64, float32:
		tmps := fmt.Sprintf("%f", nv)
		if Contains(tmps, ".") {
			tmps = strings.TrimRight(tmps, "0")
			tmps = strings.TrimRight(tmps, ".")
		}

		return tmps
	case int64:
		return strconv.FormatInt(nv, 10)
	case int:
		return IntToStr(nv)
	case string:
		return nv
	case []byte:
		return string(nv)
	case []rune:
		return string(nv)
	case *strings.Builder:
		return nv.String()
	case *bytes.Buffer:
		return string(nv.Bytes())
	case strings.Builder:
		return nv.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}

var ToStr = TKX.ToStr

func (pA *TK) ToFloat(v interface{}, defaultA ...float64) (result float64) {
	var defaultT float64

	if defaultA == nil || len(defaultA) < 1 {
		defaultT = 0.0
	} else {
		defaultT = defaultA[0]
	}

	defer func() {
		r := recover()

		if r != nil {
			result = defaultT
			return
		}
	}()

	switch nv := v.(type) {
	case bool:
		if v.(bool) {
			result = float64(1)
		} else {
			result = float64(0)
		}

		return
	case int:
		result = float64(v.(int))
		return
	case int8:
		result = float64(v.(int8))
		return
	case int16:
		result = float64(v.(int16))
		return
	case int32:
		result = float64(v.(int32))
		return
	case int64:
		result = float64(v.(int64))
		return
	case uint:
		result = float64(v.(uint))
		return
	case uint8:
		result = float64(v.(uint8))
		return
	case uint16:
		result = float64(v.(uint16))
		return
	case uint32:
		result = float64(v.(uint32))
		return
	case uint64:
		result = float64(v.(uint64))
		return
	case float64:
		result = v.(float64)
		return
	case float32:
		result = float64(v.(float32))
		return
	case string:
		nT, errT := strconv.ParseFloat(strings.TrimSpace(v.(string)), 64)
		if errT != nil {
			result = defaultT
			return
		}

		result = nT
		return
	case *big.Int:
		result = float64(nv.Int64())
		return
	case *big.Float:
		fv, _ := nv.Float64()
		result = fv
		return
	default:
		nT, errT := strconv.ParseFloat(fmt.Sprintf("%v", v), 64)
		if errT != nil {
			result = defaultT
			return
		}

		result = nT
		return
	}
}

var ToFloat = TKX.ToFloat

func (pA *TK) ToInt(v interface{}, defaultA ...int) (result int) {
	var defaultT int

	if defaultA == nil || len(defaultA) < 1 {
		defaultT = 0
	} else {
		defaultT = defaultA[0]
	}

	defer func() {
		r := recover()

		if r != nil {
			result = defaultT
			return
		}
	}()

	switch nv := v.(type) {
	case bool:
		if nv {
			result = 1
		} else {
			result = 0
		}

		return
	case int:
		result = nv
		return
	case int8:
		result = int(nv)
		return
	case int16:
		result = int(nv)
		return
	case int32:
		result = int(nv)
		return
	case int64:
		result = int(nv)
		return
	case uint:
		result = int(nv)
		return
	case uint8:
		result = int(nv)
		return
	case uint16:
		result = int(nv)
		return
	case uint32:
		result = int(nv)
		return
	case uint64:
		result = int(nv)
		return
	case float64:
		result = int(nv)
		return
	case float32:
		result = int(nv)
		return
	case string:
		nT, errT := strconv.ParseInt(strings.TrimSpace(nv), 10, 0)
		if errT != nil {
			result = defaultT
			return
		}

		result = int(nT)
		return
	case *big.Int:
		result = int(nv.Int64())
		return
	case *big.Float:
		rs, _ := nv.Int64()
		result = int(rs)
		return
	case time.Duration:
		result = int(nv)
		return
	case time.Weekday:
		result = int(nv)
		return
	default:
		nT, errT := strconv.ParseInt(fmt.Sprintf("%v", v), 10, 0)
		if errT != nil {
			result = defaultT
			return
		}

		result = int(nT)
		return
	}
}

var ToInt = TKX.ToInt

func (pA *TK) ToRune(v interface{}, defaultA ...rune) (result rune) {
	var defaultT rune

	if defaultA == nil || len(defaultA) < 1 {
		defaultT = 0
	} else {
		defaultT = defaultA[0]
	}

	defer func() {
		r := recover()

		if r != nil {
			result = defaultT
			return
		}
	}()

	switch v.(type) {
	case bool:
		if v.(bool) {
			result = 1
		} else {
			result = 0
		}

		return
	case int:
		result = rune(v.(int))
		return
	case int8:
		result = rune(v.(int8))
		return
	case int16:
		result = rune(v.(int16))
		return
	case int32:
		result = v.(int32)
		return
	case int64:
		result = rune(v.(int64))
		return
	case uint:
		result = rune(v.(int))
		return
	case uint8:
		result = rune(v.(uint8))
		return
	case uint16:
		result = rune(v.(uint16))
		return
	case uint32:
		result = rune(v.(uint32))
		return
	case uint64:
		result = rune(v.(uint64))
		return
	case float64:
		result = rune(v.(float64))
		return
	case float32:
		result = rune(v.(float32))
		return
	case string:
		nT, errT := strconv.ParseInt(strings.TrimSpace(v.(string)), 10, 0)
		if errT != nil {
			result = defaultT
			return
		}

		result = rune(nT)
		return
	default:
		nT, errT := strconv.ParseInt(fmt.Sprintf("%v", v), 10, 0)
		if errT != nil {
			result = defaultT
			return
		}

		result = rune(nT)
		return
	}
}

var ToRune = TKX.ToRune

func (pA *TK) StrToIntQuick(strA string) (int, error) {
	var n uint64
	var v byte

	lenT := len(strA)

	firstT := true

	minusT := false

	for i := 0; i < lenT; i++ {
		d := strA[i]

		if firstT && d == '-' {
			minusT = true
			firstT = false
			continue
		}

		if '0' <= d && d <= '9' {
			v = d - '0'
			firstT = false
		} else {
			return 0, Errf("invalid int")
		}

		n *= uint64(10)
		n += uint64(v)
	}

	if minusT {
		return -int(n), nil
	}

	return int(n), nil
}

var StrToIntQuick = TKX.StrToIntQuick

// StrToIntWithDefaultValue 字符串转整数，如果有问题则返回默认数值
func (pA *TK) StrToIntWithDefaultValue(strA string, defaultA ...int) int {
	defaultT := -1

	if (defaultA != nil) && (len(defaultA) > 0) {
		defaultT = defaultA[0]
	}

	nT, errT := strconv.ParseInt(strA, 10, 0)
	if errT != nil {
		return defaultT
	}

	return int(nT)
}

var StrToIntWithDefaultValue = TKX.StrToIntWithDefaultValue

// StrToInt 字符串转整数
func (pA *TK) StrToInt(strA string, defaultA ...int) int {
	defaultT := -1

	if (defaultA != nil) && (len(defaultA) > 0) {
		defaultT = defaultA[0]
	}

	nT, errT := strconv.ParseInt(strA, 10, 0)
	if errT != nil {
		return defaultT
	}

	return int(nT)
}

var StrToInt = TKX.StrToInt

// StrToIntE 字符串转整数
func (pA *TK) StrToIntE(strA string) (int, error) {
	nT, errT := strconv.ParseInt(strA, 10, 0)

	return int(nT), errT
}

var StrToIntE = TKX.StrToIntE

func (pA *TK) ToIntI(valueA interface{}, defaultA int) int {
	nT, errT := strconv.ParseInt(fmt.Sprintf("%d", valueA), 10, 0)
	if errT != nil {
		return defaultA
	}

	return int(nT)
}

var ToIntI = TKX.ToIntI

func (pA *TK) StrToInt64(strA string, defaultA ...int64) int64 {
	var defaultT int64 = -1

	if (defaultA != nil) && (len(defaultA) > 0) {
		defaultT = defaultA[0]
	}

	nT, errT := strconv.ParseInt(strA, 10, 0)
	if errT != nil {
		return defaultT
	}

	return nT
}

var StrToInt64 = TKX.StrToInt64

func (pA *TK) StrToInt64WithDefaultValue(strA string, defaultA int64) int64 {
	nT, errT := strconv.ParseInt(strA, 10, 64)
	if errT != nil {
		return defaultA
	}

	return nT
}

var StrToInt64WithDefaultValue = TKX.StrToInt64WithDefaultValue

func (pA *TK) StrToIntPositive(strA string) int {
	nT, errT := strconv.ParseInt(strA, 10, 0)
	if errT != nil {
		return -1
	}

	return int(nT)
}

var StrToIntPositive = TKX.StrToIntPositive

func (pA *TK) StrToHex(strA string) string {
	return hex.EncodeToString([]byte(strA))
}

var StrToHex = TKX.StrToHex

func (pA *TK) ToHex(vA interface{}) string {
	var rs string

	switch vT := vA.(type) {
	case byte:
		rs = fmt.Sprintf("%x", vT)

		if len(rs)%2 != 0 {
			rs = "0" + rs
		}
	case int:
		rs = fmt.Sprintf("%x", vT)

		if len(rs)%2 != 0 {
			rs = "0" + rs
		}
	case []byte:
		rs = fmt.Sprintf("%x", vT)
	case bool:
		if vT {
			rs = "01"
		}
		rs = "00"
	case string:
		rs = fmt.Sprintf("%x", vT)
	default:
		rs1 := DataToBytes(vA)

		if !IsError(rs1) {
			rs = fmt.Sprintf("%x", rs1)
		} else {
			rs = fmt.Sprintf("%x", vA)
		}
	}

	return rs
}

var ToHex = TKX.ToHex

func (pA *TK) StrToFloat64WithDefaultValue(strA string, defaultA float64) float64 {
	nT, errT := strconv.ParseFloat(strA, 64)
	if errT != nil {
		return defaultA
	}

	return nT
}

var StrToFloat64WithDefaultValue = TKX.StrToFloat64WithDefaultValue

func (pA *TK) StrToFloat64(strA string, defaultA ...float64) float64 {
	var defaultT float64 = -1

	if (defaultA != nil) && (len(defaultA) > 0) {
		defaultT = defaultA[0]
	}

	nT, errT := strconv.ParseFloat(strA, 64)

	if errT != nil {
		return defaultT
	}

	return nT
}

var StrToFloat64 = TKX.StrToFloat64

func (pA *TK) StrToFloat64E(strA string) (float64, error) {
	nT, errT := strconv.ParseFloat(strA, 64)

	return nT, errT
}

var StrToFloat64E = TKX.StrToFloat64E

func (pA *TK) Float64ToStr(floatA float64) string {
	tmps := fmt.Sprintf("%f", floatA)
	if Contains(tmps, ".") {
		tmps = strings.TrimRight(tmps, "0")
		tmps = strings.TrimRight(tmps, ".")
	}

	return tmps
}

var Float64ToStr = TKX.Float64ToStr

func (pA *TK) StrToTimeCompact(strA string, defaultA time.Time) time.Time {
	t, err := time.Parse(TimeFormatCompact, strA)
	if err != nil {
		return defaultA
	}

	return t
}

var StrToTimeCompact = TKX.StrToTimeCompact

func (pA *TK) StrToTimeCompactNoError(strA string) time.Time {
	t, _ := time.Parse(TimeFormatCompact, strA)

	return t
}

var StrToTimeCompactNoError = TKX.StrToTimeCompactNoError

func (pA *TK) FormatStringSliceSlice(sliceA [][]string, sepA string, lineSepA string) string {
	var bufT strings.Builder

	for i, v := range sliceA {
		if i != 0 {
			bufT.WriteString(lineSepA)
		}

		for ii, vv := range v {
			if ii != 0 {
				bufT.WriteString(sepA)
			}

			bufT.WriteString(vv)
		}
	}

	return bufT.String()
}

var FormatStringSliceSlice = TKX.FormatStringSliceSlice

// number related

func (pA *TK) AbsInt(c int) int {
	if c < 0 {
		return -c
	} else {
		return c
	}
}

var AbsInt = TKX.AbsInt

func (pA *TK) Abs(c interface{}) interface{} {
	switch nv := c.(type) {
	case int:
		if nv < 0 {
			return -nv
		} else {
			return nv
		}
	case int8:
		if nv < 0 {
			return -nv
		} else {
			return nv
		}
	case int16:
		if nv < 0 {
			return -nv
		} else {
			return nv
		}
	case int32:
		if nv < 0 {
			return -nv
		} else {
			return nv
		}
	case int64:
		if nv < 0 {
			return -nv
		} else {
			return nv
		}
	case uint:
		if nv < 0 {
			return -nv
		} else {
			return nv
		}
	case uint8:
		if nv < 0 {
			return -nv
		} else {
			return nv
		}
	case uint16:
		if nv < 0 {
			return -nv
		} else {
			return nv
		}
	case uint32:
		if nv < 0 {
			return -nv
		} else {
			return nv
		}
	case uint64:
		if nv < 0 {
			return -nv
		} else {
			return nv
		}
	case float64:
		if nv < 0.0 {
			return -nv
		} else {
			return nv
		}
	case float32:
		if nv < 0.0 {
			return -nv
		} else {
			return nv
		}
	}

	return fmt.Errorf("unknown type: %T", c)
}

var Abs = TKX.Abs

// IntToKMGT convert a number to "3.21K", "1.2G", etc, formatA like "%.02f"
// if sizeA < 1024, formatA is ignored
func (pA *TK) IntToKMGT(sizeA interface{}, formatA ...interface{}) string {
	formatT := "%.02f"
	if len(formatA) > 0 {
		formatT = ToStr(formatA[0])
	}

	var sizeT int = ToInt(sizeA, 0)

	if sizeT >= 1099511627776 {
		return fmt.Sprintf(formatT+"T", float64(sizeT)/1099511627776)
	} else if sizeT >= 1073741824 {
		return fmt.Sprintf(formatT+"G", float64(sizeT)/1073741824)
	} else if sizeT >= 1048576 {
		return fmt.Sprintf(formatT+"M", float64(sizeT)/1048576)
	} else if sizeT >= 1024 {
		return fmt.Sprintf(formatT+"K", float64(sizeT)/1024)
	} else {
		return fmt.Sprintf("%dB", sizeT)
	}
}

var IntToKMGT = TKX.IntToKMGT

// IntToWYZ convert a number to "3.21万", "1.2亿", etc, formatA like "%.2f"
// if sizeA < 1024, formatA is ignored
func (pA *TK) IntToWYZ(sizeA int, formatA string) string {
	if formatA == "" {
		formatA = "%.2f"
	}

	if sizeA >= 1000000000000 {
		return fmt.Sprintf(formatA+"兆", float64(sizeA)/1000000000000)
	} else if sizeA >= 100000000 {
		return fmt.Sprintf(formatA+"亿", float64(sizeA)/100000000)
	} else if sizeA >= 10000 {
		return fmt.Sprintf(formatA+"万", float64(sizeA)/10000)
	} else {
		return fmt.Sprintf("%d", sizeA)
	}
}

var IntToWYZ = TKX.IntToWYZ

// 日志相关

func (pA *TK) SetLogFile(fileNameA string) {
	logFileG = fileNameA
}

var SetLogFile = TKX.SetLogFile

func (pA *TK) LogWithTime(formatA string, argsA ...interface{}) {
	if EndsWith(formatA, "\n") {
		AppendStringToFile(fmt.Sprintf(fmt.Sprintf("[%v] ", time.Now())+formatA, argsA...), logFileG)
	} else {
		AppendStringToFile(fmt.Sprintf(fmt.Sprintf("[%v] ", time.Now())+formatA+"\n", argsA...), logFileG)
	}
}

var LogWithTime = TKX.LogWithTime

func (pA *TK) LogWithTimeCompact(formatA string, argsA ...interface{}) {
	if EndsWith(formatA, "\n") {
		AppendStringToFile(fmt.Sprintf(fmt.Sprintf("[%v] ", time.Now().Format(TimeFormatCompact2))+formatA, argsA...), logFileG)
	} else {
		AppendStringToFile(fmt.Sprintf(fmt.Sprintf("[%v] ", time.Now().Format(TimeFormatCompact2))+formatA+"\n", argsA...), logFileG)
	}
}

var LogWithTimeCompact = TKX.LogWithTimeCompact

// 文件操作相关函数 file related

// IfFileExists 判断文件是否存在
func (pA *TK) IfFileExists(fileNameA string) bool {
	_, err := os.Stat(fileNameA)
	return err == nil || os.IsExist(err)
}

var IfFileExists = TKX.IfFileExists

func (pA *TK) GetFileInfo(filePathA string) (map[string]string, error) {
	fi, errT := os.Stat(filePathA)
	if errT != nil && !os.IsExist(errT) {
		return nil, errT
	}

	absPathT, errT := filepath.Abs(filePathA)
	if errT != nil {
		return nil, errT
	}

	mapT := map[string]string{"Path": filePathA, "Abs": absPathT, "Name": filepath.Base(filePathA), "Ext": filepath.Ext(filePathA), "Size": Int64ToStr(fi.Size()), "IsDir": BoolToStr(fi.IsDir()), "Time": FormatTime(fi.ModTime(), TimeFormatCompact), "Mode": fmt.Sprintf("%v", fi.Mode())}

	return mapT, nil
}

var GetFileInfo = TKX.GetFileInfo

// IsFile if is file
func (pA *TK) IsFile(fileNameA string) bool {
	f, errT := os.Open(fileNameA)
	if errT != nil {
		return false
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return false
	}

	if mode := fi.Mode(); mode.IsRegular() {
		return true
	}

	return false
}

var IsFile = TKX.IsFile

// IsDirectory if is directory
func (pA *TK) IsDirectory(dirNameA string) bool {
	f, err := os.Open(dirNameA)
	if err != nil {
		return false
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return false
	}

	if mode := fi.Mode(); mode.IsDir() {
		return true
	}

	return false
}

var IsDirectory = TKX.IsDirectory

func (pA *TK) GetFilePathSeperator() string {
	osT := runtime.GOOS
	if osT == "windows" {
		return "\\"
	} else {
		return "/"
	}
}

var GetFilePathSeperator = TKX.GetFilePathSeperator

func (pA *TK) GetLastComponentOfFilePath(pathA string, sepA ...string) string {
	var sepT string
	if len(sepA) > 0 {
		sepT = sepA[0]
		listT := strings.Split(pathA, sepT)

		return listT[len(listT)-1]
	} else {
		sepT = GetFilePathSeperator()
	}

	if EndsWith(pathA, sepT) {
		return ""
	} else {
		return filepath.Base(pathA)
	}
}

var GetLastComponentOfFilePath = TKX.GetLastComponentOfFilePath

func (pA *TK) GetLastComponentBySeparator(strA string, sepA ...string) string {
	var sepT string
	if len(sepA) > 0 {
		sepT = sepA[0]
	} else {
		sepT = GetFilePathSeperator()
	}

	listT := strings.Split(strA, sepT)

	return listT[len(listT)-1]
}

var GetLastComponentBySeparator = TKX.GetLastComponentBySeparator

func (pA *TK) GetDirOfFilePath(pathA string) string {
	return filepath.Dir(pathA)
}

var GetDirOfFilePath = TKX.GetDirOfFilePath

func (pA *TK) RemoveFileExt(filePathA string) string {
	extT := filepath.Ext(filePathA)
	return filePathA[:len(filePathA)-len(extT)]
}

var RemoveFileExt = TKX.RemoveFileExt

func (pA *TK) GetFileExt(filePathA string) string {
	return filepath.Ext(filePathA)
}

var GetFileExt = TKX.GetFileExt

func (pA *TK) RemoveLastSubString(strA string, subStrA string) string {
	if EndsWith(strA, subStrA) {
		return strA[:len(strA)-len(subStrA)]
	}
	return strA
}

var RemoveLastSubString = TKX.RemoveLastSubString

func (pA *TK) RemoveFirstSubString(strA string, subStrA string) string {
	if StartsWith(strA, subStrA) {
		return strA[len(subStrA):]
	}
	return strA
}

var RemoveFirstSubString = TKX.RemoveFirstSubString

func (pA *TK) AddLastSubString(strA string, subStrA string) string {
	if !EndsWith(strA, subStrA) {
		return strA + subStrA
	}
	return strA
}

var AddLastSubString = TKX.AddLastSubString

func (pA *TK) RemoveFile(filePathT string) error {
	if IsDirectory(filePathT) {
		return Errf("%v is a directory", filePathT)
	}

	errT := os.Remove(filePathT)

	if errT != nil {
		return errT
	}

	if IfFileExists(filePathT) {
		return Errf("failed to remove file: %v", filePathT)
	}

	return nil
}

var RemoveFile = TKX.RemoveFile

func (pA *TK) CreateFile(filePathT string, optsA ...string) error {
	overwriteT := IfSwitchExistsWhole(optsA, "-overwrite")

	if !overwriteT {
		if IfFileExists(filePathT) {
			return Errf("file already exists")
		}
	}

	fileT, errT := os.Create(filePathT)

	if errT != nil {
		return errT
	}

	fileT.Close()

	return nil
}

var CreateFile = TKX.CreateFile

func (pA *TK) OpenFile(filePathT string, optsA ...string) interface{} {
	flagStrT := GetSwitch(optsA, "-flag=")

	flagT := ToInt(flagStrT, 0)

	if IfSwitchExistsWhole(optsA, "-readOnly") {
		flagT = flagT | os.O_RDONLY
	} else {
		flagT = flagT | os.O_RDWR
	}

	if IfSwitchExistsWhole(optsA, "-create") {
		flagT = flagT | os.O_CREATE
	}

	if IfSwitchExistsWhole(optsA, "-append") {
		flagT = flagT | os.O_APPEND
	}

	if IfSwitchExistsWhole(optsA, "-truncate") {
		flagT = flagT | os.O_TRUNC
	}

	permT := OctetToInt(GetSwitch(optsA, "-perm=", "0777"), 0777)

	fileT, errT := os.OpenFile(filePathT, flagT, fs.FileMode(permT))

	if errT != nil {
		return errT
	}

	return fileT
}

var OpenFile = TKX.OpenFile

func (pA *TK) RenameFile(filePathT string, destFilePathT string, optsA ...string) error {
	if !IfFileExists(filePathT) {
		return Errf("file not exists")
	}

	errT := os.Rename(filePathT, destFilePathT)

	if errT != nil {
		return errT
	}

	return nil
}

var RenameFile = TKX.RenameFile

//
// func (pA *TK) CreateTempFile(dirT string, patternT string, optsA ...string) error {
// 	f, err := os.CreateTemp(dirT, patternT)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer os.Remove(f.Name()) // clean up

// 	if _, err := f.Write([]byte("content")); err != nil {
// 		log.Fatal(err)
// 	}
// 	if err := f.Close(); err != nil {
// 		log.Fatal(err)
// 	}
// }

// var CreateTempFile = TKX.CreateTempFile

func (pA *TK) GetFileSize(filePathA string) (int64, error) {
	fi, err := os.Stat(filePathA)
	if err != nil {
		return -1, err
	}

	sizeT := fi.Size()

	return sizeT, nil
}

var GetFileSize = TKX.GetFileSize

func (pA *TK) GetFileSizeCompact(filePathA string) int64 {
	fi, err := os.Stat(filePathA)
	if err != nil {
		return -1
	}

	sizeT := fi.Size()

	return sizeT
}

var GetFileSizeCompact = TKX.GetFileSizeCompact

func (pA *TK) GenerateFileListInDir(dirA string, patternA string, verboseA bool) []string {
	strListT := make([]string, 0, 100)

	pathT, errT := filepath.Abs(dirA)

	if errT != nil {
		return nil
	}

	errT = filepath.Walk(pathT, func(path string, f os.FileInfo, err error) error {
		if verboseA {
			Pln(path)
		}

		if f == nil {
			return err
		}

		// Pl("pathT: %v -> path: %v", pathT, path)

		// if f.IsDir() { // && path != "." && path != pathT {
		if f.IsDir() {
			if path != "." && path != pathT {
				return filepath.SkipDir
			} else {
				return nil
			}
		}

		matchedT, errTI := filepath.Match(patternA, filepath.Base(path))
		if errTI == nil {
			if matchedT {
				strListT = append(strListT, path)
			}
		}

		return nil
	})

	if errT != nil {
		if verboseA {
			Pl("Search directory failed: %v", errT.Error())
		}

		return nil
	}

	return strListT
}

var GenerateFileListInDir = TKX.GenerateFileListInDir

func (pA *TK) GenerateFileListRecursively(dirA string, patternA string, verboseA bool) []string {
	strListT := make([]string, 0, 100)

	errT := filepath.Walk(dirA, func(path string, f os.FileInfo, err error) error {
		if verboseA {
			Pln(path)
		}

		if f == nil {
			return err
		}

		if f.IsDir() {
			return nil
		}

		matchedT, errTI := filepath.Match(patternA, filepath.Base(path))
		if errTI == nil {
			if matchedT {
				strListT = append(strListT, path)
				// Pl("append path: %v", path)
			}
		} else {
			// Pl("matching failed: %v", errTI.Error())
		}

		return nil
	})

	if errT != nil {
		Pl("Search directory failed: %v", errT.Error())
		return nil
	}

	return strListT
}

var GenerateFileListRecursively = TKX.GenerateFileListRecursively

func (pA *TK) GenerateFileListRecursivelyWithExclusive(dirA string, patternA string, exclusivePatternA string, verboseA bool) []string {
	strListT := make([]string, 0, 100)

	errT := filepath.Walk(dirA, func(path string, f os.FileInfo, err error) error {
		if verboseA {
			Pln(path)
		}

		if f == nil {
			return err
		}

		if f.IsDir() {
			return nil
		}

		matchedT, errTI := filepath.Match(patternA, filepath.Base(path))
		if errTI == nil {
			if matchedT {
				if exclusivePatternA != "" {
					matched2T, err2TI := filepath.Match(exclusivePatternA, filepath.Base(path))
					if err2TI == nil {
						if matched2T {
							return nil
						}
					}
				}

				strListT = append(strListT, path)
			}
		} else {
			Pl("matching failed: %v", errTI.Error())
		}

		return nil
	})

	if errT != nil {
		Pl("Search directory failed: %v", errT.Error())
		return nil
	}

	return strListT
}

var GenerateFileListRecursivelyWithExclusive = TKX.GenerateFileListRecursivelyWithExclusive

type FileInfo struct {
	Name  string
	IsDir string
	Size  string
	Mode  string
	Time  string
}

type MSSArraySorter struct {
	KeyM  string
	DescM bool
	AryM  []map[string]string
}

func (s MSSArraySorter) Len() int {
	return len(s.AryM)
}

func (s MSSArraySorter) Swap(i, j int) {
	s.AryM[i], s.AryM[j] = s.AryM[j], s.AryM[i]
}

func (s MSSArraySorter) Less(i, j int) bool {
	if s.DescM {
		return s.AryM[i][s.KeyM] > s.AryM[j][s.KeyM]
	}

	return s.AryM[i][s.KeyM] < s.AryM[j][s.KeyM]
}

func SortMSSArray(aryA []map[string]string, keyA string, descA bool) {
	sortStructT := MSSArraySorter{
		KeyM:  keyA,
		DescM: descA,
		AryM:  aryA,
	}

	sort.Sort(sortStructT)
}

type MSIArraySorter struct {
	KeyM  string
	DescM bool
	AryM  []map[string]interface{}
}

func (s MSIArraySorter) Len() int {
	return len(s.AryM)
}

func (s MSIArraySorter) Swap(i, j int) {
	s.AryM[i], s.AryM[j] = s.AryM[j], s.AryM[i]
}

func (s MSIArraySorter) Less(i, j int) bool {
	rsT := LessI(s.AryM[i][s.KeyM], s.AryM[j][s.KeyM])
	if s.DescM {
		return !rsT
	}

	return rsT
}

func SortMSIArray(aryA []map[string]interface{}, keyA string, descA bool) {
	sortStructT := MSIArraySorter{
		KeyM:  keyA,
		DescM: descA,
		AryM:  aryA,
	}

	sort.Sort(sortStructT)
}

func (pA *TK) GetFileList(dirA string, argsA ...string) []map[string]string {

	verboseT := IfSwitchExistsWhole(argsA, "-verbose")
	if verboseT {
		Pl("dir: %v, params: %v", dirA, argsA)
	}

	recursiveT := IfSwitchExistsWhole(argsA, "-recursive")

	compactT := IfSwitchExistsWhole(argsA, "-compact")

	minSizeT := int64(GetSwitchWithDefaultIntValue(argsA, "-minSize=", -1))
	maxSizeT := int64(GetSwitchWithDefaultIntValue(argsA, "-maxSize=", MAX_INT))

	limitT := GetSwitchWithDefaultIntValue(argsA, "-limit=", 0)

	sortT := GetSwitch(argsA, "-sort=", "")
	sortKeyT := GetSwitch(argsA, "-sortKey=", "Time")

	withDirectoryT := IfSwitchExistsWhole(argsA, "-withDir")
	dirOnlyT := IfSwitchExistsWhole(argsA, "-dirOnly")

	patternT := GetSwitch(argsA, "-pattern=", "*")

	exclusivePatternT := GetSwitch(argsA, "-exclusive=", "")

	pathT, errT := filepath.Abs(dirA)

	if errT != nil {
		return nil
	}

	mapListT := make([]map[string]string, 0, 100)

	countT := 0

	errT = filepath.Walk(dirA, func(path string, f os.FileInfo, err error) error {
		if limitT > 0 && countT >= limitT {
			return nil
		}

		if verboseT {
			Pln(path)
		}

		if f == nil {
			return err
		}

		if f.IsDir() {
			if recursiveT {
				if withDirectoryT || dirOnlyT {
					matchedT, errTI := filepath.Match(patternT, filepath.Base(path))
					if errTI == nil {
						if matchedT {
							if exclusivePatternT != "" {
								matched2T, err2TI := filepath.Match(exclusivePatternT, filepath.Base(path))
								if err2TI == nil {
									if !matched2T {
										absPathT, err1T := filepath.Abs(path)

										if err1T == nil {
											if compactT {
												mapListT = append(mapListT, map[string]string{"Abs": absPathT, "Size": Int64ToStr(f.Size()), "IsDir": "true"})
												countT++
												if limitT > 0 && countT >= limitT {
													return filepath.SkipAll
												}

											} else {
												mapListT = append(mapListT, map[string]string{"Path": path, "Abs": absPathT, "Ext": filepath.Ext(path), "Name": filepath.Base(path), "Size": Int64ToStr(f.Size()), "IsDir": "true", "Time": FormatTime(f.ModTime(), TimeFormatCompact), "Mode": fmt.Sprintf("%v", f.Mode())})
												countT++
												if limitT > 0 && countT >= limitT {
													return filepath.SkipAll
												}

											}
										}
									}
								}
							} else {
								absPathT, err1T := filepath.Abs(path)

								if err1T == nil {
									if compactT {
										mapListT = append(mapListT, map[string]string{"Abs": absPathT, "IsDir": "true"})
										countT++
										if limitT > 0 && countT >= limitT {
											return filepath.SkipAll
										}

									} else {
										mapListT = append(mapListT, map[string]string{"Path": path, "Abs": absPathT, "Ext": filepath.Ext(path), "Name": filepath.Base(path), "Size": Int64ToStr(f.Size()), "IsDir": "true", "Time": FormatTime(f.ModTime(), TimeFormatCompact), "Mode": fmt.Sprintf("%v", f.Mode())})
										countT++
										if limitT > 0 && countT >= limitT {
											return filepath.SkipAll
										}

									}
								}
							}
						}
					}
				}

				return nil
			} else {
				if path != "." && path != pathT {
					if withDirectoryT || dirOnlyT {
						matchedT, errTI := filepath.Match(patternT, filepath.Base(path))
						if errTI == nil {
							if matchedT {
								if exclusivePatternT != "" {
									matched2T, err2TI := filepath.Match(exclusivePatternT, filepath.Base(path))
									if err2TI == nil {
										if !matched2T {
											absPathT, err1T := filepath.Abs(path)

											if err1T == nil {
												if compactT {
													mapListT = append(mapListT, map[string]string{"Abs": absPathT, "Size": Int64ToStr(f.Size()), "IsDir": "true"})
													countT++
													if limitT > 0 && countT >= limitT {
														return filepath.SkipAll
													}

												} else {
													mapListT = append(mapListT, map[string]string{"Path": path, "Abs": absPathT, "Ext": filepath.Ext(path), "Name": filepath.Base(path), "Size": Int64ToStr(f.Size()), "IsDir": "true", "Time": FormatTime(f.ModTime(), TimeFormatCompact), "Mode": fmt.Sprintf("%v", f.Mode())})
													countT++
													if limitT > 0 && countT >= limitT {
														return filepath.SkipAll
													}

												}
											}
										}
									}
								} else {
									absPathT, err1T := filepath.Abs(path)

									if err1T == nil {
										if compactT {
											mapListT = append(mapListT, map[string]string{"Abs": absPathT, "Size": Int64ToStr(f.Size()), "IsDir": "true"})
											countT++
											if limitT > 0 && countT >= limitT {
												return filepath.SkipAll
											}

										} else {
											mapListT = append(mapListT, map[string]string{"Path": path, "Abs": absPathT, "Ext": filepath.Ext(path), "Name": filepath.Base(path), "Size": Int64ToStr(f.Size()), "IsDir": "true", "Time": FormatTime(f.ModTime(), TimeFormatCompact), "Mode": fmt.Sprintf("%v", f.Mode())})
											countT++
											if limitT > 0 && countT >= limitT {
												return filepath.SkipAll
											}

										}
									}
								}
							}
						}
					}

					return filepath.SkipDir
				} else {
					return nil
				}
			}
		}

		if dirOnlyT {
			return nil
		}

		matchedT, errTI := filepath.Match(patternT, filepath.Base(path))
		if errTI == nil {
			if matchedT && f.Size() >= minSizeT && f.Size() <= maxSizeT {
				if exclusivePatternT != "" {
					matched2T, err2TI := filepath.Match(exclusivePatternT, filepath.Base(path))
					if err2TI == nil {
						if matched2T {
							return nil
						}
					}
				}

				absPathT, err1T := filepath.Abs(path)

				if err1T != nil {
					return nil
				}

				if compactT {
					mapListT = append(mapListT, map[string]string{"Abs": absPathT, "Size": Int64ToStr(f.Size()), "IsDir": BoolToStr(f.IsDir())})
					countT++
					if limitT > 0 && countT >= limitT {
						return filepath.SkipAll
					}

				} else {
					mapListT = append(mapListT, map[string]string{"Path": path, "Abs": absPathT, "Ext": filepath.Ext(path), "Name": filepath.Base(path), "Size": Int64ToStr(f.Size()), "IsDir": BoolToStr(f.IsDir()), "Time": FormatTime(f.ModTime(), TimeFormatCompact), "Mode": fmt.Sprintf("%v", f.Mode())})
					countT++
					if limitT > 0 && countT >= limitT {
						return filepath.SkipAll
					}

				}

			}
		} else {
			if verboseT {
				Pl("matching failed: %v", errTI.Error())
			}
		}

		return nil
	})

	if errT != nil {
		if verboseT {
			Pl("Search directory failed: %v", errT.Error())
		}

		return nil
	}

	if sortT == "" {
		return mapListT
	}

	sortStructT := MSSArraySorter{
		KeyM:  sortKeyT,
		DescM: sortT == "desc",
		AryM:  mapListT,
	}

	sort.Sort(sortStructT)

	return mapListT
}

var GetFileList = TKX.GetFileList

func (pA *TK) Ls(dirA string) []string {
	return GenerateFileListInDir(dirA, "*", false)
}

var Ls = TKX.Ls

func (pA *TK) Lsr(dirA string) []string {
	return GenerateFileListRecursivelyWithExclusive(dirA, "*", "", false)
}

var Lsr = TKX.Lsr

func (pA *TK) GetAvailableFileName(fileNameA string) string {
	fileNameT := fileNameA

	for i := 1; true; i++ {
		if !IfFileExists(fileNameT) {
			break
		}

		fileNameT = fmt.Sprintf("%s_%d%s", RemoveFileExt(fileNameA), i, filepath.Ext(fileNameA))
	}

	return fileNameT
}

var GetAvailableFileName = TKX.GetAvailableFileName

// LoadStringFromFile 从文件中读取整个内容到字符串中
func (pA *TK) LoadStringFromFile(fileNameA string) string {
	if !IfFileExists(fileNameA) {
		return GenerateErrorString("file not exists")
	}

	fileT, err := os.Open(fileNameA)
	if err != nil {
		return GenerateErrorString(err.Error())
	}

	defer fileT.Close()

	fileContentT, err := io.ReadAll(fileT)
	if err != nil {
		return GenerateErrorString(err.Error())
	}

	return string(fileContentT)
}

var LoadStringFromFile = TKX.LoadStringFromFile

func (pA *TK) LoadText(fileNameA string) interface{} {
	if !IfFileExists(fileNameA) {
		return fmt.Errorf("file not exists")
	}

	fileT, err := os.Open(fileNameA)
	if err != nil {
		return err
	}

	defer fileT.Close()

	fileContentT, err := io.ReadAll(fileT)
	if err != nil {
		return err
	}

	return string(fileContentT)
}

var LoadText = TKX.LoadText

// LoadStringFromFileWithDefault 从文件中读取整个内容到字符串中，出现问题时返回默认字符串
func (pA *TK) LoadStringFromFileWithDefault(fileNameA string, defaultA string) string {
	if !IfFileExists(fileNameA) {
		return defaultA
	}

	fileT, errT := os.Open(fileNameA)
	if errT != nil {
		return defaultA
	}

	defer fileT.Close()

	fileContentT, errT := io.ReadAll(fileT)
	if errT != nil {
		return defaultA
	}

	return string(fileContentT)
}

var LoadStringFromFileWithDefault = TKX.LoadStringFromFileWithDefault

func (pA *TK) LoadStringFromFileE(fileNameA string) (string, error) {
	if !IfFileExists(fileNameA) {
		return "", fmt.Errorf("file not exists")
	}

	fileT, err := os.Open(fileNameA)
	if err != nil {
		return "", err
	}

	defer fileT.Close()

	fileContentT, err := io.ReadAll(fileT)
	if err != nil {
		return "", err
	}

	return string(fileContentT), nil
}

var LoadStringFromFileE = TKX.LoadStringFromFileE

func (pA *TK) LoadStringFromFileB(fileNameA string) (string, bool) {
	if !IfFileExists(fileNameA) {
		return "file not exists", false
	}

	fileT, err := os.Open(fileNameA)
	if err != nil {
		return err.Error(), false
	}

	defer fileT.Close()

	fileContentT, err := io.ReadAll(fileT)
	if err != nil {
		return err.Error(), false
	}

	return string(fileContentT), true
}

var LoadStringFromFileB = TKX.LoadStringFromFileB

// LoadBytes LoadBytes, no numA or numA < 0 indicates read all
func (pA *TK) LoadBytes(fileNameA string, numA ...int) []byte {
	if !IfFileExists(fileNameA) {
		return nil
	}

	fileT, err := os.Open(fileNameA)
	if err != nil {
		return nil
	}

	defer fileT.Close()

	if numA == nil || len(numA) < 1 || numA[0] <= 0 {
		fileContentT, err := io.ReadAll(fileT)
		if err != nil {
			return nil
		}

		return fileContentT
	}

	bufT := make([]byte, numA[0])
	nn, err := fileT.Read(bufT)
	if (err != nil) || (nn != len(bufT)) {
		return nil
	}

	return bufT
}

var LoadBytes = TKX.LoadBytes

// LoadBytesFromFileE LoadBytes, no numA or numA[0] < 0 indicates read all
func (pA *TK) LoadBytesFromFileE(fileNameA string, numA ...int) ([]byte, error) {
	if !IfFileExists(fileNameA) {
		return nil, Errf("file not exists")
	}

	fileT, errT := os.Open(fileNameA)
	if errT != nil {
		return nil, errT
	}

	defer fileT.Close()

	if numA == nil || len(numA) < 1 || numA[0] <= 0 {
		fileContentT, errT := io.ReadAll(fileT)
		if errT != nil {
			return nil, errT
		}

		return fileContentT, nil
	}

	bufT := make([]byte, numA[0])

	nnT, errT := fileT.Read(bufT)
	if errT != nil {
		return nil, errT
	}

	if nnT != len(bufT) {
		return nil, Errf("read bytes not identical")
	}

	return bufT, nil
}

var LoadBytesFromFileE = TKX.LoadBytesFromFileE

// LoadBytesFromFile LoadBytes, no numA or numA[0] < 0 indicates read all
func (pA *TK) LoadBytesFromFile(fileNameA string, numA ...int) interface{} {
	if !IfFileExists(fileNameA) {
		return Errf("file not exists")
	}

	fileT, errT := os.Open(fileNameA)
	if errT != nil {
		return errT
	}

	defer fileT.Close()

	if numA == nil || len(numA) < 1 || numA[0] <= 0 {
		fileContentT, errT := io.ReadAll(fileT)
		if errT != nil {
			return errT
		}

		return fileContentT
	}

	bufT := make([]byte, numA[0])

	nnT, errT := fileT.Read(bufT)
	if errT != nil {
		return errT
	}

	if nnT != len(bufT) {
		return Errf("read bytes not identical")
	}

	return bufT
}

var LoadBytesFromFile = TKX.LoadBytesFromFile

func (pA *TK) SaveBytesToFile(bytesA []byte, fileA string) string {
	file, err := os.Create(fileA)
	if err != nil {
		return GenerateErrorString(err.Error())
	}

	defer file.Close()

	wFile := bufio.NewWriter(file)
	_, err = wFile.Write(bytesA)

	if err != nil {
		return GenerateErrorString(err.Error())
	}

	wFile.Flush()

	return ""
}

var SaveBytesToFile = TKX.SaveBytesToFile

func (pA *TK) SaveBytesToFileE(bytesA []byte, fileA string) error {
	file, err := os.Create(fileA)
	if err != nil {
		return err
	}

	defer file.Close()

	wFile := bufio.NewWriter(file)
	_, err = wFile.Write(bytesA)

	if err != nil {
		return err
	}

	wFile.Flush()

	return nil
}

var SaveBytesToFileE = TKX.SaveBytesToFileE

// SaveStringToFile 保存字符串到文件
func (pA *TK) SaveStringToFile(strA string, fileA string) string {
	file, err := os.Create(fileA)
	if err != nil {
		return GenerateErrorString(err.Error())
	}

	defer file.Close()
	wFile := bufio.NewWriter(file)
	_, errT := wFile.WriteString(strA)
	
	if errT != nil {
		return ErrorToString(errT)
	}
	
	wFile.Flush()

	return ""
}

var SaveStringToFile = TKX.SaveStringToFile

func (pA *TK) SaveStringToFileE(strA string, fileA string) error {
	file, err := os.Create(fileA)
	if err != nil {
		return err
	}

	defer file.Close()
	wFile := bufio.NewWriter(file)
	_, errT := wFile.WriteString(strA)
	if errT != nil {
		return errT
	}
	
	wFile.Flush()

	return nil
}

var SaveStringToFileE = TKX.SaveStringToFileE
var SaveText = TKX.SaveStringToFileE

func (pA *TK) AppendStringToFile(strA string, fileA string) string {
	fileT, err := os.OpenFile(fileA, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		//TXPlErr(err)
		return GenerateErrorString(err.Error())
	}

	defer fileT.Close()

	writerT := bufio.NewWriter(fileT)

	_, errT := writerT.WriteString(strA)

	if errT != nil {
		return ErrorToString(errT)
	}
	
	writerT.Flush()

	return ""
}

var AppendStringToFile = TKX.AppendStringToFile

func (pA *TK) LoadStringList(fileNameA string) ([]string, string) {
	if !IfFileExists(fileNameA) {
		return nil, "file not exists"
	}

	fileT, err := os.Open(fileNameA)
	if err != nil {
		return nil, err.Error()
	}

	defer fileT.Close()

	fileContentT, err := io.ReadAll(fileT)
	if err != nil {
		return nil, err.Error()
	}

	stringList := SplitLines(string(fileContentT))

	return stringList, ""
}

var LoadStringList = TKX.LoadStringList

func (pA *TK) LoadStringListFromFile(fileNameA string) ([]string, error) {
	if !IfFileExists(fileNameA) {
		return nil, fmt.Errorf("file not exists")
	}

	fileT, err := os.Open(fileNameA)
	if err != nil {
		return nil, err
	}

	defer fileT.Close()

	fileContentT, err := io.ReadAll(fileT)
	if err != nil {
		return nil, err
	}

	stringList := SplitLines(string(fileContentT))

	return stringList, nil
}

var LoadStringListFromFile = TKX.LoadStringListFromFile

func (pA *TK) LoadStringListBuffered(fileNameA string, trimA bool, skipEmptyA bool) ([]string, error) {
	if !IfFileExists(fileNameA) {
		return nil, Errf("file not exists", fileNameA)
	}

	fileT, errT := os.Open(fileNameA)
	if errT != nil {
		return nil, errT
	}

	defer fileT.Close()

	bufT := make([]string, 0, 1000)

	readerT := bufio.NewReader(fileT)

	for true {
		strT, errT := readerT.ReadString('\n')
		if errT != nil {
			if errT == io.EOF {
				if trimA {
					strT = Trim(strT)
				}

				if skipEmptyA {
					if strT != "" {
						bufT = append(bufT, strT)
					}
				} else {
					bufT = append(bufT, strT)
				}

			} else {
				return bufT, errT
			}

			break
		}

		if trimA {
			strT = Trim(strT)
		}

		if skipEmptyA {
			if strT != "" {
				bufT = append(bufT, strT)
			}
		} else {
			bufT = append(bufT, strT)
		}

	}

	return bufT, nil
}

var LoadStringListBuffered = TKX.LoadStringListBuffered

func (pA *TK) SaveStringList(strListA []string, fileA string) string {
	if strListA == nil {
		return GenerateErrorString("invalid parameter")
	}

	file, err := os.Create(fileA)
	if err != nil {
		return GenerateErrorString(err.Error())
	}

	defer file.Close()

	wFile := bufio.NewWriter(file)
	_, errT := wFile.WriteString(JoinLines(strListA))
	if errT != nil {
		return ErrorToString(errT)
	}
	
	wFile.Flush()

	return ""
}

var SaveStringList = TKX.SaveStringList

func (pA *TK) SaveStringListWin(strListA []string, fileA string) string {
	if strListA == nil {
		return GenerateErrorString("invalid parameter")
	}

	file, err := os.Create(fileA)
	if err != nil {
		return GenerateErrorString(err.Error())
	}

	defer file.Close()

	wFile := bufio.NewWriter(file)
	_, errT := wFile.WriteString(JoinLinesBySeparator(strListA, "\r\n"))
	if errT != nil {
		return ErrorToString(errT)
	}
	
	wFile.Flush()

	return ""
}

var SaveStringListWin = TKX.SaveStringListWin

func (pA *TK) SaveStringListBufferedByRange(strListA []string, fileA string, sepA string, startA int, endA int) string {
	if strListA == nil {
		return GenerateErrorString("invalid parameter")
	}

	if strListA == nil {
		return GenerateErrorString("empty list")
	}

	lenT := len(strListA)

	if startA < 0 || endA >= lenT {
		return GenerateErrorString("invalid range")
	}

	var errT error

	file, errT := os.Create(fileA)
	if errT != nil {
		return GenerateErrorString(errT.Error())
	}

	defer file.Close()

	wFile := bufio.NewWriter(file)

	for i := startA; i <= endA; i++ {
		_, errT = wFile.WriteString(strListA[i])
		if errT != nil {
			return GenerateErrorString(errT.Error())
		}

		if i != endA {
			_, errT = wFile.WriteString(sepA)
			if errT != nil {
				return GenerateErrorString(errT.Error())
			}
		}
	}

	wFile.Flush()

	return ""
}

var SaveStringListBufferedByRange = TKX.SaveStringListBufferedByRange

func (pA *TK) SaveStringListBuffered(strListA []string, fileA string, sepA string) string {
	if strListA == nil {
		return GenerateErrorString("invalid parameter")
	}

	if strListA == nil {
		return GenerateErrorString("empty list")
	}

	lenT := len(strListA)

	var errT error

	file, errT := os.Create(fileA)
	if errT != nil {
		return GenerateErrorString(errT.Error())
	}

	defer file.Close()

	wFile := bufio.NewWriter(file)

	for i := 0; i < lenT; i++ {
		_, errT = wFile.WriteString(strListA[i])
		if errT != nil {
			return GenerateErrorString(errT.Error())
		}

		if i != (lenT - 1) {
			_, errT = wFile.WriteString(sepA)
			if errT != nil {
				return GenerateErrorString(errT.Error())
			}
		}
	}

	wFile.Flush()

	return ""
}

var SaveStringListBuffered = TKX.SaveStringListBuffered

func (pA *TK) LoadStringListRemoveEmpty(fileNameA string) []string {
	if !IfFileExists(fileNameA) {
		return nil
	}

	fileT, err := os.Open(fileNameA)
	if err != nil {
		return nil
	}

	defer fileT.Close()

	fileContentT, err := io.ReadAll(fileT)
	if err != nil {
		return nil
	}

	stringList := SplitLinesRemoveEmpty(string(fileContentT))

	return stringList
}

var LoadStringListRemoveEmpty = TKX.LoadStringListRemoveEmpty

func (pA *TK) LoadStringListAsMap(fileNameA string) map[string]int {
	rs, errStr := LoadStringList(fileNameA)

	if errStr != "" || rs == nil {
		return nil
	}

	rs1 := make(map[string]int, 0)
	for _, v := range rs {
		rs1[v] = 1
	}

	return rs1
}

var LoadStringListAsMap = TKX.LoadStringListAsMap

func (pA *TK) LoadStringListAsMapRemoveEmpty(fileNameA string) map[string]int {
	rs, errStr := LoadStringList(fileNameA)

	if errStr != "" || rs == nil {
		return nil
	}

	rs1 := make(map[string]int, 0)
	for _, v := range rs {
		if Trim(v) == "" {
			continue
		}

		rs1[v] = 1
	}

	return rs1
}

var LoadStringListAsMapRemoveEmpty = TKX.LoadStringListAsMapRemoveEmpty

func (pA *TK) LoadJSONMapStringFloat64ArrayFromFile(fileNameA string) map[string][]float64 {
	if !IfFileExists(fileNameA) {
		return nil
	}

	strT := LoadStringFromFile(fileNameA)
	if IsErrorString(strT) {
		return nil
	}

	return JSONToMapStringFloat64Array(strT)
}

var LoadJSONMapStringFloat64ArrayFromFile = TKX.LoadJSONMapStringFloat64ArrayFromFile

// ReadLineFromBufioReader return result string, error and if reach EOF
func (pA *TK) ReadLineFromBufioReader(readerA *bufio.Reader) (string, bool, error) {
	if readerA == nil {
		return "", false, Errf("nil reader")
	}

	strT, errT := readerA.ReadString('\n')

	if errT != nil {
		if errT == io.EOF {
			return strT, true, nil
		}

		return strT, false, errT
	}

	return strT, false, nil

}

var ReadLineFromBufioReader = TKX.ReadLineFromBufioReader

//func (pA *TK) ReadPackageFromReader(readerA io.Reader, remainStateA int, remainA string, startA string, endA string, bufLenA int) interface{} {
//	if readerA == nil {
//		return fmt.Errorf("nil reader")
//	}
//	
//	if bufLenA <= 0 {
//		bufLenA = 100
//	}
//	
//	bufT := make([]byte, bufLenA)
//	
//	oneBufT := make([]byte, 0, bufLenA)
//	
//	if len(remainA) > 0 {
//		oneBufT = append(oneBufT, ([]byte(remainA))...)
//	}
//	
//	stateT := remainStateA // 0: start, 1: has prefix
//	
//	startBytesT := []byte(startA)
//	
//	lenStartBytesT := len(startBytesT)
//	
//	for true {
//		n, errT := readerA.Read(bufT)
//		
//		if errT != nil {
////			if errT == io.EOF {
////				if stateT == 0 {
////					return io.EOF
////				}
////				
////				if stateT == 1 {
////					return io.EOF
////				}
////			}
//			
//			return errT
//		}
//		
//		if n <= 0 {
//			continue
//		}
//		
//		if stateT == 0 {
//			idxT := bytes.Index(bufT, startBytesT)
//			
//			if idxT < 0 {
//				continue
//			}
//			
//			if n > lenStartBytesT {
//					
//			}
//			
//			stateT = 
//		}
//	}
//
//	strT, errT := readerA.ReadString('\n')
//
//	if errT != nil {
//		if errT == io.EOF {
//			return strT, true, nil
//		}
//
//		return strT, false, errT
//	}
//
//	return strT, false, nil
//
//}
//
//var ReadPackageFromReader = TKX.ReadPackageFromReader

func (pA *TK) RestoreLineEnds(strA string, replacementA string) string {
	rs := strings.Replace(strA, replacementA, "\n", -1)
	return rs
}

var RestoreLineEnds = TKX.RestoreLineEnds

// 二进制数据相关 bytes related

func (pA *TK) BytesStartsWith(dataA []byte, subA interface{}) bool {
	// Plo(dataA, subA)
	if dataA == nil || len(dataA) < 1 {
		return false
	}

	switch nv := subA.(type) {
	case []byte:
		if len(dataA) < len(nv) {
			return false
		}

		return bytes.HasPrefix(dataA, nv)
	case []rune:
		nv1 := []byte(string(nv))
		if len(dataA) < len(nv1) {
			return false
		}

		return bytes.HasPrefix(dataA, nv1)
	case string:
		nv1 := []byte(nv)
		if len(dataA) < len(nv1) {
			return false
		}

		return bytes.HasPrefix(dataA, nv1)
	case []interface{}:
		lenT := len(nv)

		if len(dataA) < lenT {
			return false
		}

		bufT := make([]byte, lenT)

		for i := 0; i < lenT; i++ {
			bufT[i] = ToByte(nv[i])
		}

		return bytes.HasPrefix(dataA, bufT)
	}

	return false
}

var BytesStartsWith = TKX.BytesStartsWith

func (pA *TK) BytesEndsWith(dataA []byte, subA interface{}) bool {
	if dataA == nil || len(dataA) < 1 {
		return false
	}

	switch nv := subA.(type) {
	case []byte:
		if len(dataA) < len(nv) {
			return false
		}

		return bytes.HasSuffix(dataA, nv)
	case []rune:
		nv1 := []byte(string(nv))
		if len(dataA) < len(nv1) {
			return false
		}

		return bytes.HasSuffix(dataA, nv1)
	case string:
		nv1 := []byte(nv)
		if len(dataA) < len(nv1) {
			return false
		}

		return bytes.HasSuffix(dataA, nv1)
	case []interface{}:
		lenT := len(nv)

		if len(dataA) < lenT {
			return false
		}

		bufT := make([]byte, lenT)

		for i := 0; i < lenT; i++ {
			bufT[i] = ToByte(nv[i])
		}

		return bytes.HasSuffix(dataA, bufT)
	}

	return false
}

var BytesEndsWith = TKX.BytesEndsWith

func (pA *TK) BytesContains(dataA []byte, subA interface{}) bool {
	if dataA == nil || len(dataA) < 1 {
		return false
	}

	switch nv := subA.(type) {
	case []byte:
		if len(dataA) < len(nv) {
			return false
		}

		return bytes.Contains(dataA, nv)
	case []rune:
		nv1 := []byte(string(nv))
		if len(dataA) < len(nv1) {
			return false
		}

		return bytes.Contains(dataA, nv1)
	case string:
		nv1 := []byte(nv)
		if len(dataA) < len(nv1) {
			return false
		}

		return bytes.Contains(dataA, nv1)
	case []interface{}:
		lenT := len(nv)

		if len(dataA) < lenT {
			return false
		}

		bufT := make([]byte, lenT)

		for i := 0; i < lenT; i++ {
			bufT[i] = ToByte(nv[i])
		}

		return bytes.Contains(dataA, bufT)
	}

	return false
}

var BytesContains = TKX.BytesContains

func (pA *TK) BytesIndex(dataA []byte, subA interface{}) int {
	if dataA == nil || len(dataA) < 1 {
		return -1
	}

	switch nv := subA.(type) {
	case []byte:
		if len(dataA) < len(nv) {
			return -1
		}

		return bytes.Index(dataA, nv)
	case []rune:
		nv1 := []byte(string(nv))
		if len(dataA) < len(nv1) {
			return -1
		}

		return bytes.Index(dataA, nv1)
	case string:
		nv1 := []byte(nv)
		if len(dataA) < len(nv1) {
			return -1
		}

		return bytes.Index(dataA, nv1)
	case []interface{}:
		lenT := len(nv)

		if len(dataA) < lenT {
			return -1
		}

		bufT := make([]byte, lenT)

		for i := 0; i < lenT; i++ {
			bufT[i] = ToByte(nv[i])
		}

		return bytes.Index(dataA, bufT)
	}

	return -1
}

var BytesIndex = TKX.BytesIndex

// 双行列表相关 dual list related

func (pA *TK) LoadDualLineList(fileNameA string) ([][]string, string) {
	rs, err := LoadStringList(fileNameA)

	if err != "" {
		return nil, err
	}

	lenT := len(rs) / 2

	bufT := make([][]string, lenT)

	var bufP []string

	for i := 0; i < lenT; i++ {
		bufP = make([]string, 2)

		bufP[0] = rs[i*2]
		bufP[1] = rs[i*2+1]

		bufT[i] = bufP
	}

	return bufT, ""
}

var LoadDualLineList = TKX.LoadDualLineList

func (pA *TK) SaveDualLineList(listA [][]string, fileNameA string) string {
	if listA == nil {
		return GenerateErrorString("nil list")
	}

	lenT := len(listA)

	bufT := make([]string, lenT*2)

	for i := 0; i < lenT; i++ {
		bufT[i*2] = listA[i][0]
		bufT[i*2+1] = listA[i][1]
	}

	return SaveStringList(bufT, fileNameA)
}

var SaveDualLineList = TKX.SaveDualLineList

func (pA *TK) RemoveDuplicateInDualLineList(listA [][]string) [][]string {
	if listA == nil {
		return nil
	}

	listT := listA

	i := 0

	for true {
		lenT := len(listT)

		if lenT <= 1 {
			break
		}

		if i >= lenT {
			break
		}

		found := false

		for j := i + 1; j < lenT; j++ {
			if listT[j][0] == listT[i][0] {
				found = true
				break
			}
		}

		if found {
			listT = append(listT[:i], listT[i+1:]...)
			continue
		}

		i++
	}

	return listT
}

var RemoveDuplicateInDualLineList = TKX.RemoveDuplicateInDualLineList

func (pA *TK) AppendDualLineList(listA [][]string, fileNameA string) string {
	if listA == nil {
		return GenerateErrorString("nil list")
	}

	var listT [][]string
	var err string

	if IfFileExists(fileNameA) {
		listT, err = LoadDualLineList(fileNameA)

		if err != "" {
			return GenerateErrorStringF("failed to open file: %v", err)
		}

		listT = RemoveDuplicateInDualLineList(append(listT, listA...))

		// lenTiA := len(listA)
		// lenTi := len(listT)

		// for i := 0; i < lenTiA; i++ {
		// 	found := -1

		// 	for j := 0; j < lenTi; j++ {
		// 		if listA[i][0] == listT[j][0] {
		// 			found = j
		// 			break
		// 		}
		// 	}

		// 	if found < 0 {
		// 		listT = append(listT, listA[i])
		// 	} else {
		// 		listT[found][0] = listA[i][0]
		// 		listT[found][1] = listA[i][1]
		// 	}
		// }
	} else {
		listT = listA
	}

	lenT := len(listT)

	bufT := make([]string, lenT*2)

	for i := 0; i < lenT; i++ {
		bufT[i*2] = listT[i][0]
		bufT[i*2+1] = listT[i][1]
	}

	return SaveStringList(bufT, fileNameA)
}

var AppendDualLineList = TKX.AppendDualLineList

// SimpleMap related
// in a simplemap structure, key/value pairs are in form as KEY=VALUE
// "=" in keys should be replaced as `EQ`
// line-ends in values such as "\n" should be replaced as #CR#
// comments could be used after ####

func (pA *TK) LoadSimpleMapFromFile(fileNameA string) map[string]string {
	if !IfFileExists(fileNameA) {
		return nil
	}

	strListT, _ := LoadStringList(fileNameA)

	if strListT == nil {
		return nil
	}

	mapT := make(map[string]string)
	for i := range strListT {
		lineT := strListT[i]
		lineT = strings.SplitN(lineT, "####", 2)[0]
		lineListT := strings.SplitN(lineT, "=", 2)
		if (lineListT == nil) || (len(lineListT) < 2) {
			continue
		}
		mapT[Replace(lineListT[0], "`EQ`", "=")] = RestoreLineEnds(lineListT[1], "#CR#")
	}

	return mapT
}

var LoadSimpleMapFromFile = TKX.LoadSimpleMapFromFile

func (pA *TK) LoadSimpleMapFromFileE(fileNameA string) (map[string]string, error) {
	if !IfFileExists(fileNameA) {
		return nil, fmt.Errorf("file not exists")
	}

	fc, errT := LoadStringFromFileE(fileNameA)

	if errT != nil {
		return nil, errT
	}

	return LoadSimpleMapFromStringE(fc)
}

var LoadSimpleMapFromFileE = TKX.LoadSimpleMapFromFileE

func (pA *TK) SimpleMapToString(mapA interface{}) string {

	map1, ok := mapA.(map[string]string)

	var kk string

	if ok {
		strListT := make([]string, 0, len(map1)+1)

		for k, v := range map1 {
			kk = Replace(k, "=", "`EQ`")
			strListT = append(strListT, kk+"="+ReplaceLineEnds(v, "#CR#"))
		}

		return JoinLines(strListT)
	}

	map2, ok := mapA.(map[string]interface{})

	if ok {
		strListT := make([]string, 0, len(map2)+1)

		for k, v := range map2 {
			vs := ToStr(v)
			kk = Replace(k, "=", "`EQ`")
			strListT = append(strListT, kk+"="+ReplaceLineEnds(vs, "#CR#"))
		}

		return JoinLines(strListT)
	}

	return ""
}

var SimpleMapToString = TKX.SimpleMapToString

func (pA *TK) LoadSimpleMapFromString(strA string) map[string]string {
	strListT := SplitLines(strA)

	if strListT == nil {
		return nil
	}

	mapT := make(map[string]string)
	for i := range strListT {
		lineT := strListT[i]

		lineT = strings.SplitN(lineT, "####", 2)[0]

		lineListT := strings.SplitN(lineT, "=", 2)
		if (lineListT == nil) || (len(lineListT) < 2) {
			continue
		}
		mapT[Replace(lineListT[0], "`EQ`", "=")] = RestoreLineEnds(lineListT[1], "#CR#")
	}

	return mapT
}

var LoadSimpleMapFromString = TKX.LoadSimpleMapFromString

func (pA *TK) LoadSimpleMapFromStringE(strA string) (map[string]string, error) {
	strListT := SplitLines(strA)

	if strListT == nil {
		return nil, fmt.Errorf("string nil")
	}

	mapT := make(map[string]string)
	for i := range strListT {
		lineT := strListT[i]

		lineT = strings.SplitN(lineT, "####", 2)[0]

		lineListT := strings.SplitN(lineT, "=", 2)
		if (lineListT == nil) || (len(lineListT) < 2) {
			continue
		}
		mapT[Replace(lineListT[0], "`EQ`", "=")] = RestoreLineEnds(lineListT[1], "#CR#")
	}

	return mapT, nil
}

var LoadSimpleMapFromStringE = TKX.LoadSimpleMapFromStringE

func (pA *TK) ReplaceLineEnds(strA string, replacementA string) string {
	rs := strings.Replace(strA, "\r", "", -1)
	rs = strings.Replace(rs, "\n", replacementA, -1)
	return rs
}

var ReplaceLineEnds = TKX.ReplaceLineEnds

func (pA *TK) SaveSimpleMapToFile(mapA map[string]string, fileA string) string {
	fileT, errT := os.Create(fileA)
	if errT != nil {
		return GenerateErrorString(errT.Error())
	}

	defer fileT.Close()

	wFile := bufio.NewWriter(fileT)

	var kk string

	if mapA != nil {
		for k, v := range mapA {
			kk = Replace(k, "=", "`EQ`")
			wFile.WriteString(kk + "=" + ReplaceLineEnds(v, "#CR#") + "\n")
		}
	}

	wFile.Flush()

	return ""
}

var SaveSimpleMapToFile = TKX.SaveSimpleMapToFile

func (pA *TK) AppendSimpleMapFromFile(mapA map[string]string, fileNameA string) string {
	if !IfFileExists(fileNameA) {
		return "file not exists"
	}

	strListT, errStrT := LoadStringList(fileNameA)

	if errStrT != "" {
		return "fail to load file content"
	}

	for i := range strListT {
		lineT := strListT[i]

		lineT = strings.SplitN(lineT, "####", 2)[0]

		lineListT := strings.SplitN(lineT, "=", 2)
		if (lineListT == nil) || (len(lineListT) < 2) {
			continue
		}

		mapA[Replace(lineListT[0], "`EQ`", "=")] = RestoreLineEnds(lineListT[1], "#CR#")
	}

	return ""
}

var AppendSimpleMapFromFile = TKX.AppendSimpleMapFromFile

func (pA *TK) LoadSimpleMapFromDir(dirA string) map[string]string {
	if !IfFileExists(dirA) {
		return nil
	}

	if !IsDirectory(dirA) {
		return nil
	}

	mapT := make(map[string]string)
	files := GenerateFileListRecursively(dirA, "*.txt", false)
	if files == nil {
		return nil
	}

	for _, v := range files {
		AppendSimpleMapFromFile(mapT, v)
	}

	return mapT
}

var LoadSimpleMapFromDir = TKX.LoadSimpleMapFromDir

// GetLinesFromFile at least will return []string{}, avoid nil result
func (pA *TK) GetLinesFromFile(fileNameA string, startA int, endA int, optionsA ...string) ([]string, error) {
	failForRangeT := false

	if !IsNilOrEmpty(optionsA) {
		if IfSwitchExistsWhole(optionsA, "-index") {
			startA++
			endA++
		}

		if IfSwitchExistsWhole(optionsA, "-range") {
			failForRangeT = true
		}
	}

	if failForRangeT {
		if startA < 1 {
			return []string{}, Errf("invalid start index: %v", startA)
		}
	}

	f, errT := os.Open(fileNameA)
	if errT != nil {
		return []string{}, Errf("failed to open file: %v", errT)
	}

	defer f.Close()

	r := bufio.NewReader(f)

	i := 0

	var bufT = make([]string, 0, endA-startA+2)

	for true {
		line, errT := r.ReadString('\n')

		i++

		// pl("process %v", i)

		if i < startA {
			continue
		}

		if i > endA {
			break
		}

		if errT != nil {
			if errT == io.EOF {
				if failForRangeT {
					return bufT, Errf("invalid end index: %v", endA)
				}

				break
			}

			return bufT, Errf("failed to read content: %v, line: %v", errT, line)
		}

		// pl("%v", line)
		bufT = append(bufT, strings.TrimRight(line, "\r\n"))

	}

	return bufT, nil
}

var GetLinesFromFile = TKX.GetLinesFromFile

// 编码解码相关 encode/decode

func (pA *TK) EncodeToBase64(bufA []byte) string {
	return base64.StdEncoding.EncodeToString(bufA)
}

var EncodeToBase64 = TKX.EncodeToBase64

func (pA *TK) EncodeToBase64ByRawURL(bufA []byte) string {
	return base64.RawURLEncoding.EncodeToString(bufA)
}

var EncodeToBase64ByRawURL = TKX.EncodeToBase64ByRawURL

func (pA *TK) ToBase64(vA interface{}) string {
	if vA == nil {
		return ""
	}

	b1, ok := vA.([]byte)
	if ok {
		return base64.StdEncoding.EncodeToString(b1)
	}

	b2, ok := vA.(string)
	if ok {
		return base64.StdEncoding.EncodeToString([]byte(b2))
	}

	b3, ok := vA.(*bytes.Buffer)
	if ok {
		return base64.StdEncoding.EncodeToString(b3.Bytes())
	}

	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v", vA)))
}

var ToBase64 = TKX.ToBase64

func (pA *TK) EncodeHTML(strA string) string {
	return html.EscapeString(strA)
}

var EncodeHTML = TKX.EncodeHTML

func (pA *TK) DecodeHTML(strA string) string {
	return html.UnescapeString(strA)
}

var DecodeHTML = TKX.DecodeHTML

func (pA *TK) DecodeFromBase64(strA string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(strA)
}

var DecodeFromBase64 = TKX.DecodeFromBase64

func (pA *TK) DecodeFromBase64ByRawURL(strA string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(strA)
}

var DecodeFromBase64ByRawURL = TKX.DecodeFromBase64ByRawURL

func (pA *TK) FromBase64(strA string) interface{} {
	bufT, errT := base64.StdEncoding.DecodeString(strA)

	if errT != nil {
		return errT
	}

	return bufT
}

var FromBase64 = TKX.FromBase64

// EncodeToXMLString 转换字符串XML格式编码的字符串，例如：字符串“<as>\"!sdsdsgfde345344对方对方对法国</as>” 会编码为 “&lt;as&gt;&#34;!sdsdsgfde345344对方对方对法国&lt;/as&gt;”
func (pA *TK) EncodeToXMLString(strA string) string {
	var bufT strings.Builder

	errT := xml.EscapeText(&bufT, []byte(strA))

	if errT != nil {
		return strA
	}

	return bufT.String()
}

var EncodeToXMLString = TKX.EncodeToXMLString

// ToJSON use fast method
func (pA *TK) ToJSON(objA interface{}) (string, error) {
	// var json = jsoniter.ConfigCompatibleWithStandardLibrary
	// var json = jsoniter.ConfigFastest
	rs, errT := jsoniter.Marshal(objA)

	// if errT != nil {
	// 	return GenerateErrorString(errT.Error())
	// }

	return string(rs), errT
}

var ToJSON = TKX.ToJSON

type UnaStruct1 func() (io.ReadCloser, error)

func (a UnaStruct1) MarshalJSON() ([]byte, error) {

	return []byte{}, nil
}

type UnaStruct2 func(actionA string, objA interface{}, dataA interface{}, paramsA ...interface{}) interface{}

func (a UnaStruct2) MarshalJSON() ([]byte, error) {

	return []byte{}, nil
}

func (a TXDelegate) MarshalJSON() ([]byte, error) {

	return []byte("\"N/A(TXDelegate)\""), nil
}

func (a QuickVarDelegate) MarshalJSON() ([]byte, error) {

	return []byte("\"N/A(QuickVarDelegate)\""), nil
}

// func (a func() (io.ReadCloser, error)) MarshalJSON() ([]byte, error) {

// 	return []byte{}, nil
// }

func (pA *TK) Sdump(objsA ...interface{}) string {
	return spew.Sdump(objsA...)
}

var Sdump = TKX.Sdump

func (pA *TK) Sdumpf(formatA string, objsA ...interface{}) string {
	return spew.Sprintf(formatA, objsA...)
}

var Sdumpf = TKX.Sdumpf

func (pA *TK) Dump(objsA ...interface{}) {
	spew.Dump(objsA...)
}

var Dump = TKX.Dump

func (pA *TK) Dumpf(formatA string, objsA ...interface{}) {
	spew.Printf(formatA+"\n", objsA...)
}

var Dumpf = TKX.Dumpf

func (pA *TK) ToJSONX(objA interface{}, optsA ...string) string {
	var errT error

	ifDefaultT := IfSwitchExists(optsA, "-default=")

	indentT := false
	if IfSwitchExistsWhole(optsA, "-indent") {
		indentT = true
	}

	var jsonT jsoniter.API

	if optsA == nil || len(optsA) < 1 {
		jsonT = jsoniter.ConfigDefault
	} else if IfSwitchExistsWhole(optsA, "-sort") {
		jsonT = jsoniter.ConfigCompatibleWithStandardLibrary
	} else if IfSwitchExistsWhole(optsA, "-fast") {
		jsonT = jsoniter.ConfigFastest
	} else {
		jsonT = jsoniter.ConfigDefault
	}

	var rs []byte

	if indentT {
		rs, errT = jsonT.MarshalIndent(objA, "", "  ")
	} else {
		rs, errT = jsonT.Marshal(objA)

	}

	if errT != nil {

		rs, errT = json.Marshal(objA)

		if errT != nil {
			if ifDefaultT {
				return GetSwitchWithDefaultValue(optsA, "-default=", GenerateErrorString(errT.Error()))
			}

			return GenerateErrorString(errT.Error())
		}

	}

	return string(rs)
}

var ToJSONX = TKX.ToJSONX

func (pA *TK) ToJSONWithDefault(objA interface{}, defaultA string) string {
	rs, errT := jsoniter.Marshal(objA)

	if errT != nil {
		return defaultA
	}

	return string(rs)
}

var ToJSONWithDefault = TKX.ToJSONWithDefault

func (pA *TK) FormatJson(jsonA string, optsA ...string) string {
	if len(optsA) > 0 {
		optsT := &pretty.Options{
			Width: ToInt(GetSwitch(optsA, "-width=", "80"), 80),
			Prefix: GetSwitch(optsA, "-prefix=", ""),
			Indent: GetSwitch(optsA, "-indent=", "  "),
			SortKeys: IfSwitchExists(optsA, "-sort"),
		}
		
		rs := pretty.PrettyOptions([]byte(jsonA), optsT)

		return string(rs)
	}
	
	rs := pretty.Pretty([]byte(jsonA))

	return string(rs)
}

var FormatJson = TKX.FormatJson

func (pA *TK) CompactJson(jsonA string, optsA ...string) string {
	rs := pretty.Ugly([]byte(jsonA))

	return string(rs)
}

var CompactJson = TKX.CompactJson

// ToJSONIndent use fast method
func (pA *TK) ToJSONIndent(objA interface{}) (string, error) {
	// var json = jsoniter.ConfigCompatibleWithStandardLibrary
	// var json = jsoniter.ConfigFastest
	rs, errT := jsoniter.MarshalIndent(objA, "", "  ")

	// if errT != nil {
	// 	return GenerateErrorString(errT.Error())
	// }

	return string(rs), errT
}

var ToJSONIndent = TKX.ToJSONIndent

func (pA *TK) ToJSONIndentWithDefault(objA interface{}, defaultA string) string {
	rs, errT := jsoniter.MarshalIndent(objA, "", "  ")

	if errT != nil {
		return defaultA
	}

	return string(rs)
}

var ToJSONIndentWithDefault = TKX.ToJSONIndentWithDefault

func (pA *TK) ReadAllString(vA interface{}) string {
	nv1, ok := vA.([]byte)

	if ok {
		return string(nv1)
	}

	nv2, ok := vA.(io.Reader)

	if ok {
		rsT, errT := io.ReadAll(nv2)

		if errT != nil {
			return ErrorToString(errT)
		}

		return string(rsT)
	}

	return ErrStrf("unsupported type: %T(%v)", vA, vA)
}

var ReadAllString = TKX.ReadAllString

func (pA *TK) Close(vA interface{}) error {
	nv1, ok := vA.(io.Closer)

	if ok {
		return nv1.Close()
	}

	return fmt.Errorf("unsupported type: %T(%v)", vA, vA)
}

var Close = TKX.Close

// FromJson fast JSON decode
func (pA *TK) FromJSON(jsonA string) (interface{}, error) {
	var rs interface{}

	errT := jsoniter.Unmarshal([]byte(jsonA), &rs)

	if errT != nil {
		return nil, errT
	}

	return rs, nil
}

var FromJSON = TKX.FromJSON

func (pA *TK) FromJSONX(vA interface{}) interface{} {
	jsonA := ToStr(vA)

	// if ok {
	// 	return Errf("string type required")

	// }

	rsT, errT := FromJSON(jsonA)

	if errT != nil {
		return errT
	}

	return rsT
}

var FromJSONX = TKX.FromJSONX

func (pA *TK) FromJSONWithDefault(jsonA string, defaultA ...interface{}) interface{} {
	var rs interface{}

	errT := jsoniter.Unmarshal([]byte(jsonA), &rs)

	if errT != nil {
		if len(defaultA) > 0 {
			return defaultA[0]
		}

		return nil
	}

	return rs
}

var FromJSONWithDefault = TKX.FromJSONWithDefault

func (pA *TK) MSSFromJSON(jsonA string) (map[string]string, error) {
	var rs map[string]string

	errT := jsoniter.Unmarshal([]byte(jsonA), &rs)

	if errT != nil {
		return nil, errT
	}

	return rs, nil
}

var MSSFromJSON = TKX.MSSFromJSON

func (pA *TK) LoadJSONFromFile(filePathA string, bufA interface{}) error {
	fcT, errT := LoadBytesFromFileE(filePathA, -1)

	if errT != nil {
		return errT
	}

	errT = jsoniter.Unmarshal(fcT, bufA)

	if errT != nil {
		return errT
	}

	return nil

}

var LoadJSONFromFile = TKX.LoadJSONFromFile

func (pA *TK) LoadJSONFromString(strA string, bufA interface{}) error {
	errT := jsoniter.Unmarshal([]byte(strA), bufA)

	if errT != nil {
		return errT
	}

	return nil

}

var LoadJSONFromString = TKX.LoadJSONFromString

func (pA *TK) SaveJSONToFile(objA interface{}, filePathA string) error {
	rs, errT := jsoniter.Marshal(objA)

	if errT != nil {
		return errT
	}

	errT = SaveStringToFileE(string(rs), filePathA)

	return errT
}

var SaveJSONToFile = TKX.SaveJSONToFile

func (pA *TK) SaveJSONIndentToFile(objA interface{}, filePathA string) error {
	rs, errT := jsoniter.MarshalIndent(objA, "", "  ")

	if errT != nil {
		return errT
	}

	errT = SaveStringToFileE(string(rs), filePathA)

	return errT
}

var SaveJSONIndentToFile = TKX.SaveJSONIndentToFile

func (pA *TK) LoadMSSFromJSONFile(filePathA string) (map[string]string, error) {
	fcT, errT := LoadStringFromFileE(filePathA)

	if errT != nil {
		return nil, errT
	}

	var rs map[string]string

	errT = jsoniter.Unmarshal([]byte(fcT), &rs)

	if errT != nil {
		return nil, errT
	}

	return rs, nil
}

var LoadMSSFromJSONFile = TKX.LoadMSSFromJSONFile

func (pA *TK) SaveMSSToJSONFile(mapA map[string]string, filePathA string) error {
	rs1, errT := ToJSONIndent(mapA)

	if errT != nil {
		return errT
	}

	rs := SaveStringToFileE(rs1, filePathA)

	return rs

}

var SaveMSSToJSONFile = TKX.SaveMSSToJSONFile

// GetJSONNode return jsoniter.Any type as interface{}
func (pA *TK) GetJSONNode(jsonA string, pathA ...interface{}) interface{} {
	aryT := make([]interface{}, 0, len(pathA))

	var typeT reflect.Type
	for _, v := range pathA {
		typeT = reflect.TypeOf(v)

		if typeT.Kind() == reflect.Int64 {
			aryT = append(aryT, int(v.(int64)))
		} else if typeT.Kind() == reflect.String && v.(string) == "*" {
			aryT = append(aryT, int32('*'))
		} else {
			aryT = append(aryT, v)
		}
	}

	rs := jsoniter.Get([]byte(jsonA), aryT...)

	return rs.GetInterface()
}

var GetJSONNode = TKX.GetJSONNode

// GetJSONSubNode return jsoniter.Any type as interface{}
func (pA *TK) GetJSONSubNode(jsonNodeA jsoniter.Any, pathA ...interface{}) interface{} {
	aryT := make([]interface{}, 0, len(pathA))

	var typeT reflect.Type
	for _, v := range pathA {
		typeT = reflect.TypeOf(v)

		if typeT.Kind() == reflect.Int64 {
			aryT = append(aryT, int(v.(int64)))
		} else if typeT.Kind() == reflect.String && v.(string) == "*" {
			aryT = append(aryT, int32('*'))
		} else {
			aryT = append(aryT, v)
		}
	}

	rs := jsonNodeA.Get(aryT...)

	return rs.GetInterface()
}

var GetJSONSubNode = TKX.GetJSONSubNode

// GetJSONNodeAny return jsoniter.Any type
// func Get(data []byte, path ...interface{}) Any takes interface{} as path. If string, it will lookup json map. If int, it will lookup json array. If '*', it will map to each element of array or each key of map.
func (pA *TK) GetJSONNodeAny(jsonA string, pathA ...interface{}) jsoniter.Any {
	aryT := make([]interface{}, 0, len(pathA))

	var typeT reflect.Type
	for _, v := range pathA {
		typeT = reflect.TypeOf(v)

		if typeT.Kind() == reflect.Int64 {
			aryT = append(aryT, int(v.(int64)))
		} else if typeT.Kind() == reflect.String && v.(string) == "*" {
			aryT = append(aryT, int32('*'))
		} else {
			aryT = append(aryT, v)
		}
	}

	rs := jsoniter.Get([]byte(jsonA), aryT...)

	return rs
}

var GetJSONNodeAny = TKX.GetJSONNodeAny

// pathA refer to github.com/tidwall/gjson
func (pA *TK) GetJSONNodeString(jsonA string, pathA string) string {
//	return jsoniter.Get([]byte(jsonA), pathA...).ToString()
	return gjson.Get(jsonA, pathA).String()
}

var GetJSONNodeString = TKX.GetJSONNodeString

// pathA refer to github.com/tidwall/gjson
func (pA *TK) GetJSONNodeStrings(jsonA string, pathA string) interface{} {
	rs := gjson.Get(jsonA, pathA).Array()
	
	lenT := len(rs)
	
	aryT := make([]string, 0, lenT) 
	
	for i := 0; i < lenT; i ++ {
		aryT = append(aryT, rs[i].String())
	}
	
	return aryT
}

var GetJSONNodeStrings = TKX.GetJSONNodeStrings

func (pA *TK) GetJSONSubNodeAny(jsonNodeA jsoniter.Any, pathA ...interface{}) jsoniter.Any {
	aryT := make([]interface{}, 0, len(pathA))

	var typeT reflect.Type
	for _, v := range pathA {
		typeT = reflect.TypeOf(v)

		if typeT.Kind() == reflect.Int64 {
			aryT = append(aryT, int(v.(int64)))
		} else if typeT.Kind() == reflect.String && v.(string) == "*" {
			aryT = append(aryT, int32('*'))
		} else {
			aryT = append(aryT, v)
		}
	}

	rs := jsonNodeA.Get(aryT...)

	return rs
}

var GetJSONSubNodeAny = TKX.GetJSONSubNodeAny

// ObjectToJSON 任意对象转换为JSON字符串
func (pA *TK) ObjectToJSON(objA interface{}) string {
	bufferT, errT := json.Marshal(objA)
	if errT != nil {
		return GenerateErrorStringF("failed: %s", errT.Error())
	}

	return string(bufferT)
}

var ObjectToJSON = TKX.ObjectToJSON

func (pA *TK) ObjectToJSONIndent(objA interface{}) string {
	bufferT, errT := json.MarshalIndent(objA, "", "")
	if errT != nil {
		return GenerateErrorStringF("failed: %s", errT.Error())
	}

	return string(bufferT)
}

var ObjectToJSONIndent = TKX.ObjectToJSONIndent

func (pA *TK) JSONToMapStringFloat64Array(objStrA string) map[string][]float64 {
	var rMapT map[string][]float64
	errT := json.Unmarshal([]byte(objStrA), &rMapT)
	if errT != nil {
		return nil
	}

	return rMapT
}

var JSONToMapStringFloat64Array = TKX.JSONToMapStringFloat64Array

func (pA *TK) JSONToMapStringString(objStrA string) map[string]string {
	var rMapT map[string]string
	errT := json.Unmarshal([]byte(objStrA), &rMapT)
	if errT != nil {
		return nil
	}

	return rMapT
}

var JSONToMapStringString = TKX.JSONToMapStringString

func (pA *TK) JSONToMapStringAny(objStrA string) map[string]interface{} {
	var rMapT map[string]interface{}
	errT := json.Unmarshal([]byte(objStrA), &rMapT)
	if errT != nil {
		return nil
	}

	return rMapT
}

var JSONToMapStringAny = TKX.JSONToMapStringAny

func (pA *TK) JSONToMapStringStringArray(objStrA string) []map[string]string {
	var rMapT []map[string]string
	errT := json.Unmarshal([]byte(objStrA), &rMapT)
	if errT != nil {
		return nil
	}

	return rMapT
}

var JSONToMapStringStringArray = TKX.JSONToMapStringStringArray

func (pA *TK) JSONToObject(objStrA string) interface{} {
	var rs interface{}
	errT := json.Unmarshal([]byte(objStrA), &rs)
	if errT != nil {
		return nil
	}

	return rs
}

var JSONToObject = TKX.JSONToObject

func (pA *TK) JSONToObjectE(objStrA string) (interface{}, error) {
	var rs interface{}

	errT := json.Unmarshal([]byte(objStrA), &rs)
	if errT != nil {
		return nil, errT
	}

	return rs, nil
}

var JSONToObjectE = TKX.JSONToObjectE

func (pA *TK) SafelyGetStringForKeyWithDefault(mapA map[string]string, keyA string, defaultA ...string) string {
	defaultT := ""
	if (defaultA != nil) && (len(defaultA) > 0) {
		defaultT = defaultA[0]
	}

	if mapA == nil {
		return defaultT
	}

	v, ok := mapA[keyA]
	if !ok {
		return defaultT
	}

	return v
}

var SafelyGetStringForKeyWithDefault = TKX.SafelyGetStringForKeyWithDefault

func (pA *TK) GetMSIStringWithDefault(mapA map[string]interface{}, keyA string, defaultA string) (result string) {
	defer func() {
		r := recover()

		if r != nil {
			result = defaultA
			return
		}
	}()

	if mapA == nil {
		result = defaultA
		return
	}

	v, ok := mapA[keyA]
	if !ok {
		result = defaultA
		return
	}

	result = v.(string)
	return
}

var GetMSIStringWithDefault = TKX.GetMSIStringWithDefault

func (pA *TK) SafelyGetFloat64ForKeyWithDefault(mapA map[string]string, keyA string, defaultA float64) float64 {
	if mapA == nil {
		return defaultA
	}

	v, ok := mapA[keyA]
	if !ok {
		return defaultA
	}

	return StrToFloat64WithDefaultValue(v, defaultA)
}

var SafelyGetFloat64ForKeyWithDefault = TKX.SafelyGetFloat64ForKeyWithDefault

func (pA *TK) SafelyGetIntForKeyWithDefault(mapA map[string]string, keyA string, defaultA int) int {
	if mapA == nil {
		return defaultA
	}

	v, ok := mapA[keyA]
	if !ok {
		return defaultA
	}

	return StrToIntWithDefaultValue(v, defaultA)
}

var SafelyGetIntForKeyWithDefault = TKX.SafelyGetIntForKeyWithDefault

func (pA *TK) JSONToStringArray(objStrA string) []string {
	var rArrayT []string
	errT := json.Unmarshal([]byte(objStrA), &rArrayT)
	if errT != nil {
		return nil
	}

	return rArrayT
}

var JSONToStringArray = TKX.JSONToStringArray

func (pA *TK) EncodeStringSimple(strA string) string {
	lenT := len(strA)

	hexCount := 0
	for i := 0; i < lenT; i++ {
		v := strA[i]
		if !(((v >= '0') && (v <= '9')) || ((v >= 'a') && (v <= 'z'))) {
			// Prl("v=%v, %v, %v", v, 'a', (((v >= '0') && (v <= '9')) || ((v >= 'a') && (v <= 'z'))))
			hexCount++
		}
	}

	if hexCount == 0 {
		return strA
	}

	t := make([]byte, lenT+2*hexCount)
	j := 0

	for i := 0; i < lenT; i++ {
		switch v := strA[i]; {
		case !(((v >= '0') && (v <= '9')) || ((v >= 'a') && (v <= 'z'))):
			t[j] = '%'
			t[j+1] = "0123456789ABCDEF"[v>>4]
			t[j+2] = "0123456789ABCDEF"[v&15]
			j += 3
		default:
			t[j] = strA[i]
			j++
		}
	}

	return string(t)
}

var EncodeStringSimple = TKX.EncodeStringSimple

func (pA *TK) EncodeStringUnderline(strA string) string {
	lenT := len(strA)

	var sbuf strings.Builder

	tableStrT := "0123456789ABCDEF"

	for i := 0; i < lenT; i++ {
		v := strA[i]

		if !(((v >= '0') && (v <= '9')) || ((v >= 'a') && (v <= 'z')) || ((v >= 'A') && (v <= 'Z'))) {
			sbuf.WriteByte('_')
			sbuf.WriteByte(tableStrT[v>>4])
			sbuf.WriteByte(tableStrT[v&15])
		} else {
			sbuf.WriteByte(strA[i])
		}
	}

	return sbuf.String()
}

var EncodeStringUnderline = TKX.EncodeStringUnderline

func (pA *TK) EncodeStringCustom(strA string, paddingA byte) string {
	if paddingA == 0 {
		paddingA = '_'
	}

	lenT := len(strA)

	var sbuf strings.Builder

	tableStrT := "0123456789ABCDEF"

	for i := 0; i < lenT; i++ {
		v := strA[i]

		if !(((v >= '0') && (v <= '9')) || ((v >= 'a') && (v <= 'z')) || ((v >= 'A') && (v <= 'Z'))) {
			sbuf.WriteByte(paddingA)
			sbuf.WriteByte(tableStrT[v>>4])
			sbuf.WriteByte(tableStrT[v&15])
		} else {
			sbuf.WriteByte(strA[i])
		}
	}

	return sbuf.String()
}

var EncodeStringCustom = TKX.EncodeStringCustom

func (pA *TK) EncodeStringCustomEx(strA string, paddingsA ...byte) string {
	var paddingA byte

	if paddingsA == nil || len(paddingsA) == 0 {
		paddingA = '_'
	} else {
		paddingA = paddingsA[0]
	}

	lenT := len(strA)

	var sbuf strings.Builder

	tableStrT := "0123456789ABCDEF"

	for i := 0; i < lenT; i++ {
		v := strA[i]

		if !(((v >= '0') && (v <= '9')) || ((v >= 'a') && (v <= 'z'))) {
			sbuf.WriteByte(paddingA)
			sbuf.WriteByte(tableStrT[v>>4])
			sbuf.WriteByte(tableStrT[v&15])
		} else {
			sbuf.WriteByte(strA[i])
		}
	}

	return sbuf.String()
}

var EncodeStringCustomEx = TKX.EncodeStringCustomEx

func (pA *TK) ishex(c byte) bool {
	switch {
	case '0' <= c && c <= '9':
		return true
	case 'a' <= c && c <= 'f':
		return true
	case 'A' <= c && c <= 'F':
		return true
	}
	return false
}

var ishex = TKX.ishex

func (pA *TK) unhex(c byte) byte {
	switch {
	case '0' <= c && c <= '9':
		return c - '0'
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10
	}
	return 0
}

var unhex = TKX.unhex

func (pA *TK) DecodeStringSimple(s string) string {
	// Count %, check that they're well-formed.
	n := 0
	// hasPlus := false

	for i := 0; i < len(s); {
		switch s[i] {
		case '%':
			n++

			if i+2 >= len(s) || !ishex(s[i+1]) || !ishex(s[i+2]) {
				// s = s[i:]
				// if len(s) > 3 {
				// 	s = s[:3]
				// }
				return s
			}

			i += 3

		default:
			i++
		}
	}

	// if n == 0 && !hasPlus {
	// 	return GenerateErrorString("invalid format")
	// }

	t := make([]byte, len(s)-2*n)
	j := 0
	for i := 0; i < len(s); {
		switch s[i] {
		case '%':
			t[j] = unhex(s[i+1])<<4 | unhex(s[i+2])
			j++
			i += 3
		default:
			t[j] = s[i]
			j++
			i++
		}
	}
	return string(t)
}

var DecodeStringSimple = TKX.DecodeStringSimple

func (pA *TK) DecodeStringUnderline(s string) string {
	var bufT strings.Builder

	lenT := len(s)

	for i := 0; i < lenT; {
		if s[i] == '_' {
			if i+2 >= lenT {
				return s
			}
			bufT.WriteByte(unhex(s[i+1])<<4 | unhex(s[i+2]))

			i += 3
		} else {
			bufT.WriteByte(s[i])
			i++
		}
	}

	return bufT.String()
}

var DecodeStringUnderline = TKX.DecodeStringUnderline

func (pA *TK) DecodeStringCustom(s string, paddingsA ...byte) string {
	var paddingA byte

	if paddingsA == nil || len(paddingsA) == 0 {
		paddingA = '_'
	} else {
		paddingA = paddingsA[0]
	}

	var bufT strings.Builder

	lenT := len(s)

	for i := 0; i < lenT; {
		if s[i] == paddingA {
			if i+2 >= lenT {
				return s
			}
			bufT.WriteByte(unhex(s[i+1])<<4 | unhex(s[i+2]))

			i += 3
		} else {
			bufT.WriteByte(s[i])
			i++
		}
	}

	return bufT.String()
}

var DecodeStringCustom = TKX.DecodeStringCustom

func (pA *TK) MD5Encrypt(strA string) string {
	tmpb := md5.Sum([]byte(strA))

	tmpbb := tmpb[:]

	return hex.EncodeToString(tmpbb)
}

var MD5Encrypt = TKX.MD5Encrypt

func (pA *TK) MD5EncryptToBytes(strA string) []byte {
	tmpb := md5.Sum([]byte(strA))

	tmpbb := tmpb[:]

	return tmpbb
}

var MD5EncryptToBytes = TKX.MD5EncryptToBytes

// 加密解密相关

func (pA *TK) BytesToHex(bufA []byte, optsA ...string) string {
	// Pl("%#v %#v", bufA, optsA)
	if len(optsA) > 0 {
		if IfSwitchExistsWhole(optsA, "-split") {
			var sb strings.Builder

			for i, v := range bufA {
				if i > 0 {
					sb.WriteString(" ")
				}

				sb.WriteString(fmt.Sprintf("%02X", v))
			}

			return sb.String()
		}
	}

	return strings.ToUpper(hex.EncodeToString(bufA))
}

var BytesToHex = TKX.BytesToHex

func (pA *TK) BytesToHexX(bufA []byte) string {
	var sb strings.Builder

	for i, v := range bufA {
		if i > 0 {
			sb.WriteString(" ")
		}

		sb.WriteString(fmt.Sprintf("%02X", v))
	}

	return sb.String()
}

var BytesToHexX = TKX.BytesToHexX

func (pA *TK) HexToBytes(strA string) []byte {
	if strings.HasPrefix(strA, "HEX_") {
		strA = strA[4:]
	}

	buf, err := hex.DecodeString(strA)
	if err != nil {
		return nil
	}

	return buf
}

var HexToBytes = TKX.HexToBytes

func (pA *TK) HexToStr(strA string) string {
	if strings.HasPrefix(strA, "HEX_") {
		strA = strA[4:]
	}

	buf, err := hex.DecodeString(strA)
	if err != nil {
		return ErrToStr(err)
	}

	return string(buf)
}

var HexToStr = TKX.HexToStr

// HexToInt return 0 if failed
func (pA *TK) HexToInt(strA string, optsA ...string) int {
	defaultT := GetSwitchWithDefaultIntValue(optsA, "-default=", 0)

	lenT := len(strA)

	if lenT < 1 {
		return defaultT
	}

	if lenT%2 != 0 {
		strA = "0" + strA
	}

	buf, err := hex.DecodeString(strA)
	if err != nil {
		return defaultT
	}

	// endianStrT := Trim(GetSwitch(optsA, "-endian=", ""))

	// endianT := GetSystemEndianInt()

	// if StartsWith(endianStrT, "B") {
	// 	endianT = 1
	// } else if StartsWith(endianStrT, "L") {
	// 	endianT = 0
	// }

	// if endianT < 1 {

	// }

	lenT = len(buf)

	if lenT < 1 {
		return defaultT
	}

	m := 1

	rs := 0

	for i := lenT - 1; i >= 0; i-- {
		rs += int(buf[i]) * m

		m *= 256
	}

	return rs
}

var HexToInt = TKX.HexToInt

func (pA *TK) OctetToInt(strA string, defaultA ...int) int {
	defaultT := -1

	if len(defaultA) > 0 {
		defaultT = defaultA[0]
	}

	vT, errT := strconv.ParseInt(strA, 8, 64)

	if errT != nil {
		return defaultT
	}

	return int(vT)
}

var OctetToInt = TKX.OctetToInt

// HexToUInt return 0 if failed
func (pA *TK) HexToUInt(strA string, optsA ...string) uint64 {
	defaultT := GetSwitchWithDefaultIntValue(optsA, "-default=", 0)

	lenT := len(strA)

	if lenT < 1 {
		return uint64(defaultT)
	}

	if lenT%2 != 0 {
		strA = "0" + strA
	}

	buf, err := hex.DecodeString(strA)
	if err != nil {
		return uint64(defaultT)
	}

	// endianStrT := Trim(GetSwitch(optsA, "-endian=", ""))

	// endianT := GetSystemEndianInt()

	// if StartsWith(endianStrT, "B") {
	// 	endianT = 1
	// } else if StartsWith(endianStrT, "L") {
	// 	endianT = 0
	// }

	// if endianT < 1 {

	// }

	lenT = len(buf)

	if lenT < 1 {
		return uint64(defaultT)
	}

	var m uint64 = 1

	var rs uint64 = 0

	for i := lenT - 1; i >= 0; i-- {
		rs += uint64(buf[i]) * m

		m *= 256
	}

	return rs
}

var HexToUInt = TKX.HexToUInt

func (pA *TK) GetRandomByte() byte {
	Randomize()

	randT := rand.Intn(256)

	return byte(randT)
}

var GetRandomByte = TKX.GetRandomByte

func (pA *TK) EncryptDataByTXDEE(srcDataA []byte, codeA string) []byte {
	if srcDataA == nil {
		return nil
	}

	dataLen := len(srcDataA)
	if dataLen < 1 {
		return srcDataA
	}

	codeT := codeA
	if codeT == "" {
		codeT = "topxeq"
	}

	codeBytes := []byte(codeT)
	codeLen := len(codeBytes)

	bufB := make([]byte, dataLen+4)

	bufB[0] = GetRandomByte()
	bufB[1] = GetRandomByte()

	for i := 0; i < dataLen; i++ {
		bufB[2+i] = srcDataA[i] + codeBytes[i%codeLen] + byte(i+1) + bufB[1]
	}

	bufB[dataLen+4-2] = GetRandomByte()
	bufB[dataLen+4-1] = GetRandomByte()

	return bufB
}

var EncryptDataByTXDEE = TKX.EncryptDataByTXDEE

func (pA *TK) SumBytes(srcDataA []byte) byte {
	if srcDataA == nil {
		return 0
	}

	lenT := len(srcDataA)

	var b byte = 0

	for i := 0; i < lenT; i++ {
		b += srcDataA[i]
	}

	return b
}

var SumBytes = TKX.SumBytes

func (pA *TK) EncryptDataByTXDEF(srcDataA []byte, codeA ...string) []byte {
	if srcDataA == nil {
		return nil
	}

	dataLen := len(srcDataA)
	if dataLen < 1 {
		return srcDataA
	}

	codeT := ""

	addHeadT := false

	if (codeA != nil) && (len(codeA) > 0) {
		codeT = GetParam(codeA, 0, "") // codeA[0]

		// if len(codeA) > 1 {
		addHeadT = IfSwitchExists(codeA, "-addHead")
		// }

		codeT = GetSwitch(codeA, "-code=", codeT)
	}

	if codeT == "" {
		codeT = "topxeq"
	}

	// Pl("code: %v, addHead: %v", codeT, addHeadT)

	codeBytes := []byte(codeT)
	codeLen := len(codeBytes)

	sumT := int(SumBytes(codeBytes))

	addLenT := int((sumT % 5) + 2)
	encIndexT := sumT % addLenT
	// Plvsr(addLenT, encIndexT)

	headT := []byte("//TXDEF#")

	headLenT := 0

	var bufB []byte

	if addHeadT {
		headLenT = len(headT)
	}

	bufB = make([]byte, dataLen+addLenT+headLenT)

	for i := 0; i < headLenT; i++ {
		bufB[i] = headT[i]
	}

	for i := 0; i < addLenT; i++ {
		bufB[i+headLenT] = GetRandomByte()
	}

	for i := 0; i < dataLen; i++ {
		bufB[addLenT+i+headLenT] = srcDataA[i] + codeBytes[i%codeLen] + byte(i+1) + bufB[encIndexT+headLenT]
	}

	return bufB
}

var EncryptDataByTXDEF = TKX.EncryptDataByTXDEF

const TXDEF_BUFFER_LEN = 1000

func (pA *TK) EncryptStreamByTXDEF(readerA io.Reader, codeA string, writerA io.Writer) error {
	if readerA == nil {
		return Errf("reader nil")
	}

	if writerA == nil {
		return Errf("writer nil")
	}

	codeT := codeA
	if codeT == "" {
		codeT = "topxeq"
	}

	codeBytes := []byte(codeT)
	codeLen := len(codeBytes)

	// if codeLen > TXDEF_BUFFER_LEN {
	// 	return Errf("code length longer than buffer length")
	// }

	sumT := int(SumBytes(codeBytes))

	addLenT := int((sumT % 5) + 2)
	encIndexT := sumT % addLenT
	// Plvsr(addLenT, encIndexT)

	var idxByte byte

	addBufT := make([]byte, addLenT)

	for i := 0; i < addLenT; i++ {
		addBufT[i] = GetRandomByte()

		if i == encIndexT {
			idxByte = addBufT[i]
		}

	}

	_, errT := writerA.Write(addBufT)

	if errT != nil {
		return errT
	}

	bufArrayT := make([]byte, 0, TXDEF_BUFFER_LEN+8)

	bufT := bytes.NewBuffer(bufArrayT)

	bufArrayReaderT := make([]byte, TXDEF_BUFFER_LEN)

	i := 0

	breakFlagT := false

	for {
		bytesT, errT := readerA.Read(bufArrayReaderT)

		if errT != nil {
			if errT == io.EOF {
				breakFlagT = true
			} else {
				return errT
			}

		}

		for j := 0; j < bytesT; j++ {
			bufT.WriteByte(bufArrayReaderT[j] + codeBytes[i%codeLen] + byte(i+1) + idxByte)

			i++

			if bufT.Len() >= TXDEF_BUFFER_LEN {
				_, errT := writerA.Write(bufT.Bytes()[:TXDEF_BUFFER_LEN])

				if errT != nil {
					return errT
				}

				bufT.Reset()
			}
		}

		if breakFlagT {
			if bufT.Len() > 0 {
				_, errT := writerA.Write(bufT.Bytes())

				if errT != nil {
					return errT
				}

				bufT.Reset()
			}

			break
		}
	}

	return nil
}

var EncryptStreamByTXDEF = TKX.EncryptStreamByTXDEF

func (pA *TK) DecryptStreamByTXDEF(readerA io.Reader, codeA string, writerA io.Writer) error {
	if readerA == nil {
		return Errf("reader nil")
	}

	if writerA == nil {
		return Errf("writer nil")
	}

	codeT := codeA
	if codeT == "" {
		codeT = "topxeq"
	}

	codeBytes := []byte(codeT)
	codeLen := len(codeBytes)

	sumT := int(SumBytes(codeBytes))

	addLenT := int((sumT % 5) + 2)
	encIndexT := sumT % addLenT
	// Plvsr(addLenT, encIndexT)

	var idxByte byte

	addBufT := make([]byte, addLenT)

	bytesT, errT := readerA.Read(addBufT)

	if errT != nil {
		if errT != io.EOF {
			return errT
		}
	}

	if bytesT != addLenT {
		return Errf("failed to read header")
	}

	idxByte = addBufT[encIndexT]

	bufArrayT := make([]byte, 0, TXDEF_BUFFER_LEN+8)

	bufT := bytes.NewBuffer(bufArrayT)

	bufArrayReaderT := make([]byte, TXDEF_BUFFER_LEN)

	i := 0

	breakFlagT := false

	for {
		bytesT, errT := readerA.Read(bufArrayReaderT)

		if errT != nil {
			if errT == io.EOF {
				breakFlagT = true
			} else {
				return errT
			}

		}

		for j := 0; j < bytesT; j++ {
			bufT.WriteByte(bufArrayReaderT[j] - codeBytes[i%codeLen] - byte(i+1) - idxByte)

			i++

			if bufT.Len() >= TXDEF_BUFFER_LEN {
				_, errT := writerA.Write(bufT.Bytes()[:TXDEF_BUFFER_LEN])

				if errT != nil {
					return errT
				}

				bufT.Reset()
			}
		}

		if breakFlagT {
			if bufT.Len() > 0 {
				_, errT := writerA.Write(bufT.Bytes())

				if errT != nil {
					return errT
				}

				bufT.Reset()
			}

			break
		}
	}

	return nil
}

var DecryptStreamByTXDEF = TKX.DecryptStreamByTXDEF

func (pA *TK) DecryptDataByTXDEE(srcDataA []byte, codeA string) []byte {
	if srcDataA == nil {
		return nil
	}

	dataLen := len(srcDataA)
	if dataLen < 4 {
		return nil
	}

	dataLen -= 4

	codeT := codeA
	if codeT == "" {
		codeT = "topxeq"
	}

	codeBytes := []byte(codeT)
	codeLen := len(codeBytes)

	bufB := make([]byte, dataLen)

	for i := 0; i < dataLen; i++ {
		bufB[i] = srcDataA[2+i] - codeBytes[i%codeLen] - byte(i+1) - srcDataA[1]
	}

	return bufB
}

var DecryptDataByTXDEE = TKX.DecryptDataByTXDEE

func (pA *TK) IsDataEncryptedByTXDEF(srcDataA []byte) bool {
	return bytes.HasPrefix(srcDataA, []byte("//TXDEF#"))
}

var IsDataEncryptedByTXDEF = TKX.IsDataEncryptedByTXDEF

func (pA *TK) IsStringEncryptedByTXDEF(strA string) bool {
	return strings.HasPrefix(strA, "//TXDEF#")
}

var IsStringEncryptedByTXDEF = TKX.IsStringEncryptedByTXDEF

func (pA *TK) DecryptDataByTXDEF(srcDataA []byte, codeA ...string) []byte {
	if srcDataA == nil {
		return nil
	}

	codeT := ""

	if (codeA != nil) && (len(codeA) > 0) {
		codeT = GetParam(codeA, 0, "") // codeA[0]

		codeT = GetSwitch(codeA, "-code=", codeT)
	}

	if codeT == "" {
		codeT = "topxeq"
	}

	codeBytes := []byte(codeT)
	codeLen := len(codeBytes)

	sumT := int(SumBytes(codeBytes))

	addLenT := int((sumT % 5) + 2)
	encIndexT := sumT % addLenT

	if bytes.HasPrefix(srcDataA, []byte("//TXDEF#")) {
		srcDataA = srcDataA[8:]
	}

	dataLen := len(srcDataA)
	if dataLen < addLenT {
		return nil
	}

	dataLen -= addLenT

	bufB := make([]byte, dataLen)

	for i := 0; i < dataLen; i++ {
		bufB[i] = srcDataA[addLenT+i] - codeBytes[i%codeLen] - byte(i+1) - srcDataA[encIndexT]
	}

	return bufB
}

var DecryptDataByTXDEF = TKX.DecryptDataByTXDEF

func (pA *TK) EncryptStringByTXTE(strA string, codeA string) string {
	if strA == "" {
		return ""
	}

	codeT := codeA
	if codeT == "" {
		codeT = "topxeq"
	}

	sBufT := []byte(strA)
	codeButT := []byte(codeT)

	sDataLen := len(sBufT)
	codeLenT := len(codeButT)

	dBufT := make([]byte, sDataLen)

	for i := 0; i < sDataLen; i++ {
		dBufT[i] = sBufT[i] + codeButT[i%codeLenT] + byte(i+1)
	}

	return strings.ToUpper(hex.EncodeToString(dBufT))

}

var EncryptStringByTXTE = TKX.EncryptStringByTXTE

func (pA *TK) DecryptStringByTXTE(strA string, codeA string) string {
	if strA == "" {
		return ""
	}

	codeT := codeA
	if codeT == "" {
		codeT = "topxeq"
	}

	sBufT, errT := hex.DecodeString(strA)
	if errT != nil {
		return ""
	}
	codeButT := []byte(codeT)

	sDataLen := len(sBufT)
	codeLenT := len(codeButT)

	dBufT := make([]byte, sDataLen)

	for i := 0; i < sDataLen; i++ {
		dBufT[i] = sBufT[i] - codeButT[i%codeLenT] - byte(i+1)
	}

	return string(dBufT)

}

var DecryptStringByTXTE = TKX.DecryptStringByTXTE

func (pA *TK) EncryptStringByTXDEE(strA, codeA string) string {
	if strA == "" {
		return ""
	}

	dataDT := EncryptDataByTXDEE([]byte(strA), codeA)
	if dataDT == nil {
		return GenerateErrorStringF("encrypting failed")
	}

	return strings.ToUpper(hex.EncodeToString(dataDT))
}

var EncryptStringByTXDEE = TKX.EncryptStringByTXDEE

func (pA *TK) DecryptStringByTXDEE(strA, codeA string) string {
	if strA == "" {
		return ""
	}

	var sBufT []byte
	var errT error

	if StartsWith(strA, "740404") {
		sBufT, errT = hex.DecodeString(strA[6:])
	} else {
		sBufT, errT = hex.DecodeString(strA)
	}

	if errT != nil {
		return GenerateErrorStringF("decrypting failed")
	}

	dataDT := DecryptDataByTXDEE(sBufT, codeA)
	if dataDT == nil {
		return GenerateErrorStringF("decrypting failed")
	}

	return string(dataDT)
}

var DecryptStringByTXDEE = TKX.DecryptStringByTXDEE

func (pA *TK) EncryptStringByTXDEF(strA string, codeA ...string) string {
	if strA == "" {
		return ""
	}

	codeT := ""

	if (codeA != nil) && (len(codeA) > 0) {
		codeT = codeA[0]
	}

	dataDT := EncryptDataByTXDEF([]byte(strA), codeT)
	if dataDT == nil {
		return GenerateErrorStringF("encrypting failed")
	}

	return strings.ToUpper(hex.EncodeToString(dataDT))
}

var EncryptStringByTXDEF = TKX.EncryptStringByTXDEF

func (pA *TK) DecryptStringByTXDEF(strA string, codeA ...string) string {
	if strA == "" {
		return ""
	}

	codeT := ""

	if (codeA != nil) && (len(codeA) > 0) {
		codeT = codeA[0]
	}

	var sBufT []byte
	var errT error

	if StartsWith(strA, "740404") {
		sBufT, errT = hex.DecodeString(strA[6:])
	} else if StartsWith(strA, "//TXDEF#") {
		sBufT, errT = hex.DecodeString(strA[8:])
	} else {
		sBufT, errT = hex.DecodeString(strA)
	}

	if errT != nil {
		return GenerateErrorStringF("decrypting failed: %v", errT)
	}

	dataDT := DecryptDataByTXDEF(sBufT, codeT)
	if dataDT == nil {
		return GenerateErrorStringF("decrypting failed")
	}

	return string(dataDT)
}

var DecryptStringByTXDEF = TKX.DecryptStringByTXDEF

func (pA *TK) EncryptFileByTXDEF(fileNameA, codeA, outputFileA string) error {
	if !IfFileExists(fileNameA) {
		return Errf("")
	}

	srcStatT, errT := os.Stat(fileNameA)
	if errT != nil {
		return Errf("error os.Stat src %s: %s", fileNameA, errT.Error())
	}

	codeT := codeA
	if codeT == "" {
		codeT = "topxeq"
	}

	outputFileT := outputFileA
	if outputFileT == "" {
		outputFileT = fileNameA + ".txdef"
	}

	fileContenT, errT := os.ReadFile(fileNameA)
	if errT != nil {
		return errT
	}

	writeContentT := EncryptDataByTXDEF(fileContenT, codeT)
	if writeContentT == nil {
		return Errf("encrypt data failed")
	}

	errT = os.WriteFile(outputFileT, writeContentT, srcStatT.Mode())
	if errT != nil {
		return errT
	}

	return nil
}

var EncryptFileByTXDEF = TKX.EncryptFileByTXDEF

func (pA *TK) EncryptFileByTXDEFWithHeader(fileNameA, codeA, outputFileA string) error {
	if !IfFileExists(fileNameA) {
		return Errf("")
	}

	srcStatT, errT := os.Stat(fileNameA)
	if errT != nil {
		return Errf("error os.Stat src %s: %s", fileNameA, errT.Error())
	}

	codeT := codeA
	if codeT == "" {
		codeT = "topxeq"
	}

	outputFileT := outputFileA
	if outputFileT == "" {
		outputFileT = fileNameA + ".txdef"
	}

	fileContenT, errT := os.ReadFile(fileNameA)
	if errT != nil {
		return errT
	}

	writeContentT := EncryptDataByTXDEF(fileContenT, codeT)
	if writeContentT == nil {
		return Errf("encrypt data failed")
	}

	bufT := []byte("//TXDEF#")

	bufT = append(bufT, writeContentT...)

	errT = os.WriteFile(outputFileT, bufT, srcStatT.Mode())
	if errT != nil {
		return errT
	}

	return nil
}

var EncryptFileByTXDEFWithHeader = TKX.EncryptFileByTXDEFWithHeader

func (pA *TK) EncryptFileByTXDEFStream(fileNameA, codeA, outputFileA string) error {
	if !IfFileExists(fileNameA) {
		return Errf("")
	}

	srcStatT, errT := os.Stat(fileNameA)
	if errT != nil {
		return Errf("error os.Stat src %s: %s", fileNameA, errT.Error())
	}

	codeT := codeA
	if codeT == "" {
		codeT = "topxeq"
	}

	outputFileT := outputFileA
	if outputFileT == "" {
		outputFileT = fileNameA + ".txdef"
	}

	fileT, errT := os.Open(fileNameA)
	if errT != nil {
		return Errf("failed to open source file: %v", errT)
	}

	defer fileT.Close()

	file2T, errT := os.OpenFile(outputFileT, os.O_CREATE, srcStatT.Mode())
	if errT != nil {
		return Errf("failed to create target file: %v", errT)
	}

	defer file2T.Close()

	bufT := bufio.NewWriter(file2T)

	errT = EncryptStreamByTXDEF(fileT, codeT, bufT)

	if errT != nil {
		return Errf("failed to encrypt: %v", errT)
	}

	errT = bufT.Flush()

	if errT != nil {
		return Errf("failed to flush output file: %v", errT)
	}

	return nil
}

var EncryptFileByTXDEFStream = TKX.EncryptFileByTXDEFStream

func (pA *TK) DecryptFileByTXDEFStream(fileNameA, codeA, outputFileA string) error {
	if !IfFileExists(fileNameA) {
		return Errf("")
	}

	srcStatT, errT := os.Stat(fileNameA)
	if errT != nil {
		return Errf("error os.Stat src %s: %s", fileNameA, errT.Error())
	}

	codeT := codeA
	if codeT == "" {
		codeT = "topxeq"
	}

	outputFileT := outputFileA
	if outputFileT == "" {
		outputFileT = fileNameA + ".untxdef"
	}

	fileT, errT := os.Open(fileNameA)
	if errT != nil {
		return Errf("failed to open source file: %v", errT)
	}

	defer fileT.Close()

	file2T, errT := os.OpenFile(outputFileT, os.O_CREATE, srcStatT.Mode())
	if errT != nil {
		return Errf("failed to create target file: %v", errT)
	}

	defer file2T.Close()

	bufT := bufio.NewWriter(file2T)

	errT = DecryptStreamByTXDEF(fileT, codeT, bufT)

	if errT != nil {
		return Errf("failed to decrypt: %v", errT)
	}

	errT = bufT.Flush()

	if errT != nil {
		return Errf("failed to flush output file: %v", errT)
	}

	return nil
}

var DecryptFileByTXDEFStream = TKX.DecryptFileByTXDEFStream

func (pA *TK) ErrorToString(errA error) string {
	if errA == nil {
		return ""
	}

	return "TXERROR:" + errA.Error()
}

var ErrorToString = TKX.ErrorToString

func (pA *TK) ErrorToEmptyString(vA interface{}) string {
	if vA == nil {
		return ""
	}

	_, ok := vA.(error)
	if ok {
		return ""
	}

	return ToStr(vA)
}

var ErrorToEmptyString = TKX.ErrorToEmptyString

func (pA *TK) ErrToErrStr(errA error) string {
	if errA == nil {
		return ""
	}

	return GenerateErrorString(errA.Error())
}

var ErrToErrStr = TKX.ErrToErrStr

func (pA *TK) EncryptFileByTXDEFS(fileNameA, codeA, outputFileA string) string {
	return ErrorToString(EncryptFileByTXDEF(fileNameA, codeA, outputFileA))
}

var EncryptFileByTXDEFS = TKX.EncryptFileByTXDEFS

func (pA *TK) EncryptFileByTXDEFStreamS(fileNameA, codeA, outputFileA string) string {
	return ErrorToString(EncryptFileByTXDEFStream(fileNameA, codeA, outputFileA))
}

var EncryptFileByTXDEFStreamS = TKX.EncryptFileByTXDEFStreamS

func (pA *TK) DecryptFileByTXDEF(fileNameA, codeA, outputFileA string) error {
	if !IfFileExists(fileNameA) {
		return Errf("file not exists")
	}

	srcStatT, errT := os.Stat(fileNameA)
	if errT != nil {
		return Errf("error os.Stat src %s: %s", fileNameA, errT.Error())
	}

	codeT := codeA
	if codeT == "" {
		codeT = "topxeq"
	}

	outputFileT := outputFileA
	if outputFileT == "" {
		outputFileT = fileNameA + ".untxdef"
	}

	fileContenT, errT := os.ReadFile(fileNameA)
	if errT != nil {
		return errT
	}

	if bytes.HasPrefix(fileContenT, []byte("//TXDEF#")) {
		fileContenT = fileContenT[8:]
	}

	writeContentT := DecryptDataByTXDEF(fileContenT, codeT)
	if writeContentT == nil {
		return Errf("decrypt data failed")
	}

	errT = os.WriteFile(outputFileT, writeContentT, srcStatT.Mode())
	if errT != nil {
		return errT
	}

	return nil

}

var DecryptFileByTXDEF = TKX.DecryptFileByTXDEF

func (pA *TK) DecryptFileByTXDEFS(fileNameA, codeA, outputFileA string) string {
	return ErrorToString(DecryptFileByTXDEF(fileNameA, codeA, outputFileA))
}

var DecryptFileByTXDEFS = TKX.DecryptFileByTXDEFS

func (pA *TK) DecryptFileByTXDEFStreamS(fileNameA, codeA, outputFileA string) string {
	return ErrorToString(DecryptFileByTXDEFStream(fileNameA, codeA, outputFileA))
}

var DecryptFileByTXDEFStreamS = TKX.DecryptFileByTXDEFStreamS

func (pA *TK) Pkcs7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	//	Pl("padding: %v", padding)
	return append(ciphertext, padtext...)
}

var Pkcs7Padding = TKX.Pkcs7Padding

// func Pkcs7Unpad(b []byte, blocksize int) ([]byte, error) {
// 	if blocksize <= 0 {
// 		return nil, Errf("ErrInvalidBlockSize")
// 	}
// 	if b == nil || len(b) == 0 {
// 		return nil, Errf("ErrInvalidPKCS7Data")
// 	}
// 	if len(b)%blocksize != 0 {
// 		return nil, Errf("ErrInvalidPKCS7Padding")
// 	}
// 	c := b[len(b)-1]
// 	n := int(c)
// 	if n == 0 || n > len(b) {
// 		return nil, Errf("ErrInvalidPKCS7Padding")
// 	}
// 	for i := 0; i < n; i++ {
// 		if b[len(b)-n+i] != c {
// 			return nil, Errf("ErrInvalidPKCS7Padding")
// 		}
// 	}
// 	return b[:len(b)-n], nil
// }

func (pA *TK) AESEncrypt(src, key []byte) ([]byte, error) {
	//	key = toMD5(key)
	keyT := key
	if len(keyT) > 16 {
		keyT = keyT[0:16]
	}
	block, err := aes.NewCipher(keyT)

	if err != nil {
		return nil, err
	}
	bs := block.BlockSize()
	//	Printf("Src: %v\n", src)
	//	Printf("Key: %v\n", key)
	//	Printf("Block size: %v\n", bs)
	//	src = zeroPadding(src, bs)
	src = Pkcs7Padding(src, bs)
	//	Pl("After padding: %v", src)
	if len(src)%bs != 0 {
		return nil, errors.New("Need a multiple of the blocksize")
	}
	out := make([]byte, len(src))
	dst := out
	iv := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	for len(src) > 0 {
		//		Pl("Encing: %v\n", src[:bs])
		for j := 0; j < bs; j++ {
			src[j] ^= iv[j]
		}
		//		Pl("EncingXORed: %v", src[:bs])
		block.Encrypt(dst, src[:bs])
		src = src[bs:]

		for j := 0; j < bs; j++ {
			iv[j] = dst[j]
		}
		dst = dst[bs:]
	}
	return out, nil
}

var AESEncrypt = TKX.AESEncrypt

// func Unpad(src []byte) ([]byte, error) {
// 	length := len(src)
// 	unpadding := int(src[length-1])

// 	if unpadding > length {
// 		return nil, errors.New("unpad error. This could happen when incorrect encryption key is used")
// 	}

// 	return src[:(length - unpadding)], nil
// }

func (pA *TK) AESDecrypt(src, key []byte) ([]byte, error) {
	//	key = toMD5(key)
	keyT := key
	if len(keyT) > 16 {
		keyT = keyT[0:16]
	}

	block, err := aes.NewCipher(keyT)

	if err != nil {
		return nil, err
	}

	bs := block.BlockSize()
	//	Printf("Src: %v\n", src)
	//	Printf("Key: %v\n", key)
	// Pl("Block size: %v", bs)
	//	src = zeroPadding(src, bs)
	// beforeLen := len(src)
	// // src = Pkcs7Padding(src, bs)
	// afterLen := len(src)
	// diffLen := afterLen - beforeLen
	// Pl("beforeLen: %v, afterLen: %v, diffLen: %v", beforeLen, afterLen, diffLen)
	//	Pl("After padding: %v", src)
	if len(src)%bs != 0 {
		return nil, errors.New("Need a multiple of the blocksize")
	}

	out := make([]byte, len(src))
	dst := out

	iv := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	for len(src) > 0 {
		//		Pl("EncingXORed: %v", src[:bs])
		block.Decrypt(dst, src[:bs])

		//		Pl("Encing: %v\n", src[:bs])
		for j := 0; j < bs; j++ {
			dst[j] ^= iv[j]
		}

		for j := 0; j < bs; j++ {
			iv[j] = src[j]
		}

		src = src[bs:]

		dst = dst[bs:]
	}

	// Pl("out 1: %#v", out)
	outLenT := len(out)
	unpadLenT := int(out[outLenT-1])

	if unpadLenT <= outLenT {
		for i := 0; i < unpadLenT; i++ {
			if out[outLenT-1-i] != byte(unpadLenT) {
				return out, nil
			}
		}

		out = out[:outLenT-unpadLenT]
	}

	// Pl("out len: %v", len(out))
	// out = out[:len(out)-diffLen]
	return out, nil
}

var AESDecrypt = TKX.AESDecrypt

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func (pA *TK) AESEncryptCBC(src, key []byte) (bytesR []byte, errR error) {
	defer func() {
		r := recover()

		if r != nil {
			bytesR = nil
			errR = fmt.Errorf("%v(%v, %v)", r, src, key)
			return
		}
	}()

	block, _ := aes.NewCipher(key)
	blockSize := block.BlockSize()
	src = pkcs5Padding(src, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	encrypted := make([]byte, len(src))
	blockMode.CryptBlocks(encrypted, src) 
	return encrypted, nil
}

var AESEncryptCBC = TKX.AESEncryptCBC

func (pA *TK) AESDecryptCBC(src, key []byte) (bytesR []byte, errR error) {
	defer func() {
		r := recover()

		if r != nil {
			bytesR = nil
			errR = fmt.Errorf("%v", r)
			return
		}
	}()

	block, err := aes.NewCipher(key)
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	decrypted := make([]byte, len(src))
	blockMode.CryptBlocks(decrypted, src) 
	decrypted = pkcs5UnPadding(decrypted) 
	
	return decrypted, err
}

var AESDecryptCBC = TKX.AESDecryptCBC

// func AESDecrypt(src, key []byte) ([]byte, error) {
// 	//	key = toMD5(key)
// 	keyT := key
// 	if len(keyT) > 16 {
// 		keyT = keyT[0:16]
// 	}

// 	block, err := aes.NewCipher(keyT)

// 	if err != nil {
// 		return nil, err
// 	}

// 	bs := block.BlockSize()
// 	//	Printf("Src: %v\n", src)
// 	//	Printf("Key: %v\n", key)
// 	//	Printf("Block size: %v\n", bs)
// 	//	src = zeroPadding(src, bs)
// 	beforeLen := len(src)
// 	src = Pkcs7Padding(src, bs)
// 	// src, errT := Pkcs7Unpad(src, bs)
// 	// if errT != nil {
// 	// 	return nil, errT
// 	// }
// 	Pl("src: %#v", src)
// 	afterLen := len(src)

// 	diffLen := afterLen - beforeLen
// 	lenT := len(src)

// 	if lenT < 1 {
// 		return nil, errors.New("Zero length")
// 	}

// 	//	Pl("beforeLen: %v, afterLen: %v, diffLen: %v", beforeLen, afterLen, diffLen)
// 	//	Pl("After padding: %v", src)
// 	if lenT%bs != 0 {
// 		return nil, errors.New("Need a multiple of the blocksize")
// 	}

// 	out := make([]byte, lenT)
// 	dst := out
// 	iv := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
// 	for lenT > 0 {
// 		//		Pl("EncingXORed: %v", src[:bs])
// 		block.Decrypt(dst, src[:bs])

// 		//		Pl("Encing: %v\n", src[:bs])
// 		for j := 0; j < bs; j++ {
// 			dst[j] ^= iv[j]
// 		}

// 		for j := 0; j < bs; j++ {
// 			iv[j] = src[j]
// 		}

// 		src = src[bs:]

// 		dst = dst[bs:]
// 	}
// 	// diffLen := int(src[lenT-1])
// 	// Pl("diffLen: %v", diffLen)
// 	out = out[:len(out)-diffLen]
// 	// out, errT = Pkcs7Unpad(out, bs)

// 	// if errT != nil {
// 	// 	return nil, errT
// 	// }

// 	return out, nil
// }

// URL相关 url related

func (pA *TK) AnalyzeURLParams(strA string) map[string]string {
	rMapT := make(map[string]string)

	tmpL := strings.Split(strA, "__")

	for i := range tmpL {
		tmpL2 := strings.SplitN(tmpL[i], "=", 2)

		if len(tmpL2) < 2 {
			continue
		}

		rMapT[tmpL2[0]] = tmpL2[1]
	}

	return rMapT
}

var AnalyzeURLParams = TKX.AnalyzeURLParams

func (pA *TK) UrlEncode(strA string) string {
	return url.QueryEscape(strA)
}

var UrlEncode = TKX.UrlEncode

func (pA *TK) UrlEncode2(strA string) string {
	u, err := url.Parse(strA)
	if err != nil {
		return GenerateErrorString("parsing URL failed")
	}
	return u.String()
}

var UrlEncode2 = TKX.UrlEncode2

func (pA *TK) UrlDecode(strA string) string {
	rStrT, errT := url.QueryUnescape(strA)
	if errT != nil {
		return GenerateErrorString(errT.Error())
	}
	return rStrT
}

var UrlDecode = TKX.UrlDecode

// JoinURL concat a base URL and a relative URL
func (pA *TK) JoinURL(urlBaseA string, urlNextA string) string {
	u, err := url.Parse(urlNextA)
	if err != nil {
		return GenerateErrorString(err.Error())
	}

	base, err := url.Parse(urlBaseA)
	if err != nil {
		return GenerateErrorString(err.Error())
	}
	return base.ResolveReference(u).String()
}

var JoinURL = TKX.JoinURL

func (pA *TK) FormToMap(formA url.Values) map[string]string {
	mapT := make(map[string]string, 0)

	if formA == nil {
		return mapT
	}

	for k, v := range formA {
		mapT[k] = v[0]
	}

	return mapT
}

var FormToMap = TKX.FormToMap

// debug related

var DebugModeG bool = false
var debugLockG sync.Mutex
var debugBufferG bytes.Buffer

func (pA *TK) AddDebug(strA string) {
	if DebugModeG {
		debugLockG.Lock()
		debugBufferG.WriteString(strA + "\n")
		debugLockG.Unlock()
	}
}

var AddDebug = TKX.AddDebug

func (pA *TK) AddDebugF(formatA string, argsA ...interface{}) {
	if !DebugModeG {
		return
	}

	debugLockG.Lock()

	debugBufferG.WriteString(fmt.Sprintf(fmt.Sprintf("[%v] ", GetNowTimeStringFormal())+formatA+"\n", argsA...))

	debugLockG.Unlock()
}

var AddDebugF = TKX.AddDebugF

func (pA *TK) ClearDebug() {
	debugBufferG.Truncate(0)
}

var ClearDebug = TKX.ClearDebug

func (pA *TK) GetDebug() string {
	return debugBufferG.String()
}

var GetDebug = TKX.GetDebug

// http/web service related

func (pA *TK) NewStaticWebHandler(pathA string) interface{} {
	if !IfFileExists(pathA) || !IsDirectory(pathA) {
		return fmt.Errorf("dir not exists")
	}

	var staticFST http.Handler = http.FileServer(http.Dir(pathA))

	serveStaticDirHandlerT := func(w http.ResponseWriter, r *http.Request) {
		old := r.URL.Path

		// tk.Pl("urlPath: %v", r.URL.Path)

		name := filepath.Join(pathA, path.Clean(old))

		// tk.Pl("name: %v", name)

		info, err := os.Lstat(name)
		if err == nil {
			if !info.IsDir() {
				staticFST.ServeHTTP(w, r)
				// http.ServeFile(w, r, name)
			} else {
				if IfFileExists(filepath.Join(name, "index.html")) {
					staticFST.ServeHTTP(w, r)
				} else {
					http.NotFound(w, r)
				}
			}
		} else {
			http.NotFound(w, r)
		}
	}

	return serveStaticDirHandlerT
}

var NewStaticWebHandler = TKX.NewStaticWebHandler

func (pA *TK) SetResponseHeader(resA http.ResponseWriter, keyA string, valueA string) error {
	resA.Header().Set(keyA, valueA)

	return nil
}

var SetResponseHeader = TKX.SetResponseHeader

func (pA *TK) WriteResponseHeader(resA http.ResponseWriter, argsA ...interface{}) error {
	defaultT := "200"
	if len(argsA) > 0 {
		defaultT = fmt.Sprintf("%v", argsA[0])
	}

	resA.WriteHeader(StrToIntWithDefaultValue(defaultT, 200))

	return nil
}

var WriteResponseHeader = TKX.WriteResponseHeader

func (pA *TK) WriteResponse(resA http.ResponseWriter, strA string) error {
	_, errT := resA.Write([]byte(strA))

	return errT
}

var WriteResponse = TKX.WriteResponse

func (pA *TK) DownloadPageUTF8(urlA string, postDataA url.Values, customHeaders string, timeoutSecsA time.Duration, optsA ...string) string {
	client := &http.Client{
		//CheckRedirect: redirectPolicyFunc,
		Timeout: time.Second * timeoutSecsA,
	}

	var urlT string
	if !StartsWithIgnoreCase(urlA, "http") {
		urlT = "http://" + urlA
	} else {
		urlT = urlA
	}

	var respT *http.Response
	var errT error
	var req *http.Request

	if !IsEmptyTrim(customHeaders) {
		if postDataA == nil {
			req, errT = http.NewRequest("GET", urlT, nil)
		} else {
			req, errT = http.NewRequest("POST", urlT, nil)
			req.PostForm = postDataA
		}

		headersT := SplitLines(customHeaders)

		for i := 0; i < len(headersT); i++ {
			lineT := headersT[i]
			aryT := strings.Split(lineT, ":")
			if len(aryT) < 2 {
				continue
			}
			req.Header.Add(aryT[0], Trim(Replace(aryT[1], "`", ":")))
			//TXPl("%s=%s", aryT[0], aryT[1])
		}

		if IfSwitchExistsWhole(optsA, "-verbose") {
			Pl("REQ: %v", req)
		}

		respT, errT = client.Do(req)
	} else {
		if IfSwitchExistsWhole(optsA, "-verbose") {
			Pl("URL: %v", urlT)
		}

		if postDataA == nil {
			respT, errT = client.Get(urlT)
		} else {
			if IfSwitchExistsWhole(optsA, "-verbose") {
				Pl("POST data: %v", postDataA)
			}

			respT, errT = client.PostForm(urlT, postDataA)
		}
	}

	if errT == nil {
		defer respT.Body.Close()
		if respT.StatusCode != 200 {
			if IfSwitchExistsWhole(optsA, "-detail") {
				Pl("response status: %v (%v)", respT.StatusCode, respT)
			}

			return GenerateErrorStringF("response status: %v", respT.StatusCode)
		} else {
			body, errT := io.ReadAll(respT.Body)

			if errT == nil {
				return string(body)
			} else {
				return GenerateErrorString(errT.Error())
			}
		}
	} else {
		return GenerateErrorString(errT.Error())
	}
}

var DownloadPageUTF8 = TKX.DownloadPageUTF8

func (pA *TK) DownloadWebPage(urlA string, postDataA map[string]string, customHeadersA map[string]string, optsA ...string) string {
	timeoutStrT := GetSwitch(optsA, "-timeout=", "15")

	encodingT := GetSwitch(optsA, "-encoding=", "")

	timeoutSecsT := time.Second * time.Duration(StrToInt(timeoutStrT, 15))

	client := &http.Client{
		//CheckRedirect: redirectPolicyFunc,
		Timeout: timeoutSecsT,
	}

	var urlT string
	if !StartsWithIgnoreCase(urlA, "http") {
		urlT = "http://" + urlA
	} else {
		urlT = urlA
	}

	var respT *http.Response
	var errT error
	var req *http.Request

	reqTypeT := GetSwitch(optsA, "-method=", "GET")

	if postDataA != nil && len(postDataA) > 0 && reqTypeT == "GET" {
		reqTypeT = "POST"
	}

	if reqTypeT == "POST" {
		req, errT = http.NewRequest(reqTypeT, urlT, bytes.NewBufferString(MapToPostData(postDataA).Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
		// req.PostForm = MapToPostData(postDataA)
	} else {
		req, errT = http.NewRequest(reqTypeT, urlT, nil)
	}

	if customHeadersA != nil && len(customHeadersA) > 0 {
		for k, v := range customHeadersA {
			req.Header.Add(k, v)
		}
	}

	if IfSwitchExistsWhole(optsA, "-verbose") {
		Pl("REQ: %v", req)
	}

	respT, errT = client.Do(req)

	if errT == nil {
		defer respT.Body.Close()
		if respT.StatusCode != 200 {
			if IfSwitchExistsWhole(optsA, "-detail") {
				Pl("response status: %v (%v)", respT.StatusCode, respT)
			}

			return GenerateErrorStringF("response status: %v", respT.StatusCode)
		} else {
			body, errT := io.ReadAll(respT.Body)

			if errT != nil {
				return GenerateErrorString(errT.Error())
			}

			if (encodingT == "") || (strings.ToLower(encodingT) == "utf-8") {
				return string(body)
			} else {
				return ConvertToUTF8(body, encodingT)
			}
		}
	} else {
		return GenerateErrorString(errT.Error())
	}
}

var DownloadWebPage = TKX.DownloadWebPage

// DownloadWebPageX deprecated, recommend to use getWeb instead
func (pA *TK) DownloadWebPageX(urlA string, optsA ...interface{}) string {
	timeoutStrT := GetSwitchI(optsA, "-timeout=", "15")

	encodingT := GetSwitchI(optsA, "-encoding=", "")

	timeoutSecsT := time.Second * time.Duration(StrToInt(timeoutStrT, 15))

	proxyT := GetSwitchI(optsA, "-proxy=", "")

	client := &http.Client{
		//CheckRedirect: redirectPolicyFunc,
		Timeout: timeoutSecsT,
	}

	if proxyT != "" {
		urli := url.URL{}
		urlproxy, errT := urli.Parse(proxyT)

		if errT != nil {
			return ErrorToString(errT)
		}

		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(urlproxy),
		}
	}

	var urlT string
	if !StartsWithIgnoreCase(urlA, "http") {
		urlT = "http://" + urlA
	} else {
		urlT = urlA
	}

	var respT *http.Response
	var errT error
	var req *http.Request

	reqTypeT := GetSwitchI(optsA, "-method=", "GET")

	var postDataT url.Values = nil

	postStrT := GetSwitchI(optsA, "-post=", "")

	var postBytesT []byte = nil

	if postStrT != "" {
		postObjT, errT := FromJSON(postStrT)

		if errT == nil {
			postMapT, ok := postObjT.(map[string]interface{})

			if ok {
				postDataT = MapToPostDataI(postMapT)
			}
		}
	}

	if postDataT == nil {
		for _, v := range optsA {
			rv, ok := v.(url.Values)

			if ok {
				postDataT = rv
				break
			}

			rv2, ok := v.(map[string]string)

			if ok {
				postDataT = MapToPostData(rv2)
				break
			}

			rv3, ok := v.(map[string]interface{})

			if ok {
				postDataT = MapToPostDataI(rv3)
				break
			}

			rv4, ok := v.([]byte)

			if ok {
				postBytesT = rv4
				break
			}
		}
	}

	postBodyT := GetSwitchI(optsA, "-postBody=", "")

	if (postDataT != nil && len(postDataT) > 0) || postBodyT != "" || postBytesT != nil {
		if reqTypeT != "PUT" {
			reqTypeT = "POST"
		}

	}

	postFileT := GetSwitchI(optsA, "-postFile=", "")

	if reqTypeT == "POST" || reqTypeT == "PUT" {
		if postBytesT != nil {
			req, errT = http.NewRequest(reqTypeT, urlT, bytes.NewReader(postBytesT))

		} else if postBodyT != "" {
			req, errT = http.NewRequest(reqTypeT, urlT, strings.NewReader(postBodyT))

		} else if postFileT != "" {
			// file1, err := os.Open(postFileT)
			// if err != nil {
			// 	return GenerateErrorStringF("failed to load file content: %v", err)
			// }

			// defer file1.Close()

			bufT := LoadBytes(postFileT)

			// req, errT = http.NewRequest(reqTypeT, urlT, bufio.NewReader(file1))
			req, errT = http.NewRequest(reqTypeT, urlT, bytes.NewReader(bufT))

		} else {
			req, errT = http.NewRequest(reqTypeT, urlT, bytes.NewBufferString(postDataT.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
		}
		// req.PostForm = MapToPostData(postDataA)
	} else {
		req, errT = http.NewRequest(reqTypeT, urlT, nil)
	}

	customHeadersStrT := Trim(GetSwitchI(optsA, "-headers=", ""))

	if customHeadersStrT != "" {
		headerObjT, errT := FromJSON(customHeadersStrT)

		if errT == nil {
			headerMapT, ok := headerObjT.(map[string]interface{})

			if ok {
				for k, v := range headerMapT {
					s, ok := v.(string)
					if ok {
						req.Header.Add(k, s)
					}
				}
			}
		}
	}

	if IfSwitchExistsWholeI(optsA, "-verbose") {
		Pl("REQ: %v", req)
	}

	respT, errT = client.Do(req)

	if errT == nil {
		defer respT.Body.Close()
		if respT.StatusCode != 200 {
			if IfSwitchExistsWholeI(optsA, "-detail") {
				body, errT := io.ReadAll(respT.Body)

				if errT != nil {
					body = []byte(errT.Error())
				}
				Pl("response status: %v (%v) body: %v", respT.StatusCode, respT, string(body))
			}

			return GenerateErrorStringF("response status: %v", respT.StatusCode)
		} else {
			body, errT := io.ReadAll(respT.Body)

			if errT != nil {
				return GenerateErrorString(errT.Error())
			}

			if (encodingT == "") || (strings.ToLower(encodingT) == "utf-8") {
				return string(body)
			} else {
				return ConvertToUTF8(body, encodingT)
			}
		}
	} else {
		return GenerateErrorString(errT.Error())
	}
}

var DownloadWebPageX = TKX.DownloadWebPageX

func (pA *TK) GetWebHead(urlA string, optsA ...interface{}) interface{} {
	respT, errT := http.Head(urlA)
	if errT != nil {
		return errT
	}

	return respT
}

var GetWebHead = TKX.GetWebHead

func (pA *TK) UrlExists(urlA string, optsA ...interface{}) interface{} {
	respT, errT := http.Head(urlA)
	if errT != nil {
		return errT
	}

	if respT.StatusCode != http.StatusOK {
		return false
	}

	return true
}

var UrlExists = TKX.UrlExists

func (pA *TK) GetWeb(urlA string, optsA ...interface{}) interface{} {
	// Pl("%#v", optsA)
	timeoutStrT := GetSwitchI(optsA, "-timeout=", "15")

	encodingT := GetSwitchI(optsA, "-encoding=", "")

	proxyT := GetSwitchI(optsA, "-proxy=", "")

	timeoutSecsT := time.Second * time.Duration(StrToInt(timeoutStrT, 15))

	client := &http.Client{
		//CheckRedirect: redirectPolicyFunc,
		Timeout: timeoutSecsT, // time.Second * timeoutSecsT,
	}

	if proxyT != "" {
		urli := url.URL{}
		urlproxy, errT := urli.Parse(proxyT)

		if errT != nil {
			return errT
		}

		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(urlproxy),
		}
	}

	var urlT string
	if !StartsWithIgnoreCase(urlA, "http") {
		urlT = "http://" + urlA
	} else {
		urlT = urlA
	}

	var respT *http.Response
	var errT error
	var req *http.Request

	reqTypeT := GetSwitchI(optsA, "-method=", "GET")

	var postDataT url.Values = nil

	postStrT := GetSwitchI(optsA, "-post=", "")

	var postBytesT []byte = nil

	if postStrT != "" {
		if strings.HasPrefix(postStrT, "HEX_") {
			postStrT = HexToStr(postStrT)
		}

		postObjT, errT := FromJSON(postStrT)

		if errT == nil {
			postMapT, ok := postObjT.(map[string]interface{})

			if ok {
				postDataT = MapToPostDataI(postMapT)
			}
		}
	}

	if postDataT == nil {
		for _, v := range optsA {
			rv, ok := v.(url.Values)

			if ok {
				postDataT = rv
				break
			}

			rv2, ok := v.(map[string]string)

			if ok {
				postDataT = MapToPostData(rv2)
				break
			}

			rv3, ok := v.(map[string]interface{})

			if ok {
				postDataT = MapToPostDataI(rv3)
				break
			}

			rv4, ok := v.([]byte)

			if ok {
				postBytesT = rv4
				break
			}
		}
	}

	postBodyT := GetSwitchI(optsA, "-postBody=", "")

	if strings.HasPrefix(postBodyT, "HEX_") {
		postBodyT = HexToStr(postBodyT)
	}

	if (postDataT != nil && len(postDataT) > 0) || postBodyT != "" || postBytesT != nil {
		if reqTypeT != "PUT" {
			reqTypeT = "POST"
		}

	}

	postFileT := GetSwitchI(optsA, "-postFile=", "")

	if reqTypeT == "POST" || reqTypeT == "PUT" {
		if postBytesT != nil {
			req, errT = http.NewRequest(reqTypeT, urlT, bytes.NewReader(postBytesT))

		} else if postBodyT != "" {
			req, errT = http.NewRequest(reqTypeT, urlT, strings.NewReader(postBodyT))

		} else if postFileT != "" {
			// file1, err := os.Open(postFileT)
			// if err != nil {
			// 	return GenerateErrorStringF("failed to load file content: %v", err)
			// }

			// defer file1.Close()

			bufT := LoadBytes(postFileT)

			// req, errT = http.NewRequest(reqTypeT, urlT, bufio.NewReader(file1))
			req, errT = http.NewRequest(reqTypeT, urlT, bytes.NewReader(bufT))

		} else {
			req, errT = http.NewRequest(reqTypeT, urlT, bytes.NewBufferString(postDataT.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
		}
		// req.PostForm = MapToPostData(postDataA)
	} else {
		req, errT = http.NewRequest(reqTypeT, urlT, nil)
	}

	if errT != nil {
		return errT
	}

	customHeadersStrT := Trim(GetSwitchI(optsA, "-headers=", ""))

	if customHeadersStrT != "" {
		if strings.HasPrefix(customHeadersStrT, "HEX_") {
			customHeadersStrT = HexToStr(customHeadersStrT)
		}

		headerObjT, errT := FromJSON(customHeadersStrT)

		if errT == nil {
			headerMapT, ok := headerObjT.(map[string]interface{})

			if ok {
				for k, v := range headerMapT {
					s, ok := v.(string)
					if ok {
						req.Header.Add(k, s)
					}
				}
			}
		}
	}

	if IfSwitchExistsWholeI(optsA, "-verbose") {
		Pl("REQ: %v", req)
	}

	respT, errT = client.Do(req)

	if errT == nil {
		defer respT.Body.Close()
		if respT.StatusCode != 200 {
			if IfSwitchExistsWholeI(optsA, "-detail") {
				body, errT := io.ReadAll(respT.Body)

				if errT != nil {
					body = []byte(errT.Error())
				}
				Pl("response status: %v (%v) body: %v", respT.StatusCode, respT, string(body))
			}

			return fmt.Errorf("response status: %v", respT.StatusCode)
		} else {
			body, errT := io.ReadAll(respT.Body)

			if errT != nil {
				return errT
			}

			if IfSwitchExistsWholeI(optsA, "-detail") {
				Pl("response status: %v (%v) body: %v", respT.StatusCode, respT, string(body))
			}

			if IfSwitchExistsWholeI(optsA, "-bytes") {
				return body
			}

			if (encodingT == "") || (strings.ToLower(encodingT) == "utf-8") {
				if IfSwitchExistsWholeI(optsA, "-object") {
					var rs interface{}

					errT := jsoniter.Unmarshal(body, &rs)

					if errT != nil {
						return fmt.Errorf("failed to parse JSON object: %v", errT)
					}

					return rs
				}

				return string(body)
			} else {
				return ConvertToUTF8(body, encodingT)
			}
		}
	} else {
		return errT
	}
}

var GetWeb = TKX.GetWeb

func (pA *TK) GetWebResponse(urlA string, optsA ...interface{}) interface{} {
	// Pl("%#v", optsA)
	timeoutStrT := GetSwitchI(optsA, "-timeout=", "15")

	proxyT := GetSwitchI(optsA, "-proxy=", "")

	timeoutSecsT := time.Second * time.Duration(StrToInt(timeoutStrT, 15))

	client := &http.Client{
		//CheckRedirect: redirectPolicyFunc,
		Timeout: timeoutSecsT, // time.Second * timeoutSecsT,
	}

	if proxyT != "" {
		urli := url.URL{}
		urlproxy, errT := urli.Parse(proxyT)

		if errT != nil {
			return errT
		}

		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(urlproxy),
		}
	}

	var urlT string
	if !StartsWithIgnoreCase(urlA, "http") {
		urlT = "http://" + urlA
	} else {
		urlT = urlA
	}

	var respT *http.Response
	var errT error
	var req *http.Request

	reqTypeT := GetSwitchI(optsA, "-method=", "GET")

	var postDataT url.Values = nil

	postStrT := GetSwitchI(optsA, "-post=", "")

	var postBytesT []byte = nil

	if postStrT != "" {
		if strings.HasPrefix(postStrT, "HEX_") {
			postStrT = HexToStr(postStrT)
		}

		postObjT, errT := FromJSON(postStrT)

		if errT == nil {
			postMapT, ok := postObjT.(map[string]interface{})

			if ok {
				postDataT = MapToPostDataI(postMapT)
			}
		}
	}

	if postDataT == nil {
		for _, v := range optsA {
			rv, ok := v.(url.Values)

			if ok {
				postDataT = rv
				break
			}

			rv2, ok := v.(map[string]string)

			if ok {
				postDataT = MapToPostData(rv2)
				break
			}

			rv3, ok := v.(map[string]interface{})

			if ok {
				postDataT = MapToPostDataI(rv3)
				break
			}

			rv4, ok := v.([]byte)

			if ok {
				postBytesT = rv4
				break
			}
		}
	}

	postBodyT := GetSwitchI(optsA, "-postBody=", "")

	if strings.HasPrefix(postBodyT, "HEX_") {
		postBodyT = HexToStr(postBodyT)
	}

	if (postDataT != nil && len(postDataT) > 0) || postBodyT != "" || postBytesT != nil {
		if reqTypeT != "PUT" {
			reqTypeT = "POST"
		}

	}

	postFileT := GetSwitchI(optsA, "-postFile=", "")

	if reqTypeT == "POST" || reqTypeT == "PUT" {
		if postBytesT != nil {
			req, errT = http.NewRequest(reqTypeT, urlT, bytes.NewReader(postBytesT))

		} else if postBodyT != "" {
			req, errT = http.NewRequest(reqTypeT, urlT, strings.NewReader(postBodyT))

		} else if postFileT != "" {
			// file1, err := os.Open(postFileT)
			// if err != nil {
			// 	return GenerateErrorStringF("failed to load file content: %v", err)
			// }

			// defer file1.Close()

			bufT := LoadBytes(postFileT)

			// req, errT = http.NewRequest(reqTypeT, urlT, bufio.NewReader(file1))
			req, errT = http.NewRequest(reqTypeT, urlT, bytes.NewReader(bufT))

		} else {
			req, errT = http.NewRequest(reqTypeT, urlT, bytes.NewBufferString(postDataT.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
		}
		// req.PostForm = MapToPostData(postDataA)
	} else {
		req, errT = http.NewRequest(reqTypeT, urlT, nil)
	}

	if errT != nil {
		return errT
	}

	customHeadersStrT := Trim(GetSwitchI(optsA, "-headers=", ""))

	if customHeadersStrT != "" {
		if strings.HasPrefix(customHeadersStrT, "HEX_") {
			customHeadersStrT = HexToStr(customHeadersStrT)
		}

		headerObjT, errT := FromJSON(customHeadersStrT)

		if errT == nil {
			headerMapT, ok := headerObjT.(map[string]interface{})

			if ok {
				for k, v := range headerMapT {
					s, ok := v.(string)
					if ok {
						req.Header.Add(k, s)
					}
				}
			}
		}
	}

	if IfSwitchExistsWholeI(optsA, "-verbose") {
		Pl("REQ: %v", req)
	}

	respT, errT = client.Do(req)

	if errT != nil {
		return fmt.Errorf("failed to get response: %v", errT)
	}

	return respT
}

var GetWebResponse = TKX.GetWebResponse

func (pA *TK) GetWebResponseBody(urlA string, optsA ...interface{}) interface{} {
	// Pl("%#v", optsA)
	timeoutStrT := GetSwitchI(optsA, "-timeout=", "15")

	proxyT := GetSwitchI(optsA, "-proxy=", "")

	timeoutSecsT := time.Second * time.Duration(StrToInt(timeoutStrT, 15))

	client := &http.Client{
		//CheckRedirect: redirectPolicyFunc,
		Timeout: timeoutSecsT, // time.Second * timeoutSecsT,
	}

	if proxyT != "" {
		urli := url.URL{}
		urlproxy, errT := urli.Parse(proxyT)

		if errT != nil {
			return errT
		}

		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(urlproxy),
		}
	}

	var urlT string
	if !StartsWithIgnoreCase(urlA, "http") {
		urlT = "http://" + urlA
	} else {
		urlT = urlA
	}

	var respT *http.Response
	var errT error
	var req *http.Request

	reqTypeT := GetSwitchI(optsA, "-method=", "GET")

	var postDataT url.Values = nil

	postStrT := GetSwitchI(optsA, "-post=", "")

	var postBytesT []byte = nil

	if postStrT != "" {
		if strings.HasPrefix(postStrT, "HEX_") {
			postStrT = HexToStr(postStrT)
		}

		postObjT, errT := FromJSON(postStrT)

		if errT == nil {
			postMapT, ok := postObjT.(map[string]interface{})

			if ok {
				postDataT = MapToPostDataI(postMapT)
			}
		}
	}

	if postDataT == nil {
		for _, v := range optsA {
			rv, ok := v.(url.Values)

			if ok {
				postDataT = rv
				break
			}

			rv2, ok := v.(map[string]string)

			if ok {
				postDataT = MapToPostData(rv2)
				break
			}

			rv3, ok := v.(map[string]interface{})

			if ok {
				postDataT = MapToPostDataI(rv3)
				break
			}

			rv4, ok := v.([]byte)

			if ok {
				postBytesT = rv4
				break
			}
		}
	}

	postBodyT := GetSwitchI(optsA, "-postBody=", "")

	if strings.HasPrefix(postBodyT, "HEX_") {
		postBodyT = HexToStr(postBodyT)
	}

	if (postDataT != nil && len(postDataT) > 0) || postBodyT != "" || postBytesT != nil {
		if reqTypeT != "PUT" {
			reqTypeT = "POST"
		}

	}

	postFileT := GetSwitchI(optsA, "-postFile=", "")

	if reqTypeT == "POST" || reqTypeT == "PUT" {
		if postBytesT != nil {
			req, errT = http.NewRequest(reqTypeT, urlT, bytes.NewReader(postBytesT))

		} else if postBodyT != "" {
			req, errT = http.NewRequest(reqTypeT, urlT, strings.NewReader(postBodyT))

		} else if postFileT != "" {
			// file1, err := os.Open(postFileT)
			// if err != nil {
			// 	return GenerateErrorStringF("failed to load file content: %v", err)
			// }

			// defer file1.Close()

			bufT := LoadBytes(postFileT)

			// req, errT = http.NewRequest(reqTypeT, urlT, bufio.NewReader(file1))
			req, errT = http.NewRequest(reqTypeT, urlT, bytes.NewReader(bufT))

		} else {
			req, errT = http.NewRequest(reqTypeT, urlT, bytes.NewBufferString(postDataT.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
		}
		// req.PostForm = MapToPostData(postDataA)
	} else {
		req, errT = http.NewRequest(reqTypeT, urlT, nil)
	}

	if errT != nil {
		return errT
	}

	customHeadersStrT := Trim(GetSwitchI(optsA, "-headers=", ""))

	if customHeadersStrT != "" {
		if strings.HasPrefix(customHeadersStrT, "HEX_") {
			customHeadersStrT = HexToStr(customHeadersStrT)
		}

		headerObjT, errT := FromJSON(customHeadersStrT)

		if errT == nil {
			headerMapT, ok := headerObjT.(map[string]interface{})

			if ok {
				for k, v := range headerMapT {
					s, ok := v.(string)
					if ok {
						req.Header.Add(k, s)
					}
				}
			}
		}
	}

	if IfSwitchExistsWholeI(optsA, "-verbose") {
		Pl("REQ: %v", req)
	}

	respT, errT = client.Do(req)

	if errT != nil {
		return fmt.Errorf("failed to get response: %v", errT)
	}

	// if respT.StatusCode != 200 {
	// 	respT.Body.Close()
	// 	return fmt.Errorf("response status: %v", respT.StatusCode)
	// }

	if IfSwitchExistsWholeI(optsA, "-withLen") {
		return []interface{}{respT.Body, respT.ContentLength}
	}

	return respT.Body
}

var GetWebResponseBody = TKX.GetWebResponseBody

func (pA *TK) DownloadWebBytes(urlA string, postDataA map[string]string, customHeadersA map[string]string, optsA ...string) ([]byte, map[string]string, error) {
	timeoutStrT := GetSwitch(optsA, "-timeout=", "15")

	// encodingT := GetSwitch(optsA, "-encoding=", "")

	timeoutSecsT := time.Second * time.Duration(StrToInt(timeoutStrT, 15))

	client := &http.Client{
		//CheckRedirect: redirectPolicyFunc,
		Timeout: timeoutSecsT,
	}

	var urlT string
	if !StartsWithIgnoreCase(urlA, "http") {
		urlT = "http://" + urlA
	} else {
		urlT = urlA
	}

	var respT *http.Response
	var errT error
	var req *http.Request

	reqTypeT := GetSwitch(optsA, "-method=", "GET")

	if postDataA != nil && len(postDataA) > 0 && reqTypeT == "GET" {
		reqTypeT = "POST"
	}

	req, errT = http.NewRequest(reqTypeT, urlT, nil)

	if postDataA != nil {
		req.PostForm = MapToPostData(postDataA)
	}

	if customHeadersA != nil && len(customHeadersA) > 0 {
		for k, v := range customHeadersA {
			req.Header.Add(k, v)
		}
	}

	if IfSwitchExistsWhole(optsA, "-verbose") {
		Pl("REQ: %v", req)
	}

	respT, errT = client.Do(req)

	if errT == nil {
		defer respT.Body.Close()
		if respT.StatusCode != 200 {
			if IfSwitchExistsWhole(optsA, "-detail") {
				body, errT := io.ReadAll(respT.Body)

				if errT != nil {
					body = []byte(errT.Error())
				}
				Pl("response status: %v (%v) body: %v", respT.StatusCode, respT, string(body))
			}

			return nil, nil, Errf("response status: %v", respT.StatusCode)
		} else {
			body, errT := io.ReadAll(respT.Body)

			if errT != nil {
				return nil, nil, errT
			}

			mapT := make(map[string]string, 0)
			formT := respT.Header

			for k, v := range formT {
				mapT[k] = v[0]
			}

			return body, mapT, nil
		}
	} else {
		return nil, nil, errT
	}
}

var DownloadWebBytes = TKX.DownloadWebBytes

// DownloadPage download page with any encoding and convert to UTF-8
func (pA *TK) DownloadPage(urlA string, originalEncodingA string, postDataA url.Values, customHeaders string, timeoutSecsA time.Duration) string {
	client := &http.Client{
		Timeout: time.Second * timeoutSecsA,
	}

	var urlT string
	if !strings.HasPrefix(strings.ToLower(urlA), "http") {
		urlT = "http://" + urlA
	} else {
		urlT = urlA
	}

	var respT *http.Response
	var errT error
	var req *http.Request

	if Trim(customHeaders) != "" {
		if postDataA == nil {
			req, errT = http.NewRequest("GET", urlT, nil)
		} else {
			req, errT = http.NewRequest("POST", urlT, nil)
			req.PostForm = postDataA
		}

		headersT := SplitLines(customHeaders)

		for i := 0; i < len(headersT); i++ {
			lineT := headersT[i]
			aryT := strings.Split(lineT, ":")
			if len(aryT) < 2 {
				continue
			}
			req.Header.Add(aryT[0], Replace(aryT[1], "`", ":"))
		}

		respT, errT = client.Do(req)
	} else {
		if postDataA == nil {
			respT, errT = client.Get(urlT)
		} else {
			respT, errT = client.PostForm(urlT, postDataA)
		}
	}

	if errT == nil {
		defer respT.Body.Close()
		if respT.StatusCode != 200 {
			return GenerateErrorStringF("response status: %v", respT.StatusCode)
		} else {
			body, errT := io.ReadAll(respT.Body)

			if errT == nil {
				if (originalEncodingA == "") || (strings.ToLower(originalEncodingA) == "utf-8") {
					return string(body)
				} else {
					return ConvertToUTF8(body, originalEncodingA)
				}
			} else {
				return GenerateErrorString(errT.Error())
			}
		}
	} else {
		return GenerateErrorString(errT.Error())
	}

}

var DownloadPage = TKX.DownloadPage

// HttpRequest download page with any encoding and convert to UTF-8
func (pA *TK) HttpRequest(urlA string, methodA string, originalEncodingA string, postDataA url.Values, customHeaders string, timeoutSecsA time.Duration, optsA ...string) string {
	client := &http.Client{
		Timeout: time.Second * timeoutSecsA,
	}

	var urlT string
	if !strings.HasPrefix(strings.ToLower(urlA), "http") {
		urlT = "http://" + urlA
	} else {
		urlT = urlA
	}

	var respT *http.Response
	var errT error
	var req *http.Request

	if postDataA == nil {
		req, errT = http.NewRequest(methodA, urlT, nil)
		// req.PostForm = postDataA
	} else {
		req, errT = http.NewRequest(methodA, urlT, strings.NewReader(postDataA.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	headersT := SplitLines(customHeaders)

	for i := 0; i < len(headersT); i++ {
		lineT := headersT[i]
		aryT := strings.Split(lineT, ":")
		if len(aryT) < 2 {
			continue
		}
		req.Header.Add(aryT[0], Replace(aryT[1], "`", ":"))
	}

	if IfSwitchExistsWhole(optsA, "-verbose") {
		Pl("REQ: %v", req)
	}

	respT, errT = client.Do(req)

	if errT == nil {
		defer respT.Body.Close()

		if IfSwitchExistsWhole(optsA, "-verbose") {
			Pl("response status: %v (%v)", respT.StatusCode, respT)
		}

		if respT.StatusCode != 200 {
			return GenerateErrorStringF("response status: %v", respT.StatusCode)
		} else {
			body, errT := io.ReadAll(respT.Body)

			if errT == nil {
				if (originalEncodingA == "") || (strings.ToLower(originalEncodingA) == "utf-8") {
					return string(body)
				} else {
					return ConvertToUTF8(body, originalEncodingA)
				}
			} else {
				return GenerateErrorString(errT.Error())
			}
		}
	} else {
		return GenerateErrorString(errT.Error())
	}

}

var HttpRequest = TKX.HttpRequest

func (pA *TK) MapToPostData(postDataA map[string]string) url.Values {
	postDataT := url.Values{}

	for k, v := range postDataA {
		postDataT.Set(k, v)
	}

	return postDataT
}

var MapToPostData = TKX.MapToPostData

func (pA *TK) MapToPostDataI(postDataA map[string]interface{}) url.Values {
	postDataT := url.Values{}

	for k, v := range postDataA {
		s, ok := v.(string)

		if ok {
			postDataT.Set(k, s)
		}
	}

	return postDataT
}

var MapToPostDataI = TKX.MapToPostDataI

func (pA *TK) DownloadPageByMap(urlA string, originalEncodingA string, postDataA map[string]string, customHeaders string, timeoutSecsA time.Duration) string {
	if postDataA == nil {
		return DownloadPage(urlA, originalEncodingA, nil, customHeaders, timeoutSecsA)
	}

	postDataT := url.Values{}

	for k, v := range postDataA {

		postDataT.Set(k, v)

	}

	return DownloadPage(urlA, originalEncodingA, postDataT, customHeaders, timeoutSecsA)
}

var DownloadPageByMap = TKX.DownloadPageByMap

func (pA *TK) GetLastComponentOfUrl(urlA string) string {
	urlT, errT := url.Parse(urlA)
	if errT != nil {
		return GenerateErrorStringF(errT.Error())
	}

	splitRT := strings.Split(urlT.Path, "/")

	return splitRT[len(splitRT)-1]
}

var GetLastComponentOfUrl = TKX.GetLastComponentOfUrl

func (pA *TK) DownloadFile(urlA, dirA, fileNameA string, argsA ...string) string {

	ifRenameT := IfSwitchExistsWhole(argsA, "-rename")

	var urlT string
	var fileNameT string = fileNameA

	if (!StartsWithIgnoreCase(urlA, "http://")) && (!StartsWithIgnoreCase(urlA, "https://")) {
		urlT = "http://" + urlA
	} else {
		urlT = urlA
	}

	respT, errT := http.Get(urlT)
	if errT != nil {
		return GenerateErrorStringF("failed to get URL: %v", errT.Error())
	}

	if respT.StatusCode != 200 {
		return GenerateErrorStringF("failed to get URL: status code = %v", respT.StatusCode)
	}

	defer respT.Body.Close()

	if fileNameT == "" {
		fileNameT = GetLastComponentOfUrl(urlT)
	}

	if dirA != "" {
		fileNameT = filepath.Join(dirA, fileNameT)
	}

	if ifRenameT {
		fileNameT = GetAvailableFileName(fileNameT)
	}

	fileT, errT := os.Create(fileNameT)
	if errT != nil {
		return GenerateErrorStringF("failed to create file: %v", errT.Error())
	}
	defer fileT.Close()

	// if respT.ContentLength == -1 {
	// 	tmpBuf, _ := io.ReadAll(respT.Body)
	// 	return GenerateErrorStringF("failed to get http response content length: %v\n%#v", string(tmpBuf), respT)
	// }

	bufT := make([]byte, 1000000)

	// var writeCountT int = 0

	for {
		n, errT := respT.Body.Read(bufT)

		if n == 0 && errT.Error() == "EOF" {
			break
		}

		if (errT != nil) && (errT.Error() != "EOF") {
			return GenerateErrorStringF("failed to download: %v", errT.Error())
		}

		_, errT = fileT.Write(bufT[:n])
		if errT != nil {
			return GenerateErrorStringF("failed to writer file: %v", errT)
		}

		// writeCountT += countT
	}

	// valid download exe
	// fi, errT := fileT.Stat()
	// if errT == nil {
	// 	if fi.Size() != respT.ContentLength {
	// 		return GenerateErrorStringF("Downloaded file size error")
	// 	}
	// }
	return fileNameT
}

var DownloadFile = TKX.DownloadFile

func (pA *TK) DownloadFileWithProgress(urlA, dirA, fileNameA string, funcA func(interface{}) interface{}, argsA ...string) string {

	ifRenameT := IfSwitchExistsWhole(argsA, "-rename")

	var urlT string
	var fileNameT string = fileNameA

	if (!StartsWithIgnoreCase(urlA, "http://")) && (!StartsWithIgnoreCase(urlA, "https://")) {
		urlT = "http://" + urlA
	} else {
		urlT = urlA
	}

	respT, errT := http.Get(urlT)
	if errT != nil {
		return GenerateErrorStringF("failed to get URL: %v", errT.Error())
	}

	if respT.StatusCode != 200 {
		return GenerateErrorStringF("failed to get URL: status code = %v", respT.StatusCode)
	}

	defer respT.Body.Close()

	if fileNameT == "" {
		fileNameT = GetLastComponentOfUrl(urlT)
	}

	if dirA != "" {
		fileNameT = filepath.Join(dirA, fileNameT)
	}

	if ifRenameT {
		fileNameT = GetAvailableFileName(fileNameT)
	}

	fileT, errT := os.Create(fileNameT)
	if errT != nil {
		return GenerateErrorStringF("failed to create file: %v", errT.Error())
	}
	defer fileT.Close()

	// if respT.ContentLength == -1 {
	// 	tmpBuf, _ := io.ReadAll(respT.Body)
	// 	return GenerateErrorStringF("failed to get http response content length: %v\n%#v", string(tmpBuf), respT)
	// }

	lengthT := respT.ContentLength

	var multiWriterT io.Writer

	if funcA != nil {
		var countingWriterT io.Writer

		ifPercentT := IfSwitchExistsWhole(argsA, "-percent")

		if ifPercentT {
			countingWriterT = NewCountingWriter(funcA, "-percent", lengthT)
		} else {
			countingWriterT = NewCountingWriter(funcA)
		}

		multiWriterT = io.MultiWriter(fileT, countingWriterT)
	} else {
		multiWriterT = fileT
	}

	bufT := make([]byte, 1000000)

	// var writeCountT int = 0

	for {
		n, errT := respT.Body.Read(bufT)

		if n == 0 && errT.Error() == "EOF" {
			break
		}

		if (errT != nil) && (errT.Error() != "EOF") {
			return GenerateErrorStringF("failed to download: %v", errT.Error())
		}

		_, errT = multiWriterT.Write(bufT[:n])
		if errT != nil {
			return GenerateErrorStringF("failed to writer file: %v", errT)
		}

		// writeCountT += countT
	}

	// valid download exe
	// fi, errT := fileT.Stat()
	// if errT == nil {
	// 	if fi.Size() != respT.ContentLength {
	// 		return GenerateErrorStringF("Downloaded file size error")
	// 	}
	// }
	return fileNameT
}

var DownloadFileWithProgress = TKX.DownloadFileWithProgress

func (pA *TK) DownloadBytes(urlA string) ([]byte, error) {

	var urlT string

	if !StartsWithIgnoreCase(urlA, "http://") {
		urlT = "http://" + urlA
	} else {
		urlT = urlA
	}

	respT, errT := http.Get(urlT)
	if errT != nil {
		return nil, Errf("failed to get URL: %v", errT.Error())
	}

	if respT.StatusCode != 200 {
		return nil, Errf("failed to get URL: status code = %v", respT.StatusCode)
	}

	defer respT.Body.Close()

	bufT, errT := io.ReadAll(respT.Body)

	if errT != nil {
		return nil, Errf("failed to get http response body: %v", errT)
	}

	return bufT, nil
}

var DownloadBytes = TKX.DownloadBytes

func (pA *TK) DownloadBytesWithProgress(urlA string, funcA func(interface{}) interface{}) ([]byte, error) {

	var urlT string

	if (!StartsWithIgnoreCase(urlA, "http://")) && (!StartsWithIgnoreCase(urlA, "https://")) {
		urlT = "http://" + urlA
	} else {
		urlT = urlA
	}

	respT, errT := http.Get(urlT)
	if errT != nil {
		return nil, Errf("failed to get URL: %v", errT.Error())
	}

	if respT.StatusCode != 200 {
		return nil, Errf("failed to get URL: status code = %v", respT.StatusCode)
	}

	defer respT.Body.Close()

	lengthT := 0

	tmps := strings.TrimSpace(respT.Header.Get("Content-Length"))

	if tmps != "" {
		lengthT = StrToIntWithDefaultValue(tmps, 0)
	}

	countingWriterT := NewCountingWriter(funcA, "-percent", fmt.Sprintf("-total=%v", lengthT))

	teeReaderT := io.TeeReader(respT.Body, countingWriterT)

	bufT, errT := io.ReadAll(teeReaderT)

	if errT != nil {
		return nil, Errf("failed to get http response body: %v", errT)
	}

	return bufT, nil
}

var DownloadBytesWithProgress = TKX.DownloadBytesWithProgress

// PostRequest : another POST request sender
func (pA *TK) PostRequest(urlA, reqBodyA string, timeoutSecsA time.Duration) (string, error) {

	req, err := http.NewRequest("POST", urlA, strings.NewReader(reqBodyA))

	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json; encoding=utf-8")

	client := &http.Client{
		//CheckRedirect: redirectPolicyFunc,
		Timeout: time.Second * timeoutSecsA,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

var PostRequest = TKX.PostRequest

// PostRequestX : Post Request with custom headers
func (pA *TK) PostRequestX(urlA, reqBodyA string, customHeadersA string, timeoutSecsA time.Duration, optsA ...string) (string, error) {

	methodT := "POST"
	if IfSwitchExistsWhole(optsA, "-forceGet") {
		methodT = "GET"
	}

	req, err := http.NewRequest(methodT, urlA, strings.NewReader(reqBodyA))

	if err != nil {
		return "", err
	}

	headersT := SplitLines(customHeadersA)

	for i := 0; i < len(headersT); i++ {
		lineT := headersT[i]
		aryT := strings.SplitN(lineT, ":", 2)
		if len(aryT) < 2 {
			continue
		}
		req.Header.Add(aryT[0], Replace(aryT[1], "TX_COLONS_XT", ":"))
		// Pl("%s=%s", aryT[0], aryT[1])
	}

	if IfSwitchExistsWhole(optsA, "-verbose") {
		Pl("POST data: %v", reqBodyA)
		Pl("REQ: %v", req)
	}

	client := &http.Client{
		//CheckRedirect: redirectPolicyFunc,
		Timeout: time.Second * timeoutSecsA,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if IfSwitchExistsWhole(optsA, "-detail") {
		Pl("response status: %v (%v)", resp.StatusCode, resp)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

var PostRequestX = TKX.PostRequestX

func (pA *TK) PrepareMultiPartFileFromBytes(fieldNameA, fileNameA string, bytesA []byte) (string, []byte, error) {
	var b bytes.Buffer
    w := multipart.NewWriter(&b)
	
	r := bytes.NewReader([]byte(bytesA))
	
	var fw io.Writer
	
	var err error

	if fw, err = w.CreateFormFile(fieldNameA, fileNameA); err != nil {
		return "", nil, fmt.Errorf("failed to create field: (%v) %v", fieldNameA, err)
	}
	
	if _, err = io.Copy(fw, r); err != nil {
		return "", nil, fmt.Errorf("failed to set field: (%v) %v", fieldNameA, err)
	}
	
    w.Close()
	
	return w.FormDataContentType(), b.Bytes(), nil
}

var PrepareMultiPartFileFromBytes = TKX.PrepareMultiPartFileFromBytes

// PutRequestX : Put Request with custom headers
func (pA *TK) PutRequestX(urlA, reqBodyA string, customHeadersA string, timeoutSecsA time.Duration, optsA ...string) (string, error) {

	req, err := http.NewRequest("PUT", urlA, strings.NewReader(reqBodyA))

	if err != nil {
		return "", err
	}

	headersT := SplitLines(customHeadersA)

	for i := 0; i < len(headersT); i++ {
		lineT := headersT[i]
		aryT := strings.Split(lineT, ":")
		if len(aryT) < 2 {
			continue
		}
		req.Header.Add(aryT[0], Replace(aryT[1], "`", ":"))
		// Pl("%s=%s", aryT[0], aryT[1])
	}

	if IfSwitchExistsWhole(optsA, "-verbose") {
		Pl("PUT data: %v", reqBodyA)
	}

	if IfSwitchExistsWhole(optsA, "-verbose") {
		Pl("REQ: %v", req)
	}

	client := &http.Client{
		//CheckRedirect: redirectPolicyFunc,
		Timeout: time.Second * timeoutSecsA,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if IfSwitchExistsWhole(optsA, "-detail") {
		Pl("response status: %v (%v)", resp.StatusCode, resp)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

var PutRequestX = TKX.PutRequestX

// PostRequestBytesX : PostRequest with custom headers
func (pA *TK) PostRequestBytesX(urlA string, reqBodyA []byte, customHeadersA string, timeoutSecsA time.Duration) ([]byte, error) {

	req, err := http.NewRequest("POST", urlA, bytes.NewReader(reqBodyA))

	if err != nil {
		return nil, err
	}

	headersT := SplitLines(customHeadersA)

	for i := 0; i < len(headersT); i++ {
		lineT := headersT[i]
		aryT := strings.Split(lineT, ":")

		if len(aryT) < 2 {
			continue
		}

		req.Header.Add(aryT[0], Replace(aryT[1], "`", ":"))
		// Pl("%s=%s", aryT[0], aryT[1])
	}

	client := &http.Client{
		//CheckRedirect: redirectPolicyFunc,
		Timeout: time.Second * timeoutSecsA,
	}

	resp, err := client.Do(req)
	if err != nil {
		return []byte(Spr("%#v", resp)), err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return body, err
	}

	return body, nil
}

var PostRequestBytesX = TKX.PostRequestBytesX

// RequestX : Network(http) Request with custom headers
func (pA *TK) RequestX(urlA, methodA, reqBodyA string, customHeadersA string, timeoutSecsA time.Duration, optsA ...string) (string, error) {

	// methodA: GET, POST, PUT, etc
	req, err := http.NewRequest(methodA, urlA, strings.NewReader(reqBodyA))

	if err != nil {
		return "", err
	}

	headersT := SplitLines(customHeadersA)

	for i := 0; i < len(headersT); i++ {
		lineT := headersT[i]
		aryT := strings.Split(lineT, ":")
		if len(aryT) < 2 {
			continue
		}
		req.Header.Add(aryT[0], Replace(aryT[1], "`", ":"))
		// Pl("%s=%s", aryT[0], aryT[1])
	}

	if IfSwitchExistsWhole(optsA, "-verbose") {
		Pl("REQUEST data: %v", reqBodyA)
	}

	if IfSwitchExistsWhole(optsA, "-verbose") {
		Pl("REQ: %v", req)
	}

	client := &http.Client{
		//CheckRedirect: redirectPolicyFunc,
		Timeout: time.Second * timeoutSecsA,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if IfSwitchExistsWhole(optsA, "-detail") {
		Pl("response status: %v (%v)", resp.StatusCode, resp)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

var RequestX = TKX.RequestX

// PostRequestBytesX : PostRequest with custom headers
func (pA *TK) PostRequestBytesWithMSSHeaderX(urlA string, reqBodyA []byte, customHeadersA map[string]string, timeoutSecsA time.Duration) ([]byte, error) {

	req, err := http.NewRequest("POST", urlA, bytes.NewReader(reqBodyA))

	if err != nil {
		return nil, err
	}

	if customHeadersA != nil {
		for k, v := range customHeadersA {
			req.Header.Add(k, v)
		}
	}

	client := &http.Client{
		//CheckRedirect: redirectPolicyFunc,
		Timeout: time.Second * timeoutSecsA,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

var PostRequestBytesWithMSSHeaderX = TKX.PostRequestBytesWithMSSHeaderX

// PostRequestBytesWithCookieX : PostRequest with custom headers
func (pA *TK) PostRequestBytesWithCookieX(urlA string, reqBodyA []byte, customHeadersA string, jarA *cookiejar.Jar, timeoutSecsA time.Duration) ([]byte, *cookiejar.Jar, error) {

	req, err := http.NewRequest("POST", urlA, bytes.NewReader(reqBodyA))

	if err != nil {
		return nil, nil, err
	}

	headersT := SplitLines(customHeadersA)

	for i := 0; i < len(headersT); i++ {
		lineT := headersT[i]
		aryT := strings.Split(lineT, ":")

		if len(aryT) < 2 {
			continue
		}

		req.Header.Add(aryT[0], Replace(aryT[1], "`", ":"))
		// Pl("%s=%s", aryT[0], aryT[1])
	}

	jarT := jarA

	if jarT == nil {
		jarT, _ = cookiejar.New(nil)
	}

	client := &http.Client{
		//CheckRedirect: redirectPolicyFunc,
		Timeout: time.Second * timeoutSecsA,
		Jar:     jarT,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	// Pl("%#v", client.Jar)

	// cookiesT := resp.Cookies()

	// for i, v := range cookiesT {
	// 	Pl("%v, %#v", i, v)
	// }

	return body, jarT, nil
}

var PostRequestBytesWithCookieX = TKX.PostRequestBytesWithCookieX

func (pA *TK) GetFormValueWithDefaultValue(reqA *http.Request, keyA string, defaultA ...string) string {
	defaultT := ""

	if len(defaultA) > 0 {
		defaultT = defaultA[0]
	}

	valueT, ok := reqA.Form[keyA]
	if ok {
		return valueT[0]
	} else {
		return defaultT
	}
}

var GetFormValueWithDefaultValue = TKX.GetFormValueWithDefaultValue

func (pA *TK) IfFormValueExists(reqA *http.Request, keyA string) bool {
	_, ok := reqA.Form[keyA]
	if ok {
		return true
	} else {
		return false
	}
}

var IfFormValueExists = TKX.IfFormValueExists

func (pA *TK) GenerateJSONPResponse(statusA string, valueA string, reqA *http.Request, argsA ...string) string {
	_, valueOnlyT := reqA.Form["valueonly"]

	if valueOnlyT {
		if _, withErrorT := reqA.Form["witherror"]; withErrorT {
			if statusA != "success" {
				return GenerateErrorString(valueA)
			}
		}

		return valueA
	} else {
		mT := make(map[string]string)
		mT["Status"] = statusA
		mT["Value"] = valueA

		if argsA != nil && len(argsA) > 0 {
			lenT := len(argsA) / 2

			for i := 0; i < lenT; i++ {
				mT[argsA[i*2]] = argsA[i*2+1]
			}
		}

		returnValue, ifReturnValue := reqA.Form["returnvalue"]

		if ifReturnValue {
			mT["ReturnValue"] = returnValue[0]
		}

		name, ok := reqA.Form["callback"]
		if ok {
			if valueOnlyT {
				return fmt.Sprintf("%v(%v);", name[0], valueA)
			} else {
				return fmt.Sprintf("%v(%v);", name[0], ObjectToJSON(mT))
			}
		}

		return fmt.Sprintf("%v", ObjectToJSON(mT))
	}
}

var GenerateJSONPResponse = TKX.GenerateJSONPResponse

func (pA *TK) GenerateJSONPResponseMix(statusA string, valueA string, reqA *http.Request, mapA map[string]string) string {
	_, valueOnlyT := reqA.Form["valueonly"]

	if valueOnlyT {
		if _, withErrorT := reqA.Form["witherror"]; withErrorT {
			if statusA != "success" {
				return GenerateErrorString(valueA)
			}
		}

		return valueA
	} else {
		mT := make(map[string]string)
		mT["Status"] = statusA
		mT["Value"] = valueA

		if mapA != nil {
			for k, v := range mapA {
				mT[k] = v
			}
		}

		returnValue, ifReturnValue := reqA.Form["returnvalue"]

		if ifReturnValue {
			mT["ReturnValue"] = returnValue[0]
		}

		name, ok := reqA.Form["callback"]
		if ok {
			if valueOnlyT {
				return fmt.Sprintf("%v(%v);", name[0], valueA)
			} else {
				return fmt.Sprintf("%v(%v);", name[0], ObjectToJSON(mT))
			}
		}

		return fmt.Sprintf("%v", ObjectToJSON(mT))
	}
}

var GenerateJSONPResponseMix = TKX.GenerateJSONPResponseMix

func (pA *TK) GenerateJSONPResponseWithMore(statusA string, valueA string, reqA *http.Request, argsA ...string) string {
	var valueOnlyT bool = false

	if reqA != nil {
		_, valueOnlyT = reqA.Form["valueonly"]
	}

	if valueOnlyT {
		if _, withErrorT := reqA.Form["witherror"]; withErrorT {
			if statusA != "success" {
				return GenerateErrorString(valueA)
			}
		}

		return valueA
	} else {
		mT := make(map[string]string)
		mT["Status"] = statusA
		mT["Value"] = valueA

		if argsA != nil && len(argsA) > 0 {
			lenT := len(argsA) / 2

			for i := 0; i < lenT; i++ {
				mT[argsA[i*2]] = argsA[i*2+1]
			}
		}

		var returnValue []string
		var ifReturnValue bool = false

		if reqA != nil {
			returnValue, ifReturnValue = reqA.Form["returnvalue"]
		}

		if ifReturnValue {
			mT["ReturnValue"] = returnValue[0]
		}

		var name []string
		var ok bool = false

		if reqA != nil {
			name, ok = reqA.Form["callback"]
		}

		if ok {
			if valueOnlyT {
				return fmt.Sprintf("%v(%v);", name[0], valueA)
			} else {
				return fmt.Sprintf("%v(%v);", name[0], ObjectToJSON(mT))
			}
		}

		return fmt.Sprintf("%v", ObjectToJSON(mT))
	}
}

var GenerateJSONPResponseWithMore = TKX.GenerateJSONPResponseWithMore

func (pA *TK) GenerateJSONPResponseWithObject(statusA string, valueA string, objectA string, reqA *http.Request) string {
	_, valueOnlyT := reqA.Form["valueonly"]
	_, objectOnlyT := reqA.Form["objectonly"]

	if objectOnlyT {
		return objectA
	} else if valueOnlyT {
		return valueA
	} else {
		mT := make(map[string]string)
		mT["Status"] = statusA
		mT["Value"] = valueA
		mT["Object"] = objectA

		returnValue, ifReturnValue := reqA.Form["returnvalue"]

		if ifReturnValue {
			mT["ReturnValue"] = returnValue[0]
		}

		name, ok := reqA.Form["callback"]
		if ok {
			if valueOnlyT {
				return fmt.Sprintf("%v(%v);", name[0], valueA)
			} else {
				return fmt.Sprintf("%v(%v);", name[0], ObjectToJSON(mT))
			}
		}

		return fmt.Sprintf("%v", ObjectToJSON(mT))
	}
}

var GenerateJSONPResponseWithObject = TKX.GenerateJSONPResponseWithObject

func (pA *TK) GenerateJSONPResponseWith2Object(statusA string, valueA string, objectA string, object2A string, reqA *http.Request) string {
	_, valueOnlyT := reqA.Form["valueonly"]
	_, objectOnlyT := reqA.Form["objectonly"]

	if objectOnlyT {
		return objectA
	} else if valueOnlyT {
		return valueA
	} else {
		mT := make(map[string]string)
		mT["Status"] = statusA
		mT["Value"] = valueA
		mT["Object"] = objectA
		mT["Object2"] = object2A

		returnValue, ifReturnValue := reqA.Form["returnvalue"]

		if ifReturnValue {
			mT["ReturnValue"] = returnValue[0]
		}

		name, ok := reqA.Form["callback"]
		if ok {
			if valueOnlyT {
				return fmt.Sprintf("%v(%v);", name[0], valueA)
			} else {
				return fmt.Sprintf("%v(%v);", name[0], ObjectToJSON(mT))
			}
		}

		return fmt.Sprintf("%v", ObjectToJSON(mT))
	}
}

var GenerateJSONPResponseWith2Object = TKX.GenerateJSONPResponseWith2Object

func (pA *TK) GenerateJSONPResponseWith3Object(statusA string, valueA string, objectA string, object2A string, object3A string, reqA *http.Request) string {
	_, valueOnlyT := reqA.Form["valueonly"]
	_, objectOnlyT := reqA.Form["objectonly"]

	if objectOnlyT {
		return objectA
	} else if valueOnlyT {
		return valueA
	} else {
		mT := make(map[string]string)
		mT["Status"] = statusA
		mT["Value"] = valueA
		mT["Object"] = objectA
		mT["Object2"] = object2A
		mT["Object3"] = object3A

		returnValue, ifReturnValue := reqA.Form["returnvalue"]

		if ifReturnValue {
			mT["ReturnValue"] = returnValue[0]
		}

		name, ok := reqA.Form["callback"]
		if ok {
			if valueOnlyT {
				return fmt.Sprintf("%v(%v);", name[0], valueA)
			} else {
				return fmt.Sprintf("%v(%v);", name[0], ObjectToJSON(mT))
			}
		}

		return fmt.Sprintf("%v", ObjectToJSON(mT))
	}
}

var GenerateJSONPResponseWith3Object = TKX.GenerateJSONPResponseWith3Object

func (pA *TK) GetSuccessValue(strA string) string {
	rv := JSONToMapStringString(strA)
	if rv == nil {
		return GenerateErrorString("invalid json data")
	}

	statusT, ok := rv["Status"]
	if !ok {
		return GenerateErrorString("invalid json data")
	}

	if statusT != "success" {
		return GenerateErrorString("status not success")
	}

	valueT, ok := rv["Value"]
	if !ok {
		return GenerateErrorString("invalid json data")
	}

	return valueT
}

var GetSuccessValue = TKX.GetSuccessValue

func (pA *TK) JSONResponseToHTML(jsonA string) string {
	rv := JSONToMapStringString(jsonA)
	if rv == nil {
		return ""
	}

	statusT := rv["Status"]

	colorT := "auto"

	if statusT == "success" {
		statusT = "操作成功"
	} else if statusT == "fail" {
		statusT = "操作失败"
		colorT = "red"
	} else {
		statusT = "未知操作状态"
		colorT = "orange"
	}

	valueT := rv["Value"]

	tmplT := `<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<meta http-equiv="content-type" content="text/html; charset=UTF-8" />
<meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no">
</head>
<body>
<div style="display:table; width: 100%;">
	<div style="display: table-cell; vertical-align: middle; text-align:center; height: 500px;">
		<div style="margin-left:auto; margin-right:auto; font-size: 1.5em; font-weight: bold; color: TX_COLOR_XT;"><span>TX_STATUS_XT</span></div>
		<div style="margin-left:auto; margin-right:auto; margin-top: 3.0em;"><span>TX_VALUE_XT</span></div>
	</div>
</div>
</body>
</html>

	`

	tmplT = Replace(tmplT, "TX_STATUS_XT", EncodeHTML(statusT))
	tmplT = Replace(tmplT, "TX_VALUE_XT", EncodeHTML(valueT))
	tmplT = Replace(tmplT, "TX_COLOR_XT", colorT)

	return tmplT
}

var JSONResponseToHTML = TKX.JSONResponseToHTML

// 数学相关 math related

// max integer
const MAX_INT = int(^uint(0) >> 1)

// Max 取多个数值的最大值，数据类型以第一个参数为准
func (pA *TK) Max(nA ...interface{}) interface{} {
	lenT := len(nA)
	if lenT < 1 {
		return fmt.Errorf("no values")
	}

	switch nv := nA[0].(type) {
	case bool:
		if nv {
			return true
		}

		for i := 1; i < lenT; i++ {
			if ToBool(nA[i]) {
				return true
			}
		}

		return false
	case byte:
		maxT := nv

		for i := 1; i < lenT; i++ {
			vT := ToByte(nA[i])

			if vT > maxT {
				maxT = vT
			}
		}

		return maxT
	case rune:
		maxT := nv

		for i := 1; i < lenT; i++ {
			vT := ToRune(nA[i])

			if vT > maxT {
				maxT = vT
			}
		}

		return maxT
	case int:
		maxT := nv

		for i := 1; i < lenT; i++ {
			vT := ToInt(nA[i])

			if vT > maxT {
				maxT = vT
			}
		}

		return maxT
	case int64:
		maxT := nv

		for i := 1; i < lenT; i++ {
			vT := int64(ToInt(nA[i]))

			if vT > maxT {
				maxT = vT
			}
		}

		return maxT
	case float32:
		maxT := nv

		for i := 1; i < lenT; i++ {
			vT := float32(ToFloat(nA[i]))

			if vT > maxT {
				maxT = vT
			}
		}

		return maxT
	case float64:
		maxT := nv

		for i := 1; i < lenT; i++ {
			vT := ToFloat(nA[i])

			if vT > maxT {
				maxT = vT
			}
		}

		return maxT
	case string:
		maxT := nv

		for i := 1; i < lenT; i++ {
			vT := ToStr(nA[i])

			if vT > maxT {
				maxT = vT
			}
		}

		return maxT
	case UndefinedStruct:
		for i := 1; i < lenT; i++ {
			if ToBool(nA[i]) {
				return true
			}
		}

		return false
	default:
		return fmt.Errorf("unsupported type: %T(%v)", nA[0], nA[0])
	}
}

var Max = TKX.Max

// Min 取多个数值的最小值，数据类型以第一个参数为准
func (pA *TK) Min(nA ...interface{}) interface{} {
	lenT := len(nA)
	if lenT < 1 {
		return fmt.Errorf("no values")
	}

	switch nv := nA[0].(type) {
	case bool:
		if !nv {
			return false
		}

		for i := 1; i < lenT; i++ {
			if !ToBool(nA[i]) {
				return false
			}
		}

		return true
	case byte:
		minT := nv

		for i := 1; i < lenT; i++ {
			vT := ToByte(nA[i])

			if vT < minT {
				minT = vT
			}
		}

		return minT
	case rune:
		minT := nv

		for i := 1; i < lenT; i++ {
			vT := ToRune(nA[i])

			if vT < minT {
				minT = vT
			}
		}

		return minT
	case int:
		minT := nv

		for i := 1; i < lenT; i++ {
			vT := ToInt(nA[i])

			if vT < minT {
				minT = vT
			}
		}

		return minT
	case int64:
		minT := nv

		for i := 1; i < lenT; i++ {
			vT := int64(ToInt(nA[i]))

			if vT < minT {
				minT = vT
			}
		}

		return minT
	case float32:
		minT := nv

		for i := 1; i < lenT; i++ {
			vT := float32(ToFloat(nA[i]))

			if vT < minT {
				minT = vT
			}
		}

		return minT
	case float64:
		minT := nv

		for i := 1; i < lenT; i++ {
			vT := ToFloat(nA[i])

			if vT < minT {
				minT = vT
			}
		}

		return minT
	case string:
		minT := nv

		for i := 1; i < lenT; i++ {
			vT := ToStr(nA[i])

			if vT < minT {
				minT = vT
			}
		}

		return minT
	case UndefinedStruct:
		return false
	default:
		return fmt.Errorf("unsupported type: %T(%v)", nA[0], nA[0])
	}
}

var Min = TKX.Min

func (pA *TK) Ceil(nA interface{}) interface{} {
	switch nv := nA.(type) {
	case bool:
		if nv {
			return 1
		}
		return 0
	case byte:
		return int(nv)
	case rune:
		return int(nv)
	case int:
		return nv
	case int64:
		return int(nv)
	case float32:
		return int(math.Ceil(float64(nv)))
	case float64:
		return int(math.Ceil(nv))
	case UndefinedStruct:
		return 0
	default:
		return fmt.Errorf("unsupported type: %T(%v)", nA, nA)
	}
}

var Ceil = TKX.Ceil

func (pA *TK) Floor(nA interface{}) interface{} {
	switch nv := nA.(type) {
	case bool:
		if nv {
			return 1
		}
		return 0
	case byte:
		return int(nv)
	case rune:
		return int(nv)
	case int:
		return nv
	case int64:
		return int(nv)
	case float32:
		return int(math.Floor(float64(nv)))
	case float64:
		return int(math.Floor(nv))
	case UndefinedStruct:
		return 0
	default:
		return fmt.Errorf("unsupported type: %T(%v)", nA, nA)
	}
}

var Floor = TKX.Floor

func (pA *TK) Round(nA interface{}) interface{} {
	switch nv := nA.(type) {
	case bool:
		if nv {
			return 1
		}
		return 0
	case byte:
		return int(nv)
	case rune:
		return int(nv)
	case int:
		return nv
	case int64:
		return int(nv)
	case float32:
		return int(math.Round(float64(nv)))
	case float64:
		return int(math.Round(nv))
	case UndefinedStruct:
		return 0
	default:
		return fmt.Errorf("unsupported type: %T(%v)", nA, nA)
	}
}

var Round = TKX.Round

func (pA *TK) GetNegativeResult(nA interface{}) interface{} {
	switch nv := nA.(type) {
	case bool:
		return !nv
	case byte:
		return -nv
	case rune:
		return -nv
	case int:
		return -nv
	case int64:
		return -nv
	case float32:
		return -nv
	case float64:
		return -nv
	case UndefinedStruct:
		return true
	default:
		return fmt.Errorf("unknown type(1-): %T(%v)", nA, nA)
	}
}

var GetNegativeResult = TKX.GetNegativeResult

func (pA *TK) GetLogicalNotResult(nA interface{}) interface{} {
	if nA == nil {
		return true
	}

	switch nv := nA.(type) {
	case bool:
		return !nv
	case UndefinedStruct:
		return true
	}

	return false
}

var GetLogicalNotResult = TKX.GetLogicalNotResult

func (pA *TK) GetBitNotResult(nA interface{}) interface{} {
	// ^x    bitwise complement    is m ^ x  with m = "all bits set to 1" for unsigned x
	// and  m = -1 for signed x
	// so ^x is bitwise not of x

	if nA == nil {
		return true
	}

	var v3 interface{}

	switch nv := nA.(type) {
	case bool:
		v3 = !nv
	case byte:
		v3 = ^nv
	case rune:
		v3 = ^nv
	case int:
		v3 = ^nv
	case []byte:
		bufT := make([]byte, len(nv))
		for i := 0; i < len(nv); i++ {
			bufT[i] = ^(nv[i])
		}

		v3 = bufT
	case string:
		buf, err := hex.DecodeString(nv)
		if err != nil {
			return fmt.Errorf("hex convertion failed")
		}

		for i := 0; i < len(buf); i++ {
			buf[i] = ^(buf[i])
		}

		v3 = hex.EncodeToString(buf)
	default:
		return fmt.Errorf("unknown type(*): %T(%v)", nA, nA)
	}

	return v3
}

var GetBitNotResult = TKX.GetBitNotResult

func (pA *TK) GetMultiplyResult(n1A interface{}, n2A interface{}) (result interface{}) {
	defer func() {
		r := recover()

		if r != nil {
			result = fmt.Errorf("failed: %v(%v, %v)", r, n1A, n2A)
			return
		}
	}()

	switch nv := n1A.(type) {
	case byte:
		return nv * n2A.(byte)
	case rune:
		return nv * n2A.(rune)
	case int:
		return nv * n2A.(int)
	case int64:
		return nv * n2A.(int64)
	case float32:
		return nv * n2A.(float32)
	case float64:
		return nv * n2A.(float64)
	case string:
		return strings.Repeat(nv, ToInt(n2A))
	default:
		return fmt.Errorf("unknown type(*): %T(%v)", n1A, n1A)
	}
}

var GetMultiplyResult = TKX.GetMultiplyResult

func (pA *TK) GetDivResult(n1A interface{}, n2A interface{}) (result interface{}) {
	defer func() {
		r := recover()

		if r != nil {
			result = fmt.Errorf("failed: %v(%v, %v)", r, n1A, n2A)
			return
		}
	}()

	switch nv := n1A.(type) {
	case byte:
		return nv / n2A.(byte)
	case rune:
		return nv / n2A.(rune)
	case int:
		return nv / n2A.(int)
	case int64:
		return nv / n2A.(int64)
	case float32:
		return nv / n2A.(float32)
	case float64:
		return nv / n2A.(float64)
	default:
		return fmt.Errorf("unknown type(/): %T(%v)", n1A, n1A)
	}
}

var GetDivResult = TKX.GetDivResult

func (pA *TK) GetModResult(n1A interface{}, n2A interface{}) (result interface{}) {
	defer func() {
		r := recover()

		if r != nil {
			result = fmt.Errorf("failed: %v(%v, %v)", r, n1A, n2A)
			return
		}
	}()

	switch nv := n1A.(type) {
	case byte:
		return nv % n2A.(byte)
	case rune:
		return nv % n2A.(rune)
	case int:
		return nv % n2A.(int)
	case int64:
		return nv % n2A.(int64)
	// case float32:
	// 	return nv % n2A.(float32)
	// case float64:
	// 	return nv % n2A.(float64)
	default:
		return fmt.Errorf("unknown type(%%): %T(%v)", n1A, n1A)
	}
}

var GetModResult = TKX.GetModResult

func (pA *TK) GetMinusResult(n1A interface{}, n2A interface{}) (result interface{}) {
	defer func() {
		r := recover()

		if r != nil {
			result = fmt.Errorf("failed: %v(%v, %v)", r, n1A, n2A)
			return
		}
	}()

	switch nv := n1A.(type) {
	case byte:
		return nv - n2A.(byte)
	case rune:
		return nv - n2A.(rune)
	case int:
		return nv - n2A.(int)
	case int64:
		return nv - n2A.(int64)
	case float32:
		return nv - n2A.(float32)
	case float64:
		return nv - n2A.(float64)
	case time.Time:
		rsT := ToTime(n2A)

		if IsError(rsT) {
			t2 := ToInt(n2A, MAX_INT)

			if t2 == MAX_INT {
				return fmt.Errorf("time conversion failed: %T -> %T", n1A, n2A)
			}

			return nv.Add(time.Duration(-t2) * time.Millisecond)
		} else {
			return ToInt(nv.Sub(rsT.(time.Time)) / time.Millisecond)
		}

	default:
		return fmt.Errorf("unknown type(-): %T(%v)", n1A, n1A)
	}
}

var GetMinusResult = TKX.GetMinusResult

func (pA *TK) GetAddResult(n1A interface{}, n2A interface{}) (result interface{}) {
	defer func() {
		r := recover()

		if r != nil {
			result = fmt.Errorf("failed: %v(%v, %v)", r, n1A, n2A)
			return
		}
	}()

	switch nv := n1A.(type) {
	case byte:
		return nv + n2A.(byte)
	case rune:
		return nv + n2A.(rune)
	case int:
		return nv + n2A.(int)
	case int64:
		return nv + n2A.(int64)
	case float32:
		return nv + n2A.(float32)
	case float64:
		return nv + n2A.(float64)
	case string:
		return nv + n2A.(string)
	case []byte:
		return append(nv, n2A.([]byte)...)
	case time.Time:
		return nv.Add(time.Duration(time.Millisecond * time.Duration(ToInt(n2A))))
	default:
		return fmt.Errorf("unknown type(+): %T(%v)", n1A, n1A)
	}
}

var GetAddResult = TKX.GetAddResult

func (pA *TK) GetAddResultForce(n1A interface{}, n2A interface{}) (result interface{}) {
	defer func() {
		r := recover()

		if r != nil {
			result = fmt.Errorf("failed: %v(%v, %v)", r, n1A, n2A)
			return
		}
	}()

	if n1A == nil {
		return n2A
	}

	switch nv := n1A.(type) {
	case byte:
		return nv + ToByte(n2A)
	case rune:
		return nv + ToRune(n2A)
	case int:
		return nv + ToInt(n2A, 0)
	case int64:
		return nv + int64(ToInt(n2A))
	case float32:
		return nv + float32(ToFloat(n2A))
	case float64:
		return nv + ToFloat(n2A)
	case string:
		return nv + ToStr(n2A)
	case []byte:
		tmpBufT := DataToBytes(n2A)
		if IsError(tmpBufT) {
			return tmpBufT
		}

		return append(nv, tmpBufT.([]byte)...)
	case time.Time:
		return nv.Add(time.Duration(time.Millisecond * time.Duration(ToInt(n2A, 0))))
	default:
		return fmt.Errorf("unknown type(+): %T(%v)", n1A, n1A)
	}
}

var GetAddResultForce = TKX.GetAddResultForce

func (pA *TK) GetAddsResult(valuesA ...interface{}) (result interface{}) {
	defer func() {
		r := recover()

		if r != nil {
			result = fmt.Errorf("failed: %v(%v)", r, valuesA)
			return
		}
	}()

	var valueT interface{}

	for i, v := range valuesA {
		if i == 0 {
			valueT = v

			continue
		}

		valueT = GetAddResultForce(valueT, v)
	}

	return valueT

}

var GetAddsResult = TKX.GetAddsResult

func (pA *TK) GetBitANDResult(n1A interface{}, n2A interface{}) (result interface{}) {
	defer func() {
		r := recover()

		if r != nil {
			result = fmt.Errorf("failed: %v(%v, %v)", r, n1A, n2A)
			return
		}
	}()

	switch nv := n1A.(type) {
	case byte:
		return nv & n2A.(byte)
	case rune:
		return nv & n2A.(rune)
	case int:
		return nv & n2A.(int)
	case int64:
		return nv & n2A.(int64)
	default:
		return fmt.Errorf("unknown type(+): %T(%v)", n1A, n1A)
	}
}

var GetBitANDResult = TKX.GetBitANDResult

func (pA *TK) GetBitORResult(n1A interface{}, n2A interface{}) (result interface{}) {
	defer func() {
		r := recover()

		if r != nil {
			result = fmt.Errorf("failed: %v(%v, %v)", r, n1A, n2A)
			return
		}
	}()

	switch nv := n1A.(type) {
	case byte:
		return nv | n2A.(byte)
	case rune:
		return nv | n2A.(rune)
	case int:
		return nv | n2A.(int)
	case int64:
		return nv | n2A.(int64)
	default:
		return fmt.Errorf("unknown type(+): %T(%v)", n1A, n1A)
	}
}

var GetBitORResult = TKX.GetBitORResult

func (pA *TK) GetBitXORResult(n1A interface{}, n2A interface{}) (result interface{}) {
	defer func() {
		r := recover()

		if r != nil {
			result = fmt.Errorf("failed: %v(%v, %v)", r, n1A, n2A)
			return
		}
	}()

	switch nv := n1A.(type) {
	case byte:
		return nv ^ n2A.(byte)
	case rune:
		return nv ^ n2A.(rune)
	case int:
		return nv ^ n2A.(int)
	case int64:
		return nv ^ n2A.(int64)
	default:
		return fmt.Errorf("unknown type(+): %T(%v)", n1A, n1A)
	}
}

var GetBitXORResult = TKX.GetBitXORResult

func (pA *TK) GetBitANDNOTResult(n1A interface{}, n2A interface{}) (result interface{}) {
	defer func() {
		r := recover()

		if r != nil {
			result = fmt.Errorf("failed: %v(%v, %v)", r, n1A, n2A)
			return
		}
	}()

	switch nv := n1A.(type) {
	case byte:
		return nv &^ n2A.(byte)
	case rune:
		return nv &^ n2A.(rune)
	case int:
		return nv &^ n2A.(int)
	case int64:
		return nv &^ n2A.(int64)
	default:
		return fmt.Errorf("unknown type(+): %T(%v)", n1A, n1A)
	}
}

var GetBitANDNOTResult = TKX.GetBitANDNOTResult

func (pA *TK) GetGTResult(n1A interface{}, n2A interface{}) (result interface{}) {
	defer func() {
		r := recover()

		if r != nil {
			result = fmt.Errorf("failed: %v(%v, %v)", r, n1A, n2A)
			return
		}
	}()

	switch nv := n1A.(type) {
	case byte:
		return nv > n2A.(byte)
	case rune:
		return nv > n2A.(rune)
	case int:
		return nv > n2A.(int)
	case int64:
		return nv > n2A.(int64)
	case float32:
		return nv > n2A.(float32)
	case float64:
		return nv > n2A.(float64)
	case string:
		return nv > n2A.(string)
	case time.Time:
		return nv.After(n2A.(time.Time))
	default:
		return fmt.Errorf("unknown type(>): %T(%v)", n1A, n1A)
	}
}

var GetGTResult = TKX.GetGTResult

func (pA *TK) GetGETResult(n1A interface{}, n2A interface{}) (result interface{}) {
	defer func() {
		r := recover()

		if r != nil {
			result = fmt.Errorf("failed: %v(%v, %v)", r, n1A, n2A)
			return
		}
	}()

	switch nv := n1A.(type) {
	case byte:
		return nv >= n2A.(byte)
	case rune:
		return nv >= n2A.(rune)
	case int:
		return nv >= n2A.(int)
	case int64:
		return nv >= n2A.(int64)
	case float32:
		return nv >= n2A.(float32)
	case float64:
		return nv >= n2A.(float64)
	case string:
		return nv >= n2A.(string)
	case time.Time:
		return !nv.Before(n2A.(time.Time))
	default:
		return fmt.Errorf("unknown type(>=): %T(%v)", n1A, n1A)
	}
}

var GetGETResult = TKX.GetGETResult

func (pA *TK) GetLTResult(n1A interface{}, n2A interface{}) (result interface{}) {
	defer func() {
		r := recover()

		if r != nil {
			result = fmt.Errorf("failed: %v(%v, %v)", r, n1A, n2A)
			return
		}
	}()

	switch nv := n1A.(type) {
	case byte:
		return nv < n2A.(byte)
	case rune:
		return nv < n2A.(rune)
	case int:
		return nv < n2A.(int)
	case int64:
		return nv < n2A.(int64)
	case float32:
		return nv < n2A.(float32)
	case float64:
		return nv < n2A.(float64)
	case string:
		return nv < n2A.(string)
	case time.Time:
		return nv.Before(n2A.(time.Time))
	default:
		return fmt.Errorf("unknown type(<): %T(%v)", n1A, n1A)
	}
}

var GetLTResult = TKX.GetLTResult

func (pA *TK) GetLETResult(n1A interface{}, n2A interface{}) (result interface{}) {
	defer func() {
		r := recover()

		if r != nil {
			result = fmt.Errorf("failed: %v(%v, %v)", r, n1A, n2A)
			return
		}
	}()

	switch nv := n1A.(type) {
	case byte:
		return nv <= n2A.(byte)
	case rune:
		return nv <= n2A.(rune)
	case int:
		return nv <= n2A.(int)
	case int64:
		return nv <= n2A.(int64)
	case float32:
		return nv <= n2A.(float32)
	case float64:
		return nv <= n2A.(float64)
	case string:
		return nv <= n2A.(string)
	case time.Time:
		return !nv.After(n2A.(time.Time))
	default:
		return fmt.Errorf("unknown type(<=): %T(%v)", n1A, n1A)
	}
}

var GetLETResult = TKX.GetLETResult

func (pA *TK) GetEQResult(n1A interface{}, n2A interface{}) (result interface{}) {
	defer func() {
		r := recover()

		if r != nil {
			result = fmt.Errorf("failed: %v(%v, %v)", r, n1A, n2A)
			return
		}
	}()

	switch nv := n1A.(type) {
	case bool:
		nv2, ok := n2A.(bool)
		if !ok {
			return false
		}
		return nv == nv2
	case byte:
		nv2, ok := n2A.(byte)
		if !ok {
			return false
		}
		return nv == nv2
	case rune:
		nv2, ok := n2A.(rune)
		if !ok {
			return false
		}
		return nv == nv2
	case int:
		nv2, ok := n2A.(int)
		if !ok {
			return false
		}
		return nv == nv2
	case int64:
		nv2, ok := n2A.(int64)
		if !ok {
			return false
		}
		return nv == nv2
	case float32:
		nv2, ok := n2A.(float32)
		if !ok {
			return false
		}
		return nv == nv2
	case float64:
		nv2, ok := n2A.(float64)
		if !ok {
			return false
		}
		return nv == nv2
	case string:
		nv2, ok := n2A.(string)
		if !ok {
			return false
		}
		return nv == nv2
	case time.Time:
		nv2, ok := n2A.(time.Time)
		if !ok {
			return false
		}
		return !nv.Before(nv2) && !nv.After(nv2)
	case interface{}:
		return nv == n2A.(interface{})
	default:
		return fmt.Errorf("unknown type(==): %T(%v)", n1A, n1A)
	}
}

var GetEQResult = TKX.GetEQResult

func (pA *TK) GetNEQResult(n1A interface{}, n2A interface{}) (result interface{}) {
	defer func() {
		r := recover()

		if r != nil {
			result = fmt.Errorf("failed: %v(%v, %v)", r, n1A, n2A)
			return
		}
	}()

	switch nv := n1A.(type) {
	case bool:
		n2v, ok := n2A.(bool)
		if !ok {
			return true
		}

		return nv != n2v
	case byte:
		n2v, ok := n2A.(byte)
		if !ok {
			return true
		}

		return nv != n2v
	case rune:
		n2v, ok := n2A.(rune)
		if !ok {
			return true
		}

		return nv != n2v
	case int:
		n2v, ok := n2A.(int)
		if !ok {
			return true
		}

		return nv != n2v
	case int64:
		n2v, ok := n2A.(int64)
		if !ok {
			return true
		}

		return nv != n2v
	case float32:
		n2v, ok := n2A.(float32)
		if !ok {
			return true
		}

		return nv != n2v
	case float64:
		n2v, ok := n2A.(float64)
		if !ok {
			return true
		}

		return nv != n2v
	case string:
		n2v, ok := n2A.(string)
		if !ok {
			return true
		}

		return nv != n2v
	case time.Time:
		n2v, ok := n2A.(time.Time)
		if !ok {
			return true
		}

		return !nv.Before(n2v) || !nv.After(n2v)
	case interface{}:
		return nv != n2A
	default:
		return fmt.Errorf("unknown type(!=): %T(%v)", n1A, n1A)
	}
}

var GetNEQResult = TKX.GetNEQResult

func (pA *TK) GetANDResult(n1A interface{}, n2A interface{}) (result interface{}) {
	defer func() {
		r := recover()

		if r != nil {
			result = fmt.Errorf("failed: %v(%v, %v)", r, n1A, n2A)
			return
		}
	}()

	switch nv := n1A.(type) {
	case bool:
		if !nv {
			return false
		}

		return n2A.(bool)

		// return nv && n2A.(bool)
	default:
		return fmt.Errorf("unknown type(!=): %T(%v)", n1A, n1A)
	}
}

var GetANDResult = TKX.GetANDResult

func (pA *TK) GetORResult(n1A interface{}, n2A interface{}) (result interface{}) {
	defer func() {
		r := recover()

		if r != nil {
			result = fmt.Errorf("failed: %v(%v, %v)", r, n1A, n2A)
			return
		}
	}()

	switch nv := n1A.(type) {
	case bool:
		return nv || n2A.(bool)
	default:
		return fmt.Errorf("unknown type(!=): %T(%v)", n1A, n1A)
	}
}

var GetORResult = TKX.GetORResult

func (pA *TK) AdjustFloat(nA float64, digitA ...int) float64 {
	digitT := 10
	if len(digitA) > 0 {
		digitT = digitA[0]
	}

	return math.Round(nA*math.Pow10(digitT)) / math.Pow10(digitT)
}

var AdjustFloat = TKX.AdjustFloat

// func NewBigFloatZero() *big.Float {
// 	r := big.NewFloat(0.0)
// 	r.SetPrec(256)
// 	return r
// }

// func BigMul(a, b *big.Float) *big.Float {
// 	return NewBigFloatZero().Mul(a, b)
// }

func (pA *TK) AdjustBigFloat(nA *big.Float, digitA ...int) *big.Float {
	digitT := 10
	if len(digitA) > 0 {
		digitT = digitA[0]
	}

	if digitT < 0 {
		digitT = 0
	}

	result := big.NewFloat(0).Set(nA)

	for i := 0; i < digitT; i++ {
		result = big.NewFloat(0).Mul(result, big.NewFloat(10))
	}

	n1, _ := result.Int64()
	f1 := big.NewFloat(0).SetInt64(n1)

	f2 := big.NewFloat(0).Quo(f1, big.NewFloat(math.Pow10(digitT)))

	return f2
}

var AdjustBigFloat = TKX.AdjustBigFloat

func (pA *TK) LimitPrecision(nA interface{}, digitA ...int) error {
	digitT := 10
	if len(digitA) > 0 {
		digitT = digitA[0]
	}

	switch t := nA.(type) {
	case *float64:
		vT := *(nA.(*float64))
		*(nA.(*float64)) = math.Round(vT*math.Pow10(digitT)) / math.Pow10(digitT)
	case *[]float64:
		pT := (nA.(*[]float64))

		lenT := len(*pT)

		for i := 0; i < lenT; i++ {
			(*pT)[i] = math.Round((*pT)[i]*math.Pow10(digitT)) / math.Pow10(digitT)
		}

	default:
		return Errf("%v", "unknown type: %v", t)
	}

	return nil

}

var LimitPrecision = TKX.LimitPrecision

func (pA *TK) Float32ArrayToFloat64Array(aryA []float32) []float64 {
	if aryA == nil {
		return nil
	}

	rs := make([]float64, len(aryA))

	for i := 0; i < len(aryA); i++ {
		rs[i] = float64(aryA[i])
	}

	return rs
}

var Float32ArrayToFloat64Array = TKX.Float32ArrayToFloat64Array

func (pA *TK) GenerateRandomFloats(sizeA int) []float64 {
	bufT := make([]float64, sizeA)

	Randomize()

	for i := 0; i < sizeA; i++ {
		bufT[i] = rand.Float64()
	}

	return bufT
}

var GenerateRandomFloats = TKX.GenerateRandomFloats

func (pA *TK) CalCosineSimilarityBetweenFloatsBig(f1, f2 []float64) float64 {
	if f1 == nil || f2 == nil {
		return -1
	}

	l1 := len(f1)
	l2 := len(f2)

	if l1 != l2 {
		return -1
	}

	var rr *big.Float = new(big.Float)
	var f1r *big.Float = new(big.Float)
	var f2r *big.Float = new(big.Float)

	for i := 0; i < l1; i++ {
		f1b := new(big.Float).SetFloat64(f1[i])
		f2b := new(big.Float).SetFloat64(f2[i])
		rr.Add(rr, new(big.Float).Mul(f1b, f2b))
		f1r.Add(f1r, new(big.Float).Mul(f1b, f1b))
		f2r.Add(f2r, new(big.Float).Mul(f2b, f2b))
	}

	tmprs1 := f1r.Sqrt(f1r)
	tmprs2 := f2r.Sqrt(f2r)

	tmprsr := new(big.Float).Mul(tmprs1, tmprs2)

	rs, _ := (rr.Quo(rr, tmprsr)).Float64()

	return rs
}

var CalCosineSimilarityBetweenFloatsBig = TKX.CalCosineSimilarityBetweenFloatsBig

// 数据库相关 database related

// GetDBConnection must close it manually
func (pA *TK) GetDBConnection(driverA string, pathT string) *sql.DB {
	dbT, errT := sql.Open(driverA, pathT)

	if errT != nil {
		return nil
	}

	errT = dbT.Ping()

	if errT != nil {
		dbT.Close()
		return nil
	}

	return dbT
}

var GetDBConnection = TKX.GetDBConnection

func (pA *TK) FormatSQLValue(strA string) string {
	strT := strings.Replace(strA, "\r", "\\r", -1)
	strT = strings.Replace(strT, "\n", "\\n", -1)
	strT = strings.Replace(strT, "'", "''", -1)

	return strT
}

var FormatSQLValue = TKX.FormatSQLValue

func (pA *TK) ListToSQLList(vA interface{}) string {
	bufT := new(strings.Builder)

	bufT.WriteString("(")

	switch nv := vA.(type) {
	case []string:
		for i, v := range nv {
			if i > 0 {
				bufT.WriteString(",")
			}
			bufT.WriteString("'")
			v = strings.Replace(v, "'", "''", -1)
			bufT.WriteString(v)
			bufT.WriteString("'")
		}
	case []interface{}:
		for i, v := range nv {
			if i > 0 {
				bufT.WriteString(",")
			}
			bufT.WriteString("'")
			vs := strings.Replace(ToStr(v), "'", "''", -1)
			bufT.WriteString(vs)
			bufT.WriteString("'")
		}
	}

	bufT.WriteString(")")

	return bufT.String()
}

var ListToSQLList = TKX.ListToSQLList

// GetDBRowCount 获取类似select count(*)的结果
func (pA *TK) GetDBRowCount(dbA *sql.DB, sqlA string) (int, error) {
	if dbA == nil {
		return 0, fmt.Errorf("DB pointer nil")
	}

	var c int

	errT := dbA.QueryRow(sqlA).Scan(&c)

	if errT == sql.ErrNoRows {
		return 0, fmt.Errorf("no rows")
	}

	return c, nil
}

var GetDBRowCount = TKX.GetDBRowCount

// GetDBRowCountCompact 获取类似select count(*)的结果
// return < 0 if fail
func (pA *TK) GetDBRowCountCompact(dbA *sql.DB, sqlA string) int {
	c, errT := GetDBRowCount(dbA, sqlA)

	if errT != nil {
		return -1
	}

	return c
}

var GetDBRowCountCompact = TKX.GetDBRowCountCompact

// GetDBResultString 获取类似select a from ...的结果
func (pA *TK) GetDBResultString(dbA *sql.DB, sqlA string) (string, error) {
	if dbA == nil {
		return "", fmt.Errorf("DB pointer nil")
	}

	var s string

	errT := dbA.QueryRow(sqlA).Scan(&s)

	if errT == sql.ErrNoRows {
		return "", fmt.Errorf("no rows")
	}

	return s, nil
}

var GetDBResultString = TKX.GetDBResultString

// GetDBResultArray 获取类似select a from ...的多行结果
func (pA *TK) GetDBResultArray(dbA *sql.DB, sqlA string) ([][]string, error) {
	if dbA == nil {
		return nil, fmt.Errorf("DB pointer nil")
	}

	rowsT, errT := dbA.Query(sqlA)
	if errT != nil {
		return nil, errT
	}

	defer rowsT.Close()

	columnsT, errT := rowsT.Columns()

	if errT != nil {
		return nil, errT
	}

	// columnsTypesT, errT := rowsT.ColumnTypes()

	columnCountT := len(columnsT)

	if columnCountT < 1 {
		return nil, Errf("zero columns")
	}

	sliceT := make([][]string, 0)

	for rowsT.Next() {
		subSliceT := make([]interface{}, columnCountT)
		subSlicePointerT := make([]interface{}, columnCountT)

		for j := 0; j < columnCountT; j++ {
			subSlicePointerT[j] = &subSliceT[j]
		}

		errT := rowsT.Scan(subSlicePointerT...)

		if errT != nil {
			return nil, errT
		}

		// subSliceT[j] = tmps

		errT = rowsT.Err()
		if errT != nil {
			return nil, errT
		}

		subStringSliceT := make([]string, columnCountT)

		for j := 0; j < columnCountT; j++ {
			subStringSliceT[j] = Spr("%v", subSliceT[j])
		}

		sliceT = append(sliceT, subStringSliceT)

	}

	return sliceT, nil
}

var GetDBResultArray = TKX.GetDBResultArray

func (pA *TK) RecordsToMapArray(recA interface{}) []map[string]string {
	if recA == nil {
		return nil
	}

	nv1, ok := recA.([][]string)

	if ok {
		lenT := len(nv1)

		if lenT < 1 {
			return nil
		}

		lineLenT := len(nv1[0])

		aryT := make([]map[string]string, lenT-1)

		for i := 1; i < lenT; i++ {
			mapT := make(map[string]string, lenT)

			curLineLenT := len(nv1[i])
			for j := 0; j < lineLenT; j++ {
				if j >= curLineLenT {
					mapT[nv1[0][j]] = ""
					continue
				}
				mapT[nv1[0][j]] = nv1[i][j]
			}

			aryT[i-1] = mapT
		}

		return aryT
	}

	nv2, ok := recA.([][]interface{})

	if ok {
		lenT := len(nv2)

		if lenT < 1 {
			return nil
		}

		lineLenT := len(nv2[0])

		aryT := make([]map[string]string, lenT-1)

		keysT := make([]string, lineLenT)

		for i0 := 0; i0 < lineLenT; i0++ {
			keysT[i0] = ToStr(nv2[0][i0])
		}

		for i := 1; i < lenT; i++ {
			mapT := make(map[string]string, lenT)

			curLineLenT := len(nv2[i])
			for j := 0; j < lineLenT; j++ {
				if j >= curLineLenT {
					mapT[keysT[j]] = ""
					continue
				}
				mapT[keysT[j]] = ToStr(nv2[i][j])
			}

			aryT[i-1] = mapT
		}

		return aryT
	}

	nv3, ok := recA.([]interface{})

	if ok {
		lenT := len(nv3)

		if lenT < 1 {
			return nil
		}

		nv3a, ok := nv3[0].([]interface{})

		if ok {
			// Pl("weird type: %T(%#v)", nv3[0], nv3[0])
			// return nil
			lineLenT := len(nv3a)

			aryT := make([]map[string]string, lenT-1)

			keysT := make([]string, lineLenT)

			for i0 := 0; i0 < lineLenT; i0++ {
				keysT[i0] = ToStr(nv3a[i0])
			}

			for i := 1; i < lenT; i++ {
				mapT := make(map[string]string, lenT)

				curLineLenT := len(nv3[i].([]interface{}))
				for j := 0; j < lineLenT; j++ {
					if j >= curLineLenT {
						mapT[keysT[j]] = ""
						continue
					}
					mapT[keysT[j]] = ToStr(nv3[i].([]interface{})[j])
				}

				aryT[i-1] = mapT
			}

			return aryT
		}

		nv3b, ok := nv3[0].([]string)

		if ok {
			lineLenT := len(nv3b)

			aryT := make([]map[string]string, lenT-1)

			keysT := make([]string, lineLenT)

			for i0 := 0; i0 < lineLenT; i0++ {
				keysT[i0] = ToStr(nv3b[i0])
			}

			for i := 1; i < lenT; i++ {
				mapT := make(map[string]string, lenT)

				curLineLenT := len(nv3[i].([]interface{}))
				for j := 0; j < lineLenT; j++ {
					if j >= curLineLenT {
						mapT[keysT[j]] = ""
						continue
					}
					mapT[keysT[j]] = ToStr(nv3[i].([]string)[j])
				}

				aryT[i-1] = mapT
			}

			return aryT
		}

	}

	// Pl("unsupported type: %T(%#v)", recA, recA)
	return nil
}

var RecordsToMapArray = TKX.RecordsToMapArray

func (pA *TK) RecordsToOrderedMapArray(recA interface{}) interface{} {
	if recA == nil {
		return fmt.Errorf("empty records")
	}

	nv1, ok := recA.([][]string)

	if ok {
		lenT := len(nv1)

		if lenT < 1 {
			return nil
		}

		lineLenT := len(nv1[0])

		aryT := make([]*OrderedMap, lenT-1)

		for i := 1; i < lenT; i++ {
			mapT := NewOrderedMap()

			curLineLenT := len(nv1[i])
			for j := 0; j < lineLenT; j++ {
				if j >= curLineLenT {
					mapT.Set(nv1[0][j], "")
					continue
				}
				mapT.Set(nv1[0][j], nv1[i][j])
			}

			aryT[i-1] = mapT
		}

		return aryT
	}

	nv2, ok := recA.([][]interface{})

	if ok {
		lenT := len(nv2)

		if lenT < 1 {
			return nil
		}

		lineLenT := len(nv2[0])

		aryT := make([]*OrderedMap, lenT-1)

		keysT := make([]string, lineLenT)

		for i0 := 0; i0 < lineLenT; i0++ {
			keysT[i0] = ToStr(nv2[0][i0])
		}

		for i := 1; i < lenT; i++ {
			mapT := NewOrderedMap()

			curLineLenT := len(nv2[i])
			for j := 0; j < lineLenT; j++ {
				if j >= curLineLenT {
					mapT.Set(keysT[j], "")
					continue
				}
				mapT.Set(keysT[j], ToStr(nv2[i][j]))
			}

			aryT[i-1] = mapT
		}

		return aryT
	}

	nv3, ok := recA.([]interface{})

	if ok {
		lenT := len(nv3)

		if lenT < 1 {
			return nil
		}

		nv3a, ok := nv3[0].([]interface{})

		if ok {
			// Pl("weird type: %T(%#v)", nv3[0], nv3[0])
			// return nil
			lineLenT := len(nv3a)

			aryT := make([]*OrderedMap, lenT-1)

			keysT := make([]string, lineLenT)

			for i0 := 0; i0 < lineLenT; i0++ {
				keysT[i0] = ToStr(nv3a[i0])
			}

			for i := 1; i < lenT; i++ {
				mapT := NewOrderedMap()

				curLineLenT := len(nv3[i].([]interface{}))
				for j := 0; j < lineLenT; j++ {
					if j >= curLineLenT {
						mapT.Set(keysT[j], "")
						continue
					}
					mapT.Set(keysT[j], ToStr(nv3[i].([]interface{})[j]))
				}

				aryT[i-1] = mapT
			}

			return aryT
		}

		nv3b, ok := nv3[0].([]string)

		if ok {
			lineLenT := len(nv3b)

			aryT := make([]*OrderedMap, lenT-1)

			keysT := make([]string, lineLenT)

			for i0 := 0; i0 < lineLenT; i0++ {
				keysT[i0] = ToStr(nv3b[i0])
			}

			for i := 1; i < lenT; i++ {
				mapT := NewOrderedMap()

				curLineLenT := len(nv3[i].([]interface{}))
				for j := 0; j < lineLenT; j++ {
					if j >= curLineLenT {
						mapT.Set(keysT[j], "")
						continue
					}
					mapT.Set(keysT[j], ToStr(nv3[i].([]string)[j]))
				}

				aryT[i-1] = mapT
			}

			return aryT
		}

	}

	// Pl("unsupported type: %T(%#v)", recA, recA)
	return fmt.Errorf("unsupported type: %T(%#v)", recA, recA)
}

var RecordsToOrderedMapArray = TKX.RecordsToOrderedMapArray

func (pA *TK) RecordsToMapStringMap(recA interface{}, keyA string) interface{} {
	if recA == nil {
		return fmt.Errorf("empty records")
	}

	nv1, ok := recA.([][]string)

	if ok {
		lenT := len(nv1)

		if lenT < 1 {
			return nil
		}

		lineLenT := len(nv1[0])

		aryT := make(map[string]map[string]string, lenT-1)

		for i := 1; i < lenT; i++ {
			mapT := make(map[string]string)

			for j := 0; j < lineLenT; j++ {
				mapT[nv1[0][j]] = nv1[i][j]
			}

			aryT[mapT[keyA]] = mapT
		}

		return aryT
	}

	nv2, ok := recA.([][]interface{})

	if ok {
		lenT := len(nv2)

		if lenT < 1 {
			return nil
		}

		lineLenT := len(nv2[0])

		aryT := make(map[string]map[string]string, lenT-1)

		keysT := make([]string, lineLenT)

		for i0 := 0; i0 < lineLenT; i0++ {
			keysT[i0] = ToStr(nv2[0][i0])
		}

		for i := 1; i < lenT; i++ {
			mapT := make(map[string]string)

			for j := 0; j < lineLenT; j++ {
				mapT[keysT[j]] = ToStr(nv2[i][j])
			}

			aryT[mapT[keyA]] = mapT
		}

		return aryT
	}

	nv3, ok := recA.([]interface{})

	if ok {
		lenT := len(nv3)

		if lenT < 1 {
			return nil
		}

		nv3a, ok := nv3[0].([]interface{})

		if ok {
			// Pl("weird type: %T(%#v)", nv3[0], nv3[0])
			// return nil
			lineLenT := len(nv3a)

			aryT := make(map[string]map[string]string, lenT-1)

			keysT := make([]string, lineLenT)

			for i0 := 0; i0 < lineLenT; i0++ {
				keysT[i0] = ToStr(nv3a[i0])
			}

			for i := 1; i < lenT; i++ {
				mapT := make(map[string]string)

				for j := 0; j < lineLenT; j++ {
					mapT[keysT[j]] = ToStr(nv3[i].([]interface{})[j])
				}

				aryT[mapT[keyA]] = mapT
			}

			return aryT
		}

		nv3b, ok := nv3[0].([]string)

		if ok {
			lineLenT := len(nv3b)

			aryT := make(map[string]map[string]string, lenT-1)

			keysT := make([]string, lineLenT)

			for i0 := 0; i0 < lineLenT; i0++ {
				keysT[i0] = ToStr(nv3b[i0])
			}

			for i := 1; i < lenT; i++ {
				mapT := make(map[string]string)

				for j := 0; j < lineLenT; j++ {
					mapT[keysT[j]] = ToStr(nv3[i].([]string)[j])
				}

				aryT[mapT[keyA]] = mapT
			}

			return aryT
		}

	}

	// Pl("unsupported type: %T(%#v)", recA, recA)
	return fmt.Errorf("unsupported type: %T(%#v)", recA, recA)
}

var RecordsToMapStringMap = TKX.RecordsToMapStringMap

func (pA *TK) RecordsToOrderedMapStringMap(recA interface{}, keyA string) interface{} {
	if recA == nil {
		return fmt.Errorf("empty records")
	}

	nv1, ok := recA.([][]string)

	if ok {
		lenT := len(nv1)

		if lenT < 1 {
			return nil
		}

		lineLenT := len(nv1[0])

		aryT := make(map[string]*OrderedMap)

		for i := 1; i < lenT; i++ {
			mapT := NewOrderedMap()

			for j := 0; j < lineLenT; j++ {
				mapT.Set(nv1[0][j], nv1[i][j])
			}

			aryT[mapT.GetString(keyA)] = mapT
		}

		return aryT
	}

	nv2, ok := recA.([][]interface{})

	if ok {
		lenT := len(nv2)

		if lenT < 1 {
			return nil
		}

		lineLenT := len(nv2[0])

		aryT := make(map[string]*OrderedMap, lenT-1)

		keysT := make([]string, lineLenT)

		for i0 := 0; i0 < lineLenT; i0++ {
			keysT[i0] = ToStr(nv2[0][i0])
		}

		for i := 1; i < lenT; i++ {
			mapT := NewOrderedMap()

			for j := 0; j < lineLenT; j++ {
				mapT.Set(keysT[j], ToStr(nv2[i][j]))
			}

			aryT[mapT.GetString(keyA)] = mapT
		}

		return aryT
	}

	nv3, ok := recA.([]interface{})

	if ok {
		lenT := len(nv3)

		if lenT < 1 {
			return nil
		}

		nv3a, ok := nv3[0].([]interface{})

		if ok {
			// Pl("weird type: %T(%#v)", nv3[0], nv3[0])
			// return nil
			lineLenT := len(nv3a)

			aryT := make(map[string]*OrderedMap, lenT-1)

			keysT := make([]string, lineLenT)

			for i0 := 0; i0 < lineLenT; i0++ {
				keysT[i0] = ToStr(nv3a[i0])
			}

			for i := 1; i < lenT; i++ {
				mapT := NewOrderedMap()

				for j := 0; j < lineLenT; j++ {
					mapT.Set(keysT[j], ToStr(nv3[i].([]interface{})[j]))
				}

				aryT[mapT.GetString(keyA)] = mapT
			}

			return aryT
		}

		nv3b, ok := nv3[0].([]string)

		if ok {
			lineLenT := len(nv3b)

			aryT := make(map[string]*OrderedMap, lenT-1)

			keysT := make([]string, lineLenT)

			for i0 := 0; i0 < lineLenT; i0++ {
				keysT[i0] = ToStr(nv3b[i0])
			}

			for i := 1; i < lenT; i++ {
				mapT := NewOrderedMap()

				for j := 0; j < lineLenT; j++ {
					mapT.Set(keysT[j], ToStr(nv3[i].([]string)[j]))
				}

				aryT[mapT.GetString(keyA)] = mapT
			}

			return aryT
		}

	}

	// Pl("unsupported type: %T(%#v)", recA, recA)
	return fmt.Errorf("unsupported type: %T(%#v)", recA, recA)
}

var RecordsToOrderedMapStringMap = TKX.RecordsToOrderedMapStringMap

// 文本编码相关 encoding related

// ConvertToGB18030 转换UTF-8字符串为GB18030编码
func (pA *TK) ConvertToGB18030(srcA string) string {
	encoderT := mahonia.NewEncoder("GB18030")

	return encoderT.ConvertString(srcA)
	// dst := make([]byte, len(srcA)*2)
	// transformer := simplifiedchinese.GB18030.NewEncoder()
	// nDst, _, err := transformer.Transform(dst, []byte(srcA), true)
	// if err != nil {
	// 	return GenerateErrorStringF("encoding failed")
	// }
	// return string(dst[:nDst])
}

var ConvertToGB18030 = TKX.ConvertToGB18030

// ConvertToGB18030Bytes 转换UTF-8字符串为GB18030编码的字节切片
// func ConvertToGB18030Bytes(srcA string) []byte {
// 	dst := make([]byte, len(srcA)*2)

// 	transformer := simplifiedchinese.GB18030.NewEncoder()

// 	nDst, _, err := transformer.Transform(dst, []byte(srcA), true)

// 	if err != nil {
// 		return nil
// 	}

//		return dst[:nDst]
//	}
func (pA *TK) ConvertToGB18030Bytes(srcA string) []byte {

	encoderT := mahonia.NewEncoder("GB18030")

	tmps := encoderT.ConvertString(srcA)

	// destT := []byte(tmps)

	return []byte(tmps)
}

var ConvertToGB18030Bytes = TKX.ConvertToGB18030Bytes

// func ConvertToUTF8(srcA []byte, srcEncA string) string {
// 	srcEncT := srcEncA

//		switch srcEncT {
//		case "", "GB18030", "gb18030", "GBK", "gbk", "GB2312", "gb2312":
//			dst := make([]byte, len(srcA)*2)
//			transformer := simplifiedchinese.GB18030.NewDecoder()
//			nDst, _, err := transformer.Transform(dst, srcA, true)
//			if err != nil {
//				return GenerateErrorStringF("encoding failed: %v", err.Error())
//			}
//			return string(dst[:nDst])
//		case "utf-8", "UTF-8":
//			return string(srcA)
//		case "windows-1252", "windows1252":
//			dst := make([]byte, len(srcA)*2)
//			transformer := charmap.Windows1252.NewDecoder()
//			nDst, _, err := transformer.Transform(dst, srcA, true)
//			if err != nil {
//				return GenerateErrorStringF("encoding failed: %v", srcEncA)
//			}
//			return string(dst[:nDst])
//		default:
//			return GenerateErrorStringF("unknown encoding: %v", srcEncA)
//		}
//	}
//
// ConvertToUTF8 转换GB18030编码等字符串(字节形式)为UTF-8字符串
func (pA *TK) ConvertToUTF8(srcA []byte, srcEncA string) string {
	srcEncT := srcEncA

	if srcEncT == "" {
		srcEncT = "GB18030"
	}

	decodeT := mahonia.NewDecoder(srcEncT)

	_, cdataT, errT := decodeT.Translate(srcA, true)

	if errT != nil {
		return GenerateErrorStringF("encoding failed: %v", errT.Error())
	}

	return string(cdataT)

}

var ConvertToUTF8 = TKX.ConvertToUTF8

func (pA *TK) Utf8ToGb(strA string) interface{} {
	dst := make([]byte, len(strA) * 2)
	tr := simplifiedchinese.GB18030.NewEncoder()
	
	nDst, _, err := tr.Transform(dst, []byte(strA), true)
	
	if err != nil {
		return err
	}
	
	return string(dst[:nDst])
}

var Utf8ToGb = TKX.Utf8ToGb

// ConvertStringToUTF8 转换GB18030编码等字符串为UTF-8字符串
func (pA *TK) ConvertStringToUTF8(srcA string, srcEncA string) (result string) {
	defer func() {
		r := recover()
		if r != nil {
			result = srcA
			return
		}
	}()

	srcEncT := srcEncA

	if srcEncT == "" {
		srcEncT = "GB18030"
	}
	
	if strings.ToLower(srcEncT) == "gb2312" {
		srcEncT = "GB18030"
	}

	decodeT := mahonia.NewDecoder(srcEncT)

	return decodeT.ConvertString(srcA)

}

var ConvertStringToUTF8 = TKX.ConvertStringToUTF8

// XML related

func (pA *TK) ReshapeXML(xmlA string) string {
	var errT error

	treeT := etree.NewDocument()

	if treeT == nil {
		return GenerateErrorStringF("create XML tree failed")
	}

	errT = treeT.ReadFromString(xmlA)

	if errT != nil {
		return GenerateErrorStringF("invalid XML: %v", errT)
	}

	treeT.Indent(2)

	outputT, errT := treeT.WriteToString()

	if errT != nil {
		return GenerateErrorStringF("failed to reshape XML: %v", errT)
	}

	return outputT

}

var ReshapeXML = TKX.ReshapeXML

func (pA *TK) FlattenXML(xmlA string, nodeA string) string {
	var errT error

	treeT := etree.NewDocument()

	if treeT == nil {
		return GenerateErrorStringF("create XML tree failed")
	}

	errT = treeT.ReadFromString(xmlA)

	if errT != nil {
		return GenerateErrorStringF("invalid XML: %v", errT)
	}

	rootNodeT := treeT.FindElement("//" + nodeA)

	if rootNodeT == nil {
		return GenerateErrorStringF("node not found: %v", nodeA)
	}

	nodesT := rootNodeT.ChildElements()

	var bufT strings.Builder

	for i, v := range nodesT {
		if i > 0 {
			bufT.WriteString("\n")
		}

		bufT.WriteString(Spr("%v: %v", v.Tag, v.Text()))
	}

	return bufT.String()

}

var FlattenXML = TKX.FlattenXML

func (pA *TK) GetMSSFromXML(xmlA string, nodeA string) (map[string]string, error) {
	var errT error

	treeT := etree.NewDocument()

	if treeT == nil {
		return nil, Errf("create XML tree failed")
	}

	errT = treeT.ReadFromString(xmlA)

	if errT != nil {
		return nil, Errf("invalid XML: %v", errT)
	}

	rootNodeT := treeT.FindElement("//" + nodeA)

	if rootNodeT == nil {
		return nil, Errf("node not found: %v", nodeA)
	}

	nodesT := rootNodeT.ChildElements()

	mapT := make(map[string]string, len(nodesT))
	for _, jv := range nodesT {
		mapT[jv.Tag] = jv.Text()
	}

	return mapT, nil
}

var GetMSSFromXML = TKX.GetMSSFromXML

func (pA *TK) ToXML(vA interface{}, optsA ...interface{}) interface{} {
	var errT error

	strT, ok := vA.(string)
	if ok {
		return EncodeToXMLString(strT)
	}

	defaultT := GetSwitchI(optsA, "-default=", "TXERROR:no default")

	rootNameT := GetSwitchI(optsA, "-root=", "")

	etreeNodeT, ok := vA.(*etree.Element)
	if ok {
		indentT := StrToInt(GetSwitchI(optsA, "-indent=", "2"), 2)
		newDocT := etree.NewDocument()
		newDocT.SetRoot(etreeNodeT)
		newDocT.Indent(indentT)
		xmlStrT, errT := newDocT.WriteToString()

		if errT != nil {
			if !IsErrStr(defaultT) {
				return defaultT
			}

			return errT
		}

		return xmlStrT
	}

	if !IfSwitchExistsWholeI(optsA, "-nomsi") {
		// msiT, ok := vA.(map[string]interface{})

		// if !ok {
		// 	if rootNameT == "" {
		// 		rootNameT = "root"
		// 	}

		// 	msiT = map[string]interface{}{
		// 		rootNameT: vA,
		// 	}
		// }

		indentT := StrToInt(GetSwitchI(optsA, "-indent=", "0"), 0)
		// ifIndentT := IfSwitchExistsWholeI(optsA, "-indent")
		ifCdataT := IfSwitchExistsWholeI(optsA, "-cdata")
		rootAttrsT := GetSwitchI(optsA, "-rootAttr=", "")

		var rootAttrMapT map[string]string = nil
		if rootAttrsT != "" {
			rootAttrMapT = JSONToMapStringString(rootAttrsT)
		}

		config := NewXMLFromAny(vA)

		if indentT > 0 {
			config.WithIndent("", strings.Repeat(" ", indentT))
		}

		if ifCdataT {
			config.AsCData()
		}

		if rootNameT != "" {
			config.WithRoot(rootNameT, rootAttrMapT)
		}

		rsT, errT := config.MarshalToString()

		if errT != nil {
			if !IsErrStr(defaultT) {
				return defaultT
			}

			return errT
		}

		return rsT

	}

	bufT, errT := xml.Marshal(vA)
	if errT != nil {
		if !IsErrStr(defaultT) {
			return defaultT
		}

		return errT
	}

	return string(bufT)
}

var ToXML = TKX.ToXML

func (pA *TK) GetNodeStringFromXML(xmlA string, nodeA string) (string, error) {
	var errT error

	treeT := etree.NewDocument()

	if treeT == nil {
		return "", Errf("create XML tree failed")
	}

	errT = treeT.ReadFromString(xmlA)

	if errT != nil {
		return "", Errf("invalid XML: %v", errT)
	}

	stringNodeT := treeT.FindElement("//" + nodeA)

	if stringNodeT == nil {
		return "", Errf("node not found: %v", nodeA)
	}

	// Pl("xmlnode: %v, %v", stringNodeT, stringNodeT.FullTag())

	return stringNodeT.Text(), nil
}

var GetNodeStringFromXML = TKX.GetNodeStringFromXML

func (pA *TK) GetNodesStringFromXML(xmlA string, nodeA string) interface{} {
	var errT error

	treeT := etree.NewDocument()

	if treeT == nil {
		return Errf("create XML tree failed")
	}

	errT = treeT.ReadFromString(xmlA)

	if errT != nil {
		return Errf("invalid XML: %v", errT)
	}

	stringNodesT := treeT.FindElements("//" + nodeA)

	aryT := make([]string, 0)

	for _, v := range stringNodesT {
		aryT = append(aryT, v.Text())
	}

	// Pl("xmlnode: %v, %v", stringNodeT, stringNodeT.FullTag())

	return aryT
}

var GetNodesStringFromXML = TKX.GetNodesStringFromXML

func (pA *TK) FromXMLX(xmlA string, nodeA ...interface{}) interface{} {
	var errT error

	treeT := etree.NewDocument()

	if treeT == nil {
		return Errf("create XML tree failed")
	}

	errT = treeT.ReadFromString(xmlA)

	if errT != nil {
		return Errf("invalid XML: %v", errT)
	}

	if len(nodeA) < 1 {
		return treeT.Root()
	}

	strT, ok := nodeA[0].(string)

	if !ok {
		return treeT.Root()
	}

	nodeT := treeT.FindElement(strT)

	if nodeT == nil {
		return Errf("node not found: %v", nodeA)
	}

	// Pl("xmlnode: %v, %v", stringNodeT, stringNodeT.FullTag())

	return nodeT
}

var FromXMLX = TKX.FromXMLX

func (pA *TK) GetXmlDocument(xmlA string) interface{} {
	var errT error

	treeT := etree.NewDocument()

	if treeT == nil {
		return Errf("create XML tree failed")
	}

	xmlT := strings.TrimSpace(xmlA)

	if xmlT == "" {
		return treeT
	}

	errT = treeT.ReadFromString(xmlT)

	if errT != nil {
		return Errf("invalid XML: %v", errT)
	}

	return treeT
}

var GetXmlDocument = TKX.GetXmlDocument

func (pA *TK) GetMSSArrayFromXML(xmlA string, nodeA string) ([]map[string]string, error) {
	var errT error

	treeT := etree.NewDocument()

	if treeT == nil {
		return nil, Errf("create XML tree failed")
	}

	errT = treeT.ReadFromString(xmlA)

	if errT != nil {
		return nil, Errf("invalid XML: %v", errT)
	}

	rootNodeT := treeT.FindElement("//" + nodeA)

	if rootNodeT == nil {
		return nil, Errf("node not found: %v", nodeA)
	}

	nodesT := rootNodeT.ChildElements()

	aryT := make([]map[string]string, 0)

	for _, v := range nodesT {
		internalNodesT := v.ChildElements()

		mapT := make(map[string]string, len(internalNodesT))
		for _, jv := range internalNodesT {
			mapT[jv.Tag] = jv.Text()
		}

		aryT = append(aryT, mapT)
	}

	return aryT, nil
}

var GetMSSArrayFromXML = TKX.GetMSSArrayFromXML

// GetXMLNode if no labelsA, return the root node, else return the specific node
// example: tk.GetXMLNode("... XML content", "envelop", "body", "anode")
func (pA *TK) GetXMLNode(xmlA string, labelsA ...string) (*xmlx.Node, error) {
	return xmlx.GetXMLNode(xmlA, labelsA...)
}

var GetXMLNode = TKX.GetXMLNode

// // decode xml to map[string]interface{}

// const (
// 	attrPrefix = "@"
// 	textPrefix = "#text"
// )

// var (
// 	//ErrInvalidDocument invalid document err
// 	ErrInvalidDocument = errors.New("invalid document")

// 	//ErrInvalidRoot data at the root level is invalid err
// 	ErrInvalidRoot = errors.New("data at the root level is invalid")
// )

// type node struct {
// 	Parent  *node
// 	Value   map[string]interface{}
// 	Attrs   []xml.Attr
// 	Label   string
// 	Text    string
// 	HasMany bool
// }

// // Decoder instance
// type XMLDecoder struct {
// 	r          io.Reader
// 	attrPrefix string
// 	textPrefix string
// }

// // NewXMLDecoder create new decoder instance
// func NewXMLDecoder(reader io.Reader) *XMLDecoder {
// 	return NewXMLDecoderWithPrefix(reader, attrPrefix, textPrefix)
// }

// // NewXMLDecoderWithPrefix create new decoder instance with custom attribute prefix and text prefix
// func NewXMLDecoderWithPrefix(reader io.Reader, attrPrefix, textPrefix string) *XMLDecoder {
// 	return &XMLDecoder{r: reader, attrPrefix: attrPrefix, textPrefix: textPrefix}
// }

// //Decode xml string to map[string]interface{}
// func (d *XMLDecoder) Decode() (map[string]interface{}, error) {
// 	decoder := xml.NewDecoder(d.r)
// 	n := &node{}
// 	stack := make([]*node, 0)

// 	for {
// 		token, err := decoder.Token()
// 		if err != nil && err != io.EOF {
// 			return nil, err
// 		}

// 		if token == nil {
// 			break
// 		}

// 		switch tok := token.(type) {
// 		case xml.StartElement:
// 			{
// 				n = &node{
// 					Label:  tok.Name.Local,
// 					Parent: n,
// 					Value:  map[string]interface{}{tok.Name.Local: map[string]interface{}{}},
// 					Attrs:  tok.Attr,
// 				}

// 				setAttrs(n, &tok, d.attrPrefix)
// 				stack = append(stack, n)

// 				if n.Parent != nil {
// 					n.Parent.HasMany = true
// 				}

// 				Pl("node add: %#v", stack[len(stack)-1])

// 			}

// 		case xml.CharData:
// 			Pl("CharData: %#v", string(token.(xml.CharData)))
// 			data := strings.TrimSpace(string(tok))
// 			if len(stack) > 0 {
// 				stack[len(stack)-1].Text = data
// 				Pl("insert")
// 			} else if len(data) > 0 {
// 				Pl("insert return")
// 				return nil, ErrInvalidRoot
// 			}
// 			Pl("insert not")

// 		case xml.EndElement:
// 			{
// 				length := len(stack)
// 				stack, n = stack[:length-1], stack[length-1]

// 				if !n.HasMany {
// 					if len(n.Attrs) > 0 {
// 						m := n.Value[n.Label].(map[string]interface{})
// 						m[d.textPrefix] = n.Text
// 					} else {
// 						n.Value[n.Label] = n.Text
// 						Pl("n.Value: %v, n.Label: %v, n.Text: %v", n.Value, n.Label, n.Text)
// 					}
// 				}

// 				if len(stack) == 0 {
// 					return n.Value, nil
// 				}

// 				setNodeValue(n)
// 				n = n.Parent

// 				Pl("n = n.Parent")
// 			}
// 		case xml.ProcInst:
// 			{
// 				tt := token.(xml.ProcInst)
// 				Pl("tt: %v", string(tt.Inst))
// 			}
// 		default:
// 			Pl("token: %#v, tok: %#v", token, tok)
// 		}

// 	}

// 	return nil, ErrInvalidDocument
// }

// func setAttrs(n *node, tok *xml.StartElement, attrPrefix string) {
// 	if len(tok.Attr) > 0 {
// 		m := make(map[string]interface{})
// 		for _, attr := range tok.Attr {
// 			if len(attr.Name.Space) > 0 {
// 				m[attrPrefix+attr.Name.Space+":"+attr.Name.Local] = attr.Value
// 			} else {
// 				m[attrPrefix+attr.Name.Local] = attr.Value
// 			}
// 		}
// 		n.Value[tok.Name.Local] = m
// 	}
// }

// func setNodeValue(n *node) {
// 	if v, ok := n.Parent.Value[n.Parent.Label]; ok {
// 		m := v.(map[string]interface{})
// 		if v, ok = m[n.Label]; ok {
// 			switch item := v.(type) {
// 			case string:
// 				Pl("string item: %v, v: %#v", item, v)
// 				m[n.Label] = []string{item, n.Value[n.Label].(string)}
// 			case []string:
// 				m[n.Label] = append(item, n.Value[n.Label].(string))
// 			case map[string]interface{}:
// 				Pl("map[string]interface{} item: %v, v: %#v", item, v)
// 				vm := getMap(n)
// 				if vm != nil {
// 					m[n.Label] = []map[string]interface{}{item, vm}
// 				}
// 			case []map[string]interface{}:
// 				Pl("[]map[string]interface{} item: %v, v: %#v", item, v)
// 				vm := getMap(n)
// 				if vm != nil {
// 					m[n.Label] = append(item, vm)
// 				}
// 			default:
// 				Pl("item: %v, v: %#v", item, v)
// 			}
// 		} else {
// 			m[n.Label] = n.Value[n.Label]
// 		}

// 	} else {
// 		n.Parent.Value[n.Parent.Label] = n.Value[n.Label]
// 	}
// }

// func getMap(node *node) map[string]interface{} {
// 	if v, ok := node.Value[node.Label]; ok {
// 		switch v.(type) {
// 		case string:
// 			return map[string]interface{}{node.Label: v}
// 		case map[string]interface{}:
// 			return node.Value[node.Label].(map[string]interface{})
// 		}
// 	}

// 	return nil
// }

func (pA *TK) FromXML(xmlA string) (interface{}, error) {
	return GetXMLNode(xmlA)
	// decoder := NewXMLDecoder(strings.NewReader(xmlA))
	// result, err := decoder.Decode()

	// if err != nil {
	// 	return nil, err
	// }

	// return result, nil
}

var FromXML = TKX.FromXML

//func (pA *TK) FromXMLX2(xmlA string) interface{} {
//	var result map[string]interface{}
//
//	err := xml.Unmarshal([]byte(xmlA), &result)
//
//	if err != nil {
//		return err
//	}
//
//	return result
//}
//
//var FromXMLX2 = TKX.FromXMLX2

func (pA *TK) FromXMLWithDefault(xmlA string, defaultA interface{}) interface{} {
	// decoder := NewXMLDecoder(strings.NewReader(xmlA))
	// result, err := decoder.Decode()

	// if err != nil {
	// 	return defaultA
	// }

	// return result

	result, err := GetXMLNode(xmlA)

	if err != nil {
		return defaultA
	}

	return result
}

var FromXMLWithDefault = TKX.FromXMLWithDefault

// 事件相关 event related

// SimpleEvent 简捷的事件结构
type SimpleEvent struct {
	Type  string
	Value string
}

// Init 为SimpleEvent初始化数据
func (p *SimpleEvent) Init(typeA string, valueA string) {
	p.Type = typeA
	p.Value = valueA
}

// CreateSimpleEvent 创建一个SimpleEvent对象，并为其赋初值
func (pA *TK) CreateSimpleEvent(typeA string, valueA string) *SimpleEvent {
	p := &SimpleEvent{}

	p.Type = typeA
	p.Value = valueA

	return p
}

var CreateSimpleEvent = TKX.CreateSimpleEvent

// HTML related

func (pA *TK) RemoveHtmlTags(strA string) string {
	reT := regexp.MustCompile("<[^>].*?>")
	rStrT := reT.ReplaceAllString(strA, "")

	rStrT = Replace(rStrT, "\r\n", "\n")
	reT2 := regexp.MustCompile("^\\s*?$")
	rStrT2 := reT2.ReplaceAllString(rStrT, "")
	rStrT2 = Replace(rStrT2, "\n\n", "\n")

	return rStrT2
}

var RemoveHtmlTags = TKX.RemoveHtmlTags

func (pA *TK) RemoveHtmlTagsX(strA string, optionsA ...string) string {
	if Trim(strA) == "" {
		return strA
	}

	rStrT := RegReplace(strA, "<script[^>]*?>[\\S\\s]*?</script>", "")

	rStrT = RegReplace(rStrT, "<[^>]*?>", "")

	rStrT = strings.Replace(rStrT, "&nbsp;", " ", -1)
	rStrT = strings.Replace(rStrT, "&middot;", "·", -1)
	rStrT = Trim(rStrT)

	if IfSwitchExistsWhole(optionsA, "-removeWhiteSpace") {
		rStrT = strings.Replace(rStrT, "\r", "", -1)
		rStrT = strings.Replace(rStrT, "\n", "", -1)
		rStrT = RegReplace(rStrT, "\\s+", " ")
	}

	if IfSwitchExistsWhole(optionsA, "-replaceComma") {
		rStrT = strings.Replace(rStrT, ",", "`", -1)
	}

	return rStrT
}

var RemoveHtmlTagsX = TKX.RemoveHtmlTagsX

func (pA *TK) HTMLToText(htmlA string, optionsA ...string) (result string) {
	defer func() {
		r := recover()
		if r != nil {
			result = htmlA
			return
		}
	}()

	if Trim(htmlA) == "" {
		result = ""
		return
	}

	typeT := GetSwitchWithDefaultValue(optionsA, "-type=", "")

	if typeT == "tx" {
		rs := RegReplace(htmlA, "(?i:[\\S\\s]*?<body[^>]*?>([\\s\\S]*?)</body>[\\S\\s]*)", "$1")
		rs = RegReplace(rs, "(?i:<script[^>]*?>[\\s\\S]*?</script>)", "")
		for RegContains(rs, "<[^>]*?>([\\s\\S]*?)</[^>]*?>") {
			rs = RegReplace(rs, "<[^>]*?>([\\s\\S]*?)</[^>]*?>", "$1")
		}
		result = RegReplace(rs, "<[^>]*?>", "")
		return
	}

	docT, errT := html.Parse(strings.NewReader(htmlA))
	if errT != nil {
		result = htmlA
		return
	} else {
		_, textT, simplified, flattened, cleaned, errT := sandblast.ExtractEx(docT, sandblast.KeepLinks)
		if errT != nil {
			result = htmlA
			return
		} else {
			switch typeT {
			case "", "text":
				result = textT
				return
			case "simple", "s":
				result = simplified.DebugString()
				return
			case "flatten", "f":
				result = flattened.DebugString()
				return
			case "cleaned", "c":
				result = cleaned.DebugString()
				return
			case "beautified", "b":
				result = RemoveHtmlTags(cleaned.DebugString())
				return
			case "x":
				reT := regexp.MustCompile("<[^>]*?>\\[\\d*?\\]")
				var tmpRT string
				if flattened == nil {
					tmpRT = reT.ReplaceAllString(textT, "")
				} else {
					tmpRT = reT.ReplaceAllString(flattened.DebugString(), "")
				}
				reT = regexp.MustCompile("<[^>]*?>\\[(.*?)\\([^)]*?\\)\\]")
				tmpRT = reT.ReplaceAllString(tmpRT, "$1")
				reT = regexp.MustCompile("^\\s*?(\\S)")
				tmpRT = reT.ReplaceAllString(tmpRT, "$1")
				result = tmpRT
				return
			default:
				result = textT
				return
			}
		}
	}
}

var HTMLToText = TKX.HTMLToText

// Markdown Related

type myIDs struct {
}

func (s *myIDs) Generate(value []byte, kind mdast.NodeKind) []byte {
	s1 := string(value)

	// s1 = RegReplaceX(s1, `^- `, "--")

	s1 = RegReplaceX(s1, `\s+`, "-")

	s1 = strings.ToLower(s1)

	s1 = RegReplaceX(s1, `[（）\/]`, "")

	s1 = RegReplaceX(s1, `\*\*([^\*]+?)\*\*`, "${1}")

	// s1 = UrlEncode2(s1)

	// Pl("kind: %v, value: %#v", kind, s1)
	return []byte(s1)
}

func (s *myIDs) Put(value []byte) {
}

func (pA *TK) RenderMarkdown(markdownA string) string {
	ctx := mdparser.NewContext(mdparser.WithIDs(&myIDs{}))
	markdownT := goldmark.New(
		goldmark.WithExtensions(
			mdextension.GFM,
			mdhighlighting.Highlighting,
		),
		goldmark.WithParserOptions(
			mdparser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			mdhtml.WithHardWraps(),
			mdhtml.WithXHTML(),
			mdhtml.WithUnsafe(),
		),
	)

	var buf bytes.Buffer
	if err := markdownT.Convert([]byte(markdownA), &buf, mdparser.WithContext(ctx)); err != nil {
		return ErrToStr(err)
	}

	return string(buf.Bytes())
}

var RenderMarkdown = TKX.RenderMarkdown

// reflect related

func (pA *TK) ReflectHasMethod(vA interface{}, nameA string) (result bool) {
	defer func() {
		r := recover()

		if r != nil {
			result = false
			return
		}
	}()

	var rv1 reflect.Value

	switch nv := vA.(type) {
	case *interface{}:
		// Pl("hereree1")
		rv1 = reflect.ValueOf(GetRef(nv))
	default:
		rv1 = reflect.ValueOf(vA)

	}

	// Pl("rv1: %T %#v %v", rv1, rv1, rv1)

	rv2 := rv1.MethodByName(nameA)

	if rv2.IsZero() {
		result = false
		return
	}

	result = true
	return
}

var ReflectHasMethod = TKX.ReflectHasMethod

func (pA *TK) ReflectCallMethod(vA interface{}, nameA string, argsA ...interface{}) (result interface{}) {
	defer func() {
		r := recover()

		if r != nil {
			result = fmt.Errorf("%v([%T/%v/%v->%v] %v)", r, vA, reflect.TypeOf(vA), reflect.TypeOf(vA).Kind(), nameA, vA)
			return
		}
	}()

	var rv1 reflect.Value

	switch nv := vA.(type) {
	case *interface{}:
		// Pl("hereree1")
		rv1 = reflect.ValueOf(GetRef(nv))
	default:
		rv1 = reflect.ValueOf(vA)

	}

	// Pl("rv1: %T %#v %v", rv1, rv1, rv1)

	rv2 := rv1.MethodByName(nameA)

	if rv2.IsZero() {
		result = fmt.Errorf("unknown method: %v（%T/%v）", nameA, vA, vA)
		return
	}

	// Pl("rv2: %T %#v %v", rv2, rv2, rv2)

	lenT := len(argsA)

	sl := make([]reflect.Value, 0, lenT)

	for i := 0; i < lenT; i++ {
		sl = append(sl, reflect.ValueOf(argsA[i]))
	}

	rrvT := rv2.Call(sl)

	rvr := make([]interface{}, 0)

	for _, v9 := range rrvT {
		rvr = append(rvr, v9.Interface())
	}

	rLenT := len(rvr)

	if rLenT < 1 {
		result = nil
		return
	} else if rLenT > 1 {
		result = rvr
		return
	}

	result = rvr[0]
	return
}

var ReflectCallMethod = TKX.ReflectCallMethod

func (pA *TK) ReflectCallMethodQuick(vA interface{}, nameA string, argsA ...interface{}) []interface{} {
	var rv1 reflect.Value = reflect.ValueOf(vA)

	// Pl("rv1: %T %#v %v", rv1, rv1, rv1)

	rv2 := rv1.MethodByName(nameA)

	if !rv2.IsValid() {
		return nil
	}

	if rv2.IsZero() {
		return nil
	}

	// Pl("rv2: %T %#v %v", rv2, rv2, rv2)

	lenT := len(argsA)

	sl := make([]reflect.Value, 0, lenT)

	for i := 0; i < lenT; i++ {
		sl = append(sl, reflect.ValueOf(argsA[i]))
	}

	rrvT := rv2.Call(sl)

	rvr := make([]interface{}, 0)

	for _, v9 := range rrvT {
		rvr = append(rvr, v9.Interface())
	}

	return rvr
}

var ReflectCallMethodQuick = TKX.ReflectCallMethodQuick

func (pA *TK) ReflectCallMethodCompact(vA interface{}, nameA string, argsA ...interface{}) interface{} {
	var rv1 reflect.Value = reflect.ValueOf(vA)

	// Pl("rv1: %T %#v %v", rv1, rv1, rv1)

	rv2 := rv1.MethodByName(nameA)

	if !rv2.IsValid() || rv2.IsZero() {
		return nil
	}

	// Pl("rv2: %T %#v %v", rv2, rv2, rv2)

	lenT := len(argsA)

	sl := make([]reflect.Value, 0, lenT)

	for i := 0; i < lenT; i++ {
		sl = append(sl, reflect.ValueOf(argsA[i]))
	}

	rrvT := rv2.Call(sl)

	rvr := make([]interface{}, 0)

	for _, v9 := range rrvT {
		rvr = append(rvr, v9.Interface())
	}

	if len(rvr) < 1 {
		return nil
	} else if len(rvr) < 2 {
		return rvr[0]
	}

	return rvr
}

var ReflectCallMethodCompact = TKX.ReflectCallMethodCompact

func (pA *TK) ReflectCallFuncQuick(vA interface{}, argsA ...interface{}) []interface{} {
	var rv1 reflect.Value = reflect.ValueOf(vA)

	lenT := len(argsA)

	sl := make([]reflect.Value, 0, lenT)

	for i := 0; i < lenT; i++ {
		sl = append(sl, reflect.ValueOf(argsA[i]))
	}

	rrvT := rv1.Call(sl)

	rvr := make([]interface{}, 0)

	for _, v9 := range rrvT {
		rvr = append(rvr, v9.Interface())
	}

	return rvr
}

var ReflectCallFuncQuick = TKX.ReflectCallFuncQuick

func (pA *TK) ReflectCallFuncCompact(vA interface{}, argsA ...interface{}) interface{} {
	var rv1 reflect.Value = reflect.ValueOf(vA)

	lenT := len(argsA)

	sl := make([]reflect.Value, 0, lenT)

	for i := 0; i < lenT; i++ {
		sl = append(sl, reflect.ValueOf(argsA[i]))
	}

	rrvT := rv1.Call(sl)

	rvr := make([]interface{}, 0)

	for _, v9 := range rrvT {
		rvr = append(rvr, v9.Interface())
	}

	if len(rvr) < 1 {
		return nil
	} else if len(rvr) < 2 {
		return rvr[0]
	}

	return rvr
}

var ReflectCallFuncCompact = TKX.ReflectCallFuncCompact

func (pA *TK) ReflectGetMember(vA interface{}, argsA ...interface{}) (result interface{}) {
	defer func() {
		r := recover()

		if r != nil {
			result = fmt.Errorf("%v([%T/%v/%v->%v] %v)", r, vA, reflect.TypeOf(vA), reflect.TypeOf(vA).Kind(), argsA, vA)
			return
		}
	}()

	var vr interface{} = vA

	for _, v2 := range argsA {

		typeT := reflect.TypeOf(vr)

		kindT := typeT.Kind()

		// tk.Pl("kind: %v", kindT)

		if kindT == reflect.Ptr {
			vrf := reflect.ValueOf(vr).Elem()

			kindT = vrf.Kind()

			// tk.Pl("vrf: %v, kind: %v", vrf, kindT)

			if kindT == reflect.Struct {
				vr = vrf.Interface()
			}
		}

		if kindT == reflect.Struct {
			rv1 := reflect.ValueOf(vr)

			rvMidT := rv1.FieldByName(ToStr(v2))

			if !rvMidT.IsValid() || rvMidT.IsZero() {
				// rvMidT = rv1.MethodByName(ToStr(v2))

				// if !rvMidT.IsValid() || rvMidT.IsZero() {
				result = fmt.Errorf("nil or unknown member: %v（%T/%v）.%v", vr, vr, kindT, v2)
				return
				// }
			}

			rv2 := rvMidT.Interface()

			vr = rv2
			continue
		}

		result = fmt.Errorf("unknown member: %v（%T/%v）.%v", vr, vr, kindT, v2)
		return
	}

	result = vr
	return
}

var ReflectGetMember = TKX.ReflectGetMember

func (pA *TK) ReflectGetMethod(vA interface{}, argsA ...interface{}) (result interface{}) {
	defer func() {
		r := recover()

		if r != nil {
			result = fmt.Errorf("%v([%T/%v/%v->%v] %v)", r, vA, reflect.TypeOf(vA), reflect.TypeOf(vA).Kind(), argsA, vA)
			return
		}
	}()

	var vr interface{} = vA

	for _, v2 := range argsA {

		typeT := reflect.TypeOf(vr)

		kindT := typeT.Kind()

		// tk.Pl("kind: %v", kindT)

		if kindT == reflect.Ptr {
			vrf := reflect.ValueOf(vr).Elem()

			kindT = vrf.Kind()

			// tk.Pl("vrf: %v, kind: %v", vrf, kindT)

			if kindT == reflect.Struct {
				vr = vrf.Interface()
			}
		}

		if kindT == reflect.Struct {
			rv1 := reflect.ValueOf(vr)

			rvMidT := rv1.MethodByName(ToStr(v2))

			if !rvMidT.IsValid() || rvMidT.IsZero() {
				result = fmt.Errorf("nil or unknown method: %v（%T/%v）.%v", vr, vr, kindT, v2)
				return
			}

			rv2 := rvMidT.Interface()

			vr = rv2
			continue
		}

		result = fmt.Errorf("unknown method: %v（%T/%v）.%v", vr, vr, kindT, v2)
		return
	}

	result = vr
	return
}

var ReflectGetMethod = TKX.ReflectGetMethod

func (pA *TK) ReflectCallMethodSlice(vA interface{}, nameA string, argsA ...interface{}) (result interface{}) {
	defer func() {
		r := recover()

		if r != nil {
			result = fmt.Errorf("%v([%T/%v/%v->%v] %v)", r, vA, reflect.TypeOf(vA), reflect.TypeOf(vA).Kind(), nameA, vA)
			return
		}
	}()

	rv1 := reflect.ValueOf(vA)

	rv2 := rv1.MethodByName(nameA)

	lenT := len(argsA)

	sl := make([]reflect.Value, 0, lenT)

	for i := 0; i < lenT-1; i++ {
		sl = append(sl, reflect.ValueOf(argsA[i]))
	}

	if lenT > 0 {
		tmpSl := make([]reflect.Value, 0, lenT)
		for i := lenT - 1; i < lenT; i++ {
			tmpSl = append(tmpSl, reflect.ValueOf(argsA[i]))
		}
		sl = append(sl, tmpSl...)
	}

	rrvT := rv2.CallSlice(sl)

	rvr := make([]interface{}, 0)

	for _, v9 := range rrvT {
		rvr = append(rvr, v9.Interface())
	}

	rLenT := len(rvr)

	if rLenT < 1 {
		result = nil
		return
	} else if rLenT > 1 {
		result = rvr
		return
	}

	result = rvr[0]
	return
}

var ReflectCallMethodSlice = TKX.ReflectCallMethodSlice

// Hans Related

func (pA *TK) ToPinYin(strA string, argsA ...string) interface{} {
	pinyin1 := pinyin.NewArgs()

	if IfSwitchExistsWhole(argsA, "-first") {
		pinyin1.Style = pinyin.FirstLetter
	}

	if IfSwitchExistsWhole(argsA, "-tone") {
		pinyin1.Style = pinyin.Tone
	}

	if IfSwitchExistsWhole(argsA, "-digitTone") {
		pinyin1.Style = pinyin.Tone3
	}

	if IfSwitchExistsWhole(argsA, "-digitTone2") {
		pinyin1.Style = pinyin.Tone2
	}

	if IfSwitchExistsWhole(argsA, "-ascOnly") {
		pinyin1.Fallback = func(r rune, a pinyin.Args) []string {
			if r > 255 {
				return []string{}
			}

			return []string{string(r)}
		}
	} else if IfSwitchExistsWhole(argsA, "-pinYinOnly") {
	} else {
		pinyin1.Fallback = func(r rune, a pinyin.Args) []string {
			return []string{string(r)}
		}
	}

	var rs [][]string = pinyin.Pinyin(strA, pinyin1)

	sepT := GetSwitch(argsA, "-sep=", "")

	if IfSwitchExistsWhole(argsA, "-raw") {
		return rs
	}

	sary1 := make([]string, 0, len(rs))

	for _, v := range rs {
		lineStrT := strings.Join(v, sepT)
		sary1 = append(sary1, lineStrT)
	}

	return strings.Join(sary1, sepT)
}

var ToPinYin = TKX.ToPinYin

// Misc Related

func (pA *TK) Pass() {

}

var Pass = TKX.Pass

func (pA *TK) WrapError(vA interface{}, errA error) interface{} {
	if errA == nil {
		return vA
	}

	return errA
}

var WrapError = TKX.WrapError

func (pA *TK) StatusToString(statusA string) string {
	var mapT map[string]string

	errT := jsoniter.Unmarshal([]byte(statusA), &mapT)

	if errT != nil {
		return "TXERROR:failed to parse JSON: " + errT.Error()
	}

	statusT, ok := mapT["Status"]

	if !ok {
		return "TXERROR:invalid data format"
	}

	valueT, ok := mapT["Value"]

	if !ok {
		return "TXERROR:invalid data format"
	}

	if statusT == "success" {
		return valueT
	}

	return "TXERROR:" + valueT
}

var StatusToString = TKX.StatusToString

func (pA *TK) StatusToMap(statusA string) interface{} {
	var mapT map[string]string

	errT := jsoniter.Unmarshal([]byte(statusA), &mapT)

	if errT != nil {
		return errT
	}

	return mapT
}

var StatusToMap = TKX.StatusToMap

// 10: self|20:large|30:medium -> {"10": " self", "20": "large", "30": "medium"}
func (pA *TK) SimpleStrToMap(strA string) map[string]string {
	var mapT map[string]string = make(map[string]string)

	listT := strings.Split(strA, "|")
	
	for _, v := range listT {
		list1T := strings.SplitN(v, ":", 2)
		
		if len(list1T) < 2 {
			continue
		}
		
//		k := strings.TrimSpace(list1T[0])
//		v := strings.TrimSpace(list1T[1])
		
		mapT[list1T[0]] = list1T[1]
	}

	return mapT
}

var SimpleStrToMap = TKX.SimpleStrToMap

// 10: self|20:large|30:medium -> {"10": " self", "20": "large", "30": "medium"}
func (pA *TK) ReverseStrMap(mapA map[string]string) map[string]string {
	var mapT map[string]string = make(map[string]string, len(mapA))

	for k, v := range mapA {
		mapT[v] = k
	}

	return mapT
}

var ReverseStrMap = TKX.ReverseStrMap

func (pA *TK) StringsToJson(strsA ...string) string {
	lenT := len(strsA)

	if lenT < 2 {
		return "{}"
	}

	len2T := lenT / 2

	mapT := map[string]string{}

	for i := 0; i < len2T; i++ {
		mapT[strsA[i*2]] = strsA[i*2+1]
	}

	rs, errT := jsoniter.Marshal(mapT)

	if errT != nil {
		return ErrorToString(errT)
	}

	return string(rs)
}

var StringsToJson = TKX.StringsToJson

func (pA *TK) DealRef(refA *interface{}, cmdA string, argsA ...interface{}) interface{} {
	vT := *refA

	if cmdA == "++" {
		switch nv := vT.(type) {
		case int:
			rs := nv + 1
			*refA = rs
			return rs
		case int64:
			rs := nv + 1
			*refA = rs
			return rs
		case byte:
			rs := nv + 1
			*refA = rs
			return rs
		}

	} else if cmdA == "--" {
		switch nv := vT.(type) {
		case int:
			rs := nv - 1
			*refA = rs
			return rs
		case int64:
			rs := nv - 1
			*refA = rs
			return rs
		case byte:
			rs := nv - 1
			*refA = rs
			return rs
		}

	}

	return fmt.Errorf("unsupported type action: %T.%v", vT, cmdA)
}

var DealRef = TKX.DealRef

func (pA *TK) ObjectsToS(aryA ...interface{}) []string {
	if aryA == nil {
		return nil
	}

	rs := make([]string, 0, len(aryA))

	for _, v := range aryA {
		rs = append(rs, ToStr(v))
	}

	return rs
}

var ObjectsToS = TKX.ObjectsToS

type SortStruct struct {
	Value    interface{}
	Key      string
	Desc     bool
	LenFunc  (func() int)
	LessFunc (func(i, j int) bool)
	SwapFunc (func(i, j int))
}

func (p *SortStruct) Len() int {
	if p.LenFunc != nil {
		return p.LenFunc()
	}

	if p == nil {
		return 0
	}

	switch nv := p.Value.(type) {
	case int, int8, int16, int32, int64:
		return 0
	case uint, uint8, uint16, uint32, uint64:
		return 0
	case float32, float64:
		return 0
	case complex64, complex128:
		return 0
	case bool:
		return 0
	case string:
		return len(nv)
	case []int:
		return len(nv)
	case []int8:
		return len(nv)
	case []int16:
		return len(nv)
	case []int32:
		return len(nv)
	case []int64:
		return len(nv)
	case []uint:
		return len(nv)
	case []uint8:
		return len(nv)
	case []uint16:
		return len(nv)
	case []uint32:
		return len(nv)
	case []uint64:
		return len(nv)
	case []float64:
		return len(nv)
	case []string:
		return len(nv)
	case []interface{}:
		return len(nv)
	case []map[string]string:
		return len(nv)
	case []map[string]interface{}:
		return len(nv)
	}

	return 0
}

func (pA *TK) LessI(v1, v2 interface{}, optsA ...interface{}) bool {
	// Pl("%#v, %#v", v1, v2)
	if v1 == nil {
		if v2 == nil {
			return false
		}

		return true
	} else if v2 == nil {
		return false
	}

	switch nv1 := v1.(type) {
	case bool:
		nv1i := BoolToInt(nv1)
		switch nv2 := v2.(type) {
		case bool:
			return nv1i < BoolToInt(nv2)
		case int:
			return nv1i < nv2
		case int8:
			return nv1i < int(nv2)
		case int16:
			return nv1i < int(nv2)
		case int32:
			return nv1i < int(nv2)
		case int64:
			return nv1i < int(nv2)
		case uint:
			return nv1i < int(nv2)
		case uint8:
			return nv1i < int(nv2)
		case uint16:
			return nv1i < int(nv2)
		case uint32:
			return nv1i < int(nv2)
		case uint64:
			return nv1i < int(nv2)
		case float32:
			return float64(nv1i) < float64(nv2)
		case float64:
			return float64(nv1i) < nv2
		case string:
			return float64(nv1i) < ToFloat(nv2)
		}
	case int:
		switch nv2 := v2.(type) {
		case bool:
			return nv1 < BoolToInt(nv2)
		case int:
			return nv1 < nv2
		case int8:
			return nv1 < int(nv2)
		case int16:
			return nv1 < int(nv2)
		case int32:
			return nv1 < int(nv2)
		case int64:
			return nv1 < int(nv2)
		case uint:
			return nv1 < int(nv2)
		case uint8:
			return nv1 < int(nv2)
		case uint16:
			return nv1 < int(nv2)
		case uint32:
			return nv1 < int(nv2)
		case uint64:
			return nv1 < int(nv2)
		case float32:
			return float64(nv1) < float64(nv2)
		case float64:
			return float64(nv1) < nv2
		case string:
			return float64(nv1) < ToFloat(nv2)
		}
	case int8:
		switch nv2 := v2.(type) {
		case bool:
			return int(nv1) < BoolToInt(nv2)
		case int:
			return int(nv1) < nv2
		case int8:
			return nv1 < nv2
		case int16:
			return int(nv1) < int(nv2)
		case int32:
			return int(nv1) < int(nv2)
		case int64:
			return int(nv1) < int(nv2)
		case uint:
			return int(nv1) < int(nv2)
		case uint8:
			return int(nv1) < int(nv2)
		case uint16:
			return int(nv1) < int(nv2)
		case uint32:
			return int(nv1) < int(nv2)
		case uint64:
			return int(nv1) < int(nv2)
		case float32:
			return float64(nv1) < float64(nv2)
		case float64:
			return float64(nv1) < nv2
		case string:
			return float64(nv1) < ToFloat(nv2)
		}
	case int16:
		switch nv2 := v2.(type) {
		case bool:
			return int(nv1) < BoolToInt(nv2)
		case int:
			return int(nv1) < nv2
		case int8:
			return int(nv1) < int(nv2)
		case int16:
			return int(nv1) < int(nv2)
		case int32:
			return int(nv1) < int(nv2)
		case int64:
			return int(nv1) < int(nv2)
		case uint:
			return int(nv1) < int(nv2)
		case uint8:
			return int(nv1) < int(nv2)
		case uint16:
			return int(nv1) < int(nv2)
		case uint32:
			return int(nv1) < int(nv2)
		case uint64:
			return int(nv1) < int(nv2)
		case float32:
			return float64(nv1) < float64(nv2)
		case float64:
			return float64(nv1) < nv2
		case string:
			return float64(nv1) < ToFloat(nv2)
		}
	case int32:
		switch nv2 := v2.(type) {
		case bool:
			return int(nv1) < BoolToInt(nv2)
		case int:
			return int(nv1) < nv2
		case int8:
			return int(nv1) < int(nv2)
		case int16:
			return int(nv1) < int(nv2)
		case int32:
			return int(nv1) < int(nv2)
		case int64:
			return int(nv1) < int(nv2)
		case uint:
			return int(nv1) < int(nv2)
		case uint8:
			return int(nv1) < int(nv2)
		case uint16:
			return int(nv1) < int(nv2)
		case uint32:
			return int(nv1) < int(nv2)
		case uint64:
			return int(nv1) < int(nv2)
		case float32:
			return float64(nv1) < float64(nv2)
		case float64:
			return float64(nv1) < nv2
		case string:
			return float64(nv1) < ToFloat(nv2)
		}
	case int64:
		switch nv2 := v2.(type) {
		case bool:
			return int(nv1) < BoolToInt(nv2)
		case int:
			return int(nv1) < nv2
		case int8:
			return int(nv1) < int(nv2)
		case int16:
			return int(nv1) < int(nv2)
		case int32:
			return int(nv1) < int(nv2)
		case int64:
			return int(nv1) < int(nv2)
		case uint:
			return int(nv1) < int(nv2)
		case uint8:
			return int(nv1) < int(nv2)
		case uint16:
			return int(nv1) < int(nv2)
		case uint32:
			return int(nv1) < int(nv2)
		case uint64:
			return int(nv1) < int(nv2)
		case float32:
			return float64(nv1) < float64(nv2)
		case float64:
			return float64(nv1) < nv2
		case string:
			return float64(nv1) < ToFloat(nv2)
		}
	case uint:
		switch nv2 := v2.(type) {
		case bool:
			return int(nv1) < BoolToInt(nv2)
		case int:
			return int(nv1) < nv2
		case int8:
			return int(nv1) < int(nv2)
		case int16:
			return int(nv1) < int(nv2)
		case int32:
			return int(nv1) < int(nv2)
		case int64:
			return int(nv1) < int(nv2)
		case uint:
			return int(nv1) < int(nv2)
		case uint8:
			return int(nv1) < int(nv2)
		case uint16:
			return int(nv1) < int(nv2)
		case uint32:
			return int(nv1) < int(nv2)
		case uint64:
			return int(nv1) < int(nv2)
		case float32:
			return float64(nv1) < float64(nv2)
		case float64:
			return float64(nv1) < nv2
		case string:
			return float64(nv1) < ToFloat(nv2)
		}
	case uint8:
		switch nv2 := v2.(type) {
		case bool:
			return int(nv1) < BoolToInt(nv2)
		case int:
			return int(nv1) < nv2
		case int8:
			return int(nv1) < int(nv2)
		case int16:
			return int(nv1) < int(nv2)
		case int32:
			return int(nv1) < int(nv2)
		case int64:
			return int(nv1) < int(nv2)
		case uint:
			return int(nv1) < int(nv2)
		case uint8:
			return int(nv1) < int(nv2)
		case uint16:
			return int(nv1) < int(nv2)
		case uint32:
			return int(nv1) < int(nv2)
		case uint64:
			return int(nv1) < int(nv2)
		case float32:
			return float64(nv1) < float64(nv2)
		case float64:
			return float64(nv1) < nv2
		case string:
			return float64(nv1) < ToFloat(nv2)
		}
	case uint16:
		switch nv2 := v2.(type) {
		case bool:
			return int(nv1) < BoolToInt(nv2)
		case int:
			return int(nv1) < nv2
		case int8:
			return int(nv1) < int(nv2)
		case int16:
			return int(nv1) < int(nv2)
		case int32:
			return int(nv1) < int(nv2)
		case int64:
			return int(nv1) < int(nv2)
		case uint:
			return int(nv1) < int(nv2)
		case uint8:
			return int(nv1) < int(nv2)
		case uint16:
			return int(nv1) < int(nv2)
		case uint32:
			return int(nv1) < int(nv2)
		case uint64:
			return int(nv1) < int(nv2)
		case float32:
			return float64(nv1) < float64(nv2)
		case float64:
			return float64(nv1) < nv2
		case string:
			return float64(nv1) < ToFloat(nv2)
		}
	case uint32:
		switch nv2 := v2.(type) {
		case bool:
			return int(nv1) < BoolToInt(nv2)
		case int:
			return int(nv1) < nv2
		case int8:
			return int(nv1) < int(nv2)
		case int16:
			return int(nv1) < int(nv2)
		case int32:
			return int(nv1) < int(nv2)
		case int64:
			return int(nv1) < int(nv2)
		case uint:
			return int(nv1) < int(nv2)
		case uint8:
			return int(nv1) < int(nv2)
		case uint16:
			return int(nv1) < int(nv2)
		case uint32:
			return int(nv1) < int(nv2)
		case uint64:
			return int(nv1) < int(nv2)
		case float32:
			return float64(nv1) < float64(nv2)
		case float64:
			return float64(nv1) < nv2
		case string:
			return float64(nv1) < ToFloat(nv2)
		}
	case uint64:
		switch nv2 := v2.(type) {
		case bool:
			return int(nv1) < BoolToInt(nv2)
		case int:
			return int(nv1) < nv2
		case int8:
			return int(nv1) < int(nv2)
		case int16:
			return int(nv1) < int(nv2)
		case int32:
			return int(nv1) < int(nv2)
		case int64:
			return int(nv1) < int(nv2)
		case uint:
			return int(nv1) < int(nv2)
		case uint8:
			return int(nv1) < int(nv2)
		case uint16:
			return int(nv1) < int(nv2)
		case uint32:
			return int(nv1) < int(nv2)
		case uint64:
			return int(nv1) < int(nv2)
		case float32:
			return float64(nv1) < float64(nv2)
		case float64:
			return float64(nv1) < nv2
		case string:
			return float64(nv1) < ToFloat(nv2)
		}
	case float32:
		switch nv2 := v2.(type) {
		case bool:
			return float64(nv1) < float64(BoolToInt(nv2))
		case int:
			return float64(nv1) < float64(nv2)
		case int8:
			return float64(nv1) < float64(nv2)
		case int16:
			return float64(nv1) < float64(nv2)
		case int32:
			return float64(nv1) < float64(nv2)
		case int64:
			return float64(nv1) < float64(nv2)
		case uint:
			return float64(nv1) < float64(nv2)
		case uint8:
			return float64(nv1) < float64(nv2)
		case uint16:
			return float64(nv1) < float64(nv2)
		case uint32:
			return float64(nv1) < float64(nv2)
		case uint64:
			return float64(nv1) < float64(nv2)
		case float32:
			return nv1 < nv2
		case float64:
			return float64(nv1) < nv2
		case string:
			return float64(nv1) < ToFloat(nv2)
		}
	case float64:
		switch nv2 := v2.(type) {
		case bool:
			return nv1 < float64(BoolToInt(nv2))
		case int:
			return nv1 < float64(nv2)
		case int8:
			return nv1 < float64(nv2)
		case int16:
			return nv1 < float64(nv2)
		case int32:
			return nv1 < float64(nv2)
		case int64:
			return nv1 < float64(nv2)
		case uint:
			return nv1 < float64(nv2)
		case uint8:
			return nv1 < float64(nv2)
		case uint16:
			return nv1 < float64(nv2)
		case uint32:
			return nv1 < float64(nv2)
		case uint64:
			return nv1 < float64(nv2)
		case float32:
			return nv1 < float64(nv2)
		case float64:
			return nv1 < nv2
		case string:
			return nv1 < ToFloat(nv2)
		}
	case string:
		switch nv2 := v2.(type) {
		case bool:
			return ToFloat(nv1) < float64(BoolToInt(nv2))
		case int:
			return ToFloat(nv1) < float64(nv2)
		case int8:
			return ToFloat(nv1) < float64(nv2)
		case int16:
			return ToFloat(nv1) < float64(nv2)
		case int32:
			return ToFloat(nv1) < float64(nv2)
		case int64:
			return ToFloat(nv1) < float64(nv2)
		case uint:
			return ToFloat(nv1) < float64(nv2)
		case uint8:
			return ToFloat(nv1) < float64(nv2)
		case uint16:
			return ToFloat(nv1) < float64(nv2)
		case uint32:
			return ToFloat(nv1) < float64(nv2)
		case uint64:
			return ToFloat(nv1) < float64(nv2)
		case float32:
			return ToFloat(nv1) < float64(nv2)
		case float64:
			return ToFloat(nv1) < nv2
		case string:
			return nv1 < nv2
		}
	case map[string]string:
		nv2, ok := v2.(map[string]string)

		if ok {
			keyT := GetSwitchI(optsA, "-key=", "")
			if keyT != "" {
				return nv1[keyT] < nv2[keyT]
			}
		}
	case map[string]interface{}:
		nv2, ok := v2.(map[string]interface{})

		if ok {
			keyT := GetSwitchI(optsA, "-key=", "")
			if keyT != "" {
				return TKX.LessI(nv1[keyT], nv2[keyT])
			}
		}
	case []interface{}:
		nv2, ok := v2.([]interface{})

		if ok {
			keyT := ToInt(GetSwitchI(optsA, "-key=", ""))
			if keyT >= 0 {
				return TKX.LessI(nv1[keyT], nv2[keyT])
			}
		}
		// case interface{}:
		// 	nv2, ok := v2.(interface{})

		// 	if ok {
		// 		keyT := GetSwitchI(optsA, "-key=", "")
		// 		if keyT != "" {
		// 			return TKX.LessI(nv1[keyT], nv2[keyT])
		// 		}
		// 	}
		// default:
		// 	Pl("v2: (%T) %v", v2, v2)

	}

	return false

	// type1 := fmt.Sprintf("%T", v1)
	// type2 := fmt.Sprintf("%T", v2)

	// if type1 == type2 {
	// 	if StartsWith(type1, "bool") {
	// 		return v1.(bool) < v2.(bool)
	// 	} else if StartsWith(type1, "int") {
	// 		return ToInt(v1) < ToInt(v2)
	// 	} else if StartsWith(type1, "float") {
	// 		return ToFloat(v1) < ToFloat(v2)
	// 	} else if StartsWith(type1, "string") {
	// 		return ToFloat(v1) < ToFloat(v2)
	// 	}
	// }

	// if type1 == "string" {
	// 	if type2 == "string" {
	// 		return v1.(string) < v2.(string)
	// 	}
	// }

	// return false
}

var LessI = TKX.LessI

func (p *SortStruct) Less(i, j int) bool {
	if p.LessFunc != nil {
		return p.LessFunc(i, j)
	}

	if p == nil {
		return false
	}

	switch nv := p.Value.(type) {
	case int, int8, int16, int32, int64:
		return false
	case uint, uint8, uint16, uint32, uint64:
		return false
	case float32, float64:
		return false
	case complex64, complex128:
		return false
	case bool:
		return false
	case string:
		return false
	case []int:
		return BitXor(p.Desc, nv[i] < nv[j]).(bool)
	case []int8:
		return BitXor(p.Desc, nv[i] < nv[j]).(bool)
	case []int16:
		return BitXor(p.Desc, nv[i] < nv[j]).(bool)
	case []int32:
		return BitXor(p.Desc, nv[i] < nv[j]).(bool)
	case []int64:
		return BitXor(p.Desc, nv[i] < nv[j]).(bool)
	case []uint:
		return BitXor(p.Desc, nv[i] < nv[j]).(bool)
	case []uint8:
		return BitXor(p.Desc, nv[i] < nv[j]).(bool)
	case []uint16:
		return BitXor(p.Desc, nv[i] < nv[j]).(bool)
	case []uint32:
		return BitXor(p.Desc, nv[i] < nv[j]).(bool)
	case []uint64:
		return BitXor(p.Desc, nv[i] < nv[j]).(bool)
	case []float64:
		return BitXor(p.Desc, nv[i] < nv[j]).(bool)
	case []string:
		return BitXor(p.Desc, nv[i] < nv[j]).(bool)
	case []interface{}:
		return BitXor(p.Desc, LessI(nv[i], nv[j], "-key="+p.Key)).(bool)
	case []map[string]string:
		return BitXor(p.Desc, LessI(nv[i], nv[j], "-key="+p.Key)).(bool)
	case []map[string]interface{}:
		return BitXor(p.Desc, LessI(nv[i], nv[j], "-key="+p.Key)).(bool)
	}

	return false
}

func (p *SortStruct) Swap(i, j int) {
	if p.SwapFunc != nil {
		p.SwapFunc(i, j)
		return
	}

	if p == nil {
		return
	}

	switch nv := p.Value.(type) {
	case int, int8, int16, int32, int64:
		return
	case uint, uint8, uint16, uint32, uint64:
		return
	case float32, float64:
		return
	case complex64, complex128:
		return
	case bool:
		return
	case string:
		return
	case []int:
		nv[i], nv[j] = nv[j], nv[i]
		return
	case []int8:
		nv[i], nv[j] = nv[j], nv[i]
		return
	case []int16:
		nv[i], nv[j] = nv[j], nv[i]
		return
	case []int32:
		nv[i], nv[j] = nv[j], nv[i]
		return
	case []int64:
		nv[i], nv[j] = nv[j], nv[i]
		return
	case []uint:
		nv[i], nv[j] = nv[j], nv[i]
		return
	case []uint8:
		nv[i], nv[j] = nv[j], nv[i]
		return
	case []uint16:
		nv[i], nv[j] = nv[j], nv[i]
		return
	case []uint32:
		nv[i], nv[j] = nv[j], nv[i]
		return
	case []uint64:
		nv[i], nv[j] = nv[j], nv[i]
		return
	case []float64:
		nv[i], nv[j] = nv[j], nv[i]
		return
	case []string:
		nv[i], nv[j] = nv[j], nv[i]
		return
	case []interface{}:
		nv[i], nv[j] = nv[j], nv[i]
		return
	case []map[string]string:
		nv[i], nv[j] = nv[j], nv[i]
		return
	case []map[string]interface{}:
		nv[i], nv[j] = nv[j], nv[i]
		return
	}

	return
}

func (pA *TK) Sort(vA interface{}, optsA ...interface{}) interface{} {
	if vA == nil {
		return vA
	}

	keyT := GetSwitchI(optsA, "-key=", "")
	ifDescT := IfSwitchExistsWholeI(optsA, "-desc")

	var lenFuncT (func() int) = nil
	var lessFuncT (func(i, j int) bool) = nil
	var swapFuncT (func(i, j int)) = nil

	for _, v := range optsA {
		nv1, ok := v.((func() int))
		if ok {
			lenFuncT = nv1
		}

		nv2, ok := v.((func(i, j int) bool))
		if ok {
			lessFuncT = nv2
		}

		nv3, ok := v.((func(i, j int)))
		if ok {
			swapFuncT = nv3
		}
	}

	// switch nv := vA.(type) {
	// case []map[string]string:
	// 	SortMSSArray(nv, keyT, ifDescT)
	// 	return nv
	// 	// case []map[string]interface{}:
	// 	// 	SortMSIArray(nv, keyT, ifDescT)
	// 	// 	return nv
	// }

	sortStructT := &SortStruct{
		Value:    vA,
		Key:      keyT,
		Desc:     ifDescT,
		LenFunc:  lenFuncT,
		LessFunc: lessFuncT,
		SwapFunc: swapFuncT,
	}

	sort.Sort(sortStructT)

	return sortStructT.Value
}

var SortX = TKX.Sort

func (pA *TK) IsNil(v interface{}) bool {
	if v == nil {
		return true
	}

	_, ok := v.(string)

	if ok {
		return false
	}

	tmps := fmt.Sprintf("%v", v)

	if tmps == "<nil>" {
		return true
	}

	// tmps = fmt.Sprintf("%#v", v)
	// Pl("tmps:%v", tmps)
	// if EndsWith(tmps, "(nil)") {
	// 	return true
	// }

	return false
}

var IsNil = TKX.IsNil

func (pA *TK) IsNilOrEmpty(v interface{}) bool {
	if v == nil {
		return true
	}

	switch v.(type) {
	case nil:
		return true
	case string:
		if v.(string) == "" {
			return true
		}
	case []string:
		if len(v.([]string)) < 1 {
			return true
		}
	default:
		tmps := fmt.Sprintf("%v", v)

		if tmps == "<nil>" {
			return true
		}

	}

	return false
}

var IsNilOrEmpty = TKX.IsNilOrEmpty

func (pA *TK) IsNilOrErrX(v interface{}) bool {
	if v == nil {
		return true
	}

	switch nv := v.(type) {
	case nil:
		return true
	case string:
		if strings.HasPrefix(nv, "TXERROR:") {
			return true
		}
	case error:
		if nv != nil {
			return true
		}
	case *UndefinedStruct:
		return true
	case UndefinedStruct:
		return true
	default:
		tmps := fmt.Sprintf("%v", v)

		if tmps == "<nil>" {
			return true
		}

	}

	return false
}

var IsNilOrErrX = TKX.IsNilOrErrX

func (pA *TK) TrimSafely(vA interface{}, defaultA ...string) string {
	var defaultT string = ""
	if len(defaultA) > 0 {
		defaultT = defaultA[0]
	}

	switch nv := vA.(type) {
	case string:
		return strings.TrimSpace(nv)
	}

	return defaultT
}

var TrimSafely = TKX.TrimSafely

func (pA *TK) TrimEx(vA interface{}, argsA ...string) string {
	if vA == nil {
		return ""
	}

	if vA == Undefined {
		return ""
	}

	if nv, ok := vA.(string); ok {
		return Trim(nv, argsA...)
	}

	return Trim(fmt.Sprintf("%v", vA), argsA...)
}

var TrimEx = TKX.TrimEx

func (pA *TK) IsError(vA interface{}) bool {
	if vA == nil {
		return false
	}

	_, ok := vA.(error)
	if ok {
		return true
	}

	return false
}

var IsError = TKX.IsError

func (pA *TK) IsErrX(vA interface{}) bool {
	if vA == nil {
		return false
	}

	_, ok := vA.(error)

	if ok {
		return true
	}

	nv2, ok := vA.(string)

	if ok {
		return IsErrStr(nv2)
	}

	return false
}

var IsErrX = TKX.IsErrX

func (pA *TK) GetErrStrX(vA interface{}) string {
	if vA == nil {
		return ""
	}

	nv1, ok := vA.(error)

	if ok {
		if nv1 == nil {
			return ""
		}

		return nv1.Error()
	}

	nv2, ok := vA.(string)

	if ok {
		if IsErrStr(nv2) {
			return GetErrStr(nv2)
		} else {
			return nv2
		}
	}

	return fmt.Sprintf("%v", vA)
}

var GetErrStrX = TKX.GetErrStrX

func (pA *TK) IfThenElse(condA bool, thenA interface{}, elseA interface{}) interface{} {
	if condA {
		return thenA
	}

	return elseA
}

var IfThenElse = TKX.IfThenElse

func (pA *TK) GenerateToken(appCodeA string, userIDA string, roleA string, optsA ...string) string {
	secretT := GetSwitch(optsA, "-secret=", "is_Token")

	if secretT == "" {
		secretT = "is_Token"
	}

	strT := appCodeA + "|" + userIDA + "|" + roleA + "|" + GetNowTimeString()

	return EncryptStringByTXDEF(strT, secretT)
}

var GenerateToken = TKX.GenerateToken

func (pA *TK) CheckToken(tokenA string, optsA ...string) string {
	appT := GetSwitch(optsA, "-app=", "")

	secretT := GetSwitch(optsA, "-secret=", "is_Token")

	tokenExpireT := time.Duration(ToInt(GetSwitch(optsA, "-expire=", "1440")))

	strT := DecryptStringByTXDEF(tokenA, secretT)

	if StartsWith(tokenA, "RAW:") {
		strT = tokenA[4:]
	}

	listT := Split(strT, "|")

	if len(listT) < 4 {
		return ErrStrf("invalid token(len)")
	}

	if appT != "" {
		if appT != listT[0] {
			return ErrStrf("invalid token(app)")
		}

	}

	timeT, errT := StrToTimeByFormat(listT[3], TimeFormatCompact)
	if errT != nil {
		return ErrStrf("invalid token(time)")
	}

	expectTimeT := timeT.Add(time.Minute * time.Duration(tokenExpireT))

	if expectTimeT.Before(time.Now()) {
		return ErrStrf("token expired")
	}

	return listT[0] + "|" + listT[1] + "|" + listT[2] + "|" + ""
}

var CheckToken = TKX.CheckToken

// GenerateQR default -level=1, 4 indicates more error tolerance
func (pA *TK) GenerateQR(contentA string, optsA ...string) (barcode.Barcode, error) {
	var qrCode barcode.Barcode

	levelT := GetSwitchWithDefaultIntValue(optsA, "-level=", 1)

	switch levelT {
	case 0:
		qrCode, _ = qr.Encode(contentA, qr.L, qr.Auto)
	case 1:
		qrCode, _ = qr.Encode(contentA, qr.M, qr.Auto)
	case 2:
		qrCode, _ = qr.Encode(contentA, qr.Q, qr.Auto)
	case 3:
		qrCode, _ = qr.Encode(contentA, qr.H, qr.Auto)
	default:
		qrCode, _ = qr.Encode(contentA, qr.M, qr.Auto)
	}

	qrCode, _ = barcode.Scale(qrCode, GetSwitchWithDefaultIntValue(optsA, "-width=", 500), GetSwitchWithDefaultIntValue(optsA, "-height=", 500))

	return qrCode, nil
}

var GenerateQR = TKX.GenerateQR

func (pA *TK) Compress(dataA interface{}, argsA ...interface{}) interface{} {
	//NoCompression      = flate.NoCompression
	//BestSpeed          = flate.BestSpeed
	//BestCompression    = flate.BestCompression
	//DefaultCompression = flate.DefaultCompression
	//gzip.NewWriterLevel()

	var inputBufT []byte

	switch nv := dataA.(type) {
	case []byte:
		inputBufT = nv
	case string:
		inputBufT = []byte(nv)
	default:
		return fmt.Errorf("failed to compress, unsupported type: %T(%v)", dataA, dataA)
	}

	var buf bytes.Buffer

	gzipWriter := gzip.NewWriter(&buf)

	_, err := gzipWriter.Write(inputBufT)
	if err != nil {
		gzipWriter.Close()
		return err
	}

	errT := gzipWriter.Close()
	if errT != nil {
		return err
	}

	return buf.Bytes()
}

var Compress = TKX.Compress

func (pA *TK) Uncompress(dataA interface{}, argsA ...interface{}) interface{} {
	var inputBufT []byte

	switch nv := dataA.(type) {
	case []byte:
		inputBufT = nv
	// case string:
	// 	inputBufT = []byte(inputBufT)
	default:
		return fmt.Errorf("failed to uncompress, unsupported type: %T(%v)", dataA, dataA)
	}

	bytesReader := bytes.NewReader(inputBufT)

	gzipReader, err := gzip.NewReader(bytesReader)
	if err != nil {
		return err
	}

	defer gzipReader.Close()

	buf := new(bytes.Buffer)

	_, err = buf.ReadFrom(gzipReader)

	if err != nil {
		return err
	}

	return buf.Bytes()
}

var Uncompress = TKX.Uncompress

func (pA *TK) CompressText(textA string) string {
	dataT := Compress([]byte(textA))

	_, ok := dataT.(error)
	if ok {
		return textA
	}

	return EncodeToBase64(dataT.([]byte))
}

var CompressText = TKX.CompressText

func (pA *TK) UncompressText(textA string) string {
	dataT, errT := DecodeFromBase64(textA)
	if errT != nil {
		return textA
	}

	dataOutT := Uncompress(dataT)

	_, ok := dataOutT.(error)
	if ok {
		return textA
	}

	return string(dataOutT.([]byte))
}

var UncompressText = TKX.UncompressText

func (pA *TK) GetZipArchiver(argsA ...string) *archiver.Zip {
	z := &archiver.Zip{
		// CompressionLevel:       3,
		OverwriteExisting: IfSwitchExistsWhole(argsA, "-overwrite"),
		// MkdirAll:               true,
		// SelectiveCompression:   true,
		ImplicitTopLevelFolder: IfSwitchExistsWhole(argsA, "-implicitTop"),
		ContinueOnError:        IfSwitchExistsWhole(argsA, "-continueOnError"),
	}

	return z
}

var GetZipArchiver = TKX.GetZipArchiver
var GetZip = TKX.GetZipArchiver

func (pA *TK) GetEncZipReader(filePathA string) (*zipx.ReadCloser, error) {
	return zipx.OpenReader(filePathA)
}

var GetEncZipReader = TKX.GetEncZipReader

func (pA *TK) GetEncZipWriter(wA io.Writer) *zipx.Writer {
	return zipx.NewWriter(wA)
}

var GetEncZipWriter = TKX.GetEncZipWriter

func (pA *TK) TableToMSSArray(tableA [][]string) []map[string]string {
	if tableA == nil {
		return []map[string]string{}
	}

	lenT := len(tableA)

	if lenT < 1 {
		return []map[string]string{}
	}

	inLenT := len(tableA[0])

	bufT := make([]map[string]string, 0, lenT)

	for i, v := range tableA {
		if i == 0 {
			continue
		}

		inBufT := make(map[string]string, inLenT)

		for j, jv := range v {
			inBufT[tableA[0][j]] = jv
		}

		bufT = append(bufT, inBufT)
	}

	return bufT

}

var TableToMSSArray = TKX.TableToMSSArray

func (pA *TK) TableToOrderedMapArray(tableA [][]string) []*OrderedMap {
	if tableA == nil {
		return []*OrderedMap{}
	}

	lenT := len(tableA)

	if lenT < 1 {
		return []*OrderedMap{}
	}

	// inLenT := len(tableA[0])

	bufT := make([]*OrderedMap, 0, lenT)

	for i, v := range tableA {
		if i == 0 {
			continue
		}

		inBufT := NewOrderedMap()

		for j, jv := range v {
			inBufT.Set(tableA[0][j], jv)
			// Pl("%v -> %v -> %v -> %#v", tableA[0][j], jv, inBufT, inBufT)
		}

		bufT = append(bufT, inBufT)
	}

	return bufT

}

var TableToOrderedMapArray = TKX.TableToOrderedMapArray

func (pA *TK) TableToMSSMap(tableA [][]string, keyA string) map[string]map[string]string {
	if tableA == nil {
		return map[string]map[string]string{}
	}

	lenT := len(tableA)

	if lenT < 1 {
		return map[string]map[string]string{}
	}

	inLenT := len(tableA[0])

	bufT := make(map[string]map[string]string, 0)

	for i, v := range tableA {
		if i == 0 {
			continue
		}

		inBufT := make(map[string]string, inLenT)

		for j, jv := range v {
			inBufT[tableA[0][j]] = jv
		}

		bufT[inBufT[keyA]] = inBufT
	}

	return bufT

}

var TableToMSSMap = TKX.TableToMSSMap

func (pA *TK) TableToMSSOrderedMap(tableA [][]string, keyA string) map[string]*OrderedMap {
	if tableA == nil {
		return map[string]*OrderedMap{}
	}

	lenT := len(tableA)

	if lenT < 1 {
		return map[string]*OrderedMap{}
	}

	// inLenT := len(tableA[0])

	bufT := make(map[string]*OrderedMap, 0)

	for i, v := range tableA {
		if i == 0 {
			continue
		}

		inBufT := NewOrderedMap()

		for j, jv := range v {
			inBufT.Set(tableA[0][j], jv)
		}

		bufT[inBufT.GetString(keyA)] = inBufT
	}

	return bufT

}

var TableToMSSOrderedMap = TKX.TableToMSSOrderedMap

func (pA *TK) TableToMSSMapArray(tableA [][]string, keyA string) map[string][]map[string]string {
	if tableA == nil {
		return map[string][]map[string]string{}
	}

	lenT := len(tableA)

	if lenT < 1 {
		return map[string][]map[string]string{}
	}

	inLenT := len(tableA[0])

	bufT := make(map[string][]map[string]string, 0)

	for i, v := range tableA {
		if i == 0 {
			continue
		}

		inBufT := make(map[string]string, inLenT)

		for j, jv := range v {
			inBufT[tableA[0][j]] = jv
		}

		_, ok := bufT[inBufT[keyA]]
		if !ok {
			bufT[inBufT[keyA]] = make([]map[string]string, 0, lenT)
		}

		bufT[inBufT[keyA]] = append(bufT[inBufT[keyA]], inBufT)
	}

	return bufT

}

var TableToMSSMapArray = TKX.TableToMSSMapArray

func (pA *TK) TableToMSSJSON(tableA [][]string) string {
	if tableA == nil {
		return "[]"
	}

	lenT := len(tableA)

	if lenT < 1 {
		return "[]"
	}

	inLenT := len(tableA[0])

	bufT := make([]map[string]string, 0, lenT)

	for i, v := range tableA {
		if i == 0 {
			continue
		}

		inBufT := make(map[string]string, inLenT)

		for j, jv := range v {
			inBufT[tableA[0][j]] = jv
		}

		bufT = append(bufT, inBufT)
	}

	return ToJSONX(bufT, "-default=", "-sort")

}

var TableToMSSJSON = TKX.TableToMSSJSON

func (pA *TK) GetUUID1() string {
	uuidT, errT := uuid.NewV1()
	if errT != nil {
		return GenerateErrorStringF("failed to generate UUID: %v", errT)
	}

	return uuidT.String()
}

var GetUUID1 = TKX.GetUUID1

func (pA *TK) GetUUID4() string {
	u1 := uuid.Must(uuid.NewV4())
	return u1.String()
}

var GetUUID4 = TKX.GetUUID4

func (pA *TK) GetUUID() string {
	uuidT, errT := uuid.NewV1()
	if errT != nil {
		return GenerateErrorStringF("failed to generate UUID: %v", errT)
	}

	return uuidT.String()
}

var GetUUID = TKX.GetUUID

func (pA *TK) IsFloat64NearlyEqual(a, b float64) bool {

	if math.Abs(a-b) < 0.000001 {
		return true
	}

	return false
}

var IsFloat64NearlyEqual = TKX.IsFloat64NearlyEqual

// SetValue set a value to a pointer
func (pA *TK) SetValue(p interface{}, v interface{}) error {
	// tk.Pl("%#v", reflect.TypeOf(p).Kind())
	// p = v

	srcRef := reflect.ValueOf(v)
	vp := reflect.ValueOf(p)
	// Pl("srcRef: %v, vp: %v, Elem: %v", srcRef, vp, vp.Elem())
	vp.Elem().Set(srcRef)

	return nil
}

var SetValue = TKX.SetValue

// GetValue get a value from a pointer
func (pA *TK) GetValue(p interface{}) interface{} {
	vp := reflect.Indirect(reflect.ValueOf(p))
	return vp.Interface()
}

var GetValue = TKX.GetValue

func (pA *TK) GetPointer(p interface{}) interface{} {
	vp := reflect.Indirect(reflect.ValueOf(p))
	return vp.Pointer()
}

var GetPointer = TKX.GetPointer

func (pA *TK) GetRef(p interface{}) interface{} {
	// Pl("def2: %T %v", p, p)

	switch nnv := p.(type) {
	case *bool:
		return nnv
	case *byte:
		return nnv
	case *rune:
		return nnv
	case *int:
		return nnv
	case *int64:
		return nnv
	case *float32:
		return nnv
	case *float64:
		return nnv
	case *string:
		return nnv
	case *[]byte:
		return nnv
	case *[]rune:
		return nnv
	case *[]int:
		return nnv
	case *[]int64:
		return nnv
	case *[]float32:
		return nnv
	case *[]float64:
		return nnv
	case *[]string:
		return nnv
	case *map[string]string:
		return nnv
	case *map[string]interface{}:
		return nnv
	case *bytes.Buffer:
		return nnv
	case *strings.Builder:
		return nnv
	case *interface{}:
		// Pl("here3")
		v2 := *nnv

		switch nv := v2.(type) {
		case bool:
			pr := (*bool)(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(nnv)) + 8))))
			return pr
		case byte:
			pr := (*byte)(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(nnv)) + 8))))
			return pr
		case rune:
			pr := (*rune)(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(nnv)) + 8))))
			return pr
		case int:
			pr := (*int)(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(nnv)) + 8))))
			return pr
		case int64:
			pr := (*int64)(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(nnv)) + 8))))
			return pr
		case float32:
			pr := (*float32)(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(nnv)) + 8))))
			return pr
		case float64:
			pr := (*float64)(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(nnv)) + 8))))
			// fmt.Printf("nv(%T)=%v, &nv=%v\n", nv, nv, &nv)
			return pr
		case string:
			pr := (*string)(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(nnv)) + 8))))
			return pr
		case []byte:
			pr := (*[]byte)(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(nnv)) + 8))))
			return pr
		case []rune:
			pr := (*[]rune)(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(nnv)) + 8))))
			return pr
		case []int:
			pr := (*[]int)(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(nnv)) + 8))))
			return pr
		case []int64:
			pr := (*[]int64)(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(nnv)) + 8))))
			return pr
		case []float32:
			pr := (*[]float32)(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(nnv)) + 8))))
			return pr
		case []float64:
			pr := (*[]float64)(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(nnv)) + 8))))
			return pr
		case []string:
			pr := (*[]string)(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(nnv)) + 8))))
			return pr
		case map[string]string:
			pr := (*map[string]string)(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(nnv)) + 8))))
			return pr
		case map[string]interface{}:
			pr := (*map[string]interface{})(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(nnv)) + 8))))
			return pr
		case bytes.Buffer:
			pr := (*bytes.Buffer)(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(nnv)) + 8))))
			return pr
		case strings.Builder:
			pr := (*strings.Builder)(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(nnv)) + 8))))
			return pr

		default:
			fmt.Printf("nv(%T)=%v, &nv=%v\n", nv, nv, &nv)
		}

		return nnv
	default:
		return nnv

	}
}

var GetRef = TKX.GetRef

func (pA *TK) GetAddr(p interface{}) interface{} {
	vp := reflect.Indirect(reflect.ValueOf(p))
	return vp.Addr()
}

var GetAddr = TKX.GetAddr

func (pA *TK) SetByRef(ppA interface{}, vA interface{}) (result error) {
	defer func() {
		r := recover()

		if r != nil {
			result = fmt.Errorf("faile to set by ref: %v(%v)", r, string(debug.Stack()))

			return
		}
	}()

	// Pln("---")
	// Plv(ppA)
	// Plv(vA)

	switch nv := ppA.(type) {
	case *bool:
		*nv = ToBool(vA)
		return nil

	case *byte:
		*nv = ToByte(vA, 0)
		return nil

	case *int:
		*nv = ToInt(vA, 0)
		return nil

	case *int64:
		*nv = int64(ToInt(vA, 0))
		return nil

	case *rune:
		*nv = ToRune(vA, 0)
		return nil

	case *uint:
		*nv = uint(ToInt(vA, 0))
		return nil

	case *uint64:
		*nv = uint64(ToInt(vA, 0))
		return nil

	case *float32:
		*nv = float32(ToFloat(vA, 0.0))
		return nil

	case *float64:
		*nv = ToFloat(vA, 0.0)
		return nil

	case *string:
		*nv = ToStr(vA)
		return nil

	case *error:
		*nv = fmt.Errorf("%v", ToStr(vA))
		return nil

	case *interface{}:
		*nv = vA
		return nil

	}

	valueT := reflect.ValueOf(ppA)

	kindT := valueT.Kind()

	if kindT != reflect.Pointer {
		return fmt.Errorf("not pointer type")
	}

	elemT := valueT.Elem()

	if !elemT.CanSet() {
		return fmt.Errorf("value cannot be set")
	}

	elemT.Set(reflect.ValueOf(vA))

	return nil
}

var SetByRef = TKX.SetByRef

func (pA *TK) GetRefValue(ppA interface{}) (result interface{}, err error) {
	defer func() {
		r := recover()

		if r != nil {
			result = nil
			err = fmt.Errorf("faile to get ref value: %v(%v)", r, string(debug.Stack()))
			return
		}
	}()

	switch nv := ppA.(type) {
	case *FlexRef:
		return nv.GetValue(), nil
	case *MapRef:
		return nv.GetValue(), nil
	case *bool:
		return *nv, nil
	case *byte:
		return *nv, nil
	case *int:
		return *nv, nil
	case *int64:
		return *nv, nil
	case *rune:
		return *nv, nil
	case *uint:
		return *nv, nil
	case *uint64:
		return *nv, nil
	case *float32:
		return *nv, nil
	case *float64:
		return *nv, nil
	case *string:
		return *nv, nil
	case *error:
		return *nv, nil
	case *[]string:
		return *nv, nil
	case *[]interface{}:
		return *nv, nil
	case *map[string]string:
		return *nv, nil
	case *map[string]interface{}:
		return *nv, nil
	case *interface{}:
		return *nv, nil
	}

	valueT := reflect.ValueOf(ppA)

	kindT := valueT.Kind()

	if kindT != reflect.Pointer {
		return nil, fmt.Errorf("not pointer type")
	}

	elemT := valueT.Elem()

	if !elemT.CanInterface() {
		return nil, fmt.Errorf("value cannot convert to interface")
	}

	return elemT.Interface(), nil
}

var GetRefValue = TKX.GetRefValue

func (pA *TK) GetRefValueQuick(ppA interface{}) (result interface{}) {
	defer func() {
		r := recover()

		if r != nil {
			result = fmt.Errorf("faile to get ref value: %v(%v)", r, string(debug.Stack()))
			return
		}
	}()

	// Plvx(ppA)

	switch nv := ppA.(type) {
	case *FlexRef:
		return nv.GetValue()
	case *MapRef:
		return nv.GetValue()
	case *bool:
		return *nv
	case *byte:
		// Pl("byte: %v", *nv)
		return *nv
	case *int:
		return *nv
	case *rune:
		return *nv
	case *int64:
		return *nv
	case *uint:
		return *nv
	case *uint64:
		return *nv
	case *float32:
		return *nv
	case *float64:
		return *nv
	case *string:
		return *nv
	case *error:
		return *nv
	case *[]string:
		return *nv
	case *[]interface{}:
		return *nv
	case *map[string]string:
		return *nv
	case *map[string]interface{}:
		return *nv
	case *interface{}:
		return *nv
	}

	valueT := reflect.ValueOf(ppA)

	kindT := valueT.Kind()

	if kindT != reflect.Pointer {
		return fmt.Errorf("not pointer type")
	}

	elemT := valueT.Elem()

	if !elemT.CanInterface() {
		return fmt.Errorf("value cannot convert to interface")
	}

	return elemT.Interface()
}

var GetRefValueQuick = TKX.GetRefValueQuick

func (pA *TK) CheckErrorFunc(errA error, funcA func()) {
	if errA != nil {
		PlErr(errA)

		if funcA != nil {
			funcA()
		}

		os.Exit(1)
	}

}

var CheckErrorFunc = TKX.CheckErrorFunc

func (pA *TK) CheckError(errA interface{}, funcsA ...(func())) {

	errT, ok := errA.(error)
	if !ok {
		return
	}

	if errT != nil {
		PlErr(errT)

		if funcsA != nil {
			for _, v := range funcsA {
				v()
			}
		}

		os.Exit(1)
	}

}

var CheckError = TKX.CheckError

func (pA *TK) CheckErrX(errA interface{}, funcsA ...(func())) {

	if IsErrX(errA) {
		Pl("Error: %v", GetErrStrX(errA))

		for _, v := range funcsA {
			v()
		}

		os.Exit(1)
	}

}

var CheckErrX = TKX.CheckErrX

func (pA *TK) CheckErrorString(strA string, funcsA ...(func())) {
	if IsErrorString(strA) {
		PlErrString(strA)

		if funcsA != nil {
			for _, v := range funcsA {
				v()
			}
		}

		os.Exit(1)
	}

}

var CheckErrorString = TKX.CheckErrorString

func (pA *TK) TypeOfValue(vA interface{}) string {
	return fmt.Sprintf("%T", vA)
}

var TypeOfValue = TKX.TypeOfValue

func (pA *TK) TypeOfValueReflect(vA interface{}) string {
	rs := reflect.TypeOf(vA)
	return rs.String()
}

var TypeOfValueReflect = TKX.TypeOfValueReflect

func (pA *TK) KindOfValueReflect(vA interface{}) string {
	rs := reflect.TypeOf(vA)
	return rs.Kind().String()
}

var KindOfValueReflect = TKX.KindOfValueReflect

func (pA *TK) GetClipText() string {
	textT, errT := clipboard.ReadAll()
	if errT != nil {
		return GenerateErrorStringF("could not get text from clipboard: %v", errT.Error())
	}

	return textT
}

var GetClipText = TKX.GetClipText

func (pA *TK) GetClipboardTextWithDefault(defaultA string) string {
	textT, errT := clipboard.ReadAll()
	if errT != nil {
		return defaultA
	} else {
		return textT
	}

}

var GetClipboardTextWithDefault = TKX.GetClipboardTextWithDefault

func (pA *TK) GetClipboardTextDefaultEmpty() string {
	textT, errT := clipboard.ReadAll()
	if errT != nil {
		return ""
	} else {
		return textT
	}

}

var GetClipboardTextDefaultEmpty = TKX.GetClipboardTextDefaultEmpty

func (pA *TK) SetClipText(textA string) error {
	return clipboard.WriteAll(textA)
}

var SetClipText = TKX.SetClipText

func (pA *TK) GetTextFromFileOrClipboard(fileT string, defaultA string) string {
	if IsEmptyTrim(fileT) {
		return GetClipboardTextWithDefault(defaultA)
	}

	if !IfFileExists(fileT) {
		return GetClipboardTextWithDefault(defaultA)
	}

	return LoadStringFromFileWithDefault(fileT, defaultA)
}

var GetTextFromFileOrClipboard = TKX.GetTextFromFileOrClipboard

// RemoveItemsInArray
func (pA *TK) RemoveItemsInArray(aryA interface{}, startA int, endA int) interface{} {

	aryT, ok := aryA.([]interface{})

	if ok {
		if startA < 0 || startA >= len(aryT) {
			// Pl("Runtime error: %v", "index out of range")
			return nil
		}

		if endA < 0 || endA >= len(aryT) {
			// Pl("Runtime error: %v", "index out of range")
			return nil
		}

		rs := make([]interface{}, 0, len(aryT)-(endA+1-startA))

		rs = append(rs, aryT[:startA]...)
		rs = append(rs, aryT[endA+1:]...)

		return rs
	}

	aryST, ok := aryA.([]string)

	if ok {
		if startA < 0 || startA >= len(aryST) {
			return nil
		}

		if endA < 0 || endA >= len(aryST) {
			return nil
		}

		rs := make([]string, 0, len(aryST)-(endA+1-startA))

		rs = append(rs, aryST[:startA]...)
		rs = append(rs, aryST[endA+1:]...)

		return rs
	}

	arySAT, ok := aryA.([][]string)

	if ok {
		if startA < 0 || startA >= len(arySAT) {
			return nil
		}

		if endA < 0 || endA >= len(arySAT) {
			return nil
		}

		rs := make([][]string, 0, len(arySAT)-(endA+1-startA))

		rs = append(rs, arySAT[:startA]...)
		rs = append(rs, arySAT[endA+1:]...)

		return rs
	}

	arySIT, ok := aryA.([][]interface{})

	if ok {
		if startA < 0 || startA >= len(arySIT) {
			return nil
		}

		if endA < 0 || endA >= len(arySIT) {
			return nil
		}

		rs := make([][]interface{}, 0, len(arySIT)-(endA+1-startA))

		rs = append(rs, arySIT[:startA]...)
		rs = append(rs, arySIT[endA+1:]...)

		return rs
	}

	aryIT, ok := aryA.([]int)

	if ok {
		if startA < 0 || startA >= len(aryIT) {
			return nil
		}

		if endA < 0 || endA >= len(aryIT) {
			return nil
		}

		rs := make([]int, 0, len(aryIT)-(endA+1-startA))

		rs = append(rs, aryIT[:startA]...)
		rs = append(rs, aryIT[endA+1:]...)

		return rs
	}

	aryI64T, ok := aryA.([]int64)

	if ok {
		if startA < 0 || startA >= len(aryI64T) {
			return nil
		}

		if endA < 0 || endA >= len(aryI64T) {
			return nil
		}

		rs := make([]int64, 0, len(aryI64T)-(endA+1-startA))

		rs = append(rs, aryI64T[:startA]...)
		rs = append(rs, aryI64T[endA+1:]...)

		return rs
	}

	aryFT, ok := aryA.([]float64)

	if ok {
		if startA < 0 || startA >= len(aryFT) {
			return nil
		}

		if endA < 0 || endA >= len(aryFT) {
			return nil
		}

		rs := make([]float64, 0, len(aryFT)-(endA+1-startA))

		rs = append(rs, aryFT[:startA]...)
		rs = append(rs, aryFT[endA+1:]...)

		return rs
	}

	aryBT, ok := aryA.([]bool)

	if ok {
		if startA < 0 || startA >= len(aryBT) {
			return nil
		}

		if endA < 0 || endA >= len(aryBT) {
			return nil
		}

		rs := make([]bool, 0, len(aryBT)-(endA+1-startA))

		rs = append(rs, aryBT[:startA]...)
		rs = append(rs, aryBT[endA+1:]...)

		return rs
	}

	aryByT, ok := aryA.([]byte)

	if ok {
		if startA < 0 || startA >= len(aryByT) {
			return nil
		}

		if endA < 0 || endA >= len(aryByT) {
			return nil
		}

		rs := make([]byte, 0, len(aryByT)-(endA+1-startA))

		rs = append(rs, aryByT[:startA]...)
		rs = append(rs, aryByT[endA+1:]...)

		return rs
	}

	return nil
	// if idxT == 0 {
	// 	return ayrA[idxT + 1:]
	// }

	// if idxT == len(aryA) - 1 {
	// 	return ayrA[0:len(aryA) - 1]
	// }

	// return append(aryA[:idxA], aryA[idxA+1:]...)

}

var RemoveItemsInArray = TKX.RemoveItemsInArray

func (pA *TK) RemoveStringInArray(sliceA []string, idxA int) []string {
	if idxA >= 0 && idxA < len(sliceA) {
		copy(sliceA[idxA:], sliceA[idxA+1:])
		sliceA = sliceA[0 : len(sliceA)-1]
	}

	return sliceA
}

var RemoveStringInArray = TKX.RemoveStringInArray

func (pA *TK) RemoveIntInArray(sliceA []int, idxA int) []int {
	if idxA >= 0 && idxA < len(sliceA) {
		copy(sliceA[idxA:], sliceA[idxA+1:])
		sliceA = sliceA[0 : len(sliceA)-1]
	}

	return sliceA
}

var RemoveIntInArray = TKX.RemoveIntInArray

func (pA *TK) BitXor(p interface{}, v interface{}) interface{} {
	switch nv := p.(type) {
	case bool:
		if nv {
			return !(v.(bool))
		}

		return v.(bool)
	case int:
		return p.(int) ^ v.(int)
	case int64:
		return p.(int64) ^ v.(int64)
	case int32:
		return p.(int32) ^ v.(int32)
	case int16:
		return p.(int16) ^ v.(int16)
	case int8:
		return p.(int8) ^ v.(int8)
	case uint64:
		return p.(uint64) ^ v.(uint64)
	case uint32:
		return p.(uint32) ^ v.(uint32)
	case uint16:
		return p.(uint16) ^ v.(uint16)
	case uint8:
		return p.(uint8) ^ v.(uint8)
	case uint:
		return p.(uint) ^ v.(uint)
	}

	return 0
}

var BitXor = TKX.BitXor

func (pA *TK) ToPointerStringArray(aryA []string) *[]string {
	return &aryA
}

var ToPointerStringArray = TKX.ToPointerStringArray

func (pA *TK) ToPointerFloat64Array(aryA []float64) *[]float64 {
	return &aryA
}

var ToPointerFloat64Array = TKX.ToPointerFloat64Array

// ParseHexColor inspired by gg
func (pA *TK) ParseHexColor(x string) (r, g, b, a int) {
	x = strings.TrimPrefix(x, "#")
	a = 255
	if len(x) == 3 {
		format := "%1x%1x%1x"
		fmt.Sscanf(x, format, &r, &g, &b)
		r |= r << 4
		g |= g << 4
		b |= b << 4
	}
	if len(x) == 6 {
		format := "%02x%02x%02x"
		fmt.Sscanf(x, format, &r, &g, &b)
	}
	if len(x) == 8 {
		format := "%02x%02x%02x%02x"
		fmt.Sscanf(x, format, &r, &g, &b, &a)
	}
	return
}

var ParseHexColor = TKX.ParseHexColor

func (pA *TK) HexToColor(strA string, optsA ...string) color.Color {
	r, g, b, a := ParseHexColor(strA)

	formatT := GetSwitch(optsA, "-format=", "")

	if formatT == "NRGBA" {
		return color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	}

	return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
}

var HexToColor = TKX.HexToColor

func (pA *TK) ColorToHex(colorA color.Color, optsA ...string) string {

	r, g, b, a := colorA.RGBA()

	if IfSwitchExists(optsA, "-withAlpha") {

		return fmt.Sprintf("%02X%02X%02X%02X", uint8(r), uint8(g), uint8(b), uint8(a))
	}

	return fmt.Sprintf("%02X%02X%02X", uint8(r), uint8(g), uint8(b))
}

var ColorToHex = TKX.ColorToHex

// DeepClone deep copies original and returns the copy as an interface.
func (pA *TK) DeepClone(original interface{}) (copy interface{}) {
	if original == nil {
		return nil
	}
	value := reflect.ValueOf(original)
	return pA.deepCopy(value).Interface()
}

var DeepClone = TKX.DeepClone

// DeepCopyFromTo deep copies original and assigns the copy to the copy argument (pointer).
func (pA *TK) DeepCopyFromTo(original, copy interface{}) error {
	if original == nil {
		copy = nil
		return nil
	} else if copy == nil { // TODO try to initialize it here
		return fmt.Errorf("FromTo: copy target is nil, it should be a valid pointer")
		// copyValue := reflect.New(value.Type().Elem()).Elem()
		// copy = copyValue.Interface()
	}
	copyValue := reflect.ValueOf(copy)
	if copyValue.Kind() != reflect.Ptr {
		return fmt.Errorf("FromTo: copy target type %T and not a pointer", copy)
	}
	value := reflect.ValueOf(original)
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			copy = nil // TODO return typed nil
			return nil
		}
		value = value.Elem()
	}
	copyValue.Elem().Set(pA.deepCopy(value))
	return nil
}

var DeepCopyFromTo = TKX.DeepCopyFromTo

func (pA *TK) deepCopy(original reflect.Value) reflect.Value {
	switch original.Kind() {
	case reflect.Slice:
		return pA.deepCopySlice(original)
	case reflect.Map:
		return pA.deepCopyMap(original)
	case reflect.Ptr:
		return pA.deepCopyPointer(original)
	case reflect.Struct:
		return pA.deepCopyStruct(original)
	case reflect.Chan:
		return pA.deepCopyChan(original)
	case reflect.Array:
		return pA.deepCopyArray(original)
	default:
		return forceCopyValue(original)
	}
}

// forceCopyValue simply creates a new pointer and sets its value to the original.
func (pA *TK) forceCopyValue(original reflect.Value) reflect.Value {
	originalType := original.Type()
	newPointer := reflect.New(originalType)
	newPointer.Elem().Set(original)
	return newPointer.Elem()
}

var forceCopyValue = TKX.forceCopyValue

func (pA *TK) deepCopySlice(original reflect.Value) reflect.Value {
	if original.IsNil() {
		return original
	}
	copy := reflect.MakeSlice(original.Type(), 0, 0)
	for i := 0; i < original.Len(); i++ {
		elementCopy := pA.deepCopy(original.Index(i))
		copy = reflect.Append(copy, elementCopy)
	}
	return copy
}

func (pA *TK) deepCopyArray(original reflect.Value) reflect.Value {
	if original.Len() == 0 {
		// it cannot be changed anyway, so we can return the original
		return original
	}
	elementType := original.Index(0).Type()
	arrayType := reflect.ArrayOf(original.Len(), elementType)
	newPointer := reflect.New(arrayType)
	copy := newPointer.Elem()
	for i := 0; i < original.Len(); i++ {
		subCopy := pA.deepCopy(original.Index(i))
		copy.Index(i).Set(subCopy)
	}
	return copy
}

func (pA *TK) deepCopyMap(original reflect.Value) reflect.Value {
	if original.IsNil() {
		return original
	}
	keyType := original.Type().Key()
	valueType := original.Type().Elem()
	mapType := reflect.MapOf(keyType, valueType)
	copy := reflect.MakeMap(mapType)
	for _, key := range original.MapKeys() {
		value := pA.deepCopy(original.MapIndex(key))
		copy.SetMapIndex(key, value)
	}
	return copy
}

func (pA *TK) deepCopyPointer(original reflect.Value) reflect.Value {
	if original.IsNil() {
		return original
	}
	element := original.Elem()
	copy := reflect.New(element.Type())
	copyElement := pA.deepCopy(element)
	copy.Elem().Set(copyElement)
	return copy
}

var deepCopyPointer = TKX.deepCopyPointer

func (pA *TK) deepCopyStruct(original reflect.Value) reflect.Value {
	copy := reflect.New(original.Type()).Elem()
	copy.Set(original)
	for i := 0; i < original.NumField(); i++ {
		fieldValue := copy.Field(i)
		fieldValue = reflect.NewAt(fieldValue.Type(), unsafe.Pointer(fieldValue.UnsafeAddr())).Elem()
		fieldValue.Set(pA.deepCopy(fieldValue))
	}
	return copy
}

func (pA *TK) deepCopyChan(original reflect.Value) reflect.Value {
	return reflect.MakeChan(original.Type(), original.Cap())
}

type PlainAuth struct {
	identity, username, password string
	host                         string
}

type ServerInfo struct {
	Name string   // SMTP server name
	TLS  bool     // using TLS, with valid certificate for Name
	Auth []string // advertised authentication mechanisms
}

// PlainAuth get plain auth
func (pA *TK) GetPlainAuth(identity, username, password, host string) smtp.Auth {
	return &PlainAuth{identity, username, password, host}
}

var GetPlainAuth = TKX.GetPlainAuth

func (pA *TK) isLocalhost(name string) bool {
	return name == "localhost" || name == "127.0.0.1" || name == "::1"
}

var isLocalhost = TKX.isLocalhost

func (a *PlainAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	// Must have TLS, or else localhost server.
	// Note: If TLS is not true, then we can't trust ANYTHING in ServerInfo.
	// In particular, it doesn't matter if the server advertises PLAIN auth.
	// That might just be the attacker saying
	// "it's ok, you can trust me with your password."
	// if !server.TLS && !isLocalhost(server.Name) {
	// 	return "", nil, errors.New("unencrypted connection")
	// }
	if server.Name != a.host {
		return "", nil, errors.New("wrong host name")
	}
	resp := []byte(a.identity + "\x00" + a.username + "\x00" + a.password)
	return "PLAIN", resp, nil
}

func (a *PlainAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	command := string(fromServer)
	command = strings.TrimSpace(command)
	command = strings.TrimSuffix(command, ":")
	command = strings.ToLower(command)

	if more {
		if command == "username" {
			return []byte(fmt.Sprintf("%s", a.username)), nil
		} else if command == "password" {
			return []byte(fmt.Sprintf("%s", a.password)), nil
		} else {
			// We've already sent everything.
			return nil, fmt.Errorf("unexpected server challenge: %s", command)
		}
	}
	return nil, nil
}

type LoginAuth struct {
	username, password string
}

func (pA *TK) GetLoginAuth(username, password string) smtp.Auth {
	return &LoginAuth{username, password}
}

var GetLoginAuth = TKX.GetLoginAuth

func (a *LoginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", nil, nil
}

func (a *LoginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	command := string(fromServer)
	command = strings.TrimSpace(command)
	command = strings.TrimSuffix(command, ":")
	command = strings.ToLower(command)

	if more {
		if command == "username" {
			return []byte(fmt.Sprintf("%s", a.username)), nil
		} else if command == "password" {
			return []byte(fmt.Sprintf("%s", a.password)), nil
		} else {
			// We've already sent everything.
			return nil, fmt.Errorf("unexpected server challenge: %s", command)
		}
	}
	return nil, nil
}

// network related

func (pA *TK) IsHttps(reqA *http.Request) bool {
	if reqA.TLS != nil || reqA.Header.Get("X-Forwarded-Proto") == "https" {
		return true
	}

	return false
}

var IsHttps = TKX.IsHttps

func (pA *TK) Listen(netA string, addrA string) interface{} {
	rsT, errT := net.Listen(netA, addrA)
	if errT != nil {
		return fmt.Errorf("failed to listen: %v", errT)
	}

	return rsT
}

var Listen = TKX.Listen

// sock5 related

func (pA *TK) StartSocksServer(optionsA ...string) error {
	ipT := GetSwitchWithDefaultValue(optionsA, "-ip=", "0.0.0.0")
	portT := GetSwitchWithDefaultValue(optionsA, "-port=", "7480")
	passwordT := GetSwitchWithDefaultValue(optionsA, "-password=", "acb123!@#")
	verboseT := IfSwitchExistsWhole(optionsA, "-verbose")

	lenT := len(passwordT)

	if lenT < 16 {
		passwordT = passwordT + strings.Repeat("z", 16-lenT)
	} else if lenT > 16 {
		passwordT = passwordT[0:16]
	}

	remote, err := net.Listen("tcp", fmt.Sprintf("%s:%s", ipT, portT))

	if err != nil {
		return err
	}

	if verboseT {
		Pl("Start server on %v:%v", ipT, portT)
	}

	for {
		conn, err := remote.Accept()
		if err != nil {
			return Errf("accept err: %v", err)
		} else {
			if verboseT {
				Pl("Client connected: %v", conn.RemoteAddr())
			}
			socks.OpenRemoteTunnel(conn, passwordT)
		}
	}

}

var StartSocksServer = TKX.StartSocksServer

func (pA *TK) StartSocksClient(optionsA ...string) error {
	remoteIpT := GetSwitchWithDefaultValue(optionsA, "-remoteIp=", "0.0.0.0")
	remotePortT := GetSwitchWithDefaultValue(optionsA, "-remotePort=", "7480")
	localIpT := GetSwitchWithDefaultValue(optionsA, "-localIp=", "0.0.0.0")
	localPortT := GetSwitchWithDefaultValue(optionsA, "-localPort=", "7481")
	passwordT := GetSwitchWithDefaultValue(optionsA, "-password=", "acb123!@#")
	verboseT := IfSwitchExistsWhole(optionsA, "-verbose")

	lenT := len(passwordT)

	if lenT < 16 {
		passwordT = passwordT + strings.Repeat("z", 16-lenT)
	} else if lenT > 16 {
		passwordT = passwordT[0:16]
	}

	local, err := net.Listen("tcp", fmt.Sprintf("%s:%s", localIpT, localPortT))
	if err != nil {
		return err
	}

	if verboseT {
		Pl("Start socks proxy on %v:%v, remote server: %v: %v", localIpT, localPortT, remoteIpT, remotePortT)
	}

	for {
		conn, err := local.Accept()
		if err != nil {
			return Errf("accept err: %v", err)
		} else {
			if verboseT {
				Pl("Client connected: %v", conn.RemoteAddr())
			}
			socks.OpenLocalTunnel(conn, fmt.Sprintf("%s:%s", remoteIpT, remotePortT), passwordT)
		}
	}

}

var StartSocksClient = TKX.StartSocksClient

// Transparent proxy related

func (pA *TK) copyWR(w io.WriteCloser, r io.Reader) {
	_, err := io.Copy(w, r)
	if err != nil {
		Pl("failed to copy: %v", err)
	}
	w.Close()
}

var copyWR = TKX.copyWR

func (pA *TK) forwardConn(lc net.Conn, server string, verboseA bool) {
	defer lc.Close()

	rc, err := net.Dial("tcp", server)
	if err != nil {
		if verboseA {
			Pl("failed to dial: %v", err)
		}
		return
	}

	if verboseA {
		Pl("Forwarding connection to %v", server)
	}
	go copyWR(rc, lc)

	copyWR(lc, rc)

	if verboseA {
		Pl("Terminated:  %s -> %s ", lc.RemoteAddr(), server)
	}
}

var forwardConn = TKX.forwardConn

func (pA *TK) StartTransparentProxy(local, server string, optionsA ...string) error {
	verboseT := IfSwitchExistsWhole(optionsA, "-verbose")

	l, err := net.Listen("tcp", local)
	if err != nil {
		if verboseT {
			Pl("failed to listen: %v", err)
		}

		return Errf("failed to listen: %v", err)
	}

	defer l.Close()

	if verboseT {
		Pl("Listening on %v", l.Addr())
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			if verboseT {
				Pl("failed to accept: %v", err)
			}
		}

		if verboseT {
			Pl("New connection from %v", conn.RemoteAddr())
		}

		go forwardConn(conn, server, verboseT)
	}
}

var StartTransparentProxy = TKX.StartTransparentProxy

func (pA *TK) StartTransparentProxy2(localA, remoteA string, optionsA ...string) error {
	verboseT := IfSwitchExistsWhole(optionsA, "-verbose")

	listener, err := net.Listen("tcp", localA)
	if err != nil {
		if verboseT {
			Pl("Failed to setup listener: %v", err)
		}
		return Errf("Failed to setup listener: %v", err)
	}

	defer listener.Close()

	if verboseT {
		Pl("Listen on: %v:%v", listener.Addr(), "")
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			if verboseT {
				Pl("Failed to accept connection: %v", err)
			}
			return Errf("ERROR: failed to accept listener: %v", err)
		}

		// if verboseT {
		// 	Pl("Accepted connection on %v from %v", conn.LocalAddr(), conn.RemoteAddr())
		// }

		go func(connA net.Conn, remote1A string) {
			client, err := net.Dial("tcp", remote1A)
			if err != nil {
				if verboseT {
					Pl("Dial failed: %v", err)
				}
				return
			}

			if verboseT {
				Pl("Connected on %v from %v", connA.LocalAddr(), conn.RemoteAddr())
			}

			go func() {
				defer client.Close()
				defer conn.Close()
				io.Copy(client, conn)
			}()
			go func() {
				defer client.Close()
				defer conn.Close()
				io.Copy(conn, client)
			}()

		}(conn, remoteA)
	}
}

var StartTransparentProxy2 = TKX.StartTransparentProxy2

// var CreateTXCollection = TKX.CreateTXCollection
// var Trim = TKX.Trim
// var RegMatch = TKX.RegMatch
// var FindFirstDiffIndex = TKX.FindFirstDiffIndex
// var StartsWith = TKX.StartsWith
// var IsErrorString = TKX.IsErrorString
// var GetErrorString = TKX.GetErrorString
// var RegReplace = TKX.RegReplace
// var GenerateErrorString = TKX.GenerateErrorString
// var EncodeStringSimple = TKX.EncodeStringSimple
// var Replace = TKX.Replace
// var ObjectToJSON = TKX.ObjectToJSON
// var SplitLines = TKX.SplitLines
// var IsErrStr = TKX.IsErrStr
// var Pl = TKX.Pl
// var Pl = TKX.Pl

// Client holds the net conn and read/write buffer objects.
type Connection struct {
	conn   io.ReadWriteCloser
	Reader *textproto.Reader
	Writer *textproto.Writer
}

// NewConnection initializes a connection.
func NewConnection(conn io.ReadWriteCloser) *Connection {
	return &Connection{
		conn,
		textproto.NewReader(bufio.NewReader(conn)),
		textproto.NewWriter(bufio.NewWriter(conn)),
	}
}

// Close closes a connection.
func (c *Connection) Close() error {
	return c.conn.Close()
}

// Cmd sends the given command on the connection.
func (c *Connection) Cmd(format string, args ...interface{}) (string, error) {
	if err := c.Writer.PrintfLine(format, args...); err != nil {
		return "", fmt.Errorf("failed to write with format and args: %w", err)
	}

	return c.ReadLine()
}

// ReadLine reads a single line from the buffer.
func (c *Connection) ReadLine() (string, error) {
	line, err := c.Reader.ReadLine()
	if err != nil {
		return "", fmt.Errorf("failed to read line: %w", err)
	}

	if len(line) < 1 {
		return "", errors.New("empty response")
	}

	if IsPop3Err(line) {
		return line, fmt.Errorf("something went wrong: %s", line)
	}

	return line, nil
}

// ReadLines reads from the buffer until it hits the message end dot (".").
func (c *Connection) ReadLines() (lines []string, err error) {
	for {
		line, err := c.ReadLine()
		if err != nil {
			return nil, fmt.Errorf("failed to read line: %w", err)
		}

		// Look for a dot to indicate the end of a message
		// from the server.
		if line == "." {
			break
		}
		lines = append(lines, line)
	}
	return
}

// POP3 replies as extracted from rfc1939 section 9.
// const (
// 	OK  = "+OK"
// 	ERR = "-ERR"
// )

// IsPop3OK checks to see if the reply from the server contains +OK.
func IsPop3OK(s string) bool {
	if strings.Fields(s)[0] != "+OK" {
		return false
	}
	return true
}

// IsPop3Err checks to see if the reply from the server contains +Err.
func IsPop3Err(s string) bool {
	if strings.Fields(s)[0] != "-ERR" {
		return false
	}
	return true
}

// MessageList represents the metadata returned by the server for a
// message stored in the maildrop.
type MessageList struct {
	// Non unique id reported by the server
	ID int

	// Size of the message
	Size int
}

const (
	CommandReset = "RSET"

	// CommandStat is a command to retrieve statistics about mailbox.
	CommandStat = "STAT"

	// CommandDelete is a command to delete message from POP3 server.
	CommandDelete = "DELE"

	// CommandList is a command to get list of messages from POP3 server.
	CommandList = "LIST"

	// CommandUIDL is a command to get list of messages from POP3 server.
	CommandUIDL = "UIDL"

	// CommandNoop is a ping-like command that tells POP3 to do nothing.
	// (i.e. send something line pong-response).
	CommandNoop = "NOOP"

	// CommandPassword is a command to send user password to POP3 server.
	CommandPassword = "PASS"

	// CommandQuit is a command to tell POP3 server that you are quitting.
	CommandQuit = "QUIT"

	// CommandRetrieve is a command to retrieve POP3 message from server.
	CommandRetrieve = "RETR"

	// CommandUser is a command to send user login to POP3 server.
	CommandUser = "USER"
)

// Client for POP3.
type Pop3Client struct {
	conn *Connection
}

func (pA *TK) ConnectPop3(addrA string, userNameA string, passwordA string, optsA ...string) (*Pop3Client, error) {
	secureT := !IfSwitchExists(optsA, "-insecure")

	var clientT *Pop3Client = nil
	var errT error

	if secureT {
		clientT, errT = DialPop3TLS(addrA)

	} else {
		clientT, errT = DialPop3(addrA)
	}

	if errT != nil {
		return nil, errT
	}

	errT = clientT.Authorization(userNameA, passwordA)

	if errT != nil {
		return nil, errT
	}

	return clientT, nil
}

var ConnectPop3 = TKX.ConnectPop3

func (c *Pop3Client) GetMessage(idxT int, optsA ...string) (map[string]string, error) {
	msgT, errT := c.Retr(idxT)
	if errT != nil {
		return nil, errT
	}

	// // A list of headers is retrieved via Envelope.GetHeaderKeys().
	// headers := message.GetHeaderKeys()
	// sort.Strings(headers)

	// // Print each header, key and value.
	// for _, header := range headers {
	// 	fmt.Printf("%s: %v\n", header, message.GetHeader(header))
	// }

	rsMapT := make(map[string]string, 10)

	rsMapT["Content-Type"] = msgT.GetHeader("Content-Type")
	rsMapT["Subject"] = msgT.GetHeader("Subject")
	rsMapT["Date"] = msgT.GetHeader("Date")
	rsMapT["From"] = msgT.GetHeader("From")
	rsMapT["Message-Id"] = msgT.GetHeader("Message-Id")
	rsMapT["Mime-Version"] = msgT.GetHeader("Mime-Version")
	rsMapT["Received"] = msgT.GetHeader("Received")
	rsMapT["To"] = msgT.GetHeader("To")
	rsMapT["Cc"] = msgT.GetHeader("Cc")
	rsMapT["X-Mailer"] = msgT.GetHeader("X-Mailer")
	rsMapT["X-Originating-Ip"] = msgT.GetHeader("X-Originating-Ip")
	rsMapT["X-Priority"] = msgT.GetHeader("X-Priority")
	rsMapT["Text"] = msgT.Text
	rsMapT["HTML"] = msgT.HTML

	return rsMapT, nil
}

// DialPop3 opens new connection and creates a new POP3 client.
func (pA *TK) DialPop3(addr string) (c *Pop3Client, err error) {
	var conn net.Conn
	if conn, err = net.Dial("tcp", addr); err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}

	return NewPop3Client(conn)
}

var DialPop3 = TKX.DialPop3

// DialPop3TLS opens new TLS connection and creates a new POP3 Pop3Client.
func (pA *TK) DialPop3TLS(addr string) (c *Pop3Client, err error) {
	var conn *tls.Conn
	if conn, err = tls.Dial("tcp", addr, &tls.Config{InsecureSkipVerify: true}); err != nil {
		return nil, fmt.Errorf("failed to dial tls: %w", err)
	}
	return NewPop3Client(conn)
}

var DialPop3TLS = TKX.DialPop3TLS

// NewPop3Client creates a new POP3 Pop3Client.
func (pA *TK) NewPop3Client(conn net.Conn) (*Pop3Client, error) {
	c := &Pop3Client{
		conn: NewConnection(conn),
	}

	// Make sure we receive the server greeting
	line, err := c.conn.ReadLine()
	if err != nil {
		return nil, fmt.Errorf("failed to read line: %w", err)
	}

	if !IsPop3OK(line) {
		return nil, fmt.Errorf("server did not response with +OK: %s", line)
	}

	return c, nil
}

var NewPop3Client = TKX.NewPop3Client

// Authorization logs into POP3 server with login and password.
func (c *Pop3Client) Authorization(user, pass string) error {
	if _, err := c.conn.Cmd("%s %s", CommandUser, user); err != nil {
		return fmt.Errorf("failed at USER command: %w", err)
	}

	if _, err := c.conn.Cmd("%s %s", CommandPassword, pass); err != nil {
		return fmt.Errorf("failed at PASS command: %w", err)
	}

	return c.Noop()
}

// Quit sends the QUIT message to the POP3 server and closes the connection.
func (c *Pop3Client) Quit() error {
	if _, err := c.conn.Cmd(CommandQuit); err != nil {
		return fmt.Errorf("failed at QUIT command: %w", err)
	}

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}

	return nil
}

// Noop will do nothing however can prolong the end of a connection.
func (c *Pop3Client) Noop() error {
	if _, err := c.conn.Cmd(CommandNoop); err != nil {
		return fmt.Errorf("failed at NOOP command: %w", err)
	}

	return nil
}

// Stat retrieves a drop listing for the current maildrop, consisting of the
// number of messages and the total size (in octets) of the maildrop.
// In the event of an error, all returned numeric values will be 0.
func (c *Pop3Client) Stat() (count, size int, err error) {
	line, err := c.conn.Cmd(CommandStat)
	if err != nil {
		return
	}

	if len(strings.Fields(line)) != 3 {
		return 0, 0, fmt.Errorf("invalid response returned from server: %s", line)
	}

	// Number of messages in maildrop
	count, err = strconv.Atoi(strings.Fields(line)[1])
	if err != nil {
		return
	}
	if count == 0 {
		return
	}

	// Total size of messages in bytes
	size, err = strconv.Atoi(strings.Fields(line)[2])
	if err != nil {
		return
	}
	if size == 0 {
		return
	}
	return
}

func (c *Pop3Client) GetCount() (count int, err error) {
	line, err := c.conn.Cmd(CommandStat)
	if err != nil {
		return
	}

	if len(strings.Fields(line)) < 2 {
		return 0, fmt.Errorf("invalid response returned from server: %s", line)
	}

	// Number of messages in maildrop
	count, err = strconv.Atoi(strings.Fields(line)[1])
	if err != nil {
		return
	}
	if count == 0 {
		return
	}

	return
}

// ListAll returns a MessageList object which contains all messages in the maildrop.
func (c *Pop3Client) ListAll() (list []MessageList, err error) {
	if _, err = c.conn.Cmd(CommandList); err != nil {
		return
	}

	lines, err := c.conn.ReadLines()
	if err != nil {
		return
	}

	for _, v := range lines {
		id, err := strconv.Atoi(strings.Fields(v)[0])
		if err != nil {
			return nil, err
		}

		size, err := strconv.Atoi(strings.Fields(v)[1])
		if err != nil {
			return nil, err
		}
		list = append(list, MessageList{id, size})
	}
	return
}

func (c *Pop3Client) ListMessages(argsA ...string) (list []map[string]string, err error) {
	nT := ""
	if len(argsA) > 0 {
		nT = " " + argsA[0]
	}

	if _, err = c.conn.Cmd("LIST" + nT); err != nil {
		return
	}

	lines, err := c.conn.ReadLines()
	if err != nil {
		return
	}

	for _, v := range lines {
		// id, err := strconv.Atoi(strings.Fields(v)[0])
		// if err != nil {
		// 	return nil, err
		// }

		// size, err := strconv.Atoi(strings.Fields(v)[1])
		// if err != nil {
		// 	return nil, err
		// }
		list = append(list, map[string]string{"Index": strings.Fields(v)[0], "Size": strings.Fields(v)[0]})
	}
	return
}

// ListUIDAll returns a map[string]string object which contains all messages in the maildrop.
func (c *Pop3Client) ListUIDAll() (list []map[string]string, err error) {
	if _, err = c.conn.Cmd(CommandUIDL); err != nil {
		return
	}

	lines, err := c.conn.ReadLines()
	if err != nil {
		return
	}

	for _, v := range lines {
		// list = append(list, map[string]string{"Data": v})
		// indexT, err := strconv.Atoi(strings.Fields(v)[0])
		// if err != nil {
		// 	return nil, err
		// }

		// uidT, err := strconv.Atoi(strings.Fields(v)[1])
		// if err != nil {
		// 	return nil, err
		// }

		list = append(list, map[string]string{"Index": strings.Fields(v)[0], "Uid": strings.Fields(v)[1]})
	}
	return
}

func (c *Pop3Client) GetUID(idxA int) (uid string, err error) {
	if _, err = c.conn.Cmd(CommandUIDL, IntToStr(idxA)); err != nil {
		return
	}

	lines, err := c.conn.ReadLines()
	if err != nil {
		return
	}

	for _, v := range lines {
		// list = append(list, map[string]string{"Data": v})
		// indexT, err := strconv.Atoi(strings.Fields(v)[0])
		// if err != nil {
		// 	return nil, err
		// }

		// uidT, err := strconv.Atoi(strings.Fields(v)[1])
		// if err != nil {
		// 	return nil, err
		// }
		uid = strings.Fields(v)[1]
		return
	}
	return
}

// SendCommand returns a []string object which contains all response lines.
func (c *Pop3Client) SendCommand(cmdA string, argsA ...interface{}) ([]string, error) {
	if _, err := c.conn.Cmd(cmdA, argsA...); err != nil {
		return nil, err
	}

	lines, err := c.conn.ReadLines()
	if err != nil {
		return nil, err
	}

	return lines, nil
}

// Rset will unmark any messages that have being marked for deletion in
// the current session.
func (c *Pop3Client) Rset() error {
	if _, err := c.conn.Cmd(CommandReset); err != nil {
		return fmt.Errorf("failed at RSET command: %w", err)
	}
	return nil
}

// Retr downloads the given message and returns it as a mail.Message object.
func (c *Pop3Client) Retr(msg int) (*enmime.Envelope, error) {
	if _, err := c.conn.Cmd("%s %d", CommandRetrieve, msg); err != nil {
		return nil, fmt.Errorf("failed at RETR command: %w", err)
	}

	message, err := enmime.ReadEnvelope(c.conn.Reader.DotReader())
	if err != nil {
		return nil, fmt.Errorf("failed to read message: %w", err)
	}

	return message, nil
}

// Dele will delete the given message from the maildrop.
// Changes will only take affect after the Quit command is issued.
func (c *Pop3Client) Dele(msg int) error {
	if _, err := c.conn.Cmd("%s %d", CommandDelete, msg); err != nil {
		return fmt.Errorf("failed at DELE command: %w", err)
	}
	return nil
}

const fileScheme = "file"

// URI represents the full URI for a file.
type URI string

func (uri URI) IsFile() bool {
	return strings.HasPrefix(string(uri), "file://")
}

// Filename returns the file path for the given URI.
// It is an error to call this on a URI that is not a valid filename.
func (uri URI) Filename() string {
	filename, err := filename(uri)
	if err != nil {
		panic(err)
	}
	return filepath.FromSlash(filename)
}

func filename(uri URI) (string, error) {
	if uri == "" {
		return "", nil
	}
	u, err := url.ParseRequestURI(string(uri))
	if err != nil {
		return "", err
	}
	if u.Scheme != fileScheme {
		return "", fmt.Errorf("only file URIs are supported, got %q from %q", u.Scheme, uri)
	}
	// If the URI is a Windows URI, we trim the leading "/" and uppercase
	// the drive letter, which will never be case sensitive.
	if isWindowsDriveURIPath(u.Path) {
		u.Path = strings.ToUpper(string(u.Path[1])) + u.Path[2:]
	}
	return u.Path, nil
}

func URIFromURI(s string) URI {
	if !strings.HasPrefix(s, "file://") {
		return URI(s)
	}

	if !strings.HasPrefix(s, "file:///") {
		// VS Code sends URLs with only two slashes, which are invalid. golang/go#39789.
		s = "file:///" + s[len("file://"):]
	}
	// Even though the input is a URI, it may not be in canonical form. VS Code
	// in particular over-escapes :, @, etc. Unescape and re-encode to canonicalize.
	path, err := url.PathUnescape(s[len("file://"):])
	if err != nil {
		panic(err)
	}

	// File URIs from Windows may have lowercase drive letters.
	// Since drive letters are guaranteed to be case insensitive,
	// we change them to uppercase to remain consistent.
	// For example, file:///c:/x/y/z becomes file:///C:/x/y/z.
	if isWindowsDriveURIPath(path) {
		path = path[:1] + strings.ToUpper(string(path[1])) + path[2:]
	}
	u := url.URL{Scheme: fileScheme, Path: path}
	return URI(u.String())
}

func CompareURI(a, b URI) int {
	if equalURI(a, b) {
		return 0
	}
	if a < b {
		return -1
	}
	return 1
}

func equalURI(a, b URI) bool {
	if a == b {
		return true
	}
	// If we have the same URI basename, we may still have the same file URIs.
	if !strings.EqualFold(path.Base(string(a)), path.Base(string(b))) {
		return false
	}
	fa, err := filename(a)
	if err != nil {
		return false
	}
	fb, err := filename(b)
	if err != nil {
		return false
	}
	// Stat the files to check if they are equal.
	infoa, err := os.Stat(filepath.FromSlash(fa))
	if err != nil {
		return false
	}
	infob, err := os.Stat(filepath.FromSlash(fb))
	if err != nil {
		return false
	}
	return os.SameFile(infoa, infob)
}

// URIFromPath returns a span URI for the supplied file path.
// It will always have the file scheme.
func URIFromPath(path string) URI {
	if path == "" {
		return ""
	}
	// Handle standard library paths that contain the literal "$GOROOT".
	// TODO(rstambler): The go/packages API should allow one to determine a user's $GOROOT.
	// const prefix = "$GOROOT"
	// if len(path) >= len(prefix) && strings.EqualFold(prefix, path[:len(prefix)]) {
	// 	suffix := path[len(prefix):]
	// 	path = runtime.GOROOT() + suffix
	// }
	if !isWindowsDrivePath(path) {
		if abs, err := filepath.Abs(path); err == nil {
			path = abs
		}
	}
	// Check the file path again, in case it became absolute.
	if isWindowsDrivePath(path) {
		path = "/" + strings.ToUpper(string(path[0])) + path[1:]
	}
	path = filepath.ToSlash(path)
	u := url.URL{
		Scheme: fileScheme,
		Path:   path,
	}
	return URI(u.String())
}

func (pA *TK) PathToURI(path string) string {
	if path == "" {
		return ""
	}

	if !isWindowsDrivePath(path) {
		if abs, err := filepath.Abs(path); err == nil {
			path = abs
		}
	}

	if isWindowsDrivePath(path) {
		path = "/" + strings.ToUpper(string(path[0])) + path[1:]
	}

	path = filepath.ToSlash(path)
	u := url.URL{
		Scheme: fileScheme,
		Path:   path,
	}

	return u.String()
}

var PathToURI = TKX.PathToURI

func getUrifilename(uri string) string {
	if uri == "" {
		return ""
	}
	u, err := url.ParseRequestURI(uri)
	if err != nil {
		return ErrToStr(err)
	}
	if u.Scheme != fileScheme {
		return ErrStrf("only file URIs are supported, got %q from %q", u.Scheme, uri)
	}
	// If the URI is a Windows URI, we trim the leading "/" and uppercase
	// the drive letter, which will never be case sensitive.
	if isWindowsDriveURIPath(u.Path) {
		u.Path = strings.ToUpper(string(u.Path[1])) + u.Path[2:]
	}
	return u.Path
}

func (pA *TK) URIToPath(uri string) string {
	if uri == "" {
		return ""
	}
	u, err := url.ParseRequestURI(uri)
	if err != nil {
		return ErrToStr(err)
	}
	if u.Scheme != fileScheme {
		return ErrStrf("only file URIs are supported, got %q from %q", u.Scheme, uri)
	}
	// If the URI is a Windows URI, we trim the leading "/" and uppercase
	// the drive letter, which will never be case sensitive.
	if isWindowsDriveURIPath(u.Path) {
		u.Path = strings.ToUpper(string(u.Path[1])) + u.Path[2:]
	}

	return filepath.FromSlash(u.Path)
}

var URIToPath = TKX.URIToPath

// isWindowsDrivePath returns true if the file path is of the form used by
// Windows. We check if the path begins with a drive letter, followed by a ":".
// For example: C:/x/y/z.
func isWindowsDrivePath(path string) bool {
	if len(path) < 3 {
		return false
	}
	return unicode.IsLetter(rune(path[0])) && path[1] == ':'
}

// isWindowsDriveURI returns true if the file URI is of the format used by
// Windows URIs. The url.Parse package does not specially handle Windows paths
// (see golang/go#6027), so we check if the URI path has a drive prefix (e.g. "/C:").
func isWindowsDriveURIPath(uri string) bool {
	if len(uri) < 4 {
		return false
	}
	return uri[0] == '/' && unicode.IsLetter(rune(uri[1])) && uri[2] == ':'
}

// start of encoding xml from map[string]interface{}

type Indentation struct {
	Prefix string
	Indent string
}

type Root struct {
	Name          *xml.Name
	XMLAttributes *[]xml.Attr
	Attributes    map[string]string
}

type StructMap struct {
	CData  bool
	Map    map[string]interface{}
	Indent *Indentation
	Root   *Root
}

type StructAny struct {
	CData  bool
	Value  interface{}
	Indent *Indentation
	Root   *Root
}

type xmlMapEntry struct {
	XMLName    xml.Name
	Value      interface{} `xml:",innerxml"`
	CDataValue interface{} `xml:",cdata"`
}

// Initializes the builder. Required to do anything with this library
func (pA *TK) NewXMLFromMSI(input map[string]interface{}) *StructMap {
	return &StructMap{Map: input}
}

var NewXMLFromMSI = TKX.NewXMLFromMSI

func (pA *TK) NewXMLFromAny(input interface{}) *StructAny {
	return &StructAny{Value: input}
}

var NewXMLFromAny = TKX.NewXMLFromAny

// Add indentation to your xml
func (smap *StructMap) WithIndent(prefix string, indent string) *StructMap {
	smap.Indent = &Indentation{Prefix: prefix, Indent: indent}
	return smap
}

func (p *StructAny) WithIndent(prefix string, indent string) *StructAny {
	p.Indent = &Indentation{Prefix: prefix, Indent: indent}
	return p
}

// Add a root node to your xml
func (smap *StructMap) WithRoot(name string, attributes map[string]string) *StructMap {
	attr := []xml.Attr{}
	for k, v := range attributes {
		attr = append(attr, xml.Attr{Name: xml.Name{Local: k}, Value: v})
	}
	smap.Root = &Root{Name: &xml.Name{Local: name}, XMLAttributes: &attr, Attributes: attributes}
	return smap
}

func (p *StructAny) WithRoot(name string, attributes map[string]string) *StructAny {
	attr := []xml.Attr{}
	for k, v := range attributes {
		attr = append(attr, xml.Attr{Name: xml.Name{Local: k}, Value: v})
	}
	p.Root = &Root{Name: &xml.Name{Local: name}, XMLAttributes: &attr, Attributes: attributes}
	return p
}

// Add CDATA tags to all nodes
func (smap *StructMap) AsCData() *StructMap {
	smap.CData = true
	return smap
}

func (p *StructAny) AsCData() *StructAny {
	p.CData = true
	return p
}

// Prints your configuration in json format
func (smap *StructMap) Print() *StructMap {
	var indent interface{} = smap.Indent
	var root interface{}
	if smap.Indent != nil {
		indent = map[string]int{"indent_spaces": len(*&smap.Indent.Indent), "prefix_spaces": len(*&smap.Indent.Prefix)}
	}
	if root = smap.Root; root != nil {
		root = map[string]interface{}{"name": *&smap.Root.Name.Local, "attributes": smap.Root.Attributes}
	}
	b, _ := json.MarshalIndent(map[string]interface{}{"root": root, "cdata": smap.CData, "indent": indent}, " ", "  ")
	fmt.Println(string(b))
	return smap
}

// Builds XML as bytes
func (smap *StructMap) Marshal() ([]byte, error) {
	if smap.Indent == nil {
		return xml.Marshal(smap)
	} else {
		return xml.MarshalIndent(smap, smap.Indent.Prefix, smap.Indent.Indent)
	}
}

func (p *StructAny) Marshal() ([]byte, error) {
	if p.Indent == nil {
		return xml.Marshal(p)
	} else {
		return xml.MarshalIndent(p, p.Indent.Prefix, p.Indent.Indent)
	}
}

// Builds XML as string
func (smap *StructMap) MarshalToString() (string, error) {
	xmlBytes, err := smap.Marshal()
	return string(xmlBytes), err
}

func (p *StructAny) MarshalToString() (string, error) {
	xmlBytes, err := p.Marshal()
	return string(xmlBytes), err
}

func (smap StructMap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if len(smap.Map) == 0 {
		return nil
	}
	if smap.Root != nil {
		start = xml.StartElement{Name: *smap.Root.Name, Attr: *smap.Root.XMLAttributes}
		if err := e.EncodeToken(start); err != nil {
			return err
		}
	}

	for k, v := range smap.Map {
		if err := handleChildren(e, k, v, smap.CData); err != nil {
			return err
		}
	}

	if smap.Root != nil {
		return e.EncodeToken(start.End())
	}
	return nil
}

func (p StructAny) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if p.Value == nil {
		return nil
	}

	if p.Root != nil {
		start = xml.StartElement{Name: *p.Root.Name, Attr: *p.Root.XMLAttributes}
		if err := e.EncodeToken(start); err != nil {
			return err
		}
	}

	// for k, v := range p.Value {
	err := handleChildren(e, "", p.Value, p.CData)
	if err != nil {
		return err
	}
	// }

	if p.Root != nil {
		return e.EncodeToken(start.End())
	}

	return nil
}

func handleChildren(e *xml.Encoder, fieldName string, v interface{}, cdata bool) error {
	var errT error

	if reflect.TypeOf(v) == nil {
		if fieldName == "" {
			return nil
		}

		return e.Encode(xmlMapEntry{XMLName: xml.Name{Local: fieldName}, Value: ""})
	} else if reflect.TypeOf(v).Kind() == reflect.Map {
		if fieldName != "" {
			errT = e.EncodeToken(xml.StartElement{Name: xml.Name{Local: fieldName}})

			if errT != nil {
				return errT
			}
		}

		for key, val := range v.(map[string]interface{}) {
			if key == "xml_child_name" {
				continue
			}

			errT = handleChildren(e, key, val, cdata)

			if errT != nil {
				return errT
			}
		}

		if fieldName != "" {
			return e.EncodeToken(xml.EndElement{Name: xml.Name{Local: fieldName}})
		}

		return nil

	} else if reflect.TypeOf(v).Kind() == reflect.Slice {
		if fieldName != "" {
			e.EncodeToken(xml.StartElement{Name: xml.Name{Local: fieldName}})
		}
		// childName := fieldName + "_child"
		childName := "item"
		vs, ok := v.([]interface{})
		if ok {
			// if _, hasChildName := vs[0].(map[string]interface{})["xml_child_name"]; hasChildName {
			// 	childName = vs[0].(map[string]interface{})["xml_child_name"].(string)
			// }
			for _, elem := range vs {
				handleChildren(e, childName, elem, cdata)
			}
		} else {
			if _, hasChildName := v.([]map[string]interface{})[0]["xml_child_name"]; hasChildName {
				childName = v.([]map[string]interface{})[0]["xml_child_name"].(string)
			}
			for _, elem := range v.([]map[string]interface{}) {
				handleChildren(e, childName, elem, cdata)
			}
		}
		if fieldName != "" {
			return e.EncodeToken(xml.EndElement{Name: xml.Name{Local: fieldName}})
		}

		return nil
	}

	if cdata {
		if fieldName != "" {
			return e.Encode(xmlMapEntry{XMLName: xml.Name{Local: fieldName}, CDataValue: v})
		}

		return e.Encode(v)
	} else {
		if fieldName != "" {
			return e.Encode(xmlMapEntry{XMLName: xml.Name{Local: fieldName}, Value: fmt.Sprintf("%v", v)})
		}

		return e.Encode(v)
	}
}

// Data struct related

type CountingWriter struct {
	Count      int
	Total      int
	IfPercent  bool
	IfRaw      bool
	WritebackI *int
	WritebackS *string
	WritebackA *interface{}
	Callback   func(interface{}) interface{}
	Lock       *sync.Mutex
}

// func NewCountingWriter(lockA *sync.Mutex, writebackA *string, writebackIA *interface{}, callbackA func(interface{}) interface{}) io.Writer {
func NewCountingWriter(argsA ...interface{}) io.Writer {
	vT := &CountingWriter{}
	argsT := make([]string, 0, len(argsA))
	for _, v := range argsA {
		if v == nil {
			continue
		}

		if nv, ok := v.(string); ok {
			argsT = append(argsT, nv)
			continue
		}

		if nv, ok := v.(int); ok {
			vT.Total = nv
			continue
		}

		if nv, ok := v.(int64); ok {
			vT.Total = int(nv)
			continue
		}

		if nv, ok := v.(*int); ok {
			vT.WritebackI = nv
			continue
		}

		if nv, ok := v.(*string); ok {
			vT.WritebackS = nv
			continue
		}

		if nv, ok := v.(*interface{}); ok {
			vT.WritebackA = nv
			continue
		}

		if nv, ok := v.(func(interface{}) interface{}); ok {
			vT.Callback = nv
			continue
		}

		if nv, ok := v.(*sync.Mutex); ok {
			vT.Lock = nv
			continue
		}
	}

	if IfSwitchExistsWhole(argsT, "-percent") {
		vT.IfPercent = true
	}

	if IfSwitchExistsWhole(argsT, "-raw") {
		vT.IfRaw = true
	}

	totalT := GetSwitch(argsT, "-total=", "")

	if totalT != "" {
		vT.Total = StrToIntWithDefaultValue(totalT, 0)
	}

	return vT // &CountingWriter{Lock: lockA, Writeback: writebackA, WritebackI: writebackIA, Callback: callbackA, Count: 0}
}

func (pA *CountingWriter) Reset() {
	if pA.Lock != nil {
		pA.Lock.Lock()
	}

	pA.Count = 0

	if pA.Lock != nil {
		pA.Lock.Unlock()
	}
}

// func (pA *CountingWriter) SetCallback(funcA func(interface{}) interface{}) {
// 	pA.Callback = funcA
// }

// func (pA *CountingWriter) SetWriteback(writebackA *string) {
// 	pA.Writeback = writebackA
// }

func (pA *CountingWriter) Write(p []byte) (n int, err error) {
	lenT := len(p)

	if pA.Lock != nil {
		pA.Lock.Lock()
	}

	pA.Count += lenT

	var infoT string

	percentT := ""

	if pA.IfPercent {
		if pA.Total > 0 {
			percentT = fmt.Sprintf("(%v%%)", pA.Count*100/pA.Total)
		}
	}

	if !pA.IfRaw {
		infoT = fmt.Sprintf("%v/%v%v", IntToKMGT(pA.Count), IntToKMGT(pA.Total), percentT)
	} else {
		infoT = fmt.Sprintf("%v/%v%v", pA.Count, pA.Total, percentT)
	}

	if pA.Callback != nil {
		pA.Callback(infoT)
	}

	if pA.WritebackS != nil {
		*(pA.WritebackS) = infoT
	}

	if pA.WritebackI != nil {
		if pA.IfPercent {
			if pA.Total >= 0 {
				*(pA.WritebackI) = pA.Count * 100 / pA.Total
			} else {
				*(pA.WritebackI) = 0
			}
		} else {
			*(pA.WritebackI) = pA.Count
		}
	}

	if pA.WritebackA != nil {
		*(pA.WritebackA) = infoT
	}

	if pA.Lock != nil {
		pA.Lock.Unlock()
	}

	return lenT, nil
}

// SyncQueue
type SyncQueue struct {
	Items *doublylinkedlist.List
	Lock  sync.Mutex
	Cap   int
}

// no size or <= 0 indicates no size limit
func (pA *TK) NewSyncQueue(sizeA ...int) *SyncQueue {
	sizeT := 0
	if len(sizeA) > 0 {
		sizeT = sizeA[0]
	}

	rs := &SyncQueue{Items: doublylinkedlist.New(), Cap: sizeT}

	return rs
}

var NewSyncQueue = TKX.NewSyncQueue

func (p *SyncQueue) Add(vA interface{}) {
	p.Lock.Lock()

	p.Items.Add(vA)

	if p.Cap > 0 {
		if p.Items.Size() > p.Cap {
			p.Items.Remove(0)
		}
	}

	p.Lock.Unlock()
}

func (p *SyncQueue) RemoveLast() {
	p.Lock.Lock()

	lenT := p.Items.Size()

	if lenT > 0 {
		p.Items.Remove(lenT - 1)
	}

	p.Lock.Unlock()
}

func (p *SyncQueue) ReplaceFirst(vA interface{}) {
	p.Lock.Lock()

	lenT := p.Items.Size()

	if lenT > 0 {
		p.Items.Set(0, vA)
	}

	p.Lock.Unlock()
}

func (p *SyncQueue) ReplaceLast(vA interface{}) {
	p.Lock.Lock()

	lenT := p.Items.Size()

	if lenT > 0 {
		p.Items.Set(lenT-1, vA)
	}

	p.Lock.Unlock()
}

func (p *SyncQueue) Clear() {
	p.Lock.Lock()

	p.Items.Clear()

	p.Lock.Unlock()
}

func (p *SyncQueue) ClearAdd(vA interface{}) {
	p.Lock.Lock()

	p.Items.Clear()

	p.Items.Add(vA)

	p.Lock.Unlock()
}

func (p *SyncQueue) QuickGet() interface{} {
	var result interface{} = nil
	var b bool

	p.Lock.Lock()

	result, b = p.Items.Get(0)

	if b {
		p.Items.Remove(0)
	}

	p.Lock.Unlock()

	if b {
		return result
	}

	return nil
}

func (p *SyncQueue) QuickGetList() []interface{} {
	var result []interface{} = nil
	var b bool
	var vT interface{}

	p.Lock.Lock()

	lenT := p.Items.Size()

	result = make([]interface{}, lenT)

	for i := 0; i < lenT; i++ {
		vT, b = p.Items.Get(i)

		if b {
			result[i] = vT
		} else {
			result[i] = nil
		}
	}

	p.Lock.Unlock()

	return result
}

func (p *SyncQueue) Get() (interface{}, bool) {
	var result interface{} = nil
	var b bool

	p.Lock.Lock()

	result, b = p.Items.Get(0)

	if b {
		// fmt.Printf("result: %v, b: %v\n", result, b)
		p.Items.Remove(0)
	}

	p.Lock.Unlock()

	return result, b
}

func (p *SyncQueue) GetCompact() interface{} {
	var result interface{} = nil
	var b bool

	p.Lock.Lock()

	result, b = p.Items.Get(0)

	if b {
		// fmt.Printf("result: %v, b: %v\n", result, b)
		p.Items.Remove(0)
	} else {
		result = nil
	}

	p.Lock.Unlock()

	return result
}

func (p *SyncQueue) PeekLast() (interface{}, bool) {
	var result interface{} = nil
	var b bool

	p.Lock.Lock()

	lenT := p.Items.Size()

	if lenT < 1 {
		p.Lock.Unlock()
		return nil, false
	}

	result, b = p.Items.Get(lenT - 1)

	// if b {
	// 	p.Items.Remove(0)
	// }

	p.Lock.Unlock()

	return result, b
}

func (p *SyncQueue) PeekLastCompact() interface{} {
	var result interface{} = nil
	var b bool

	p.Lock.Lock()

	lenT := p.Items.Size()

	if lenT < 1 {
		p.Lock.Unlock()
		return nil
	}

	result, b = p.Items.Get(lenT - 1)

	p.Lock.Unlock()

	if !b {
		return nil
	}

	return result
}

func (p *SyncQueue) PeekFirstCompact() interface{} {
	var result interface{} = nil
	var b bool

	p.Lock.Lock()

	lenT := p.Items.Size()

	if lenT < 1 {
		p.Lock.Unlock()
		return nil
	}

	result, b = p.Items.Get(0)

	p.Lock.Unlock()

	if !b {
		return nil
	}

	return result
}

func (p *SyncQueue) PeekCompact(idxA int) interface{} {
	var result interface{} = nil
	var b bool

	p.Lock.Lock()

	lenT := p.Items.Size()

	if idxA < 0 || idxA > (lenT-1) {
		p.Lock.Unlock()
		return nil
	}

	result, b = p.Items.Get(idxA)

	p.Lock.Unlock()

	if !b {
		return nil
	}

	return result
}

func (p *SyncQueue) Size() int {
	var result int

	p.Lock.Lock()

	result = p.Items.Size()

	p.Lock.Unlock()

	return result
}

type SimpleStack struct {
	Items        []interface{}
	Pointer      int
	DefaultValue interface{}
}

// if present, the 1st argument is the initial capacity of the stack, the 2nd is the default value return while use pop or peek actions when there are no values in the stack
func (p *SimpleStack) Reset(argsA ...interface{}) {
	sizeT := 10

	p.DefaultValue = nil

	lenT := len(argsA)

	if lenT > 0 {
		sizeT = ToInt(argsA[0])
	}

	if lenT > 1 {
		p.DefaultValue = argsA[1]
	}

	p.Items = make([]interface{}, 0, sizeT)

	p.Pointer = 0
}

func (p *SimpleStack) Clear() {
	p.Items = make([]interface{}, 0)

	p.Pointer = 0
}

func (p *SimpleStack) Push(vA interface{}) {
	// if p.Items == Undefined {
	// 	p.Reset()
	// }

	lenT := len(p.Items)

	if p.Pointer >= lenT {
		p.Items = append(p.Items, vA)
	} else {
		p.Items[p.Pointer] = vA
	}

	p.Pointer++

	// p.Items = append(p.Items, vA)
}

func (p *SimpleStack) Pop() interface{} {
	if p.Pointer < 1 {
		return p.DefaultValue
	}

	p.Pointer--
	rs := p.Items[p.Pointer]
	return rs
	// if p.Items == nil {
	// 	p.Reset()

	// 	return nil
	// }

	// lenT := len(p.Items)

	// if lenT < 1 {
	// 	return nil
	// }

	// rs := p.Items[lenT-1]

	// p.Items = p.Items[0 : lenT-1]

	// return rs
}

func (p *SimpleStack) Peek() interface{} {
	if p.Pointer < 1 {
		return p.DefaultValue
	}

	// p.StackPointerM--
	rs := p.Items[p.Pointer-1]

	return rs
	// if p.Items == nil {
	// 	p.Reset()

	// 	return nil
	// }

	// lenT := len(p.Items)

	// if lenT < 1 {
	// 	return nil
	// }

	// return p.Items[lenT-1]
}

// use -1 to get the topmost item in the stack, -2 as the second topmost item, and so on
func (p *SimpleStack) PeekLayer(idxA int) interface{} {
	if idxA < 0 {
		idxA = p.Size() - idxA
	}

	if p.Pointer < 1 {
		return p.DefaultValue
	}

	if idxA >= p.Pointer {
		return p.DefaultValue
	}

	rs := p.Items[idxA]

	return rs
}

// use index -1 to set the topmost item in the stack, -2 as the second topmost item, and so on
func (p *SimpleStack) SetByIndex(idxA int, vA interface{}) error {
	if idxA < 0 {
		idxA = p.Size() - idxA
	}

	if p.Pointer < 1 {
		return fmt.Errorf("invalid pointer: %v", p.Pointer)
	}

	if idxA >= p.Pointer {
		return fmt.Errorf("index out of range: %v/%v", idxA, p.Pointer)
	}

	p.Items[idxA] = vA

	return nil
}

func (p *SimpleStack) Size() int {
	return p.Pointer
}

func (p *SimpleStack) Reverse() *SimpleStack {
	lenT := p.Size()

	rs := &SimpleStack{Items: make([]interface{}, 0, lenT), DefaultValue: p.DefaultValue, Pointer: 0}

	for i := lenT - 1; i >= 0; i-- {
		rs.Push(p.Items[i])
	}

	return rs
}

// if present, the 1st argument is the initial capacity of the stack, the 2nd is the default value return while use pop or peek actions when there are no values in the stack
func (pA *TK) NewSimpleStack(argsA ...interface{}) *SimpleStack {
	rs := &SimpleStack{}

	rs.Reset(argsA...)

	return rs
}

var NewSimpleStack = TKX.NewSimpleStack

type SyncStack struct {
	Items        []interface{}
	Pointer      int
	DefaultValue interface{}

	Lock sync.Mutex
}

// if present, the 1st argument is the initial capacity of the stack, the 2nd is the default value return while use pop or peek actions when there are no values in the stack
func (p *SyncStack) Reset(argsA ...interface{}) {
	sizeT := 10

	lenT := len(argsA)

	if lenT > 0 {
		sizeT = ToInt(argsA[0])
	}

	p.Lock.Lock()

	p.DefaultValue = nil

	if lenT > 1 {
		p.DefaultValue = argsA[1]
	}

	p.Items = make([]interface{}, 0, sizeT)

	p.Pointer = 0

	p.Lock.Unlock()
}

func (p *SyncStack) Push(vA interface{}) {
	p.Lock.Lock()

	lenT := len(p.Items)

	if p.Pointer >= lenT {
		p.Items = append(p.Items, vA)
	} else {
		p.Items[p.Pointer] = vA
	}

	p.Pointer++

	p.Lock.Unlock()
}

func (p *SyncStack) Pop() interface{} {
	p.Lock.Lock()
	if p.Pointer < 1 {
		rs := p.DefaultValue
		p.Lock.Unlock()
		return rs
	}

	p.Pointer--
	rs := p.Items[p.Pointer]
	p.Lock.Unlock()
	return rs
}

func (p *SyncStack) Peek() interface{} {
	p.Lock.Lock()
	if p.Pointer < 1 {
		rs := p.DefaultValue
		p.Lock.Unlock()
		return rs
	}

	rs := p.Items[p.Pointer-1]
	p.Lock.Unlock()

	return rs
}

func (p *SyncStack) PeekLayer(idxA int) interface{} {
	p.Lock.Lock()
	if idxA < 0 {
		idxA = p.Pointer - idxA
	}

	if p.Pointer < 1 {
		rs := p.DefaultValue
		p.Lock.Unlock()
		return rs
	}

	if idxA >= p.Pointer {
		rs := p.DefaultValue
		p.Lock.Unlock()
		return rs
	}

	rs := p.Items[idxA]

	p.Lock.Unlock()
	return rs
}

func (p *SyncStack) Size() int {
	p.Lock.Lock()

	rs := p.Pointer

	p.Lock.Unlock()
	return rs
}

func (pA *TK) NewSyncStack(argsA ...interface{}) *SyncStack {
	rs := &SyncStack{}

	rs.Reset(argsA...)

	return rs
}

var NewSyncStack = TKX.NewSyncStack

func (pA *TK) NewObject(argsA ...interface{}) interface{} {
	lenT := len(argsA)

	if lenT < 1 {
		return Errf("not enough parameters")
	}

	typeT := ToLower(ToStr(argsA[0]))

	// need lower case
	switch typeT {
	case "nil":
		return nil
	case "undefined":
		return new(UndefinedStruct)
	case "bool":
		return new(bool)
	case "byte":
		return new(byte)
	case "int":
		return new(int)
	case "int64":
		return new(int64)
	case "rune", "int32":
		return new(rune)
	case "uint":
		return new(uint)
	case "uint64":
		return new(uint64)
	case "float", "float64":
		return new(float64)
	case "float32":
		return new(float32)
	case "string":
		return new(string)
		// case "string":
		// 	if lenT > 1 {
		// 		return CreateStringSimple(ToStr(argsA[1]))
		// 	}

		// 	return CreateStringEmpty()
	case "any", "interface{}":
		return new(interface{})
	case "bytesbuffer":
		if lenT > 1 {
			nv, ok := argsA[1].([]byte)

			if ok {
				return bytes.NewBuffer(nv)
			}

			return bytes.NewBufferString(ToStr(argsA[1]))
		}
		// var bufT bytes.Buffer
		return new(bytes.Buffer)
	case "stringbuffer", "stringbuilder", "strings.builder":

		var bufT = new(strings.Builder)
		if lenT > 1 {
			bufT.WriteString(ToStr(argsA[1]))
		}

		return bufT
	case "stringreader":
		if lenT > 1 {
			return strings.NewReader(ToStr(argsA[1]))
		}

		return strings.NewReader("")
	case "stringring":
		if lenT > 1 {
			return NewStringRing(ToInt(argsA[1]))
		}

		return NewStringRing()
	case "bytequeue":
		if lenT > 1 {
			return NewByteQueue(ToInt(argsA[1]))
		}

		return NewByteQueue()
	case "anyqueue":
		if lenT > 1 {
			return NewAnyQueue(ToInt(argsA[1]))
		}

		return NewAnyQueue()
	case "array":
		return &[]interface{}{}
	case "intarray":
		return &[]int{}
	case "floatarray":
		return &[]float64{}
	case "strarray", "stringarray":
		return &[]string{}
	case "map", "mapsi", "mapsa":
		return &map[string]interface{}{}
	case "mapss":
		return &map[string]string{}
	case "simplestack":
		return NewSimpleStack()
	case "stack":
		return linkedliststack.New()
	case "list", "arraylist":
		return arraylist.New()
	case "linklist", "linkedlist":
		return doublylinkedlist.New()
	case "set", "hashset":
		return hashset.New()
	case "treeset":
		return linkedhashset.New()
	case "syncqueue":
		return NewSyncQueue()
	case "seq":
		if lenT > 1 {
			typeT := ToStr(argsA[1])

			if typeT == "-global" {
				return AutoSeq
			}
			// return fmt.Errorf(ToStr(argsA[1]), argsA[2:]...)
		}

		return NewSeq()
	case "orderedmap":
		return NewOrderedMap()
	case "error", "err":
		if lenT > 1 {
			return fmt.Errorf(ToStr(argsA[1]), argsA[2:]...)
		}

		return fmt.Errorf("")
	case "errorstring", "errstr":
		if lenT > 1 {
			return GenerateErrorStringF(ToStr(argsA[1]), argsA[2:]...)
		}

		return GenerateErrorStringF("")
	case "txresult":
		if lenT > 2 {
			return NewTXResult(ToStr(argsA[1]), ToStr(argsA[2]))
		}

		if lenT > 1 {
			return NewTXResult(ToStr(argsA[1]), "")
		}

		return NewTXResult("success", "")
	case "flex", "simpleflexobject":
		return NewSimpleFlexObject(argsA[1:]...)
	case "tree", "btree":
		if lenT < 2 {
			return Errf("not enough parameters")
		}

		return btree.NewWithIntComparator(ToInt(argsA[1]))
	case "statichttphandler", "staticwebhandler":
		if lenT < 2 {
			return Errf("not enough parameters")
		}

		return NewStaticWebHandler(ToStr(argsA[1]))
	case "image":
		// if lenT < 2 {
		// 	return Errf("not enough parameters")
		// }

		return NewImage(AnyArrayToStringArray(argsA[1:])...)
	case "memfile":
		return NewMemFile(nil)
	case "xmltree", "xmldocument":
		treeT := etree.NewDocument()

		if treeT == nil {
			return Errf("create XML tree failed")
		}

		if lenT < 2 {
			return treeT
		}

		xmlT := strings.TrimSpace(ToStr(argsA[1]))

		if xmlT == "" {
			return treeT
		}

		errT := treeT.ReadFromString(xmlT)

		if errT != nil {
			return Errf("invalid XML: %v", errT)
		}

		return treeT
	case "xmlelement":
		if lenT < 2 {
			return Errf("not enough parameters")
		}

		elmentT := etree.NewElement(ToStr(argsA[1]))

		if elmentT == nil {
			return Errf("create XML element failed")
		}

		return elmentT
	case "tk", "tkc":
		return NewTK()
	}

	return Errf("unknown object type: %v", typeT)
}

var NewObject = TKX.NewObject

var MimeTypes = map[string]string{
	"":                         "application/octet-stream",
	"123":                      "application/vnd.lotus-1-2-3",
	"3dml":                     "text/vnd.in3d.3dml",
	"3ds":                      "image/x-3ds",
	"3g2":                      "video/3gpp2",
	"3gp":                      "video/3gpp",
	"3gpp":                     "video/3gpp",
	"7z":                       "application/x-7z-compressed",
	"aab":                      "application/x-authorware-bin",
	"aac":                      "audio/x-aac",
	"aam":                      "application/x-authorware-map",
	"aas":                      "application/x-authorware-seg",
	"abw":                      "application/x-abiword",
	"ac":                       "application/pkix-attr-cert",
	"acc":                      "application/vnd.americandynamics.acc",
	"ace":                      "application/x-ace-compressed",
	"acu":                      "application/vnd.acucobol",
	"acutc":                    "application/vnd.acucorp",
	"adp":                      "audio/adpcm",
	"aep":                      "application/vnd.audiograph",
	"afm":                      "application/x-font-type1",
	"afp":                      "application/vnd.ibm.modcap",
	"ahead":                    "application/vnd.ahead.space",
	"ai":                       "application/postscript",
	"aif":                      "audio/x-aiff",
	"aifc":                     "audio/x-aiff",
	"aiff":                     "audio/x-aiff",
	"air":                      "application/vnd.adobe.air-application-installer-package+zip",
	"ait":                      "application/vnd.dvb.ait",
	"ami":                      "application/vnd.amiga.ami",
	"apk":                      "application/vnd.android.package-archive",
	"apng":                     "image/apng",
	"appcache":                 "text/cache-manifest",
	"application":              "application/x-ms-application",
	"apr":                      "application/vnd.lotus-approach",
	"arc":                      "application/x-freearc",
	"arj":                      "application/x-arj",
	"arw":                      "image/x-sony-arw",
	"asc":                      "application/pgp-signature",
	"asf":                      "video/x-ms-asf",
	"asm":                      "text/x-asm",
	"aso":                      "application/vnd.accpac.simply.aso",
	"asx":                      "video/x-ms-asf",
	"atc":                      "application/vnd.acucorp",
	"atom":                     "application/atom+xml",
	"atomcat":                  "application/atomcat+xml",
	"atomsvc":                  "application/atomsvc+xml",
	"atx":                      "application/vnd.antix.game-component",
	"au":                       "audio/basic",
	"avi":                      "video/x-msvideo",
	"aw":                       "application/applixware",
	"azf":                      "application/vnd.airzip.filesecure.azf",
	"azs":                      "application/vnd.airzip.filesecure.azs",
	"azv":                      "image/vnd.airzip.accelerator.azv",
	"azw":                      "application/vnd.amazon.ebook",
	"bat":                      "application/x-msdownload",
	"bcpio":                    "application/x-bcpio",
	"bdf":                      "application/x-font-bdf",
	"bdm":                      "application/vnd.syncml.dm+wbxml",
	"bdoc":                     "application/x-bdoc",
	"bed":                      "application/vnd.realvnc.bed",
	"bh2":                      "application/vnd.fujitsu.oasysprs",
	"bin":                      "application/octet-stream",
	"blb":                      "application/x-blorb",
	"blorb":                    "application/x-blorb",
	"bmi":                      "application/vnd.bmi",
	"bmp":                      "image/x-ms-bmp",
	"book":                     "application/vnd.framemaker",
	"box":                      "application/vnd.previewsystems.box",
	"boz":                      "application/x-bzip2",
	"bpk":                      "application/octet-stream",
	"btif":                     "image/prs.btif",
	"buffer":                   "application/octet-stream",
	"bz":                       "application/x-bzip",
	"bz2":                      "application/x-bzip2",
	"c":                        "text/x-c",
	"c11amc":                   "application/vnd.cluetrust.cartomobile-config",
	"c11amz":                   "application/vnd.cluetrust.cartomobile-config-pkg",
	"c4d":                      "application/vnd.clonk.c4group",
	"c4f":                      "application/vnd.clonk.c4group",
	"c4g":                      "application/vnd.clonk.c4group",
	"c4p":                      "application/vnd.clonk.c4group",
	"c4u":                      "application/vnd.clonk.c4group",
	"cab":                      "application/vnd.ms-cab-compressed",
	"caf":                      "audio/x-caf",
	"cap":                      "application/vnd.tcpdump.pcap",
	"car":                      "application/vnd.curl.car",
	"cat":                      "application/vnd.ms-pki.seccat",
	"cb7":                      "application/x-cbr",
	"cba":                      "application/x-cbr",
	"cbr":                      "application/x-cbr",
	"cbt":                      "application/x-cbr",
	"cbz":                      "application/x-cbr",
	"cc":                       "text/x-c",
	"cco":                      "application/x-cocoa",
	"cct":                      "application/x-director",
	"ccxml":                    "application/ccxml+xml",
	"cdbcmsg":                  "application/vnd.contact.cmsg",
	"cdf":                      "application/x-netcdf",
	"cdkey":                    "application/vnd.mediastation.cdkey",
	"cdmia":                    "application/cdmi-capability",
	"cdmic":                    "application/cdmi-container",
	"cdmid":                    "application/cdmi-domain",
	"cdmio":                    "application/cdmi-object",
	"cdmiq":                    "application/cdmi-queue",
	"cdx":                      "chemical/x-cdx",
	"cdxml":                    "application/vnd.chemdraw+xml",
	"cdy":                      "application/vnd.cinderella",
	"cer":                      "application/pkix-cert",
	"cfs":                      "application/x-cfs-compressed",
	"cgm":                      "image/cgm",
	"chat":                     "application/x-chat",
	"chm":                      "application/vnd.ms-htmlhelp",
	"chrt":                     "application/vnd.kde.kchart",
	"cif":                      "chemical/x-cif",
	"cii":                      "application/vnd.anser-web-certificate-issue-initiation",
	"cil":                      "application/vnd.ms-artgalry",
	"cla":                      "application/vnd.claymore",
	"class":                    "application/java-vm",
	"clkk":                     "application/vnd.crick.clicker.keyboard",
	"clkp":                     "application/vnd.crick.clicker.palette",
	"clkt":                     "application/vnd.crick.clicker.template",
	"clkw":                     "application/vnd.crick.clicker.wordbank",
	"clkx":                     "application/vnd.crick.clicker",
	"clp":                      "application/x-msclip",
	"cmc":                      "application/vnd.cosmocaller",
	"cmdf":                     "chemical/x-cmdf",
	"cml":                      "chemical/x-cml",
	"cmp":                      "application/vnd.yellowriver-custom-menu",
	"cmx":                      "image/x-cmx",
	"cod":                      "application/vnd.rim.cod",
	"coffee":                   "text/coffeescript",
	"com":                      "application/x-msdownload",
	"conf":                     "text/plain",
	"cpio":                     "application/x-cpio",
	"cpp":                      "text/x-c",
	"cpt":                      "application/mac-compactpro",
	"cr2":                      "image/x-canon-cr2",
	"crd":                      "application/x-mscardfile",
	"crl":                      "application/pkix-crl",
	"crt":                      "application/x-x509-ca-cert",
	"crw":                      "image/x-canon-crw",
	"crx":                      "application/x-chrome-extension",
	"cryptonote":               "application/vnd.rig.cryptonote",
	"csh":                      "application/x-csh",
	"csl":                      "application/vnd.citationstyles.style+xml",
	"csml":                     "chemical/x-csml",
	"csp":                      "application/vnd.commonspace",
	"css":                      "text/css",
	"cst":                      "application/x-director",
	"csv":                      "text/csv",
	"cu":                       "application/cu-seeme",
	"curl":                     "text/vnd.curl",
	"cww":                      "application/prs.cww",
	"cxt":                      "application/x-director",
	"cxx":                      "text/x-c",
	"dae":                      "model/vnd.collada+xml",
	"daf":                      "application/vnd.mobius.daf",
	"dart":                     "application/vnd.dart",
	"dataless":                 "application/vnd.fdsn.seed",
	"davmount":                 "application/davmount+xml",
	"dbk":                      "application/docbook+xml",
	"dcr":                      "application/x-director",
	"dcurl":                    "text/vnd.curl.dcurl",
	"dd2":                      "application/vnd.oma.dd2+xml",
	"ddd":                      "application/vnd.fujixerox.ddd",
	"deb":                      "application/x-debian-package",
	"def":                      "text/plain",
	"deploy":                   "application/octet-stream",
	"der":                      "application/x-x509-ca-cert",
	"dfac":                     "application/vnd.dreamfactory",
	"dgc":                      "application/x-dgc-compressed",
	"dic":                      "text/x-c",
	"dir":                      "application/x-director",
	"dis":                      "application/vnd.mobius.dis",
	"disposition-notification": "message/disposition-notification",
	"dist":                     "application/octet-stream",
	"distz":                    "application/octet-stream",
	"djv":                      "image/vnd.djvu",
	"djvu":                     "image/vnd.djvu",
	"dll":                      "application/x-msdownload",
	"dmg":                      "application/x-apple-diskimage",
	"dmp":                      "application/vnd.tcpdump.pcap",
	"dms":                      "application/octet-stream",
	"dna":                      "application/vnd.dna",
	"dng":                      "image/x-adobe-dng",
	"doc":                      "application/msword",
	"docm":                     "application/vnd.ms-word.document.macroenabled.12",
	"docx":                     "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	"dot":                      "application/msword",
	"dotm":                     "application/vnd.ms-word.template.macroenabled.12",
	"dotx":                     "application/vnd.openxmlformats-officedocument.wordprocessingml.template",
	"dp":                       "application/vnd.osgi.dp",
	"dpg":                      "application/vnd.dpgraph",
	"dra":                      "audio/vnd.dra",
	"drle":                     "image/dicom-rle",
	"dsc":                      "text/prs.lines.tag",
	"dssc":                     "application/dssc+der",
	"dtb":                      "application/x-dtbook+xml",
	"dtd":                      "application/xml-dtd",
	"dts":                      "audio/vnd.dts",
	"dtshd":                    "audio/vnd.dts.hd",
	"dump":                     "application/octet-stream",
	"dvb":                      "video/vnd.dvb.file",
	"dvi":                      "application/x-dvi",
	"dwf":                      "model/vnd.dwf",
	"dwg":                      "image/vnd.dwg",
	"dxf":                      "image/vnd.dxf",
	"dxp":                      "application/vnd.spotfire.dxp",
	"dxr":                      "application/x-director",
	"ear":                      "application/java-archive",
	"ecelp4800":                "audio/vnd.nuera.ecelp4800",
	"ecelp7470":                "audio/vnd.nuera.ecelp7470",
	"ecelp9600":                "audio/vnd.nuera.ecelp9600",
	"ecma":                     "application/ecmascript",
	"edm":                      "application/vnd.novadigm.edm",
	"edx":                      "application/vnd.novadigm.edx",
	"efif":                     "application/vnd.picsel",
	"ei6":                      "application/vnd.pg.osasli",
	"elc":                      "application/octet-stream",
	"emf":                      "image/emf",
	"eml":                      "message/rfc822",
	"emma":                     "application/emma+xml",
	"emz":                      "application/x-msmetafile",
	"eol":                      "audio/vnd.digital-winds",
	"eot":                      "application/vnd.ms-fontobject",
	"eps":                      "application/postscript",
	"epub":                     "application/epub+zip",
	"erf":                      "image/x-epson-erf",
	"es":                       "application/ecmascript",
	"es3":                      "application/vnd.eszigno3+xml",
	"esa":                      "application/vnd.osgi.subsystem",
	"esf":                      "application/vnd.epson.esf",
	"et3":                      "application/vnd.eszigno3+xml",
	"etx":                      "text/x-setext",
	"eva":                      "application/x-eva",
	"evy":                      "application/x-envoy",
	"exe":                      "application/x-msdownload",
	"exi":                      "application/exi",
	"exr":                      "image/aces",
	"ext":                      "application/vnd.novadigm.ext",
	"ez":                       "application/andrew-inset",
	"ez2":                      "application/vnd.ezpix-album",
	"ez3":                      "application/vnd.ezpix-package",
	"f":                        "text/x-fortran",
	"f4v":                      "video/x-f4v",
	"f77":                      "text/x-fortran",
	"f90":                      "text/x-fortran",
	"fbs":                      "image/vnd.fastbidsheet",
	"fcdt":                     "application/vnd.adobe.formscentral.fcdt",
	"fcs":                      "application/vnd.isac.fcs",
	"fdf":                      "application/vnd.fdf",
	"fe_launch":                "application/vnd.denovo.fcselayout-link",
	"fg5":                      "application/vnd.fujitsu.oasysgp",
	"fgd":                      "application/x-director",
	"fh":                       "image/x-freehand",
	"fh4":                      "image/x-freehand",
	"fh5":                      "image/x-freehand",
	"fh7":                      "image/x-freehand",
	"fhc":                      "image/x-freehand",
	"fig":                      "application/x-xfig",
	"fits":                     "image/fits",
	"flac":                     "audio/x-flac",
	"fli":                      "video/x-fli",
	"flo":                      "application/vnd.micrografx.flo",
	"flv":                      "video/x-flv",
	"flw":                      "application/vnd.kde.kivio",
	"flx":                      "text/vnd.fmi.flexstor",
	"fly":                      "text/vnd.fly",
	"fm":                       "application/vnd.framemaker",
	"fnc":                      "application/vnd.frogans.fnc",
	"for":                      "text/x-fortran",
	"fpx":                      "image/vnd.fpx",
	"frame":                    "application/vnd.framemaker",
	"fsc":                      "application/vnd.fsc.weblaunch",
	"fst":                      "image/vnd.fst",
	"ftc":                      "application/vnd.fluxtime.clip",
	"fti":                      "application/vnd.anser-web-funds-transfer-initiation",
	"fvt":                      "video/vnd.fvt",
	"fxp":                      "application/vnd.adobe.fxp",
	"fxpl":                     "application/vnd.adobe.fxp",
	"fzs":                      "application/vnd.fuzzysheet",
	"g2w":                      "application/vnd.geoplan",
	"g3":                       "image/g3fax",
	"g3w":                      "application/vnd.geospace",
	"gac":                      "application/vnd.groove-account",
	"gam":                      "application/x-tads",
	"gbr":                      "application/rpki-ghostbusters",
	"gca":                      "application/x-gca-compressed",
	"gdl":                      "model/vnd.gdl",
	"gdoc":                     "application/vnd.google-apps.document",
	"geo":                      "application/vnd.dynageo",
	"geojson":                  "application/geo+json",
	"gex":                      "application/vnd.geometry-explorer",
	"ggb":                      "application/vnd.geogebra.file",
	"ggt":                      "application/vnd.geogebra.tool",
	"ghf":                      "application/vnd.groove-help",
	"gif":                      "image/gif",
	"gim":                      "application/vnd.groove-identity-message",
	"glb":                      "model/gltf-binary",
	"gltf":                     "model/gltf+json",
	"gml":                      "application/gml+xml",
	"gmx":                      "application/vnd.gmx",
	"gnumeric":                 "application/x-gnumeric",
	"gph":                      "application/vnd.flographit",
	"gpx":                      "application/gpx+xml",
	"gqf":                      "application/vnd.grafeq",
	"gqs":                      "application/vnd.grafeq",
	"gram":                     "application/srgs",
	"gramps":                   "application/x-gramps-xml",
	"gre":                      "application/vnd.geometry-explorer",
	"grv":                      "application/vnd.groove-injector",
	"grxml":                    "application/srgs+xml",
	"gsf":                      "application/x-font-ghostscript",
	"gsheet":                   "application/vnd.google-apps.spreadsheet",
	"gslides":                  "application/vnd.google-apps.presentation",
	"gtar":                     "application/x-gtar",
	"gtm":                      "application/vnd.groove-tool-message",
	"gtw":                      "model/vnd.gtw",
	"gv":                       "text/vnd.graphviz",
	"gxf":                      "application/gxf",
	"gxt":                      "application/vnd.geonext",
	"gz":                       "application/gzip",
	"h":                        "text/x-c",
	"h261":                     "video/h261",
	"h263":                     "video/h263",
	"h264":                     "video/h264",
	"hal":                      "application/vnd.hal+xml",
	"hbci":                     "application/vnd.hbci",
	"hbs":                      "text/x-handlebars-template",
	"hdd":                      "application/x-virtualbox-hdd",
	"hdf":                      "application/x-hdf",
	"heic":                     "image/heic",
	"heics":                    "image/heic-sequence",
	"heif":                     "image/heif",
	"heifs":                    "image/heif-sequence",
	"hh":                       "text/x-c",
	"hjson":                    "application/hjson",
	"hlp":                      "application/winhlp",
	"hpgl":                     "application/vnd.hp-hpgl",
	"hpid":                     "application/vnd.hp-hpid",
	"hps":                      "application/vnd.hp-hps",
	"hqx":                      "application/mac-binhex40",
	"htc":                      "text/x-component",
	"htke":                     "application/vnd.kenameaapp",
	"htm":                      "text/html",
	"html":                     "text/html",
	"hvd":                      "application/vnd.yamaha.hv-dic",
	"hvp":                      "application/vnd.yamaha.hv-voice",
	"hvs":                      "application/vnd.yamaha.hv-script",
	"i2g":                      "application/vnd.intergeo",
	"icc":                      "application/vnd.iccprofile",
	"ice":                      "x-conference/x-cooltalk",
	"icm":                      "application/vnd.iccprofile",
	"ico":                      "image/x-icon",
	"ics":                      "text/calendar",
	"ief":                      "image/ief",
	"ifb":                      "text/calendar",
	"ifm":                      "application/vnd.shana.informed.formdata",
	"iges":                     "model/iges",
	"igl":                      "application/vnd.igloader",
	"igm":                      "application/vnd.insors.igm",
	"igs":                      "model/iges",
	"igx":                      "application/vnd.micrografx.igx",
	"iif":                      "application/vnd.shana.informed.interchange",
	"img":                      "application/octet-stream",
	"imp":                      "application/vnd.accpac.simply.imp",
	"ims":                      "application/vnd.ms-ims",
	"in":                       "text/plain",
	"ini":                      "text/plain",
	"ink":                      "application/inkml+xml",
	"inkml":                    "application/inkml+xml",
	"install":                  "application/x-install-instructions",
	"iota":                     "application/vnd.astraea-software.iota",
	"ipfix":                    "application/ipfix",
	"ipk":                      "application/vnd.shana.informed.package",
	"irm":                      "application/vnd.ibm.rights-management",
	"irp":                      "application/vnd.irepository.package+xml",
	"iso":                      "application/x-iso9660-image",
	"itp":                      "application/vnd.shana.informed.formtemplate",
	"ivp":                      "application/vnd.immervision-ivp",
	"ivu":                      "application/vnd.immervision-ivu",
	"jad":                      "text/vnd.sun.j2me.app-descriptor",
	"jade":                     "text/jade",
	"jam":                      "application/vnd.jam",
	"jar":                      "application/java-archive",
	"jardiff":                  "application/x-java-archive-diff",
	"java":                     "text/x-java-source",
	"jisp":                     "application/vnd.jisp",
	"jls":                      "image/jls",
	"jlt":                      "application/vnd.hp-jlyt",
	"jng":                      "image/x-jng",
	"jnlp":                     "application/x-java-jnlp-file",
	"joda":                     "application/vnd.joost.joda-archive",
	"jp2":                      "image/jp2",
	"jpe":                      "image/jpeg",
	"jpeg":                     "image/jpeg",
	"jpf":                      "image/jpx",
	"jpg":                      "image/jpeg",
	"jpg2":                     "image/jp2",
	"jpgm":                     "video/jpm",
	"jpgv":                     "video/jpeg",
	"jpm":                      "video/jpm",
	"jpx":                      "image/jpx",
	"js":                       "application/javascript",
	"json":                     "application/json",
	"json5":                    "application/json5",
	"jsonld":                   "application/ld+json",
	"jsonml":                   "application/jsonml+json",
	"jsx":                      "text/jsx",
	"k25":                      "image/x-kodak-k25",
	"kar":                      "audio/midi",
	"karbon":                   "application/vnd.kde.karbon",
	"kdc":                      "image/x-kodak-kdc",
	"keynote":                  "application/vnd.apple.keynote",
	"kfo":                      "application/vnd.kde.kformula",
	"kia":                      "application/vnd.kidspiration",
	"kml":                      "application/vnd.google-earth.kml+xml",
	"kmz":                      "application/vnd.google-earth.kmz",
	"kne":                      "application/vnd.kinar",
	"knp":                      "application/vnd.kinar",
	"kon":                      "application/vnd.kde.kontour",
	"kpr":                      "application/vnd.kde.kpresenter",
	"kpt":                      "application/vnd.kde.kpresenter",
	"kpxx":                     "application/vnd.ds-keypoint",
	"ksp":                      "application/vnd.kde.kspread",
	"ktr":                      "application/vnd.kahootz",
	"ktx":                      "image/ktx",
	"ktz":                      "application/vnd.kahootz",
	"kwd":                      "application/vnd.kde.kword",
	"kwt":                      "application/vnd.kde.kword",
	"lasxml":                   "application/vnd.las.las+xml",
	"latex":                    "application/x-latex",
	"lbd":                      "application/vnd.llamagraphics.life-balance.desktop",
	"lbe":                      "application/vnd.llamagraphics.life-balance.exchange+xml",
	"les":                      "application/vnd.hhe.lesson-player",
	"less":                     "text/less",
	"lha":                      "application/x-lzh-compressed",
	"link66":                   "application/vnd.route66.link66+xml",
	"list":                     "text/plain",
	"list3820":                 "application/vnd.ibm.modcap",
	"listafp":                  "application/vnd.ibm.modcap",
	"litcoffee":                "text/coffeescript",
	"lnk":                      "application/x-ms-shortcut",
	"log":                      "text/plain",
	"lostxml":                  "application/lost+xml",
	"lrf":                      "application/octet-stream",
	"lrm":                      "application/vnd.ms-lrm",
	"ltf":                      "application/vnd.frogans.ltf",
	"lua":                      "text/x-lua",
	"luac":                     "application/x-lua-bytecode",
	"lvp":                      "audio/vnd.lucent.voice",
	"lwp":                      "application/vnd.lotus-wordpro",
	"lzh":                      "application/x-lzh-compressed",
	"m13":                      "application/x-msmediaview",
	"m14":                      "application/x-msmediaview",
	"m1v":                      "video/mpeg",
	"m21":                      "application/mp21",
	"m2a":                      "audio/mpeg",
	"m2v":                      "video/mpeg",
	"m3a":                      "audio/mpeg",
	"m3u":                      "audio/x-mpegurl",
	"m3u8":                     "application/vnd.apple.mpegurl",
	"m4a":                      "audio/x-m4a",
	"m4p":                      "application/mp4",
	"m4u":                      "video/vnd.mpegurl",
	"m4v":                      "video/x-m4v",
	"ma":                       "application/mathematica",
	"mads":                     "application/mads+xml",
	"mag":                      "application/vnd.ecowin.chart",
	"maker":                    "application/vnd.framemaker",
	"man":                      "text/troff",
	"manifest":                 "text/cache-manifest",
	"map":                      "application/json",
	"mar":                      "application/octet-stream",
	"markdown":                 "text/markdown",
	"mathml":                   "application/mathml+xml",
	"mb":                       "application/mathematica",
	"mbk":                      "application/vnd.mobius.mbk",
	"mbox":                     "application/mbox",
	"mc1":                      "application/vnd.medcalcdata",
	"mcd":                      "application/vnd.mcd",
	"mcurl":                    "text/vnd.curl.mcurl",
	"md":                       "text/markdown",
	"mdb":                      "application/x-msaccess",
	"mdi":                      "image/vnd.ms-modi",
	"me":                       "text/troff",
	"mesh":                     "model/mesh",
	"meta4":                    "application/metalink4+xml",
	"metalink":                 "application/metalink+xml",
	"mets":                     "application/mets+xml",
	"mfm":                      "application/vnd.mfmp",
	"mft":                      "application/rpki-manifest",
	"mgp":                      "application/vnd.osgeo.mapguide.package",
	"mgz":                      "application/vnd.proteus.magazine",
	"mid":                      "audio/midi",
	"midi":                     "audio/midi",
	"mie":                      "application/x-mie",
	"mif":                      "application/vnd.mif",
	"mime":                     "message/rfc822",
	"mj2":                      "video/mj2",
	"mjp2":                     "video/mj2",
	"mjs":                      "application/javascript",
	"mk3d":                     "video/x-matroska",
	"mka":                      "audio/x-matroska",
	"mkd":                      "text/x-markdown",
	"mks":                      "video/x-matroska",
	"mkv":                      "video/x-matroska",
	"mlp":                      "application/vnd.dolby.mlp",
	"mmd":                      "application/vnd.chipnuts.karaoke-mmd",
	"mmf":                      "application/vnd.smaf",
	"mml":                      "text/mathml",
	"mmr":                      "image/vnd.fujixerox.edmics-mmr",
	"mng":                      "video/x-mng",
	"mny":                      "application/x-msmoney",
	"mobi":                     "application/x-mobipocket-ebook",
	"mods":                     "application/mods+xml",
	"mov":                      "video/quicktime",
	"movie":                    "video/x-sgi-movie",
	"mp2":                      "audio/mpeg",
	"mp21":                     "application/mp21",
	"mp2a":                     "audio/mpeg",
	"mp3":                      "audio/mpeg",
	"mp4":                      "video/mp4",
	"mp4a":                     "audio/mp4",
	"mp4s":                     "application/mp4",
	"mp4v":                     "video/mp4",
	"mpc":                      "application/vnd.mophun.certificate",
	"mpd":                      "application/dash+xml",
	"mpe":                      "video/mpeg",
	"mpeg":                     "video/mpeg",
	"mpg":                      "video/mpeg",
	"mpg4":                     "video/mp4",
	"mpga":                     "audio/mpeg",
	"mpkg":                     "application/vnd.apple.installer+xml",
	"mpm":                      "application/vnd.blueice.multipass",
	"mpn":                      "application/vnd.mophun.application",
	"mpp":                      "application/vnd.ms-project",
	"mpt":                      "application/vnd.ms-project",
	"mpy":                      "application/vnd.ibm.minipay",
	"mqy":                      "application/vnd.mobius.mqy",
	"mrc":                      "application/marc",
	"mrcx":                     "application/marcxml+xml",
	"mrw":                      "image/x-minolta-mrw",
	"ms":                       "text/troff",
	"mscml":                    "application/mediaservercontrol+xml",
	"mseed":                    "application/vnd.fdsn.mseed",
	"mseq":                     "application/vnd.mseq",
	"msf":                      "application/vnd.epson.msf",
	"msg":                      "application/vnd.ms-outlook",
	"msh":                      "model/mesh",
	"msi":                      "application/x-msdownload",
	"msl":                      "application/vnd.mobius.msl",
	"msm":                      "application/octet-stream",
	"msp":                      "application/octet-stream",
	"msty":                     "application/vnd.muvee.style",
	"mts":                      "model/vnd.mts",
	"mus":                      "application/vnd.musician",
	"musicxml":                 "application/vnd.recordare.musicxml+xml",
	"mvb":                      "application/x-msmediaview",
	"mwf":                      "application/vnd.mfer",
	"mxf":                      "application/mxf",
	"mxl":                      "application/vnd.recordare.musicxml",
	"mxml":                     "application/xv+xml",
	"mxs":                      "application/vnd.triscape.mxs",
	"mxu":                      "video/vnd.mpegurl",
	"n-gage":                   "application/vnd.nokia.n-gage.symbian.install",
	"n3":                       "text/n3",
	"nb":                       "application/mathematica",
	"nbp":                      "application/vnd.wolfram.player",
	"nc":                       "application/x-netcdf",
	"ncx":                      "application/x-dtbncx+xml",
	"nef":                      "image/x-nikon-nef",
	"nfo":                      "text/x-nfo",
	"ngdat":                    "application/vnd.nokia.n-gage.data",
	"nitf":                     "application/vnd.nitf",
	"nlu":                      "application/vnd.neurolanguage.nlu",
	"nml":                      "application/vnd.enliven",
	"nnd":                      "application/vnd.noblenet-directory",
	"nns":                      "application/vnd.noblenet-sealer",
	"nnw":                      "application/vnd.noblenet-web",
	"npx":                      "image/vnd.net-fpx",
	"nsc":                      "application/x-conference",
	"nsf":                      "application/vnd.lotus-notes",
	"ntf":                      "application/vnd.nitf",
	"numbers":                  "application/vnd.apple.numbers",
	"nzb":                      "application/x-nzb",
	"oa2":                      "application/vnd.fujitsu.oasys2",
	"oa3":                      "application/vnd.fujitsu.oasys3",
	"oas":                      "application/vnd.fujitsu.oasys",
	"obd":                      "application/x-msbinder",
	"obj":                      "application/x-tgif",
	"oda":                      "application/oda",
	"odb":                      "application/vnd.oasis.opendocument.database",
	"odc":                      "application/vnd.oasis.opendocument.chart",
	"odf":                      "application/vnd.oasis.opendocument.formula",
	"odft":                     "application/vnd.oasis.opendocument.formula-template",
	"odg":                      "application/vnd.oasis.opendocument.graphics",
	"odi":                      "application/vnd.oasis.opendocument.image",
	"odm":                      "application/vnd.oasis.opendocument.text-master",
	"odp":                      "application/vnd.oasis.opendocument.presentation",
	"ods":                      "application/vnd.oasis.opendocument.spreadsheet",
	"odt":                      "application/vnd.oasis.opendocument.text",
	"oga":                      "audio/ogg",
	"ogg":                      "audio/ogg",
	"ogv":                      "video/ogg",
	"ogx":                      "application/ogg",
	"omdoc":                    "application/omdoc+xml",
	"onepkg":                   "application/onenote",
	"onetmp":                   "application/onenote",
	"onetoc":                   "application/onenote",
	"onetoc2":                  "application/onenote",
	"opf":                      "application/oebps-package+xml",
	"opml":                     "text/x-opml",
	"oprc":                     "application/vnd.palm",
	"orf":                      "image/x-olympus-orf",
	"org":                      "text/x-org",
	"osf":                      "application/vnd.yamaha.openscoreformat",
	"osfpvg":                   "application/vnd.yamaha.openscoreformat.osfpvg+xml",
	"otc":                      "application/vnd.oasis.opendocument.chart-template",
	"otf":                      "font/otf",
	"otg":                      "application/vnd.oasis.opendocument.graphics-template",
	"oth":                      "application/vnd.oasis.opendocument.text-web",
	"oti":                      "application/vnd.oasis.opendocument.image-template",
	"otp":                      "application/vnd.oasis.opendocument.presentation-template",
	"ots":                      "application/vnd.oasis.opendocument.spreadsheet-template",
	"ott":                      "application/vnd.oasis.opendocument.text-template",
	"ova":                      "application/x-virtualbox-ova",
	"ovf":                      "application/x-virtualbox-ovf",
	"owl":                      "application/rdf+xml",
	"oxps":                     "application/oxps",
	"oxt":                      "application/vnd.openofficeorg.extension",
	"p":                        "text/x-pascal",
	"p10":                      "application/pkcs10",
	"p12":                      "application/x-pkcs12",
	"p7b":                      "application/x-pkcs7-certificates",
	"p7c":                      "application/pkcs7-mime",
	"p7m":                      "application/pkcs7-mime",
	"p7r":                      "application/x-pkcs7-certreqresp",
	"p7s":                      "application/pkcs7-signature",
	"p8":                       "application/pkcs8",
	"pac":                      "application/x-ns-proxy-autoconfig",
	"pages":                    "application/vnd.apple.pages",
	"pas":                      "text/x-pascal",
	"paw":                      "application/vnd.pawaafile",
	"pbd":                      "application/vnd.powerbuilder6",
	"pbm":                      "image/x-portable-bitmap",
	"pcap":                     "application/vnd.tcpdump.pcap",
	"pcf":                      "application/x-font-pcf",
	"pcl":                      "application/vnd.hp-pcl",
	"pclxl":                    "application/vnd.hp-pclxl",
	"pct":                      "image/x-pict",
	"pcurl":                    "application/vnd.curl.pcurl",
	"pcx":                      "image/x-pcx",
	"pdb":                      "application/x-pilot",
	"pde":                      "text/x-processing",
	"pdf":                      "application/pdf",
	"pef":                      "image/x-pentax-pef",
	"pem":                      "application/x-x509-ca-cert",
	"pfa":                      "application/x-font-type1",
	"pfb":                      "application/x-font-type1",
	"pfm":                      "application/x-font-type1",
	"pfr":                      "application/font-tdpfr",
	"pfx":                      "application/x-pkcs12",
	"pgm":                      "image/x-portable-graymap",
	"pgn":                      "application/x-chess-pgn",
	"pgp":                      "application/pgp-encrypted",
	"php":                      "application/x-httpd-php",
	"pic":                      "image/x-pict",
	"pkg":                      "application/octet-stream",
	"pki":                      "application/pkixcmp",
	"pkipath":                  "application/pkix-pkipath",
	"pkpass":                   "application/vnd.apple.pkpass",
	"pl":                       "application/x-perl",
	"plb":                      "application/vnd.3gpp.pic-bw-large",
	"plc":                      "application/vnd.mobius.plc",
	"plf":                      "application/vnd.pocketlearn",
	"pls":                      "application/pls+xml",
	"pm":                       "application/x-perl",
	"pml":                      "application/vnd.ctc-posml",
	"png":                      "image/png",
	"pnm":                      "image/x-portable-anymap",
	"portpkg":                  "application/vnd.macports.portpkg",
	"pot":                      "application/vnd.ms-powerpoint",
	"potm":                     "application/vnd.ms-powerpoint.template.macroenabled.12",
	"potx":                     "application/vnd.openxmlformats-officedocument.presentationml.template",
	"ppam":                     "application/vnd.ms-powerpoint.addin.macroenabled.12",
	"ppd":                      "application/vnd.cups-ppd",
	"ppm":                      "image/x-portable-pixmap",
	"pps":                      "application/vnd.ms-powerpoint",
	"ppsm":                     "application/vnd.ms-powerpoint.slideshow.macroenabled.12",
	"ppsx":                     "application/vnd.openxmlformats-officedocument.presentationml.slideshow",
	"ppt":                      "application/vnd.ms-powerpoint",
	"pptm":                     "application/vnd.ms-powerpoint.presentation.macroenabled.12",
	"pptx":                     "application/vnd.openxmlformats-officedocument.presentationml.presentation",
	"pqa":                      "application/vnd.palm",
	"prc":                      "application/x-pilot",
	"pre":                      "application/vnd.lotus-freelance",
	"prf":                      "application/pics-rules",
	"ps":                       "application/postscript",
	"psb":                      "application/vnd.3gpp.pic-bw-small",
	"psd":                      "image/vnd.adobe.photoshop",
	"psf":                      "application/x-font-linux-psf",
	"pskcxml":                  "application/pskc+xml",
	"pti":                      "image/prs.pti",
	"ptid":                     "application/vnd.pvi.ptid1",
	"pub":                      "application/x-mspublisher",
	"pvb":                      "application/vnd.3gpp.pic-bw-var",
	"pwn":                      "application/vnd.3m.post-it-notes",
	"pya":                      "audio/vnd.ms-playready.media.pya",
	"pyv":                      "video/vnd.ms-playready.media.pyv",
	"qam":                      "application/vnd.epson.quickanime",
	"qbo":                      "application/vnd.intu.qbo",
	"qfx":                      "application/vnd.intu.qfx",
	"qps":                      "application/vnd.publishare-delta-tree",
	"qt":                       "video/quicktime",
	"qwd":                      "application/vnd.quark.quarkxpress",
	"qwt":                      "application/vnd.quark.quarkxpress",
	"qxb":                      "application/vnd.quark.quarkxpress",
	"qxd":                      "application/vnd.quark.quarkxpress",
	"qxl":                      "application/vnd.quark.quarkxpress",
	"qxt":                      "application/vnd.quark.quarkxpress",
	"ra":                       "audio/x-realaudio",
	"raf":                      "image/x-fuji-raf",
	"ram":                      "audio/x-pn-realaudio",
	"raml":                     "application/raml+yaml",
	"rar":                      "application/x-rar-compressed",
	"ras":                      "image/x-cmu-raster",
	"raw":                      "image/x-panasonic-raw",
	"rcprofile":                "application/vnd.ipunplugged.rcprofile",
	"rdf":                      "application/rdf+xml",
	"rdz":                      "application/vnd.data-vision.rdz",
	"rep":                      "application/vnd.businessobjects",
	"res":                      "application/x-dtbresource+xml",
	"rgb":                      "image/x-rgb",
	"rif":                      "application/reginfo+xml",
	"rip":                      "audio/vnd.rip",
	"ris":                      "application/x-research-info-systems",
	"rl":                       "application/resource-lists+xml",
	"rlc":                      "image/vnd.fujixerox.edmics-rlc",
	"rld":                      "application/resource-lists-diff+xml",
	"rm":                       "application/vnd.rn-realmedia",
	"rmi":                      "audio/midi",
	"rmp":                      "audio/x-pn-realaudio-plugin",
	"rms":                      "application/vnd.jcp.javame.midlet-rms",
	"rmvb":                     "application/vnd.rn-realmedia-vbr",
	"rnc":                      "application/relax-ng-compact-syntax",
	"rng":                      "application/xml",
	"roa":                      "application/rpki-roa",
	"roff":                     "text/troff",
	"rp9":                      "application/vnd.cloanto.rp9",
	"rpm":                      "application/x-redhat-package-manager",
	"rpss":                     "application/vnd.nokia.radio-presets",
	"rpst":                     "application/vnd.nokia.radio-preset",
	"rq":                       "application/sparql-query",
	"rs":                       "application/rls-services+xml",
	"rsd":                      "application/rsd+xml",
	"rss":                      "application/rss+xml",
	"rtf":                      "text/rtf",
	"rtx":                      "text/richtext",
	"run":                      "application/x-makeself",
	"s":                        "text/x-asm",
	"s3m":                      "audio/s3m",
	"saf":                      "application/vnd.yamaha.smaf-audio",
	"sass":                     "text/x-sass",
	"sbml":                     "application/sbml+xml",
	"sc":                       "application/vnd.ibm.secure-container",
	"scd":                      "application/x-msschedule",
	"scm":                      "application/vnd.lotus-screencam",
	"scq":                      "application/scvp-cv-request",
	"scs":                      "application/scvp-cv-response",
	"scss":                     "text/x-scss",
	"scurl":                    "text/vnd.curl.scurl",
	"sda":                      "application/vnd.stardivision.draw",
	"sdc":                      "application/vnd.stardivision.calc",
	"sdd":                      "application/vnd.stardivision.impress",
	"sdkd":                     "application/vnd.solent.sdkm+xml",
	"sdkm":                     "application/vnd.solent.sdkm+xml",
	"sdp":                      "application/sdp",
	"sdw":                      "application/vnd.stardivision.writer",
	"sea":                      "application/x-sea",
	"see":                      "application/vnd.seemail",
	"seed":                     "application/vnd.fdsn.seed",
	"sema":                     "application/vnd.sema",
	"semd":                     "application/vnd.semd",
	"semf":                     "application/vnd.semf",
	"ser":                      "application/java-serialized-object",
	"setpay":                   "application/set-payment-initiation",
	"setreg":                   "application/set-registration-initiation",
	"sfd-hdstx":                "application/vnd.hydrostatix.sof-data",
	"sfs":                      "application/vnd.spotfire.sfs",
	"sfv":                      "text/x-sfv",
	"sgi":                      "image/sgi",
	"sgl":                      "application/vnd.stardivision.writer-global",
	"sgm":                      "text/sgml",
	"sgml":                     "text/sgml",
	"sh":                       "application/x-sh",
	"shar":                     "application/x-shar",
	"shex":                     "text/shex",
	"shf":                      "application/shf+xml",
	"shtml":                    "text/html",
	"sid":                      "image/x-mrsid-image",
	"sig":                      "application/pgp-signature",
	"sil":                      "audio/silk",
	"silo":                     "model/mesh",
	"sis":                      "application/vnd.symbian.install",
	"sisx":                     "application/vnd.symbian.install",
	"sit":                      "application/x-stuffit",
	"sitx":                     "application/x-stuffitx",
	"skd":                      "application/vnd.koan",
	"skm":                      "application/vnd.koan",
	"skp":                      "application/vnd.koan",
	"skt":                      "application/vnd.koan",
	"sldm":                     "application/vnd.ms-powerpoint.slide.macroenabled.12",
	"sldx":                     "application/vnd.openxmlformats-officedocument.presentationml.slide",
	"slim":                     "text/slim",
	"slm":                      "text/slim",
	"slt":                      "application/vnd.epson.salt",
	"sm":                       "application/vnd.stepmania.stepchart",
	"smf":                      "application/vnd.stardivision.math",
	"smi":                      "application/smil+xml",
	"smil":                     "application/smil+xml",
	"smv":                      "video/x-smv",
	"smzip":                    "application/vnd.stepmania.package",
	"snd":                      "audio/basic",
	"snf":                      "application/x-font-snf",
	"so":                       "application/octet-stream",
	"spc":                      "application/x-pkcs7-certificates",
	"spf":                      "application/vnd.yamaha.smaf-phrase",
	"spl":                      "application/x-futuresplash",
	"spot":                     "text/vnd.in3d.spot",
	"spp":                      "application/scvp-vp-response",
	"spq":                      "application/scvp-vp-request",
	"spx":                      "audio/ogg",
	"sql":                      "application/x-sql",
	"sr2":                      "image/x-sony-sr2",
	"src":                      "application/x-wais-source",
	"srf":                      "image/x-sony-srf",
	"srt":                      "application/x-subrip",
	"sru":                      "application/sru+xml",
	"srx":                      "application/sparql-results+xml",
	"ssdl":                     "application/ssdl+xml",
	"sse":                      "application/vnd.kodak-descriptor",
	"ssf":                      "application/vnd.epson.ssf",
	"ssml":                     "application/ssml+xml",
	"st":                       "application/vnd.sailingtracker.track",
	"stc":                      "application/vnd.sun.xml.calc.template",
	"std":                      "application/vnd.sun.xml.draw.template",
	"stf":                      "application/vnd.wt.stf",
	"sti":                      "application/vnd.sun.xml.impress.template",
	"stk":                      "application/hyperstudio",
	"stl":                      "application/vnd.ms-pki.stl",
	"str":                      "application/vnd.pg.format",
	"stw":                      "application/vnd.sun.xml.writer.template",
	"styl":                     "text/stylus",
	"stylus":                   "text/stylus",
	"sub":                      "text/vnd.dvb.subtitle",
	"sus":                      "application/vnd.sus-calendar",
	"susp":                     "application/vnd.sus-calendar",
	"sv4cpio":                  "application/x-sv4cpio",
	"sv4crc":                   "application/x-sv4crc",
	"svc":                      "application/vnd.dvb.service",
	"svd":                      "application/vnd.svd",
	"svg":                      "image/svg+xml",
	"svgz":                     "image/svg+xml",
	"swa":                      "application/x-director",
	"swf":                      "application/x-shockwave-flash",
	"swi":                      "application/vnd.aristanetworks.swi",
	"sxc":                      "application/vnd.sun.xml.calc",
	"sxd":                      "application/vnd.sun.xml.draw",
	"sxg":                      "application/vnd.sun.xml.writer.global",
	"sxi":                      "application/vnd.sun.xml.impress",
	"sxm":                      "application/vnd.sun.xml.math",
	"sxw":                      "application/vnd.sun.xml.writer",
	"t":                        "text/troff",
	"t3":                       "application/x-t3vm-image",
	"t38":                      "image/t38",
	"taglet":                   "application/vnd.mynfc",
	"tao":                      "application/vnd.tao.intent-module-archive",
	"tap":                      "image/vnd.tencent.tap",
	"tar":                      "application/x-tar",
	"tcap":                     "application/vnd.3gpp2.tcap",
	"tcl":                      "application/x-tcl",
	"teacher":                  "application/vnd.smart.teacher",
	"tei":                      "application/tei+xml",
	"teicorpus":                "application/tei+xml",
	"tex":                      "application/x-tex",
	"texi":                     "application/x-texinfo",
	"texinfo":                  "application/x-texinfo",
	"text":                     "text/plain",
	"tfi":                      "application/thraud+xml",
	"tfm":                      "application/x-tex-tfm",
	"tfx":                      "image/tiff-fx",
	"tga":                      "image/x-tga",
	"thmx":                     "application/vnd.ms-officetheme",
	"tif":                      "image/tiff",
	"tiff":                     "image/tiff",
	"tk":                       "application/x-tcl",
	"tmo":                      "application/vnd.tmobile-livetv",
	"torrent":                  "application/x-bittorrent",
	"tpl":                      "application/vnd.groove-tool-template",
	"tpt":                      "application/vnd.trid.tpt",
	"tr":                       "text/troff",
	"tra":                      "application/vnd.trueapp",
	"trm":                      "application/x-msterminal",
	"ts":                       "video/mp2t",
	"tsd":                      "application/timestamped-data",
	"tsv":                      "text/tab-separated-values",
	"ttc":                      "font/collection",
	"ttf":                      "font/ttf",
	"ttl":                      "text/turtle",
	"twd":                      "application/vnd.simtech-mindmapper",
	"twds":                     "application/vnd.simtech-mindmapper",
	"txd":                      "application/vnd.genomatix.tuxedo",
	"txf":                      "application/vnd.mobius.txf",
	"txt":                      "text/plain",
	"u32":                      "application/x-authorware-bin",
	"u8dsn":                    "message/global-delivery-status",
	"u8hdr":                    "message/global-headers",
	"u8mdn":                    "message/global-disposition-notification",
	"u8msg":                    "message/global",
	"udeb":                     "application/x-debian-package",
	"ufd":                      "application/vnd.ufdl",
	"ufdl":                     "application/vnd.ufdl",
	"ulx":                      "application/x-glulx",
	"umj":                      "application/vnd.umajin",
	"unityweb":                 "application/vnd.unity",
	"uoml":                     "application/vnd.uoml+xml",
	"uri":                      "text/uri-list",
	"uris":                     "text/uri-list",
	"urls":                     "text/uri-list",
	"ustar":                    "application/x-ustar",
	"utz":                      "application/vnd.uiq.theme",
	"uu":                       "text/x-uuencode",
	"uva":                      "audio/vnd.dece.audio",
	"uvd":                      "application/vnd.dece.data",
	"uvf":                      "application/vnd.dece.data",
	"uvg":                      "image/vnd.dece.graphic",
	"uvh":                      "video/vnd.dece.hd",
	"uvi":                      "image/vnd.dece.graphic",
	"uvm":                      "video/vnd.dece.mobile",
	"uvp":                      "video/vnd.dece.pd",
	"uvs":                      "video/vnd.dece.sd",
	"uvt":                      "application/vnd.dece.ttml+xml",
	"uvu":                      "video/vnd.uvvu.mp4",
	"uvv":                      "video/vnd.dece.video",
	"uvva":                     "audio/vnd.dece.audio",
	"uvvd":                     "application/vnd.dece.data",
	"uvvf":                     "application/vnd.dece.data",
	"uvvg":                     "image/vnd.dece.graphic",
	"uvvh":                     "video/vnd.dece.hd",
	"uvvi":                     "image/vnd.dece.graphic",
	"uvvm":                     "video/vnd.dece.mobile",
	"uvvp":                     "video/vnd.dece.pd",
	"uvvs":                     "video/vnd.dece.sd",
	"uvvt":                     "application/vnd.dece.ttml+xml",
	"uvvu":                     "video/vnd.uvvu.mp4",
	"uvvv":                     "video/vnd.dece.video",
	"uvvx":                     "application/vnd.dece.unspecified",
	"uvvz":                     "application/vnd.dece.zip",
	"uvx":                      "application/vnd.dece.unspecified",
	"uvz":                      "application/vnd.dece.zip",
	"vbox":                     "application/x-virtualbox-vbox",
	"vbox-extpack":             "application/x-virtualbox-vbox-extpack",
	"vcard":                    "text/vcard",
	"vcd":                      "application/x-cdlink",
	"vcf":                      "text/x-vcard",
	"vcg":                      "application/vnd.groove-vcard",
	"vcs":                      "text/x-vcalendar",
	"vcx":                      "application/vnd.vcx",
	"vdi":                      "application/x-virtualbox-vdi",
	"vhd":                      "application/x-virtualbox-vhd",
	"vis":                      "application/vnd.visionary",
	"viv":                      "video/vnd.vivo",
	"vmdk":                     "application/x-virtualbox-vmdk",
	"vob":                      "video/x-ms-vob",
	"vor":                      "application/vnd.stardivision.writer",
	"vox":                      "application/x-authorware-bin",
	"vrml":                     "model/vrml",
	"vsd":                      "application/vnd.visio",
	"vsf":                      "application/vnd.vsf",
	"vss":                      "application/vnd.visio",
	"vst":                      "application/vnd.visio",
	"vsw":                      "application/vnd.visio",
	"vtf":                      "image/vnd.valve.source.texture",
	"vtt":                      "text/vtt",
	"vtu":                      "model/vnd.vtu",
	"vxml":                     "application/voicexml+xml",
	"w3d":                      "application/x-director",
	"wad":                      "application/x-doom",
	"wadl":                     "application/vnd.sun.wadl+xml",
	"war":                      "application/java-archive",
	"wasm":                     "application/wasm",
	"wav":                      "audio/x-wav",
	"wax":                      "audio/x-ms-wax",
	"wbmp":                     "image/vnd.wap.wbmp",
	"wbs":                      "application/vnd.criticaltools.wbs+xml",
	"wbxml":                    "application/vnd.wap.wbxml",
	"wcm":                      "application/vnd.ms-works",
	"wdb":                      "application/vnd.ms-works",
	"wdp":                      "image/vnd.ms-photo",
	"weba":                     "audio/webm",
	"webapp":                   "application/x-web-app-manifest+json",
	"webm":                     "video/webm",
	"webmanifest":              "application/manifest+json",
	"webp":                     "image/webp",
	"wg":                       "application/vnd.pmi.widget",
	"wgt":                      "application/widget",
	"wks":                      "application/vnd.ms-works",
	"wm":                       "video/x-ms-wm",
	"wma":                      "audio/x-ms-wma",
	"wmd":                      "application/x-ms-wmd",
	"wmf":                      "image/wmf",
	"wml":                      "text/vnd.wap.wml",
	"wmlc":                     "application/vnd.wap.wmlc",
	"wmls":                     "text/vnd.wap.wmlscript",
	"wmlsc":                    "application/vnd.wap.wmlscriptc",
	"wmv":                      "video/x-ms-wmv",
	"wmx":                      "video/x-ms-wmx",
	"wmz":                      "application/x-msmetafile",
	"woff":                     "font/woff",
	"woff2":                    "font/woff2",
	"wpd":                      "application/vnd.wordperfect",
	"wpl":                      "application/vnd.ms-wpl",
	"wps":                      "application/vnd.ms-works",
	"wqd":                      "application/vnd.wqd",
	"wri":                      "application/x-mswrite",
	"wrl":                      "model/vrml",
	"wsc":                      "message/vnd.wfa.wsc",
	"wsdl":                     "application/wsdl+xml",
	"wspolicy":                 "application/wspolicy+xml",
	"wtb":                      "application/vnd.webturbo",
	"wvx":                      "video/x-ms-wvx",
	"x32":                      "application/x-authorware-bin",
	"x3d":                      "model/x3d+xml",
	"x3db":                     "model/x3d+binary",
	"x3dbz":                    "model/x3d+binary",
	"x3dv":                     "model/x3d+vrml",
	"x3dvz":                    "model/x3d+vrml",
	"x3dz":                     "model/x3d+xml",
	"x3f":                      "image/x-sigma-x3f",
	"xaml":                     "application/xaml+xml",
	"xap":                      "application/x-silverlight-app",
	"xar":                      "application/vnd.xara",
	"xbap":                     "application/x-ms-xbap",
	"xbd":                      "application/vnd.fujixerox.docuworks.binder",
	"xbm":                      "image/x-xbitmap",
	"xdf":                      "application/xcap-diff+xml",
	"xdm":                      "application/vnd.syncml.dm+xml",
	"xdp":                      "application/vnd.adobe.xdp+xml",
	"xdssc":                    "application/dssc+xml",
	"xdw":                      "application/vnd.fujixerox.docuworks",
	"xenc":                     "application/xenc+xml",
	"xer":                      "application/patch-ops-error+xml",
	"xfdf":                     "application/vnd.adobe.xfdf",
	"xfdl":                     "application/vnd.xfdl",
	"xht":                      "application/xhtml+xml",
	"xhtml":                    "application/xhtml+xml",
	"xhvml":                    "application/xv+xml",
	"xif":                      "image/vnd.xiff",
	"xla":                      "application/vnd.ms-excel",
	"xlam":                     "application/vnd.ms-excel.addin.macroenabled.12",
	"xlc":                      "application/vnd.ms-excel",
	"xlf":                      "application/x-xliff+xml",
	"xlm":                      "application/vnd.ms-excel",
	"xls":                      "application/vnd.ms-excel",
	"xlsb":                     "application/vnd.ms-excel.sheet.binary.macroenabled.12",
	"xlsm":                     "application/vnd.ms-excel.sheet.macroenabled.12",
	"xlsx":                     "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	"xlt":                      "application/vnd.ms-excel",
	"xltm":                     "application/vnd.ms-excel.template.macroenabled.12",
	"xltx":                     "application/vnd.openxmlformats-officedocument.spreadsheetml.template",
	"xlw":                      "application/vnd.ms-excel",
	"xm":                       "audio/xm",
	"xml":                      "text/xml",
	"xo":                       "application/vnd.olpc-sugar",
	"xop":                      "application/xop+xml",
	"xpi":                      "application/x-xpinstall",
	"xpl":                      "application/xproc+xml",
	"xpm":                      "image/x-xpixmap",
	"xpr":                      "application/vnd.is-xpr",
	"xps":                      "application/vnd.ms-xpsdocument",
	"xpw":                      "application/vnd.intercon.formnet",
	"xpx":                      "application/vnd.intercon.formnet",
	"xsd":                      "application/xml",
	"xsl":                      "application/xml",
	"xslt":                     "application/xslt+xml",
	"xsm":                      "application/vnd.syncml+xml",
	"xspf":                     "application/xspf+xml",
	"xul":                      "application/vnd.mozilla.xul+xml",
	"xvm":                      "application/xv+xml",
	"xvml":                     "application/xv+xml",
	"xwd":                      "image/x-xwindowdump",
	"xyz":                      "chemical/x-xyz",
	"xz":                       "application/x-xz",
	"yaml":                     "text/yaml",
	"yang":                     "application/yang",
	"yin":                      "application/yin+xml",
	"yml":                      "text/yaml",
	"ymp":                      "text/x-suse-ymp",
	"z1":                       "application/x-zmachine",
	"z2":                       "application/x-zmachine",
	"z3":                       "application/x-zmachine",
	"z4":                       "application/x-zmachine",
	"z5":                       "application/x-zmachine",
	"z6":                       "application/x-zmachine",
	"z7":                       "application/x-zmachine",
	"z8":                       "application/x-zmachine",
	"zaz":                      "application/vnd.zzazz.deck+xml",
	"zip":                      "application/zip",
	"zir":                      "application/vnd.zul",
	"zirz":                     "application/vnd.zul",
	"zmm":                      "application/vnd.handheld-entertainment+xml",
}

func (pA *TK) GetMimeTypeByExt(extensionA string) string {
	if strings.HasPrefix(extensionA, ".") {
		extensionA = extensionA[1:]
	}

	contentTypeT, ok := MimeTypes[strings.ToLower(extensionA)]

	if !ok {
		return "application/octet-stream"
	}

	return contentTypeT
}

var GetMimeTypeByExt = TKX.GetMimeTypeByExt

func (pA *TK) Len(vA interface{}) int {

	lenT := -1

	switch nv := vA.(type) {
	case string:
		lenT = len(nv)
		break
	case []byte:
		lenT = len(nv)
		break
	case []rune:
		lenT = len(nv)
		break
	case []int:
		lenT = len(nv)
		break
	case []float64:
		lenT = len(nv)
		break
	case []interface{}:
		lenT = len(nv)
		break
	case []string:
		lenT = len(nv)
		break
	case []map[string]string:
		lenT = len(nv)
		break
	case []map[string]interface{}:
		lenT = len(nv)
		break
	case map[string]interface{}:
		lenT = len(nv)
		break
	case map[string]string:
		lenT = len(nv)
		break
	case map[string]int:
		lenT = len(nv)
		break
	case *ByteQueue:
		lenT = nv.Size()
		break
	case *AnyQueue:
		lenT = nv.Size()
		break
	case *OrderedMap:
		lenT = nv.Len()
		break
	default:
		valueT := reflect.ValueOf(vA)

		kindT := valueT.Kind()

		if kindT == reflect.Array || kindT == reflect.Slice || kindT == reflect.Chan || kindT == reflect.Map {
			lenT = valueT.Len()
		}

	}

	return lenT
}

var Len = TKX.Len

// func Counter(stopA int) Iterator {
// 	return &StepCounterIterator{Current: 0, Step: 1, Stop: stopA}
// }

// func StepCounter(startA, stopA, stepA int) Iterator {
// 	return &StepCounterIterator{Current: startA, Step: stepA, Stop: stopA}
// }

// type StepCounterIterator struct {
// 	Current int
// 	Step    int
// 	Stop    int
// }

// func (iter *StepCounterIterator) Next() (interface{}, bool) {
// 	if iter.Current >= iter.Stop {
// 		return 0, false
// 	}

// 	item := iter.Current
// 	iter.Current += iter.Step

// 	return item, true
// }

// func (iter *StepCounterIterator) QuickNext() interface{} {
// 	if iter.Current >= iter.Stop {
// 		return nil
// 	}

// 	item := iter.Current
// 	iter.Current += iter.Step

// 	return item
// }

// func StepFloatCounter(startA, stopA, stepA float64) Iterator {
// 	return &StepFloatCounterIterator{Current: startA, Step: stepA, Stop: stopA}
// }

// type StepFloatCounterIterator struct {
// 	Current float64
// 	Step    float64
// 	Stop    float64
// }

// func (iter *StepFloatCounterIterator) Next() (interface{}, bool) {
// 	if iter.Current >= iter.Stop {
// 		return 0, false
// 	}

// 	item := iter.Current
// 	iter.Current += iter.Step

// 	return item, true
// }

// func (iter *StepFloatCounterIterator) QuickNext() interface{} {
// 	if iter.Current >= iter.Stop {
// 		return nil
// 	}

// 	item := iter.Current
// 	iter.Current += iter.Step

// 	return item
// }

// dirA: the condition direction: 0(default): >=, 1: <=, 2: >, 3: <, 4: ==
// startA == -1, indicates startA = len(dataA) - 1
// stopA == -1, indicates stopA = len(dataA)
// iter2 := tk.NewCompactIterator("abc123", 0, -1, 1, 0)
// for {
// 	item, ok := iter2.Next()
// 	if !ok {
// 		break
// 	}
// 	fmt.Println(item)
// }
// output in order: a, b, c, 1, 2, 3, each in a line
//
// iter3 := tk.NewCompactIterator("abc123", -1, 0, -1, 3)
//
// for {
// 	item, ok := iter3.Next()
// 	if !ok {
// 		break
// 	}
// 	fmt.Println(item)
// }
// output in order: 3, 2, 1, c, b, a, each in a line

// func NewCompactIterator(dataA interface{}, startA, stopA, stepA interface{}, dirA ...int) Iterator {

type Iterator interface {
	// IsValid() bool
	HasNext() bool
	// count, index/key, value, if has value
	Next() (int, interface{}, interface{}, bool)
	QuickNext() interface{}
	QuickNextWithIndex() interface{}
}

// NewCompactIterator(dataA, startA, stopA, stepA, dir(ection)A)
// for dirA: 0(default): >=, 1: <=, 2: >, 3: <, 4: ==, means the condition operator to terminate the range loop
// in most cases, for positive direction, it should be 0, for negative range, it will be 1
func (pA *TK) NewCompactIterator(dataA interface{}, argsA ...interface{}) Iterator {
	var startT interface{} = 0
	if len(argsA) > 0 {
		startT = argsA[0]
	}

	var stopT interface{} = -1
	if len(argsA) > 1 {
		stopT = argsA[1]
	}

	var stepT interface{} = 1
	if len(argsA) > 2 {
		stepT = argsA[2]
	}

	dirT := 0
	if len(argsA) > 3 {
		dirT = ToInt(argsA[3], 0)
	}

	switch nv := dataA.(type) {
	case int:
		// stopT := ToInt(stopT, 0)

		// if stopT <= -1 {
		// 	stopT = nv
		// }

		return &CompactIterator{Type: 1, Current: ToInt(startT, 0), Step: ToInt(stepT, 1), Stop: nv, Direction: dirT}
	case float64:
		// startT := ToFloat(startT, 0.0)

		// if startT <= -1 {
		// 	startT = 0.0
		// }

		// stopT := ToFloat(stopT, 0.0)

		// if stopT <= -1 {
		// 	stopT = nv
		// }

		return &CompactIterator{Type: 2, Current: ToFloat(startT, 0.0), Step: ToFloat(stepT, 1.0), Stop: nv, Direction: dirT}
	case string:
		startT := ToInt(startT, 0)

		if startT <= -1 {
			startT = len(nv) - 1
		}

		stopT := ToInt(stopT, 0)

		if stopT <= -1 {
			stopT = len(nv)
		}

		return &CompactIterator{Type: 21, Data: nv, Current: startT, Step: ToInt(stepT, 1), Stop: stopT, Direction: dirT}

	case []byte:
		startT := ToInt(startT, 0)

		if startT <= -1 {
			startT = len(nv) - 1
		}

		stopT := ToInt(stopT, 0)

		if stopT <= -1 {
			stopT = len(nv)
		}

		return &CompactIterator{Type: 22, Data: nv, Current: startT, Step: ToInt(stepT, 1), Stop: stopT, Direction: dirT}

	case []rune:
		startT := ToInt(startT, 0)

		if startT <= -1 {
			startT = len(nv) - 1
		}

		stopT := ToInt(stopT, 0)

		if stopT <= -1 {
			stopT = len(nv)
		}

		return &CompactIterator{Type: 23, Data: nv, Current: startT, Step: ToInt(stepT, 1), Stop: stopT, Direction: dirT}

	case []int:
		startT := ToInt(startT, 0)

		if startT <= -1 {
			startT = len(nv) - 1
		}

		stopT := ToInt(stopT, 0)

		if stopT <= -1 {
			stopT = len(nv)
		}

		return &CompactIterator{Type: 31, Data: nv, Current: startT, Step: ToInt(stepT, 1), Stop: stopT, Direction: dirT}

	case []float64:
		startT := ToInt(startT, 0)

		if startT <= -1 {
			startT = len(nv) - 1
		}

		stopT := ToInt(stopT, 0)

		if stopT <= -1 {
			stopT = len(nv)
		}

		return &CompactIterator{Type: 32, Data: nv, Current: startT, Step: ToInt(stepT, 1), Stop: stopT, Direction: dirT}

	case []string:
		startT := ToInt(startT, 0)

		if startT <= -1 {
			startT = len(nv) - 1
		}

		stopT := ToInt(stopT, 0)

		if stopT <= -1 {
			stopT = len(nv)
		}

		return &CompactIterator{Type: 33, Data: nv, Current: startT, Step: ToInt(stepT, 1), Stop: stopT, Direction: dirT}

	case []interface{}:
		startT := ToInt(startT, 0)

		if startT <= -1 {
			startT = len(nv) - 1
		}

		stopT := ToInt(stopT, 0)

		if stopT <= -1 {
			stopT = len(nv)
		}

		return &CompactIterator{Type: 91, Data: nv, Current: startT, Step: ToInt(stepT, 1), Stop: stopT, Direction: dirT}

	case []map[string]string:
		startT := ToInt(startT, 0)

		if startT <= -1 {
			startT = len(nv) - 1
		}

		stopT := ToInt(stopT, 0)

		if stopT <= -1 {
			stopT = len(nv)
		}

		return &CompactIterator{Type: 92, Data: nv, Current: startT, Step: ToInt(stepT, 1), Stop: stopT, Direction: dirT}

	case []map[string]interface{}:
		startT := ToInt(startT, 0)

		if startT <= -1 {
			startT = len(nv) - 1
		}

		stopT := ToInt(stopT, 0)

		if stopT <= -1 {
			stopT = len(nv)
		}

		return &CompactIterator{Type: 93, Data: nv, Current: startT, Step: ToInt(stepT, 1), Stop: stopT, Direction: dirT}

	case []*OrderedMap:
		startT := ToInt(startT, 0)

		if startT <= -1 {
			startT = len(nv) - 1
		}

		stopT := ToInt(stopT, 0)

		if stopT <= -1 {
			stopT = len(nv)
		}

		return &CompactIterator{Type: 94, Data: nv, Current: startT, Step: ToInt(stepT, 1), Stop: stopT, Direction: dirT}

	case map[string]string:
		startT := ToInt(startT, 0)

		if startT <= -1 {
			startT = len(nv) - 1
		}

		stopT := ToInt(stopT, 0)

		if stopT <= -1 {
			stopT = len(nv)
		}

		keysT := make([]string, 0, len(nv))

		for k := range nv {
			keysT = append(keysT, k)
		}

		return &CompactIterator{Type: 51, Data: nv, Current: startT, Step: ToInt(stepT, 1), Stop: stopT, Direction: dirT, Keys: keysT}

	case map[string]int:
		startT := ToInt(startT, 0)

		if startT <= -1 {
			startT = len(nv) - 1
		}

		stopT := ToInt(stopT, 0)

		if stopT <= -1 {
			stopT = len(nv)
		}

		keysT := make([]string, 0, len(nv))

		for k := range nv {
			keysT = append(keysT, k)
		}

		return &CompactIterator{Type: 52, Data: nv, Current: startT, Step: ToInt(stepT, 1), Stop: stopT, Direction: dirT, Keys: keysT}

	case map[string]float64:
		startT := ToInt(startT, 0)

		if startT <= -1 {
			startT = len(nv) - 1
		}

		stopT := ToInt(stopT, 0)

		if stopT <= -1 {
			stopT = len(nv)
		}

		keysT := make([]string, 0, len(nv))

		for k := range nv {
			keysT = append(keysT, k)
		}

		return &CompactIterator{Type: 53, Data: nv, Current: startT, Step: ToInt(stepT, 1), Stop: stopT, Direction: dirT, Keys: keysT}

	case map[string]interface{}:
		startT := ToInt(startT, 0)

		if startT <= -1 {
			startT = len(nv) - 1
		}

		stopT := ToInt(stopT, 0)

		if stopT <= -1 {
			stopT = len(nv)
		}

		keysT := make([]string, 0, len(nv))

		for k := range nv {
			keysT = append(keysT, k)
		}

		return &CompactIterator{Type: 81, Data: nv, Current: startT, Step: ToInt(stepT, 1), Stop: stopT, Direction: dirT, Keys: keysT}

	case map[int]int:
		startT := ToInt(startT, 0)

		if startT <= -1 {
			startT = len(nv) - 1
		}

		stopT := ToInt(stopT, 0)

		if stopT <= -1 {
			stopT = len(nv)
		}

		keysT := make([]int, 0, len(nv))

		for k := range nv {
			keysT = append(keysT, k)
		}

		return &CompactIterator{Type: 61, Data: nv, Current: startT, Step: ToInt(stepT, 1), Stop: stopT, Direction: dirT, Keys: keysT}

	case map[int]string:
		startT := ToInt(startT, 0)

		if startT <= -1 {
			startT = len(nv) - 1
		}

		stopT := ToInt(stopT, 0)

		if stopT <= -1 {
			stopT = len(nv)
		}

		keysT := make([]int, 0, len(nv))

		for k := range nv {
			keysT = append(keysT, k)
		}

		return &CompactIterator{Type: 62, Data: nv, Current: startT, Step: ToInt(stepT, 1), Stop: stopT, Direction: dirT, Keys: keysT}

	case *OrderedMap:
		startT := ToInt(startT, 0)

		lenT := nv.Len()

		if startT <= -1 {
			startT = lenT - 1
		}

		stopT := ToInt(stopT, 0)

		if stopT <= -1 {
			stopT = lenT
		}

		keysT := nv.GetKeys()

		return &CompactIterator{Type: 82, Data: nv, Current: startT, Step: ToInt(stepT, 1), Stop: stopT, Direction: dirT, Keys: keysT}

	}

	valueT := reflect.ValueOf(dataA)

	kindT := valueT.Kind()

	if kindT == reflect.Array || kindT == reflect.Slice {
		len1T := valueT.Len()
		startT := ToInt(startT, 0)

		if startT <= -1 {
			startT = len1T - 1
		}

		stopT := ToInt(stopT, 0)

		if stopT <= -1 {
			stopT = len1T
		}

		return &CompactIterator{Type: 97, Data: valueT, Current: startT, Step: ToInt(stepT, 1), Stop: stopT, Direction: dirT}

	} else if kindT == reflect.Map {
		len1T := valueT.Len()
		startT := ToInt(startT, 0)

		if startT <= -1 {
			startT = len1T - 1
		}

		stopT := ToInt(stopT, 0)

		if stopT <= -1 {
			stopT = len1T
		}

		return &CompactIterator{Type: 98, Data: valueT, Current: startT, Step: ToInt(stepT, 1), Stop: stopT, Direction: dirT, Keys: valueT.MapKeys()}

	}

	return nil // &CompactIterator{Type: 0, Data: dataA}
}

var NewCompactIterator = TKX.NewCompactIterator

type CompactIterator struct {
	Type      int // 0: unknown, 1: int, 2: float64, 21: string, 22: []byte, 23: []rune, 31: []int, 32: []float64, 33: []string, 51: map[string]string, 52: map[string]int, 53: map[string]float64, 61: map[int]int, 62: map[int]string, 81: map[string]interface{}, 82: *OrderedMap, 91: []interface{}, 92: []map[string]string, 93: []map[string]interface{}, 94: []*OrderedMap, 97: reflect array/slice, 98: reflect map
	Direction int // 0(default): >=, 1: <=, 2: >, 3: <, 4: ==

	Count int

	Current interface{}
	Step    interface{}
	Stop    interface{}

	Data interface{}

	Keys interface{}
	// Slice interface{}
}

// func (p *CompactIterator) IsValid() bool {
// 	return p.Type != 0
// }

func (p *CompactIterator) HasNext() bool {
	if p.Type == 2 { //float64
		cv := p.Current.(float64)
		stopT := p.Stop.(float64)

		if p.Direction < 1 {
			if cv >= stopT {
				return false
			}
		} else if p.Direction == 1 {
			if cv <= stopT {
				return false
			}
		} else if p.Direction == 2 {
			if cv > stopT {
				return false
			}
		} else if p.Direction == 3 {
			if cv < stopT {
				return false
			}
		} else if p.Direction == 4 {
			if cv == stopT {
				return false
			}
		}

		return true
	}

	cv := p.Current.(int)
	stopT := p.Stop.(int)

	if p.Direction < 1 {
		if cv >= stopT {
			return false
		}
	} else if p.Direction == 1 {
		if cv <= stopT {
			return false
		}
	} else if p.Direction == 2 {
		if cv > stopT {
			return false
		}
	} else if p.Direction == 3 {
		if cv < stopT {
			return false
		}
	} else if p.Direction == 4 {
		if cv == stopT {
			return false
		}
	}

	return true

}

func (p *CompactIterator) Next() (int, interface{}, interface{}, bool) {
	nilT := Undefined

	if !p.HasNext() {
		return p.Count, nilT, nilT, false
	}

	// Pln("p.Type:", p.Type)

	switch p.Type {
	case 1: // int
		// cv := p.Current.(int)
		// stopT := p.Stop.(int)
		// nilT := Undefined

		// if p.Direction < 1 {
		// 	if cv >= stopT {
		// 		return p.Count, nilT, false
		// 	}
		// } else if p.Direction == 1 {
		// 	if cv <= stopT {
		// 		return p.Count, nilT, false
		// 	}
		// } else if p.Direction == 2 {
		// 	if cv > stopT {
		// 		return p.Count, nilT, false
		// 	}
		// } else if p.Direction == 3 {
		// 	if cv < stopT {
		// 		return p.Count, nilT, false
		// 	}
		// } else if p.Direction == 4 {
		// 	if cv == stopT {
		// 		return p.Count, nilT, false
		// 	}
		// }

		item := p.Current
		p.Current = item.(int) + p.Step.(int)
		p.Count++

		return p.Count - 1, p.Count - 1, item, true
	case 2: // float64
		// cv := p.Current.(float64)
		// stopT := p.Stop.(float64)
		// nilT := Undefined

		// if p.Direction < 1 {
		// 	if cv >= stopT {
		// 		return nilT, false
		// 	}
		// } else if p.Direction == 1 {
		// 	if cv <= stopT {
		// 		return nilT, false
		// 	}
		// } else if p.Direction == 2 {
		// 	if cv > stopT {
		// 		return nilT, false
		// 	}
		// } else if p.Direction == 3 {
		// 	if cv < stopT {
		// 		return nilT, false
		// 	}
		// } else if p.Direction == 4 {
		// 	if cv == stopT {
		// 		return nilT, false
		// 	}
		// }

		// dataP := p.Data.(float64)

		item := p.Current
		p.Current = item.(float64) + p.Step.(float64)
		p.Count++

		return p.Count - 1, p.Count - 1, item, true

	case 21: // string
		// cv := p.Current.(int)
		// stopT := p.Stop.(int)
		// nilT := Undefined

		// if p.Direction < 1 {
		// 	if cv >= stopT {
		// 		return nilT, false
		// 	}
		// } else if p.Direction == 1 {
		// 	if cv <= stopT {
		// 		return nilT, false
		// 	}
		// } else if p.Direction == 2 {
		// 	if cv > stopT {
		// 		return nilT, false
		// 	}
		// } else if p.Direction == 3 {
		// 	if cv < stopT {
		// 		return nilT, false
		// 	}
		// } else if p.Direction == 4 {
		// 	if cv == stopT {
		// 		return nilT, false
		// 	}
		// }

		cv := p.Current.(int)
		dataP := p.Data.(string)

		item := dataP[cv]
		p.Current = cv + p.Step.(int)
		p.Count++

		return p.Count - 1, p.Count - 1, item, true

	case 22: // []byte
		cv := p.Current.(int)
		// stopT := p.Stop.(int)
		// nilT := Undefined
		dataP := p.Data.([]byte)

		// if p.Direction < 1 {
		// 	if cv >= stopT {
		// 		return nilT, false
		// 	}
		// } else if p.Direction == 1 {
		// 	if cv <= stopT {
		// 		return nilT, false
		// 	}
		// } else if p.Direction == 2 {
		// 	if cv > stopT {
		// 		return nilT, false
		// 	}
		// } else if p.Direction == 3 {
		// 	if cv < stopT {
		// 		return nilT, false
		// 	}
		// } else if p.Direction == 4 {
		// 	if cv == stopT {
		// 		return nilT, false
		// 	}
		// }

		item := dataP[cv]
		p.Current = cv + p.Step.(int)
		p.Count++

		return p.Count - 1, p.Count - 1, item, true

	case 23: // []rune
		cv := p.Current.(int)
		// stopT := p.Stop.(int)
		// nilT := Undefined
		dataP := p.Data.([]rune)

		// if p.Direction < 1 {
		// 	if cv >= stopT {
		// 		return nilT, false
		// 	}
		// } else if p.Direction == 1 {
		// 	if cv <= stopT {
		// 		return nilT, false
		// 	}
		// } else if p.Direction == 2 {
		// 	if cv > stopT {
		// 		return nilT, false
		// 	}
		// } else if p.Direction == 3 {
		// 	if cv < stopT {
		// 		return nilT, false
		// 	}
		// } else if p.Direction == 4 {
		// 	if cv == stopT {
		// 		return nilT, false
		// 	}
		// }

		item := dataP[cv]
		p.Current = cv + p.Step.(int)
		p.Count++

		return p.Count - 1, p.Count - 1, item, true

	case 31: // []int
		cv := p.Current.(int)
		// stopT := p.Stop.(int)
		// nilT := Undefined
		dataP := p.Data.([]int)

		// if p.Direction < 1 {
		// 	if cv >= stopT {
		// 		return nilT, false
		// 	}
		// } else if p.Direction == 1 {
		// 	if cv <= stopT {
		// 		return nilT, false
		// 	}
		// } else if p.Direction == 2 {
		// 	if cv > stopT {
		// 		return nilT, false
		// 	}
		// } else if p.Direction == 3 {
		// 	if cv < stopT {
		// 		return nilT, false
		// 	}
		// } else if p.Direction == 4 {
		// 	if cv == stopT {
		// 		return nilT, false
		// 	}
		// }

		item := dataP[cv]
		p.Current = cv + p.Step.(int)
		p.Count++

		return p.Count - 1, p.Count - 1, item, true

	case 32: // []float64
		cv := p.Current.(int)
		// stopT := p.Stop.(int)
		// nilT := Undefined
		dataP := p.Data.([]float64)

		// if p.Direction < 1 {
		// 	if cv >= stopT {
		// 		return nilT, false
		// 	}
		// } else if p.Direction == 1 {
		// 	if cv <= stopT {
		// 		return nilT, false
		// 	}
		// } else if p.Direction == 2 {
		// 	if cv > stopT {
		// 		return nilT, false
		// 	}
		// } else if p.Direction == 3 {
		// 	if cv < stopT {
		// 		return nilT, false
		// 	}
		// } else if p.Direction == 4 {
		// 	if cv == stopT {
		// 		return nilT, false
		// 	}
		// }

		item := dataP[cv]
		p.Current = cv + p.Step.(int)
		p.Count++

		return p.Count - 1, p.Count - 1, item, true

	case 33: // []string
		cv := p.Current.(int)
		// stopT := p.Stop.(int)
		// nilT := Undefined
		dataP := p.Data.([]string)

		// if p.Direction < 1 {
		// 	if cv >= stopT {
		// 		return nilT, false
		// 	}
		// } else if p.Direction == 1 {
		// 	if cv <= stopT {
		// 		return nilT, false
		// 	}
		// } else if p.Direction == 2 {
		// 	if cv > stopT {
		// 		return nilT, false
		// 	}
		// } else if p.Direction == 3 {
		// 	if cv < stopT {
		// 		return nilT, false
		// 	}
		// } else if p.Direction == 4 {
		// 	if cv == stopT {
		// 		return nilT, false
		// 	}
		// }

		item := dataP[cv]
		p.Current = cv + p.Step.(int)
		p.Count++

		return p.Count - 1, p.Count - 1, item, true

	case 51: // map[string]string
		cv := p.Current.(int)
		// stopT := p.Stop.(int)
		// nilT := Undefined
		dataP := p.Data.(map[string]string)

		// if p.Direction < 1 {
		// 	if cv >= stopT {
		// 		return nilT, false
		// 	}
		// } else if p.Direction == 1 {
		// 	if cv <= stopT {
		// 		return nilT, false
		// 	}
		// } else if p.Direction == 2 {
		// 	if cv > stopT {
		// 		return nilT, false
		// 	}
		// } else if p.Direction == 3 {
		// 	if cv < stopT {
		// 		return nilT, false
		// 	}
		// } else if p.Direction == 4 {
		// 	if cv == stopT {
		// 		return nilT, false
		// 	}
		// }

		keyT := p.Keys.([]string)[cv]

		item := dataP[keyT]
		p.Current = cv + p.Step.(int)
		p.Count++

		return p.Count - 1, keyT, item, true

	case 52: // map[string]int
		cv := p.Current.(int)
		// stopT := p.Stop.(int)
		// nilT := Undefined
		dataP := p.Data.(map[string]int)

		// if p.Direction < 1 {
		// 	if cv >= stopT {
		// 		return nilT, false
		// 	}
		// } else if p.Direction == 1 {
		// 	if cv <= stopT {
		// 		return nilT, false
		// 	}
		// } else if p.Direction == 2 {
		// 	if cv > stopT {
		// 		return nilT, false
		// 	}
		// } else if p.Direction == 3 {
		// 	if cv < stopT {
		// 		return nilT, false
		// 	}
		// } else if p.Direction == 4 {
		// 	if cv == stopT {
		// 		return nilT, false
		// 	}
		// }

		keyT := p.Keys.([]string)[cv]

		item := dataP[keyT]
		p.Current = cv + p.Step.(int)
		p.Count++

		return p.Count - 1, keyT, item, true

	case 53: // map[string]float64
		cv := p.Current.(int)
		dataP := p.Data.(map[string]float64)

		keyT := p.Keys.([]string)[cv]

		item := dataP[keyT]
		p.Current = cv + p.Step.(int)
		p.Count++

		return p.Count - 1, keyT, item, true

	case 61: // map[int]int
		cv := p.Current.(int)
		dataP := p.Data.(map[int]int)

		keyT := p.Keys.([]int)[cv]

		item := dataP[keyT]
		p.Current = cv + p.Step.(int)
		p.Count++

		return p.Count - 1, keyT, item, true

	case 62: // map[int]string
		cv := p.Current.(int)
		dataP := p.Data.(map[int]string)

		keyT := p.Keys.([]int)[cv]

		item := dataP[keyT]
		p.Current = cv + p.Step.(int)
		p.Count++

		return p.Count - 1, keyT, item, true
	case 81: // map[string]interface{}
		cv := p.Current.(int)
		dataP := p.Data.(map[string]interface{})

		keyT := p.Keys.([]string)[cv]

		item := dataP[keyT]
		p.Current = cv + p.Step.(int)
		p.Count++

		return p.Count - 1, keyT, item, true

	case 82: // *OrderedMap
		cv := p.Current.(int)
		dataP := p.Data.(*OrderedMap)

		keyT := p.Keys.([]interface{})[cv]

		item, _ := dataP.Get(keyT)
		p.Current = cv + p.Step.(int)
		p.Count++

		return p.Count - 1, keyT, item, true

	case 91: // []interface{}
		cv := p.Current.(int)
		dataP := p.Data.([]interface{})

		item := dataP[cv]
		p.Current = cv + p.Step.(int)
		p.Count++

		return p.Count - 1, p.Count - 1, item, true

	case 92: // []map[string]string
		cv := p.Current.(int)
		dataP := p.Data.([]map[string]string)

		item := dataP[cv]
		p.Current = cv + p.Step.(int)
		p.Count++

		return p.Count - 1, p.Count - 1, item, true

	case 93: // []map[string]interface{}
		cv := p.Current.(int)
		dataP := p.Data.([]map[string]interface{})

		item := dataP[cv]
		p.Current = cv + p.Step.(int)
		p.Count++

		return p.Count - 1, p.Count - 1, item, true

	case 94: // []*orderedMap
		cv := p.Current.(int)
		dataP := p.Data.([]*OrderedMap)

		item := dataP[cv]
		p.Current = cv + p.Step.(int)
		p.Count++

		return p.Count - 1, p.Count - 1, item, true

	case 97: // reflect array/slice
		cv := p.Current.(int)
		dataP := p.Data.(reflect.Value)

		item := dataP.Index(cv).Interface()
		p.Current = cv + p.Step.(int)
		p.Count++

		return p.Count - 1, p.Count - 1, item, true

	case 98: // reflect map
		cv := p.Current.(int)
		dataP := p.Data.(reflect.Value)

		keyT := p.Keys.([]reflect.Value)[cv]

		item := dataP.MapIndex(keyT)
		p.Current = cv + p.Step.(int)
		p.Count++

		return p.Count - 1, keyT, item, true

	}

	return p.Count, nilT, nilT, false

}

func (p *CompactIterator) QuickNext() interface{} {
	_, _, rs, b := p.Next()

	if !b {
		return Undefined
	}

	return rs
}

func (p *CompactIterator) QuickNextWithIndex() interface{} {
	count, indexKey, rs, b := p.Next()

	if !b {
		return Undefined
	}

	return []interface{}{indexKey, rs, count, b}
}

// bluetooth related

// func (pA *TK) BluetoothDiscoverDevice(funcA func(argsA ...interface{}) interface{}) error {

// 	var bluetoothAdapter = bluetooth.DefaultAdapter

// 	// errT := bluetoothAdapter.Enable()

// 	// if errT != nil {
// 	// 	return fmt.Errorf("enable Bluetooth function failed: %v", errT)
// 	// }

// 	errT := bluetoothAdapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
// 		funcA(device.Address.String(), device.RSSI, device.LocalName(), device.Bytes())
// 	})

// 	if errT != nil {
// 		return fmt.Errorf("scan failed: %v", errT)
// 	}

// 	return nil

// }

// var BluetoothDiscoverDevice = TKX.BluetoothDiscoverDevice

// orderedMap
// fmt.Println("## Iterating over pairs from oldest to newest: ##")
// for pair := om.Oldest(); pair != nil; pair = pair.Next() {
// 	fmt.Printf("%s => %s\n", pair.Key, pair.Value)
// }

// fmt.Println("## Iterating over the 2 newest pairs: ##")
// i := 0
// for pair := om.Newest(); pair != nil; pair = pair.Prev() {
// 	fmt.Printf("%s => %s\n", pair.Key, pair.Value)
// 	i++
// 	if i >= 2 {
// 		break
// 	}
// }

type OrderedMapPair struct {
	Key   interface{}
	Value interface{}

	element *list.Element
}

type OrderedMap struct {
	pairs map[interface{}]*OrderedMapPair
	list  *list.List
}

// New creates a new OrderedMap.
func (pA *TK) NewOrderedMap() *OrderedMap {
	return &OrderedMap{
		pairs: make(map[interface{}]*OrderedMapPair),
		list:  list.New(),
	}
}

var NewOrderedMap = TKX.NewOrderedMap

func (pA *TK) ToOrderedMap(vA interface{}) interface{} {
	if vA == nil {
		return fmt.Errorf("nil input")
	}

	var rs *OrderedMap

	switch nv := vA.(type) {
	case map[string]string:
		rs = NewOrderedMap()

		for k, v := range nv {
			rs.Set(k, v)
		}

		return rs
	case map[string]bool:
		rs = NewOrderedMap()

		for k, v := range nv {
			rs.Set(k, v)
		}

		return rs
	case map[string]byte:
		rs = NewOrderedMap()

		for k, v := range nv {
			rs.Set(k, v)
		}

		return rs
	case map[string]rune:
		rs = NewOrderedMap()

		for k, v := range nv {
			rs.Set(k, v)
		}

		return rs
	case map[string]int:
		rs = NewOrderedMap()

		for k, v := range nv {
			rs.Set(k, v)
		}

		return rs
	case map[string]int64:
		rs = NewOrderedMap()

		for k, v := range nv {
			rs.Set(k, v)
		}

		return rs
	case map[string]float64:
		rs = NewOrderedMap()

		for k, v := range nv {
			rs.Set(k, v)
		}

		return rs
	case map[int]string:
		rs = NewOrderedMap()

		for k, v := range nv {
			rs.Set(k, v)
		}

		return rs
	case map[int]int:
		rs = NewOrderedMap()

		for k, v := range nv {
			rs.Set(k, v)
		}

		return rs
	case map[int]float64:
		rs = NewOrderedMap()

		for k, v := range nv {
			rs.Set(k, v)
		}

		return rs
	case map[string]interface{}:
		rs = NewOrderedMap()

		for k, v := range nv {
			rs.Set(k, v)
		}

		return rs

	}

	return fmt.Errorf("unsupported type: %T(%v)", vA, vA)
}

var ToOrderedMap = TKX.ToOrderedMap

func (pA *OrderedMap) Dump() string {
	var bufT = new(strings.Builder)

	bufT.WriteString("pairs: [")

	for k, v := range pA.pairs {
		bufT.WriteString(fmt.Sprintf(" <key: %v, value: %v>", k, v))
	}

	bufT.WriteString("] list: [")

	cntT := 0

	for e := pA.list.Front(); e != nil; e = e.Next() {
		v := e.Value.(*OrderedMapPair)
		bufT.WriteString(fmt.Sprintf(" <%v key: %v, value: %v>", cntT, v.Key, v.Value))

		cntT++
	}

	bufT.WriteString("]")

	return bufT.String()
}

func (pA *OrderedMap) SortStringKeys(argsA ...string) error {
	keysT := pA.GetStringKeys()

	descT := IfSwitchExists(argsA, "-desc")

	if descT {
		sort.Sort(sort.Reverse(sort.StringSlice(keysT)))
	} else {
		sort.Strings(keysT)
	}

	pA.list.Init()
	for _, v := range keysT {
		// pair := &OrderedMapPair{
		// 	Key:   v,
		// 	Value: pA.pairs[v].Value,
		// }
		// pair.element = om.list.PushBack(pair)
		pA.pairs[v].element = pA.list.PushBack(pA.pairs[v])
	}

	return nil
}

func (pA *OrderedMap) SortStringKeysByFunc(funcA func(i, j int) bool) error {
	keysT := pA.GetStringKeys()

	sort.SliceStable(keysT, funcA)

	pA.list.Init()
	for _, v := range keysT {
		pA.pairs[v].element = pA.list.PushBack(pA.pairs[v])
	}

	return nil
}

// Get looks for the given key, and returns the value associated with it,
// or nil if not found. The boolean it returns says whether the key is present in the map.
func (om *OrderedMap) Get(key interface{}) (interface{}, bool) {
	if pair, present := om.pairs[key]; present {
		return pair.Value, present
	}
	return nil, false
}

func (om *OrderedMap) GetCompact(key interface{}) interface{} {
	if pair, present := om.pairs[key]; present {
		return pair.Value
	}

	return fmt.Errorf("key not found")
}

func (om *OrderedMap) GetQuick(key interface{}) interface{} {
	if pair, present := om.pairs[key]; present {
		return pair.Value
	}

	return Undefined
}

func (om *OrderedMap) GetString(key interface{}) string {
	if pair, present := om.pairs[key]; present {
		return ToStr(pair.Value)
	}
	return ""
}

func (om *OrderedMap) GetByIndex(idxA int) (interface{}, bool) {

	if idxA < 0 || idxA >= om.list.Len() {
		return nil, false
	}

	cntT := 0

	for e := om.list.Front(); e != nil; e = e.Next() {
		if cntT == idxA {
			return e.Value.(*OrderedMapPair).Value, true
		}

		cntT++
	}

	return nil, false
}

func (om *OrderedMap) GetKeyByIndex(idxA int) (interface{}, bool) {
	if idxA < 0 || idxA >= om.list.Len() {
		return nil, false
	}

	cntT := 0

	for e := om.list.Front(); e != nil; e = e.Next() {
		if cntT == idxA {
			return e.Value.(*OrderedMapPair).Key, true
		}

		cntT++
	}

	return nil, false
}

func (om *OrderedMap) GetPairByIndex(idxA int) *OrderedMapPair {
	if idxA < 0 || idxA >= om.list.Len() {
		return nil
	}

	cntT := 0

	for e := om.list.Front(); e != nil; e = e.Next() {
		if cntT == idxA {
			return e.Value.(*OrderedMapPair)
		}

		cntT++
	}

	return nil
}

func (om *OrderedMap) GetItemByIndex(idxA int) []interface{} {
	if idxA < 0 || idxA >= om.list.Len() {
		return nil
	}

	cntT := 0

	for e := om.list.Front(); e != nil; e = e.Next() {
		if cntT == idxA {
			return []interface{}{e.Value.(*OrderedMapPair).Key, e.Value.(*OrderedMapPair).Value}
		}

		cntT++
	}

	return nil
}

// Load is an alias for Get, mostly to present an API similar to `sync.Map`'s.
func (om *OrderedMap) Load(key interface{}) (interface{}, bool) {
	return om.Get(key)
}

// GetPair looks for the given key, and returns the pair associated with it,
// or nil if not found. The Pair struct can then be used to iterate over the ordered map
// from that point, either forward or backward.
func (om *OrderedMap) GetPair(key interface{}) *OrderedMapPair {
	return om.pairs[key]
}

// Set sets the key-value pair, and returns what `Get` would have returned
// on that key prior to the call to `Set`.
func (om *OrderedMap) Set(key interface{}, value interface{}) (interface{}, bool) {
	if pair, present := om.pairs[key]; present {
		oldValue := pair.Value
		pair.Value = value
		return oldValue, true
	}

	pair := &OrderedMapPair{
		Key:   key,
		Value: value,
	}

	pair.element = om.list.PushBack(pair)
	om.pairs[key] = pair

	return nil, false
}

func (om *OrderedMap) SetByIndex(idxA int, value interface{}) error {
	if idxA < 0 || idxA >= om.list.Len() {
		return fmt.Errorf("index out of range: %v/%v", idxA, om.list.Len())
	}

	cntT := 0

	for e := om.list.Front(); e != nil; e = e.Next() {
		if cntT == idxA {
			e.Value.(*OrderedMapPair).Value = value
			return nil
		}

		cntT++
	}

	return fmt.Errorf("not found")
}

// Store is an alias for Set, mostly to present an API similar to `sync.Map`'s.
func (om *OrderedMap) Store(key interface{}, value interface{}) (interface{}, bool) {
	return om.Set(key, value)
}

// Delete removes the key-value pair, and returns what `Get` would have returned
// on that key prior to the call to `Delete`.
func (om *OrderedMap) Delete(key interface{}) (interface{}, bool) {
	if pair, present := om.pairs[key]; present {
		om.list.Remove(pair.element)
		delete(om.pairs, key)
		return pair.Value, true
	}

	return nil, false
}

func (om *OrderedMap) DeleteQuick(key interface{}) error {
	if pair, present := om.pairs[key]; present {
		om.list.Remove(pair.element)
		delete(om.pairs, key)
		return nil
	}

	return fmt.Errorf("not found")
}

func (om *OrderedMap) DeleteByIndex(idxA int) error {

	if idxA < 0 || idxA >= om.list.Len() {
		return fmt.Errorf("index out of range: %v/%v", idxA, om.list.Len())
	}

	cntT := 0

	for e := om.list.Front(); e != nil; e = e.Next() {
		if cntT == idxA {
			delete(om.pairs, e.Value.(*OrderedMapPair).Key)
			om.list.Remove(e)
			return nil
		}

		cntT++
	}

	return fmt.Errorf("not found")
}

func (om *OrderedMap) Remove(key interface{}) (interface{}, bool) {
	return om.Delete(key)
}

// Len returns the length of the ordered map.
func (om *OrderedMap) Len() int {
	return len(om.pairs)
}

// Oldest returns a pointer to the oldest pair. It's meant to be used to iterate on the ordered map's
// pairs from the oldest to the newest, e.g.:
// for pair := orderedMap.Oldest(); pair != nil; pair = pair.Next() { fmt.Printf("%v => %v\n", pair.Key, pair.Value) }
func (om *OrderedMap) Oldest() *OrderedMapPair {
	return listElementToPair(om.list.Front())
}

// Newest returns a pointer to the newest pair. It's meant to be used to iterate on the ordered map's
// pairs from the newest to the oldest, e.g.:
// for pair := orderedMap.Oldest(); pair != nil; pair = pair.Next() { fmt.Printf("%v => %v\n", pair.Key, pair.Value) }
func (om *OrderedMap) Newest() *OrderedMapPair {
	return listElementToPair(om.list.Back())
}

// Next returns a pointer to the next pair.
func (p *OrderedMapPair) Next() *OrderedMapPair {
	return listElementToPair(p.element.Next())
}

// Previous returns a pointer to the previous pair.
func (p *OrderedMapPair) Prev() *OrderedMapPair {
	return listElementToPair(p.element.Prev())
}

func listElementToPair(element *list.Element) *OrderedMapPair {
	if element == nil {
		return nil
	}
	return element.Value.(*OrderedMapPair)
}

// KeyNotFoundError may be returned by functions in this package when they're called with keys that are not present
// in the map.
// type KeyNotFoundError struct {
// 	MissingKey interface{}
// }

// var _ error = &KeyNotFoundError{}

// func (e *KeyNotFoundError) Error() string {
// 	return fmt.Sprintf("missing key: %v", e.MissingKey)
// }

// MoveAfter moves the value associated with key to its new position after the one associated with markKey.
// Returns an error iff key or markKey are not present in the map.
func (om *OrderedMap) MoveAfter(key, markKey interface{}) error {
	elements, err := om.getElements(key, markKey)
	if err != nil {
		return err
	}
	om.list.MoveAfter(elements[0], elements[1])
	return nil
}

// MoveBefore moves the value associated with key to its new position before the one associated with markKey.
// Returns an error iff key or markKey are not present in the map.
func (om *OrderedMap) MoveBefore(key, markKey interface{}) error {
	elements, err := om.getElements(key, markKey)
	if err != nil {
		return err
	}
	om.list.MoveBefore(elements[0], elements[1])
	return nil
}

func (om *OrderedMap) getElements(keys ...interface{}) ([]*list.Element, error) {
	elements := make([]*list.Element, len(keys))
	for i, k := range keys {
		pair, present := om.pairs[k]
		if !present {
			return nil, fmt.Errorf("missing key: %v", k) // &KeyNotFoundError{k}
		}
		elements[i] = pair.element
	}
	return elements, nil
}

func (om *OrderedMap) GetKeys() []interface{} {
	rsT := make([]interface{}, len(om.pairs))

	cntT := 0
	for e := om.list.Front(); e != nil; e = e.Next() {
		rsT[cntT] = e.Value.(*OrderedMapPair).Key
		cntT++
	}

	return rsT
}

func (om *OrderedMap) GetStringKeys() []string {
	rsT := make([]string, len(om.pairs))

	cntT := 0
	for e := om.list.Front(); e != nil; e = e.Next() {
		rsT[cntT] = ToStr(e.Value.(*OrderedMapPair).Key)
		cntT++
	}

	return rsT
}

func (om *OrderedMap) GetValues() []interface{} {
	rsT := make([]interface{}, len(om.pairs))

	cntT := 0
	for e := om.list.Front(); e != nil; e = e.Next() {
		rsT[cntT] = e.Value.(*OrderedMapPair).Value
		cntT++
	}

	return rsT
}

func (om *OrderedMap) GetItems() []interface{} {
	rsT := make([]interface{}, len(om.pairs))

	cntT := 0
	for e := om.list.Front(); e != nil; e = e.Next() {
		valueT := e.Value.(*OrderedMapPair)
		rsT[cntT] = []interface{}{valueT.Key, valueT.Value}
		cntT++
	}

	return rsT
}

func (om *OrderedMap) ToMap() map[string]interface{} {
	rsT := make(map[string]interface{}, len(om.pairs))

	for e := om.list.Front(); e != nil; e = e.Next() {
		valueT := e.Value.(*OrderedMapPair)
		rsT[ToStr(valueT.Key)] = valueT.Value
	}

	return rsT
}

func (om *OrderedMap) ToMapAnyKey() map[interface{}]interface{} {
	rsT := make(map[interface{}]interface{}, len(om.pairs))

	for e := om.list.Front(); e != nil; e = e.Next() {
		valueT := e.Value.(*OrderedMapPair)
		rsT[valueT.Key] = valueT.Value
	}

	return rsT
}

// MoveToBack moves the value associated with key to the back of the ordered map.
// Returns an error iff key is not present in the map.
func (om *OrderedMap) MoveToBack(key interface{}) error {
	pair, present := om.pairs[key]
	if !present {
		return fmt.Errorf("missing key: %v", key)
	}
	om.list.MoveToBack(pair.element)
	return nil
}

// MoveToFront moves the value associated with key to the front of the ordered map.
// Returns an error iff key is not present in the map.
func (om *OrderedMap) MoveToFront(key interface{}) error {
	pair, present := om.pairs[key]
	if !present {
		return fmt.Errorf("missing key: %v", key)
	}
	om.list.MoveToFront(pair.element)
	return nil
}

func (om OrderedMap) String() string {
	return ToJSONX(&om)
}

var (
	_ json.Marshaler   = &OrderedMap{}
	_ json.Unmarshaler = &OrderedMap{}
)

// MarshalJSON implements the json.Marshaler interface.
func (om *OrderedMap) MarshalJSON() ([]byte, error) { //nolint:funlen
	if om == nil || om.list == nil {
		return []byte("null"), nil
	}

	writer := jwriter.Writer{}
	writer.RawByte('{')

	for pair, firstIteration := om.Oldest(), true; pair != nil; pair = pair.Next() {
		// Pl("pair: %v", pair)
		if firstIteration {
			firstIteration = false
		} else {
			writer.RawByte(',')
		}

		switch key := any(pair.Key).(type) {
		case string:
			writer.String(key)
		case encoding.TextMarshaler:
			writer.RawByte('"')
			writer.Raw(key.MarshalText())
			writer.RawByte('"')
		case int:
			writer.IntStr(key)
		case int8:
			writer.Int8Str(key)
		case int16:
			writer.Int16Str(key)
		case int32:
			writer.Int32Str(key)
		case int64:
			writer.Int64Str(key)
		case uint:
			writer.UintStr(key)
		case uint8:
			writer.Uint8Str(key)
		case uint16:
			writer.Uint16Str(key)
		case uint32:
			writer.Uint32Str(key)
		case uint64:
			writer.Uint64Str(key)
		default:

			// this switch takes care of wrapper types around primitive types, such as
			// type myType string
			switch keyValue := reflect.ValueOf(key); keyValue.Type().Kind() {
			case reflect.String:
				writer.String(keyValue.String())
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				writer.Int64Str(keyValue.Int())
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				writer.Uint64Str(keyValue.Uint())
			default:
				return nil, fmt.Errorf("unsupported key type: %T", key)
			}
		}

		writer.RawByte(':')
		// the error is checked at the end of the function
		writer.Raw(json.Marshal(pair.Value)) //nolint:errchkjson
	}

	writer.RawByte('}')

	return dumpWriter(&writer)
}

func dumpWriter(writer *jwriter.Writer) ([]byte, error) {
	if writer.Error != nil {
		return nil, writer.Error
	}

	var buf bytes.Buffer
	buf.Grow(writer.Size())
	if _, err := writer.DumpTo(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (om *OrderedMap) UnmarshalJSON(data []byte) error {
	if om.list == nil {
		// om.initialize(0)
		om.list = list.New()
		om.pairs = make(map[interface{}]*OrderedMapPair)
	}

	return jsonparser.ObjectEach(
		data,
		func(keyData []byte, valueData []byte, dataType jsonparser.ValueType, offset int) error {
			if dataType == jsonparser.String {
				// jsonparser removes the enclosing quotes; we need to restore them to make a valid JSON
				valueData = data[offset-len(valueData)-2 : offset]
			}

			var key any
			var value any

			switch typedKey := any(&key).(type) {
			case *string:
				s, err := decodeUTF8(keyData)
				if err != nil {
					return err
				}
				*typedKey = s
			case encoding.TextUnmarshaler:
				if err := typedKey.UnmarshalText(keyData); err != nil {
					return err
				}
			case *int, *int8, *int16, *int32, *int64, *uint, *uint8, *uint16, *uint32, *uint64:
				if err := json.Unmarshal(keyData, typedKey); err != nil {
					return err
				}
			default:
				// this switch takes care of wrapper types around primitive types, such as
				// type myType string
				switch reflect.TypeOf(key).Kind() {
				case reflect.String:
					s, err := decodeUTF8(keyData)
					if err != nil {
						return err
					}

					convertedKeyData := reflect.ValueOf(s).Convert(reflect.TypeOf(key))
					reflect.ValueOf(&key).Elem().Set(convertedKeyData)
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
					reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					if err := json.Unmarshal(keyData, &key); err != nil {
						return err
					}
				default:
					return fmt.Errorf("unsupported key type: %T", key)
				}
			}

			if err := json.Unmarshal(valueData, &value); err != nil {
				return err
			}

			om.Set(key, value)
			return nil
		})
}

func decodeUTF8(input []byte) (string, error) {
	remaining, offset := input, 0
	runes := make([]rune, 0, len(remaining))

	for len(remaining) > 0 {
		r, size := utf8.DecodeRune(remaining)
		if r == utf8.RuneError && size <= 1 {
			return "", fmt.Errorf("not a valid UTF-8 string (at position %d): %s", offset, string(input))
		}

		runes = append(runes, r)
		remaining = remaining[size:]
		offset += size
	}

	return string(runes), nil
}

func (pA *TK) WeixinPaySignString(valuesA interface{}, keyA string, argsA ...string) string {
	if valuesA == nil {
		return ""
	}

	var bufT strings.Builder

	var keysT []string

	// valuesT := make([]string, 0, lenT)

	methodT := GetSwitch(argsA, "-method=", "")

	var hmacT string

	nv1, ok := valuesA.(url.Values)

	if ok {
		lenT := len(nv1)

		keysT = make([]string, 0, lenT)

		for k := range nv1 {
			keysT = append(keysT, k)
			// valuesT = append(valuesT, nv1[k])
		}

		sort.Strings(keysT)

		for i, k := range keysT {
			vi := nv1[k]

			if len(vi) < 1 {
				continue
			}

			if i > 0 {
				bufT.WriteString("&")
			}

			bufT.WriteString(k)
			bufT.WriteString("=")
			bufT.WriteString(vi[0])
		}

		strT := bufT.String() + "&key=" + keyA

		if methodT == "sha256" {
			h := sha256.New()
			h.Write([]byte(strT))
			d := h.Sum(nil)
			hmacT = string(d)
		} else {
			hmacT = strings.ToUpper(MD5Encrypt(strT))
		}

		return hmacT
	}

	nv2, ok := valuesA.(map[string]string)

	if ok {
		lenT := len(nv1)

		keysT = make([]string, 0, lenT)

		for k := range nv2 {
			keysT = append(keysT, k)
			// valuesT = append(valuesT, nv1[k])
		}

		sort.Strings(keysT)

		for i, k := range keysT {
			vi := nv2[k]

			if len(vi) < 1 {
				continue
			}

			if i > 0 {
				bufT.WriteString("&")
			}

			bufT.WriteString(k)
			bufT.WriteString("=")
			bufT.WriteString(vi)
		}

		strT := bufT.String() + "&key=" + keyA

		if methodT == "sha256" {
			h := sha256.New()
			h.Write([]byte(strT))
			d := h.Sum(nil)
			hmacT = string(d)
		} else {
			hmacT = strings.ToUpper(MD5Encrypt(strT))
		}

		return hmacT

	}

	nv3, ok := valuesA.(map[string]interface{})

	if ok {
		lenT := len(nv1)

		keysT = make([]string, 0, lenT)

		for k := range nv3 {
			keysT = append(keysT, k)
			// valuesT = append(valuesT, nv1[k])
		}

		sort.Strings(keysT)

		for i, k := range keysT {
			vi := ToStr(nv3[k])

			if len(vi) < 1 {
				continue
			}

			if i > 0 {
				bufT.WriteString("&")
			}

			bufT.WriteString(k)
			bufT.WriteString("=")
			bufT.WriteString(vi)
		}

		strT := bufT.String() + "&key=" + keyA

		if methodT == "sha256" {
			h := sha256.New()
			h.Write([]byte(strT))
			d := h.Sum(nil)
			hmacT = string(d)
		} else {
			hmacT = strings.ToUpper(MD5Encrypt(strT))
		}

		return hmacT
	}

	return ""

	// strT := bufT.String() + "&key=" + keyA

	// // hmacT := hmac.New(sha256.New, []byte(keyA))
	// // hmacT := hmac.New(md5.New, []byte(strT))
	// hmacT := strings.ToUpper(MD5Encrypt(strT))

	// // hmacT.Write([]byte(strT))
	// return hmacT

	// return strings.ToUpper(hex.EncodeToString(hmacT.Sum([]byte(""))))
}

var WeixinPaySignString = TKX.WeixinPaySignString

func (pA *TK) SignRSAWithSHA256(strA string, keyA string) string {
	bytesToSign := []byte(strA)

	block, err8 := pem.Decode([]byte("-----BEGIN RSA PRIVATE KEY-----\n" + keyA + "\n-----END RSA PRIVATE KEY-----")) //-----BEGIN RSA PRIVATE KEY----
	if err8 != nil {
		// tk.Pl("Error trying decode endorser private key: %v - %v", block, err8)
		// return ""
	}

	key, _ := x509.ParsePKCS1PrivateKey(block.Bytes)

	h := sha256.New()
	h.Write(bytesToSign)
	d := h.Sum(nil)

	signature, err := rsa.SignPKCS1v15(cryptorand.Reader, key, crypto.SHA256, d)
	if err != nil {
		// tk.Pl("Error trying decode endorser private key")
		return ErrStrf("Error trying decode endorser private key: %v", err)
	}

	signatureString := base64.StdEncoding.EncodeToString(signature)

	return signatureString
}

var SignRSAWithSHA256 = TKX.SignRSAWithSHA256

func (pA *TK) SignWithHS512(strA string, keyA []byte) string {
	h := hmac.New(sha512.New, keyA)
	h.Write([]byte(strA))
	return base64.RawURLEncoding.EncodeToString([]byte(h.Sum(nil)))
}

var SignWithHS512 = TKX.SignWithHS512

func (pA *TK) GenerateJwtToken(objA interface{}, keyA []byte, optsA ...string) string {
	json1T, errT := json.Marshal(objA)
	
	if errT != nil {
		return "TXERROR:"+errT.Error()
	}
	
	var headStrT string
	
	if IfSwitchExists(optsA, "-noType") {
		headStrT =  `{"alg":"HS512"}`
	} else {
		headStrT = `{"alg":"HS512","typ":"JWT"}`
	}

	tmps := base64.RawURLEncoding.EncodeToString([]byte(headStrT)) + "." + base64.RawURLEncoding.EncodeToString([]byte(json1T))
	
	signedT := SignWithHS512(tmps, keyA)
	
	return tmps+"."+signedT
}

var GenerateJwtToken = TKX.GenerateJwtToken

func (pA *TK) ParseJwtToken(strA string, keyA []byte, optsA ...string) interface{} {
	listT := strings.SplitN(strA, ".", 3)
	
	if len(listT) < 3 {
		return fmt.Errorf("not enough parts: %v", len(listT))
	}
	
	bufT, errT := base64.RawURLEncoding.DecodeString(listT[0])

	if errT != nil {
		return fmt.Errorf("failed to decode base64-raw-url: %v", errT)
	}

	headStrT := string(bufT)
	
	bufT, errT = base64.RawURLEncoding.DecodeString(listT[1])

	if errT != nil {
		return fmt.Errorf("failed to decode base64-raw-url: %v", errT)
	}

	bodyStrT := string(bufT)
	
	tmps := base64.RawURLEncoding.EncodeToString([]byte(headStrT)) + "." + base64.RawURLEncoding.EncodeToString([]byte(bodyStrT))
	
	if !IfSwitchExists(optsA, "-noValidate") {
		signedT := SignWithHS512(tmps, keyA)
		
		if signedT != listT[2] {
			return fmt.Errorf("invalid signature")
		}
	}
	
	return JSONToMapStringAny(bodyStrT)
}

var ParseJwtToken = TKX.ParseJwtToken

func (pA *TK) AlipaySignString(valuesA interface{}, keyA string) string {
	if valuesA == nil {
		return ""
	}

	var bufT strings.Builder

	nv1, ok := valuesA.(url.Values)

	if ok {
		keysT := make([]string, 0, len(nv1))

		for k := range nv1 {
			keysT = append(keysT, k)
		}

		sort.Strings(keysT)

		for i, k := range keysT {
			vi := nv1[k]

			if len(vi) < 1 {
				continue
			}

			if i > 0 {
				bufT.WriteString("&")
			}

			bufT.WriteString(k)
			bufT.WriteString("=")
			bufT.WriteString(vi[0])
		}

		return SignRSAWithSHA256(bufT.String(), keyA)
	}

	nv2, ok := valuesA.(map[string]string)

	if ok {
		keysT := make([]string, 0, len(nv2))

		for k := range nv2 {
			keysT = append(keysT, k)
		}

		sort.Strings(keysT)

		for i, k := range keysT {
			vi := nv2[k]

			if len(vi) < 1 {
				continue
			}

			if i > 0 {
				bufT.WriteString("&")
			}

			bufT.WriteString(k)
			bufT.WriteString("=")
			bufT.WriteString(vi)
		}

		return SignRSAWithSHA256(bufT.String(), keyA)
	}

	nv3, ok := valuesA.(map[string]interface{})

	if ok {
		keysT := make([]string, 0, len(nv3))

		for k := range nv3 {
			keysT = append(keysT, k)
		}

		sort.Strings(keysT)

		for i, k := range keysT {
			vi := ToStr(nv3[k])

			if len(vi) < 1 {
				continue
			}

			if i > 0 {
				bufT.WriteString("&")
			}

			bufT.WriteString(k)
			bufT.WriteString("=")
			bufT.WriteString(vi)
		}

		return SignRSAWithSHA256(bufT.String(), keyA)
	}

	return ""

}

var AlipaySignString = TKX.AlipaySignString

// image related

func (pA *TK) LoadImageFromBytes(dataA []byte, argsA ...string) interface{} {
	var imgT image.Image = nil
	var errT error

	readerT := bytes.NewReader(dataA)

	typeT := strings.TrimSpace(GetSwitch(argsA, "-type=", ""))

	typeT = strings.ToLower(strings.TrimPrefix(typeT, "."))
	// Pl("LoadImageFromBytes: %#v", typeT)

	switch typeT {
	case "", "jpg", "jpeg":
		imgT, errT = jpeg.Decode(readerT)

		if errT != nil {
			return errT
		}
	case "png":
		imgT, errT = png.Decode(readerT)

		if errT != nil {
			return errT
		}
	case "gif":
		imgT, errT = gif.Decode(readerT)

		if errT != nil {
			return errT
		}
	case "bmp":
		imgT, errT = bmp.Decode(readerT)

		if errT != nil {
			return errT
		}
	default:
		return fmt.Errorf("unsupported image type")
	}

	return imgT
}

var LoadImageFromBytes = TKX.LoadImageFromBytes

func (pA *TK) LoadImageFromFile(pathA string, argsA ...string) interface{} {
	var imgT image.Image = nil
	var errT error

	var fileT *os.File

	fileT, errT = os.Open(pathA)

	if errT != nil {
		return errT
	}

	defer fileT.Close()

	readerT := bufio.NewReader(fileT)

	extT := strings.ToLower(filepath.Ext(pathA))

	typeT := strings.TrimSpace(GetSwitch(argsA, "-type=", extT))

	typeT = strings.Trim(typeT, ".")

	switch typeT {
	case "", "jpg", "jpeg":
		imgT, errT = jpeg.Decode(readerT)

		if errT != nil {
			return errT
		}
	case "png":
		imgT, errT = png.Decode(readerT)

		if errT != nil {
			return errT
		}
	case "gif":
		imgT, errT = gif.Decode(readerT)

		if errT != nil {
			return errT
		}
	case "bmp":
		imgT, errT = bmp.Decode(readerT)

		if errT != nil {
			return errT
		}
	default:
		return fmt.Errorf("unsupported image type")
	}

	return imgT
}

var LoadImageFromFile = TKX.LoadImageFromFile

func (p *TK) SaveImageToFile(imageA image.Image, filePathA string, formatA ...string) error {
	fileT, errT := os.Create(filePathA)
	if errT != nil {
		return errT
	}
	defer fileT.Close()

	var formatT string

	if formatA == nil || len(formatA) < 1 {
		formatT = ".png"
	} else {
		formatT = strings.ToLower(formatA[0])
	}

	formatT = strings.Trim(formatT, ".")

	switch formatT {
	case "", "png":
		errT = png.Encode(fileT, imageA)
	case "jpg", "jpeg":
		errT = jpeg.Encode(fileT, imageA, nil)
	case "gif":
		errT = gif.Encode(fileT, imageA, nil)
	case "bmp":
		errT = bmp.Encode(fileT, imageA)
	default:
		errT = png.Encode(fileT, imageA)
	}

	return errT
}

var SaveImageToFile = TKX.SaveImageToFile

func (p *TK) SaveImageToBytes(imageA image.Image, formatA ...string) []byte {
	var bufT = new(bytes.Buffer)

	var formatT string

	var errT error

	if formatA == nil || len(formatA) < 1 {
		formatT = ".png"
	} else {
		formatT = strings.ToLower(strings.TrimSpace(formatA[0]))
	}

	formatT = strings.Trim(formatT, ".")

	switch formatT {
	case "", "png":
		errT = png.Encode(bufT, imageA)
	case "jpg", "jpeg":
		errT = jpeg.Encode(bufT, imageA, nil)
	case "gif":
		errT = gif.Encode(bufT, imageA, nil)
	case "bmp":
		errT = bmp.Encode(bufT, imageA)
	default:
		errT = png.Encode(bufT, imageA)
	}

	if errT != nil {
		return nil
	}

	return bufT.Bytes()
}

var SaveImageToBytes = TKX.SaveImageToBytes

type InterpolationFunction int

// InterpolationFunction constants
const (
	// Nearest-neighbor interpolation
	NearestNeighbor InterpolationFunction = iota
	// Bilinear interpolation
	Bilinear
	// Bicubic interpolation (with cubic hermite spline)
	Bicubic
	// Mitchell-Netravali interpolation
	MitchellNetravali
	// Lanczos interpolation (a=2)
	Lanczos2
	// Lanczos interpolation (a=3)
	Lanczos3
)

func nearest(in float64) float64 {
	if in >= -0.5 && in < 0.5 {
		return 1
	}
	return 0
}

func linear(in float64) float64 {
	in = math.Abs(in)
	if in <= 1 {
		return 1 - in
	}
	return 0
}

func cubic(in float64) float64 {
	in = math.Abs(in)
	if in <= 1 {
		return in*in*(1.5*in-2.5) + 1.0
	}
	if in <= 2 {
		return in*(in*(2.5-0.5*in)-4.0) + 2.0
	}
	return 0
}

func mitchellnetravali(in float64) float64 {
	in = math.Abs(in)
	if in <= 1 {
		return (7.0*in*in*in - 12.0*in*in + 5.33333333333) * 0.16666666666
	}
	if in <= 2 {
		return (-2.33333333333*in*in*in + 12.0*in*in - 20.0*in + 10.6666666667) * 0.16666666666
	}
	return 0
}

func sinc(x float64) float64 {
	x = math.Abs(x) * math.Pi
	if x >= 1.220703e-4 {
		return math.Sin(x) / x
	}
	return 1
}

func lanczos2(in float64) float64 {
	if in > -2 && in < 2 {
		return sinc(in) * sinc(in*0.5)
	}
	return 0
}

func lanczos3(in float64) float64 {
	if in > -3 && in < 3 {
		return sinc(in) * sinc(in*0.3333333333333333)
	}
	return 0
}

// kernal, returns an InterpolationFunctions taps and kernel.
func (i InterpolationFunction) kernel() (int, func(float64) float64) {
	switch i {
	case Bilinear:
		return 2, linear
	case Bicubic:
		return 4, cubic
	case MitchellNetravali:
		return 4, mitchellnetravali
	case Lanczos2:
		return 4, lanczos2
	case Lanczos3:
		return 6, lanczos3
	default:
		// Default to NearestNeighbor.
		return 2, nearest
	}
}

// values <1 will sharpen the image
var blur = 1.0

// Calculates scaling factors using old and new image dimensions.
func calcFactors(width, height uint, oldWidth, oldHeight float64) (scaleX, scaleY float64) {
	if width == 0 {
		if height == 0 {
			scaleX = 1.0
			scaleY = 1.0
		} else {
			scaleY = oldHeight / float64(height)
			scaleX = scaleY
		}
	} else {
		scaleX = oldWidth / float64(width)
		if height == 0 {
			scaleY = scaleX
		} else {
			scaleY = oldHeight / float64(height)
		}
	}
	return
}

func createWeights8(dy, filterLength int, blur, scale float64, kernel func(float64) float64) ([]int16, []int, int) {
	filterLength = filterLength * int(math.Max(math.Ceil(blur*scale), 1))
	filterFactor := math.Min(1./(blur*scale), 1)

	coeffs := make([]int16, dy*filterLength)
	start := make([]int, dy)
	for y := 0; y < dy; y++ {
		interpX := scale * (float64(y) + 0.5)
		start[y] = int(interpX) - filterLength/2 + 1
		interpX -= float64(start[y])
		for i := 0; i < filterLength; i++ {
			in := (interpX - float64(i)) * filterFactor
			coeffs[y*filterLength+i] = int16(kernel(in) * 256)
		}
	}

	return coeffs, start, filterLength
}

// range [-65536,65536]
func createWeights16(dy, filterLength int, blur, scale float64, kernel func(float64) float64) ([]int32, []int, int) {
	filterLength = filterLength * int(math.Max(math.Ceil(blur*scale), 1))
	filterFactor := math.Min(1./(blur*scale), 1)

	coeffs := make([]int32, dy*filterLength)
	start := make([]int, dy)
	for y := 0; y < dy; y++ {
		interpX := scale * (float64(y) + 0.5)
		start[y] = int(interpX) - filterLength/2 + 1
		interpX -= float64(start[y])
		for i := 0; i < filterLength; i++ {
			in := (interpX - float64(i)) * filterFactor
			coeffs[y*filterLength+i] = int32(kernel(in) * 65536)
		}
	}

	return coeffs, start, filterLength
}

func createWeightsNearest(dy, filterLength int, blur, scale float64) ([]bool, []int, int) {
	filterLength = filterLength * int(math.Max(math.Ceil(blur*scale), 1))
	filterFactor := math.Min(1./(blur*scale), 1)

	coeffs := make([]bool, dy*filterLength)
	start := make([]int, dy)
	for y := 0; y < dy; y++ {
		interpX := scale * (float64(y) + 0.5)
		start[y] = int(interpX) - filterLength/2 + 1
		interpX -= float64(start[y])
		for i := 0; i < filterLength; i++ {
			in := (interpX - float64(i)) * filterFactor
			if in >= -0.5 && in < 0.5 {
				coeffs[y*filterLength+i] = true
			} else {
				coeffs[y*filterLength+i] = false
			}
		}
	}

	return coeffs, start, filterLength
}

type imageWithSubImage interface {
	image.Image
	SubImage(image.Rectangle) image.Image
}

func makeSlice(img imageWithSubImage, i, n int) image.Image {
	return img.SubImage(image.Rect(img.Bounds().Min.X, img.Bounds().Min.Y+i*img.Bounds().Dy()/n, img.Bounds().Max.X, img.Bounds().Min.Y+(i+1)*img.Bounds().Dy()/n))
}

func nearestRGBA(in *image.RGBA, out *image.RGBA, scale float64, coeffs []bool, offset []int, filterLength int) {
	newBounds := out.Bounds()
	maxX := in.Bounds().Dx() - 1

	for x := newBounds.Min.X; x < newBounds.Max.X; x++ {
		row := in.Pix[x*in.Stride:]
		for y := newBounds.Min.Y; y < newBounds.Max.Y; y++ {
			var rgba [4]float32
			var sum float32
			start := offset[y]
			ci := y * filterLength
			for i := 0; i < filterLength; i++ {
				if coeffs[ci+i] {
					xi := start + i
					switch {
					case xi < 0:
						xi = 0
					case xi >= maxX:
						xi = 4 * maxX
					default:
						xi *= 4
					}
					rgba[0] += float32(row[xi+0])
					rgba[1] += float32(row[xi+1])
					rgba[2] += float32(row[xi+2])
					rgba[3] += float32(row[xi+3])
					sum++
				}
			}

			xo := (y-newBounds.Min.Y)*out.Stride + (x-newBounds.Min.X)*4
			out.Pix[xo+0] = floatToUint8(rgba[0] / sum)
			out.Pix[xo+1] = floatToUint8(rgba[1] / sum)
			out.Pix[xo+2] = floatToUint8(rgba[2] / sum)
			out.Pix[xo+3] = floatToUint8(rgba[3] / sum)
		}
	}
}

func nearestRGBA64(in *image.RGBA64, out *image.RGBA64, scale float64, coeffs []bool, offset []int, filterLength int) {
	newBounds := out.Bounds()
	maxX := in.Bounds().Dx() - 1

	for x := newBounds.Min.X; x < newBounds.Max.X; x++ {
		row := in.Pix[x*in.Stride:]
		for y := newBounds.Min.Y; y < newBounds.Max.Y; y++ {
			var rgba [4]float32
			var sum float32
			start := offset[y]
			ci := y * filterLength
			for i := 0; i < filterLength; i++ {
				if coeffs[ci+i] {
					xi := start + i
					switch {
					case xi < 0:
						xi = 0
					case xi >= maxX:
						xi = 8 * maxX
					default:
						xi *= 8
					}
					rgba[0] += float32(uint16(row[xi+0])<<8 | uint16(row[xi+1]))
					rgba[1] += float32(uint16(row[xi+2])<<8 | uint16(row[xi+3]))
					rgba[2] += float32(uint16(row[xi+4])<<8 | uint16(row[xi+5]))
					rgba[3] += float32(uint16(row[xi+6])<<8 | uint16(row[xi+7]))
					sum++
				}
			}

			xo := (y-newBounds.Min.Y)*out.Stride + (x-newBounds.Min.X)*8
			value := floatToUint16(rgba[0] / sum)
			out.Pix[xo+0] = uint8(value >> 8)
			out.Pix[xo+1] = uint8(value)
			value = floatToUint16(rgba[1] / sum)
			out.Pix[xo+2] = uint8(value >> 8)
			out.Pix[xo+3] = uint8(value)
			value = floatToUint16(rgba[2] / sum)
			out.Pix[xo+4] = uint8(value >> 8)
			out.Pix[xo+5] = uint8(value)
			value = floatToUint16(rgba[3] / sum)
			out.Pix[xo+6] = uint8(value >> 8)
			out.Pix[xo+7] = uint8(value)
		}
	}
}

func nearestGray(in *image.Gray, out *image.Gray, scale float64, coeffs []bool, offset []int, filterLength int) {
	newBounds := out.Bounds()
	maxX := in.Bounds().Dx() - 1

	for x := newBounds.Min.X; x < newBounds.Max.X; x++ {
		row := in.Pix[x*in.Stride:]
		for y := newBounds.Min.Y; y < newBounds.Max.Y; y++ {
			var gray float32
			var sum float32
			start := offset[y]
			ci := y * filterLength
			for i := 0; i < filterLength; i++ {
				if coeffs[ci+i] {
					xi := start + i
					switch {
					case xi < 0:
						xi = 0
					case xi >= maxX:
						xi = maxX
					}
					gray += float32(row[xi])
					sum++
				}
			}

			offset := (y-newBounds.Min.Y)*out.Stride + (x - newBounds.Min.X)
			out.Pix[offset] = floatToUint8(gray / sum)
		}
	}
}

func floatToUint8(x float32) uint8 {
	// Nearest-neighbor values are always
	// positive no need to check lower-bound.
	if x > 0xfe {
		return 0xff
	}
	return uint8(x)
}

func floatToUint16(x float32) uint16 {
	if x > 0xfffe {
		return 0xffff
	}
	return uint16(x)
}

// ycc is an in memory YCbCr image.  The Y, Cb and Cr samples are held in a
// single slice to increase resizing performance.
type ycc struct {
	// Pix holds the image's pixels, in Y, Cb, Cr order. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*3].
	Pix []uint8
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect image.Rectangle
	// SubsampleRatio is the subsample ratio of the original YCbCr image.
	SubsampleRatio image.YCbCrSubsampleRatio
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (p *ycc) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*3
}

func (p *ycc) Bounds() image.Rectangle {
	return p.Rect
}

func (p *ycc) ColorModel() color.Model {
	return color.YCbCrModel
}

func (p *ycc) At(x, y int) color.Color {
	if !(image.Point{x, y}.In(p.Rect)) {
		return color.YCbCr{}
	}
	i := p.PixOffset(x, y)
	return color.YCbCr{
		p.Pix[i+0],
		p.Pix[i+1],
		p.Pix[i+2],
	}
}

func (p *ycc) Opaque() bool {
	return true
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *ycc) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	if r.Empty() {
		return &ycc{SubsampleRatio: p.SubsampleRatio}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	return &ycc{
		Pix:            p.Pix[i:],
		Stride:         p.Stride,
		Rect:           r,
		SubsampleRatio: p.SubsampleRatio,
	}
}

// newYCC returns a new ycc with the given bounds and subsample ratio.
func newYCC(r image.Rectangle, s image.YCbCrSubsampleRatio) *ycc {
	w, h := r.Dx(), r.Dy()
	buf := make([]uint8, 3*w*h)
	return &ycc{Pix: buf, Stride: 3 * w, Rect: r, SubsampleRatio: s}
}

// YCbCr converts ycc to a YCbCr image with the same subsample ratio
// as the YCbCr image that ycc was generated from.
func (p *ycc) YCbCr() *image.YCbCr {
	ycbcr := image.NewYCbCr(p.Rect, p.SubsampleRatio)
	var off int

	switch ycbcr.SubsampleRatio {
	case image.YCbCrSubsampleRatio422:
		for y := ycbcr.Rect.Min.Y; y < ycbcr.Rect.Max.Y; y++ {
			yy := (y - ycbcr.Rect.Min.Y) * ycbcr.YStride
			cy := (y - ycbcr.Rect.Min.Y) * ycbcr.CStride
			for x := ycbcr.Rect.Min.X; x < ycbcr.Rect.Max.X; x++ {
				xx := (x - ycbcr.Rect.Min.X)
				yi := yy + xx
				ci := cy + xx/2
				ycbcr.Y[yi] = p.Pix[off+0]
				ycbcr.Cb[ci] = p.Pix[off+1]
				ycbcr.Cr[ci] = p.Pix[off+2]
				off += 3
			}
		}
	case image.YCbCrSubsampleRatio420:
		for y := ycbcr.Rect.Min.Y; y < ycbcr.Rect.Max.Y; y++ {
			yy := (y - ycbcr.Rect.Min.Y) * ycbcr.YStride
			cy := (y/2 - ycbcr.Rect.Min.Y/2) * ycbcr.CStride
			for x := ycbcr.Rect.Min.X; x < ycbcr.Rect.Max.X; x++ {
				xx := (x - ycbcr.Rect.Min.X)
				yi := yy + xx
				ci := cy + xx/2
				ycbcr.Y[yi] = p.Pix[off+0]
				ycbcr.Cb[ci] = p.Pix[off+1]
				ycbcr.Cr[ci] = p.Pix[off+2]
				off += 3
			}
		}
	case image.YCbCrSubsampleRatio440:
		for y := ycbcr.Rect.Min.Y; y < ycbcr.Rect.Max.Y; y++ {
			yy := (y - ycbcr.Rect.Min.Y) * ycbcr.YStride
			cy := (y/2 - ycbcr.Rect.Min.Y/2) * ycbcr.CStride
			for x := ycbcr.Rect.Min.X; x < ycbcr.Rect.Max.X; x++ {
				xx := (x - ycbcr.Rect.Min.X)
				yi := yy + xx
				ci := cy + xx
				ycbcr.Y[yi] = p.Pix[off+0]
				ycbcr.Cb[ci] = p.Pix[off+1]
				ycbcr.Cr[ci] = p.Pix[off+2]
				off += 3
			}
		}
	default:
		// Default to 4:4:4 subsampling.
		for y := ycbcr.Rect.Min.Y; y < ycbcr.Rect.Max.Y; y++ {
			yy := (y - ycbcr.Rect.Min.Y) * ycbcr.YStride
			cy := (y - ycbcr.Rect.Min.Y) * ycbcr.CStride
			for x := ycbcr.Rect.Min.X; x < ycbcr.Rect.Max.X; x++ {
				xx := (x - ycbcr.Rect.Min.X)
				yi := yy + xx
				ci := cy + xx
				ycbcr.Y[yi] = p.Pix[off+0]
				ycbcr.Cb[ci] = p.Pix[off+1]
				ycbcr.Cr[ci] = p.Pix[off+2]
				off += 3
			}
		}
	}
	return ycbcr
}

// imageYCbCrToYCC converts a YCbCr image to a ycc image for resizing.
func imageYCbCrToYCC(in *image.YCbCr) *ycc {
	w, h := in.Rect.Dx(), in.Rect.Dy()
	r := image.Rect(0, 0, w, h)
	buf := make([]uint8, 3*w*h)
	p := ycc{Pix: buf, Stride: 3 * w, Rect: r, SubsampleRatio: in.SubsampleRatio}
	var off int

	switch in.SubsampleRatio {
	case image.YCbCrSubsampleRatio422:
		for y := in.Rect.Min.Y; y < in.Rect.Max.Y; y++ {
			yy := (y - in.Rect.Min.Y) * in.YStride
			cy := (y - in.Rect.Min.Y) * in.CStride
			for x := in.Rect.Min.X; x < in.Rect.Max.X; x++ {
				xx := (x - in.Rect.Min.X)
				yi := yy + xx
				ci := cy + xx/2
				p.Pix[off+0] = in.Y[yi]
				p.Pix[off+1] = in.Cb[ci]
				p.Pix[off+2] = in.Cr[ci]
				off += 3
			}
		}
	case image.YCbCrSubsampleRatio420:
		for y := in.Rect.Min.Y; y < in.Rect.Max.Y; y++ {
			yy := (y - in.Rect.Min.Y) * in.YStride
			cy := (y/2 - in.Rect.Min.Y/2) * in.CStride
			for x := in.Rect.Min.X; x < in.Rect.Max.X; x++ {
				xx := (x - in.Rect.Min.X)
				yi := yy + xx
				ci := cy + xx/2
				p.Pix[off+0] = in.Y[yi]
				p.Pix[off+1] = in.Cb[ci]
				p.Pix[off+2] = in.Cr[ci]
				off += 3
			}
		}
	case image.YCbCrSubsampleRatio440:
		for y := in.Rect.Min.Y; y < in.Rect.Max.Y; y++ {
			yy := (y - in.Rect.Min.Y) * in.YStride
			cy := (y/2 - in.Rect.Min.Y/2) * in.CStride
			for x := in.Rect.Min.X; x < in.Rect.Max.X; x++ {
				xx := (x - in.Rect.Min.X)
				yi := yy + xx
				ci := cy + xx
				p.Pix[off+0] = in.Y[yi]
				p.Pix[off+1] = in.Cb[ci]
				p.Pix[off+2] = in.Cr[ci]
				off += 3
			}
		}
	default:
		// Default to 4:4:4 subsampling.
		for y := in.Rect.Min.Y; y < in.Rect.Max.Y; y++ {
			yy := (y - in.Rect.Min.Y) * in.YStride
			cy := (y - in.Rect.Min.Y) * in.CStride
			for x := in.Rect.Min.X; x < in.Rect.Max.X; x++ {
				xx := (x - in.Rect.Min.X)
				yi := yy + xx
				ci := cy + xx
				p.Pix[off+0] = in.Y[yi]
				p.Pix[off+1] = in.Cb[ci]
				p.Pix[off+2] = in.Cr[ci]
				off += 3
			}
		}
	}
	return &p
}

func nearestGray16(in *image.Gray16, out *image.Gray16, scale float64, coeffs []bool, offset []int, filterLength int) {
	newBounds := out.Bounds()
	maxX := in.Bounds().Dx() - 1

	for x := newBounds.Min.X; x < newBounds.Max.X; x++ {
		row := in.Pix[x*in.Stride:]
		for y := newBounds.Min.Y; y < newBounds.Max.Y; y++ {
			var gray float32
			var sum float32
			start := offset[y]
			ci := y * filterLength
			for i := 0; i < filterLength; i++ {
				if coeffs[ci+i] {
					xi := start + i
					switch {
					case xi < 0:
						xi = 0
					case xi >= maxX:
						xi = 2 * maxX
					default:
						xi *= 2
					}
					gray += float32(uint16(row[xi+0])<<8 | uint16(row[xi+1]))
					sum++
				}
			}

			offset := (y-newBounds.Min.Y)*out.Stride + (x-newBounds.Min.X)*2
			value := floatToUint16(gray / sum)
			out.Pix[offset+0] = uint8(value >> 8)
			out.Pix[offset+1] = uint8(value)
		}
	}
}

// Keep value in [0,255] range.
func clampUint8(in int32) uint8 {
	if in < 0 {
		return 0
	}
	if in > 255 {
		return 255
	}
	return uint8(in)
}

// Keep value in [0,65535] range.
func clampUint16(in int64) uint16 {
	if in < 0 {
		return 0
	}
	if in > 65535 {
		return 65535
	}
	return uint16(in)
}

func nearestYCbCr(in *ycc, out *ycc, scale float64, coeffs []bool, offset []int, filterLength int) {
	newBounds := out.Bounds()
	maxX := in.Bounds().Dx() - 1

	for x := newBounds.Min.X; x < newBounds.Max.X; x++ {
		row := in.Pix[x*in.Stride:]
		for y := newBounds.Min.Y; y < newBounds.Max.Y; y++ {
			var p [3]float32
			var sum float32
			start := offset[y]
			ci := y * filterLength
			for i := 0; i < filterLength; i++ {
				if coeffs[ci+i] {
					xi := start + i
					switch {
					case xi < 0:
						xi = 0
					case xi >= maxX:
						xi = 3 * maxX
					default:
						xi *= 3
					}
					p[0] += float32(row[xi+0])
					p[1] += float32(row[xi+1])
					p[2] += float32(row[xi+2])
					sum++
				}
			}

			xo := (y-newBounds.Min.Y)*out.Stride + (x-newBounds.Min.X)*3
			out.Pix[xo+0] = floatToUint8(p[0] / sum)
			out.Pix[xo+1] = floatToUint8(p[1] / sum)
			out.Pix[xo+2] = floatToUint8(p[2] / sum)
		}
	}
}

func nearestGeneric(in image.Image, out *image.RGBA64, scale float64, coeffs []bool, offset []int, filterLength int) {
	newBounds := out.Bounds()
	maxX := in.Bounds().Dx() - 1

	for x := newBounds.Min.X; x < newBounds.Max.X; x++ {
		for y := newBounds.Min.Y; y < newBounds.Max.Y; y++ {
			var rgba [4]float32
			var sum float32
			start := offset[y]
			ci := y * filterLength
			for i := 0; i < filterLength; i++ {
				if coeffs[ci+i] {
					xi := start + i
					switch {
					case xi < 0:
						xi = 0
					case xi >= maxX:
						xi = maxX
					}
					r, g, b, a := in.At(xi+in.Bounds().Min.X, x+in.Bounds().Min.Y).RGBA()
					rgba[0] += float32(r)
					rgba[1] += float32(g)
					rgba[2] += float32(b)
					rgba[3] += float32(a)
					sum++
				}
			}

			offset := (y-newBounds.Min.Y)*out.Stride + (x-newBounds.Min.X)*8
			value := floatToUint16(rgba[0] / sum)
			out.Pix[offset+0] = uint8(value >> 8)
			out.Pix[offset+1] = uint8(value)
			value = floatToUint16(rgba[1] / sum)
			out.Pix[offset+2] = uint8(value >> 8)
			out.Pix[offset+3] = uint8(value)
			value = floatToUint16(rgba[2] / sum)
			out.Pix[offset+4] = uint8(value >> 8)
			out.Pix[offset+5] = uint8(value)
			value = floatToUint16(rgba[3] / sum)
			out.Pix[offset+6] = uint8(value >> 8)
			out.Pix[offset+7] = uint8(value)
		}
	}
}

func resizeGeneric(in image.Image, out *image.RGBA64, scale float64, coeffs []int32, offset []int, filterLength int) {
	newBounds := out.Bounds()
	maxX := in.Bounds().Dx() - 1

	for x := newBounds.Min.X; x < newBounds.Max.X; x++ {
		for y := newBounds.Min.Y; y < newBounds.Max.Y; y++ {
			var rgba [4]int64
			var sum int64
			start := offset[y]
			ci := y * filterLength
			for i := 0; i < filterLength; i++ {
				coeff := coeffs[ci+i]
				if coeff != 0 {
					xi := start + i
					switch {
					case xi < 0:
						xi = 0
					case xi >= maxX:
						xi = maxX
					}
					r, g, b, a := in.At(xi+in.Bounds().Min.X, x+in.Bounds().Min.Y).RGBA()
					rgba[0] += int64(coeff) * int64(r)
					rgba[1] += int64(coeff) * int64(g)
					rgba[2] += int64(coeff) * int64(b)
					rgba[3] += int64(coeff) * int64(a)
					sum += int64(coeff)
				}
			}

			offset := (y-newBounds.Min.Y)*out.Stride + (x-newBounds.Min.X)*8
			value := clampUint16(rgba[0] / sum)
			out.Pix[offset+0] = uint8(value >> 8)
			out.Pix[offset+1] = uint8(value)
			value = clampUint16(rgba[1] / sum)
			out.Pix[offset+2] = uint8(value >> 8)
			out.Pix[offset+3] = uint8(value)
			value = clampUint16(rgba[2] / sum)
			out.Pix[offset+4] = uint8(value >> 8)
			out.Pix[offset+5] = uint8(value)
			value = clampUint16(rgba[3] / sum)
			out.Pix[offset+6] = uint8(value >> 8)
			out.Pix[offset+7] = uint8(value)
		}
	}
}

func resizeRGBA(in *image.RGBA, out *image.RGBA, scale float64, coeffs []int16, offset []int, filterLength int) {
	newBounds := out.Bounds()
	maxX := in.Bounds().Dx() - 1

	for x := newBounds.Min.X; x < newBounds.Max.X; x++ {
		row := in.Pix[x*in.Stride:]
		for y := newBounds.Min.Y; y < newBounds.Max.Y; y++ {
			var rgba [4]int32
			var sum int32
			start := offset[y]
			ci := y * filterLength
			for i := 0; i < filterLength; i++ {
				coeff := coeffs[ci+i]
				if coeff != 0 {
					xi := start + i
					switch {
					case xi < 0:
						xi = 0
					case xi >= maxX:
						xi = 4 * maxX
					default:
						xi *= 4
					}
					rgba[0] += int32(coeff) * int32(row[xi+0])
					rgba[1] += int32(coeff) * int32(row[xi+1])
					rgba[2] += int32(coeff) * int32(row[xi+2])
					rgba[3] += int32(coeff) * int32(row[xi+3])
					sum += int32(coeff)
				}
			}

			xo := (y-newBounds.Min.Y)*out.Stride + (x-newBounds.Min.X)*4
			out.Pix[xo+0] = clampUint8(rgba[0] / sum)
			out.Pix[xo+1] = clampUint8(rgba[1] / sum)
			out.Pix[xo+2] = clampUint8(rgba[2] / sum)
			out.Pix[xo+3] = clampUint8(rgba[3] / sum)
		}
	}
}

func resizeRGBA64(in *image.RGBA64, out *image.RGBA64, scale float64, coeffs []int32, offset []int, filterLength int) {
	newBounds := out.Bounds()
	maxX := in.Bounds().Dx() - 1

	for x := newBounds.Min.X; x < newBounds.Max.X; x++ {
		row := in.Pix[x*in.Stride:]
		for y := newBounds.Min.Y; y < newBounds.Max.Y; y++ {
			var rgba [4]int64
			var sum int64
			start := offset[y]
			ci := y * filterLength
			for i := 0; i < filterLength; i++ {
				coeff := coeffs[ci+i]
				if coeff != 0 {
					xi := start + i
					switch {
					case xi < 0:
						xi = 0
					case xi >= maxX:
						xi = 8 * maxX
					default:
						xi *= 8
					}
					rgba[0] += int64(coeff) * int64(uint16(row[xi+0])<<8|uint16(row[xi+1]))
					rgba[1] += int64(coeff) * int64(uint16(row[xi+2])<<8|uint16(row[xi+3]))
					rgba[2] += int64(coeff) * int64(uint16(row[xi+4])<<8|uint16(row[xi+5]))
					rgba[3] += int64(coeff) * int64(uint16(row[xi+6])<<8|uint16(row[xi+7]))
					sum += int64(coeff)
				}
			}

			xo := (y-newBounds.Min.Y)*out.Stride + (x-newBounds.Min.X)*8
			value := clampUint16(rgba[0] / sum)
			out.Pix[xo+0] = uint8(value >> 8)
			out.Pix[xo+1] = uint8(value)
			value = clampUint16(rgba[1] / sum)
			out.Pix[xo+2] = uint8(value >> 8)
			out.Pix[xo+3] = uint8(value)
			value = clampUint16(rgba[2] / sum)
			out.Pix[xo+4] = uint8(value >> 8)
			out.Pix[xo+5] = uint8(value)
			value = clampUint16(rgba[3] / sum)
			out.Pix[xo+6] = uint8(value >> 8)
			out.Pix[xo+7] = uint8(value)
		}
	}
}

func resizeGray(in *image.Gray, out *image.Gray, scale float64, coeffs []int16, offset []int, filterLength int) {
	newBounds := out.Bounds()
	maxX := in.Bounds().Dx() - 1

	for x := newBounds.Min.X; x < newBounds.Max.X; x++ {
		row := in.Pix[(x-newBounds.Min.X)*in.Stride:]
		for y := newBounds.Min.Y; y < newBounds.Max.Y; y++ {
			var gray int32
			var sum int32
			start := offset[y]
			ci := y * filterLength
			for i := 0; i < filterLength; i++ {
				coeff := coeffs[ci+i]
				if coeff != 0 {
					xi := start + i
					switch {
					case xi < 0:
						xi = 0
					case xi >= maxX:
						xi = maxX
					}
					gray += int32(coeff) * int32(row[xi])
					sum += int32(coeff)
				}
			}

			offset := (y-newBounds.Min.Y)*out.Stride + (x - newBounds.Min.X)
			out.Pix[offset] = clampUint8(gray / sum)
		}
	}
}

func resizeGray16(in *image.Gray16, out *image.Gray16, scale float64, coeffs []int32, offset []int, filterLength int) {
	newBounds := out.Bounds()
	maxX := in.Bounds().Dx() - 1

	for x := newBounds.Min.X; x < newBounds.Max.X; x++ {
		row := in.Pix[x*in.Stride:]
		for y := newBounds.Min.Y; y < newBounds.Max.Y; y++ {
			var gray int64
			var sum int64
			start := offset[y]
			ci := y * filterLength
			for i := 0; i < filterLength; i++ {
				coeff := coeffs[ci+i]
				if coeff != 0 {
					xi := start + i
					switch {
					case xi < 0:
						xi = 0
					case xi >= maxX:
						xi = 2 * maxX
					default:
						xi *= 2
					}
					gray += int64(coeff) * int64(uint16(row[xi+0])<<8|uint16(row[xi+1]))
					sum += int64(coeff)
				}
			}

			offset := (y-newBounds.Min.Y)*out.Stride + (x-newBounds.Min.X)*2
			value := clampUint16(gray / sum)
			out.Pix[offset+0] = uint8(value >> 8)
			out.Pix[offset+1] = uint8(value)
		}
	}
}

func resizeYCbCr(in *ycc, out *ycc, scale float64, coeffs []int16, offset []int, filterLength int) {
	newBounds := out.Bounds()
	maxX := in.Bounds().Dx() - 1

	for x := newBounds.Min.X; x < newBounds.Max.X; x++ {
		row := in.Pix[x*in.Stride:]
		for y := newBounds.Min.Y; y < newBounds.Max.Y; y++ {
			var p [3]int32
			var sum int32
			start := offset[y]
			ci := y * filterLength
			for i := 0; i < filterLength; i++ {
				coeff := coeffs[ci+i]
				if coeff != 0 {
					xi := start + i
					switch {
					case xi < 0:
						xi = 0
					case xi >= maxX:
						xi = 3 * maxX
					default:
						xi *= 3
					}
					p[0] += int32(coeff) * int32(row[xi+0])
					p[1] += int32(coeff) * int32(row[xi+1])
					p[2] += int32(coeff) * int32(row[xi+2])
					sum += int32(coeff)
				}
			}

			xo := (y-newBounds.Min.Y)*out.Stride + (x-newBounds.Min.X)*3
			out.Pix[xo+0] = clampUint8(p[0] / sum)
			out.Pix[xo+1] = clampUint8(p[1] / sum)
			out.Pix[xo+2] = clampUint8(p[2] / sum)
		}
	}
}

func resizeNearest(width, height uint, scaleX, scaleY float64, img image.Image, interp InterpolationFunction) image.Image {
	taps, _ := interp.kernel()
	cpus := runtime.NumCPU()
	wg := sync.WaitGroup{}

	switch input := img.(type) {
	case *image.RGBA:
		// 8-bit precision
		temp := image.NewRGBA(image.Rect(0, 0, input.Bounds().Dy(), int(width)))
		result := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))

		// horizontal filter, results in transposed temporary image
		coeffs, offset, filterLength := createWeightsNearest(temp.Bounds().Dy(), taps, blur, scaleX)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(temp, i, cpus).(*image.RGBA)
			go func() {
				defer wg.Done()
				nearestRGBA(input, slice, scaleX, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()

		// horizontal filter on transposed image, result is not transposed
		coeffs, offset, filterLength = createWeightsNearest(result.Bounds().Dy(), taps, blur, scaleY)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(result, i, cpus).(*image.RGBA)
			go func() {
				defer wg.Done()
				nearestRGBA(temp, slice, scaleY, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()
		return result
	case *image.YCbCr:
		// 8-bit precision
		// accessing the YCbCr arrays in a tight loop is slow.
		// converting the image to ycc increases performance by 2x.
		temp := newYCC(image.Rect(0, 0, input.Bounds().Dy(), int(width)), input.SubsampleRatio)
		result := newYCC(image.Rect(0, 0, int(width), int(height)), input.SubsampleRatio)

		coeffs, offset, filterLength := createWeightsNearest(temp.Bounds().Dy(), taps, blur, scaleX)
		in := imageYCbCrToYCC(input)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(temp, i, cpus).(*ycc)
			go func() {
				defer wg.Done()
				nearestYCbCr(in, slice, scaleX, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()

		coeffs, offset, filterLength = createWeightsNearest(result.Bounds().Dy(), taps, blur, scaleY)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(result, i, cpus).(*ycc)
			go func() {
				defer wg.Done()
				nearestYCbCr(temp, slice, scaleY, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()
		return result.YCbCr()
	case *image.RGBA64:
		// 16-bit precision
		temp := image.NewRGBA64(image.Rect(0, 0, input.Bounds().Dy(), int(width)))
		result := image.NewRGBA64(image.Rect(0, 0, int(width), int(height)))

		// horizontal filter, results in transposed temporary image
		coeffs, offset, filterLength := createWeightsNearest(temp.Bounds().Dy(), taps, blur, scaleX)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(temp, i, cpus).(*image.RGBA64)
			go func() {
				defer wg.Done()
				nearestRGBA64(input, slice, scaleX, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()

		// horizontal filter on transposed image, result is not transposed
		coeffs, offset, filterLength = createWeightsNearest(result.Bounds().Dy(), taps, blur, scaleY)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(result, i, cpus).(*image.RGBA64)
			go func() {
				defer wg.Done()
				nearestGeneric(temp, slice, scaleY, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()
		return result
	case *image.Gray:
		// 8-bit precision
		temp := image.NewGray(image.Rect(0, 0, input.Bounds().Dy(), int(width)))
		result := image.NewGray(image.Rect(0, 0, int(width), int(height)))

		// horizontal filter, results in transposed temporary image
		coeffs, offset, filterLength := createWeightsNearest(temp.Bounds().Dy(), taps, blur, scaleX)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(temp, i, cpus).(*image.Gray)
			go func() {
				defer wg.Done()
				nearestGray(input, slice, scaleX, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()

		// horizontal filter on transposed image, result is not transposed
		coeffs, offset, filterLength = createWeightsNearest(result.Bounds().Dy(), taps, blur, scaleY)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(result, i, cpus).(*image.Gray)
			go func() {
				defer wg.Done()
				nearestGray(temp, slice, scaleY, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()
		return result
	case *image.Gray16:
		// 16-bit precision
		temp := image.NewGray16(image.Rect(0, 0, input.Bounds().Dy(), int(width)))
		result := image.NewGray16(image.Rect(0, 0, int(width), int(height)))

		// horizontal filter, results in transposed temporary image
		coeffs, offset, filterLength := createWeightsNearest(temp.Bounds().Dy(), taps, blur, scaleX)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(temp, i, cpus).(*image.Gray16)
			go func() {
				defer wg.Done()
				nearestGray16(input, slice, scaleX, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()

		// horizontal filter on transposed image, result is not transposed
		coeffs, offset, filterLength = createWeightsNearest(result.Bounds().Dy(), taps, blur, scaleY)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(result, i, cpus).(*image.Gray16)
			go func() {
				defer wg.Done()
				nearestGray16(temp, slice, scaleY, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()
		return result
	default:
		// 16-bit precision
		temp := image.NewRGBA64(image.Rect(0, 0, img.Bounds().Dy(), int(width)))
		result := image.NewRGBA64(image.Rect(0, 0, int(width), int(height)))

		// horizontal filter, results in transposed temporary image
		coeffs, offset, filterLength := createWeightsNearest(temp.Bounds().Dy(), taps, blur, scaleX)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(temp, i, cpus).(*image.RGBA64)
			go func() {
				defer wg.Done()
				nearestGeneric(img, slice, scaleX, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()

		// horizontal filter on transposed image, result is not transposed
		coeffs, offset, filterLength = createWeightsNearest(result.Bounds().Dy(), taps, blur, scaleY)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(result, i, cpus).(*image.RGBA64)
			go func() {
				defer wg.Done()
				nearestRGBA64(temp, slice, scaleY, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()
		return result
	}

}

func (p *TK) CenterEnlargeImage(imgA image.Image, optsA ...string) image.Image {
	widthT := imgA.Bounds().Dx()
	heightT := imgA.Bounds().Dy()

	toWidthT := ToInt(GetSwitch(optsA, "-width=", "256"), 256)
	toHeightT := ToInt(GetSwitch(optsA, "-height=", "256"), 256)

	if toWidthT <= widthT && toHeightT <= heightT {
		return imgA
	}

	colorStrT := GetSwitch(optsA, "-color=", "#000000")

	if colorStrT == "" {
		colorStrT = "#000000"
	}

	r, g, b, a := ParseHexColor(colorStrT)

	colorT := color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)}

	destImageT := image.NewRGBA(image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: toWidthT, Y: toHeightT}})

	fillStartX := (toWidthT - widthT) / 2
	fillEndX := fillStartX + widthT

	fillStartY := (toHeightT - heightT) / 2
	fillEndY := fillStartY + heightT

	// Pln("fillStartX: ", fillStartX, "fillEndX: ", fillEndX, "fillStartY: ", fillStartY, "fillEndY: ", fillEndY)

	for y := 0; y < toHeightT; y++ {
		for x := 0; x < toWidthT; x++ {
			if x >= fillStartX && y >= fillStartY && x <= fillEndX && y <= fillEndY {
				destImageT.Set(x, y, imgA.At(x-fillStartX, y-fillStartY))
			} else {
				destImageT.Set(x, y, colorT)
			}
		}
	}

	return destImageT
}

var CenterEnlargeImage = TKX.CenterEnlargeImage

func (p *TK) ResizeImageX(imgA image.Image, optsA ...string) image.Image {
	widthT := ToInt(GetSwitch(optsA, "-width=", "256"), 256)
	heightT := ToInt(GetSwitch(optsA, "-height=", "256"), 256)

	return ResizeImage(widthT, heightT, imgA)
}

var ResizeImageX = TKX.ResizeImageX

func (p *TK) ResizeImage(widthA, heightA int, img image.Image, interpA ...InterpolationFunction) image.Image {
	width := uint(widthA)
	height := uint(heightA)
	scaleX, scaleY := calcFactors(width, height, float64(img.Bounds().Dx()), float64(img.Bounds().Dy()))
	if width == 0 {
		width = uint(0.7 + float64(img.Bounds().Dx())/scaleX)
	}
	if height == 0 {
		height = uint(0.7 + float64(img.Bounds().Dy())/scaleY)
	}

	var interp InterpolationFunction

	if interpA == nil || len(interpA) < 1 {
		interp = Lanczos3
	} else {
		interp = interpA[0]
	}

	if interp == NearestNeighbor {
		return resizeNearest(width, height, scaleX, scaleY, img, interp)
	}

	taps, kernel := interp.kernel()
	cpus := runtime.NumCPU()
	wg := sync.WaitGroup{}

	// Generic access to image.Image is slow in tight loops.
	// The optimal access has to be determined from the concrete image type.
	switch input := img.(type) {
	case *image.RGBA:
		// 8-bit precision
		temp := image.NewRGBA(image.Rect(0, 0, input.Bounds().Dy(), int(width)))
		result := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))

		// horizontal filter, results in transposed temporary image
		coeffs, offset, filterLength := createWeights8(temp.Bounds().Dy(), taps, blur, scaleX, kernel)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(temp, i, cpus).(*image.RGBA)
			go func() {
				defer wg.Done()
				resizeRGBA(input, slice, scaleX, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()

		// horizontal filter on transposed image, result is not transposed
		coeffs, offset, filterLength = createWeights8(result.Bounds().Dy(), taps, blur, scaleY, kernel)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(result, i, cpus).(*image.RGBA)
			go func() {
				defer wg.Done()
				resizeRGBA(temp, slice, scaleY, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()
		return result
	case *image.YCbCr:
		// 8-bit precision
		// accessing the YCbCr arrays in a tight loop is slow.
		// converting the image to ycc increases performance by 2x.
		temp := newYCC(image.Rect(0, 0, input.Bounds().Dy(), int(width)), input.SubsampleRatio)
		result := newYCC(image.Rect(0, 0, int(width), int(height)), input.SubsampleRatio)

		coeffs, offset, filterLength := createWeights8(temp.Bounds().Dy(), taps, blur, scaleX, kernel)
		in := imageYCbCrToYCC(input)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(temp, i, cpus).(*ycc)
			go func() {
				defer wg.Done()
				resizeYCbCr(in, slice, scaleX, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()

		coeffs, offset, filterLength = createWeights8(result.Bounds().Dy(), taps, blur, scaleY, kernel)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(result, i, cpus).(*ycc)
			go func() {
				defer wg.Done()
				resizeYCbCr(temp, slice, scaleY, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()
		return result.YCbCr()
	case *image.RGBA64:
		// 16-bit precision
		temp := image.NewRGBA64(image.Rect(0, 0, input.Bounds().Dy(), int(width)))
		result := image.NewRGBA64(image.Rect(0, 0, int(width), int(height)))

		// horizontal filter, results in transposed temporary image
		coeffs, offset, filterLength := createWeights16(temp.Bounds().Dy(), taps, blur, scaleX, kernel)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(temp, i, cpus).(*image.RGBA64)
			go func() {
				defer wg.Done()
				resizeRGBA64(input, slice, scaleX, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()

		// horizontal filter on transposed image, result is not transposed
		coeffs, offset, filterLength = createWeights16(result.Bounds().Dy(), taps, blur, scaleY, kernel)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(result, i, cpus).(*image.RGBA64)
			go func() {
				defer wg.Done()
				resizeGeneric(temp, slice, scaleY, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()
		return result
	case *image.Gray:
		// 8-bit precision
		temp := image.NewGray(image.Rect(0, 0, input.Bounds().Dy(), int(width)))
		result := image.NewGray(image.Rect(0, 0, int(width), int(height)))

		// horizontal filter, results in transposed temporary image
		coeffs, offset, filterLength := createWeights8(temp.Bounds().Dy(), taps, blur, scaleX, kernel)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(temp, i, cpus).(*image.Gray)
			go func() {
				defer wg.Done()
				resizeGray(input, slice, scaleX, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()

		// horizontal filter on transposed image, result is not transposed
		coeffs, offset, filterLength = createWeights8(result.Bounds().Dy(), taps, blur, scaleY, kernel)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(result, i, cpus).(*image.Gray)
			go func() {
				defer wg.Done()
				resizeGray(temp, slice, scaleY, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()
		return result
	case *image.Gray16:
		// 16-bit precision
		temp := image.NewGray16(image.Rect(0, 0, input.Bounds().Dy(), int(width)))
		result := image.NewGray16(image.Rect(0, 0, int(width), int(height)))

		// horizontal filter, results in transposed temporary image
		coeffs, offset, filterLength := createWeights16(temp.Bounds().Dy(), taps, blur, scaleX, kernel)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(temp, i, cpus).(*image.Gray16)
			go func() {
				defer wg.Done()
				resizeGray16(input, slice, scaleX, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()

		// horizontal filter on transposed image, result is not transposed
		coeffs, offset, filterLength = createWeights16(result.Bounds().Dy(), taps, blur, scaleY, kernel)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(result, i, cpus).(*image.Gray16)
			go func() {
				defer wg.Done()
				resizeGray16(temp, slice, scaleY, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()
		return result
	default:
		// 16-bit precision
		temp := image.NewRGBA64(image.Rect(0, 0, img.Bounds().Dy(), int(width)))
		result := image.NewRGBA64(image.Rect(0, 0, int(width), int(height)))

		// horizontal filter, results in transposed temporary image
		coeffs, offset, filterLength := createWeights16(temp.Bounds().Dy(), taps, blur, scaleX, kernel)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(temp, i, cpus).(*image.RGBA64)
			go func() {
				defer wg.Done()
				resizeGeneric(img, slice, scaleX, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()

		// horizontal filter on transposed image, result is not transposed
		coeffs, offset, filterLength = createWeights16(result.Bounds().Dy(), taps, blur, scaleY, kernel)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(result, i, cpus).(*image.RGBA64)
			go func() {
				defer wg.Done()
				resizeRGBA64(temp, slice, scaleY, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()
		return result
	}
}

var ResizeImage = TKX.ResizeImage

// ResizeImageQuick set width or height < 1, will keep ratio, otherwise as new width and height
func (pA *TK) ResizeImageQuick(imgA image.Image, newWidthA int, newHeightA int) image.Image {
	srcBounds := imgA.Bounds()
	srcWidth := srcBounds.Dx()
	srcHeight := srcBounds.Dy()

	// Calculate the aspect ratio of the source image
	srcAspectRatio := float64(srcWidth) / float64(srcHeight)

	if newWidthA < 1 && newHeightA < 1 {
		newWidthA = srcWidth
		newHeightA = srcHeight
	} else if newWidthA < 1 {
		newWidthA = int(float64(newHeightA) * srcAspectRatio)
	} else if newHeightA < 1 {
		newHeightA = int(float64(newWidthA) / srcAspectRatio)
	} else {
		imgT := image.NewRGBA(image.Rect(0, 0, newWidthA, newHeightA))

		draw.CatmullRom.Scale(imgT, imgT.Rect, imgA, srcBounds, draw.Over, nil)

		return imgT
	}

	// Calculate the target width and height while maintaining the aspect ratio
	var targetWidth, targetHeight int
	if float64(newWidthA)/float64(newHeightA) > srcAspectRatio {
		targetWidth = int(float64(newHeightA) * srcAspectRatio)
		targetHeight = newHeightA
	} else {
		targetWidth = newWidthA
		targetHeight = int(float64(newWidthA) / srcAspectRatio)
	}

	imgT := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))

	draw.CatmullRom.Scale(imgT, imgT.Rect, imgA, srcBounds, draw.Over, nil)

	return imgT
}

var ResizeImageQuick = TKX.ResizeImageQuick

// this function is highly inspired from github.com/1-ashraful-islam/image2ascii, thanks
func (pA *TK) ImageToAscii(imageA image.Image, widthA int, heightA int) interface{} {

	// Resize the image to half the height as we're going to use half-blocks to represent two vertical pixels.
	resizedImgT := ResizeImageQuick(imageA, widthA, heightA) // Adjust width and height as needed

	asciiImageT := make([]string, 0)

	for y := resizedImgT.Bounds().Min.Y; y < resizedImgT.Bounds().Max.Y; y += 2 {
		var line string
		for x := resizedImgT.Bounds().Min.X; x < resizedImgT.Bounds().Max.X; x++ {
			topPixelColor := resizedImgT.At(x, y)
			bottomPixelColor := color.RGBA{255, 255, 255, 255} // Default to white if bottom pixel exceeds image bounds
			if y+1 < resizedImgT.Bounds().Max.Y {
				bottomPixelColor, _ = resizedImgT.At(x, y+1).(color.RGBA)
			}

			// Extract RGBA components for top and bottom pixels
			r1, g1, b1, _ := topPixelColor.RGBA()
			r2, g2, b2, _ := bottomPixelColor.RGBA()

			// Print the half-block with top color as foreground and bottom color as background
			line += fmt.Sprintf("\x1b[38;2;%d;%d;%dm\x1b[48;2;%d;%d;%dm▀\x1b[0m", r1>>8, g1>>8, b1>>8, r2>>8, g2>>8, b2>>8)
		}

		asciiImageT = append(asciiImageT, line)
	}

	return asciiImageT
}

var ImageToAscii = TKX.ImageToAscii

func (pA *TK) ThumbImage(imgA image.Image, maxWidthA int, maxHeightA int) image.Image {
	if imgA == nil {
		return nil
	}

	maxWidth := uint(maxWidthA)
	maxHeight := uint(maxHeightA)

	origBounds := imgA.Bounds()
	origWidth := uint(origBounds.Dx())
	origHeight := uint(origBounds.Dy())
	newWidth, newHeight := origWidth, origHeight

	// Return original image if it have same or smaller size as constraints
	if maxWidth >= origWidth && maxHeight >= origHeight {
		return imgA
	}

	// Preserve aspect ratio
	if origWidth > maxWidth {
		newHeight = uint(origHeight * maxWidth / origWidth)
		if newHeight < 1 {
			newHeight = 1
		}
		newWidth = maxWidth
	}

	if newHeight > maxHeight {
		newWidth = uint(newWidth * maxHeight / newHeight)
		if newWidth < 1 {
			newWidth = 1
		}
		newHeight = maxHeight
	}

	return ResizeImage(int(newWidth), int(newHeight), imgA, MitchellNetravali)

}

var ThumbImage = TKX.ThumbImage

func (pA *TK) NewImage(argsA ...string) image.Image {
	widthT := ToInt(GetSwitch(argsA, "-width=", "100"), 100)
	heightT := ToInt(GetSwitch(argsA, "-height=", "100"), 100)

	colorStrT := GetSwitch(argsA, "-color=", "#000000")

	if colorStrT == "" {
		colorStrT = "#000000"
	}

	imgT := image.NewRGBA(image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: widthT, Y: heightT}})

	r, g, b, a := ParseHexColor(colorStrT)

	colorT := color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)}

	for x := 0; x < widthT; x++ {
		for y := 0; y < heightT; y++ {
			imgT.Set(x, y, colorT)
		}
	}

	return imgT
}

var NewImage = TKX.NewImage

// base on github.com/gezhengbin888/similarity_compare

func (pA *TK) CalTextSimilarity(str1A string, str2A string) float64 {
	if str1A == str2A {
		return 1.0
	}

	return (calcScoreMain(str1A, str2A) + calcScoreMain(str1A, str2A)) / 2
}

var CalTextSimilarity = TKX.CalTextSimilarity

func calcScoreMain(s_new string, s_old string) float64 {
	newLen := utf8.RuneCountInString(s_new)
	oldLen := utf8.RuneCountInString(s_old)
	newSLen := len(s_new)
	oldSLen := len(s_old)
	score := 0.0
	sOneNew := ""
	for i := 0; i < newLen; i++ {
		prefix := []rune(s_new)[i : i+1]
		// 将子串之前的字符串转换成[]rune
		sOneNew += string(prefix)
		strIndex := strings.Index(s_old, sOneNew)
		if strIndex == -1 {
			if utf8.RuneCountInString(sOneNew) >= 2 { //没找到且是连续2个字符以上的 截取前n-1个字符计算

				prefix1 := []rune(sOneNew)[0 : utf8.RuneCountInString(sOneNew)-1]
				// 将子串之前的字符串转换成[]rune
				contailStr := string(prefix1)
				score = calcScore(contailStr, score, strings.LastIndex(s_old, contailStr), strings.LastIndex(s_new, contailStr), newLen, oldLen, newSLen, oldSLen)
				i = i - 1
			}
			sOneNew = ""
		} else {
			if i == (newLen - 1) { //找到且是最后一个字符  直接拿来计算
				//fmt.Println(s_old,s_new, sOneNew,strings.Index(s_old,sOneNew), strings.Index(s_new,sOneNew) )
				score = calcScore(sOneNew, score, strings.LastIndex(s_old, sOneNew), strings.LastIndex(s_new, sOneNew), newLen, oldLen, newSLen, oldSLen)
			}
		}
	}
	return score
}

func calcScore(contailStr string, score float64, strIndexOld int, strIndexNew int, newLen int, oldLen int, newSLen int, oldSLen int) float64 {
	contailStrLen := utf8.RuneCountInString(contailStr)
	/*积分增加项= 长度匹配量 * 长度占比 越长越接近1*/
	asd := float64(contailStrLen) / (float64((newLen + oldLen)) / 2) * (float64(contailStrLen) / float64(contailStrLen+1))
	//fmt.Println(utils.Typeof(asd))
	//.Println(asd)
	score += asd
	pian := math.Abs(float64(strIndexOld - strIndexNew))
	if pian != 0 {
		xi := pian / float64((newSLen+oldSLen)/2)
		asff := (asd * xi) / 6 //宽松度 越大越不敏感 越大分数越大  越不严格
		if asff < asd {
			score -= asff
		}
	}
	return score
}

// based on github.com/dzstudio/similar-text

// return the len of longest string both in str1 and str2 and the positions in str1 and str2
func SimilarStr(str1 []rune, str2 []rune) (int, int, int) {
	var maxLen, tmp, pos1, pos2 = 0, 0, 0, 0
	len1, len2 := len(str1), len(str2)

	for p := 0; p < len1; p++ {
		for q := 0; q < len2; q++ {
			tmp = 0
			for p+tmp < len1 && q+tmp < len2 && str1[p+tmp] == str2[q+tmp] {
				tmp++
			}
			if tmp > maxLen {
				maxLen, pos1, pos2 = tmp, p, q
			}
		}

	}

	return maxLen, pos1, pos2
}

// return the total length of longest string both in str1 and str2
func SimilarChar(str1 []rune, str2 []rune) int {
	maxLen, pos1, pos2 := SimilarStr(str1, str2)
	total := maxLen

	if maxLen != 0 {
		if pos1 > 0 && pos2 > 0 {
			total += SimilarChar(str1[:pos1], str2[:pos2])
		}
		if pos1+maxLen < len(str1) && pos2+maxLen < len(str2) {
			total += SimilarChar(str1[pos1+maxLen:], str2[pos2+maxLen:])
		}
	}

	return total
}

// return a int value in [0, 100], which stands for match level
func SimilarText(str1 string, str2 string) int {
	txt1, txt2 := []rune(str1), []rune(str2)
	if len(txt1) == 0 || len(txt2) == 0 {
		return 0
	}
	return SimilarChar(txt1, txt2) * 200 / (len(txt1) + len(txt2))
}

func (pA *TK) GetTextSimilarity(str1A string, str2A string, methodA ...int) float64 {
	methodT := 0

	if len(methodA) > 0 {
		methodT = methodA[0]
	}

	if methodT == 1 {
		return float64(SimilarText(str1A, str2A)) / 100.0
	}

	return CalTextSimilarity(str1A, str2A)
}

var GetTextSimilarity = TKX.GetTextSimilarity

// otp related

func (pA *TK) GenerateOtpCode(secretA string, optsA ...string) string {
	if len(optsA) < 1 {
		s1, errT := totp.GenerateCode(secretA, time.Now())

		if errT != nil {
			return ErrorToString(errT)
		}

		return s1
	}

	timeT := ToTime(optsA[0])

	if IsErrX(timeT) {
		return "TXERROR:" + GetErrStrX(timeT)
	}

	s1, errT := totp.GenerateCode(secretA, timeT.(time.Time))

	if errT != nil {
		return ErrorToString(errT)
	}

	return s1
}

var GenerateOtpCode = TKX.GenerateOtpCode

func (pA *TK) ValidateOtpCode(passCodeA string, secretA string, optsA ...string) bool {
	return totp.Validate(passCodeA, secretA)
}

var ValidateOtpCode = TKX.ValidateOtpCode

// memory file, conform to os.File
type MemFile struct {
	m sync.Mutex
	b []byte
	i int
}

// New creates and initializes a new MemFile using b as its initial contents.
// The new File takes ownership of b.
func NewMemFile(b []byte) *MemFile {
	return &MemFile{b: b}
}

// Read reads up to len(b) bytes from the MemFile.
// It returns the number of bytes read and any error encountered.
// At end of file, Read returns (0, io.EOF).
func (fb *MemFile) Read(b []byte) (int, error) {
	fb.m.Lock()
	defer fb.m.Unlock()

	n, err := fb.readAt(b, int64(fb.i))
	fb.i += n
	return n, err
}

func (fb *MemFile) ReadAll() ([]byte, error) {
	fb.m.Lock()
	defer fb.m.Unlock()

	return fb.b, nil
}

// ReadAt reads len(b) bytes from the MemFile starting at byte offset.
// It returns the number of bytes read and the error, if any.
// At end of file, that error is io.EOF.
func (fb *MemFile) ReadAt(b []byte, offset int64) (int, error) {
	fb.m.Lock()
	defer fb.m.Unlock()
	return fb.readAt(b, offset)
}
func (fb *MemFile) readAt(b []byte, off int64) (int, error) {
	if off < 0 || int64(int(off)) < off {
		return 0, errors.New("invalid argument")
	}
	if off > int64(len(fb.b)) {
		return 0, io.EOF
	}
	n := copy(b, fb.b[off:])
	if n < len(b) {
		return n, io.EOF
	}
	return n, nil
}

// Write writes len(b) bytes to the MemFile.
// It returns the number of bytes written and an error, if any.
// If the current file offset is past the io.EOF, then the space in-between are
// implicitly filled with zero bytes.
func (fb *MemFile) Write(b []byte) (int, error) {
	fb.m.Lock()
	defer fb.m.Unlock()

	n, err := fb.writeAt(b, int64(fb.i))
	fb.i += n
	return n, err
}

// WriteAt writes len(b) bytes to the MemFile starting at byte offset.
// It returns the number of bytes written and an error, if any.
// If offset lies past io.EOF, then the space in-between are implicitly filled
// with zero bytes.
func (fb *MemFile) WriteAt(b []byte, offset int64) (int, error) {
	fb.m.Lock()
	defer fb.m.Unlock()
	return fb.writeAt(b, offset)
}
func (fb *MemFile) writeAt(b []byte, off int64) (int, error) {
	if off < 0 || int64(int(off)) < off {
		return 0, errors.New("invalid argument")
	}
	if off > int64(len(fb.b)) {
		fb.truncate(off)
	}
	n := copy(fb.b[off:], b)
	fb.b = append(fb.b, b[n:]...)
	return len(b), nil
}

// Seek sets the offset for the next Read or Write on file with offset,
// interpreted according to whence: 0 means relative to the origin of the file,
// 1 means relative to the current offset, and 2 means relative to the end.
func (fb *MemFile) Seek(offset int64, whence int) (int64, error) {
	fb.m.Lock()
	defer fb.m.Unlock()

	var abs int64
	switch whence {
	case io.SeekStart:
		abs = offset
	case io.SeekCurrent:
		abs = int64(fb.i) + offset
	case io.SeekEnd:
		abs = int64(len(fb.b)) + offset
	default:
		return 0, errors.New("invalid argument")
	}
	if abs < 0 {
		return 0, errors.New("invalid argument")
	}
	fb.i = int(abs)
	return abs, nil
}

// Truncate changes the size of the file. It does not change the I/O offset.
func (fb *MemFile) Truncate(n int64) error {
	fb.m.Lock()
	defer fb.m.Unlock()
	return fb.truncate(n)
}
func (fb *MemFile) truncate(n int64) error {
	switch {
	case n < 0 || int64(int(n)) < n:
		return errors.New("invalid argument")
	case n <= int64(len(fb.b)):
		fb.b = fb.b[:n]
		return nil
	default:
		fb.b = append(fb.b, make([]byte, int(n)-len(fb.b))...)
		return nil
	}
}

func (fb *MemFile) Close() error {
	fb.m.Lock()
	defer fb.m.Unlock()

	fb.b = nil

	return nil
}

// Bytes returns the full contents of the MemFile.
// The result in only valid until the next Write, WriteAt, or Truncate call.
func (fb *MemFile) Bytes() []byte {
	fb.m.Lock()
	defer fb.m.Unlock()
	return fb.b
}

// HTML template

func (pA *TK) ProcessHtmlTemplate(templateNameA string, baseDirA string, dataA interface{}, optsA ...string) interface{} {
	tmplFile := filepath.Join(baseDirA, templateNameA)

	if !IfFileExists(tmplFile) {
		return fmt.Errorf("file not exists: %v", tmplFile)
	}

	tmpl, errT := template.ParseFiles(tmplFile)

	if errT != nil {
		return errT
	}

	var sb = new(strings.Builder)

	errT = tmpl.ExecuteTemplate(sb, templateNameA, dataA)

	if errT != nil {
		return errT
	}

	return sb.String()
}

var ProcessHtmlTemplate = TKX.ProcessHtmlTemplate

// big float
func (pA *TK) BigFloatPower(f1 *big.Float, f2 *big.Float) *big.Float {
	return bigfloat.Pow(f1, f2)
}

var BigFloatPower = TKX.BigFloatPower

func (pA *TK) BigFloatExp(f1 *big.Float) *big.Float {
	return bigfloat.Exp(f1)
}

var BigFloatExp = TKX.BigFloatExp

func (pA *TK) BigFloatLog(f1 *big.Float) *big.Float {
	return bigfloat.Log(f1)
}

var BigFloatLog = TKX.BigFloatLog

func (pA *TK) DealString(strA string, optsA ...string) string {
	if strings.HasPrefix(strA, "HEX_") {
		strA = strA[4:]

		buf, errT := hex.DecodeString(strA)
		if errT != nil {
			return ErrStrf("failed decode hex: %v", errT)
		}

		return string(buf)
	} else if strings.HasPrefix(strA, "//TXHEX#") {
		strA = strA[8:]

		buf, errT := hex.DecodeString(strA)
		if errT != nil {
			return ErrStrf("failed decode hex: %v", errT)
		}

		return string(buf)
	} else if strings.HasPrefix(strA, "//TXTE#") {
		codeT := ""

		if len(optsA) > 0 {
			codeT = optsA[0]
		}

		return DecryptStringByTXTE(strA[7:], codeT)
	} else if strings.HasPrefix(strA, "//TXDEF#") || strings.HasPrefix(strA, "740404") {
		codeT := ""

		if len(optsA) > 0 {
			codeT = optsA[0]
		}

		return DecryptStringByTXDEF(strA, codeT)
	} else if strings.HasPrefix(strA, "//TXRR#") {
		strA = strA[7:]

		codeT := "char"

		if len(optsA) > 0 {
			codeT = optsA[0]
		}

		if strings.HasPrefix(strA, "//TXDEF#") {
			strA = DecryptStringByTXDEF(strings.TrimSpace(strA[8:]), codeT)
		}

		strA = strings.ReplaceAll(strA, "TX_secureCode_XT", EncryptStringByTXDEF(codeT, "char"))

		if strings.HasPrefix(strA, "http") {
			strA = ToStr(GetWeb(strings.TrimSpace(strA)))
		}

		return strA
	}

	return strA
}

var DealString = TKX.DealString

// steganography related

// modified from github.com/auyer/steganography

// EncodeNRGBA encodes a given string into the input image using least significant bit encryption (LSB steganography)
// The minnimum image size is 24 pixels for one byte. For each additional byte, it is necessary 3 more pixels.
/*
	Input:
		writeBuffer *bytes.Buffer : the destination of the encoded image bytes
		pictureInputFile image.NRGBA : image data used in encoding
		message []byte : byte slice of the message to be encoded
	Output:
		bytes buffer ( io.writter ) to create file, or send data.
*/
func EncodeNRGBA(writeBuffer *bytes.Buffer, rgbImage *image.NRGBA, message []byte) error {

	var messageLength = uint32(len(message))

	var width = rgbImage.Bounds().Dx()
	var height = rgbImage.Bounds().Dy()
	var c color.NRGBA
	var bit byte
	var ok bool
	//var encodedImage image.Image
	if MaxEncodeSize(rgbImage) < messageLength+4 {
		return errors.New("message too large for image")
	}

	one, two, three, four := splitToBytes(messageLength)

	message = append([]byte{four}, message...)
	message = append([]byte{three}, message...)
	message = append([]byte{two}, message...)
	message = append([]byte{one}, message...)

	ch := make(chan byte, 100)

	go getNextBitFromString(message, ch)

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {

			c = rgbImage.NRGBAAt(x, y) // get the color at this pixel

			/*  RED  */
			bit, ok = <-ch
			if !ok { // if we don't have any more bits left in our message

				rgbImage.SetNRGBA(x, y, c)
				break
			}
			setLSB(&c.R, bit)

			/*  GREEN  */
			bit, ok = <-ch
			if !ok {
				rgbImage.SetNRGBA(x, y, c)
				break
			}
			setLSB(&c.G, bit)

			/*  BLUE  */
			bit, ok = <-ch
			if !ok {
				rgbImage.SetNRGBA(x, y, c)
				break
			}
			setLSB(&c.B, bit)

			rgbImage.SetNRGBA(x, y, c)
		}
	}

	err := png.Encode(writeBuffer, rgbImage)
	return err
}

// Encode encodes a given string into the input image using least significant bit encryption (LSB steganography)
// The minnimum image size is 23 pixels
// It wraps EncodeNRGBA making the conversion from image.Image to image.NRGBA
/*
	Input:
		writeBuffer *bytes.Buffer : the destination of the encoded image bytes
		message []byte : byte slice of the message to be encoded
		pictureInputFile image.Image : image data used in encoding
	Output:
		bytes buffer ( io.writter ) to create file, or send data.
*/
func EncodeBytesInImageBuffer(bufA *bytes.Buffer, imageA *image.Image, contentA []byte, optsA ...string) error {
	if imageA == nil {
		imageT := NewImage(optsA...)

		rgbImageT := imageToNRGBA(imageT)

		return EncodeNRGBA(bufA, rgbImageT, contentA)
	}

	rgbImageT := imageToNRGBA(*imageA)

	return EncodeNRGBA(bufA, rgbImageT, contentA)
}

func (p *TK) EncodeBytesInImage(contentA []byte, imageA image.Image, optsA ...string) interface{} {
	bufT := new(bytes.Buffer)
	rgbImageT := imageToNRGBA(imageA)

	errT := EncodeNRGBA(bufT, rgbImageT, contentA)

	if errT != nil {
		return errT
	}

	return bufT.Bytes()
}

var EncodeBytesInImage = TKX.EncodeBytesInImage

func (p *TK) DecodeBytesFromImage(imageA image.Image, optsA ...string) []byte {
	sizeOfMessageT := GetMessageSizeFromImage(imageA)

	bufT := Decode(sizeOfMessageT, imageA)

	return bufT
}

var DecodeBytesFromImage = TKX.DecodeBytesFromImage

// decodeNRGBA gets messages from pictures using LSB steganography, decode the message from the picture and return it as a sequence of bytes
/*
	Input:
		startOffset uint32 : number of bytes used to declare size of message
		msgLen uint32 : size of the message to be decoded
		pictureInputFile image.NRGBA : image data used in decoding
	Output:
		message []byte decoded from image
*/
func decodeNRGBA(startOffset uint32, msgLen uint32, rgbImage *image.NRGBA) (message []byte) {

	var byteIndex uint32
	var bitIndex uint32

	width := rgbImage.Bounds().Dx()
	height := rgbImage.Bounds().Dy()

	var c color.NRGBA
	var lsb byte

	message = append(message, 0)

	// iterate through every pixel in the image and stitch together the message bit by bit
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {

			c = rgbImage.NRGBAAt(x, y) // get the color of the pixel

			/*  RED  */
			lsb = getLSB(c.R)                                                    // get the least significant bit from the red component of this pixel
			message[byteIndex] = setBitInByte(message[byteIndex], bitIndex, lsb) // add this bit to the message
			bitIndex++

			if bitIndex > 7 { // when we have filled up a byte, move on to the next byte
				bitIndex = 0
				byteIndex++

				if byteIndex >= msgLen+startOffset {
					return message[startOffset : msgLen+startOffset]
				}

				message = append(message, 0)
			}

			/*  GREEN  */
			lsb = getLSB(c.G)
			message[byteIndex] = setBitInByte(message[byteIndex], bitIndex, lsb)
			bitIndex++

			if bitIndex > 7 {

				bitIndex = 0
				byteIndex++

				if byteIndex >= msgLen+startOffset {
					return message[startOffset : msgLen+startOffset]
				}

				message = append(message, 0)
			}

			/*  BLUE  */
			lsb = getLSB(c.B)
			message[byteIndex] = setBitInByte(message[byteIndex], bitIndex, lsb)
			bitIndex++

			if bitIndex > 7 {
				bitIndex = 0
				byteIndex++

				if byteIndex >= msgLen+startOffset {
					return message[startOffset : msgLen+startOffset]
				}

				message = append(message, 0)
			}
		}
	}
	return
}

// decode gets messages from pictures using LSB steganography, decode the message from the picture and return it as a sequence of bytes
// It wraps EncodeNRGBA making the conversion from image.Image to image.NRGBA
/*
	Input:
		startOffset uint32 : number of bytes used to declare size of message
		msgLen uint32 : size of the message to be decoded
		pictureInputFile image.Image : image data used in decoding
	Output:
		message []byte decoded from image
*/
func decode(startOffset uint32, msgLen uint32, pictureInputFile image.Image) (message []byte) {

	rgbImage := imageToNRGBA(pictureInputFile)
	return decodeNRGBA(startOffset, msgLen, rgbImage)

}

// Decode gets messages from pictures using LSB steganography, decode the message from the picture and return it as a sequence of bytes
// It wraps EncodeNRGBA making the conversion from image.Image to image.NRGBA
/*
	Input:
		msgLen uint32 : size of the message to be decoded
		pictureInputFile image.Image : image data used in decoding
	Output:
		message []byte decoded from image
*/
func Decode(msgLen uint32, pictureInputFile image.Image) (message []byte) {
	return decode(4, msgLen, pictureInputFile) // the offset of 4 skips the "header" where message length is defined

}

// MaxEncodeSize given an image will find how many bytes can be stored in that image using least significant bit encoding
// ((width * height * 3) / 8 ) - 4
// The result must be at least 4,
func MaxEncodeSize(img image.Image) uint32 {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	eval := ((width * height * 3) / 8) - 4
	if eval < 4 {
		eval = 0
	}
	return uint32(eval)
}

// GetMessageSizeFromImage gets the size of the message from the first four bytes encoded in the image
func GetMessageSizeFromImage(pictureInputFile image.Image) (size uint32) {

	sizeAsByteArray := decode(0, 4, pictureInputFile)
	size = combineToInt(sizeAsByteArray[0], sizeAsByteArray[1], sizeAsByteArray[2], sizeAsByteArray[3])
	return
}

// getNextBitFromString each call will return the next subsequent bit in the string
func getNextBitFromString(byteArray []byte, ch chan byte) {

	var offsetInBytes int
	var offsetInBitsIntoByte int
	var choiceByte byte

	lenOfString := len(byteArray)

	for {
		if offsetInBytes >= lenOfString {
			close(ch)
			return
		}

		choiceByte = byteArray[offsetInBytes]
		ch <- getBitFromByte(choiceByte, offsetInBitsIntoByte)

		offsetInBitsIntoByte++

		if offsetInBitsIntoByte >= 8 {
			offsetInBitsIntoByte = 0
			offsetInBytes++
		}
	}
}

// getLSB given a byte, will return the least significant bit of that byte
func getLSB(b byte) byte {
	if b%2 == 0 {
		return 0
	}
	return 1
}

// setLSB given a byte will set that byte's least significant bit to a given value (where true is 1 and false is 0)
func setLSB(b *byte, bit byte) {
	if bit == 1 {
		*b = *b | 1
	} else if bit == 0 {
		var mask byte = 0xFE
		*b = *b & mask
	}
}

// getBitFromByte given a bit will return a bit from that byte
func getBitFromByte(b byte, indexInByte int) byte {
	b = b << uint(indexInByte)
	var mask byte = 0x80

	var bit = mask & b

	if bit == 128 {
		return 1
	}
	return 0
}

// setBitInByte sets a specific bit in a byte to a given value and returns the new byte
func setBitInByte(b byte, indexInByte uint32, bit byte) byte {
	var mask byte = 0x80
	mask = mask >> uint(indexInByte)

	if bit == 0 {
		mask = ^mask
		b = b & mask
	} else if bit == 1 {
		b = b | mask
	}
	return b
}

// combineToInt given four bytes, will return the 32 bit unsigned integer which is the composition of those four bytes (one is MSB)
func combineToInt(one, two, three, four byte) (ret uint32) {
	ret = uint32(one)
	ret = ret << 8
	ret = ret | uint32(two)
	ret = ret << 8
	ret = ret | uint32(three)
	ret = ret << 8
	ret = ret | uint32(four)
	return
}

// splitToBytes given an unsigned integer, will split this integer into its four bytes
func splitToBytes(x uint32) (one, two, three, four byte) {
	one = byte(x >> 24)
	var mask uint32 = 255

	two = byte((x >> 16) & mask)
	three = byte((x >> 8) & mask)
	four = byte(x & mask)
	return
}

// imageToNRGBA converts image.Image to image.NRGBA
func imageToNRGBA(src image.Image) *image.NRGBA {
	bounds := src.Bounds()

	var m *image.NRGBA
	var width, height int

	width = bounds.Dx()
	height = bounds.Dy()

	m = image.NewNRGBA(image.Rect(0, 0, width, height))

	draw.Draw(m, m.Bounds(), src, bounds.Min, draw.Src)
	return m
}

// docx related

// GetXmlContent extracts the XML content from a .docx file.
func GetXmlContent(docxFile string) ([]byte, error) {
	r, err := zip.OpenReader(docxFile)
	if err != nil {
		return nil, fmt.Errorf("error opening docx file: %w", err)
	}
	defer r.Close()

	for _, f := range r.File {
		if f.Name == "word/document.xml" {
			rc, err := f.Open()
			if err != nil {
				return nil, fmt.Errorf("error opening document.xml: %w", err)
			}
			defer rc.Close()

			content, err := io.ReadAll(rc) // Changed from ioutil.ReadAll to io.ReadAll
			if err != nil {
				return nil, fmt.Errorf("error reading document.xml: %w", err)
			}

			return content, nil
		}
	}

	return nil, fmt.Errorf("document.xml not found")
}

// GetTextByParagraph parses XML content and extracts text by paragraphs.
func GetTextByParagraph(content []byte) ([]string, error) {
	decoder := xml.NewDecoder(bytes.NewReader(content))
	var resultStrings []string
	var builder strings.Builder 
	inParagraph := false

	for {
		t, err := decoder.Token()
		if t == nil {
			break
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("error decoding XML: %w", err)
		}

		switch se := t.(type) {
		case xml.StartElement:
//			Pl("<ELM START> %#v", se)
			if se.Name.Local == "p" {
				inParagraph = true
				builder.Reset() 
			}
		case xml.CharData:
//			Pl("<ELM CharData> %#v", se)
			if inParagraph {
				builder.WriteString(string(se)) 
			}
		case xml.EndElement:
//			Pl("<ELM END> %#v", se)
//			Pl("<P END> %#v", builder.String())
			if se.Name.Local == "p" { //  && builder.Len() > 0
				resultStrings = append(resultStrings, SplitLines(builder.String())...)
				inParagraph = false
			}
		default:
//			Pl("<ELM OTHER> %#v", se)
			break
		
		}
	}

	return resultStrings, nil
}

func (pA *TK) DocxToText(filePathA string, optsA ...string) interface{} {
	xmlContent, err := GetXmlContent(filePathA)
    if err != nil {
        return err
    }

    texts, err := GetTextByParagraph(xmlContent)
    if err != nil {
        return err
    }

    return texts
}

var DocxToText = TKX.DocxToText

func (pA *TK) DocxBytesToText(bytesA []byte, optsA ...string) interface{} {
	
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	readerT := bytes.NewReader(bytesA)

	fsys, err := archives.FileSystem(ctx, "a.docx", readerT)
	if err != nil {
		return fmt.Errorf("failed to open archive: %v", err)
	}
	
	f, err := fsys.Open("word/document.xml")
	if err != nil {
		return fmt.Errorf("failed to get document.xml: %v", err)
	}
	
	defer f.Close()
	
	xmlContent, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("failed to read document.xml: %v", err)
	}

    texts, err := GetTextByParagraph(xmlContent)
    if err != nil {
		return fmt.Errorf("failed to get text paragraph: %v", err)
    }

    return texts
}

var DocxBytesToText = TKX.DocxBytesToText

func (pA *TK) ReplacePatternsInDocxBytes(bytesA []byte, replacesA []string, optsA ...string) interface{} {
	doc, err := docxrepl.OpenBytes(bytesA)
	if err != nil {
		return err
	}
	
	replaceMap := make(map[string]interface{})
	
	lenT := len(replacesA) / 2
	
	for i := 0; i < lenT; i ++ {
		replaceMap[replacesA[i * 2]] = replacesA[i * 2 + 1]
	}
	
	openDelimiter := GetSwitch(optsA, "-startRune=", "{")
	closeDelimiter := GetSwitch(optsA, "-endRune=", "}")
	
	err = docxrepl.SetDelimiter(openDelimiter, closeDelimiter)
	if err != nil {
		return err
	}

	if IfSwitchExists(optsA, "-refreshOnly") {
		return bytesA
	}

	err = doc.ReplaceAll(replaceMap)
	if err != nil {
		if !strings.Contains(err.Error(), "not all placeholders were replaced") {
			return err
		}
	}

	var bufT bytes.Buffer

	err = doc.Write(&bufT)

	if err != nil {
		return err
	}

	return bufT.Bytes()
}

var ReplacePatternsInDocxBytes = TKX.ReplacePatternsInDocxBytes

func (pA *TK) GetPlaceholderInDocxBytes(bytesA []byte, optsA ...string) interface{} {
	doc, err := docxrepl.OpenBytes(bytesA)
	if err != nil {
		return err
	}
	
	return doc.PlaceholderTexts()
}

var GetPlaceholderInDocxBytes = TKX.GetPlaceholderInDocxBytes

func (pA *TK) GetFileListInZip(filePathA string, optsA ...string) interface{} {
	r, err := zip.OpenReader(filePathA)
	if err != nil {
		return err
	}
	
	defer r.Close()
	
	listT := make([]string, 0, len(r.File))
	
	encT := GetSwitch(optsA, "-encode=", "")

	for _, f := range r.File {
		if f.NonUTF8 || f.Flags == 0 || !utf8.ValidString(f.Name) {
			listT = append(listT, ConvertStringToUTF8(f.Name, encT))
		} else {
			listT = append(listT, f.Name)
		}
	}
	
	return listT
}

var GetFileListInZip = TKX.GetFileListInZip

func (pA *TK) GetFileListInZipBytes(bytesA []byte, optsA ...string) interface{} {
	readerT := bytes.NewReader([]byte(bytesA))
	
	if readerT == nil {
		return fmt.Errorf("failed to init reader")
	}
	
	r, err := zip.NewReader(readerT, int64(len(bytesA)))
	if err != nil {
		return err
	}
	
	listT := make([]string, 0, len(r.File))
	
	encT := GetSwitch(optsA, "-encode=", "")

	for _, f := range r.File {
		if f.NonUTF8 || f.Flags == 0 || !utf8.ValidString(f.Name) {
			listT = append(listT, ConvertStringToUTF8(f.Name, encT))
		} else {
			listT = append(listT, f.Name)
		}
	}
	
	return listT
}

var GetFileListInZipBytes = TKX.GetFileListInZipBytes

func (pA *TK) IsFileNameUtf8InZipBytes(bytesA []byte, optsA ...string) interface{} {
	readerT := bytes.NewReader([]byte(bytesA))
	
	if readerT == nil {
		return fmt.Errorf("failed to init reader")
	}
	
	r, err := zip.NewReader(readerT, int64(len(bytesA)))
	if err != nil {
		return err
	}
	
	for _, f := range r.File {
		if f.NonUTF8 || f.Flags == 0 || !utf8.ValidString(f.Name) {
			return false
		} else {
			return true
		}
	}
	
	return true
}

var IsFileNameUtf8InZipBytes = TKX.IsFileNameUtf8InZipBytes

func (pA *TK) GetFileListInArchive(filePathA string, optsA ...string) interface{} {

//	if strings.HasSuffix(strings.ToLower(filePathA), ".zip") {
//		return GetFileListInZip(filePathA, optsA...)
//	}
	
//	if formatT == "zip" {
//		return GetFileListInZip(filePathA, optsA...)
//	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	formatT := GetSwitch(optsA, "-format=", "")

	var fsys fs.FS
	var err error
	
//	extractor = archives.Zip{TextEncoding: simplifiedchinese.GB18030}

	if strings.HasSuffix(strings.ToLower(filePathA), "zip") || formatT == "zip" {
		fsys = &archives.ArchiveFS{
			Path: filePathA,
			Format: archives.Zip{TextEncoding: simplifiedchinese.GB18030},
		}
	} else {
		if formatT != "" {
			fsys, err = archives.FileSystem(ctx, filePathA, nil)
			if err != nil {
				return fmt.Errorf("failed to open archive: %v", err)
			}
		} else {
			format, _, identifyErr := archives.Identify(context.Background(), filePathA, nil)
			if identifyErr != nil {
				return fmt.Errorf("identify format: %w", identifyErr)
			}
			
			if strings.HasSuffix(format.Extension(), "zip") {
				fsys = &archives.ArchiveFS{
					Path: filePathA,
					Format: archives.Zip{TextEncoding: simplifiedchinese.GB18030},
				}
			} else {
				fsys, err = archives.FileSystem(ctx, filePathA, nil)
				if err != nil {
					return fmt.Errorf("failed to open archive: %v", err)
				}
			}
		}

	}

//	fsys, err := archives.FileSystem(ctx, filePathA, nil)
//	if err != nil {
//		return fmt.Errorf("failed to open archive: %v", err)
//	}
//	f, err := fsys.Open(filePathA)
//	if err != nil {
//		return err
//	}
//	
//	defer f.Close()
	
	listT := make([]string, 0, 10)
	
//	if dir, ok := fsys.(fs.ReadDirFile); ok {
//		// 0 gets all entries, but you can pass > 0 to paginate
//		entries, err := dir.ReadDir(0)
//		if err != nil {
//			return err
//		}
//		
//		for _, e := range entries {
//			listT = append(listT, e.Name())
////			fmt.Println(e.Extension())
//		}
//	}
//	
//	return listT
	
	err1 := fs.WalkDir(fsys, ".", func(pathA string, deA fs.DirEntry, errA error) error {
		if errA != nil {
			return fmt.Errorf("failed to walk dir: %v", errA)
		}
		
//		if path == ".git" {
//			return fs.SkipDir
//		}

		if !deA.IsDir() {
			encT := GetSwitch(optsA, "-encode=", "")
			
			if !utf8.ValidString(pathA) {
				listT = append(listT, ConvertStringToUTF8(pathA, encT))
			} else {
				listT = append(listT, pathA)
			}
		}

		return nil
	})
	
	if err1 != nil {
		return fmt.Errorf("failed after walking dir: %v(%#v)", err1, listT)
	}

	return listT
}

var GetFileListInArchive = TKX.GetFileListInArchive

func (pA *TK) GetFileListInArchiveBytes(bytesA []byte, optsA ...string) interface{} {

	formatT := GetSwitch(optsA, "-format=", "")

//	if formatT == "zip" {
//		return GetFileListInZipBytes(bytesA, optsA...)
//	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	fakeFileNameT := ""
	
	if formatT != "" {
		fakeFileNameT = "a." + formatT
	}

	var fsys fs.FS
	var err error

	if formatT == "zip" {
		fsys = &archives.ArchiveFS{
			Stream: io.NewSectionReader(bytes.NewReader(bytesA), 0, int64(len(bytesA))),
			Format: archives.Zip{TextEncoding: simplifiedchinese.GB18030},
		}
	} else {
		if formatT != "" {
			fsys, err = archives.FileSystem(ctx, fakeFileNameT, bytes.NewReader(bytesA))
			if err != nil {
				return fmt.Errorf("failed to open archive: %v", err)
			}
		} else {
			format, _, identifyErr := archives.Identify(context.Background(), fakeFileNameT, bytes.NewReader(bytesA))
			if identifyErr != nil {
				return fmt.Errorf("identify format: %w", identifyErr)
			}
			
			if strings.HasSuffix(format.Extension(), "zip") {
				fsys = &archives.ArchiveFS{
					Stream: io.NewSectionReader(bytes.NewReader(bytesA), 0, int64(len(bytesA))),
					Format: archives.Zip{TextEncoding: simplifiedchinese.GB18030},
				}
			} else {
				fsys, err = archives.FileSystem(ctx, fakeFileNameT, bytes.NewReader(bytesA))
				if err != nil {
					return fmt.Errorf("failed to open archive: %v", err)
				}
			}
		}

	}

	
//	f, err := fsys.Open(filePathA)
//	if err != nil {
//		return err
//	}
//	
//	defer f.Close()
	
	listT := make([]string, 0, 10)
	
//	if dir, ok := fsys.(fs.ReadDirFile); ok {
//		// 0 gets all entries, but you can pass > 0 to paginate
//		entries, err := dir.ReadDir(0)
//		if err != nil {
//			return err
//		}
//		
//		for _, e := range entries {
//			listT = append(listT, e.Name())
////			fmt.Println(e.Extension())
//		}
//	}
//	
//	return listT
	
	err1 := fs.WalkDir(fsys, ".", func(pathA string, deA fs.DirEntry, errA error) error {
		if errA != nil {
			return fmt.Errorf("failed to walk dir: %v", errA)
		}
		
//		if path == ".git" {
//			return fs.SkipDir
//		}

		if !deA.IsDir() {
			encT := GetSwitch(optsA, "-encode=", "")
			
			if !utf8.ValidString(pathA) {
				listT = append(listT, ConvertStringToUTF8(pathA, encT))
			} else {
				listT = append(listT, pathA)
			}
		}

		return nil
	})
	
	if err1 != nil {
		return fmt.Errorf("failed after walking dir: %v(%#v)", err1, listT)
	}

	return listT
}

var GetFileListInArchiveBytes = TKX.GetFileListInArchiveBytes

func (pA *TK) LoadBytesInArchive(archivePathA string, filePathA string, optsA ...string) interface{} {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	formatT := GetSwitch(optsA, "-format=", "")

	var fsys fs.FS
	var err error
	
//	extractor = archives.Zip{TextEncoding: simplifiedchinese.GB18030}

	if strings.HasSuffix(strings.ToLower(archivePathA), "zip") || formatT == "zip" {
		fsys = &archives.ArchiveFS{
			Path: archivePathA,
			Format: archives.Zip{TextEncoding: simplifiedchinese.GB18030},
		}
	} else {
		if formatT != "" {
			fsys, err = archives.FileSystem(ctx, archivePathA, nil)
			if err != nil {
				return fmt.Errorf("failed to open archive: %v", err)
			}
		} else {
			format, _, identifyErr := archives.Identify(context.Background(), archivePathA, nil)
			if identifyErr != nil {
				return fmt.Errorf("identify format: %w", identifyErr)
			}
			
			if strings.HasSuffix(format.Extension(), "zip") {
				fsys = &archives.ArchiveFS{
					Path: archivePathA,
					Format: archives.Zip{TextEncoding: simplifiedchinese.GB18030},
				}
			} else {
				fsys, err = archives.FileSystem(ctx, archivePathA, nil)
				if err != nil {
					return fmt.Errorf("failed to open archive: %v", err)
				}
			}
		}

//		fsys, err = archives.FileSystem(ctx, fakeFileNameT, bytes.NewReader(bytesA))
//		if err != nil {
//			return fmt.Errorf("failed to open archive: %v", err)
//		}
	}
//	fsys, err := archives.FileSystem(ctx, archivePathA, nil)
//	if err != nil {
//		return fmt.Errorf("failed to open archive: %v", err)
//	}

	fileT, err := fsys.Open(filePathA)
	if err != nil {
		return err
	}
	
	defer fileT.Close()
	
	limitT := ToInt(GetSwitch(optsA, "-limit=", "0"), 0)
	
	if limitT <= 0 {
		fileContentT, err := io.ReadAll(fileT)
		if err != nil {
			return err
		}

		return fileContentT
	}

	bufT := make([]byte, limitT)
	nn, err := fileT.Read(bufT)
	if err != nil {
		fmt.Errorf("failed to read archive: %v", err)
	}
	
	if nn != len(bufT) {
		fmt.Errorf("failed to read archive: %v", err)
	}

	return bufT
}

var LoadBytesInArchive = TKX.LoadBytesInArchive

func (pA *TK) ExtractFileInArchive(archivePathA string, filePathA string, toPathA string, optsA ...string) interface{} {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	formatT := GetSwitch(optsA, "-format=", "")

	var fsys fs.FS
	var err error
	
//	extractor = archives.Zip{TextEncoding: simplifiedchinese.GB18030}

	if strings.HasSuffix(strings.ToLower(archivePathA), "zip") || formatT == "zip" {
		fsys = &archives.ArchiveFS{
			Path: archivePathA,
			Format: archives.Zip{TextEncoding: simplifiedchinese.GB18030},
		}
	} else {
		if formatT != "" {
			fsys, err = archives.FileSystem(ctx, archivePathA, nil)
			if err != nil {
				return fmt.Errorf("failed to open archive: %v", err)
			}
		} else {
			format, _, identifyErr := archives.Identify(context.Background(), archivePathA, nil)
			if identifyErr != nil {
				return fmt.Errorf("identify format: %w", identifyErr)
			}
			
			if strings.HasSuffix(format.Extension(), "zip") {
				fsys = &archives.ArchiveFS{
					Path: archivePathA,
					Format: archives.Zip{TextEncoding: simplifiedchinese.GB18030},
				}
			} else {
				fsys, err = archives.FileSystem(ctx, archivePathA, nil)
				if err != nil {
					return fmt.Errorf("failed to open archive: %v", err)
				}
			}
		}

//		fsys, err = archives.FileSystem(ctx, fakeFileNameT, bytes.NewReader(bytesA))
//		if err != nil {
//			return fmt.Errorf("failed to open archive: %v", err)
//		}
	}

//	fsys, err := archives.FileSystem(ctx, archivePathA, nil)
//	if err != nil {
//		return fmt.Errorf("failed to open archive: %v", err)
//	}

	fileT, err := fsys.Open(filePathA)
	if err != nil {
		return err
	}
	
	defer fileT.Close()

	forceT := IfSwitchExists(optsA, "-force")

	if !forceT {
		if IfFileExists(toPathA) {
			return fmt.Errorf("file %s already exists", toPathA)
		}
	}

	destination, err := os.Create(toPathA)
	if err != nil {
		return fmt.Errorf("failed to create dest file: %v", err)
	}

	defer destination.Close()

	bufferSizeT := ToInt(GetSwitch(optsA, "-bufferSize=", ""), 1000000)

	if bufferSizeT <= 0 {
		bufferSizeT = 1000000
	}

	buf := make([]byte, bufferSizeT)
	for {
		n, err := fileT.Read(buf)

		if err != nil && err != io.EOF {
			return err
		}

		if n == 0 {
			break
		}

		if _, err := destination.Write(buf[:n]); err != nil {
			return err
		}
	}

	return ""
}

var ExtractFileInArchive = TKX.ExtractFileInArchive

func (pA *TK) GetFileContentInArchiveBytes(bytesA []byte, filePathA string, optsA ...string) interface{} {

	formatT := GetSwitch(optsA, "-format=", "")

	fakeFileNameT := ""
	
	if formatT != "" {
		fakeFileNameT = "a." + formatT
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	var fsys fs.FS
	var err error
	
//	extractor = archives.Zip{TextEncoding: simplifiedchinese.GB18030}

	if formatT == "zip" {
		fsys = &archives.ArchiveFS{
			Stream: io.NewSectionReader(bytes.NewReader(bytesA), 0, int64(len(bytesA))),
			Format: archives.Zip{TextEncoding: simplifiedchinese.GB18030},
		}
	} else {
		if formatT != "" {
			fsys, err = archives.FileSystem(ctx, fakeFileNameT, bytes.NewReader(bytesA))
			if err != nil {
				return fmt.Errorf("failed to open archive: %v", err)
			}
		} else {
			format, _, identifyErr := archives.Identify(context.Background(), fakeFileNameT, bytes.NewReader(bytesA))
			if identifyErr != nil {
				return fmt.Errorf("identify format: %w", identifyErr)
			}
			
			if strings.HasSuffix(format.Extension(), "zip") {
				fsys = &archives.ArchiveFS{
					Stream: io.NewSectionReader(bytes.NewReader(bytesA), 0, int64(len(bytesA))),
					Format: archives.Zip{TextEncoding: simplifiedchinese.GB18030},
				}
//				filePathA = Utf8ToGb(filePathA).(string)
			} else {
				fsys, err = archives.FileSystem(ctx, fakeFileNameT, bytes.NewReader(bytesA))
				if err != nil {
					return fmt.Errorf("failed to open archive: %v", err)
				}
			}
		}

//		fsys, err = archives.FileSystem(ctx, fakeFileNameT, bytes.NewReader(bytesA))
//		if err != nil {
//			return fmt.Errorf("failed to open archive: %v", err)
//		}
	}

	fileT, err := fsys.Open(filePathA)
	if err != nil {
		return err
	}
	
	defer fileT.Close()

	fileContentT, err := io.ReadAll(fileT)
	if err != nil {
		return GenerateErrorStringTX(err.Error())
	}

	if err != nil && err != io.EOF {
		return err
	}

	return fileContentT
}

var GetFileContentInArchiveBytes = TKX.GetFileContentInArchiveBytes

func (pA *TK) ExtractArchive(archivePathA string, toPathA string, optsA ...string) interface{} {

	archiveFile, openErr := os.Open(archivePathA)
	if openErr != nil {
		return fmt.Errorf("open archive %s: %w", archivePathA, openErr)
	}
	defer archiveFile.Close()
	
	errT := EnsureMakeDirsE(toPathA)
	if errT != nil {
		return fmt.Errorf("failed to create dest dir: %v(%v)", errT, toPathA)
	}
	
	if !IfFileExists(toPathA) {
		return fmt.Errorf("failed to create dest dir: %v", toPathA)
	}

	format, input, identifyErr := archives.Identify(context.Background(), archivePathA, archiveFile)
	if identifyErr != nil {
		return fmt.Errorf("identify format: %w", identifyErr)
	}
	
//	Pl("format: %v, %v", format.Extension(), format.MediaType())
	
	var extractor archives.Extractor
	var ok bool
	
	if strings.HasSuffix(format.Extension(), "zip") {
//		fsys = &archives.ArchiveFS{
//			Stream: io.NewSectionReader(bytes.NewReader(bytesA), 0, int64(len(bytesA))),
//			Format: archives.Zip{},
//		}
		encT := GetSwitch(optsA, "-encode=", "")
		
		if encT == "utf8" || encT == "utf-8" {
			extractor, ok = format.(archives.Extractor)
			if !ok {
				return fmt.Errorf("unsupported format for extraction")
			}
		} else {
			extractor = archives.Zip{TextEncoding: simplifiedchinese.GB18030}
		}
		
	} else {
		extractor, ok = format.(archives.Extractor)
		if !ok {
			return fmt.Errorf("unsupported format for extraction")
		}
	}

	
	forceT := IfSwitchExists(optsA, "-force")
	
	noFileDirT := IfSwitchExists(optsA, "-noFileDir")

	err := extractor.Extract(context.Background(), input, func(ctx context.Context, f archives.FileInfo) error {
		// do something with the file here; or, if you only want a specific file or directory,
		// just return until you come across the desired f.NameInArchive value(s)
		
		var dstPath string
		
		if noFileDirT {
			dstPath = filepath.Join(toPathA, filepath.Base(f.NameInArchive))
		} else {
			dstPath = filepath.Join(toPathA, f.NameInArchive)
		}
		
		if f.IsDir() {
			errT := EnsureMakeDirsE(dstPath)
			if errT != nil {
				return fmt.Errorf("failed to make directory: %v", dstPath)
			}
			return nil
		}
		
		if !forceT {
			if IfFileExists(dstPath) {
				return fmt.Errorf("file already exists: %v", dstPath)
			}
		}
		
		if f.LinkTarget != "" {
			return nil
		}
		
		dirT := filepath.Dir(dstPath)
		errT := EnsureMakeDirsE(dirT)
		if errT != nil {
			return fmt.Errorf("failed to make directory: %v", dirT)
		}
		
		reader, openErr := f.Open()
		if openErr != nil {
			return fmt.Errorf("open file: %w", openErr)
		}
		defer reader.Close()

		dstFile, createErr := os.OpenFile(dstPath, os.O_CREATE|os.O_WRONLY, f.Mode())
		if createErr != nil {
			return fmt.Errorf("create file: %w", createErr)
		}
		defer dstFile.Close()

		if _, copyErr := io.Copy(dstFile, reader); copyErr != nil {
			return fmt.Errorf("copy: %w", copyErr)
		}
		
		return nil
	})
	
	if err != nil {
		return fmt.Errorf("failed to extract file: %w", err)
	}
	
	return ""
}

var ExtractArchive = TKX.ExtractArchive

func (pA *TK) IsUtf8Str(strA string) bool {
	return utf8.ValidString(strA)
}

var IsUtf8Str = TKX.IsUtf8Str

func (pA *TK) IsUtf8Bytes(bytesA []byte) bool {
	return utf8.Valid(bytesA)
}

var IsUtf8Bytes = TKX.IsUtf8Bytes

// multi-line input
type errMsg error

type TextAreaModel struct {
	Textarea textarea.Model
	err      error
	ReturnMode	string
	Title	string
	Bottom string
	ReturnStatus int // -1: exit, 1 done
	Quitting bool
	KeepOnQuit bool
	Delegate func(...interface{}) interface{}
}

func (m TextAreaModel) Init() bubbletea.Cmd {
	return textarea.Blink
}

func (m TextAreaModel) Update(msg bubbletea.Msg) (bubbletea.Model, bubbletea.Cmd) {
	var cmds []bubbletea.Cmd
	var cmd bubbletea.Cmd
	
	switch msg := msg.(type) {
	case bubbletea.KeyMsg:
//		fmt.Printf("msg got: %#v\n", msg)

		if m.Delegate != nil {
			
			rs := m.Delegate("key", int(msg.Type), msg.Runes, msg.Alt, msg.Paste, msg.String()).(string)
			
//			fmt.Printf("callback return: %v\n", rs)

			if strings.HasPrefix(rs, "TXERROR:") {
				bubbletea.Printf("input callback error: %v", rs[8:])
			}
			
			if rs == "clear" {
				m.Textarea.SetValue("")
				cmd = m.Textarea.Focus()
				cmds = append(cmds, cmd)
				return m, bubbletea.Batch(cmds...)
			} else if rs == "exit" {
				m.ReturnStatus = -1 // "exit"
				m.Quitting = true
				return m, bubbletea.Quit
			} else if rs == "return" {
				m.ReturnStatus = 1 // "done"
				m.Quitting = true
				return m, bubbletea.Quit
			} else if rs == "blur" {
				m.Textarea.Blur()
			}
		}


		switch msg.Type {
//		case bubbletea.KeyEsc:
//			if m.Textarea.Focused() {
//				m.Textarea.Blur()
//			}
		case bubbletea.KeyTab:
			if !m.Textarea.Focused() {
				cmd = m.Textarea.Focus()
				cmds = append(cmds, cmd)
			}
			m.Textarea, cmd = m.Textarea.Update(bubbletea.KeyMsg{Runes: []rune{9}})
			cmds = append(cmds, cmd)
			return m, bubbletea.Batch(cmds...)
		case bubbletea.KeyCtrlQ:
			if m.Delegate == nil {
				m.ReturnStatus = 1 // "done"
				m.Quitting = true
				return m, bubbletea.Quit
			}
		case bubbletea.KeyCtrlX:
			if m.Delegate == nil {
				m.ReturnStatus = -1 // "exit"
				m.Quitting = true
				return m, bubbletea.Quit
			}
		default:
			if !m.Textarea.Focused() {
				cmd = m.Textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.Textarea, cmd = m.Textarea.Update(msg)
	cmds = append(cmds, cmd)
	return m, bubbletea.Batch(cmds...)
}

func (m TextAreaModel) View() string {
	if m.Quitting {
		if !m.KeepOnQuit {
			return ""
		}
	}
	
	return fmt.Sprintf(
		"%v\n\n%s\n\n%s", 
		m.Title,
		m.Textarea.View(),
		m.Bottom,
	) + "\n\n"
}

// rs := GetMultiLineInput("-placeholder=", "-title=Please enter...", "-bottom=Press Ctrl-Q to finish, Ctrl-X to exit.", "-width=100", "-height=10")
func (pA *TK) GetMultiLineInput(deleA func(...interface{}) interface{}, optsA ...string) string {
	placeholderT := GetSwitch(optsA, "-placeholder=", "")
	returnModeT := GetSwitch(optsA, "-returnMode=", "")
	titleT := GetSwitch(optsA, "-title=", "")
	titleT = GetSwitch(optsA, "-prompt=", titleT)
	bottomT := GetSwitch(optsA, "-bottom=", "")
	widthT := ToInt(GetSwitch(optsA, "-width=", "-1"), -1)
	heightT := ToInt(GetSwitch(optsA, "-height=", "-1"), -1)
	ifShowLineNumbersT := GetSwitch(optsA, "-showLineNumbers=", "true") != "false"
	charLimitT := ToInt(GetSwitch(optsA, "-charLimit=", "0"), 0)
	maxWidthT := ToInt(GetSwitch(optsA, "-maxWidth=", "0"), 0)
	maxHeightT := ToInt(GetSwitch(optsA, "-maxHeight=", "0"), 0)
	keepT := IfSwitchExists(optsA, "-keep")
	textT := GetSwitch(optsA, "-text=", "")

	p := bubbletea.NewProgram(func () TextAreaModel {
		ti := textarea.New()
		ti.Placeholder = placeholderT

		if heightT >= 0 {
			ti.SetHeight(heightT)
		}
		
		if widthT >= 0 {
			ti.SetWidth(widthT)
		}
		
		ti.ShowLineNumbers = ifShowLineNumbersT
		ti.CharLimit = charLimitT
		ti.MaxWidth = maxWidthT
		ti.MaxHeight = maxHeightT
		
		ti.SetValue(textT)

		ti.Focus()

		return TextAreaModel{
			Textarea: ti,
			Title: titleT,
			Bottom: bottomT,
			err: nil,
			Delegate: deleA,
			KeepOnQuit: keepT,
			ReturnMode: returnModeT,
		}
	}())

	rs, err := p.Run()
	
	if err != nil {
		return "TXERROR:" + fmt.Sprintf("%v", err)
	}
	
//	fmt.Printf("%#v\n", rs)
	
	rs1 := rs.(TextAreaModel)
	
	if rs1.ReturnStatus < 0 {
		return "TXERROR:" + "cancel"
	}
	
//	fmt.Printf("%#v\n", rs1)
	
	return rs1.Textarea.Value()
}

var GetMultiLineInput = TKX.GetMultiLineInput

// excel column name
const (
	asciiAlphabetStart = 65
	alphabetCount      = 26
	colMin             = 1
)

func (pA *TK) ConvertExcelColIndexToAlphabet(col int) (string, error) {
	if col < colMin {
		return "", fmt.Errorf("argument is out of range [%d]", col)
	}
	
	var colName string
	tmp := col

	for tmp > 0 {
		index := tmp - colMin
		remaining := index / alphabetCount
		charIndex := int(math.Mod(float64(index), alphabetCount))
		colName = string(rune(charIndex+asciiAlphabetStart)) + colName
		tmp = remaining
	}

	return colName, nil
}

var ConvertExcelColIndexToAlphabet = TKX.ConvertExcelColIndexToAlphabet

// cui related
var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type TeaTableModel struct {
	Table table.Model
	Mode string // view, sel
	ReturnStatus int // -1: exit, 1 done
	Quitting bool
	KeepOnQuit bool
	Delegate func(...interface{}) interface{}
}

func (m TeaTableModel) Init() bubbletea.Cmd { return nil }

func (m TeaTableModel) Update(msg bubbletea.Msg) (bubbletea.Model, bubbletea.Cmd) {
	var cmd bubbletea.Cmd
	switch msg := msg.(type) {
	case bubbletea.KeyMsg:
		if m.Delegate != nil {
			
			rs := m.Delegate("key", int(msg.Type), msg.Runes, msg.Alt, msg.Paste, msg.String(), m.Table.Cursor(), m.Table.SelectedRow()).(string)
			
//			fmt.Printf("callback return: %v\n", rs)

			if strings.HasPrefix(rs, "TXERROR:") {
				bubbletea.Printf("input callback error: %v", rs[8:])
			}
			
			if rs == "show" {
				return m, bubbletea.Batch(
					bubbletea.Printf("[%v] %v", m.Table.Cursor(), ToJSONX(m.Table.SelectedRow())),
				)
			} else if rs == "select" {
				m.ReturnStatus = 1
				m.Quitting = true
				return m, bubbletea.Quit
			} else if rs == "exit" {
				m.ReturnStatus = -1
				m.Quitting = true
				return m, bubbletea.Quit
			} else if rs == "return" {
				m.ReturnStatus = 1
				m.Quitting = true
				return m, bubbletea.Quit
			} else if rs == "blur" {
				m.Table.Blur()
			}
		}

		switch msg.String() {
		case "esc", "ctrl+c":
			if m.Table.Focused() {
				m.Table.Blur()
				
				if m.Delegate == nil {
					m.ReturnStatus = -1
					m.Quitting = true
					return m, bubbletea.Quit
				}
			} else {
				m.Table.Focus()
			}
		case "q":
			if m.Delegate == nil {
				m.ReturnStatus = 1
				m.Quitting = true
				return m, bubbletea.Quit
			}
		case "enter":
			if m.Mode == "sel" {
				m.ReturnStatus = 1
				m.Quitting = true
				return m, bubbletea.Quit
			}
			
			if m.Delegate == nil {
				return m, bubbletea.Batch(
					bubbletea.Printf("[%v] %v", m.Table.Cursor(), ToJSONX(m.Table.SelectedRow())),
				)
			}
		}
	}
	m.Table, cmd = m.Table.Update(msg)
	return m, cmd
}

func (m TeaTableModel) View() string {
	if m.Quitting {
		if !m.KeepOnQuit {
			return ""
		}
	}
	
	return baseStyle.Render(m.Table.View()) + "\n"
}

func (pA *TK) ShowTableCompact(dataA [][]string, deleA func(...interface{}) interface{}, optsA ...string) string {
	modeT := GetSwitch(optsA, "-mode=", "")
	widthT := ToInt(GetSwitch(optsA, "-width=", "100"), 100)
	heightT := ToInt(GetSwitch(optsA, "-height=", "7"), 7)
	keepT := IfSwitchExists(optsA, "-keep")

	lenT := len(dataA)
	
	if lenT < 1 {
//		return "TXERROR:empty data"
		dataA = [][]string{}
	}
	
	columnsT := []table.Column{
//		{Title: "Rank", Width: 4},
//		{Title: "City", Width: 10},
//		{Title: "Country", Width: 10},
//		{Title: "Population", Width: 10},
	}
	
	colConfigStrT := strings.TrimSpace(GetSwitch(optsA, "-cols=", ""))
	
	rowsT := []table.Row{
//		{"1", "Tokyo", "Japan", "37,274,000"},
//		{"2", "Delhi", "India", "32,065,760"},
	}
	
	if colConfigStrT == "" {
		if len(dataA) < 1 {
			columnsT = []table.Column{}
		} else {
			len2 := len(dataA[0])
			
			for i := 0; i < len2; i ++ {
				columnsT = append(columnsT, table.Column{Title: dataA[0][i], Width: 10})
			}
			
			for i := 1; i < lenT; i ++ {
				rowsT = append(rowsT, dataA[i])
			}
		}
	} else {
		colDefAryT := JSONToMapStringStringArray(colConfigStrT)
		
		if colDefAryT == nil {
			return "TXERROR:invalid column definition"
		}
		
		len2 := len(colDefAryT)
		
		for i := 0; i < len2; i ++ {
			columnsT = append(columnsT, table.Column{Title: colDefAryT[i]["Title"], Width: ToInt(colDefAryT[i]["Width"], 0)})
		}
		
		for i := 0; i < lenT; i ++ {
			rowsT = append(rowsT, dataA[i])
		}
	}

	t := table.New(
		table.WithColumns(columnsT),
		table.WithRows(rowsT),
		table.WithFocused(true),
		table.WithHeight(heightT),
		table.WithWidth(widthT),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	m := TeaTableModel{
		Table: t,
		Mode: modeT,
		KeepOnQuit: keepT,
		Delegate: deleA,
	}
	
	rs, err := bubbletea.NewProgram(m).Run()
	
	if err != nil {
//		fmt.Println("Error running program:", err)
		return fmt.Sprintf("TXERROR:%v", err)
	}
	
	rs1 := rs.(TeaTableModel)
	
	if rs1.ReturnStatus < 0 {
		return "TXERROR:" + "cancel"
	}
	
//	fmt.Printf("%#v\n", rs1)
	
	return ToJSONX(map[string]interface{}{"Idx": rs1.Table.Cursor(), "Item": rs1.Table.SelectedRow()})
}

var ShowTableCompact = TKX.ShowTableCompact

func (pA *TK) StrDiff(s1 string, s2 string, optsA ...string) interface{} {
	dmp := diffmatchpatch.New()
	
	asRuneT := IfSwitchExists(optsA, "-asRune")

	formatT := GetSwitch(optsA, "-format=", "")

	var diffs []diffmatchpatch.Diff
	
	if asRuneT {
		diffs = dmp.DiffMain(s1, s2, false)
	} else {
		diffs = dmp.DiffMainRunes([]rune(s1), []rune(s2), false)
	}
	
	var rs interface{}
	
	switch formatT {
	case "", "html":
		rs = dmp.DiffPrettyHtml(diffs)
	case "text":
		rs = dmp.DiffPrettyText(diffs)
	case "delta":
		rs = dmp.DiffToDelta(diffs)
	case "json":
		rs = ToJSONX(diffs, "-sort")
	case "map", "object":
		rs = JSONToObject(ToJSONX(diffs))
	case "raw":
		rs = diffs
	default:
		rs = dmp.DiffPrettyHtml(diffs)
	}
	
//	fmt.Printf("%v\n", diffs)
//	fmt.Printf("%#v\n", diffs)
//	fmt.Printf("%v\n", dmp.DiffText1(diffs))
//	fmt.Printf("%v\n", dmp.DiffText2(diffs))
//	fmt.Printf("%v\n", dmp.DiffPrettyHtml(diffs))
//	fmt.Printf("%v\n", dmp.DiffPrettyText(diffs))
//	fmt.Printf("%v\n", dmp.DiffToDelta(diffs))
	
	return rs
}

var StrDiff = TKX.StrDiff

// cron related
func (pA *TK) IsValidCronExpr(exprA string) bool {
	return gronx.IsValid(exprA)
}

var IsValidCronExpr = TKX.IsValidCronExpr

func (pA *TK) IsCronExprDue(exprA string, optsA ...interface{}) interface{} {
	gron := gronx.New()
	
	if len(optsA) < 1 {
		rs, err := gron.IsDue(exprA)
		
		if err != nil {
			return err
		}
		
		return rs
	}
	
	timeT := ToTime(optsA[0], time.Now()).(time.Time)
	
	rs, err := gron.IsDue(exprA, timeT)
	
	if err != nil {
		return err
	}
	
	return rs
	
}

var IsCronExprDue = TKX.IsCronExprDue

func (pA *TK) SplitCronExpr(exprA string) interface{} {
//	rs, err := gronx.Segments(exprA)
//	
//	if err != nil {
//		return err
//	}

	rs := RegSplitX(strings.TrimSpace(exprA), `\s+`)
	
	if len(rs) < 5 {
		return fmt.Errorf("invalid crontab expression")
	}
	
	ary1 := rs[0:5]

	ary2 := rs[5:]

	return []string{strings.Join(ary1, " "), strings.Join(ary2, " ")}
}

var SplitCronExpr = TKX.SplitCronExpr

var taskerG *tasker.Tasker = nil
var taskerMutexG sync.Mutex


func (pA *TK) ResetTasker(optsA ...string) error {
	taskerMutexG.Lock()
	
	if taskerG != nil {
		taskerG.Stop()
	}
	
	optionsT := tasker.Option{
//		Verbose: ,
		// optional: defaults to local
//		Tz:      "Asia/Bangkok",
		// optional: defaults to stderr log stream
//		Out:     "/full/path/to/output-file",
	}
	
	if IfSwitchExists(optsA, "-verbose") {
		optionsT.Verbose = true
	}
	
	tzT := GetSwitch(optsA, "-tz=", "")
	
	if tzT != "" {
		optionsT.Tz = tzT
	}
	
	logFileT := GetSwitch(optsA, "-log=", "")
	
	if logFileT != "" {
		optionsT.Out = logFileT
	}
	
	taskerG = tasker.New(optionsT)
	
//	fmt.Printf("taskerG: %#v (%v)\n", taskerG, taskerG.Running())
	
	taskerMutexG.Unlock()

	return nil
}

var ResetTasker = TKX.ResetTasker

func (pA *TK) StopTasker() error {
	taskerMutexG.Lock()

	if taskerG != nil {
		if taskerG.Running() {
			taskerG.Stop()
		} else {
			taskerMutexG.Unlock()
			return fmt.Errorf("tasker not running")
		}
	} else {
		taskerMutexG.Unlock()
		return fmt.Errorf("tasker not initialzed")
	}

	taskerMutexG.Unlock()
	return nil
}

var StopTasker = TKX.StopTasker

func (pA *TK) RunTasker() error {
	taskerMutexG.Lock()

	if taskerG != nil {
		if taskerG.Running() {
			taskerG.Stop()
		}

		go taskerG.Run()
	} else {
		taskerMutexG.Unlock()
		return fmt.Errorf("tasker not initialzed")
	}

	taskerMutexG.Unlock()
	return nil
}

var RunTasker = TKX.RunTasker

// call ResetTaskerFirst, after this function, taskerG will stop
func (pA *TK) AddSimpleTask(exprA string, funcA func(argsA ...interface{}) interface{}) error {
	taskerMutexG.Lock()

	if taskerG != nil {
		if taskerG.Running() {
//			taskerG.Stop()
		}
	} else {
		taskerMutexG.Unlock()
		return fmt.Errorf("tasker not initialzed")
	}
	
	taskerG.Task(exprA, func(ctx context.Context) (int, error) {
		rs := funcA()
		
//		fmt.Printf("rs: %v\n", rs)
		
		if IsError(rs) {
//			fmt.Printf("err rs: %v\n", rs)
			return 1, fmt.Errorf("%v", rs)
		}

		// then return exit code and error, for eg: if everything okay
//		fmt.Printf("0 rs: %v\n", rs)
		return 0, nil
	}, false)
	
	taskerMutexG.Unlock()

	return nil
}

var AddSimpleTask = TKX.AddSimpleTask

func (pA *TK) AddShellTask(exprA string, shellCmdA string) (result error) {
	defer func() {
		r := recover()

		if r != nil {
			result = fmt.Errorf("failed to set map value: %v", r)
			return
		}
	}()

	taskerMutexG.Lock()

	if taskerG != nil {
		if taskerG.Running() {
//			taskerG.Stop()
		}
	} else {
		taskerMutexG.Unlock()
		return fmt.Errorf("tasker not initialzed")
	}
	
	taskerG.Task(exprA, taskerG.Taskify(shellCmdA, tasker.Option{}), false)
	
//	taskerG.Task(exprA, func(ctx context.Context) (int, error) {
//		aryT := ParseCommandLineCompact(shellCmdA)
//		
//		rs := SystemCmd(aryT[0], aryT[1:]...)
//		
//		if IsError(rs) {
//			return 1, fmt.Errorf("%v", rs)
//		}
//
//		// then return exit code and error, for eg: if everything okay
//		return 0, nil
//	})
	
	taskerMutexG.Unlock()

	return nil
}

var AddShellTask = TKX.AddShellTask

func (pA *TK) SetImageOpacity(imageA image.Image, opacityA float64) image.Image {
	if opacityA < 0 {
		opacityA = 0
	} else if opacityA > 1 {
		opacityA = 1
	}

	bounds := imageA.Bounds()
	mask := image.NewAlpha(bounds)
	
	r := int(opacityA * 255.0)
	
	for x := 0; x < bounds.Dx(); x++ {
		for y := 0; y < bounds.Dy(); y++ {
//			r := mapValues(alpha, 1, 0, 0, 255)
			mask.SetAlpha(x, y, color.Alpha{uint8(r)})
		}
	}

	mskWatermark := image.NewRGBA(bounds)
	draw.DrawMask(mskWatermark, bounds, imageA, image.ZP, mask, image.ZP, draw.Over)
	return mskWatermark
}

var SetImageOpacity = TKX.SetImageOpacity

func (pA *TK) AddWatermarkToImage(imageA image.Image, watermarkA image.Image, optsA ...string) image.Image {
	repeatT := IfSwitchExists(optsA, "-repeat")
	
	opacityT := ToFloat(GetSwitch(optsA, "-opacity="), 1.0)
	
	if opacityT < 0 {
		opacityT = 0
	} else if opacityT > 1 {
		opacityT = 1
	}

	bounds := watermarkA.Bounds()
	mask := image.NewAlpha(bounds)
	
	r := int(opacityT * 255.0)
	
	for x := 0; x < bounds.Dx(); x++ {
		for y := 0; y < bounds.Dy(); y++ {
//			r := mapValues(alpha, 1, 0, 0, 255)
			mask.SetAlpha(x, y, color.Alpha{uint8(r)})
		}
	}

	mskWatermark := image.NewRGBA(bounds)
	drawOri.DrawMask(mskWatermark, bounds, watermarkA, image.ZP, mask, image.ZP, draw.Over)
	
	widthT := ToFloat(imageA.Bounds().Max.X - imageA.Bounds().Min.X)
	if widthT < 0.0 {
		widthT = 0.0
	}
	
	heightT := ToFloat(imageA.Bounds().Max.Y - imageA.Bounds().Min.Y)
	if heightT < 0.0 {
		heightT = 0.0
	}
	
	dc := gg.NewContextForImage(imageA)

	if repeatT {
//	pattern1 := gg.NewSurfacePattern(mskWatermark, gg.RepeatBoth)
		dc.MoveTo(0, 0)
		
		dc.LineTo(widthT, 0)
		dc.LineTo(widthT, heightT)
		dc.LineTo(0, heightT)
		dc.LineTo(0, 0)
		
		dc.ClosePath()
	
		dc.SetFillStyle(gg.NewSurfacePattern(mskWatermark, gg.RepeatBoth))

		dc.Fill()
	} else {
//		dc.SetFillStyle(gg.NewSurfacePattern(mskWatermark, gg.RepeatNone))

		x := ToInt(GetSwitch(optsA, "-x=", "0"), 0)
		y := ToInt(GetSwitch(optsA, "-y=", "0"), 0)
		
		dc.DrawImage(mskWatermark, x, y)
		
		dc.Stroke()
	}
	
	

	return dc.Image()
}

var AddWatermarkToImage = TKX.AddWatermarkToImage

// return error in delegate func to stop the ticker thread
func (pA *TK) RunTickerFunc(secondsA float64, funcA func(argsA ...interface{}) interface{}) error {
	if funcA == nil {
		return fmt.Errorf("nil func")
	}
	
	funcT := func() {
		c := time.Tick(time.Duration(secondsA*1000) * time.Millisecond)
		
		for next := range c {
			rs := funcA(next)
			
	//		fmt.Printf("rs: %v\n", rs)
			
			if IsError(rs) {
	//			fmt.Printf("err rs: %v\n", rs)
				return
			}
		}
	}
	
	go funcT()
	
	return nil
}

var RunTickerFunc = TKX.RunTickerFunc



