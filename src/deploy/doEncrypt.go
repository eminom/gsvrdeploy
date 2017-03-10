
package main

import (
    "github.com/xxtea/xxtea-go/xxtea"
    "io/ioutil"
    "os"
    _ "fmt"
)

func getEncryptedBuffer(path string)[]byte{
    c, err := ioutil.ReadFile(path)
    if nil != err {
        panic(err)
    }
    if string(c[:len(xxtea_sig)]) == xxtea_sig {
        //fmt.Printf("%v encrypted already.\n", path)
        return c
    }
    //fmt.Printf("%v encrypted.\n", path)
    return xxtea.Encrypt(c, []byte(xxtea_key))
}

func doXXTeaToUno(from string, to string){
    theEnc := getEncryptedBuffer(from)
    fout, err := os.Create(to)
    if nil != err {
        panic(err)
    }
    defer fout.Close()
    fout.Write([]byte(xxtea_sig))
    fout.Write(theEnc)
}



