package main

import (
	"crypto/md5"
	"fmt"
	"github.com/OneOfOne/xxhash"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func (self *calcPackOut) toStr() string {
	return fmt.Sprintf("%v\t%v\t%v\t%v", self.hashName, self.size, self.xxhash, self.relpath)
}

func calcFileMd5(path string) string {
	fin, err := os.Open(path)
	if nil != err {
		panic(err)
	}
	defer fin.Close()
	h := md5.New()
	_, e2 := io.Copy(h, fin)
	if nil != e2 {
		panic(err)
	}
	// MD5:总是32位宽的16进制的数字;
	return fmt.Sprintf("%x", h.Sum(nil))
}

func calcFileXXHash(path string) string {
	fin, err := os.Open(path)
	if nil != err {
		panic(err)
	}
	defer fin.Close()
	h := xxhash.NewS32(10241024) // New32(注意是10进制)
	_, e2 := io.Copy(h, fin)
	if nil != e2 {
		panic(e2)
	}
	return fmt.Sprintf("%08x", h.Sum32())
}

//定制标准ENTRY
func calcFileEnt(abs_path string, relpath string) interface{} {
	fi, err := os.Stat(abs_path)
	if nil != err {
		panic(err)
	}
	md5 := calcFileMd5(abs_path)
	xxhash := calcFileXXHash(abs_path)
	relpath = strings.Replace(relpath, "\\", "/", -1)
	extName := filepath.Ext(fi.Name())
	outName := fmt.Sprintf("%s%s", md5, extName)
	return &calcPackOut{
		hashName: outName,
		size:     int(fi.Size()),
		xxhash:   xxhash,
		relpath:  relpath,
		abs_path: abs_path,
	}
}

func hacerLaCama(p interface{}) interface{} {
	fpath := p.(*calcPackIn)
	return calcFileEnt(fpath.abs_path, fpath.relpath)
}
