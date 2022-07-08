package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"syscall"
	"time"
)

type Ntof struct {
	token  string
	cron   *cron.Cron
	client *Client
}

func newNtof() *Ntof {
	return &Ntof{
		cron:   cron.New(cron.WithSeconds()),
		client: NewClient(),
	}
}

const (
	GoodSIdShangWu = 2
	GoodSIdXiaWu   = 3

	BaseURL     = "http://h5.lingyuan888.com/index.php/WebApi"
	LoginURL    = "/User/login"
	GoodListURL = "/Goods/goods"
	BuyURL      = "/Order/buy"
)

func (n *Ntof) Login(username, password string) error {
	params := map[string]string{
		"uphone":   username,
		"password": password,
	}

	out := struct {
		Token    string
		NickName string
	}{}
	err := n.client.PostForm(BaseURL+LoginURL, params, &out)
	if err != nil {
		return err
	}

	fmt.Println("login success, token:", out.Token)

	n.token = out.Token
	return nil
}

type Good struct {
	Id            string `json:"id"`
	GNum          string `json:"gnum"`
	GName         string `json:"gname"`
	OriginalPrice string `json:"orignialprice"`
	CurPrice      string `json:"curprice"`
	CTime         string `json:"ctime"`
	CId           string `json:"cid"`
	GDes          string `json:"gdes"`
	GStatus       string `json:"gstatus"` // 1 卖光了
	OwnUname      string `json:"ownuname"`
	OnSale        string `json:"onsale"`
	SName         string `json:"sname"`
	Status        string `json:"status"`
}

func (n *Ntof) GoodList(id int) ([]*Good, error) {
	var data struct {
		Total  string  `json:"total"`
		Offset int     `json:"offset"`
		Goods  []*Good `json:"goods"`
	}
	err := n.client.Get(fmt.Sprintf(BaseURL+GoodListURL+"?page=1&count=100&sid=%d&token=%s", id, n.token), &data)
	if err != nil {
		return nil, err
	}
	return data.Goods, nil
}

func (n *Ntof) Buy(showid, gid string, token string) error {
	params := map[string]string{
		"gid":    gid,
		"showid": showid,
		"token":  token,
	}
	var out interface{}
	err := n.client.PostForm(BaseURL+BuyURL, params, out)
	if err != nil {
		return err
	}

	fmt.Println("buy response", out)
	return nil
}

var ntof *Ntof

func main() {
	ntof = newNtof()
	username := os.Getenv("USER")
	password := os.Getenv("PASSWORD")
	err := ntof.Login(username, password)
	if err != nil {
		log.Fatalln(err)
	}

	var (
		allGoods   []*Good
		initialCap float64
		marketCap  float64
		players    = make(map[string]int)
	)

	goods, err := ntof.GoodList(GoodSIdShangWu)
	if err != nil {
		log.Println(err)
		return
	}
	allGoods = append(allGoods, goods...)
	goods, err = ntof.GoodList(GoodSIdXiaWu)
	if err != nil {
		log.Println(err)
		return
	}
	allGoods = append(allGoods, goods...)

	for _, good := range allGoods {
		if _, ok := players[good.OwnUname]; !ok {
			players[good.OwnUname] = 0
		}
		players[good.OwnUname]++

		ori, _ := strconv.ParseFloat(good.OriginalPrice, 10)
		cur, _ := strconv.ParseFloat(good.CurPrice, 10)
		marketCap += cur
		initialCap += ori
	}

	fmt.Println("Initial:", initialCap)
	fmt.Println("MarketCap:", marketCap)
	fmt.Println("Players:", len(players))

	for player, amount := range players {
		fmt.Println(fmt.Sprintf("%s: %d", player, amount))
	}

	if os.Getenv("RUN") == "1" {
		runJob()
	}

	ntof.cron.AddFunc("55 29 11 * * *", runJob)
	ntof.cron.AddFunc("55 59 14 * * *", runJob)
	ntof.cron.Start()
	defer ntof.cron.Stop()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)

	select {
	case sig := <-sigChan:
		fmt.Println("received shutdown", "signal", sig)
	}

	fmt.Println("Graceful shutdown successful")
}

func runJob() {
	username := os.Getenv("USER")
	password := os.Getenv("PASSWORD")
	err := ntof.Login(username, password)
	if err != nil {
		log.Fatalln(err)
	}

	var start = time.Now()

	sid := GoodSIdShangWu
	hour := time.Now().Hour()
	if hour > 12 {
		sid = GoodSIdXiaWu
	}

	goods, err := ntof.GoodList(sid)
	if err != nil {
		log.Println(err)
		return
	}

	sort.Slice(goods, func(i, j int) bool {
		cur1, _ := strconv.ParseFloat(goods[i].CurPrice, 10)
		cur2, _ := strconv.ParseFloat(goods[j].CurPrice, 10)
		return cur1 > cur2
	})

	for {
		for _, good := range goods {
			// 卖光了
			if good.GStatus != "7" {
				continue
			}

			go func() {
				err := ntof.Buy(strconv.Itoa(sid), good.Id, ntof.token)
				if err != nil {
					fmt.Println(err)
				}
			}()
		}

		time.Sleep(50 * time.Millisecond)
		if time.Now().After(start.Add(1 * time.Minute)) {
			fmt.Println("timeout exit")
			break
		}
	}
}
