package product

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"path"
	"text/template"
	"time"
)

type ProductReportFilter struct {
	NameFilter			string `json:"productName"`
	ManufacturerFilter 	string `json:"manufacturer"`
	SKUFilter			string `json:"sku"`
}

func handleProductReport(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var productFilder ProductReportFilter
		err := json.NewDecoder(r.Body).Decode(&productFilder)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		products, err := searchForProductData(productFilder)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		t := template.New("report.gotmpl").Funcs(template.FuncMap{"mod": func(i, x int) bool {return i%x == 0}})
		t, err = t.ParseFiles(path.Join("templates", "report.gotmpl"))
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var tmpl bytes.Buffer
		if len(products) <= 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		err = t.Execute(&tmpl, products)
		rdr := bytes.NewReader(tmpl.Bytes())
		w.Header().Set("Content-Disposition", "Attachement")
		http.ServeContent(w, r , "report.html", time.Now(), rdr)
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}