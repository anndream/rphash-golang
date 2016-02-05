package main;

import (
  "fmt"
  "io/ioutil"
  "github.com/wenkesj/rphash/api"
  "github.com/wenkesj/rphash/parse"
);

var numberOfClusters = 6;

const (
  exampleInputFileName = "input.json";
  exampleOutputFileName = "output.json";
  exampleDataLabel = "people";
);

func main() {
  parser := parse.NewParser();
  bytes, _ := ioutil.ReadFile(exampleInputFileName);
  jsonData := parser.BytesToJSON(bytes);
  data := parser.JSONToFloat64Matrix(exampleDataLabel, jsonData);
  cluster := api.NewRPHash(data, numberOfClusters);
  cluster.Run();
  centroids := cluster.GetCentroids();
  jsonCentroids := parser.Float64MatrixToJSON(exampleDataLabel, centroids);
  jsonBytes := parser.JSONToBytes(jsonCentroids);
  err := ioutil.WriteFile(exampleOutputFileName, jsonBytes, 0644);
  if err != nil {
    panic(err);
  }
};
