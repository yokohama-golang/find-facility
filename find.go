package main

import (
	"fmt"
	"github.com/comail/colog"
	"github.com/sclevine/agouti"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

func screenshot(page *agouti.Page) {
	log.Printf("trace: page.Screenshot(golang.png)")
	if err := page.Screenshot("golang.png"); err != nil {
		log.Fatalf("error: Failed to screenshot: %v", err)
	}
}

func submitByID(page *agouti.Page, id string) {
	log.Printf("trace: page.FindByID(%s).Submit()", id)
	if err := page.FindByID(id).Submit(); err != nil {
		screenshot(page)
		log.Fatalf("error: Failed to submit: %v", err)
	}
}

func clickByID(page *agouti.Page, id string) {
	log.Printf("trace: page.FindByID(%s).Click()", id)
	if err := page.FindByID(id).Click(); err != nil {
		screenshot(page)
		log.Fatalf("error: Failed to click: %v", err)
	}
}

func fillByID(page *agouti.Page, id, text string) {
	log.Printf("trace: page.FindByID(%s).Fill(%s)", id, text)
	if err := page.FindByID(id).Fill(text); err != nil {
		screenshot(page)
		log.Fatalf("error: Failed to fill: %v", err)
	}
}

func allByLabel(page *agouti.Page, label string) {
	log.Print("trace: page.AllByLabel().Click()", label)
	if err := page.AllByLabel(label).Click(); err != nil {
		screenshot(page)
		log.Fatalf("error: Failed to click: %v", err)
	}
}

func visitWelcome(page *agouti.Page) {
	clickByID(page, "btnNormal")
}

func visitSystemMenu(page *agouti.Page) {
	clickByID(page, "rbtnYoyaku")
}

func visitFacilitySearch(page *agouti.Page) {
	clickByID(page, "dgTable_ctl09_chkShisetsu")
	clickByID(page, "ucPCFooter_btnForward")
}

func visitDateAndTimeSelection(page *agouti.Page, year, month, day int) {
	fillByID(page, "txtYear", strconv.Itoa(year))
	fillByID(page, "txtMonth", strconv.Itoa(month))
	fillByID(page, "txtDay", strconv.Itoa(day))
	clickByID(page, "rbtnWeek")
	clickByID(page, "chkSat")
	clickByID(page, "ucPCFooter_btnForward")
}

type reserverStatusPerFacilityVisitor func(name, maxHumans, status string, statusSelection *agouti.Selection)

func iterateReservedStatusPerFacility(page *agouti.Page, id string, f reserverStatusPerFacilityVisitor) {
	log.Printf("trace: page.FindByID(%s).All(tr)", id)
	trs := page.FindByID(id).All("tr")
	ntrs, _ := trs.Count()
	for i := 0; i < ntrs; i++ {
		tr := trs.At(i)
		tds := tr.All("td")

		name, _ := tds.At(0).Text()
		maxHumans, _ := tds.At(1).Text()
		statusSelection := tds.At(2)
		status, _ := statusSelection.Text()
		status = strings.TrimSpace(status)
		f(name, maxHumans, status, statusSelection)
	}
}

func visitReservedStatusPerFacility(page *agouti.Page) {
	f := func(name, maxHumans, status string, statusSelection *agouti.Selection) {
		switch status {
		case "△":
		case "○":
		default:
			return
		}
		log.Print("debug: Selecting " + name + "," + maxHumans + "," + status)
		statusSelection.Click()
	}
	iterateReservedStatusPerFacility(page, "dlRepeat_ctl00_tpItem_dgTable", f)
	clickByID(page, "ucPCFooter_btnForward")
}

type reserverStatusPerHourVisitor func(date, name, maxHumans string, status []string)

func iterateReservedStatusPerHour(page *agouti.Page, id string, f reserverStatusPerHourVisitor) error {
	log.Printf("trace: page.FindByID(%s).All(tr)", id)
	trs := page.FindByID(id).All("tr")
	ntrs, err := trs.Count()
	if err != nil {
		return fmt.Errorf("Not found id: %s", id)
	}
	tr := trs.At(0)
	tds := tr.All("td")
	date, _ := tds.At(0).Text()
	date = strings.Replace(date, "\n", "", -1)

	for i := 1; i < ntrs; i++ {
		const NAME_INDEX = 0
		const MAX_HUMANS_INDEX = 1
		const BEGIN_STATUS_INDEX = 2
		const LAST_STATUS_INDEX = 14
		const SIZE_STATUS = 12

		tr := trs.At(i)
		tds := tr.All("td")

		name, _ := tds.At(NAME_INDEX).Text()
		maxHumans, _ := tds.At(MAX_HUMANS_INDEX).Text()
		status := make([]string, 0, SIZE_STATUS)
		for j := BEGIN_STATUS_INDEX; j < LAST_STATUS_INDEX; j++ {
			s, _ := tds.At(j).Text()
			s = strings.TrimSpace(s)
			status = append(status, s)
		}
		f(date, name, maxHumans, status)
	}
	return nil
}

func visitReservedStatusPerHour(page *agouti.Page) {
	const (
		FROM_9 = iota
		FROM_10
		FROM_11
		FROM_12
		FROM_13
		FROM_14
		FROM_15
		FROM_16
		FROM_17
		FROM_18
		FROM_19
		FROM_20
	)
	var timeOfHope = []int{FROM_13, FROM_14, FROM_15, FROM_16}
	f := func(date, name, maxHumans string, status []string) {
		for _, t := range timeOfHope {
			if status[t] != "○" {
				return
			}
		}
		log.Printf("Found empty facility:  %s,%s,%s", date, name, maxHumans)
	}
	for i := 0; ; i++ {
		id := fmt.Sprintf("dlRepeat_ctl%02d_tpItem_dgTable", i)
		err := iterateReservedStatusPerHour(page, id, f)
		if err != nil {
			return
		}
	}
}

func gotoSystemMenu(page *agouti.Page) {
	clickByID(page, "ucPCFooter_btnToMenu")
}

func findFacility(from, to time.Time) {
	log.Printf("from:%s to:%s", from, to)
	driver := agouti.PhantomJS()
	if err := driver.Start(); err != nil {
		log.Fatalf("error: Failed to start driver: %v", err)
	}
	defer driver.Stop()

	page, err := driver.NewPage()
	log.Printf("trace: page, err := driver.NewPage()")
	if err != nil {
		screenshot(page)
		log.Fatalf("Failed to open page: %v", err)
	}
	const url = "https://yoyaku.e-kanagawa.lg.jp/Kanagawa/Web/Wg_KoukyouShisetsuYoyakuMoushikomi.aspx"

	log.Printf("trace: page.Navigate(%s)", url)
	if err := page.Navigate(url); err != nil {
		screenshot(page)
		log.Fatalf("Failed to navigate: %v", err)
	}

	visitWelcome(page)

	for t := from; t.Before(to); t = t.AddDate(0, 0, 7) {
		visitSystemMenu(page)
		visitFacilitySearch(page)
		visitDateAndTimeSelection(page, t.Year(), int(t.Month()), t.Day())
		visitReservedStatusPerFacility(page)
		visitReservedStatusPerHour(page)
		gotoSystemMenu(page)
	}
}

func main() {
	colog.SetDefaultLevel(colog.LInfo)
	colog.SetMinLevel(colog.LInfo)
	colog.Register()

	const START_DAYS = 7 * 4
	const INTERVAL_DAYS = 7 * 8
	now := time.Now()
	from := now.AddDate(0, 0, START_DAYS)
	to := from.AddDate(0, 0, INTERVAL_DAYS*2)
	wg := sync.WaitGroup{}
	////////////////////////////////////////////////////////////////////////
	// $ time go run find.go (4 weeks version)                            //
	// go: finding github.com/comail/colog latest                         //
	// [  info ] 2019/04/05 08:47:13 Finished.                            //
	//     go run find.go  74.61s user 4.91s system 26% cpu 5:02.34 total //
	////////////////////////////////////////////////////////////////////////

	////////////////////////////////////////////////////////////////////////////////////////////////
	// $ time go run find.go (16 weeks version)                                                   //
	// go: finding github.com/comail/colog latest                                                 //
	// [  info ] 2019/04/05 08:59:59 Found empty facility:  2019年6月29日（土）,ホール,260人      //
	// [  info ] 2019/04/05 09:01:40 Found empty facility:  2019年7月13日（土）,３０１会議室,90人 //
	// [  info ] 2019/04/05 09:09:19 Found empty facility:  2019年8月24日（土）,６０２会議室,18人 //
	// [  info ] 2019/04/05 09:09:30 Finished.                                                    //
	//     go run find.go  286.52s user 16.97s system 25% cpu 19:33.50 total                      //
	////////////////////////////////////////////////////////////////////////////////////////////////

	///////////////////////////////////////////////////////////////////////////////////////////
	// $ time go run find.go  (8 weeks version per goroutine, 2 goroutine)                   //
	// go: finding github.com/comail/colog latest                                            //
	// [  info ] 2019/04/08 01:37:24 Found empty facility:  2019年6月29日（土）,ホール,260人 //
	// [  info ] 2019/04/08 01:37:26 Finished.                                               //
	//     go run find.go  271.17s user 16.07s system 44% cpu 10:44.46 total                 //
	///////////////////////////////////////////////////////////////////////////////////////////
	for t := from; t.Before(to); t = t.AddDate(0, 0, INTERVAL_DAYS) {
		wg.Add(1)
		go func(from, to time.Time) {
			findFacility(from, to)
			wg.Done()
		}(t, t.AddDate(0, 0, INTERVAL_DAYS))
	}
	wg.Wait()
	log.Printf("Finished.")
}