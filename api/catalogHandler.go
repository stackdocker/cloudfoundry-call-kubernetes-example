/*
Copyright 2015 All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package api

import (
    "net/http"
    "encoding/json"
    "io"
    "bytes"
    "log"
)

type catalogHandler interface {
    authenticate(w http.ResponseWriter, r *http.Request) bool
}

func (c *CatalogV2) authenticate(w http.ResponseWriter, r *http.Request) bool {
    if r.URL.User != nil && r.URL.User.Username() != "" {
        log.Print("Basic Authentication: ", r.URL.User.Username())
        if password, _ := r.URL.User.Password(); password != "" {
            log.Print(" : ", password)
        }
        log.Println()
    }
    return true
}

var catalog *CatalogV2 = &CatalogV2 {
    Services : []ServiceV2 {
        ServiceV2 {
            Id: "9c372bbc-1e7b-472b-bcb6-eeda5b21eb35",    
            Name: "redis-cluster-managed-by-kubernetes",
            Description: `The Redis is a high reliable and scalable cluster deployed upon` + 
                ` Kubernetes v1, it failovers in master/slave, and load balancing with multiple` + 
                ` sentinel nodes`,
            Bindable: true,
            Tags: []string{"redis", "cluster", "k-v", "database"},
            Plans: []ServicePlanV2 {
                ServicePlanV2 {
                    Id: "8cfbbaf5-efdb-41c1-89ab-f797185f7818",
                    Name: "demo",
                    Description: "this is a redis cluster demo",
                    Free: true,
                },
            },
        },
    },
}

/*
    curl -H "X-Broker-API-Version: 2.6" http://username:password@broker-url/v2/catalog
*/

func HandleCatalog(w http.ResponseWriter, r *http.Request) {
    if r.Method != "GET" {
        http.Error(w, "Api only support GET method.", http.StatusMethodNotAllowed)
        return
    }
    if v := r.Header["X-Broker-API-Version"]; len(v) > 0 && v[0] != "2.6" {
        http.Error(w, "Unmatched API version.", http.StatusPreconditionFailed)
        return
    }
    
    if !catalog.authenticate(w, r) {
        http.Error(w, "Not authorized", http.StatusUnauthorized)
        return
    }
    
    //enc := json.NewEncoder(w)
    //enc.Encode(catalog)
    js, err := json.Marshal(catalog)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    //w.Write(js)
    io.Copy(w, bytes.NewBuffer(js))
}