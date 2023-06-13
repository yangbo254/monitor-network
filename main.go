package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/aliyun-log-go-sdk/producer"
	_ "github.com/joho/godotenv/autoload"
)

func GetUri(srcUrl string) (*url.URL, error) {
	return url.Parse(srcUrl)
}

func httpGetUseHttpProxy(descUrl string, uri *url.URL) (int64, error) {

	client := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(uri),
		},
	}
	start := time.Now()
	resp, err := client.Get(descUrl)
	if err != nil {
		return 0, err
	}
	elapsed := time.Since(start)
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return elapsed.Milliseconds(), nil
	}

	return 0, fmt.Errorf("resp status code error :%v", resp.Request.Response.Status)
}

var producerInstance *producer.Producer

type Callback struct {
}

func (callback *Callback) Success(result *producer.Result) {
	// attemptList := result.GetReservedAttempts() // 遍历获得所有的发送记录
	// for _, attempt := range attemptList {
	// 	fmt.Println(attempt)
	// }
}

func (callback *Callback) Fail(result *producer.Result) {
	fmt.Println(result.IsSuccessful())        // 获得发送日志是否成功
	fmt.Println(result.GetErrorCode())        // 获得最后一次发送失败错误码
	fmt.Println(result.GetErrorMessage())     // 获得最后一次发送失败信息
	fmt.Println(result.GetReservedAttempts()) // 获得producerBatch 每次尝试被发送的信息
	fmt.Println(result.GetRequestId())        // 获得最后一次发送失败请求Id
	fmt.Println(result.GetTimeStampMs())      // 获得最后一次发送失败请求时间
}

func reportInit() {
	producerConfig := producer.GetDefaultProducerConfig()
	producerConfig.Endpoint = os.Getenv("Endpoint")
	producerConfig.AccessKeyID = os.Getenv("AccessKeyID")
	producerConfig.AccessKeySecret = os.Getenv("AccessKeySecret")
	producerInstance = producer.InitProducer(producerConfig)
	producerInstance.Start()
}

func report(src, cn, hk, us string) {
	callBack := &Callback{}
	log := producer.GenerateLog(uint32(time.Now().Unix()), map[string]string{"CN": cn, "HK": hk, "US": us})
	fmt.Println(log.GetContents())
	_ = producerInstance.SendLogWithCallBack(os.Getenv("Project"), os.Getenv("Logstore"), "topic", src, log, callBack)
	// if err != nil {
	// 	fmt.Println(err)
	// }
}

var cliUseMysql = flag.Bool("usemysql", false, "whether to use mysql")

func main() {
	var (
		httpProxyArray         []string
		descCN, descHK, descSG string
	)
	flag.Parse()
	if *cliUseMysql {
		DbInit()
	}

	reportInit()
	StartHttpEngine()
	for {
		httpProxyEnv := os.Getenv("HTTP_PROXY_LIST")
		if len(httpProxyEnv) != 0 {
			httpProxyArray = strings.Split(httpProxyEnv, ",")
		}
		descCN = os.Getenv("DESC_CN")
		descHK = os.Getenv("DESC_HK")
		descSG = os.Getenv("DESC_SG")

		if len(descCN) == 0 {
			descCN = "https://www.shanghai.gov.cn/"
		}
		if len(descHK) == 0 {
			descHK = "https://www.google.com.hk/"
		}
		if len(descSG) == 0 {
			descSG = "https://www.zaobao.com.sg/"
		}

		funcCheck := func(desc string, proxyUrl *url.URL) int64 {
			if elapsed, err := httpGetUseHttpProxy(descCN, proxyUrl); err == nil {
				return elapsed
			} else {
				return 0
			}
		}
		for _, proxyUrl := range httpProxyArray {
			var (
				elapsedCN int64
				elapsedHK int64
				elapsedUS int64
				Host      string
			)
			Host = "errorProxy"
			pUrl, err := GetUri(proxyUrl)
			if err == nil {
				Host = pUrl.Host
				elapsedCN = funcCheck(descCN, pUrl)
				elapsedHK = funcCheck(descHK, pUrl)
				elapsedUS = funcCheck(descSG, pUrl)
			}
			DbSetVistorRecords(time.Now().Unix(), Host, elapsedCN, elapsedHK, elapsedUS)
			report(Host, strconv.FormatInt(elapsedCN, 10), strconv.FormatInt(elapsedHK, 10), strconv.FormatInt(elapsedUS, 10))
		}

		time.Sleep(time.Millisecond * 30 * 1000)
	}

}
