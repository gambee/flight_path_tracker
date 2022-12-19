# Flight Path Tracker

#### Quickstart

You should be able to build a docker image and run that image locally on port 8080
with `make run_image`.

To build and run the go program only (i.e. not in docker) run `make run`.

To run the prewritten tests (ensure the server is running), run `make test`.

If your terminal doesn't scroll, or you would like it to be easier to navigate the
output of the tests, run one of the following: `make test | less`, `make test | more`, `make test > some_file`.


#### Example Request (POST)

`
[["PDX","LAX"],["ATL","EFG"],["LAX","ATL"]]
`
#### Example Response
`
{
    "flights": [
        [
            "PDX",
            "LAX"
        ],
        [
            "ATL",
            "EFG"
        ],
        [
            "LAX",
            "ATL"
        ]
    ],
    "flightChain": [
        "PDX",
        "LAX",
        "ATL",
        "EFG"
    ],
    "layovers": [
        "LAX",
        "ATL"
    ],
    "start": "PDX",
    "end": "EFG"
}
`
