package handlers

import (
    "io"
    "net/http"
    "encoding/json"
    "github.com/gambee/flight_path_tracker/internal/flightpath"
)


func Calculate(w http.ResponseWriter, r *http.Request) {
    var fp flightpath.FlightPath
    if bdy, err := io.ReadAll(r.Body); err != nil {
        w.WriteHeader(400)
        w.Write([]byte(err.Error()))
    } else {
        if err := json.Unmarshal(bdy, &fp); err != nil {
            w.WriteHeader(400)
            w.Write([]byte(err.Error()))
        } else if err := fp.BuildChain(); err != nil {
            w.WriteHeader(400)
            w.Write([]byte(err.Error()))
        } else if bs, err := json.MarshalIndent(fp, "",  "    "); err != nil {
            // NOTE: I only use MarshalIndent here for easier displaying of tests
            //       Would not use in normal production code.
            w.WriteHeader(400)
            w.Write([]byte(err.Error()))
            // Also, for how many times these two lines are repeated, I would 
            // likely abstract error handling better.
        } else {
            w.Write(bs)
        }
    }
}
