package main

import (
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/proxy"
	"github.com/satori/go.uuid"
	//"matou-sakura/config"
	"fmt"
	"time"
	//"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"math/rand"
	"os"
)

type ImageDownloader struct{

}

func (imgDL *ImageDownloader) Download(imgUrl string) error {
	//fmt.Println("To download url: ", imgUrl)
	uuid, _ := uuid.NewV4()
	saveDir := "./download/"
	imgName := saveDir + uuid.String() + ".jpg"
	dl_proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse("http://127.0.0.1:1080")
	}
	transport := &http.Transport{Proxy: dl_proxy}
	client := &http.Client{Transport: transport}
	resp, err := client.Get(imgUrl)
	if err != nil {
		fmt.Println("Http error when download image : ", imgUrl, "\nError msg: ", err, " UUID: ", uuid)
		return nil
	 }
	defer resp.Body.Close()
	file, err_file := os.Create(imgName)
	if err_file != nil {
		fmt.Println("Error when create file: ", imgName, "\nError msg: ", err_file, " UUID: ", uuid)
		return nil
	}
	defer file.Close()
	//io.Copy(file, resp.Body)
	content, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		fmt.Println("Error in ioutil : ", imgName, "\nError msg: ", err2, " UUID: ", uuid)
		return nil
	}
	file.Write(content)
	//fmt.Println("download img ", imgUrl, " finish")
	return nil
}

var imageDownloader = &ImageDownloader{}

type YandeCollector struct {
	MainUrl string
	SearchKeyWord string
	YandeColly *colly.Collector
	task chan string
	done chan struct{}
}

var YandeHandler = &YandeCollector{}

func (yc *YandeCollector) Prepare(keyword string) error {
	yc.task = make(chan string, 20)
	yc.done = make(chan struct{})
	for i := 0; i < 3; i++ {
		go func(){
			for imgUrl := range yc.task {
				imageDownloader.Download(imgUrl)
				num := rand.Int31n(5)
				time.Sleep(time.Duration(num)*time.Second)
			}
			yc.done <- struct{}{}
		}()
	}
	yc.MainUrl = "https://yande.re"
	yc.SearchKeyWord = keyword
	yc.YandeColly = colly.NewCollector(
		colly.AllowedDomains("yande.re"),
		//colly.MaxDepth(2),
		colly.Async(true),
	)
	rp, err := proxy.RoundRobinProxySwitcher("http://127.0.0.1:1080")
	if err != nil {
		fmt.Println("Colly Proxy set error")
		return nil
	}
	yc.YandeColly.SetProxyFunc(rp)
	yc.YandeColly.Limit(&colly.LimitRule{
		DomainGlob: "*",
		Parallelism: 2,
		RandomDelay: 5*time.Second,
	})
	yc.YandeColly.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL)
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36")
	})
	yc.YandeColly.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited: ", r.Request.URL)
	})
	//yc.YandeColly.OnHTML("a[class=thumb][href]", func(e *colly.HTMLElement) {
		// Goquery 筛选出图片预览页面
	//	imgPage := e.Attr("href")
	//	imgPageUrl := yc.MainUrl + imgPage
	//	e.Request.Visit(imgPageUrl)
	//})
	yc.YandeColly.OnHTML("a[class^=directlink][class$=largeimg][href]", func(e *colly.HTMLElement) {
		// 前往大图进行下载
		downloadUrl := e.Attr("href")
		yc.task <- downloadUrl
		//go imageDownloader.Download(downloadUrl, yc.task)

	})
	yc.YandeColly.OnHTML("a[class=next_page][rel=next][href]", func(e *colly.HTMLElement) {
		// 下一页
		relateURL := e.Attr("href")
		nextPage := yc.MainUrl + relateURL
		e.Request.Visit(nextPage)
	})

	yc.YandeColly.OnError(func(r *colly.Response, e error) {
		fmt.Println("Request URL: ", r.Request.URL, " failed with response: ", r, "\nError", e)
		fmt.Println("Retrying url: ", r.Request.URL)
		r.Request.Retry()
	})
	return nil
}

func (yc *YandeCollector) Start() {
	searchUrl := "https://yande.re/post?page=9&tags=" + yc.SearchKeyWord + "+"
	yc.YandeColly.Visit(searchUrl)
	yc.YandeColly.Wait()
}

func (yc *YandeCollector) End() {
	close(yc.task)
	for i := 0; i < 3; i++ {
		<-yc.done
	}
}

func main() {
	searchKeyWord := "matou_sakura"
	//confVip := config.LoadConfig()
	//fmt.Println(confVip.GetString("proxy.addr"))
	YandeHandler.Prepare(searchKeyWord)
	YandeHandler.Start()
	YandeHandler.End()
}
