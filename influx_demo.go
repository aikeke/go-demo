package main

import (
	"fmt"
	"log"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
)

// influxdb demo

func connInflux() client.Client {
	cli, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "https://xxx:8086",
		Username: "dev",
		Password: "",
	})
	if err != nil {
		log.Fatal(err)
	}
	return cli
}

// query
func queryDB(cli client.Client, cmd string) (res []client.Result, err error) {
	fmt.Println(cmd)
	q := client.Query{
		Command:  cmd,
		Database: "device_shadow",
	}
	if response, err := cli.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}

// insert
func writesPoints(cli client.Client) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "device_shadow",
		Precision: "ns", //精度，默认ns
	})
	if err != nil {
		log.Fatal(err)
	}
	tags := map[string]string{
		"attributeCode":    "Battery",
		"attributeId":      "1024006",
		"attributeName":    "电量",
		"deviceEntityCode": "this-deviceistest",
		"deviceEntityId":   "1208244",
		"deviceEntityName": "爱牵挂智能手表",
		"deviceTypeCode":   "aqg_smart_watch",
		"deviceTypeId":     "1018002",
		"deviceTypeName":   "爱牵挂智能手表",
		"spaceId":          "2052062",
	}
	fields := map[string]interface{}{
		"projectSpaceId": 22026322,
		"value":          "",
		"valueFloat":     0.0,
		"valueNumber":    0,
	}

	pt, err := client.NewPoint("shadow_update_short", tags, fields, time.Unix(0, 1655200657400249948))
	if err != nil {
		log.Fatal(err)
	}
	bp.AddPoint(pt)
	err = cli.Write(bp)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("insert success")
}

func main() {
	conn := connInflux()
	fmt.Println(conn)

	// insert
	writesPoints(conn)

	// // 获取10条数据并展示
	// qs := fmt.Sprintf("SELECT * FROM %s WHERE  time >=1651334400000000000 limit 10", "shadow_update_short")
	// res, err := queryDB(conn, qs)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// for _, row := range res[0].Series[0].Values {
	// 	for j, value := range row {
	// 		log.Printf("j:%d value:%v\n", j, value)
	// 	}
	// }
}
