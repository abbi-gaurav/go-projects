package main

import "github.com/tidwall/gjson"

//https://github.com/tidwall/gjson#path-syntax

const json = `{
				"stats": [
					{
						"name": "http.inbound_0.0.0.0_8080.downstream_cx_active",
						"value": 65
					}
					]
				}`

func extraActiveConnectionCount() {
	activeCounts := gjson.Get(json, "stats.#(name==http.inbound_0.0.0.0_8080.downstream_cx_active).value")
	println(activeCounts.String())
}
