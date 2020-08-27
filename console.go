// +build !aws
// +build !ifttt_ut

package main

import (
	"github.com/comail/colog"
	"log"
	"sync"
	"time"
)

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
