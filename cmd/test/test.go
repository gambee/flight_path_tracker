package main

import (
    "fmt"
    "bytes"
    "net/http"
    "io/ioutil"
)

const (
    cType = `application/json`
    calcURL = `http://localhost:8080/calculate`
)

var plds []string = []string{
    // Success expected for the following
    `[["PDX", "LAX"],["ATL","EFG"],["LAX","ATL"]]`,
    `[["PDX", "LAX"],["ATL","EFG"],["LAX","ATL"],["ABC","PDX"]]`,

    // Errors expected for the following
    `[["PDX","LAX"],["ATL","EFG"],["XYZ","LAX"]]`, // multiple arriving at LAX
    `[["PDX","LAX"],["ATL","EFG"],["PDX","XYZ"]]`, // multiple departing PDX
    `[["PDX","LAX"],["ATL","PDX"],["LAX","ATL"]]`, // full loop
    `[["XYZ","ABC"],["PDX","LAX"],["ATL","PDX"],["LAX","ATL"]]`, // partial loop
}


func main() {
    for _, v := range plds {
        if e := TestCalculate(v); e != nil {
            fmt.Println(e)
        }
    }
}

func TestCalculate(payload string) error {
    pl := bytes.NewReader([]byte(payload))
    if resp, err := http.Post(calcURL, cType, pl); err != nil {
        return err
    } else {
        if bdy, err := ioutil.ReadAll(resp.Body); err != nil {
            panic(err)
        } else {
            fmt.Printf("request:\n%s\nresponse_code: %d\nresponse:\n%s\n\n", payload, resp.StatusCode, bdy)
        }
    }
    return nil
}

