// github issue: n/a
// expected value: [1, 2, 3]
// expected type: ARRAY

a := [
    "1",
    "22",
    "333",
]

b := a.map(func(x) { len(x) } )
