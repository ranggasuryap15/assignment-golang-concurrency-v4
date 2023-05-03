package main

import (
	"fmt"
	"time"
)

type RowData struct {
	RankWebsite int
	Domain      string
	TLD         string
	IDN_TLD     string
	Valid       bool
	RefIPs      int
}

func GetTLD(domain string) (TLD string, IDN_TLD string) {
	var ListIDN_TLD = map[string]string{
		".com": ".co.id",
		".org": ".org.id",
		".gov": ".go.id",
	}

	for i := len(domain) - 1; i >= 0; i-- {
		if domain[i] == '.' {
			TLD = domain[i:]
			break
		}
	}

	if _, ok := ListIDN_TLD[TLD]; ok {
		return TLD, ListIDN_TLD[TLD]
	} else {
		return TLD, TLD
	}
}

func ProcessGetTLD(website RowData, ch chan RowData, chErr chan error) {
	TLD, IDN_TLD := GetTLD(website.Domain)
	
	// jika website domain kosong
	if website.Domain == "" {
		chErr <- fmt.Errorf("domain name is empty")
		return
	}

	// jika valid bernilai false
	if !website.Valid {
		chErr <- fmt.Errorf("domain not valid")
		return
	}

	if website.RefIPs == -1 {
		chErr <- fmt.Errorf("domain RefIPs not valid")
		return
	}

	ch <- RowData{
		RankWebsite: website.RankWebsite,
		Domain: website.Domain,
		TLD: TLD,
		IDN_TLD: IDN_TLD,
		Valid: website.Valid,
		RefIPs: website.RefIPs,
	}
}

// Gunakan variable ini sebagai goroutine di fungsi FilterAndGetDomain
var FuncProcessGetTLD = ProcessGetTLD

func FilterAndFillData(TLD string, data []RowData) ([]RowData, error) {
	ch := make(chan RowData, len(data))
	errCh := make(chan error)
	rowData := []RowData{}
	
	for _, website := range data {
		go FuncProcessGetTLD(website, ch, errCh)
		
		if website.Domain == "" || !website.Valid || website.RefIPs == -1 {
			return rowData, <-errCh
		}
		time.Sleep(250 * time.Millisecond)
		select {
		case value := <-ch:
			rowData = append(rowData, value)
		}
	}
	
	return rowData, nil
}

// gunakan untuk melakukan debugging
func main() {
	rows, err := FilterAndFillData(".com", []RowData{
		{1, "google.com", "", "", true, 100},
		{2, "facebook.com", "", "", true, 100},
		{3, "golang.org", "", "", true, 100},
	})

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(rows)
}
